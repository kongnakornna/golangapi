# 🎯 **Software Business Owner, Product Owner และ Agile Requirements Gathering**

## **1. บทบาทที่แตกต่างแต่เชื่อมโยงกัน**

| **บทบาท** | **ความรับผิดชอบหลัก** | **มุมมอง** | **Decision Authority** |
|-----------|---------------------|------------|-----------------------|
| **Software Business Owner** | - กำหนดวิสัยทัศน์และกลยุทธ์<br>- รับผิดชอบ ROI<br>- ดูแลงบประมาณ<br>- ตัดสินใจระดับธุรกิจ | **Why** (ทำไมต้องทำ)<br>**What** (ทำอะไรถึงจะสำเร็จ) | ระดับธุรกิจและงบประมาณ |
| **Product Owner** | - จัดลำดับความสำคัญของ backlog<br>- กำหนด requirement<br>- รับผิดชอบ product value<br>- สื่อสารกับ stakeholders | **What** (คุณสมบัติอะไร)<br>**Who** (ใครคือผู้ใช้)<br>**Value** (ค่าที่ผู้ใช้ได้รับ) | ระดับ product feature และ priority |
| **Agile Team** | - วิธีการ implement<br>- Technical feasibility<br>- Estimation<br>- การส่งมอบจริง | **How** (ทำอย่างไร)<br>**When** (ใช้เวลานานเท่าไร) | ระดับ implementation |

## **2. Agile Requirements Gathering Framework**

### **The Agile Requirements Pyramid**

```
          ⎫ Vision & Strategy (Business Owner)
          ⎬ Goals & OKRs (Business Owner + Product Owner)
    WHY → ⎭ User Problems & Needs (Product Owner)
          ⎫ Features & Epics (Product Owner + Team)
    WHAT →⎬ User Stories (Product Owner + Team)
          ⎭ Acceptance Criteria (Team + PO)
    HOW →  Tasks & Technical Details (Team)
```

## **3. เทคนิคการเก็บ Requirement แบบ Agile**

### **3.1 Discovery Phase (ก่อน Sprint)**

**เครื่องมือที่แนะนำ:**
1. **User Story Mapping** (Jeff Patton)
   ```
   Example: E-commerce Platform
   
   User Activities:
   │─ Browse Products
   │  ├─ Search items
   │  ├─ Filter by category
   │  └─ Sort results
   │
   │─ Make Purchase
   │  ├─ Add to cart
   │  ├─ Enter shipping info
   │  └─ Complete payment
   │
   │─ Track Order
   │  ├─ View order status
   │  └─ Contact support
   ```

2. **Impact Mapping** (Gojko Adzic)
   ```
   Goal: Increase checkout conversion by 20%
   │
   ├── Actors: First-time buyers, Returning customers
   │
   ├── Impacts: 
   │   ├── Reduce abandonment
   │   ├── Increase trust
   │   └── Simplify process
   │
   └── Deliverables:
       ├── Guest checkout
       ├── Multiple payment options
       └── Progress indicator
   ```

3. **Jobs To Be Done (JTBD) Framework**
   ```
   When: [สถานการณ์]
   I want to: [แรงจูงใจ]
   So I can: [ผลลัพธ์ที่ต้องการ]
   
   Example:
   เมื่อ: ฉันกำลังรีบและอยากได้อาหารมาที่บ้าน
   ฉันต้องการ: สั่งอาหารผ่านมือถือได้ใน 2 นาที
   เพื่อ: ฉันจะได้มีเวลาทำงานต่อ
   ```

### **3.2 Refinement Phase (ระหว่าง Sprint)**

**User Story Template (INVEST Principle):**
```markdown
**As a** [type of user]
**I want** [some goal]
**So that** [some reason]

**Acceptance Criteria:**
1. Given [context]
   When [action]
   Then [result]
   
2. Given [context]
   When [action]
   Then [result]

**Definition of Done:**
- [ ] Unit tests written
- [ ] Integration tests passed
- [ ] Code reviewed
- [ ] Documentation updated
- [ ] Deployed to staging
```

**ตัวอย่างจริง:**
```markdown
**As a** customer
**I want** to save items to a wishlist
**So that** I can purchase them later

**Acceptance Criteria:**
1. Given I'm logged in
   When I click "Add to Wishlist" on a product
   Then the item appears in my wishlist

2. Given I have items in my wishlist
   When I view my wishlist
   Then I can see all saved items with prices

3. Given I have an item in my wishlist
   When the item goes out of stock
   Then I receive an email notification
```

## **4. Collaboration Framework ของ 3 บทบาท**

### **4.1 Sprint Cycle Collaboration**

```
        [Business Owner]
              │
              ▼  Provides vision, market data, business constraints
        ┌─────────────┐
        │ Product     │←── Feedback loop with customers
        │  Owner      │
        └─────────────┘
              │
              ▼  Creates & prioritizes backlog
        ┌─────────────┐
        │ Agile Team  │←── Technical feasibility feedback
        │ (Dev, QA,  │
        │   UX)      │
        └─────────────┘
```

