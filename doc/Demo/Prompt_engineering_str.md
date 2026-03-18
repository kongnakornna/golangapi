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

ทั้งสองภาษามีหลักการเดียวกัน แต่ต้องคำนึงถึงลักษณะเฉพาะของภาษาและบริบทวัฒนธรรมในการเขียน prompt ที่มีประสิทธิภาพ