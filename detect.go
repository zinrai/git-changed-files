package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func run(dir, ext string, max int, ref string) error {
	base := ref
	if isNullRef(base) {
		var err error
		base, err = emptyTreeHash()
		if err != nil {
			return err
		}
	}

	cmd := exec.Command("git", "diff", "--name-only", "--diff-filter=d", base, "HEAD", "--", dir)
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("git diff failed: %w", err)
	}

	files := splitLines(string(out))
	files = filterByDirectory(files, dir)

	if ext != "" {
		files = filterByExtension(files, ext)
	}

	if max > 0 && len(files) > max {
		return fmt.Errorf("found %d changed files, exceeding the limit of %d\n%s", len(files), max, strings.Join(files, "\n"))
	}

	for _, f := range files {
		fmt.Println(f)
	}

	return nil
}

// isNullRef reports whether ref is an all-zeros SHA. Git hooks and workflow
// engines pass it as the old rev when no comparison base exists, such as on
// branch creation or force push.
func isNullRef(ref string) bool {
	if ref == "" {
		return false
	}
	for _, c := range ref {
		if c != '0' {
			return false
		}
	}
	return true
}

// emptyTreeHash returns the hash of the empty tree for the repository's
// object format, so a null ref can be diffed as "everything at HEAD".
func emptyTreeHash() (string, error) {
	out, err := exec.Command("git", "hash-object", "-t", "tree", os.DevNull).Output()
	if err != nil {
		return "", fmt.Errorf("git hash-object failed: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func filterByDirectory(files []string, dir string) []string {
	cleaned := filepath.Clean(dir)
	if cleaned == "." {
		return files
	}

	prefix := cleaned + string(filepath.Separator)
	var result []string
	for _, f := range files {
		if strings.HasPrefix(filepath.Clean(f), prefix) {
			result = append(result, f)
		}
	}
	return result
}

func filterByExtension(files []string, ext string) []string {
	var result []string
	for _, f := range files {
		if filepath.Ext(f) == ext {
			result = append(result, f)
		}
	}
	return result
}

func splitLines(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	return strings.Split(s, "\n")
}
