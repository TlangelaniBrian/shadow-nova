#!/bin/bash

################################################################################
# Git History Cleanup Script
#
# Purpose: Remove sensitive files (.env, secrets) from git history
#
# WARNING: This script will rewrite git history. All contributors must
#          re-clone the repository after this operation.
#
# Prerequisites:
#   - git-filter-repo installed (brew install git-filter-repo)
#   - Clean working directory (no uncommitted changes)
#   - Backup of repository recommended
#
# Usage: ./scripts/cleanup-git-history.sh
################################################################################

set -e  # Exit on error
set -u  # Exit on undefined variable

# Color output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Script configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
BACKUP_DIR="$HOME/git-backups"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
BACKUP_NAME="shadow-nova-backup-$TIMESTAMP"

################################################################################
# Helper Functions
################################################################################

print_header() {
    echo -e "\n${BLUE}===================================================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}===================================================================${NC}\n"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ $1${NC}"
}

confirm_action() {
    local prompt="$1"
    local response
    echo -e "${YELLOW}$prompt (yes/no)${NC}"
    read -r response
    if [[ "$response" != "yes" ]]; then
        print_error "Operation cancelled by user"
        exit 1
    fi
}

################################################################################
# Pre-flight Checks
################################################################################

preflight_checks() {
    print_header "Running Pre-flight Checks"

    # Check if we're in a git repository
    if ! git rev-parse --git-dir > /dev/null 2>&1; then
        print_error "Not a git repository"
        exit 1
    fi
    print_success "In git repository"

    # Check for uncommitted changes
    if ! git diff-index --quiet HEAD -- 2>/dev/null; then
        print_error "You have uncommitted changes. Please commit or stash them first."
        git status --short
        exit 1
    fi
    print_success "Working directory is clean"

    # Check if git-filter-repo is installed
    if ! command -v git-filter-repo &> /dev/null; then
        print_error "git-filter-repo is not installed"
        echo ""
        echo "Install it using one of these methods:"
        echo "  macOS:   brew install git-filter-repo"
        echo "  Ubuntu:  apt-get install git-filter-repo"
        echo "  pip:     pip3 install git-filter-repo"
        echo ""
        exit 1
    fi
    print_success "git-filter-repo is installed"

    # Check current branch
    CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
    print_info "Current branch: $CURRENT_BRANCH"

    # Check for remotes
    if git remote | grep -q "origin"; then
        REMOTE_URL=$(git remote get-url origin)
        print_info "Remote origin: $REMOTE_URL"
    else
        print_warning "No remote 'origin' configured"
    fi
}

################################################################################
# Backup Creation
################################################################################

create_backup() {
    print_header "Creating Backup"

    # Create backup directory if it doesn't exist
    mkdir -p "$BACKUP_DIR"

    # Create backup using git clone
    print_info "Backing up to: $BACKUP_DIR/$BACKUP_NAME"
    git clone "$REPO_ROOT" "$BACKUP_DIR/$BACKUP_NAME"

    print_success "Backup created at: $BACKUP_DIR/$BACKUP_NAME"
    print_info "You can restore from this backup if needed"
}

################################################################################
# History Cleanup
################################################################################

cleanup_history() {
    print_header "Cleaning Git History"

    print_warning "This will remove the following files from ALL git history:"
    echo "  - .env"
    echo "  - .env.local"
    echo "  - .env.production"
    echo "  - .env.development"
    echo "  - .env.*"
    echo "  - **/backend/.env"
    echo "  - **/frontend/.env"
    echo ""

    confirm_action "Are you absolutely sure you want to proceed?"

    # Create paths file for git-filter-repo
    PATHS_FILE=$(mktemp)
    cat > "$PATHS_FILE" << 'EOF'
.env
.env.local
.env.production
.env.development
.env.staging
.env.test
backend/.env
backend/.env.local
frontend/.env
frontend/.env.local
EOF

    print_info "Running git-filter-repo..."

    # Run git-filter-repo
    git filter-repo --invert-paths --paths-from-file "$PATHS_FILE" --force

    # Cleanup temp file
    rm "$PATHS_FILE"

    print_success "Git history cleaned successfully"
}

