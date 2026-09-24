##### Review Checklist Template
### Code Review Checklist
---------------

## Code Quality
- [ ] Code is readable and maintainable
- [ ] No code duplication (DRY principle)
- [ ] Functions are small and focused (< 30 lines)
- [ ] Meaningful variable/function names
- [ ] Proper error handling
- [ ] No commented-out code (unless with explanation)

## Functionality
- [ ] Implements all requirements from ticket
- [ ] Handles edge cases appropriately
- [ ] No regression in existing functionality
- [ ] Performance is acceptable

## Security
- [ ] Input validation implemented
- [ ] No sensitive data exposure
- [ ] Authentication/authorization checks
- [ ] SQL injection/XSS prevention

## Testing
- [ ] Unit tests cover main functionality
- [ ] Tests for edge cases
- [ ] Integration tests if needed
- [ ] Test coverage is adequate

## Architecture
- [ ] Follows existing patterns
- [ ] Separation of concerns
- [ ] Reasonable module dependencies
- [ ] No over-engineering

## Documentation
- [ ] Code comments where necessary
- [ ] README updates if needed
- [ ] API documentation if applicable
- [ ] Configuration changes documented
  
###### แปลภาษาไทย ใช้โครงสร้างเดิม

### Review Checklist Template

### Code Review Checklist

## Code Quality
- [ ] Code is readable and maintainable
- [ ] No code duplication (DRY principle)
- [ ] Functions are small and focused (< 30 lines)
- [ ] Meaningful variable/function names
- [ ] Proper error handling
- [ ] No commented-out code (unless with explanation)

## Functionality
- [ ] Implements all requirements from ticket
- [ ] Handles edge cases appropriately
- [ ] No regression in existing functionality
- [ ] Performance is acceptable

## Security
- [ ] Input validation implemented
- [ ] No sensitive data exposure
- [ ] Authentication/authorization checks
- [ ] SQL injection/XSS prevention

## Testing
- [ ] Unit tests cover main functionality
- [ ] Tests for edge cases
- [ ] Integration tests if needed
- [ ] Test coverage is adequate

## Architecture
- [ ] Follows existing patterns
- [ ] Separation of concerns
- [ ] Reasonable module dependencies
- [ ] No over-engineering

## Documentation
- [ ] Code comments where necessary
- [ ] README updates if needed
- [ ] API documentation if applicable
- [ ] Configuration changes documented

---------------

###### แปลภาษาไทย ใช้โครงสร้างเดิม

### แบบตรวจสอบการตรวจโค้ด

## คุณภาพโค้ด
- [ ] โค้ดอ่านเข้าใจง่ายและบำรุงรักษาได้
- [ ] ไม่มีการซ้ำซ้อนโค้ด (หลักการ DRY)
- [ ] ฟังก์ชันมีขนาดเล็กและเจาะจง (< 30 บรรทัด)
- [ ] ชื่อตัวแปร/ฟังก์ชันมีความหมายชัดเจน
- [ ] การจัดการข้อผิดพลาดอย่างเหมาะสม
- [ ] ไม่มีโค้ดที่ถูกคอมเมนต์ทิ้งไว้ (เว้นแต่จะมีคำอธิบาย)

## การทำงาน
- [ ] ดำเนินการตามข้อกำหนดทั้งหมดจาก ticket
- [ ] จัดการกับกรณีขอบเขตอย่างเหมาะสม
- [ ] ไม่เกิดการถดถอยในการทำงานเดิม
- [ ] มีประสิทธิภาพในระดับที่ยอมรับได้

## ความปลอดภัย
- [ ] มีการตรวจสอบความถูกต้องของอินพุต
- [ ] ไม่มีการเปิดเผยข้อมูลสำคัญ
- [ ] มีการตรวจสอบการยืนยันตัวตน/สิทธิ์
- [ ] ป้องกันการโจมตี SQL injection/XSS

