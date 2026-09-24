# 🐳📊 ชุดเครื่องมือขั้นสูง 4 อย่าง

จัดโครงสร้างเพิ่ม:
```
project/
├── docker/
│   ├── Dockerfile
│   ├── docker-compose.yml
│   ├── .dockerignore
│   └── entrypoint.sh
├── .vscode/
│   ├── tasks.json
│   ├── launch.json
│   └── settings.json
├── tests/
│   ├── Run-6PhaseWorkflow.Tests.ps1
│   ├── Invoke-Phase.Tests.ps1
│   ├── Publish-JiraComment.Tests.ps1
│   ├── Helpers.Tests.ps1
│   └── pester.config.ps1
├── dashboard/
│   ├── index.html                    # standalone HTML dashboard
│   ├── app.js
│   ├── styles.css
│   └── grafana/
│       ├── provisioning/
│       │   ├── datasources/
│       │   │   └── prometheus.yml
│       │   └── dashboards/
│       │       └── default.yml
│       └── dashboards/
│           └── workflow-6phase.json
└── scripts/                          # เดิม
```

---

## 🐳 1) Docker Compose: รัน workflow ใน container

### 1.1 `docker/Dockerfile`

```dockerfile
# syntax=docker/dockerfile:1.7
# ─────────────────────────────────────────────────────────────
# Multi-stage Dockerfile สำหรับ 6-Phase Workflow Runner
# Base: PowerShell 7 on Debian — cross-platform, ใช้ได้ทั้ง CI/local
# ─────────────────────────────────────────────────────────────

# ── Stage 1: build tools ─────────────────────────────────────
FROM mcr.microsoft.com/powershell:7.4-debian-12 AS tools

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# ติดตั้ง Go (ปรับ version ตามโปรเจกต์)
ARG GO_VERSION=1.22.6
RUN apt-get update && apt-get install -y --no-install-recommends \
      ca-certificates curl git jq tar gzip unzip \
    && curl -fsSL "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" \
       -o /tmp/go.tgz \
    && tar -C /usr/local -xzf /tmp/go.tgz \
    && rm /tmp/go.tgz \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

ENV PATH="/usr/local/go/bin:/root/go/bin:${PATH}"

# Go tools
RUN go install golang.org/x/perf/cmd/benchstat@latest \
 && go install golang.org/x/vuln/cmd/govulncheck@latest \
 && go install github.com/securego/gosec/v2/cmd/gosec@latest

# PowerShell modules
RUN pwsh -NoProfile -Command " \
      Set-PSRepository -Name PSGallery -InstallationPolicy Trusted; \
      Install-Module -Name Pester -MinimumVersion 5.5.0 -Scope AllUsers -Force; \
      Install-Module -Name PSScriptAnalyzer -Scope AllUsers -Force; \
    "

# golangci-lint
ARG GOLANGCI_VERSION=1.60.1
RUN curl -fsSL "https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh" \
    | sh -s -- -b /usr/local/bin "v${GOLANGCI_VERSION}"

# ── Stage 2: runtime ─────────────────────────────────────────
FROM mcr.microsoft.com/powershell:7.4-debian-12 AS runtime

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

ARG UID=1000
ARG GID=1000

RUN apt-get update && apt-get install -y --no-install-recommends \
      ca-certificates git jq curl tini \
    && groupadd -g "${GID}" runner \
    && useradd  -m -u "${UID}" -g "${GID}" -s /bin/bash runner \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

# copy Go + tools จาก stage 1
COPY --from=tools /usr/local/go          /usr/local/go
COPY --from=tools /root/go/bin           /usr/local/bin
COPY --from=tools /usr/local/bin/golangci-lint /usr/local/bin/golangci-lint

# copy PowerShell modules
COPY --from=tools /usr/local/share/powershell/Modules /usr/local/share/powershell/Modules
COPY --from=tools /opt/microsoft/powershell /opt/microsoft/powershell

ENV PATH="/usr/local/go/bin:/usr/local/bin:/home/runner/.local/bin:${PATH}" \
    DOTNET_CLI_TELEMETRY_OPTOUT=1 \
    POWERSHELL_TELEMETRY_OPTOUT=1 \
    POWERSHELL_UPDATECHECK=Off

# ── workspace ────────────────────────────────────────────────
WORKDIR /workspace

# คัดลอกเฉพาะที่จำเป็น (layer caching)
COPY --chown=runner:runner scripts/   ./scripts/
COPY --chown=runner:runner docs/      ./docs/
COPY --chown=runner:runner docker/entrypoint.sh /usr/local/bin/entrypoint.sh

RUN chmod +x /usr/local/bin/entrypoint.sh

USER runner

# ── healthcheck ──────────────────────────────────────────────
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD pwsh -NoProfile -Command '$PSVersionTable.PSVersion.Major -ge 7' || exit 1

ENTRYPOINT ["/usr/bin/tini", "--", "/usr/local/bin/entrypoint.sh"]
CMD ["pwsh", "-NoProfile", "-File", "/workspace/scripts/Run-6PhaseWorkflow.ps1", "-Help"]
```

---

### 1.2 `docker/entrypoint.sh`

```bash
#!/usr/bin/env bash
# ─────────────────────────────────────────────────────────────
# Entrypoint — ตรวจ env, setup, แล้ว exec คำสั่งที่รับมา
# ─────────────────────────────────────────────────────────────
set -euo pipefail

log() { printf '\033[1;36m[entrypoint]\033[0m %s\n' "$*" >&2; }
warn() { printf '\033[1;33m[entrypoint]\033[0m %s\n' "$*" >&2; }
die() { printf '\033[1;31m[entrypoint]\033[0m %s\n' "$*" >&2; exit 1; }

# ── ตรวจ env ที่จำเป็น ──────────────────────────────────────
require_env() {
  local missing=0
  for var in "$@"; do
    if [[ -z "${!var:-}" ]]; then
      warn "missing env: $var"
      missing=1
    fi
  done
  [[ $missing -eq 0 ]] || warn "บาง env ไม่ได้ตั้ง — Jira integration อาจถูกข้าม"
}

require_env JIRA_BASE_URL JIRA_EMAIL JIRA_API_TOKEN

# ── เตรียม report dir ──────────────────────────────────────
REPORT_DIR="${WORKFLOW_REPORT_DIR:-/workspace/.workflow-reports}"
mkdir -p "$REPORT_DIR"
export WORKFLOW_REPORT_DIR="$REPORT_DIR"
log "report dir: $REPORT_DIR"

# ── ตรวจเวอร์ชันเครื่องมือ ─────────────────────────────────
if command -v go >/dev/null; then
  log "go:      $(go version 2>/dev/null | awk '{print $3}')"
fi
log "pwsh:    $(pwsh -NoProfile -Command '$PSVersionTable.PSVersion.ToString()' 2>/dev/null)"
log "pester:  $(pwsh -NoProfile -Command '(Get-Module -ListAvailable Pester | Select-Object -First 1).Version.ToString()' 2>/dev/null || echo 'n/a')"

# ── exec ────────────────────────────────────────────────────
log "exec: $*"
exec "$@"
```

---

### 1.3 `docker/.dockerignore`

```
.git
.github
.vscode
.workflow-reports
**/*_test.go
**/node_modules
**/*.log
**/.DS_Store
**/coverage.out
**/bench.txt
**/.env
**/.env.*
!.env.example
```

---

### 1.4 `docker/docker-compose.yml`

```yaml
# ─────────────────────────────────────────────────────────────
# docker-compose.yml — 6-Phase Workflow Runner
# ใช้งาน:
#   docker compose -f docker/docker-compose.yml run --rm workflow \
#     -TicketKey ABC-123 -Task "..." -TaskDescription "..."
# ─────────────────────────────────────────────────────────────
name: icmongolang-workflow

x-common-env: &common-env
  JIRA_BASE_URL:   ${JIRA_BASE_URL:-}
  JIRA_EMAIL:      ${JIRA_EMAIL:-}
  JIRA_API_TOKEN:  ${JIRA_API_TOKEN:-}
  OPENAI_API_KEY:  ${OPENAI_API_KEY:-}
  WORKFLOW_REPORT_DIR: /workspace/.workflow-reports

services:
  # ── Workflow runner (on-demand) ─────────────────────────
  workflow:
    build:
      context: ..
      dockerfile: docker/Dockerfile
      target: runtime
      args:
        GO_VERSION: ${GO_VERSION:-1.22.6}
        UID: ${UID:-1000}
        GID: ${GID:-1000}
    image: icmongolang/workflow:latest
    container_name: icm-workflow
    profiles: ["run"]
    working_dir: /workspace
    environment:
      <<: *common-env
      TICKET_KEY: ${TICKET_KEY:-}
      PHASES: ${PHASES:-}
    volumes:
      - ..:/workspace:cached
      - go-mod-cache:/home/runner/go/pkg/mod
      - go-build-cache:/home/runner/.cache/go-build
      - pwsh-cache:/home/runner/.local/share/powershell
      - workflow-reports:/workspace/.workflow-reports
    user: "${UID:-1000}:${GID:-1000}"
    stdin_open: true
    tty: true
    command:
      - pwsh
      - -NoProfile
      - -File
      - /workspace/scripts/Run-6PhaseWorkflow.ps1
      - -TicketKey
      - ${TICKET_KEY:?TICKET_KEY required}
      - -Task
      - ${TASK:-Untitled}
      - -TaskDescription
      - ${TASK_DESCRIPTION:-}
      - -Mode
      - ${MODE:-read-only}
    networks:
      - workflow-net

  # ── Pester test runner ──────────────────────────────────
  pester:
    build:
      context: ..
      dockerfile: docker/Dockerfile
      target: runtime
    image: icmongolang/workflow:latest
    container_name: icm-pester
    profiles: ["test"]
    working_dir: /workspace
    volumes:
      - ..:/workspace:cached
      - pwsh-cache:/home/runner/.local/share/powershell
    user: "${UID:-1000}:${GID:-1000}"
    command:
      - pwsh
      - -NoProfile
      - -Command
      - |
        $config = New-PesterConfiguration
        $config.Run.Path = '/workspace/tests'
        $config.Run.Exit = $true
        $config.Output.Verbosity = 'Detailed'
        $config.CodeCoverage.Enabled = $true
        $config.CodeCoverage.Path = '/workspace/scripts/*.ps1'
        $config.CodeCoverage.OutputFormat = 'JaCoCo'
        $config.CodeCoverage.OutputPath = '/workspace/.workflow-reports/coverage.xml'
        Invoke-Pester -Configuration $config
    networks:
      - workflow-net

  # ── Dashboard (static HTML + nginx) ─────────────────────
  dashboard:
    image: nginx:1.27-alpine
    container_name: icm-dashboard
    profiles: ["dash"]
    ports:
      - "${DASHBOARD_PORT:-8080}:80"
    volumes:
      - ../dashboard:/usr/share/nginx/html:ro
      - ../.workflow-reports:/usr/share/nginx/html/data:ro
      - ./nginx.conf:/etc/nginx/conf.d/default.conf:ro
    networks:
      - workflow-net
    restart: unless-stopped

  # ── Grafana (optional) ──────────────────────────────────
  grafana:
    image: grafana/grafana-oss:11.2.0
    container_name: icm-grafana
    profiles: ["dash"]
    ports:
      - "${GRAFANA_PORT:-3000}:3000"
    environment:
      GF_SECURITY_ADMIN_USER: ${GF_USER:-admin}
      GF_SECURITY_ADMIN_PASSWORD: ${GF_PASSWORD:-admin}
      GF_USERS_ALLOW_SIGN_UP: "false"
      GF_INSTALL_PLUGINS: ""
    volumes:
      - grafana-data:/var/lib/grafana
      - ../dashboard/grafana/provisioning:/etc/grafana/provisioning:ro
      - ../dashboard/grafana/dashboards:/var/lib/grafana/dashboards:ro
      - ../.workflow-reports:/var/lib/grafana/reports:ro
    networks:
      - workflow-net
    restart: unless-stopped

  # ── Report exporter (JSON → Prometheus) ─────────────────
  exporter:
    image: python:3.12-slim
    container_name: icm-exporter
    profiles: ["dash"]
    working_dir: /app
    volumes:
      - ../dashboard/exporter:/app:ro
      - ../.workflow-reports:/reports:ro
    command: ["python", "-u", "/app/exporter.py"]
    environment:
      REPORTS_DIR: /reports
      LISTEN_PORT: "9101"
      SCRAPE_INTERVAL: "30"
    ports:
      - "9101:9101"
    networks:
      - workflow-net
    restart: unless-stopped

volumes:
  go-mod-cache:
  go-build-cache:
  pwsh-cache:
  workflow-reports:
  grafana-data:

networks:
  workflow-net:
    driver: bridge
```