### **4.2 Meeting Structure ที่ได้ผล**

| **Meeting** | **ความถี่** | **ผู้เข้าร่วม** | **วัตถุประสงค์** | **Output** |
|------------|------------|----------------|-----------------|------------|
| **Strategic Roadmap Review** | ไตรมาส | Business Owner, PO, Key Stakeholders | Align product with business goals | Updated roadmap |
| **Product Vision Workshop** | 6 เดือน | Business Owner, PO, All stakeholders | Refresh vision and strategy | Vision document |
| **Backlog Refinement** | สัปดาห์ | PO, Dev Team | Clarify and estimate stories | Ready stories |
| **Sprint Planning** | ทุก Sprint | PO, Dev Team | Plan next sprint | Sprint backlog |
| **Sprint Review** | ทุก Sprint | All stakeholders | Demo and gather feedback | Feedback items |
| **Sprint Retrospective** | ทุก Sprint | Dev Team, PO | Improve process | Action items |

## **5. Requirement Prioritization Techniques**

### **5.1 RICE Scoring Model**
```
Score = (Reach × Impact × Confidence) / Effort

Example: Feature - "Push Notifications"
- Reach: 500 users/month
- Impact: 3 (scale 0.25-3)
- Confidence: 80% → 0.8
- Effort: 2 person-months

RICE Score = (500 × 3 × 0.8) / 2 = 600
```

### **5.2 Value vs Effort Matrix**
```
                 Effort
            Low        High
         ┌─────────┬─────────┐
   High  │  Quick  │ Major   │
Value    │  Wins   │ Projects│
         ├─────────┼─────────┤
   Low   │ Fill-Ins│ Avoid   │
         └─────────┴─────────┘
```

### **5.3 Kano Model**
```python
class KanoModel:
    def categorize_feature(self, user_survey_data):
        # Features are categorized as:
        # 1. Basic (Must-have) - ไม่มีแล้วผู้ใช้ไม่พอใจ
        # 2. Performance (More is better) - ยิ่งมีมาก ยิ่งพอใจ
        # 3. Delighters (Unexpected) - ไม่คาดหวัง แต่ดีใจที่มี
        # 4. Indifferent - ไม่สนใจว่าจะมีหรือไม่
        # 5. Reverse - มีแล้วไม่พอใจ
        
        # ตัวอย่างสำหรับ Mobile Banking App:
        features = {
            "login": "Basic",
            "transfer_money": "Basic", 
            "bill_payment": "Performance",
            "spending_analytics": "Performance",
            "voice_command": "Delighter",
            "AI_financial_advice": "Delighter"
        }
        return features
```

## **6. Communication Tools และ Artifacts**

### **6.1 Living Documents ที่ควรมี**

| **Document** | **เจ้าของ** | **Audience** | **Update Frequency** |
|-------------|------------|--------------|---------------------|
| **Product Vision** | Business Owner | All | 6-12 เดือน |
| **Product Roadmap** | Product Owner | Stakeholders | ทุกไตรมาส |
| **Product Backlog** | Product Owner | Dev Team | ต่อเนื่อง |
| **Release Plan** | Product Owner | Business Owner, Dev | ทุก Release |
| **User Personas** | Product Owner | All | เมื่อมีข้อมูลใหม่ |
| **Metrics Dashboard** | Product Owner | Business Owner | สัปดาห์/เดือน |

### **6.2 Collaboration Tools**
```
1. Requirement Management:
   - Jira + Confluence
   - Productboard
   - Aha!

2. User Research:
   - UserTesting
   - Hotjar
   - SurveyMonkey

3. Prototyping:
   - Figma
   - Miro
   - InVision

4. Analytics:
   - Mixpanel
   - Google Analytics
   - Amplitude
```

## **7. กรณีศึกษา: การเก็บ Requirement สำหรับฟีเจอร์ใหม่**

### **Scenario: เพิ่มฟีเจอร์ "Split Bill" ใน Food Delivery App**

**Business Owner Perspective:**
```
วัตถุประสงค์ทางธุรกิจ:
1. เพิ่ม Average Order Value 15%
2. ดึงกลุ่มผู้ใช้ Gen Z ให้ใช้งานมากขึ้น
3. เพิ่ม Competitive Advantage

KPIs ที่ต้องติดตาม:
- Adoption rate ของ Split Bill
- Impact ต่อ AOV
- Customer satisfaction score
```

