#!/bin/bash

BIN_DIR="$HOME/bin"
COMMON_PROFILE="$HOME/.common_profile"

echo "🔧 Setting up Termux DevTools (Shared Config)..."

# 1. Create ~/bin if missing
if [ ! -d "$BIN_DIR" ]; then
    mkdir -p "$BIN_DIR"
fi

# 2. Compile repo-dump
if [ -f "repo-dump.go" ]; then
    echo "🔨 Compiling repo-dump..."
    go build -o "$BIN_DIR/repo-dump" repo-dump.go
    echo "✅ Installed: repo-dump"
fi

# 3. Compile release-trigger
if [ -f "trigger-release.go" ]; then
    echo "🔨 Compiling release-trigger..."
    go build -o "$BIN_DIR/release-trigger" trigger-release.go
    echo "✅ Installed: release-trigger"
fi

# 4. Create/Update ~/.common_profile
# We check if the PATH export already exists to avoid duplicates
if ! grep -q "$BIN_DIR" "$COMMON_PROFILE" 2>/dev/null; then
    echo "" >> "$COMMON_PROFILE"
    echo "# Shared Environment Variables" >> "$COMMON_PROFILE"
    echo "export PATH=\"$BIN_DIR:\$PATH\"" >> "$COMMON_PROFILE"
    echo "✅ Created/Updated $COMMON_PROFILE"
else
    echo "🔗 $COMMON_PROFILE already configured."
fi

# 5. Link to Shell Configs
# Function to safely append the source command
link_shell_config() {
    local rc_file="$1"
    local source_cmd="[ -f \"$COMMON_PROFILE\" ] && source \"$COMMON_PROFILE\""
    
    if [ -f "$rc_file" ]; then
        if ! grep -q "common_profile" "$rc_file"; then
            echo "" >> "$rc_file"
            echo "# Source shared configuration" >> "$rc_file"
            echo "$source_cmd" >> "$rc_file"
            echo "🔗 Linked to $rc_file"
        else
            echo "🔗 $rc_file already linked."
        fi
    fi
}

link_shell_config "$HOME/.bashrc"
link_shell_config "$HOME/.zshrc"

echo "🎉 Setup Complete!"
echo "👉 Run 'source ~/.zshrc' to apply changes immediately."