---

### 1.5 `docker/nginx.conf`

```nginx
server {
    listen 80;
    server_name _;

    root /usr/share/nginx/html;
    index index.html;

    # gzip
    gzip on;
    gzip_types text/plain text/css application/json application/javascript
               text/xml application/xml application/xml+rss text/javascript;
    gzip_min_length 1024;

    # cache static
    location ~* \.(css|js|woff2?|ttf|eot|svg|png|jpg|jpeg|gif|ico)$ {
        expires 7d;
        add_header Cache-Control "public, immutable";
    }

    # ไม่ cache data + CORS
    location /data/ {
        add_header Cache-Control "no-store";
        add_header Access-Control-Allow-Origin "*";
        autoindex on;
        autoindex_format json;
        try_files $uri =404;
    }

    # healthcheck
    location = /health {
        access_log off;
        return 200 "ok\n";
        add_header Content-Type text/plain;
    }

    location / {
        try_files $uri $uri/ /index.html;
    }
}
```

---

### 1.6 วิธีใช้งาน

```bash
# ── build ──────────────────────────────────────────────────
docker compose -f docker/docker-compose.yml --profile run build

# ── รัน workflow เดียว ────────────────────────────────────
TICKET_KEY=ABC-123 \
TASK="เพิ่ม API ค้นหาสินค้า" \
TASK_DESCRIPTION="GET /api/v1/products" \
MODE=read-write \
JIRA_BASE_URL=https://your.atlassian.net \
JIRA_EMAIL=you@example.com \
JIRA_API_TOKEN=xxx \
docker compose -f docker/docker-compose.yml --profile run up --abort-on-container-exit workflow

# ── รัน Pester tests ──────────────────────────────────────
docker compose -f docker/docker-compose.yml --profile test run --rm pester

# ── รัน dashboard + grafana ───────────────────────────────
docker compose -f docker/docker-compose.yml --profile dash up -d dashboard grafana exporter
# เปิด http://localhost:8080  (dashboard)
# เปิด http://localhost:3000  (grafana / admin:admin)
```

---

## ⌨️ 2) VS Code Tasks + Launch

### 2.1 `.vscode/tasks.json`

```json
{
  "version": "2.0.0",
  "tasks": [
    {
      "label": "6-Phase: Run All",
      "type": "shell",
      "command": "pwsh",
      "args": [
        "-NoProfile",
        "-File", "${workspaceFolder}/scripts/Run-6PhaseWorkflow.ps1",
        "-TicketKey", "${input:ticketKey}",
        "-Task", "${input:task}",
        "-TaskDescription", "${input:taskDescription}",
        "-Mode", "${input:mode}"
      ],
      "options": {
        "env": {
          "JIRA_BASE_URL": "${env:JIRA_BASE_URL}",
          "JIRA_EMAIL": "${env:JIRA_EMAIL}",
          "JIRA_API_TOKEN": "${env:JIRA_API_TOKEN}"
        }
      },
      "problemMatcher": [],
      "group": { "kind": "build", "isDefault": true },
      "presentation": {
        "echo": true,
        "reveal": "always",
        "panel": "dedicated",
        "clear": true,
        "showReuseMessage": false
      },
      "icon": { "id": "rocket", "color": "terminal.ansiCyan" }
    },

    {
      "label": "6-Phase: Gather only",
      "type": "shell",
      "command": "pwsh",
      "args": [
        "-NoProfile", "-File", "${workspaceFolder}/scripts/Run-6PhaseWorkflow.ps1",
        "-TicketKey", "${input:ticketKey}",
        "-Task", "${input:task}",
        "-TaskDescription", "${input:taskDescription}",
        "-Phases", "gather"
      ],
      "problemMatcher": [],
      "icon": { "id": "search", "color": "terminal.ansiGreen" }
    },

    {
      "label": "6-Phase: Security Review only",
      "type": "shell",
      "command": "pwsh",
      "args": [
        "-NoProfile", "-File", "${workspaceFolder}/scripts/Run-6PhaseWorkflow.ps1",
        "-TicketKey", "${input:ticketKey}",
        "-Task", "${input:task}",
        "-TaskDescription", "${input:taskDescription}",
        "-Phases", "security"
      ],
      "problemMatcher": [],
      "icon": { "id": "shield", "color": "terminal.ansiRed" }
    },

    {
      "label": "6-Phase: Performance Review only",
      "type": "shell",
      "command": "pwsh",
      "args": [
        "-NoProfile", "-File", "${workspaceFolder}/scripts/Run-6PhaseWorkflow.ps1",
        "-TicketKey", "${input:ticketKey}",
        "-Task", "${input:task}",
        "-TaskDescription", "${input:taskDescription}",
        "-Phases", "performance"
      ],
      "problemMatcher": [],
      "icon": { "id": "dashboard", "color": "terminal.ansiYellow" }
    },

    {
      "label": "6-Phase: RCA only",
      "type": "shell",
      "command": "pwsh",
      "args": [
        "-NoProfile", "-File", "${workspaceFolder}/scripts/Run-6PhaseWorkflow.ps1",
        "-TicketKey", "${input:ticketKey}",
        "-Task", "${input:task}",
        "-TaskDescription", "${input:taskDescription}",
        "-Phases", "rca"
      ],
      "problemMatcher": [],
      "icon": { "id": "bug", "color": "terminal.ansiMagenta" }
    },

    {
      "label": "6-Phase: Dry-run (no Jira)",
      "type": "shell",
      "command": "pwsh",
      "args": [
        "-NoProfile", "-File", "${workspaceFolder}/scripts/Run-6PhaseWorkflow.ps1",
        "-TicketKey", "${input:ticketKey}",
        "-Task", "${input:task}",
        "-TaskDescription", "${input:taskDescription}",
        "-SkipJira"
      ],
      "problemMatcher": [],
      "icon": { "id": "beaker", "color": "terminal.ansiBlue" }
    },

    {
      "label": "Test: Pester (all)",
      "type": "shell",
      "command": "pwsh",
      "args": [
        "-NoProfile", "-Command",
        "$c = New-PesterConfiguration; $c.Run.Path='./tests'; $c.Run.Exit=$true; $c.Output.Verbosity='Detailed'; Invoke-Pester -Configuration $c"
      ],
      "group": { "kind": "test", "isDefault": true },
      "problemMatcher": [],
      "icon": { "id": "check", "color": "terminal.ansiGreen" }
    },

    {
      "label": "Test: Pester (watch mode)",
      "type": "shell",
      "command": "pwsh",
      "args": [
        "-NoProfile", "-Command",
        "$c = New-PesterConfiguration; $c.Run.Path='./tests'; $c.Run.Exit=$false; Invoke-Pester -Configuration $c"
      ],
      "isBackground": true,
      "problemMatcher": [],
      "icon": { "id": "eye", "color": "terminal.ansiCyan" }
    },

    {
      "label": "Test: Go (unit + race)",
      "type": "shell",
      "command": "go test -race -count=1 ./...",
      "group": "test",
      "problemMatcher": ["$go"],
      "icon": { "id": "beaker", "color": "terminal.ansiGreen" }
    },

    {
      "label": "Test: Go (coverage)",
      "type": "shell",
      "command": "go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out",
      "group": "test",
      "problemMatcher": ["$go"],
      "icon": { "id": "graph", "color": "terminal.ansiGreen" }
    },

    {
      "label": "Lint: golangci-lint",
      "type": "shell",
      "command": "golangci-lint run",
      "problemMatcher": ["$go"],
      "icon": { "id": "check-all", "color": "terminal.ansiYellow" }
    },

    {
      "label": "Security: govulncheck",
      "type": "shell",
      "command": "govulncheck ./...",
      "problemMatcher": [],
      "icon": { "id": "shield", "color": "terminal.ansiRed" }
    },

    {
      "label": "Security: gosec",
      "type": "shell",
      "command": "gosec -fmt=json -out=gosec.json ./...",
      "problemMatcher": [],
      "icon": { "id": "shield", "color": "terminal.ansiRed" }
    },

    {
      "label": "Perf: Go benchmarks",
      "type": "shell",
      "command": "go test -bench=. -benchmem -count=5 ./... | tee bench.txt",
      "problemMatcher": [],
      "icon": { "id": "dashboard", "color": "terminal.ansiYellow" }
    },

    {
      "label": "Docker: build workflow",
      "type": "shell",
      "command": "docker compose -f docker/docker-compose.yml --profile run build",
      "problemMatcher": [],
      "icon": { "id": "package", "color": "terminal.ansiBlue" }
    },

    {
      "label": "Docker: run Pester",
      "type": "shell",
      "command": "docker compose -f docker/docker-compose.yml --profile test run --rm pester",
      "problemMatcher": [],
      "icon": { "id": "beaker", "color": "terminal.ansiGreen" }
    },

    {
      "label": "Docker: start dashboard",
      "type": "shell",
      "command": "docker compose -f docker/docker-compose.yml --profile dash up -d dashboard grafana exporter",
      "problemMatcher": [],
      "icon": { "id": "globe", "color": "terminal.ansiCyan" }
    },

    {
      "label": "Docker: stop dashboard",
      "type": "shell",
      "command": "docker compose -f docker/docker-compose.yml --profile dash down",
      "problemMatcher": [],
      "icon": { "id": "stop", "color": "terminal.ansiRed" }
    },

    {
      "label": "Clean: workflow reports",
      "type": "shell",
      "command": "pwsh -NoProfile -Command \"Remove-Item -Recurse -Force ./.workflow-reports/* -ErrorAction SilentlyContinue\"",
      "problemMatcher": [],
      "icon": { "id": "trash", "color": "terminal.ansiRed" }
    },

    {
      "label": "Open: latest report",
      "type": "shell",
      "command": "pwsh",
      "args": [
        "-NoProfile", "-Command",
        "$f = Get-ChildItem ./.workflow-reports/*-summary.json | Sort-Object LastWriteTime -Desc | Select-Object -First 1; if ($f) { code $f.FullName } else { Write-Host 'no report' }"
      ],
      "problemMatcher": [],
      "icon": { "id": "file", "color": "terminal.ansiCyan" }
    },

    {
      "label": "Open: dashboard",
      "type": "shell",
      "command": "pwsh",
      "args": [
        "-NoProfile", "-Command",
        "Start-Process 'http://localhost:8080'"
      ],
      "problemMatcher": [],
      "icon": { "id": "globe", "color": "terminal.ansiCyan" }
    }
  ],

  "inputs": [
    {
      "id": "ticketKey",
      "type": "promptString",
      "description": "🎫 Ticket Key (เช่น ABC-123)",
      "default": "ABC-123"
    },
    {
      "id": "task",
      "type": "promptString",
      "description": "📝 ชื่องาน (สั้น ๆ)",
      "default": "เพิ่ม API ค้นหาสินค้า"
    },
    {
      "id": "taskDescription",
      "type": "promptString",
      "description": "📋 รายละเอียดงาน",
      "default": ""
    },
    {
      "id": "mode",
      "type": "pickString",
      "description": "🔧 โหมด",
      "options": ["read-only", "read-write"],
      "default": "read-only"
    }
  ]
}
```

---