**Product Owner Activities:**
```markdown
1. **User Research:**
   - Interviews with 20 frequent users
   - Competitive analysis (ดูว่าแอพอื่นทำอย่างไร)
   - Survey to 1000 users

2. **Define Requirements:**
   As a user dining with friends
   I want to split the bill easily
   So that I don't have to calculate manually

3. **Prioritization (ใช้ RICE):**
   - Reach: 10,000 users/month
   - Impact: 2.5 (moderate-high)
   - Confidence: 70%
   - Effort: 3 person-sprints
   Score: (10000 × 2.5 × 0.7) / 3 = 5,833
```

**Agile Team Collaboration:**
```yaml
Sprint 1:
  - Technical feasibility study
  - Create API design
  - Build prototype

Sprint 2: 
  - Develop core splitting logic
  - Create UI components
  - Integration with payment gateway

Sprint 3:
  - Testing and bug fixing
  - Performance optimization
  - Prepare for release
```

## **8. Common Pitfalls และวิธีหลีกเลี่ยง**

### **Pitfall 1: Business Owner ลงรายละเอียดเกินไป**
**Solution:**  
Business Owner ควรบอก **"อะไร"** และ **"ทำไม"**  
แต่ไม่ใช่ **"อย่างไร"**

### **Pitfall 2: Product Owner ไม่เข้าใจ technical constraints**
**Solution:**  
Regular tech talks และ involving tech lead early

### **Pitfall 3: Agile team ไม่เข้าใจ business context**
**Solution:**  
Share customer feedback, analytics data, business metrics

### **Pitfall 4: Requirements เปลี่ยนบ่อยเกินไป**
**Solution:**  
ใช้ rolling wave planning, keep options open, build flexibly

## **9. Success Metrics สำหรับ Requirements Process**

| **Metric** | **Target** | **How to Measure** |
|------------|------------|-------------------|
| **Requirement Clarity Score** | >4/5 | Survey team after refinement |
| **Change Request Rate** | <10% | (# changes after sprint start)/(total stories) |
| **Cycle Time** | 2-4 สัปดาห์ | Idea → Production |
| **Business Value Delivered** | Align with OKRs | Quarterly business review |
| **Team Satisfaction** | >80% | Retrospective feedback |

## **10. Best Practices สรุป**

### **สำหรับ Business Owner:**
1. **Communicate the Why** - ทีมต้องเข้าใจเหตุผลทางธุรกิจ
2. **Trust but Verify** - มอบอำนาจแต่ติดตามผลลัพธ์
3. **Provide Context** - Share market data, competition, strategy
4. **Be Available** - สำหรับคำถามสำคัญและตัดสินใจ

### **สำหรับ Product Owner:**
1. **Be the Bridge** - ระหว่าง business และ technical
2. **Focus on Value** - ทุก feature ต้องตอบคำถาม "ผู้ใช้ได้ประโยชน์อะไร"
3. **Embrace Learning** - Requirements เป็น hypothesis ที่ต้อง validate
4. **Prioritize Ruthlessly** - You can't do everything

### **สำหรับ Agile Team:**
1. **Ask Why** - เข้าใจ business context
2. **Provide Feasibility Feedback** - Early and often
3. **Focus on Quality** - Sustainable pace สำคัญกว่า speed
4. **Collaborate Continuously** - Don't wait for handoffs

## **📊 Template: Product Requirement Document (Lightweight)**

```markdown
# [Feature Name]

## 1. Business Context
- **Business Goal:** [อะไรที่ธุรกิจต้องการ]
- **Success Metrics:** [วัดผลอย่างไร]
- **Target Users:** [ใครคือผู้ใช้]

## 2. Problem Statement
- **Current State:** [ปัจจุบันเกิดปัญหาอะไร]
- **Pain Points:** [ผู้ใช้ไม่พอใจอะไร]
- **Opportunity:** [โอกาสคืออะไร]

## 3. Proposed Solution
- **User Stories:** [ผู้ใช้จะได้อะไร]
- **Key Features:** [ฟีเจอร์สำคัญ]
- **Non-goals:** [อะไรที่เราจะไม่ทำ]

## 4. Technical Considerations
- **Architecture Impact:** [กระทบระบบอย่างไร]
- **Dependencies:** [ต้องพึ่งพาอะไร]
- **Risks:** [ความเสี่ยงที่ต้องจัดการ]

## 5. Go-to-Market
- **Release Plan:** [วางแผน rollout อย่างไร]
- **Training Needed:** [ต้องฝึกอบรมไหม]
- **Support Plan:** [จะ support อย่างไร]
```

## **🎯 คำคมสรุป**

> "Business Owner ต้องรู้ว่าทำไม  
> Product Owner ต้องรู้ว่าอะไร  
> Agile Team ต้องรู้ว่าทำอย่างไร  
> และทั้งหมดต้องรู้ว่า **ใคร** คือผู้ใช้ที่แท้จริง"

**ความสำเร็จเกิดจาก:**  
**วิสัยทัศน์ที่ชัดเจน × Requirements ที่สมบูรณ์ × การทำงานร่วมกันอย่างมีประสิทธิภาพ**