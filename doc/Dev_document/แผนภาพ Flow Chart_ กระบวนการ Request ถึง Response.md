<img src="https://r2cdn.perplexity.ai/pplx-full-logo-primary-dark%402x.png" style="height:64px;margin-right:32px"/>

## แผนภาพ Flow Chart: กระบวนการ Request ถึง Response ของ NestJS API

แผนภาพด้านล่างแสดงกระบวนการทำงานแบบละเอียดของ NestJS API ตั้งแต่รับ HTTP Request จนถึงส่ง Response กลับไปยัง Client โดยผ่าน Middleware, Guards, Interceptors, Pipes, Controller, Service และ Database[^1][^2]

![แผนภาพ Flow Chart แสดงกระบวนการ Request ถึง Response ของ NestJS API - ตั้งแต่การรับ HTTP Request ผ่าน Middleware, Guards, Interceptors, Validation, Controller, Service, Database จนถึงการส่ง Response กลับไปยัง Client](https://ppl-ai-code-interpreter-files.s3.amazonaws.com/web/direct-files/2b884aa2db3bb486bc4cace0ac04c358/23f6364f-b7b1-451a-ae9f-56440b708bc2/70421ecb.png)

แผนภาพ Flow Chart แสดงกระบวนการ Request ถึง Response ของ NestJS API - ตั้งแต่การรับ HTTP Request ผ่าน Middleware, Guards, Interceptors, Validation, Controller, Service, Database จนถึงการส่ง Response กลับไปยัง Client

## อธิบายแต่ละขั้นตอนในกระบวนการ

### Layer 1: Middleware (เริ่มต้นของ Request)

**Logger Middleware**
Middleware ทำงานเป็นลำดับแรกสุดก่อน routing ใช้สำหรับ logging, CORS, compression หรือการแปลง request body  ในตัวอย่างนี้ Logger Middleware บันทึกข้อมูล request เช่น HTTP method, URL, timestamp และ client IP address[^1][^2][^3]

```typescript
// ตัวอย่างการใช้งาน
@Injectable()
export class LoggerMiddleware implements NestMiddleware {
  use(req: Request, res: Response, next: NextFunction) {
    console.log(`${req.method} ${req.originalUrl} - ${new Date().toISOString()}`);
    next(); // ส่งต่อไปยังขั้นตอนถัดไป
  }
}
```


### Layer 2: Guards (Authentication \& Authorization)

**JWT Auth Guard**
ตรวจสอบความถูกต้องของ JWT token ใน Authorization header  หากไม่มี token หรือ token ไม่ถูกต้อง จะส่ง 401 Unauthorized Error ทันที  Guards ทำงานหลัง middleware แต่ก่อน interceptors และมีสิทธิ์ปฏิเสธ request ได้[^1][^2][^4]

**Roles Guard**
ตรวจสอบว่า user มี role ที่จำเป็นหรือไม่  เช่น endpoint สำหรับ admin เท่านั้น หากไม่มีสิทธิ์จะส่ง 403 Forbidden Error[^5][^1]

### Layer 3: Interceptors (Before Request)

**Cache Interceptor**
ตรวจสอบว่ามีข้อมูลใน Redis cache หรือไม่  ถ้ามี (cache hit) จะส่งข้อมูลกลับทันทีโดยไม่ต้องผ่าน controller, service หรือ database ช่วยลด response time อย่างมาก[^6][^7]

**Logging Interceptor**
เริ่มจับเวลาเพื่อวัด performance ของ request  Interceptors สามารถทำงานทั้งก่อนและหลังการเรียก route handler[^1][^2]

### Layer 4: Pipes (Validation)

**Validation Pipe**
ตรวจสอบและแปลง (transform) request data ตาม DTO (Data Transfer Object)  ใช้ decorators จาก `class-validator` เช่น `@IsString()`, `@IsEmail()`, `@IsNotEmpty()`  หากข้อมูลไม่ถูกต้องจะส่ง 400 Bad Request Error พร้อมรายละเอียด validation errors[^1][^5]

```typescript
// ตัวอย่าง DTO
export class CreateBookingDto {
  @IsNotEmpty()
  @IsUUID()
  customerId: string;

  @IsNotEmpty()
  @IsDateString()
  bookingDate: string;

  @IsString()
  serviceType: string;
}
```


### Layer 5: Controller (Route Handler)

**BookingsController**
รับ request ที่ผ่านการตรวจสอบทั้งหมดแล้ว  Controller มีหน้าที่รับ HTTP request และส่งต่อไปยัง Service layer  ไม่ควรมี business logic ใน controller[^3][^5]

```typescript
@Post()
@UseGuards(JwtAuthGuard, RolesGuard)
@Roles('user', 'admin')
async create(@Body() createBookingDto: CreateBookingDto) {
  return this.bookingsService.create(createBookingDto);
}
```


### Layer 6: Service Layer (Business Logic)

**BookingsService**
ประมวลผล business logic ทั้งหมด  เช่น:[^3][^5]

- ตรวจสอบ availability ของช่วงเวลา
- Generate unique ID
- คำนวณราคา
- ตรวจสอบเงื่อนไข business rules
- เรียกใช้ repository เพื่อติดต่อ database[^8]


### Layer 7: Repository \& Database

**TypeORM Repository**
จัดการการเข้าถึง database โดยใช้ TypeORM  รองรับ query builder, relations, transactions และ caching  ข้อมูลถูกบันทึกลง PostgreSQL database และส่งกลับมายัง service[^8][^9]

```typescript
// ตัวอย่างการใช้ Repository
const booking = this.bookingsRepository.create(createBookingDto);
const savedBooking = await this.bookingsRepository.save(booking);
```


### Layer 8: Response Processing

**Transform Interceptor**
แปลง response ให้อยู่ในรูปแบบมาตรฐาน  เช่น wrap ด้วย metadata อย่าง status code, message, timestamp:[^1][^2]

```typescript
{
  "statusCode": 201,
  "message": "Success",
  "data": { /* booking data */ },
  "timestamp": "2025-11-26T12:09:00.000Z"
}
```

**Logging Interceptor (After)**
บันทึกเวลาที่ใช้ในการประมวลผล request  เช่น "POST /api/v1/bookings 201 - 245ms"[^1]

**Cache Update**
อัพเดทข้อมูลใหม่เข้า Redis cache เพื่อให้ request ถัดไปได้ใช้ข้อมูลจาก cache[^6]

### Error Handling Flow

**Exception Filter**
จับ exceptions ที่เกิดขึ้นจากทุก layer  แปลงเป็น HTTP response ที่มีรูปแบบสอดคล้องกัน  ส่ง appropriate HTTP status code (400, 401, 403, 404, 500) พร้อม error message ที่เป็นมิตรกับผู้ใช้[^1][^2]

```typescript
{
  "statusCode": 400,
  "message": "Validation failed",
  "errors": [
    "customerId must be a UUID",
    "bookingDate must be a valid date string"
  ],
  "timestamp": "2025-11-26T12:09:00.000Z"
}
```


## สรุปลำดับการทำงาน

| Step | Component | Purpose | Can Block Request? |
| :-- | :-- | :-- | :-- |
| 1 | Middleware | Logging, preprocessing | Yes (rarely) |
| 2 | JWT Auth Guard | Verify authentication | Yes (401) |
| 3 | Roles Guard | Check authorization | Yes (403) |
| 4 | Cache Interceptor | Check cache hit | No (shortcuts) |
| 5 | Logging Interceptor | Start timer | No |
| 6 | Validation Pipe | Validate input | Yes (400) |
| 7 | Controller | Route handling | No |
| 8 | Service | Business logic | Can throw errors |
| 9 | Repository | Database access | Can throw errors |
| 10 | Transform Interceptor | Format response | No |
| 11 | Exception Filter | Handle errors | Only on errors |

ระบบนี้ใช้ layered architecture ที่แยก concerns ชัดเจน ทำให้ code มีความเป็นระเบียบ ง่ายต่อการ maintain และ test[^1][^2][^5]

<div align="center">⁂</div>

[^1]: https://victorokoye.hashnode.dev/maximizing-performance-in-nestjs-with-middlewares-guards-interceptors-and-pipes

[^2]: https://www.devcentrehouse.eu/blogs/nestjs-middleware-vs-guards-vs-interceptors/

[^3]: https://www.geeksforgeeks.org/javascript/folder-structure-of-a-nestjs-project/

[^4]: https://www.freecodecamp.org/news/how-to-add-jwt-based-authentication-in-nest-js/

[^5]: https://dev.to/tak089/nestjs-roadmap-for-2025-5jj

[^6]: https://www.linkedin.com/posts/nandhakumar1488_nestjs-redis-nodejs-activity-7393299511475036161-hEf-

[^7]: https://betterprogramming.pub/optimizing-your-apps-performance-by-using-redis-as-a-caching-layer-dc48ce7551ab

[^8]: https://dev.to/refifauzan/nestjs-with-typeorm-and-postgresql-3466

[^9]: https://docs.nestjs.com/techniques/database

