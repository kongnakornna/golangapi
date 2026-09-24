 ## การออกแบบ Git Flow สำหรับโปรเจกต์ระบบศูนย์บริการรถยนต์
## Template

- สำหรับโปรเจกต์นี้ที่มีทั้ง Backend (NestJS + Spring Boot) และ Frontend (ReactJS) แนะนำให้ใช้ **Modified Git Flow** ที่ผสมผสานความเข้มงวดของ Git Flow กับความคล่องตัวของ GitHub Flow เพื่อรองรับ microservices architecture[^1][^2]

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
├── dev (integration)
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
| **dev** | `dev` | Integration branch | main (via release) |
| **feature** | `feature/[ticket]-[description]` | New features | dev |
| **bugfix** | `bugfix/[ticket]-[description]` | Bug fixes | dev |
| **release** | `release/v[major].[minor].[patch]` | Release preparation | main + dev |
| **hotfix** | `hotfix/[ticket]-[description]` | Production fixes | main + dev |

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
│                         MAIN BRANCH                          │
│                    (Production Ready)                        │
└────┬──────────────────────────────────────────────┬──────────┘
     │                                              │
     │ Hotfix                                   Release Merge
     │                                              │
┌────┴───────────┐                           ┌──────┴──────────┐
│ hotfix/fix-xxx │                           │ release/v1.0.0  │
└────┬───────────┘                           └──────┬──────────┘
     │                                              │
     │ Merge back                                   │ Merge
     │                                              │
┌────┴──────────────────────────────────────────────┴─────────┐
│                      DEVELOP BRANCH    dev                  │
│                  (Integration Branch)                       │
└────┬─────────┬────────┬────────┬────────┬────────┬──────────┘
     │         │        │        │        │        │
     │         │        │        │        │        │
┌────┴────┐ ┌──┴──┐  ┌──┴───┐ ┌──┴────┐┌──┴────┐ ┌─┴─────────┐
│ feature │ │feat.│  │bugfix│ │feature││bugfix │ │ feature   │
│ /booking│ │/auth│  │/cache│ │/kafka ││/valid.│ │ /payment  │
└─────────┘ └─────┘  └──────┘ └───────┘└───────┘ └───────────┘
```


## Complete Workflow

### 1. Feature Development Workflow

```bash
# 1. Create feature branch from dev
git checkout dev
git pull origin dev
git checkout -b feature/CAR-123-booking-api

# 2. Work on feature (commit often)
git add .
git commit -m "feat(booking): add booking API endpoint"
git commit -m "feat(booking): implement validation logic"
git commit -m "test(booking): add unit tests for booking service"

# 3. Keep feature branch updated with dev
git fetch origin
git rebase origin/dev

# 4. Push feature branch
git push origin feature/CAR-123-booking-api

# 5. Create Pull Request (PR)
# - Go to GitHub/GitLab
# - Create PR from feature/CAR-123-booking-api → dev
# - Add reviewers
# - Link JIRA ticket

# 6. After PR approval, merge to dev
# - Squash and merge (clean history)
# - Delete feature branch
```


### 2. Release Workflow

```bash
# 1. Create release branch from dev
git checkout dev
git pull origin dev
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

# 5. Merge back to dev
git checkout dev
git merge --no-ff release/v1.0.0
git push origin dev

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

# 5. Merge to dev
git checkout dev
git merge --no-ff hotfix/CAR-301-fix-payment-timeout
git push origin dev

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


## Description
Brief description of changes

## Type of Change
- [ ] Bug fix (non-breaking change)
- [x] New feature (non-breaking change)
- [x] Breaking change
- [ ] Documentation update

## Related Issues
Closes #123
Related to #456

## Changes Made
- Added booking API endpoint
- Implemented JWT authentication
- Updated database schema

## Testing
- [x] Unit tests pass
- [ ] Integration tests pass
- [ ] Manual testing completed

## Checklist
- [ ] Code follows project style guidelines
- [x] Self-review completed
- [x] Comments added for complex logic
- [ ] Documentation updated
- [ ] No new warnings generated

## Screenshots (if applicable)
[Add screenshots here]

## Deployment Notes
- Requires database migration
- Environment variables need updating
 


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
Branch: dev
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
    branches: [main, dev]
  pull_request:
    branches: [main, dev]

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

  # Deploy to staging (on dev)
  deploy-staging:
    needs: [nestjs-test, spring-boot-test, frontend-test]
    if: github.ref == 'refs/heads/dev'
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

- ❌ Don't commit directly to main/dev
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

- Git Flow นี้ให้ความสมดุลระหว่าง structure และ flexibility เหมาะสำหรับทีมขนาดกลางถึงใหญ่ที่ต้องการ quality control แต่ยังคงความคล่องตัวในการพัฒนา[^2][^1]


# รูปแบบการตั้งชื่อ branch สำหรับฟีเจอร์ขนาดเล็ก

สำหรับการตั้งชื่อ Branch ที่เป็น **"ฟีเจอร์ขนาดเล็ก" (Small Feature)** หรือการปรับแก้เล็กๆ น้อยๆ ที่อาจจะไม่ถึงขั้นเรียกว่าเป็น "Full Feature" ใหญ่ๆ นั้น หลักการสำคัญคือ **"อย่าสร้าง Prefix ใหม่เยอะเกินความจำเป็น"** เพื่อไม่ให้ทีมสับสนครับ

นี่คือ 3 แนวทางที่แนะนำ เรียงตามความนิยมและความเหมาะสม:

### 1. ใช้ `feature/` เหมือนเดิม (แนะนำสูงสุด ⭐️)

แม้จะเป็นงานเล็กๆ ก็ควรนับเป็น Feature เพื่อความสม่ำเสมอ (Consistency) ในการตั้งชื่อและการค้นหา

* **หลักการ:** ใช้ Pattern เดิมแต่เน้นคำอธิบายที่ **กระชับ** และเจาะจง
* **Pattern:** `feature/<TICKET>-<specific-action>`
* **ตัวอย่าง:**
    * `feature/CART-201-add-delete-btn` (แค่เพิ่มปุ่มลบปุ่มเดียว)
    * `feature/USER-305-change-font-size` (เปลี่ยนขนาดฟอนต์)
    * `feature/AUTH-102-hide-password` (ซ่อนรหัสผ่าน)


### 2. ใช้ `tweak/` (สำหรับงานจุกจิก/ปรับแต่ง)

ถ้าทีมรู้สึกว่าคำว่า `feature` ดูยิ่งใหญ่ไปสำหรับงานแก้สีปุ่ม หรือขยับ Layout นิดหน่อย การใช้ `tweak` จะสื่อความหมายได้ดีกว่าว่า "เป็นการปรับแต่งเล็กน้อย"

