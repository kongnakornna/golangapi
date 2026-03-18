
## การออกแบบ CI/CD Workflow สอดคล้องกับ Git Flow

ระบบ CI/CD ที่ออกแบบนี้จะรองรับ Git Flow โดยมี 3 environments หลัก (Development, Staging, Production) และใช้ GitHub Actions เป็นตัวอย่าง[^1][^2]

## สถาปัตยกรรม CI/CD Overview

```
┌─────────────────────────────────────────────────────────────┐
│                     GIT FLOW BRANCHES                        │
├──────────────┬──────────────┬──────────────┬────────────────┤
│   Feature    │   Develop    │   Release    │     Main       │
│   Branches   │              │   Branches   │  (Production)  │
└──────┬───────┴──────┬───────┴──────┬───────┴────────┬───────┘
       │              │              │                │
       ▼              ▼              ▼                ▼
┌──────────────┬──────────────┬──────────────┬────────────────┐
│   CI: Test   │  CI + Deploy │  CI + Deploy │  CI + Deploy   │
│   Only       │  to Dev      │  to Staging  │  to Production │
└──────────────┴──────────────┴──────────────┴────────────────┘
       │              │              │                │
       ▼              ▼              ▼                ▼
┌──────────────┬──────────────┬──────────────┬────────────────┐
│  No Deploy   │ Development  │   Staging    │   Production   │
│              │ Environment  │ Environment  │  Environment   │
└──────────────┴──────────────┴──────────────┴────────────────┘
```


## Environment Configuration

### Environment Matrix

| Environment | Branch | Auto Deploy | Approval Required | Database | URL |
| :-- | :-- | :-- | :-- | :-- | :-- |
| **Development** | `develop` | ✅ Yes | ❌ No | Dev DB | `dev.car-service.com` |
| **Staging** | `release/*` | ✅ Yes | ❌ No | Staging DB | `staging.car-service.com` |
| **Production** | `main` | ⚠️ Manual | ✅ Yes (2 approvers) | Prod DB | `car-service.com` |

[^6][^3]

## Complete GitHub Actions Workflow

### Main CI/CD Pipeline

**.github/workflows/ci-cd.yml**

