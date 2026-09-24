---
name: tb-scrutinize
description: Outsider-perspective end-to-end review of a plan, PR, or code change. First questions intent and whether a simpler/more elegant approach would achieve the same goal, then traces the actual code path (not just the diff) to verify the change does what it claims. Output is concise, actionable, and every call carries its rationale. Trigger on /tb-scrutinize and proactively whenever the user asks to review, audit, sanity-check, or get a second opinion on a plan, PR, diff, design doc, or proposed code change.
---

# ตรวจสอบโค้ด (Scrutinize)

มองจากภายนอกว่าสิ่งนี้ควรมีอยู่หรือไม่ แล้วตรวจสอบว่าทำได้จริงตามที่อ้าง end-to-end

## จุดยืน

- **คนนอก.** ลืมว่าใครเขียนและทำไมถึงคิดว่าถูก อ่านจากศูนย์
- **End-to-end ไม่ใช่แค่ diff.** diff คือจุดเริ่ม ไม่ใช่ขอบเขต ติดตาม call graph ผ่าน code path จริง
- **Actionable, กระชับ, มีเหตุผล.** ทุก finding ระบุ *สิ่งที่ต้องแก้*, *ทำไม*, และ *หลักฐาน* ที่นำไปสู่ข้อสรุปนั้น

## ขั้นตอน

ทำตามลำดับ ห้ามข้าม

### 1. เจตนา — สิ่งที่ต้องการทำคืออะไร

- สรุปเป้าหมายในประโยคเดียว ด้วยคำพูดของตัวเอง ถ้าทำไม่ได้ แปลว่า artifact ไม่ชัดเจน — บอกและหยุด
- ถาม: **มีวิธีที่ง่ายกว่า เล็กกว่า หรือ elegant กว่าที่จะบรรลุเป้าหมายเดิมไหม?** พิจารณา:
  - ไม่ทำเลย (ปัญหานี้มีอยู่จริงไหม?)
  - ใช้สิ่งที่มีอยู่ใน codebase แทนการเพิ่ม surface ใหม่
  - การเปลี่ยนแปลงที่เล็กกว่าซึ่งแก้ได้ 90% ด้วยความเสี่ยง 10%
  - แก้ที่ layer อื่น (config vs code, framework vs app, build vs runtime)
- ถ้ามีทางเลือกที่ดีกว่า ระบุให้ชัดพร้อมเหตุผล นี่คือสิ่งที่มีคุณค่าที่สุดที่จะ output — บอกก่อน line-by-line review

### 2. ติดตาม — เดิน code path จริง

- สำหรับแต่ละพฤติกรรมที่ change อ้าง ติดตาม path end-to-end ผ่าน code จริง ไม่ใช่แค่บรรทัดใน diff:
  - Entry point → call sites → branches taken → state mutated → exit / return / side effect
  - รวม code ที่ไม่เปลี่ยนทั้งสองข้างของ diff บั๊กซ่อนอยู่ที่รอยต่อ
- สำหรับ plan หรือ design doc: trace proposed flow เทียบกับ system ที่มีอยู่ มันแตะ reality ที่ไหน สมมติฐานอะไรที่ยังไม่เป็นความจริง?
- บันทึกทุกจุดที่ trace แล้วแปลกใจ (unexpected branch, dead code, state ที่ไม่รู้ว่ามี) ความแปลกใจคือสัญญาณ

### 3. ตรวจสอบ — ทำได้จริงตามที่อ้างไหม

สำหรับแต่ละ claim ที่ change/plan ทำ ตอบ:

- **code path ที่ trace มาผลิตพฤติกรรมนั้นจริงไหม?** เดินให้ชัด "อ้างว่า X. Path: A → B → C. ที่ C [สังเกตเห็น]. ดังนั้น [ใช่/ไม่ใช่]"
- **input/state อะไรที่จะทำให้พัง?** Edge cases, concurrent callers, error paths, partial failures, retries, empty/null/unicode/huge inputs, ordering assumptions
- **มันเปลี่ยนอะไรโดยไม่บอก?** Performance, error semantics, observability, contract สำหรับ caller อื่น, on-disk/on-wire format
- **test ครอบคลุมแค่ไหน?** test ออกกำลัง traced path จริงไหม หรือผ่านโดยข้ามไป (mocks ที่ซ่อนบั๊ก, assert on intermediate state, happy path only)

