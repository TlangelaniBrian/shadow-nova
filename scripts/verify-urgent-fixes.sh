#!/bin/bash

################################################################################
# Urgent Fixes Verification Script
#
# Purpose: Verify all urgent security and bug fixes have been applied
#
# Checks:
#   1. No secrets in current files
#   2. Admin middleware exists and is properly implemented
#   3. Token naming consistency (auth_token vs token)
#   4. Database indexes exist
#   5. github_username column added
#   6. .env files properly gitignored
#
# Usage: ./scripts/verify-urgent-fixes.sh
################################################################################

set -u  # Exit on undefined variable

# Color output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Script configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Counters
TOTAL_CHECKS=0
PASSED_CHECKS=0
FAILED_CHECKS=0
WARNING_CHECKS=0

################################################################################
# Helper Functions
################################################################################

print_header() {
    echo -e "\n${BLUE}===================================================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}===================================================================${NC}\n"
}

print_section() {
    echo -e "\n${CYAN}--- $1 ---${NC}\n"
}

print_check() {
    echo -e "${CYAN}[CHECK]${NC} $1"
    ((TOTAL_CHECKS++))
}

print_pass() {
    echo -e "${GREEN}  ✓ PASS:${NC} $1"
    ((PASSED_CHECKS++))
}

print_fail() {
    echo -e "${RED}  ✗ FAIL:${NC} $1"
    ((FAILED_CHECKS++))
}

print_warning() {
    echo -e "${YELLOW}  ⚠ WARN:${NC} $1"
    ((WARNING_CHECKS++))
}

print_info() {
    echo -e "${BLUE}  ℹ INFO:${NC} $1"
}

################################################################################
# Check 1: Secrets in Current Files
################################################################################

check_no_secrets_in_files() {
    print_section "Check 1: No Secrets in Current Files"

    print_check "Scanning for exposed secrets in tracked files"

    cd "$REPO_ROOT"

    # Patterns to search for
    SECRET_PATTERNS=(
        "CLIENT_ID.*=.*[A-Za-z0-9]"
        "CLIENT_SECRET.*=.*[A-Za-z0-9]"
        "API_KEY.*=.*[A-Za-z0-9]"
        "JWT_SECRET.*=.*[A-Za-z0-9]"
        "DATABASE_URL.*=.*postgresql://"
        "GOOGLE_CLIENT_ID.*=.*[A-Za-z0-9]"
        "GITHUB_CLIENT_ID.*=.*[A-Za-z0-9]"
        "GEMINI_API_KEY.*=.*[A-Za-z0-9]"
    )

    SECRETS_FOUND=0

    for pattern in "${SECRET_PATTERNS[@]}"; do
        # Search in tracked files only
        if git grep -i "$pattern" -- '*.js' '*.ts' '*.json' '*.md' '*.yml' '*.yaml' 2>/dev/null | grep -v "example" | grep -q .; then
            print_fail "Found pattern '$pattern' in tracked files"
            git grep -n -i "$pattern" -- '*.js' '*.ts' '*.json' '*.md' '*.yml' '*.yaml' 2>/dev/null | grep -v "example" | head -5
            ((SECRETS_FOUND++))
        fi
    done

    if [ $SECRETS_FOUND -eq 0 ]; then
        print_pass "No secrets found in tracked files"
    fi

    print_check "Verifying .env files are not tracked"

    if git ls-files | grep -q "\.env$"; then
        print_fail ".env files are being tracked by git"
        git ls-files | grep "\.env$"
    else
        print_pass "No .env files in git tracking"
    fi

    print_check "Verifying .env is in .gitignore"

    if [ -f "$REPO_ROOT/.gitignore" ] && grep -q "^\.env$" "$REPO_ROOT/.gitignore"; then
        print_pass ".env is in .gitignore"
    else
        print_fail ".env is NOT in .gitignore"
    fi
}

################################################################################
# Check 2: Admin Middleware
################################################################################

