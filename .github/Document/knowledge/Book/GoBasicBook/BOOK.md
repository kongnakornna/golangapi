# คู่มือภาษา Golang แบบเรียบง่าย โดย  คงนคร จันทะคุณ  
---
 
```mermaid
graph TB
    subgraph CompileTime [คอมไพล์]
        A[source.go] --> B[Go Compiler<br/>gc + linker]
        B --> C[Static Binary<br/>executable]
    end

    subgraph Runtime [รันไทม์]
        C --> D[Go Runtime]
        D --> E[Goroutine Scheduler<br/>M:P:G Model]
        D --> F[Garbage Collector<br/>concurrent tri-color]
        D --> G[Memory Allocator<br/>mcache / mspan / mheap]
    end

    subgraph SchedulerDetail [Goroutine Scheduler]
        E --> H[G - Goroutine<br/>lightweight thread]
        E --> I[P - Processor<br/>logical CPU]
        E --> J[M - Machine<br/>OS thread]
        H -.-> I
        I -.-> J
        J --> K[CPU Cores]
    end

    subgraph Communication [การสื่อสาร]
        L[Goroutine A] -->|ch <- data| M[Channel<br/>chan T]
        M -->|data <- ch| N[Goroutine B]
        O[select] --> M
        P[sync.Mutex] --> L
        P --> N
    end

    subgraph GC [Garbage Collector]
        F --> Q[Mark phase<br/>concurrent]
        Q --> R[Sweep phase<br/>concurrent]
        R --> S[STW only short<br/> pause]
    end

    Runtime --> Communication
    Runtime --> GC
```

  
```
+-----------------------+        +------------------------+
|      source.go        |        |   Go Runtime (exec)    |
+----------+------------+        +-----------+------------+
           |                                 |
           v                                 v
+----------+------------+        +-----------+------------+
|   Go Compiler (gc)    |        |  Goroutine Scheduler   |
|   + static linking    |        |  (M:P:G Model)         |
+----------+------------+        +-----------+------------+
           |                                 |
           v                                 v
+----------+------------+        +-----------+------------+
|   Static Binary       |------->|   M (OS Threads)       |
|   (single executable) |        |   P (Logical Processors)|
+-----------------------+        |   G (Goroutines)       |
                                 +------------------------+
                                          |
                                          v
+---------------------------------------------------------+
|               Concurrency & Sync                         |
|  +------------+    +------------+    +------------+     |
|  | Goroutine  |--->|  Channel   |--->| Goroutine  |     |
|  +------------+    |  (buffered/|    +------------+     |
|                    |   unbuff)  |                       |
|                    +------------+                       |
|  +---------------------------------------------------+  |
|  | select, sync.Mutex, WaitGroup                     |  |
|  +---------------------------------------------------+  |
+---------------------------------------------------------+
                                          |
                                          v
+---------------------------------------------------------+
|            Garbage Collector (concurrent)               |
|  Tri-color marking + concurrent sweep, short STW        |
+---------------------------------------------------------+
```

---
 

## Domain-Driven Design with Go
----
```mermaid
graph TD
    subgraph Presentation
        A[HTTP Handlers<br/>gin/echo]
        B[gRPC Services]
        C[CLI Commands]
    end

    subgraph Application
        D[Application Services<br/>/application]
        E[DTOs / Commands / Queries]
        F[Use Case Orchestration]
    end

    subgraph Domain
        G[Entities<br/>struct with behavior]
        H[Value Objects<br/>immutable structs]
        I[Aggregates<br/>root entity + invariants]
        J[Domain Services<br/>stateless operations]
        K[Domain Events<br/>event recording/dispatch]
        L[Repository Interfaces<br/>domain layer contracts]
    end

    subgraph Infrastructure
        M[Repository Implementations<br/>SQL, MongoDB, etc.]
        N[Event Bus<br/>pub/sub, Kafka, etc.]
        O[External Services<br/>HTTP clients, SDKs]
        P[Persistence<br/>GORM, sqlx, mongodb driver]
    end

    A --> D
    B --> D
    C --> D
    D --> F
    F --> G
    F --> J
    G --> I
    H --> I
    I --> L
    I -->|fires| K
    K --> N
    L -.->|implements| M
    M --> P
    D --> E
    D -->|uses| O
```

---

### คำอธิบายโครงสร้าง DDD ใน Go

#### 1. **Presentation Layer** (ชั้นนำเสนอ)
- **บทบาท**: จัดการการรับส่งข้อมูลจากผู้ใช้ (HTTP, gRPC, CLI) และแปลงเป็นคำสั่งหรือแบบสอบถามที่ Application Layer เข้าใจ
- **ใน Go**: มักใช้ไลบรารี เช่น `net/http`, `gin`, `echo`, `grpc-go`  
- **หลักการ**: ไม่มีตรรกะทางธุรกิจ เน้นการแปลง request → DTO และเรียก Application Service

#### 2. **Application Layer** (ชั้นแอปพลิเคชัน)
- **บทบาท**: ประสานงานระหว่าง Presentation กับ Domain ควบคุม workflow ของ Use Case  
- **ใน Go**: แพ็กเกจ `application` หรือ `usecase` ประกอบด้วย **Application Services**  
  - รับ DTO/Command  
  - เรียก Domain Services / Aggregates  
  - จัดการ Unit of Work (Transaction)  
  - ส่ง Domain Events ผ่าน Event Bus  
- **ตัวอย่าง**:  
  ```go
  type OrderService struct {
      repo      domain.OrderRepository
      payment   domain.PaymentService
      eventBus  domain.EventBus
  }

  func (s *OrderService) PlaceOrder(ctx context.Context, cmd PlaceOrderCommand) error {
      // load aggregate
      order, _ := s.repo.FindByID(cmd.OrderID)
      // domain logic
      order.Place(cmd.Items)
      // persist
      s.repo.Save(order)
      // publish events
      s.eventBus.Publish(order.Events()...)
      return nil
  }
  ```

#### 3. **Domain Layer** (ชั้นโดเมน) — **หัวใจของ DDD**
- **Entities**: struct ที่มี identity และพฤติกรรม เช่น  
  ```go
  type Order struct {
      id      OrderID
      status  OrderStatus
      items   []OrderItem
      // behavior
      func (o *Order) Cancel() { ... }
      func (o *Order) AddItem(item OrderItem) error { ... }
  }
  ```
- **Value Objects**: struct ไม่มี identity เปรียบเทียบด้วยค่า เช่น `Address`, `Money`  
- **Aggregates**: กลุ่มของ entities/value objects ที่มี root entity (aggregate root) เป็นตัวควบคุมความสอดคล้อง  
- **Domain Services**: มีพฤติกรรมที่ไม่ผูกกับ aggregate หรือ entity ใดโดยเฉพาะ  
- **Domain Events**: ใช้ record การเปลี่ยนแปลงที่สำคัญ  
- **Repository Interfaces**: กำหนด interface สำหรับการเก็บและค้นหา aggregate root อยู่ใน domain layer

#### 4. **Infrastructure Layer** (ชั้นโครงสร้างพื้นฐาน)
- **Repository Implementations**: ใช้ GORM, sqlx, MongoDB driver เพื่อ implement interface ที่ domain กำหนด  
- **Event Bus**: ใช้ระบบเช่น Kafka, RabbitMQ หรือ in‑memory bus  
- **External Services**: HTTP clients, SDKs ต่างๆ  
- **หลักการ**: Infrastructure ต้องขึ้นกับ domain (Dependency Inversion) ไม่ใช่ในทางกลับกัน

---

### การนำ DDD ไปใช้กับ Go: ข้อแนะนำเฉพาะภาษา

1. **จัดโครงสร้างโปรเจกต์ตามโมดูล**  
   ```
   /cmd
     /myapp
       main.go
   /internal
     /domain
       /order
         order.go (entity, value objects)
         repository.go (interface)
         events.go
     /application
       order_service.go
     /infrastructure
       /persistence
         order_repo_mysql.go
       /bus
         event_bus_kafka.go
     /presentation
       /http
         order_handler.go
   /pkg
     /shared
       errors.go, uuid.go, etc.
   ```

2. **ใช้ interface เพื่อ Dependency Inversion**  
   - Application layer รับ domain interface  
   - Infrastructure ถูก inject ผ่าน constructor  
   - ใช้ `wire` (Google Wire) หรือ manual DI สำหรับการประกอบ dependencies

3. **จัดการ Transaction**  
   - นิยมใช้ `Unit of Work` pattern: application service เริ่ม transaction ผ่าน interface เช่น  
     ```go
     type UnitOfWork interface {
         Begin(ctx context.Context) error
         Commit() error
         Rollback() error
         OrderRepository() OrderRepository
     }
     ```
   - ใน Go มักใช้ `context` เพื่อส่ง transaction object (เช่น `*sql.Tx`) ไปยัง repository methods

4. **Domain Events**  
   - ใช้ channel หรือ event bus ภายใน memory ก่อน แล้วค่อยเพิ่ม infrastructure  
   - ตัวอย่าง:
     ```go
     type DomainEvent interface {
         OccurredAt() time.Time
     }

     type OrderPlaced struct { ... }

     type EventBus interface {
         Publish(event DomainEvent) error
         Subscribe(handler func(DomainEvent)) error
     }
     ```

5. **Value Objects กับ immutability**  
   - ใช้ struct พร้อม private fields และ constructor functions  
   - เปรียบเทียบด้วย `==` หรือ implement `Equals` method

---

### สรุปประโยชน์ของการใช้ DDD ร่วมกับ Go

- **ความชัดเจนของโดเมน**: โค้ดสะท้อนภาษาธุรกิจ (Ubiquitous Language)  
- **การแยกหน้าที่**: แต่ละ layer มีความรับผิดชอบชัดเจน ลด coupling  
- **ทดสอบง่าย**: domain layer ไม่ขึ้นกับ infrastructure สามารถ unit test ด้วย mock  
- **ปรับเปลี่ยน infrastructure ได้**: เปลี่ยนฐานข้อมูลหรือ event bus โดยไม่กระทบโดเมน  
- **Go เหมาะสม**: struct, interface, package system ช่วยให้จัดระเบียบตาม bounded context ได้ดี และการทำงาน concurrency ผ่าน goroutine ช่วยให้จัดการ domain events ได้มีประสิทธิภาพ

 

---

### 1. Aggregate (กลุ่มวัตถุที่มีความสอดคล้อง)

**Aggregate** คือกลุ่มของวัตถุ (Entities + Value Objects) ที่ถูกยึดเข้าด้วยกันโดย **Aggregate Root** (รูท) ซึ่งเป็นตัวเดียวที่อนุญาตให้เข้าถึงหรือแก้ไขวัตถุอื่นภายในกลุ่มจากภายนอก การออกแบบ Aggregate ช่วยรักษาความถูกต้องของข้อมูล (invariants) และลดความซับซ้อนในการจัดการธุรกรรม

**หลักการสำคัญ:**
- **Aggregate Root** มี identity และเป็นจุดเดียวที่ภายนอกเข้าถึงได้
- การเปลี่ยนแปลงใด ๆ ภายใน Aggregate ต้องทำผ่าน Root เท่านั้น
- ภายใน Aggregate เดียวกันต้องมีความสอดคล้องกันในทันที (consistency boundary)
- ระหว่าง Aggregates ควรใช้ **eventual consistency** ผ่าน Domain Events

**ตัวอย่างใน Go:**

```go
// aggregate root: Order
type Order struct {
    id      OrderID
    status  OrderStatus
    items   []OrderItem   // value object
    total   Money
}

func (o *Order) AddItem(product Product, quantity int) error {
    if o.status != Draft {
        return errors.New("cannot add item to non-draft order")
    }
    // invariant: total must not exceed limit
    newTotal := o.total.Add(product.Price.Mul(quantity))
    if newTotal.GreaterThan(MaxOrderTotal) {
        return errors.New("order total exceeds limit")
    }
    o.items = append(o.items, NewOrderItem(product, quantity))
    o.total = newTotal
    o.addDomainEvent(OrderItemAdded{OrderID: o.id, ProductID: product.ID})
    return nil
}

// repository interface รับเฉพาะ aggregate root
type OrderRepository interface {
    Save(order *Order) error
    FindByID(id OrderID) (*Order, error)
}
```

---

### 2. Event Storming (เทคนิคค้นพบโดเมน)

**Event Storming** เป็นเวิร์กช็อปแบบมีส่วนร่วมที่ช่วยให้ทีม (นักพัฒนา, นักธุรกิจ, ผู้เชี่ยวชาญโดเมน) ระบุและเข้าใจโดเมนผ่านเหตุการณ์สำคัญที่เกิดขึ้นในระบบ โดยใช้โน้ตสีต่าง ๆ บนกระดาน

**สัญลักษณ์ทั่วไป:**
- **สีส้ม** – Domain Events (สิ่งที่เกิดขึ้นแล้ว) เช่น `OrderPlaced`, `PaymentReceived`
- **สีน้ำเงิน** – Commands (คำสั่งที่ทำให้เกิดเหตุการณ์) เช่น `PlaceOrder`, `RefundPayment`
- **สีเหลือง** – Aggregates (กลุ่มข้อมูลที่ถูกคำสั่งเรียก) เช่น `Order`, `Customer`
- **สีม่วง** – External Systems / Policies (ระบบภายนอกหรือกฎ)
- **สีเขียว** – Read Models / Views (ข้อมูลที่แสดงผล)

**ขั้นตอนคร่าว ๆ:**
1. ระบุ Domain Events (อดีต) โดยเรียงตามลำดับเวลา
2. ระบุ Commands ที่ทำให้เกิดเหตุการณ์นั้น
3. จับคู่ Command กับ Aggregate (ผู้รับผิดชอบ)
4. เพิ่ม Policies / Rules และ External Systems
5. ระบุ Read Models ที่จำเป็นต่อการแสดงผล

Event Storming นำไปสู่การกำหนด **Bounded Context** และ **Aggregates** ที่ชัดเจน ก่อนเริ่มเขียนโค้ด

---

### 3. CQRS (Command Query Responsibility Segregation)

CQRS แยกโมเดลการ **เขียน** (Command) และ **อ่าน** (Query) ออกจากกัน ทำให้สามารถปรับแต่งให้เหมาะสมกับแต่ละฝั่งได้ เช่น ใช้ฐานข้อมูลแบบ Event Sourcing สำหรับ Command และฐานข้อมูลแบบ Materialized View สำหรับ Query

**โครงสร้างทั่วไป:**

```mermaid
graph LR
    subgraph Client
        A[User Interface]
    end

    subgraph Command Side
        B[Command Handler]
        C[Aggregate]
        D[Event Store]
        E[Event Bus]
    end

    subgraph Query Side
        F[Query Handler]
        G[Read Database]
        H[Projection]
    end

    A -->|Command| B
    B --> C
    C -->|produces events| D
    D --> E
    E -->|project| H
    H --> G
    A -->|Query| F
    F --> G
    G -->|result| A
```

**ใน Go สามารถจัดโครงสร้างได้ดังนี้:**

#### แยก Command และ Query Models
```go
// Command models (เขียน)
type PlaceOrderCommand struct {
    OrderID string
    Items   []OrderItemDTO
}

// Query models (อ่าน)
type OrderView struct {
    OrderID    string
    Status     string
    TotalPrice float64
    Items      []OrderItemView
}
```

#### Command Handlers
```go
type PlaceOrderHandler struct {
    repo      domain.OrderRepository
    eventBus  domain.EventBus
}

func (h *PlaceOrderHandler) Handle(ctx context.Context, cmd PlaceOrderCommand) error {
    order, err := domain.NewOrder(cmd.OrderID)
    if err != nil {
        return err
    }
    for _, item := range cmd.Items {
        order.AddItem(item.ProductID, item.Quantity)
    }
    if err := h.repo.Save(order); err != nil {
        return err
    }
    // Publish events for projection
    for _, event := range order.Events() {
        h.eventBus.Publish(event)
    }
    return nil
}
```

#### Query Handlers (อ่านจาก Read Database)
```go
type OrderQueryHandler struct {
    db *sql.DB // หรือ ORM
}

func (h *OrderQueryHandler) GetOrder(ctx context.Context, orderID string) (*OrderView, error) {
    var view OrderView
    err := h.db.QueryRowContext(ctx, "SELECT ... FROM order_views WHERE id = ?", orderID).Scan(&view)
    return &view, err
}
```

#### Projection (สร้าง Read Model จาก Events)
```go
type OrderProjection struct {
    db *sql.DB
}

func (p *OrderProjection) HandleEvent(event domain.DomainEvent) {
    switch e := event.(type) {
    case OrderPlaced:
        p.db.Exec("INSERT INTO order_views (id, status, total) VALUES (?, ?, ?)", e.OrderID, "Placed", e.Total)
    case OrderItemAdded:
        p.db.Exec("INSERT INTO order_items_view ...")
    }
}
```

---

### 4. การนำ CQRS ไปใช้ใน Go อย่างมีประสิทธิภาพ

- **ใช้ Interface ในการแยก**: Command handlers, Query handlers, Projection handlers แต่ละตัวเป็น struct ที่ implements interface ต่างกัน ทำให้ test และ replace ได้ง่าย
- **Event Store**: ใน Go สามารถใช้ database เช่น PostgreSQL พร้อม JSONB เก็บ events, หรือใช้ไลบรารี `github.com/eventials/goeventsourcing` แต่แนะนำให้เริ่มจาก in‑memory ก่อนค่อยเพิ่ม complexity
- **Read Database**: อาจใช้ฐานข้อมูลแยก (SQL, NoSQL) หรือ cache เช่น Redis สำหรับการ query ที่รวดเร็ว
- **Concurrency**: Goroutine + channel ใช้ในการประมวลผล projection แบบ asynchronous

**ข้อควรระวัง:**
- CQRS เพิ่มความซับซ้อน เหมาะสำหรับระบบที่ต้องการความยืดหยุ่นสูง (microservices, high scalability)
- ไม่จำเป็นต้องใช้ Event Sourcing เสมอไป CQRS สามารถแยกโมเดลอ่าน-เขียนโดยใช้ฐานข้อมูลเดียวกันได้ (แต่แยกตารางหรือ schema)
- การจัดการ eventual consistency ต้องออกแบบ UX ให้เหมาะสม

---

### สรุปการเพิ่มองค์ประกอบ

| Concept | บทบาทใน DDD with Go |
|---------|---------------------|
| **Aggregate** | กลุ่มข้อมูลและพฤติกรรมที่มีความสอดคล้อง ใช้ interface repository ในการโหลด/บันทึกผ่าน aggregate root |
| **Event Storming** | กระบวนการออกแบบร่วมกันเพื่อค้นหา domain events, commands, aggregates, read models ก่อนพัฒนา |
| **CQRS** | แยก command (เขียน) และ query (อ่าน) ใช้ event bus และ projection เพื่อสร้าง read model |

ทั้งสามแนวคิดช่วยให้การพัฒนา DDD ใน Go มีระเบียบแบบแผน รองรับความซับซ้อนของโดเมน และขยายขนาดได้ง่าย หากเริ่มต้นด้วย bounded context เล็ก ๆ ก่อนแล้วค่อยเพิ่ม CQRS เมื่อจำเป็นก็จะลดความเสี่ยงได้ดี

---
 
 
# คู่มือการเรียนรู้ภาษา Go อย่างครอบคลุม ตั้งแต่พื้นฐานจนถึงระดับสูง
---
## 📚 สารบัญ

### ภาคที่ 1: ปฐมบทกับการเขียนโปรแกรม
- บทที่ 1: ความรู้เบื้องต้นเกี่ยวกับการเขียนโปรแกรมคอมพิวเตอร์
- บทที่ 2: รู้จักกับภาษา Go
- บทที่ 3: พื้นฐานการใช้งาน Terminal
- บทที่ 4: เตรียมสภาพแวดล้อมสำหรับพัฒนา
- บทที่ 5: สร้างแอปพลิเคชันแรกของคุณ

### ภาคที่ 2: พื้นฐานภาษาและโครงสร้างข้อมูล
- บทที่ 6: ระบบเลขฐานสองและฐานสิบ
- บทที่ 7: เลขฐานสิบหก, ฐานแปด, ASCII, UTF8, Unicode และ Runes
- บทที่ 8: ตัวแปร, ค่าคงที่ และชนิดข้อมูลพื้นฐาน
- บทที่ 9: คำสั่งควบคุมการทำงาน
- บทที่ 10: ฟังก์ชัน
- บทที่ 11: แพคเกจและการนำเข้า
- บทที่ 12: การเริ่มต้นทำงานของแพคเกจ
- บทที่ 13: การสร้างชนิดข้อมูลใหม่ (Types)
- บทที่ 14: เมธอด (Methods)
- บทที่ 15: พอยน์เตอร์ (Pointer)
- บทที่ 16: อินเทอร์เฟซ (Interfaces)

### ภาคที่ 3: การจัดการโปรเจกต์และโครงสร้างข้อมูลขั้นสูง
- บทที่ 17: Go Modules - การจัดการโปรเจกต์สมัยใหม่
- บทที่ 18: Go Module Proxies
- บทที่ 19: การทดสอบหน่วย (Unit Tests)
- บทที่ 20: อาเรย์ (Arrays)
- บทที่ 21: สไลซ์ (Slices)
- บทที่ 22: แมพ (Maps)
- บทที่ 23: การจัดการข้อผิดพลาด (Errors)

### ภาคที่ 4: การพัฒนาแอปพลิเคชันเชิงปฏิบัติ
- บทที่ 24: ฟังก์ชันนิรนาม (Anonymous functions) และ Closure
- บทที่ 25: การจัดการข้อมูล JSON และ XML
- บทที่ 26: พื้นฐานการสร้าง HTTP Server
- บทที่ 27: Enum, Iota และ Bitmask
- บทที่ 28: วันที่และเวลา
- บทที่ 29: การจัดเก็บข้อมูล: ไฟล์และฐานข้อมูล
- บทที่ 30: การทำงานพร้อมกัน (Concurrency)
- บทที่ 31: การบันทึกเหตุการณ์ (Logging)
- บทที่ 32: เทมเพลต (Templates)
- บทที่ 33: การจัดการค่า Configuration

### ภาคที่ 5: สู่การเป็นนักพัฒนา Go มืออาชีพ
- บทที่ 34: การวัดประสิทธิภาพ (Benchmarks)
- บทที่ 35: สร้าง HTTP Client
- บทที่ 36: การวิเคราะห์โปรไฟล์ (Program Profiling)
- บทที่ 37: การจัดการ Context
- บทที่ 38: Generics - การเขียนโค้ดแบบยืดหยุ่น
- บทที่ 39: Go กับกระบวนทัศน์ OOP?
- บทที่ 40: การอัปเกรดหรือดาวน์เกรดเวอร์ชัน Go
- บทที่ 41: คำแนะนำในการออกแบบโค้ดที่ดี
- บทที่ 42: ชีทสรุป (Cheatsheet)

---

## ภาคที่ 1: ปฐมบทกับการเขียนโปรแกรม

### บทที่ 1: ความรู้เบื้องต้นเกี่ยวกับการเขียนโปรแกรมคอมพิวเตอร์

#### 1.1 การเขียนโปรแกรมคืออะไร?
การเขียนโปรแกรม (Programming) คือกระบวนการสร้างชุดคำสั่งที่ใช้ควบคุมการทำงานของคอมพิวเตอร์ให้ทำงานตามที่เราต้องการ โดยใช้ภาษาเฉพาะที่คอมพิวเตอร์สามารถเข้าใจได้ ภาษาที่มนุษย์ใช้เขียนเรียกว่า “ภาษาคอมพิวเตอร์ระดับสูง” (High-Level Language) เช่น Go, Python, Java ซึ่งจากนั้นจะถูกแปลงเป็นภาษาเครื่อง (Machine Language) ที่เป็นเลขฐานสอง (0 และ 1) ที่ซีพียูสามารถประมวลผลได้

#### 1.2 โครงสร้างพื้นฐานของโปรแกรม
โปรแกรมคอมพิวเตอร์โดยทั่วไปประกอบด้วย:
- **ข้อมูล (Data)** : ตัวเลข, ข้อความ, รายการต่างๆ
- **การประมวลผล (Processing)** : การดำเนินการกับข้อมูล เช่น การคำนวณ การเปรียบเทียบ
- **การควบคุมการทำงาน (Control Flow)** : การตัดสินใจ (if-else), การวนซ้ำ (loop)
- **การจัดเก็บ (Storage)** : หน่วยความจำ, ไฟล์, ฐานข้อมูล
- **อินพุต/เอาท์พุต (I/O)** : การรับข้อมูลจากผู้ใช้ หรือแสดงผล

#### 1.3 ตัวแปลภาษาและคอมไพเลอร์
- **คอมไพเลอร์ (Compiler)** : แปลงซอร์สโค้ดทั้งหมดเป็นไฟล์ได้ก่อนรัน (เช่น Go, C)
- **อินเทอร์พรีเตอร์ (Interpreter)** : แปลงและรันทีละคำสั่ง (เช่น Python, JavaScript)

Go เป็นภาษาแบบคอมไพล์ (compiled) ซึ่งมีข้อดีคือทำงานเร็วและสร้างไฟล์ binary ที่รันได้ทันทีโดยไม่ต้องพึ่งพาสิ่งแวดล้อมอื่น (ยกเว้นระบบปฏิบัติการ)

#### 1.4 กระบวนทัศน์การเขียนโปรแกรม
Go รองรับการเขียนโปรแกรมแบบ:
- **Procedural** : ใช้ฟังก์ชันและลำดับขั้นตอน
- **Concurrent** : ทำงานพร้อมกันด้วย goroutine
- **Functional** (บางส่วน) : ฟังก์ชันเป็น first-class citizen
- **ไม่ใช่ OOP แบบคลาสสิก** : ใช้ struct และ interface แทน inheritance
-----------
### 1. Golang Procedural คืออะไร
**Procedural** หรือการเขียนโปรแกรมแบบเชิงกระบวนการ เป็นรูปแบบที่เน้นการเขียนฟังก์ชัน (function) และเรียกใช้ตามลำดับขั้นตอน โกแลง (Go) รองรับการเขียนแบบนี้โดยใช้ฟังก์ชันเป็นหน่วยหลัก โค้ดจะถูกแบ่งออกเป็นฟังก์ชันย่อย ๆ ที่ทำงานเฉพาะอย่าง และเรียกใช้ตามลำดับที่กำหนด ถึงแม้ Go จะไม่ใช่ภาษา procedural ล้วน ๆ แต่ก็สามารถเขียนในสไตล์นี้ได้อย่างดี

---

### 2. Concurrent คืออะไร
**Concurrency** (การทำงานร่วมกัน) คือความสามารถของโปรแกรมในการจัดการงานหลาย ๆ งานพร้อมกัน โดยงานเหล่านี้อาจเริ่มทำงาน สลับการทำงาน หรือจบลงในเวลาที่เหลื่อมกัน โดยไม่จำเป็นต้องทำงานพร้อมกันจริง ๆ (parallel) Go มีเครื่องมือสำหรับ concurrency โดยตรง ได้แก่ **goroutine** และ **channel** ทำให้การเขียนโปรแกรม concurrent ง่ายและมีประสิทธิภาพ

---

### 3. goroutine คืออะไร
**goroutine** คือเธรด (thread) ขนาดเบาที่ถูกจัดการโดย runtime ของ Go ใช้สำหรับทำงานแบบ concurrent การสร้าง goroutine ทำได้ง่ายโดยใช้คีย์เวิร์ด `go` หน้าฟังก์ชัน:

```go
go myFunction()
```

Goroutine มีน้ำหนักเบา ใช้สแต็กเริ่มต้นเพียงไม่กี่กิโลไบต์ และสามารถทำงานหลายหมื่นตัวพร้อมกันได้ในโปรแกรมเดียว

---

### 4. Functional คืออะไร
**Functional programming** (การเขียนโปรแกรมเชิงฟังก์ชัน) เป็นกระบวนทัศน์ที่เน้นการใช้ฟังก์ชันบริสุทธิ์ (pure function) ไม่มีการเปลี่ยนแปลงสถานะภายนอก และใช้ฟังก์ชันเป็นพลเมืองชั้นหนึ่ง (first-class citizen) Go รองรับคุณสมบัติบางอย่างของ functional เช่น การส่งฟังก์ชันเป็นค่าพารามิเตอร์ การคืนค่าฟังก์ชัน แต่ไม่มีการรับประกันความบริสุทธิ์ของฟังก์ชันหรือโครงสร้างข้อมูลแบบ immutable โดยตรง

---

### 5. first-class citizen คืออะไร
**First-class citizen** (พลเมืองชั้นหนึ่ง) หมายถึงสิ่งที่สามารถทำได้เช่นเดียวกับค่าอื่น ๆ ในภาษา เช่น ถูกกำหนดให้กับตัวแปร ถูกส่งเป็นอาร์กิวเมนต์ให้ฟังก์ชัน ถูกคืนค่าจากฟังก์ชัน หรือถูกเก็บในโครงสร้างข้อมูล ใน Go **ฟังก์ชัน** ถือเป็น first-class citizen:

```go
var fn func(int) int = func(x int) int { return x * x }
```

นอกจากนี้ ชนิดข้อมูลต่าง ๆ เช่น struct, interface ก็ถือเป็น first-class citizen เช่นกัน

---

### 6. OOP คืออะไร
**OOP** (Object-Oriented Programming) เป็นกระบวนทัศน์ที่เน้นการสร้างวัตถุ (object) ซึ่งรวมข้อมูลและพฤติกรรมเข้าด้วยกัน หลักการสำคัญคือ encapsulation, inheritance, และ polymorphism Go ไม่ใช่ภาษา OOP แบบดั้งเดิมเพราะไม่มี class แต่ใช้ **struct** เป็นตัวเก็บข้อมูล และใช้ **method** ที่ผูกกับ struct รวมถึง **interface** เพื่อสร้าง polymorphism การสืบทอด (inheritance) ทำได้โดยการ **composition** (การแทรก struct) แทน

---

### 7. struct คืออะไร
**struct** เป็นชนิดข้อมูลที่ใช้รวมฟิลด์ (field) หลาย ๆ ชนิดเข้าด้วยกัน คล้ายกับ class ในภาษาอื่น แต่ไม่มี method ในตัว เราสามารถกำหนด method ให้กับ struct ได้โดยประกาศฟังก์ชันที่มี **receiver** เป็น struct นั้น

```go
type Person struct {
    Name string
    Age  int
}

func (p Person) Greet() string {
    return "Hello, " + p.Name
}
```

---

### 8. interface คืออะไร
**interface** คือชนิดข้อมูลที่กำหนดชุดของ method signature (ชื่อฟังก์ชัน, พารามิเตอร์, ผลลัพธ์) โดยไม่ระบุการนำไปใช้ ถ้าชนิดใด (เช่น struct) มี method ครบตามที่ interface กำหนด ชนิดนั้นจะถูกพิจารณาว่า **implement** interface นั้นโดยอัตโนมัติ (implicit implementation) ทำให้เกิด polymorphism แบบหนึ่งใน Go

```go
type Greeter interface {
    Greet() string
}

// Person มี method Greet() string ดังนั้น Person จึงเป็น Greeter โดยอัตโนมัติ
```

---

### 9. inheritance คืออะไร
**Inheritance** (การสืบทอด) คือกลไกที่ class หนึ่งสามารถสืบทอดสมาชิก (ฟิลด์และ method) จากอีก class หนึ่ง เพื่อนำมาใช้หรือขยายความสามารถ Go **ไม่มี** inheritance แบบ class-based แต่นิยมใช้ **composition** โดยการฝัง struct (embedded struct) เพื่อให้ได้ผลลัพธ์คล้ายการสืบทอด แต่ยังคงความยืดหยุ่นและหลีกเลี่ยงปัญหาที่เกิดจาก inheritance ที่ซับซ้อน

```go
type Animal struct {
    Name string
}

type Dog struct {
    Animal  // ฝัง struct Animal เข้ามา
    Breed string
}
```

Dog จะสามารถเข้าถึงฟิลด์ Name และ method (ถ้ามี) ของ Animal ได้โดยตรง
--------

#### 1.5 ขั้นตอนการพัฒนาโปรแกรม
1. เขียนซอร์สโค้ด (.go)
2. คอมไพล์ (go build)
3. ทดสอบรัน (go run)
4. แก้ไขข้อผิดพลาด (debug)
5. จัดการแพคเกจ (go mod)
6. ทดสอบหน่วย (go test)

---

### บทที่ 2: รู้จักกับภาษา Go

#### 2.1 ประวัติความเป็นมา
ภาษา Go (หรือ Golang) ถูกพัฒนาโดย Google เริ่มต้นในปี 2007 โดย Robert Griesemer, Rob Pike, และ Ken Thompson เปิดตัวเป็นโอเพนซอร์สในปี 2009 จุดประสงค์เพื่อแก้ปัญหาที่เกิดขึ้นในภาษา C++ และ Java ในระบบขนาดใหญ่ของ Google เช่น การคอมไพล์ที่ช้า, การจัดการการทำงานพร้อมกันที่ซับซ้อน, และความยุ่งยากในการบำรุงรักษา

#### 2.2 จุดเด่นของภาษา Go
- **เรียบง่ายและอ่านง่าย** : ไวยากรณ์กระชับ ไม่มีฟีเจอร์ที่ซับซ้อนเกินจำเป็น
- **คอมไพล์เร็ว** : สามารถคอมไพล์โปรเจกต์ขนาดใหญ่ได้ในไม่กี่วินาที
- **การจัดการหน่วยความจำอัตโนมัติ** : มี garbage collector ที่มีประสิทธิภาพ
- **Concurrency ระดับภาษา** : goroutine และ channel ทำให้เขียนโปรแกรม concurrent ได้ง่าย
- **Static typing** : ตรวจสอบชนิดข้อมูลตั้งแต่ตอนคอมไพล์ ป้องกันข้อผิดพลาด
- **เครื่องมือที่ครบครัน** : go fmt, go test, go mod, go vet, go doc
- **สามารถคอมไพล์ข้ามแพลตฟอร์ม** (cross-compile) ไปยัง Windows, Linux, macOS, ARM เป็นต้น

#### 2.3 โครงสร้างภาษา Go
- ไม่มี class แต่ใช้ struct และ method
- ไม่มี inheritance แต่ใช้ composition และ interface
- ไม่มี exception handling แต่ใช้ error return value
- มี garbage collection
- มี pointer แต่ไม่มี pointer arithmetic

-----------
## โครงสร้างภาษา Go

ภาษา Go (หรือ Golang) มีโครงสร้างพื้นฐานที่เรียบง่ายและชัดเจน ดังนี้

### 1. **Package**
- ทุกไฟล์ `.go` ต้องขึ้นต้นด้วย `package` เพื่อระบุแพ็กเกจที่สังกัด
- โปรแกรมหลักต้องใช้ `package main` และมีฟังก์ชัน `main()` เป็นจุดเริ่มต้น

### 2. **Import**
- ใช้ `import` เพื่อนำแพ็กเกจภายนอกหรือมาตรฐานเข้ามาใช้งาน
- รองรับการ import แบบกลุ่มหรือแบบ single

### 3. **Function**
- ประกาศด้วย `func` ตามด้วยชื่อ พารามิเตอร์ และชนิดค่าที่คืน
- สามารถคืนค่าหลายค่าได้ (multiple return values)

### 4. **Struct & Interface**
- `struct` ใช้รวบรวมข้อมูล (คล้ายคลาสแบบ data-only)
- `interface` กำหนดชุดของ method signature ใช้สร้าง polymorphism
--------------
`interface` ในภาษา Go คือการกำหนด **ชุดของ method signature** (ชื่อ method, พารามิเตอร์, และค่าที่คืน) โดยไม่ต้องระบุการทำงานจริง ชนิดใดก็ตาม (เช่น `struct`) ที่มี method ครบตามที่ interface กำหนดไว้ จะถูกพิจารณาว่า **implement** interface นั้นโดยอัตโนมัติ (implicit implementation)

**Polymorphism** (ภาวะหลายรูปแบบ) คือความสามารถในการเขียนโค้ดที่ทำงานกับค่าได้หลายชนิดผ่านทาง interface เดียวกัน เนื่องจากตัวแปรชนิด interface สามารถเก็บค่าได้ทุกชนิดที่ implement interface นั้น ส่งผลให้เราไม่ต้องรู้ชนิดที่แท้จริงของข้อมูลในขณะเขียนโค้ด เพียงแค่เรียก method ผ่าน interface ก็จะได้พฤติกรรมที่เหมาะสมตามชนิดของข้อมูลที่ถูกเก็บไว้

### ตัวอย่าง