* **ความหมาย:** การบิด/ดัดแปลง/ปรับแต่ง (ไม่ใช่การสร้างใหม่)

```
*   **Pattern:** `tweak/<TICKET>-<description>`
```

* **ตัวอย่าง:**
    * `tweak/UI-501-adjust-padding` (ขยับช่องว่าง)
    * `tweak/UX-112-wording-change` (แก้คำผิด/แก้ข้อความ)
    * `tweak/CSS-99-dark-mode-color` (ปรับสี Dark mode นิดหน่อย)


### 3. ใช้ `ui/` หรือ `ux/` (เน้นงานหน้าบ้านโดยเฉพาะ)

สำหรับทีมที่มี Frontend หรือ Designer แยกชัดเจน อาจใช้ Prefix นี้เพื่อบอกว่าเป็นงานที่ไม่กระทบ Logic หรือ Database เลย

* **ความหมาย:** งานที่เกี่ยวกับหน้าตาล้วนๆ

```
*   **Pattern:** `ui/<TICKET>-<description>`
```

* **ตัวอย่าง:**
    * `ui/HOME-404-hero-banner` (เปลี่ยนรูป Banner)
    * `ui/MENU-202-icon-update` (เปลี่ยนไอคอน)

***

### 💡 เคล็ดลับ: การตั้งชื่อเมื่อเป็น "ส่วนย่อย" ของฟีเจอร์ใหญ่

บางครั้งฟีเจอร์เล็กๆ นั้นเป็นส่วนหนึ่งของโปรเจกต์ใหญ่ (เช่น ทำระบบ Login แต่แยก Branch มาทำแค่ปุ่ม "Lium Password")
แนะนำให้ใช้ **"Slash Grouping"** (การใช้ / ซ้อนกัน) เพื่อจัดกลุ่มใน Git GUI (เช่น SourceTree, GitKraken จะเห็นเป็นโฟลเดอร์)

* **Pattern:** `feature/<MAJOR-FEATURE>/<SUB-TASK>`
* **ตัวอย่าง:**
    * `feature/login/forget-password-link`
    * `feature/login/google-button`
    * `feature/dashboard/chart-widget`
    * `feature/dashboard/export-csv-btn`

**ข้อดี:** เวลาดูในโปรแกรม Git จะเห็น Branch ถูกจัดกลุ่มเป็นระเบียบ ไม่สะเปะสะปะครับ

### สรุปคำแนะนำ

* ถ้าอยากให้เป็นระเบียบที่สุด -> ใช้ **`feature/`** ทั้งหมด (Work สุดสำหรับ Automation)
* ถ้าอยากแยกความรู้สึกของงาน -> ใช้ **`tweak/`** สำหรับงานแก้เล็กๆ น้อยๆ
* ถ้าเป็นงานส่วนย่อยของงานใหญ่ -> ใช้ **`feature/main/sub`**


# การออกแบบ Git Flow

## **ภาพรวมของ Git Flow**
Git Flow เป็นรูปแบบการทำงานกับ Git ที่กำหนดโครงสร้างสาขา (branch) และขั้นตอนการทำงานอย่างชัดเจน เหมาะสำหรับโครงการที่มีการ release เป็นช่วง ๆ

## **โครงสร้าง Branch หลัก**

### **1. Branch ถาวร (Permanent Branches)**
- **`main`/`master`** - เก็บโค้ดที่พร้อม production เสมอ
- **`dev`** - เก็บโค้ดสำหรับการพัฒนาอย่างต่อเนื่อง

### **2. Branch ชั่วคราว (Temporary Branches)**
- **Feature branches** - สำหรับฟีเจอร์ใหม่
- **Release branches** - สำหรับเตรียม release
- **Hotfix branches** - สำหรับแก้ไขด่วนใน production

## **Workflow แบบละเอียด**

### **ขั้นตอนการพัฒนา Feature ใหม่**
```bash
# 1. สร้าง feature branch จาก dev
git checkout dev
git pull origin dev
git checkout -b feature/ชื่อ-feature

# 2. พัฒนาและ commit
git add .
git commit -m "เพิ่มฟีเจอร์..."

# 3. เสร็จแล้ว merge กลับไป dev
git checkout dev
git pull origin dev
git merge --no-ff feature/ชื่อ-feature
git push origin dev

# 4. ลบ feature branch
git branch -d feature/ชื่อ-feature
```

### **ขั้นตอนการเตรียม Release**
```bash
# 1. สร้าง release branch จาก dev
git checkout dev
git checkout -b release/1.2.0

# 2. ปรับปรุงสำหรับ release (version numbers, docs)
# 3. ทดสอบและแก้ไข bugs ใน branch นี้

# 4. เมื่อพร้อม merge ไป main และ dev
git checkout main
git merge --no-ff release/1.2.0
git tag -a v1.2.0 -m "Release version 1.2.0"

git checkout dev
git merge --no-ff release/1.2.0

# 5. ลบ release branch
git branch -d release/1.2.0
```

### **ขั้นตอน Hotfix**
```bash
# 1. สร้าง hotfix branch จาก main
git checkout main
git checkout -b hotfix/เร่งด่วน

# 2. แก้ไขปัญหา
# 3. Merge กลับไป main และ dev
git checkout main
git merge --no-ff hotfix/เร่งด่วน
git tag -a v1.2.1 -m "Hotfix for critical issue"

git checkout dev
git merge --no-ff hotfix/เร่งด่วน

# 4. ลบ hotfix branch
git branch -d hotfix/เร่งด่วน
```

## **กฎสำคัญ**
1. **ห้าม push ตรงไป main/dev** ต้องผ่าน merge request/pull request
2. **ใช้ --no-ff** เพื่อคงประวัติการ merge
3. **ตั้งชื่อ branch ตาม convention**
4. **Review code** ก่อน merge เสมอ

## **Conventions การตั้งชื่อ**
- Feature: `feature/login-system`
- Release: `release/1.2.0`
- Hotfix: `hotfix/login-error`
- Bugfix: `bugfix/issue-123`

## **เครื่องมือช่วยเหลือ**
- **Git Flow CLI** - อัตโนมัติ workflow
- **GitKraken, Sourcetree** - GUI ที่ support Git Flow
- **GitLab/GitHub** - ตั้งค่า merge request rules

## **ข้อดีและข้อเสีย**

### **ข้อดี**
- โครงสร้างชัดเจน
- เหมาะกับทีมขนาดใหญ่
- การ release ทำได้เป็นระบบ
- แยก環境การพัฒนาชัดเจน

### **ข้อเสีย**
- ซับซ้อนสำหรับโปรเจคเล็ก
- เรียนรู้ curve สูง
- อาจ overkill สำหรับบางโปรเจค

