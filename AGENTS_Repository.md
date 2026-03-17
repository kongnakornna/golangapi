# คู่มือการใช้งาน Repository (Repo)

## โครงสร้างโปรเจกต์และการจัดระเบียบโมดูล
จุดเริ่มต้นของ service อยู่ที่ `cmd/app/main.go`  
การประกอบแอปพลิเคชันอยู่ที่ `internal/apps/app/bootstrap/`  
การประกอบเส้นทาง (routing) อยู่ที่ `internal/apps/app/router/`  
โมดูลหลัก (domain) อยู่ที่ `internal/core/` (เช่น `auth`, `user`, `health`)  
แต่ละโมดูลจะมีโครงสร้างย่อย เช่น `handler`, `service`, `repository`, `dto`, `model`, `routes`  
โครงสร้างพื้นฐาน (infrastructure) อยู่ที่ `internal/platform/` (เช่น `config`, `db`)  
เครื่องมือสำหรับ transport layer อยู่ที่ `internal/transport/` (เช่น `httpx`, `middleware`)  
ส่วนประกอบที่สามารถนำมาใช้ซ้ำได้ (reusable components) อยู่ที่ `pkg/` (cache, jwt, logger, queue, transaction, utils)  
เอกสาร API ที่ถูกสร้างจะอยู่ที่ `api/app/`  
การตั้งค่าอยู่ที่ `configs/`  
SQL migrations อยู่ที่ `migrations/`  
ทรัพยากรสำหรับ deploy อยู่ที่ `deploy/`  
สคริปต์ต่างๆ อยู่ที่ `scripts/`

## คำสั่ง build, test และพัฒนา
- `./scripts/dev.sh` ใช้สร้าง Swagger และใช้ `air` สำหรับ hot reload (ดู `.air.toml`)
- `go run cmd/app/main.go` สำหรับรัน service โดยตรง
- `./scripts/swagger.sh` ใช้สร้าง Swagger ใหม่ไปยัง `api/app/`
- `./scripts/build.sh` ใช้ build binary สำหรับหลาย ๆ แพลตฟอร์มไปยัง `build/`
- `go test ./...` ใช้รัน tests ทั้งหมด; สามารถเพิ่ม `-race` หรือ `-coverprofile=coverage.out`
- `docker compose -f deploy/docker/docker-compose.yaml up -d` ใช้เริ่ม PostgreSQL และ Redis

## รูปแบบการเขียนโค้ดและการตั้งชื่อ
ใช้ `gofmt` (ย่อหน้า (indent) ด้วย tab)  
ชื่อ package เป็นตัวพิมพ์เล็ก, ชื่อไฟล์ใช้ snake_case (เช่น `user_service.go`, `auth_handler.go`)  
identifier ที่ export ใช้ PascalCase, ที่ไม่ export ใช้ camelCase  
เมื่อเพิ่ม type ใหม่ ควรใส่ไว้ในโมดูลหลักที่เกี่ยวข้องก่อน และรักษาการแบ่งชั้น (layer) handler/service/repository ให้ชัดเจน

## ข้อกำหนดการทดสอบ
tests จะอยู่ไดเรกทอรีเดียวกับโค้ด, ชื่อไฟล์ `*_test.go`, ชื่อฟังก์ชัน `TestXxx`  
โปรเจกต์ใช้ `testify` (ดูตัวอย่างได้ที่ `internal/core/user/service/user_service_test.go`)  
เมื่อแก้ไข logic ของ service สามารถรัน: `go test ./internal/core/user/service/`

## ข้อกำหนดการ commit และ PR
ข้อความ commit ให้ใช้รูปแบบ prefix ตามที่เคยใช้ (เช่น `feat: ...`, `fix: ...`)  
PR ควรมีคำอธิบายสั้น ๆ, คำสั่งสำหรับทดสอบ, และระบุว่ามีการอัปเดตการตั้งค่า/migration/Swagger หรือไม่ (หากมีการเปลี่ยนแปลงที่ `api/app/` กรุณาระบุด้วย)

