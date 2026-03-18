<img src="https://r2cdn.perplexity.ai/pplx-full-logo-primary-dark%402x.png" style="height:64px;margin-right:32px"/>

## การจัดการ Hotfix และ Rollback เมื่อ Deployment ล้มเหลว

การจัดการ hotfix และ rollback ที่มีประสิทธิภาพเป็นสิ่งสำคัญสำหรับระบบ production  ต่อไปนี้คือกลยุทธ์และขั้นตอนที่ครบถ้วน[^1][^2]

## Hotfix Strategy

### 1. Hotfix Workflow

```
Production Issue Detected
         ↓
    Assess Severity
         ↓
    ┌────────────────┐
    │ Critical Bug?  │
    └────┬──────┬────┘
         │      │
    Yes  │      │  No → Create normal bugfix
         ↓
Create Hotfix Branch from main
         ↓
Fix → Test → Review → Merge
         ↓
Deploy to Production (Fast-track)
         ↓
Merge back to develop
         ↓
Post-mortem Analysis
```


### 2. Hotfix Branch Strategy

**สร้าง Hotfix Branch**

```bash
# 1. Create hotfix branch from main (production)
git checkout main
git pull origin main
git checkout -b hotfix/v1.0.1-fix-payment-timeout

# 2. Apply the fix
# Edit files...
git add .
git commit -m "fix(payment): increase timeout from 10s to 30s

- Update payment gateway timeout configuration
- Add retry logic for timeout scenarios
- Add monitoring alert for payment failures

Fixes: #TICKET-301
Severity: Critical
Impact: Payment failures affecting 15% of transactions"

# 3. Bump patch version
npm version patch  # 1.0.0 → 1.0.1

# 4. Push hotfix branch
git push origin hotfix/v1.0.1-fix-payment-timeout
```


### 3. Fast-Track Hotfix Pipeline

**.github/workflows/hotfix.yml**

