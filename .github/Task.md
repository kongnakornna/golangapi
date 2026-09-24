http://tms.thaibev.com/?tokenData

---------------------------
 ภาษาไทย
เทมเพลต Task List สำหรับงาน Vibe Coding
# Task: [ชื่องาน]
วันที่: [YYYY-MM-DD]
สถานะ: [ ] ยังไม่เริ่ม / [ ] กำลังทำ / [ ] รอรีวิว / [ ] เสร็จ

## ขอบเขต (Scope)
- [ ] งานที่ต้องทำ 1
- [ ] งานที่ต้องทำ 2
- [ ] งานที่ต้องทำ 3

## สิ่งที่ยังไม่ทำ (Out of scope)
- [ ] ไม่ทำ X
- [ ] ไม่ทำ Y

## Spec สรุป
- เป้าหมาย:
- ผู้ใช้:
- ข้อมูลที่เกี่ยวข้อง:

## ขั้นตอนการทำ (Workflow steps)
- [ ] 1. อัปเดต CLAUDE.md (ถ้าจำเป็น)
- [ ] 2. ให้ Claude วางแผน
- [ ] 3. ให้ Claude สร้าง/แก้ไข
- [ ] 4. ทดสอบใน local
- [ ] 5. ดู diff และ commit
- [ ] 6. Push และเปิด PR (ถ้าใช้)
- [ ] 7. ทดสอบ Preview URL
- [ ] 8. ขอรีวิว (ถ้ามีคนอื่น)
- [ ] 9. Merge และทดสอบ production
- [ ] 10. เขียน release note

## เกณฑ์เสร็จ (Definition of Done)
- [ ] Acceptance criteria ทั้งหมดผ่าน
- [ ] ทดสอบ manual แล้ว
- [ ] Diff ผ่านการตรวจสอบ
- [ ] ไม่มี secret หลุด
- [ ] Deploy ขึ้น production/preview แล้ว

## หมายเหตุ

--------------------------------------
--------------------------------------

ตัวอย่างที่กรอกแล้ว
# Task: แก้ไข duplicate email error ใน waitlist form
วันที่: 2026-05-27
สถานะ: [x] เสร็จ

## ขอบเขต
- [x] แก้ไขให้แสดงข้อความ "อีเมลนี้ลงทะเบียนแล้ว" แทน error จาก database
- [x] เก็บ error เดิมไว้ใน log (ไม่แสดงให้ผู้ใช้เห็น)
- [x] ทดสอบบน local และ preview

## เกณฑ์เสร็จ
- [x] ผู้ใช้เห็นข้อความที่เข้าใจได้
- [x] ไม่มี technical error โผล่
- [x] ยังคง insert ข้อมูลใหม่ได้
English
Task List Template for Vibe Coding
# Task: [Task name]
Date: [YYYY-MM-DD]
Status: [ ] Not started / [ ] In progress / [ ] Review / [ ] Done

## Scope
- [ ] To-do item 1
- [ ] To-do item 2
- [ ] To-do item 3

## Out of scope
- [ ] Not doing X
- [ ] Not doing Y

## Spec summary
- Goal:
- Users:
- Related data:

## Workflow steps
- [ ] 1. Update CLAUDE.md (if needed)
- [ ] 2. Ask Claude to plan
- [ ] 3. Ask Claude to build/edit
- [ ] 4. Test locally
- [ ] 5. Review diff and commit
- [ ] 6. Push and open PR (if used)
- [ ] 7. Test Preview URL
- [ ] 8. Request review (if team)
- [ ] 9. Merge and test production
- [ ] 10. Write release note

## Definition of Done
- [ ] All acceptance criteria met
- [ ] Manual testing done
- [ ] Diff reviewed
- [ ] No secrets exposed
- [ ] Deployed to production/preview

## Notes
Filled Example
# Task: Fix duplicate email error in waitlist form
Date: 2026-05-27
Status: [x] Done

## Scope
- [x] Show "This email is already on the waitlist" message instead of database error
- [x] Keep original error in log (not shown to user)
- [x] Test on local and preview

## Definition of Done
- [x] User sees understandable message
- [x] No technical error shown
- [x] New submissions still work


