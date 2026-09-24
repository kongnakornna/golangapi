TASK: xxxxx
### รายละเอียด 
 

### หลักการทำงาน (Concept / Behavior spec)
 - ข้อกำหนด (Requirements)
    xxx
 - เป้าหมาย (Target)
   xxxx
 - ขอบเขต (Scope)
     xxx
 - วาดโครงสร้าง Foder และ ไฟล์ของ Module
 - ออกแบบ Workflow การทำงาน
 - เพื่ออธิบาย Workflow กระบวนการ เพื่อ ทำความเข้าใจ
 - Affected กับส่วนไหนบ้าง
 - แสดงรายการ ให้เห็นว่า จะทำงานกระบวการอย่างไร
 - ให้ความสำคัญกับ Performance for application ประสิทธิภาพของแอปพลิเคชัน การทำงานของโปรแกรม ว่ารวดเร็ว เสถียร และมีประสิทธิภาพมาก
 - ทำตามหลัก Test-Driven Development (TDD)
 - ข้อห้าม
  - พยามทำความเข้าใจ โครงสร้างเดิม แล้วนำมา ใช้งาน พยายามไม่แก้โครงสร้างเดิม 
  - หากมี Function เดิม ที่ใช้งานร่วมกัน ให้นำมาใช้
  -   หากมี Function ใหม่ มี Performance ดีกว่าเดิม ให้ทำใหม่
  - ข้อควรระวัง ภาษา Golang ให้ตรวจสอบ panic ด้วย ต้องระวังปัญหา panic   
  - ข้อดี
  - ข้อเสีย
- กระบวนการทดสอบผล check list
- สรุป
- รายงานสรุปผลการดำเนิดการ อะไรไปบ้าง
- หากเป็น Backend ให้สร้าง Swagger: API Documentation & Design Tools for Teams
- หากเป็น API ให้สร้าง Postman
- หากเป็น Fonfend ให้ทำคู่มือ ด้วย
 
------------------------------
skill
1. สั่ง copilot  
   - copilot
   - /tb-brainstorming
   - แล้ว เอาโจทย์มา วาง เพื่อให้  AI วิเคาะห์ 
   - /tb-writing-plans
   -  ถาม และ ให้ ทางเลือก แสดงแนวทางตัดสินใจ ไปจน กว่าจะ สมบุรณ์
   - /tb-executing-plans
   - เริ่ม ทำงาน ตาม plans


/tb-scrutinize 
- Affected กับส่วนไหนบ้าง
- panic  
- สิ่งที่ต้องการ code ใน  branches  นี้ ต้องไม่เปลียนแปลง  พถติกรรม ใน branches  dev  code ใน branches  dev ต้องใช้งานได้ โดยไม่มีผลกระทบ 
- แค่นำ code นี้ไป เพิ่มเข้าไป โดยไม่มีผลกระทบ code ใน BRANCHES dev  
- รายงานสรุปผลการดำเนิดการ อะไรไปบ้าง และบอกกระบวนการทดสอบผล


 

มี context พอแล้ว สรุปให้ฟัง:
────────────────────
Context ปัจจุบัน
•  /team-calibrate/detail  → route group มีแค่  OAuthMiddleware  (validate JWT เท่านั้น)
• Handler ( GetPerformanceTeamCalibrate ) รับ  jsonParam  query → parse  userIDs  → ถ้าว่างใช้ login user แทน
• ไม่มี access control ว่าใครส่ง  userIDs  ของคนอื่นได้ไหม
•  GetTargetByUserID(userID)  → คืน raw SQL WHERE clause สำหรับ filter employee ตาม RBP target
• Screen name proxy =  "Proxy"  (มีใช้ใน router.go line 1126)
• Concurrent employee IDs (คนเดียวกันหลาย position) → ใช้  GetConcurrentEmployeeIDs 

────────────────────
