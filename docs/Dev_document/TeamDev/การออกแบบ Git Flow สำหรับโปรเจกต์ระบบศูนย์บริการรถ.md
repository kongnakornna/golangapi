
## การออกแบบ Git Flow สำหรับโปรเจกต์ระบบศูนย์บริการรถยนต์

สำหรับโปรเจกต์นี้ที่มีทั้ง Backend (NestJS + Spring Boot) และ Frontend (ReactJS) แนะนำให้ใช้ **Modified Git Flow** ที่ผสมผสานความเข้มงวดของ Git Flow กับความคล่องตัวของ GitHub Flow เพื่อรองรับ microservices architecture[^1][^2]

## โครงสร้าง Repository

### Option 1: Monorepo (แนะนำ)

```
car-service-system/
├── backend/
│   ├── nestjs-api/          # NestJS API Service
│   ├── spring-boot-api/     # Spring Boot + Kafka Service
│   └── shared/              # Shared types, utilities
├── frontend/
│   └── react-app/           # ReactJS Application
├── infrastructure/
│   ├── docker/              # Docker configurations
│   ├── kubernetes/          # K8s manifests
│   └── terraform/           # Infrastructure as Code
├── docs/                    # Documentation
└── scripts/                 # Deployment scripts
```

**ข้อดีของ Monorepo:**

- เปลี่ยนแปลงหลายส่วนใน PR เดียว (เช่น เพิ่ม API endpoint และ UI พร้อมกัน)
- Dependency management ง่ายกว่า
- Integration testing สะดวก
- Version control แบบ atomic
[^3][^4]


### Option 2: Multi-Repo (สำหรับทีมใหญ่)

```
Repositories:
- car-service-nestjs-api
- car-service-spring-boot-api  
- car-service-frontend
- car-service-infrastructure
```


## Git Flow Structure

### Branch Strategy

```
main (production)
├── develop (integration)
│   ├── feature/booking-system
│   ├── feature/repair-tracking
│   ├── feature/payment-integration
│   └── feature/notification-service
├── release/v1.0.0
│   └── release/v1.1.0
└── hotfix/fix-payment-bug
```


### Branch Types และ Naming Conventions

| Branch Type | Naming Pattern | Purpose | Merge To |
| :-- | :-- | :-- | :-- |
| **main** | `main` | Production-ready code | - |
| **develop** | `develop` | Integration branch | main (via release) |
| **feature** | `feature/[ticket]-[description]` | New features | develop |
| **bugfix** | `bugfix/[ticket]-[description]` | Bug fixes | develop |
| **release** | `release/v[major].[minor].[patch]` | Release preparation | main + develop |
| **hotfix** | `hotfix/[ticket]-[description]` | Production fixes | main + develop |

[^5][^6]

### Branch Naming Examples

```bash
# Features
feature/CAR-123-booking-api
feature/CAR-124-jwt-authentication
feature/CAR-125-kafka-integration
feature/CAR-126-booking-form-ui

# Bug fixes
bugfix/CAR-201-fix-booking-validation
bugfix/CAR-202-cache-invalidation-issue

# Releases
release/v1.0.0
release/v1.1.0

# Hotfixes
hotfix/CAR-301-fix-payment-timeout
hotfix/CAR-302-security-patch
```


## Workflow Diagram

```
┌──────────────────────────────────────────────────────────────┐
│                         MAIN BRANCH                           │
│                    (Production Ready)                         │
└────┬──────────────────────────────────────────────┬──────────┘
     │                                               │
     │ Hotfix                                   Release Merge
     │                                               │
┌────┴───────────┐                           ┌──────┴──────────┐
│ hotfix/fix-xxx │                           │ release/v1.0.0  │
└────┬───────────┘                           └──────┬──────────┘
     │                                               │
     │ Merge back                                    │ Merge
     │                                               │
┌────┴───────────────────────────────────────────────┴─────────┐
│                      DEVELOP BRANCH                           │
│                  (Integration Branch)                         │
└────┬────────┬────────┬────────┬────────┬────────┬───────────┘
     │        │        │        │        │        │
     │        │        │        │        │        │
┌────┴────┐ ┌┴────┐  ┌┴─────┐ ┌┴──────┐┌┴──────┐ ┌┴──────────┐
│ feature │ │feat.│  │bugfix│ │feature││bugfix │ │ feature   │
│ /booking│ │/auth│  │/cache│ │/kafka ││/valid.│ │ /payment  │
└─────────┘ └─────┘  └──────┘ └───────┘└───────┘ └───────────┘
```


## Complete Workflow

### 1. Feature Development Workflow

