# 🚀 ชุดเครื่องมืออัตโนมัติ 3 อย่าง

จัดโครงสร้าง:
```
project/
├── scripts/
│   ├── Run-6PhaseWorkflow.ps1        # orchestrator หลัก
│   ├── Invoke-Phase.ps1              # helper เรียกแต่ละเฟส
│   ├── Publish-JiraComment.ps1       # helper โพสต์ Jira
│   └── config/
│       └── workflow.config.json
├── .github/
│   ├── workflows/
│   │   ├── 6phase-on-label.yml
│   │   ├── security-review.yml
│   │   └── performance-review.yml
│   ├── ISSUE_TEMPLATE/
│   │   ├── feature.md
│   │   └── bug.md
│   ├── pull_request_template.md
│   └── copilot-instructions.md
└── docs/opencode/                    # skills (จากรอบก่อน)
```

---

## 📜 1) PowerShell: Auto-run 6 เฟส + Jira

### 1.1 `scripts/config/workflow.config.json`

```json
{
  "phases": [
    {
      "id": "gather",
      "name": "Gather",
      "emoji": "📥",
      "skill": "gen-gather",
      "required": true
    },
    {
      "id": "plan",
      "name": "Plan",
      "emoji": "🗺️",
      "skill": "gen-plan",
      "required": true,
      "dependsOn": ["gather"]
    },
    {
      "id": "execute",
      "name": "Execute",
      "emoji": "⚙️",
      "skill": "gen-execute",
      "required": true,
      "dependsOn": ["plan"],
      "mode": "read-write"
    },
    {
      "id": "security",
      "name": "Security Review",
      "emoji": "🔒",
      "skill": "gen-security-review",
      "required": true,
      "dependsOn": ["execute"],
      "failFast": true
    },
    {
      "id": "performance",
      "name": "Performance Review",
      "emoji": "⚡",
      "skill": "gen-performance-review",
      "required": true,
      "dependsOn": ["execute"],
      "failFast": true
    },
    {
      "id": "rca",
      "name": "Root Cause Analysis",
      "emoji": "🔍",
      "skill": "gen-rca",
      "required": false,
      "dependsOn": ["security", "performance"]
    }
  ],
  "jira": {
    "enabled": true,
    "commentTemplate": "h3. {emoji} Phase: {name}\n*Ticket:* {ticketKey}\n*Status:* {status}\n*เวลา:* {timestamp}\n\n{report}",
    "attachReports": true,
    "transitionOnComplete": null
  },
  "output": {
    "reportDir": "./.workflow-reports",
    "format": "markdown"
  },
  "opencode": {
    "command": "opencode",
    "args": ["run", "--model", "big-pickle"],
    "timeoutSec": 600
  }
}
```

---

### 1.2 `scripts/Invoke-Phase.ps1`

```powershell
<#
.SYNOPSIS
  เรียก agent ให้ทำงาน 1 เฟส แล้วคืนผลลัพธ์

.PARAMETER PhaseId
  รหัสเฟส (gather, plan, execute, security, performance, rca)

.PARAMETER TicketKey
  รหัส Jira ticket

.PARAMETER Task
  ชื่องาน

.PARAMETER TaskDescription
  รายละเอียดงาน

.PARAMETER ContextJson
  JSON ของ output จากเฟสก่อนหน้า (optional)

.PARAMETER ConfigPath
  path ไปยัง workflow.config.json

.EXAMPLE
  .\Invoke-Phase.ps1 -PhaseId gather -TicketKey ABC-123 -Task "..." -TaskDescription "..."
#>
[CmdletBinding()]
param(
  [Parameter(Mandatory)][string]$PhaseId,
  [Parameter(Mandatory)][string]$TicketKey,
  [Parameter(Mandatory)][string]$Task,
  [Parameter(Mandatory)][string]$TaskDescription,
  [string]$ContextJson = '{}',
  [string]$ConfigPath = "$PSScriptRoot/config/workflow.config.json"
)

$ErrorActionPreference = 'Stop'
Set-StrictMode -Version Latest

# --- โหลด config ---
if (-not (Test-Path $ConfigPath)) {
  throw "Config not found: $ConfigPath"
}
$config = Get-Content $ConfigPath -Raw | ConvertFrom-Json
$phase = $config.phases | Where-Object { $_.id -eq $PhaseId }
if (-not $phase) {
  throw "Unknown phase: $PhaseId"
}

# --- เตรียม prompt ---
$skillFile = Join-Path $PSScriptRoot "..\docs\opencode\phases\$($phase.skill).md"
if (-not (Test-Path $skillFile)) {
  throw "Skill file not found: $skillFile"
}

$promptTemplate = @"
# สั่งงาน: $($phase.skill)

## Context จากเฟสก่อนหน้า
$ContextJson

## งาน
- Ticket: $TicketKey
- Task: $Task
- Description: $TaskDescription
- Mode: $($phase.mode ?? 'read-only')

กรุณาทำตาม skill นี้และตอบเป็น Markdown ตาม template เท่านั้น
"@

# --- บันทึก prompt ---
$reportDir = $config.output.reportDir
if (-not (Test-Path $reportDir)) {
  New-Item -ItemType Directory -Path $reportDir -Force | Out-Null
}
$promptFile = Join-Path $reportDir "$TicketKey-$PhaseId-prompt.md"
Set-Content -Path $promptFile -Value $promptTemplate -Encoding UTF8

# --- เรียก opencode ---
Write-Host "▶ [$($phase.emoji)] Running phase: $($phase.name)" -ForegroundColor Cyan
$startTime = Get-Date

$tempOut = Join-Path $reportDir "$TicketKey-$PhaseId-output.md"
try {
  $args = @(
    'run',
    '--model', 'big-pickle',
    '--prompt-file', $promptFile,
    '--output', $tempOut
  )

  $proc = Start-Process -FilePath $config.opencode.command `
    -ArgumentList $args `
    -NoNewWindow -PassThru -Wait `
    -RedirectStandardError (Join-Path $reportDir "$TicketKey-$PhaseId.err")

  if ($proc.ExitCode -ne 0) {
    throw "Phase '$PhaseId' failed with exit code $($proc.ExitCode)"
  }

  if (-not (Test-Path $tempOut)) {
    throw "Phase '$PhaseId' produced no output"
  }

  $report = Get-Content $tempOut -Raw
  $duration = (Get-Date) - $startTime

  # --- วิเคราะห์ status ---
  $status = if ($report -match '🔴|fail|❌') { 'FAIL' }
            elseif ($report -match '🟡|blocked') { 'WARN' }
            else { 'PASS' }

  $result = [PSCustomObject]@{
    phaseId    = $PhaseId
    phaseName  = $phase.name
    emoji      = $phase.emoji
    status     = $status
    durationMs = [int]$duration.TotalMilliseconds
    report     = $report
    reportPath = $tempOut
    startedAt  = $startTime.ToString('o')
    endedAt    = (Get-Date).ToString('o')
    error      = $null
  }

  # --- บันทึก JSON ---
  $jsonPath = Join-Path $reportDir "$TicketKey-$PhaseId-result.json"
  $result | ConvertTo-Json -Depth 10 | Set-Content -Path $jsonPath -Encoding UTF8

  Write-Host "  ✅ $($phase.name) → $status ($([int]$duration.TotalSeconds)s)" `
    -ForegroundColor $(if ($status -eq 'PASS') { 'Green' } else { 'Yellow' })

  return $result
}
catch {
  $duration = (Get-Date) - $startTime
  $errResult = [PSCustomObject]@{
    phaseId    = $PhaseId
    phaseName  = $phase.name
    emoji      = $phase.emoji
    status     = 'ERROR'
    durationMs = [int]$duration.TotalMilliseconds
    report     = ''
    reportPath = $null
    startedAt  = $startTime.ToString('o')
    endedAt    = (Get-Date).ToString('o')
    error      = $_.Exception.Message
  }
  Write-Host "  ❌ $($phase.name) → ERROR: $($_.Exception.Message)" -ForegroundColor Red
  return $errResult
}
```