## **ทางเลือกอื่น**
- **GitHub Flow** - ง่ายกว่า เหมาะกับ continuous deployment
- **GitLab Flow** - ผสมระหว่าง Git Flow และ GitHub Flow
- **Trunk-based Development** - เหมาะกับทีมที่ mature ด้าน CI/CD

## **ตัวอย่าง Visualization**
```
main    *---*--------*-----------* (tags)
         \         / \         /
dev   *---*---*---*---*---*
           \ /       \ /
feature     *         *
```

เลือกใช้ Git Flow เมื่อต้องการควบคุมกระบวนการ release อย่างเป็นระบบ หรือทำงานในทีมขนาดใหญ่ที่มีการพัฒนาควบคู่กันหลายฟีเจอร์พร้อมกัน


# TASK LIST และ CHECKLIST Templates สำหรับระบบศูนย์บริการรถยนต์

## 1. TASK LIST Template

### English Version

#### Feature Development Task List

## TASK: [Feature Name]
**Task ID:** CAR-[Number]
**Priority:** High/Medium/Low
**Sprint:** Sprint [Number]
**Estimated Effort:** [X] Story Points / [Y] Days
**Assigned To:** [Developer Name]
**Due Date:** YYYY-MM-DD

### 📋 Description
[Brief description of the feature or task]

### 🎯 Acceptance Criteria
- [x] Criteria 1 - [Specific condition]
- [x] Criteria 2 - [Specific condition]
- [ ] Criteria 3 - [Specific condition]
- [ ] Criteria 4 - [Specific condition]

### 🔧 Technical Requirements
**Backend (NestJS/Spring Boot):**
- [x] API endpoints design
- [x] Database schema changes
- [x] Business logic implementation
- [ ] Kafka integration (if applicable)
- [ ] Authentication/Authorization

**Frontend (ReactJS):**
- [x] UI component design
- [x] API integration
- [x] State management
- [ ] Form validation
- [ ] Responsive design

**Infrastructure:**
- [x] Docker configuration updates
- [x] Environment variables
- [ ] Kubernetes manifests
- [ ] Monitoring setup

### 📝 Development Tasks
**Phase 1: Setup**
- [x] Create feature branch: `feature/CAR-[ID]-short-description`
- [x] Update dependencies (if needed)
- [ ] Set up devment environment

**Phase 2: Implementation**
- [x] Core functionality implementation
- [x] Unit tests (coverage > 80%)
- [ ] Integration tests
- [ ] Error handling
- [ ] Logging implementation

**Phase 3: Code Quality**
- [x] Code review preparation
- [x] Linting and formatting
- [ ] Performance optimization
- [ ] Security review

**Phase 4: Documentation**
- [x] Update API documentation (Swagger/OpenAPI)
- [x] Update README files
- [ ] Add inline comments for complex logic
- [ ] Update CHANGELOG

### 🧪 Testing Tasks
**Automated Testing:**
- [x] Unit tests written and passing
- [x] Integration tests written and passing
- [ ] E2E tests (if applicable)
- [ ] Load testing (if applicable)

**Manual Testing:**
- [x] Feature testing in devment environment
- [x] Cross-browser testing (for frontend)
- [ ] Mobile responsiveness testing
- [ ] API endpoint testing (Postman/curl)

**Integration Testing:**
- [x] Test with dependent services
- [x] Test Kafka message flow
- [ ] Test database migrations
- [ ] Test with authentication

### 🚀 Deployment Tasks
**Development Environment:**
- [x] Build and deploy to dev environment
- [ ] Verify deployment success
- [x] Test in integrated environment

**Integration:**
- [ ] Create Pull Request to dev branch
- [x] Address review comments
- [x] Merge to dev branch
- [ ] Deploy to staging (if configured)

### 🔗 Dependencies
**Blocks:**
- [Dependency 1 - Task ID]
- [Dependency 2 - Task ID]

**Blocked By:**
- [Blocking Task 1 - Task ID]
- [Blocking Task 2 - Task ID]

**Related Tasks:**
- [Related Task 1 - Task ID]
- [Related Task 2 - Task ID]

### 📊 Status Tracking
**Progress:** [0-100]%
**Last Updated:** YYYY-MM-DD
**Current Status:** 
- [ ] Not Started
- [ ] In Progress
- [ ] Code Review
- [ ] Testing
- [ ] Ready for Deployment
- [ ] Completed

### 💬 Notes & Comments
[Additional notes, assumptions, or important information]
---------------------------------------------------------------------


#### Release Preparation Task List

## RELEASE PREPARATION: v[major].[minor].[patch]
**Release Manager:** [Name]
**Planned Release Date:** YYYY-MM-DD
**Change Window:** HH:MM - HH:MM

### 📋 Release Scope
**Features Included:**
- [Feature 1] - CAR-[ID]
- [Feature 2] - CAR-[ID]
- [Feature 3] - CAR-[ID]

**Bug Fixes Included:**
- [Bug Fix 1] - CAR-[ID]
- [Bug Fix 2] - CAR-[ID]

**Breaking Changes (if any):**
- [Breaking Change 1]
- [Breaking Change 2]

### 🔍 Pre-release Checklist
**Code Quality:**
- [ ] All features merged to dev branch
- [ ] Code review completed for all changes
- [ ] No open high-priority bugs
- [ ] Security scanning completed and passed

**Testing:**
- [ ] Integration tests passing (100%)
- [ ] E2E tests passing (if applicable)
- [ ] Performance testing completed
- [ ] Load testing completed
- [ ] UAT completed and signed-off

**Documentation:**
- [ ] API documentation updated
- [ ] User documentation updated
- [ ] CHANGELOG updated
- [ ] Release notes prepared

**Infrastructure:**
- [ ] Database migrations ready
- [ ] Environment configurations updated
- [ ] Monitoring alerts configured
- [ ] Backup procedures verified

### 📝 Release Process Tasks
**Step 1: Preparation (Day -2)**
- [ ] Create release branch: `release/v[major].[minor].[patch]`
- [ ] Update version numbers in all services
- [ ] Final integration testing on release branch
- [ ] Prepare deployment scripts

**Step 2: Staging Deployment (Day -1)**
- [ ] Deploy to staging environment
- [ ] Run smoke tests on staging
- [ ] Conduct final UAT on staging
- [ ] Get stakeholder sign-off

**Step 3: Production Deployment (Release Day)**
- [ ] Merge release branch to main
- [ ] Create git tag: `v[major].[minor].[patch]`
- [ ] Deploy to production environment
- [ ] Verify production deployment
- [ ] Run post-deployment smoke tests

**Step 4: Post-release (Day +1)**
- [ ] Merge release branch back to dev
- [ ] Delete release branch
- [ ] Update documentation with release notes
- [ ] Conduct post-release review meeting

