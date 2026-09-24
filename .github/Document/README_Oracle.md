# Oracle 

การตั้งค่าให้ `OracleServiceXE` เริ่มทำงานทันทีที่เปิด Windows 11 สามารถทำได้หลายวิธีตามที่อธิบายด้านล่างครับ

---

## ✅ วิธีที่ 1: ตั้งค่าแบบกราฟิก (GUI) ผ่าน Services

วิธีนี้ทำได้ง่ายที่สุด ไม่ต้องใช้คำสั่งอะไรเลย

1. กดปุ่ม `Windows + R` พิมพ์ `services.msc` แล้วกด Enter
2. ในหน้าต่าง Services ให้เลื่อนหา **`OracleServiceXE`**
3. คลิกขวาที่ `OracleServiceXE` เลือก **Properties**
4. ในแท็บ General ดูที่ช่อง **Startup type** ให้เปลี่ยนจาก `Manual` (หรือค่าปัจจุบัน) เป็น **`Automatic`**
5. คลิก **Apply** → **OK**

การตั้งค่า `Automatic` หมายความว่า Windows จะเริ่ม Service นี้โดยอัตโนมัติทันทีที่ระบบ boot เสร็จ

---

## ✅ วิธีที่ 2: ตั้งค่าผ่าน Command Line (CMD แบบ Administrator)

เปิด Command Prompt ในฐานะ Administrator แล้วใช้คำสั่ง `sc`:

```cmd
sc config OracleServiceXE start= auto
```
> **หมายเหตุสำคัญ:** ต้องมี space หลังเครื่องหมาย `=` นะครับ (เขียนว่า `start= auto` ไม่ใช่ `start=auto`)

คำสั่งนี้จะเปลี่ยน Startup type เป็น **Automatic** (หรือเรียกอีกอย่างว่า `auto` start type)

---

## ✅ วิธีที่ 3: ตั้งค่าผ่าน PowerShell (Run as Administrator)

เปิด PowerShell ในฐานะ Administrator แล้วใช้คำสั่ง:

```powershell
Set-Service -Name "OracleServiceXE" -StartupType 'Automatic'
```
หรือแบบสั้น:
```powershell
Set-Service OracleServiceXE -StartupType Automatic
```

คำสั่งนี้ยังสามารถใช้เพื่อเริ่ม Service ได้ทันทีด้วย โดยใส่พารามิเตอร์ `-Status running` เพิ่ม:
```powershell
Set-Service -Name "OracleServiceXE" -Status running -StartupType Automatic
```


---

## ⚠️ ข้อควรทราบเพิ่มเติม

- โดยค่าเริ่มต้น Oracle Database XE จะถูกติดตั้งให้ **เริ่มอัตโนมัติอยู่แล้ว** เมื่อเปิดระบบ ถ้าคุณจำเป็นต้องใช้คำสั่ง `net start` ทุกครั้ง แสดงว่ามีการปรับเปลี่ยน Startup type เป็น `Manual` หรืออาจเกิดปัญหากับ Service ครับ
- หากต้องการลดภาระของระบบ (ประหยัด RAM/CPU) แนะนำให้ตั้ง Startup type เป็น **Manual** แล้วค่อยใช้คำสั่ง `net start` เปิดเมื่อต้องการใช้งาน
- นอกจาก `OracleServiceXE` แล้ว ถ้าคุณใช้งาน Oracle ผ่าน Network ก็ควรตั้งค่า **`OracleXETNSListener`** ให้เริ่มอัตโนมัติเช่นกัน

---

## 🔍 ตรวจสอบผลลัพธ์

หลังจากตั้งค่าเสร็จแล้ว ให้ดูที่ Startup type ของ Service ได้ผ่านทางหน้าต่าง Services (`services.msc`) หรือใช้คำสั่ง:
```cmd
sc qc OracleServiceXE
```
หาคำว่า `START_TYPE` ถ้าขึ้น `AUTO_START` แสดงว่าตั้งค่าสำเร็จแล้วครับ