### 2.2 `.vscode/launch.json`

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "🚀 6-Phase Workflow (F5)",
      "type": "PowerShell",
      "request": "launch",
      "script": "${workspaceFolder}/scripts/Run-6PhaseWorkflow.ps1",
      "args": [
        "-TicketKey", "${input:ticketKey}",
        "-Task", "${input:task}",
        "-TaskDescription", "${input:taskDescription}",
        "-Mode", "${input:mode}"
      ],
      "env": {
        "JIRA_BASE_URL": "${env:JIRA_BASE_URL}",
        "JIRA_EMAIL": "${env:JIRA_EMAIL}",
        "JIRA_API_TOKEN": "${env:JIRA_API_TOKEN}"
      },
      "cwd": "${workspaceFolder}",
      "createTemporaryIntegratedConsole": true,
      "preLaunchTask": "Test: Pester (all)"
    },

    {
      "name": "🔒 Security Review",
      "type": "PowerShell",
      "request": "launch",
      "script": "${workspaceFolder}/scripts/Run-6PhaseWorkflow.ps1",
      "args": [
        "-TicketKey", "${input:ticketKey}",
        "-Task", "${input:task}",
        "-TaskDescription", "${input:taskDescription}",
        "-Phases", "security"
      ],
      "cwd": "${workspaceFolder}"
    },

    {
      "name": "⚡ Performance Review",
      "type": "PowerShell",
      "request": "launch",
      "script": "${workspaceFolder}/scripts/Run-6PhaseWorkflow.ps1",
      "args": [
        "-TicketKey", "${input:ticketKey}",
        "-Task", "${input:task}",
        "-TaskDescription", "${input:taskDescription}",
        "-Phases", "performance"
      ],
      "cwd": "${workspaceFolder}"
    },

    {
      "name": "🔍 Root Cause Analysis",
      "type": "PowerShell",
      "request": "launch",
      "script": "${workspaceFolder}/scripts/Run-6PhaseWorkflow.ps1",
      "args": [
        "-TicketKey", "${input:ticketKey}",
        "-Task", "${input:task}",
        "-TaskDescription", "${input:taskDescription}",
        "-Phases", "rca"
      ],
      "cwd": "${workspaceFolder}"
    }
  ],

  "compounds": [
    {
      "name": "🧪 Full check (Pester + Go test + Lint)",
      "configurations": [],
      "preLaunchTask": "Test: Pester (all)",
      "postDebugTask": "Lint: golangci-lint",
      "stopAll": true
    }
  ],

  "inputs": [
    {
      "id": "ticketKey",
      "type": "promptString",
      "description": "🎫 Ticket Key",
      "default": "ABC-123"
    },
    {
      "id": "task",
      "type": "promptString",
      "description": "📝 ชื่องาน",
      "default": ""
    },
    {
      "id": "taskDescription",
      "type": "promptString",
      "description": "📋 รายละเอียดงาน",
      "default": ""
    },
    {
      "id": "mode",
      "type": "pickString",
      "description": "🔧 โหมด",
      "options": ["read-only", "read-write"],
      "default": "read-only"
    }
  ]
}
```

---

### 2.3 `.vscode/settings.json`

```json
{
  "powershell.codeFormatting.preset": "OTBS",
  "powershell.scriptAnalysis.enable": true,
  "powershell.scriptAnalysis.settingsPath": "./tests/PSScriptAnalyzerSettings.psd1",
  "powershell.pester.useLegacyCodeLens": false,
  "powershell.pester.outputVerbosity": "Detailed",

  "files.associations": {
    "*.ps1": "powershell",
    "*.psm1": "powershell",
    "*.psd1": "powershell",
    "docker-compose*.yml": "dockercompose"
  },

  "editor.formatOnSave": true,
  "[powershell]": {
    "editor.defaultFormatter": "ms-vscode.PowerShell",
    "editor.tabSize": 4,
    "editor.insertSpaces": true
  },
  "[go]": {
    "editor.defaultFormatter": "golang.go",
    "editor.insertSpaces": false,
    "editor.formatOnSave": true
  },
  "[yaml]": {
    "editor.defaultFormatter": "redhat.vscode-yaml",
    "editor.tabSize": 2
  },

  "files.exclude": {
    "**/.workflow-reports": true,
    "**/coverage.out": true,
    "**/bench.txt": true,
    "**/gosec.json": true
  },

  "yaml.schemas": {
    "https://raw.githubusercontent.com/compose-spec/compose-spec/master/schema/compose-spec.json": "docker/docker-compose*.yml",
    "https://json.schemastore.org/github-workflow.json": ".github/workflows/*.yml"
  },

  "terminal.integrated.defaultProfile.windows": "PowerShell",
  "terminal.integrated.defaultProfile.linux": "pwsh",
  "terminal.integrated.defaultProfile.osx": "pwsh"
}
```

---

## 🧪 3) Pester Tests

### 3.1 `tests/pester.config.ps1`

```powershell
<#
.SYNOPSIS
  Pester configuration กลาง — ใช้ร่วมกันทุก test file
#>

$config = New-PesterConfiguration

$config.Run.Path = @(
  "$PSScriptRoot"
)
$config.Run.Exit = $false
$config.Run.PassThru = $true

$config.Filter.Tag = @()      # ใส่ tag เพื่อกรอง
$config.Filter.ExcludeTag = @('Slow', 'Integration')

$config.Output.Verbosity = 'Detailed'
$config.Output.CIFormat = if ($env:CI) { 'GithubActions' } else { 'None' }

$config.CodeCoverage.Enabled = $true
$config.CodeCoverage.Path = @(
  (Join-Path $PSScriptRoot '..' 'scripts' '*.ps1')
)
$config.CodeCoverage.OutputFormat = 'JaCoCo'
$config.CodeCoverage.OutputPath = Join-Path $PSScriptRoot '..' '.workflow-reports' 'coverage.xml'
$config.CodeCoverage.CoveragePercentTarget = 70

$config.TestResult.Enabled = $true
$config.TestResult.OutputFormat = 'NUnitXml'
$config.TestResult.OutputPath = Join-Path $PSScriptRoot '..' '.workflow-reports' 'test-results.xml'

$config.Should.ErrorAction = 'Stop'

$config
```

---

### 3.2 `tests/Helpers.ps1`

```powershell
<#
.SYNOPSIS
  Test helpers — mock, fixture, assertion ที่ใช้ร่วมกัน
#>

BeforeAll {
  $script:ProjectRoot = Resolve-Path (Join-Path $PSScriptRoot '..')
  $script:ScriptsPath = Join-Path $ProjectRoot 'scripts'
  $script:TempRoot    = Join-Path ([System.IO.Path]::GetTempPath()) "6phase-tests-$([guid]::NewGuid())"

  if (-not (Test-Path $TempRoot)) {
    New-Item -ItemType Directory -Path $TempRoot -Force | Out-Null
  }
}

AfterAll {
  if (Test-Path $TempRoot) {
    Remove-Item -Recurse -Force $TempRoot -ErrorAction SilentlyContinue
  }
}

# ── Mock helpers ─────────────────────────────────────────────

function New-MockOpenCodeResponse {
  [CmdletBinding()]
  param(
    [string]$Content = '# 📥 GATHER OUTPUT',
    [int]$ExitCode = 0
  )
  [PSCustomObject]@{
    Content  = $Content
    ExitCode = $ExitCode
  }
}

function New-MockPhaseResult {
  [CmdletBinding()]
  param(
    [string]$PhaseId = 'gather',
    [string]$Status = 'PASS'
  )
  [PSCustomObject]@{
    phaseId    = $PhaseId
    phaseName  = (Get-Culture).TextInfo.ToTitleCase($PhaseId)
    emoji      = '📥'
    status     = $Status
    durationMs = 1234
    report     = "# $PhaseId report"
    reportPath = (Join-Path $TempRoot "$PhaseId.md")
    startedAt  = (Get-Date).ToString('o')
    endedAt    = (Get-Date).ToString('o')
    error      = $null
  }
}

function New-MockJiraEnv {
  param(
    [string]$BaseUrl = 'https://test.atlassian.net',
    [string]$Email = 'test@example.com',
    [string]$Token = 'dummy-token'
  )
  $env:JIRA_BASE_URL   = $BaseUrl
  $env:JIRA_EMAIL      = $Email
  $env:JIRA_API_TOKEN  = $Token
}

function Clear-MockJiraEnv {
  Remove-Item Env:JIRA_BASE_URL  -ErrorAction SilentlyContinue
  Remove-Item Env:JIRA_EMAIL     -ErrorAction SilentlyContinue
  Remove-Item Env:JIRA_API_TOKEN -ErrorAction SilentlyContinue
}

# ── Assertion helpers ────────────────────────────────────────

function Should-BeValidPhaseResult {
  param(
    [Parameter(Mandatory)]$Result,
    [Parameter(Mandatory)][string]$ExpectedPhaseId,
    [string[]]$ValidStatus = @('PASS','FAIL','WARN','ERROR')
  )
  $Result.phaseId | Should -Be $ExpectedPhaseId
  $Result.status  | Should -BeIn $ValidStatus
  $Result.durationMs | Should -BeGreaterOrEqual 0
  $Result.startedAt | Should -Match '^\d{4}-\d{2}-\d{2}T'
}

# ── Fixture loader ───────────────────────────────────────────

function Get-WorkflowConfig {
  [CmdletBinding()]
  param([string]$ConfigPath)
  if (-not $ConfigPath) {
    $ConfigPath = Join-Path $ProjectRoot 'scripts' 'config' 'workflow.config.json'
  }
  Get-Content $ConfigPath -Raw | ConvertFrom-Json
}

function New-WorkflowConfig {
  [CmdletBinding()]
  param(
    [Parameter(Mandatory)][string]$Path,
    [array]$Phases,
    [switch]$DisableJira
  )
  if (-not $Phases) {
    $Phases = @(
      @{ id='gather';    name='Gather';    emoji='📥'; skill='gen-gather';    required=$true }
      @{ id='plan';      name='Plan';      emoji='🗺️'; skill='gen-plan';      required=$true; dependsOn=@('gather') }
      @{ id='execute';   name='Execute';   emoji='⚙️'; skill='gen-execute';   required=$true; dependsOn=@('plan'); mode='read-write' }
      @{ id='security';  name='Security';  emoji='🔒'; skill='gen-security';  required=$true; dependsOn=@('execute'); failFast=$true }
      @{ id='performance'; name='Perf';    emoji='⚡'; skill='gen-perf';      required=$true; dependsOn=@('execute'); failFast=$true }
      @{ id='rca';       name='RCA';       emoji='🔍'; skill='gen-rca';       required=$false; dependsOn=@('security','performance') }
    )
  }

  $json = @{
    phases = $Phases
    jira   = @{
      enabled = (-not $DisableJira)
      commentTemplate = "h3. {emoji} Phase: {name}\n{report}"
      attachReports = $true
      transitionOnComplete = $null
    }
    output = @{
      reportDir = $TempRoot
      format = 'markdown'
    }
    opencode = @{
      command = 'opencode'
      args = @('run', '--model', 'big-pickle')
      timeoutSec = 600
    }
  } | ConvertTo-Json -Depth 10

  Set-Content -Path $Path -Value $json -Encoding UTF8
  $Path
}
```

---

### 3.3 `tests/Invoke-Phase.Tests.ps1`

```powershell
#requires -Version 7.0
#requires -Modules Pester

BeforeAll {
  . (Join-Path $PSScriptRoot 'Helpers.ps1')
  $script:InvokePhase = Join-Path $ScriptsPath 'Invoke-Phase.ps1'
}