---

### 1.3 `scripts/Publish-JiraComment.ps1`

```powershell
<#
.SYNOPSIS
  โพสต์ comment ไปยัง Jira ticket (REST API v3)

.PARAMETER TicketKey
.PARAMETER Emoji
.PARAMETER PhaseName
.PARAMETER Status
.PARAMETER Report
.PARAMETER AttachPath
  path ของไฟล์ที่จะแนบ (optional)

.NOTES
  ต้องตั้ง env: JIRA_BASE_URL, JIRA_EMAIL, JIRA_API_TOKEN
#>
[CmdletBinding()]
param(
  [Parameter(Mandatory)][string]$TicketKey,
  [Parameter(Mandatory)][string]$Emoji,
  [Parameter(Mandatory)][string]$PhaseName,
  [Parameter(Mandatory)][string]$Status,
  [Parameter(Mandatory)][string]$Report,
  [string]$AttachPath = $null
)

$ErrorActionPreference = 'Stop'
Set-StrictMode -Version Latest

# --- ตรวจ env ---
if (-not $env:JIRA_BASE_URL)    { throw 'JIRA_BASE_URL not set' }
if (-not $env:JIRA_EMAIL)       { throw 'JIRA_EMAIL not set' }
if (-not $env:JIRA_API_TOKEN)   { throw 'JIRA_API_TOKEN not set' }

$auth = [Convert]::ToBase64String(
  [Text.Encoding]::ASCII.GetBytes("$($env:JIRA_EMAIL):$($env:JIRA_API_TOKEN)")
)
$headers = @{
  Authorization = "Basic $auth"
  'Content-Type' = 'application/json'
}

# --- ตรวจ token ---
try {
  $me = Invoke-RestMethod -Uri "$($env:JIRA_BASE_URL)/rest/api/3/myself" `
    -Headers $headers -Method Get
  Write-Verbose "Jira auth OK: $($me.emailAddress)"
} catch {
  throw "Jira auth failed: $($_.Exception.Message)"
}

# --- สร้าง ADF (Atlassian Document Format) ---
$timestamp = (Get-Date).ToString('yyyy-MM-dd HH:mm:ss')

# แปลง Markdown → simple text (Jira ADF)
$reportLines = $Report -split "`n"
$content = @()

# Heading
$content += @{
  type = 'heading'
  attrs = @{ level = 3 }
  content = @(@{ type = 'text'; text = "$Emoji Phase: $PhaseName" })
}

# Metadata
$content += @{
  type = 'paragraph'
  content = @(
    @{ type = 'text'; text = 'Ticket: '; marks = @(@{ type = 'strong' }) }
    @{ type = 'text'; text = $TicketKey }
    @{ type = 'text'; text = '  |  Status: '; marks = @(@{ type = 'strong' }) }
    @{ type = 'text'; text = $Status }
    @{ type = 'text'; text = '  |  Time: '; marks = @(@{ type = 'strong' }) }
    @{ type = 'text'; text = $timestamp }
  )
}

# Report body
$content += @{
  type = 'codeBlock'
  attrs = @{ language = 'markdown' }
  content = @(@{ type = 'text'; text = $Report })
}

$payload = @{
  body = @{
    type = 'doc'
    version = 1
    content = $content
  }
} | ConvertTo-Json -Depth 20

# --- โพสต์ comment ---
$commentUrl = "$($env:JIRA_BASE_URL)/rest/api/3/issue/$TicketKey/comment"
try {
  $resp = Invoke-RestMethod -Uri $commentUrl -Method Post -Headers $headers -Body $payload
  Write-Host "  💬 Jira comment posted (id=$($resp.id))" -ForegroundColor Green
} catch {
  Write-Warning "Jira comment failed: $($_.Exception.Message)"
  if ($_.ErrorDetails.Message) {
    Write-Warning $_.ErrorDetails.Message
  }
}