```yaml
name: Hotfix Deployment

on:
  push:
    branches:
      - 'hotfix/**'

env:
  SEVERITY: critical

jobs:
  # ==================== VALIDATE HOTFIX ====================
  validate-hotfix:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Extract hotfix info
        id: hotfix
        run: |
          BRANCH_NAME=${GITHUB_REF#refs/heads/}
          VERSION=$(echo $BRANCH_NAME | sed 's/hotfix\///')
          echo "version=$VERSION" >> $GITHUB_OUTPUT
          echo "branch=$BRANCH_NAME" >> $GITHUB_OUTPUT

      - name: Verify hotfix branch from main
        run: |
          git fetch origin main
          MERGE_BASE=$(git merge-base HEAD origin/main)
          MAIN_HEAD=$(git rev-parse origin/main)
          
          if [ "$MERGE_BASE" != "$MAIN_HEAD" ]; then
            echo "❌ Hotfix must be created from main branch"
            exit 1
          fi

      - name: Check for breaking changes
        run: |
          # Ensure no schema changes in hotfix
          if git diff origin/main --name-only | grep -E "migration|schema"; then
            echo "⚠️  Warning: Database changes detected in hotfix"
            echo "This requires manual approval"
          fi

  # ==================== EXPEDITED TESTING ====================
  expedited-tests:
    needs: validate-hotfix
    runs-on: ubuntu-latest
    timeout-minutes: 15  # Fast timeout for hotfixes

    steps:
      - uses: actions/checkout@v4

      - name: Setup environment
        uses: actions/setup-node@v4
        with:
          node-version: '18'

      - name: Install dependencies
        run: npm ci
        working-directory: backend/nestjs-api

      - name: Run critical tests only
        run: |
          # Run only tests related to the fix
          npm test -- --testPathPattern="payment" --coverage
        working-directory: backend/nestjs-api

      - name: Security scan
        uses: aquasecurity/trivy-action@master
        with:
          scan-type: 'fs'
          severity: 'CRITICAL,HIGH'

  # ==================== BUILD HOTFIX ====================
  build-hotfix:
    needs: expedited-tests
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v4

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Login to registry
        uses: docker/login-action@v3
        with:
          username: ${{ secrets.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKER_PASSWORD }}

      - name: Extract version
        id: version
        run: |
          VERSION=$(echo ${GITHUB_REF#refs/heads/hotfix/})
          echo "version=$VERSION" >> $GITHUB_OUTPUT

      - name: Build and push hotfix image
        uses: docker/build-push-action@v5
        with:
          context: backend/nestjs-api
          push: true
          tags: |
            your-org/car-service-nestjs:${{ steps.version.outputs.version }}
            your-org/car-service-nestjs:hotfix-latest
          labels: |
            hotfix=true
            severity=critical
            version=${{ steps.version.outputs.version }}

  # ==================== DEPLOY TO STAGING FOR VALIDATION ====================
  deploy-staging-validation:
    needs: build-hotfix
    runs-on: ubuntu-latest
    environment:
      name: staging-hotfix
      url: https://hotfix-staging.car-service.com

    steps:
      - uses: actions/checkout@v4

      - name: Setup kubectl
        uses: azure/setup-kubectl@v3

      - name: Deploy to isolated staging
        run: |
          # Deploy to isolated namespace for hotfix validation
          kubectl create namespace hotfix-validation --dry-run=client -o yaml | kubectl apply -f -
          
          kubectl set image deployment/nestjs-api \
            nestjs-api=your-org/car-service-nestjs:${{ needs.validate-hotfix.outputs.version }} \
            -n hotfix-validation

      - name: Wait for deployment
        run: |
          kubectl rollout status deployment/nestjs-api -n hotfix-validation --timeout=5m

      - name: Run smoke tests
        run: |
          # Test the specific fix
          npm run test:hotfix:payment
          
      - name: Load test the fix
        run: |
          # Simulate production load
          k6 run --vus 100 --duration 2m tests/payment-load-test.js

  # ==================== PRODUCTION DEPLOYMENT ====================
  deploy-production:
    needs: deploy-staging-validation
    runs-on: ubuntu-latest
    environment:
      name: production-hotfix
      url: https://car-service.com

    steps:
      - uses: actions/checkout@v4

      - name: Notify team - Deployment starting
        uses: 8398a7/action-slack@v3
        with:
          status: custom
          custom_payload: |
            {
              text: "🚨 HOTFIX DEPLOYMENT STARTING",
              attachments: [{
                color: 'warning',
                text: `Version: ${{ needs.validate-hotfix.outputs.version }}\nSeverity: Critical\nDeploying by: ${{ github.actor }}`
              }]
            }
          webhook_url: ${{ secrets.SLACK_WEBHOOK }}

      - name: Create pre-deployment backup
        run: |
          # Backup current state
          kubectl get all -n production -o yaml > backup-$(date +%s).yaml
          
          # Backup database
          ./scripts/backup-database.sh production hotfix-backup

      - name: Setup kubectl
        uses: azure/setup-kubectl@v3

      - name: Configure kubectl
        run: |
          echo "${{ secrets.KUBE_CONFIG_PROD }}" | base64 -d > kubeconfig
          export KUBECONFIG=kubeconfig

      - name: Deploy hotfix with canary strategy
        run: |
          # Deploy to 10% of traffic first
          kubectl set image deployment/nestjs-api-canary \
            nestjs-api=your-org/car-service-nestjs:${{ needs.validate-hotfix.outputs.version }} \
            -n production
          
          # Wait for canary to be healthy
          kubectl rollout status deployment/nestjs-api-canary -n production

      - name: Monitor canary for 5 minutes
        run: |
          ./scripts/monitor-canary.sh 300
          
          # Check metrics
          ERROR_RATE=$(./scripts/get-error-rate.sh)
          
          if (( $(echo "$ERROR_RATE > 0.01" | bc -l) )); then
            echo "❌ High error rate detected: $ERROR_RATE"
            echo "ROLLBACK_NEEDED=true" >> $GITHUB_ENV
            exit 1
          fi

      - name: Roll out to all instances
        if: env.ROLLBACK_NEEDED != 'true'
        run: |
          kubectl set image deployment/nestjs-api \
            nestjs-api=your-org/car-service-nestjs:${{ needs.validate-hotfix.outputs.version }} \
            -n production
          
          kubectl rollout status deployment/nestjs-api -n production

      - name: Verify production deployment
        run: |
          # Run production smoke tests
          ./scripts/production-smoke-tests.sh
          
          # Verify the specific fix
          curl -f https://car-service.com/api/v1/health
          ./scripts/verify-payment-fix.sh

      - name: Create Git tag
        run: |
          git tag -a ${{ needs.validate-hotfix.outputs.version }} \
            -m "Hotfix: ${{ needs.validate-hotfix.outputs.version }}"
          git push origin ${{ needs.validate-hotfix.outputs.version }}

      - name: Merge hotfix to main
        run: |
          git checkout main
          git merge --no-ff hotfix/${{ needs.validate-hotfix.outputs.version }}
          git push origin main

      - name: Merge hotfix to develop
        run: |
          git checkout develop
          git merge --no-ff hotfix/${{ needs.validate-hotfix.outputs.version }}
          git push origin develop

      - name: Delete hotfix branch
        run: |
          git push origin --delete hotfix/${{ needs.validate-hotfix.outputs.version }}

      - name: Notify team - Success
        if: success()
        uses: 8398a7/action-slack@v3
        with:
          status: success
          text: |
            ✅ HOTFIX DEPLOYED SUCCESSFULLY
            Version: ${{ needs.validate-hotfix.outputs.version }}
            Time: $(date)
          webhook_url: ${{ secrets.SLACK_WEBHOOK }}

  # ==================== ROLLBACK ON FAILURE ====================
  rollback-on-failure:
    needs: deploy-production
    if: failure()
    runs-on: ubuntu-latest

    steps:
      - name: Emergency rollback
        run: |
          kubectl rollout undo deployment/nestjs-api -n production
          kubectl rollout undo deployment/nestjs-api-canary -n production

      - name: Restore database
        run: |
          ./scripts/restore-database.sh production hotfix-backup

      - name: Alert team
        uses: 8398a7/action-slack@v3
        with:
          status: failure
          text: |
            🚨 HOTFIX DEPLOYMENT FAILED - ROLLED BACK
            Version: ${{ needs.validate-hotfix.outputs.version }}
            System restored to previous state
          webhook_url: ${{ secrets.SLACK_WEBHOOK }}
```