### 🚨 Rollback Plan
**Triggers for Rollback:**
- [ ] Critical errors affecting core functionality
- [ ] Security vulnerabilities discovered
- [ ] Performance degradation > 50%
- [ ] Data corruption detected

**Rollback Procedure:**
1. [Step 1 of rollback]
2. [Step 2 of rollback]
3. [Step 3 of rollback]

### 📞 Communication Plan
**Pre-release Communication:**
- [ ] Notify stakeholders 3 days before release
- [ ] Schedule release meeting
- [ ] Prepare maintenance page (if needed)

**During Release Communication:**
- [ ] Start of deployment notification
- [ ] Status updates every 30 minutes
- [ ] Completion notification

**Post-release Communication:**
- [ ] Success/failure notification
- [ ] Known issues communication
- [ ] User notification (if UI changes)

### 📋 Sign-off Sheet
**Product Owner Approval:** _________________ Date: _________
**Tech Lead Approval:** ______________________ Date: _________
**DevOps Approval:** ________________________ Date: _________
**Release Manager Sign-off:** ________________ Date: _________

**Release Status:** 
- [ ] Successfully Deployed
- [ ] Partially Deployed (with notes)
- [ ] Rolled Back
- [ ] Postponed

**Notes:** ____________________________________________________
---------------------------------------------------------------------


### ภาษาไทย

#### รายการงานพัฒนาคุณสมบัติใหม่

## งาน: [ชื่อคุณสมบัติ]
**รหัสงาน:** CAR-[หมายเลข]
**ลำดับความสำคัญ:** สูง/ปานกลาง/ต่ำ
**สปรินท์:** สปรินท์ [หมายเลข]
**ความพยายามที่ประมาณการ:** [X] สตอรี่พอยต์ / [Y] วัน
**ผู้รับผิดชอบ:** [ชื่อผู้พัฒนา]
**กำหนดส่ง:** YYYY-MM-DD

### 📋 คำอธิบาย
[คำอธิบายสั้นๆ ของคุณสมบัติหรืองาน]

### 🎯 เกณฑ์การยอมรับ
- [x] เกณฑ์ที่ 1 - [เงื่อนไขเฉพาะ]
- [x] เกณฑ์ที่ 2 - [เงื่อนไขเฉพาะ]
- [ ] เกณฑ์ที่ 3 - [เงื่อนไขเฉพาะ]
- [ ] เกณฑ์ที่ 4 - [เงื่อนไขเฉพาะ]

### 🔧 ข้อกำหนดทางเทคนิค
**แบ็กเอนด์ (NestJS/Spring Boot):**
- [x] การออกแบบ API endpoints
- [x] การเปลี่ยนแปลงโครงสร้างฐานข้อมูล
- [ ] การพัฒนาตรรกะธุรกิจ
- [ ] การบูรณาการ Kafka (ถ้ามี)
- [ ] การยืนยันตัวตน/การอนุญาต

**ฟรอนต์เอนด์ (ReactJS):**
- [ ] การออกแบบส่วนประกอบ UI
- [x] การบูรณาการ API
- [x] การจัดการสถานะ
- [x] การตรวจสอบความถูกต้องของฟอร์ม
- [ ] การออกแบบตอบสนอง

**โครงสร้างพื้นฐาน:**
- [x] การอัปเดตการตั้งค่า Docker
- [x] ตัวแปรสภาพแวดล้อม
- [x] Kubernetes manifests
- [ ] การตั้งค่าการเฝ้าติดตาม

### 📝 งานการพัฒนา
**ขั้นตอนที่ 1: การตั้งค่า**
- [ ] สร้างสาขา feature: `feature/CAR-[ID]-คำอธิบายสั้น`
- [ ] อัปเดตการพึ่งพา (ถ้าจำเป็น)
- [ ] ตั้งค่าสภาพแวดล้อมการพัฒนา

**ขั้นตอนที่ 2: การดำเนินการ**
- [ ] การพัฒนาหน้าที่หลัก
- [x] ทดสอบหน่วย (ความครอบคลุม > 80%)
- [x] ทดสอบการรวมระบบ
- [ ] การจัดการข้อผิดพลาด
- [ ] การบันทึกเหตุการณ์

**ขั้นตอนที่ 3: คุณภาพโค้ด**
- [ ] การเตรียมพร้อมสำหรับการตรวจสอบโค้ด
- [x] การตรวจสอบรูปแบบโค้ดและการจัดรูปแบบ
- [x] การปรับปรุงประสิทธิภาพ
- [ ] การตรวจสอบความปลอดภัย

**ขั้นตอนที่ 4: เอกสาร**
- [x] อัปเดตเอกสาร API (Swagger/OpenAPI)
- [x] อัปเดตไฟล์ README
- [ ] เพิ่มคอมเมนต์ในโค้ดสำหรับตรรกะที่ซับซ้อน
- [ ] อัปเดต CHANGELOG

### 🧪 งานการทดสอบ
**การทดสอบอัตโนมัติ:**
- [ ] ทดสอบหน่วยเขียนและผ่านแล้ว
- [x] ทดสอบการรวมระบบเขียนและผ่านแล้ว
- [x] ทดสอบ E2E (ถ้ามี)
- [ ] ทดสอบโหลด (ถ้ามี)

**การทดสอบด้วยตนเอง:**
- [ ] ทดสอบคุณสมบัติในสภาพแวดล้อมการพัฒนา
- [ ] ทดสอบข้ามเบราว์เซอร์ (สำหรับฟรอนต์เอนด์)
- [ ] ทดสอบการตอบสนองบนมือถือ
- [ ] ทดสอบ API endpoint (Postman/curl)

**การทดสอบการรวมระบบ:**
- [x] ทดสอบกับบริการที่พึ่งพา
- [x] ทดสอบการไหลของข้อความ Kafka
- [x] ทดสอบการย้ายฐานข้อมูล
- [ ] ทดสอบกับการยืนยันตัวตน

### 🚀 งานการดีพลอย
**สภาพแวดล้อมการพัฒนา:**
- [x] สร้างและดีพลอยไปยังสภาพแวดล้อม dev
- [ ] ยืนยันความสำเร็จการดีพลอย
- [ ] ทดสอบในสภาพแวดล้อมที่รวมระบบ

**การรวมระบบ:**
- [x] สร้าง Pull Request ไปยังสาขา dev
- [ ] แก้ไขความคิดเห็นจากการตรวจสอบ
- [ ] รวมโค้ดกับสาขา dev
- [ ] ดีพลอยไปยัง staging (ถ้าตั้งค่าไว้)