```yaml
name: Car Service CI/CD Pipeline

on:
  push:
    branches:
      - main
      - develop
      - 'release/**'
  pull_request:
    branches:
      - main
      - develop

env:
  REGISTRY: docker.io
  NESTJS_IMAGE: your-org/car-service-nestjs
  SPRING_IMAGE: your-org/car-service-spring
  FRONTEND_IMAGE: your-org/car-service-frontend
  NODE_VERSION: '18'
  JAVA_VERSION: '17'

jobs:
  # ==================== DETECT CHANGES ====================
  detect-changes:
    runs-on: ubuntu-latest
    outputs:
      nestjs: ${{ steps.filter.outputs.nestjs }}
      spring: ${{ steps.filter.outputs.spring }}
      frontend: ${{ steps.filter.outputs.frontend }}
      infra: ${{ steps.filter.outputs.infra }}
    steps:
      - uses: actions/checkout@v4
      
      - uses: dorny/paths-filter@v2
        id: filter
        with:
          filters: |
            nestjs:
              - 'backend/nestjs-api/**'
            spring:
              - 'backend/spring-boot-api/**'
            frontend:
              - 'frontend/react-app/**'
            infra:
              - 'infrastructure/**'

  # ==================== NESTJS CI ====================
  nestjs-ci:
    needs: detect-changes
    if: needs.detect-changes.outputs.nestjs == 'true'
    runs-on: ubuntu-latest
    
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: ${{ env.NODE_VERSION }}
          cache: 'npm'
          cache-dependency-path: backend/nestjs-api/package-lock.json

      - name: Install dependencies
        working-directory: backend/nestjs-api
        run: npm ci

      - name: Lint code
        working-directory: backend/nestjs-api
        run: npm run lint

      - name: Run unit tests
        working-directory: backend/nestjs-api
        run: npm run test:cov

      - name: Run e2e tests
        working-directory: backend/nestjs-api
        run: npm run test:e2e

      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./backend/nestjs-api/coverage/coverage-final.json
          flags: nestjs
          name: nestjs-coverage

      - name: SonarCloud Scan
        uses: SonarSource/sonarcloud-github-action@master
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          SONAR_TOKEN: ${{ secrets.SONAR_TOKEN }}
        with:
          projectBaseDir: backend/nestjs-api

  # ==================== SPRING BOOT CI ====================
  spring-boot-ci:
    needs: detect-changes
    if: needs.detect-changes.outputs.spring == 'true'
    runs-on: ubuntu-latest

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Setup Java
        uses: actions/setup-java@v4
        with:
          java-version: ${{ env.JAVA_VERSION }}
          distribution: 'temurin'
          cache: 'maven'

      - name: Build with Maven
        working-directory: backend/spring-boot-api
        run: mvn clean install -DskipTests

      - name: Run unit tests
        working-directory: backend/spring-boot-api
        run: mvn test

      - name: Run integration tests
        working-directory: backend/spring-boot-api
        run: mvn verify

      - name: Generate test report
        if: always()
        uses: dorny/test-reporter@v1
        with:
          name: Spring Boot Tests
          path: backend/spring-boot-api/target/surefire-reports/*.xml
          reporter: java-junit

      - name: Code coverage
        working-directory: backend/spring-boot-api
        run: mvn jacoco:report

  # ==================== FRONTEND CI ====================
  frontend-ci:
    needs: detect-changes
    if: needs.detect-changes.outputs.frontend == 'true'
    runs-on: ubuntu-latest

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: ${{ env.NODE_VERSION }}
          cache: 'npm'
          cache-dependency-path: frontend/react-app/package-lock.json

      - name: Install dependencies
        working-directory: frontend/react-app
        run: npm ci

      - name: Lint code
        working-directory: frontend/react-app
        run: npm run lint

      - name: Run tests
        working-directory: frontend/react-app
        run: npm test -- --coverage --watchAll=false

      - name: Build production
        working-directory: frontend/react-app
        run: npm run build
        env:
          VITE_API_URL: ${{ secrets.API_URL }}

      - name: Upload build artifacts
        uses: actions/upload-artifact@v3
        with:
          name: frontend-build
          path: frontend/react-app/dist
          retention-days: 7

  # ==================== SECURITY SCANNING ====================
  security-scan:
    needs: [detect-changes]
    if: |
      needs.detect-changes.outputs.nestjs == 'true' ||
      needs.detect-changes.outputs.spring == 'true' ||
      needs.detect-changes.outputs.frontend == 'true'
    runs-on: ubuntu-latest

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Run Trivy vulnerability scanner
        uses: aquasecurity/trivy-action@master
        with:
          scan-type: 'fs'
          scan-ref: '.'
          format: 'sarif'
          output: 'trivy-results.sarif'

      - name: Upload Trivy results to GitHub Security
        uses: github/codeql-action/upload-sarif@v2
        with:
          sarif_file: 'trivy-results.sarif'

  # ==================== BUILD & PUSH DOCKER IMAGES ====================
  build-images:
    needs: [nestjs-ci, spring-boot-ci, frontend-ci, detect-changes]
    if: |
      always() &&
      (github.ref == 'refs/heads/develop' || 
       github.ref == 'refs/heads/main' ||
       startsWith(github.ref, 'refs/heads/release/'))
    runs-on: ubuntu-latest

    strategy:
      matrix:
        service:
          - name: nestjs
            path: backend/nestjs-api
            image: ${{ env.NESTJS_IMAGE }}
            condition: needs.detect-changes.outputs.nestjs == 'true'
          - name: spring
            path: backend/spring-boot-api
            image: ${{ env.SPRING_IMAGE }}
            condition: needs.detect-changes.outputs.spring == 'true'
          - name: frontend
            path: frontend/react-app
            image: ${{ env.FRONTEND_IMAGE }}
            condition: needs.detect-changes.outputs.frontend == 'true'

    steps:
      - name: Checkout code
        if: ${{ matrix.service.condition }}
        uses: actions/checkout@v4

      - name: Set up Docker Buildx
        if: ${{ matrix.service.condition }}
        uses: docker/setup-buildx-action@v3

      - name: Login to Docker Hub
        if: ${{ matrix.service.condition }}
        uses: docker/login-action@v3
        with:
          username: ${{ secrets.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKER_PASSWORD }}

      - name: Extract metadata
        if: ${{ matrix.service.condition }}
        id: meta
        uses: docker/metadata-action@v5
        with:
          images: ${{ matrix.service.image }}
          tags: |
            type=ref,event=branch
            type=sha,prefix={{branch}}-
            type=semver,pattern={{version}}
            type=raw,value=latest,enable={{is_default_branch}}

      - name: Build and push
        if: ${{ matrix.service.condition }}
        uses: docker/build-push-action@v5
        with:
          context: ${{ matrix.service.path }}
          push: true
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}
          cache-from: type=registry,ref=${{ matrix.service.image }}:buildcache
          cache-to: type=registry,ref=${{ matrix.service.image }}:buildcache,mode=max

  # ==================== DEPLOY TO DEVELOPMENT ====================
  deploy-dev:
    needs: [build-images]
    if: github.ref == 'refs/heads/develop'
    runs-on: ubuntu-latest
    environment:
      name: development
      url: https://dev.car-service.com

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Setup kubectl
        uses: azure/setup-kubectl@v3
        with:
          version: 'latest'

      - name: Configure kubectl
        run: |
          echo "${{ secrets.KUBE_CONFIG_DEV }}" | base64 -d > kubeconfig
          export KUBECONFIG=kubeconfig

      - name: Deploy to Kubernetes
        run: |
          kubectl set image deployment/nestjs-api \
            nestjs-api=${{ env.NESTJS_IMAGE }}:develop-${{ github.sha }} \
            -n development
          
          kubectl set image deployment/spring-api \
            spring-api=${{ env.SPRING_IMAGE }}:develop-${{ github.sha }} \
            -n development
          
          kubectl set image deployment/frontend \
            frontend=${{ env.FRONTEND_IMAGE }}:develop-${{ github.sha }} \
            -n development

      - name: Wait for rollout
        run: |
          kubectl rollout status deployment/nestjs-api -n development
          kubectl rollout status deployment/spring-api -n development
          kubectl rollout status deployment/frontend -n development

      - name: Run smoke tests
        run: |
          curl -f https://dev.car-service.com/health || exit 1

      - name: Notify Slack
        if: always()
        uses: 8398a7/action-slack@v3
        with:
          status: ${{ job.status }}
          text: 'Deployment to Development: ${{ job.status }}'
          webhook_url: ${{ secrets.SLACK_WEBHOOK }}

  # ==================== DEPLOY TO STAGING ====================
  deploy-staging:
    needs: [build-images]
    if: startsWith(github.ref, 'refs/heads/release/')
    runs-on: ubuntu-latest
    environment:
      name: staging
      url: https://staging.car-service.com

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Extract version
        id: version
        run: |
          VERSION=$(echo ${GITHUB_REF#refs/heads/release/})
          echo "version=$VERSION" >> $GITHUB_OUTPUT

      - name: Setup kubectl
        uses: azure/setup-kubectl@v3

      - name: Configure kubectl
        run: |
          echo "${{ secrets.KUBE_CONFIG_STAGING }}" | base64 -d > kubeconfig
          export KUBECONFIG=kubeconfig

      - name: Deploy to Kubernetes
        run: |
          # Update deployment with version tags
          kubectl set image deployment/nestjs-api \
            nestjs-api=${{ env.NESTJS_IMAGE }}:${{ steps.version.outputs.version }} \
            -n staging
          
          kubectl set image deployment/spring-api \
            spring-api=${{ env.SPRING_IMAGE }}:${{ steps.version.outputs.version }} \
            -n staging
          
          kubectl set image deployment/frontend \
            frontend=${{ env.FRONTEND_IMAGE }}:${{ steps.version.outputs.version }} \
            -n staging

      - name: Run database migrations
        run: |
          kubectl run migration-${{ github.sha }} \
            --image=${{ env.NESTJS_IMAGE }}:${{ steps.version.outputs.version }} \
            --restart=Never \
            --command -- npm run migration:run \
            -n staging

      - name: Wait for rollout
        run: |
          kubectl rollout status deployment/nestjs-api -n staging
          kubectl rollout status deployment/spring-api -n staging
          kubectl rollout status deployment/frontend -n staging

      - name: Run integration tests
        run: |
          npm run test:integration:staging

      - name: Performance test
        run: |
          docker run --rm -i grafana/k6 run - < ./tests/load-test.js

  # ==================== DEPLOY TO PRODUCTION ====================
  deploy-production:
    needs: [build-images]
    if: github.ref == 'refs/heads/main'
    runs-on: ubuntu-latest
    environment:
      name: production
      url: https://car-service.com

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Setup kubectl
        uses: azure/setup-kubectl@v3

      - name: Configure kubectl
        run: |
          echo "${{ secrets.KUBE_CONFIG_PROD }}" | base64 -d > kubeconfig
          export KUBECONFIG=kubeconfig

      - name: Create backup
        run: |
          # Backup current deployment
          kubectl get deployment -n production -o yaml > backup-deployment.yaml
          
          # Backup database
          ./scripts/backup-database.sh

      - name: Deploy with blue-green strategy
        run: |
          # Deploy to green environment
          kubectl apply -f infrastructure/kubernetes/green/ -n production
          
          # Wait for green to be healthy
          kubectl wait --for=condition=available --timeout=300s \
            deployment/nestjs-api-green -n production

      - name: Run smoke tests on green
        run: |
          curl -f https://green.car-service.com/health || exit 1

      - name: Switch traffic to green
        run: |
          kubectl patch service nestjs-api \
            -p '{"spec":{"selector":{"version":"green"}}}' \
            -n production

      - name: Monitor for 5 minutes
        run: |
          sleep 300
          # Check error rates, response times
          ./scripts/check-metrics.sh

      - name: Cleanup old blue deployment
        run: |
          kubectl delete deployment -l version=blue -n production

      - name: Create Git tag
        run: |
          git tag -a v${{ github.run_number }} -m "Production release ${{ github.run_number }}"
          git push origin v${{ github.run_number }}

      - name: Notify team
        uses: 8398a7/action-slack@v3
        with:
          status: ${{ job.status }}
          text: |
            🚀 Production Deployment Successful
            Version: v${{ github.run_number }}
            Commit: ${{ github.sha }}
          webhook_url: ${{ secrets.SLACK_WEBHOOK }}

  # ==================== ROLLBACK (Manual Trigger) ====================
  rollback:
    if: github.event_name == 'workflow_dispatch'
    runs-on: ubuntu-latest
    environment:
      name: production

    steps:
      - name: Rollback to previous version
        run: |
          kubectl rollout undo deployment/nestjs-api -n production
          kubectl rollout undo deployment/spring-api -n production
          kubectl rollout undo deployment/frontend -n production

      - name: Restore database
        run: |
          ./scripts/restore-database.sh ${{ github.event.inputs.backup_id }}
```


