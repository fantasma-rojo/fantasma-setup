package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

func main() {
	var buffer bytes.Buffer
	root, _ := os.Getwd()
	projectName := filepath.Base(root)
	fileCount := 0

	fmt.Printf("📂 Scanning: %s\n", projectName)
	fmt.Println("----------------------------------------")

	// 1. Walk the Directory
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil { return err }

		// SKIP DIRECTORIES
		if d.IsDir() {
			name := d.Name()
			// Skip hidden dirs (like .git) and build artifacts
			if strings.HasPrefix(name, ".") || name == "target" || name == "node_modules" || name == "vendor" {
				return filepath.SkipDir
			}
			return nil
		}

		// SKIP EXTENSIONS
		ext := strings.ToLower(filepath.Ext(path))
		skipExts := map[string]bool{
			".ko": true, ".txz": true, ".o": true, ".a": true, 
			".png": true, ".jpg": true, ".jpeg": true, ".exe": true, 
			".bin": true, ".lock": true, ".zip": true, ".gz": true,
		}
		if skipExts[ext] { return nil }

		// SKIP SELF
		if filepath.Base(path) == "repo-dump" || filepath.Base(path) == "repo-dump.go" {
			return nil
		}

		// READ CONTENT
		content, err := os.ReadFile(path)
		if err != nil { 
			fmt.Printf("⚠️  Read Error: %s\n", filepath.Base(path))
			return nil 
		}

		// UTF-8 CHECK
		if !utf8.Valid(content) { return nil }

		// ADD TO BUFFER
		relPath, _ := filepath.Rel(root, path)
		fmt.Printf("📄 Adding: %s\n", relPath)
		
		buffer.WriteString(fmt.Sprintf("\n========== FILE: %s ==========\n", relPath))
		buffer.Write(content)
		buffer.WriteString("\n")
		fileCount++

		return nil
	})

	if err != nil {
		fmt.Printf("❌ Walk Error: %v\n", err)
		return
	}

	if fileCount == 0 {
		fmt.Println("⚠️  No text files found to dump.")
		return
	}

	fmt.Println("----------------------------------------")
	fmt.Printf("📦 Total Files: %d | Size: %d bytes\n", fileCount, buffer.Len())

	// 2. Try Clipboard
	if tryClipboard(buffer) {
		fmt.Println("✅ Success! Copied to Termux clipboard.")
	} else {
		// 3. Fallback to File
		saveToFile(buffer)
	}
}

func tryClipboard(data bytes.Buffer) bool {
	fmt.Print("📋 Attempting to copy to clipboard... ")
	
	// Check if tool exists
	_, err := exec.LookPath("termux-clipboard-set")
	if err != nil {
		fmt.Println("FAILED (Tool not found)")
		fmt.Println("   (Run 'pkg install termux-api' to fix)")
		return false
	}

	cmd := exec.Command("termux-clipboard-set")
	cmd.Stdin = &data
	
	// Capture stderr to see why it fails
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("FAILED\n❌ Error: %v\n   %s\n", err, stderr.String())
		return false
	}

	return true
}

func saveToFile(data bytes.Buffer) {
	filename := "repo_dump.txt"
	fmt.Printf("💾 Fallback: Saving to '%s'...", filename)
	
	err := os.WriteFile(filename, data.Bytes(), 0644)
	if err != nil {
		fmt.Printf("FAILED\n❌ Write Error: %v\n", err)
		return
	}
	
	absPath, _ := filepath.Abs(filename)
	fmt.Printf("DONE\n✅ Saved to: %s\n", absPath)
}
