#!/bin/bash

# GitHub Workflow Status Checker
# Usage: ./scripts/check-workflows.sh

set -e

REPO="kringen/homelab"
WORKFLOWS=("ci.yml" "deploy.yml")

echo "🔍 Checking GitHub workflow status for $REPO"
echo "=============================================="

for workflow in "${WORKFLOWS[@]}"; do
    echo ""
    echo "📋 $workflow:"
    
    # Check if gh CLI is available
    if command -v gh > /dev/null; then
        # Get latest workflow run status
        gh run list --workflow="$workflow" --limit=1 --json status,conclusion,headBranch,createdAt --template '{{range .}}Status: {{.status}} | Conclusion: {{.conclusion}} | Branch: {{.headBranch}} | Created: {{.createdAt}}{{end}}'
    else
        echo "   ⚠️  GitHub CLI (gh) not installed. Install with: brew install gh"
        echo "   📖 View manually: https://github.com/$REPO/actions/workflows/$workflow"
    fi
done

echo ""
echo "🔗 View all workflows: https://github.com/$REPO/actions"
echo ""
echo "💡 Tips:"
echo "   - Install GitHub CLI: brew install gh (macOS) or apt install gh (Ubuntu)"
echo "   - Authenticate: gh auth login"
echo "   - View workflow details: gh run view [run-id]"
