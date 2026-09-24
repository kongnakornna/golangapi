# 📦 Universal Module Template — OpenCode + Big Pickle

> Template สำหรับสร้าง Module ใหม่ในอนาคต — ใช้ได้กับทุก domain ไม่ใช่แค่ PDPA
> Framework: R–T–F, T–A–G, B–A–B, C–A–R–E, R–I–S–E

---

## 📋 Table of Contents

1. [โครงสร้าง Template](#1-โครงสร้าง-template)
2. [Framework Reference](#2-framework-reference)
3. [Template Files (Master)](#3-template-files-master)
4. [Scripts อัตโนมัติ](#4-scripts-อัตโนมัติ)
5. [ตัวอย่างการใช้งานจริง](#5-ตัวอย่างการใช้งานจริง)
6. [Prompt สำหรับ Big Pickle / Claude / ChatGPT](#6-prompt-สำหรับ-big-pickle--claude--chatgpt)

---

## 1. โครงสร้าง Template

```
module-template/
├── README.md                              # วิธีใช้ template
├── FRAMEWORKS.md                          # คู่มือ 5 frameworks
├── AGENTS.md.template                     # Project rules (placeholder)
├── .opencode/
│   ├── config.json.template
│   └── commands/
│       ├── 00-init.md.template
│       ├── 01-domain.md.template
│       ├── 02-application.md.template
│       ├── 03-infrastructure.md.template
│       ├── 04-interface.md.template
│       ├── 05-deploy.md.template
│       ├── 06-test.md.template
│       ├── 07-review.md.template
│       ├── 08-fix.md.template
│       ├── 09-docs.md.template
│       ├── 10-seed.md.template
│       ├── 11-monitor.md.template
│       ├── 12-security.md.template
│       ├── 13-migrate.md.template
│       ├── 14-bench.md.template
│       └── 15-rollback.md.template
├── prompts/
│   ├── big-pickle.md.template             # Single prompt สำหรับ Big Pickle
│   ├── claude.md.template                 # Full prompt สำหรับ Claude
│   └── chatgpt-split.md.template          # Split สำหรับ ChatGPT (3 ส่วน)
├── scripts/
│   ├── new-module.sh                      # Bootstrap โปรเจกต์ใหม่
│   ├── fill-template.sh                   # แทนที่ placeholder
│   └── validate.sh                        # ตรวจสอบหลังสร้าง
└── examples/
    └── pdpa/                              # ตัวอย่างที่กรอกแล้ว (reference)
        ├── config.json
        ├── AGENTS.md
        └── commands/*.md
```

---

## 2. Framework Reference

### `FRAMEWORKS.md`

```markdown
# 5 Prompt Frameworks — คู่มืออ้างอิง

> ที่มา: True Digital Academy — Khizer Abbas (Growth Marketing)

## 1. R–T–F: Role – Task – Format

**ใช้เมื่อ**: ต้องการผลลัพธ์ที่ format ชัดเจน

| องค์ประกอบ | คำอธิบาย | ตัวอย่าง |
|-----------|---------|----------|
| **R**ole | ให้ AI เป็นใคร | "Senior Go Developer" |
| **T**ask | ทำอะไร | "สร้าง PDPA Domain Layer" |
| **F**ormat | รูปแบบ output | "ไฟล์ .go เต็ม, ไม่มี pseudo-code" |

**เทมเพลต**:
```
ทำตัวเป็น [ROLE]
ทำงาน [TASK]
ขอผลลัพธ์ในรูปแบบ [FORMAT]
```

**ใช้ใน**: ทุก command file (โครงหลัก)

---

## 2. T–A–G: Task – Action – Goal

**ใช้เมื่อ**: ต้องการให้ AI เข้าใจเจตนา + เป้าหมาย

| องค์ประกอบ | คำอธิบาย |
|-----------|---------|
| **T**ask | งานที่ต้องทำ |
| **A**ction | วิธีการทำ |
| **G**oal | จุดหมายปลายทาง |

**เทมเพลต**:
```
งาน: [TASK]
วิธีทำ: [ACTION]
เป้าหมาย: [GOAL — measurable]
```

**ใช้ใน**: `01-domain.md`, `02-application.md`

---

## 3. B–A–B: Before – After – Bridge

**ใช้เมื่อ**: ต้องการเปลี่ยนสถานะ (migration, refactor, deploy)

| องค์ประกอบ | คำอธิบาย |
|-----------|---------|
| **B**efore | ปัญหา/สถานะปัจจุบัน |
| **A**fter | ภาพที่อยากได้ |
| **B**ridge | ขั้นตอนการเปลี่ยน |

**เทมเพลต**:
```
ตอนนี้: [BEFORE — ปัญหา]
อยากได้: [AFTER — เป้าหมาย]
ขอขั้นตอน: [BRIDGE — step-by-step]
```

**ใช้ใน**: `05-deploy.md`, `13-migrate.md`, `15-rollback.md`

---

## 4. C–A–R–E: Context – Action – Result – Example

**ใช้เมื่อ**: งานซับซ้อนที่ต้องมีบริบท + ตัวอย่าง

| องค์ประกอบ | คำอธิบาย |
|-----------|---------|
| **C**ontext | ฉากหลัง |
| **A**ction | สิ่งที่ให้ทำ |
| **R**esult | ผลลัพธ์ที่ต้องการ |
| **E**xample | ตัวอย่างใกล้เคียง |

**เทมเพลต**:
```
บริบท: [CONTEXT]
ให้ทำ: [ACTION]
ผลลัพธ์: [RESULT]
ตัวอย่าง: [EXAMPLE — reference implementation]
```

**ใช้ใน**: `03-infrastructure.md`, `04-interface.md`, `12-security.md`

---

## 5. R–I–S–E: Role – Input – Steps – Expectation

**ใช้เมื่อ**: งานใหญ่ที่ต้องให้ข้อมูลครบ + ขั้นตอนชัด

| องค์ประกอบ | คำอธิบาย |
|-----------|---------|
| **R**ole | AI เป็นใคร |
| **I**nput | ข้อมูลที่ให้ |
| **S**teps | ขั้นตอน |
| **E**xpectation | ผลลัพธ์สุดท้าย |

**เทมเพลต**:
```
คุณคือ: [ROLE]
ข้อมูล: [INPUT — tech stack, schema, rules]
ขั้นตอน: [STEPS — numbered]
คาดหวัง: [EXPECTATION — definition of done]
```

**ใช้ใน**: `AGENTS.md`, `00-init.md`, `08-fix.md`, `15-rollback.md`

---

## 📊 Mapping: Framework × Command

| Command | หลัก | รอง |
|---------|------|-----|
| 00-init | R–T–F | R–I–S–E |
| 01-domain | R–T–F | T–A–G + R–I–S–E |
| 02-application | R–T–F | T–A–G + C–A–R–E |
| 03-infrastructure | R–T–F | C–A–R–E + R–I–S–E |
| 04-interface | R–T–F | C–A–R–E |
| 05-deploy | B–A–B | R–I–S–E |
| 06-test | R–T–F | — |
| 07-review | R–I–S–E | — |
| 08-fix | R–I–S–E | C–A–R–E |
| 09-docs | R–T–F | C–A–R–E + B–A–B |
| 10-seed | R–T–F | B–A–B |
| 11-monitor | R–T–F | C–A–R–E |
| 12-security | R–I–S–E | C–A–R–E |
| 13-migrate | B–A–B | R–I–S–E |
| 14-bench | R–T–F | R–I–S–E |
| 15-rollback | R–I–S–E | B–A–B |

---

## 🎯 เลือก Framework ตามสถานการณ์

| สถานการณ์ | Framework ที่เหมาะ |
|-----------|-------------------|
| งานเดียวจบ ต้องการ output ชัด | R–T–F |
| ต้องการให้ AI เข้าใจเป้าหมาย | T–A–G |
| เปลี่ยนสถานะ / migrate / deploy | B–A–B |
| งานซับซ้อน ต้องมีตัวอย่าง | C–A–R–E |
| งานใหญ่ ต้องให้ข้อมูลครบ | R–I–S–E |
| **งาน production** | **R–I–S–E + C–A–R–E + R–T–F** |
| **Prototype เร็ว** | **R–T–F เดี่ยว** |
```

---

## 3. Template Files (Master)

### 3.1 `AGENTS.md.template`

```markdown
# AGENTS.md — {{MODULE_NAME}} Module

> **Framework**: R–I–S–E
> **Project**: {{PROJECT_NAME}}
> **Module**: {{MODULE_NAME}}
> **Created**: {{CREATED_DATE}}

---

## 🎭 ROLE

You are a **Principal Software Architect + Senior {{LANGUAGE}} Developer** 
specializing in:
- Clean Architecture & Domain-Driven Design
- {{DOMAIN_DESCRIPTION}}
- {{INDUSTRY_COMPLIANCE}} (e.g., PDPA, GDPR, HIPAA)
- Event-Driven Architecture ({{MESSAGE_BROKER}})
{{#if BLOCKCHAIN}}- Blockchain-backed audit trails ({{BLOCKCHAIN_NETWORK}}){{/if}}

You write **production-ready {{LANGUAGE}} {{LANGUAGE_VERSION}}** code. 
No pseudo-code. No TODOs.

---

## 📥 INPUT — Project Context

### Tech Stack
| Component | Version | Purpose |
|-----------|---------|---------|
| Language | {{LANGUAGE_VERSION}} | Primary |
| Framework | {{WEB_FRAMEWORK}} | HTTP |
| ORM | {{ORM}} | Data access |
| Database | {{PRIMARY_DB}} | Persistence |
{{#each ADDITIONAL_SERVICES}}
| {{name}} | {{version}} | {{purpose}} |
{{/each}}

### Module Structure
```
{{MODULE_PATH}}/
├── domain/          # Pure business logic (NO external deps)
├── application/     # Use cases (depends only on domain)
├── infrastructure/  # External adapters
└── interfaces/      # HTTP, WebSocket, CLI
```

### Dependency Direction (STRICT)
```
interfaces → application → domain ← infrastructure
```

### Database Schema ({{TABLE_COUNT}} tables, prefix `{{TABLE_PREFIX}}`)
{{#each TABLES}}
- `{{table_name}}` → {{description}}
{{/each}}

### Use Cases ({{USE_CASE_COUNT}} total)
{{#each USE_CASES}}
{{@index}}. {{name}} — {{description}}
{{/each}}

### Business Rules (INVIOLABLE)
{{#each BUSINESS_RULES}}
- {{this}}
{{/each}}

---

## 🪜 STEPS — Standard Workflow

### For Every Command
1. **Read** this file + relevant source files
2. **Plan** in Plan mode (Shift+Tab) if non-trivial
3. **Build** in Build mode
4. **Verify**: `{{BUILD_CMD}}` exits 0
5. **Lint**: `{{LINT_CMD}}` exits 0
6. **Test**: `{{TEST_CMD}}` passes
7. **Report**: files + build + self-check

### Implementation Order
```
Phase 0 → Bootstrap (go.mod, .env, docker-compose)
Phase 1 → Domain Layer
Phase 2 → Application Layer
Phase 3 → Infrastructure Layer
Phase 4 → Interface Layer + Entry Points
Phase 5 → Deployment
Phase 6 → Testing
Phase 7 → Self-review
```

---

## ✅ EXPECTATION — Definition of Done

### Every Change Must Satisfy
- [ ] `{{BUILD_CMD}}` exits 0
- [ ] `{{LINT_CMD}}` exits 0
- [ ] `{{TEST_CMD}}` passes
- [ ] Test coverage > {{COVERAGE_THRESHOLD}}%
- [ ] No pseudo-code, no `// TODO` in critical paths
- [ ] Domain layer has **ZERO** infrastructure imports
- [ ] Every entity has `NewXxx()` constructor
- [ ] Every Use Case has `Execute(ctx, input)` signature
- [ ] Every public method takes `context.Context`
- [ ] Every state change writes to audit trail

---

## 📏 Coding Standards

### Style
- {{STYLE_GUIDE}} (e.g., idiomatic Go, gofmt)
- Constructor injection
- Interface-based dependencies
- Context propagation
- Error wrapping with `%w`

### Naming
| Type | Pattern | Example |
|------|---------|---------|
| Entity | `Xxx` | `{{EXAMPLE_ENTITY}}` |
| Value Object | `Xxx` | `{{EXAMPLE_VO}}` |
| Use Case | `XxxUseCase` | `{{EXAMPLE_UC}}` |
| Repository | `XxxRepository` | `{{EXAMPLE_REPO}}` |
| Domain Error | `ErrXxx` | `{{EXAMPLE_ERR}}` |

### Forbidden
- ❌ Global variables / singletons
- ❌ Business logic in infrastructure
- ❌ Domain importing infrastructure
- ❌ Bare `errors.New()` outside `domain/errors/`
- ❌ `panic()` in non-startup code
- ❌ Pseudo-code or `// TODO`

---

## 🚨 Emergency Fallbacks

### If Build Fails
1. Read error carefully
2. `{{BUILD_CMD}} 2>&1 | head -50`
3. Fix one file at a time
4. Re-run until clean

### If Test Fails
1. `{{TEST_CMD}} -v -run TestName`
2. Add `t.Logf()` for debugging
3. Never skip or comment out tests

### If Unsure
1. Check this AGENTS.md
2. Check sibling files for patterns
3. **ASK the user** before inventing new patterns

---

## 🎯 Success Criteria

- ✅ `docker-compose up -d` brings all services healthy
- ✅ Migration creates {{TABLE_COUNT}} tables
- ✅ API server starts on `:{{API_PORT}}`
- ✅ All workers connect to {{MESSAGE_BROKER}}
- ✅ Full flow works end-to-end
{{#if BLOCKCHAIN}}- ✅ Blockchain proof recorded for critical events{{/if}}
```

---

### 3.2 `.opencode/config.json.template`

```json
{
  "$schema": "https://opencode.ai/config.json",
  "instructions": [
    "AGENTS.md",
    ".opencode/commands/*.md"
  ],
  "model": "{{MODEL}}",
  "permissions": {
    "edit": true,
    "bash": true,
    "read": true,
    "write": true
  },
  "commands": {
    "init":          ".opencode/commands/00-init.md",
    "domain":        ".opencode/commands/01-domain.md",
    "app":           ".opencode/commands/02-application.md",
    "infra":         ".opencode/commands/03-infrastructure.md",
    "interface":     ".opencode/commands/04-interface.md",
    "deploy":        ".opencode/commands/05-deploy.md",
    "test":          ".opencode/commands/06-test.md",
    "review":        ".opencode/commands/07-review.md",
    "fix":           ".opencode/commands/08-fix.md",
    "docs":          ".opencode/commands/09-docs.md",
    "seed":          ".opencode/commands/10-seed.md",
    "monitor":       ".opencode/commands/11-monitor.md",
    "security":      ".opencode/commands/12-security.md",
    "migrate":       ".opencode/commands/13-migrate.md",
    "bench":         ".opencode/commands/14-bench.md",
    "rollback":      ".opencode/commands/15-rollback.md"
  },
  "variables": {
    "MODULE_NAME": "{{MODULE_NAME}}",
    "MODULE_PATH": "{{MODULE_PATH}}",
    "TABLE_PREFIX": "{{TABLE_PREFIX}}",
    "API_PORT": "{{API_PORT}}"
  }
}
```

---

### 3.3 `.opencode/commands/00-init.md.template`

```markdown
---
description: Bootstrap {{MODULE_NAME}} Module — {{LANGUAGE}} project skeleton
framework: R–T–F + R–I–S–E
---

# R–T–F: Bootstrap {{MODULE_NAME}}

## ROLE
Senior {{LANGUAGE}} Developer setting up a Clean Architecture project.

## TASK
Create the project skeleton for {{MODULE_NAME}} with {{LANGUAGE}} {{LANGUAGE_VERSION}}.

## FORMAT
Full runnable files. No pseudo-code.

---

# R–I–S–E

## ROLE
You are bootstrapping a new module. Assume {{MODULE_PATH}} doesn't exist yet.

## INPUT
- Language: {{LANGUAGE}} {{LANGUAGE_VERSION}}
- Module path: `{{MODULE_PATH}}`
- Dependencies:
{{#each DEPENDENCIES}}
  - {{this}}
{{/each}}

## STEPS

### 1. Create `{{DEPENDENCY_FILE}}`
Include all required dependencies:
{{#each DEPENDENCIES}}
- {{this}}
{{/each}}

### 2. Create `.env.example` and `.env`
{{#each ENV_VARS}}
# ===== {{group}} =====
{{key}}={{default}}
{{/each}}

### 3. Create directory tree
```
{{MODULE_PATH}}/
├── domain/
│   ├── entity/
│   ├── value_object/
│   ├── repository/
│   ├── service/
│   ├── event/
│   └── errors/
├── application/
├── infrastructure/
│   ├── persistence/
│   ├── messaging/
│   ├── search/
│   ├── services/
│   └── scheduler/
└── interfaces/
    ├── http/
    ├── websocket/
    └── localization/
        └── locales/
```

### 4. Create `.gitignore`
Standard {{LANGUAGE}} ignores + `.env`

### 5. Run
```bash
{{SETUP_CMD}}
{{BUILD_CMD}}
```

## EXPECTATION
- [ ] Dependency file exists
- [ ] Directory tree created
- [ ] `.env.example` + `.env` exist
- [ ] `.gitignore` covers `.env`
- [ ] `{{BUILD_CMD}}` exits 0
- [ ] Report: files + build + self-check

## 🚀 BEGIN
1. Create dependency file
2. Create .env
3. Create directory tree
4. Run setup
5. Report
```

---

### 3.4 `.opencode/commands/01-domain.md.template`

```markdown
---
description: สร้าง {{MODULE_NAME}} Domain Layer
framework: R–T–F + T–A–G + R–I–S–E
---

# R–T–F

## ROLE
Senior {{LANGUAGE}} Developer specializing in Domain-Driven Design.

## TASK
Create the complete **Domain Layer** for {{MODULE_NAME}}.

## FORMAT
Full files (no snippets). Final report with build status.

---

# T–A–G

## TASK
Build {{DOMAIN_ARTIFACT_COUNT}} domain artifacts.

## ACTION
Follow R–I–S–E steps below exactly.

## GOAL
`{{BUILD_CMD}} ./{{MODULE_PATH}}/domain/...` exits 0.
Zero imports from application/infrastructure/interfaces.

---

# R–I–S–E

## INPUT (Preconditions)
- `AGENTS.md` read
- Dependency file exists (from `/init`)
- Directory tree exists

## STEPS

### Step 1 — Value Objects
Create in `{{MODULE_PATH}}/domain/value_object/`:

{{#each VALUE_OBJECTS}}
#### `{{file_name}}.go`
```{{language}}
// {{description}}
{{code}}
```
{{/each}}

### Step 2 — Domain Errors
Create `{{MODULE_PATH}}/domain/errors/errors.go`:
```{{language}}
{{errors_code}}
```

### Step 3 — Entities
{{#each ENTITIES}}
#### `{{MODULE_PATH}}/domain/entity/{{file_name}}.go`
```{{language}}
{{code}}
```
{{/each}}

### Step 4 — Repository Interfaces
{{#each REPOSITORIES}}
#### `{{MODULE_PATH}}/domain/repository/{{file_name}}.go`
```{{language}}
{{code}}
```
{{/each}}

### Step 5 — Domain Services
{{#each SERVICES}}
#### `{{MODULE_PATH}}/domain/service/{{file_name}}.go`
```{{language}}
{{code}}
```
{{/each}}

### Step 6 — Domain Events
{{#if EVENTS}}
Create `{{MODULE_PATH}}/domain/event/event.go`:
```{{language}}
{{events_code}}
```
{{/if}}

## EXPECTATION
- [ ] {{VALUE_OBJECT_COUNT}} value objects
- [ ] {{ENTITY_COUNT}} entities
- [ ] {{REPO_COUNT}} repository interfaces
- [ ] {{SERVICE_COUNT}} domain services
- [ ] 1 errors file
- [ ] `{{BUILD_CMD}}` exits 0
- [ ] Zero imports from upper layers

## VERIFY
```bash
{{BUILD_CMD}} ./{{MODULE_PATH}}/domain/...
grep -r "infrastructure" {{MODULE_PATH}}/domain/ && echo "FAIL" || echo "PASS"
```

## REPORT
```
✅ Domain Layer complete
Files: {{TOTAL_FILES}}
Build: PASS
Boundary: PASS
Next: run /app
```
```

---

### 3.5 Command files อื่นๆ (อ้างอิง)

Command files `02` ถึง `15` ใช้โครงสร้างเดียวกันกับที่สร้างไว้ในคำตอบก่อนหน้า โดยแทนที่ด้วย placeholder:

| Command | Placeholder หลัก |
|---------|------------------|
| `02-application.md` | `{{USE_CASES}}`, `{{DTOs}}`, `{{PORTS}}` |
| `03-infrastructure.md` | `{{MODELS}}`, `{{REPO_IMPLS}}`, `{{CONSUMERS}}` |
| `04-interface.md` | `{{HANDLERS}}`, `{{ROUTES}}`, `{{ENTRY_POINTS}}` |
| `05-deploy.md` | `{{DOCKER_SERVICES}}`, `{{MIGRATION_SQL}}` |
| `06-test.md` | `{{TEST_FILES}}`, `{{COVERAGE_THRESHOLD}}` |
| `07-review.md` | `{{CHECKS}}`, `{{BOUNDARY_CHECKS}}` |
| `08-fix.md` | `{{ERROR_CATEGORIES}}` |
| `09-docs.md` | `{{ENDPOINTS}}`, `{{DIAGRAMS}}` |
| `10-seed.md` | `{{SEED_ENTITIES}}`, `{{FAKE_DATA}}` |
| `11-monitor.md` | `{{METRICS}}`, `{{HEALTH_CHECKS}}` |
| `12-security.md` | `{{THREAT_MODEL}}`, `{{FIXES}}` |
| `13-migrate.md` | `{{LEGACY_SCHEMA}}`, `{{MAPPING}}` |
| `14-bench.md` | `{{BENCHMARKS}}`, `{{LOAD_TESTS}}` |
| `15-rollback.md` | `{{ROLLBACK_STEPS}}`, `{{RUNBOOKS}}` |

---

## 4. Scripts อัตโนมัติ

### `scripts/new-module.sh`

```bash
#!/usr/bin/env bash
# Bootstrap โปรเจกต์ใหม่จาก template
# Usage: ./new-module.sh <project-name> <module-name>

set -euo pipefail

PROJECT_NAME="${1:-myproject}"
MODULE_NAME="${2:-mymodule}"
TARGET_DIR="${3:-./${PROJECT_NAME}}"

if [ -d "$TARGET_DIR" ]; then
  echo "❌ Directory $TARGET_DIR already exists"
  exit 1
fi

echo "📦 Creating project: $PROJECT_NAME (module: $MODULE_NAME)"
mkdir -p "$TARGET_DIR"
cp -r module-template/. "$TARGET_DIR/"

cd "$TARGET_DIR"

# Fill placeholders
./scripts/fill-template.sh "$PROJECT_NAME" "$MODULE_NAME"

echo "✅ Project created at $TARGET_DIR"
echo ""
echo "Next steps:"
echo "  1. cd $TARGET_DIR"
echo "  2. Edit AGENTS.md to fill in business rules"
echo "  3. Edit .opencode/commands/00-init.md with tech stack"
echo "  4. opencode"
echo "  5. /init"
```

### `scripts/fill-template.sh`

```bash
#!/usr/bin/env bash
# แทนที่ placeholder ทั่วไปในไฟล์ .template
# Usage: ./fill-template.sh <project-name> <module-name>

set -euo pipefail

PROJECT_NAME="${1:-myproject}"
MODULE_NAME="${2:-mymodule}"
MODULE_PATH="${3:-internal/modules/${MODULE_NAME}}"
TABLE_PREFIX="${4:-${MODULE_NAME}_}"
CREATED_DATE="$(date +%Y-%m-%d)"

# Default replacements
declare -A REPLACEMENTS=(
  ["{{PROJECT_NAME}}"]="$PROJECT_NAME"
  ["{{MODULE_NAME}}"]="$MODULE_NAME"
  ["{{MODULE_PATH}}"]="$MODULE_PATH"
  ["{{TABLE_PREFIX}}"]="$TABLE_PREFIX"
  ["{{CREATED_DATE}}"]="$CREATED_DATE"
  ["{{MODEL}}"]="opencode/big-pickle"
  ["{{LANGUAGE}}"]="Go"
  ["{{LANGUAGE_VERSION}}"]="1.21"
  ["{{WEB_FRAMEWORK}}"]="Gin"
  ["{{ORM}}"]="GORM"
  ["{{PRIMARY_DB}}"]="PostgreSQL 15"
  ["{{MESSAGE_BROKER}}"]="Kafka"
  ["{{API_PORT}}"]="8080"
  ["{{COVERAGE_THRESHOLD}}"]="70"
  ["{{BUILD_CMD}}"]="go build"
  ["{{LINT_CMD}}"]="golangci-lint run"
  ["{{TEST_CMD}}"]="go test ./... -race"
  ["{{SETUP_CMD}}"]="go mod tidy"
  ["{{DEPENDENCY_FILE}}"]="go.mod"
  ["{{STYLE_GUIDE}}"]="idiomatic Go, gofmt"
)

# Apply to all .template files
find . -name "*.template" -type f | while read -r file; do
  echo "  filling: $file"
  tmpfile="${file}.tmp"
  cp "$file" "$tmpfile"

  for placeholder in "${!REPLACEMENTS[@]}"; do
    value="${REPLACEMENTS[$placeholder]}"
    # macOS compatible sed
    sed -i.bak "s|${placeholder}|${value}|g" "$tmpfile" 2>/dev/null || \
      sed -i "s|${placeholder}|${value}|g" "$tmpfile"
  done

  # Remove .template extension
  target="${file%.template}"
  mv "$tmpfile" "$target"
  rm -f "${tmpfile}.bak"
done

echo "✅ Placeholders filled"
```

### `scripts/validate.sh`

```bash
#!/usr/bin/env bash
# ตรวจสอบหลังสร้างโปรเจกต์
set -euo pipefail

echo "🔍 Validating project..."
FAIL=0

# 1. Files exist
for f in AGENTS.md .opencode/config.json; do
  if [ ! -f "$f" ]; then
    echo "❌ Missing: $f"
    FAIL=1
  fi
done

# 2. No leftover placeholders
if grep -r "{{" --include="*.md" --include="*.json" --include="*.go" . 2>/dev/null | grep -v ".template"; then
  echo "⚠️ Unfilled placeholders found (see above)"
  FAIL=1
fi

# 3. Commands count
CMD_COUNT=$(ls .opencode/commands/*.md 2>/dev/null | wc -l)
echo "  Commands: $CMD_COUNT"
if [ "$CMD_COUNT" -lt 16 ]; then
  echo "⚠️ Expected 16 commands, found $CMD_COUNT"
fi

# 4. Build
if command -v go >/dev/null 2>&1; then
  go build ./... && echo "  Build: PASS" || { echo "  Build: FAIL"; FAIL=1; }
fi

if [ $FAIL -eq 0 ]; then
  echo "✅ Validation passed"
else
  echo "❌ Validation failed"
  exit 1
fi
```

---

## 5. ตัวอย่างการใช้งานจริง

### 🎬 Scenario: สร้างระบบ E-Commerce Module

```bash
# Step 1: Bootstrap
./scripts/new-module.sh ecommerce-platform order

# Step 2: Edit config
cd ecommerce-platform
vim AGENTS.md
# แทนที่ business rules:
#   - Order expires in 24 hours
#   - Cannot cancel shipped order
#   - CANCELLED → refund within 7 days
#   - PAID → SHIPPED → DELIVERED → COMPLETED

# Step 3: OpenCode
opencode

# Step 4: Run commands
/init           # Bootstrap project
/domain         # Domain Layer (Order, OrderItem, Payment)
/app            # Use Cases (CreateOrder, PayOrder, ShipOrder, etc.)
/infra          # GORM, Kafka, Redis
/interface      # HTTP handlers
/deploy         # Docker + migration
/test           # Tests
/review         # Self-review
/security       # Audit
/docs           # Swagger
```

### 📋 Fill-in Template สำหรับ Order Module

**Input data** (ที่ต้องกรอก):

```yaml
MODULE_NAME: order
MODULE_PATH: internal/modules/order
TABLE_PREFIX: order_
API_PORT: 8081

ENTITIES:
  - Order (ID, UserID, Status, Total, CreatedAt)
  - OrderItem (ID, OrderID, ProductID, Qty, Price)
  - Payment (ID, OrderID, Status, Amount, Method)
  - Shipment (ID, OrderID, TrackingNo, Carrier, ShippedAt)

VALUE_OBJECTS:
  - OrderStatus (PENDING, PAID, SHIPPED, DELIVERED, CANCELLED)
  - PaymentStatus (PENDING, SUCCESS, FAILED, REFUNDED)
  - PaymentMethod (CREDIT_CARD, PROMPTPAY, COD)

USE_CASES:
  1. CreateOrder
  2. PayOrder
  3. CancelOrder
  4. ShipOrder
  5. DeliverOrder
  6. GetOrderStatus
  7. ListUserOrders

BUSINESS_RULES:
  - Order expires in 24 hours if unpaid
  - Cannot cancel SHIPPED or DELIVERED
  - Refund within 7 days of CANCELLED
  - Payment must succeed before SHIPPED
```

**Output**: Command files ที่แทนที่ placeholder แล้ว

---

## 6. Prompt สำหรับ Big Pickle / Claude / ChatGPT

### 6.1 `prompts/big-pickle.md.template`

```markdown
╔══════════════════════════════════════════════════════════╗
║  R–T–F: Role – Task – Format                             ║
╚══════════════════════════════════════════════════════════╝

ROLE:
You are a Principal Software Architect + Senior {{LANGUAGE}} Developer 
specializing in Clean Architecture, DDD, and {{DOMAIN_DESCRIPTION}}.

TASK:
Build a production-ready {{MODULE_NAME}} Module in {{LANGUAGE}} {{LANGUAGE_VERSION}}.

FORMAT:
- Full runnable code per file
- Path header before each code block
- Order: Domain → Application → Infrastructure → Interface

╔══════════════════════════════════════════════════════════╗
║  R–I–S–E: Role – Input – Steps – Expectation             ║
╚══════════════════════════════════════════════════════════╝

ROLE:
You are a Principal Architect + Senior {{LANGUAGE}} Developer.

INPUT:
{{#each INPUT_SECTIONS}}
### {{title}}
{{content}}
{{/each}}

STEPS:
Phase 1: Foundation (dependency file, .env, docker-compose)
Phase 2: Domain Layer (VOs, entities, repos, services, errors)
Phase 3: Application Layer ({{USE_CASE_COUNT}} use cases, DTOs)
Phase 4: Infrastructure Layer (models, repos, cache, kafka, services)
Phase 5: Interface Layer (handlers, routes, entry points)
Phase 6: Testing (unit + integration)
Phase 7: Self-review

EXPECTATION:
- ✅ Compiles with `{{BUILD_CMD}}`
- ✅ `docker-compose up` works
- ✅ Test coverage > {{COVERAGE_THRESHOLD}}%
- ✅ Zero boundary violations
- ✅ Self-check checklist at end

╔══════════════════════════════════════════════════════════╗
║  🧠 REASONING (Big Pickle Specific)                       ║
╚══════════════════════════════════════════════════════════╝

Before code:
1. List all assumptions
2. Explain schema FK relationships
3. Explain layer dependency direction
4. Identify trade-offs
5. Plan implementation order

🚀 BEGIN
Start with Architecture Decision Summary → Phase 1 → 7.
```

### 6.2 `prompts/claude.md.template`

```markdown
# Claude — Full Prompt for {{MODULE_NAME}}

ใช้เมื่อ: ต้องการ output ทั้งหมดในครั้งเดียว (Claude context ใหญ่)

---

[Copy ทั้ง Prompt Big Pickle ด้านบน]

**เพิ่มเติมสำหรับ Claude**:

## 📚 Extended Context

Claude, คุณมี context window ใหญ่พอที่จะเข้าใจทั้งหมดในครั้งเดียว 
กรุณา:

1. วิเคราะห์ requirement ทั้งหมดก่อน
2. ระบุ assumptions ที่ชัดเจน
3. สร้าง Mermaid diagrams สำหรับ architecture
4. สร้างโค้ดทุก layer ตามลำดับ
5. ให้ reasoning ระหว่างทาง
6. จบด้วย self-check checklist

## 🎯 Output Checklist

- [ ] Domain Layer ({{DOMAIN_FILE_COUNT}} files)
- [ ] Application Layer ({{APP_FILE_COUNT}} files)
- [ ] Infrastructure Layer ({{INFRA_FILE_COUNT}} files)
- [ ] Interface Layer ({{INTERFACE_FILE_COUNT}} files)
- [ ] Entry Points ({{ENTRY_COUNT}} binaries)
- [ ] Config (go.mod, .env, docker-compose.yml)
- [ ] Migration SQL
- [ ] Writing Plan 7 phases
- [ ] Mermaid diagrams
```

### 6.3 `prompts/chatgpt-split.md.template`

```markdown
# ChatGPT — Split Prompts for {{MODULE_NAME}}

ใช้เมื่อ: ChatGPT มี context limit → แบ่งเป็น 3 ส่วน

---

## 📩 Prompt 1 (Domain Layer)

[Copy ส่วน R–T–F + R–I–S–E + INPUT ด้านบน]

**ขอบเขต**: สร้างเฉพาะ Domain Layer
- Value Objects
- Entities
- Repository Interfaces
- Domain Services
- Domain Errors
- Domain Events

**Output**: โค้ดเต็ม + build verification

---

## 📩 Prompt 2 (Application + Infrastructure)

[Copy ส่วน R–T–F + R–I–S–E ด้านบน + เพิ่ม]

**ก่อนหน้า**: Domain Layer สร้างแล้ว (ดูโค้ดด้านบน)

**ขอบเขต**: สร้าง Application + Infrastructure
- Use Cases ทั้งหมด
- DTOs
- Ports (local interfaces)
- GORM Models + Repos
- Cache/Producer/Consumers
- Services (email, LLM, blockchain)

**Output**: โค้ดเต็ม + build verification

---

## 📩 Prompt 3 (Interface + Entry Points)

[Copy ส่วน R–T–F + R–I–S–E ด้านบน + เพิ่ม]

**ก่อนหน้า**: Domain + Application + Infrastructure สร้างแล้ว

**ขอบเขต**: สร้าง Interface + Entry Points
- HTTP Handlers
- Routes
- WebSocket Hub
- i18n
- main.go ทั้งหมด

**Output**: โค้ดเต็ม + build + run verification
```

---

## 📖 README.md (ของ Template เอง)

```markdown
# Module Template — Universal

Template สำหรับสร้าง Module ใหม่ในภาษา Go (หรือภาษาอื่นๆ) 
ตาม Clean Architecture + DDD

## 🚀 Quick Start

### ครั้งแรก: Setup template

```bash
# 1. Clone หรือ download template
git clone <template-repo> module-template
cd module-template

# 2. ทดสอบว่าทำงานได้
./scripts/new-module.sh demo demo-module /tmp/demo
```

### สร้างโปรเจกต์ใหม่

```bash
# วิธีที่ 1: ใช้ script
./scripts/new-module.sh myproject mymodule

# วิธีที่ 2: Copy manual
cp -r module-template my-project
cd my-project
./scripts/fill-template.sh my-project my-module
```

## 📚 เอกสาร

- [FRAMEWORKS.md](FRAMEWORKS.md) — 5 prompt frameworks
- [AGENTS.md](AGENTS.md) — project rules
- [.opencode/commands/](.opencode/commands/) — 16 commands

## 🎯 Workflow

```
1. ./scripts/new-module.sh <project> <module>
2. แก้ AGENTS.md (business rules)
3. แก้ .opencode/commands/00-init.md (tech stack)
4. opencode
5. /init → /domain → /app → /infra → /interface
6. /test → /review → /security → /docs
7. /deploy
```

## 🔧 16 Commands

| # | Command | Framework | Purpose |
|---|---------|-----------|---------|
| 00 | `/init` | R–T–F + R–I–S–E | Bootstrap |
| 01 | `/domain` | R–T–F + T–A–G | Domain Layer |
| 02 | `/app` | R–T–F + T–A–G + C–A–R–E | Application |
| 03 | `/infra` | R–T–F + C–A–R–E | Infrastructure |
| 04 | `/interface` | R–T–F + C–A–R–E | HTTP + WS |
| 05 | `/deploy` | B–A–B | Docker + Migration |
| 06 | `/test` | R–T–F | Unit Tests |
| 07 | `/review` | R–I–S–E | Self Review |
| 08 | `/fix` | R–I–S–E + C–A–R–E | Debug |
| 09 | `/docs` | R–T–F + C–A–R–E | Swagger |
| 10 | `/seed` | R–T–F + B–A–B | Test Data |
| 11 | `/monitor` | R–T–F + C–A–R–E | Observability |
| 12 | `/security` | R–I–S–E + C–A–R–E | Hardening |
| 13 | `/migrate` | B–A–B + R–I–S–E | Legacy Migration |
| 14 | `/bench` | R–T–F + R–I–S–E | Performance |
| 15 | `/rollback` | R–I–S–E + B–A–B | Recovery |

## 🎨 Customization

### เพิ่ม Domain ใหม่
1. Copy `.opencode/commands/00-init.md.template` → `00-init-<domain>.md`
2. แก้ INPUT section ให้เหมาะกับ domain
3. เพิ่มใน `config.json`

### เพิ่ม Command ใหม่
1. สร้างไฟล์ `.opencode/commands/16-<name>.md`
2. ใช้ structure เดียวกับ commands อื่นๆ
3. Update `config.json`

### เปลี่ยนภาษา
1. Copy all `.template` files
2. แทนที่ `{{LANGUAGE}}` = Rust, Python, TypeScript ฯลฯ
3. ปรับ commands ให้เหมาะกับ ecosystem

## 📊 หลักการทำงาน

1. **R–I–S–E ฝังใน AGENTS.md** — input ครบ + steps ชัด + expectation วัดได้
2. **R–T–F ในทุก command** — role/task/format ชัดเจน
3. **B–A–B สำหรับ migration** — before/after/bridge
4. **C–A–R–E สำหรับ complex tasks** — context + example
5. **T–A–G สำหรับ objectives** — task + action + goal

## ⚠️ ข้อควรระวัง

- Template นี้ **ไม่ใช่ boilerplate code** — เป็น prompt engineering template
- ต้อง **กรอก placeholder** ก่อนใช้ (หรือใช้ script)
- **ตรวจสอบ `go build`** หลังทุก command
- **Review ก่อน commit** — AI อาจ hallucinate

## 📝 License

MIT
```

---

## 📊 สรุปทั้งหมด — Template Package

```
📦 module-template/
├── 📄 README.md                            # วิธีใช้
├── 📄 FRAMEWORKS.md                        # 5 frameworks
├── 📄 AGENTS.md.template                   # Project rules
├── 📁 .opencode/
│   ├── 📄 config.json.template
│   └── 📁 commands/                        # 16 commands
│       ├── 00-init.md.template
│       ├── 01-domain.md.template
│       ├── 02-application.md.template
│       ├── 03-infrastructure.md.template
│       ├── 04-interface.md.template
│       ├── 05-deploy.md.template
│       ├── 06-test.md.template
│       ├── 07-review.md.template
│       ├── 08-fix.md.template
│       ├── 09-docs.md.template
│       ├── 10-seed.md.template
│       ├── 11-monitor.md.template
│       ├── 12-security.md.template
│       ├── 13-migrate.md.template
│       ├── 14-bench.md.template
│       └── 15-rollback.md.template
├── 📁 prompts/                             # Prompts หลาย AI
│   ├── big-pickle.md.template
│   ├── claude.md.template
│   └── chatgpt-split.md.template
├── 📁 scripts/                             # Automation
│   ├── new-module.sh
│   ├── fill-template.sh
│   └── validate.sh
└── 📁 examples/                            # ตัวอย่างที่กรอกแล้ว
    └── pdpa/
        ├── config.json
        ├── AGENTS.md
        └── commands/*.md
```

---

## 🎯 วิธีนำไปใช้กับโปรเจกต์ใหม่ (Step-by-Step)

```bash
# ═══════════════════════════════════════════════════
# ครั้งที่ 1: Setup template
# ═══════════════════════════════════════════════════
git clone <template-repo> ~/templates/module-template
cd ~/templates/module-template
chmod +x scripts/*.sh

# ═══════════════════════════════════════════════════
# ครั้งที่ 2+: สร้างโปรเจกต์ใหม่
# ═══════════════════════════════════════════════════
~/templates/module-template/scripts/new-module.sh \
  inventory-system \
  stock

cd inventory-system

# ═══════════════════════════════════════════════════
# ปรับแต่ง 3 ไฟล์
# ═══════════════════════════════════════════════════
vim AGENTS.md                          # business rules
vim .opencode/commands/00-init.md      # tech stack
vim .opencode/commands/01-domain.md    # entities

# ═══════════════════════════════════════════════════
# รัน OpenCode
# ═══════════════════════════════════════════════════
opencode

# ใน TUI:
/init
/domain
/app
/infra
/interface
/test
/review
/security
/docs
/deploy

# ═══════════════════════════════════════════════════
# ตรวจสอบ
# ═══════════════════════════════════════════════════
~/templates/module-template/scripts/validate.sh
```

---

## 💡 ตัวอย่างการกรอก Placeholder

### ก่อน (Template)

```markdown
# AGENTS.md — {{MODULE_NAME}} Module

## ROLE
You are a Principal Software Architect + Senior {{LANGUAGE}} Developer 
specializing in {{DOMAIN_DESCRIPTION}}.
```

### หลัง (Filled)

```markdown
# AGENTS.md — Stock Module

## ROLE
You are a Principal Software Architect + Senior Go Developer 
specializing in inventory management and warehouse operations.
```

---

## 🎁 สรุปสิ่งที่ได้

| ไฟล์ | จำนวน | ประโยชน์ |
|------|-------|---------|
| `AGENTS.md.template` | 1 | Project rules ที่ใช้ซ้ำได้ |
| `.opencode/commands/*.md.template` | 16 | Commands ครบ workflow |
| `prompts/*.md.template` | 3 | Prompts สำหรับ Big Pickle/Claude/ChatGPT |
| `scripts/*.sh` | 3 | Automate setup + validate |
| `FRAMEWORKS.md` | 1 | คู่มือ 5 frameworks |
| `examples/pdpa/` | 1 set | ตัวอย่างจริง |

**ผลลัพธ์**: สร้าง Module ใหม่ได้ใน **5 นาที** (แทนที่จะเป็น 5 ชั่วโมง) พร้อม pattern ที่สม่ำเสมอทุกครั้ง

ถ้าต้องการให้ผม **สร้างไฟล์ .template เต็มทุกไฟล์** (16 commands + 3 prompts + scripts) พร้อม content ครบ กรุณาบอกได้เลยครับ — จะทำให้เสร็จทีละไฟล์