--------------------------------------
--------------------------------------

เทมเพลตรายการตรวจสอบ / Checklist Template
ภาษาไทย
Checklist A: ก่อนเริ่มโปรเจกต์ใหม่
# Pre‑Project Checklist

## ไอเดียและขอบเขต
- [ ] อธิบาย version แรกได้ใน 3 ประโยค
- [ ] ฟีเจอร์หลักไม่เกิน 3 อย่าง
- [ ] รู้ว่าใครคือผู้ใช้
- [ ] รู้ว่าผู้ใช้ต้องทำ flow อะไรให้สำเร็จ
- [ ] รู้ว่าต้องเก็บข้อมูลอะไรบ้าง

## ความเสี่ยง
- [ ] ไม่เก็บข้อมูลอ่อนไหว (ถ้าไม่จำเป็น)
- [ ] ถ้าแอปพัง ธุรกิจไม่หยุดทันที
- [ ] สามารถปิดหรือ rollback ได้
- [ ] รู้ว่าจะให้ใครช่วย review ถ้าปลอดภัย

## เทคนิค
- [ ] มีโฟลเดอร์โปรเจกต์แล้ว
- [ ] มี Git init หรือ repo พร้อม
- [ ] Claude Code เข้าถึงโปรเจกต์ได้
- [ ] มี CLAUDE.md อย่างน้อย minimal
Checklist B: ก่อน Deploy ขึ้น Production
# Pre‑Deploy Checklist

## โค้ดและ Git
- [ ] ไม่มี `.env` หรือ `.env.local` ถูก commit
- [ ] ไม่มี secret/API key ใน source code
- [ ] Git status สะอาด (ทุกอย่าง commit แล้ว)
- [ ] diff ผ่านการตรวจสอบ

## Build และ Environment
- [ ] `npm run build` หรือคำสั่ง build ทำงานใน local
- [ ] Environment variables ทั้งหมดถูกตั้งใน Vercel (Production)
- [ ] ไม่มีตัวแปรที่ใช้ค่า test ค้างอยู่

## ฐานข้อมูล (ถ้ามี)
- [ ] RLS เปิดอยู่
- [ ] Public user insert ได้ (ถ้าต้องการ) แต่ select/update/delete ไม่ได้
- [ ] ไม่มี service_role key ใน frontend

## การทดสอบ
- [ ] หน้าเปิดได้บน desktop และ mobile
- [ ] ฟอร์ม validation ทำงาน
- [ ] Success และ error state ชัดเจน
- [ ] ไม่มีข้อมูล test หลงเหลือใน production

## หลัง Deploy
- [ ] ทดสอบ production URL แล้ว
- [ ] rollback plan พร้อม (รู้วิธีย้อนกลับ)
- [ ] มีการจด deploy note
Checklist C: การ Debug
# Debug Evidence Checklist

ก่อนขอให้ Claude แก้ไข ให้เก็บหลักฐานให้ครบ:

## ข้อมูลทั่วไป
- [ ] Environment: local / preview / production
- [ ] URL หรือหน้าที่เกิดปัญหา
- [ ] เวลาที่เกิด (ถ้ามี)

## พฤติกรรม
- [ ] Expected behavior (ควรเกิดอะไร)
- [ ] Actual behavior (เกิดอะไรขึ้นจริง)
- [ ] Steps to reproduce (ทำซ้ำยังไง)

## หลักฐาน
- [ ] Browser console error (copy ข้อความ)
- [ ] Network tab status และ response
- [ ] Terminal error log (ถ้าเป็น local)
- [ ] Vercel build/deployment log (ถ้า deploy fail)
- [ ] Supabase error หรือ policy info (ถ้าเกี่ยวข้อง)

## การเปลี่ยนแปลงล่าสุด
- [ ] รายการไฟล์ที่เพิ่งเปลี่ยน
- [ ] commit ล่าสุด
- [ ] deployment ล่าสุด
Checklist D: Security (10 ข้อก่อนเปิดให้คนอื่นใช้)
# Security Launch Checklist

## Data
- [ ] เก็บเฉพาะข้อมูลที่จำเป็นจริง ๆ
- [ ] ไม่มีข้อมูล sensitive ที่ไม่จำเป็น