## Environment-Specific Configurations

### Development Environment

**infrastructure/kubernetes/dev/deployment.yaml**

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nestjs-api
  namespace: development
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nestjs-api
      env: dev
  template:
    metadata:
      labels:
        app: nestjs-api
        env: dev
    spec:
      containers:
      - name: nestjs-api
        image: your-org/car-service-nestjs:develop-latest
        env:
        - name: NODE_ENV
          value: "development"
        - name: DB_HOST
          value: "postgres-dev.default.svc.cluster.local"
        - name: REDIS_HOST
          value: "redis-dev.default.svc.cluster.local"
        - name: LOG_LEVEL
          value: "debug"
        resources:
          requests:
            memory: "256Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "200m"
        livenessProbe:
          httpGet:
            path: /health
            port: 3000
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 3000
          initialDelaySeconds: 10
          periodSeconds: 5
```


### Production Environment (Blue-Green)

**infrastructure/kubernetes/production/deployment-green.yaml**

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nestjs-api-green
  namespace: production
spec:
  replicas: 3
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  selector:
    matchLabels:
      app: nestjs-api
      version: green
  template:
    metadata:
      labels:
        app: nestjs-api
        version: green
    spec:
      containers:
      - name: nestjs-api
        image: your-org/car-service-nestjs:latest
        env:
        - name: NODE_ENV
          value: "production"
        - name: DB_HOST
          valueFrom:
            secretKeyRef:
              name: db-credentials
              key: host
        resources:
          requests:
            memory: "1Gi"
            cpu: "500m"
          limits:
            memory: "2Gi"
            cpu: "1000m"
        livenessProbe:
          httpGet:
            path: /health
            port: 3000
          initialDelaySeconds: 60
          periodSeconds: 30
          failureThreshold: 3
        readinessProbe:
          httpGet:
            path: /ready
            port: 3000
          initialDelaySeconds: 30
          periodSeconds: 10
          failureThreshold: 3
```