## การตั้งค่าและความปลอดภัย
สำหรับการพัฒนาในเครื่อง ให้คัดลอก `configs/config.example.yaml` เป็น `configs/config.yaml` และสามารถใช้ environment variable ที่ขึ้นต้นด้วย `APP_` เพื่อแทนที่ค่าได้  
อย่า commit คีย์สำคัญ (secret); สำหรับ production ให้ใช้ `configs/config.production.yaml` หรือ environment variable ในการแทนที่  
หากมีการแก้ไข route หรือโครงสร้าง request/response กรุณารัน `./scripts/swagger.sh` เพื่ออัปเดตเอกสาร


### คู่มือเกณฑ์การทำ Code Review
---------------
####  1. หลักการพื้นฐาน

### 1.1 จุดประสงค์ของ Code Review
- **ปรับปรุงคุณภาพโค้ด** - ทำให้โค้ดอ่านง่าย, บำรุงรักษาง่าย
- **แบ่งปันความรู้** - เพิ่มความเข้าใจในระบบ across team
- **ตรวจสอบความถูกต้อง** - ลด bugs และ security issues
- **รักษามาตรฐาน** - ให้โค้ดสอดคล้องกับ coding standards ของทีม

### 1.2 Mindset ที่ควรมี
- **เป็นผู้ช่วย ไม่ใช่ผู้พิพากษา** - ใช้คำถามมากกว่าการสั่ง
- **เคารพผู้เขียน** - โฟกัสที่โค้ด ไม่ใช่ตัวบุคคล
- **เปิดใจรับฟัง** - ทั้งผู้ review และผู้เขียน
- **เน้นการเรียนรู้** - ทุก review เป็นโอกาสในการเรียนรู้

## 2. เกณฑ์การประเมินหลัก

### 2.1 ความถูกต้อง (Correctness)
- [ ] **Logic ถูกต้อง** - โค้ดทำงานตาม requirement จริงหรือไม่
- [ ] **Edge cases** - จัดการกับ corner cases อย่างเหมาะสม
- [ ] **Error handling** - มีการจัดการข้อผิดพลาดที่ครอบคลุม
- [ ] **Business logic** - ตรงตาม business rules

### 2.2 ความปลอดภัย (Security)
- [ ] **Input validation** - validate input ทุกช่องทาง
- [ ] **Authentication/Authorization** - ตรวจสอบสิทธิ์อย่างเหมาะสม
- [ ] **Data protection** - ไม่ expose sensitive data
- [ ] **SQL injection/XSS** - ป้องกัน vulnerability พื้นฐาน

### 2.3 ประสิทธิภาพ (Performance)
- [ ] **Algorithm efficiency** - ใช้ algorithm ที่เหมาะสม
- [ ] **Database queries** - query มีประสิทธิภาพ (ใช้ index, ไม่มี N+1)
- [ ] **Memory usage** - ไม่มี memory leak
- [ ] **Response time** - อยู่ในเกณฑ์ที่ยอมรับได้

### 2.4 การทดสอบ (Testing)
- [ ] **Unit tests** - มี test coverage ที่เหมาะสม
- [ ] **Test cases** - ครอบคลุมทั้ง happy path และ edge cases
- [ ] **Test readability** - test อ่านเข้าใจง่าย
- [ ] **Integration tests** - (ถ้าจำเป็น) สำหรับ critical flows

### 2.5 การออกแบบ (Design)
- [ ] **Separation of concerns** - แต่ละ module/class มีหน้าที่ชัดเจน
- [ ] **SOLID principles** - ใช้ principles พื้นฐานอย่างเหมาะสม
- [ ] **Design patterns** - ใช้ patterns เมื่อเหมาะสม (แต่ไม่ over-engineer)
- [ ] **Dependencies** - การพึ่งพาระหว่าง components มีเหตุผล

## 3. คุณภาพของโค้ด

