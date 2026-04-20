package main

import "testing"

func TestFilterByDirectory(t *testing.T) {
	files := []string{
		"configs/app/prod.json",
		"configs/app/staging.json",
		"configs/db/prod.json",
		"other/file.txt",
	}

	t.Run("filters by directory", func(t *testing.T) {
		result := filterByDirectory(files, "configs/app")
		if len(result) != 2 {
			t.Fatalf("len = %d, want 2", len(result))
		}
		if result[0] != "configs/app/prod.json" {
			t.Errorf("result[0] = %q, want %q", result[0], "configs/app/prod.json")
		}
		if result[1] != "configs/app/staging.json" {
			t.Errorf("result[1] = %q, want %q", result[1], "configs/app/staging.json")
		}
	})

	t.Run("filters by different directory", func(t *testing.T) {
		result := filterByDirectory(files, "configs/db")
		if len(result) != 1 {
			t.Fatalf("len = %d, want 1", len(result))
		}
	})

	t.Run("no match", func(t *testing.T) {
		result := filterByDirectory(files, "nonexistent")
		if len(result) != 0 {
			t.Errorf("len = %d, want 0", len(result))
		}
	})

	t.Run("empty input", func(t *testing.T) {
		result := filterByDirectory([]string{}, "configs/app")
		if len(result) != 0 {
			t.Errorf("len = %d, want 0", len(result))
		}
	})

	t.Run("does not match partial directory name", func(t *testing.T) {
		files := []string{
			"configs/appdata/prod.json",
			"configs/app/prod.json",
		}
		result := filterByDirectory(files, "configs/app")
		if len(result) != 1 {
			t.Fatalf("len = %d, want 1", len(result))
		}
		if result[0] != "configs/app/prod.json" {
			t.Errorf("result[0] = %q, want %q", result[0], "configs/app/prod.json")
		}
	})

	t.Run("dot directory returns all files", func(t *testing.T) {
		files := []string{
			"configs/app/prod.json",
			"data/file.txt",
		}
		result := filterByDirectory(files, ".")
		if len(result) != 2 {
			t.Errorf("len = %d, want 2", len(result))
		}
	})
}

func TestFilterByExtension(t *testing.T) {
	files := []string{
		"data/001.json",
		"data/002.json",
		"data/notes.txt",
		"data/config.yaml",
	}

	t.Run("filters .json", func(t *testing.T) {
		result := filterByExtension(files, ".json")
		if len(result) != 2 {
			t.Fatalf("len = %d, want 2", len(result))
		}
	})

	t.Run("filters .txt", func(t *testing.T) {
		result := filterByExtension(files, ".txt")
		if len(result) != 1 {
			t.Fatalf("len = %d, want 1", len(result))
		}
	})

	t.Run("no match", func(t *testing.T) {
		result := filterByExtension(files, ".xml")
		if len(result) != 0 {
			t.Errorf("len = %d, want 0", len(result))
		}
	})

	t.Run("empty input", func(t *testing.T) {
		result := filterByExtension([]string{}, ".json")
		if len(result) != 0 {
			t.Errorf("len = %d, want 0", len(result))
		}
	})
}

func TestSplitLines(t *testing.T) {
	t.Run("multiple lines", func(t *testing.T) {
		result := splitLines("a.json\nb.json\nc.json")
		if len(result) != 3 {
			t.Fatalf("len = %d, want 3", len(result))
		}
	})

	t.Run("trailing newline", func(t *testing.T) {
		result := splitLines("a.json\nb.json\n")
		if len(result) != 2 {
			t.Fatalf("len = %d, want 2", len(result))
		}
	})

	t.Run("single line", func(t *testing.T) {
		result := splitLines("a.json")
		if len(result) != 1 {
			t.Fatalf("len = %d, want 1", len(result))
		}
		if result[0] != "a.json" {
			t.Errorf("result[0] = %q, want %q", result[0], "a.json")
		}
	})

	t.Run("empty string", func(t *testing.T) {
		result := splitLines("")
		if len(result) != 0 {
			t.Errorf("len = %d, want 0", len(result))
		}
	})

	t.Run("whitespace only", func(t *testing.T) {
		result := splitLines("  \n  \n  ")
		if len(result) != 0 {
			t.Errorf("len = %d, want 0", len(result))
		}
	})
}
