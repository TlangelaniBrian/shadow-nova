# Shadow Nova - AWS Deployment Guide

This guide outlines how to deploy the Shadow Nova application to AWS using a student account.

## Prerequisites

- AWS Account (Student/Free Tier)
- AWS CLI installed and configured
- Docker installed

## 1. Frontend Deployment (AWS Amplify)

AWS Amplify is the easiest way to deploy Vue.js applications.

1.  **Push your code to GitHub/GitLab/Bitbucket**.
2.  **Log in to AWS Console** and go to **AWS Amplify**.
3.  Click **"New App"** -> **"Host web app"**.
4.  Connect your repository.
5.  Amplify will automatically detect the settings:
    - **Build command**: `pnpm run build`
    - **Output directory**: `dist`
6.  Click **"Save and Deploy"**.

## 2. Backend Deployment (AWS App Runner or EC2)

For a student account, **App Runner** is easiest (container-based), but **EC2** is cheaper (free tier eligible).

### Option A: App Runner (Recommended for Containers)

1.  **Push your Backend Docker Image** to Amazon ECR (Elastic Container Registry).
    ```bash
    aws ecr create-repository --repository-name shadow-nova-backend
    docker build -t shadow-nova-backend -f Dockerfile.backend .
    aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin <your-account-id>.dkr.ecr.us-east-1.amazonaws.com
    docker tag shadow-nova-backend:latest <your-account-id>.dkr.ecr.us-east-1.amazonaws.com/shadow-nova-backend:latest
    docker push <your-account-id>.dkr.ecr.us-east-1.amazonaws.com/shadow-nova-backend:latest
    ```
2.  Go to **AWS App Runner** console.
3.  Create service -> Source: **Container Registry**.
4.  Select your image URI.
5.  Configure settings:
    - Port: `3000`
    - Environment Variables: `DATABASE_URL` (see below).

### Option B: EC2 (Free Tier)

1.  Launch an **t2.micro** or **t3.micro** instance (Ubuntu).
2.  SSH into the instance.
3.  Install Docker.
4.  Run the container:
    ```bash
    docker run -d -p 3000:3000 -e DATABASE_URL=... <your-image>
    ```

## 3. Database (Amazon RDS)

1.  Go to **RDS Console**.
2.  Create Database -> **PostgreSQL**.
3.  Select **"Free Tier"** template.
4.  Settings:
    - DB Instance Class: `db.t3.micro`
    - Storage: 20GB (Autoscaling off)
    - Public Access: **No** (unless you need to connect from local PC, then Yes + Security Group IP whitelist).
5.  Create Database.
6.  Copy the **Endpoint** URL.
7.  Construct your connection string: `postgres://username:password@endpoint:5432/shadownova`.

## 4. Connecting It All

1.  Update your Backend App Runner/EC2 environment variable `DATABASE_URL` with the RDS connection string.
2.  Update your Frontend environment variable `VITE_API_URL` to point to your Backend URL.