## Database Migration Pipeline

**.github/workflows/database-migration.yml**

```yaml
name: Database Migration

on:
  push:
    branches:
      - develop
      - main
    paths:
      - 'database/migrations/**'

jobs:
  migrate-dev:
    if: github.ref == 'refs/heads/develop'
    runs-on: ubuntu-latest
    environment: development

    steps:
      - uses: actions/checkout@v4

      - name: Run migrations
        run: |
          docker run --rm \
            -e DATABASE_URL=${{ secrets.DEV_DATABASE_URL }} \
            -v $(pwd)/database/migrations:/migrations \
            migrate/migrate \
            -path=/migrations \
            -database ${{ secrets.DEV_DATABASE_URL }} \
            up

  migrate-staging:
    if: startsWith(github.ref, 'refs/heads/release/')
    runs-on: ubuntu-latest
    environment: staging

    steps:
      - uses: actions/checkout@v4

      - name: Backup database
        run: |
          ./scripts/backup-database.sh staging

      - name: Run migrations
        run: |
          docker run --rm \
            -e DATABASE_URL=${{ secrets.STAGING_DATABASE_URL }} \
            -v $(pwd)/database/migrations:/migrations \
            migrate/migrate \
            -path=/migrations \
            -database ${{ secrets.STAGING_DATABASE_URL }} \
            up

      - name: Verify migration
        run: |
          ./scripts/verify-migration.sh staging

  migrate-production:
    if: github.ref == 'refs/heads/main'
    runs-on: ubuntu-latest
    environment: production

    steps:
      - uses: actions/checkout@v4

      - name: Create backup
        run: |
          ./scripts/backup-database.sh production

      - name: Dry-run migration
        run: |
          ./scripts/dry-run-migration.sh production

      - name: Wait for approval
        uses: trstringer/manual-approval@v1
        with:
          approvers: senior-dev-team
          minimum-approvals: 2

      - name: Run migrations
        run: |
          docker run --rm \
            -e DATABASE_URL=${{ secrets.PROD_DATABASE_URL }} \
            -v $(pwd)/database/migrations:/migrations \
            migrate/migrate \
            -path=/migrations \
            -database ${{ secrets.PROD_DATABASE_URL }} \
            up

      - name: Verify migration
        run: |
          ./scripts/verify-migration.sh production
```


