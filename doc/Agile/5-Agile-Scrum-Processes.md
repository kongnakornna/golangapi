# Scrum Processes: Product Backlog กระบวนการสครัม: แบ็กล็อกผลิตภัณฑ์

---

## **บทนำ (Introduction)**
Scrum เป็นกรอบงาน (Framework) สำหรับการพัฒนาซอฟต์แวร์แบบ Agile ที่เน้นการทำงานเป็นทีม การปรับตัวต่อการเปลี่ยนแปลง และการส่งมอบคุณค่าให้กับลูกค้าอย่างต่อเนื่อง กระบวนการ Scrum ประกอบด้วยบทบาท (Roles) งานสร้าง (Artifacts) และกิจกรรม (Events) ที่เชื่อมโยงกัน โดย **Product Backlog** เป็นงานสร้างหลักและหัวใจสำคัญที่ขับเคลื่อนการพัฒนาใน Scrum

---

## **บทนิยาม (Definition)**
**Product Backlog (แบ็กล็อกผลิตภัณฑ์)** คือ รายการของงานหรืองานที่จำเป็นสำหรับการพัฒนาผลิตภัณฑ์ ซึ่งถูกจัดลำดับตามความสำคัญ คุณค่า ความเสี่ยง และความจำเป็น โดยเป็นแหล่งข้อมูลเดียว (Single Source of Truth) ที่รวบรวมความต้องการ ทุกอย่างสำหรับผลิตภัณฑ์นั้น **Product Owner (เจ้าของผลิตภัณฑ์)** เป็นผู้มีอำนาจและรับผิดชอบแต่เพียงผู้เดียวในการจัดการ Product Backlog ทั้งในด้านเนื้อหา การจัดลำดับความสำคัญ การทำให้พร้อมใช้งาน และการสื่อสารกับผู้มีส่วนได้ส่วนเสีย

---

## **บทหัวข้อ (Topics)**
### 1. ลักษณะของ Product Backlog (Characteristics)
- **เป็นรายการที่มีชีวิต (Living Artifact)**: ไม่เคยสมบูรณ์ (Never Complete) และเปลี่ยนแปลงได้ตลอดเวลา
- **มีการจัดลำดับ (Ordered)**: เรียงลำดับจากบนลงล่าง โดยรายการบนสุดมีความสำคัญและพร้อมสำหรับการพัฒนาใน Sprint ถัดไปมากที่สุด
- **มีรายละเอียดที่ปรับเปลี่ยนได้ (Dynamic & Emergent)**: รายละเอียดของรายการจะชัดเจนมากขึ้นเรื่อยๆ ยิ่งรายการอยู่ด้านบนยิ่งต้องมีความชัดเจนสูง (Refined)
- **ครอบคลุมทุกประเภทของงาน**: รวมถึงคุณลักษณะใหม่ (Features), การปรับปรุง (Enhancements), การแก้ไขข้อบกพร่อง (Bug Fixes), งานวิจัย (Spikes), งานทางเทคนิค (Technical Debt) และอื่นๆ

### 2. คุณลักษณะของรายการใน Product Backlog (Attributes of a Product Backlog Item - PBI)
- **คำอธิบาย (Description)**: มักเขียนในรูปแบบ User Story (As a [role], I want [feature], so that [benefit]) หรือรูปแบบอื่นๆ ที่สื่อความต้องการได้ชัดเจน
- **ลำดับความสำคัญ (Order)**: ตำแหน่งในรายการ ซึ่งสะท้อนถึงคุณค่าและความเร่งด่วน
- **การประมาณขนาด (Estimate)**: มักใช้ Story Points, วันทำงาน หรือขนาดสัมพัทธ์ เพื่อช่วยทีมในการวางแผน
- **คุณค่า (Value)**: ประโยชน์เชิงธุรกิจหรือต่อผู้ใช้ที่จะได้รับ
- **(Optional) เกณฑ์การยอมรับ (Acceptance Criteria)**: เงื่อนไขที่ต้องเป็นจริงเพื่อให้ถือว่างานนั้น "เสร็จสมบูรณ์" (Done)

