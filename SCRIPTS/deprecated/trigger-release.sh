#!/bin/bash

# Ensure we are in a git repo
if [ ! -d ".git" ]; then
    echo "❌ Error: Not a git repository."
    exit 1
fi

# Fetch tags to ensure we have the latest list
git fetch --tags > /dev/null 2>&1

# Get the latest tag, default to v0.0.0 if none exists
LATEST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")

# Logic to increment Semantic Version (vX.Y.Z -> vX.Y.Z+1)
VERSION=${LATEST_TAG#v}
IFS='.' read -r -a PARTS <<< "$VERSION"
MAJOR=${PARTS[0]}
MINOR=${PARTS[1]}
PATCH=${PARTS[2]}
NEW_PATCH=$((PATCH + 1))
NEW_TAG="v$MAJOR.$MINOR.$NEW_PATCH"

echo "🚀 Release Trigger"
echo "------------------"
echo "Current Version: $LATEST_TAG"
echo "Next Version:    $NEW_TAG"
echo ""
read -p "Push release $NEW_TAG? [y/N] " -n 1 -r
echo ""

if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "📦 Tagging release..."
    git tag -a "$NEW_TAG" -m "Release $NEW_TAG via CLI"
    
    echo "⬆️  Pushing to origin..."
    git push origin "$NEW_TAG"
    
    echo "✅ Done! GitHub Action triggered."
else
    echo "❌ Cancelled."
fi
