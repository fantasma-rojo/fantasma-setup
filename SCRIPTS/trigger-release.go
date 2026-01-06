package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	// 1. Ensure we are in a git repo
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		fmt.Println("❌ Error: Not a git repository.")
		os.Exit(1)
	}

	fmt.Println("🔄 Fetching latest tags...")
	exec.Command("git", "fetch", "--tags").Run()

	// 2. Get the latest tag
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	output, err := cmd.Output()
	
	currentTag := "v0.0.0"
	if err == nil {
		currentTag = strings.TrimSpace(string(output))
	}

	// 3. Parse Semantic Version (vX.Y.Z)
	// Remove 'v' prefix
	version := strings.TrimPrefix(currentTag, "v")
	parts := strings.Split(version, ".")

	var major, minor, patch int

	if len(parts) >= 3 {
		major, _ = strconv.Atoi(parts[0])
		minor, _ = strconv.Atoi(parts[1])
		patch, _ = strconv.Atoi(parts[2])
	}

	// 4. Increment Patch Version
	newPatch := patch + 1
	nextTag := fmt.Sprintf("v%d.%d.%d", major, minor, newPatch)

	// 5. User Confirmation
	fmt.Println("🚀 Release Trigger")
	fmt.Println("------------------")
	fmt.Printf("Current Version: %s\n", currentTag)
	fmt.Printf("Next Version:    %s\n\n", nextTag)

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Push release %s? [y/N]: ", nextTag)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))

	if input == "y" || input == "yes" {
		fmt.Println("📦 Tagging release...")
		if err := exec.Command("git", "tag", "-a", nextTag, "-m", "Release "+nextTag).Run(); err != nil {
			fmt.Printf("❌ Failed to create tag: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("⬆️  Pushing to origin...")
		if err := exec.Command("git", "push", "origin", nextTag).Run(); err != nil {
			fmt.Printf("❌ Failed to push tag: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("✅ Done! GitHub Action triggered.")
	} else {
		fmt.Println("❌ Cancelled.")
	}
}