### 3. การสร้างและการบริหารจัดการ Product Backlog (Creation & Management)
- **เริ่มต้นด้วย "User Stories"**: แปลงความต้องการจากผู้มีส่วนได้ส่วนเสียเป็น User Stories ที่เข้าใจง่าย
- **การปรับปรุงให้ละเอียด (Backlog Refinement)**: เป็นกิจกรรมต่อเนื่องที่ทีม Scrum และ Product Owner ร่วมกัน
    - แบ่งงานใหญ่ให้เล็กลง (Split)
    - ชี้แจงรายละเอียดให้ชัดเจน (Clarify)
    - ประเมินขนาด (Estimate)
    - ปรับลำดับความสำคัญใหม่ (Reorder)
- **ความรับผิดชอบของ Product Owner**:
    - สื่อสารวิสัยทัศน์และเป้าหมายของผลิตภัณฑ์
    - สร้างและสื่อสารรายการใน Product Backlog
    - จัดลำดับรายการใน Product Backlog
    - ทำให้มั่นใจว่า Product Backlink มีความโปร่งใส ชัดเจน และทีมเข้าใจ

### 4. ตัวอย่างลำดับขั้นตอนการพัฒนา (Example Development Flow)
1. **สร้างบทนำ** (Create Introduction)
2. **สร้างบทนิยาม** (Create Definitions)
3. **สร้างบทหัวข้อ** (Create Topics)
4. **ออกแบบคู่มือ** (Design Manual)
5. **ออกแบบเวิร์กโฟลว์** (Design Workflow)
6. **ออกแบบเทมเพลตรายการงาน (Task List Template)**
7. **ออกแบบเทมเพลตรายการตรวจสอบ (Checklist Template)**
    - *หมายเหตุ: ควรเรียงลำดับหมายเลขให้ต่อเนื่อง (เช่น 7, 8 แทนที่จะเป็น 7, 8, 7, 8)*
8. **พัฒนาไฟล์ Excel สำหรับเทมเพลต** (Develop Excel Template Files)

### 5. ข้อควรพิจารณาเพิ่มเติม (Additional Considerations)
- **ภาษาที่ใช้ (Language)**: ผลิตภัณฑ์นี้ควรได้รับการออกแบบให้รองรับ **2 ภาษา** ได้แก่ **ภาษาอังกฤษ (English)** และ **ภาษาไทย (Thai)** เพื่อให้สอดคล้องกับผู้ใช้
- **การจัดลำดับความสำคัญ**: Product Owner ต้องตัดสินใจจัดลำดับรายการต่างๆ ข้างต้น โดยพิจารณาจาก
    - คุณค่าต่อผู้ใช้/ธุรกิจ (Value)
    - ความเสี่ยง (Risk)
    - ความจำเป็นทางเทคนิคหรือการพึ่งพากัน (Dependencies)
    - ความเร่งด่วน (Urgency)

---

## **Workflow การออกแบบ (ออกแบบเวิร์กโฟลว์)**


![A person smiling in a garden](assets/Workflow.png "Workflow")
---

## **TASK LIST Template (เทมเพลตรายการงาน)**

**Product Backlog Item (PBI) / User Story:** [ชื่อเรื่อง PBI/User Story ที่นี่]
**ID:** [PB-001] **Story Points:** [5] **ลำดับความสำคัญ:** [สูง]