# --- แนบไฟล์ ---
if ($AttachPath -and (Test-Path $AttachPath)) {
  $attachHeaders = @{
    Authorization = "Basic $auth"
    'X-Atlassian-Token' = 'no-check'
  }
  $attachUrl = "$($env:JIRA_BASE_URL)/rest/api/3/issue/$TicketKey/attachments"
  try {
    $form = @{ file = Get-Item $AttachPath }
    $r = Invoke-RestMethod -Uri $attachUrl -Method Post -Headers $attachHeaders -Form $form
    Write-Host "  📎 Attached: $($AttachPath)" -ForegroundColor Green
  } catch {
    Write-Warning "Jira attach failed: $($_.Exception.Message)"
  }
}
```

---

### 1.4 `scripts/Run-6PhaseWorkflow.ps1` (orchestrator หลัก)

```powershell
<#
.SYNOPSIS
  รัน 6 เฟสอัตโนมัติ แล้วโพสต์ผลไป Jira

.DESCRIPTION
  Orchestrator ที่:
  1. รันเฟสตามลำดับ (gather → plan → execute → security → performance → rca)
  2. ส่ง context จากเฟสก่อนหน้าไปยังเฟสถัดไป
  3. ถ้า security/performance fail → abort (failFast)
  4. โพสต์ผลแต่ละเฟสไป Jira
  5. สรุปผลรวม

.PARAMETER TicketKey
.PARAMETER Task
.PARAMETER TaskDescription
.PARAMETER Mode
  read-only | read-write (default: read-only)
.PARAMETER Phases
  รันเฉพาะบางเฟส เช่น -Phases gather,plan
.PARAMETER SkipJira
  ข้ามการโพสต์ Jira
.PARAMETER WhatIf

.EXAMPLE
  .\Run-6PhaseWorkflow.ps1 `
    -TicketKey "ABC-123" `
    -Task "เพิ่ม API ค้นหาสินค้า" `
    -TaskDescription "GET /api/v1/products?category=..." `
    -Mode read-write
#>
[CmdletBinding(SupportsShouldProcess)]
param(
  [Parameter(Mandatory)][string]$TicketKey,
  [Parameter(Mandatory)][string]$Task,
  [Parameter(Mandatory)][string]$TaskDescription,
  [ValidateSet('read-only','read-write')]
  [string]$Mode = 'read-only',
  [string[]]$Phases = @(),
  [switch]$SkipJira,
  [string]$ConfigPath = "$PSScriptRoot/config/workflow.config.json"
)

$ErrorActionPreference = 'Stop'
Set-StrictMode -Version Latest
$startTime = Get-Date

Write-Host ""
Write-Host "╔══════════════════════════════════════════════════════════╗" -ForegroundColor Cyan
Write-Host "║  🚀 6-Phase Workflow Runner                              ║" -ForegroundColor Cyan
Write-Host "╚══════════════════════════════════════════════════════════╝" -ForegroundColor Cyan
Write-Host "  Ticket: $TicketKey"
Write-Host "  Task:   $Task"
Write-Host "  Mode:   $Mode"
Write-Host ""

# --- โหลด config ---
$config = Get-Content $ConfigPath -Raw | ConvertFrom-Json

# --- เลือกเฟสที่จะรัน ---
$phasesToRun = if ($Phases.Count -gt 0) {
  $config.phases | Where-Object { $_.id -in $Phases }
} else {
  $config.phases
}

if (-not $phasesToRun) {
  throw "No phases selected"
}

# --- ตรวจ dependency ---
foreach ($p in $phasesToRun) {
  if ($p.dependsOn) {
    foreach ($dep in $p.dependsOn) {
      if ($dep -notin $phasesToRun.id) {
        throw "Phase '$($p.id)' depends on '$dep' which is not in the run list"
      }
      # ตรวจว่า dependency มาก่อนหน้า
      $depIdx  = [array]::IndexOf($phasesToRun.id, $dep)
      $currIdx = [array]::IndexOf($phasesToRun.id, $p.id)
      if ($depIdx -gt $currIdx) {
        throw "Phase order invalid: '$dep' must come before '$($p.id)'"
      }
    }
  }
}

# --- Results ---
$results = [ordered]@{}
$contextJson = @{
  ticketKey = $TicketKey
  task      = $Task
} | ConvertTo-Json -Compress

# --- Loop เฟส ---
foreach ($phase in $phasesToRun) {
  Write-Host "────────────────────────────────────────────────────────" -ForegroundColor DarkGray

  if ($phase.mode -eq 'read-write' -and $Mode -eq 'read-only') {
    Write-Host "  ⏭️  Skip $($phase.name) (mode=read-only)" -ForegroundColor DarkYellow
    continue
  }

  if ($PSCmdlet.ShouldProcess($phase.name, 'Run phase')) {
    $result = & "$PSScriptRoot/Invoke-Phase.ps1" `
      -PhaseId $phase.id `
      -TicketKey $TicketKey `
      -Task $Task `
      -TaskDescription $TaskDescription `
      -ContextJson $contextJson `
      -ConfigPath $ConfigPath

    $results[$phase.id] = $result

    # --- โพสต์ Jira ---
    if (-not $SkipJira -and $config.jira.enabled) {
      if ($PSCmdlet.ShouldProcess($TicketKey, 'Post Jira comment')) {
        & "$PSScriptRoot/Publish-JiraComment.ps1" `
          -TicketKey $TicketKey `
          -Emoji $result.emoji `
          -PhaseName $result.phaseName `
          -Status $result.status `
          -Report $result.report `
          -AttachPath $result.reportPath
      }
    }

    # --- FailFast ---
    if ($phase.failFast -and $result.status -in @('FAIL','ERROR')) {
      Write-Host ""
      Write-Host "  🛑 FailFast triggered at '$($phase.name)'" -ForegroundColor Red
      Write-Host "     หยุด workflow — ต้องแก้ไขก่อนดำเนินการต่อ" -ForegroundColor Red
      break
    }

    # --- ส่ง context ต่อไป ---
    $contextJson = @{
      ticketKey      = $TicketKey
      task           = $Task
      previousPhases = $results
    } | ConvertTo-Json -Depth 10 -Compress
  }
}