## Secrets & Keys
- [ ] ไม่มี secret ในโค้ด, GitHub หรือ frontend
- [ ] ใช้ publishable key (ไม่ใช่ service_role) ใน browser
- [ ] .env.example มีแค่ placeholder
- [ ] Vercel env vars ตั้งครบ (ไม่ใช่ค่าทดสอบ)

## Database & RLS
- [ ] RLS เปิดอยู่
- [ ] Public user แทรก (insert) ได้เท่าที่จำเป็น
- [ ] Public user ไม่สามารถอ่านข้อมูลทั้งหมดได้
- [ ] ถ้ามี dashboard: มี authentication และ authorization ชัดเจน

## Input & Error
- [ ] form validate ข้อมูล
- [ ] error message ไม่เปิดเผยความลับทางเทคนิค

## Operations
- [ ] มี rollback plan
- [ ] มี owner และ review cadence
- [ ] Launch gate: [ ] Green / [ ] Yellow / [ ] Red
English
Checklist A: Before Starting a New Project
# Pre‑Project Checklist

## Idea & Scope
- [ ] Can describe version 1 in 3 sentences
- [ ] No more than 3 main features
- [ ] Know who the user is
- [ ] Know the main user flow
- [ ] Know what data needs to be stored

## Risks
- [ ] Not collecting sensitive data (unless necessary)
- [ ] If app breaks, business doesn't stop immediately
- [ ] Can pause or rollback
- [ ] Know who can help review if needed

## Technical
- [ ] Project folder exists
- [ ] Git init or repo ready
- [ ] Claude Code can access the project
- [ ] At least a minimal CLAUDE.md exists
Checklist B: Before Deploying to Production
# Pre‑Deploy Checklist

## Code & Git
- [ ] No `.env` or `.env.local` committed
- [ ] No secrets/API keys in source code
- [ ] Git status is clean (all changes committed)
- [ ] Diff reviewed

## Build & Environment
- [ ] `npm run build` (or equivalent) works locally
- [ ] All environment variables set in Vercel (Production)
- [ ] No leftover test values

## Database (if any)
- [ ] RLS is enabled
- [ ] Public can insert only (if needed), not select/update/delete
- [ ] No service_role key in frontend

## Testing
- [ ] Page works on desktop and mobile
- [ ] Form validation works
- [ ] Success and error states are clear
- [ ] No test data left in production

## Post‑Deploy
- [ ] Production URL tested
- [ ] Rollback plan ready
- [ ] Deploy note recorded
Checklist C: Debugging
# Debug Evidence Checklist

Before asking Claude to fix, collect this evidence:

## General
- [ ] Environment: local / preview / production
- [ ] URL or page where issue occurred
- [ ] Time of occurrence (if relevant)

## Behavior
- [ ] Expected behavior
- [ ] Actual behavior
- [ ] Steps to reproduce

## Evidence
- [ ] Browser console error (copy text)
- [ ] Network tab status and response
- [ ] Terminal error log (if local)
- [ ] Vercel build/deployment log (if deploy fails)
- [ ] Supabase error or policy info (if relevant)

## Recent changes
- [ ] List of recently changed files
- [ ] Last commit
- [ ] Last deployment
Checklist D: Security (10 items before real users)
# Security Launch Checklist

## Data
- [ ] Collect only necessary data
- [ ] No unnecessary sensitive data

## Secrets & Keys
- [ ] No secrets in code, GitHub, or frontend
- [ ] Using publishable key (not service_role) in browser
- [ ] .env.example contains only placeholders
- [ ] Vercel env vars fully set (not test values)

## Database & RLS
- [ ] RLS is enabled
- [ ] Public can insert only what's needed
- [ ] Public cannot read all data
- [ ] If dashboard exists: authentication & authorization are clear

## Input & Error
- [ ] Form validates input
- [ ] Error messages don't expose technical secrets

## Operations
- [ ] Rollback plan exists
- [ ] Owner and review cadence defined
- [ ] Launch gate: [ ] Green / [ ] Yellow / [ ] Red


--------------------------------------
--------------------------------------

# Vibe Coding Templates