```go
package main

import "fmt"

// กำหนด interface ที่มี method signature
type Speaker interface {
    Speak() string
}

// Struct Dog โคงสร้างหลัก
type Dog struct {
    Name string
}

// Dog implement Speaker  นำโคงสร้างหลัก Dog มา ใช้งาน ใน function Dog
func (d Dog) Speak() string {
    return "woof! I'm " + d.Name
}

// Struct Cat 
type Cat struct {
    Name string
}

// Cat implement Speaker นำโคงสร้างหลัก Cat มา ใช้งาน ใน function Cat
func (c Cat) Speak() string {
    return "meow! I'm " + c.Name
}

// ฟังก์ชันที่รับ Speaker ทำให้ทำงานได้กับทั้ง Dog และ Cat
func Introduce(s Speaker) {
    fmt.Println(s.Speak())
}
// เรียกใช้งาน ฟังก์ชัน มาทำงาน ใน main  ฟังก์ชันที่รับ Introduce จาก โครงสร้าง ฟังก์ชัน Dog  ฟังก์ชัน Cat 
func main() {
    dog := Dog{Name: "Max"}
    cat := Cat{Name: "Luna"}

    Introduce(dog) // woof! I'm Max  => ผลลัพธ์
    Introduce(cat) // meow! I'm Luna => ผลลัพธ์ 
}
```

### สรุปประโยชน์
- **Polymorphism** ทำให้สามารถเขียนฟังก์ชันหรือโครงสร้างที่ใช้งานได้กับหลายชนิด โดยไม่ต้องแก้ไขโค้ดเดิมเมื่อเพิ่มชนิดใหม่ (แค่ให้ชนิดใหม่ implement interface เดิม)
- **Decoupling** โค้ดที่เรียกใช้ interface จะไม่ขึ้นอยู่กับโครงสร้างของชนิดข้อมูลโดยตรง ทำให้ทดสอบ (unit test) ได้ง่ายขึ้นโดยใช้ mock objects
- **ยืดหยุ่น** เหมาะกับการออกแบบระบบที่ต้องการขยายความสามารถโดยไม่กระทบส่วนที่ใช้ interface
--------------
### 5. **Method**
- ฟังก์ชันที่มี **receiver** (เช่น `func (p Person) Greet() string`) ทำให้ struct มีพฤติกรรม =>  receiver เสมือน "ตัวรับ" ของ method นั้น ทำงานคล้ายกับ this หรือ self ในภาษาอื่น เช่น C,java ,javascript
---------------
**Receiver** ในภาษา Go คือตัวแปรพิเศษที่ผูก method เข้ากับชนิดข้อมูล (type) โดยเฉพาะ โดย method จะถูกประกาศด้วย syntax `func (r T) MethodName(...) ...` โดย `r` คือ receiver ซึ่งเปรเสมือน "ตัวรับ" ของ method นั้น ทำงานคล้ายกับ `this` หรือ `self` ในภาษาอื่น แต่มีความยืดหยุ่นมากกว่า

## ไวยากรณ์พื้นฐาน
```go
type MyType struct {
    Value int
}

// value receiver "รับ" ค่าผ่าน method  มาทำงาน โดยระบบ ประเภทข้อมูล แบบ MyType
func (m MyType) GetValue() int {
    return m.Value
}

// pointer receiver
func (m *MyType) SetValue(v int) {
    m.Value = v
}
```
------------------------
**Pointer Receiver** คือการประกาศ method โดยให้ receiver เป็น **pointer** ระบุตำแหน่ง ไปยังชนิดข้อมูล (type) นั้น ๆ เช่น `func (t *T) MethodName()` แทนที่จะเป็นตัวแปรธรรมดา (`t T`)

## ทำงานอย่างไร
- เมื่อ method มี pointer receiver ตัวแปร `t` ใน method จะเป็น pointer ที่ชี้ไปยัง instance โคร้งสร้างดั้งเดิม (original instance)
- การเปลี่ยนแปลงค่าผ่าน `t.field` จะส่งผลต่อ instance ต้นฉบับโดยตรง เพราะเป็นการแก้ไขผ่าน address
- ไม่มีการคัดลอก struct ทั้งก้อน ทำให้มีประสิทธิภาพเมื่อ struct มีขนาดใหญ่

## ตัวอย่าง
```go
package main

import "fmt"
// ชือโครงสร้างหลัก Person
type Person struct {
    Name string  //ตัวแปร  Name
    Age  int  //ตัวแปร  Age
}

// Pointer receiver – สามารถแก้ไขค่าต้นฉบับได้
// รับค่า (p *Person) จาก โครงสร้างหลัก Person มาใช้ใน ชือฟังก์ชั่น  HaveBirthday
func (p *Person) HaveBirthday() {
    p.Age++ // แก้ไข field ของ instance ดั้งเดิม  เพิ่มจำนวน เพิ่มใหม่เข้าไป  1 เช่น 30+1=31
}

// Value receiver – ทำงานกับสำเนา
func (p Person) Greet() string {
    return "Hello, " + p.Name
}

func main() {
    p := Person{Name: "Alice", Age: 30}
    p.HaveBirthday()
    fmt.Println(p.Age) // 31 (เปลี่ยนแปลงจริง)
    
    fmt.Println(p.Greet()) // Hello, Alice
}
```

## เปรียบเทียบ Pointer Receiver กับ Value Receiver

| คุณสมบัติ | Pointer Receiver (`*T`) | Value Receiver (`T`) |
|----------|------------------------|---------------------|
| การแก้ไขค่า | ✅ แก้ไขค่าเดิมได้ | ❌ ทำงานกับสำเนา |
| การคัดลอก | ไม่คัดลอก struct (ส่งแค่ address) | คัดลอก struct ทั้งก้อน |
| ประสิทธิภาพ | ดีสำหรับ struct ใหญ่ | อาจช้าถ้า struct ใหญ่ |
| ความปลอดภัย | เปลี่ยนแปลงได้ | immutable (ไม่เปลี่ยนต้นฉบับ) |
| การ implement interface | ถ้า method เป็น pointer receiver จะ implement โดย pointer type เท่านั้น | ทั้ง value และ pointer implement ได้ (Go จะ dereference อัตโนมัติ) |

## เมื่อใดควรใช้ Pointer Receiver
1. **ต้องการเปลี่ยนแปลง state** ของ struct (setter, update)
2. **struct มีขนาดใหญ่** – เพื่อลดการคัดลอก
3. **ต้องการความสม่ำเสมอ** – ถ้า method ใด method หนึ่งใน type ใช้ pointer receiver ควรใช้ pointer receiver สำหรับทุก method ของ type นั้น (เพื่อให้ interface เข้าใจง่าย)
4. **nil receiver handling** – pointer receiver สามารถตรวจสอบและจัดการกับ receiver ที่เป็น nil ได้ (value receiver ไม่สามารถ)

## ข้อควรระวัง
- หากคุณผูก method กับ **value receiver** (`func (t T)`) แล้วเรียกใช้ผ่าน pointer (`&t`), Go จะแปลงเป็น value ให้อัตโนมัติ (เพราะเข้าถึง field ได้)
- แต่ถ้าผูก method กับ **pointer receiver** (`func (t *T)`) แล้วเรียกใช้ผ่าน value (`t`), Go จะแปลงเป็น pointer ให้อัตโนมัติด้วย (สะดวก)
- อย่างไรก็ตาม **interface** มีข้อจำกัด: ถ้า interface มี method ที่เป็น pointer receiver เฉพาะ pointer type เท่านั้นที่ implement interface นั้น

```go
type Walker interface {
    Walk() string
}

type Dog struct{}

// pointer receiver
func (d *Dog) Walk() string { return "walking" }

var w Walker
// w = Dog{}   // ผิด: Dog does not implement Walker (Walk method has pointer receiver)
w = &Dog{}     // ถูกต้อง
```

## สรุป
- **Pointer receiver** ช่วยให้ method สามารถแก้ไข instance ต้นฉบับได้ และลดการคัดลอกข้อมูล
- ใช้ pointer receiver เมื่อต้องการเปลี่ยนแปลง struct หรือเมื่อ struct มีขนาดใหญ่
- ใช้ value receiver เมื่อ method เป็น pure function (ไม่เปลี่ยนแปลง state) หรือ struct มีขนาดเล็ก
- การเลือกชนิด receiver มีผลต่อ interface implementation และควรใช้ให้สอดคล้องกันทั้ง type
--------------------------
## ประเภทของ Receiver

### 1. **Value Receiver** (`func (t T) Method()`)
- รับ **สำเนา** ของค่าเดิม
- การเปลี่ยนแปลงภายใน method **ไม่มีผล** กับค่าต้นฉบับ
- ปลอดภัยในการใช้งานกับข้อมูลที่ไม่ต้องการเปลี่ยนแปลง
- ทำงานกับ interface ได้ดี (ถ้า interface ต้องการ method นี้)

### 2. **Pointer Receiver** (`func (t *T) Method()`)
- รับ **พอยน์เตอร์** ไปยังค่าเดิม
- สามารถ **แก้ไข** ค่าต้นฉบับได้
- เมื่อ method ถูกเรียกด้วย pointer receiver จะไม่เกิดการคัดลอกโครงสร้างใหญ่ ทำให้มีประสิทธิภาพสำหรับ struct ขนาดใหญ่

## ตัวอย่างเปรียบเทียบ
```go
package main

import "fmt"

type Counter struct {
    count int
}

// value receiver - ไม่เปลี่ยนค่าต้นฉบับ
func (c Counter) IncrementValue() {
    c.count++
    fmt.Println("Inside value receiver:", c.count)
}

// pointer receiver - เปลี่ยนค่าต้นฉบับ
func (c *Counter) IncrementPointer() {
    c.count++
    fmt.Println("Inside pointer receiver:", c.count)
}

func main() {
    c := Counter{count: 0}
    
    c.IncrementValue()   // Inside value receiver: 1
    fmt.Println(c.count) // 0 (ไม่เปลี่ยนแปลง)
    
    c.IncrementPointer() // Inside pointer receiver: 1
    fmt.Println(c.count) // 1 (เปลี่ยนแปลง)
}
```

## การเลือกใช้ Value Receiver หรือ Pointer Receiver

| พิจารณาจาก | ควรใช้ Pointer Receiver | ควรใช้ Value Receiver |
|------------|------------------------|----------------------|
| **ต้องการแก้ไขค่า** | ✅ จำเป็น | ❌ ไม่สามารถ |
| **ขนาดของ struct** | ใหญ่ (ลดการคัดลอก) | เล็ก (copy cost ต่ำ) |
| **ความสอดคล้อง** | ต้องการ method ทั้งหมดเป็น pointer เพื่อให้ interface เข้าใจง่าย | method ทั้งหมดเป็น value ก็ได้ |
| **nil safety** | สามารถจัดการกับ receiver ที่เป็น nil ได้ | ไม่สามารถ (panic ถ้าเป็น nil) |
| **การ implement interface** | ถ้า interface มี method หนึ่งเป็น pointer receiver ต้องใช้ pointer กับตัวแปรนั้น | interface สามารถรับทั้ง value และ pointer ได้ถ้า method ทั้งหมดเป็น value receiver |

## ความสัมพันธ์กับ Interface

- หาก interface กำหนด method ไว้เป็น **value receiver** (`func (T) Method()`), ค่า type `T` และ `*T` **ทั้งคู่** จะ implement interface นั้น
- หาก interface กำหนด method ไว้เป็น **pointer receiver** (`func (*T) Method()`), **เฉพาะ** pointer (`*T`) เท่านั้นที่ implement interface

```go
type Speaker interface {
    Speak() string
}

type Dog struct{}

// value receiver
func (d Dog) Speak() string { return "woof" }

var s Speaker
s = Dog{}   // ใช้ได้
s = &Dog{}  // ใช้ได้เช่นกัน (Go จะ dereference อัตโนมัติ)
```

```go
type Walker interface {
    Walk() string
}

type Cat struct{}

// pointer receiver
func (c *Cat) Walk() string { return "walking" }

var w Walker
// w = Cat{}   // ผิด: Cat does not implement Walker (Walk method has pointer receiver)
w = &Cat{}     // ถูกต้อง
```

## สรุป
- Receiver เป็นกลไกหลักของ Go ในการผูก method กับ type
- Value receiver ใช้เมื่อต้องการอ่านค่า หรือเมื่อ type มีขนาดเล็ก
- Pointer receiver ใช้เมื่อต้องการแก้ไขค่า หรือ type มีขนาดใหญ่
- การเลือกชนิด receiver ส่งผลต่อการ implement interface และประสิทธิภาพ
- ไม่มี "class" ใน Go แต่ receiver + struct + interface ให้ความสามารถเทียบเท่า OOP ในรูปแบบที่ยืดหยุ่นและเรียบง่าย
---------------
### 6. **Goroutine & Channel**
- ใช้ `go` เพื่อเรียกฟังก์ชันทำงานแบบ concurrent (goroutine)
- `channel` ใช้สื่อสารระหว่าง goroutine แบบ type-safe

### 7. **Control Structures**
- `if`, `for`, `switch`, `defer` เป็นต้น
- ไม่มี `while` หรือ `do-while` แต่ใช้ `for` แทนได้ทั้งหมด

### 8. **Error Handling**
- ไม่มี exception; ใช้การคืนค่า error เป็นค่าสุดท้าย และตรวจสอบด้วย `if err != nil`

---

## การออกแบบ Workflow ใน Go

Workflow (กระแสการทำงาน) ใน Go สามารถออกแบบได้หลายรูปแบบขึ้นอยู่กับความต้องการ: sequential, concurrent, parallel, event-driven, หรือ state machine โดยใช้เครื่องมือหลัก ๆ ของ Go ดังนี้

### 1. **Sequential Workflow**
- เหมาะกับงานที่ต้องทำตามลำดับทีละขั้น
- เขียนเป็นฟังก์ชันเรียกต่อเนื่องกัน ตรวจสอบ error ทุกขั้นตอน

```go
func ProcessOrder(order Order) error {
    if err := validate(order); err != nil {
        return err
    }
    if err := save(order); err != nil {
        return err
    }
    return notify(order)
}
```

### 2. **Concurrent Workflow กับ Goroutine และ Channel**
- ใช้เมื่อมีหลายงานที่สามารถทำพร้อมกันได้ เช่น การเรียก API หลายแหล่ง หรือการประมวลผลข้อมูลแบบขนาน
- **Pattern: Worker Pool** – ควบคุมจำนวน goroutine ที่ทำงานพร้อมกัน

```go
func ProcessItems(items []Item) {
    jobs := make(chan Item)
    results := make(chan Result)
    
    // เริ่ม worker จำนวน N ตัว
    for w := 0; w < numWorkers; w++ {
        go worker(jobs, results)
    }
    
    // ส่งงาน
    for _, item := range items {
        jobs <- item
    }
    close(jobs)
    
    // เก็บผลลัพธ์
    for r := range results {
        // จัดการผลลัพธ์
    }
}
```

### 3. **Pipeline Pattern**
- ใช้ channel เชื่อมต่อหลายขั้นตอน (stage) โดยแต่ละขั้นตอนรับ input channel และส่ง output channel
- แต่ละ stage ทำงานแบบ concurrent

```go
func stage1(input <-chan Data) <-chan ProcessedData { ... }
func stage2(input <-chan ProcessedData) <-chan Result { ... }

func main() {
    input := make(chan Data)
    go func() { /* เติมข้อมูลลง input */ }()
    
    out1 := stage1(input)
    out2 := stage2(out1)
    
    for res := range out2 {
        fmt.Println(res)
    }
}
```

### 4. **Fan-in / Fan-out**
- **Fan-out**: กระจายงานไปยังหลาย goroutine (เช่น หลาย worker)
- **Fan-in**: รวมผลลัพธ์จากหลาย goroutine มาไว้ใน channel เดียว

### 5. **Workflow with Context**
- ใช้ `context.Context` เพื่อควบคุม timeout, cancellation และค่า shared ระหว่าง goroutine
- เป็นมาตรฐานในการออกแบบ workflow ที่ต้องหยุดทำงานเมื่อเกิด timeout หรือได้รับคำสั่งยกเลิก

```go
func ProcessWithTimeout(ctx context.Context, data Data) error {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    
    ch := make(chan Result)
    go func() {
        // ทำงานหนัก
        ch <- doWork(data)
    }()
    
    select {
    case res := <-ch:
        // จัดการผลลัพธ์
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}
```

### 6. **State Machine Workflow**
- เหมาะกับกระบวนการที่มีหลายสถานะ (เช่น ออเดอร์, การยืนยันตัวตน)
- มักใช้ struct ที่เก็บสถานะปัจจุบัน และ method สำหรับเปลี่ยนสถานะ พร้อมตรวจสอบความถูกต้อง

### 7. **Error Handling ใน Workflow**
- ใช้ error propagation ผ่าน channel (ส่ง error เป็นค่าคู่กับผลลัพธ์) หรือใช้ errgroup จาก `golang.org/x/sync/errgroup` เพื่อจัดการ goroutine พร้อมกันและรอผลลัพธ์/error แรก

### 8. **Testing Workflow**
- การออกแบบ workflow ที่ดีควรแยก business logic ออกจาก infrastructure (เช่น database, HTTP client) เพื่อให้สามารถทดสอบ unit test ได้ง่าย โดยใช้ interface และ dependency injection

---

### สรุป
โครงสร้างของ Go เน้นความเรียบง่าย ใช้ package เป็นหน่วยย่อย และมีเครื่องมือ concurrency ที่ทรงพลัง (goroutine, channel) การออกแบบ workflow ใน Go จึงมักใช้ประโยชน์จากเครื่องมือเหล่านี้เพื่อสร้างระบบที่ขนานงานได้ดี จัดการ error ได้ชัดเจน และควบคุมเวลา/การยกเลิกผ่าน context ซึ่งช่วยให้โค้ดอ่านง่าย บำรุงรักษาได้ และมีประสิทธิภาพสูง
-----------
#### 2.4 ตัวอย่าง Hello World ใน Go
```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

#### 2.5 ใครใช้ Go บ้าง?
- **Google** : ระบบ backend, Kubernetes, Docker
- **Uber** : ระบบการจับคู่การเดินทาง
- **Dropbox** : เปลี่ยนจาก Python มาเป็น Go สำหรับระบบจัดเก็บไฟล์
- **Netflix** : ส่วนของ proxy และ caching
- **Cloudflare** : เครื่องมือโครงสร้างพื้นฐาน

---

### บทที่ 3: พื้นฐานการใช้งาน Terminal

#### 3.1 Terminal คืออะไร?
Terminal (หรือ command line, console) เป็นเครื่องมือที่ให้เราสั่งงานคอมพิวเตอร์ผ่านข้อความ แทนการใช้ GUI การพัฒนา Go มักใช้ terminal ในการรันคำสั่งต่างๆ เช่น go build, go run, go test, git เป็นต้น

#### 3.2 คำสั่งพื้นฐาน (Unix/Linux/macOS)
- `pwd` : แสดงไดเรกทอรีปัจจุบัน
- `ls` : แสดงรายการไฟล์ (ใช้ `ls -la` แสดงรายละเอียด)
- `cd <path>` : เปลี่ยนไดเรกทอรี
- `mkdir <name>` : สร้างโฟลเดอร์
- `touch <file>` : สร้างไฟล์
- `rm <file>` : ลบไฟล์ (ใช้ `rm -rf` ลบโฟลเดอร์)
- `cp <source> <dest>` : คัดลอกไฟล์
- `mv <source> <dest>` : ย้ายหรือเปลี่ยนชื่อ
- `cat <file>` : แสดงเนื้อหาไฟล์
- `echo <text>` : แสดงข้อความ
- `grep <pattern> <file>` : ค้นหาข้อความ

#### 3.3 คำสั่งพื้นฐานสำหรับ Windows (Command Prompt หรือ PowerShell)
- `cd` : เปลี่ยนไดเรกทอรี
- `dir` : แสดงรายการไฟล์
- `mkdir` : สร้างโฟลเดอร์
- `del` : ลบไฟล์
- `copy` : คัดลอก
- `move` : ย้าย
- `type` : แสดงเนื้อหาไฟล์

#### 3.4 การตั้งค่า PATH
หลังจากติดตั้ง Go แล้ว เราต้องให้ terminal สามารถเรียกใช้คำสั่ง `go` ได้ โดยเพิ่มไดเรกทอรีของ Go (เช่น `/usr/local/go/bin` หรือ `C:\Go\bin`) ลงในตัวแปรสภาพแวดล้อม PATH

#### 3.5 การใช้ go command
- `go version` : ตรวจสอบเวอร์ชัน Go
- `go env` : แสดงตัวแปรสภาพแวดล้อมของ Go
- `go run <file>` : คอมไพล์และรันโปรแกรม
- `go build` : คอมไพล์เป็น binary
- `go mod init` : เริ่มต้น go module
- `go get` : ดาวน์โหลดแพคเกจ (deprecated ใช้ `go install` หรือ `go mod download`)
- `go test` : รัน test

---

### บทที่ 4: เตรียมสภาพแวดล้อมสำหรับพัฒนา

#### 4.1 การติดตั้ง Go
1. ดาวน์โหลดจาก https://go.dev/dl/
2. เลือกเวอร์ชันที่เหมาะสมกับ OS ของคุณ
3. ติดตั้งตามขั้นตอน
4. ตรวจสอบการติดตั้ง: เปิด terminal แล้วพิมพ์ `go version`

#### 4.2 ตัวแปรสภาพแวดล้อมที่สำคัญ
- `GOROOT` : ตำแหน่งที่ติดตั้ง Go (มักตั้งค่าให้อัตโนมัติ)
- `GOPATH` : ตำแหน่ง workspace ของ Go (เลิกใช้แล้ว ถ้าใช้ modules)
- `GOBIN` : ตำแหน่งที่ติดตั้ง binaries ที่ `go install`
- `GOOS`/`GOARCH` : ใช้ cross-compile

#### 4.3 การเลือก Editor/IDE
- **Visual Studio Code** : มี extension "Go" โดย官方 ให้ linting, autocomplete, debugging
- **GoLand** : IDE เฉพาะของ JetBrains
- **Vim/Neovim** : ใช้ plugin vim-go
- **Sublime Text** : มี GoSublime

#### 4.4 การตั้งค่า Go Modules
Go Modules เป็นระบบจัดการ dependencies เริ่มใช้ตั้งแต่ Go 1.11 โดยค่าเริ่มต้นใช้งานได้ตั้งแต่ Go 1.16 เป็นต้นไป
```bash
mkdir myproject
cd myproject
go mod init example.com/myproject
```
ไฟล์ `go.mod` จะถูกสร้างขึ้น

#### 4.5 workspace structure (แบบเดิมกับแบบใหม่)
- **แบบ GOPATH** (ก่อน Go 1.11): ต้องวางโค้ดใน `$GOPATH/src`
- **แบบ Modules** (ปัจจุบัน): สามารถสร้างโปรเจกต์ได้ทุกที่ โดยมี go.mod กำกับ

#### 4.6 การติดตั้งเครื่องมือเสริม
- **gopls** : language server (ติดตั้งโดยอัตโนมัติใน VS Code)
- **dlv** : delve debugger (`go install github.com/go-delve/delve/cmd/dlv@latest`)
- **golangci-lint** : linter (`go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`)

---

### บทที่ 5: สร้างแอปพลิเคชันแรกของคุณ

#### 5.1 สร้างโปรเจกต์
```bash
mkdir hello
cd hello
go mod init hello
```

#### 5.2 เขียนโค้ด
สร้างไฟล์ `main.go`:
```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, Go!")
}
```

#### 5.3 รันโปรแกรม
```bash
go run main.go
```
หรือคอมไพล์แล้วรัน:
```bash
go build -o hello
./hello   # (Linux/macOS) หรือ hello.exe (Windows)
```

#### 5.4 เพิ่มฟังก์ชัน
```go
package main

import "fmt"

func greet(name string) string {
    return "Hello, " + name
}

func main() {
    message := greet("Gopher")
    fmt.Println(message)
}
```
ทดสอบ `go run main.go` ควรได้ "Hello, Gopher"

#### 5.5 การอ่านค่าจาก command line
```go
package main

import (
    "fmt"
    "os"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Please provide a name")
        return
    }
    name := os.Args[1]
    fmt.Printf("Hello, %s!\n", name)
}
```
รัน: `go run main.go Somchai`

#### 5.6 การจัดรูปแบบโค้ด
ใช้คำสั่ง `go fmt` เพื่อจัดรูปแบบโค้ดให้เป็นมาตรฐาน:
```bash
go fmt ./...
```

#### 5.7 ข้อผิดพลาดที่พบบ่อย
- ใช้ตัวแปรแต่ไม่ได้ใช้ -> Go จะคอมไพล์ไม่ผ่าน (ยกเว้นตัวแปร _)
- import แพคเกจแต่ไม่ได้ใช้ -> ต้องลบออกหรือใช้ _
- ใส่เครื่องหมาย `,` หรือ `;` ไม่ถูกต้อง

---

## ภาคที่ 2: พื้นฐานภาษาและโครงสร้างข้อมูล

### บทที่ 6: ระบบเลขฐานสองและฐานสิบ

#### 6.1 ความสำคัญของระบบเลข
คอมพิวเตอร์ทำงานด้วยสัญญาณไฟฟ้า 2 สถานะ คือ ON (1) และ OFF (0) ดังนั้นการแทนค่าทั้งหมดในคอมพิวเตอร์จึงใช้เลขฐานสอง

#### 6.2 เลขฐานสิบ (Decimal)
เป็นระบบที่มนุษย์ใช้ประจำ มีตัวเลข 0-9 (ฐาน 10) เช่น 123 = 1×10² + 2×10¹ + 3×10⁰

#### 6.3 เลขฐานสอง (Binary)
ใช้ตัวเลข 0 และ 1 (ฐาน 2) เช่น 1101₂ = 1×2³ + 1×2² + 0×2¹ + 1×2⁰ = 13₁₀

#### 6.4 การแปลงเลขฐานสองเป็นฐานสิบ
วิธี: คูณแต่ละหลักด้วย 2 ยกกำลังตำแหน่ง (เริ่มจากขวาเป็น 0)
- 1010₂ = (1×8) + (0×4) + (1×2) + (0×1) = 10₁₀

#### 6.5 การแปลงเลขฐานสิบเป็นฐานสอง
วิธี: หารด้วย 2 ไปเรื่อยๆ เก็บเศษ แล้วอ่านเศษจากล่างขึ้นบน
- 13 ÷ 2 = 6 เศษ 1
- 6 ÷ 2 = 3 เศษ 0
- 3 ÷ 2 = 1 เศษ 1
- 1 ÷ 2 = 0 เศษ 1
→ อ่านจากล่างขึ้นบน: 1101₂

#### 6.6 การแสดงเลขฐานสองใน Go
Go รองรับ literal ของเลขฐานสองโดยใช้ prefix `0b` หรือ `0B`:
```go
var x int = 0b1101 // 13
fmt.Printf("%b\n", x) // แสดงเป็น binary: 1101
```

#### 6.7 หน่วยของข้อมูล
- 1 bit = 0 หรือ 1
- 1 byte = 8 bits
- 1 kilobyte (KB) = 1024 bytes
- 1 megabyte (MB) = 1024 KB
- Go มีชนิดข้อมูลตามขนาด: uint8 (1 byte), uint16, uint32, uint64

#### 6.8 ระบบฐานอื่นๆ
- ฐานแปด (octal): ใช้ prefix `0o` หรือ `0O` (เช่น 0o17 = 15)
- ฐานสิบหก (hexadecimal): ใช้ prefix `0x` (เช่น 0xFF = 255)

---

### บทที่ 7: เลขฐานสิบหก, ฐานแปด, ASCII, UTF8, Unicode และ Runes

#### 7.1 เลขฐานสิบหก (Hexadecimal)
ใช้ตัวเลข 0-9 และตัวอักษร A-F (10-15) ฐาน 16 นิยมใช้แทนสี, address ในหน่วยความจำ
- 0xFF = 255, 0x10 = 16
- ใน Go: `color := 0xFF00FF`

#### 7.2 เลขฐานแปด (Octal)
ใช้ตัวเลข 0-7 ฐาน 8
- 0o10 = 8, 0o777 = 511

#### 7.3 ASCII
American Standard Code for Information Interchange เป็นรหัส 7 bits แทนตัวอักษรภาษาอังกฤษ, ตัวเลข, เครื่องหมาย 128 ตัว
- 'A' = 65, 'a' = 97, '0' = 48

#### 7.4 Unicode
เป็นมาตรฐานสากลที่กำหนดรหัสให้กับอักขระทุกภาษา (มากกว่า 1 ล้านรหัส) โดยแต่ละรหัสเรียกว่า "code point" เช่น U+0041 คือ 'A'

#### 7.5 UTF-8
UTF-8 เป็นการเข้ารหัส Unicode แบบ variable-length (1-4 bytes) ที่ Go ใช้เป็นค่าเริ่มต้นสำหรับสตริง
- ตัวอักษร ASCII ใช้ 1 byte
- ตัวอักษรไทย 3 bytes

#### 7.6 Rune ใน Go
`rune` เป็น alias ของ `int32` ใช้แทน Unicode code point ใน Go
```go
var ch rune = 'CD' // รหัส Unicode 0xE01
fmt.Printf("%c %U\n", ch, ch) //   U+0E01
```
การวนลูปผ่าน string แบบ rune:
```go
s := "hello"
for i, r := range s {
    fmt.Printf("%d: %c (%U)\n", i, r, r)
}
```

#### 7.7 การจัดการ string และ bytes
String ใน Go เป็น read-only slice of bytes (ไม่ใช่ runes) การนับจำนวน runes ใช้ `utf8.RuneCountInString`
```go
import "unicode/utf8"
s := "สวัสดี"
fmt.Println(len(s)) // 12 bytes
fmt.Println(utf8.RuneCountInString(s)) // 4 runes
```

#### 7.8 การแปลงระหว่าง string และ []rune
```go
s := "ไทย"
runes := []rune(s)   // แปลงเป็น slice of runes
s2 := string(runes)  // แปลงกลับ
```

---

### บทที่ 8: ตัวแปร, ค่าคงที่ และชนิดข้อมูลพื้นฐาน

#### 8.1 ตัวแปร (Variables)
ตัวแปรคือที่เก็บข้อมูลในหน่วยความจำ มีชื่อและชนิดข้อมูล
```go
var name string = "John"
var age int = 30
var isActive bool = true
```
รูปแบบ: `var <ชื่อ> <ชนิด> = <ค่าเริ่มต้น>`

การประกาศแบบสั้นด้วย `:=` (ใช้ภายในฟังก์ชันเท่านั้น):
```go
name := "John"   // Go infer type เป็น string
age := 30
```

การประกาศหลายตัวแปร:
```go
var x, y int = 1, 2
var (
    a string = "hello"
    b int = 42
)
```

#### 8.2 ค่าคงที่ (Constants)
ค่าคงที่ประกาศด้วย `const` ไม่สามารถเปลี่ยนแปลงค่าได้
```go
const Pi = 3.14159
const (
    StatusOK = 200
    StatusNotFound = 404
)
```
ค่าคงที่สามารถเป็นตัวเลข, สตริง, หรือ boolean เท่านั้น

#### 8.3 ชนิดข้อมูลพื้นฐาน
**Numeric types:**
- Integers: `int`, `int8`, `int16`, `int32`, `int64` (int ขนาดตามระบบ 32/64 bit)
- Unsigned: `uint`, `uint8` (byte), `uint16`, `uint32`, `uint64`
- Floating point: `float32`, `float64`
- Complex: `complex64`, `complex128`

**Boolean:** `bool` (true/false)

**String:** `string` (immutable sequence of bytes)

**Rune:** `rune` (int32) ใช้แทน Unicode code point

**Byte:** `byte` (uint8) ใช้แทน byte

#### 8.4 ค่าเริ่มต้น (Zero values)
ตัวแปรทุกตัวเมื่อประกาศโดยไม่กำหนดค่า จะมีค่าเริ่มต้น:
- ตัวเลข: 0
- bool: false
- string: "" (empty string)
- pointer, slice, map, channel, interface, function: nil

#### 8.5 การแปลงชนิดข้อมูล (Type conversion)
Go ต้องแปลงชนิดอย่างชัดเจน ไม่มี implicit conversion
```go
var i int = 42
var f float64 = float64(i)   // แปลง int -> float64
var u uint = uint(f)          // แปลง float64 -> uint
```

#### 8.6 การกำหนดชื่อตัวแปร
- ขึ้นต้นด้วยตัวอักษรหรือ _
- ตามด้วยตัวอักษร ตัวเลข หรือ _
- case-sensitive
- นิยมใช้ camelCase (เช่น userName)
- ตัวพิมพ์ใหญ่ต้น = exported (public) ถ้าอยู่ในแพคเกจ

#### 8.7 ตัวดำเนินการพื้นฐาน
- คณิตศาสตร์: `+`, `-`, `*`, `/`, `%`
- การเปรียบเทียบ: `==`, `!=`, `<`, `>`, `<=`, `>=`
- ตรรกะ: `&&`, `||`, `!`
- Bitwise: `&`, `|`, `^`, `&^`, `<<`, `>>`

---

### บทที่ 9: คำสั่งควบคุมการทำงาน

#### 9.1 if-else
```go
if score >= 80 {
    fmt.Println("A")
} else if score >= 70 {
    fmt.Println("B")
} else {
    fmt.Println("C")
}
```
สามารถมีคำสั่งสั้น ๆ ก่อนเงื่อนไข:
```go
if val := getValue(); val > 10 {
    fmt.Println("big")
}
```

#### 9.2 switch
Go มี switch ที่ยืดหยุ่น:
```go
switch day {
case "Monday":
    fmt.Println("Start of week")
case "Friday":
    fmt.Println("TGIF")
default:
    fmt.Println("Other day")
}
```
switch แบบไม่มี expression (เทียบกับ true):
```go
switch {
case score >= 80:
    fmt.Println("A")
case score >= 70:
    fmt.Println("B")
default:
    fmt.Println("C")
}
```
fallthrough (ใช้ในกรณีต้องการให้ทำงาน case ถัดไป):
```go
switch num {
case 1:
    fmt.Println("one")
    fallthrough
case 2:
    fmt.Println("two")
}
// ถ้า num = 1 จะพิมพ์ one และ two
```

#### 9.3 for loop
Go มีแค่ `for` (ไม่มี while, do-while)
```go
// for แบบคลาสสิก
for i := 0; i < 10; i++ {
    fmt.Println(i)
}

// while-like
x := 0
for x < 10 {
    fmt.Println(x)
    x++
}

// infinite loop
for {
    // ใช้ break เพื่อออก
}
```

#### 9.4 range
ใช้กับ array, slice, map, string:
```go
numbers := []int{2,4,6}
for index, value := range numbers {
    fmt.Println(index, value)
}

// ถ้าไม่ต้องการ index ใช้ _
for _, value := range numbers {
    fmt.Println(value)
}
```

#### 9.5 break, continue, goto
- `break` : ออกจาก loop
- `continue` : ข้ามไป iteration ถัดไป
- `goto` : กระโดดไปยัง label (ไม่แนะนำให้ใช้มาก)

```go
for i := 0; i < 5; i++ {
    if i == 2 {
        continue
    }
    if i == 4 {
        break
    }
    fmt.Println(i)
}
// output: 0 1 3
```

#### 9.6 label และการ break ออกจาก outer loop

**Outer loop** (หรือ loop ภายนอก) คือลูปที่อยู่ชั้นนอกสุดในโครงสร้างของลูปซ้อนกัน (nested loops) หมายถึงลูปที่มีลูปอื่นอยู่ข้างใน ใช้สำหรับควบคุมการวนรอบที่มีหลายมิติ เช่น การวนสมาชิกในอาร์เรย์ 2 มิติ หรือการทำซ้ำงานที่ต้องมีลูปย่อยหลายรอบในแต่ละรอบของลูปหลัก

## ตัวอย่างในภาษา Go

```go
package main

import "fmt"

func main() {
    // outer loop วนแถว (row)
    for i := 0; i < 3; i++ {
        // inner loop วนคอลัมน์ (column)
        for j := 0; j < 3; j++ {
            fmt.Printf("(%d,%d) ", i, j)
        }
        fmt.Println()
    }
}
```

**ผลลัพธ์:**
```
(0,0) (0,1) (0,2) 
(1,0) (1,1) (1,2) 
(2,0) (2,1) (2,2) 
```

## การควบคุม outer loop ด้วย `break` และ `continue`

ใน Go สามารถใช้ **label** เพื่อระบุว่าต้องการ `break` หรือ `continue` ที่ outer loop แทนที่จะเป็นแค่ inner loop:

```go
package main

import "fmt"

