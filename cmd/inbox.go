package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/dramxx/later/config"
	"github.com/dramxx/later/gist"
)

func buildContent(nonEntryLines, entryLines []string) string {
	var lines []string
	for _, line := range nonEntryLines {
		if line != "" {
			lines = append(lines, line)
		}
	}
	for _, line := range entryLines {
		lines = append(lines, line)
	}
	if len(lines) == 0 {
		return ""
	}
	return strings.Join(lines, "\n") + "\n"
}

func Inbox() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	content, err := gist.GetInbox(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	lines := strings.Split(content, "\n")
	if len(lines) == 0 || (len(lines) == 1 && lines[0] == "") {
		lines = nil
	}
	var entryLines []string
	var nonEntryLines []string

	for _, line := range lines {
		if strings.HasPrefix(line, "[") {
			entryLines = append(entryLines, line)
		} else {
			nonEntryLines = append(nonEntryLines, line)
		}
	}

	if len(os.Args) >= 3 && os.Args[2] == "--clear" {
		newContent := buildContent(nonEntryLines, nil)

		err = gist.UpdateInbox(cfg, newContent)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✓ cleared")
		return
	}

	if len(os.Args) >= 3 && os.Args[2] == "--pop" {
		if len(os.Args) < 4 {
			fmt.Println("Usage: later inbox --pop <n> [n...]")
			os.Exit(1)
		}

		indices := []int{}
		seen := make(map[int]bool)
		for _, arg := range os.Args[3:] {
			n, err := strconv.Atoi(arg)
			if err != nil || n < 1 || n > len(entryLines) {
				fmt.Fprintf(os.Stderr, "Error: invalid index: %s\n", arg)
				os.Exit(1)
			}
			if seen[n] {
				fmt.Fprintf(os.Stderr, "Error: duplicate index: %d\n", n)
				os.Exit(1)
			}
			seen[n] = true
			indices = append(indices, n)
		}

		removed := 0
		for i := len(entryLines) - 1; i >= 0; i-- {
			for _, idx := range indices {
				if i == idx-1 {
					entryLines = append(entryLines[:i], entryLines[i+1:]...)
					removed++
					break
				}
			}
		}

		newContent := buildContent(nonEntryLines, entryLines)

		err = gist.UpdateInbox(cfg, newContent)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if removed == 1 {
			fmt.Println("✓ removed 1 entry")
		} else {
			fmt.Printf("✓ removed %d entries\n", removed)
		}
		return
	}

	if len(entryLines) == 0 {
		fmt.Println("inbox is empty")
		return
	}

	for i, line := range entryLines {
		fmt.Printf("%d  %s\n", i+1, line)
	}
}