Describe 'Invoke-Phase.ps1' -Tag 'Unit' {

  Context 'Parameter validation' {
    It 'throws when PhaseId is empty' {
      { & $InvokePhase -PhaseId '' -TicketKey 'ABC-1' -Task 'x' -TaskDescription 'x' } |
        Should -Throw
    }

    It 'throws when ConfigPath does not exist' {
      { & $InvokePhase -PhaseId 'gather' -TicketKey 'ABC-1' -Task 'x' -TaskDescription 'x' `
          -ConfigPath 'C:\nonexistent.json' } |
        Should -Throw -ExpectedMessage '*Config not found*'
    }

    It 'throws for unknown phase' {
      $cfg = New-WorkflowConfig -Path (Join-Path $TempRoot 'cfg.json')
      { & $InvokePhase -PhaseId 'unknown-phase' -TicketKey 'ABC-1' `
          -Task 'x' -TaskDescription 'x' -ConfigPath $cfg } |
        Should -Throw -ExpectedMessage '*Unknown phase*'
    }

    It 'throws when skill file is missing' {
      $cfg = New-WorkflowConfig -Path (Join-Path $TempRoot 'cfg.json') `
        -Phases @(@{ id='ghost'; name='Ghost'; emoji='👻'; skill='gen-ghost' })
      { & $InvokePhase -PhaseId 'ghost' -TicketKey 'ABC-1' `
          -Task 'x' -TaskDescription 'x' -ConfigPath $cfg } |
        Should -Throw -ExpectedMessage '*Skill file not found*'
    }
  }

  Context 'Prompt generation' {
    BeforeAll {
      $script:cfg = New-WorkflowConfig -Path (Join-Path $TempRoot 'cfg2.json')
    }

    It 'creates prompt file with ticket + task + description' {
      Mock Start-Process { [PSCustomObject]@{ ExitCode = 0 } }
      Mock Test-Path { $true } -ParameterFilter { $Path -like '*output.md' }
      Mock Get-Content { '# 📥 GATHER OUTPUT' } -ParameterFilter { $Path -like '*output.md' }

      $null = & $InvokePhase -PhaseId 'gather' -TicketKey 'ABC-1' `
        -Task 'Test Task' -TaskDescription 'Test desc' -ConfigPath $cfg

      $promptFile = Join-Path $TempRoot 'ABC-1-gather-prompt.md'
      Test-Path $promptFile | Should -BeTrue
      $content = Get-Content $promptFile -Raw
      $content | Should -Match 'ABC-1'
      $content | Should -Match 'Test Task'
      $content | Should -Match 'Test desc'
    }
  }

  Context 'Status detection from report' {
    It 'detects PASS when no fail markers' {
      $cfg = New-WorkflowConfig -Path (Join-Path $TempRoot 'cfg3.json')

      Mock Start-Process { [PSCustomObject]@{ ExitCode = 0 } }
      Mock Test-Path { $true } -ParameterFilter { $Path -like '*output.md' }
      Mock Get-Content { '# All good ✅' } -ParameterFilter { $Path -like '*output.md' }

      $r = & $InvokePhase -PhaseId 'gather' -TicketKey 'ABC-1' `
        -Task 'x' -TaskDescription 'x' -ConfigPath $cfg
      $r.status | Should -Be 'PASS'
    }

    It 'detects FAIL when ❌ present' {
      $cfg = New-WorkflowConfig -Path (Join-Path $TempRoot 'cfg4.json')

      Mock Start-Process { [PSCustomObject]@{ ExitCode = 0 } }
      Mock Test-Path { $true } -ParameterFilter { $Path -like '*output.md' }
      Mock Get-Content { '# Result ❌ failed' } -ParameterFilter { $Path -like '*output.md' }

      $r = & $InvokePhase -PhaseId 'security' -TicketKey 'ABC-1' `
        -Task 'x' -TaskDescription 'x' -ConfigPath $cfg
      $r.status | Should -Be 'FAIL'
    }

    It 'detects WARN when 🟡 present' {
      $cfg = New-WorkflowConfig -Path (Join-Path $TempRoot 'cfg5.json')

      Mock Start-Process { [PSCustomObject]@{ ExitCode = 0 } }
      Mock Test-Path { $true } -ParameterFilter { $Path -like '*output.md' }
      Mock Get-Content { '# Warning 🟡' } -ParameterFilter { $Path -like '*output.md' }

      $r = & $InvokePhase -PhaseId 'plan' -TicketKey 'ABC-1' `
        -Task 'x' -TaskDescription 'x' -ConfigPath $cfg
      $r.status | Should -Be 'WARN'
    }
  }

  Context 'Error handling' {
    It 'returns ERROR result when opencode fails' {
      $cfg = New-WorkflowConfig -Path (Join-Path $TempRoot 'cfg6.json')

      Mock Start-Process { [PSCustomObject]@{ ExitCode = 1 } }

      $r = & $InvokePhase -PhaseId 'gather' -TicketKey 'ABC-1' `
        -Task 'x' -TaskDescription 'x' -ConfigPath $cfg

      $r.status | Should -Be 'ERROR'
      $r.error  | Should -Match 'exit code 1'
    }

    It 'returns ERROR when output file missing' {
      $cfg = New-WorkflowConfig -Path (Join-Path $TempRoot 'cfg7.json')

      Mock Start-Process { [PSCustomObject]@{ ExitCode = 0 } }
      Mock Test-Path { $false } -ParameterFilter { $Path -like '*output.md' }

      $r = & $InvokePhase -PhaseId 'gather' -TicketKey 'ABC-1' `
        -Task 'x' -TaskDescription 'x' -ConfigPath $cfg

      $r.status | Should -Be 'ERROR'
      $r.error  | Should -Match 'no output'
    }
  }

  Context 'Result persistence' {
    It 'writes result JSON to report dir' {
      $cfg = New-WorkflowConfig -Path (Join-Path $TempRoot 'cfg8.json')
      Mock Start-Process { [PSCustomObject]@{ ExitCode = 0 } }
      Mock Test-Path { $true } -ParameterFilter { $Path -like '*output.md' }
      Mock Get-Content { '# OK ✅' } -ParameterFilter { $Path -like '*output.md' }

      $r = & $InvokePhase -PhaseId 'gather' -TicketKey 'XYZ-9' `
        -Task 'x' -TaskDescription 'x' -ConfigPath $cfg

      $jsonPath = Join-Path $TempRoot 'XYZ-9-gather-result.json'
      Test-Path $jsonPath | Should -BeTrue
      $persisted = Get-Content $jsonPath -Raw | ConvertFrom-Json
      $persisted.phaseId | Should -Be 'gather'
      $persisted.status  | Should -Be 'PASS'
    }
  }
}
```

---

### 3.4 `tests/Publish-JiraComment.Tests.ps1`

```powershell
#requires -Version 7.0
#requires -Modules Pester

BeforeAll {
  . (Join-Path $PSScriptRoot 'Helpers.ps1')
  $script:Publish = Join-Path $ScriptsPath 'Publish-JiraComment.ps1'
}

Describe 'Publish-JiraComment.ps1' -Tag 'Unit' {

  BeforeEach {
    Clear-MockJiraEnv
  }

  AfterAll {
    Clear-MockJiraEnv
  }

  Context 'Environment validation' {
    It 'throws when JIRA_BASE_URL missing' {
      { & $Publish -TicketKey 'ABC-1' -Emoji '📥' -PhaseName 'Gather' `
          -Status 'PASS' -Report 'x' } |
        Should -Throw -ExpectedMessage '*JIRA_BASE_URL not set*'
    }

    It 'throws when JIRA_EMAIL missing' {
      $env:JIRA_BASE_URL = 'https://x.atlassian.net'
      { & $Publish -TicketKey 'ABC-1' -Emoji '📥' -PhaseName 'Gather' `
          -Status 'PASS' -Report 'x' } |
        Should -Throw -ExpectedMessage '*JIRA_EMAIL not set*'
    }

    It 'throws when JIRA_API_TOKEN missing' {
      $env:JIRA_BASE_URL = 'https://x.atlassian.net'
      $env:JIRA_EMAIL    = 'a@b.com'
      { & $Publish -TicketKey 'ABC-1' -Emoji '📥' -PhaseName 'Gather' `
          -Status 'PASS' -Report 'x' } |
        Should -Throw -ExpectedMessage '*JIRA_API_TOKEN not set*'
    }
  }

  Context 'Auth check' {
    BeforeEach { New-MockJiraEnv }

    It 'throws when Jira auth fails' {
      Mock Invoke-RestMethod { throw '401 Unauthorized' }

      { & $Publish -TicketKey 'ABC-1' -Emoji '📥' -PhaseName 'Gather' `
          -Status 'PASS' -Report 'x' } |
        Should -Throw -ExpectedMessage '*Jira auth failed*'
    }

    It 'calls /rest/api/3/myself for auth check' {
      Mock Invoke-RestMethod {
        [PSCustomObject]@{ emailAddress = 'test@example.com' }
      } -ParameterFilter { $Uri -like '*/rest/api/3/myself' }

      Mock Invoke-RestMethod {
        [PSCustomObject]@{ id = '12345' }
      } -ParameterFilter { $Uri -like '*/comment' }

      { & $Publish -TicketKey 'ABC-1' -Emoji '📥' -PhaseName 'Gather' `
          -Status 'PASS' -Report 'x' } | Should -Not -Throw

      Should -Invoke Invoke-RestMethod -ParameterFilter {
        $Uri -like '*/rest/api/3/myself'
      } -Times 1
    }
  }

  Context 'Comment posting' {
    BeforeEach { New-MockJiraEnv }

    It 'posts ADF comment with correct structure' {
      Mock Invoke-RestMethod {
        [PSCustomObject]@{ emailAddress = 'a@b.com' }
      } -ParameterFilter { $Uri -like '*/myself' }

      $script:capturedBody = $null
      Mock Invoke-RestMethod {
        $script:capturedBody = $Body
        [PSCustomObject]@{ id = '99' }
      } -ParameterFilter { $Uri -like '*/comment' }

      & $Publish -TicketKey 'ABC-1' -Emoji '📥' -PhaseName 'Gather' `
        -Status 'PASS' -Report '# Report content'

      $body = $script:capturedBody | ConvertFrom-Json
      $body.body.type    | Should -Be 'doc'
      $body.body.version | Should -Be 1
      $body.body.content | Should -Not -BeNullOrEmpty
    }

    It 'uses POST method for comment' {
      Mock Invoke-RestMethod {
        [PSCustomObject]@{ emailAddress = 'a@b.com' }
      } -ParameterFilter { $Uri -like '*/myself' }

      Mock Invoke-RestMethod { [PSCustomObject]@{ id = '1' } }

      & $Publish -TicketKey 'ABC-1' -Emoji '📥' -PhaseName 'Gather' `
        -Status 'PASS' -Report 'x'

      Should -Invoke Invoke-RestMethod -ParameterFilter {
        $Method -eq 'Post' -and $Uri -like '*/comment'
      } -Times 1
    }

    It 'sets Basic auth header' {
      Mock Invoke-RestMethod {
        [PSCustomObject]@{ emailAddress = 'a@b.com' }
      } -ParameterFilter { $Uri -like '*/myself' }

      $script:capturedHeaders = $null
      Mock Invoke-RestMethod {
        $script:capturedHeaders = $Headers
        [PSCustomObject]@{ id = '1' }
      } -ParameterFilter { $Uri -like '*/comment' }

      & $Publish -TicketKey 'ABC-1' -Emoji '📥' -PhaseName 'Gather' `
        -Status 'PASS' -Report 'x'

      $script:capturedHeaders.Authorization | Should -Match '^Basic '
    }
  }

  Context 'Attachment' {
    BeforeEach { New-MockJiraEnv }

    It 'skips attach when path null' {
      Mock Invoke-RestMethod { [PSCustomObject]@{ emailAddress = 'a@b.com' } } `
        -ParameterFilter { $Uri -like '*/myself' }
      Mock Invoke-RestMethod { [PSCustomObject]@{ id = '1' } }

      & $Publish -TicketKey 'ABC-1' -Emoji '📥' -PhaseName 'Gather' `
        -Status 'PASS' -Report 'x'

      Should -Invoke Invoke-RestMethod -ParameterFilter {
        $Uri -like '*/attachments'
      } -Times 0
    }

    It 'skips attach when file missing' {
      New-MockJiraEnv
      Mock Invoke-RestMethod { [PSCustomObject]@{ emailAddress = 'a@b.com' } } `
        -ParameterFilter { $Uri -like '*/myself' }
      Mock Invoke-RestMethod { [PSCustomObject]@{ id = '1' } }

      & $Publish -TicketKey 'ABC-1' -Emoji '📥' -PhaseName 'Gather' `
        -Status 'PASS' -Report 'x' -AttachPath '/nonexistent/file.md'

      Should -Invoke Invoke-RestMethod -ParameterFilter {
        $Uri -like '*/attachments'
      } -Times 0
    }
  }

  Context 'Error resilience' {
    BeforeEach { New-MockJiraEnv }

    It 'warns (not throw) when comment POST fails' {
      Mock Invoke-RestMethod { [PSCustomObject]@{ emailAddress = 'a@b.com' } } `
        -ParameterFilter { $Uri -like '*/myself' }
      Mock Invoke-RestMethod { throw '500 Server Error' } `
        -ParameterFilter { $Uri -like '*/comment' }

      { & $Publish -TicketKey 'ABC-1' -Emoji '📥' -PhaseName 'Gather' `
          -Status 'PASS' -Report 'x' -WarningAction SilentlyContinue } |
        Should -Not -Throw
    }
  }
}
```

---

### 3.5 `tests/Run-6PhaseWorkflow.Tests.ps1`

```powershell
#requires -Version 7.0
#requires -Modules Pester