func main() {
    // ตั้ง label ชื่อ "outer"
outer:
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            if i == 1 && j == 1 {
                break outer // ออกจาก outer loop ทั้งหมด
            }
            fmt.Printf("(%d,%d) ", i, j)
        }
        fmt.Println()
    }
}
```

**ผลลัพธ์:**
```
(0,0) (0,1) (0,2) 
(1,0) 
```
(เมื่อเจอ (1,1) จะออกจาก outer loop ทันที)

## การใช้งานทั่วไป
- วนซ้ำข้อมูลในโครงสร้างหลายมิติ (matrix, slice of slices)
- การทำอัลกอริทึมที่ต้องตรวจสอบเงื่อนไขข้ามหลายชั้น (เช่น การค้นหาด้วย brute force)
- การทำงานที่ต้องทำซ้ำเป็นรอบใหญ่และภายในแต่ละรอบมีขั้นตอนย่อยหลายขั้น

## ข้อควรระวัง
- การใช้ `break` โดยไม่มี label จะออกแค่ inner loop เท่านั้น
- ควรตั้งชื่อ label ให้สื่อความหมาย (เช่น `outer`, `mainLoop`) เพื่อให้โค้ดอ่านง่าย
- การใช้ `goto` แทน label สำหรับ `break`/`continue` ไม่แนะนำให้ใช้ในกรณีนี้

สรุป: **outer loop** คือลูปหลักที่ครอบลูปอื่นอยู่ภายใน ใช้สำหรับควบคุมการวนรอบในระดับสูงสุดของโครงสร้างซ้อนกัน

---

### บทที่ 10: ฟังก์ชัน

#### 10.1 การประกาศฟังก์ชัน
```go
func add(x int, y int) int {
    return x + y
}
```
ถ้าพารามิเตอร์ชนิดเดียวกัน สามารถย่อ:
```go
func add(x, y int) int {
    return x + y
}
```
ฟังก์ชัน `add` ในภาษา Go เป็นตัวอย่างพื้นฐานของฟังก์ชันที่รับพารามิเตอร์สองตัวและคืนค่าผลรวม โดยมีรายละเอียดดังนี้:

### โครงสร้างของฟังก์ชัน
```go
func add(x int, y int) int {
    return x + y
}
```
- **`func`** – คำสงวนสำหรับประกาศฟังก์ชัน  
- **`add`** – ชื่อฟังก์ชัน  
- **`(x int, y int)`** – รายการพารามิเตอร์ ระบุชื่อและชนิดข้อมูล (ทั้งคู่เป็น `int`)  
- **`int`** – ชนิดของค่าที่คืนกลับ  
- **`return x + y`** – คำนวณผลรวมของ `x` และ `y` แล้วส่งกลับไปยังผู้เรียก

### ตัวอย่างการใช้งาน

#### 1. เรียกใช้โดยตรงใน `main`
```go
package main

import "fmt"

func add(x int, y int) int {
    return x + y
}

func main() {
    result := add(5, 3)
    fmt.Println(result) // 8
}
```

#### 2. นำผลลัพธ์ไปใช้ต่อ
```go
a := 10
b := 20
sum := add(a, b) + 5   // 30 + 5 = 35
fmt.Println(sum)       // 35
```

#### 3. ส่งฟังก์ชันเป็นค่า (first-class function)
```go
func apply(fn func(int, int) int, a, b int) int {
    return fn(a, b)
}

func main() {
    fmt.Println(apply(add, 7, 8)) // 15
}
```

#### 4. ใช้ในนิพจน์แบบสั้น (short variable declaration)
```go
total := add(100, 200) // total = 300
```
```go
package main

import "fmt"

// add returns the sum of two integers
func add(x, y int) int {
    return x + y
}

// addFloat is a generic version that works with any numeric type (requires Go 1.18+)
func addGeneric[T int | float64](x, y T) T {
    return x + y
}

func main() {
    // ----- Using add with ints -----
    total := 10          // total is int
    // Correct: use = for reassignment, not :=
    total = add(5, 3)    // total becomes 8
    fmt.Println("total (int):", total)

    // ----- Using add with float64 (explicit conversion) -----
    var f float64
    // Error previously: cannot use add(5,3) (int) as float64
    // Fix: convert the int result to float64
    f = float64(add(5, 3))
    fmt.Println("f (float64):", f)

    // ----- Alternative: using generic function with floats directly -----
    f2 := addGeneric(2.5, 3.7) // f2 is float64
    fmt.Println("f2 (generic float):", f2)

    // ----- If you need to mix int and float, convert explicitly -----
    var mixed float64
    mixed = float64(add(5, 3)) + 2.5
    fmt.Println("mixed:", mixed)

    // ----- Short variable declaration with a new variable -----
    sum := add(10, 20) // new variable, okay
    fmt.Println("sum:", sum)
}
```
### ข้อสังเกต
- ถ้าพารามิเตอร์มีชนิดเดียวกันสามารถเขียนย่อเป็น `func add(x, y int) int` ได้
- การคืนค่าเป็นแบบ **explicit return** (ระบุค่าโดยตรง) ซึ่งเป็นวิธีที่นิยมใน Go
- ฟังก์ชันนี้เป็น pure function เพราะไม่มีผลข้างเคียง ขึ้นกับค่าอินพุตอย่างเดียว

 
#### 10.2 การคืนค่าหลายค่า
Go รองรับการคืนค่าหลายค่าได้โดยตรง:
```go
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("division by zero")
    }
    return a / b, nil
}
```
เรียกใช้:
```go
result, err := divide(10, 2)
if err != nil {
    fmt.Println("Error:", err)
} else {
    fmt.Println("Result:", result)
}
```

#### 10.3 Named return values
สามารถตั้งชื่อค่าที่คืนได้:
```go
func split(sum int) (x, y int) {
    x = sum * 4 / 9
    y = sum - x
    return // naked return
}
```
ฟังก์ชัน `sum` ในภาษา Go ใช้ **variadic parameter** (`...int`) เพื่อรับจำนวนอาร์กิวเมนต์ `int` ตั้งแต่ 0 ตัวขึ้นไป โดยอาร์กิวเมนต์ทั้งหมดจะถูกจัดเก็บใน `nums` ซึ่งเป็น slice ของ `int`

### การทำงาน:
1. **ประกาศตัวแปร `total`** เริ่มต้นที่ `0`
2. **วนลูปด้วย `range`** เพื่อเข้าถึงแต่ละค่าของ `nums` โดยไม่สนใจ index (ใช้ `_` แทน)
   - แต่ละรอบ `n` จะได้รับค่าของสมาชิกใน slice
   - `total += n` บวกสะสมเข้า `total`
3. **คืนค่า `total`** ซึ่งเป็นผลรวมของอาร์กิวเมนต์ทั้งหมด

### ตัวอย่างการเรียกใช้:
```go
fmt.Println(sum(1, 2, 3))   // 6
fmt.Println(sum())          // 0 (ไม่มีอาร์กิวเมนต์)
```

### หมายเหตุ:
- หากมี slice อยู่แล้ว สามารถส่งเข้าไปในฟังก์ชันโดยใช้ `...` เช่น `numbers := []int{1,2,3}; sum(numbers...)`
- Variadic parameter ต้องเป็นพารามิเตอร์ตัวสุดท้ายของฟังก์ชัน
---------------------
ฟังก์ชัน `split` ในภาษา Go ทำงานโดยใช้ **named return values** (`x, y int`) และ **naked return** (คืนค่าโดยไม่ระบุตัวแปร) เพื่อส่งค่ากลับไปยังผู้เรียก

### กระบวนการทำงานทีละขั้นตอน:
1. **รับพารามิเตอร์** `sum` (ชนิด `int`)
2. **คำนวณค่า `x`** ด้วย `sum * 4 / 9`  
   - เนื่องจากตัวเลขทั้งหมดเป็น `int` การหารเป็นแบบ integer division (ตัดเศษทิ้ง)
3. **คำนวณค่า `y`** ด้วย `sum - x`  
   - `y` จะเป็นค่าที่เหลือเมื่อนำ `sum` ลบด้วย `x`
4. **คืนค่า** ด้วย `return` แบบเปล่า (naked return)  
   - ระบบจะคืนค่าปัจจุบันของ `x` และ `y` อัตโนมัติ ตามลำดับที่ประกาศไว้

### ตัวอย่างการเรียกใช้:
```go
a, b := split(10)
fmt.Println(a, b) // ผลลัพธ์: 4 6
```

**การคำนวณ:**
- `x = 10 * 4 / 9 = 40 / 9 = 4` (เศษ 4 ถูกตัดทิ้ง)
- `y = 10 - 4 = 6`

### ข้อควรรู้:
- Naked return เหมาะกับฟังก์ชันสั้น ๆ ที่มีชื่อตัวแปรคืนค่าชัดเจน ช่วยให้โค้ดกระชับ
- หากฟังก์ชันยาวหรือซับซ้อน ควรเขียน `return x, y` แบบเต็มเพื่อความเข้าใจง่าย
- การใช้ integer division ต้องระวังเรื่องการปัดเศษลง (floor) ซึ่งอาจไม่ตรงกับที่คาดหวังถ้าต้องการค่าทศนิยม
---------------------

**ข้อควรระวัง**: naked return อาจลดความชัดเจน

#### 10.4 ฟังก์ชันแบบ variadic (รับพารามิเตอร์ไม่จำกัด)
```go
func sum(nums ...int) int {
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}
// เรียกใช้: sum(1,2,3,4)
```

ฟังก์ชัน `sum(nums ...int) int` รับอาร์กิวเมนต์จำนวนเท่าใดก็ได้ (variadic parameter) ซึ่งจะถูกเก็บใน `nums` เป็น slice ของ `int` ฟังก์ชันใช้ `for range` วนบวกค่าทุกตัวแล้วคืนผลรวม

**ตัวอย่างการเรียกใช้:**
```go
s1 := sum(1, 2, 3, 4)       // 10
s2 := sum()                 // 0 (ไม่มีอาร์กิวเมนต์)
s3 := sum(5, 10)            // 15

numbers := []int{2, 4, 6}
s4 := sum(numbers...)       // 12 (ส่ง slice โดยใช้ ...)

fmt.Println("s4:", s4)

```
#### 10.5 ฟังก์ชันเป็นค่า (first-class functions)
ฟังก์ชันสามารถถูกกำหนดให้กับตัวแปร, ส่งเป็นพารามิเตอร์, คืนค่าเป็นผลลัพธ์
```go
var fn func(int) int = func(x int) int { return x * 2 }
result := fn(5) // 10
```
ในภาษา Go ฟังก์ชันเป็น **first-class citizen** สามารถกำหนดให้ตัวแปรเก็บฟังก์ชันได้  

```go
var fn func(int) int = func(x int) int { return x * 2 }
```

- `var fn func(int) int` ประกาศตัวแปร `fn` ชนิด **function type** รับ `int` คืน `int`  
- `= func(x int) int { return x * 2 }` กำหนด **anonymous function** (ฟังก์ชันนิรนาม) ที่รับ `x` แล้วคืน `x*2`  
- จากนั้นเรียกใช้ผ่าน `fn(5)` ได้ผลลัพธ์ `10`  

### ตัวอย่างเพิ่มเติม
```go
// ประกาศแบบสั้น (type inference)
square := func(n int) int { return n * n }
fmt.Println(square(4)) // 16

// ส่งฟังก์ชันเป็น argument
func apply(fn func(int) int, val int) int {
    return fn(val)
}
result := apply(func(x int) int { return x + 10 }, 7) // 17
```

ประโยชน์: ช่วยให้เขียนโค้ดที่ยืดหยุ่น (higher-order functions, callback)

# 10.6 defer

`defer` คือคำสั่งในภาษา Go ที่ใช้ **เลื่อนการทำงานของฟังก์ชัน** ออกไปจนกว่าฟังก์ชันรอบนอก (ฟังก์ชันที่ประกาศ `defer`) จะจบการทำงาน ไม่ว่าจะจบแบบปกติ (return) หรือเกิด panic (panic คือ หยุดการทำงาน)
`defer` ทำให้ฟังก์ชันถูกเรียกหลังจากฟังก์ชัน enclosing จบการทำงาน (ใช้สำหรับ cleanup)
### หลักการทำงาน
- คำสั่งที่อยู่ภายใต้ `defer` จะถูกเก็บไว้ในสแต็ก (stack) และทำงานในลำดับ **LIFO** (Last In, First Out) – ตัวที่ประกาศทีหลังจะถูกเรียกก่อน
- นิยมใช้เพื่อปิดทรัพยากร (ไฟล์, database connection, mutex unlock) เพื่อป้องกันการรั่วไหล

### ตัวอย่าง
```go
func readFile() {
    f, err := os.Open("data.txt")
    if err != nil {
        return
    }
    defer f.Close() // ปิดไฟล์เมื่อฟังก์ชันจบ

    // อ่านข้อมูล...
}
```
### สรุป
`defer` ช่วยให้โค้ดสะอาด รับประกันการ cleanup แม้มี `return` ก่อนถึงคำสั่งปิด หรือเกิด panic (panic คือ หยุดการทำงาน)
`panic` เป็น built-in function ในภาษา Go ที่ใช้หยุดการทำงานปกติของโปรแกรม (คล้ายกับ exception ในภาษาอื่น) เมื่อ panic เกิดขึ้น:

1. การทำงานของฟังก์ชันปัจจุบันจะหยุดทันที  
2. เริ่ม defer statements (ในฟังก์ชันปัจจุบัน) ตามลำดับ LIFO  
3. กลับไปยังฟังก์ชันผู้เรียก และทำซ้ำ (unwind stack)  
4. หากไม่มี `recover` จับไว้ โปรแกรมจะจบการทำงานพร้อมแสดง stack trace  

**ตัวอย่าง:**
```go
func mayPanic() {
    defer fmt.Println("defer 1")
    defer fmt.Println("defer 2")
    panic("เกิดข้อผิดพลาดร้ายแรง")
    fmt.Println("ไม่") // จะไม่ถูกทำงานต่อไป

func main() {
    mayPanic()
    fmt.Println("ไม่ถึงบรรทัดนี้") // จะไม่ถูกทำงานต่อไป
}
```
ผลลัพธ์:
```
defer 2
defer 1
panic: เกิดข้อผิดพลาดร้ายแรง
...
```

**การกู้คืนด้วย `recover`**  
`recover` จะคืนค่าที่ส่งไปใน `panic` (ถ้ามี) และสามารถทำให้โปรแกรมทำงานต่อได้ โดยต้องเรียกใน `defer`:

```go
func safeCall() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("กู้คืนจาก panic:", r)
        }
    }()
    panic("error")
}

func main() {
    safeCall()
    fmt.Println("โปรแกรมทำงานต่อ")
}
```

**สรุป**  
- `panic` หยุดการทำงานปกติและ unwind stack  (การคลายสแต็ก)
- `defer` ยังคงถูกเรียกแม้ panic  
  **การคลายสแต็ก (unwind stack)** คือกระบวนการที่โปรแกรมย้อนกลับไปตามลำดับการเรียกฟังก์ชัน (call stack) โดย คำสั่ง `defer` ในแต่ละฟังก์ชันก่อนที่จะยุติฟังก์ชันนั้น ๆ มักเกิดขึ้นเมื่อเกิด `panic` หรือสิ้นสุดการทำงานของ goroutine
- ใช้ `recover` ใน `defer` เพื่อจัดการ panic และให้โปรแกรมดำเนินต่อไปได้
  **Unwind stack** (การคลายสแต็ก) คือกระบวนการที่เกิดขึ้นเมื่อเกิด `panic` ใน Go โดยโปรแกรมจะ **ย้อนกลับขึ้นไปตามลำดับการเรียกใช้ฟังก์ชัน** (call stack) และ  (execute or execution)  `defer` statements ที่ถูกค้างไว้ในแต่ละฟังก์ชันก่อนที่จะยุติการทำงานของฟังก์ชันนั้น ๆ

**Execute** หรือ **Execution** ในบริบทการเขียนโปรแกรม หมายถึง **การรันคำสั่ง** หรือ **การทำให้โปรแกรมทำงานตามที่เขียนไว้**  

- **Execute** (verb): ดำเนินการรันคำสั่ง, ฟังก์ชัน, หรือโปรแกรม  
- **Execution** (noun): กระบวนการที่โปรแกรมหรือคำสั่งกำลังทำงาน  

เมื่อเกิด `panic` การ **execute** ปกติจะหยุด แต่ **deferred functions** ยังคงถูก **execute** ระหว่าง unwind stack  

 
---

### แผนภาพการไหลของข้อมูล (Data Flow) เมื่อเกิด panic
แบ่งเป็น 3 แผนภาพหลัก:
---

## 1. แผนภาพการไหลของข้อมูล (Data Flow) เมื่อเกิด panic  (หยุดการทำงาน)

```mermaid
graph TD
    subgraph main
        M1[defer f1]
        M2[call level1]
        M3[รับ panic]
        M4{มี recover?}
        M5[handle]
        M6[exit]
    end

    subgraph level1
        L1[defer f2]
        L2[call level2]
        L3[รับ panic]
        L4{มี recover?}
        L5[handle]
        L6[unwind]
    end

    subgraph level2
        N1[defer f3]
        N2[call level3]
        N3[panic value]
        N4{มี recover?}
        N5[handle]
        N6[unwind]
    end

    M2 --> L2
    L2 --> N2
    N2 --> N3
    N3 -->|panic ส่งขึ้น| L3
    L3 --> L4
    L4 -->|ไม่มี| L6
    L6 -->|panic ส่งขึ้น| M3
    M3 --> M4
    M4 -->|ไม่มี| M6
    M4 -->|มี| M5
    L4 -->|มี| L5
    N4 -->|มี| N5

    style N3 fill:#f99,stroke:#333
    style L3 fill:#f99,stroke:#333
    style M3 fill:#f99,stroke:#333
```

**คำอธิบาย**:  
- แผนภาพนี้แสดงการเกิด panic ที่ `level2` (หรือ `level3`)  
- ค่า panic จะไหลขึ้นบนผ่านแต่ละฟังก์ชัน  
- แต่ละระดับสามารถมี `recover` ใน `defer` เพื่อจับ panic และหยุดการส่งต่อ  

---

## 2. แผนภาพ Call Stack และการไหลของ panic (แบบลำดับเวลา)

```mermaid
sequenceDiagram
    participant main
    participant L1 as level1
    participant L2 as level2
    participant L3 as level3

    main->>L1: เรียก level1()
    activate L1
    L1->>L2: เรียก level2()
    activate L2
    L2->>L3: เรียก level3()
    activate L3
    L3-->>L3: panic เกิดขึ้น (ค่า "error")
    Note over L3: execute defer (level3)
    L3-->>L2: panic ส่งขึ้น
    deactivate L3
    Note over L2: execute defer (level2)
    L2-->>L1: panic ส่งขึ้น
    deactivate L2
    Note over L1: execute defer (level1)
    L1-->>main: panic ส่งขึ้น
    deactivate L1
    Note over main: execute defer (main)<br/>โปรแกรมหยุด (ถ้าไม่มี recover)
```

---

## 3. แผนภาพแบบมี `recover` หยุดการ unwind (การคลายสแต็ก)

```mermaid
sequenceDiagram
    participant main
    participant L1 as level1
    participant L2 as level2
    participant L3 as level3

    main->>L1: เรียก level1()
    activate L1
    L1->>L2: เรียก level2()
    activate L2
    L2->>L3: เรียก level3()
    activate L3
    L3-->>L3: panic เกิดขึ้น (ค่า "error")
    Note over L3: execute defer (level3)
    L3-->>L2: panic ส่งขึ้น
    deactivate L3
    Note over L2: defer มี recover()<br/>จับ panic และหยุด unwind
    L2-->>L1: return ปกติ (ไม่มี panic)
    deactivate L2
    L1-->>main: return ปกติ
    deactivate L1
    Note over main: โปรแกรมทำงานต่อ
```
 
--- 

### คำอธิบาย
- **execute defer (levelX)** = การ( execution ) ฟังก์ชันที่ถูก `defer` ไว้ในขณะที่เกิด `panic` และกำลัง unwind stack  
- ลูกศรแนวตั้ง (`|`) แสดงลำดับการเรียกฟังก์ชันลงไป  
- ลูกศรแนวนอน (`<-----`) แสดงการส่ง panic value ขึ้นไปยัง caller พร้อมกับการ unwind  
- เมื่อ panic เกิดขึ้นใน `level3` การ execute ปกติจะหยุด แต่ `defer` ในแต่ละระดับยังคงถูก execute ก่อนที่ panic จะถูกส่งขึ้นไป

### ลำดับการไหลของข้อมูล (Data Flow)

1. **panic เกิดที่ `level3`** → ค่า panic (เช่น `"error"`) ถูกสร้างขึ้น  
2. **unwind stack เริ่ม** – การทำงานปกติใน `level3` หยุด; `defer` ทั้งหมดใน `level3` ทำงาน  
3. panic value **ไหลขึ้น** ไปยัง `level2`  
4. `level2` ( execution )  `defer` ของตัวเอง (LIFO)  
5. panic value **ไหลขึ้น** ไปยัง `level1`  
6. `level1` ( execution )  `defer`  
7. panic value **ไหลขึ้น** ไปยัง `main`  
8. `main` ( execution )  `defer`  
9. ถ้าไม่มี `recover()` ที่ไหน โปรแกรมจบพร้อมแสดง panic value และ stack trace  

### ถ้ามี `recover()` ใน `defer` ณ จุดใด

- panic value จะถูก **จับ** ไว้ที่จุดนั้น  
- unwind stack **หยุด** ที่ฟังก์ชันนั้น  
- การทำงานปกติจะดำเนินต่อหลังจาก `defer` ที่เรียก `recover()`  

---  

1. **ต้นทาง panic** เกิดขึ้นใน `level3()` – ค่า panic (เช่น `"error"`) ถูกสร้างขึ้น  
2. **เริ่มการคลายสแต็ก (unwind stack)** – การทำงานปกติใน `level3()` หยุดทันที; `defer` ใน `level3()` ถูกเรียก (ถ้ามี)  
3. **ค่า panic** ถูกส่งขึ้นไปยัง `level2()`  
4. **`level2()`**   `defer` ของตัวเอง (ตามลำดับ LIFO) จากนั้นถ้าไม่มี `recover` ค่า panic จะถูกส่งขึ้นไปต่อ  
5. **ค่า panic** ถูกส่งขึ้นไปยัง `level1()` แล้วก็ `main()`  
6. ใน `main()` (หรือฟังก์ชันใด ๆ) ถ้ามี `recover()` อยู่ใน `defer` จะสามารถ **จับ** ค่า panic ไว้ได้ และ **หยุดการคลายสแต็ก** การทำงานจะดำเนินต่อตามปกติ  

### สรุป
- ค่า panic เป็น **ข้อมูล** ที่ไหลขึ้นไปตาม call stack  
- ทุก `defer` สามารถตรวจสอบค่า panic ผ่าน `recover()`  
- ถ้า `defer` ใดเรียก `recover()` การไหลของ panic จะหยุดที่ระดับนั้น และโปรแกรมทำงานต่อ  

 

*ตัวอย่าง:*  
```go
fmt.Println("This will execute") // บรรทัดนี้จะถูก execute
```
### กระบวนการ unwind stack  (การคลายสแต็ก)
1. เกิด `panic` ที่จุดใดจุดหนึ่ง  
2. การ  (execute or execution) ของฟังก์ชันปัจจุบันหยุดทันที  
3. `defer` ทั้งหมดในฟังก์ชันปัจจุบันถูกเรียก (แบบ LIFO)  
4. ควบคุมกลับไปยังฟังก์ชันที่เรียก (caller)  
5. ทำซ้ำข้อ 2-4 จนถึง `main` หรือจนกว่าจะเจอ `recover`  
-   (execute or execution)  กระบวนการที่โปรเซสเซอร์หรือ runtime ของ Go ทำตามคำสั่งในโค้ด (เช่น การทำงานของฟังก์ชัน, การวนลูป, การคืนค่า)
  

**การคลายสแต็ก (stack unwinding)** คือกระบวนการที่เกิดขึ้นเมื่อเกิด `panic` ใน Go โดยโปรแกรมจะ **ย้อนกลับไปตามลำดับการเรียกฟังก์ชัน (call stack)** จากฟังก์ชันที่เกิด panic ขึ้นไปยังฟังก์ชันที่เรียกมันเรื่อยๆ จนถึง `main` (หรือจนกว่าจะเจอ `recover`)  

ระหว่างการคลายสแต็ก:
- การทำงานปกติของฟังก์ชันปัจจุบัน **หยุดทันที**
- คำสั่ง `defer` **ทั้งหมด** ในฟังก์ชันนั้นจะถูกเรียก (execute) ตามลำดับ LIFO (ประกาศทีหลังเรียกก่อน)
- จากนั้น panic value จะถูกส่งขึ้นไปยังฟังก์ชันที่เรียก (caller) และทำซ้ำขั้นตอนเดิม

หากไม่มี `recover()` ใดๆ จับ panic ไว้ กระบวนการจะดำเนินไปจนถึง `main` แล้วโปรแกรมจะหยุดพร้อมแสดง stack trace  
หากมี `recover()` ภายใน `defer` ของฟังก์ชันใด ณ จุดนั้น panic จะถูกจับไว้ และ **การคลายสแต็กจะหยุด** โปรแกรมจะทำงานต่อจากฟังก์ชันนั้นตามปกติ

### ตัวอย่างภาพจากแผนภาพ (แบบย่อ)
```
main() → level1() → level2() → level3()
                                   |
                              panic เกิด
                                   ↓
                         execute defer (level3)
                                   ↓
                         panic ส่งขึ้น → level2()
                                   ↓
                         execute defer (level2)
                                   ↓
                         panic ส่งขึ้น → level1()
                                   ↓
                         execute defer (level1)
                                   ↓
                         panic ส่งขึ้น → main()
                                   ↓
                         execute defer (main)
                                   ↓
                         ถ้าไม่มี recover → โปรแกรมหยุด
```
## `recover` คืออะไร

`recover` เป็น built-in function ในภาษา Go ที่ใช้ **กู้คืนโปรแกรมจากภาวะ panic** ช่วยให้โปรแกรมไม่ต้องหยุดทำงานทั้งระบบเมื่อเกิดข้อผิดพลาดร้ายแรง โดยสามารถจับค่า panic ที่เกิดขึ้นและทำให้โปรแกรมทำงานต่อไปได้ตามปกติ

---

### หลักการทำงาน

- `recover` จะทำงานได้ **เฉพาะเมื่อถูกเรียกภายใน `defer` function** เท่านั้น  
- หากเรียกนอก `defer` หรือไม่มี panic เกิดขึ้น `recover()` จะคืนค่า `nil`  
- เมื่อ panic เกิดขึ้นและมี `defer` ที่เรียก `recover()` ค่า panic จะถูกจับไว้ และ **การคลายสแต็ก (stack unwinding) จะหยุด** ที่ฟังก์ชันนั้น  
- `recover()` คืนค่าที่เป็น `interface{}` ซึ่งคือค่าที่ส่งเข้าไปใน `panic` (เช่น ข้อความ error, struct, ฯลฯ)

---

### ตัวอย่างการใช้งาน

```go
package main

import "fmt"

func safeDivision(a, b int) {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("กู้คืนจาก panic:", r)
            fmt.Println("โปรแกรมทำงานต่อ")
        }
    }()

    if b == 0 {
        panic("ตัวหารเป็นศูนย์")
    }
    fmt.Println("ผลลัพธ์:", a/b)
}

func main() {
    safeDivision(10, 2) // ทำงานปกติ → ผลลัพธ์: 5
    safeDivision(10, 0) // เกิด panic แต่ recover จับได้
    fmt.Println("จบโปรแกรม")
}
```

**ผลลัพธ์:**
```
ผลลัพธ์: 5
กู้คืนจาก panic: ตัวหารเป็นศูนย์
โปรแกรมทำงานต่อ
จบโปรแกรม
```

---

### ข้อควรระวัง

1. **ต้องอยู่ใน `defer` เท่านั้น**  
   - การเรียก `recover()` นอก `defer` จะไม่มีผลและคืนค่า `nil` เสมอ

2. **ใช้เฉพาะเมื่อจำเป็น**  
   - โดยทั่วไปควรจัดการ error ผ่าน `error` return value มากกว่าใช้ `panic`/`recover`  
   - `panic`/`recover` เหมาะสำหรับกรณีที่ไม่ควรเกิดขึ้น (programmer error, initialization failure ที่ไม่สามารถดำเนินต่อได้)

3. **`recover` คืนค่าเป็น `interface{}`**  
   - ต้องทำ type assertion หรือ type switch เพื่อใช้งานค่าที่ได้

---

### ความสัมพันธ์กับ `panic` และ `defer`

```
panic เกิดขึ้น → คลายสแต็ก (unwind stack) → execute defer → ถ้ามี recover() → จับค่า panic → หยุด unwind → โปรแกรมทำงานต่อ
```

**สรุป:** `recover` คือกลไกในการจับ panic ช่วยให้โปรแกรม Go มีความยืดหยุ่นในการจัดการกับข้อผิดพลาดร้ายแรงโดยไม่ต้องหยุดทำงานทั้งระบบ เหมาะสำหรับใช้ใน library หรือ middleware ที่ต้องการป้องกันการ crash ของโปรแกรมหลัก

---

 
**สรุป:** การคลายสแต็กคือการเดินกลับขึ้น call stack พร้อม deferred functions เพื่อให้โอกาสในการ cleanup หรือกู้คืน (recover) ก่อนที่โปรแกรมจะจบ
 
### ตัวอย่าง
```go
func level3() {
    defer fmt.Println("level3 defer")
    panic("error")
    fmt.Println("level3 end") // ไม่ถูก  (execute or execution) 
}

func level2() {
    defer fmt.Println("level2 defer")
    level3()
    fmt.Println("level2 end") // ไม่ถูก  (execute or execution) 
}

func level1() {
    defer fmt.Println("level1 defer")
    level2()
    fmt.Println("level1 end") // ไม่ถูก  (execute or execution) 
}

func main() {
    level1()
    fmt.Println("main end") // ไม่ถูก  (execute or execution) 
}
```
**ผลลัพธ์ (และลำดับ unwind):**
```
level3 defer
level2 defer
level1 defer
panic: error
...
```
แม้ `panic` จะเกิดขึ้นใน `level3` แต่ `defer` ใน `level2` และ `level1` ยังคงถูกเรียก เพราะโปรแกรม unwind stack ขึ้นไปเรื่อย ๆ

### สรุป
- **Unwind stack** คือการเดินกลับขึ้นไปตาม call stack พร้อม  (execute or execution)  deferred functions  
- เกิดจาก `panic` หรือ `runtime.Goexit()` (แต่ `Goexit` ไม่ unwind ถึง main)  
- `recover` สามารถหยุด unwind ได้หากถูกเรียกใน `defer` ของฟังก์ชันที่เกิด panic
```go
func readFile() {
    f, err := os.Open("file.txt")
    if err != nil {
        return
    }
    defer f.Close() // จะถูกเรียกเมื่อ readFile จบ
    // อ่านไฟล์...
}
```
defer จะทำงานในลำดับ LIFO (stack)

ฟังก์ชัน `readFile` ใช้ `os.Open` เพื่อเปิดไฟล์ หากเกิดข้อผิดพลาด (`err != nil`) จะ `return` ทันทีโดยไม่ดำเนินการต่อ หากเปิดสำเร็จ จะใช้ `defer f.Close()` เพื่อกำหนดให้ปิดไฟล์เมื่อฟังก์ชัน `readFile` ทำงานเสร็จ (ไม่ว่าจะจบแบบปกติหรือเกิด panic) panic คือ การหยุดหการทำงาน 

### หลักการทำงานของ `defer`
- คำสั่ง `defer` จะเลื่อนการทำงานของฟังก์ชันที่ตามมาให้เรียกเมื่อฟังก์ชันปัจจุบัน (ที่ประกาศ `defer`) จบการทำงาน
- ลำดับการเรียก `defer` เป็นแบบ LIFO (Last In, First Out) – ประกาศทีหลังจะถูกเรียกก่อน
- เหมาะสำหรับการปิดทรัพยากร (ไฟล์, connection) เพื่อป้องกันการรั่วไหล

----------------------------------
### ตัวอย่างการใช้งาน
----------------------------------
เพื่อให้เนื้อหาของคุณสมบูรณ์ยิ่งขึ้น ผมขอเสนอแผนภาพประกอบเพิ่มเติมในรูปแบบ **Mermaid** ซึ่งสามารถแทรกใน Markdown หรือเอกสารที่รองรับได้ ช่วยให้เห็นภาพการทำงานของ first-class functions, defer, และ panic/unwind stack ได้ชัดเจนขึ้น

---

## 1. แผนภาพ First-Class Functions

```mermaid
graph TD
    A[ฟังก์ชันเป็น first-class citizen] --> B[กำหนดให้ตัวแปร]
    A --> C[ส่งเป็นพารามิเตอร์]
    A --> D[คืนค่าเป็นผลลัพธ์]

    B --> B1["var fn func(int) int = func(x int) int { return x*2 }"]
    B1 --> B2["fn(5) → 10"]

    C --> C1["apply(fn func(int) int, val int) int { return fn(val) }"]
    C1 --> C2["apply(func(x int) int { return x+10 }, 7) → 17"]

    D --> D1["makeMultiplier(factor int) func(int) int"]
    D1 --> D2["double := makeMultiplier(2)"]
    D2 --> D3["double(5) → 10"]
```

---

## 2. แผนภาพการทำงานของ `defer` (LIFO)

```mermaid
sequenceDiagram
    participant main
    participant A as func A()
    participant B as func B()

    main->>A: เรียก A()
    activate A
    A->>A: defer fmt.Println("A: defer 1")
    A->>A: defer fmt.Println("A: defer 2")
    A->>B: เรียก B()
    activate B
    B->>B: defer fmt.Println("B: defer 1")
    B->>B: defer fmt.Println("B: defer 2")
    B-->>A: return
    deactivate B
    Note right of A: ลำดับการเรียก defer ใน B: <br/> B: defer 2 → B: defer 1
    A-->>main: return
    deactivate A
    Note right of main: ลำดับการเรียก defer ใน A: <br/> A: defer 2 → A: defer 1
```

---

## 3. แผนภาพ Panic, Recover และ Stack Unwinding

### 3.1 การ unwind stack เมื่อเกิด panic (ไม่มี recover)

```mermaid
sequenceDiagram
    participant main
    participant L1 as level1()
    participant L2 as level2()
    participant L3 as level3()

    main->>L1: call level1()
    activate L1
    L1->>L2: call level2()
    activate L2
    L2->>L3: call level3()
    activate L3
    L3--xL3: panic("error")
    Note over L3: execute defer ใน level3
    L3-->>L2: panic ส่งขึ้น (unwind)
    deactivate L3
    Note over L2: execute defer ใน level2
    L2-->>L1: panic ส่งขึ้น
    deactivate L2
    Note over L1: execute defer ใน level1
    L1-->>main: panic ส่งขึ้น
    deactivate L1
    Note over main: execute defer ใน main<br/>โปรแกรมหยุด พิมพ์ stack trace
```

### 3.2 การใช้ `recover` หยุดการ unwind

```mermaid
sequenceDiagram
    participant main
    participant L1 as level1()
    participant L2 as level2()
    participant L3 as level3()

    main->>L1: call level1()
    activate L1
    L1->>L2: call level2()
    activate L2
    L2->>L3: call level3()
    activate L3
    L3--xL3: panic("error")
    Note over L3: execute defer ใน level3
    L3-->>L2: panic ส่งขึ้น
    deactivate L3
    Note over L2: defer func() { recover() } จับ panic
    Note over L2: หยุด unwind ที่ level2
    L2-->>L1: return ปกติ (ไม่มี panic)
    deactivate L2
    L1-->>main: return ปกติ
    deactivate L1
    Note over main: โปรแกรมทำงานต่อ
