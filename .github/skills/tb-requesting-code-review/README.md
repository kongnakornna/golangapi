# tb-requesting-code-review

dispatch subagent มา review โค้ดก่อนปัญหาลุกลาม — reviewer ได้ context ที่จัดมาให้พอดี ไม่ใช่ประวัติทั้ง session

> **ที่มา:** จาก official **superpowers** requesting-code-review skill

## ทำอะไร

- หา git SHA (base/head) แล้ว dispatch code reviewer subagent ตาม template `code-reviewer.md`
- จัด context ให้ reviewer เฉพาะส่วนที่ต้องดู (ไม่เทประวัติ session เข้าไป)
- ใช้เป็น final whole-branch review ใน flow ของ tb-subagent-driven-development และจุด checkpoint ของ tb-executing-plans

## Trigger เมื่อไหร่

หลังทำ task เสร็จ, หลังจบ feature ใหญ่, หรือก่อน merge เข้า main — เมื่ออยากให้มีสายตาที่สองตรวจงาน

## ติดตั้ง

### 1. ดึง repo ลงเครื่อง (ครั้งแรก)

```bash
git clone https://gitlab.thaibevapp.com/thaibev-development-guide/codestyle/ai-skill/development.git
```

ถ้ามีอยู่แล้ว อัปเดตด้วย:

```bash
cd development && git pull
```

### 2. Copy folder `tb-requesting-code-review` ไปวางที่ global skills folder ของ Copilot

- **macOS:** `~/.copilot/skills/tb-requesting-code-review`
- **Windows:** `%USERPROFILE%\.copilot\skills\tb-requesting-code-review`

```bash
# macOS
cp -R tb-requesting-code-review ~/.copilot/skills/
```
```powershell
# Windows (PowerShell)
Copy-Item -Recurse tb-requesting-code-review ~\.copilot\skills\
```

> **หมายเหตุ:** `tb-subagent-driven-development` อ้างไฟล์ `code-reviewer.md` ของ skill นี้
> ผ่าน relative path — ถ้าใช้ flow นั้นต้อง copy skill นี้ไปไว้ใน `~/.copilot/skills/` ด้วย

## Requirements

- GitHub Copilot (VS Code)
- git (หา SHA สำหรับ review)