## Rollback Strategies

### 1. Kubernetes Automatic Rollback

**deployment.yaml with automatic rollback**

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nestjs-api
  namespace: production
spec:
  replicas: 3
  progressDeadlineSeconds: 600  # 10 minutes timeout
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0  # Zero downtime
  minReadySeconds: 30  # Wait 30s before considering pod ready
  
  template:
    metadata:
      labels:
        app: nestjs-api
    spec:
      containers:
      - name: nestjs-api
        image: your-org/car-service-nestjs:latest
        
        # Health checks for automatic rollback
        livenessProbe:
          httpGet:
            path: /health
            port: 3000
          initialDelaySeconds: 60
          periodSeconds: 10
          failureThreshold: 3  # Fail after 3 attempts
          successThreshold: 1
          timeoutSeconds: 5
        
        readinessProbe:
          httpGet:
            path: /ready
            port: 3000
          initialDelaySeconds: 30
          periodSeconds: 5
          failureThreshold: 3
          successThreshold: 1
          timeoutSeconds: 3
        
        # Startup probe for slow starting apps
        startupProbe:
          httpGet:
            path: /health
            port: 3000
          initialDelaySeconds: 0
          periodSeconds: 10
          failureThreshold: 30  # 5 minutes total
```


### 2. Manual Rollback Commands

**Rollback Scripts**

**scripts/rollback-production.sh**

```bash
#!/bin/bash

set -e

NAMESPACE="production"
DEPLOYMENT_NAME="nestjs-api"

echo "🔄 Starting rollback process..."