### 🔗 การพึ่งพา
**บล็อกงานอื่น:**
- [การพึ่งพาที่ 1 - รหัสงาน]
- [การพึ่งพาที่ 2 - รหัสงาน]

**ถูกบล็อกโดย:**
- [งานที่บล็อก 1 - รหัสงาน]
- [งานที่บล็อก 2 - รหัสงาน]

**งานที่เกี่ยวข้อง:**
- [งานที่เกี่ยวข้อง 1 - รหัสงาน]
- [งานที่เกี่ยวข้อง 2 - รหัสงาน]

### 📊 การติดตามสถานะ
**ความคืบหน้า:** [0-100]%
**อัปเดตล่าสุด:** YYYY-MM-DD
**สถานะปัจจุบัน:**
- [x] ยังไม่เริ่ม
- [ ] กำลังดำเนินการ
- [x] กำลังตรวจสอบโค้ด
- [ ] กำลังทดสอบ
- [ ] พร้อมสำหรับการดีพลอย
- [ ] เสร็จสิ้น

### 💬 หมายเหตุและความคิดเห็น
[หมายเหตุเพิ่มเติม สมมติฐาน หรือข้อมูลสำคัญ]
---------------------------------------------------------------------


#### รายการงานเตรียมการปล่อยเวอร์ชัน

## การเตรียมการปล่อยเวอร์ชัน: v[หลัก].[รอง].[แก้ไข]
**ผู้จัดการการปล่อย:** [ชื่อ]
**วันที่ปล่อยที่วางแผน:** YYYY-MM-DD
**ช่วงเวลาการเปลี่ยนแปลง:** HH:MM - HH:MM

### 📋 ขอบเขตการปล่อย
**คุณสมบัติที่รวมอยู่:**
- [คุณสมบัติ 1] - CAR-[ID]
- [คุณสมบัติ 2] - CAR-[ID]
- [คุณสมบัติ 3] - CAR-[ID]

**การแก้ไขข้อบกพร่องที่รวมอยู่:**
- [การแก้ไขข้อบกพร่อง 1] - CAR-[ID]
- [การแก้ไขข้อบกพร่อง 2] - CAR-[ID]

**การเปลี่ยนแปลงที่ทำลายความเข้ากันได้ (ถ้ามี):**
- [การเปลี่ยนแปลงที่ทำลายความเข้ากันได้ 1]
- [การเปลี่ยนแปลงที่ทำลายความเข้ากันได้ 2]

### 🔍 รายการตรวจสอบก่อนการปล่อย
**คุณภาพโค้ด:**
- [x] คุณสมบัติทั้งหมดรวมกับสาขา dev แล้ว
- [ ] การตรวจสอบโค้ดเสร็จสิ้นสำหรับการเปลี่ยนแปลงทั้งหมด
- [ ] ไม่มีข้อบกพร่องระดับสูงที่ยังไม่แก้ไข
- [ ] การตรวจสอบความปลอดภัยเสร็จสิ้นและผ่าน

**การทดสอบ:**
- [ ] ทดสอบการรวมระบบผ่าน (100%)
- [x] ทดสอบ E2E ผ่าน (ถ้ามี)
- [ ] ทดสอบประสิทธิภาพเสร็จสิ้น
- [ ] ทดสอบโหลดเสร็จสิ้น
- [ ] UAT เสร็จสิ้นและได้รับการลงนาม

**เอกสาร:**
- [ ] อัปเดตเอกสาร API
- [x] อัปเดตเอกสารผู้ใช้
- [x] อัปเดต CHANGELOG
- [ ] เตรียมบันทึกการปล่อยเวอร์ชัน

**โครงสร้างพื้นฐาน:**
- [ ] การย้ายฐานข้อมูลพร้อม
- [ ] อัปเดตการกำหนดค่าสภาพแวดล้อม
- [ ] กำหนดค่าการแจ้งเตือนการเฝ้าติดตาม
- [ ] ตรวจสอบขั้นตอนการสำรองข้อมูล

### 📝 งานกระบวนการปล่อย
**ขั้นตอนที่ 1: การเตรียมการ (วัน -2)**
- [ ] สร้างสาขา release: `release/v[หลัก].[รอง].[แก้ไข]`
- [ ] อัปเดตหมายเลขเวอร์ชันในบริการทั้งหมด
- [ ] ทดสอบการรวมระบบสุดท้ายบนสาขา release
- [ ] เตรียมสคริปต์การดีพลอย

**ขั้นตอนที่ 2: การดีพลอย staging (วัน -1)**
- [ ] ดีพลอยไปยังสภาพแวดล้อม staging
- [ ] ทดสอบพื้นฐานบน staging
- [ ] ดำเนินการ UAT สุดท้ายบน staging
- [ ] ได้รับการลงนามจากผู้มีส่วนได้ส่วนเสีย

**ขั้นตอนที่ 3: การดีพลอย production (วันปล่อย)**
- [ ] รวมสาขา release กับ main
- [ ] สร้างแท็ก git: `v[หลัก].[รอง].[แก้ไข]`
- [ ] ดีพลอยไปยังสภาพแวดล้อม production
- [ ] ยืนยันการดีพลอย production
- [ ] ทดสอบพื้นฐานหลังการดีพลอย

**ขั้นตอนที่ 4: หลังการปล่อย (วัน +1)**
- [ ] รวมสาขา release กลับไปยัง dev
- [ ] ลบสาขา release
- [ ] อัปเดตเอกสารด้วยบันทึกการปล่อยเวอร์ชัน
- [ ] ดำเนินการประชุมทบทวนหลังการปล่อย

### 🚨 แผนการย้อนกลับ
**ทริกเกอร์สำหรับการย้อนกลับ:**
- [ ] ข้อผิดพลาดร้ายแรงที่กระทบต่อหน้าที่หลัก
- [ ] พบช่องโหว่ความปลอดภัย
- [ ] ประสิทธิภาพลดลง > 50%
- [ ] ตรวจพบข้อมูลเสียหาย

**ขั้นตอนการย้อนกลับ:**
1. [ขั้นตอนที่ 1 ของการย้อนกลับ]
2. [ขั้นตอนที่ 2 ของการย้อนกลับ]
3. [ขั้นตอนที่ 3 ของการย้อนกลับ]

### 📞 แผนการสื่อสาร
**การสื่อสารก่อนการปล่อย:**
- [ ] แจ้งผู้มีส่วนได้ส่วนเสีย 3 วันก่อนการปล่อย
- [ ] จัดการประชุมการปล่อย
- [ ] เตรียมหน้าบำรุงรักษา (ถ้าจำเป็น)

**การสื่อสารระหว่างการปล่อย:**
- [ ] การแจ้งเตือนเริ่มต้นการดีพลอย
- [ ] อัปเดตสถานะทุก 30 นาที
- [ ] การแจ้งเตือนความเสร็จสิ้น