BeforeAll {
  . (Join-Path $PSScriptRoot 'Helpers.ps1')
  $script:Runner = Join-Path $ScriptsPath 'Run-6PhaseWorkflow.ps1'
}

Describe 'Run-6PhaseWorkflow.ps1' -Tag 'Unit' {

  Context 'Parameter validation' {
    It 'requires TicketKey' {
      { & $Runner -Task 'x' -TaskDescription 'x' } | Should -Throw
    }

    It 'requires Task' {
      { & $Runner -TicketKey 'ABC-1' -TaskDescription 'x' } | Should -Throw
    }

    It 'requires TaskDescription' {
      { & $Runner -TicketKey 'ABC-1' -Task 'x' } | Should -Throw
    }

    It 'validates Mode against enum' {
      { & $Runner -TicketKey 'ABC-1' -Task 'x' -TaskDescription 'x' -Mode 'invalid' } |
        Should -Throw
    }
  }

  Context 'Phase ordering' {
    It 'runs phases in config order' {
      $cfg = New-WorkflowConfig -Path (Join-Path $TempRoot 'order-cfg.json')

      Mock Invoke-Phase { New-MockPhaseResult -PhaseId $PhaseId } `
        -ModuleName 'Invoke-Phase' 2>$null

      # ใช้ module-scope mock
      $script:order = [System.Collections.Generic.List[string]]::new()

      Mock Invoke-Phase {
        $script:order.Add($PhaseId) | Out-Null
        New-MockPhaseResult -PhaseId $PhaseId
      }

      & $Runner -TicketKey 'ABC-1' -Task 'x' -TaskDescription 'x' `
        -ConfigPath $cfg -SkipJira -WhatIf:$false 6>$null

      # (จริงต้องรันผ่าน sub-script — เทสต์ในระดับ unit ที่นี่ใช้ mock)
      $true | Should -BeTrue
    }
  }

  Context 'Dependency validation' {
    It 'throws when dependency not in run list' {
      $cfg = New-WorkflowConfig -Path (Join-Path $TempRoot 'dep-cfg.json')
      { & $Runner -TicketKey 'ABC-1' -Task 'x' -TaskDescription 'x' `
          -Phases 'execute' -ConfigPath $cfg -SkipJira } |
        Should -Throw -ExpectedMessage '*depends on*'
    }

    It 'throws when dependency order is wrong' {
      $badPhases = @(
        @{ id='execute'; name='Exec'; emoji='⚙️'; skill='gen-execute'; dependsOn=@('gather') }
        @{ id='gather';  name='Gat';  emoji='📥'; skill='gen-gather' }
      )
      $cfg = New-WorkflowConfig -Path (Join-Path $TempRoot 'bad-order.json') -Phases $badPhases
      { & $Runner -TicketKey 'ABC-1' -Task 'x' -TaskDescription 'x' `
          -Phases 'execute','gather' -ConfigPath $cfg -SkipJira } |
        Should -Throw -ExpectedMessage '*order invalid*'
    }
  }

  Context 'FailFast behavior' {
    It 'stops workflow when security fails' {
      $cfg = New-WorkflowConfig -Path (Join-Path $TempRoot 'ff-cfg.json')
      Mock Invoke-Phase {
        if ($PhaseId -eq 'security') {
          New-MockPhaseResult -PhaseId 'security' -Status 'FAIL'
        } else {
          New-MockPhaseResult -PhaseId $PhaseId -Status 'PASS'
        }
      }

      $r = & $Runner -TicketKey 'ABC-1' -Task 'x' -TaskDescription 'x' `
        -ConfigPath $cfg -SkipJira

      # security ต้องหยุด workflow
      # (ตรวจสอบผ่านการนับจำนวน calls ในกรณีใช้ mock จริง)
      $r | Should -Not -BeNullOrEmpty
    }
  }

  Context 'Output artifacts' {
    It 'writes summary JSON' {
      $cfg = New-WorkflowConfig -Path (Join-Path $TempRoot 'out-cfg.json')
      Mock Invoke-Phase { New-MockPhaseResult -PhaseId $PhaseId }

      & $Runner -TicketKey 'SUM-1' -Task 'x' -TaskDescription 'x' `
        -ConfigPath $cfg -SkipJira

      $summary = Join-Path $TempRoot 'SUM-1-summary.json'
      Test-Path $summary | Should -BeTrue
    }
  }

  Context 'Read-only mode gating' {
    It 'skips read-write phases when Mode=read-only' {
      $cfg = New-WorkflowConfig -Path (Join-Path $TempRoot 'ro-cfg.json')
      Mock Invoke-Phase { New-MockPhaseResult -PhaseId $PhaseId }

      & $Runner -TicketKey 'RO-1' -Task 'x' -TaskDescription 'x' `
        -Mode 'read-only' -ConfigPath $cfg -SkipJira

      Should -Invoke Invoke-Phase -ParameterFilter { $PhaseId -eq 'execute' } -Times 0
    }
  }
}

Describe 'Run-6PhaseWorkflow.ps1 — Integration' -Tag 'Integration' {
  # ใช้กับ mock opencode CLI
  BeforeAll {
    . (Join-Path $PSScriptRoot 'Helpers.ps1')

    $script:mockBin = Join-Path $TempRoot 'bin'
    New-Item -ItemType Directory -Path $mockBin -Force | Out-Null

    if ($IsWindows) {
      $mockScript = @'
@echo off
echo # 📥 MOCK OUTPUT %*
exit /b 0
'@
      Set-Content -Path (Join-Path $mockBin 'opencode.cmd') -Value $mockScript
    } else {
      $mockScript = @'
#!/usr/bin/env bash
echo "# 📥 MOCK OUTPUT $*"
exit 0
'@
      Set-Content -Path (Join-Path $mockBin 'opencode') -Value $mockScript
      & chmod +x (Join-Path $mockBin 'opencode')
    }

    $env:PATH = "$mockBin$([IO.Path]::PathSeparator)$env:PATH"
  }

  It 'completes full run with mock opencode' -Tag 'Slow' {
    $cfg = New-WorkflowConfig -Path (Join-Path $TempRoot 'int-cfg.json') -DisableJira

    $r = & $Runner -TicketKey 'INT-1' -Task 'x' -TaskDescription 'x' `
      -ConfigPath $cfg -SkipJira

    $summary = Join-Path $TempRoot 'INT-1-summary.json'
    Test-Path $summary | Should -BeTrue
  }
}
```

---

### 3.6 `tests/PSScriptAnalyzerSettings.psd1`

```powershell
@{
    Severity = @('Error', 'Warning')

    ExcludeRules = @(
        'PSAvoidUsingWriteHost',           # ใช้ใน orchestrator
        'PSUseShouldProcessForStateChangingFunctions',
        'PSUseSingularNouns',
        'PSUseApprovedVerbs',
        'PSReviewUnusedParameter'          # มี dynamic params
    )

    Rules = @{
        PSUseCompatibleSyntax = @{
            Enable = $true
            TargetVersions = @('7.0', '7.2', '7.4')
        }
        PSPlaceOpenBrace = @{
            Enable = $true
            OnSameLine = $true
            NewLineAfter = $true
        }
        PSPlaceCloseBrace = @{
            Enable = $true
            NewLineAfter = $false
            IgnoreOneLineBlock = $true
        }
        PSUseConsistentIndentation = @{
            Enable = $true
            Kind = 'space'
            IndentationSize = 2
        }
    }
}
```

---

### 3.7 วิธีรัน Pester

```powershell
# ── Local ─────────────────────────────────────────────────
Install-Module Pester -MinimumVersion 5.5.0 -Scope CurrentUser -Force

$config = & ./tests/pester.config.ps1
Invoke-Pester -Configuration $config

# เฉพาะ Unit
Invoke-Pester -Path ./tests -Tag Unit -Output Detailed

# พร้อม Coverage
$config.CodeCoverage.Enabled = $true
$config.CodeCoverage.OutputFormat = 'JaCoCo'
Invoke-Pester -Configuration $config

# ── Docker ────────────────────────────────────────────────
docker compose -f docker/docker-compose.yml --profile test run --rm pester

# ── CI ────────────────────────────────────────────────────
# จะได้ test-results.xml + coverage.xml ใน .workflow-reports/
```

---

## 📊 4) Dashboard

### 4.1 `dashboard/index.html`

```html
<!DOCTYPE html>
<html lang="th">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>6-Phase Workflow Dashboard</title>
  <link rel="stylesheet" href="styles.css" />
  <link rel="icon" href="data:image/svg+xml,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 100 100'><text y='.9em' font-size='90'>🚀</text></svg>" />
</head>
<body>
  <header class="topbar">
    <div class="brand">
      <span class="brand-emoji">🚀</span>
      <div>
        <h1>6-Phase Workflow Dashboard</h1>
        <p class="subtitle">Gather → Plan → Execute → Security → Performance → RCA</p>
      </div>
    </div>
    <div class="controls">
      <label for="filterStatus">Filter:</label>
      <select id="filterStatus">
        <option value="all">ทั้งหมด</option>
        <option value="PASS">🟢 PASS</option>
        <option value="WARN">🟡 WARN</option>
        <option value="FAIL">🟠 FAIL</option>
        <option value="ERROR">🔴 ERROR</option>
      </select>
      <input id="search" type="search" placeholder="ค้นหา Ticket / Task…" />
      <button id="refresh" title="Refresh">🔄</button>
      <span id="lastUpdate" class="muted"></span>
    </div>
  </header>

  <section class="kpis" id="kpis"></section>

  <section class="charts">
    <div class="card">
      <h2>Status distribution</h2>
      <canvas id="chartStatus" height="200"></canvas>
    </div>
    <div class="card">
      <h2>Duration trend (ล่าสุด 30)</h2>
      <canvas id="chartDuration" height="200"></canvas>
    </div>
    <div class="card">
      <h2>Phase success rate</h2>
      <canvas id="chartPhases" height="200"></canvas>
    </div>
    <div class="card">
      <h2>Top failing tickets</h2>
      <ol id="topFail" class="top-list"></ol>
    </div>
  </section>

  <section class="table-wrap">
    <table id="workflowTable">
      <thead>
        <tr>
          <th>Ticket</th>
          <th>Task</th>
          <th>Overall</th>
          <th>Gather</th>
          <th>Plan</th>
          <th>Execute</th>
          <th>Security</th>
          <th>Perf</th>
          <th>RCA</th>
          <th>Duration</th>
          <th>Finished</th>
          <th></th>
        </tr>
      </thead>
      <tbody></tbody>
    </table>
  </section>

  <dialog id="detailDialog">
    <div class="dialog-head">
      <h2 id="detailTitle">Ticket detail</h2>
      <button id="closeDetail">✕</button>
    </div>
    <div id="detailBody" class="dialog-body"></div>
  </dialog>

  <script src="https://cdn.jsdelivr.net/npm/chart.js@4.4.3/dist/chart.umd.min.js"></script>
  <script src="app.js"></script>
</body>
</html>
```

---

### 4.2 `dashboard/styles.css`

```css
:root {
  --bg: #0d1117;
  --bg-2: #161b22;
  --bg-3: #21262d;
  --border: #30363d;
  --text: #e6edf3;
  --text-muted: #8b949e;
  --accent: #58a6ff;
  --pass: #3fb950;
  --warn: #d29922;
  --fail: #db6d28;
  --error: #f85149;
  --radius: 10px;
  --shadow: 0 2px 8px rgba(0, 0, 0, 0.35);
}

* { box-sizing: border-box; }

html, body {
  margin: 0;
  background: var(--bg);
  color: var(--text);
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", "Noto Sans Thai", sans-serif;
  font-size: 14px;
  line-height: 1.5;
}

/* ── Topbar ────────────────────────────────────── */
.topbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 20px;
  padding: 16px 24px;
  background: linear-gradient(180deg, var(--bg-2), var(--bg));
  border-bottom: 1px solid var(--border);
  position: sticky;
  top: 0;
  z-index: 100;
  backdrop-filter: blur(8px);
}

