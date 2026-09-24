# tb-executing-plans

รัน implementation plan ที่เขียนไว้แล้วทีละ task พร้อม review checkpoint ระหว่างทาง

> **ที่มา:** จาก official **superpowers** executing-plans skill

## ทำอะไร

- รับไฟล์ implementation plan ที่เขียนไว้ แล้ว execute ทีละขั้นใน session แยก
- มี checkpoint ให้ review ระหว่างทาง ไม่รวดเดียวจบ
- เหมาะกับงานที่วาง plan ไว้ล่วงหน้า (คู่กับ tb-writing-plans)

## Trigger เมื่อไหร่

เมื่อมี implementation plan (ไฟล์ที่เขียนไว้แล้ว) และต้องการลงมือ execute ตาม plan นั้น

## ติดตั้ง

### 1. ดึง repo ลงเครื่อง (ครั้งแรก)

```bash
git clone https://gitlab.thaibevapp.com/thaibev-development-guide/codestyle/ai-skill/development.git
```

ถ้ามีอยู่แล้ว อัปเดตด้วย:

```bash
cd development && git pull
```

### 2. Copy folder `tb-executing-plans` ไปวางที่ global skills folder ของ Copilot

- **macOS:** `~/.copilot/skills/tb-executing-plans`
- **Windows:** `%USERPROFILE%\.copilot\skills\tb-executing-plans`

```bash
# macOS
cp -R tb-executing-plans ~/.copilot/skills/
```
```powershell
# Windows (PowerShell)
Copy-Item -Recurse tb-executing-plans ~\.copilot\skills\
```

## Requirements

- GitHub Copilot (VS Code)