```
 ----------------
 **การคลายสแต็ก (stack unwinding)** คือกระบวนการที่เกิดขึ้นเมื่อเกิด `panic` ใน Go โดยโปรแกรมจะ **ย้อนกลับไปตามลำดับการเรียกฟังก์ชัน (call stack)** จากฟังก์ชันที่เกิด panic ขึ้นไปยังฟังก์ชันที่เรียกมันเรื่อยๆ จนถึง `main` (หรือจนกว่าจะเจอ `recover`)  

ระหว่างการคลายสแต็ก:
- การทำงานปกติของฟังก์ชันปัจจุบัน **หยุดทันที**
- คำสั่ง `defer` **ทั้งหมด** ในฟังก์ชันนั้นจะถูกเรียก (execute) ตามลำดับ LIFO (ประกาศทีหลังเรียกก่อน)
- จากนั้น panic value จะถูกส่งขึ้นไปยังฟังก์ชันที่เรียก (caller) และทำซ้ำขั้นตอนเดิม

หากไม่มี `recover()` ใดๆ จับ panic ไว้ กระบวนการจะดำเนินไปจนถึง `main` แล้วโปรแกรมจะหยุดพร้อมแสดง stack trace  
หากมี `recover()` ภายใน `defer` ของฟังก์ชันใด ณ จุดนั้น panic จะถูกจับไว้ และ **การคลายสแต็กจะหยุด** โปรแกรมจะทำงานต่อจากฟังก์ชันนั้นตามปกติ

### ตัวอย่างภาพจากแผนภาพ (แบบย่อ)
```
main() → level1() → level2() → level3()
                                   |
                              panic เกิด
                                   ↓
                         execute defer (level3)
                                   ↓
                         panic ส่งขึ้น → level2()
                                   ↓
                         execute defer (level2)
                                   ↓
                         panic ส่งขึ้น → level1()
                                   ↓
                         execute defer (level1)
                                   ↓
                         panic ส่งขึ้น → main()
                                   ↓
                         execute defer (main)
                                   ↓
                         ถ้าไม่มี recover → โปรแกรมหยุด
```

**สรุป:** การคลายสแต็กคือการเดินกลับขึ้น call stack พร้อม deferred functions เพื่อให้โอกาสในการ cleanup หรือกู้คืน (recover) ก่อนที่โปรแกรมจะจบ
 
```go
package main

import (
    "fmt"
    "os"
)

func readFile() {
    f, err := os.Open("file.txt")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer f.Close() // รับประกันว่าจะปิดไฟล์เมื่อฟังก์ชันจบ

    // อ่านข้อมูลจากไฟล์...
    data := make([]byte, 100)
    n, _ := f.Read(data)
    fmt.Printf("Read %d bytes: %s\n", n, string(data[:n]))
}

func main() {
    readFile()
}
```

### สรุป
`defer` ช่วยให้โค้ดสะอาดและปลอดภัย เพราะรับประกันการทำความสะอาด (เช่น ปิดไฟล์) แม้ในกรณีที่มีการ `return` ก่อนถึงคำสั่งปิดหรือเกิด panic

----------------------
# recover
## `recover` คืออะไร

`recover` เป็น built-in function ในภาษา Go ที่ใช้ **กู้คืนโปรแกรมจากภาวะ panic** ช่วยให้โปรแกรมไม่ต้องหยุดทำงานทั้งระบบเมื่อเกิดข้อผิดพลาดร้ายแรง โดยสามารถจับค่า panic ที่เกิดขึ้นและทำให้โปรแกรมทำงานต่อไปได้ตามปกติ

---

### หลักการทำงาน

- `recover` จะทำงานได้ **เฉพาะเมื่อถูกเรียกภายใน `defer` function** เท่านั้น  
- หากเรียกนอก `defer` หรือไม่มี panic เกิดขึ้น `recover()` จะคืนค่า `nil`  
- เมื่อ panic เกิดขึ้นและมี `defer` ที่เรียก `recover()` ค่า panic จะถูกจับไว้ และ **การคลายสแต็ก (stack unwinding) จะหยุด** ที่ฟังก์ชันนั้น  
- `recover()` คืนค่าที่เป็น `interface{}` ซึ่งคือค่าที่ส่งเข้าไปใน `panic` (เช่น ข้อความ error, struct, ฯลฯ)

---

### ตัวอย่างการใช้งาน

```go
package main

import "fmt"

func safeDivision(a, b int) {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("กู้คืนจาก panic:", r)
            fmt.Println("โปรแกรมทำงานต่อ")
        }
    }()

    if b == 0 {
        panic("ตัวหารเป็นศูนย์")
    }
    fmt.Println("ผลลัพธ์:", a/b)
}

func main() {
    safeDivision(10, 2) // ทำงานปกติ → ผลลัพธ์: 5
    safeDivision(10, 0) // เกิด panic แต่ recover จับได้
    fmt.Println("จบโปรแกรม")
}
```

**ผลลัพธ์:**
```
ผลลัพธ์: 5
กู้คืนจาก panic: ตัวหารเป็นศูนย์
โปรแกรมทำงานต่อ
จบโปรแกรม
```

---

### ข้อควรระวัง

1. **ต้องอยู่ใน `defer` เท่านั้น**  
   - การเรียก `recover()` นอก `defer` จะไม่มีผลและคืนค่า `nil`(ค่าว่าง) เสมอ

2. **ใช้เฉพาะเมื่อจำเป็น**  
   - โดยทั่วไปควรจัดการ error ผ่าน `error` return value มากกว่าใช้ `panic`/`recover`  
   - `panic`/`recover` เหมาะสำหรับกรณีที่ไม่ควรเกิดขึ้น (programmer error, initialization failure ที่ไม่สามารถดำเนินต่อได้)

3. **`recover` คืนค่าเป็น `interface{}`**  
   - ต้องทำ type assertion หรือ type switch เพื่อใช้งานค่าที่ได้

---

### ความสัมพันธ์กับ `panic` และ `defer`

```
panic เกิดขึ้น → คลายสแต็ก (unwind stack) → execute defer → ถ้ามี recover() → จับค่า panic → หยุด unwind → โปรแกรมทำงานต่อ
```

**สรุป:** `recover` คือกลไกในการจับ panic ช่วยให้โปรแกรม Go มีความยืดหยุ่นในการจัดการกับข้อผิดพลาดร้ายแรงโดยไม่ต้องหยุดทำงานทั้งระบบ เหมาะสำหรับใช้ใน library หรือ middleware ที่ต้องการป้องกันการ crash ของโปรแกรมหลัก

---
 

เพื่อให้เห็นภาพการทำงานของ **`recover`** อย่างชัดเจน ผมขอเสนอแผนภาพ 3 แบบ ดังนี้

---

## 1. แผนภาพลำดับ (Sequence Diagram) – เปรียบเทียบแบบมี/ไม่มี `recover`

### กรณี **ไม่มี recover** (โปรแกรมหยุด)

```mermaid
sequenceDiagram
    participant main
    participant L1 as level1()
    participant L2 as level2()
    participant L3 as level3()

    main->>L1: call level1()
    activate L1
    L1->>L2: call level2()
    activate L2
    L2->>L3: call level3()
    activate L3
    L3-->>L3: panic("error")
    Note over L3: execute defer (L3)
    L3-->>L2: panic ส่งขึ้น
    deactivate L3
    Note over L2: execute defer (L2)
    L2-->>L1: panic ส่งขึ้น
    deactivate L2
    Note over L1: execute defer (L1)
    L1-->>main: panic ส่งขึ้น
    deactivate L1
    Note over main: execute defer (main)<br/>โปรแกรมหยุด
```

### กรณี **มี recover** (โปรแกรมทำงานต่อ)

```mermaid
sequenceDiagram
    participant main
    participant L1 as level1()
    participant L2 as level2()
    participant L3 as level3()

    main->>L1: call level1()
    activate L1
    L1->>L2: call level2()
    activate L2
    L2->>L3: call level3()
    activate L3
    L3-->>L3: panic("error")
    Note over L3: execute defer (L3)
    L3-->>L2: panic ส่งขึ้น
    deactivate L3
    Note over L2: defer มี recover()<br/>จับ panic และหยุด unwind
    L2-->>L1: return ปกติ (ไม่มี panic)
    deactivate L2
    L1-->>main: return ปกติ
    deactivate L1
    Note over main: โปรแกรมทำงานต่อ
```

---

## 2. แผนภาพการไหล (Flowchart) – กระบวนการตัดสินใจของ `recover`
 
---

## 1. Mermaid Flowchart (ปรับปรุง)

```mermaid
graph TD
    Start([เกิด panic]) --> A{ฟังก์ชันนี้มี defer หรือไม่?}
    
    A -->|ไม่มี defer| B[ส่ง panic ไปยัง caller]
    A -->|มี defer| C[execute defer ทั้งหมด<br/>ตามลำดับ LIFO]
    
    C --> D{ใน defer มี recover()?}
    D -->|มี| E[recover จับ panic value]
    E --> F[หยุดการคลายสแต็กที่ฟังก์ชันนี้]
    F --> G[ฟังก์ชัน return ปกติ]
    G --> H([โปรแกรมทำงานต่อ])
    
    D -->|ไม่มี recover| B
    
    B --> I{caller คือ main?}
    I -->|ใช่| J[โปรแกรมหยุด<br/>พิมพ์ stack trace]
    I -->|ไม่ใช่| K[ย้ายไปยัง caller]
    K --> A
    
    style E fill:#9f9,stroke:#333
    style J fill:#f99,stroke:#333
    style H fill:#9cf,stroke:#333
```

**คำอธิบายเพิ่มเติม:**
- หลังจาก `recover` จับ panic ได้ การคลายสแต็กจะหยุดทันที และฟังก์ชันจะ return ตามปกติโดยไม่ส่ง panic ต่อ
- หากไม่มี `recover` ใน defer ใด ๆ panic จะถูกส่งขึ้นไปยัง caller ซ้ำจนกว่าจะถึง `main` หรือพบ `recover`

---

## 2. ASCII Flowchart (แบบข้อความ)

```
                     ┌─────────────────────┐
                     │     เกิด panic      │
                     └──────────┬──────────┘
                                │
                                ▼
                     ┌─────────────────────┐
                     │ มี defer ในฟังก์ชัน  │
                     │   ปัจจุบันหรือไม่?   │
                     └──────────┬──────────┘
                                │
            ┌───────────────────┴───────────────────┐
            │ไม่มี                                  │มี
            ▼                                       ▼
┌───────────────────────┐            ┌───────────────────────────┐
│ ส่ง panic ไปยัง caller │            │ execute defer ทั้งหมด     │
│ (เริ่ม unwind stack)   │            │ ตามลำดับ LIFO             │
└───────────┬───────────┘            └─────────────┬─────────────┘
            │                                      │
            │                                      ▼
            │                         ┌───────────────────────────┐
            │                         │ ใน defer มี recover()?    │
            │                         └─────────────┬─────────────┘
            │                                      │
            │                 ┌────────────────────┴────────────────────┐
            │                 │ไม่มี                                  │มี
            │                 ▼                                       ▼
            │    ┌────────────────────────┐            ┌───────────────────────────┐
            │    │ ส่ง panic ไปยัง caller │            │ recover จับ panic value    │
            │    └───────────┬────────────┘            └─────────────┬─────────────┘
            │                │                                       │
            │                │                                       ▼
            │                │                         ┌───────────────────────────┐
            │                │                         │ หยุดการคลายสแต็ก         │
            │                │                         │ ที่ฟังก์ชันนี้             │
            │                │                         └─────────────┬─────────────┘
            │                │                                       │
            │                │                                       ▼
            │                │                         ┌───────────────────────────┐
            │                │                         │ ฟังก์ชัน return ปกติ      │
            │                │                         └─────────────┬─────────────┘
            │                │                                       │
            │                │                                       ▼
            │                │                         ┌───────────────────────────┐
            │                │                         │    โปรแกรมทำงานต่อ        │
            │                │                         └───────────────────────────┘
            │                │
            ▼                ▼
┌─────────────────────────────────────┐
│      caller คือ main หรือไม่?        │
└─────────────────┬───────────────────┘
                  │
        ┌─────────┴─────────┐
        │ใช่                │ไม่ใช่
        ▼                   ▼
┌─────────────────┐   ┌─────────────────┐
│ โปรแกรมหยุด    │   │ ย้ายไปยัง caller │
│ พิมพ์ stack trace│   │ (กลับไปถามว่า   │
└─────────────────┘   │ มี defer อีกไหม)│
                      └────────┬────────┘
                               │
                               └──────→ (กลับไปที่จุดเริ่มต้น)
```
  
แผนภาพเหล่านี้ช่วยให้เห็นภาพรวมของ **`recover`** ได้ชัดเจนยิ่งขึ้น โดยเฉพาะบทบาทในการ **หยุดการคลายสแต็ก** และทำให้โปรแกรมดำเนินต่อไปได้อย่างปลอดภัย
---------------------
 
#### 10.7 panic และ recover
- `panic` : หยุดการทำงานปกติและเริ่ม unwind stack
- `recover` : ใช้ใน defer เพื่อจับ panic และควบคุมการทำงานต่อ

```go
func safeDivide(a, b int) {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered from panic:", r)
        }
    }()
    if b == 0 {
        panic("division by zero")
    }
    fmt.Println(a / b)
}
```
**ข้อควรระวัง**: panic ควรใช้ในกรณีผิดปกติรุนแรงเท่านั้น ไม่ใช่แทน error handling

#### 10.8 ฟังก์ชัน init
แต่ละแพคเกจ ได้รับ `init()` ฟังก์ชัน ซึ่งจะถูกเรียกอัตโนมัติเมื่อแพคเกจถูกโหลด (ก่อน main)
```go
func init() {
    fmt.Println("initializing package")
}
```
ในการทำธุรกรรมที่เกี่ยวข้องกับการสั่งซื้อสินค้า การตัดสต็อก และการออกใบเสร็จ เราต้องการให้ข้อมูลทุกส่วนถูกบันทึกอย่างสมบูรณ์หรือไม่บันทึกเลย หากเกิดข้อผิดพลาดขึ้นที่ขั้นตอนใดขั้นตอนหนึ่ง ระบบต้องสามารถ **ย้อนกลับ (rollback)** การเปลี่ยนแปลงทั้งหมดเพื่อรักษาความถูกต้องของข้อมูล

## การใช้ GORM Transaction เพื่อ Rollback

GORM มีฟังก์ชัน `db.Transaction` ที่ช่วยให้เราสามารถรวมหลายคำสั่ง SQL ไว้ใน transaction เดียวกันได้ โดยหากฟังก์ชันที่ส่งเข้าไปคืนค่า `error` GORM จะทำการ rollback โดยอัตโนมัติ ถ้าคืน `nil` จะ commit

### โครงสร้างตารางตัวอย่าง
```go
type Order struct {
    ID        uint
    UserID    uint
    Total     float64
    CreatedAt time.Time
}

type Stock struct {
    ProductID uint `gorm:"primaryKey"`
    Quantity  int
}

type Receipt struct {
    ID        uint
    OrderID   uint
    Amount    float64
    IssuedAt  time.Time
}
```

### ฟังก์ชัน PlaceOrder แบบ Transaction
```go
func PlaceOrder(db *gorm.DB, userID uint, items []CartItem) error {
    // เริ่ม transaction
    return db.Transaction(func(tx *gorm.DB) error {
        // 1. คำนวณราคารวม และตรวจสอบสต็อกพร้อม lock
        var total float64
        for _, item := range items {
            var stock Stock
            // Lock แถว stock เพื่อป้องกัน race condition
            if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
                Where("product_id = ?", item.ProductID).
                First(&stock).Error; err != nil {
                return err // สินค้าไม่มีในระบบ
            }
            if stock.Quantity < item.Quantity {
                return errors.New("สินค้าไม่พอ")
            }
            // หักสต็อก (จะบันทึกภายหลัง)
            stock.Quantity -= item.Quantity
            if err := tx.Save(&stock).Error; err != nil {
                return err
            }
            // คำนวณราคารวม (สมมุติมีฟังก์ชัน getPrice)
            total += getPrice(item.ProductID) * float64(item.Quantity)
        }

        // 2. สร้าง order
        order := Order{UserID: userID, Total: total}
        if err := tx.Create(&order).Error; err != nil {
            return err
        }

        // 3. สร้าง receipt
        receipt := Receipt{OrderID: order.ID, Amount: total}
        if err := tx.Create(&receipt).Error; err != nil {
            return err
        }

        // ทุกอย่างสำเร็จ -> commit อัตโนมัติ
        return nil
    })
}
```

### กระบวนการทำงานเมื่อเกิดข้อผิดพลาด
- หากขั้นตอนใด (เช่น การ lock stock หรือการหักสต็อก หรือการ insert order/receipt) คืน error กลับมา ฟังก์ชันที่ส่งให้ `Transaction` จะคืน error นั้น
- GORM จะทำการ rollback ทุกคำสั่งที่ได้ดำเนินการไปแล้วใน transaction เดียวกัน (เช่น order ที่สร้างไปแล้วจะถูกลบ, stock ที่หักไปแล้วจะถูกคืนค่า)
- โปรแกรมภายนอกจะได้รับ error และสามารถแจ้งผู้ใช้ว่าเกิดข้อผิดพลาด โดยไม่มีข้อมูลค้างในฐานข้อมูล

### การป้องกัน Race Condition ด้วย Lock
- `Clauses(clause.Locking{Strength: "UPDATE"})` จะเพิ่ม `FOR UPDATE` ใน SQL เพื่อ lock แถว stock ขณะที่เราอ่านค่า
- เมื่อ transaction ยังไม่ commit แถวที่ lock จะไม่ให้ transaction อื่นอ่านหรือเขียนได้ (ขึ้นอยู่กับ isolation level) ทำให้การหักสต็อกปลอดภัยจากการทำงานพร้อมกัน

### ข้อควรระวัง
- **Transaction ควรสั้นที่สุด** หลีกเลี่ยงการทำงานที่ใช้เวลานานภายใน transaction (เช่น การเรียก API ภายนอก) เพราะจะผูก lock ไว้นาน
- **จัดการ error อย่างเหมาะสม** หากเกิด error ควรแจ้งให้ผู้ใช้ทราบถึงสาเหตุ (เช่น สินค้าไม่พอ)
- **เลือกใช้ isolation level** หากต้องการปรับระดับความเข้มงวด สามารถตั้งค่าได้ผ่าน `tx.Exec("SET TRANSACTION ISOLATION LEVEL ...")` ก่อนเริ่ม transaction

## สรุป
การใช้ GORM transaction ช่วยให้เราสามารถทำ rollback ได้โดยอัตโนมัติเมื่อเกิดความผิดพลาด การเพิ่ม lock ที่แถว stock ช่วยให้ข้อมูลสอดคล้องแม้มีผู้ใช้หลายคนสั่งซื้อพร้อมกัน วิธีนี้ทำให้การทำงานที่ต้องอาศัยความถูกต้องของข้อมูลหลายส่วน (order, stock, receipt) มีความปลอดภัยและเชื่อถือได้
---

### บทที่ 11: แพคเกจและการนำเข้า

#### 11.1 แพคเกจ (Packages)
Go จัดระเบียบโค้ดเป็นแพคเกจ แต่ละไฟล์ .go ต้องขึ้นต้นด้วย `package <name>` ชื่อแพคเกจควรเป็นตัวพิมพ์เล็ก

- `package main` : สำหรับโปรแกรมที่รันได้ (executable)
- แพคเกจอื่นๆ : สำหรับ library

#### 11.2 การนำเข้า (import)
```go
import "fmt"
import "math/rand"
```
หรือแบบ grouped:
```go
import (
    "fmt"
    "math/rand"
)
```

#### 11.3 การตั้งชื่อให้กับ import (alias)
```go
import (
    "fmt"
    r "math/rand"   // alias
)
```

#### 11.4 Blank import
ใช้ `_` เพื่อนำเข้าเฉพาะ side-effect (เช่น เรียก init) โดยไม่ต้องใช้ฟังก์ชัน
```go
import _ "image/png" // ลงทะเบียนตัวถอดรหัส PNG
```

#### 11.5 การเข้าถึงสมาชิก (exported vs unexported)
- ตัวพิมพ์ใหญ่: exported (public) สามารถเข้าถึงได้จากแพคเกจอื่น
- ตัวพิมพ์เล็ก: unexported (private) เข้าถึงได้เฉพาะภายในแพคเกจ

```go
// mymath/mymath.go
package mymath

func Add(a, b int) int { // Exported
    return a + b
}

func sub(a, b int) int { // Unexported
    return a - b
}
```

#### 11.6 การจัดโครงสร้างโปรเจกต์ด้วยแพคเกจ
```
myproject/
├── go.mod
├── main.go
├── utils/
│   └── stringutils.go
└── models/
    └── user.go
```
ใน `main.go`:
```go
import (
    "myproject/utils"
    "myproject/models"
)
```

#### 11.7 การจัดระเบียบแพคเกจมาตรฐาน
- `fmt` : การจัดรูปแบบ I/O
- `io` / `os` : การอ่านเขียนไฟล์
- `net/http` : HTTP client/server
- `encoding/json` : JSON
- `sync` : concurrency primitives
- `time` : เวลา
- `context` : context management

#### 11.8 การสร้างแพคเกจของตัวเอง
สร้างไดเรกทอรีใหม่ภายในโมดูล เขียนไฟล์ .go ด้วย package name ที่ตรงกับชื่อโฟลเดอร์

---

### บทที่ 12: การเริ่มต้นทำงานของแพคเกจ

#### 12.1 ลำดับการเริ่มต้น
1. แพคเกจที่ถูก import จะถูกเริ่มต้นก่อน (แบบ recursive)
2. ตัวแปรระดับแพคเกจถูกกำหนดค่า (initialized)
3. ฟังก์ชัน `init()` ถูกเรียก (เรียงตามลำดับการประกาศในไฟล์)
4. ถ้าเป็น `package main` ฟังก์ชัน `main()` จะถูกเรียก

#### 12.2 ฟังก์ชัน init
สามารถมีได้หลาย init ในแพคเกจเดียวกัน หรือแม้แต่หลาย init ในไฟล์เดียวกัน
```go
var config string

func init() {
    config = loadConfig()
    fmt.Println("init called")
}
```

#### 12.3 การควบคุมลำดับ init
ถ้าต้องการให้ init ทำงานก่อน ควรนำเข้าแพคเกจนั้นโดย blank import

#### 12.4 ตัวอย่างการเริ่มต้น
ไฟล์ `a.go`:
```go
package main

import "fmt"

var a = b + 1
var b = 1

func init() {
    fmt.Println("init 1", a, b)
}

func init() {
    fmt.Println("init 2")
}

func main() {
    fmt.Println("main", a, b)
}
```
ผลลัพธ์:
```
init 1 2 1
init 2
main 2 1
```
ตัวแปรจะถูกกำหนดค่าตามลำดับการประกาศในไฟล์

#### 12.5 การใช้ init เพื่อตั้งค่าแพคเกจ
เหมาะสำหรับการตั้งค่า connection pools, การลงทะเบียน driver, การโหลด configuration
```go
package database

import "database/sql"

var db *sql.DB

func init() {
    // เปิด connection pool
    db, _ = sql.Open("mysql", "user:pass@/dbname")
}
```

#### 12.6 การหลีกเลี่ยงการใช้งาน init มากเกินไป
init ทำให้ยากต่อการทดสอบและการควบคุมลำดับ หากเป็นไปได้ควรใช้ฟังก์ชันเริ่มต้นที่เรียกอย่างชัดเจน (เช่น New) แทน

---

### บทที่ 13: การสร้างชนิดข้อมูลใหม่ (Types)

#### 13.1 type keyword
ใช้สร้างชนิดข้อมูลใหม่จากชนิดเดิม
```go
type Celsius float64
type Fahrenheit float64

func (c Celsius) ToFahrenheit() Fahrenheit {
    return Fahrenheit(c*9/5 + 32)
}
```
`Celsius` และ `Fahrenheit` เป็นคนละชนิดกัน แม้ underlying type เหมือนกัน ต้องแปลงอย่างชัดเจน

#### 13.2 struct
struct คือการรวมตัวแปรหลายๆ ตัวเข้าด้วยกัน
```go
type Person struct {
    Name string
    Age  int
    Email string
}
```
การสร้าง instance:
```go
p1 := Person{"John", 30, "john@example.com"} // ตามลำดับ field
p2 := Person{Name: "Jane", Age: 25}          // ระบุชื่อ field
p3 := new(Person) // สร้าง pointer
p3.Name = "Bob"
```

#### 13.3 struct fields และ tags
Tags ใช้สำหรับ metadata เช่น การ encode/decode JSON
```go
type User struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Password string `json:"-"` // ไม่แสดงใน JSON
}
```

#### 13.4 การสืบทอด (embedding) แทน inheritance
Go ไม่มี inheritance แต่ใช้ composition ผ่าน embedding
```go
type Animal struct {
    Name string
}

func (a Animal) Speak() string {
    return "???"
}

type Dog struct {
    Animal   // embedded field
    Breed string
}

func (d Dog) Speak() string {
    return "Woof!"
}
```
Dog จะมี field Name และ method Speak (สามารถ override ได้)

#### 13.5 ชนิดข้อมูลแบบ type alias
Alias (Go 1.9+) ทำให้มีชื่อใหม่ที่เทียบเท่ากับชนิดเดิม ใช้ในกรณี refactoring
```go
type MyInt = int // alias, MyInt และ int ใช้แทนกันได้
```
ความแตกต่าง: type definition (type MyInt int) สร้างชนิดใหม่, alias แค่ชื่ออื่น

#### 13.6 ชนิดข้อมูลแบบ interface
จะกล่าวในบทที่ 16

---

### บทที่ 14: เมธอด (Methods)

#### 14.1 การกำหนดเมธอด
เมธอดคือฟังก์ชันที่มี receiver (ตัวรับ)
```go
type Rectangle struct {
    Width, Height float64
}

func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func (r *Rectangle) Scale(f float64) {
    r.Width *= f
    r.Height *= f
}
```

#### 14.2 Pointer receiver vs Value receiver
- **Value receiver**: ทำงานกับ copy เปลี่ยนแปลงไม่กระทบตัวเดิม
- **Pointer receiver**: ทำงานกับตัวแปรเดิม เปลี่ยนแปลงได้ และประหยัดหน่วยความจำสำหรับ struct ใหญ่

กฎ: ถ้าเมธอดต้องแก้ไข receiver ให้ใช้ pointer receiver

#### 14.3 เมธอดกับ non-struct
สามารถกำหนดเมธอดให้กับชนิดใดก็ได้ (ยกเว้นชนิดจากแพคเกจอื่นและชนิดพื้นฐานที่ไม่ใช่ alias)
```go
type MyInt int

func (m MyInt) Double() MyInt {
    return m * 2
}
```

#### 14.4 การเรียกเมธอด
```go
r := Rectangle{10, 5}
area := r.Area()      // value receiver
r.Scale(2)            // pointer receiver
```

#### 14.5 การแปลงระหว่าง value และ pointer อัตโนมัติ
Go จะแปลงให้อัตโนมัติเมื่อเรียกเมธอด:
- ถ้ามี pointer receiver แต่เรียกจาก value: Go จะส่ง address ให้ (ถ้า value addressable)
- ถ้ามี value receiver แต่เรียกจาก pointer: Go จะ dereference ให้

#### 14.6 เมธอดและการสืบทอด (embedding)
เมื่อ embed struct ไปอีก struct เมธอดของ struct ที่ถูก embed จะถูกเลื่อนขึ้นมา
```go
type Animal struct{}

func (a Animal) Speak() string {
    return "..."
}

type Dog struct {
    Animal
}

d := Dog{}
d.Speak() // เรียก Speak ของ Animal
```

#### 14.7 เมธอดและ interface
เมธอดเป็นสิ่งที่ทำให้ type implement interface ได้

---

### บทที่ 15: พอยน์เตอร์ (Pointer)

#### 15.1 พอยน์เตอร์คืออะไร?
พอยน์เตอร์คือตัวแปรที่เก็บ address ของตัวแปรอื่น
```go
var x int = 10
var p *int = &x   // p ชี้ไปที่ x
fmt.Println(*p)   // 10 (dereference)
*p = 20           // เปลี่ยนค่า x ผ่าน pointer
```

#### 15.2 การประกาศ pointer
- `&` : address-of operator
- `*` : dereference operator (และใช้ใน type declaration)

#### 15.3 Zero value ของ pointer
pointer ที่ไม่ได้กำหนดค่าจะเป็น `nil`

#### 15.4 การใช้ pointer กับฟังก์ชัน
เพื่อให้ฟังก์ชันสามารถแก้ไขตัวแปรนอกฟังก์ชันได้
```go
func zeroVal(val int) {
    val = 0
}

func zeroPtr(ptr *int) {
    *ptr = 0
}

x := 5
zeroVal(x)      // x ยังเป็น 5
zeroPtr(&x)     // x กลายเป็น 0
```

#### 15.5 การใช้ pointer กับ struct
เพื่อประสิทธิภาพ (ไม่ต้อง copy struct ทั้งตัว) และการแก้ไข
```go
func (p *Person) UpdateAge(age int) {
    p.Age = age
}
```
#### 15.6 new() และ make()
- `new(T)` : จัดสรรหน่วยความจำสำหรับ type T, คืน pointer *T, zero value
- `make(T, ...)` : ใช้กับ slice, map, channel เท่านั้น; คืน initialized (ไม่ใช่ zero) value

#### 15.7 Pointer arithmetic
Go ไม่สนับสนุน pointer arithmetic (ต่างจาก C) เพื่อความปลอดภัย

#### 15.8 ข้อควรระวัง
- nil pointer dereference จะทำให้ panic
- ควรใช้ pointer เมื่อจำเป็นเท่านั้น (การแก้ไข, ประสิทธิภาพ)

---
**nil pointer dereference** (การอ้างถึง pointer ที่เป็น nil) คือความผิดพลาดที่เกิดขึ้นเมื่อโปรแกรมพยายามเข้าถึงข้อมูลผ่าน pointer (ตัวแปรที่ชี้ไปยังตำแหน่งหน่วยความจำ) แต่ pointer นั้นไม่ได้ชี้ไปยังออบเจ็กต์หรือพื้นที่หน่วยความจำที่ถูกต้อง – มันมีค่าเป็น `nil` (null, ศูนย์) การพยายามอ่านหรือเขียนผ่าน pointer ดังกล่าวจึงไม่มีความหมายและทำให้โปรแกรมหยุดทำงานทันที

---

### nil pointer dereference คือ?.
ใน Go การ dereference nil pointer จะทำให้เกิด **panic** แบบ runtime:
```go
var p *int // p เป็น nil pointer
fmt.Println(*p) // panic: runtime error: invalid memory address or nil pointer dereference
```
หรือเมื่อเรียก method บน receiver ที่เป็น nil:
```go
type MyStruct struct { Name string }
func (m *MyStruct) Say() { fmt.Println(m.Name) }

var s *MyStruct // s เป็น nil
s.Say() // panic: nil pointer dereference (ถ้า method อ่านค่า m.Name)
```

### เกิดในภาษาอื่น
- **C/C++**: segmentation fault (core dumped)
- **Java**: NullPointerException
- **Python**: AttributeError หรือ TypeError (ขึ้นกับบริบท)
- **JavaScript**: TypeError: Cannot read property 'x' of null

---

### สาเหตุทั่วไป
- ประกาศตัวแปร pointer แต่ไม่ initialize (`var p *int`)
- คืนค่า nil จากฟังก์ชัน (เช่น `return nil, err`) แล้วนำมาใช้โดยไม่ตรวจสอบ
- ลบ element จากโครงสร้างข้อมูลแล้ว pointer เก่ายังถูกนำไปใช้
- รับค่าจาก map หรือ type assertion โดยไม่ตรวจสอบ ok

---

### ผลกระทบ
- โปรแกรม crash ทันที (หรือ panic)
- สูญเสียข้อมูลที่ยังไม่ถูกบันทึก
- อาจเป็นช่องโหว่ด้านความปลอดภัย (denial of service)

---

### วิธีป้องกันและแก้ไข

#### 1. ตรวจสอบ nil ก่อนใช้งาน
```go
if p != nil {
    fmt.Println(*p)
} else {
    // จัดการกรณี nil
}
```

#### 2. ใช้ zero value หรือ default object แทน pointer
```go
type Config struct {
    Timeout int
}
// ใช้ค่าเริ่มต้นแทน nil
config := Config{Timeout: 30}
```

#### 3. สำหรับ map, slice, channel – ตรวจสอบการ初始化
```go
var m map[string]int
m["key"] = 1 // panic: assignment to entry in nil map
// ควรใช้ make
m := make(map[string]int)
```

#### 4. ใช้ comma ok idiom สำหรับ type assertion หรือ map access
```go
if value, ok := myMap["key"]; ok {
    // ใช้ value
}
```

#### 5. ใช้ `recover()` กับ defer เพื่อดักจับ panic (ใน Go)
```go
defer func() {
    if r := recover(); r != nil {
        log.Println("recovered from:", r)
    }
}()
```

#### 6. ใช้เครื่องมือ static analysis
- `go vet` ช่วยตรวจสอบ potential nil dereference บางกรณี
- `nilaway` (จาก Uber) ตรวจจับ nil pointer bug ได้ละเอียดขึ้น

---

### สรุป
**nil pointer dereference** คือการพยายามเข้าถึงข้อมูลผ่าน pointer ที่ไม่มีค่า (nil) ซึ่งเป็นสาเหตุหลักของ panic และ program crash การป้องกันที่ดีที่สุดคือการตรวจสอบ nil ก่อนใช้งาน และออกแบบโค้ดให้ลดการพึ่งพา pointer ที่อาจเป็น nil ตั้งแต่แรก

คำว่า **Panic** ในบริบทของโปรแกรมหรือระบบคอมพิวเตอร์หมายถึง **สถานการณ์ที่โปรแกรมหรือระบบพบข้อผิดพลาดร้ายแรงจนไม่สามารถทำงานต่อไปได้อย่างปลอดภัย จึงต้องหยุดการทำงานทันที (Crash) เพื่อป้องกันความเสียหายหรือข้อมูลเสียหาย**

หากคุณกำลังหมายถึง **Panic ในภาษา Go (Golang)** ซึ่งเป็นภาษาที่ใช้คำนี้บ่อยที่สุด คำตอบมีดังนี้

---

### 1. Panic คืออะไร (ใน Go)

ในภาษา Go `panic` คือ built-in function ที่ทำให้โปรแกรมหยุดการทำงานปกติทันที และเริ่ม unwind the stack (คลายกองซ้อนฟังก์ชัน) พร้อมกับรัน `defer` statements ทั้งหมด

**ลักษณะของ Panic:**
- เกิดจากโปรแกรมเมอร์เรียก `panic()` โดยตรง
- เกิดจาก Runtime error เช่น การเข้าถึง slice index นอกขอบเขต (`index out of range`), การเรียก method บน nil pointer, การส่งค่า nil ไปยัง map

**ตัวอย่าง:**
```go
package main

func main() {
    defer func() {
        fmt.Println("defer ทำงานก่อน panic")
    }()
    
    panic("เกิดข้อผิดพลาดร้ายแรง")
    fmt.Println("บรรทัดนี้จะไม่ถูกทำงาน")
}
```

เมื่อ panic เกิดขึ้น โปรแกรมจะหยุดและแสดง stack trace เพื่อให้开发者追踪ต้นตอของปัญหา

---

### 2. แก้ไขอย่างไร

การแก้ไข panic มีหลักการหลักๆ ดังนี้:

#### 2.1 ใช้ `recover()` เพื่อดักจับ Panic
ใน Go มีฟังก์ชัน `recover()` ที่ใช้สำหรับดักจับ panic และทำให้โปรแกรมกลับมาทำงานต่อได้ (คล้าย try-catch ในภาษาอื่น)

**วิธีใช้:**
```go
package main

import "fmt"

func riskyFunction() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("ดัก panic ได้:", r)
            // log error หรือทำการ cleanup
        }
    }()
    
    // โค้ดที่อาจ panic
    var slice []int
    fmt.Println(slice[5]) // ทำให้ panic: index out of range
}

