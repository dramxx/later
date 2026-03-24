package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/dramxx/later/config"
)

func Config() {
	if len(os.Args) >= 3 {
		switch os.Args[2] {
		case "--init":
			initConfig()
			return
		case "--path":
			configPath, err := config.EnsureFile()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			fmt.Println(configPath)
			return
		default:
			fmt.Println("Usage: later config [--init|--path]")
			os.Exit(1)
		}
	}

	configPath, err := config.EnsureFile()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	cmd, err := editorCommand(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Printf("Config file: %s\n", configPath)
		fmt.Println("Run 'later config --init' to set it up without an editor.")
		os.Exit(1)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err = cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening editor: %v\n", err)
		fmt.Printf("Config file: %s\n", configPath)
		fmt.Println("Run 'later config --init' to set it up without an editor.")
		os.Exit(1)
	}
}

func initConfig() {
	configPath, err := config.EnsureFile()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("GitHub token (gist scope only):")
	token, err := readRequiredValue(reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Private gist ID:")
	gistID, err := readRequiredValue(reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	err = config.Save(&config.Config{Gist: config.GistConfig{Token: token, GistID: gistID}})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ config saved to %s\n", configPath)
}

func readRequiredValue(reader *bufio.Reader) (string, error) {
	value, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", err
	}

	value = strings.TrimSpace(value)
	if value == "" {
		return "", fmt.Errorf("value cannot be empty")
	}

	return value, nil
}

func editorCommand(configPath string) (*exec.Cmd, error) {
	if editor := strings.TrimSpace(os.Getenv("VISUAL")); editor != "" {
		return exec.Command(editor, configPath), nil
	}

	if editor := strings.TrimSpace(os.Getenv("EDITOR")); editor != "" {
		return exec.Command(editor, configPath), nil
	}

	var candidates [][]string
	switch runtime.GOOS {
	case "windows":
		candidates = [][]string{{"notepad.exe", configPath}}
	case "darwin":
		candidates = [][]string{{"open", "-e", configPath}}
	default:
		candidates = [][]string{{"xdg-open", configPath}, {"nano", configPath}, {"vim", configPath}, {"vi", configPath}}
	}

	for _, candidate := range candidates {
		if _, err := exec.LookPath(candidate[0]); err == nil {
			return exec.Command(candidate[0], candidate[1:]...), nil
		}
	}

	return nil, fmt.Errorf("no editor found")
}