### 3.1 การอ่านเข้าใจ (Readability)
- [ ] **การตั้งชื่อ** - ตัวแปร, ฟังก์ชัน, class ชื่อสื่อความหมาย
- [ ] **ความซับซ้อน** - ฟังก์ชันไม่ยาวเกินไป (แนะนำ < 20-30 lines)
- [ ] **Comment** - มี comment เมื่อจำเป็น (อธิบาย why ไม่ใช่ what)
- [ ] **Consistency** - สอดคล้องกับ style ของโปรเจค

### 3.2 การบำรุงรักษา (Maintainability)
- [ ] **Duplication** - ไม่มี code ซ้ำซ้อน (DRY principle)
- [ ] **Complexity** - cyclomatic complexity ไม่สูงเกินไป
- [ ] **Modularity** - แยกส่วนที่ reuse ได้
- [ ] **Configuration** - hard-coded values น้อยที่สุด

### 3.3 Coding Standards
- [ ] **Formatting** - ตาม convention ของภาษาและทีม
- [ ] **Language features** - ใช้ features ใหม่เมื่อเหมาะสม
- [ ] **Best practices** - ตาม community standards
- [ ] **Linting rules** - ผ่าน linting rules ที่กำหนด

## 4. กระบวนการ Review

### 4.1 สำหรับผู้ Review
**ควรทำ:**
- Review ภายใน timeline ที่กำหนด (แนะนำ < 24 ชม.)
- ให้ feedback ที่เป็น constructive
- ชื่นชมจุดที่ดีของโค้ด
- ถามคำถามเมื่อไม่เข้าใจ
- ตรวจสอบทั้ง implementation และ tests

**ไม่ควรทำ:**
- คิดแทนผู้เขียน (micromanage)
- จับผิดเรื่อง formatting เล็กน้อย
- บังคับให้ refactor โดยไม่มีเหตุผลสำคัญ
- ล่าช้าโดยไม่จำเป็น

### 4.2 สำหรับผู้ขอ Review
**ควรเตรียม:**
- คำอธิบายการเปลี่ยนแปลง (PR description)
- Link ไปถึง requirement/ticket
- Test instructions (ถ้าจำเป็น)
- ระบุจุดที่ต้องการ feedback เป็นพิเศษ

**ควรตอบสนอง:**
- พิจารณา feedback อย่างเปิดใจ
- อธิบาย reasoning เมื่อไม่เห็นด้วย
- ขอ clarification เมื่อไม่เข้าใจ feedback
- ขอบคุณ reviewer

## 5. Checklist อย่างรวดเร็ว

### 5.1 Pre-review
- [ ] โค้ด compile/run ได้
- [ ] Tests ผ่านทั้งหมด
- [ ] CI/CD pipeline ผ่าน
- [ ] Documentation updated (ถ้าจำเป็น)

### 5.2 During Review
- [ ] Security issues
- [ ] Performance problems
- [ ] Major bugs
- [ ] Architectural concerns
- [ ] Test coverage

### 5.3 Before Approval
- [ ] All comments addressed/resolved
- [ ] No regression introduced
- [ ] Meets acceptance criteria
- [ ] Ready for deployment

## 6. การให้ Feedback ที่มีประสิทธิภาพ

### 6.1 เทคนิคการเขียน comment
- **ใช้คำถาม**: "คุณคิดว่า approach นี้จะทำงานกับ case X อย่างไร?"
- **ให้ตัวอย่าง**: "อาจลองใช้วิธีนี้: `const result = data.filter(x => x.active)`"
- **อ้างอิง standards**: "ตาม style guide เราใช้ camelCase สำหรับตัวแปร"
- **ระบุระดับความสำคัญ**: 
  - MUST (critical): ต้องแก้ก่อน merge
  - SHOULD (important): แนะนำให้แก้
  - COULD (minor): optional improvement