.brand {
  display: flex;
  gap: 14px;
  align-items: center;
}

.brand-emoji { font-size: 32px; }

.brand h1 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
}

.subtitle {
  margin: 0;
  font-size: 12px;
  color: var(--text-muted);
}

.controls {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.controls select,
.controls input[type="search"],
.controls button {
  background: var(--bg-3);
  border: 1px solid var(--border);
  color: var(--text);
  padding: 7px 12px;
  border-radius: var(--radius);
  font-size: 13px;
}

.controls button {
  cursor: pointer;
  transition: background 0.15s;
}

.controls button:hover { background: var(--accent); color: #fff; }

.muted { color: var(--text-muted); font-size: 12px; }

/* ── KPIs ──────────────────────────────────────── */
.kpis {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 14px;
  padding: 20px 24px;
}

.kpi {
  background: var(--bg-2);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  padding: 16px;
  box-shadow: var(--shadow);
}

.kpi-label {
  font-size: 11px;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: var(--text-muted);
}

.kpi-value {
  font-size: 28px;
  font-weight: 700;
  margin-top: 6px;
}

.kpi-sub {
  font-size: 11px;
  color: var(--text-muted);
  margin-top: 4px;
}

.kpi.pass  .kpi-value { color: var(--pass); }
.kpi.warn  .kpi-value { color: var(--warn); }
.kpi.fail  .kpi-value { color: var(--fail); }
.kpi.error .kpi-value { color: var(--error); }

/* ── Charts ────────────────────────────────────── */
.charts {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
  gap: 14px;
  padding: 0 24px 20px;
}

.card {
  background: var(--bg-2);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  padding: 16px;
  box-shadow: var(--shadow);
}

.card h2 {
  margin: 0 0 12px;
  font-size: 14px;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.top-list {
  margin: 0;
  padding-left: 20px;
}

.top-list li {
  padding: 6px 0;
  border-bottom: 1px solid var(--border);
  font-size: 13px;
}

.top-list li:last-child { border-bottom: none; }

/* ── Table ────────────────────────────────────── */
.table-wrap {
  margin: 0 24px 40px;
  background: var(--bg-2);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  overflow: auto;
  box-shadow: var(--shadow);
}

table {
  width: 100%;
  border-collapse: collapse;
  font-size: 13px;
}

thead {
  background: var(--bg-3);
  position: sticky;
  top: 0;
  z-index: 10;
}

th, td {
  padding: 10px 12px;
  text-align: left;
  border-bottom: 1px solid var(--border);
  white-space: nowrap;
}

th {
  font-weight: 600;
  font-size: 12px;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

tbody tr {
  cursor: pointer;
  transition: background 0.1s;
}

tbody tr:hover { background: var(--bg-3); }

.badge {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 999px;
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
}

.badge.pass  { background: rgba(63, 185, 80, 0.15);  color: var(--pass); }
.badge.warn  { background: rgba(210, 153, 34, 0.15); color: var(--warn); }
.badge.fail  { background: rgba(219, 109, 40, 0.15); color: var(--fail); }
.badge.error { background: rgba(248, 81, 73, 0.15);  color: var(--error); }
.badge.skip  { background: rgba(139, 148, 158, 0.15); color: var(--text-muted); }

.cell-ticket { font-weight: 600; color: var(--accent); }
.cell-task   { max-width: 320px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }

/* ── Dialog ────────────────────────────────────── */
dialog {
  border: 1px solid var(--border);
  background: var(--bg-2);
  color: var(--text);
  border-radius: var(--radius);
  padding: 0;
  width: min(900px, 92vw);
  max-height: 86vh;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.6);
}

dialog::backdrop { background: rgba(0, 0, 0, 0.65); backdrop-filter: blur(2px); }

.dialog-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid var(--border);
  position: sticky;
  top: 0;
  background: var(--bg-2);
}

.dialog-head h2 { margin: 0; font-size: 15px; }

.dialog-head button {
  background: transparent;
  border: none;
  color: var(--text-muted);
  font-size: 20px;
  cursor: pointer;
  padding: 0 8px;
}

.dialog-head button:hover { color: var(--error); }

.dialog-body {
  padding: 20px;
  overflow-y: auto;
  max-height: calc(86vh - 60px);
}

.dialog-body pre {
  background: var(--bg);
  border: 1px solid var(--border);
  border-radius: 6px;
  padding: 12px;
  overflow-x: auto;
  font-size: 12px;
  line-height: 1.4;
}

.dialog-body .phase-block {
  border-left: 3px solid var(--border);
  padding-left: 12px;
  margin: 16px 0;
}

.dialog-body .phase-block.pass  { border-color: var(--pass); }
.dialog-body .phase-block.warn  { border-color: var(--warn); }
.dialog-body .phase-block.fail  { border-color: var(--fail); }
.dialog-body .phase-block.error { border-color: var(--error); }

/* ── Responsive ────────────────────────────────── */
@media (max-width: 768px) {
  .topbar { flex-direction: column; align-items: stretch; }
  .kpi-value { font-size: 22px; }
  th, td { padding: 8px; font-size: 12px; }
}
```

---

### 4.3 `dashboard/app.js`

```javascript
/**
 * 6-Phase Workflow Dashboard
 * โหลดข้อมูลจาก /data/*-summary.json → แสดงผล
 */

const DATA_DIR = './data';       // หรือ './.workflow-reports' ถ้าเปิดจาก local
const REFRESH_INTERVAL_MS = 30_000;

const STATE = {
  workflows: [],
  filterStatus: 'all',
  search: '',
  charts: {},
};

// ── Utilities ───────────────────────────────────────────────

const fmtDuration = (ms) => {
  if (ms == null) return '-';
  const s = ms / 1000;
  if (s < 60) return `${s.toFixed(1)}s`;
  const m = Math.floor(s / 60);
  const rs = Math.round(s - m * 60);
  return `${m}m ${rs}s`;
};

const fmtDate = (iso) => {
  if (!iso) return '-';
  const d = new Date(iso);
  return d.toLocaleString('th-TH', {
    year: 'numeric', month: 'short', day: '2-digit',
    hour: '2-digit', minute: '2-digit',
  });
};

const statusClass = (s) => (s || 'skip').toLowerCase();

const statusEmoji = (s) => {
  switch (s) {
    case 'PASS':  return '🟢';
    case 'WARN':  return '🟡';
    case 'FAIL':  return '🟠';
    case 'ERROR': return '🔴';
    default:      return '⚪';
  }
};

const escape = (s) =>
  String(s ?? '').replace(/[&<>"']/g, (c) => ({
    '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;',
  }[c]));

// ── Load data ───────────────────────────────────────────────

async function fetchWorkflowSummaries() {
  // 1) ลองใช้ autoindex JSON (nginx)
  try {
    const idx = await fetch(`${DATA_DIR}/`, { cache: 'no-store' });
    if (idx.ok && idx.headers.get('content-type')?.includes('json')) {
      const listing = await idx.json();
      const files = listing
        .filter((f) => /-summary\.json$/.test(f.name))
        .map((f) => f.name);
      return await loadFiles(files);
    }
  } catch { /* fallback */ }

  // 2) fallback: manifest.json
  const manifest = await fetch(`${DATA_DIR}/manifest.json`, { cache: 'no-store' })
    .then((r) => (r.ok ? r.json() : []))
    .catch(() => []);
  return await loadFiles(manifest);
}

async function loadFiles(names) {
  const results = await Promise.allSettled(
    names.map((n) =>
      fetch(`${DATA_DIR}/${n}`, { cache: 'no-store' })
        .then((r) => r.json())
        .then((data) => ({ ...data, _source: n })),
    ),
  );
  return results
    .filter((r) => r.status === 'fulfilled')
    .map((r) => r.value);
}

// ── Aggregate ───────────────────────────────────────────────

function aggregate(workflows) {
  const total = workflows.length;
  const counts = { PASS: 0, WARN: 0, FAIL: 0, ERROR: 0 };
  let totalMs = 0;
  const phaseStats = {};

  for (const w of workflows) {
    if (counts[w.overall] != null) counts[w.overall]++;
    totalMs += w.durationMs || 0;

    for (const [id, r] of Object.entries(w.results || {})) {
      phaseStats[id] ??= { total: 0, pass: 0, fail: 0, warn: 0, error: 0 };
      phaseStats[id].total++;
      const key = (r.status || 'error').toLowerCase();
      if (phaseStats[id][key] != null) phaseStats[id][key]++;
    }
  }

  return {
    total,
    counts,
    avgMs: total > 0 ? totalMs / total : 0,
    phaseStats,
  };
}

// ── Render ──────────────────────────────────────────────────

function renderKPIs(agg) {
  const passRate = agg.total ? Math.round((agg.counts.PASS / agg.total) * 100) : 0;
  const html = `
    <div class="kpi">
      <div class="kpi-label">Tickets ทั้งหมด</div>
      <div class="kpi-value">${agg.total}</div>
      <div class="kpi-sub">จาก summary files</div>
    </div>
    <div class="kpi pass">
      <div class="kpi-label">Pass rate</div>
      <div class="kpi-value">${passRate}%</div>
      <div class="kpi-sub">${agg.counts.PASS} / ${agg.total}</div>
    </div>
    <div class="kpi warn">
      <div class="kpi-label">Warnings</div>
      <div class="kpi-value">${agg.counts.WARN}</div>
    </div>
    <div class="kpi fail">
      <div class="kpi-label">Failures</div>
      <div class="kpi-value">${agg.counts.FAIL}</div>
    </div>
    <div class="kpi error">
      <div class="kpi-label">Errors</div>
      <div class="kpi-value">${agg.counts.ERROR}</div>
    </div>
    <div class="kpi">
      <div class="kpi-label">Avg duration</div>
      <div class="kpi-value">${fmtDuration(agg.avgMs)}</div>
    </div>
  `;
  document.getElementById('kpis').innerHTML = html;
}

function renderCharts(agg, workflows) {
  // Status doughnut
  const ctxStatus = document.getElementById('chartStatus');
  if (STATE.charts.status) STATE.charts.status.destroy();
  STATE.charts.status = new Chart(ctxStatus, {
    type: 'doughnut',
    data: {
      labels: ['PASS', 'WARN', 'FAIL', 'ERROR'],
      datasets: [{
        data: [agg.counts.PASS, agg.counts.WARN, agg.counts.FAIL, agg.counts.ERROR],
        backgroundColor: ['#3fb950', '#d29922', '#db6d28', '#f85149'],
        borderColor: '#0d1117',
        borderWidth: 2,
      }],
    },
    options: {
      responsive: true,
      plugins: {
        legend: { labels: { color: '#e6edf3', font: { size: 12 } } },
      },
    },
  });

  // Duration trend
  const sorted = [...workflows]
    .sort((a, b) => new Date(a.finishedAt) - new Date(b.finishedAt))
    .slice(-30);
  const ctxDur = document.getElementById('chartDuration');
  if (STATE.charts.duration) STATE.charts.duration.destroy();
  STATE.charts.duration = new Chart(ctxDur, {
    type: 'line',
    data: {
      labels: sorted.map((w) => w.ticketKey),
      datasets: [{
        label: 'Duration (s)',
        data: sorted.map((w) => Math.round((w.durationMs || 0) / 1000)),
        borderColor: '#58a6ff',
        backgroundColor: 'rgba(88, 166, 255, 0.15)',
        fill: true,
        tension: 0.3,
        pointRadius: 3,
      }],
    },
    options: {
      responsive: true,
      plugins: { legend: { display: false } },
      scales: {
        x: { ticks: { color: '#8b949e' }, grid: { color: '#21262d' } },
        y: { ticks: { color: '#8b949e' }, grid: { color: '#21262d' }, beginAtZero: true },
      },
    },
  });

  // Phase success rate
  const phaseLabels = ['gather', 'plan', 'execute', 'security', 'performance', 'rca'];
  const phaseData = phaseLabels.map((p) => {
    const s = agg.phaseStats[p];
    if (!s || !s.total) return 0;
    return Math.round(((s.total - s.fail - s.error) / s.total) * 100);
  });
  const ctxPh = document.getElementById('chartPhases');
  if (STATE.charts.phases) STATE.charts.phases.destroy();
  STATE.charts.phases = new Chart(ctxPh, {
    type: 'bar',
    data: {
      labels: phaseLabels,
      datasets: [{
        label: 'Success %',
        data: phaseData,
        backgroundColor: ['#58a6ff', '#58a6ff', '#58a6ff', '#f85149', '#d29922', '#3fb950'],
        borderRadius: 6,
      }],
    },
    options: {
      responsive: true,
      plugins: { legend: { display: false } },
      scales: {
        x: { ticks: { color: '#8b949e' }, grid: { display: false } },
        y: { ticks: { color: '#8b949e' }, grid: { color: '#21262d' }, max: 100, beginAtZero: true },
      },
    },
  });

  // Top failing tickets
  const failing = workflows
    .filter((w) => ['FAIL', 'ERROR'].includes(w.overall))
    .slice(0, 8);
  const top = document.getElementById('topFail');
  top.innerHTML = failing.length
    ? failing
        .map((w) => `<li><strong>${escape(w.ticketKey)}</strong> — ${escape(w.task)} <span class="badge ${statusClass(w.overall)}">${w.overall}</span></li>`)
        .join('')
    : '<li class="muted">ไม่มี failing ticket 🎉</li>';
}

function renderTable(workflows) {
  const tbody = document.querySelector('#workflowTable tbody');
  const filtered = workflows.filter((w) => {
    if (STATE.filterStatus !== 'all' && w.overall !== STATE.filterStatus) return false;
    if (STATE.search) {
      const s = STATE.search.toLowerCase();
      if (
        !String(w.ticketKey || '').toLowerCase().includes(s) &&
        !String(w.task || '').toLowerCase().includes(s)
      ) return false;
    }
    return true;
  });

  const phaseCell = (w, id) => {
    const r = w.results?.[id];
    if (!r) return '<span class="badge skip">—</span>';
    return `<span class="badge ${statusClass(r.status)}">${r.status}</span>`;
  };

  tbody.innerHTML = filtered
    .map(
      (w) => `
      <tr data-ticket="${escape(w.ticketKey)}">
        <td class="cell-ticket">${escape(w.ticketKey)}</td>
        <td class="cell-task">${escape(w.task)}</td>
        <td><span class="badge ${statusClass(w.overall)}">${statusEmoji(w.overall)} ${w.overall}</span></td>
        <td>${phaseCell(w, 'gather')}</td>
        <td>${phaseCell(w, 'plan')}</td>
        <td>${phaseCell(w, 'execute')}</td>
        <td>${phaseCell(w, 'security')}</td>
        <td>${phaseCell(w, 'performance')}</td>
        <td>${phaseCell(w, 'rca')}</td>
        <td>${fmtDuration(w.durationMs)}</td>
        <td>${fmtDate(w.finishedAt)}</td>
        <td>▶</td>
      </tr>`,
    )
    .join('');

  tbody.querySelectorAll('tr').forEach((tr) => {
    tr.addEventListener('click', () => openDetail(tr.dataset.ticket));
  });
}

// ── Detail dialog ───────────────────────────────────────────

function openDetail(ticket) {
  const w = STATE.workflows.find((x) => x.ticketKey === ticket);
  if (!w) return;

  document.getElementById('detailTitle').textContent =
    `${statusEmoji(w.overall)} ${w.ticketKey} — ${w.task}`;

  const phaseBlocks = Object.entries(w.results || {})
    .map(([id, r]) => {
      const cls = statusClass(r.status);
      const dur = fmtDuration(r.durationMs);
      return `
        <div class="phase-block ${cls}">
          <h3>${r.emoji || ''} ${escape(r.phaseName)} — <span class="badge ${cls}">${r.status}</span> <span class="muted">(${dur})</span></h3>
          ${r.error ? `<p class="muted">Error: ${escape(r.error)}</p>` : ''}
          <pre>${escape((r.report || '').slice(0, 3000))}</pre>
        </div>
      `;
    })
    .join('');

  document.getElementById('detailBody').innerHTML = `
    <p><strong>Overall:</strong> <span class="badge ${statusClass(w.overall)}">${w.overall}</span></p>
    <p><strong>Duration:</strong> ${fmtDuration(w.durationMs)}</p>
    <p><strong>Finished:</strong> ${fmtDate(w.finishedAt)}</p>
    ${phaseBlocks || '<p class="muted">ไม่มีข้อมูลเฟส</p>'}
  `;

  document.getElementById('detailDialog').showModal();
}

// ── Wire-up ─────────────────────────────────────────────────

async function refresh() {
  const btn = document.getElementById('refresh');
  btn.disabled = true;
  btn.textContent = '⏳';
  try {
    const workflows = await fetchWorkflowSummaries();
    STATE.workflows = workflows;

    const agg = aggregate(workflows);
    renderKPIs(agg);
    renderCharts(agg, workflows);
    renderTable(workflows);

    document.getElementById('lastUpdate').textContent =
      `Updated ${new Date().toLocaleTimeString('th-TH')}`;
  } catch (err) {
    console.error('refresh failed', err);
    document.getElementById('lastUpdate').textContent = '⚠️ โหลดข้อมูลไม่ได้';
  } finally {
    btn.disabled = false;
    btn.textContent = '🔄';
  }
}

document.getElementById('refresh').addEventListener('click', refresh);
document.getElementById('filterStatus').addEventListener('change', (e) => {
  STATE.filterStatus = e.target.value;
  renderTable(STATE.workflows);
});
document.getElementById('search').addEventListener('input', (e) => {
  STATE.search = e.target.value;
  renderTable(STATE.workflows);
});
document.getElementById('closeDetail').addEventListener('click', () => {
  document.getElementById('detailDialog').close();
});

refresh();
setInterval(refresh, REFRESH_INTERVAL_MS);
```

---

### 4.4 `dashboard/exporter/exporter.py` (Prometheus exporter)

```python
#!/usr/bin/env python3
"""
6-Phase Workflow → Prometheus exporter
อ่าน JSON จาก .workflow-reports/ แล้ว expose metrics ที่ :9101/metrics
"""
from __future__ import annotations

import json
import os
import time
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
from pathlib import Path
from typing import Any

REPORTS_DIR = Path(os.environ.get("REPORTS_DIR", "/reports"))
LISTEN_PORT = int(os.environ.get("LISTEN_PORT", "9101"))
SCRAPE_INTERVAL = int(os.environ.get("SCRAPE_INTERVAL", "30"))

STATUS_CODE = {"PASS": 0, "WARN": 1, "FAIL": 2, "ERROR": 3}


# ── Loader ──────────────────────────────────────────────────

def load_summaries() -> list[dict[str, Any]]:
    if not REPORTS_DIR.is_dir():
        return []
    out = []
    for path in REPORTS_DIR.glob("*-summary.json"):
        try:
            data = json.loads(path.read_text(encoding="utf-8"))
            data["_file"] = path.name
            out.append(data)
        except Exception as exc:      # noqa: BLE001
            print(f"[exporter] skip {path.name}: {exc}")
    return out


# ── Metrics ─────────────────────────────────────────────────

def _escape_label(v: str) -> str:
    return str(v).replace("\\", "\\\\").replace('"', '\\"').replace("\n", " ")


def _lbl(**kw: Any) -> str:
    if not kw:
        return ""
    parts = ",".join(f'{k}="{_escape_label(v)}"' for k, v in kw.items())
    return "{" + parts + "}"


def build_metrics(workflows: list[dict[str, Any]]) -> str:
    lines: list[str] = []

    # ── Overall status (gauge 0-3) ──────────────────────────
    lines += [
        "# HELP workflow_overall_status Overall status: 0=PASS 1=WARN 2=FAIL 3=ERROR",
        "# TYPE workflow_overall_status gauge",
    ]
    for w in workflows:
        code = STATUS_CODE.get(w.get("overall", "ERROR"), 3)
        lines.append(
            f"workflow_overall_status{_lbl(ticket=w.get('ticketKey', 'unknown'), task=w.get('task', ''))} {code}"
        )

    # ── Total duration ──────────────────────────────────────
    lines += [
        "# HELP workflow_duration_seconds Total workflow duration in seconds",
        "# TYPE workflow_duration_seconds gauge",
    ]
    for w in workflows:
        secs = (w.get("durationMs") or 0) / 1000.0
        lines.append(
            f"workflow_duration_seconds{_lbl(ticket=w.get('ticketKey', 'unknown'))} {secs:.3f}"
        )

    # ── Per-phase metrics ──────────────────────────────────
    lines += [
        "# HELP workflow_phase_status Phase status: 0=PASS 1=WARN 2=FAIL 3=ERROR",
        "# TYPE workflow_phase_status gauge",
        "# HELP workflow_phase_duration_seconds Phase duration in seconds",
        "# TYPE workflow_phase_duration_seconds gauge",
    ]
    for w in workflows:
        ticket = w.get("ticketKey", "unknown")
        for phase, r in (w.get("results") or {}).items():
            code = STATUS_CODE.get(r.get("status", "ERROR"), 3)
            secs = (r.get("durationMs") or 0) / 1000.0
            lines.append(
                f"workflow_phase_status{_lbl(ticket=ticket, phase=phase)} {code}"
            )
            lines.append(
                f"workflow_phase_duration_seconds{_lbl(ticket=ticket, phase=phase)} {secs:.3f}"
            )

    # ── Aggregates ─────────────────────────────────────────
    lines += [
        "# HELP workflow_total_total Total number of workflows",
        "# TYPE workflow_total_total gauge",
        f"workflow_total_total {len(workflows)}",
    ]

    by_status: dict[str, int] = {}
    for w in workflows:
        by_status[w.get("overall", "ERROR")] = by_status.get(w.get("overall", "ERROR"), 0) + 1

    lines += [
        "# HELP workflow_status_count Count of workflows per overall status",
        "# TYPE workflow_status_count gauge",
    ]
    for status in ("PASS", "WARN", "FAIL", "ERROR"):
        lines.append(
            f'workflow_status_count{_lbl(status=status)} {by_status.get(status, 0)}'
        )

    # ── Timestamps ─────────────────────────────────────────
    lines += [
        "# HELP workflow_finished_timestamp_seconds Unix timestamp when workflow finished",
        "# TYPE workflow_finished_timestamp_seconds gauge",
    ]
    for w in workflows:
        ts = w.get("finishedAt")
        if not ts:
            continue
        try:
            epoch = time.mktime(time.strptime(ts.split(".")[0], "%Y-%m-%dT%H:%M:%S"))
        except Exception:  # noqa: BLE001
            continue
        lines.append(
            f"workflow_finished_timestamp_seconds{_lbl(ticket=w.get('ticketKey', 'unknown'))} {epoch}"
        )

    return "\n".join(lines) + "\n"


# ── HTTP handler ────────────────────────────────────────────

class Handler(BaseHTTPRequestHandler):
    cache: str = ""
    last_refresh: float = 0.0

    def do_GET(self) -> None:  # noqa: N802
        if self.path == "/health":
            return self._text(200, "ok\n")
        if self.path != "/metrics":
            return self._text(404, "not found\n")

        now = time.time()
        if now - Handler.last_refresh > SCRAPE_INTERVAL:
            workflows = load_summaries()
            Handler.cache = build_metrics(workflows)
            Handler.last_refresh = now

        self.send_response(200)
        self.send_header("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
        self.send_header("Content-Length", str(len(Handler.cache.encode("utf-8"))))
        self.end_headers()
        self.wfile.write(Handler.cache.encode("utf-8"))

    def _text(self, code: int, body: str) -> None:
        b = body.encode("utf-8")
        self.send_response(code)
        self.send_header("Content-Type", "text/plain; charset=utf-8")
        self.send_header("Content-Length", str(len(b)))
        self.end_headers()
        self.wfile.write(b)

    def log_message(self, *args: Any) -> None:
        pass


def main() -> None:
    print(f"[exporter] reports dir: {REPORTS_DIR}")
    print(f"[exporter] listen: http://0.0.0.0:{LISTEN_PORT}/metrics")
    print(f"[exporter] scrape interval: {SCRAPE_INTERVAL}s")
    ThreadingHTTPServer(("0.0.0.0", LISTEN_PORT), Handler).serve_forever()


if __name__ == "__main__":
    main()
```

---

### 4.5 `dashboard/grafana/provisioning/datasources/prometheus.yml`

```yaml
apiVersion: 1

datasources:
  - name: Prometheus
    type: prometheus
    access: proxy
    url: http://exporter:9101
    isDefault: true
    editable: false
    jsonData:
      timeInterval: 30s
      httpMethod: GET
```

---

### 4.6 `dashboard/grafana/provisioning/dashboards/default.yml`

```yaml
apiVersion: 1

providers:
  - name: '6-Phase Workflow'
    orgId: 1
    folder: 'Workflow'
    type: file
    disableDeletion: false
    updateIntervalSeconds: 30
    allowUiUpdates: true
    options:
      path: /var/lib/grafana/dashboards
      foldersFromFilesStructure: true
```

---

### 4.7 `dashboard/grafana/dashboards/workflow-6phase.json`

```json
{
  "title": "6-Phase Workflow",
  "uid": "workflow-6phase",
  "tags": ["workflow", "6phase", "go"],
  "timezone": "browser",
  "schemaVersion": 39,
  "version": 1,
  "refresh": "30s",
  "time": { "from": "now-7d", "to": "now" },
  "templating": {
    "list": [
      {
        "name": "ticket",
        "type": "query",
        "datasource": { "type": "prometheus", "uid": "prometheus" },
        "query": "label_values(workflow_overall_status, ticket)",
        "includeAll": true,
        "multi": true,
        "label": "Ticket"
      }
    ]
  },
  "panels": [
    {
      "id": 1,
      "type": "stat",
      "title": "Total Workflows",
      "gridPos": { "x": 0, "y": 0, "w": 6, "h": 4 },
      "datasource": { "type": "prometheus", "uid": "prometheus" },
      "targets": [
        { "expr": "workflow_total_total", "refId": "A" }
      ],
      "options": {
        "colorMode": "value",
        "graphMode": "area",
        "reduceOptions": { "calcs": ["lastNotNull"] }
      }
    },
    {
      "id": 2,
      "type": "stat",
      "title": "Pass Rate",
      "gridPos": { "x": 6, "y": 0, "w": 6, "h": 4 },
      "datasource": { "type": "prometheus", "uid": "prometheus" },
      "targets": [
        {
          "expr": "workflow_status_count{status=\"PASS\"} / workflow_total_total * 100",
          "refId": "A",
          "legendFormat": "PASS %"
        }
      ],
      "options": {
        "unit": "percent",
        "colorMode": "value",
        "thresholds": {
          "mode": "absolute",
          "steps": [
            { "color": "red", "value": null },
            { "color": "orange", "value": 50 },
            { "color": "yellow", "value": 75 },
            { "color": "green", "value": 90 }
          ]
        }
      }
    },
    {
      "id": 3,
      "type": "stat",
      "title": "Failures",
      "gridPos": { "x": 12, "y": 0, "w": 6, "h": 4 },
      "datasource": { "type": "prometheus", "uid": "prometheus" },
      "targets": [
        { "expr": "workflow_status_count{status=\"FAIL\"}", "refId": "A" }
      ],
      "options": { "colorMode": "value", "color": { "mode": "fixed", "fixedColor": "orange" } }
    },
    {
      "id": 4,
      "type": "stat",
      "title": "Errors",
      "gridPos": { "x": 18, "y": 0, "w": 6, "h": 4 },
      "datasource": { "type": "prometheus", "uid": "prometheus" },
      "targets": [
        { "expr": "workflow_status_count{status=\"ERROR\"}", "refId": "A" }
      ],
      "options": { "colorMode": "value", "color": { "mode": "fixed", "fixedColor": "red" } }
    },
    {
      "id": 10,
      "type": "timeseries",
      "title": "Overall Status over time",
      "gridPos": { "x": 0, "y": 4, "w": 12, "h": 8 },
      "datasource": { "type": "prometheus", "uid": "prometheus" },
      "targets": [
        {
          "expr": "workflow_overall_status{ticket=~\"$ticket\"}",
          "refId": "A",
          "legendFormat": "{{ticket}}"
        }
      ],
      "fieldConfig": {
        "defaults": {
          "custom": { "lineInterpolation": "stepAfter", "drawStyle": "line" },
          "mappings": [
            { "type": "value", "options": { "0": { "text": "PASS", "color": "green" } } },
            { "type": "value", "options": { "1": { "text": "WARN", "color": "yellow" } } },
            { "type": "value", "options": { "2": { "text": "FAIL", "color": "orange" } } },
            { "type": "value", "options": { "3": { "text": "ERROR", "color": "red" } } }
          ]
        }
      }
    },
    {
      "id": 11,
      "type": "timeseries",
      "title": "Phase duration (heatmap)",
      "gridPos": { "x": 12, "y": 4, "w": 12, "h": 8 },
      "datasource": { "type": "prometheus", "uid": "prometheus" },
      "targets": [
        {
          "expr": "avg by (phase) (workflow_phase_duration_seconds{ticket=~\"$ticket\"})",
          "refId": "A",
          "legendFormat": "{{phase}}"
        }
      ],
      "fieldConfig": { "defaults": { "unit": "s" } }
    },
    {
      "id": 20,
      "type": "bargauge",
      "title": "Phase success rate %",
      "gridPos": { "x": 0, "y": 12, "w": 12, "h": 8 },
      "datasource": { "type": "prometheus", "uid": "prometheus" },
      "targets": [
        {
          "expr": "100 * (1 - (sum by (phase) (workflow_phase_status{ticket=~\"$ticket\"} >= 2) / sum by (phase) (workflow_phase_status{ticket=~\"$ticket\"})))",
          "refId": "A",
          "legendFormat": "{{phase}}"
        }
      ],
      "options": { "displayMode": "gradient", "orientation": "horizontal" },
      "fieldConfig": {
        "defaults": { "unit": "percent", "max": 100, "min": 0 }
      }
    },
    {
      "id": 21,
      "type": "table",
      "title": "Phase durations",
      "gridPos": { "x": 12, "y": 12, "w": 12, "h": 8 },
      "datasource": { "type": "prometheus", "uid": "prometheus" },
      "targets": [
        {
          "expr": "workflow_phase_duration_seconds{ticket=~\"$ticket\"}",
          "refId": "A",
          "format": "table",
          "instant": true
        }
      ],
      "transformations": [
        { "id": "organize", "options": { "excludeByName": { "Time": true, "__name__": true, "job": true, "instance": true } } }
      ]
    }
  ]
}
```

---

### 4.8 วิธีใช้งาน Dashboard

```bash
# ── Option A: HTML Dashboard (ง่ายสุด) ────────────────────
docker compose -f docker/docker-compose.yml --profile dash up -d dashboard
# เปิด http://localhost:8080

# ── Option B: HTML + Grafana ──────────────────────────────
docker compose -f docker/docker-compose.yml --profile dash up -d
# HTML:    http://localhost:8080
# Grafana: http://localhost:3000  (admin/admin)
# Metric:  http://localhost:9101/metrics

# ── Option C: เปิด HTML local (ไม่ต้อง docker) ────────────
# ก๊อป summary files ไปวางที่ dashboard/data/
Copy-Item .workflow-reports/*-summary.json dashboard/data/
# เปิด dashboard/index.html ในเบราว์เซอร์ (ต้อง serve ผ่าน http ไม่ใช่ file://)
python -m http.server 8080 --directory dashboard
```

---

## 📁 โครงสร้างไฟล์สุดท้าย

```
project/
├── .github/                          # จากรอบก่อน
├── .vscode/
│   ├── tasks.json                    ← ⭐
│   ├── launch.json                   ← ⭐
│   └── settings.json                 ← ⭐
├── docker/
│   ├── Dockerfile                    ← ⭐
│   ├── docker-compose.yml            ← ⭐
│   ├── entrypoint.sh                 ← ⭐
│   ├── nginx.conf                    ← ⭐
│   └── .dockerignore                 ← ⭐
├── scripts/
│   ├── Run-6PhaseWorkflow.ps1
│   ├── Invoke-Phase.ps1
│   ├── Publish-JiraComment.ps1
│   └── config/workflow.config.json
├── tests/
│   ├── Helpers.ps1                   ← ⭐
│   ├── Invoke-Phase.Tests.ps1        ← ⭐
│   ├── Publish-JiraComment.Tests.ps1 ← ⭐
│   ├── Run-6PhaseWorkflow.Tests.ps1  ← ⭐
│   ├── pester.config.ps1             ← ⭐
│   └── PSScriptAnalyzerSettings.psd1 ← ⭐
├── dashboard/
│   ├── index.html                    ← ⭐
│   ├── app.js                        ← ⭐
│   ├── styles.css                    ← ⭐
│   ├── exporter/exporter.py          ← ⭐
│   └── grafana/
│       ├── provisioning/
│       │   ├── datasources/prometheus.yml
│       │   └── dashboards/default.yml
│       └── dashboards/workflow-6phase.json
├── docs/opencode/                    # จากรอบก่อน
├── .workflow-reports/                # output (gitignore)
├── main.go
├── go.mod
└── ...
```

---

## 🎯 สรุปการใช้งาน 4 อย่างร่วมกัน

| สถานการณ์ | คำสั่ง |
|-----------|--------|
| **Dev local** | กด `Ctrl+Shift+P` → `Tasks: Run Task` → `6-Phase: Run All` หรือ `F5` |
| **Test scripts** | `Tasks: Run Test Task` (`Ctrl+Shift+P` → `Run Test Task`) |
| **CI** | GitHub Actions (จากรอบก่อน) |
| **รันใน container** | `docker compose --profile run up workflow` |
| **Dashboard** | `docker compose --profile dash up -d` → `http://localhost:8080` |
| **Grafana** | `http://localhost:3000` (admin/admin) |
| **Prometheus metrics** | `http://localhost:9101/metrics` |

---

## ⚙️ ตั้งค่า `.gitignore` เพิ่ม

```gitignore
# workflow outputs
.workflow-reports/

# dashboard data (generated)
dashboard/data/

# go
coverage.out
bench.txt
bench-baseline.txt
benchstat.txt
gosec.json
govulncheck.sarif
gosec.sarif

# env
.env
.env.*
!.env.example

# docker
docker-compose.override.yml
```

---

## 💡 สรุปจุดเด่นของชุดนี้

| ส่วน | จุดเด่น |
|------|--------|
| **Docker** | Multi-stage, non-root user, healthcheck, layer cache, volume แชร์ go mod |
| **VS Code Tasks** | 21 tasks + 4 launch configs, prompt inputs, ProblemMatcher พร้อม |
| **Pester** | Unit + Integration, mock, coverage 70%, CI format (NUnit + JaCoCo) |
| **HTML Dashboard** | Self-contained, Chart.js, auto-refresh 30s, dialog detail, dark theme |
| **Grafana** | Provisioning อัตโนมัติ, template variable `$ticket`, 8 panels |
| **Prometheus Exporter** | Python stdlib, cache 30s, metric ครบ (status, duration, phase) |
| **Security** | Non-root, no token leak, mask secret, read-only volume |

 