package main

import (
	"fmt"
	"os"

	"github.com/dramxx/later/cmd"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	if os.Args[1] == "-h" || os.Args[1] == "--help" {
		printUsage()
		return
	}

	switch os.Args[1] {
	case "send":
		cmd.Send()
	case "inbox":
		cmd.Inbox()
	case "config":
		cmd.Config()
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: later <command>")
	fmt.Println("Commands: send, inbox, config")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  later config --init")
	fmt.Println("  later send https://example.com")
	fmt.Println("  later inbox")
	fmt.Println("  later inbox --clear")
	fmt.Println("  later inbox --pop 1")
	fmt.Println("  later config")
}