**การสื่อสารหลังการปล่อย:**
- [ ] การแจ้งเตือนความสำเร็จ/ล้มเหลว
- [ ] การสื่อสารปัญหาที่ทราบ
- [ ] การแจ้งเตือนผู้ใช้ (ถ้ามีการเปลี่ยนแปลง UI)

### 📋 แบบฟอร์มการลงนาม
**การอนุมัติโดย Product Owner:** _____________ วันที่: _________
**การอนุมัติโดย Tech Lead:** __________________ วันที่: _________
**การอนุมัติโดย DevOps:** ____________________ วันที่: _________
**การลงนามโดยผู้จัดการการปล่อย:** ___________ วันที่: _________

**สถานะการปล่อย:**
- [x] ดีพลอยสำเร็จ
- [x] ดีพลอยบางส่วน (พร้อมหมายเหตุ)
- [ ] ย้อนกลับแล้ว
- [ ] เลื่อนออกไป

**หมายเหตุ:** ____________________________________________________
---------------------------------------------------------------------
 
## 2. CHECKLIST Template

### English Version

#### Code Review Checklist
 
## CODE REVIEW CHECKLIST
**Reviewer:** [Name]
**Developer:** [Name]
**PR Link:** [GitHub/GitLab Link]
**Task ID:** CAR-[Number]
**Review Date:** YYYY-MM-DD

### 🔍 Code Quality
**Readability & Maintainability:**
- [ ] Code is clean and well-organized
- [ ] Meaningful variable and function names
- [ ] Consistent coding style throughout
- [ ] No dead or commented-out code
- [ ] Proper separation of concerns

**Error Handling:**
- [ ] All errors are handled appropriately
- [ ] Meaningful error messages
- [ ] Graceful degradation where needed
- [ ] Proper logging of errors

**Performance:**
- [ ] No obvious performance bottlenecks
- [ ] Efficient database queries
- [ ] Proper use of caching (if applicable)
- [ ] Memory usage optimized

### 🧪 Testing
**Test Coverage:**
- [ ] Unit tests cover critical paths
- [ ] Integration tests included
- [ ] Test coverage meets minimum requirement (80%)
- [ ] Edge cases are tested

**Test Quality:**
- [ ] Tests are meaningful and not trivial
- [ ] Tests use proper setup/teardown
- [ ] Mocking is used appropriately
- [ ] Tests are independent of each other

### 🔒 Security
**Input Validation:**
- [ ] All user inputs are validated
- [ ] SQL injection prevention
- [ ] XSS prevention measures
- [ ] File upload validation (if applicable)

**Authentication & Authorization:**
- [ ] Proper authentication checks
- [ ] Role-based access control
- [ ] Session management
- [ ] Token validation

**Sensitive Data:**
- [ ] No hard-coded secrets
- [ ] Sensitive data encrypted
- [ ] Proper password handling
- [ ] API keys secured

### 📝 Documentation
**Code Documentation:**
- [ ] Complex logic has comments
- [ ] Public APIs documented
- [ ] README updated if needed
- [ ] Architecture decisions documented

**API Documentation:**
- [ ] Swagger/OpenAPI updated
- [ ] Request/response examples
- [ ] Error responses documented
- [ ] Authentication requirements documented

### 🛠 Technical Implementation
**Backend (NestJS/Spring Boot):**
- [ ] Proper use of framework patterns
- [ ] Database layer abstraction
- [ ] Business logic separation
- [ ] API response formats consistent

**Frontend (ReactJS):**
- [ ] Component structure follows guidelines
- [ ] State management appropriate
- [ ] Proper use of hooks
- [ ] Responsive design implemented

**Database:**
- [ ] Proper indexes created
- [ ] No N+1 query problems
- [ ] Data integrity maintained
- [ ] Migration scripts provided (if needed)

### 🔄 Integration
**API Integration:**
- [ ] API contracts maintained
- [x] Versioning handled properly
- [x] Backward compatibility considered
- [ ] Rate limiting implemented (if needed)

**Service Communication:**
- [ ] Kafka topics/messages properly defined
- [ ] Error handling in async communication
- [ ] Retry logic implemented
- [ ] Message serialization/deserialization correct

### 📊 Review Outcome
**Review Status:**
- [ ] ✅ Approved
- [x] 🔄 Approved with minor comments
- [x] ⚠️ Requires changes before approval
- [x] ❌ Rejected

**Required Actions:**
- [Action 1]
- [Action 2]
- [Action 3]

**Comments:**
[Detailed comments about the code]

**Follow-up Date:** YYYY-MM-DD
**Reviewer Signature:** _______________________
---------------------------------------------------------------------
 

#### Deployment Checklist
 
## DEPLOYMENT CHECKLIST
**Environment:** [Production/Staging]
**Version:** v[major].[minor].[patch]
**Deployment Date:** YYYY-MM-DD
**Deployment Time:** HH:MM
**Deployment Lead:** [Name]
**Change Window:** HH:MM - HH:MM

### 🔍 Pre-deployment Verification
**Code Verification:**
- [ ] Code is tagged with correct version
- [ ] All tests passing in CI/CD pipeline
- [ ] No unresolved high-severity bugs
- [ ] Security scan passed
- [ ] Performance tests passed

**Configuration Verification:**
- [ ] Environment variables verified
- [ ] Database connection strings correct
- [ ] External service URLs verified
- [ ] Feature flags properly configured

**Infrastructure Verification:**
- [ ] Sufficient resources available
- [ ] Load balancer configured
- [ ] Auto-scaling groups ready
- [ ] Monitoring tools active

### 💾 Backup Procedures
**Data Backup:**
- [ ] Database backup completed
- [ ] File storage backup completed
- [ ] Configuration backup completed
- [ ] Backup integrity verified

**Application Backup:**
- [ ] Previous version artifacts archived
- [ ] Docker images tagged and stored
- [ ] Deployment scripts backed up
- [ ] Rollback plan documented and tested

### 🚀 Deployment Execution
**Step 1: Pre-deployment**
- [ ] Notify stakeholders deployment starting
- [ ] Put monitoring on high alert
- [ ] Disable auto-scaling (if needed)
- [ ] Prepare deployment commands

**Step 2: Deployment**
- [ ] Deploy to 10% of instances (Canary)
- [ ] Wait 5 minutes and monitor
- [ ] Check application health
- [ ] Verify all services running

**Step 3: Full Deployment**
- [ ] Deploy to remaining instances
- [ ] Verify load balancer routing
- [ ] Check database connections
- [ ] Verify service discovery

**Step 4: Post-deployment Verification**
- [ ] Smoke tests passed
- [ ] API endpoints responding
- [ ] UI accessible and functional
- [ ] Background jobs running

