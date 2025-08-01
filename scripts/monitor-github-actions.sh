#!/bin/bash

# GitHub Actions Status Monitor
# Monitors the status of GitHub Actions workflows

echo "🔍 Monitoring GitHub Actions Status..."
echo "Repository: dewisartika8/CICD-Status-Notifier-Bot"
echo "Branch: staging"
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Check if gh CLI is available
if ! command -v gh &> /dev/null; then
    echo -e "${YELLOW}[INFO]${NC} GitHub CLI (gh) not installed. Opening browser instead..."
    echo ""
    echo "🌐 Manual check:"
    echo "   https://github.com/dewisartika8/CICD-Status-Notifier-Bot/actions"
    echo ""
    echo "📋 Expected workflows:"
    echo "   ✓ Quick CI Test (quick-ci.yml)"
    echo "   ✓ Simple CI Testing (simple-ci.yml)" 
    echo "   ✓ Local Testing CI/CD Pipeline (local-ci-cd.yml)"
    echo ""
    echo "🎯 What to look for:"
    echo "   - All workflows should show green checkmarks ✅"
    echo "   - Build steps should complete successfully"
    echo "   - Docker builds should work"
    echo "   - No red X marks ❌"
    echo ""
    echo "⏱️  Expected completion time: ~5-10 minutes"
    exit 0
fi

# If gh CLI is available, get status
echo "📊 Fetching workflow status..."
gh run list --repo dewisartika8/CICD-Status-Notifier-Bot --branch staging --limit 5

echo ""
echo "🔄 Watching for latest run..."
gh run watch --repo dewisartika8/CICD-Status-Notifier-Bot

echo ""
echo -e "${GREEN}[SUCCESS]${NC} Monitoring complete!"