# --- สรุปผลรวม ---
$duration = (Get-Date) - $startTime
Write-Host ""
Write-Host "╔══════════════════════════════════════════════════════════╗" -ForegroundColor Cyan
Write-Host "║  📊 สรุปผลรวม                                            ║" -ForegroundColor Cyan
Write-Host "╚══════════════════════════════════════════════════════════╝" -ForegroundColor Cyan

$summary = foreach ($kv in $results.GetEnumerator()) {
  [PSCustomObject]@{
    Phase    = "$($kv.Value.emoji) $($kv.Value.phaseName)"
    Status   = $kv.Value.status
    Duration = "$([math]::Round($kv.Value.durationMs/1000, 1))s"
  }
}
$summary | Format-Table -AutoSize

$overallStatus = if ($results.Values.status -contains 'ERROR') { '🔴 ERROR' }
                 elseif ($results.Values.status -contains 'FAIL') { '🟠 FAIL' }
                 elseif ($results.Values.status -contains 'WARN') { '🟡 WARN' }
                 else { '🟢 PASS' }

Write-Host "  Overall:  $overallStatus"
Write-Host "  Total:    $([math]::Round($duration.TotalSeconds, 1))s"
Write-Host "  Reports:  $($config.output.reportDir)"
Write-Host ""

# --- บันทึก summary ---
$summaryPath = Join-Path $config.output.reportDir "$TicketKey-summary.json"
@{
  ticketKey   = $TicketKey
  task        = $Task
  overall     = $overallStatus
  durationMs  = [int]$duration.TotalMilliseconds
  results     = $results
  finishedAt  = (Get-Date).ToString('o')
} | ConvertTo-Json -Depth 20 | Set-Content $summaryPath -Encoding UTF8