func main() {
    riskyFunction()
    fmt.Println("โปรแกรมยังทำงานต่อได้")
}
```

**ข้อควรระวัง:** 
- `recover()` จะทำงานก็ต่อเมื่อถูกเรียกภายใน `defer` function เท่านั้น
- ควรใช้ `recover()` ใน goroutine แยกกัน เพราะ panic ใน goroutine จะทำให้ goroutine นั้นหยุด แต่ไม่ crash โปรแกรมหลัก (main goroutine)

#### 2.2 ป้องกัน Panic ก่อนเกิด
แนวทางที่ดีที่สุดคือการป้องกันไม่ให้ panic เกิดขึ้นตั้งแต่แรก:

1. **ตรวจสอบ nil pointer ก่อนใช้งาน**
   ```go
   if myStruct != nil {
       myStruct.Method()
   }
   ```

2. **ตรวจสอบ index ก่อนเข้าถึง slice/array**
   ```go
   if index >= 0 && index < len(mySlice) {
       value := mySlice[index]
   }
   ```

3. **ใช้ comma ok idiom สำหรับ map**
   ```go
   if value, ok := myMap[key]; ok {
       // ใช้ value ได้อย่างปลอดภัย
   }
   ```

4. **ตรวจสอบ error แทนการใช้ panic สำหรับข้อผิดพลาดทั่วไป**
   ```go
   file, err := os.Open("file.txt")
   if err != nil {
       // จัดการ error อย่างเหมาะสม แทนการ panic
       log.Println(err)
       return
   }
   ```

#### 2.3 ใช้เครื่องมือวิเคราะห์
- **Go vet**: `go vet` ช่วยหา potential panics บางประเภท
- **Static analysis tools**: เช่น `nilaway`, `errcheck` ช่วยตรวจสอบการจัดการ nil และ error

---

### 3. Panic ในบริบทอื่น

ถ้าคุณหมายถึง **Kernel Panic** (ใน Linux, macOS, Windows BSOD):
- **สาเหตุ**: 硬件故障, ไดรเวอร์ผิดพลาด, RAM เสีย, OS kernel bug
- **วิธีแก้**: 
  - บันทึก error code (บนหน้าจอ)
  - รีสตาร์ทเครื่อง
  - ตรวจสอบ硬件 (Memtest, SMART disk)
  - อัปเดตไดรเวอร์หรือ OS
  - ดู log จาก `/var/log/kern.log` หรือ `dmesg`

ถ้าคุณหมายถึง **Panic Attack** (อาการตื่นตระหนกทางจิตวิทยา):
- ไม่เกี่ยวกับการเขียนโปรแกรม แต่เป็นอาการทางจิตเวช ควรพบจิตแพทย์หรือผู้เชี่ยวชาญ

---

### สรุป
- **Panic** = ข้อผิดพลาดร้ายแรงที่ทำให้โปรแกรมหยุดทำงาน
- **แก้ไข** = ใช้ `recover()` ดักจับ (ใน Go) + ป้องกันด้วยการตรวจสอบค่าก่อนใช้งาน + ใช้ error handling ที่ดี
- **หลักการ** = ใช้ panic สำหรับ "สถานการณ์ที่ไม่ควรเกิดขึ้น" (programmer error) และใช้ error สำหรับ "สถานการณ์ที่คาดการณ์ได้"
 

### บทที่ 16: อินเทอร์เฟซ (Interfaces)

#### 16.1 อินเทอร์เฟซคืออะไร?
อินเทอร์เฟซคือชุดของ method signatures ชนิดข้อมูลใดๆ ที่ implement method ครบทุก method ในอินเทอร์เฟซ จะถือว่าอินสแตนซ์นั้น implement อินเทอร์เฟซนั้นโดยปริยาย (implicit implementation)
```go
type Speaker interface {
    Speak() string
}

type Dog struct{}

func (d Dog) Speak() string {
    return "Woof!"
}

var s Speaker = Dog{}
fmt.Println(s.Speak())
```

#### 16.2 Empty interface
`interface{}` (หรือ `any` ใน Go 1.18+) เป็นอินเทอร์เฟซที่ไม่มี method ใดๆ ทุกชนิด implement มันได้
```go
var anything interface{}
anything = 42
anything = "hello"
anything = Dog{}
```
ใช้เมื่อต้องการรับค่าทุกชนิด (คล้าย dynamic type)

#### 16.3 Type assertion
ใช้เพื่อดึงค่า concrete type ออกจาก interface
```go
var s interface{} = "hello"
str, ok := s.(string)
if ok {
    fmt.Println(str)
}
```
หากไม่ตรวจสอบ ok จะ panic ถ้าชนิดไม่ตรง

#### 16.4 Type switch
ตรวจสอบชนิดของ interface
```go
switch v := s.(type) {
case string:
    fmt.Println("string:", v)
case int:
    fmt.Println("int:", v)
default:
    fmt.Println("unknown")
}
```

#### 16.5 การ embed interface
```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

type ReadWriter interface {
    Reader
    Writer
}
```

#### 16.6 Interface และ nil
ตัวแปร interface มีค่า nil เมื่อทั้ง type และ value เป็น nil
```go
var s Speaker // nil interface
fmt.Println(s == nil) // true
```
แต่ถ้าเก็บ pointer ที่เป็น nil ไว้:
```go
var d *Dog = nil
var s Speaker = d // s ไม่เป็น nil เพราะ type เป็น *Dog
```

#### 16.7 ข้อดีของ interface
- ช่วยให้โค้ดยืดหยุ่น เปลี่ยน implementation ได้ง่าย
- สนับสนุน dependency injection
- เขียน test ได้ง่าย (ใช้ mock)

#### 16.8 Interface ที่ใช้บ่อย
- `io.Reader`, `io.Writer`
- `http.Handler`
- `error` (มี method Error() string)

---

## ภาคที่ 3: การจัดการโปรเจกต์และโครงสร้างข้อมูลขั้นสูง

### บทที่ 17: Go Modules - การจัดการโปรเจกต์สมัยใหม่

#### 17.1 Go Modules คืออะไร?
Go Modules เป็นระบบจัดการ dependencies อย่างเป็นทางการ เริ่มตั้งแต่ Go 1.11 ทำให้ไม่ต้องพึ่งพา GOPATH

#### 17.2 การเริ่มต้นโมดูล
```bash
go mod init example.com/myproject
```
สร้างไฟล์ `go.mod`:
```
module example.com/myproject