```bash
# 1. Create feature branch from develop
git checkout develop
git pull origin develop
git checkout -b feature/CAR-123-booking-api

# 2. Work on feature (commit often)
git add .
git commit -m "feat(booking): add booking API endpoint"
git commit -m "feat(booking): implement validation logic"
git commit -m "test(booking): add unit tests for booking service"

# 3. Keep feature branch updated with develop
git fetch origin
git rebase origin/develop

# 4. Push feature branch
git push origin feature/CAR-123-booking-api

# 5. Create Pull Request (PR)
# - Go to GitHub/GitLab
# - Create PR from feature/CAR-123-booking-api → develop
# - Add reviewers
# - Link JIRA ticket

# 6. After PR approval, merge to develop
# - Squash and merge (clean history)
# - Delete feature branch
```


### 2. Release Workflow

```bash
# 1. Create release branch from develop
git checkout develop
git pull origin develop
git checkout -b release/v1.0.0

# 2. Update version numbers
# - package.json (frontend)
# - pom.xml / build.gradle (Spring Boot)
# - package.json (NestJS)
npm version 1.0.0
git commit -m "chore: bump version to 1.0.0"

# 3. Final testing and bug fixes on release branch
git commit -m "fix(release): minor UI adjustments"
git commit -m "docs: update changelog for v1.0.0"

# 4. Merge to main (production)
git checkout main
git merge --no-ff release/v1.0.0
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin main --tags

# 5. Merge back to develop
git checkout develop
git merge --no-ff release/v1.0.0
git push origin develop

# 6. Delete release branch
git branch -d release/v1.0.0
git push origin --delete release/v1.0.0
```


### 3. Hotfix Workflow

```bash
# 1. Create hotfix branch from main
git checkout main
git pull origin main
git checkout -b hotfix/CAR-301-fix-payment-timeout

# 2. Fix the critical bug
git add .
git commit -m "fix(payment): increase timeout to 30 seconds"
git commit -m "test(payment): add timeout test case"

# 3. Bump patch version
npm version patch  # 1.0.0 → 1.0.1

# 4. Merge to main
git checkout main
git merge --no-ff hotfix/CAR-301-fix-payment-timeout
git tag -a v1.0.1 -m "Hotfix: Fix payment timeout"
git push origin main --tags

# 5. Merge to develop
git checkout develop
git merge --no-ff hotfix/CAR-301-fix-payment-timeout
git push origin develop

# 6. Delete hotfix branch
git branch -d hotfix/CAR-301-fix-payment-timeout
git push origin --delete hotfix/CAR-301-fix-payment-timeout
```


## Commit Message Convention

### Conventional Commits Format

```
<type>(<scope>): <subject>

<body>

<footer>
```


### Commit Types

```bash
feat:     # New feature
fix:      # Bug fix
docs:     # Documentation changes
style:    # Code style changes (formatting, semicolons)
refactor: # Code refactoring
perf:     # Performance improvements
test:     # Adding or updating tests
chore:    # Build process or auxiliary tool changes
ci:       # CI/CD configuration changes
```


### Commit Examples

```bash
# Feature
git commit -m "feat(booking): add booking creation API endpoint"
git commit -m "feat(auth): implement JWT authentication"

# Bug fix
git commit -m "fix(booking): resolve date validation issue"
git commit -m "fix(cache): correct Redis cache invalidation logic"

# Documentation
git commit -m "docs(api): update API documentation for bookings"

# Testing
git commit -m "test(booking): add integration tests for booking service"

# Breaking changes
git commit -m "feat(api)!: change booking API response format

BREAKING CHANGE: Response now includes metadata field"
```


## Pull Request (PR) Guidelines

### PR Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix (non-breaking change)
- [ ] New feature (non-breaking change)
- [ ] Breaking change
- [ ] Documentation update

## Related Issues
Closes #123
Related to #456

## Changes Made
- Added booking API endpoint
- Implemented JWT authentication
- Updated database schema

## Testing
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Manual testing completed

## Checklist
- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] Comments added for complex logic
- [ ] Documentation updated
- [ ] No new warnings generated

## Screenshots (if applicable)
[Add screenshots here]

## Deployment Notes
- Requires database migration
- Environment variables need updating
```


### PR Review Process

1. **Developer** creates PR with detailed description
2. **Automated Checks** run (CI/CD pipeline)
    - Linting
    - Unit tests
    - Integration tests
    - Code coverage
3. **Code Review** by 2+ team members
4. **Changes Requested** (if needed)
5. **Approval** from reviewers
6. **Merge** to target branch
7. **Automatic Deployment** (if configured)

[^2][^5]

## Branch Protection Rules

### Main Branch Protection

```yaml
Branch: main
Rules:
  - Require pull request reviews: 2 approvals
  - Require status checks to pass: true
    - CI/CD pipeline
    - Unit tests
    - Integration tests
    - Code coverage ≥ 80%
  - Require branches to be up to date: true
  - Include administrators: true
  - Restrict who can push: Only via PR
  - Require linear history: true
```


### Develop Branch Protection

```yaml
Branch: develop
Rules:
  - Require pull request reviews: 1 approval
  - Require status checks to pass: true
    - Linting
    - Unit tests
  - Require branches to be up to date: true
  - Allow force pushes: false