## การทดสอบ
- [ ] การทดสอบหน่วยครอบคลุมการทำงานหลัก
- [ ] มีการทดสอบสำหรับกรณีขอบเขต
- [ ] มีการทดสอบบูรณาการหากจำเป็น
- [ ] ความครอบคลุมของการทดสอบเพียงพอ

## สถาปัตยกรรม
- [ ] ปฏิบัติตามรูปแบบที่มีอยู่
- [ ] การแยกความรับผิดชอบ
- [ ] การพึ่งพาโมดูลที่สมเหตุสมผล
- [ ] ไม่มีการออกแบบที่ซับซ้อนเกินจำเป็น

## เอกสารประกอบ
- [ ] มีคอมเมนต์ในโค้ดตามความจำเป็น
- [ ] อัปเดต README หากจำเป็น
- [ ] มีเอกสารประกอบ API หากเกี่ยวข้อง
- [ ] บันทึกการเปลี่ยนแปลงการกำหนดค่า


### PR/Code Review Template

## PR Title: [Brief descriptive title]

### 📋 Ticket/Issue
- **JIRA/Ticket:** [Link to ticket]
- **Type:** [Bug Fix/Feature/Refactor/Hotfix]

### 🎯 Changes Summary
[Brief description of what this PR does]

### 🔍 Testing Done
- [x] Unit tests added/updated
- [x] Integration tests
- [x] Manual testing scenarios:
  - [x] Test case 1
  - [x] Test case 2
- [x] Tested on: [Dev/Staging environment]

### 📸 Screenshots/Recordings (if UI changes)
[Attach screenshots or screen recordings]

### ✅ Checklist Before Review
- [x] Code follows team conventions
- [x] No console.log/debug statements
- [x] Documentation updated
- [x] No breaking changes introduced
- [x] Self-review completed

### ❓ Questions/Areas for Review
[Specific areas you want reviewers to focus on]

### 🚀 Deployment Notes
[Any special deployment instructions]

###### แปลภาษาไทย ใช้โครงสร้างเดิม

---------------
# แบบฟอร์มการตรวจสอบ PR/โค้ด

## ชื่อ PR: [ชื่อสั้นๆ ที่สื่อความหมาย]

### 📋 ทิกเก็ต/ปัญหา
- **JIRA/ทิกเก็ต:** [ลิงก์ไปยังทิกเก็ต]
- **ประเภท:** [การแก้ไขบั๊ก/ฟีเจอร์/รีแฟคเตอร์/การแก้ไขด่วน]

### 🎯 สรุปการเปลี่ยนแปลง
[คำอธิบายสั้นๆ ว่าประเด็นขอรวมโค้ดนี้ทำอะไร]

### 🔍 การทดสอบที่ดำเนินการ
- [x] เพิ่ม/อัปเดตการทดสอบหน่วย
- [x] การทดสอบบูรณาการ
- [x] สถานการณ์การทดสอบด้วยตนเอง:
  - [x] กรณีทดสอบที่ 1
  - [x] กรณีทดสอบที่ 2
- [x] ทดสอบบน: [สภาพแวดล้อม Dev/Staging]

### 📸 ภาพหน้าจอ/การบันทึกวิดีโอ (หากมีการเปลี่ยนแปลง UI)
[แนบภาพหน้าจอหรือการบันทึกหน้าจอ]

### ✅ ตรวจสอบก่อนส่งทบทวน
- [ ] โค้ดปฏิบัติตามข้อตกลงของทีม
- [ ] ไม่มี console.log/คำสั่งดีบั๊กค้างอยู่
- [ ] อัปเดตเอกสารแล้ว
- [ ] ไม่มีการเปลี่ยนแปลงที่กระทบการทำงานเดิม
- [ ] ทบทวนด้วยตนเองแล้ว

### ❓ คำถาม/ประเด็นสำหรับการทบทวน
[พื้นที่เฉพาะที่ต้องการให้ผู้ตรวจสอบโฟกัส]

### 🚀 หมายเหตุการติดตั้ง
[คำแนะนำพิเศษในการติดตั้ง]


