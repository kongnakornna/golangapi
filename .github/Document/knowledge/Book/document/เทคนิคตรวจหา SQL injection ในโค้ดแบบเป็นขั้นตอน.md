<img src="https://r2cdn.perplexity.ai/pplx-full-logo-primary-dark%402x.png" style="height:64px;margin-right:32px"/>

# เทคนิคตรวจหา SQL injection ในโค้ดแบบเป็นขั้นตอน

การตรวจหา **SQL Injection (SQLi)** เป็นทักษะสำคัญของ Reviewer เพราะถ้าหลุดไปคือความเสียหายใหญ่หลวง นี่คือเทคนิคการตรวจแบบเป็นขั้นตอน (Step-by-Step) ที่คุณสามารถทำตามได้ทันที โดยเน้นบริบทของ **Node.js / NestJS / TypeORM** ที่คุณใช้งานอยู่ครับ

***

### ขั้นตอนที่ 1: ค้นหาจุดที่มีการ "ประกอบร่าง" คำสั่ง SQL

Reviewer ต้องกวาดสายตาหา "จุดเสี่ยง" ที่มีการสร้าง SQL String ขึ้นมาเอง โดยค้นหา Keyword เหล่านี้ในโค้ด:

* **Keywords:** `query(`, `execute(`, `raw(`, `createQueryBuilder`
* **Symbols:** เครื่องหมายบวก `+` หรือ Template Literal `${...}` ที่อยู่ใน String SQL


### ขั้นตอนที่ 2: ตรวจสอบที่มาของตัวแปร (Trace the Input)

เมื่อเจอจุดเสี่ยงแล้ว ให้ไล่ดูว่าตัวแปรที่เอามาใส่ มาจากไหน?

* **มาจาก User โดยตรง?** (`req.body`, `req.query`, `req.params`, `args`) -> 🚨 **อันตรายสูงสุด**
* **มาจาก Internal Logic?** (ค่าคงที่, ตัวแปรที่ระบบสร้างเอง) -> ✅ ปลอดภัย


### ขั้นตอนที่ 3: เช็ครูปแบบการเขียน (Syntax Analysis)

ดูว่า Developer เขียน SQL แบบไหน?

#### ❌ แบบที่ 1: String Concatenation (อันตรายมาก!)

เอา String มาต่อกันดื้อๆ ถ้าเจอแบบนี้ **Reject ทันที**

```typescript
// ❌ DANGER: User ส่งค่า "admin' OR '1'='1" มาก็พังหมด
const sql = "SELECT * FROM users WHERE name = '" + req.body.name + "'";
await connection.query(sql);
```


#### ❌ แบบที่ 2: Template Literal (อันตรายเหมือนกัน!)

ถึงจะดูทันสมัย แต่ถ้าใส่ตัวแปรตรงๆ ก็ไม่รอด

```typescript
// ❌ DANGER: เหมือนข้างบน แค่เขียนสวยกว่า
const sql = `SELECT * FROM users WHERE name = '${req.body.name}'`;
await connection.query(sql);
```


#### ✅ แบบที่ 3: Parameterized Query (ปลอดภัย)

ใช้ตัวแทน (`?` หรือ `$1` หรือ `:name`) แล้วส่งค่าแยกต่างหาก

```typescript
// ✅ SAFE: Database จะมอง input เป็น "ข้อมูล" ไม่ใช่ "คำสั่ง"
await connection.query('SELECT * FROM users WHERE name = $1', [req.body.name]);
```


***

### ขั้นตอนที่ 4: ตรวจสอบการใช้ ORM (TypeORM Check)

แม้จะใช้ ORM ก็อาจพลาดได้ ให้เช็คดังนี้:

#### ❌ จุดตายของ TypeORM: `user input` ในฟังก์ชันดิบ

ถ้า Developer เผลอส่ง Object ที่คุมไม่ได้เข้าไปในฟังก์ชันบางตัว

```typescript
// ❌ DANGER: TypeORM เวอร์ชันเก่า (ก่อน 0.3.0) มีช่องโหว่ถ้าส่ง req.body เข้าไปตรงๆ
// User อาจส่ง { "id": 1, "password": { "$gt": "" } } เพื่อ bypass ได้
userRepository.findOne(req.body); 
```

