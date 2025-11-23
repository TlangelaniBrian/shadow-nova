# Shadow Nova - CI/CD Guide

## GitHub Actions Workflow

Shadow Nova uses GitHub Actions for continuous integration and deployment.

### Workflow Overview

The `.github/workflows/build-deploy.yml` workflow handles:

1. **Build Frontend** (Node + pnpm)
2. **Build Backend** (Go)
3. **Docker Image Creation** (Both services)
4. **Multi-Environment Deployment** (DEV → Staging → Production)

### Environments

- **DEV**: Automatic deployment on push to `develop`
- **Staging**: Automatic deployment after DEV success
- **Production**: Manual approval required, triggered on merge to `main`

### Workflow Triggers

```yaml
# Manual trigger
workflow_dispatch

# Automatic on push
push:
  branches:
    - main
    - develop

# On PR merge
pull_request:
  types: [closed]
  branches:
    - main
    - develop
```

### Setup Required

#### 1. GitHub Secrets

Go to **Settings → Secrets and variables → Actions** and add:

**Not strictly required** (uses GitHub Container Registry):

- Workflow uses `GITHUB_TOKEN` automatically
- Images pushed to `ghcr.io/your-username/shadow-nova`

#### 2. GitHub Environments

Create environments in **Settings → Environments**:

- `dev` - No approval required
- `staging` - Optional: require 1 reviewer
- `production` - **Required**: Add protection rules
  - Required reviewers: 2
  - Prevent self-review

#### 3. Enable GitHub Packages

1. Go to **Settings → Actions → General**
2. Under **Workflow permissions**, select:
   - ✅ Read and write permissions
   - ✅ Allow GitHub Actions to create and approve pull requests

### Deployment Integration

The workflow creates deployment placeholders. Integrate with your deployment tool:

#### ArgoCD Example

```yaml
deploy-dev:
  steps:
    - name: Deploy via ArgoCD
      uses: argoproj/argocd-action@v1
      with:
        version: latest
        command: app sync shadow-nova-dev
```

#### Kubernetes Example

```yaml
deploy-dev:
  steps:
    - name: Deploy to Kubernetes
      uses: azure/k8s-deploy@v4
      with:
        manifests: |
          k8s/deployment.yaml
        images: |
          ${{ needs.docker-frontend.outputs.image }}
          ${{ needs.docker-backend.outputs.image }}
```

#### Rancher Example

```yaml
deploy-dev:
  steps:
    - name: Deploy to Rancher
      env:
        RANCHER_URL: ${{ secrets.RANCHER_URL }}
        RANCHER_TOKEN: ${{ secrets.RANCHER_TOKEN }}
      run: |
        kubectl --kubeconfig=kubeconfig.yaml \
          set image deployment/frontend \
          frontend=${{ needs.docker-frontend.outputs.image }}
```

### Docker Images

Images are tagged with:

- Branch name (e.g., `main`, `develop`)
- Commit SHA (e.g., `main-a1b2c3d`)
- Version from `package.json` (semver)

**Example:**

```
ghcr.io/your-username/shadow-nova/frontend:main
ghcr.io/your-username/shadow-nova/frontend:main-a1b2c3d
ghcr.io/your-username/shadow-nova/frontend:1.0.0
```

### Build Process

#### Frontend Build

1. Enable corepack for pnpm
2. Install dependencies with `pnpm install --frozen-lockfile`
3. Build with `pnpm run build:docker`
4. Create multi-stage Docker image
5. Push to GitHub Container Registry

#### Backend Build

1. Set up Go 1.23
2. Download and verify dependencies
3. Run unit tests
4. Build static binary with `CGO_ENABLED=0`
5. Create minimal Docker image (scratch-based)
6. Push to GitHub Container Registry

### Running Workflow

#### Manual Trigger

1. Go to **Actions** tab
2. Select **Build and Deploy Shadow Nova**
3. Click **Run workflow**
4. Select branch
5. Click **Run workflow**

#### Automatic Triggers

- **Push to develop**: Builds and deploys to DEV + Staging
- **Merge PR to main**: Builds, deploys to Production, creates release

### Security Features

✅ **nginx unprivileged** - Frontend runs on port 8080 (non-root)
✅ **Minimal images** - Backend uses `scratch`, frontend uses Alpine
✅ **Build caching** - GitHub Actions cache for faster builds
✅ **Multi-stage builds** - Reduces final image size
✅ **Signed commits** - Optional: require commit signing

### Local Testing

Test Docker builds locally before pushing:

```bash
# Build frontend
docker build -f Dockerfile.frontend -t shadow-nova-frontend:local ./frontend

# Build backend
docker build -f Dockerfile.backend -t shadow-nova-backend:local ./backend

# Run with docker-compose
docker-compose up --build
```

### Monitoring Workflow

View workflow runs:

1. Go to **Actions** tab
2. Click on a workflow run
3. View job logs and artifacts

### Artifacts

The workflow saves build artifacts:

- **frontend-dist**: Built frontend (7 days retention)
- **backend-binary**: Compiled Go binary (7 days retention)

Download artifacts from the workflow run page.

### Troubleshooting

#### Build Fails: "pnpm: command not found"

→ Corepack step failed. Check Node.js version (20.x required)

#### Docker push fails: "denied"

→ Check **Settings → Actions → General → Workflow permissions**

#### Deployment doesn't trigger

→ Check branch protection rules and PR merge status

#### Tests fail

→ Run tests locally: `cd backend && go test -v ./...`

### Cost Optimization

GitHub Actions is free for public repos. For private repos:

- 2,000 minutes/month free
- Build caching reduces minutes used
- Matrix strategies can speed up builds

### Next Steps

1. ✅ Commit workflow file
2. ✅ Configure GitHub environments
3. ✅ Enable GitHub Packages
4. ⬜ Add deployment integration (ArgoCD/Kubernetes/Rancher)
5. ⬜ Set up monitoring (Sentry, Datadog, etc.)
6. ⬜ Configure branch protection rules
