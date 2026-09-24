# tb-finishing-a-development-branch

ปิดงานบน development branch — ยืนยัน test ผ่าน แล้วเสนอตัวเลือก merge / PR / เก็บไว้ / ทิ้ง พร้อมจัดการ cleanup ให้

> **ที่มา:** จาก official **superpowers** finishing-a-development-branch skill

## ทำอะไร

- รัน test ให้ผ่านก่อนเสมอ (test ไม่ผ่าน = ไม่เสนอ option)
- ตรวจ environment (repo ปกติ / worktree / detached HEAD) แล้วเลือกเมนูให้เหมาะ
- เสนอ option ชัด ๆ: merge local / push + เปิด PR / เก็บ branch ไว้ / ทิ้งงาน (ต้องพิมพ์ยืนยันก่อนทิ้ง)
- จัดการ cleanup worktree/branch ตามตัวเลือกที่เลือก
- เป็นปลายทางของ flow ทั้ง tb-subagent-driven-development และ tb-executing-plans

## Trigger เมื่อไหร่

เมื่อ implement เสร็จ test ผ่านหมด และต้องตัดสินใจว่าจะ integrate งานยังไง

## ติดตั้ง

### 1. ดึง repo ลงเครื่อง (ครั้งแรก)

```bash
git clone https://gitlab.thaibevapp.com/thaibev-development-guide/codestyle/ai-skill/development.git
```

ถ้ามีอยู่แล้ว อัปเดตด้วย:

```bash
cd development && git pull
```

### 2. Copy folder `tb-finishing-a-development-branch` ไปวางที่ global skills folder ของ Copilot

- **macOS:** `~/.copilot/skills/tb-finishing-a-development-branch`
- **Windows:** `%USERPROFILE%\.copilot\skills\tb-finishing-a-development-branch`

```bash
# macOS
cp -R tb-finishing-a-development-branch ~/.copilot/skills/
```
```powershell
# Windows (PowerShell)
Copy-Item -Recurse tb-finishing-a-development-branch ~\.copilot\skills\
```

## Requirements

- GitHub Copilot (VS Code)
- git (merge / worktree / branch cleanup)