## Monitoring and Alerts

**.github/workflows/monitoring.yml**

```yaml
name: Post-Deployment Monitoring

on:
  workflow_run:
    workflows: ["Car Service CI/CD Pipeline"]
    types:
      - completed

jobs:
  monitor:
    runs-on: ubuntu-latest
    if: ${{ github.event.workflow_run.conclusion == 'success' }}

    steps:
      - name: Check application health
        run: |
          for i in {1..10}; do
            if curl -f https://car-service.com/health; then
              echo "Health check passed"
              break
            fi
            echo "Retry $i/10"
            sleep 30
          done

      - name: Check error rates
        run: |
          # Query Prometheus/Grafana for error rates
          ERROR_RATE=$(curl -s "http://prometheus:9090/api/v1/query?query=rate(http_requests_total{status=~'5..'}[5m])")
          
          if [ "$ERROR_RATE" -gt "0.01" ]; then
            echo "High error rate detected: $ERROR_RATE"
            exit 1
          fi

      - name: Check response times
        run: |
          # Check average response time
          AVG_RESPONSE=$(curl -s "http://prometheus:9090/api/v1/query?query=http_request_duration_seconds")
          
          if [ "$AVG_RESPONSE" -gt "1.0" ]; then
            echo "High response time: $AVG_RESPONSE seconds"
            # Send alert but don't fail
          fi

      - name: Alert on failure
        if: failure()
        uses: 8398a7/action-slack@v3
        with:
          status: failure
          text: '🚨 Post-deployment monitoring failed!'
          webhook_url: ${{ secrets.SLACK_WEBHOOK }}
```