#### Prompt engineering 

จากขอมูลที่ให้มาออกแบบระบบระบบ  backeend  nustjs +fonrtent  next js react ให้ด้วย  วาดโครงสร้างระบบด้วย

1.สร้างบทนำ
2.สร้างบทนิยาม
3.สร้างบทหัวข้อ
5.ออกแบบคู่มือ
6.ออกแบบ workflow
7.TASK LIST Template
8.CHECKLIST Template
ออกแบบ 2 ภาษา
1.Enhlish
2.ภาษาไทย

Prompt engineering คืออะไร ออกแบบ Prompt template ยกตัวอย่าง ออกแบบระบบ ระบบ ERP-CRM-IoT Monitoring System
อธิบาย 2 ภาษา
1.Enhlish
2.ภาษาไทย

#### Prompt engineering คืออะไร? แนวคิดสร้าง AI ที่มีประสิทธิภาพ
- 1.สร้างบทนำ
- 2.สร้างบทนิยาม
- 3.สร้างบทหัวข้อ
- 5.ออกแบบคู่มือ
- 6.ออกแบบ workflow
- 7.TASK LIST Template
- 8.CHECKLIST Template
- ออกแบบ 2 ภาษา
- 1.Enhlish
- 2.ภาษาไทย
 ---
 ## **Prompt Engineering และแนวคิดการสร้าง AI ที่มีประสิทธิภาพ**

## 1. บทนำ (Introduction)

Prompt Engineering เป็นศาสตร์และศิลป์ในการสื่อสารกับระบบ AI (โดยเฉพาะโมเดลภาษาขนาดใหญ่) เพื่อให้ได้ผลลัพธ์ที่แม่นยำ มีประสิทธิภาพ และตรงกับความต้องการมากที่สุด การออกแบบคำสั่งหรือคำถามที่ดี ไม่เพียงแค่ช่วยให้ได้คำตอบที่ถูกต้อง แต่ยังสามารถดึงศักยภาพสูงสุดของ AI ออกมาใช้ได้อย่างเต็มที่ คู่มือนี้จัดทำขึ้นเพื่อให้ผู้ใช้เข้าใจแนวคิดพื้นฐาน พร้อมด้วยเครื่องมือและเทมเพลตที่สามารถนำไปปรับใช้ได้จริงในการทำงานกับ AI อย่างเป็นระบบ

**Introduction**
Prompt Engineering is both a science and an art of communicating with AI systems (especially Large Language Models) to obtain the most accurate, efficient, and relevant results. Well-designed prompts or questions not only yield correct answers but also unlock the full potential of AI. This guide is designed to help users grasp fundamental concepts, along with practical tools and templates that can be systematically applied in AI-related work.

---

## 2. นิยาม (Definition)

**Prompt Engineering** คือ กระบวนการออกแบบ ปรับแต่ง และปรับปรุงข้อความนำเข้า (Prompt) ที่ส่งไปให้ระบบ AI เพื่อชี้นำให้โมเดลสร้างผลลัพธ์ที่เฉพาะเจาะจงตามเป้าหมาย ซึ่งรวมถึงเทคนิคการกำหนดบทบาท โครงสร้างคำสั่ง การให้ตัวอย่าง และการกำหนดเงื่อนไขต่างๆ

**Definition**
**Prompt Engineering** is the process of designing, refining, and optimizing the input text (Prompt) given to an AI system to guide the model in generating specific, targeted outputs. This encompasses techniques such as role assignment, command structuring, providing examples, and setting constraints.

---

## 3. หัวข้อหลัก (Core Topics)

3.1 **องค์ประกอบของ Prompt ที่ดี**
    - บทบาท (Role)
    - บริบท (Context)
    - งานหรือคำชี้แจง (Task/Instruction)
    - ข้อจำกัดหรือรูปแบบ (Constraints/Format)
    - ตัวอย่าง (Examples)