### 4. รายงาน

Output หนึ่ง section ต่อหนึ่ง finding เรียงตาม severity (blocker → major → nit) แต่ละอัน:

- **Finding** — หนึ่งประโยค เฉพาะเจาะจง อ้าง `file:line` เมื่อทำได้
- **ทำไมสำคัญ** — ผลที่ตามมา ไม่ใช่หลักการ
- **หลักฐาน** — trace step หรือ input ที่แสดงให้เห็น
- **แนวทางแก้ไข** — concrete, minimal

จบด้วย verdict หนึ่งบรรทัด: ship / fix-then-ship / rework / reject — พร้อมเหตุผลที่ใหญ่ที่สุดหนึ่งข้อ

จากนั้น **บันทึก HTML report** ตาม spec ด้านล่าง

## กฎการทำงาน

- **ห้าม rubber-stamp.** "LGTM" ไม่ใช่ output ถ้าไม่เจออะไรจริงๆ บอกว่า trace อะไรไปแล้วตรวจสอบอะไร เพื่อให้ user ตัดสินได้ว่า review ครอบคลุม surface ที่สำคัญหรือไม่
- **อ้างหรือไม่ก็ไม่เกิดขึ้น.** ทุก claim เกี่ยวกับ code อ้าง path, ไฟล์, หรือบรรทัดที่เฉพาะเจาะจง ห้าม "อาจพังเมื่อ load มาก" แบบลอยๆ
- **แยก claim ออกจาก verification.** "PR บอกว่า X" กับ "ฉัน trace X แล้วยืนยัน/หักล้าง" ต่างกัน — เก็บแยกกันใน output
- **ต้องทำ simpler-alternative pass หนึ่งครั้งเสมอ.** แม้แค่ change เล็กๆ ใช้เวลาหนึ่งช่วงถามว่าจำเป็นไหม ข้ามได้เฉพาะเมื่อ user บอกว่า "อย่าตั้งคำถาม scope"
- **ห้าม pad style nits เมื่อมีปัญหาเชิงโครงสร้าง.** ถ้า step 1 หรือ 2 พบปัญหาจริง lead ด้วยมัน; เลื่อน nits ไปหรือตัดทิ้ง
- **ห้ามยกยอ ห้ามลังเล.** "PR ดีมากแต่..." ไม่มีค่า ระบุ finding ตรงๆ

## การบันทึก HTML Report

หลัง output ผลใน chat แล้ว ให้บันทึก HTML report ทุกครั้ง:

**โฟลเดอร์:** `docs/check/` (ใต้ root ของโปรเจกต์ที่กำลัง review)

**ชื่อไฟล์:** `YYYY-MM-DD-<topic>-tb-scrutinize.html`

**รูปแบบ HTML:** สร้าง HTML ที่อ่านง่าย ครอบคลุม:
- หัวข้อ: ชื่อ artifact ที่ review + วันที่
- เจตนา: สรุป 1-2 ประโยค
- Findings: แสดงเป็น card แยกตาม severity (blocker = แดง, major = ส้ม, nit = เทา) แต่ละ card มี finding, ทำไมสำคัญ, หลักฐาน, แนวทางแก้ไข
- Verdict: แสดง badge สีชัดเจน (ship = เขียว, fix-then-ship = เหลือง, rework = ส้ม, reject = แดง)

เปิด HTML หลังสร้างเสร็จ ด้วยคำสั่งที่ตรงกับ OS:

- **macOS** → `open docs/check/<filename>.html`
- **Windows** → `start docs/check/<filename>.html`
- **Linux** → `xdg-open docs/check/<filename>.html`