check_admin_middleware() {
    print_section "Check 2: Admin Middleware"

    print_check "Admin middleware file exists"

    MIDDLEWARE_FILE="$REPO_ROOT/backend/src/middleware/admin.ts"

    if [ ! -f "$MIDDLEWARE_FILE" ]; then
        print_fail "Admin middleware file not found at: $MIDDLEWARE_FILE"
        return
    fi
    print_pass "Admin middleware file exists"

    print_check "Admin middleware exports requireAdmin function"

    if grep -q "export.*requireAdmin" "$MIDDLEWARE_FILE"; then
        print_pass "requireAdmin function is exported"
    else
        print_fail "requireAdmin function not found or not exported"
    fi

    print_check "Admin middleware checks user role"

    if grep -q "role.*admin\|admin.*role" "$MIDDLEWARE_FILE"; then
        print_pass "Admin role check found in middleware"
    else
        print_fail "Admin role check not found in middleware"
    fi

    print_check "Admin middleware is used in routes"

    ROUTES_USING_ADMIN=$(grep -r "requireAdmin" "$REPO_ROOT/backend/src/routes" 2>/dev/null | wc -l)

    if [ "$ROUTES_USING_ADMIN" -gt 0 ]; then
        print_pass "Admin middleware is used in $ROUTES_USING_ADMIN route file(s)"
    else
        print_warning "Admin middleware not found in any route files"
    fi
}

################################################################################
# Check 3: Token Naming Consistency
################################################################################

check_token_consistency() {
    print_section "Check 3: Token Naming Consistency"

    print_check "Checking for 'auth_token' vs 'token' inconsistencies"

    cd "$REPO_ROOT"

    # Count occurrences of different token naming patterns
    AUTH_TOKEN_COUNT=$(grep -r "auth_token" backend/src frontend/src 2>/dev/null | grep -v "node_modules" | wc -l || echo 0)
    TOKEN_COUNT=$(grep -r "\"token\"\|'token'\|\.token\s*:" backend/src frontend/src 2>/dev/null | grep -v "node_modules" | grep -v "auth_token" | wc -l || echo 0)

    print_info "Found 'auth_token' references: $AUTH_TOKEN_COUNT"
    print_info "Found 'token' references: $TOKEN_COUNT"

    if [ "$AUTH_TOKEN_COUNT" -gt 0 ] && [ "$TOKEN_COUNT" -gt 0 ]; then
        print_warning "Mixed token naming found - consider standardizing"
        echo ""
        echo "  Examples of 'auth_token':"
        grep -rn "auth_token" backend/src frontend/src 2>/dev/null | grep -v "node_modules" | head -3
        echo ""
        echo "  Examples of 'token':"
        grep -rn "\"token\"\|'token'\|\.token\s*:" backend/src frontend/src 2>/dev/null | grep -v "node_modules" | grep -v "auth_token" | head -3
    elif [ "$AUTH_TOKEN_COUNT" -gt 0 ]; then
        print_pass "Consistently using 'auth_token'"
    elif [ "$TOKEN_COUNT" -gt 0 ]; then
        print_pass "Consistently using 'token'"
    else
        print_warning "No token references found"
    fi

    print_check "Verifying cookie/header naming consistency"

    # Check auth.ts for token handling
    AUTH_FILE="$REPO_ROOT/backend/src/routes/auth.ts"
    if [ -f "$AUTH_FILE" ]; then
        if grep -q "auth_token" "$AUTH_FILE"; then
            print_pass "auth.ts uses 'auth_token'"
        else
            print_fail "auth.ts does not use 'auth_token'"
        fi
    else
        print_warning "auth.ts not found"
    fi
}

################################################################################
# Check 4: Database Indexes
################################################################################