| Task ID (รหัสงาน) | Task Description (รายละเอียดงาน) | Assigned To (มอบหมายให้) | Estimated Effort (คน-ชม.) | Status (สถานะ) | Notes (หมายเหตุ) |
| :--- | :--- | :--- | :--- | :--- | :--- |
| T-001 | วิเคราะห์ความต้องการและเขียน Acceptance Criteria | Alice | 4 | Done |  |
| T-002 | ออกแบบ UI/UX Mockup | Bob | 8 | In Progress | รอ feedback |
| T-003 | พัฒนา Feature หลัก (Backend) | Charlie | 16 | To Do |  |
| T-004 | พัฒนา Feature หลัก (Frontend) | David | 12 | To Do |  |
| T-005 | เขียน Unit Tests | Eve | 6 | To Do |  |
| T-006 | ทดสอบ Integration | QA Team | 8 | To Do |  |
| **รวม** | | | **54** | | |

---

## **CHECKLIST Template (เทมเพลตรายการตรวจสอบ)**

**สำหรับ PBI:** [ชื่อเรื่อง PBI/User Story ที่นี่]
**ก่อนเริ่ม Sprint (Pre-Sprint Checklist)**
- [ ] PBI มีคำอธิบาย (Description) ที่ชัดเจน
- [ ] มีเกณฑ์การยอมรับ (Acceptance Criteria) ที่ทดสอบได้
- [ ] ทีมได้ทำการประเมินขนาด (Estimation) แล้ว
- [ ] งานถูกแบ่งออกเป็น Task ย่อยที่จัดการได้ (ถ้าจำเป็น)
- [ ] ข้อจำกัดหรือการพึ่งพา (Dependencies) ถูกระบุและจัดการแล้ว

**Definition of Ready (DoR) - เกณฑ์ความพร้อมสำหรับการนำไปพัฒนา**
- [ ] ความต้องการชัดเจนและทีมเข้าใจ
- [ ] มีข้อตกลงร่วมกันระหว่าง PO และทีม
- [ ] ขนาดเหมาะสมสำหรับ 1 Sprint
- [ ] เกณฑ์การยอมรับชัดเจน
- [ ] UI/UX Design พร้อม (ถ้ามี)
- [ ] ทีมมีทักษะเพียงพอหรือมีแผนที่จะได้รับทักษะนั้น

**Definition of Done (DoD) - เกณฑ์ความสำเร็จของงาน**
- [ ] โค้ดได้รับการพัฒนาแล้ว
- [ ] โค้ดผ่านการ Review แล้ว
- [ ] Unit Tests ผ่านทั้งหมดและครอบคลุม
- [ ] Integration Tests ผ่าน
- [ ] ผ่านการทดสอบจาก QA/UAT (ถ้ามี)
- [ ] เอกสารอัปเดต (ถ้าจำเป็น)
- [ ] งานถูกนำไป deploy ในสภาพแวดล้อมที่กำหนด

---
## **ลิงก์ไปยังไฟล์เทมเพลต (Links to Template Files)**
*ไฟล์ Excel เหล่านี้ควรถูกเก็บในที่ที่ทีมสามารถเข้าถึงได้ร่วมกัน (เช่น Google Drive, SharePoint, Confluence)*

- **[Product_Backlog_Task_Template.xlsx](<assets/Product_Backlog_Task_Template.xlsx>)** - เทมเพลตรายการงาน (Task List)
- **[Product_Backlog_Checklist_Template.xlsx](<assets/Product_Backlog_Checklist_Template.xlsx>)** - เทมเพลตรายการตรวจสอบ (Checklist)
---

# Flowchart: การบริหารจัดการ Product Backlog (Product Backlog Management)

![A person smiling in a garden](assets/ProductBacklog.png "Workflow")
---

## รายละเอียดขั้นตอนใน Flowchart (Flowchart Step Details)

### 1. **เริ่มต้น (Start)**
**English:** Begin Product Backlog Management Process  
**ภาษาไทย:** เริ่มต้นกระบวนการบริหารจัดการ Product Backlog

### 2. **รวบรวม Input (Gather Inputs)**
**English:** Product Owner gathers inputs from various sources  
**ภาษาไทย:** Product Owner รวบรวมข้อมูลจากแหล่งต่าง ๆ