# Get current revision
CURRENT_REVISION=$(kubectl rollout history deployment/$DEPLOYMENT_NAME -n $NAMESPACE | tail -1 | awk '{print $1}')
echo "Current revision: $CURRENT_REVISION"

# Get previous revision
PREVIOUS_REVISION=$((CURRENT_REVISION - 1))
echo "Rolling back to revision: $PREVIOUS_REVISION"

# Confirm rollback
read -p "Are you sure you want to rollback? (yes/no): " CONFIRM
if [ "$CONFIRM" != "yes" ]; then
    echo "Rollback cancelled"
    exit 0
fi

# Create backup of current state
echo "📦 Creating backup..."
kubectl get deployment $DEPLOYMENT_NAME -n $NAMESPACE -o yaml > "backup-before-rollback-$(date +%s).yaml"

# Perform rollback
echo "⏪ Rolling back..."
kubectl rollout undo deployment/$DEPLOYMENT_NAME -n $NAMESPACE

# Wait for rollback to complete
echo "⏳ Waiting for rollback to complete..."
kubectl rollout status deployment/$DEPLOYMENT_NAME -n $NAMESPACE --timeout=5m

# Verify rollback
echo "✅ Verifying rollback..."
kubectl get pods -n $NAMESPACE -l app=$DEPLOYMENT_NAME

# Check health
echo "🏥 Checking application health..."
sleep 10
HEALTH_STATUS=$(curl -s -o /dev/null -w "%{http_code}" https://car-service.com/health)

if [ "$HEALTH_STATUS" == "200" ]; then
    echo "✅ Rollback successful - Application is healthy"
else
    echo "❌ Rollback completed but health check failed"
    exit 1
fi

# Notify team
curl -X POST $SLACK_WEBHOOK \
  -H 'Content-Type: application/json' \
  -d "{
    \"text\": \"✅ Production rollback completed successfully\",
    \"attachments\": [{
      \"color\": \"good\",
      \"fields\": [
        {\"title\": \"Deployment\", \"value\": \"$DEPLOYMENT_NAME\", \"short\": true},
        {\"title\": \"Revision\", \"value\": \"$PREVIOUS_REVISION\", \"short\": true},
        {\"title\": \"Executed by\", \"value\": \"$USER\", \"short\": true}
      ]
    }]
  }"

echo "🎉 Rollback complete!"
```


### 3. Blue-Green Rollback

**Instant rollback by switching service**

```bash
#!/bin/bash

# Current: Green is live, Blue is previous version

echo "🔄 Performing instant rollback via traffic switch..."

# Switch service to point to blue (old) deployment
kubectl patch service nestjs-api -n production \
  --type='json' \
  -p='[{"op": "replace", "path": "/spec/selector/version", "value": "blue"}]'

echo "✅ Traffic switched to blue environment (previous version)"

# Verify
kubectl get service nestjs-api -n production -o jsonpath='{.spec.selector}'

# Monitor for 2 minutes
echo "📊 Monitoring application..."
sleep 120

# Check error rates
ERROR_RATE=$(./scripts/get-error-rate.sh)
echo "Error rate: $ERROR_RATE"

if (( $(echo "$ERROR_RATE < 0.01" | bc -l) )); then
    echo "✅ Rollback successful - System stable"
    
    # Scale down green (failed) deployment
    kubectl scale deployment/nestjs-api-green --replicas=0 -n production
else
    echo "⚠️  Warning: Error rate still elevated"
fi
```


### 4. Database Rollback Strategy

**scripts/database-rollback.sh**

```bash
#!/bin/bash

set -e

ENVIRONMENT=$1
BACKUP_ID=$2

if [ -z "$ENVIRONMENT" ] || [ -z "$BACKUP_ID" ]; then
    echo "Usage: ./database-rollback.sh <environment> <backup-id>"
    exit 1
fi

echo "⚠️  DATABASE ROLLBACK - This is a critical operation!"
echo "Environment: $ENVIRONMENT"
echo "Backup ID: $BACKUP_ID"

# Confirm
read -p "Type 'ROLLBACK' to confirm: " CONFIRM
if [ "$CONFIRM" != "ROLLBACK" ]; then
    echo "Cancelled"
    exit 0