go 1.21
```

#### 17.3 การเพิ่ม dependencies
เมื่อเรา import แพคเกจภายนอกในโค้ดและรัน `go build` หรือ `go mod tidy` Go จะดึง dependencies และเพิ่มลงใน go.mod พร้อมสร้าง go.sum
```bash
go get github.com/gorilla/mux
```

#### 17.4 go.mod
ประกอบด้วย:
- module path
- เวอร์ชัน Go
- require: dependencies และเวอร์ชัน
- replace: แทนที่โมดูลด้วย local path หรือเวอร์ชันอื่น
- exclude: ไม่ใช้บางเวอร์ชัน

#### 17.5 go.sum
เก็บ checksum ของ dependencies เพื่อความปลอดภัยและการตรวจสอบความถูกต้อง

#### 17.6 การอัปเดต dependencies
```bash
go get -u              # อัปเดตทุก dependency
go get -u ./...        # อัปเดตภายในโมดูล
go get github.com/gorilla/mux@v1.8.0  # อัปเดตเฉพาะ
go mod tidy            # ลบ dependencies ที่ไม่ใช้
```

#### 17.7 vendor directory
สามารถสร้าง vendor folder สำหรับเก็บ dependencies ในโปรเจกต์ (สำหรับการ build ที่ไม่ต้องเชื่อมต่อ internet)
```bash
go mod vendor
```
แล้ว build ด้วย `go build -mod=vendor`

#### 17.8 การทำงานกับ private repositories
ตั้งค่า GOPRIVATE:
```bash
go env -w GOPRIVATE=github.com/mycompany/*
```

---

### บทที่ 18: Go Module Proxies

#### 18.1 โมดูลพร็อกซีคืออะไร?
พร็อกซีทำหน้าที่เป็นตัวกลางในการให้บริการโมดูล Go ช่วยให้การดาวน์โหลด dependencies เร็วขึ้น และมั่นใจว่าโค้ดไม่ถูกแก้ไข

#### 18.2 ค่าเริ่มต้น
Go 1.13+ ใช้ proxy `https://proxy.golang.org` เป็น default พร้อม fallback ไป direct

#### 18.3 การตั้งค่า GOPROXY
```bash
go env -w GOPROXY=https://goproxy.io,direct
```
สามารถตั้งค่าเป็น `direct` เพื่อดึงโดยตรงจาก VCS

#### 18.4 การใช้ proxy แบบ private
สำหรับ private repository ควรตั้ง `GOPRIVATE` เพื่อ bypass proxy

#### 18.5 การตรวจสอบ checksum
Go ใช้ checksum database (`sum.golang.org`) ตรวจสอบ integrity

#### 18.6 การปิดใช้งาน checksum สำหรับ private
```bash
go env -w GOSUMDB=off
```

---

### บทที่ 19: การทดสอบหน่วย (Unit Tests)

#### 19.1 การเขียน test
ไฟล์ test ต้องลงท้ายด้วย `_test.go` และฟังก์ชัน test ต้องขึ้นต้นด้วย `Test` และรับ `*testing.T`
```go
// math_test.go
package math

import "testing"

func TestAdd(t *testing.T) {
    result := Add(2, 3)
    expected := 5
    if result != expected {
        t.Errorf("Add(2,3) = %d; want %d", result, expected)
    }
}
```

#### 19.2 การรัน test
```bash
go test                 # รัน test ใน current directory
go test ./...           # รันทั้งหมด
go test -v              # verbose output
go test -run TestAdd    # รันเฉพาะฟังก์ชันที่ match
```

#### 19.3 Table-driven tests
รูปแบบที่นิยมใน Go:
```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name string
        a, b int
        want int
    }{
        {"positive", 2, 3, 5},
        {"negative", -1, -2, -3},
        {"zero", 0, 0, 0},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := Add(tt.a, tt.b); got != tt.want {
                t.Errorf("Add(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
            }
        })
    }
}
```

#### 19.4 การใช้ subtests
`t.Run` ช่วยให้ test ทำงานแยก และสามารถรันย่อยได้

#### 19.5 การทดสอบ HTTP
```go
func TestHandler(t *testing.T) {
    req := httptest.NewRequest("GET", "/", nil)
    w := httptest.NewRecorder()
    handler(w, req)
    if w.Code != http.StatusOK {
        t.Errorf("expected 200, got %d", w.Code)
    }
}
```

#### 19.6 การใช้ test fixtures
สามารถใช้ `TestMain` เพื่อ setup/teardown
```go
func TestMain(m *testing.M) {
    // setup
    code := m.Run()
    // teardown
    os.Exit(code)
}
```

#### 19.7 การครอบคลุม (coverage)
```bash
go test -cover
go test -coverprofile=coverage.out
go tool cover -html=coverage.out   # ดูรายละเอียดใน browser
```

#### 19.8 การใช้ mock
แนะนำให้ใช้ interface เพื่อให้สามารถเปลี่ยน implementation เป็น mock ได้ง่าย

---

### บทที่ 20: อาเรย์ (Arrays)

#### 20.1 อาเรย์คืออะไร?
อาเรย์คือชุดข้อมูลขนาดคงที่ของชนิดเดียวกัน
```go
var arr [5]int               // array ขนาด 5, zero values
arr2 := [3]int{1,2,3}        // กำหนดค่า
arr3 := [...]int{1,2,3}      // ให้ compiler นับขนาด
```

#### 20.2 การเข้าถึงและแก้ไข
```go
arr[0] = 10
value := arr[0]
len(arr) // ความยาว
```

#### 20.3 อาเรย์เป็นค่า (value type)
ใน Go อาเรย์เป็น value type เมื่อส่งเข้า function จะถูก copy ทั้งตัว
```go
func modify(arr [3]int) {
    arr[0] = 100 // ไม่มีผลกับตัวแปรเดิม
}
```

#### 20.4 การใช้ pointer กับอาเรย์
```go
func modifyPtr(arr *[3]int) {
    arr[0] = 100 // แก้ไขตัวแปรเดิม
}
```

#### 20.5 ข้อจำกัด
ขนาดอาเรย์เป็นส่วนหนึ่งของชนิด ดังนั้น `[3]int` และ `[4]int` เป็นคนละชนิดกัน
การส่งอาเรย์ขนาดใหญ่ไปยังฟังก์ชันอาจทำให้ประสิทธิภาพลดลง (copy ข้อมูล)

#### 20.6 การวนลูป
```go
for i := 0; i < len(arr); i++ {
    fmt.Println(arr[i])
}

for i, v := range arr {
    fmt.Println(i, v)
}
```

#### 20.7 อาเรย์แบบหลายมิติ
```go
var matrix [2][3]int
matrix[0][1] = 5
```

---

### บทที่ 21: สไลซ์ (Slices)

#### 21.1 สไลซ์คืออะไร?
สไลซ์เป็นโครงสร้างข้อมูลที่ยืดหยุ่นกว่าอาเรย์ มีขนาด dynamic จริงๆ แล้วสไลซ์คือ view ของอาเรย์
```go
var s []int               // nil slice
s2 := []int{1,2,3}        // slice literal
s3 := make([]int, 5)      // length 5, capacity 5
s4 := make([]int, 5, 10)  // length 5, capacity 10
```

#### 21.2 การเพิ่มสมาชิกด้วย append
```go
s := []int{1,2}
s = append(s, 3)         // [1,2,3]
s = append(s, 4,5)       // [1,2,3,4,5]
```

#### 21.3 การ slicing
```go
arr := [5]int{1,2,3,4,5}
slice := arr[1:4]        // [2,3,4]
```
slicing ใช้ syntax `[low:high]` (low inclusive, high exclusive) มี default: low = 0, high = len

#### 21.4 capacity และ length
- `len(s)` : ความยาวปัจจุบัน
- `cap(s)` : ความจุ (จำนวน element ที่สามารถเก็บได้โดยไม่ต้อง allocate ใหม่)

เมื่อ append แล้วเกิน capacity Go จะ allocate array ใหม่ (ขนาดประมาณสองเท่า)

#### 21.5 การ copy
```go
src := []int{1,2,3}
dst := make([]int, len(src))
copy(dst, src)
```
copy คืนจำนวน element ที่คัดลอก

#### 21.6 การลบ element
ไม่มีฟังก์ชันลบโดยตรง ใช้ slicing รวม:
```go
// ลบ element index i
s = append(s[:i], s[i+1:]...)
```

#### 21.7 การ prepend
```go
s = append([]int{0}, s...)
```

#### 21.8 slice ของ slice (multidimensional)
```go
matrix := make([][]int, 3)
for i := range matrix {
    matrix[i] = make([]int, 3)
}
```

#### 21.9 การส่ง slice ไปยังฟังก์ชัน
สไลซ์ส่งเป็น reference (pointer to underlying array) ดังนั้นการแก้ไข element มีผลต่อตัวแปรเดิม

#### 21.10 ข้อควรระวัง: shared underlying array
เมื่อ slicing จาก slice เดิม อาจเกิดการแชร์อาเรย์ข้างใต้ ทำให้การแก้ไขกระทบกัน

---

### บทที่ 22: แมพ (Maps)

#### 22.1 แมพคืออะไร?
แมพคือโครงสร้างข้อมูลแบบ key-value (คล้าย dictionary) key type ต้อง comparable (==)
```go
var m map[string]int          // nil map, ใส่ค่าไม่ได้
m2 := make(map[string]int)    // empty map
m3 := map[string]int{
    "one": 1,
    "two": 2,
}
```

#### 22.2 การใช้งาน
```go
m := make(map[string]int)
m["key"] = 42
value := m["key"]
delete(m, "key")
```

#### 22.3 การตรวจสอบว่ามี key หรือไม่
```go
val, ok := m["key"]
if ok {
    fmt.Println("found:", val)
} else {
    fmt.Println("not found")
}
```

#### 22.4 การวนลูปผ่าน map
```go
for key, value := range m {
    fmt.Println(key, value)
}
```
ลำดับการวนไม่แน่นอน (randomized)

#### 22.5 Map เป็น reference type
การส่ง map ไปฟังก์ชันจะไม่ copy ข้อมูล การแก้ไขมีผลต่อต้นฉบับ

#### 22.6 nil map
nil map ไม่สามารถเขียนได้ (จะ panic) ต้องใช้ make ก่อนเสมอ

#### 22.7 การใช้ struct เป็น key
struct ที่ fields ทั้งหมด comparable สามารถเป็น key ได้
```go
type Key struct {
    A, B int
}
m := make(map[Key]string)
m[Key{1,2}] = "value"
```

#### 22.8 ประสิทธิภาพ
Map มีการ rehash เมื่อขนาดโตขึ้น ควรให้ hint ขนาดล่วงหน้าด้วย `make(map[K]V, hint)` เพื่อลดการ rehash

#### 22.9 sync.Map
ถ้าต้องการใช้ map ใน concurrent environment ควรใช้ `sync.Map` หรือใช้ mutex

---

### บทที่ 23: การจัดการข้อผิดพลาด (Errors)

#### 23.1 Error ใน Go
Go ใช้ error interface ในการจัดการข้อผิดพลาด ไม่มี exception
```go
type error interface {
    Error() string
}
```
ฟังก์ชันมักคืนค่า error เป็นค่าสุดท้าย
```go
func doSomething() (result int, err error) {
    if somethingWrong {
        return 0, errors.New("something went wrong")
    }
    return 42, nil
}
```

#### 23.2 การสร้าง error
- `errors.New("message")`
- `fmt.Errorf("format %v", arg)`

#### 23.3 การตรวจสอบ error
```go
result, err := doSomething()
if err != nil {
    // จัดการ error
    fmt.Println("Error:", err)
    return
}
// ใช้ result
```

#### 23.4 Custom error types
สามารถสร้าง struct ที่ implement error interface เพื่อเพิ่มข้อมูล
```go
type ValidationError struct {
    Field string
    Value interface{}
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation failed on %s: %v", e.Field, e.Value)
}
```

#### 23.5 การตรวจสอบชนิดของ error
ใช้ type assertion หรือ errors.As
```go
var ve ValidationError
if errors.As(err, &ve) {
    fmt.Println("Field:", ve.Field)
}
```

#### 23.6 การห่อหุ้ม error (wrapping)
ใช้ `%w` ใน fmt.Errorf เพื่อ wrap error
```go
if err != nil {
    return fmt.Errorf("failed to process: %w", err)
}
```
การ unwrap ด้วย `errors.Unwrap` หรือตรวจสอบด้วย `errors.Is`

#### 23.7 errors.Is และ errors.As
- `errors.Is(err, target)` : ตรวจสอบว่า err หรือ err ที่ถูก wrap มี target หรือไม่
- `errors.As(err, &target)` : ดึง error เฉพาะชนิดออกมา

```go
if errors.Is(err, sql.ErrNoRows) {
    // handle no rows
}
```

#### 23.8 panic vs error
ใช้ error สำหรับกรณีที่คาดการณ์ได้ (เช่น ไฟล์ไม่พบ)
ใช้ panic สำหรับกรณีที่ไม่ควรเกิดขึ้น (programmer error) หรือไม่สามารถดำเนินต่อได้

#### 23.9 การใช้ defer กับ error
สามารถใช้ defer เพื่อทำ cleanup แม้มี error

---

## ภาคที่ 4: การพัฒนาแอปพลิเคชันเชิงปฏิบัติ

### บทที่ 24: ฟังก์ชันนิรนาม (Anonymous functions) และ Closure

#### 24.1 ฟังก์ชันนิรนาม
ฟังก์ชันที่ไม่มีชื่อ สามารถกำหนดให้กับตัวแปรหรือส่งเป็น argument
```go
func() {
    fmt.Println("anonymous")
}()

add := func(a, b int) int {
    return a + b
}
result := add(2,3)
```

#### 24.2 Closure
Closure คือฟังก์ชันที่สามารถเข้าถึงตัวแปรที่อยู่นอกขอบเขตของมัน
```go
func counter() func() int {
    count := 0
    return func() int {
        count++
        return count
    }
}

c := counter()
fmt.Println(c()) // 1
fmt.Println(c()) // 2
```

#### 24.3 การใช้งานจริง
- สร้างฟังก์ชัน factory
- ใช้กับ callback (เช่น http.HandlerFunc)
- ใช้กับ goroutine (passing variables)

#### 24.4 ข้อควรระวังเกี่ยวกับ closure ใน loop
```go
func main() {
    funcs := []func(){}
    for i := 0; i < 3; i++ {
        funcs = append(funcs, func() {
            fmt.Println(i) // i เป็นตัวแปรเดียวกัน
        })
    }
    for _, f := range funcs {
        f() // output: 3 3 3
    }
}
```
วิธีแก้: ใช้ตัวแปร local copy
```go
for i := 0; i < 3; i++ {
    i := i // shadow
    funcs = append(funcs, func() {
        fmt.Println(i)
    })
}
```

---

### บทที่ 25: การจัดการข้อมูล JSON และ XML

#### 25.1 JSON encoding
ใช้ struct tags เพื่อกำหนดชื่อ field ใน JSON
```go
type Person struct {
    Name string `json:"name"`
    Age  int    `json:"age,omitempty"`
    Email string `json:"-"` // ไม่รวมใน JSON
}
```

#### 25.2 Marshal (Go -> JSON)
```go
p := Person{Name: "John", Age: 30}
data, err := json.Marshal(p)
if err != nil {
    // handle
}
fmt.Println(string(data)) // {"name":"John","age":30}
```
ใช้ `MarshalIndent` เพื่อจัดรูปแบบ

#### 25.3 Unmarshal (JSON -> Go)
```go
jsonStr := `{"name":"John","age":30}`
var p Person
err := json.Unmarshal([]byte(jsonStr), &p)
```

#### 25.4 การทำงานกับ dynamic JSON
ใช้ `map[string]interface{}` หรือ `interface{}`:
```go
var data map[string]interface{}
json.Unmarshal([]byte(jsonStr), &data)
name := data["name"].(string)
```

#### 25.5 การใช้ json.RawMessage
เพื่อการ decode แบบ lazy

#### 25.6 XML
คล้ายกับ JSON แต่ใช้แท็ก `xml`:
```go
type Person struct {
    XMLName xml.Name `xml:"person"`
    Name    string   `xml:"name"`
    Age     int      `xml:"age"`
}
data, _ := xml.Marshal(p)
```

---

### บทที่ 26: พื้นฐานการสร้าง HTTP Server

#### 26.1 HTTP Handler
`http.Handler` interface มี method `ServeHTTP(http.ResponseWriter, *http.Request)`

วิธีง่าย: ใช้ `http.HandleFunc`
```go
func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}

func main() {
    http.HandleFunc("/", helloHandler)
    http.ListenAndServe(":8080", nil)
}
```

#### 26.2 การใช้ ServeMux
```go
mux := http.NewServeMux()
mux.HandleFunc("/hello", helloHandler)
mux.HandleFunc("/bye", byeHandler)
http.ListenAndServe(":8080", mux)
```

#### 26.3 การอ่าน request data
- Query parameters: `r.URL.Query().Get("name")`
- Form data (POST): `r.ParseForm(); r.FormValue("name")`
- JSON body: `json.NewDecoder(r.Body).Decode(&data)`

#### 26.4 การตั้งค่า header และ status
```go
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusCreated)
```

#### 26.5 Routing ขั้นสูง
ใช้ third-party เช่น `gorilla/mux` หรือ `chi`:
```go
r := mux.NewRouter()
r.HandleFunc("/users/{id}", getUser).Methods("GET")
```

#### 26.6 Middleware
ฟังก์ชันที่ห่อ handler เพื่อเพิ่ม logic (logging, auth)
```go
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Println(r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
    })
}

// ใช้
mux.Use(loggingMiddleware)
```

#### 26.7 การสร้าง HTTP server แบบ custom
```go
srv := &http.Server{
    Addr:         ":8080",
    Handler:      mux,
    ReadTimeout:  10 * time.Second,
    WriteTimeout: 10 * time.Second,
}
srv.ListenAndServe()
```

---

### บทที่ 27: Enum, Iota และ Bitmask

#### 27.1 Enum ใน Go
Go ไม่มี enum ในตัว แต่ใช้ const และ iota
```go
type Weekday int

const (
    Sunday Weekday = iota
    Monday
    Tuesday
    Wednesday
    Thursday
    Friday
    Saturday
)
```
iota เริ่มที่ 0 และเพิ่มทีละ 1

#### 27.2 การกำหนดค่าเริ่มต้นไม่ใช่ 0
```go
const (
    _ = iota // ignore first
    KB = 1 << (10 * iota) // 1 << (10*1) = 1024
    MB // 1 << (10*2) = 1048576
    GB
)
```

#### 27.3 Bitmask
ใช้ bitwise operator เพื่อรวม flags
```go
type Permission uint8

const (
    Read Permission = 1 << iota // 1
    Write                       // 2
    Execute                     // 4
)

// รวม
perm := Read | Write // 3

// ตรวจสอบ
if perm&Read != 0 {
    fmt.Println("has read")
}
```

#### 27.4 การสร้าง string representation
```go
func (p Permission) String() string {
    var flags []string
    if p&Read != 0 {
        flags = append(flags, "Read")
    }
    // ...
    return strings.Join(flags, "|")
}
```

---

### บทที่ 28: วันที่และเวลา

#### 28.1 time.Time
```go
t := time.Now()
fmt.Println(t) // 2024-01-15 10:30:00 +0000 UTC

// สร้างเวลาที่กำหนด
t2 := time.Date(2024, time.January, 15, 10, 30, 0, 0, time.UTC)
```

#### 28.2 การจัดรูปแบบ
Go ใช้รูปแบบอ้างอิง: `Mon Jan 2 15:04:05 MST 2006`
```go
t := time.Now()
fmt.Println(t.Format("2006-01-02 15:04:05"))
fmt.Println(t.Format("02/01/2006"))
```

#### 28.3 การ parse string เป็น time
```go
layout := "2006-01-02 15:04:05"
t, err := time.Parse(layout, "2024-01-15 14:30:00")
```

#### 28.4 การคำนวณเวลา
```go
t := time.Now()
t2 := t.Add(24 * time.Hour)
diff := t2.Sub(t) // 24h0m0s
```

#### 28.5 time.Duration
```go
d, _ := time.ParseDuration("1h30m")
time.Sleep(d)
```

#### 28.6 Timer และ Ticker
- `time.After(d)` : คืน channel ที่จะรับค่าเมื่อครบเวลา
- `time.Ticker` : ส่งค่าทุกๆ ช่วงเวลา

```go
ticker := time.NewTicker(1 * time.Second)
go func() {
    for range ticker.C {
        fmt.Println("tick")
    }
}()
time.Sleep(5 * time.Second)
ticker.Stop()
```

#### 28.7 Timezone
ใช้ `time.LoadLocation("Asia/Bangkok")` เพื่อเปลี่ยน location

---

### บทที่ 29: การจัดเก็บข้อมูล: ไฟล์และฐานข้อมูล

#### 29.1 การอ่านไฟล์
```go
data, err := os.ReadFile("file.txt")
if err != nil {
    log.Fatal(err)
}
fmt.Println(string(data))
```

#### 29.2 การเขียนไฟล์
```go
err := os.WriteFile("output.txt", []byte("hello"), 0644)
```

#### 29.3 การใช้ bufio สำหรับอ่านทีละบรรทัด
```go
file, err := os.Open("file.txt")
if err != nil { ... }
defer file.Close()

scanner := bufio.NewScanner(file)
for scanner.Scan() {
    line := scanner.Text()
    fmt.Println(line)
}
```

#### 29.4 การใช้ database/sql
```go
import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

db, err := sql.Open("mysql", "user:pass@tcp(127.0.0.1:3306)/dbname")
defer db.Close()

rows, err := db.Query("SELECT id, name FROM users")
for rows.Next() {
    var id int
    var name string
    rows.Scan(&id, &name)
}
```

#### 29.5 การใช้ ORM (GORM)
```go
import "gorm.io/gorm"

type User struct {
    gorm.Model
    Name string
    Age  int
}

db.Create(&User{Name: "John", Age: 30})
var user User
db.First(&user, 1)
```

#### 29.6 การจัดการ connection pool
```go
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(25)
db.SetConnMaxLifetime(5 * time.Minute)
```

---

### บทที่ 30: การทำงานพร้อมกัน (Concurrency)

#### 30.1 Goroutine
Goroutine คือ lightweight thread จัดการโดย Go runtime
```go
go func() {
    fmt.Println("hello from goroutine")
}()
```

#### 30.2 Channel
Channel ใช้สื่อสารระหว่าง goroutine
```go
ch := make(chan int)
go func() {
    ch <- 42 // ส่งค่า
}()
value := <-ch // รับค่า
```

#### 30.3 Buffered channel
```go
ch := make(chan int, 2) // buffer size 2
ch <- 1
ch <- 2
```

#### 30.4 การใช้ select
รอหลาย channel
```go
select {
case msg1 := <-ch1:
    fmt.Println(msg1)
case msg2 := <-ch2:
    fmt.Println(msg2)
case <-time.After(1 * time.Second):
    fmt.Println("timeout")
default:
    fmt.Println("no message")
}
```

#### 30.5 Worker pool pattern
```go
func worker(id int, jobs <-chan int, results chan<- int) {
    for j := range jobs {
        results <- j * 2
    }
}

jobs := make(chan int, 100)
results := make(chan int, 100)

for w := 1; w <= 3; w++ {
    go worker(w, jobs, results)
}

for j := 1; j <= 5; j++ {
    jobs <- j
}
close(jobs)

for r := 1; r <= 5; r++ {
    <-results
}
```

#### 30.6 sync.WaitGroup
รอให้ goroutine ทั้งหมดทำงานเสร็จ
```go
var wg sync.WaitGroup
for i := 0; i < 5; i++ {
    wg.Add(1)
    go func(i int) {
        defer wg.Done()
        fmt.Println(i)
    }(i)
}
wg.Wait()
```

#### 30.7 Mutex
ป้องกัน race condition
```go
var mu sync.Mutex
var counter int

mu.Lock()
counter++
mu.Unlock()
```

#### 30.8 การตรวจจับ race condition
```bash
go run -race main.go
go test -race
```

#### 30.9 Atomic operations
```go
import "sync/atomic"
var counter int64
atomic.AddInt64(&counter, 1)
```

---

### บทที่ 31: การบันทึกเหตุการณ์ (Logging)

#### 31.1 log package พื้นฐาน
```go
log.Println("Info message")
log.Printf("User %s logged in", username)
log.Fatal("fatal error") // log แล้ว os.Exit(1)
```

#### 31.2 การตั้งค่า log flags
```go
log.SetFlags(log.LstdFlags | log.Lshortfile)
// 2024/01/15 10:30:00 main.go:42: message
```

#### 31.3 การสร้าง logger แยก
```go
file, _ := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
logger := log.New(file, "INFO: ", log.Ldate|log.Ltime)
logger.Println("application started")
```

#### 31.4 การใช้ structured logging (log/slog) Go 1.21+
```go
import "log/slog"

slog.Info("user login", "user", "john", "ip", "127.0.0.1")
slog.Error("database error", "error", err)
```

#### 31.5 ระดับ log
สามารถตั้งระดับด้วย slog.SetLogLoggerLevel

#### 31.6 การใช้ logging middleware ใน HTTP
```go
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        slog.Info("request",
            "method", r.Method,
            "path", r.URL.Path,
            "duration", time.Since(start),
        )
    })
}
```

---

### บทที่ 32: เทมเพลต (Templates)

#### 32.1 text/template และ html/template
Go มีแพคเกจสองตัว: `text/template` สำหรับข้อความทั่วไป, `html/template` สำหรับ HTML (มีการ escape อัตโนมัติ)

#### 32.2 การใช้งานพื้นฐาน
```go
tmpl, err := template.New("test").Parse("Hello {{.Name}}")
data := struct{ Name string }{"John"}
tmpl.Execute(os.Stdout, data)
```

#### 32.3 การทำงานกับ struct และ map
```go
type Person struct {
    Name string
    Age  int
}
data := Person{"John", 30}
tmpl.Execute(w, data)
```
ใน template: `{{.Name}}` , `{{.Age}}`

#### 32.4 คำสั่งใน template
- `{{if .Condition}} ... {{else}} ... {{end}}`
- `{{range .Items}} ... {{end}}`
- `{{with .Field}} ... {{end}}` เปลี่ยน context

#### 32.5 ฟังก์ชันใน template
```go
funcMap := template.FuncMap{
    "toUpper": strings.ToUpper,
}
tmpl := template.New("test").Funcs(funcMap)
tmpl.Parse(`{{. | toUpper}}`)
```

#### 32.6 การใช้ template files
```go
tmpl := template.Must(template.ParseFiles("index.html"))
tmpl.Execute(w, data)
```

#### 32.7 การรวม template (partials)
```go
tmpl := template.Must(template.ParseGlob("templates/*.html"))
```

#### 32.8 ความปลอดภัยของ html/template
html/template จะ escape อัตโนมัติเพื่อป้องกัน XSS ใช้ `{{. | safeHTML}}` เฉพาะเมื่อมั่นใจ

---

### บทที่ 33: การจัดการค่า Configuration

#### 33.1 การใช้ environment variables
```go
import "os"

dbHost := os.Getenv("DB_HOST")
if dbHost == "" {
    dbHost = "localhost"
}
```

#### 33.2 การใช้ flag package
```go
import "flag"

var port = flag.Int("port", 8080, "server port")
flag.Parse()
fmt.Println("port:", *port)
```

#### 33.3 การใช้ไฟล์ configuration (JSON/YAML)
```go
type Config struct {
    Server struct {
        Port int `json:"port"`
    } `json:"server"`
}

data, _ := os.ReadFile("config.json")
var cfg Config
json.Unmarshal(data, &cfg)
```

#### 33.4 การใช้ viper (popular library)
```go
import "github.com/spf13/viper"

viper.SetConfigName("config")
viper.SetConfigType("yaml")
viper.AddConfigPath(".")
viper.ReadInConfig()

port := viper.GetInt("server.port")
```

#### 33.5 การใช้ struct tags และการ validate
```go
type Config struct {
    Port int `mapstructure:"port" validate:"required,min=1024,max=65535"`
}
```

#### 33.6 การโหลดค่า order (override)
ลำดับความสำคัญ: default -> file -> env -> flag

---

## ภาคที่ 5: สู่การเป็นนักพัฒนา Go มืออาชีพ

### บทที่ 34: การวัดประสิทธิภาพ (Benchmarks)

#### 34.1 การเขียน benchmark
ไฟล์ `_test.go` ฟังก์ชันขึ้นต้นด้วย `Benchmark` รับ `*testing.B`
```go
func BenchmarkAdd(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Add(1, 2)
    }
}
```

#### 34.2 การรัน benchmark
```bash
go test -bench=.           # รันทั้งหมด
go test -bench=Add         # รันเฉพาะ Add
go test -bench=. -benchmem # แสดง memory allocation
```

#### 34.3 การ reset timer
```go
func BenchmarkSomething(b *testing.B) {
    setup()
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        // code to benchmark
    }
}
```

#### 34.4 การเปรียบเทียบประสิทธิภาพ
ใช้ `benchstat` tool

#### 34.5 การเขียน benchmark ที่มีการทำงานที่แตกต่างกัน
```go
func BenchmarkAdd(b *testing.B) {
    tests := []struct{
        name string
        a,b int
    }{
        {"small", 1,2},
        {"large", 1000000, 2000000},
    }
    for _, tt := range tests {
        b.Run(tt.name, func(b *testing.B) {
            for i := 0; i < b.N; i++ {
                Add(tt.a, tt.b)
            }
        })
    }
}
```

---

### บทที่ 35: สร้าง HTTP Client

#### 35.1 การใช้ http.Get
```go
resp, err := http.Get("https://api.example.com/data")
if err != nil { ... }
defer resp.Body.Close()
body, _ := io.ReadAll(resp.Body)
```

#### 35.2 การสร้าง custom client
```go
client := &http.Client{
    Timeout: 10 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:    10,
        IdleConnTimeout: 30 * time.Second,
    },
}
resp, err := client.Get(url)
```

#### 35.3 การส่ง POST request พร้อม JSON
```go
data := map[string]interface{}{"name": "John"}
jsonData, _ := json.Marshal(data)

req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
req.Header.Set("Content-Type", "application/json")

resp, err := client.Do(req)
```

#### 35.4 การจัดการ cookies
```go
jar, _ := cookiejar.New(nil)
client := &http.Client{Jar: jar}
```

#### 35.5 การใช้ context กับ request
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
req := req.WithContext(ctx)
resp, err := client.Do(req)
```

#### 35.6 การ retry logic
ควรใช้ library เช่น `github.com/avast/retry-go`

---

### บทที่ 36: การวิเคราะห์โปรไฟล์ (Program Profiling)

#### 36.1 การใช้ pprof
Go มี built-in profiler สำหรับ CPU, memory, goroutine, block, mutex

#### 36.2 การเปิด pprof endpoint ใน HTTP server
```go
import _ "net/http/pprof"

func main() {
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
    // ...
}
```

#### 36.3 การเก็บ CPU profile
```go
f, _ := os.Create("cpu.prof")
pprof.StartCPUProfile(f)
defer pprof.StopCPUProfile()
// run code
```

#### 36.4 การวิเคราะห์ profile
```bash
go tool pprof -http=:8080 cpu.prof
```

#### 36.5 การตรวจสอบ memory allocation
```bash
go test -memprofile mem.prof -bench .
```

#### 36.6 การวิเคราะห์ goroutine
```bash
curl http://localhost:6060/debug/pprof/goroutine?debug=1
```

---

### บทที่ 37: การจัดการ Context

#### 37.1 Context คืออะไร?
Context ใช้ส่งค่า deadline, cancellation, และ request-scoped values ผ่านขอบเขตของ API

#### 37.2 การสร้าง context
```go
ctx := context.Background()   // empty context, root
ctx := context.TODO()         // เมื่อยังไม่แน่ใจ
```

#### 37.3 Context with cancellation
```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()
// ถ้าเรียก cancel() จะทำให้ ctx.Done() ถูกปิด
```

#### 37.4 Context with timeout
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
```

#### 37.5 Context with deadline
```go
deadline := time.Now().Add(5*time.Second)
ctx, cancel := context.WithDeadline(context.Background(), deadline)
```

#### 37.6 การตรวจสอบ context
```go
select {
case <-ctx.Done():
    fmt.Println("context cancelled:", ctx.Err())
    return
default:
    // do work
}
```

#### 37.7 การส่งค่าใน context (ไม่ควรใช้มาก)
```go
type key string
ctx := context.WithValue(context.Background(), key("userID"), 123)
userID := ctx.Value(key("userID")).(int)
```

#### 37.8 การใช้ context ใน HTTP handler
```go
func handler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    // ใช้ ctx ในการทำงาน
}
```

---

### บทที่ 38: Generics - การเขียนโค้ดแบบยืดหยุ่น

#### 38.1 Generics ใน Go 1.18+
Generics ช่วยให้เขียนฟังก์ชันหรือชนิดข้อมูลที่ทำงานกับหลายชนิดได้ โดยไม่ต้องใช้ interface{} และ type assertion

#### 38.2 ฟังก์ชัน generic
```go
func Map[T any](s []T, f func(T) T) []T {
    result := make([]T, len(s))
    for i, v := range s {
        result[i] = f(v)
    }
    return result
}

// เรียกใช้
nums := []int{1,2,3}
doubled := Map(nums, func(x int) int { return x * 2 })
```

#### 38.3 Type constraints
```go
func Sum[T int | float64](a, b T) T {
    return a + b
}
```

#### 38.4 การใช้ interface เป็น constraint
```go
type Number interface {
    int | int64 | float64
}

func Max[T Number](a, b T) T {
    if a > b {
        return a
    }
    return b
}
```

#### 38.5 Generic types (structs)
```go
type Stack[T any] struct {
    items []T
}

func (s *Stack[T]) Push(item T) {
    s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() T {
    item := s.items[len(s.items)-1]
    s.items = s.items[:len(s.items)-1]
    return item
}
```

#### 38.6 constraints แพคเกจ
Go มีแพคเกจ `constraints` (experimental) สำหรับ constraints ทั่วไป

---

### บทที่ 39: Go กับกระบวนทัศน์ OOP?

#### 39.1 Go ไม่ใช่ภาษา OOP แบบดั้งเดิม
Go ไม่มี class, inheritance, polymorphism ผ่าน method overriding แต่ใช้ composition และ interface

#### 39.2 Encapsulation
ใช้ package visibility (ตัวพิมพ์เล็ก/ใหญ่)

#### 39.3 Polymorphism
ใช้ interface แทน inheritance
```go
type Shape interface {
    Area() float64
}

type Circle struct{ Radius float64 }
func (c Circle) Area() float64 { return math.Pi * c.Radius * c.Radius }

type Rectangle struct{ W, H float64 }
func (r Rectangle) Area() float64 { return r.W * r.H }

func PrintArea(s Shape) {
    fmt.Println(s.Area())
}
```

#### 39.4 Composition over inheritance
```go
type Engine struct{}
func (e Engine) Start() { ... }

type Car struct {
    Engine  // embedded
    Wheels int
}
// Car สามารถเรียก Start() ได้โดยตรง
```

#### 39.5 การ reuse โค้ด
ใช้ function, interface, embedding แทน inheritance

#### 39.6 ข้อดีของแนวทางนี้
- ลด coupling
- เข้าใจง่าย
- หลีกเลี่ยงปัญหา "diamond problem"

---

### บทที่ 40: การอัปเกรดหรือดาวน์เกรดเวอร์ชัน Go

#### 40.1 การติดตั้งหลายเวอร์ชัน
ใช้ `go install golang.org/dl/go1.20@latest` แล้ว `go1.20 download`

#### 40.2 การเปลี่ยนเวอร์ชัน (Linux/macOS)
ใช้ `go get golang.org/dl/go1.21.0` แล้ว `go1.21.0 download` หรือใช้ `gvm` (Go Version Manager)

#### 40.3 การอัปเดต go.mod
```bash
go mod edit -go=1.21
go mod tidy
```

#### 40.4 การตรวจสอบ compatibility
ใช้ `go fix` เพื่ออัปเดตโค้ดอัตโนมัติ

#### 40.5 การทดสอบกับเวอร์ชันใหม่
ใช้ CI/CD เพื่อทดสอบหลายเวอร์ชัน

#### 40.6 การดาวน์เกรด
ถ้าต้องการกลับไปใช้เวอร์ชันเก่า ให้เปลี่ยน PATH หรือใช้ version manager

---

### บทที่ 41: คำแนะนำในการออกแบบโค้ดที่ดี

#### 41.1 Go idioms
- **Keep it simple**: หลีกเลี่ยงความซับซ้อนที่ไม่จำเป็น
- **Use interfaces only when necessary**: อย่าสร้าง interface ก่อน แต่ให้สร้างเมื่อมีความจำเป็น
- **Accept interfaces, return structs**: ฟังก์ชันควรรับ interface, คืน concrete type
- **Use zero values**: ใช้ zero value ให้เป็นประโยชน์
- **Prefer composition over inheritance**

#### 41.2 การตั้งชื่อ
- **Package names**: สั้น, เป็นคำเดียว, ตัวพิมพ์เล็ก
- **Variables**: ใช้ชื่อสั้นภายในขอบเขตเล็ก, ชื่อยาวในขอบเขตใหญ่
- **Getters**: ไม่ใช้ `Get` prefix เช่น `Name()` แทน `GetName()`

#### 41.3 การจัดการ error
- ตรวจสอบ error ทันที
- ห่อ error ด้วยข้อมูลเพิ่มเติม (context)
- อย่าใช้ panic สำหรับการจัดการ error ปกติ

#### 41.4 การจัดโครงสร้างโปรเจกต์
```
/cmd/          # executables
/internal/     # private code
/pkg/          # public library code
/api/          # API definitions
```

#### 41.5 การเขียน documentation
- ใช้ comment ก่อนฟังก์ชัน/ชนิด
- ตัวอย่างใน comment จะปรากฏใน `go doc`

#### 41.6 การใช้ linter
```bash
golangci-lint run
```

#### 41.7 การทำ code review
- ตรวจสอบ error handling
- ดู race conditions
- ตรวจสอบการใช้ goroutine และ channel

---

### บทที่ 42: ชีทสรุป (Cheatsheet)

#### 42.1 คำสั่ง Go ที่ใช้บ่อย
```bash
go mod init <module>      # เริ่มต้นโมดูล
go mod tidy               # จัด dependencies
go build                  # คอมไพล์
go run <file>             # รันโดยตรง
go test ./...             # ทดสอบทั้งหมด
go fmt ./...              # จัดรูปแบบโค้ด
go vet ./...              # ตรวจสอบข้อผิดพลาดที่อาจเกิดขึ้น
go get <pkg>@<version>    # ดาวน์โหลดแพคเกจ
go install <pkg>@latest   # ติดตั้ง binary
go doc <pkg>              # ดู documentation
```

#### 42.2 ชนิดข้อมูล
| Type | Description |
|------|-------------|
| `bool` | true/false |
| `string` | immutable sequence |
| `int`, `int8`, `int16`, `int32`, `int64` | signed integers |
| `uint`, `uint8`(byte), `uint16`, `uint32`, `uint64` | unsigned integers |
| `float32`, `float64` | floating point |
| `complex64`, `complex128` | complex numbers |
| `rune` | alias for int32 (Unicode) |

#### 42.3 การประกาศตัวแปร
```go
var name string
var age int = 30
name := "John"          // short declaration (inside function)
var x, y int = 1, 2
var (                   // block
    a string
    b int
)
```

#### 42.4 โครงสร้างควบคุม
```go
// if
if x > 0 { ... } else if x < 0 { ... } else { ... }

// for
for i := 0; i < 10; i++ { ... }
for condition { ... }
for { ... break }

// switch
switch value {
case 1:
    ...
default:
    ...
}
```

#### 42.5 slice และ map
```go
// slice
s := []int{1,2,3}
s = append(s, 4)
s2 := make([]int, 5, 10)
copy(dst, src)

// map
m := make(map[string]int)
m["key"] = 1
delete(m, "key")
val, ok := m["key"]
```

#### 42.6 channel
```go
ch := make(chan int)
ch <- 42
val := <-ch
close(ch)
```

#### 42.7 การทำงานกับไฟล์
```go
data, _ := os.ReadFile("file.txt")
os.WriteFile("out.txt", []byte("data"), 0644)
```

#### 42.8 JSON
```go
type T struct {
    Name string `json:"name"`
}
data, _ := json.Marshal(t)
json.Unmarshal(data, &t)
```

#### 42.9 HTTP
```go
http.HandleFunc("/", handler)
http.ListenAndServe(":8080", nil)
```

#### 42.10 Context
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
```

---

## สรุป

คู่มือนี้ครอบคลุมเนื้อหาตั้งแต่พื้นฐานไปจนถึงระดับสูงของภาษา Go หวังว่าจะเป็นประโยชน์สำหรับผู้ที่ต้องการเริ่มต้นและพัฒนาทักษะการเขียนโปรแกรมด้วย Go อย่างจริงจัง

การฝึกฝนอย่างต่อเนื่องและการนำไปใช้ในโปรเจกต์จริงจะช่วยให้คุณเข้าใจและเชี่ยวชาญภาษา Go ได้ดียิ่งขึ้น ขอให้สนุกกับการเขียนโปรแกรมด้วย Go!


## บทที่ 43: เครื่องมือและไลบรารียอดนิยมสำหรับการพัฒนาแอปพลิเคชัน Go

### 43.1 chi – เราเตอร์และมิดเดิลแวร์

[chi](https://github.com/go-chi/chi) เป็น lightweight router ที่มีประสิทธิภาพสูง รองรับ middleware และเข้ากันได้กับ net/http มาตรฐาน

**การติดตั้ง**
```bash
go get github.com/go-chi/chi/v5
```

**ตัวอย่างพื้นฐาน**
```go
package main

import (
    "net/http"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
)

func main() {
    r := chi.NewRouter()
    
    // ใช้ middleware พื้นฐาน
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    
    // กำหนด routes
    r.Get("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello World"))
    })
    
    // route พร้อมพารามิเตอร์
    r.Get("/users/{id}", getUser)
    
    // group routes
    r.Route("/api", func(r chi.Router) {
        r.Get("/users", listUsers)
        r.Post("/users", createUser)
        
        // sub-router พร้อม middleware เฉพาะ
        r.Route("/admin", func(r chi.Router) {
            r.Use(adminOnly)
            r.Get("/dashboard", adminDashboard)
        })
    })
    
    http.ListenAndServe(":8080", r)
}

func getUser(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    w.Write([]byte("User ID: " + id))
}
```

**มิดเดิลแวร์ที่ใช้บ่อย**
- `middleware.Logger` – บันทึก request
- `middleware.Recoverer` – จับ panic
- `middleware.Timeout` – กำหนด timeout
- `middleware.Compress` – บีบอัด response
- `middleware.RealIP` – ดึง IP จริงจาก proxy

---

### 43.2 viper – การจัดการ configuration

[viper](https://github.com/spf13/viper) รองรับหลายรูปแบบ (JSON, YAML, ENV, flags) และสามารถโหลดจากไฟล์, environment variables, หรือ remote system

**การติดตั้ง**
```bash
go get github.com/spf13/viper
```

**ตัวอย่างการใช้งาน**
```go
package config

import (
    "log"
    "github.com/spf13/viper"
)

type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    Redis    RedisConfig
    JWT      JWTConfig
}

type ServerConfig struct {
    Port int
    Mode string
}

type DatabaseConfig struct {
    Host     string
    Port     int
    User     string
    Password string
    Name     string
}

func LoadConfig() (*Config, error) {
    viper.SetConfigName("config")      // ชื่อไฟล์ (ไม่รวมนามสกุล)
    viper.SetConfigType("yaml")        // yaml, json, toml, etc.
    viper.AddConfigPath(".")           // path ที่ค้นหา
    viper.AddConfigPath("/etc/app/")
    viper.AutomaticEnv()               // อ่านจาก environment variables
    
    // กำหนดค่า default
    viper.SetDefault("server.port", 8080)
    viper.SetDefault("server.mode", "debug")
    
    // อ่านไฟล์
    if err := viper.ReadInConfig(); err != nil {
        log.Printf("Config file not found: %v", err)
    }
    
    var cfg Config
    if err := viper.Unmarshal(&cfg); err != nil {
        return nil, err
    }
    
    return &cfg, nil
}
```

**ไฟล์ config.yaml ตัวอย่าง**
```yaml
server:
  port: 8080
  mode: release

database:
  host: localhost
  port: 3306
  user: root
  password: secret
  name: mydb

redis:
  addr: localhost:6379
  password: ""
  db: 0

jwt:
  secret: "your-secret-key"
  access_expiry: 15m
  refresh_expiry: 7d
```

---

### 43.3 cobra – การสร้าง CLI

[cobra](https://github.com/spf13/cobra) เป็นไลบรารีสำหรับสร้าง command-line application รองรับ commands, flags, และ subcommands

**การติดตั้ง**
```bash
go get -u github.com/spf13/cobra/cobra
```

**การเริ่มต้นโปรเจกต์**
```bash
cobra init --pkg-name mycli
```

**ตัวอย่างการเพิ่ม command**
```go
// cmd/root.go
package cmd

import (
    "fmt"
    "os"
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "mycli",
    Short: "My CLI application",
    Long:  "A sample CLI built with cobra",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Hello from CLI")
    },
}

// cmd/serve.go
var serveCmd = &cobra.Command{
    Use:   "serve",
    Short: "Start the server",
    Run: func(cmd *cobra.Command, args []string) {
        port, _ := cmd.Flags().GetInt("port")
        fmt.Printf("Starting server on port %d\n", port)
    },
}

func init() {
    serveCmd.Flags().IntP("port", "p", 8080, "port to listen on")
    rootCmd.AddCommand(serveCmd)
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
```

**การใช้งาน**
```bash
go run main.go serve --port 3000
```

---

### 43.4 gorm – ORM

[gorm](https://gorm.io) เป็น ORM ที่มีฟีเจอร์ครบถ้วน: associations, hooks, preloading, transactions, etc.

**การติดตั้ง**
```bash
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql
```

**ตัวอย่างการใช้งาน**
```go
package models

import (
    "gorm.io/gorm"
    "gorm.io/driver/mysql"
)

type User struct {
    gorm.Model
    Name     string `gorm:"size:100;not null"`
    Email    string `gorm:"uniqueIndex;size:100;not null"`
    Password string `gorm:"not null"`
    IsActive bool   `gorm:"default:true"`
}

func InitDB(cfg *Config) (*gorm.DB, error) {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
    
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    
    // Auto migrate
    db.AutoMigrate(&User{})
    
    return db, nil
}
```

**CRUD ตัวอย่าง**
```go
// Create
user := User{Name: "John", Email: "john@example.com", Password: "hashed"}
db.Create(&user)

// Read
var user User
db.First(&user, 1)                     // by primary key
db.Where("email = ?", "john@example.com").First(&user)

// Update
db.Model(&user).Update("Name", "John Doe")
db.Model(&user).Updates(User{Name: "John Doe", IsActive: false})

// Delete
db.Delete(&user)
```

**Preload Associations**
```go
type Order struct {
    gorm.Model
    UserID uint
    Amount float64
    User   User // belongs to
}

// Preload when querying
var orders []Order
db.Preload("User").Find(&orders)
```

---

### 43.5 validator – การตรวจสอบข้อมูล

[validator](https://github.com/go-playground/validator) ใช้ struct tags ในการกำหนด validation rules

**การติดตั้ง**
```bash
go get github.com/go-playground/validator/v10
```

**ตัวอย่าง**
```go
package main

import (
    "fmt"
    "github.com/go-playground/validator/v10"
)

type RegisterRequest struct {
    Name     string `validate:"required,min=3,max=50"`
    Email    string `validate:"required,email"`
    Password string `validate:"required,min=8"`
    Age      int    `validate:"gte=18,lte=99"`
}

func main() {
    validate := validator.New()
    
    req := RegisterRequest{
        Name:     "Jo",
        Email:    "invalid",
        Password: "short",
        Age:      16,
    }
    
    err := validate.Struct(req)
    if err != nil {
        for _, err := range err.(validator.ValidationErrors) {
            fmt.Printf("Field %s failed on tag %s\n", err.Field(), err.Tag())
        }
    }
}
```

**Tags ที่ใช้บ่อย**
- `required` – ต้องมีค่า
- `email` – รูปแบบ email
- `min`, `max` – ขนาดต่ำสุด/สูงสุด
- `gte`, `lte` – มากกว่าหรือเท่ากับ, น้อยกว่าหรือเท่ากับ
- `oneof=red blue` – ต้องเป็นหนึ่งในค่าที่กำหนด
- `uuid` – ต้องเป็น UUID
- `url` – ต้องเป็น URL

---

### 43.6 jwt – การจัดการ JWT

[jwt-go](https://github.com/golang-jwt/jwt) เป็นไลบรารีมาตรฐานสำหรับ JWT

**การติดตั้ง**
```bash
go get github.com/golang-jwt/jwt/v5
```

**ตัวอย่างการสร้างและตรวจสอบ token**
```go
package jwtutil

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    UserID uint `json:"user_id"`
    jwt.RegisteredClaims
}

var secretKey = []byte("your-secret-key")

func GenerateAccessToken(userID uint) (string, error) {
    claims := Claims{
        UserID: userID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(secretKey)
}

func ValidateToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return secretKey, nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }
    
    return nil, jwt.ErrInvalidKey
}
```

---

### 43.7 zap – structured logging

[zap](https://github.com/uber-go/zap) เป็น logger ที่มีความเร็วสูงและรองรับ structured logging

**การติดตั้ง**
```bash
go get go.uber.org/zap
```

**ตัวอย่าง**
```go
package logger

import (
    "go.uber.org/zap"
)

var Log *zap.Logger

func InitLogger(mode string) {
    var err error
    if mode == "production" {
        Log, err = zap.NewProduction()
    } else {
        Log, err = zap.NewDevelopment()
    }
    if err != nil {
        panic(err)
    }
    defer Log.Sync()
}

func main() {
    InitLogger("development")
    defer Log.Sync()
    
    Log.Info("User logged in",
        zap.String("user", "john"),
        zap.Int("id", 123),
        zap.Duration("duration", time.Second*2),
    )
    
    Log.Error("Database error",
        zap.Error(err),
        zap.String("query", "SELECT * FROM users"),
    )
}
```

---

### 43.8 gomail – การส่งอีเมล

[gomail](https://github.com/go-gomail/gomail) เป็นไลบรารีที่ใช้งานง่ายสำหรับ SMTP

**การติดตั้ง**
```bash
go get gopkg.in/gomail.v2
```

**ตัวอย่าง**
```go
package mail

import (
    "gopkg.in/gomail.v2"
)

type Mailer struct {
    dialer *gomail.Dialer
    from   string
}

func NewMailer(host string, port int, user, pass, from string) *Mailer {
    return &Mailer{
        dialer: gomail.NewDialer(host, port, user, pass),
        from:   from,
    }
}

func (m *Mailer) Send(to, subject, body string) error {
    msg := gomail.NewMessage()
    msg.SetHeader("From", m.from)
    msg.SetHeader("To", to)
    msg.SetHeader("Subject", subject)
    msg.SetBody("text/html", body)
    
    return m.dialer.DialAndSend(msg)
}
```

---

### 43.9 hermes – สร้าง HTML email ที่สวยงาม

[hermes](https://github.com/matcornic/hermes) ใช้สร้าง email template แบบ responsive

**การติดตั้ง**
```bash
go get github.com/matcornic/hermes/v2
```

**ตัวอย่าง**
```go
package email

import (
    "github.com/matcornic/hermes/v2"
)

type EmailGenerator struct {
    h hermes.Hermes
}

func NewEmailGenerator(appURL, appName string) *EmailGenerator {
    h := hermes.Hermes{
        Product: hermes.Product{
            Name: appName,
            Link: appURL,
            Logo: appURL + "/logo.png",
        },
    }
    return &EmailGenerator{h: h}
}

func (g *EmailGenerator) WelcomeEmail(name, verifyURL string) (string, error) {
    email := hermes.Email{
        Body: hermes.Body{
            Name: name,
            Intros: []string{
                "Welcome to our platform!",
            },
            Actions: []hermes.Action{
                {
                    Instructions: "Please click below to verify your email address:",
                    Button: hermes.Button{
                        Color: "#22BC66",
                        Text:  "Verify Email",
                        Link:  verifyURL,
                    },
                },
            },
            Outros: []string{
                "If you didn't sign up, you can ignore this email.",
            },
        },
    }
    
    return g.h.GenerateHTML(email)
}
```

---

### 43.10 air – hot-reload

[air](https://github.com/cosmtrek/air) ใช้สำหรับ reload อัตโนมัติเมื่อไฟล์เปลี่ยนแปลง เหมาะสำหรับการพัฒนา

**การติดตั้ง**
```bash
go install github.com/cosmtrek/air@latest
```

**การใช้งาน**
สร้างไฟล์ `.air.toml` ใน root ของโปรเจกต์ (หรือใช้ default) แล้วรัน:
```bash
air
```

**ตัวอย่าง .air.toml**
```toml
root = "."
tmp_dir = "tmp"

[build]
  cmd = "go build -o ./tmp/main ."
  bin = "./tmp/main"
  include_ext = ["go", "tpl", "tmpl", "html"]
  exclude_dir = ["assets", "tmp", "vendor"]
  delay = 1000
  stop_on_error = true
  send_interrupt = false
  kill_delay = 500
```

---

## บทที่ 44: การออกแบบแอปพลิเคชันด้วย Clean Architecture

### 44.1 โครงสร้างโปรเจกต์แบบ Clean Architecture

โครงสร้างที่แบ่งเป็น 3 ชั้นหลัก:
- **Delivery** – รับและส่งข้อมูล (HTTP handlers, gRPC, CLI)
- **Usecase** – business logic (interactors)
- **Repository** – การเข้าถึงข้อมูล (database, external API)

นอกจากนี้ยังมี **Models** (Entities) ที่ใช้ร่วมกันทุกชั้น

**โครงสร้างโฟลเดอร์ตัวอย่าง**
```
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── config/            # การตั้งค่า
│   ├── delivery/
│   │   ├── http/          # HTTP handlers
│   │   │   ├── handler.go
│   │   │   ├── middleware.go
│   │   │   └── routes.go
│   │   └── cli/           # (optional) CLI commands
│   ├── models/            # entities / DTOs
│   ├── repository/        # implementations
│   │   ├── user_repo.go
│   │   ├── user_repo_mock.go (สำหรับ test)
│   │   └── redis/         # redis implementation
│   └── usecase/           # business logic
│       ├── user_usecase.go
│       └── auth_usecase.go
├── pkg/                   # reusable packages
│   ├── jwt/
│   ├── mail/
│   └── redis/
├── go.mod
└── config.yaml
```

### 44.2 Models (Entities)

**models/user.go**
```go
package models

import "time"

type User struct {
    ID          uint      `json:"id"`
    Name        string    `json:"name"`
    Email       string    `json:"email"`
    Password    string    `json:"-"`          // ไม่ส่งกลับใน JSON
    IsVerified  bool      `json:"is_verified"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// สำหรับการ register
type RegisterRequest struct {
    Name     string `json:"name" validate:"required,min=3"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

// สำหรับ login
type LoginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
}
```

### 44.3 Repository Interface

**internal/repository/user_repo.go**
```go
package repository

import (
    "context"
    "your-project/internal/models"
)

type UserRepository interface {
    Create(ctx context.Context, user *models.User) error
    GetByID(ctx context.Context, id uint) (*models.User, error)
    GetByEmail(ctx context.Context, email string) (*models.User, error)
    Update(ctx context.Context, user *models.User) error
    Delete(ctx context.Context, id uint) error
}
```

**Implementation ด้วย GORM (internal/repository/user_repo_gorm.go)**
```go
package repository

import (
    "context"
    "gorm.io/gorm"
    "your-project/internal/models"
)

type userRepositoryGorm struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepositoryGorm{db: db}
}

func (r *userRepositoryGorm) Create(ctx context.Context, user *models.User) error {
    return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepositoryGorm) GetByID(ctx context.Context, id uint) (*models.User, error) {
    var user models.User
    err := r.db.WithContext(ctx).First(&user, id).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}
// ... methods อื่นๆ
```

### 44.4 Usecase

**internal/usecase/user_usecase.go**
```go
package usecase

import (
    "context"
    "errors"
    "golang.org/x/crypto/bcrypt"
    "your-project/internal/models"
    "your-project/internal/repository"
)

type UserUsecase interface {
    Register(ctx context.Context, req *models.RegisterRequest) (*models.User, error)
    GetUserByID(ctx context.Context, id uint) (*models.User, error)
}

type userUsecase struct {
    userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
    return &userUsecase{userRepo: userRepo}
}

func (u *userUsecase) Register(ctx context.Context, req *models.RegisterRequest) (*models.User, error) {
    // ตรวจสอบว่ามี email ซ้ำหรือไม่
    existing, _ := u.userRepo.GetByEmail(ctx, req.Email)
    if existing != nil {
        return nil, errors.New("email already registered")
    }
    
    // Hash password
    hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }
    
    user := &models.User{
        Name:       req.Name,
        Email:      req.Email,
        Password:   string(hashed),
        IsVerified: false,
    }
    
    if err := u.userRepo.Create(ctx, user); err != nil {
        return nil, err
    }
    
    return user, nil
}
```

### 44.5 Delivery (HTTP)

**internal/delivery/http/handler.go**
```go
package http

import (
    "encoding/json"
    "net/http"
    "your-project/internal/usecase"
    "your-project/internal/models"
    "github.com/go-chi/chi/v5"
    "github.com/go-playground/validator/v10"
)

type UserHandler struct {
    userUsecase usecase.UserUsecase
    validate    *validator.Validate
}

func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
    return &UserHandler{
        userUsecase: userUsecase,
        validate:    validator.New(),
    }
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
    var req models.RegisterRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }
    
    if err := h.validate.Struct(req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    user, err := h.userUsecase.Register(r.Context(), &req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}
```

**internal/delivery/http/routes.go**
```go
package http

import (
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
)

func SetupRoutes(userHandler *UserHandler, authHandler *AuthHandler) *chi.Mux {
    r := chi.NewRouter()
    
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    
    r.Post("/api/register", userHandler.Register)
    r.Post("/api/login", authHandler.Login)
    
    // Protected routes
    r.Group(func(r chi.Router) {
        r.Use(authHandler.AuthMiddleware)
        r.Get("/api/users/{id}", userHandler.GetUser)
        // ...
    })
    
    return r
}
```

---

## บทที่ 45: การนำฟีเจอร์ไปใช้งานแบบครบวงจร

### 45.1 JWT Authentication with Refresh Token (Redis)

**แนวคิด**
- **Access Token** (short-lived, 15 นาที) ใช้สำหรับ request ปกติ
- **Refresh Token** (long-lived, 7 วัน) ใช้ขอ access token ใหม่
- เก็บ refresh token ใน Redis พร้อม user ID และ expiration

**โครงสร้าง Redis**
```
refresh_token:<token> -> user_id
user_tokens:<user_id> -> set of token ids (สำหรับ force logout)
```

**Implementation**

**pkg/redis/client.go**
```go
package redis

import (
    "context"
    "time"
    "github.com/go-redis/redis/v8"
)

type Client struct {
    rdb *redis.Client
}

func NewRedisClient(addr, password string, db int) (*Client, error) {
    rdb := redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: password,
        DB:       db,
    })
    
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := rdb.Ping(ctx).Err(); err != nil {
        return nil, err
    }
    
    return &Client{rdb: rdb}, nil
}

func (c *Client) SetRefreshToken(ctx context.Context, userID uint, token string, expiry time.Duration) error {
    key := "refresh_token:" + token
    return c.rdb.Set(ctx, key, userID, expiry).Err()
}

func (c *Client) GetUserIDByRefreshToken(ctx context.Context, token string) (uint, error) {
    key := "refresh_token:" + token
    val, err := c.rdb.Get(ctx, key).Result()
    if err != nil {
        return 0, err
    }
    var userID uint
    fmt.Sscanf(val, "%d", &userID)
    return userID, nil
}

func (c *Client) DeleteRefreshToken(ctx context.Context, token string) error {
    key := "refresh_token:" + token
    return c.rdb.Del(ctx, key).Err()
}
```

**internal/usecase/auth_usecase.go**
```go
package usecase

import (
    "context"
    "errors"
    "time"
    "your-project/internal/models"
    "your-project/internal/repository"
    "your-project/pkg/jwtutil"
    "your-project/pkg/redis"
    "golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
    Login(ctx context.Context, email, password string) (*models.TokenResponse, error)
    RefreshAccessToken(ctx context.Context, refreshToken string) (*models.TokenResponse, error)
    Logout(ctx context.Context, userID uint, refreshToken string) error
}

type authUsecase struct {
    userRepo repository.UserRepository
    redis    *redis.Client
    jwtCfg   *jwtutil.Config
}

func NewAuthUsecase(userRepo repository.UserRepository, redis *redis.Client, jwtCfg *jwtutil.Config) AuthUsecase {
    return &authUsecase{userRepo: userRepo, redis: redis, jwtCfg: jwtCfg}
}

func (a *authUsecase) Login(ctx context.Context, email, password string) (*models.TokenResponse, error) {
    user, err := a.userRepo.GetByEmail(ctx, email)
    if err != nil {
        return nil, errors.New("invalid credentials")
    }
    
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        return nil, errors.New("invalid credentials")
    }
    
    // Generate tokens
    accessToken, err := jwtutil.GenerateAccessToken(user.ID, a.jwtCfg.AccessExpiry)
    if err != nil {
        return nil, err
    }
    
    refreshToken, err := jwtutil.GenerateRefreshToken(user.ID, a.jwtCfg.RefreshExpiry)
    if err != nil {
        return nil, err
    }
    
    // Store refresh token in Redis
    if err := a.redis.SetRefreshToken(ctx, user.ID, refreshToken, a.jwtCfg.RefreshExpiry); err != nil {
        return nil, err
    }
    
    return &models.TokenResponse{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        ExpiresIn:    int64(a.jwtCfg.AccessExpiry.Seconds()),
    }, nil
}

func (a *authUsecase) RefreshAccessToken(ctx context.Context, refreshToken string) (*models.TokenResponse, error) {
    // Validate refresh token
    claims, err := jwtutil.ValidateRefreshToken(refreshToken, a.jwtCfg.Secret)
    if err != nil {
        return nil, errors.New("invalid refresh token")
    }
    
    // Check if refresh token exists in Redis
    userID, err := a.redis.GetUserIDByRefreshToken(ctx, refreshToken)
    if err != nil {
        return nil, errors.New("refresh token revoked")
    }
    
    if userID != claims.UserID {
        return nil, errors.New("token mismatch")
    }
    
    // Generate new access token
    accessToken, err := jwtutil.GenerateAccessToken(userID, a.jwtCfg.AccessExpiry)
    if err != nil {
        return nil, err
    }
    
    return &models.TokenResponse{
        AccessToken:  accessToken,
        RefreshToken: refreshToken, // same refresh token
        ExpiresIn:    int64(a.jwtCfg.AccessExpiry.Seconds()),
    }, nil
}

func (a *authUsecase) Logout(ctx context.Context, userID uint, refreshToken string) error {
    return a.redis.DeleteRefreshToken(ctx, refreshToken)
}
```

**HTTP Middleware สำหรับตรวจสอบ access token**
```go
func (h *AuthHandler) AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Missing authorization header", http.StatusUnauthorized)
            return
        }
        
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
            return
        }
        
        claims, err := jwtutil.ValidateAccessToken(parts[1], h.jwtSecret)
        if err != nil {
            http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
            return
        }
        
        // ส่ง userID ไปยัง context
        ctx := context.WithValue(r.Context(), "userID", claims.UserID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

### 45.2 Cached User in Redis

เพื่อลดการ query ฐานข้อมูลซ้ำ เรา cache user data ใน Redis

**pkg/redis/user_cache.go**
```go
package redis

import (
    "context"
    "encoding/json"
    "time"
    "your-project/internal/models"
)

func (c *Client) GetUser(ctx context.Context, userID uint) (*models.User, error) {
    key := fmt.Sprintf("user:%d", userID)
    val, err := c.rdb.Get(ctx, key).Result()
    if err != nil {
        return nil, err
    }
    
    var user models.User
    if err := json.Unmarshal([]byte(val), &user); err != nil {
        return nil, err
    }
    return &user, nil
}

func (c *Client) SetUser(ctx context.Context, user *models.User, ttl time.Duration) error {
    key := fmt.Sprintf("user:%d", user.ID)
    data, err := json.Marshal(user)
    if err != nil {
        return err
    }
    return c.rdb.Set(ctx, key, data, ttl).Err()
}

func (c *Client) DeleteUser(ctx context.Context, userID uint) error {
    key := fmt.Sprintf("user:%d", userID)
    return c.rdb.Del(ctx, key).Err()
}
```

**Cache Layer ใน Repository**
```go
func (r *userRepositoryCached) GetByID(ctx context.Context, id uint) (*models.User, error) {
    // Try cache first
    user, err := r.cache.GetUser(ctx, id)
    if err == nil {
        return user, nil
    }
    
    // Cache miss, get from DB
    user, err = r.repo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // Store in cache
    go r.cache.SetUser(context.Background(), user, 10*time.Minute)
    
    return user, nil
}
```

### 45.3 Email Verification

**Workflow**
1. ผู้ใช้สมัครสมาชิก → สร้าง user ใน DB (is_verified = false)
2. สร้าง verification token (UUID) เก็บใน Redis หรือ DB พร้อม user ID
3. ส่งอีเมลพร้อมลิงก์ `https://yourapp.com/verify?token=xxx`
4. ผู้ใช้คลิกลิงก์ → API ตรวจสอบ token → update user.is_verified = true
5. ลบ token ออกจาก Redis

**Implementation**

**internal/usecase/verification_usecase.go**
```go
package usecase

import (
    "context"
    "crypto/rand"
    "encoding/hex"
    "time"
    "your-project/internal/models"
    "your-project/pkg/redis"
)

type VerificationUsecase interface {
    SendVerificationEmail(user *models.User) error
    VerifyEmail(ctx context.Context, token string) error
}

type verificationUsecase struct {
    userRepo   repository.UserRepository
    redis      *redis.Client
    mailer     *mail.Mailer
    emailGen   *email.EmailGenerator
    appURL     string
}

func (v *verificationUsecase) SendVerificationEmail(user *models.User) error {
    token, err := generateToken()
    if err != nil {
        return err
    }
    
    // Store token in Redis with TTL (24h)
    key := "verify_token:" + token
    if err := v.redis.SetVerifyToken(context.Background(), key, user.ID, 24*time.Hour); err != nil {
        return err
    }
    
    verifyURL := v.appURL + "/verify?token=" + token
    body, err := v.emailGen.VerifyEmail(user.Name, verifyURL)
    if err != nil {
        return err
    }
    
    return v.mailer.Send(user.Email, "Verify your email", body)
}

func (v *verificationUsecase) VerifyEmail(ctx context.Context, token string) error {
    userID, err := v.redis.GetVerifyToken(ctx, token)
    if err != nil {
        return errors.New("invalid or expired token")
    }
    
    user, err := v.userRepo.GetByID(ctx, userID)
    if err != nil {
        return err
    }
    
    user.IsVerified = true
    if err := v.userRepo.Update(ctx, user); err != nil {
        return err
    }
    
    // Delete token from Redis
    v.redis.DeleteVerifyToken(ctx, token)
    
    return nil
}

func generateToken() (string, error) {
    bytes := make([]byte, 16)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return hex.EncodeToString(bytes), nil
}
```

### 45.4 Forget/Reset Password

**Workflow**
1. ผู้ใช้ขอ reset password ด้วย email
2. ระบบสร้าง reset token (UUID) เก็บใน Redis พร้อม user ID และ TTL (1h)
3. ส่งอีเมลที่มีลิงก์ `https://yourapp.com/reset-password?token=xxx`
4. ผู้ใช้คลิกลิงก์ → หน้า form ใส่ password ใหม่
5. ส่ง request พร้อม token และ password ใหม่
6. ตรวจสอบ token, update password, ลบ token

**Implementation**

**internal/usecase/password_usecase.go**
```go
package usecase

import (
    "context"
    "errors"
    "time"
    "golang.org/x/crypto/bcrypt"
    "your-project/internal/models"
)

type PasswordUsecase interface {
    SendResetEmail(email string) error
    ResetPassword(ctx context.Context, token, newPassword string) error
}

func (p *passwordUsecase) SendResetEmail(email string) error {
    user, err := p.userRepo.GetByEmail(context.Background(), email)
    if err != nil {
        // ไม่บอกว่ามี user หรือไม่เพื่อความปลอดภัย
        return nil
    }
    
    token, _ := generateToken()
    key := "reset_token:" + token
    p.redis.SetResetToken(context.Background(), key, user.ID, 1*time.Hour)
    
    resetURL := p.appURL + "/reset-password?token=" + token
    body, _ := p.emailGen.ResetPasswordEmail(user.Name, resetURL)
    p.mailer.Send(user.Email, "Reset your password", body)
    
    return nil
}

func (p *passwordUsecase) ResetPassword(ctx context.Context, token, newPassword string) error {
    userID, err := p.redis.GetResetToken(ctx, token)
    if err != nil {
        return errors.New("invalid or expired token")
    }
    
    user, err := p.userRepo.GetByID(ctx, userID)
    if err != nil {
        return err
    }
    
    hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    user.Password = string(hashed)
    
    if err := p.userRepo.Update(ctx, user); err != nil {
        return err
    }
    
    p.redis.DeleteResetToken(ctx, token)
    return nil
}
```

---

## 45.5 Workflow สรุปสำหรับการพัฒนา

### Task List Template

#### Phase 1: Project Setup
- [ ] Initialize Go module
- [ ] Setup directory structure (cmd, internal, pkg)
- [ ] Configure viper with config.yaml
- [ ] Setup GORM and database connection
- [ ] Setup Redis client
- [ ] Setup logger (zap)
- [ ] Setup validator
- [ ] Create Makefile for common tasks (build, test, run)

#### Phase 2: Models & Repository
- [ ] Define models (User, etc.)
- [ ] Implement UserRepository interface
- [ ] Implement GORM-based user repository
- [ ] Write unit tests for repository
- [ ] Implement Redis caching layer for user
- [ ] Implement refresh token storage in Redis

#### Phase 3: Usecases
- [ ] Implement UserUsecase (Register, GetUser)
- [ ] Implement AuthUsecase (Login, Logout, Refresh)
- [ ] Implement VerificationUsecase
- [ ] Implement PasswordUsecase
- [ ] Write unit tests for usecases with mocks

#### Phase 4: Delivery (HTTP)
- [ ] Create HTTP handlers
- [ ] Setup chi router with middleware (logger, recoverer, cors)
- [ ] Implement auth middleware (JWT validation)
- [ ] Map routes
- [ ] Integrate with usecases
- [ ] Write integration tests (httptest)

#### Phase 5: Email & Templates
- [ ] Setup SMTP mailer (gomail)
- [ ] Setup hermes for email templates
- [ ] Create email templates: welcome, verify, reset password
- [ ] Implement email sending in usecases

#### Phase 6: Advanced Features
- [ ] Implement JWT access/refresh token
- [ ] Implement token revocation on logout
- [ ] Implement user caching in Redis
- [ ] Add validation for all requests
- [ ] Add rate limiting middleware (optional)
- [ ] Add graceful shutdown

#### Phase 7: Testing & Quality
- [ ] Write unit tests for all layers (target >80% coverage)
- [ ] Write integration tests for API endpoints
- [ ] Setup golangci-lint in CI
- [ ] Add pre-commit hooks for formatting and linting
- [ ] Benchmark critical paths

#### Phase 8: Deployment
- [ ] Create Dockerfile for multi-stage build
- [ ] Setup docker-compose for local development (app, db, redis)
- [ ] Configure CI/CD (GitHub Actions, GitLab CI)
- [ ] Setup air for hot-reload in development
- [ ] Configure production environment variables
- [ ] Setup logging aggregation (optional)

---

### Checklist Template

#### Code Quality Checklist
- [ ] All exported functions have comments
- [ ] No unused imports or variables (go vet)
- [ ] Code formatted with go fmt
- [ ] Error handling is explicit (no ignored errors)
- [ ] No use of panic in library code (only in main/init)
- [ ] Context is passed as first parameter
- [ ] Interfaces are small and focused
- [ ] No global state except configuration

#### Security Checklist
- [ ] Passwords hashed with bcrypt
- [ ] JWT secret loaded from environment, not hardcoded
- [ ] Refresh tokens stored in Redis, not in DB
- [ ] Access token short-lived (≤15min)
- [ ] CORS configured properly (allow only trusted origins)
- [ ] Input validation on all endpoints
- [ ] SQL injection prevented by using parameterized queries (GORM)
- [ ] No sensitive data in logs
- [ ] HTTPS enforced in production
- [ ] Rate limiting on auth endpoints

#### Performance Checklist
- [ ] Database indexes created on frequently queried columns
- [ ] User data cached in Redis
- [ ] Connection pools configured for DB and Redis
- [ ] Use of goroutines for non-blocking tasks (email sending)
- [ ] Avoid N+1 queries (use Preload in GORM)
- [ ] Benchmarks for critical paths

#### Testing Checklist
- [ ] Unit tests cover business logic (usecase)
- [ ] Repository tests with testcontainers or in-memory DB
- [ ] HTTP handler tests with httptest
- [ ] Mock external dependencies (Redis, Mailer)
- [ ] Race condition tests with `-race` flag
- [ ] Test coverage >80%

#### Deployment Checklist
- [ ] Configurable via environment variables
- [ ] Graceful shutdown (wait for existing requests)
- [ ] Health check endpoint
- [ ] Logging to stdout (for container)
- [ ] Docker image built with non-root user
- [ ] Secrets not baked into image
- [ ] Database migration runs automatically on startup (or separate step)
- [ ] Readiness and liveness probes configured
- [ ] Monitoring (Prometheus metrics) exposed

---

### ตัวอย่าง Workflow สมบูรณ์สำหรับ Feature: Register with Email Verification

```mermaid
sequenceDiagram
    participant User
    participant API
    participant DB
    participant Redis
    participant Mailer

    User->>API: POST /register (name, email, password)
    API->>API: Validate request
    API->>DB: Check if email exists
    DB-->>API: not exists
    API->>API: Hash password
    API->>DB: Create user (is_verified=false)
    API->>Redis: Store verification token (TTL 24h)
    API->>Mailer: Send verification email
    Mailer-->>User: Email with link
    API-->>User: 201 Created (wait for verification)
    
    User->>API: GET /verify?token=xxx
    API->>Redis: Get userID by token
    Redis-->>API: userID
    API->>DB: Update user.is_verified=true
    API->>Redis: Delete token
    API-->>User: 200 OK (verified)
```

---

## สรุป

คู่มือเพิ่มเติมนี้ครอบคลุมเครื่องมือและแนวปฏิบัติที่จำเป็นสำหรับการพัฒนาแอปพลิเคชัน Go ในระดับ production โดยเน้นที่:

- **โครงสร้าง Clean Architecture** ที่แยก Delivery, Usecase, Repository
- **การจัดการ JWT** ด้วย access/refresh token เก็บใน Redis
- **การ cache** user data ใน Redis
- **การส่งอีเมล** พร้อมเทมเพลตที่สวยงาม
- **เครื่องมือ** ที่ช่วยเพิ่มประสิทธิภาพการพัฒนา เช่น air สำหรับ hot-reload, cobra สำหรับ CLI, chi สำหรับ routing

การนำไปใช้ควรปรับตามความเหมาะสมของแต่ละโปรเจกต์ โดยให้ความสำคัญกับความปลอดภัย, การทดสอบ, และการบำรุงรักษาในระยะยาว

## บทที่ 46: Go-DDD: การออกแบบบริการด้วย Domain-Driven Design

---

### 46.1 บทนำสู่ Domain-Driven Design

#### 46.1.1 Domain-Driven Design (DDD) คืออะไร?

Domain-Driven Design เป็นแนวทางการออกแบบซอฟต์แวร์ที่เน้นการสร้างโมเดลที่สะท้อนความรู้ความเข้าใจ (domain knowledge) อย่างแท้จริง โดยมีหลักการสำคัญคือการทำให้ซอฟต์แวร์สอดคล้องกับความต้องการทางธุรกิจผ่านการร่วมมือกันระหว่างนักพัฒนาและผู้เชี่ยวชาญในโดเมน (domain experts)

DDD ถูกนิยามโดย Eric Evans ในหนังสือ “Domain-Driven Design: Tackling Complexity in the Heart of Software” (2003) โดยมีเป้าหมายเพื่อจัดการความซับซ้อนของซอฟต์แวร์ในองค์กรขนาดใหญ่

#### 46.1.2 หลักการสำคัญของ DDD

1. **Ubiquitous Language (ภาษาร่วม)**
   - สร้างภาษากลางที่ใช้ร่วมกันระหว่างนักพัฒนาและผู้เชี่ยวชาญโดเมน
   - ใช้ศัพท์เดียวกันในโค้ด, การสนทนา, และเอกสาร
   - ลดความเข้าใจผิดและเพิ่มความสอดคล้อง

2. **Bounded Context (บริบทที่จำกัด)**
   - แบ่งโดเมนขนาดใหญ่ออกเป็นบริทย่อยที่มีขอบเขตชัดเจน
   - แต่ละ Bounded Context มีโมเดลของตัวเองและภาษาร่วมของตัวเอง
   - ลดความซับซ้อนและความขัดแย้งของโมเดล

3. **Entities และ Value Objects**
   - **Entity**: วัตถุที่มีเอกลักษณ์ (identity) และสามารถเปลี่ยนแปลงได้ (mutable) เช่น `User`, `Order`
   - **Value Object**: วัตถุที่ไม่มีเอกลักษณ์ในตัวเอง ถูกกำหนดโดยคุณสมบัติ (immutable) เช่น `Address`, `Money`

4. **Aggregates**
   - กลุ่มของ Entities และ Value Objects ที่ถูกจัดการเป็นหน่วยเดียวกัน
   - มี Aggregate Root (รูท) เป็นตัวควบคุมความสอดคล้องของข้อมูล
   - เช่น `Order` (aggregate root) ประกอบด้วย `OrderItem` (entity) และ `Address` (value object)

5. **Domain Events**
   - เหตุการณ์สำคัญในโดเมนที่เกิดขึ้น เช่น `OrderPlaced`, `UserRegistered`
   - ใช้ในการสื่อสารระหว่าง aggregates หรือระหว่าง bounded contexts

6. **Repositories**
   - ให้ abstraction ในการเข้าถึง aggregate roots
   - ซ่อนรายละเอียดของแหล่งข้อมูล (database, cache)

7. **Domain Services**
   - ใช้สำหรับ logic ที่ไม่เหมาะจะอยู่ใน entity หรือ value object
   - เช่น `TransferService` ที่โอนเงินระหว่างบัญชี

#### 46.1.3 ทำไมต้องใช้ DDD กับ Go?

Go มีคุณสมบัติที่เข้ากันได้ดีกับ DDD:
- **Simple, readable code**: ทำให้โมเดลโดเมนอ่านง่าย
- **Interfaces**: เหมาะสำหรับกำหนด repository และ domain service abstractions
- **Composition**: สนับสนุนการสร้าง aggregates โดยการฝัง structs
- **Package system**: ช่วยจัดระเบียบ bounded contexts ได้ดี
- **Concurrency**: ใช้จัดการ domain events และการสื่อสารระหว่าง contexts

อย่างไรก็ตาม Go ไม่มี inheritance ซึ่ง DDD บางครั้งใช้ แต่เราสามารถใช้ embedding และ interface แทนได้

---

### 46.2 คู่มือการนำ DDD ไปใช้ใน Go

#### 46.2.1 โครงสร้างโปรเจกต์แบบ DDD

```
project/
├── cmd/
│   └── api/
│       └── main.go                 # entry point
├── internal/
│   ├── domain/                     # ชั้นโดเมน (core business logic)
│   │   ├── user/
│   │   │   ├── entity.go           # User entity
│   │   │   ├── value_objects.go    # Email, Password, etc.
│   │   │   ├── repository.go       # interface
│   │   │   ├── service.go          # domain services
│   │   │   └── events.go           # domain events
│   │   ├── order/
│   │   │   └── ...
│   │   └── shared/                 # shared value objects
│   ├── application/                # ชั้นแอปพลิเคชัน (use cases)
│   │   ├── user/
│   │   │   ├── register.go         # use case
│   │   │   ├── login.go
│   │   │   └── dto.go              # input/output DTOs
│   │   └── order/
│   │       └── ...
│   ├── infrastructure/             # ชั้นโครงสร้างพื้นฐาน
│   │   ├── persistence/
│   │   │   ├── gorm/
│   │   │   │   ├── user_repo.go    # implementation
│   │   │   │   └── models.go       # GORM models
│   │   │   └── redis/
│   │   ├── mail/
│   │   └── bus/                    # event bus
│   └── interfaces/                 # ชั้นติดต่อกับภายนอก
│       ├── http/
│       │   ├── handlers/
│       │   ├── middleware/
│       │   └── routes.go
│       └── cli/
├── pkg/                            # reusable packages
└── go.mod
```

#### 46.2.2 การสร้าง Domain Layer

**1. Entities (domain/user/entity.go)**
```go
package user

import (
    "time"
    "github.com/google/uuid"
)

type User struct {
    id          uuid.UUID
    email       Email          // value object
    password    Password       // value object
    name        string
    isVerified  bool
    createdAt   time.Time
    updatedAt   time.Time
}

// Constructor
func NewUser(email, password, name string) (*User, error) {
    emailVO, err := NewEmail(email)
    if err != nil {
        return nil, err
    }
    passwordVO, err := NewPassword(password)
    if err != nil {
        return nil, err
    }
    
    return &User{
        id:         uuid.New(),
        email:      *emailVO,
        password:   *passwordVO,
        name:       name,
        isVerified: false,
        createdAt:  time.Now(),
        updatedAt:  time.Now(),
    }, nil
}

// Getters (exported)
func (u *User) ID() uuid.UUID      { return u.id }
func (u *User) Email() Email       { return u.email }
func (u *User) Name() string       { return u.name }
func (u *User) IsVerified() bool   { return u.isVerified }

// Behavior methods
func (u *User) Verify() {
    u.isVerified = true
    u.updatedAt = time.Now()
    u.addDomainEvent(NewUserVerifiedEvent(u.id))
}

func (u *User) ChangePassword(old, new string) error {
    if err := u.password.Compare(old); err != nil {
        return ErrInvalidPassword
    }
    newPassword, err := NewPassword(new)
    if err != nil {
        return err
    }
    u.password = *newPassword
    u.updatedAt = time.Now()
    return nil
}
```

**2. Value Objects (domain/user/value_objects.go)**
```go
package user

import (
    "regexp"
    "errors"
    "golang.org/x/crypto/bcrypt"
)

type Email struct {
    value string
}

func NewEmail(email string) (*Email, error) {
    // validate email format
    emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
    if !emailRegex.MatchString(email) {
        return nil, errors.New("invalid email format")
    }
    return &Email{value: email}, nil
}

func (e Email) String() string { return e.value }

type Password struct {
    hash string
}

func NewPassword(plain string) (*Password, error) {
    if len(plain) < 8 {
        return nil, errors.New("password must be at least 8 characters")
    }
    hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }
    return &Password{hash: string(hash)}, nil
}

func (p *Password) Compare(plain string) error {
    return bcrypt.CompareHashAndPassword([]byte(p.hash), []byte(plain))
}
```

**3. Repository Interface (domain/user/repository.go)**
```go
package user

import (
    "context"
    "github.com/google/uuid"
)

type Repository interface {
    Save(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id uuid.UUID) (*User, error)
    FindByEmail(ctx context.Context, email Email) (*User, error)
    Delete(ctx context.Context, id uuid.UUID) error
}
```

**4. Domain Events (domain/user/events.go)**
```go
package user

import (
    "time"
    "github.com/google/uuid"
)

type DomainEvent interface {
    OccurredAt() time.Time
}

type UserRegisteredEvent struct {
    UserID    uuid.UUID `json:"user_id"`
    Email     string    `json:"email"`
    Occurred  time.Time `json:"occurred_at"`
}

func (e UserRegisteredEvent) OccurredAt() time.Time { return e.Occurred }

func NewUserRegisteredEvent(userID uuid.UUID, email string) UserRegisteredEvent {
    return UserRegisteredEvent{
        UserID:   userID,
        Email:    email,
        Occurred: time.Now(),
    }
}

// ใน entity สามารถมี slice ของ events และ method addDomainEvent
type User struct {
    // ... fields
    events []DomainEvent
}

func (u *User) addDomainEvent(event DomainEvent) {
    u.events = append(u.events, event)
}

func (u *User) Events() []DomainEvent { return u.events }
func (u *User) ClearEvents()          { u.events = nil }
```

#### 46.2.3 การสร้าง Application Layer

**Use Case (application/user/register.go)**
```go
package user

import (
    "context"
    "your-project/internal/domain/user"
    "your-project/internal/infrastructure/bus"
)

type RegisterUseCase struct {
    userRepo user.Repository
    eventBus *bus.EventBus
}

type RegisterInput struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
    Name     string `json:"name" validate:"required"`
}

type RegisterOutput struct {
    ID    string `json:"id"`
    Email string `json:"email"`
    Name  string `json:"name"`
}

func NewRegisterUseCase(repo user.Repository, eventBus *bus.EventBus) *RegisterUseCase {
    return &RegisterUseCase{
        userRepo: repo,
        eventBus: eventBus,
    }
}

func (uc *RegisterUseCase) Execute(ctx context.Context, input RegisterInput) (*RegisterOutput, error) {
    // 1. ตรวจสอบว่า email ซ้ำหรือไม่
    emailVO, _ := user.NewEmail(input.Email)
    existing, _ := uc.userRepo.FindByEmail(ctx, *emailVO)
    if existing != nil {
        return nil, ErrEmailAlreadyExists
    }
    
    // 2. สร้าง User entity
    newUser, err := user.NewUser(input.Email, input.Password, input.Name)
    if err != nil {
        return nil, err
    }
    
    // 3. บันทึกผ่าน repository
    if err := uc.userRepo.Save(ctx, newUser); err != nil {
        return nil, err
    }
    
    // 4. Dispatch domain events
    for _, event := range newUser.Events() {
        uc.eventBus.Publish(event)
    }
    newUser.ClearEvents()
    
    // 5. ส่ง output
    return &RegisterOutput{
        ID:    newUser.ID().String(),
        Email: newUser.Email().String(),
        Name:  newUser.Name(),
    }, nil
}
```

#### 46.2.4 การสร้าง Infrastructure Layer

**Repository Implementation (infrastructure/persistence/gorm/user_repo.go)**
```go
package gorm

import (
    "context"
    "gorm.io/gorm"
    "your-project/internal/domain/user"
)

type UserModel struct {
    ID         string `gorm:"primaryKey"`
    Email      string `gorm:"uniqueIndex"`
    Password   string
    Name       string
    IsVerified bool
    CreatedAt  time.Time
    UpdatedAt  time.Time
}

type UserRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) user.Repository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Save(ctx context.Context, u *user.User) error {
    model := &UserModel{
        ID:         u.ID().String(),
        Email:      u.Email().String(),
        Password:   u.PasswordHash(), // need to expose getter
        Name:       u.Name(),
        IsVerified: u.IsVerified(),
        CreatedAt:  u.CreatedAt(),
        UpdatedAt:  u.UpdatedAt(),
    }
    return r.db.WithContext(ctx).Save(model).Error
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
    var model UserModel
    err := r.db.WithContext(ctx).Where("id = ?", id.String()).First(&model).Error
    if err != nil {
        return nil, err
    }
    return model.ToDomain()
}

func (m *UserModel) ToDomain() (*user.User, error) {
    // ใช้ private constructor หรือ method เพื่อสร้าง domain object
    // อาจต้อง expose factory method ใน domain layer
}
```

**Event Bus (infrastructure/bus/event_bus.go)**
```go
package bus

import (
    "context"
    "sync"
    "your-project/internal/domain/user"
)

type EventHandler func(context.Context, user.DomainEvent) error

type EventBus struct {
    handlers map[string][]EventHandler
    mu       sync.RWMutex
}

func NewEventBus() *EventBus {
    return &EventBus{
        handlers: make(map[string][]EventHandler),
    }
}

func (b *EventBus) Subscribe(eventName string, handler EventHandler) {
    b.mu.Lock()
    defer b.mu.Unlock()
    b.handlers[eventName] = append(b.handlers[eventName], handler)
}

func (b *EventBus) Publish(event user.DomainEvent) {
    b.mu.RLock()
    handlers := b.handlers[eventName(event)]
    b.mu.RUnlock()
    
    for _, h := range handlers {
        go h(context.Background(), event) // async
    }
}
```

#### 46.2.5 การสร้าง Interfaces Layer

**HTTP Handler (interfaces/http/handlers/user_handler.go)**
```go
package handlers

import (
    "encoding/json"
    "net/http"
    "your-project/internal/application/user"
)

type UserHandler struct {
    registerUseCase *user.RegisterUseCase
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
    var input user.RegisterInput
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    output, err := h.registerUseCase.Execute(r.Context(), input)
    if err != nil {
        // map domain errors to HTTP status
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(output)
}
```

---

### 46.3 การออกแบบ Workflow ด้วย DDD

#### 46.3.1 Event Storming

Event Storming เป็นกระบวนการในการทำความเข้าใจโดเมนผ่านการระบุ events, commands, aggregates, และ bounded contexts ร่วมกันระหว่างทีมพัฒนาและผู้เชี่ยวชาญโดเมน

**ขั้นตอน Event Storming:**
1. **ระบุ Domain Events** (สีส้ม)
   - เหตุการณ์ที่เกิดขึ้นในระบบ เช่น `UserRegistered`, `OrderPlaced`, `PaymentReceived`

2. **ระบุ Commands** (สีน้ำเงิน)
   - การกระทำที่ทำให้เกิด events เช่น `RegisterUser`, `PlaceOrder`

3. **ระบุ Aggregates** (สีเหลือง)
   - วัตถุหลักที่รับ commands และสร้าง events

4. **ระบุ External Systems** (สีชมพู)
   - ระบบภายนอกที่เกี่ยวข้อง

5. **ระบุ Bounded Contexts**
   - จัดกลุ่ม events, commands, aggregates ที่เกี่ยวข้องกัน

**ตัวอย่าง Event Storming สำหรับระบบ e-commerce:**

```
[Bounded Context: User Management]
- RegisterUser (command) → UserRegistered (event)
- VerifyEmail → UserEmailVerified
- Login → UserLoggedIn

[Bounded Context: Order Management]
- PlaceOrder → OrderPlaced
- CancelOrder → OrderCancelled
- ShipOrder → OrderShipped

[Bounded Context: Payment]
- ProcessPayment → PaymentProcessed / PaymentFailed
```

#### 46.3.2 การแบ่ง Bounded Contexts

หลักการแบ่ง Bounded Context:
- แต่ละ context มีโมเดลและภาษาของตัวเอง
- การสื่อสารระหว่าง contexts ใช้ APIs หรือ messages
- หลีกเลี่ยงการแชร์โมเดลโดยตรง

**รูปแบบการสื่อสารระหว่าง Bounded Contexts:**
1. **Conformist**: context หนึ่งยอมรับโมเดลของอีก context
2. **Anti-Corruption Layer (ACL)**: แปลงโมเดลระหว่าง contexts ป้องกัน contamination
3. **Open Host Service**: เปิด API ให้ context อื่นเรียกใช้
4. **Published Language**: ใช้ภาษากลาง (เช่น JSON schema) ในการสื่อสาร

#### 46.3.3 การออกแบบ Aggregates

Aggregate เป็นหน่วยความสอดคล้อง (transactional consistency) ต้องออกแบบให้มีขอบเขตที่เหมาะสม

**หลักการออกแบบ Aggregate:**
1. **Invariants ต้องอยู่ภายใน aggregate เดียว** - กฎทางธุรกิจที่ต้องสอดคล้องกันตลอดเวลาควรอยู่ภายใน aggregate เดียว
2. **Aggregate root เป็นตัวเดียวที่เข้าถึงได้จากภายนอก**
3. **ใช้ ID references แทนการอ้างอิงโดยตรง** - เพื่อลด coupling
4. **ขนาดเล็ก** - aggregate ที่ใหญ่เกินไปจะเกิดปัญหา performance และ concurrency

**ตัวอย่าง Aggregate Design:**

```go
// Order Aggregate
type Order struct {
    id          OrderID
    customerID  CustomerID      // reference to another aggregate
    items       []OrderItem
    status      OrderStatus
    total       Money
    // ...
}

func (o *Order) AddItem(productID ProductID, quantity int, price Money) error {
    if o.status != Draft {
        return ErrOrderNotEditable
    }
    // business rule: cannot add more than 10 items
    if len(o.items) >= 10 {
        return ErrTooManyItems
    }
    o.items = append(o.items, NewOrderItem(productID, quantity, price))
    o.recalculateTotal()
    return nil
}
```

#### 46.3.4 การสื่อสารระหว่าง Bounded Contexts

**1. Synchronous (REST/gRPC):**
- ใช้เมื่อต้องการ response ทันที
- ต้องจัดการ timeout, retry, circuit breaker

**2. Asynchronous (Message Broker):**
- ใช้เมื่อต้องการ decoupling และ event-driven
- เช่น RabbitMQ, Kafka, NATS

**ตัวอย่างการสื่อสารด้วย Domain Events:**

```go
// ใน context Order Management
// เมื่อ OrderPlaced เกิดขึ้น จะส่ง event ไปยัง Payment Context
type OrderPlacedEvent struct {
    OrderID    string
    CustomerID string
    Total      Money
}

// Payment Context subscribe event นี้
func (h *PaymentHandler) HandleOrderPlaced(event OrderPlacedEvent) {
    // สร้าง Payment transaction
    payment := NewPayment(event.OrderID, event.Total)
    // ...
}
```

---

### 46.4 TASK LIST Template (สำหรับโปรเจกต์ DDD)

#### Phase 1: Domain Discovery
- [ ] จัดประชุม Event Storming กับ domain experts
- [ ] ระบุ Domain Events ทั้งหมด
- [ ] ระบุ Commands และ Aggregates
- [ ] แบ่ง Bounded Contexts
- [ ] กำหนด Ubiquitous Language สำหรับแต่ละ context
- [ ] สร้างเอกสาร Context Map

#### Phase 2: Domain Modeling
- [ ] สำหรับแต่ละ Bounded Context:
  - [ ] กำหนด Entities และ Value Objects
  - [ ] กำหนด Aggregates และ Aggregate Roots
  - [ ] กำหนด Domain Events
  - [ ] กำหนด Domain Services (ถ้ามี)
  - [ ] เขียน unit tests สำหรับ domain logic
  - [ ] Implement domain models ใน Go

#### Phase 3: Application Layer
- [ ] กำหนด Use Cases (application services)
- [ ] สร้าง DTOs (Input/Output)
- [ ] Implement Use Cases (orchestrate domain objects)
- [ ] ใช้ repository interfaces
- [ ] Dispatch domain events
- [ ] เขียน unit tests สำหรับ use cases (mock repositories)

#### Phase 4: Infrastructure Layer
- [ ] Implement Repositories (GORM, MongoDB, etc.)
- [ ] Implement Event Bus (in-memory, message broker)
- [ ] Implement external service clients (email, payment, etc.)
- [ ] Setup database migrations
- [ ] Setup caching (Redis)
- [ ] เขียน integration tests

#### Phase 5: Interfaces Layer
- [ ] Implement HTTP handlers (หรือ gRPC)
- [ ] Implement middleware (auth, logging, etc.)
- [ ] Map routes
- [ ] Implement CLI commands (ถ้ามี)
- [ ] เขียน contract tests / API tests

#### Phase 6: Cross-Cutting Concerns
- [ ] Dependency Injection (wire, fx, หรือ manual)
- [ ] Configuration management (viper)
- [ ] Logging (zap)
- [ ] Error handling strategy
- [ ] Metrics and monitoring
- [ ] Graceful shutdown

#### Phase 7: Testing & Quality
- [ ] Unit tests (domain + application) > 80% coverage
- [ ] Integration tests (infrastructure)
- [ ] End-to-end tests (critical paths)
- [ ] Performance benchmarks
- [ ] Linting and static analysis (golangci-lint)

#### Phase 8: Deployment & Operations
- [ ] Dockerfile (multi-stage)
- [ ] docker-compose for local development
- [ ] CI/CD pipeline (GitHub Actions, GitLab CI)
- [ ] Health check endpoints
- [ ] Observability (OpenTelemetry, Prometheus)

---

### 46.5 CHECKLIST Template (สำหรับตรวจสอบการออกแบบ DDD)

#### Ubiquitous Language
- [ ] ชื่อในโค้ดตรงกับภาษาที่ domain experts ใช้
- [ ] ไม่มีศัพท์เทคนิคเกินจำเป็นใน domain layer
- [ ] ทีมเข้าใจและใช้ภาษาร่วมกันอย่างสม่ำเสมอ

#### Bounded Contexts
- [ ] แต่ละ context มีขอบเขตชัดเจน
- [ ] ไม่มีการแชร์โมเดลระหว่าง contexts โดยตรง
- [ ] การสื่อสารระหว่าง contexts เป็นไปตามรูปแบบที่กำหนด (ACL, Open Host Service, etc.)

#### Entities & Value Objects
- [ ] Entities มี identity เฉพาะ (ไม่ใช้ชื่อหรือ attributes เป็น identity)
- [ ] Value Objects เป็น immutable
- [ ] Value Objects มีการ validate ใน constructor
- [ ] Entities มี behavior (methods) ที่เกี่ยวข้องกับ business logic
- [ ] ไม่มี getter/setter ที่ไม่จำเป็น (anemic model)

#### Aggregates
- [ ] Aggregate root เป็นจุดเดียวที่เข้าถึงได้จากภายนอก
- [ ] Invariants ถูกตรวจสอบภายใน aggregate
- [ ] ขนาด aggregate ไม่ใหญ่เกินไป (avoid bloated aggregates)
- [ ] การอ้างอิงถึง aggregate อื่นใช้ ID แทน direct reference
- [ ] Transactions ครอบคลุม aggregate เดียว (ไม่ครอบคลุมหลาย aggregate)

#### Domain Events
- [ ] Events ถูกตั้งชื่อในอดีตกาล (UserRegistered, OrderShipped)
- [ ] Events มีข้อมูลที่จำเป็นสำหรับการประมวลผล
- [ ] Events ถูก dispatch ผ่าน event bus
- [ ] Handlers ทำงานแบบ asynchronous (ถ้าไม่ต้องการ transactional consistency)

#### Repositories
- [ ] Repository interface อยู่ใน domain layer
- [ ] Implementation อยู่ใน infrastructure layer
- [ ] Repository methods ใช้ domain objects เป็นพารามิเตอร์และคืนค่า
- [ ] Repository ซ่อนรายละเอียดของ persistence

#### Domain Services
- [ ] ใช้เมื่อ logic ไม่เหมาะสมใน entity/value object
- [ ] Domain services มีพฤติกรรมที่เกี่ยวข้องกับโดเมน (ไม่ใช่ technical)
- [ ] รับ domain objects และคืนค่า domain objects

#### Application Layer (Use Cases)
- [ ] Use cases มีหน้าที่ orchestrate เท่านั้น (ไม่ใส่ business logic)
- [ ] Use cases มีขนาดเล็กและทำสิ่งเดียว
- [ ] Input validation ผ่าน DTOs
- [ ] ใช้ repository interfaces
- [ ] Dispatch domain events หลังจาก transaction commit

#### Infrastructure Layer
- [ ] Implementation ถูกแยกออกจาก domain ด้วย interfaces
- [ ] ใช้ dependency injection
- [ ] การจัดการ transaction ทำใน infrastructure หรือ application layer? (ต้องสอดคล้อง)

#### Testing
- [ ] Domain layer test โดยไม่พึ่ง infrastructure
- [ ] Application layer test ด้วย mock repositories
- [ ] Integration test สำหรับ infrastructure
- [ ] มี test สำหรับ domain invariants
- [ ] มี test สำหรับ domain events

#### Performance & Scalability
- [ ] Aggregate boundaries ถูกออกแบบให้รองรับการทำงานพร้อมกัน (concurrency)
- [ ] การสื่อสารระหว่าง contexts ไม่สร้าง latency สูงเกินไป
- [ ] ใช้ caching ตามความเหมาะสม

#### Documentation
- [ ] Context Map ได้รับการบันทึก
- [ ] เหตุผลในการออกแบบ aggregates และ bounded contexts มีเอกสาร
- [ ] Domain events มีการบันทึก schema

---

## สรุป

Domain-Driven Design ช่วยให้เราสร้างซอฟต์แวร์ที่สะท้อนความต้องการทางธุรกิจอย่างแท้จริง โดยเฉพาะในระบบที่มีความซับซ้อนสูง Go มีคุณสมบัติที่เหมาะสมกับการนำ DDD ไปใช้ ทั้งในเรื่องความเรียบง่าย, interfaces, และ concurrency

การนำ DDD ไปใช้ต้องอาศัยความร่วมมือระหว่างทีมพัฒนาและผู้เชี่ยวชาญโดเมน รวมถึงต้องใช้เวลาในการเรียนรู้และปรับใช้ อย่างไรก็ตาม ผลลัพธ์ที่ได้คือระบบที่ยืดหยุ่น, บำรุงรักษาง่าย, และตอบสนองต่อการเปลี่ยนแปลงทางธุรกิจได้ดี

คู่มือนี้เป็นจุดเริ่มต้นสำหรับการออกแบบบริการด้วย Go และ DDD โดยสามารถปรับใช้ตามความเหมาะสมของแต่ละโปรเจกต์

## บทที่ 47: Blueprint สำหรับโปรเจกต์ Go ระดับ Production

---

### 47.1 บทนำ

การเริ่มต้นโปรเจกต์ Go ใหม่ให้มีโครงสร้างที่แข็งแรง พร้อมขยาย และบำรุงรักษาง่าย เป็นความท้าทายที่นักพัฒนาทุกคนต้องเผชิญ โครงสร้างโปรเจกต์ที่นำเสนอในบทนี้เป็น blueprint ที่ผ่านการพิสูจน์แล้วจากโปรเจกต์จริง มีการออกแบบตามหลัก **Clean Architecture** พร้อมด้วยฟีเจอร์ครบวงจรสำหรับการพัฒนาแอปพลิเคชันระดับ production

โครงสร้างนี้แยกชั้นการทำงานอย่างชัดเจน:

- **Core (domain)** – จัดเก็บโมเดลธุรกิจหลัก (entity, repository interface, service interface)
- **Platform (infrastructure)** – จัดการการเชื่อมต่อฐานข้อมูล, Redis, logging, message queue
- **Transport (delivery)** – รับส่งข้อมูลผ่าน HTTP, middleware, response formatter
- **Apps** – จุดรวม dependencies และกำหนด routes

นอกจากนี้ยังมีฟีเจอร์พร้อมใช้ที่ช่วยลดเวลาการพัฒนา เช่น:

- Authentication (JWT access/refresh token)
- Rate limiting (IP-based)
- Health check (พร้อม dependency monitoring)
- Caching (Redis)
- Message queue (Redis pub/sub + worker pool)
- Transaction management
- Structured logging
- Graceful shutdown

โครงสร้างนี้เหมาะสำหรับโปรเจกต์ที่ต้องมีขนาดใหญ่ ทำงานพร้อมกันสูง และต้องการความน่าเชื่อถือ

---

### 47.2 คู่มือการใช้งานโครงสร้างโปรเจกต์

#### 47.2.1 โครงสร้างโฟลเดอร์หลัก

```
project-root/
├── cmd/                 # จุดเริ่มต้นของโปรแกรม (executable)
├── internal/            # โค้ดเฉพาะของแอปพลิเคชัน (ไม่ถูก import จากภายนอก)
├── pkg/                 # โค้ดที่นำกลับไปใช้ได้ (reusable) ในโปรเจกต์อื่น
├── api/                 # ไฟล์ที่เกี่ยวข้องกับ API (เช่น Swagger)
├── configs/             # ไฟล์ configuration
├── deploy/              # Docker, Kubernetes files
├── migrations/          # SQL migration files
├── scripts/             # สคริปต์ช่วยพัฒนา
└── test/                # integration tests
```

#### 47.2.2 ชั้นการทำงานภายใน `internal/`

**apps/** – จุดรวม dependencies และกำหนด routing

- `bootstrap/injection/` – ใช้สำหรับ wire dependencies (repository, service, handler) ด้วย dependency injection แบบ manual
- `router/v1/` – กำหนด routes สำหรับ API version 1 แยก public/protected

**core/** – ชั้นโดเมนของแอปพลิเคชัน (แบ่งตาม domain module)

แต่ละ module (เช่น `auth`, `user`, `iot`) มีโครงสร้างย่อย:

- `entity/` – entity หลัก (struct พร้อม behavior) และ value objects
- `repository/` – interface สำหรับ repository
- `service/` – interface สำหรับ service และการ implement business logic
- `dto/` – data transfer objects (request/response)
- `model/` – model ที่ใช้กับฐานข้อมูล (optional)
- `handler/` – HTTP handlers (รับ request, เรียก service, ส่ง response)
- `routes.go` – ลงทะเบียน routes ของ module นี้

**platform/** – ชั้น infrastructure

- `config/` – โหลด configuration (viper)
- `db/` – การเชื่อมต่อ PostgreSQL, Redis
- `cache/` – ฟังก์ชันช่วยสำหรับ Redis cache
- `queue/` – message queue (Redis pub/sub) + worker pool + dead letter queue
- `logger/` – structured logger (slog) + middleware

**transport/** – ชั้น delivery

- `middleware/` – middleware ต่างๆ (CORS, rate limit, auth, logging, recovery, security headers)
- `httpx/` – utilities สำหรับ response, request, validation
- `utils/` – helper functions (context, pagination)

#### 47.2.3 การทำงานของ dependency injection

ใน `internal/apps/app/bootstrap/injection/` มีไฟล์:

- `dependencies.go` – สร้าง dependencies ทั่วไป (config, db, redis, logger)
- `repositories.go` – สร้าง repository implementations (ส่งต่อ db)
- `services.go` – สร้าง service implementations (ส่งต่อ repository)
- `handlers.go` – สร้าง handlers (ส่งต่อ service)

ทุกอย่างถูกสร้างใน `initDependencies()` และถูกส่งไปยัง router

#### 47.2.4 การเพิ่มโมดูลใหม่ (Feature)

สมมติต้องการเพิ่มโมดูล `product`:

1. สร้างโครงสร้างใน `internal/core/product/`:
   - `entity/product.go` – entity และ behavior
   - `repository/repository.go` – interface สำหรับ repository
   - `service/service.go` – interface สำหรับ service + implementation
   - `dto/product_dto.go` – request/response structs
   - `handler/product_handler.go` – HTTP handlers
   - `routes.go` – routes สำหรับ product

2. เพิ่ม repository implementation ใน `internal/platform/db/postgres/` (หรือถ้ามีหลาย module ให้แยกไฟล์)
3. เพิ่ม service implementation ใน `internal/core/product/service/service_impl.go`
4. เพิ่ม handler ใน `internal/core/product/handler/`
5. ลงทะเบียน dependencies ใน `internal/apps/app/bootstrap/injection/`:
   - `repositories.go`: `productRepo := postgres.NewProductRepository(db)`
   - `services.go`: `productService := product.NewProductService(productRepo)`
   - `handlers.go`: `productHandler := product.NewProductHandler(productService)`
6. เพิ่ม routes ใน `internal/apps/app/router/v1/protected_routes.go` (หรือ public_routes.go)
7. อย่าลืมเพิ่ม migrations ถ้ามีการเปลี่ยนแปลง schema

---

### 47.3 การออกแบบ Workflow สำหรับการพัฒนา

#### 47.3.1 ขั้นตอนการพัฒนา Feature ใหม่

1. **Analyze** – ทำความเข้าใจ requirement, ระบุ domain models, use cases
2. **Design** – ออกแบบ entities, value objects, repository interface, service interface
3. **Implement Domain** – เขียน entity, repository interface, service interface ใน `internal/core/<module>`
4. **Implement Infrastructure** – เขียน repository implementation (GORM), cache, queue (ถ้าจำเป็น)
5. **Implement Service** – เขียน business logic ใน service implementation
6. **Implement Handler** – เขียน HTTP handlers, ตรวจสอบ input validation, ใช้ service
7. **Register Routes** – ลงทะเบียน routes ใน router
8. **Test** – เขียน unit tests (domain, service), integration tests (handler)
9. **Document** – อัปเดต API docs (Swagger) ถ้ามี

#### 47.3.2 การใช้เครื่องมือช่วยพัฒนา

- **Air** – hot-reload (`air` ใน root) ช่วยให้เห็นผลทันที
- **Docker Compose** – ใช้ `docker-compose up` ใน `deploy/docker/` เพื่อรัน Postgres, Redis
- **Migrate** – ใช้ `scripts/migrate.sh` เพื่อรัน migrations
- **Seed** – ใช้ `scripts/seed.sh` เพื่อเติมข้อมูลตัวอย่าง

#### 47.3.3 การจัดการ configuration

- ไฟล์ใน `configs/` (เช่น `config.yaml`) จะถูกโหลดด้วย viper
- Environment variables สามารถ override ค่าได้ (ตั้งค่าใน docker-compose หรือ k8s)

#### 47.3.4 การ deploy

- ใช้ Dockerfile ใน `deploy/docker/Dockerfile` (multi-stage)
- ใช้ Kubernetes manifests ใน `deploy/k8s/` สำหรับ production

---

### 47.4 TASK LIST Template (สำหรับการพัฒนา Feature ใหม่)

#### Phase 1: Domain Design
- [ ] ระบุ domain model (entity, value objects)
- [ ] กำหนด invariants (business rules)
- [ ] ระบุ use cases (methods in service)
- [ ] กำหนด events (ถ้ามี)
- [ ] ออกแบบ repository interface (methods)
- [ ] ออกแบบ DTOs (request/response)

#### Phase 2: Implementation
- [ ] สร้าง entity struct และ behavior methods
- [ ] สร้าง repository interface
- [ ] สร้าง service interface
- [ ] สร้าง DTO structs
- [ ] เขียน unit tests สำหรับ entity
- [ ] เขียน unit tests สำหรับ service (mock repository)

#### Phase 3: Infrastructure
- [ ] สร้าง repository implementation (GORM)
- [ ] สร้าง migration file (ถ้ามี)
- [ ] ตั้งค่า Redis cache (ถ้าจำเป็น)
- [ ] ตั้งค่า message queue (ถ้าจำเป็น)
- [ ] ทดสอบ repository ด้วย integration test

#### Phase 4: Delivery
- [ ] สร้าง HTTP handlers
- [ ] เพิ่ม input validation (go-playground/validator)
- [ ] สร้าง routes
- [ ] ลงทะเบียน dependencies ใน injection
- [ ] ทดสอบ handler ด้วย httptest

#### Phase 5: Integration & Documentation
- [ ] ทดสอบ end-to-end ด้วย curl/Postman
- [ ] อัปเดต Swagger docs (ถ้ามี)
- [ ] อัปเดต README (ถ้าจำเป็น)
- [ ] รัน linter (`golangci-lint run`) และแก้ไข warnings

#### Phase 6: Review & Deploy
- [ ] Code review
- [ ] ตรวจสอบ performance (ถ้ามีการ query มาก)
- [ ] รัน test coverage (`go test -cover`), ควร > 80%
- [ ] รัน race detector (`go test -race`)
- [ ] Deploy to staging
- [ ] ทดสอบใน staging
- [ ] Deploy to production

---

### 47.5 CHECKLIST Template (สำหรับการตรวจสอบคุณภาพ)

#### Domain Layer (core)
- [ ] Entity ไม่มี setter ที่ไม่จำเป็น (immutable pattern)
- [ ] Entity มี behavior ที่เกี่ยวข้อง (ไม่ใช่แค่ data holder)
- [ ] Value objects มี constructor และ validation
- [ ] Repository interface ใช้ domain objects (ไม่ใช้ model ของ ORM)
- [ ] Service interface มีชื่อชัดเจน (เช่น `UserService`, `AuthService`)
- [ ] ไม่มี business logic ซ้ำซ้อน (DRY)
- [ ] Unit test coverage > 80%

#### Infrastructure Layer (platform)
- [ ] Repository implementation ใช้ interface ที่กำหนด
- [ ] การจัดการ connection pool ถูกต้อง (max open, max idle)
- [ ] การใช้ transaction ถูกต้อง (ถ้ามี)
- [ ] Redis cache มี TTL และการ invalidate ที่เหมาะสม
- [ ] Message queue มี error handling และ dead letter queue
- [ ] Logging มี structured fields (ไม่ใช่ string concatenation)

#### Transport Layer (delivery)
- [ ] HTTP handlers มีขนาดเล็ก (ไม่เกิน 50 บรรทัด)
- [ ] Input validation ครอบคลุมทุก field
- [ ] Response มีรูปแบบมาตรฐาน (ผ่าน httpx.Response)
- [ ] มีการจัดการ error ที่เหมาะสม (map domain error -> HTTP status)
- [ ] Middleware ทำงานตามลำดับที่ถูกต้อง
- [ ] มีการทดสอบ integration (httptest)

#### Security
- [ ] JWT secret ถูกเก็บใน environment (ไม่ hardcode)
- [ ] Rate limiting เปิดใช้งานสำหรับ endpoints ที่ sensitive
- [ ] Security headers (CSP, HSTS, etc.) ถูก set ผ่าน middleware
- [ ] SQL injection ป้องกันโดย GORM (parameterized query)
- [ ] No sensitive data in logs (password, token)
- [ ] CORS กำหนด allow origins เฉพาะที่จำเป็น

#### Performance
- [ ] มีการ cache สำหรับข้อมูลที่อ่านบ่อย
- [ ] มีการใช้ connection pooling สำหรับ DB และ Redis
- [ ] มีการใช้ bulk operation แทน loop insert/update
- [ ] มีการตั้ง timeout สำหรับ HTTP client และ database queries
- [ ] มีการใช้ goroutine อย่างเหมาะสม (ไม่สร้าง goroutine มากเกินไป)

#### Reliability
- [ ] มี health check endpoints (`/health`, `/health/detailed`, `/ready`, `/live`)
- [ ] Graceful shutdown ทำงาน (รอ pending requests)
- [ ] มีการ retry สำหรับ transient failures (เช่น network)
- [ ] มี dead letter queue สำหรับ messages ที่ประมวลผลไม่สำเร็จ
- [ ] มีการ log panic และ recover ด้วย middleware

#### Maintainability
- [ ] โค้ดผ่าน `go fmt` และ `go vet`
- [ ] ไม่มี unused imports, variables
- [ ] มีการแบ่งแพคเกจตาม domain (ไม่ใช่ technical layers)
- [ ] มีการใช้ dependency injection (ไม่ใช้ global variables)
- [ ] มีเอกสารใน `docs/` สำหรับ architecture, deployment
- [ ] README มี quick start guide

#### Deployment
- [ ] Dockerfile เป็น multi-stage (ลดขนาด image)
- [ ] docker-compose ใช้งานได้ (development)
- [ ] Kubernetes manifests มี liveness/readiness probes
- [ ] Environment variables ถูกกำหนดใน deployment (ไม่ฝังใน image)
- [ ] Migration รันอัตโนมัติหรือเป็นขั้นตอนแยก (ก่อน deploy)

---

## สรุป

โครงสร้างโปรเจกต์ที่นำเสนอในบทนี้เป็น blueprint ที่รวบรวมประสบการณ์จากโปรเจกต์จริง ช่วยให้นักพัฒนาสามารถเริ่มต้นโปรเจกต์ใหม่ได้อย่างรวดเร็วโดยไม่ต้องเสียเวลากับการออกแบบโครงสร้างพื้นฐาน และยังมีฟีเจอร์พร้อมใช้ที่ช่วยให้แอปพลิเคชันมีความปลอดภัย, มีประสิทธิภาพ, และพร้อมสำหรับการขยายขนาด

การนำ blueprint นี้ไปใช้ควรปรับให้เข้ากับความต้องการของแต่ละโปรเจกต์ โดยยึดหลักการ:

- **แยกความรับผิดชอบ** (Separation of Concerns)
- **ทดสอบง่าย** (Testability)
- **บำรุงรักษาง่าย** (Maintainability)
- **ขยายได้** (Scalability)

การใช้ task list และ checklist ที่กำหนดจะช่วยให้ทีมพัฒนามีมาตรฐานเดียวกันและลดข้อผิดพลาดที่อาจเกิดขึ้น

## บทที่ 48: GORM – ORM ทรงพลังสำหรับ Go

### 48.1 บทนำ

การจัดการฐานข้อมูลเป็นหัวใจสำคัญของแอปพลิเคชันส่วนใหญ่ การเขียน SQL ด้วยมือนั้นให้ความยืดหยุ่นสูง แต่ก็อาจทำให้โค้ดยุ่งเหยิง เสี่ยงต่อข้อผิดพลาด และต้องดูแลการแปลงข้อมูลระหว่าง Go struct กับตารางด้วยตนเอง GORM (Go Object Relational Mapping) เป็น ORM ที่ได้รับความนิยมสูงสุดในภาษา Go ช่วยลดความซับซ้อนเหล่านี้ ด้วยการแมป struct กับตารางโดยอัตโนมัติ พร้อมฟังก์ชัน CRUD ที่ใช้งานง่าย การจัดการ transaction แบบมีระบบ และยังมีฟีเจอร์ขั้นสูงอย่าง query cache, queue processor ที่ช่วยให้แอปพลิเคชันทำงานได้เร็วและมั่นคงยิ่งขึ้น

บทนี้จะพาคุณสำรวจ GORM ตั้งแต่พื้นฐาน CRUD ไปจนถึงการออกแบบ session factory, การใช้ query cache, การประมวลผลคิว, และการจัดการ transaction แบบ rollback/commit อย่างเป็นระบบ พร้อมตัวอย่างโค้ดและ workflow ที่นำไปใช้ได้จริง

---

### 48.2 บทนิยาม

#### 48.2.1 CRUD
**CRUD** คือชุดการดำเนินการพื้นฐานสี่ประการในการจัดการข้อมูล:
- **C**reate – เพิ่มข้อมูลใหม่
- **R**ead – อ่าน/ดึงข้อมูล
- **U**pdate – แก้ไขข้อมูลที่มีอยู่
- **D**elete – ลบข้อมูล

ใน GORM การดำเนินการเหล่านี้ทำผ่าน `*gorm.DB` object และ method ต่างๆ เช่น `Create`, `First`, `Find`, `Save`, `Update`, `Delete`

#### 48.2.2 ORM
**ORM (Object-Relational Mapping)** คือเทคนิคที่ช่วยแปลงข้อมูลระหว่างฐานข้อมูลเชิงสัมพันธ์กับโครงสร้างข้อมูลเชิงวัตถุ (object) ในภาษาโปรแกรม โดย GORM จะทำหน้าที่:
- แปลง Go struct ให้เป็น SQL และแปลงผลลัพธ์ SQL กลับเป็น Go struct
- จัดการความสัมพันธ์ระหว่างตาราง (has one, has many, belongs to, many to many)
- รองรับ migration, transaction, hooks, และ caching

#### 48.2.3 GORM
**GORM** คือ ORM library สำหรับ Go ที่ออกแบบมาให้ใช้งานง่าย กระชับ และมีประสิทธิภาพสูง โดยมีคุณสมบัติเด่น:
- Full-featured ORM (CRUD, associations, hooks, transactions)
- ใช้งานได้กับฐานข้อมูลหลายชนิด (MySQL, PostgreSQL, SQLite, SQL Server, ClickHouse)
- มี chainable API ที่อ่านง่าย
- รองรับ eager loading (Preload)
- มี auto-migration
- สามารถขยายฟังก์ชันผ่าน plugins

#### 48.2.4 Transaction
**Transaction (ทรานแซคชัน)** คือกลุ่มของคำสั่ง SQL ที่ต้องสำเร็จทั้งหมดหรือไม่ทำเลย (All-or-Nothing) โดยมีคุณสมบัติ ACID:
- **Atomicity**: งานทั้งหมดสำเร็จหรือล้มเหลวพร้อมกัน
- **Consistency**: ข้อมูลคงความถูกต้องตามกฎธุรกิจ
- **Isolation**: ทรานแซคชันที่ทำงานพร้อมกันไม่รบกวนกัน
- **Durability**: เมื่อ commit ข้อมูลจะถูกบันทึกถาวร

GORM รองรับ transaction ผ่าน `db.Transaction()` หรือ `db.Begin()` / `Commit()` / `Rollback()`

#### 48.2.5 Cache
**Cache (แคช)** คือการเก็บข้อมูลชั่วคราวเพื่อลดการเข้าถึงฐานข้อมูลซ้ำๆ เพิ่มความเร็วในการตอบสนอง GORM มีกลไก:
- **First-level cache**: ภายใน session เดียว (ไม่เปิดเผยให้ผู้ใช้จัดการ)
- **Second-level cache**: ระดับ application, ต้องใช้ plugin เช่น `gorm-cache`
- **Query cache**: เก็บผลลัพธ์ของคำสั่ง `Find`, `First` เพื่อใช้ซ้ำ

#### 48.2.6 Queue Processor
**Queue Processor** เป็นกลไกในการประมวลผลงานที่อาจใช้เวลานานหรือเกิดบ่อยแบบ asynchronous โดย GORM สามารถทำงานร่วมกับระบบคิว (เช่น Redis, RabbitMQ) เพื่อแยกการรับ request ออกจากการทำงานหนัก ทำให้ระบบตอบสนองได้เร็วขึ้น

---

### 48.3 หัวข้อหลัก

1. **GORM CRUD** – การสร้าง, อ่าน, อัปเดต, ลบข้อมูลพื้นฐาน
2. **SessionFactory** – การสร้างและจัดการ `*gorm.DB` session สำหรับแยก context
3. **Query Cache** – การแคชผลลัพธ์ query เพื่อลดภาระฐานข้อมูล
4. **Cache Query Queue Processor** – การใช้คิวเพื่อประมวลผล query แบบ asynchronous และแคช
5. **Queue Processor Function** – การออกแบบฟังก์ชันสำหรับประมวลผลคิว
6. **SQL Queue Translate Rollback Commit** – การแปลง SQL เป็นคิวและจัดการ transaction ในระบบคิว

---

### 48.4 คู่มือการใช้งาน GORM

#### 48.4.1 การติดตั้งและการตั้งค่าพื้นฐาน

```go
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres // หรือ driver อื่นตามที่ใช้
```

**ตัวอย่างการเชื่อมต่อ PostgreSQL**

```go
package main

import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "log"
)

func main() {
    dsn := "host=localhost user=gorm password=gorm dbname=gorm port=5432 sslmode=disable TimeZone=Asia/Bangkok"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("failed to connect database")
    }

    // ใช้ db ต่อไป
}
```

#### 48.4.2 กำหนด Model (Entity)

```go
type User struct {
    ID        uint           `gorm:"primaryKey"`
    Name      string         `gorm:"size:100;not null"`
    Email     string         `gorm:"uniqueIndex;size:100;not null"`
    Age       int            `gorm:"default:0"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}
```

#### 48.4.3 CRUD Operations

**Create**

```go
user := User{Name: "สมชาย", Email: "somchai@example.com", Age: 30}
result := db.Create(&user) // สร้าง record ใหม่
fmt.Println(user.ID)       // คืนค่า ID ที่ถูกสร้าง
fmt.Println(result.Error)  // error ถ้ามี
```

**Read**

```go
// ดึง record แรกที่ตรงเงื่อนไข
var user User
db.First(&user, 1)                 // by primary key
db.First(&user, "email = ?", "somchai@example.com")

// ดึงทั้งหมด
var users []User
db.Find(&users)

// พร้อมเงื่อนไข
db.Where("age > ?", 20).Find(&users)
db.Where(&User{Name: "สมชาย"}).Find(&users)
```

**Update**

```go
// อัปเดต single column
db.Model(&user).Update("Name", "สมชาย ใหม่")

// อัปเดตหลาย columns ด้วย struct (ไม่สนใจ zero values)
db.Model(&user).Updates(User{Name: "สมชาย ใหม่", Age: 31})

// อัปเดตหลาย columns ด้วย map
db.Model(&user).Updates(map[string]interface{}{"name": "สมชาย ใหม่", "age": 31})
```

**Delete**

```go
// soft delete (ถ้ามี gorm.DeletedAt)
db.Delete(&user, 1)

// hard delete
db.Unscoped().Delete(&user, 1)
```

#### 48.4.4 SessionFactory

Session factory เป็นรูปแบบการสร้าง `*gorm.DB` ที่มี configuration คงที่ (เช่น logging, connection pool) และสามารถสร้าง session ใหม่สำหรับแต่ละ request หรือ transaction

**ตัวอย่าง session factory**

```go
package db

import (
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    "log"
    "time"
)

type SessionFactory struct {
    db *gorm.DB
}

func NewSessionFactory(dsn string) (*SessionFactory, error) {
    // ตั้งค่า logger และ connection pool
    newLogger := logger.New(
        log.New(os.Stdout, "\r\n", log.LstdFlags),
        logger.Config{
            SlowThreshold: time.Second,
            LogLevel:      logger.Info,
            Colorful:      true,
        },
    )

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: newLogger,
        NowFunc: func() time.Time { return time.Now().UTC() },
    })
    if err != nil {
        return nil, err
    }

    // ตั้งค่า connection pool
    sqlDB, _ := db.DB()
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)

    return &SessionFactory{db: db}, nil
}

// NewSession คืนค่า session ใหม่สำหรับการทำ transaction หรือ query แบบแยก context
func (sf *SessionFactory) NewSession() *gorm.DB {
    return sf.db.Session(&gorm.Session{})
}
```

**การใช้งาน**

```go
factory, _ := db.NewSessionFactory(dsn)

// สร้าง session ใหม่สำหรับ request นี้
session := factory.NewSession()
var user User
session.First(&user, 1)

// เมื่อต้องการ transaction
session.Transaction(func(tx *gorm.DB) error {
    // ใช้ tx แทน session
    return nil
})
```

#### 48.4.5 Query Cache

GORM เองไม่มี query cache ในตัว แต่สามารถใช้ plugin หรือจัดการเองผ่าน Redis หรือ memory cache

**ตัวอย่างการใช้ Redis cache กับ GORM (แบบง่าย)**

```go
import (
    "context"
    "encoding/json"
    "github.com/go-redis/redis/v8"
    "gorm.io/gorm"
    "time"
)

type CachedDB struct {
    db    *gorm.DB
    cache *redis.Client
}

func (c *CachedDB) FirstWithCache(dest interface{}, conds ...interface{}) error {
    // สร้าง cache key จากเงื่อนไข
    key := fmt.Sprintf("query:%v", conds)

    // พยายามอ่านจาก cache
    val, err := c.cache.Get(context.Background(), key).Result()
    if err == nil {
        // พบใน cache
        return json.Unmarshal([]byte(val), dest)
    }

    // ไม่พบใน cache, query ฐานข้อมูล
    if err := c.db.First(dest, conds...).Error; err != nil {
        return err
    }

    // บันทึกผลลัพธ์ลง cache (serialize)
    data, _ := json.Marshal(dest)
    c.cache.Set(context.Background(), key, data, 5*time.Minute)
    return nil
}
```

#### 48.4.6 Cache Query Queue Processor

เมื่อมี query จำนวนมากที่ต้องการ cache แบบ asynchronous (เช่น การ pre-cache ข้อมูลที่คาดว่าจะถูกเรียกบ่อย) เราสามารถใช้ queue processor ในการรับ query, ประมวลผล, และเก็บผลลัพธ์ลง cache

**โครงสร้าง**

```
[Client Request] → [API] → (อ่าน cache) → ถ้ามี → ส่งกลับ
                     ↓ ถ้าไม่มี
              [Queue (Redis)] → [Worker] → Query DB → Store Cache → Notify
```

**ตัวอย่างการใช้ Redis เป็น queue และ worker**

```go
// โครงสร้างงาน
type CacheJob struct {
    QueryKey  string        `json:"query_key"`
    TableName string        `json:"table_name"`
    Conditions []interface{} `json:"conditions"`
    TTL       time.Duration `json:"ttl"`
}

// Producer: ส่งงานเข้าระบบ queue
func (c *CachedDB) QueueQuery(ctx context.Context, job CacheJob) error {
    data, _ := json.Marshal(job)
    return c.redis.LPush(ctx, "cache_queue", data).Err()
}

// Worker: รับงานจาก queue และประมวลผล
func (c *CachedDB) StartCacheWorker(ctx context.Context) {
    for {
        result, err := c.redis.BRPop(ctx, 0, "cache_queue").Result()
        if err != nil {
            continue
        }
        var job CacheJob
        json.Unmarshal([]byte(result[1]), &job)

        // Query ฐานข้อมูลตามเงื่อนไข
        var dest interface{}
        // สมมติว่าทราบ type จาก table name
        switch job.TableName {
        case "users":
            var users []User
            c.db.Where(job.Conditions...).Find(&users)
            dest = users
        // ...
        }

        // เก็บลง cache
        data, _ := json.Marshal(dest)
        c.redis.Set(ctx, job.QueryKey, data, job.TTL)
    }
}
```

#### 48.4.7 Queue Processor Function

ฟังก์ชันสำหรับประมวลผลคิวควรถูกออกแบบให้ทำงานแบบ concurrent, มีการ retry, และจัดการ error อย่างเหมาะสม

**เทมเพลต worker**

```go
type WorkerPool struct {
    workers int
    jobs    chan CacheJob
    wg      sync.WaitGroup
    db      *gorm.DB
    redis   *redis.Client
}

func NewWorkerPool(workers int, db *gorm.DB, redis *redis.Client) *WorkerPool {
    return &WorkerPool{
        workers: workers,
        jobs:    make(chan CacheJob, 100),
        db:      db,
        redis:   redis,
    }
}

func (wp *WorkerPool) Start(ctx context.Context) {
    for i := 0; i < wp.workers; i++ {
        wp.wg.Add(1)
        go wp.worker(ctx)
    }
}

func (wp *WorkerPool) worker(ctx context.Context) {
    defer wp.wg.Done()
    for {
        select {
        case job := <-wp.jobs:
            wp.processJob(job)
        case <-ctx.Done():
            return
        }
    }
}

func (wp *WorkerPool) processJob(job CacheJob) {
    // implement query และ cache logic
    // มี retry และ error logging
}
```

#### 48.4.8 SQL Queue Translate Rollback Commit

ในระบบที่มีการประมวลผลแบบคิว งานบางอย่างอาจเกี่ยวข้องกับการเปลี่ยนแปลงฐานข้อมูล (เช่น การอัปเดตสถานะ) ซึ่งต้องมี transaction เพื่อรักษาความถูกต้อง หลักการคือ:

- **Translate SQL to Queue**: แทนที่จะ execute SQL ทันที ให้สร้าง job ที่มีข้อมูลเพียงพอที่จะ execute SQL ในภายหลัง
- **Rollback**: ถ้า job ถูกประมวลผลไม่สำเร็จหลัง commit ไปแล้ว อาจต้องมีการชดเชย (compensating action) เนื่องจากฐานข้อมูลไม่สามารถ rollback ได้ง่ายเมื่อ transaction จบแล้ว
- **Commit**: เมื่อ job ทำงานสำเร็จ ให้ commit การเปลี่ยนแปลงในฐานข้อมูล

**ตัวอย่างการใช้ queue ร่วมกับ transaction (Outbox Pattern)**

```go
type OutboxMessage struct {
    ID          uint
    AggregateID string
    EventType   string
    Payload     string
    Status      string // pending, processed, failed
    CreatedAt   time.Time
}

// ใน transaction หลัก
func (s *Service) CreateOrder(ctx context.Context, order *Order) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // 1. บันทึก order
        if err := tx.Create(order).Error; err != nil {
            return err
        }

        // 2. ใส่ OutboxMessage
        msg := OutboxMessage{
            AggregateID: fmt.Sprintf("%d", order.ID),
            EventType:   "order.created",
            Payload:     `{"order_id":` + fmt.Sprintf("%d", order.ID) + `}`,
            Status:      "pending",
        }
        if err := tx.Create(&msg).Error; err != nil {
            return err
        }
        // transaction commit จะบันทึกทั้ง order และ outbox message
        return nil
    })
}

// Worker: อ่าน outbox messages และส่งไปยัง queue
func (s *Service) ProcessOutbox(ctx context.Context) {
    var messages []OutboxMessage
    s.db.Where("status = ?", "pending").Find(&messages)
    for _, msg := range messages {
        // ส่งไปยัง message broker
        if err := s.queue.Publish(msg.EventType, msg.Payload); err == nil {
            s.db.Model(&msg).Update("status", "processed")
        } else {
            // retry logic
        }
    }
}
```

ด้วยวิธีนี้ แม้ worker จะล้มเหลว เราสามารถ retry ได้โดยไม่สูญเสียข้อมูล

---

### 48.5 การออกแบบ Workflow

#### 48.5.1 Workflow CRUD พื้นฐาน

```
[Request] → [Controller] → [Service] → [Repository] → [GORM] → [Database]
                                                              ↓
[Response] ← [Controller] ← [Service] ← [Repository] ← [GORM] ← [Result]
```

#### 48.5.2 Workflow Query Cache

```
[Client] → [API] → ตรวจสอบ cache → มี → [Response] (cache hit)
                       ↓ ไม่มี
                 Query Database → เก็บ cache → [Response] (cache miss)
```

#### 48.5.3 Workflow Cache Queue Processor

```
[API] → (ต้องการ cache ข้อมูล) → สร้าง CacheJob → Push to Queue
[Worker] → Pop from Queue → Query DB → Store Cache (TTL)
[API] → (request ครั้งถัดไป) → Read Cache → Response
```

#### 48.5.4 Workflow Transaction + Queue (Outbox)

```
1. Start Transaction
   ├─ Update Business Data (Order)
   └─ Insert Outbox Message
2. Commit Transaction
3. Worker Polls Outbox
   ├─ Publish Message to Queue
   └─ Update Outbox Status to 'processed'
4. Consumer processes message
```

---

### 48.6 Code Template สำหรับนำไปใช้

#### 48.6.1 โครงสร้างโปรเจกต์

```
project/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── core/
│   │   └── user/
│   │       ├── entity.go
│   │       ├── repository.go
│   │       ├── service.go
│   │       └── handler.go
│   ├── platform/
│   │   ├── db/
│   │   │   └── gorm/
│   │   │       ├── session_factory.go
│   │   │       └── user_repo.go
│   │   └── cache/
│   │       ├── redis.go
│   │       └── cache_worker.go
│   └── transport/
│       └── http/
│           └── handler.go
├── pkg/
│   └── queue/
│       └── redis_queue.go
├── go.mod
└── config.yaml
```

#### 48.6.2 Session Factory Template

```go
// internal/platform/db/gorm/session_factory.go
package gorm

import (
    "gorm.io/gorm"
    "gorm.io/driver/postgres"
    "time"
)

type SessionFactory struct {
    db *gorm.DB
}

func NewSessionFactory(dsn string) (*SessionFactory, error) {
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        // กำหนด logger, nowFunc ตามต้องการ
    })
    if err != nil {
        return nil, err
    }

    sqlDB, _ := db.DB()
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)

    return &SessionFactory{db: db}, nil
}

func (sf *SessionFactory) GetDB() *gorm.DB {
    return sf.db
}

func (sf *SessionFactory) NewSession() *gorm.DB {
    return sf.db.Session(&gorm.Session{})
}
```

#### 48.6.3 Repository Implementation with Cache and Queue

```go
// internal/platform/db/gorm/user_repo.go
package gorm

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
    "github.com/go-redis/redis/v8"
    "gorm.io/gorm"
    "your-project/internal/core/user"
)

type userRepository struct {
    db    *gorm.DB
    cache *redis.Client
    queue *redis.Client // same as cache for simplicity
}

func NewUserRepository(db *gorm.DB, redisClient *redis.Client) user.Repository {
    return &userRepository{
        db:    db,
        cache: redisClient,
        queue: redisClient,
    }
}

func (r *userRepository) FindByID(ctx context.Context, id uint) (*user.User, error) {
    // Try cache first
    key := fmt.Sprintf("user:%d", id)
    val, err := r.cache.Get(ctx, key).Result()
    if err == nil {
        var u user.User
        if err := json.Unmarshal([]byte(val), &u); err == nil {
            return &u, nil
        }
    }

    // Cache miss, query DB
    var u user.User
    if err := r.db.WithContext(ctx).First(&u, id).Error; err != nil {
        return nil, err
    }

    // Store in cache asynchronously via queue
    job := CacheJob{
        QueryKey: key,
        Table:    "users",
        ID:       id,
        TTL:      10 * time.Minute,
    }
    data, _ := json.Marshal(job)
    r.queue.LPush(ctx, "cache_queue", data)

    return &u, nil
}

// ... other methods
```

#### 48.6.4 Queue Worker Template

```go
// internal/platform/cache/cache_worker.go
package cache

import (
    "context"
    "encoding/json"
    "time"
    "github.com/go-redis/redis/v8"
    "gorm.io/gorm"
)

type CacheWorker struct {
    redis *redis.Client
    db    *gorm.DB
    stop  chan struct{}
}

type CacheJob struct {
    QueryKey string        `json:"query_key"`
    Table    string        `json:"table"`
    ID       uint          `json:"id"`
    TTL      time.Duration `json:"ttl"`
}

func NewCacheWorker(redis *redis.Client, db *gorm.DB) *CacheWorker {
    return &CacheWorker{
        redis: redis,
        db:    db,
        stop:  make(chan struct{}),
    }
}

func (w *CacheWorker) Start(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return
        case <-w.stop:
            return
        default:
            // Pop job from queue
            result, err := w.redis.BRPop(ctx, 0, "cache_queue").Result()
            if err != nil {
                continue
            }

            var job CacheJob
            if err := json.Unmarshal([]byte(result[1]), &job); err != nil {
                continue
            }

            // Execute query based on table
            var data interface{}
            switch job.Table {
            case "users":
                var user User
                w.db.First(&user, job.ID)
                data = user
            // add more cases
            }

            // Store in cache
            bytes, _ := json.Marshal(data)
            w.redis.Set(ctx, job.QueryKey, bytes, job.TTL)
        }
    }
}

func (w *CacheWorker) Stop() {
    close(w.stop)
}
```

#### 48.6.5 การใช้ใน main.go

```go
// cmd/api/main.go
package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "your-project/internal/platform/db/gorm"
    "your-project/internal/platform/cache"
    "your-project/internal/core/user"
    userHttp "your-project/internal/transport/http"
)

func main() {
    // Load config
    cfg := loadConfig()

    // Setup session factory
    sf, err := gorm.NewSessionFactory(cfg.Database.DSN)
    if err != nil {
        log.Fatal(err)
    }

    // Setup Redis client
    redisClient := redis.NewClient(&redis.Options{
        Addr: cfg.Redis.Addr,
    })

    // Start cache worker
    worker := cache.NewCacheWorker(redisClient, sf.GetDB())
    ctx, cancel := context.WithCancel(context.Background())
    go worker.Start(ctx)

    // Setup repositories, services, handlers
    userRepo := gorm.NewUserRepository(sf.GetDB(), redisClient)
    userService := user.NewService(userRepo)
    userHandler := userHttp.NewHandler(userService)

    // Setup HTTP server
    router := setupRouter(userHandler)

    srv := &http.Server{
        Addr:    ":8080",
        Handler: router,
    }

    // Graceful shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    cancel() // stop worker
    worker.Stop()

    ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancelShutdown()
    if err := srv.Shutdown(ctxShutdown); err != nil {
        log.Fatal("Server shutdown error:", err)
    }
}
```

---

### 48.7 สรุป

GORM เป็น ORM ที่ทรงพลังและใช้งานง่าย ช่วยให้นักพัฒนา Go จัดการฐานข้อมูลได้อย่างมีประสิทธิภาพ ด้วยฟีเจอร์ครบครันตั้งแต่ CRUD พื้นฐานไปจนถึง advanced patterns เช่น session factory, query cache, และการผสานกับระบบคิว เพื่อเพิ่มความเร็วและความน่าเชื่อถือของแอปพลิเคชัน

- **Session Factory** ช่วยแยก session และจัดการ connection pool
- **Query Cache** ลดภาระฐานข้อมูลด้วยการเก็บผลลัพธ์ซ้ำ
- **Cache Queue Processor** ทำให้การ pre-cache เป็น asynchronous ไม่กระทบ request ปัจจุบัน
- **SQL Queue Translate Rollback Commit** ผสาน transaction กับระบบคิวเพื่อความถูกต้องของข้อมูล

การใช้ GORM ร่วมกับหลักการออกแบบที่ดีจะช่วยให้คุณสร้างระบบที่มีประสิทธิภาพสูง บำรุงรักษาง่าย และพร้อมขยายตามความต้องการของธุรกิจ