### 3. **แหล่งที่มาของข้อมูล (Sources of Input)**
- **จากผู้มีส่วนได้ส่วนเสีย (From Stakeholders):** ความต้องการทางธุรกิจ
- **จากผู้ใช้ (From Users):** ความต้องการในการใช้งานจริง
- **จากทีมพัฒนา (From Development Team):** ข้อเสนอแนะทางเทคนิค
- **จากตลาดและคู่แข่ง (From Market/Competitors):** ข้อมูลตลาด
- **จากวิสัยทัศน์ผลิตภัณฑ์ (From Product Vision):** ทิศทางผลิตภัณฑ์ในระยะยาว

### 4. **กระบวนการจัดการ Product Backlog (Product Backlog Management Process)**

#### 4.1 **เพิ่มรายการใหม่ (Add New Items)**
**English:** Add newly identified items to the Product Backlog  
**ภาษาไทย:** เพิ่มรายการงานใหม่ที่พบเข้าไปใน Product Backlog

#### 4.2 **Product Backlog Item Pool**
**English:** Repository of all potential work items  
**ภาษาไทย:** ที่รวบรวมรายการงานที่เป็นไปได้ทั้งหมด

#### 4.3 **ทำ Backlog Refinement (Perform Backlog Refinement)**
**English:** Regularly refine and prepare backlog items  
**ภาษาไทย:** ปรับปรุงและเตรียมรายการงานใน Backlog อย่างสม่ำเสมอ

##### **กิจกรรมใน Refinement:**
- **ชี้แจงให้ชัดเจน (Clarify):** ทำให้รายละเอียดชัดเจน
- **ประมาณขนาดงาน (Estimate):** ประเมินความซับซ้อนของงาน
- **แบ่งย่อยงาน (Split):** แยกงานใหญ่เป็นงานย่อย
- **ตัดออก (Remove):** ลบรายการที่ไม่เกี่ยวข้องออก

#### 4.4 **จัดลำดับความสำคัญ (Prioritization)**
**English:** Product Owner orders items based on multiple criteria  
**ภาษาไทย:** Product Owner จัดลำดับรายการงานตามเกณฑ์ต่าง ๆ

##### **เกณฑ์การจัดลำดับ:**
- **ตามมูลค่า (Value):** สร้างมูลค่าสูงสุดให้ธุรกิจ
- **ตามความเสี่ยง (Risk):** ลดความเสี่ยงหรือเรียนรู้เร็ว
- **ตามการพึ่งพา (Dependencies):** งานที่ต้องทำก่อน
- **ตามความจำเป็น (Necessity):** งานที่จำเป็นต้องมี

### 5. **Product Backlog ที่มีลำดับแล้ว (Ordered Product Backlog)**
**English:** Backlog with clear priority order  
**ภาษาไทย:** Product Backlog ที่มีลำดับความสำคัญชัดเจน

### 6. **ประเมินความพร้อม (Assess Readiness)**
**English:** Top items assessed for Sprint Planning readiness  
**ภาษาไทย:** ประเมินว่างานส่วนบนของ Backlog พร้อมสำหรับการวางแผน Sprint หรือไม่

### 7. **พร้อมสำหรับ Sprint Planning (Ready for Sprint Planning)**
**English:** Items are clear, estimated, and ready for team selection  
**ภาษาไทย:** รายการงานมีความชัดเจน มีการประมาณงาน และพร้อมให้ทีมเลือกไปพัฒนา

### 8. **เลือกงานไปพัฒนา (Select Items for Development)**
**English:** Development Team selects items during Sprint Planning  
**ภาษาไทย:** ทีมพัฒนาเลือกรายการงานในช่วงการวางแผน Sprint

### 9. **ส่งมอบงาน (Deliver Work)**
**English:** Team delivers working increment at Sprint end  
**ภาษาไทย:** ทีมส่งมอบงานที่เสร็จสมบูรณ์ ณ สิ้นสุด Sprint

