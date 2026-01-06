#!/bin/bash

# Function to display usage
usage() {
    echo "Usage: $0 <github-organization>"
    echo "Example: $0 my-company-org"
    exit 1
}

# 1. Validate Inputs
ORG_NAME=$1
if [ -z "$ORG_NAME" ]; then
    usage
fi

# 2. Pre-flight Check: Is gh installed?
if ! command -v gh &> /dev/null; then
    echo "Error: GitHub CLI (gh) is not installed."
    exit 1
fi

# 3. Pre-flight Check: Auth Status & Scopes
# We check if the user has the 'read:org' scope, which is required for fetching teams.
AUTH_STATUS=$(gh auth status 2>&1)

if echo "$AUTH_STATUS" | grep -q "not logged in"; then
    echo "Error: You are not logged in. Run: gh auth login"
    exit 1
fi

# 4. Attempt to Fetch Teams with Error Handling
echo "Fetching teams for organization: $ORG_NAME..." >&2

# We capture both stdout (data) and stderr (errors) to separate variables
# using a temporary file for the data to avoid pipe masking exit codes.
TEMP_DATA=$(mktemp)
ERROR_LOG=$(mktemp)

if gh api "orgs/$ORG_NAME/teams" \
    --paginate \
    --method GET \
    --jq '.[] | [.name, .slug, .description] | @tsv' > "$TEMP_DATA" 2> "$ERROR_LOG"; then
    
    # Success Block
    echo "NAME	SLUG	DESCRIPTION"
    echo "------------------------------------------------"
    cat "$TEMP_DATA"
    
else
    # Failure Block - Diagnose the specific error
    ERROR_MSG=$(cat "$ERROR_LOG")
    
    echo "------------------------------------------------"
    echo "FAILED TO DOWNLOAD TEAM LIST"
    echo "------------------------------------------------"
    
    if echo "$ERROR_MSG" | grep -q "SAML"; then
        echo "CAUSE: SAML SSO is enforced by this organization."
        echo "FIX: Run this command to authorize your token:"
        echo "     gh auth refresh -h github.com -s read:org"
    elif echo "$ERROR_MSG" | grep -q "Not Found"; then
        echo "CAUSE: Organization '$ORG_NAME' not found OR you do not have access."
        echo "FIX: Check spelling and ensure you are a member of the org."
    elif echo "$ERROR_MSG" | grep -q "scope"; then
        echo "CAUSE: Your token is missing the 'read:org' scope."
        echo "FIX: Run: gh auth refresh -s read:org"
    else
        echo "RAW ERROR from GitHub:"
        echo "$ERROR_MSG"
    fi
    
    rm "$TEMP_DATA" "$ERROR_LOG"
    exit 1
fi

# Cleanup
rm "$TEMP_DATA" "$ERROR_LOG"