### 6.2 การจัดการความเห็นต่าง
1. **อภิปรายด้วยข้อมูล** - อ้างอิง facts, metrics, benchmarks
2. **ปรึกษา third party** - ขอความเห็นจาก senior/tech lead เมื่อจำเป็น
3. **พิจารณา trade-offs** - ทุกทางเลือกมีข้อดีข้อเสีย
4. **ตัดสินใจและเดินต่อไป** - ไม่ติดอยู่กับ perfectionism

## 7. Metrics และการปรับปรุง

### 7.1 วัดประสิทธิภาพ
- **Cycle time** - เวลาจากเปิด PR ถึง merge
- **Review depth** - จำนวน comments และ discussion
- **Defect rate** - bugs ที่พบหลัง deployment
- **Team satisfaction** - survey ความพึงพอใจ

### 7.2 การปรับปรุงกระบวนการ
- **Regular retrospectives** - พูดคุยปรับปรุงกระบวนการ
- **Pair reviewing** - สำหรับ complex changes
- **Rotation** - หมุนเวียน reviewer
- **Training** - แบ่งปัน best practices

---

### สรุป
#### Code review ที่ดีคือ **การสนทนาระหว่าง professionals** เพื่อสร้างโค้ดที่ดีที่สุดเท่าที่จะเป็นไปได้ เน้นที่:
1. **คุณภาพ** มากกว่าความเร็ว
2. **การเรียนรู้** มากกว่าการจับผิด
3. **ความร่วมมือ** มากกว่าการแข่งขัน

#### **Remember:** เราทุกคนอยู่ในทีมเดียวกัน เป้าหมายคือสร้าง software ที่ดี ไม่ใชว่าใครเก่งกว่าใคร
---------------
# Code Review Criteria Handbook

## 1. Fundamental Principles

### 1.1 Purpose of Code Review
- **Improve code quality** - Make code readable, maintainable
- **Share knowledge** - Increase system understanding across the team
- **Verify correctness** - Reduce bugs and security issues
- **Maintain standards** - Ensure code aligns with team coding standards

### 1.2 Recommended Mindset
- **Be a helper, not a judge** - Use questions more than commands
- **Respect the author** - Focus on the code, not the person
- **Be open-minded** - Both reviewers and authors should listen
- **Focus on learning** - Every review is a learning opportunity

## 2. Primary Evaluation Criteria

### 2.1 Correctness
- [ ] **Logic is correct** - Does code work according to requirements?
- [ ] **Edge cases** - Handles corner cases appropriately
- [ ] **Error handling** - Comprehensive error management
- [ ] **Business logic** - Follows business rules correctly

### 2.2 Security
- [ ] **Input validation** - Validate all input channels
- [ ] **Authentication/Authorization** - Proper permission checking
- [ ] **Data protection** - No sensitive data exposure
- [ ] **SQL injection/XSS** - Prevent basic vulnerabilities

### 2.3 Performance
- [ ] **Algorithm efficiency** - Uses appropriate algorithms
- [ ] **Database queries** - Efficient queries (use indexes, no N+1)
- [ ] **Memory usage** - No memory leaks
- [ ] **Response time** - Within acceptable limits

### 2.4 Testing
- [ ] **Unit tests** - Appropriate test coverage
- [ ] **Test cases** - Covers both happy path and edge cases
- [ ] **Test readability** - Tests are easy to understand
- [ ] **Integration tests** - (If needed) for critical flows

### 2.5 Design
- [ ] **Separation of concerns** - Each module/class has clear responsibility
- [ ] **SOLID principles** - Appropriately applies basic principles
- [ ] **Design patterns** - Uses patterns when appropriate (not over-engineered)
- [ ] **Dependencies** - Reasonable dependencies between components

## 3. Code Quality

### 3.1 Readability
- [ ] **Naming** - Variables, functions, classes have meaningful names
- [ ] **Complexity** - Functions not too long (recommended < 20-30 lines)
- [ ] **Comments** - Comments when necessary (explain why, not what)
- [ ] **Consistency** - Consistent with project style