check_database_indexes() {
    print_section "Check 4: Database Indexes"

    print_check "Checking for database index definitions"

    # Check schema file
    SCHEMA_FILE="$REPO_ROOT/backend/src/db/schema.sql"

    if [ ! -f "$SCHEMA_FILE" ]; then
        print_warning "Schema file not found at: $SCHEMA_FILE"
        return
    fi

    print_pass "Schema file found"

    # Check for specific indexes
    INDEXES_TO_CHECK=(
        "idx_users_email"
        "idx_learning_paths_user_id"
        "idx_projects_user_id"
        "idx_content_items_type"
        "idx_user_progress_user_id"
    )

    MISSING_INDEXES=0

    for index in "${INDEXES_TO_CHECK[@]}"; do
        if grep -q "$index" "$SCHEMA_FILE"; then
            print_pass "Index '$index' found"
        else
            print_fail "Index '$index' NOT found"
            ((MISSING_INDEXES++))
        fi
    done

    if [ $MISSING_INDEXES -eq 0 ]; then
        print_pass "All expected indexes are defined"
    fi

    print_check "Checking for index creation in migrations"

    MIGRATION_DIR="$REPO_ROOT/backend/src/db/migrations"

    if [ -d "$MIGRATION_DIR" ]; then
        INDEX_MIGRATIONS=$(grep -r "CREATE INDEX" "$MIGRATION_DIR" 2>/dev/null | wc -l || echo 0)
        print_info "Found $INDEX_MIGRATIONS index creation statement(s) in migrations"
    else
        print_warning "Migrations directory not found"
    fi
}

################################################################################
# Check 5: GitHub Username Column
################################################################################

check_github_username_column() {
    print_section "Check 5: GitHub Username Column"

    print_check "Checking for github_username column in schema"

    SCHEMA_FILE="$REPO_ROOT/backend/src/db/schema.sql"

    if [ ! -f "$SCHEMA_FILE" ]; then
        print_warning "Schema file not found"
        return
    fi

    if grep -A 20 "CREATE TABLE users" "$SCHEMA_FILE" | grep -q "github_username"; then
        print_pass "github_username column found in users table"
    else
        print_fail "github_username column NOT found in users table"
    fi

    print_check "Checking for github_username in TypeScript types"

    USER_TYPES=$(find "$REPO_ROOT/backend/src" -name "*.ts" -exec grep -l "github_username" {} \; 2>/dev/null | wc -l || echo 0)

    if [ "$USER_TYPES" -gt 0 ]; then
        print_pass "github_username found in $USER_TYPES TypeScript file(s)"
    else
        print_warning "github_username not found in TypeScript types"
    fi

    print_check "Checking GitHub OAuth integration"

    if grep -r "github_username.*profile\|profile.*github_username" "$REPO_ROOT/backend/src" 2>/dev/null | grep -q .; then
        print_pass "GitHub username extraction found in OAuth flow"
    else
        print_warning "GitHub username extraction not found in OAuth flow"
    fi
}

################################################################################
# Check 6: Git History
################################################################################

check_git_history() {
    print_section "Check 6: Git History"

    print_check "Checking for .env files in git history"

    cd "$REPO_ROOT"

    if git log --all --full-history --oneline -- '*.env' 2>/dev/null | grep -q .; then
        print_fail ".env files still exist in git history"
        echo ""
        echo "  Recent commits with .env files:"
        git log --all --full-history --oneline -- '*.env' 2>/dev/null | head -5
        echo ""
        print_info "Run ./scripts/cleanup-git-history.sh to remove them"
    else
        print_pass "No .env files found in git history"
    fi

    print_check "Checking for secret patterns in commit messages"

    SECRET_PATTERNS="CLIENT_ID|CLIENT_SECRET|API_KEY|JWT_SECRET"

    if git log --all --grep="$SECRET_PATTERNS" --oneline 2>/dev/null | grep -q .; then
        print_warning "Found potential secrets in commit messages"
        git log --all --grep="$SECRET_PATTERNS" --oneline 2>/dev/null | head -3
    else
        print_pass "No obvious secrets in commit messages"
    fi
}

################################################################################
# Check 7: Environment Variables
################################################################################

check_environment_setup() {
    print_section "Check 7: Environment Setup"

    print_check "Checking backend .env file exists"

    if [ -f "$REPO_ROOT/backend/.env" ]; then
        print_pass "Backend .env file exists"

        # Check for required variables
        REQUIRED_VARS=(
            "JWT_SECRET"
            "GOOGLE_CLIENT_ID"
            "GOOGLE_CLIENT_SECRET"
            "GITHUB_CLIENT_ID"
            "GITHUB_CLIENT_SECRET"
            "GEMINI_API_KEY"
            "DATABASE_URL"
        )

        MISSING_VARS=0

        for var in "${REQUIRED_VARS[@]}"; do
            if grep -q "^$var=" "$REPO_ROOT/backend/.env"; then
                print_pass "$var is set"
            else
                print_fail "$var is NOT set"
                ((MISSING_VARS++))
            fi
        done

        if [ $MISSING_VARS -eq 0 ]; then
            print_pass "All required environment variables are set"
        fi
    else
        print_fail "Backend .env file does NOT exist"
        print_info "Copy from .env.example if available"
    fi
}

