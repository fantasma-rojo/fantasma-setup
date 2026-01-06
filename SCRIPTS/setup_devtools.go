package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("❌ Could not determine home directory: %v\n", err)
		os.Exit(1)
	}

	binDir := filepath.Join(homeDir, "bin")
	commonProfile := filepath.Join(homeDir, ".common_profile")

	fmt.Println("🔧 Setting up Termux DevTools (Go Edition)...")

	// 1. Create ~/bin if missing
	if _, err := os.Stat(binDir); os.IsNotExist(err) {
		if err := os.MkdirAll(binDir, 0755); err != nil {
			fmt.Printf("❌ Failed to create %s: %v\n", binDir, err)
			os.Exit(1)
		}
		fmt.Printf("📁 Created %s\n", binDir)
	}

	// 2. Compile Tools
	compileTool("repo-dump.go", filepath.Join(binDir, "repo-dump"))
	compileTool("trigger-release.go", filepath.Join(binDir, "release-trigger"))
	compileTool("paste-to-file.go", filepath.Join(binDir, "paste-to-file"))

	// 3. Configure ~/.common_profile (Idempotent)
	configureCommonProfile(commonProfile, binDir)

	// 4. Link Shell Configs (Idempotent)
	linkShellConfig(filepath.Join(homeDir, ".bashrc"), commonProfile)
	linkShellConfig(filepath.Join(homeDir, ".zshrc"), commonProfile)

	fmt.Println("🎉 Setup Complete!")
	fmt.Println("👉 Run 'source ~/.zshrc' (or bashrc) to apply changes.")
}

// Helper: Compile a Go source file to the destination binary
func compileTool(srcFile, destPath string) {
	if _, err := os.Stat(srcFile); os.IsNotExist(err) {
		// Silently skip if source file isn't present (e.g., if you only have one tool)
		return
	}

	fmt.Printf("🔨 Compiling %s...\n", filepath.Base(destPath))
	cmd := exec.Command("go", "build", "-o", destPath, srcFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		fmt.Printf("❌ Failed to compile %s: %v\n", srcFile, err)
	} else {
		fmt.Printf("✅ Installed: %s\n", filepath.Base(destPath))
	}
}

// Helper: Create or update ~/.common_profile
func configureCommonProfile(profilePath, binDir string) {
	content := fmt.Sprintf(`
# Shared Environment Variables
export PATH="%s:$PATH"
`, binDir)

	// Check if file contains the bin path already
	if fileContains(profilePath, binDir) {
		fmt.Printf("⏭️  Skipping %s (already configured).\n", filepath.Base(profilePath))
		return
	}

	appendToFile(profilePath, content)
	fmt.Printf("✅ Updated %s\n", filepath.Base(profilePath))
}

// Helper: Link rc files to common profile
func linkShellConfig(rcPath, commonProfile string) {
	// If the shell config file doesn't exist, skip it
	if _, err := os.Stat(rcPath); os.IsNotExist(err) {
		return
	}

	sourceCmd := fmt.Sprintf(`[ -f "%s" ] && source "%s"`, commonProfile, commonProfile)
	checkStr := ".common_profile" // Simple string to check for existence

	if fileContains(rcPath, checkStr) {
		fmt.Printf("⏭️  Skipping %s (already linked).\n", filepath.Base(rcPath))
		return
	}

	appendBlock := fmt.Sprintf("\n# Source shared configuration\n%s\n", sourceCmd)
	appendToFile(rcPath, appendBlock)
	fmt.Printf("✅ Linked %s\n", filepath.Base(rcPath))
}

// Add this call inside main():
// injectReloadFunction(commonProfile)

func injectReloadFunction(profilePath string) {
	// This is the shell function we want to inject.
	// It detects the shell and sources the right file.
	reloadFunc := `
# Auto-detect shell and source the rc file
reload() {
    local shell_name=$(basename "$SHELL")
    local rc_file="$HOME/.${shell_name}rc"

    if [ -f "$rc_file" ]; then
        source "$rc_file"
        echo "♻️  Reloaded configuration ($rc_file)"
    else
        echo "❌ Could not determine config for $shell_name"
    fi
}
`
	if fileContains(profilePath, "reload()") {
		fmt.Printf("⏭️  Skipping reload function (already exists).\n")
		return
	}

	appendToFile(profilePath, reloadFunc)
	fmt.Printf("✅ Injected 'reload' function into %s\n", filepath.Base(profilePath))
}

// Utility: Check if file contains a string
func fileContains(path, search string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false // Treat read errors as "not found"
	}
	return strings.Contains(string(data), search)
}

// Utility: Append string to file
func appendToFile(path, content string) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("❌ Failed to open %s: %v\n", path, err)
		return
	}
	defer f.Close()
	if _, err := f.WriteString(content); err != nil {
		fmt.Printf("❌ Failed to write to %s: %v\n", path, err)
	}
}