**วิธีแก้:** ต้องระบุ field ชัดเจน

```typescript
// ✅ SAFE
userRepository.findOne({ where: { id: req.body.id } });
```


#### ❌ การใช้ `createQueryBuilder` แบบผิดๆ

```typescript
// ❌ DANGER: ต่อ String ใน .where()
createQueryBuilder("user")
  .where("user.name = '" + name + "'") 

// ✅ SAFE: ใช้ Parameters Object
createQueryBuilder("user")
  .where("user.name = :name", { name: name }) 
```


***

### 📝 สรุป Checklist สำหรับ Reviewer (Copy ไปใช้ได้เลย)

1. [ ] **No String Concatenation:** ห้ามใช้ `+` หรือ `${}` ต่อ SQL String เด็ดขาด
2. [ ] **Use Parameters:** ทุกค่าที่มาจาก User ต้องส่งผ่าน Parameter Array (`$1`, `?`, `:name`) เท่านั้น
3. [ ] **Validate Input:** ตัวแปรที่รับมาเป็น `string` หรือ `number` จริงไหม? (ไม่ใช่ Object แฝงตัวมา)
4. [ ] **ORM Usage:** ไม่ส่ง `req.body` เข้าไปใน ORM Method ตรงๆ ต้องแกะค่าออกมาทีละตัว (Explicit Mapping)

### 🛠 เครื่องมือช่วย (Automated Tools)

ถ้าไม่อยากตรวจตาเปล่า ให้ใช้ Tools ช่วยสแกน:

* **ESLint Plugin:** `eslint-plugin-security` (ช่วยเตือนเมื่อมีการต่อ String ใน SQL)
* **SonarQube:** แจ้งเตือนช่องโหว่ "Security Hotspot" อัตโนมัติใน CI Pipeline
<span style="display:none">[^1][^10][^11][^12][^13][^14][^15][^16][^17][^18][^19][^2][^20][^3][^4][^5][^6][^7][^8][^9]</span>

<div align="center">⁂</div>

[^1]: https://portswigger.net/web-security/sql-injection

[^2]: https://letsdefend.io/blog/how-to-detect-sql-injection

[^3]: https://owasp.org/www-project-web-security-testing-guide/latest/4-Web_Application_Security_Testing/07-Input_Validation_Testing/05-Testing_for_SQL_Injection

[^4]: https://www.w3schools.com/sql/sql_injection.asp

[^5]: https://onenine.com/5-steps-to-secure-legacy-code-from-sql-injection/

[^6]: https://www.thesmartscanner.com/vulnerability-list/sql-injection-in-typeorm

[^7]: https://www.stackhawk.com/blog/node-js-sql-injection-guide-examples-and-prevention/

[^8]: https://learn.microsoft.com/en-us/sql/relational-databases/security/sql-injection?view=sql-server-ver17

[^9]: https://www.wiz.io/vulnerability-database/cve/cve-2022-33171

[^10]: https://www.w3schools.com/nodejs/nodejs_security.asp

[^11]: https://documentation.wazuh.com/current/proof-of-concept-guide/detect-web-attack-sql-injection.html

[^12]: https://dev.to/yoshi_yoshi/typeorm-prevent-sql-injection-with-node-js-react-typescript-in-2021-1go4

[^13]: https://zerothreat.ai/blog/securing-nodejs-web-app-from-sql-injection-attacks

[^14]: https://www.cobalt.io/learning-center/pentester-guide-to-sql-injection-sqli

[^15]: https://hackerone.com/reports/319458

[^16]: https://alexlaer.hashnode.dev/nodejs-security-unleashed-nosql-injection-attack

[^17]: https://arxiv.org/html/2506.17245v1

[^18]: https://security.snyk.io/vuln/SNYK-JS-TYPEORM-13746469

[^19]: https://stackoverflow.com/questions/15778572/preventing-sql-injection-in-node-js

[^20]: https://exchange.xforce.ibmcloud.com/vulnerabilities/159093

