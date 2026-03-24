package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dramxx/later/config"
	"github.com/dramxx/later/gist"
)

func Send() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: later send <text>")
		os.Exit(1)
	}

	text := strings.Join(os.Args[2:], " ")

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

	timestamp := time.Now().Format("2006-01-02 15:04")
	newLine := fmt.Sprintf("[%s]  %s", timestamp, text)

	if content == "" {
		content = "LATER\n\n" + newLine + "\n"
	} else {
		if !strings.HasSuffix(content, "\n") {
			content += "\n"
		}
		content += newLine + "\n"
	}

	err = gist.UpdateInbox(cfg, content)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✓ saved")
}
