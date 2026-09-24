# tb-subagent-driven-development

รัน implementation plan ใน session เดียว โดย dispatch subagent ใหม่ต่อ 1 task + review หลังทุก task แล้วปิดท้ายด้วย whole-branch review

> **ที่มา:** จาก official **superpowers** subagent-driven-development skill

## ทำอะไร

- รับ implementation plan แล้ว execute ทีละ task ในsession ปัจจุบัน
- แต่ละ task ใช้ subagent ใหม่ (context สะอาด ไม่ปนกัน) → implement + test + commit + self-review
- มี task reviewer ตรวจหลังทุก task (spec compliance + code quality) และ final code review รวบทั้ง branch
- เดินงานต่อเนื่องไม่หยุดถามระหว่าง task (หยุดเฉพาะตอน BLOCKED จริง ๆ)
- เป็นตัวเลือก "แนะนำ" ที่ tb-writing-plans ส่งต่อให้ (อีกทางคือ tb-executing-plans แบบ inline)

## Trigger เมื่อไหร่

เมื่อมี implementation plan ที่ task ส่วนใหญ่อิสระต่อกัน และต้องการ execute ในsession เดียวแบบมี review checkpoint อัตโนมัติ
มักใช้ต่อจาก `tb-writing-plans`

## ติดตั้ง

### 1. ดึง repo ลงเครื่อง (ครั้งแรก)

```bash
git clone https://gitlab.thaibevapp.com/thaibev-development-guide/codestyle/ai-skill/development.git
```

ถ้ามีอยู่แล้ว อัปเดตด้วย:

```bash
cd development && git pull
```

### 2. Copy folder `tb-subagent-driven-development` ไปวางที่ global skills folder ของ Copilot

- **macOS:** `~/.copilot/skills/tb-subagent-driven-development`
- **Windows:** `%USERPROFILE%\.copilot\skills\tb-subagent-driven-development`

```bash
# macOS
cp -R tb-subagent-driven-development ~/.copilot/skills/
```
```powershell
# Windows (PowerShell)
Copy-Item -Recurse tb-subagent-driven-development ~\.copilot\skills\
```

> **สำคัญ:** skill นี้อ้างถึงไฟล์ของ skill พี่น้องด้วย relative path
> (`../tb-requesting-code-review/code-reviewer.md`) และเรียกใช้
> `tb-using-git-worktrees`, `tb-finishing-a-development-branch`,
> `tb-requesting-code-review`, `tb-test-driven-development`, `tb-writing-plans`,
> `tb-executing-plans` — ต้อง copy ทุกตัวไปไว้ใน `~/.copilot/skills/` ให้ครบ
> ไม่งั้น path จะ resolve ไม่เจอ

## Requirements

- GitHub Copilot (VS Code)
- git (ใช้ใน scripts `review-package`, `task-brief`, `sdd-workspace` และระหว่างรันแผน)
- bash (สำหรับ scripts ในโฟลเดอร์ `scripts/`)
- skill พี่น้องในชุด tb- ติดตั้งครบ (ดู note ด้านบน)