3.2 **เทคนิคขั้นสูง**
    - Chain-of-Thought (การคิดเป็นขั้นตอน)
    - Few-Shot / Zero-Shot Learning
    - Persona Pattern (การกำหนดบุคลิก)
    - Template Filling

3.3 **การปรับปรุงและทดสอบ**
    - การวนซ้ำ (Iteration)
    - A/B Testing ของ Prompt
    - การวิเคราะห์และประเมินผลลัพธ์

**Core Topics**

3.1 **Components of an Effective Prompt**
    - Role
    - Context
    - Task/Instruction
    - Constraints/Format
    - Examples

3.2 **Advanced Techniques**
    - Chain-of-Thought
    - Few-Shot / Zero-Shot Learning
    - Persona Pattern
    - Template Filling

3.3 **Refinement and Testing**
    - Iteration
    - A/B Testing of Prompts
    - Analysis and Evaluation of Outputs

---

## 4. แนวคิดสร้าง AI ที่มีประสิทธิภาพ (Concept for Building Effective AI)

(แก้ไขจากข้อ 4 ที่หายไปในคำถาม)
แนวคิดหลักอยู่ที่การมองว่า AI เป็น "ผู้ร่วมงาน" ที่ต้องสื่อสารด้วยอย่างชัดเจน มากกว่าเป็นเครื่องมือวิเศษที่ตอบทุกอย่างได้
1.  **ความชัดเจน (Clarity):** คำสั่งต้องชัดเจน กำกวมน้อยที่สุด
2.  **บริบท (Context):** ให้ข้อมูลพื้นหลังพอเหมาะเพื่อจำกัดขอบเขตคำตอบ
3.  **การแบ่งย่อย (Decomposition):** แยกงานใหญ่เป็นงานย่อยๆ ใช้ Chain-of-Thought
4.  **การวนซ้ำ (Iteration):** ปรับปรุง Prompt ตามผลลัพธ์ ค่อยๆ พัฒนา
5.  **การประเมิน (Evaluation):** กำหนดเกณฑ์ชัดเจนเพื่อวัดความสำเร็จของผลลัพธ์ AI

**Concept for Building Effective AI**
The core concept is to view AI as a "collaborator" that requires clear communication, rather than a magic tool that answers everything.
1.  **Clarity:** Instructions must be clear, with minimal ambiguity.
2.  **Context:** Provide adequate background to narrow the scope of the answer.
3.  **Decomposition:** Break down large tasks into smaller subtasks using Chain-of-Thought.
4.  **Iteration:** Refine prompts based on outputs, developing them progressively.
5.  **Evaluation:** Establish clear criteria to measure the success of AI outputs.

---

## 5. ออกแบบคู่มือ (Guide Design)

**คู่มือปฏิบัติการ Prompt Engineering 5 ขั้นตอน**
**Step 1: กำหนดเป้าหมาย** - เขียนให้ชัดเจนว่าต้องการอะไรจาก AI
**Step 2: ออกแบบโครงสร้าง Prompt** - ใช้โครงสร้าง "บทบาท > บริบท > งาน > ข้อจำกัด > ตัวอย่าง"
**Step 3: ทดลองรันและเก็บผล** - ส่ง Prompt, บันทึกผลลัพธ์
**Step 4: วิเคราะห์และปรับปรุง** - ตรวจดูว่าผลลัพธ์ตรงตามต้องการไหม ส่วนใดขาด/เกิน
**Step 5: สร้าง Template และจัดเก็บ** - นำ Prompt ที่ได้ประสิทธิภาพมาสร้างเป็นเทมเพลตสำหรับใช้ซ้ำ

**Guide Design**

**Prompt Engineering Practical Guide - 5 Steps**
**Step 1: Define Goal** - Clearly write down what you want from the AI.
**Step 2: Design Prompt Structure** - Use the structure: "Role > Context > Task > Constraints > Examples".
**Step 3: Run Experiment and Collect Results** - Send the prompt, record the output.
**Step 4: Analyze and Refine** - Check if the output meets requirements. Identify what's missing or excessive.
**Step 5: Create Template and Archive** - Turn effective prompts into reusable templates for future use.