## 1. Code Feature Template / เทมเพลตการสร้างฟีเจอร์ใหม่

**ภาษาไทย**

```markdown
# Task: [ชื่อฟีเจอร์]
วันที่: [YYYY-MM-DD]
สถานะ: [ ] ยังไม่เริ่ม / [ ] กำลังทำ / [ ] รอรีวิว / [ ] เสร็จ

## ขอบเขต (Scope)
- [ ] ฟีเจอร์หลักข้อที่ 1
- [ ] ฟีเจอร์หลักข้อที่ 2
- [ ] การจัดการ error กรณี ...

## สิ่งที่ยังไม่ทำ (Out of scope)
- [ ] ไม่ทำ UI ที่ซับซ้อนเกิน v1
- [ ] ไม่ทำระบบแจ้งเตือนอีเมล (ไว้ v2)

## Spec สรุป
- เป้าหมาย: [ผู้ใช้ทำอะไรได้ใหม่]
- ผู้ใช้: [ใครใช้]
- ข้อมูลที่เกี่ยวข้อง: [ตาราง DB, API, state]

## ขั้นตอนการทำ (Workflow steps)
- [ ] 1. อัปเดต CLAUDE.md (ถ้าจำเป็น)
- [ ] 2. ให้ Claude วางแผน: สร้าง component, API route, หรือ DB migration
- [ ] 3. ให้ Claude สร้าง/แก้ไขไฟล์
- [ ] 4. ทดสอบใน local (ทั้ง success และ error paths)
- [ ] 5. ดู diff และ commit
- [ ] 6. Push และเปิด PR (ถ้ามีทีม)
- [ ] 7. ทดสอบ Preview URL
- [ ] 8. ขอรีวิว
- [ ] 9. Merge และทดสอบ production
- [ ] 10. เขียน release note (อัปเดต doc ถ้ามี)

## เกณฑ์เสร็จ (Definition of Done)
- [ ] ฟีเจอร์ทำงานตาม spec
- [ ] ทดสอบ manual บน desktop และ mobile
- [ ] ไม่มี console error
- [ ] error handling ใช้ได้ (ผู้ใช้เห็นข้อความเข้าใจง่าย)
- [ ] ถ้ามี DB: RLS ถูกต้อง, migration ทดสอบแล้ว
- [ ] Diff ผ่านการตรวจสอบ, ไม่มี secret หลุด

## หมายเหตุ
```

**English**

```markdown
# Task: [Feature name]
Date: [YYYY-MM-DD]
Status: [ ] Not started / [ ] In progress / [ ] Review / [ ] Done

## Scope
- [ ] Main feature item 1
- [ ] Main feature item 2
- [ ] Error handling for ...

## Out of scope
- [ ] No complex UI beyond v1
- [ ] No email notifications (v2)

## Spec summary
- Goal: [What can user do]
- Users: [Who uses it]
- Related data: [DB tables, APIs, state]

## Workflow steps
- [ ] 1. Update CLAUDE.md (if needed)
- [ ] 2. Ask Claude to plan: component, API route, or DB migration
- [ ] 3. Ask Claude to create/edit files
- [ ] 4. Test locally (both success and error paths)
- [ ] 5. Review diff and commit
- [ ] 6. Push and open PR (if team)
- [ ] 7. Test Preview URL
- [ ] 8. Request review
- [ ] 9. Merge and test production
- [ ] 10. Write release note (update docs)

## Definition of Done
- [ ] Feature works per spec
- [ ] Manual testing on desktop & mobile
- [ ] No console errors
- [ ] Error handling works (user sees friendly message)
- [ ] If DB: RLS correct, migration tested
- [ ] Diff reviewed, no secrets exposed

## Notes
```

---

## 2. Review Template / เทมเพลตสำหรับตรวจสอบโค้ด

**ภาษาไทย**

