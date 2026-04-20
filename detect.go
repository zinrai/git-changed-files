package main

import (
	"flag"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

func run() error {
	dir := flag.String("dir", ".", "directory to detect changes in")
	ext := flag.String("ext", "", "filter by file extension (e.g. .json)")
	max := flag.Int("max", 0, "maximum number of changed files allowed (0 = unlimited)")
	ref := flag.String("ref", "", "git ref to compare against (required)")
	flag.Parse()

	if *ref == "" {
		return fmt.Errorf("-ref is required")
	}

	cmd := exec.Command("git", "diff", "--name-only", *ref, "HEAD", "--", *dir)
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("git diff failed: %w", err)
	}

	files := splitLines(string(out))
	files = filterByDirectory(files, *dir)

	if *ext != "" {
		files = filterByExtension(files, *ext)
	}

	if len(files) == 0 {
		return fmt.Errorf("no changed files found in %s", *dir)
	}

	if *max > 0 && len(files) > *max {
		return fmt.Errorf("found %d changed files, exceeding the limit of %d\n%s", len(files), *max, strings.Join(files, "\n"))
	}

	for _, f := range files {
		fmt.Println(f)
	}

	return nil
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
