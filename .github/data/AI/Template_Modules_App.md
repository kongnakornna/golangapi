### Name XXX  Module   
> **เป้าหมาย:**
> ที่สมบูรณ์ด้วย Clean Architecture คือ แนวทางการออกแบบซอฟต์แวร์ (Software Design Approach) ที่คิดค้นโดย Robert C. Martin (หรือ Uncle Bob) ซึ่งจัดระเบียบโครงสร้างโค้ดด้วยการแบ่งออกเป็นชั้น ๆ (Layers) เพื่อแยกส่วนตรรกะทางธุรกิจ (Business Rules) ออกจากเทคโนโลยีภายนอก เช่น ฐานข้อมูล (Database) หรือ UI  
> DDD  Domain-Driven Design คือ แนวทางและปรัชญาในการออกแบบซอฟต์แวร์ที่เน้นจำลองโครงสร้างและตรรกะของโค้ดให้สอดคล้องกับ "โดเมนธุรกิจ" (Business Domain) หรือปัญหาที่ซับซ้อนของกิจการนั้น ๆ อย่างแท้จริง  
> **เหมาะสำหรับ:** นักพัฒนาที่ต้องการระบบ  lean Architecture+ DDD  Domain-Driven Design
> รายละเอียด
   - Kafka Consumer Group (หลัก)  
   - Elasticsearch Bulk Indexer  
   - WebSocket Broadcaster  
   - **LLM Consumer** (เรียก LLM แบบ Async)  
   - **Embedding Consumer** (สร้างเวกเตอร์)  
   - **Blockchain Consumer** (บันทึกข้อมูลลง Blockchain)  
   - QR Code & Payment System
   - การติดตั้งและใช้งาน
   - Business Model
   - Mobile App
   - Prompt สำหรับการขยายระบบ

**สารบัญ**
1.ภาพรวมระบบ
2.โครงสร้าง Module 
3.Domain Layer 
   - 3.1 Entities
   - 3.2 Value Objects
   - 3.3 Repository Interfaces
   - 3.4 Domain Services
   - 3.5 Domain Errors
4.Application Layer 
   - 4.1 Use Cases
   - 4.2 DTOs
5.Infrastructure Layer 
   - 5.1 Repository Implementations
   - 5.2 JWT Implementation
   - 5.3 Bcrypt Hasher Implementation
   - 5.4 Rate Limit Middleware
6.Interface Layer 
   - 6.1 HTTP Handlers
   - 6.2 Routes
   - 6.3 Middleware
7.Database Migrations 
8.System Flow 
9.Workflow Diagram 
10.การติดตั้งและใช้งาน
11.Business Model
12.ภาคผนวก 
13.Prompt สำหรับการขยายระบบ ในอนาคต
14.หากสร้าง table database ให้ใส่ Prefix (เช่น pdpa_) ให้กับชื่อตารางใน Database
    - การใส่ Prefix (เช่น pdpa_) ให้กับชื่อตารางใน Database เป็นแนวปฏิบัติที่ดี โดยเฉพาะในระบบที่มีการรวมข้อมูลหลายโมดูลไว้ในฐานข้อมูลเดียวกันครับต่อไปนี้เป็นข้อดี และ ตัวอย่างโครงสร้างตาราง
     สำหรับระบบ PDPA (Personal Data Protection Act) ที่คุณสามารถนำไปประยุกต์ใช้ได้ทันทีครับ
     💡 ข้อดีของการใส่ Prefix pdpa_จัดหมวดหมู่ชัดเจน: แยกตารางที่เกี่ยวกับกฎหมายคุ้มครองข้อมูลส่วนบุคคลออกจากตารางทั่วไป (เช่น ตารางสินค้า หรือตารางระบบสั่งซื้อ)ปลอดภัยเวลากรณี Migrate/Backup: ช่วยให้ผู้ดูแลระบบ (DBA) เลือก Backup หรือจัดการเฉพาะข้อมูลส่วนบุคคลได้ง่ายขึ้นป้องกันชื่อซ้ำ (Name Collision): หลีกเลี่ยงปัญหาชื่อตารางซ้ำกับโมดูลอื่น เช่น ตาราง logs ของระบบทั่วไป กับ pdpa_logs สำหรับเก็บประวัติการยินยอม📋 ตัวอย่างการออกแบบตารางระบบ PDPA (pdpa_xxx)ชื่อตาราง (Table Name)คำอธิบาย (Description)ตัวอย่างการใช้งานpdpa_policiesเก็บเวอร์ชันของนโยบายความเป็นส่วนตัวPrivacy Policy v1.0, v2.0pdpa_purposesเก็บวัตถุประสงค์ในการขอข้อมูลเพื่อการตลาด, เพื่อการจัดส่งสินค้า, เพื่อวิเคราะห์ข้อมูลpdpa_consentsบันทึกการให้ความยินยอมของผู้ใช้งานUser A ยินยอมรับข้อมูลการตลาดเมื่อวันที่ DD/MM/YYYYpdpa_consent_logsประวัติการกดยินยอมหรือยกเลิก (Audit Trail)บันทึกประวัติการเปลี่ยนสถานะ (Opt-in / Opt-out) เพื่อใช้เป็นหลักฐานpdpa_user_requestsคำร้องขอใช้สิทธิ์ของเจ้าของข้อมูล (DSR)คำร้องขอเข้าถึงข้อมูล, ขอให้ลบข้อมูล (Right to be Forgotten)