```markdown
# Review: [PR หรือ Commit ที่ตรวจสอบ]
วันที่: [YYYY-MM-DD]
ผู้ตรวจสอบ: [ชื่อ]
PR link: [URL]

## สิ่งที่ตรวจสอบ
- [ ] โค้ดอ่านรู้เรื่อง, มี comment เฉพาะจุดที่ซับซ้อน
- [ ] ตรงตาม spec และ scope ที่เขียนไว้
- [ ] ไม่มี side effect กับฟีเจอร์อื่น
- [ ] การจัดการ error ครบถ้วน (ทั้ง success และ fail)
- [ ] form validation ครอบคลุม input ที่เป็นไปได้ทั้งหมด
- [ ] performance ปกติ (ไม่มี re-render รอบ, query ซ้ำซ้อน)

## ความปลอดภัย (Security)
- [ ] ไม่มี secret, API key, token ในโค้ด
- [ ] RLS policy ถูกต้อง (public insert ได้เท่าที่จำเป็น)
- [ ] ข้อมูล sensitive ไม่ถูก log หรือส่งไป frontend โดยไม่จำเป็น
- [ ] input ถูก sanitize (ป้องกัน XSS, SQL injection)

## การทดสอบ (Testing)
- [ ] ทดสอบใน local แล้ว
- [ ] ทดสอบ Preview URL แล้ว
- [ ] Edge cases: [กรณีที่ต้องระบุ]
- [ ] Mobile responsive (ถ้ามี UI)

## Issues ที่พบ
- [ ] [Issue 1: อธิบาย]
- [ ] [Issue 2: อธิบาย]

## การตัดสินใจ
- [ ] Approve – พร้อม merge
- [ ] Request changes – รอแก้ไข [จำนวน] รายการ
- [ ] Comment – ให้ feedback แต่ไม่บล็อก

## หมายเหตุ
```

**English**

```markdown
# Review: [PR or Commit being reviewed]
Date: [YYYY-MM-DD]
Reviewer: [Name]
PR link: [URL]

## Review checklist
- [ ] Code is readable, comments only for complex parts
- [ ] Meets spec and scope as written
- [ ] No side effects on other features
- [ ] Error handling is complete (both success & fail paths)
- [ ] Form validation covers all possible inputs
- [ ] Performance is normal (no excess re-renders, duplicate queries)

## Security
- [ ] No secrets, API keys, tokens in code
- [ ] RLS policies correct (public insert only where needed)
- [ ] Sensitive data not logged or sent to frontend unnecessarily
- [ ] Input is sanitized (XSS, SQL injection prevention)

## Testing
- [ ] Tested locally
- [ ] Tested on Preview URL
- [ ] Edge cases: [specify]
- [ ] Mobile responsive (if UI exists)

## Issues found
- [ ] [Issue 1: description]
- [ ] [Issue 2: description]

## Decision
- [ ] Approve – ready to merge
- [ ] Request changes – waiting for [number] fixes
- [ ] Comment – feedback given, not blocking

## Notes
```

---

## 3. Hotfix Issues Template / เทมเพลตแก้ไขด่วนสำหรับ Production

**ภาษาไทย**

```markdown
# Hotfix: [ชื่อปัญหา]
วันที่: [YYYY-MM-DD]
ความรุนแรง: [ ] วิกฤต (เว็บพัง) / [ ] สูง (ฟีเจอร์หลักพัง) / [ ] ปานกลาง
สถานะ: [ ] กำลังวิเคราะห์ / [ ] กำลังแก้ / [ ] ทดสอบ / [ ] Deploy แล้ว

## อาการปัญหา (Symptoms)
- [ ] ผู้ใช้เจออะไร: [ข้อความ error, หน้าขาว, ทำงานไม่ได้]
- [ ] Browser console error: [copy ข้อความ]
- [ ] Network response: [status code, response body]
- [ ] Log (ถ้ามี): [link หรือ snippet]

## Steps to reproduce
1. [ขั้นตอนที่ 1]
2. [ขั้นตอนที่ 2]
3. [ผลลัพธ์ที่ผิด]

## Root cause (หลังวิเคราะห์)
- [ ] สาเหตุ: [อธิบายสั้นๆ]
- [ ] ไฟล์/บรรทัดที่ผิด: [path:line]

## วิธีแก้ไข (Fix)
- [ ] แผนแก้: [อธิบาย]
- [ ] ไฟล์ที่ต้องเปลี่ยน: [รายการ]

## ขั้นตอนการทำ Hotfix
- [ ] 1. สร้าง branch จาก main: `hotfix/[ชื่อ]`
- [ ] 2. แก้ไขโค้ด
- [ ] 3. ทดสอบ local และ preview (ถ้า preview ใช้ได้)
- [ ] 4. ดู diff และ commit
- [ ] 5. Deploy ทันที (ถ้า critical อาจ bypass PR)
- [ ] 6. ทดสอบ production อีกครั้ง
- [ ] 7. แจ้งทีม (ถ้ามี)
- [ ] 8. กลับไปทำ PR ปกติเพื่อ merge กลับไปที่ main/dev (ถ้า bypass)

## เกณฑ์เสร็จ (Definition of Done)
- [ ] ปัญหาหายไปในการผลิต
- [ ] ไม่มีปัญหาใหม่เกิดขึ้น
- [ ] rollback plan พร้อม (git revert หรือ redeploy ก่อนหน้า)
- [ ] มีบันทึก hotfix ไว้ (ใน doc หรือ release note)

## การป้องกันไม่ให้เกิดซ้ำ
- [ ] เพิ่ม test กรณีนี้
- [ ] อัปเดต CLAUDE.md หรือ runbook
- [ ] ปรับปรุง error handling

## หมายเหตุ
```

