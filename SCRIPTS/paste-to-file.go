package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	// 1. Check for filename argument
	if len(os.Args) < 2 {
		fmt.Println("❌ Usage: paste-to-file <filename>")
		os.Exit(1)
	}

	targetFile := os.Args[1]

	// 2. Get Clipboard Content
	// We use the Termux API command
	cmd := exec.Command("termux-clipboard-get")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("❌ Failed to get clipboard. Is Termux API installed?\nError: %v\n", err)
		os.Exit(1)
	}

	if len(output) == 0 {
		fmt.Println("⚠️  Clipboard is empty.")
		return
	}

	// 3. Write to File (Truncate/Overwrite)
	// 0644 means readable by user, writable by user, readable by group/others
	err = os.WriteFile(targetFile, output, 0644)
	if err != nil {
		fmt.Printf("❌ Failed to write to file: %v\n", err)
		os.Exit(1)
	}

	// 4. Success Message
	absPath, _ := filepath.Abs(targetFile)
	fmt.Printf("✅ Clipboard content (%d bytes) saved to:\n   %s\n", len(output), absPath)
}