fi

# Get database credentials
DB_HOST=$(kubectl get secret db-credentials -n $ENVIRONMENT -o jsonpath='{.data.host}' | base64 -d)
DB_NAME=$(kubectl get secret db-credentials -n $ENVIRONMENT -o jsonpath='{.data.database}' | base64 -d)

echo "📦 Creating pre-rollback backup..."
pg_dump -h $DB_HOST -U postgres $DB_NAME > "pre-rollback-$(date +%s).sql"

echo "🔄 Restoring from backup..."
# Download backup from S3
aws s3 cp s3://backups/$ENVIRONMENT/$BACKUP_ID.sql ./restore.sql

# Stop application traffic
kubectl scale deployment/nestjs-api --replicas=0 -n $ENVIRONMENT

# Restore database
psql -h $DB_HOST -U postgres $DB_NAME < restore.sql

# Restart application
kubectl scale deployment/nestjs-api --replicas=3 -n $ENVIRONMENT

# Wait for pods
kubectl rollout status deployment/nestjs-api -n $ENVIRONMENT

echo "✅ Database rollback complete"

# Verify data
echo "🔍 Verifying data integrity..."
./scripts/verify-database.sh $ENVIRONMENT

# Notify
curl -X POST $SLACK_WEBHOOK \
  -H 'Content-Type: application/json' \
  -d "{\"text\": \"⚠️ Database rolled back in $ENVIRONMENT to backup $BACKUP_ID\"}"
```


## Rollback Decision Matrix

| Severity | Symptoms | Rollback Strategy | Approval Required | Estimated Time |
| :-- | :-- | :-- | :-- | :-- |
| **Critical** | Service down, >50% errors | Immediate blue-green switch | Post-facto | < 1 min |
| **High** | >10% errors, performance degraded | Kubernetes rollback | DevOps lead | 5-10 min |
| **Medium** | <10% errors, specific feature broken | Redeploy previous version | Team lead | 15-30 min |
| **Low** | Minor issues, workaround available | Schedule proper fix | Product owner | Next deploy |

## Monitoring and Alerting

**Automated Monitoring Script**

**scripts/monitor-deployment.sh**

```bash
#!/bin/bash

DEPLOYMENT=$1
NAMESPACE=$2
DURATION=${3:-300}  # Default 5 minutes

echo "📊 Monitoring deployment: $DEPLOYMENT in $NAMESPACE for ${DURATION}s"

START_TIME=$(date +%s)
FAILED_CHECKS=0
MAX_FAILURES=3

while [ $(($(date +%s) - START_TIME)) -lt $DURATION ]; do
    # Check pod status
    READY_PODS=$(kubectl get deployment $DEPLOYMENT -n $NAMESPACE -o jsonpath='{.status.readyReplicas}')
    DESIRED_PODS=$(kubectl get deployment $DEPLOYMENT -n $NAMESPACE -o jsonpath='{.spec.replicas}')
    
    # Check error rate from Prometheus
    ERROR_RATE=$(curl -s "http://prometheus:9090/api/v1/query?query=rate(http_requests_total{status=~'5..'}[1m])" | jq -r '.data.result[^0].value[^1]')
    
    # Check response time
    AVG_RESPONSE=$(curl -s "http://prometheus:9090/api/v1/query?query=histogram_quantile(0.95,rate(http_request_duration_seconds_bucket[1m]))" | jq -r '.data.result[^0].value[^1]')
    
    echo "[$(date +%T)] Pods: $READY_PODS/$DESIRED_PODS | Error Rate: $ERROR_RATE | P95 Latency: ${AVG_RESPONSE}s"
    
    # Check conditions
    if [ "$READY_PODS" != "$DESIRED_PODS" ]; then
        echo "⚠️  Warning: Not all pods are ready"
        ((FAILED_CHECKS++))
    elif (( $(echo "$ERROR_RATE > 0.05" | bc -l) )); then
        echo "⚠️  Warning: High error rate"
        ((FAILED_CHECKS++))
    elif (( $(echo "$AVG_RESPONSE > 2.0" | bc -l) )); then
        echo "⚠️  Warning: High latency"
        ((FAILED_CHECKS++))
    else
        echo "✅ All metrics healthy"
        FAILED_CHECKS=0
    fi
    
    # Trigger rollback if too many failures
    if [ $FAILED_CHECKS -ge $MAX_FAILURES ]; then
        echo "🚨 CRITICAL: Too many failed checks - Triggering automatic rollback"
        ./scripts/rollback-production.sh
        exit 1
    fi
    
    sleep 30