### 10. **ประเมินผลและเรียนรู้ (Evaluate and Learn)**
**English:** Product Owner evaluates results and gathers new insights  
**ภาษาไทย:** Product Owner ประเมินผลลัพธ์และรวบรวมการเรียนรู้ใหม่

##### **แหล่งเรียนรู้:**
- **จาก Increment:** ผลงานที่ส่งมอบจริง
- **จาก Feedback:** คำติชมจากผู้มีส่วนได้ส่วนเสีย
- **จากข้อมูลการตลาด:** ข้อมูลจากตลาดที่เปลี่ยนแปลง

### 11. **ป้อนกลับ (Feedback Loop)**
**English:** New insights feed back into the backlog  
**ภาษาไทย:** ข้อมูลใหม่ถูกป้อนกลับไปยัง Product Backlog

---

## **Workflow Checklist สำหรับการบริหารจัดการ Product Backlog**

| **ขั้นตอน** | **รายการตรวจสอบ** | **ผู้รับผิดชอบ** | **ความถี่** | **สถานะ** |
|------------|-----------------|----------------|------------|-----------|
| **1. รวบรวม** | - รวบรวมความต้องการจากผู้มีส่วนได้ส่วนเสียนัก<br>- สัมภาษณ์ผู้ใช้<br>- ศึกษาข้อมูลตลาด | Product Owner | ต่อเนื่อง | ☐ |
| **2. เขียน** | - เขียน User Story/รายการงาน<br>- กำหนด Acceptance Criteria<br>- เพิ่มรายการลงใน Backlog | Product Owner | เมื่อพบงานใหม่ | ☐ |
| **3. ปรับปรุง** | - ชี้แจงรายละเอียดงาน<br>- ประมาณขนาดงาน (Story Points)<br>- แบ่งงานใหญ่เป็นงานย่อย | Product Owner + Team | ทุกสัปดาห์ | ☐ |
| **4. จัดลำดับ** | - ประเมินมูลค่าของแต่ละงาน<br>- ประเมินความเสี่ยง<br>- พิจารณาการพึ่งพาระหว่างงาน | Product Owner | ก่อนแต่ละ Sprint | ☐ |
| **5. ติดตาม** | - อัปเดตสถานะงาน<br>- ปรับลำดับความสำคัญตามสถานการณ์<br>- ลบงานที่ล้าสมัย | Product Owner | ต่อเนื่อง | ☐ |

---

## **Key Metrics สำหรับ Product Backlog**

| **เมตริก** | **คำอธิบาย** | **เป้าหมาย** |
|-----------|-------------|-------------|
| **Backlog Health** | อัตราส่วนของงานที่พร้อมพัฒนา (Refined) ต่องานทั้งหมด | > 70% |
| **Average Age** | ระยะเวลาเฉลี่ยที่รายการงานอยู่ใน Backlog | < 60 วัน |
| **Value Distribution** | การกระจายของมูลค่าตามลำดับความสำคัญ | เรียงจากสูงไปต่ำ |
| **Completion Rate** | อัตราการส่งมอบงานจาก Backlog | สม่ำเสมอ |

---

**สรุป:** Flowchart นี้แสดงวงจรชีวิตของ Product Backlog ตั้งแต่การรวบรวมความต้องการ การปรับปรุง การจัดลำดับ ความสำคัญ ไปจนถึงการพัฒนาและเรียนรู้จากผลลัพธ์ ซึ่งเป็นกระบวนการที่เกิดขึ้นอย่างต่อเนื่องและวนซ้ำเพื่อให้มั่นใจว่าผลิตภัณฑ์จะพัฒนาตามทิศทางที่ถูกต้องและสร้างมูลค่าสูงสุด



*คู่มือนี้ถูกออกแบบให้เข้าใจง่ายและนำไปปฏิบัติได้จริง โดยอิงตาม Scrum Guide และปรับให้เหมาะกับการใช้งานทั้งภาษาอังกฤษและภาษาไทย*