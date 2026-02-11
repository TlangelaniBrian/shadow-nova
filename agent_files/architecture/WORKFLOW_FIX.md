# GitHub Workflow Fix Guide

## Current Issue

The GitHub Actions workflow is failing with:
```
The job was not started because your account is locked due to a billing issue.
```

## Root Cause

This is **not a code issue** - it's an account configuration problem. The workflow code is correct, but GitHub requires billing to be set up for Actions to run.

## How to Fix

### Option 1: Add Payment Method (Recommended)

1. Go to [GitHub Billing Settings](https://github.com/settings/billing)
2. Click "Add payment information"
3. Enter your payment method
4. GitHub Actions has a generous free tier:
   - **2,000 minutes/month** for private repos
   - **Unlimited** for public repos

### Option 2: Make Repository Public

If this is a public repository, GitHub Actions should work without billing:

1. Go to repository Settings
2. Scroll to "Danger Zone"
3. Click "Change repository visibility"
4. Select "Public"

### Option 3: Use Alternative CI/CD

If you cannot set up billing, consider:
- **GitLab CI** (3,400 free minutes/month)
- **CircleCI** (6,000 free minutes/month)
- **Travis CI** (1,000 free credits for public repos)

## Workflow Improvements Made

I've fixed several technical issues in the workflow:

### 1. Fixed Release Creation
**Before:**
```yaml
uses: actions/create-release@v1  # Deprecated
```

**After:**
```yaml
uses: softprops/action-gh-release@v2  # Current version
```

### 2. Fixed Production Deploy Condition
**Before:**
```yaml
if: github.event.pull_request.merged == true  # Won't work on push events
```

**After:**
```yaml
if: github.event_name == 'push'  # Triggers on direct pushes to main
```

### 3. Added pnpm Lockfile Check
The workflow now checks if `pnpm-lock.yaml` exists before using `--frozen-lockfile`, preventing failures on fresh clones.

## Testing the Workflow

Once billing is resolved, test with:

```bash
# Trigger workflow manually
gh workflow run build-deploy.yml

# Or push to main/develop
git push origin main
```

## Verifying Success

After fixing billing, you should see:
1. All jobs turn green in Actions tab
2. Docker images published to GitHub Container Registry
3. Release created (if commit message starts with "Release")

## Need Help?

- Check workflow runs: https://github.com/TlangelaniBrian/shadow-nova/actions
- GitHub Actions docs: https://docs.github.com/en/actions
- Billing settings: https://github.com/settings/billing