**English**

```markdown
# Hotfix: [Issue name]
Date: [YYYY-MM-DD]
Severity: [ ] Critical (site down) / [ ] High (major feature broken) / [ ] Medium
Status: [ ] Analyzing / [ ] Fixing / [ ] Testing / [ ] Deployed

## Symptoms
- [ ] What user sees: [error message, white screen, unable to act]
- [ ] Browser console error: [copy text]
- [ ] Network response: [status code, response body]
- [ ] Logs (if any): [link or snippet]

## Steps to reproduce
1. [Step 1]
2. [Step 2]
3. [Wrong result]

## Root cause (after analysis)
- [ ] Cause: [brief explanation]
- [ ] File/line: [path:line]

## Fix plan
- [ ] Solution: [description]
- [ ] Files to change: [list]

## Hotfix workflow
- [ ] 1. Create branch from main: `hotfix/[name]`
- [ ] 2. Fix code
- [ ] 3. Test locally and preview (if available)
- [ ] 4. Review diff and commit
- [ ] 5. Deploy immediately (bypass PR if critical)
- [ ] 6. Test production again
- [ ] 7. Notify team (if any)
- [ ] 8. Create normal PR to merge back to main/dev (if bypassed)

## Definition of Done
- [ ] Issue resolved in production
- [ ] No new issues introduced
- [ ] Rollback plan ready (git revert or previous redeploy)
- [ ] Hotfix documented (in doc or release note)

## Prevention
- [ ] Add test for this case
- [ ] Update CLAUDE.md or runbook
- [ ] Improve error handling

## Notes
```

---

## 4. Enhancement Template / เทมเพลตการปรับปรุงฟีเจอร์ที่มีอยู่

**ภาษาไทย**