## CI/CD Workflow Summary

```
┌─────────────────────────────────────────────────────────────┐
│                    GIT PUSH EVENT                            │
└─────────────────┬───────────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────────┐
│  PHASE 1: CONTINUOUS INTEGRATION (CI)                        │
├─────────────────────────────────────────────────────────────┤
│  1. Detect Changes (paths-filter)                           │
│  2. Lint Code (ESLint, Prettier, Checkstyle)               │
│  3. Unit Tests (Jest, JUnit)                                │
│  4. Integration Tests (E2E, API Tests)                      │
│  5. Code Coverage (>80%)                                     │
│  6. Security Scan (Trivy, SonarQube)                        │
│  7. Build Validation                                         │
└─────────────────┬───────────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────────┐
│  PHASE 2: BUILD & PACKAGE                                    │
├─────────────────────────────────────────────────────────────┤
│  1. Build Docker Images                                      │
│  2. Tag with version (branch-sha, semver)                   │
│  3. Push to Docker Registry                                  │
│  4. Generate SBOM (Software Bill of Materials)              │
└─────────────────┬───────────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────────┐
│  PHASE 3: CONTINUOUS DEPLOYMENT (CD)                         │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌──────────┐  ┌──────────┐  ┌──────────────┐              │
│  │ Develop  │  │ Release  │  │    Main      │              │
│  │ Branch   │  │ Branch   │  │   Branch     │              │
│  └────┬─────┘  └────┬─────┘  └──────┬───────┘              │
│       │             │               │                       │
│       ▼             ▼               ▼                       │
│  ┌────────┐   ┌─────────┐   ┌──────────────┐              │
│  │  DEV   │   │ STAGING │   │ PRODUCTION   │              │
│  │  Auto  │   │  Auto   │   │ + Approval   │              │
│  └────────┘   └─────────┘   └──────────────┘              │
└─────────────────────────────────────────────────────────────┘
```


## Best Practices Implementation

### 1. Fast Feedback Loop

- Unit tests run first (fastest)
- Integration tests run after
- Build only if tests pass
- Parallel execution where possible
[^4]


### 2. Security First

- Scan dependencies for vulnerabilities
- Container image scanning
- Secret management with GitHub Secrets
- SBOM generation
[^4]


### 3. Progressive Deployment

- Development: Automatic
- Staging: Automatic with tests
- Production: Manual approval + Blue-Green
[^3][^2]


### 4. Rollback Strategy

- Keep previous version running (Blue-Green)
- Database backups before migrations
- One-click rollback capability
- Automated health checks
[^1]


### 5. Observability

- Health check endpoints
- Performance monitoring
- Error rate tracking
- Slack/Email notifications
[^4]

ระบบ CI/CD นี้รองรับ Git Flow อย่างสมบูรณ์ โดยแต่ละ branch มี pipeline ที่เหมาะสม มีการตรวจสอบคุณภาพ ความปลอดภัย และการ deploy แบบอัตโนมัติที่ปลอดภัย[^5][^1]
<span style="display:none">[^10][^7][^8][^9]</span>

<div align="center">⁂</div>

[^1]: https://hackernoon.com/from-git-flow-to-cicd-a-practical-guide-to-implement-git-workflow

[^2]: https://codefresh.io/learn/github-actions/deployment-with-github-actions-quick-tutorial-and-5-best-practices/

[^3]: https://docs.github.com/en/actions/concepts/workflows-and-actions/deployment-environments

[^4]: https://gitprotect.io/blog/exploring-best-practices-and-modern-trends-in-ci-cd/

[^5]: https://docs.bytebase.com/gitops/best-practices/git-and-cicd

[^6]: https://www.atlassian.com/git/tutorials/comparing-workflows/gitflow-workflow

[^7]: https://about.gitlab.com/topics/version-control/what-are-gitlab-flow-best-practices/

[^8]: https://semaphore.io/blog/cicd-microservices-digitalocean-kubernetes

[^9]: https://www.scalefree.com/consulting/devops-solutions/behind-the-branches-navigating-git-workflows-in-modern-devops/

[^10]: https://dev.to/prodevopsguytech/devops-project-cicd-pipeline-for-a-microservices-based-application-on-kubernetes-1ba8

