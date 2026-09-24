# tb-writing-plans

แปลง spec/requirement ให้กลายเป็น **implementation plan** ละเอียดที่แตกเป็น task เล็ก ๆ แบบ TDD — ก่อนเริ่มแตะโค้ด

> **ที่มา:** custom / ดัดแปลงมาจาก official **superpowers** writing-plans skill

## ทำอะไร

- เขียนแผนโดยสมมติว่าคนทำไม่รู้จัก codebase เลย — บอกครบว่าแตะไฟล์ไหน, โค้ดอะไร, เทสยังไง
- แตกงานเป็น task เล็ก ๆ ที่ทดสอบได้เอง แต่ละ step คือ 1 action (2-5 นาที) ตามวงจร TDD
- ห้ามมี placeholder (ไม่มี TODO/"implement later") — ต้องมีโค้ดจริงและคำสั่งจริงพร้อม expected output ทุก step
- **ออกแผนเป็นไฟล์ HTML ภาษาไทย** แบ่งกรอบชัด ใส่สีตาม section แล้ว **เปิดในเบราว์เซอร์อัตโนมัติให้ review ก่อนรัน** (คงโค้ด/พาธ/คำสั่ง/ชื่อฟังก์ชันเป็นภาษาอังกฤษ)
- บันทึกแผนที่ `<root-project>/docs/plan/YYYY-MM-DD-<feature-name>.html`

## Trigger เมื่อไหร่

เมื่อมี spec หรือ requirement ของงานที่มีหลายขั้นตอน และต้องการวางแผนก่อนลงมือเขียนโค้ด
มักใช้ต่อจาก `tb-brainstorming` ที่ได้ design/spec มาแล้ว

## ติดตั้ง

### 1. ดึง repo ลงเครื่อง (ครั้งแรก)

```bash
git clone https://gitlab.thaibevapp.com/thaibev-development-guide/codestyle/ai-skill/development.git
```

ถ้ามีอยู่แล้ว อัปเดตด้วย:

```bash
cd development && git pull
```

### 2. Copy folder `tb-writing-plans` ไปวางที่ global skills folder ของ Copilot

- **macOS:** `~/.copilot/skills/tb-writing-plans`
- **Windows:** `%USERPROFILE%\.copilot\skills\tb-writing-plans`

```bash
# macOS
cp -R tb-writing-plans ~/.copilot/skills/
```
```powershell
# Windows (PowerShell)
Copy-Item -Recurse tb-writing-plans ~\.copilot\skills\
```

## Requirements

- GitHub Copilot (VS Code)
- git (ใช้หา project root และ commit ระหว่างรันแผน)
- เบราว์เซอร์ (สำหรับ review ไฟล์แผน HTML)