```markdown
# Enhancement: [ชื่อการปรับปรุง]
วันที่: [YYYY-MM-DD]
สถานะ: [ ] ยังไม่เริ่ม / [ ] กำลังทำ / [ ] รอรีวิว / [ ] เสร็จ
ปรับปรุงจาก: [ฟีเจอร์เดิมหรือไฟล์]

## ปัญหาหรือข้อจำกัดเดิม
- [ ] อธิบายสิ่งที่ยังไม่ดี: [เช่น ช้า, UX งง, error ไม่ชัดเจน]
- [ ] หลักฐาน: [screenshot, metric, log]

## การปรับปรุงที่ต้องการ
- [ ] เป้าหมาย: [เร็วขึ้น 30%, ลดขั้นตอน, error ชัดเจนขึ้น]
- [ ] การเปลี่ยนแปลงที่เสนอ:
  - [ ] เปลี่ยน logic ส่วน ...
  - [ ] เพิ่ม validation ...
  - [ ] ปรับ UI ...

## สิ่งที่ไม่เปลี่ยน (Keep unchanged)
- [ ] API signature เดิม (ถ้าใช้อยู่)
- [ ] ข้อมูลใน DB ไม่เปลี่ยนโครงสร้าง (หรือ migration รองรับ)
- [ ] UX flow หลักยังเหมือนเดิม

## ผลกระทบ (Impact)
- [ ] ฟีเจอร์ที่อาจได้รับผลกระทบ: [list]
- [ ] ต้องอัปเดต doc หรือไม่: [ ] ใช่ / [ ] ไม่
- [ ] ต้องอัปเดต test หรือไม่: [ ] ใช่ / [ ] ไม่

## ขั้นตอนการทำ
- [ ] 1. backup หรือ snapshot ก่อน (ถ้าสำคัญ)
- [ ] 2. สร้าง branch: `enhance/[ชื่อ]`
- [ ] 3. ทำการปรับปรุงทีละส่วน
- [ ] 4. ทดสอบทั้ง regression (ของเดิมไม่พัง) และ enhancement ใหม่
- [ ] 5. ทดสอบ Preview URL
- [ ] 6. ขอรีวิว (โดยเฉพาะถ้าเปลี่ยน logic สำคัญ)
- [ ] 7. Merge และทดสอบ production
- [ ] 8. อัปเดตเอกสาร (README, CLAUDE.md, หรือ wiki)

## เกณฑ์วัดความสำเร็จ (Success metrics)
- [ ] ก่อนปรับปรุง: [ค่า baseline เช่น load time 2s]
- [ ] หลังปรับปรุง: [ค่าเป้าหมาย เช่น load time <1s]
- [ ] วัดผลจริงหลัง deploy: [ ]

## เกณฑ์เสร็จ
- [ ] ทุกอย่างที่ปรับปรุงทำงานตามเป้า
- [ ] ฟีเจอร์เดิมยังทำงานปกติ (regression test ผ่าน)
- [ ] ไม่มี error ใหม่เกิดขึ้น
- [ ] Diff ผ่านการตรวจสอบ

## หมายเหตุ
```

**English**

```markdown
# Enhancement: [Improvement name]
Date: [YYYY-MM-DD]
Status: [ ] Not started / [ ] In progress / [ ] Review / [ ] Done
Based on: [Original feature or file]

## Current problem or limitation
- [ ] Describe what's not good: [e.g., slow, confusing UX, unclear error]
- [ ] Evidence: [screenshot, metric, log]

## Desired improvement
- [ ] Goal: [30% faster, fewer steps, clearer errors]
- [ ] Proposed changes:
  - [ ] Change logic in ...
  - [ ] Add validation ...
  - [ ] Adjust UI ...

## Keep unchanged
- [ ] Original API signature (if used)
- [ ] DB schema unchanged (or migration is backward compatible)
- [ ] Main UX flow remains same

## Impact
- [ ] Affected features: [list]
- [ ] Need to update docs: [ ] Yes / [ ] No
- [ ] Need to update tests: [ ] Yes / [ ] No

## Workflow steps
- [ ] 1. Backup or snapshot before starting (if critical)
- [ ] 2. Create branch: `enhance/[name]`
- [ ] 3. Implement changes incrementally
- [ ] 4. Test both regression (old still works) and new enhancement
- [ ] 5. Test Preview URL
- [ ] 6. Request review (especially if core logic changed)
- [ ] 7. Merge and test production
- [ ] 8. Update documentation (README, CLAUDE.md, or wiki)

## Success metrics
- [ ] Before: [baseline value, e.g., load time 2s]
- [ ] After: [target value, e.g., load time <1s]
- [ ] Actual measured after deploy: [ ]

## Definition of Done
- [ ] Enhancement meets goal
- [ ] Original feature still works (regression tests pass)
- [ ] No new errors introduced
- [ ] Diff reviewed

## Notes
```

---

## วิธีใช้งาน / How to use

1. **เลือก template** ที่ตรงกับงานที่ต้องทำ
2. **ก็อปปี้เนื้อหา** ไปยังไฟล์ `.md` หรือ粘贴ใน Claude / chat
3. **กรอกข้อมูล** ในช่อง `[ ]` และ `[ข้อความ]`
4. **ทำตาม workflow steps** และติ๊กเมื่อทำเสร็จ
5. **เมื่อจบงาน** ให้ติ๊กใน Definition of Done