### 📊 Monitoring & Alerting
**Performance Monitoring:**
- [ ] CPU/Memory usage normal
- [ ] Response times within SLA
- [ ] Error rates acceptable
- [ ] Database performance normal

**Business Metrics:**
- [ ] User traffic normal
- [ ] Transaction success rate normal
- [ ] Error logs being generated
- [ ] Custom metrics showing expected values

**Alert Verification:**
- [ ] Critical alerts configured
- [ ] Alert notifications working
- [ ] On-call team notified
- [ ] Escalation paths verified

### 🚨 Rollback Preparedness
**Rollback Triggers Monitored:**
- [ ] Error rate > 5%
- [ ] Response time > SLA threshold
- [ ] Critical functionality failure
- [ ] Security vulnerability detected

**Rollback Readiness:**
- [ ] Rollback script tested
- [ ] Previous version ready for deployment
- [ ] Team members on standby
- [ ] Communication plan ready

### 📞 Communication
**Pre-deployment Communication:**
- [ ] Stakeholders notified 24 hours before
- [ ] Maintenance window communicated
- [ ] Impact assessment shared

**During Deployment:**
- [ ] Status updates every 15 minutes
- [ ] Issues communicated immediately
- [ ] Expected completion time updated

**Post-deployment:**
- [ ] Success/failure notification
- [ ] Known issues communicated
- [ ] Next steps communicated

### 📋 Sign-off
**Deployment Outcome:**
- [ ] ✅ Successfully Deployed
- [ ] ⚠️ Deployed with Minor Issues
- [ ] 🔄 Partially Deployed
- [ ] ❌ Rolled Back

**Issues Encountered:**
[Description of any issues during deployment]

**Deployment Lead Sign-off:** ___________________ Time: _________
**DevOps Sign-off:** ____________________________ Time: _________
**Monitoring Team Sign-off:** ___________________ Time: _________

**Post-deployment Monitoring Period:** Next 24 hours
**Next Check-in Time:** YYYY-MM-DD HH:MM
---------------------------------------------------------------------
 

### ภาษาไทย

#### รายการตรวจสอบการตรวจสอบโค้ด
 
## รายการตรวจสอบการตรวจสอบโค้ด
**ผู้ตรวจสอบ:** [ชื่อ]
**ผู้พัฒนา:** [ชื่อ]
**ลิงก์ PR:** [GitHub/GitLab Link]
**รหัสงาน:** CAR-[หมายเลข]
**วันที่ตรวจสอบ:** YYYY-MM-DD

### 🔍 คุณภาพโค้ด
**การอ่านและบำรุงรักษา:**
- [ ] โค้ดสะอาดและจัดระเบียบดี
- [ ] ชื่อตัวแปรและฟังก์ชันมีความหมาย
- [ ] รูปแบบการเขียนโค้ดสม่ำเสมอตลอด
- [ ] ไม่มีโค้ดที่ตายหรือถูกคอมเมนต์ออก
- [ ] การแยกหน้าที่เหมาะสม

**การจัดการข้อผิดพลาด:**
- [ ] ข้อผิดพลาดทั้งหมดได้รับการจัดการอย่างเหมาะสม
- [ ] ข้อความผิดพลาดมีความหมาย
- [ ] การลดคุณภาพอย่างนุ่มนวลที่จำเป็น
- [ ] การบันทึกข้อผิดพลาดที่เหมาะสม

**ประสิทธิภาพ:**
- [ ] ไม่มีจุดคอขวดด้านประสิทธิภาพที่เห็นชัด
- [ ] คำสั่งฐานข้อมูลมีประสิทธิภาพ
- [ ] การใช้แคชอย่างเหมาะสม (ถ้ามี)
- [ ] การใช้หน่วยความจำเหมาะสม

### 🧪 การทดสอบ
**ความครอบคลุมการทดสอบ:**
- [ ] ทดสอบหน่วยครอบคลุมเส้นทางสำคัญ
- [ ] รวมการทดสอบการรวมระบบ
- [ ] ความครอบคลุมการทดสอบตรงตามข้อกำหนดขั้นต่ำ (80%)
- [ ] ทดสอบกรณีขอบเขต

**คุณภาพการทดสอบ:**
- [ ] การทดสอบมีความหมายและไม่ใช่เรื่องเล็กน้อย
- [ ] การทดสอบใช้การตั้งค่า/ล้างข้อมูลที่เหมาะสม
- [ ] การจำลองสถานการณ์ใช้อย่างเหมาะสม
- [ ] การทดสอบเป็นอิสระจากกัน

### 🔒 ความปลอดภัย
**การตรวจสอบข้อมูลนำเข้า:**
- [ ] ข้อมูลนำเข้าผู้ใช้ทั้งหมดได้รับการตรวจสอบ
- [ ] การป้องกัน SQL injection
- [ ] มาตรการป้องกัน XSS
- [ ] การตรวจสอบการอัปโหลดไฟล์ (ถ้ามี)

**การยืนยันตัวตนและการอนุญาต:**
- [ ] การตรวจสอบการยืนยันตัวตนที่เหมาะสม
- [ ] การควบคุมการเข้าถึงตามบทบาท
- [ ] การจัดการเซสชัน
- [ ] การตรวจสอบโทเค็น

**ข้อมูลลับ:**
- [ ] ไม่มีข้อมูลลับในโค้ด
- [ ] ข้อมูลลับเข้ารหัสแล้ว
- [ ] การจัดการรหัสผ่านที่เหมาะสม
- [ ] API keys ปลอดภัย

### 📝 เอกสาร
**เอกสารโค้ด:**
- [ ] ตรรกะที่ซับซ้อนมีคอมเมนต์
- [ ] APIs สาธารณะมีเอกสาร
- [ ] อัปเดต README ถ้าจำเป็น
- [ ] การตัดสินใจทางสถาปัตยกรรมมีเอกสาร

**เอกสาร API:**
- [ ] อัปเดต Swagger/OpenAPI
- [ ] ตัวอย่างคำขอ/การตอบสนอง
- [ ] การตอบสนองข้อผิดพลาดมีเอกสาร
- [ ] ข้อกำหนดการยืนยันตัวตนมีเอกสาร

### 🛠 การดำเนินการทางเทคนิค
**แบ็กเอนด์ (NestJS/Spring Boot):**
- [ ] การใช้รูปแบบเฟรมเวิร์กที่เหมาะสม
- [ ] การแยกชั้นฐานข้อมูล
- [ ] การแยกตรรกะธุรกิจ
- [ ] รูปแบบการตอบสนอง API สม่ำเสมอ