---

## 6. ออกแบบ Workflow (Workflow Design)

```
Workflow: การทำงานกับ Prompt Engineering
┌─────────────────────────────────────────────┐
│ 1. รับโจทย์ & วิเคราะห์ความต้องการ          │
│ (Define Task & Analyze Requirements)        │
└─────────────────┬───────────────────────────┘
                  ▼
┌─────────────────────────────────────────────┐
│ 2. ออกแบบ Prompt ครั้งแรก                  │
│ (Draft Initial Prompt)                      │
└─────────────────┬───────────────────────────┘
                  ▼
┌─────────────────────────────────────────────┐
│ 3. ทดสอบและประเมินผลลัพธ์                  │
│ (Test & Evaluate Output)                    │
└─────────────────┬───────────────────────────┘
                  │
         ┌────────┴────────┐
         ▼                 ▼
    ประสบความสำเร็จ   ต้องการปรับปรุง
    (Successful)      (Needs Refinement)
         │                 │
         ▼                 ▼
┌─────────────────┐ ┌──────────────────────┐
│ 6. นำไปใช้จริง  │ │ 4. ปรับปรุง Prompt  │
│ (Deploy)        │ │ (Refine Prompt)     │
└─────────────────┘ └──────────┬───────────┘
                               │
                               ▼
                ┌──────────────────────────┐
                │ 5. ทดสอบอีกครั้ง        │
                │ (Test Again)             │
                └──────────────────────────┘
```

---

## 7. TASK LIST Template

**TASK LIST: โครงสร้างสำหรับจัดการโปรเจคต์ที่ใช้ AI**
- [ ] **ขั้นตอนที่ 1: กำหนดขอบเขตและเป้าหมาย**
    - [ ] ระบุปัญหาหรืองานหลัก
    - [ ] กำหนดผลลัพธ์สุดท้ายที่คาดหวัง
    - [ ] ระบุผู้มีส่วนได้ส่วนเสียและความต้องการ
- [ ] **ขั้นตอนที่ 2: ออกแบบและพัฒนา Prompt**
    - [ ] เขียน Prompt หลัก (Master Prompt)
    - [ ] แยกย่อยเป็น Sub-prompts (หากจำเป็น)
    - [ ] เตรียมข้อมูลบริบทและตัวอย่าง
- [ ] **ขั้นตอนที่ 3: ทดสอบและประเมินผล**
    - [ ] ทดสอบกับชุดข้อมูลตัวอย่าง
    - [ ] ประเมินผลลัพธ์ด้วยเกณฑ์ที่กำหนด
    - [ ] บันทึกจุดแข็งและจุดที่ต้องปรับปรุง
- [ ] **ขั้นตอนที่ 4: ปรับใช้และตรวจสอบ**
    - [ ] นำ Prompt ไปใช้งานจริง
    - [ ] ตรวจสอบผลลัพธ์อย่างต่อเนื่อง
    - [ ] อัปเดตและปรับปรุงตามความจำเป็น

**TASK LIST Template**

**TASK LIST: Structure for Managing AI Projects**
- [ ] **Phase 1: Define Scope and Objectives**
    - [ ] Identify core problem or task
    - [ ] Define expected final outcome
    - [ ] Identify stakeholders and requirements
- [ ] **Phase 2: Design and Develop Prompts**
    - [ ] Draft Master Prompt
    - [ ] Break down into Sub-prompts (if needed)
    - [ ] Prepare context data and examples
- [ ] **Phase 3: Test and Evaluate**
    - [ ] Test with sample datasets
    - [ ] Evaluate outputs against defined criteria
    - [ ] Document strengths and areas for improvement
- [ ] **Phase 4: Deploy and Monitor**
    - [ ] Deploy prompts for actual use
    - [ ] Monitor outputs continuously
    - [ ] Update and refine as necessary