# --- โพสต์สรุปรวมไป Jira ---
if (-not $SkipJira -and $config.jira.enabled) {
  $summaryMd = ($summary | ForEach-Object {
    "| $($_.Phase) | $($_.Status) | $($_.Duration) |"
  }) -join "`n"

  $finalReport = @"
# 📊 สรุปผลรวม — $TicketKey

| Phase | Status | Duration |
|-------|--------|----------|
$summaryMd

**Overall:** $overallStatus
**Total:** $([math]::Round($duration.TotalSeconds, 1))s
"@

  & "$PSScriptRoot/Publish-JiraComment.ps1" `
    -TicketKey $TicketKey `
    -Emoji '📊' `
    -PhaseName 'Summary' `
    -Status $overallStatus `
    -Report $finalReport
}

# --- Exit code ---
if ($overallStatus -match 'ERROR|FAIL') { exit 1 } else { exit 0 }
```

---

### 1.5 การใช้งาน

```powershell
# ตั้งค่า env (ครั้งเดียว)
$env:JIRA_BASE_URL   = 'https://your-domain.atlassian.net'
$env:JIRA_EMAIL      = 'you@example.com'
$env:JIRA_API_TOKEN  = 'xxxxx'   # https://id.atlassian.com/manage-profile/security/api-tokens

# รัน 6 เฟสครบ
.\scripts\Run-6PhaseWorkflow.ps1 `
  -TicketKey "ABC-123" `
  -Task "เพิ่ม API ค้นหาสินค้า" `
  -TaskDescription "GET /api/v1/products?category=1&page=1&size=20" `
  -Mode read-write

# รันเฉพาะ gather + plan
.\scripts\Run-6PhaseWorkflow.ps1 `
  -TicketKey "ABC-123" `
  -Task "..." `
  -TaskDescription "..." `
  -Phases gather,plan

# Dry-run (ไม่โพสต์ Jira)
.\scripts\Run-6PhaseWorkflow.ps1 -TicketKey "ABC-123" -Task "..." `
  -TaskDescription "..." -SkipJira

# WhatIf (ไม่รันจริง)
.\scripts\Run-6PhaseWorkflow.ps1 -TicketKey "ABC-123" -Task "..." `
  -TaskDescription "..." -WhatIf
```

---

## 🐙 2) GitHub Actions: Trigger sub-skill ตาม label

### 2.1 `.github/workflows/6phase-on-label.yml`

```yaml
name: 6-Phase Workflow

on:
  issues:
    types: [labeled, unlabeled]
  workflow_dispatch:
    inputs:
      ticket_key:
        description: 'Ticket Key (เช่น ABC-123)'
        required: true
      phases:
        description: 'เฟส (comma-separated, เว้นว่าง = ทั้งหมด)'
        required: false
        default: ''

permissions:
  contents: write
  issues: write
  pull-requests: write
  id-token: write

concurrency:
  group: 6phase-${{ github.event.issue.number || github.event.inputs.ticket_key }}
  cancel-in-progress: false

jobs:
  detect:
    name: 🔍 Detect phases from labels
    runs-on: ubuntu-latest
    outputs:
      phases: ${{ steps.detect.outputs.phases }}
      ticket: ${{ steps.detect.outputs.ticket }}
      should_run: ${{ steps.detect.outputs.should_run }}
    steps:
      - name: Detect labels → phases
        id: detect
        uses: actions/github-script@v7
        with:
          script: |
            const labelMap = {
              'phase:gather':      'gather',
              'phase:plan':        'plan',
              'phase:execute':     'execute',
              'phase:security':    'security',
              'phase:performance': 'performance',
              'phase:rca':         'rca',
            };

            // manual dispatch
            const manualTicket = '${{ github.event.inputs.ticket_key }}';
            const manualPhases = '${{ github.event.inputs.phases }}';

            if (manualTicket) {
              core.setOutput('ticket', manualTicket);
              core.setOutput('phases', manualPhases || 'gather,plan,execute,security,performance,rca');
              core.setOutput('should_run', 'true');
              return;
            }

            const labels = context.payload.issue.labels.map(l => l.name);
            const phases = labels
              .filter(l => labelMap[l])
              .map(l => labelMap[l]);

            core.setOutput('ticket', context.payload.issue.title.match(/[A-Z]+-\d+/)?.[0] || `#${context.payload.issue.number}`);
            core.setOutput('phases', phases.join(','));
            core.setOutput('should_run', phases.length > 0 ? 'true' : 'false');

  run:
    name: 🚀 Run phases
    needs: detect
    if: needs.detect.outputs.should_run == 'true'
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
          cache: true

      - name: Setup PowerShell
        shell: pwsh
        run: |
          $PSVersionTable
          Install-Module -Name Pester -Force -SkipPublisherCheck -Scope CurrentUser -ErrorAction SilentlyContinue

      - name: Install opencode CLI
        run: |
          # ตัวอย่าง — ปรับตามช่องทางติดตั้งจริง
          npm install -g opencode || true

      - name: Run 6-phase workflow
        shell: pwsh
        env:
          JIRA_BASE_URL:    ${{ secrets.JIRA_BASE_URL }}
          JIRA_EMAIL:       ${{ secrets.JIRA_EMAIL }}
          JIRA_API_TOKEN:   ${{ secrets.JIRA_API_TOKEN }}
          OPENAI_API_KEY:   ${{ secrets.OPENAI_API_KEY }}
          TICKET_KEY:       ${{ needs.detect.outputs.ticket }}
          PHASES:           ${{ needs.detect.outputs.phases }}
        run: |
          $phases = $env:PHASES -split ',' | Where-Object { $_ }
          ./scripts/Run-6PhaseWorkflow.ps1 `
            -TicketKey $env:TICKET_KEY `
            -Task "${{ github.event.issue.title }}" `
            -TaskDescription "${{ github.event.issue.body }}" `
            -Phases $phases `
            -Mode read-write

      - name: Upload reports
        if: always()
        uses: actions/upload-artifact@v4
        with:
          name: workflow-reports-${{ needs.detect.outputs.ticket }}
          path: .workflow-reports/
          retention-days: 30

      - name: Comment on issue
        if: always() && github.event_name == 'issues'
        uses: actions/github-script@v7
        with:
          script: |
            const fs = require('fs');
            const ticket = '${{ needs.detect.outputs.ticket }}';
            const summaryPath = `.workflow-reports/${ticket}-summary.json`;
            let summary = '## 🚀 6-Phase Workflow Result\n\n';
            if (fs.existsSync(summaryPath)) {
              const data = JSON.parse(fs.readFileSync(summaryPath, 'utf8'));
              summary += `**Overall:** ${data.overall}\n`;
              summary += `**Duration:** ${Math.round(data.durationMs/1000)}s\n\n`;
              for (const [id, r] of Object.entries(data.results)) {
                summary += `- ${r.emoji} **${r.phaseName}**: ${r.status}\n`;
              }
            } else {
              summary += '⚠️ ไม่พบ summary';
            }
            await github.rest.issues.createComment({
              owner: context.repo.owner,
              repo: context.repo.repo,
              issue_number: context.issue.number,
              body: summary,
            });
```

---

### 2.2 `.github/workflows/security-review.yml`

```yaml
name: Security Review

on:
  pull_request:
    types: [opened, synchronize, reopened, labeled]
  schedule:
    - cron: '0 2 * * 1'  # ทุกจันทร์ 02:00 UTC

permissions:
  contents: read
  security-events: write
  pull-requests: write

jobs:
  govulncheck:
    name: 🛡️ govulncheck
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
      - name: Run govulncheck
        uses: golang/govulncheck-action@v1
        with:
          output-format: sarif
          output-file: govulncheck.sarif
      - name: Upload SARIF
        if: always()
        uses: github/codeql-action/upload-sarif@v3
        with:
          sarif_file: govulncheck.sarif

  gosec:
    name: 🔒 gosec
    runs-on: ubuntu-latest
    steps
      - uses: actions/checkout@v4
      - uses: securego/gosec@master
        with:
          args: '-no-fail -fmt sarif -out gosec.sarif ./...'
      - uses: github/codeql-action/upload-sarif@v3
        with:
          sarif_file: gosec.sarif

  gitleaks:
    name: 🔑 Secret scan
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      - uses: gitleaks/gitleaks-action@v2
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}

  go-race:
    name: 🏁 Race detector
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
      - run: go test -race ./... -count=1

  semgrep:
    name: 🕵️ Semgrep OWASP
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: returntocorp/semgrep-action@v1
        with:
          config: >-
            p/owasp-top-ten
            p/golang
```

---

### 2.3 `.github/workflows/performance-review.yml`

```yaml
name: Performance Review

on:
  pull_request:
    types: [opened, synchronize, labeled]
  schedule:
    - cron: '0 3 * * 2'  # ทุกอังคาร 03:00 UTC

permissions:
  contents: read
  pull-requests: write

jobs:
  benchmark:
    name: ⚡ Go benchmark
    if: |
      github.event_name == 'schedule' ||
      contains(github.event.pull_request.labels.*.name, 'phase:performance') ||
      contains(github.event.pull_request.labels.*.name, 'performance')
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod

      - name: Run benchmarks
        run: |
          go test -bench=. -benchmem -count=5 ./... | tee bench.txt

      - name: Install benchstat
        run: go install golang.org/x/perf/cmd/benchstat@latest

      - name: Compare with baseline
        run: |
          # baseline เก็บใน branch 'bench-baseline'
          if git ls-remote --heads origin bench-baseline | grep -q bench-baseline; then
            git fetch origin bench-baseline
            git show origin/bench-baseline:bench.txt > bench-baseline.txt || true
            benchstat bench-baseline.txt bench.txt | tee benchstat.txt
          else
            echo "No baseline yet — establishing baseline"
            cp bench.txt bench-baseline.txt
            benchstat bench.txt | tee benchstat.txt
          fi

      - name: Upload benchmark artifacts
        uses: actions/upload-artifact@v4
        with:
          name: benchmark-results
          path: |
            bench.txt
            benchstat.txt

      - name: Comment on PR
        if: github.event_name == 'pull_request'
        uses: actions/github-script@v7
        with:
          script: |
            const fs = require('fs');
            const body = fs.existsSync('benchstat.txt')
              ? fs.readFileSync('benchstat.txt', 'utf8').slice(0, 60000)
              : 'No results';
            await github.rest.issues.createComment({
              owner: context.repo.owner,
              repo: context.repo.repo,
              issue_number: context.issue.number,
              body: `## ⚡ Benchmark Results\n\`\`\`\n${body}\n\`\`\``,
            });