```


## CI/CD Integration

### GitHub Actions Workflow Example

```yaml
# .github/workflows/ci-cd.yml
name: CI/CD Pipeline

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main, develop]

jobs:
  # Backend NestJS
  nestjs-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
        with:
          node-version: '18'
      - run: cd backend/nestjs-api && npm ci
      - run: cd backend/nestjs-api && npm run lint
      - run: cd backend/nestjs-api && npm test
      - run: cd backend/nestjs-api && npm run test:e2e

  # Backend Spring Boot
  spring-boot-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-java@v3
        with:
          java-version: '17'
      - run: cd backend/spring-boot-api && ./mvnw test

  # Frontend React
  frontend-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
      - run: cd frontend/react-app && npm ci
      - run: cd frontend/react-app && npm run lint
      - run: cd frontend/react-app && npm test -- --coverage

  # Deploy to staging (on develop)
  deploy-staging:
    needs: [nestjs-test, spring-boot-test, frontend-test]
    if: github.ref == 'refs/heads/develop'
    runs-on: ubuntu-latest
    steps:
      - name: Deploy to Staging
        run: ./scripts/deploy-staging.sh

  # Deploy to production (on main)
  deploy-production:
    needs: [nestjs-test, spring-boot-test, frontend-test]
    if: github.ref == 'refs/heads/main'
    runs-on: ubuntu-latest
    steps:
      - name: Deploy to Production
        run: ./scripts/deploy-production.sh
```


## Monorepo Change Detection

```yaml
# Detect changes and deploy only affected services
detect-changes:
  runs-on: ubuntu-latest
  outputs:
    nestjs: ${{ steps.changes.outputs.nestjs }}
    spring: ${{ steps.changes.outputs.spring }}
    frontend: ${{ steps.changes.outputs.frontend }}
  steps:
    - uses: actions/checkout@v3
    - uses: dorny/paths-filter@v2
      id: changes
      with:
        filters: |
          nestjs:
            - 'backend/nestjs-api/**'
          spring:
            - 'backend/spring-boot-api/**'
          frontend:
            - 'frontend/react-app/**'
```


## Best Practices Summary

### Do's ✅

- **Keep branches short-lived** (< 2 weeks)
- **Commit frequently** with meaningful messages
- **Rebase before merge** to keep history clean
- **Write descriptive PR descriptions**
- **Request code reviews** from 2+ people
- **Run tests locally** before pushing
- **Update documentation** with code changes
- **Use branch protection** rules
- **Tag releases** with semantic versioning
- **Delete merged branches** to keep repo clean


### Don'ts ❌

- ❌ Don't commit directly to main/develop
- ❌ Don't push broken code
- ❌ Don't merge without PR review
- ❌ Don't use generic commit messages ("fix", "update")
- ❌ Don't leave PRs open for too long
- ❌ Don't mix unrelated changes in one commit
- ❌ Don't force push to shared branches
- ❌ Don't merge with failing tests

[^7][^1][^6]

## Team Workflow Summary

```
Developer → Create Feature Branch → Work on Feature
    ↓
Commit Changes → Push to Remote → Create PR
    ↓
Code Review → Automated Tests → Address Feedback
    ↓
Approval → Merge to Develop → Delete Feature Branch
    ↓
Weekly/Sprint Release → Create Release Branch → Test
    ↓
Merge to Main → Tag Version → Deploy to Production
    ↓
Merge Release back to Develop → Continue Development
```

Git Flow นี้ให้ความสมดุลระหว่าง structure และ flexibility เหมาะสำหรับทีมขนาดกลางถึงใหญ่ที่ต้องการ quality control แต่ยังคงความคล่องตัวในการพัฒนา[^2][^1]
<span style="display:none">[^10][^8][^9]</span>

<div align="center">⁂</div>

[^1]: https://graphite.com/guides/advanced-git-branching-strategies

[^2]: https://www.harness.io/blog/github-flow-vs-git-flow-whats-the-difference

[^3]: https://www.reddit.com/r/programming/comments/uauari/when_it_comes_to_microservices_do_you_put_each/

[^4]: https://dev.to/koseimori/implementing-continuous-delivery-for-github-monorepos-and-microservices-with-github-actions-50i8

[^5]: https://dev.to/karmpatel/git-branching-strategies-a-comprehensive-guide-24kh

[^6]: https://www.datacamp.com/tutorial/git-branching-strategy-guide

[^7]: https://sandboxtechnology.in/mastering-git-branching-merging-best-practices/

[^8]: https://www.atlassian.com/git/tutorials/comparing-workflows/gitflow-workflow

[^9]: https://docs.aws.amazon.com/prescriptive-guidance/latest/choosing-git-branch-approach/git-branching-strategies.html

[^10]: https://www.geeksforgeeks.org/git/git-flow-vs-github-flow/