---

## 8. CHECKLIST Template

**CHECKLIST: ตรวจสอบ Prompt ก่อนส่ง**
ก่อนส่ง Prompt ให้ AI โปรดตรวจสอบ:
- [ ] **บทบาท (Role):** กำหนดบทบาทของ AI ชัดเจนแล้ว (เช่น ผู้เชี่ยวชาญ, ผู้ช่วยเขียน)
- [ ] **บริบท (Context):** ให้ข้อมูลพื้นหลังเพียงพอสำหรับงานนี้แล้ว
- [ ] **งาน (Task):** คำสั่งกระชับ ชัดเจน ครบถ้วน
- [ ] **ข้อจำกัด (Constraints):** กำหนดขอบเขต (เช่น ความยาว, รูปแบบ, ภาษาที่ใช้)
- [ ] **ตัวอย่าง (Examples):** (ถ้าจำเป็น) มีตัวอย่างที่ชัดเจนเพื่อให้ AI เข้าใจรูปแบบที่ต้องการ
- [ ] **โทน (Tone):** ระบุโทนการสื่อสารที่ต้องการ (เป็นทางการ, ไม่เป็นทางการ)
- [ ] **การจัดรูปแบบ (Format):** ระบุรูปแบบผลลัพธ์ที่ต้องการ (เช่น รายการ, ย่อหน้า, JSON)
- [ ] **ได้ลองเขียนตอบด้วยตัวเองหรือยัง?** ลองคิดว่าเราจะตอบคำถามนี้อย่างไร เป็นวิธีตรวจสอบความชัดเจนที่ดี

**Post-Submission Checklist:**
- [ ] **ผลลัพธ์ตรงกับคำสั่งหลักไหม?**
- [ ] **มีข้อมูลที่ผิดพลาดหรือสร้างขึ้นมาเอง (Hallucination) หรือไม่?**
- [ ] **รูปแบบผลลัพธ์เป็นไปตามที่กำหนดไหม?**
- [ ] **ส่วนใดของผลลัพธ์ที่ดีและควรเก็บไว้พัฒนาพรอมต์ต่อไป?**

**CHECKLIST Template**

**CHECKLIST: Pre-Submission Prompt Review**
Before sending a prompt to AI, please verify:
- [ ] **Role:** AI's role is clearly defined (e.g., Expert, Writing Assistant).
- [ ] **Context:** Sufficient background information is provided for the task.
- [ ] **Task:** The instruction is concise, clear, and complete.
- [ ] **Constraints:** Scope is defined (e.g., length, format, language).
- [ ] **Examples:** (If needed) Clear examples are provided to illustrate the desired format.
- [ ] **Tone:** The desired communication tone is specified (e.g., formal, informal).
- [ ] **Format:** The desired output format is specified (e.g., list, paragraph, JSON).
- [ ] **Have you tried answering it yourself?** Consider how you would answer this question - a good clarity check.

**Post-Submission Checklist:**
- [ ] **Does the output align with the core instruction?**
- [ ] **Are there any inaccuracies or hallucinations?**
- [ ] **Is the output format as specified?**
- [ ] **What parts of the output are good and should be kept for future prompt refinement?**


## ตัวอย่าง 

หัวข้อ : ออกแบบระบบ ระบบ erp crm iot monitoring system
- บทบาท (Role)
- บริบท (Context)
- งานหรือคำชี้แจง (Task/Instruction)
- ข้อจำกัดหรือรูปแบบ (Constraints/Format)
- ตัวอย่าง (Examples)

3.2 เทคนิคขั้นสูง
- Chain-of-Thought (การคิดเป็นขั้นตอน)
- Few-Shot / Zero-Shot Learning
- Persona Pattern (การกำหนดบุคลิก)
- Template Filling

3.3 การปรับปรุงและทดสอบ
- การวนซ้ำ (Iteration)
- A/B Testing ของ Prompt
- การวิเคราะห์และประเมินผลลัพธ์