### 3.2 Maintainability
- [ ] **Duplication** - No code duplication (DRY principle)
- [ ] **Complexity** - Cyclomatic complexity not too high
- [ ] **Modularity** - Separates reusable components
- [ ] **Configuration** - Minimal hard-coded values

### 3.3 Coding Standards
- [ ] **Formatting** - Follows language and team conventions
- [ ] **Language features** - Uses new features when appropriate
- [ ] **Best practices** - Follows community standards
- [ ] **Linting rules** - Passes defined linting rules

## 4. Review Process

### 4.1 For Reviewers
**Do:**
- Review within specified timeline (recommended < 24 hours)
- Provide constructive feedback
- Praise good aspects of the code
- Ask questions when unclear
- Review both implementation and tests

**Don't:**
- Micromanage implementation
- Nitpick minor formatting issues
- Force refactoring without important reasons
- Delay unnecessarily

### 4.2 For Authors
**Prepare:**
- Change description (PR description)
- Links to requirements/tickets
- Test instructions (if needed)
- Specify areas needing special feedback

**Respond by:**
- Considering feedback with open mind
- Explaining reasoning when disagreeing
- Requesting clarification when unclear
- Thanking reviewers

## 5. Quick Checklist

### 5.1 Pre-review
- [ ] Code compiles/runs
- [ ] All tests pass
- [ ] CI/CD pipeline passes
- [ ] Documentation updated (if needed)

### 5.2 During Review
- [ ] Security issues
- [ ] Performance problems
- [ ] Major bugs
- [ ] Architectural concerns
- [ ] Test coverage

### 5.3 Before Approval
- [ ] All comments addressed/resolved
- [ ] No regression introduced
- [ ] Meets acceptance criteria
- [ ] Ready for deployment

## 6. Effective Feedback Techniques

### 6.1 Comment Writing Techniques
- **Use questions**: "How would this approach work with case X?"
- **Provide examples**: "Could try this approach: `const result = data.filter(x => x.active)`"
- **Reference standards**: "According to our style guide, we use camelCase for variables"
- **Specify priority level**:
  - MUST (critical): Must fix before merge
  - SHOULD (important): Recommended to fix
  - COULD (minor): Optional improvement

### 6.2 Handling Disagreements
1. **Discuss with data** - Reference facts, metrics, benchmarks
2. **Consult third party** - Seek opinion from senior/tech lead when needed
3. **Consider trade-offs** - Every option has pros and cons
4. **Decide and move on** - Avoid perfectionism paralysis

## 7. Metrics and Improvement

### 7.1 Performance Measurement
- **Cycle time** - Time from PR opening to merge
- **Review depth** - Number of comments and discussions
- **Defect rate** - Bugs found after deployment
- **Team satisfaction** - Satisfaction survey

### 7.2 Process Improvement
- **Regular retrospectives** - Discuss process improvements
- **Pair reviewing** - For complex changes
- **Rotation** - Rotate reviewers regularly
- **Training** - Share best practices

---

## Summary

Good code review is **a conversation between professionals** to create the best possible code. Focus on:
1. **Quality** over speed
2. **Learning** over fault-finding
3. **Collaboration** over competition

#### **Remember:** We're all on the same team. The goal is to build great software, not to prove who is better.
---------------



# Prompt Engineering: โครงสร้างและการเขียน

## 1. ภาษาไทย

### โครงสร้างพื้นฐานของ Prompt ที่ดี
```
บทบาท/บทบาทสมมติ + ภารกิจ + ข้อกำหนดเฉพาะ + รูปแบบผลลัพธ์ + เงื่อนไขเพิ่มเติม
```

### องค์ประกอบสำคัญ
1. **บทบาท (Role)**
   ```
   "คุณเป็นผู้เชี่ยวชาญด้านการตลาดดิจิทัล..."
   "ในฐานะครูสอนวิทยาศาสตร์..."
   ```