```

---

### 2.4 Labels ที่ต้องสร้างใน repo

```bash
# รันคำสั่งนี้ครั้งเดียว (gh cli)
gh label create "phase:gather"      --color "0E8A16" --description "Run phase 1: Gather"
gh label create "phase:plan"        --color "1D76DB" --description "Run phase 2: Plan"
gh label create "phase:execute"     --color "5319E7" --description "Run phase 3: Execute"
gh label create "phase:security"    --color "B60205" --description "Run phase 4: Security Review"
gh label create "phase:performance" --color "FBCA04" --description "Run phase 5: Performance Review"
gh label create "phase:rca"         --color "C2E0C6" --description "Run phase 6: Root Cause Analysis"
gh label create "performance"       --color "FBCA04" --description "Performance-sensitive change"
gh label create "security"          --color "B60205" --description "Security-sensitive change"
```

---

## 🤖 3) `.github/copilot-instructions.md`

```markdown
# 🤖 Copilot Custom Instructions — icmongolang

> ไฟล์นี้ Copilot จะอ่านอัตโนมัติทุกครั้งในโปรเจกต์นี้
> ใช้เป็น baseline สำหรับทุกการสนทนา

---

## 🎯 ตัวตนของ Agent

- **บทบาท:** Senior Go Developer + Solution Architect + QA Lead + SRE
- **ภาษา:** ตอบเป็น **ไทย** (คำศัพท์เทคนิคเป็นอังกฤษได้)
- **โหมดเริ่มต้น:** `read-only` — ต้องขออนุญาตก่อนแก้ไฟล์
- **โมเดล:** ใช้ตาม config ของ opencode (big-pickle)

---

## 🏗️ โครงสร้างโปรเจกต์

โปรเจกต์นี้เป็น **Go backend** ชื่อ `icmongolang` ใช้โครงสร้างมาตรฐาน:

```
icmongolang/
├── main.go                # entrypoint
├── go.mod / go.sum
├── cmd/                   # CLI commands / entry points
├── internal/              # โค้ดภายใน (ไม่ export)
│   ├── <module>/
│   │   ├── handler.go     # HTTP/gRPC handler
│   │   ├── service.go     # business logic
│   │   ├── repository.go  # data access
│   │   ├── model.go       # DTO / entity
│   │   └── *_test.go      # unit test
│   └── shared/            # shared utilities
├── pkg/                   # โค้ดที่ export ให้ภายนอกใช้
├── api/                   # swagger / openapi spec
├── docs/opencode/         # skill files
└── scripts/               # automation scripts
```

**ก่อนแก้ไขใด ๆ ต้องอ่าน:** `main.go`, `go.mod`, `cmd/`, `internal/`, `pkg/` เพื่อเข้าใจบริบท

---

## ✅ หลักการทำงาน (Principles)

### 1. Reuse-First
- อ่านโครงสร้างเดิมก่อน แล้วนำฟังก์ชัน/module ที่มีอยู่มาใช้ซ้ำ
- พยายาม **ไม่แก้โครงสร้างเดิม**
- ถ้ามี Function เดิมที่ใช้ร่วมกัน → นำมาใช้
- ถ้ามี Function ใหม่ที่ Performance ดีกว่าเดิม → เสนอทำใหม่ พร้อมเหตุผล + benchmark

### 2. Performance-First
- ออกแบบให้แอป **รวดเร็ว เสถียร มีประสิทธิภาพ**
- ระบุ index / caching / concurrency ทุกครั้ง
- วัดผลได้ด้วย `go test -bench=. -benchmem`
- เปรียบเทียบก่อน/หลังเสมอ

### 3. Test-Driven Development (TDD)
- **Red → Green → Refactor** บังคับ
- เขียน test ก่อน implement
- ใช้ `mockery` สำหรับ mock interface
- Coverage เป้าหมาย ≥ 80%

### 4. Zero-Panic Mindset (Go)
- ❌ ห้ามใช้ `panic()` ใน production path → return error เสมอ
- ✅ ตรวจ `nil` ก่อน dereference pointer
- ✅ ตรวจ `err != nil` ทุกครั้ง
- ✅ ใช้ `v, ok := m[k]` เมื่อไม่มั่นใจ
- ✅ ตรวจ `len(slice) > 0` ก่อน `slice[0]`
- ✅ ใช้ `v, ok := x.(T)` สำหรับ type assertion
- ✅ ใช้ `context.WithTimeout` กัน goroutine leak
- ✅ ระวัง concurrent map write — ใช้ `sync.RWMutex` / `sync.Map`
- ✅ ระวัง division by zero / integer overflow

---

## 📋 Workflow 6 เฟส (บังคับสำหรับงานใหญ่)