**ฟรอนต์เอนด์ (ReactJS):**
- [ ] โครงสร้างส่วนประกอบตามแนวทาง
- [ ] การจัดการสถานะเหมาะสม
- [ ] การใช้ hooks ที่เหมาะสม
- [ ] การออกแบบตอบสนองนำไปใช้แล้ว

**ฐานข้อมูล:**
- [ ] สร้างดัชนีที่เหมาะสม
- [ ] ไม่มีปัญหาคำสั่ง N+1
- [ ] บูรณภาพข้อมูลได้รับการรักษา
- [ ] มีสคริปต์การย้ายข้อมูล (ถ้าจำเป็น)

### 🔄 การบูรณาการ
**การบูรณาการ API:**
- [ ] สัญญา API ได้รับการรักษา
- [ ] การจัดการเวอร์ชันอย่างเหมาะสม
- [ ] พิจารณาความเข้ากันได้ย้อนหลัง
- [ ] นำการจำกัดอัตราไปใช้ (ถ้าจำเป็น)

**การสื่อสารบริการ:**
- [ ] กำหนดหัวข้อ/ข้อความ Kafka อย่างเหมาะสม
- [ ] การจัดการข้อผิดพลาดในการสื่อสารแบบอะซิงโครนัส
- [ ] นำตรรกะการลองใหม่ไปใช้
- [ ] การแปลงข้อมูลข้อความถูกต้อง

### 📊 ผลการตรวจสอบ
**สถานะการตรวจสอบ:**
- [x] ✅ อนุมัติแล้ว
- [ ] 🔄 อนุมัติพร้อมความคิดเห็นเล็กน้อย
- [ ] ⚠️ ต้องการการเปลี่ยนแปลงก่อนอนุมัติ
- [ ] ❌ ปฏิเสธ

**การดำเนินการที่ต้องการ:**
- [การดำเนินการ 1]
- [การดำเนินการ 2]
- [การดำเนินการ 3]

**ความคิดเห็น:**
[ความคิดเห็นโดยละเอียดเกี่ยวกับโค้ด]

**วันที่ติดตามผล:** YYYY-MM-DD
**ลายเซ็นผู้ตรวจสอบ:** _______________________
---------------------------------------------------------------------
 

#### รายการตรวจสอบการดีพลอย
 
## รายการตรวจสอบการดีพลอย
**สภาพแวดล้อม:** [Production/Staging]
**เวอร์ชัน:** v[หลัก].[รอง].[แก้ไข]
**วันที่ดีพลอย:** YYYY-MM-DD
**เวลาดีพลอย:** HH:MM
**ผู้นำการดีพลอย:** [ชื่อ]
**ช่วงเวลาการเปลี่ยนแปลง:** HH:MM - HH:MM

### 🔍 การตรวจสอบก่อนดีพลอย
**การตรวจสอบโค้ด:**
- [x] โค้ดติดแท็กด้วยเวอร์ชันที่ถูกต้อง
- [x] การทดสอบทั้งหมดผ่านใน CI/CD pipeline
- [ ] ไม่มีข้อบกพร่องระดับสูงที่ยังไม่แก้ไข
- [ ] การตรวจสอบความปลอดภัยผ่าน
- [ ] ทดสอบประสิทธิภาพผ่าน

**การตรวจสอบการกำหนดค่า:**
- [ ] ตรวจสอบตัวแปรสภาพแวดล้อม
- [ ] สตริงการเชื่อมต่อฐานข้อมูลถูกต้อง
- [ ] ตรวจสอบ URLs บริการภายนอก
- [ ] กำหนดค่าสถานะคุณสมบัติอย่างเหมาะสม

**การตรวจสอบโครงสร้างพื้นฐาน:**
- [ ] มีทรัพยากรเพียงพอ
- [x] กำหนดค่าตัวแบ่งเบาแรงงาน
- [x] กลุ่มปรับขนาดอัตโนมัติพร้อม
- [ ] เครื่องมือเฝ้าติดตามทำงาน

### 💾 ขั้นตอนการสำรองข้อมูล
**การสำรองข้อมูล:**
- [x] สำรองฐานข้อมูลเสร็จสิ้น
- [x] สำรองที่เก็บไฟล์เสร็จสิ้น
- [ ] สำรองการกำหนดค่าเสร็จสิ้น
- [ ] ตรวจสอบบูรณภาพการสำรองข้อมูล

**การสำรองแอปพลิเคชัน:**
- [ ] จัดเก็บอาร์ติแฟกต์เวอร์ชันก่อนหน้า
- [ ] ติดแท็กและเก็บ Docker images
- [ ] สำรองสคริปต์การดีพลอย
- [ ] มีเอกสารและทดสอบแผนการย้อนกลับ

### 🚀 การดำเนินการดีพลอย
**ขั้นตอนที่ 1: ก่อนดีพลอย**
- [ ] แจ้งผู้มีส่วนได้ส่วนเสียเริ่มดีพลอย
- [x] ตั้งค่าการเฝ้าติดตามให้อยู่ในระดับสูง
- [x] ปิดการปรับขนาดอัตโนมัติ (ถ้าจำเป็น)
- [ ] เตรียมคำสั่งดีพลอย

**ขั้นตอนที่ 2: การดีพลอย**
- [x] ดีพลอยไปยัง 10% ของอินสแตนซ์ (Canary)
- [ ] รอ 5 นาทีและเฝ้าติดตาม
- [x] ตรวจสอบสุขภาพแอปพลิเคชัน
- [ ] ยืนยันบริการทั้งหมดทำงาน

**ขั้นตอนที่ 3: การดีพลอยเต็มรูปแบบ**
- [ ] ดีพลอยไปยังอินสแตนซ์ที่เหลือ
- [x] ตรวจสอบการกำหนดเส้นทางตัวแบ่งเบาแรงงาน
- [x] ตรวจสอบการเชื่อมต่อฐานข้อมูล
- [ ] ยืนยันการค้นพบบริการ

**ขั้นตอนที่ 4: การตรวจสอบหลังดีพลอย**
- [ ] ทดสอบพื้นฐานผ่าน
- [ ] API endpoints ตอบสนอง
- [x] UI เข้าถึงได้และทำงานได้
- [x] งานพื้นหลังทำงาน

### 📊 การเฝ้าติดตามและการแจ้งเตือน
**การเฝ้าติดตามประสิทธิภาพ:**
- [ ] การใช้งาน CPU/หน่วยความจำปกติ
- [x] เวลาตอบสนองภายใน SLA
- [x] อัตราข้อผิดพลาดยอมรับได้
- [x] ประสิทธิภาพฐานข้อมูลปกติ

**เมตริกธุรกิจ:**
- [x] การจราจรผู้ใช้ปกติ
- [x] อัตราความสำเร็จ
 