################################################################################
# Check 8: Code Quality
################################################################################

check_code_quality() {
    print_section "Check 8: Code Quality Checks"

    print_check "Checking for console.log statements"

    CONSOLE_LOGS=$(grep -r "console\.log" "$REPO_ROOT/backend/src" "$REPO_ROOT/frontend/src" 2>/dev/null | grep -v "node_modules" | wc -l || echo 0)

    if [ "$CONSOLE_LOGS" -gt 5 ]; then
        print_warning "Found $CONSOLE_LOGS console.log statements (consider using proper logging)"
    else
        print_pass "Console.log usage is minimal ($CONSOLE_LOGS found)"
    fi

    print_check "Checking for TODO comments"

    TODOS=$(grep -r "TODO\|FIXME\|XXX" "$REPO_ROOT/backend/src" "$REPO_ROOT/frontend/src" 2>/dev/null | grep -v "node_modules" | wc -l || echo 0)

    if [ "$TODOS" -gt 0 ]; then
        print_info "Found $TODOS TODO/FIXME comments"
    else
        print_pass "No TODO comments found"
    fi

    print_check "Checking for error handling patterns"

    TRY_CATCH=$(grep -r "try {" "$REPO_ROOT/backend/src" 2>/dev/null | wc -l || echo 0)
    ERROR_RESPONSES=$(grep -r "sendError\|errorResponse" "$REPO_ROOT/backend/src" 2>/dev/null | wc -l || echo 0)

    print_info "Found $TRY_CATCH try-catch blocks"
    print_info "Found $ERROR_RESPONSES error response patterns"

    if [ "$TRY_CATCH" -gt 0 ] && [ "$ERROR_RESPONSES" -gt 0 ]; then
        print_pass "Error handling appears to be implemented"
    else
        print_warning "Limited error handling found"
    fi
}

################################################################################
# Summary Report
################################################################################

print_summary() {
    print_header "Verification Summary"

    echo -e "${BLUE}Total Checks:${NC}    $TOTAL_CHECKS"
    echo -e "${GREEN}Passed:${NC}          $PASSED_CHECKS"
    echo -e "${RED}Failed:${NC}          $FAILED_CHECKS"
    echo -e "${YELLOW}Warnings:${NC}        $WARNING_CHECKS"
    echo ""

    PASS_RATE=$((PASSED_CHECKS * 100 / TOTAL_CHECKS))

    if [ $FAILED_CHECKS -eq 0 ]; then
        echo -e "${GREEN}========================================${NC}"
        echo -e "${GREEN}  ✓ ALL CRITICAL CHECKS PASSED ($PASS_RATE%)${NC}"
        echo -e "${GREEN}========================================${NC}"
        exit 0
    elif [ $FAILED_CHECKS -le 3 ]; then
        echo -e "${YELLOW}========================================${NC}"
        echo -e "${YELLOW}  ⚠ SOME ISSUES FOUND ($PASS_RATE%)${NC}"
        echo -e "${YELLOW}========================================${NC}"
        exit 1
    else
        echo -e "${RED}========================================${NC}"
        echo -e "${RED}  ✗ CRITICAL ISSUES FOUND ($PASS_RATE%)${NC}"
        echo -e "${RED}========================================${NC}"
        exit 2
    fi
}

################################################################################
# Main Execution
################################################################################

main() {
    print_header "Shadow Nova - Urgent Fixes Verification"
    echo -e "${BLUE}Repository:${NC} $REPO_ROOT"
    echo -e "${BLUE}Date:${NC}       $(date '+%Y-%m-%d %H:%M:%S')"

    check_no_secrets_in_files
    check_admin_middleware
    check_token_consistency
    check_database_indexes
    check_github_username_column
    check_git_history
    check_environment_setup
    check_code_quality

    print_summary
}

# Run main function
main "$@"
