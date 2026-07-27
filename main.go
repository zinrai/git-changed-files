package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	dir := flag.String("dir", "", "directory to detect changes in (required)")
	ext := flag.String("ext", "", "filter by file extension (e.g. .json)")
	max := flag.Int("max", 0, "maximum number of changed files allowed (0 = unlimited)")
	ref := flag.String("ref", "", "git ref to compare against (required)")
	flag.Parse()

	if *dir == "" {
		fmt.Fprintln(os.Stderr, "error: -dir is required")
		os.Exit(1)
	}

	if *ref == "" {
		fmt.Fprintln(os.Stderr, "error: -ref is required")
		os.Exit(1)
	}

	if err := run(*dir, *ext, *max, *ref); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
