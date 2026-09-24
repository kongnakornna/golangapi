# tb-using-git-worktrees

เตรียม workspace แยก (isolated) ก่อนเริ่มงาน feature — ใช้ native worktree tool ของ platform ก่อน ถ้าไม่มีค่อย fallback เป็น git worktree

> **ที่มา:** จาก official **superpowers** using-git-worktrees skill

## ทำอะไร

- ตรวจก่อนว่าตอนนี้อยู่ใน workspace แยกอยู่แล้วหรือยัง (กันสร้างซ้อน)
- ถ้ายัง: ใช้ native worktree tool ก่อน ไม่มีค่อย `git worktree add` เอง
- ตั้งค่าโปรเจกต์ (install deps) + รัน test เพื่อยืนยัน baseline สะอาดก่อนลงมือ
- ใช้คู่กับ flow ของ tb-writing-plans / tb-subagent-driven-development / tb-executing-plans ที่ต้องการ workspace แยก

## Trigger เมื่อไหร่

เมื่อจะเริ่มงาน feature ที่ต้องแยกจาก workspace ปัจจุบัน หรือก่อน execute implementation plan

## ติดตั้ง

### 1. ดึง repo ลงเครื่อง (ครั้งแรก)

```bash
git clone https://gitlab.thaibevapp.com/thaibev-development-guide/codestyle/ai-skill/development.git
```

ถ้ามีอยู่แล้ว อัปเดตด้วย:

```bash
cd development && git pull
```

### 2. Copy folder `tb-using-git-worktrees` ไปวางที่ global skills folder ของ Copilot

- **macOS:** `~/.copilot/skills/tb-using-git-worktrees`
- **Windows:** `%USERPROFILE%\.copilot\skills\tb-using-git-worktrees`

```bash
# macOS
cp -R tb-using-git-worktrees ~/.copilot/skills/
```
```powershell
# Windows (PowerShell)
Copy-Item -Recurse tb-using-git-worktrees ~\.copilot\skills\
```

## Requirements

- GitHub Copilot (VS Code)
- git (worktree)