## 2.โครงสร้าง Module xxx
```
internal/modules/xxx/
            │
            ├── domain/                     # 🏛️ DOMAIN LAYER
            │   ├── entity/
            │   │   ├── xxx.go              # User Aggregate Root
            │   │   ├── xx.go               # Session Entity
            │   │   ├── xxx.go              # VerificationToken Entity
            │   │   └── xx.go               # Permission Entity
            │   │
            │   ├── value_object/
            │   │   ├── xxx.go              # Email Value Object
            │   │   ├── xx.go               # Password Value Object
            │   │   └── xx.go               # UserID Value Object
            │   │
            │   ├── repository/
            │   │   ├── xxx_repository.go   # Interface
            │   │   ├── xxx_repository.go   # Interface
            │   │   └── xxx_repository.go   # Interface
            │   │
            │   ├── service/
            │   │   ├── xxx_service.go               # Domain Service
            │   │   ├── xx_hasher.go                 # Interface
            │   │   └── xx_maker.go                  # Interface
            │   │
            │   └── errors/
            │       └── errors.go                   # Domain Errors
            │
            ├── application/                        # 🎯 APPLICATION LAYER
            │   ├── xx.go                           # Register UseCase
            │   ├── login.go                        # Login UseCase 
            │   ├── xxx.go                          # GetPermissions UseCase
            │   └── dto.go                          # Request/Response DTOs
            │
            ├── infrastructure/                            # 🔧 INFRASTRUCTURE LAYER
            │   ├── persistence/
            │   │   ├── postgres/
            │   │   │   ├── xx_repo_impl.go              # PostgreSQL Implementation
            │   │   │   ├── xx_repo_impl.go
            │   │   │   ├── xxx_repo_impl.go
            │   │   │   └── models.go                      # GORM Models
            │   │   └── redis/
            │   │       ├── session_repo_impl.go           # Redis Implementation
            │   │       └── cache_repo_impl.go
            │   │
            │   └── security/
            │       ├── jwt_maker.go                       # JWT Implementation
            │       ├── bcrypt_hasher.go                   # Bcrypt Implementation
            │       └── xx_middleware.go                   # xxx Middleware
            │
            └── interfaces/                                # 🌐 INTERFACE LAYER
                ├── http/
                │   ├── xx_handler.go                    # HTTP Handlers
                │   ├── xxx_handler.go
                │   ├── routes.go                          # Route Registration
                │   └── dto.go                             # HTTP DTOs
                └── middleware/ 
                    ├── cors.go                            # CORS Middleware
                    ├── logging.go                         # Logging Middleware
                    └── rate_limit.go                      # Rate Limit Middleware
```
