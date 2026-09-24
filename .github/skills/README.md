# ai-skill / development

รวม skills สำหรับ GitHub Copilot (VS Code) สำหรับทีม AI Dev — ใช้งานได้ทันทีโดยไม่ต้องเขียนใหม่

## Skills ที่มี

| Skill | คำอธิบาย |
|---|---|
| [tb-read-jira](./tb-read-jira/) | ดึงและสรุป Jira ticket จาก digital-and-technology.atlassian.net |
| [tb-investigate-error](./tb-investigate-error/) | สืบสวนหาสาเหตุ error/bug แล้วออกเป็นรายงาน HTML ภาษาไทย |
| [tb-brainstorming](./tb-brainstorming/) | แปลงไอเดียให้เป็น design/spec ผ่านบทสนทนา ก่อนเริ่มเขียนโค้ด (custom จาก superpowers) |
| [tb-writing-plans](./tb-writing-plans/) | แปลง spec เป็น implementation plan แบบ TDD ออกเป็นไฟล์ HTML ภาษาไทย เปิด review อัตโนมัติ (custom จาก superpowers) |
| [tb-executing-plans](./tb-executing-plans/) | รัน implementation plan ทีละ task พร้อม review checkpoint (จาก superpowers) |
| [tb-test-driven-development](./tb-test-driven-development/) | เขียน test ก่อนโค้ด ตามวินัย TDD red-green-refactor (จาก superpowers) |
| [tb-subagent-driven-development](./tb-subagent-driven-development/) | รัน plan ใน session เดียว dispatch subagent ต่อ task + review หลังทุก task (จาก superpowers) |
| [tb-using-git-worktrees](./tb-using-git-worktrees/) | เตรียม workspace แยก (worktree) ก่อนเริ่มงาน feature (จาก superpowers) |
| [tb-requesting-code-review](./tb-requesting-code-review/) | dispatch subagent มา review โค้ดก่อนปัญหาลุกลาม/ก่อน merge (จาก superpowers) |
| [tb-finishing-a-development-branch](./tb-finishing-a-development-branch/) | ปิดงาน branch — merge / PR / เก็บ / ทิ้ง พร้อม cleanup (จาก superpowers) |
| [tb-scrutinize](./tb-scrutinize/) | review plan/PR/code จากมุมคนนอก ตั้งคำถาม intent + trace code path จริง ออก HTML report (ref 9arm) |
| [grill-me](https://github.com/mattpocock/skills/blob/main/skills/productivity/grill-me/SKILL.md) | สัมภาษณ์ซักไซ้ plan/design อย่างเข้มข้น เพื่อลับให้คมก่อนลงมือ (ต้นฉบับ mattpocock) |
| [grill-with-docs](https://github.com/mattpocock/skills/blob/main/skills/engineering/grill-with-docs/SKILL.md) | ซักไซ้ plan/design พร้อมสร้างเอกสาร ADR และ glossary ไปด้วยระหว่างคุย (ต้นฉบับ mattpocock) |
| [handoff](https://github.com/mattpocock/skills/tree/main/skills/productivity/handoff) | สรุปบทสนทนาปัจจุบันเป็นเอกสาร handoff ให้ agent อื่นรับงานต่อได้ (ต้นฉบับ mattpocock) |

## วิธีติดตั้ง (GitHub Copilot บน VS Code — global)

ทีมใช้งานผ่าน **GitHub Copilot บน VS Code** ติดตั้ง skill เป็น global

### 1. Clone repo (ครั้งแรก)

```bash
git clone https://gitlab.thaibevapp.com/thaibev-development-guide/codestyle/ai-skill/development.git
cd development
```

> มีอยู่แล้ว? เข้าไปในโฟลเดอร์ repo แล้ว `git pull`

### 2. Copy skill ไปวาง

> คำสั่ง copy ด้านล่าง **ต้องรันจากในโฟลเดอร์ repo** (ที่มี `README.md` และโฟลเดอร์ `tb-*` อยู่) ไม่งั้นจะเจอ error `Cannot find path ... because it does not exist` — เช็กว่าอยู่ถูกที่ด้วย `ls` (macOS) หรือ `Get-ChildItem` (Windows) แล้วต้องเห็นโฟลเดอร์ `tb-read-jira`

ปลายทาง global skills folder:

| OS | ปลายทาง |
|---|---|
| **macOS** | `~/.copilot/skills/` |
| **Windows** | `%USERPROFILE%\.copilot\skills\` |

**copy ทุก skill ทีเดียว — macOS:**

```bash
mkdir -p ~/.copilot/skills
cp -R tb-* ~/.copilot/skills/
```

**copy ทุก skill ทีเดียว — Windows (PowerShell):**

```powershell
New-Item -ItemType Directory -Force ~\.copilot\skills
Copy-Item -Recurse -Force tb-* ~\.copilot\skills\
```

หรือ copy เฉพาะ skill ที่ต้องการ (แทน `tb-*` ด้วยชื่อ skill เช่น `tb-read-jira`)

> **อัปเดต:** `git pull` แล้วรันคำสั่ง copy ด้านบนซ้ำ (มี `-Force` / `cp -R` จะ copy ทับให้เอง)

## วิธีติดตั้ง Custom Instructions (global — VS Code & CLI)

Custom Instructions = ชุดกติกาที่ควบคุม AI ตอนช่วยเขียนโค้ด — แก้โค้ดให้น้อยที่สุด ไม่เปลี่ยนพฤติกรรมเดิมของระบบ ต้องสรุป + ขอ confirm ก่อนลงมือ และตอบกลับเป็นภาษาไทย

ติดตั้งแบบเดียวกับ skill — `git pull` แล้ว copy ไฟล์ [`.github/copilot-instructions.md`](./.github/copilot-instructions.md) ไปวางในโฟลเดอร์ที่ต้องการ (รันจากในโฟลเดอร์ repo)

แต่ละ tool อ่าน path ต่างกัน:

| Tool | Path (macOS / Linux) | Path (Windows) |
|---|---|---|
| **Copilot CLI** | `~/.copilot/instructions/` | `%USERPROFILE%\.copilot\instructions\` |
| **VS Code** | `~/Library/Application Support/Code/User/prompts/` | `%APPDATA%\Code\User\prompts\` |

### macOS / Linux

**CLI:**

```bash
mkdir -p ~/.copilot/instructions
cp .github/copilot-instructions.md ~/.copilot/instructions/thaibev.instructions.md
```

**VS Code:**

```bash
mkdir -p "$HOME/Library/Application Support/Code/User/prompts"
cp .github/copilot-instructions.md "$HOME/Library/Application Support/Code/User/prompts/thaibev.instructions.md"
```

### Windows (PowerShell)

**CLI:**

```powershell
New-Item -ItemType Directory -Force "$env:USERPROFILE\.copilot\instructions"
Copy-Item -Force .github\copilot-instructions.md "$env:USERPROFILE\.copilot\instructions\thaibev.instructions.md"
```

**VS Code:**

```powershell
New-Item -ItemType Directory -Force "$env:APPDATA\Code\User\prompts"
Copy-Item -Force .github\copilot-instructions.md "$env:APPDATA\Code\User\prompts\thaibev.instructions.md"
```

> **อัปเดต:** `git pull` แล้วรันคำสั่ง copy ด้านบนซ้ำ
>
> **ต้องการให้มีผลเฉพาะบาง project?** copy ไฟล์ไปวางที่ `.github/copilot-instructions.md` ของ project นั้น (Copilot อ่าน workspace-level ให้อัตโนมัติ)
>
> **VS Code Insiders:** เปลี่ยน `Code` → `Code - Insiders` ใน path ด้านบน

## Branching

- `master` — stable ทดสอบแล้ว ใช้งานได้
- `dev` — กำลังพัฒนา / รอ review

## วิธี Contribute

1. Fork repo นี้
2. สร้าง skill ใหม่ใน folder ของตัวเอง (`kebab-case`)
3. แต่ละ skill ต้องมี `SKILL.md` และ `README.md`
4. ส่ง Merge Request เข้า branch `dev`
5. ใช้ MR template ที่มีให้กรอกให้ครบ

ดูรายละเอียดเพิ่มเติมที่ [.gitlab/merge_request_templates/new_skill.md](.gitlab/merge_request_templates/new_skill.md)