ทุกงานที่ซับซ้อน ต้องดำเนินการตาม 6 เฟส:

1. **📥 Gather** — รวบรวมข้อมูลทั้งหมด (read-only)
2. **🗺️ Plan** — วางแผน แตก task ระบุความเสี่ยง
3. **⚙️ Execute** — ลงมือทำตามแผน TDD
4. **🔒 Security Review** — ตรวจ OWASP Top 10 + Go-specific
5. **⚡ Performance Review** — benchmark + pprof + load test
6. **🔍 Root Cause Analysis** — 5 Whys + Preventive Action

เรียกใช้ผ่าน Copilot Chat:
- `#gen-6phase-workflow` — ทั้งหมด
- `#gen-gather`, `#gen-plan`, `#gen-execute` — เฉพาะเฟส
- `#gen-security-review`, `#gen-performance-review`, `#gen-rca`

---

## ⛔ ข้อห้าม (Hard Rules)

1. **ห้าม commit / push** เว้นแต่ผู้ใช้สั่งชัดเจน
2. **ห้ามเดา** — ไม่มีข้อมูลให้ใส่ `-` + `[ต้องการข้อมูลเพิ่มเติม]`
3. **ห้ามแตะโค้ด** ในโหมด read-only
4. **ห้ามใช้ `panic()`** ใน production path
5. **ห้าม silent fix** — ทุกการแก้ต้อง log + ระบุเหตุผล
6. **ห้าม hardcode secret** — ใช้ env / vault
7. **ห้าม log PII / token / password**
8. **ห้ามแก้โครงสร้างเดิม** โดยไม่ระบุ impact analysis

---

## 🎨 รูปแบบการตอบ

### ใช้ Markdown เสมอ
- หัวข้อ: `#`, `##`, `###`
- Checklist: `- [ ]`
- ตาราง: `| col | col |`
- Diagram: **Mermaid** (` ```mermaid `)
- Code block: ระบุภาษา (` ```go `, ` ```powershell `, ` ```yaml `)

### ใช้ Emoji ตามบริบท
| ใช้กับ | Emoji |
|--------|-------|
| เฟส Gather | 📥 |
| เฟส Plan | 🗺️ |
| เฟส Execute | ⚙️ |
| เฟส Security | 🔒 |
| เฟส Performance | ⚡ |
| เฟส RCA | 🔍 |
| สำเร็จ | ✅ |
| ล้มเหลว | ❌ |
| เตือน | ⚠️ |
| ต้องข้อมูลเพิ่ม | ❓ |
| Root Cause | 🎯 |

### โครงสร้างคำตอบ
1. **สรุปสั้น** (TL;DR) — 2-3 บรรทัด
2. **รายละเอียด** — ตาม template
3. **Next step** — สิ่งที่ต้องทำต่อ

---

## 🧪 มาตรฐานการทดสอบ

### ก่อน commit ทุกครั้ง
```bash
go test ./... -v -cover
go vet ./...
golangci-lint run
go test -race ./...
govulncheck ./...
```

### Coverage
- Business logic: ≥ 80%
- Handler: ≥ 70%
- Repository: ≥ 60%

### Mock
- ใช้ `mockery` generate mock จาก interface
- ห้าม mock concrete type

---

## 📘 Backend / API

### ทุก endpoint ต้องมี
- [ ] Handler + Service + Repository แยกชั้น
- [ ] Input validation
- [ ] Error handling (ไม่ leak internals)
- [ ] Swagger annotation (`api/<module>/swagger.yaml`)
- [ ] Postman collection (`postman/<module>.postman_collection.json`)
- [ ] Unit test + integration test

### Swagger / OpenAPI
- ใช้ OpenAPI 3.0
- ระบุ schema ครบ (request, response, error)
- ตัวอย่าง request/response ทุก endpoint

### Postman
- Collection แยกตาม module
- Environment variables: `baseUrl`, `token`
- มี test script ตรวจ status + schema

---

## 🎨 Frontend (ถ้ามี)

- ต้องมี **คู่มือการใช้งาน** (installation, config, run, usage)
- ระบุ component structure
- ระบุ error / loading / empty state
- Accessibility checklist (WCAG AA)

---

## 📊 การรายงานผล

ทุกงานต้องมีรายงานสรุป:

| หัวข้อ | รายละเอียด |
|--------|-----------|
| งานที่ทำเสร็จ | - |
| ไฟล์ที่สร้างใหม่ | - |
| ไฟล์ที่แก้ไข | - |
| Test ที่เพิ่ม | - |
| Coverage | - |
| Benchmark ผลลัพธ์ | - |
| Security issues | - |
| Performance เทียบเป้า | - |
| Root Cause | - |
| ปัญหาที่พบ | - |
| งานที่ยังค้าง | - |
| Next step | - |

---

## 🔗 Skill Files ในโปรเจกต์

อ่าน skill เพิ่มเติมได้ที่:
- `docs/opencode/gen-6phase-workflow.md` — orchestrator
- `docs/opencode/phases/gen-*.md` — รายเฟส
- `docs/opencode/integrations/gen-jira-attach.md` — Jira integration
- `docs/opencode/templates/` — GitHub templates

---

## 🚀 Automation

- PowerShell orchestrator: `scripts/Run-6PhaseWorkflow.ps1`
- GitHub Actions trigger ตาม label `phase:*`
- Jira auto-comment ทุกเฟส (ถ้าตั้ง env)

### env vars ที่ต้องตั้ง
```bash
JIRA_BASE_URL=https://your-domain.atlassian.net
JIRA_EMAIL=you@example.com
JIRA_API_TOKEN=xxxxx
OPENAI_API_KEY=xxxxx
```

---

## 🎯 Definition of Done (DoD)

