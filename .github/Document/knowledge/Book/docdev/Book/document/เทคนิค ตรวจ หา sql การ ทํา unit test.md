
# เทคนิค ตรวจ หา sql การ ทํา unit test

สำหรับการ "ตรวจหา" หรือตรวจสอบ SQL ในการทำ Unit Test นั้นขึ้นอยู่กับว่าคุณต้องการตรวจสอบ **"Logic ของคำสั่ง SQL"** หรือ **"Logic ของโค้ด JavaScript/TypeScript"** ที่เรียกใช้ Database ครับ

เนื่องจากคุณใช้ **Next.js + Prisma + PostgreSQL** คำตอบจะแบ่งออกเป็น 3 ระดับตามความลึกของการทดสอบครับ:

### 1. ตรวจสอบ SQL ในระดับ Integration Test (แนะนำสูงสุด)

ในระดับนี้เท่านั้นที่ Prisma จะทำการ "สร้าง SQL จริงๆ" ออกมา หากคุณทำ Unit Test แบบ Mock (จำลอง) จะไม่มี SQL ถูกสร้างขึ้นครับ

**เทคนิค: ใช้ Prisma Event Logger เพื่อดู SQL**
วิธีนี้จะช่วยให้คุณ "เห็น" SQL ที่ถูกยิงออกไปจริงๆ และตรวจสอบ Performance ได้ (เช่น N+1 Problem)

```typescript
// ตั้งค่า PrismaClient ให้ส่ง Event เมื่อมี Query
const prisma = new PrismaClient({
  log: [
    {
      emit: 'event',
      level: 'query', // จับเฉพาะระดับ Query
    },
  ],
})

// เขียนฟังก์ชันดักจับ (Listener)
prisma.$on('query', (e) => {
  console.log('Query: ' + e.query)
  console.log('Params: ' + e.params)
  console.log('Duration: ' + e.duration + 'ms')
})

// ใน Test case ของคุณ
test('should find user', async () => {
  await prisma.user.findMany() // SQL จะถูกปรินต์ออกมาที่ Console
})
```


***

### 2. ตรวจสอบ Logic ในระดับ Unit Test (ใช้ Mock)

ในระดับนี้ เราจะ **ไม่เห็น SQL string** แต่เราจะตรวจสอบว่า **"Prisma ถูกสั่งงานด้วย parameters ที่ถูกต้องหรือไม่"** ซึ่งเป็นการตรวจ Logic แทนการตรวจ SQL ครับ

**เทคนิค: ใช้ `jest-mock-extended` หรือ `prisma-mock`**
ตรวจสอบว่าโค้ดของเราเรียกฟังก์ชัน `findMany`, `create` ด้วย `where` clause ที่ถูกต้องหรือไม่

```typescript
import { prisma } from './client' // client ที่ถูก mock แล้ว
import { mockDeep } from 'jest-mock-extended'

test('should call findMany with correct email', async () => {
  const email = 'test@example.com'
  
  // เรียกฟังก์ชันของคุณ
  await findUserByEmail(email)

  // ตรวจสอบว่า Prisma ถูกเรียกด้วย parameter ที่คาดหวังหรือไม่
  expect(prisma.user.findUnique).toHaveBeenCalledWith({
    where: { email: 'test@example.com' },
  })
})
```


***

### 3. ตรวจสอบความถูกต้องของ SQL จริง (Testcontainers)

หากคุณต้องการมั่นใจ 100% ว่า SQL นั้นทำงานได้จริง (เช่นพวก Raw Query ซับซ้อน หรือ Transaction) การใช้ Mock จะไม่เพียงพอ ต้องใช้ **Testcontainers** เพื่อรัน Database จริงใน Docker ขึ้นมาชั่วคราวขณะเทส

**เทคนิค: Testcontainers**
วิธีนี้จะรัน PostgreSQL จริงๆ ใน Docker ขึ้นมาให้เทส แล้วทำลายทิ้งเมื่อจบเทส ทำให้มั่นใจได้ว่า SQL ถูกต้องตาม Syntax ของ Postgres แน่นอน

* **ข้อดี:** เจอ Error จริงที่ Mock ไม่เจอ (เช่น constraint violation)
* **ข้อเสีย:** ช้ากว่า Unit Test ปกติ


### สรุป

