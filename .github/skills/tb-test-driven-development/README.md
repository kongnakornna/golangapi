# tb-test-driven-development

เขียน test ก่อน implementation ตามวินัย TDD (red → green → refactor)

> **ที่มา:** จาก official **superpowers** test-driven-development skill

## ทำอะไร

- บังคับเขียน test ให้ fail ก่อน แล้วค่อยเขียนโค้ดให้ผ่าน แล้วค่อย refactor
- มีไฟล์ `testing-anti-patterns.md` รวม pattern ที่ควรเลี่ยง

## Trigger เมื่อไหร่

ก่อนเขียนโค้ด feature หรือ bugfix ใด ๆ

## ติดตั้ง

### 1. ดึง repo ลงเครื่อง (ครั้งแรก)

```bash
git clone https://gitlab.thaibevapp.com/thaibev-development-guide/codestyle/ai-skill/development.git
```

ถ้ามีอยู่แล้ว อัปเดตด้วย:

```bash
cd development && git pull
```

### 2. Copy folder `tb-test-driven-development` ไปวางที่ global skills folder ของ Copilot

- **macOS:** `~/.copilot/skills/tb-test-driven-development`
- **Windows:** `%USERPROFILE%\.copilot\skills\tb-test-driven-development`

```bash
# macOS
cp -R tb-test-driven-development ~/.copilot/skills/
```
```powershell
# Windows (PowerShell)
Copy-Item -Recurse tb-test-driven-development ~\.copilot\skills\
```

## Requirements

- GitHub Copilot (VS Code)