2. **ภารกิจ (Task)**
   ```
   "เขียนเนื้อหาโพสต์ Facebook เกี่ยวกับ..."
   "วิเคราะห์ข้อมูลต่อไปนี้และสรุปประเด็นหลัก..."
   ```

3. **บริบท (Context)**
   ```
   "สำหรับธุรกิจร้านกาแฟขนาดเล็ก..."
   "เพื่อใช้สอนนักเรียนชั้นมัธยมศึกษาปีที่ 3..."
   ```

4. **รายละเอียดและข้อกำหนด (Specifications)**
   ```
   "ความยาวประมาณ 300 คำ"
   "ใช้ภาษาที่เป็นทางการ"
   "ระบุแหล่งที่มา 3 แหล่ง"
   ```

5. **รูปแบบผลลัพธ์ (Output Format)**
   ```
   "จัดรูปแบบเป็น bullet points"
   "สรุปในตาราง"
   "เขียนเป็นเรียงความ 5 ย่อหน้า"
   ```

### ตัวอย่าง Prompt ภาษาไทย
```
"ในฐานะนักโภชนาการ กรุณาอธิบายประโยชน์ของอาหารเมดิเตอร์เรเนียนสำหรับผู้สูงอายุ 
โดยเน้นที่ผลต่อสุขภาพหัวใจ ความยาวประมาณ 400 คำ ใช้ภาษาที่เข้าใจง่าย 
และสรุปเป็นข้อๆ 5 ข้อท้ายบทความ"
```

### เคล็ดลับการเขียน
- ระบุความชัดเจนมากกว่าเป็นทั่วไป
- ให้ตัวอย่างถ้าต้องการรูปแบบเฉพาะ
- กำหนดขอบเขตและข้อจำกัด
- ทดลองปรับปรุง prompt หลายครั้ง

## 2. ภาษาอังกฤษ

### Basic Prompt Structure
```
Role + Task + Context + Specifications + Output Format
```

### Key Components
1. **Role Definition**
   ```
   "You are an expert in digital marketing..."
   "As a data scientist with 10 years of experience..."
   ```

2. **Clear Task**
   ```
   "Write a product description for..."
   "Analyze the following dataset and identify trends..."
   ```

3. **Context Provision**
   ```
   "For a startup targeting Gen Z consumers..."
   "In an academic research context..."
   ```

4. **Detailed Specifications**
   ```
   "Use simple language suitable for beginners"
   "Include 5 key takeaways"
   "Limit to 500 words"
   ```

5. **Output Format**
   ```
   "Format as a JSON object"
   "Create a markdown table"
   "Structure as an executive summary"
   ```

### Example English Prompt
```
"As a financial analyst, create an investment risk assessment for renewable energy stocks. 
Consider market volatility, regulatory changes, and technological disruption. 
Present in a structured report with: 1) Executive summary, 2) Risk categories, 
3) Mitigation strategies, 4) Recommendations. Use professional tone and include data points where relevant."
```

### Prompt Writing Techniques
1. **Zero-shot Prompting**
   ```
   "Translate this paragraph to French."
   ```

2. **Few-shot Prompting** (providing examples)
   ```
   "Example 1: [input] -> [output]
    Example 2: [input] -> [output]
    Now process this: [new input]"
   ```

3. **Chain-of-Thought Prompting**
   ```
   "Explain your reasoning step by step."
   "Let's think through this problem systematically."
   ```

### Best Practices
- Be specific and unambiguous
- Use delimiters for complex inputs
- Specify the desired length and depth
- Iterate and refine based on results
- Break complex tasks into subtasks

### Advanced Techniques
- **Temperature setting** (for creativity vs. consistency)
- **System prompts** for setting behavior parameters
- **Template prompts** for reproducible results
- **Meta-prompts** for generating better prompts

ทั้งสองภาษามีหลักการเดียวกัน แต่ต้องคำนึงถึงลักษณะเฉพาะของภาษาและบริบทวัฒนธรรมในการเขียน prompt ที่มีประสิทธิภาพ****