################################################################################
# Verification
################################################################################

verify_cleanup() {
    print_header "Verifying Cleanup"

    # Check if .env files exist in history
    print_info "Searching for .env files in git history..."

    if git log --all --full-history --oneline -- '*.env' 2>/dev/null | grep -q .; then
        print_error "Found .env files still in history:"
        git log --all --full-history --oneline -- '*.env'
        return 1
    else
        print_success "No .env files found in git history"
    fi

    # Search for potential secrets in history
    print_info "Searching for potential secrets in commit messages..."

    SECRET_PATTERNS="CLIENT_ID|CLIENT_SECRET|API_KEY|JWT_SECRET|DATABASE_URL"
    if git log --all --grep="$SECRET_PATTERNS" --oneline 2>/dev/null | grep -q .; then
        print_warning "Found potential secrets in commit messages:"
        git log --all --grep="$SECRET_PATTERNS" --oneline
        echo ""
        print_info "Consider using 'git rebase -i' to reword these commits"
    else
        print_success "No obvious secrets in commit messages"
    fi

    # Show repository size comparison
    print_info "Repository statistics:"
    git count-objects -vH
}

################################################################################
# Post-Cleanup Instructions
################################################################################

show_next_steps() {
    print_header "Next Steps"

    echo -e "${YELLOW}IMPORTANT: Read these instructions carefully${NC}\n"

    echo "1. Verify the cleanup worked:"
    echo "   git log --all --full-history --oneline -- '*.env'"
    echo "   (Should return no results)"
    echo ""

    echo "2. Re-add your remote (git-filter-repo removes remotes):"
    if [[ -n "${REMOTE_URL:-}" ]]; then
        echo "   git remote add origin $REMOTE_URL"
    else
        echo "   git remote add origin <your-repo-url>"
    fi
    echo ""

    echo "3. Force push to remote (THIS WILL REWRITE HISTORY):"
    echo "   ${RED}git push origin --force --all${NC}"
    echo "   ${RED}git push origin --force --tags${NC}"
    echo ""

    echo "4. Notify all team members:"
    echo "   - History has been rewritten"
    echo "   - They must re-clone the repository"
    echo "   - They should NOT try to pull or merge"
    echo ""

    echo "5. Each team member should:"
    echo "   cd .."
    echo "   rm -rf shadow-nova"
    echo "   git clone <repo-url>"
    echo ""

    echo "6. On GitHub/GitLab:"
    echo "   - Go to Settings → Branches"
    echo "   - Temporarily disable branch protection"
    echo "   - Force push"
    echo "   - Re-enable branch protection"
    echo ""

    echo "7. Enable secret scanning:"
    echo "   - GitHub: Settings → Security → Secret scanning"
    echo "   - GitLab: Settings → Security & Compliance"
    echo ""

    echo -e "\n${GREEN}Backup location: $BACKUP_DIR/$BACKUP_NAME${NC}"
    echo "To restore from backup:"
    echo "  cd .."
    echo "  rm -rf shadow-nova"
    echo "  cp -r $BACKUP_DIR/$BACKUP_NAME shadow-nova"
    echo ""
}

################################################################################
# Main Execution
################################################################################

main() {
    print_header "Git History Cleanup Script"
    print_warning "This script will REWRITE git history"
    print_warning "All team members will need to RE-CLONE the repository"
    echo ""

    # Change to repository root
    cd "$REPO_ROOT"

    # Run preflight checks
    preflight_checks

    # Confirm one more time
    echo ""
    confirm_action "This is your LAST CHANCE to cancel. Proceed with history rewrite?"

    # Create backup
    create_backup

    # Cleanup history
    cleanup_history

    # Verify cleanup
    verify_cleanup

    # Show next steps
    show_next_steps

    print_header "Cleanup Complete"
    print_success "Git history has been cleaned"
    print_warning "Remember to follow the 'Next Steps' above"
}

# Run main function
main "$@"