งานจะถือว่าเสร็จเมื่อ:
- [ ] Unit test ผ่าน + coverage ≥ เป้า
- [ ] Security review ผ่าน (ไม่มี Critical/High)
- [ ] Performance ผ่าน SLO
- [ ] ไม่มี Root Cause ค้าง
- [ ] Swagger อัปเดต
- [ ] Postman อัปเดต
- [ ] Report สรุปเสร็จ
- [ ] Reviewer approved

---

## 💬 ตัวอย่างการสนทนา

**User:** เพิ่ม API GET /products?category=1

**Copilot:**
1. 📥 **Gather** — ผมจะสำรวจ `internal/product/` ก่อน แล้วสรุปไฟล์ที่เกี่ยวข้อง
2. ถามคำถามที่จำเป็น (pagination? filter? auth?)
3. 🗺️ **Plan** — แตก task: model → repo → service → handler → test → swagger
4. ขออนุญาตก่อนแก้ไฟล์ (read-only default)
5. ⚙️ **Execute** ตาม TDD
6. 🔒 **Security** + ⚡ **Performance** review
7. 📊 **Report**

---

## 📌 หมายเหตุ

- ไฟล์นี้เป็น **living document** — อัปเดตได้ตามความเหมาะสม
- ถ้ามี conflict ระหว่าง instructions นี้กับ skill file → **skill file ชนะ** สำหรับงานเฉพาะเจาะจง
- ถ้าไม่แน่ใจ → **ถามผู้ใช้** ดีกว่าเดา

---

**Last updated:** 2026-09-16
**Version:** 1.0.0
```

---

## 📁 3.1 โครงสร้างไฟล์สมบูรณ์

```
project/
├── .github/
│   ├── workflows/
│   │   ├── 6phase-on-label.yml
│   │   ├── security-review.yml
│   │   └── performance-review.yml
│   ├── ISSUE_TEMPLATE/
│   │   ├── feature.md
│   │   └── bug.md
│   ├── pull_request_template.md
│   └── copilot-instructions.md          ← ⭐ ไฟล์นี้
├── scripts/
│   ├── Run-6PhaseWorkflow.ps1           ← orchestrator
│   ├── Invoke-Phase.ps1                 ← helper
│   ├── Publish-JiraComment.ps1          ← helper
│   └── config/
│       └── workflow.config.json
├── docs/opencode/
│   ├── gen-6phase-workflow.md
│   ├── phases/
│   │   ├── gen-gather.md
│   │   ├── gen-plan.md
│   │   ├── gen-execute.md
│   │   ├── gen-security-review.md
│   │   ├── gen-performance-review.md
│   │   └── gen-rca.md
│   ├── integrations/
│   │   └── gen-jira-attach.md
│   └── templates/
│       ├── github-issue-feature.md
│       ├── github-issue-bug.md
│       └── github-pr.md
├── main.go
├── go.mod
├── cmd/
├── internal/
├── pkg/
└── .workflow-reports/                    ← output (gitignore)
```

---

## 🔐 3.2 ตั้งค่า Secrets (GitHub)

```bash
# GitHub repo → Settings → Secrets and variables → Actions → New repository secret
gh secret set JIRA_BASE_URL   --body "https://your-domain.atlassian.net"
gh secret set JIRA_EMAIL      --body "you@example.com"
gh secret set JIRA_API_TOKEN  --body "xxxxx"
gh secret set OPENAI_API_KEY  --body "sk-xxxxx"
```

---

## 🚦 3.3 ตัวอย่างการใช้งานรวม

### Case 1: รัน local ทั้งหมด
```powershell
# 1. ตั้ง env
$env:JIRA_BASE_URL  = 'https://your-domain.atlassian.net'
$env:JIRA_EMAIL     = 'you@example.com'
$env:JIRA_API_TOKEN = 'xxxxx'

# 2. รัน 6 เฟส
.\scripts\Run-6PhaseWorkflow.ps1 `
  -TicketKey "ABC-123" `
  -Task "เพิ่ม API ค้นหาสินค้า" `
  -TaskDescription "GET /api/v1/products?category=1&page=1&size=20" `
  -Mode read-write

# 3. ดูผล
Get-Content .\.workflow-reports\ABC-123-summary.json | ConvertFrom-Json | Format-List
```

### Case 2: ใช้ GitHub Issue
```bash
# 1. เปิด issue จาก template 'Feature Request'
# 2. ใส่ label "phase:security" + "phase:performance"
# 3. Actions จะ trigger 2 เฟสนั้นอัตโนมัติ
# 4. ผลลัพธ์จะ comment กลับที่ issue + upload artifact
```

### Case 3: ใช้ Copilot Chat
```text
# ใน VS Code:
/gen-6phase-workflow

รหัสตั๋ว: ABC-123
ชื่องาน: แก้บั๊ก pagination
...
```

### Case 4: Auto trigger ตาม PR label
```bash
# 1. เปิด PR
# 2. ใส่ label "performance"
# 3. Workflow performance-review.yml จะรัน benchmark
# 4. Comment ผลลง PR อัตโนมัติ
```

---

## 💡 สรุปจุดเด่นของชุดเครื่องมือนี้

| องค์ประกอบ | จุดเด่น |
|-----------|--------|
| **PowerShell orchestrator** | รัน 6 เฟสอัตโนมัติ, ส่ง context ระหว่างเฟส, failFast |
| **Jira integration** | โพสต์ ADF comment + แนบ report + mask token |
| **GitHub Actions** | Trigger ตาม label, benchmark อัตโนมัติ, security scan ครบ |
| **Copilot instructions** | บังคับใช้มาตรฐานทุกการสนทนาในโปรเจกต์ |
| **Idempotent** | ป้องกันโพสต์ซ้ำ, ใช้ concurrency group |
| **Secure** | ใช้ secret ทุกจุด, ไม่ hardcode, mask log |
| **Auditable** | ทุกเฟสมี report + JSON + artifact 30 วัน |

ถ้าต้องการให้ผม:


บอกได้เลยครับ 🚀