* ถ้าอยาก **"เห็น"** SQL Text: ต้องใช้ **Integration Test + Prisma Log**
* ถ้าอยาก **"ตรวจ"** Logic การ query: ใช้ **Unit Test + Mock (`toHaveBeenCalledWith`)**
* ถ้าใช้ **Raw SQL (`$queryRaw`)**: ควรใช้ Integration Test กับ Database จริง (Testcontainers) เท่านั้น เพราะ Mock จะไม่สามารถบอกได้ว่า SQL string นั้นถูก Syntax หรือไม่
<span style="display:none">[^1][^10][^11][^12][^13][^14][^15][^16][^17][^18][^19][^2][^20][^21][^22][^23][^24][^25][^26][^27][^28][^29][^3][^30][^31][^32][^33][^34][^35][^36][^37][^38][^39][^4][^40][^5][^6][^7][^8][^9]</span>

<div align="center">⁂</div>

[^1]: https://www.linkedin.com/pulse/unit-testing-sql-databases-best-practices-imran-imtiaz-ijcfe

[^2]: https://www.sqlshack.com/sql-unit-testing-best-practices/

[^3]: https://www.geeksforgeeks.org/blogs/unit-testing-best-practices/

[^4]: https://stackoverflow.com/questions/145131/whats-the-best-strategy-for-unit-testing-database-driven-applications

[^5]: https://devcom.com/tech-blog/database-unit-testing-framework-for-sql-server/

[^6]: https://stackoverflow.com/questions/310307/mocking-vs-test-db

[^7]: https://peterkellner.net/2023-09-27-using-prisma-with-typescript-for-rapid-query-testing/

[^8]: https://learn.microsoft.com/en-us/sql/ssdt/troubleshooting-sql-server-database-unit-testing-issues?view=sql-server-ver17

[^9]: https://www.datafold.com/blog/7-dbt-testing-best-practices

[^10]: https://news.ycombinator.com/item?id=42552976

[^11]: https://www.prisma.io/docs/getting-started/setup-prisma/add-to-existing-project/relational-databases/querying-the-database-typescript-sqlserver

[^12]: https://stackoverflow.com/questions/14717/identifying-sql-server-performance-problems

[^13]: https://learn.microsoft.com/en-us/sql/ssdt/walkthrough-creating-and-running-a-sql-server-unit-test?view=sql-server-ver17

[^14]: https://www.reddit.com/r/node/comments/10tdb61/why_should_i_mock_a_database_for_testing_instead/

[^15]: https://www.prisma.io/docs/orm/prisma-client/using-raw-sql/raw-queries

[^16]: https://www.sqlshack.com/10-most-common-sql-unit-testing-mistakes/

[^17]: https://stackify.com/unit-testing-basics-best-practices/

[^18]: https://circleci.com/blog/unit-testing-vs-integration-testing/

[^19]: https://www.prisma.io/docs/orm/prisma-client/using-raw-sql/typedsql

[^20]: https://www.sqlservercentral.com/articles/sql-unit-testing

[^21]: https://github.com/prisma/prisma/issues/5026

[^22]: https://www.prisma.io/docs/orm/prisma-client/testing/unit-testing

[^23]: https://stackoverflow.com/questions/71051868/how-to-mock-prisma-with-jest-mock

[^24]: https://www.prisma.io/blog/testing-series-2-xPhjjmIEsM

[^25]: https://www.prisma.io/docs/orm/prisma-client/observability-and-logging/logging

[^26]: https://stackoverflow.com/questions/73654056/how-to-see-prisma-query-values

[^27]: https://www.pingcap.com/article/step-by-step-guide-to-using-prisma-with-sql-databases/

[^28]: https://www.prisma.io/docs/orm/prisma-client/testing/integration-testing

[^29]: https://ruheni.dev/writing/nestjs-prisma-client-logging/

[^30]: https://github.com/prisma/prisma/issues/5385

[^31]: https://dev.to/eddeee888/how-to-write-tests-for-prisma-with-docker-and-jest-593i

[^32]: https://www.prisma.io/docs/orm/reference/prisma-client-reference

[^33]: https://www.prisma.io/docs/orm/prisma-client/client-extensions/middleware/logging-middleware

[^34]: https://github.com/prisma/prisma/issues/2803

[^35]: https://www.prisma.io/docs/orm/prisma-client/observability-and-logging

[^36]: https://www.prisma.io/blog/prisma-orm-v6-14-0-relationships-for-sql-views-more-robust-management-api-and-more

[^37]: https://www.prisma.io/blog/testing-series-4-OVXtDis201

[^38]: https://www.prisma.io/docs/orm/prisma-client/queries

[^39]: https://www.prisma.io/docs/orm/prisma-client/client-extensions/middleware/session-data-middleware

[^40]: https://github.com/prisma/prisma/issues/16451