done

echo "✅ Monitoring complete - Deployment appears stable"
```


## Post-Incident Checklist

```markdown
## Post-Hotfix/Rollback Checklist

### Immediate Actions (Within 1 hour)
- [ ] Verify system stability
- [ ] Check all metrics returned to normal
- [ ] Confirm no data loss
- [ ] Notify all stakeholders
- [ ] Document the incident

### Short-term Actions (Within 24 hours)
- [ ] Root cause analysis
- [ ] Update runbooks
- [ ] Review monitoring/alerting
- [ ] Plan permanent fix
- [ ] Schedule post-mortem meeting

### Long-term Actions (Within 1 week)
- [ ] Implement permanent fix
- [ ] Add tests to prevent regression
- [ ] Update deployment process
- [ ] Training session if needed
- [ ] Update disaster recovery plan

### Post-Mortem Questions
1. What triggered the issue?
2. How was it detected?
3. What was the impact?
4. How long to resolve?
5. What worked well?
6. What could be improved?
7. Action items to prevent recurrence
```


## Best Practices Summary

### Hotfix Best Practices ✅

- **Create from production** (main branch), not develop
- **Keep changes minimal** - fix only the critical issue
- **Fast-track testing** - focus on affected areas only
- **Deploy to staging first** - even for hotfixes
- **Merge back to develop** - keep branches in sync
- **Document thoroughly** - why, what, how
[^2][^1]


### Rollback Best Practices ✅

- **Automate monitoring** - detect issues early
- **Multiple rollback options** - blue-green, Kubernetes, manual
- **Practice regularly** - test rollback procedures
- **Keep previous versions** - don't delete immediately
- **Database backups** - before every deployment
- **Clear communication** - notify team immediately
[^3][^4]


### Things to Avoid ❌

- ❌ Deploying hotfix without testing
- ❌ Making multiple changes in one hotfix
- ❌ Skipping code review for hotfixes
- ❌ Forgetting to merge back to develop
- ❌ No monitoring after deployment
- ❌ Deleting old deployments too quickly

การมีระบบ hotfix และ rollback ที่ดีช่วยให้ทีมสามารถจัดการปัญหาใน production ได้อย่างรวดเร็วและปลอดภัย ลดผลกระทบต่อผู้ใช้งานให้เหลือน้อยที่สุด[^4][^5][^1]
<span style="display:none">[^10][^6][^7][^8][^9]</span>

<div align="center">⁂</div>

[^1]: https://notes.kodekloud.com/docs/AZ-400/Design-and-Implement-Deployments/Design-a-hotfix-path-plan

[^2]: https://www.growingscrummasters.com/keywords/production-hotfix/

[^3]: https://www.blinkops.com/blog/how-to-rollback-your-kubernetes-deployment

[^4]: https://snyk.io/articles/blue-green-deployment/

[^5]: https://en.wikipedia.org/wiki/Blue–green_deployment

[^6]: https://learn.microsoft.com/en-us/azure/data-factory/continuous-integration-delivery-hotfix-environment

[^7]: https://docs.fintechos.com/Platform/24.4/AdminGuide/Content/Installation/HotFixDeploy.htm

[^8]: https://www.linkedin.com/pulse/how-test-hotfix-advantages-challenges-testing-strategies-testrigor-gyrye

[^9]: https://benediktbergmann.eu/2021/02/28/how-to-handle-hotfixes-in-dataverse/

[^10]: https://github.com/kubernetes/kubernetes/issues/23211

