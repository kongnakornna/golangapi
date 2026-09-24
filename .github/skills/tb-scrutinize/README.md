# tb-scrutinize

review plan / PR / code change จากมุมมองคนนอก — ตั้งคำถามว่า "ควรมีสิ่งนี้ไหม / มีวิธีง่ายกว่าไหม"
แล้ว trace code path จริงเพื่อยืนยันว่าทำได้ตามที่อ้าง end-to-end

> **ที่มา:** ref มาจาก **9arm**

## ทำอะไร

- มองแบบคนนอก ลืมว่าใครเขียน อ่าน artifact จากศูนย์
- ตั้งคำถาม intent + เสนอทางที่ง่าย/เล็ก/elegant กว่า ก่อน review ทีละบรรทัด
- เดิน code path จริง (ไม่ใช่แค่ diff) → ออก finding เรียงตาม severity + verdict
- บันทึก **HTML report** ไป `docs/check/` (ในโปรเจกต์) แล้วเปิดอัตโนมัติ

## Trigger เมื่อไหร่

- พิมพ์ `/tb-scrutinize`
- ขอ review / audit / sanity-check / second opinion ของ plan, PR, diff, design doc, โค้ดที่จะเปลี่ยน

## ติดตั้ง

### 1. ดึง repo ลงเครื่อง (ครั้งแรก)

```bash
git clone https://gitlab.thaibevapp.com/thaibev-development-guide/codestyle/ai-skill/development.git
```

ถ้ามีอยู่แล้ว อัปเดตด้วย:

```bash
cd development && git pull
```

### 2. Copy folder `tb-scrutinize` ไปวางที่ global skills folder ของ Copilot

- **macOS:** `~/.copilot/skills/tb-scrutinize`
- **Windows:** `%USERPROFILE%\.copilot\skills\tb-scrutinize`

```bash
# macOS
cp -R tb-scrutinize ~/.copilot/skills/
```
```powershell
# Windows (PowerShell)
Copy-Item -Recurse tb-scrutinize ~\.copilot\skills\
```

## Requirements

- GitHub Copilot (VS Code)
