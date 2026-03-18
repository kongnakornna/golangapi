<img src="https://r2cdn.perplexity.ai/pplx-full-logo-primary-dark%402x.png" style="height:64px;margin-right:32px"/>

## การออกแบบระบบ Flow สำหรับ Spring Boot + Kafka (Producer \& Consumer)

ระบบศูนย์บริการรถยนต์ออนไลน์ที่ใช้ Spring Boot กับ Apache Kafka จะออกแบบเป็นสถาปัตยกรรมแบบ event-driven ที่ microservices สื่อสารกันผ่าน events  ระบบนี้แบ่งเป็น 4 กระบวนการหลักที่ทำงานร่วมกันอย่างมีประสิทธิภาพ[^1][^2][^3][^4]

## สถาปัตยกรรมระดับสูง

### Core Components

**Producers (ส่งข้อความ)**

- API Gateway Service - รับ HTTP requests จาก frontend
- Booking Service - จัดการการจองนัดหมาย
- Repair Service - จัดการงานซ่อม
- Payment Service - ประมวลผลการชำระเงิน

**Kafka Cluster**

- Topics: `booking-events`, `repair-events`, `payment-events`, `notification-events`
- Partitions: แต่ละ topic มี 10 partitions เพื่อรองรับการประมวลผลแบบ parallel
- Replication Factor: 3 เพื่อความปลอดภัยของข้อมูล

**Consumers (รับข้อความ)**

- Notification Service - ส่งการแจ้งเตือนให้ลูกค้าและพนักงาน
- Analytics Service - วิเคราะห์ข้อมูลและสร้างรายงาน
- Cache Sync Service - อัพเดท Redis cache
- Audit Service - บันทึก audit logs

[^2][^3][^4]

## Flow 1: การจองนัดหมายบริการ (Booking Flow)

### ขั้นตอนการทำงาน

**1. ลูกค้าสร้างการจอง**

```
ReactJS Frontend → POST /api/v1/bookings
                   (customerId, vehicleId, serviceType, bookingDate)
```

**2. Booking Service (Producer)**

```java
// ตรวจสอบความพร้อม
if (isSlotAvailable(bookingDate)) {
    // บันทึกลง PostgreSQL
    Booking booking = bookingRepository.save(newBooking);
    
    // อัพเดท Redis cache
    redisTemplate.opsForValue().set("booking:" + booking.getId(), booking);
    
    // สร้าง Event
    BookingCreatedEvent event = BookingCreatedEvent.builder()
        .eventId(UUID.randomUUID().toString())
        .eventType("BOOKING_CREATED")
        .bookingId(booking.getId())
        .customerId(booking.getCustomerId())
        .bookingDate(booking.getBookingDate())
        .timestamp(Instant.now())
        .build();
    
    // Publish ไปยัง Kafka
    kafkaTemplate.send("booking-events", booking.getId(), event);
}
```

**3. Consumers รับและประมวลผล**

```java
// Notification Service Consumer
@KafkaListener(topics = "booking-events", groupId = "notification-group")
public void handleBookingEvent(BookingCreatedEvent event) {
    // ส่งอีเมล/SMS ยืนยันการจอง
    notificationService.sendBookingConfirmation(event);
    
    // ส่ง push notification
    fcmService.sendNotification(event.getCustomerId(), 
        "การจองของคุณได้รับการยืนยันแล้ว");
}

// Analytics Service Consumer
@KafkaListener(topics = "booking-events", groupId = "analytics-group")
public void analyzeBooking(BookingCreatedEvent event) {
    // บันทึกข้อมูลสำหรับ dashboard
    analyticsRepository.recordBooking(event);
    
    // อัพเดทสถิติ real-time
    redisTemplate.opsForValue().increment("daily:bookings:" + LocalDate.now());
}

// Audit Service Consumer
@KafkaListener(topics = "booking-events", groupId = "audit-group")
public void auditBooking(BookingCreatedEvent event) {
    // บันทึก audit log
    auditLogRepository.save(AuditLog.from(event));
}
```


## Flow 2: กระบวนการซ่อม (Repair Flow)

### Event-Driven Repair Workflow

**1. เริ่มงานซ่อม**

```
Technician App → POST /api/v1/repairs/start
                 (bookingId, technicianId, estimatedTime)

Repair Service (Producer) →
    ├─ Save to PostgreSQL
    ├─ Update Redis cache
    └─ Publish "REPAIR_STARTED" event to Kafka
```

**2. อัพเดทสถานะระหว่างซ่อม**

```java
// สร้าง multiple events ตามสถานะ
RepairStatusEvent event = RepairStatusEvent.builder()
    .eventType("REPAIR_IN_PROGRESS")
    .status("DIAGNOSING") // → REPAIRING → TESTING
    .progressPercentage(25)
    .build();

kafkaTemplate.send("repair-events", repairId, event);
```

**3. Real-time Status Updates**

```java
// Consumer อัพเดทข้อมูลให้ลูกค้าเห็น real-time
@KafkaListener(topics = "repair-events", groupId = "status-update-group")
public void updateRepairStatus(RepairStatusEvent event) {
    // อัพเดท Redis สำหรับ real-time query
    redisTemplate.opsForHash().put(
        "repair:status:" + event.getRepairId(),
        "progress", event.getProgressPercentage()
    );
    
    // ส่ง notification ให้ลูกค้า
    webSocketService.sendToUser(event.getCustomerId(), event);
}
```

**4. ซ่อมเสร็จสิ้น**

```
Repair Service → Publish "REPAIR_COMPLETED" event
                 ↓
Payment Service Consumer → สร้างใบแจ้งหนี้
Notification Service → แจ้งเตือนลูกค้ามารับรถ
Inventory Service → อัพเดท stock อะไหล่
```


## Flow 3: การชำระเงิน (Payment Flow)

### Event Chain Pattern

**1. สร้าง Invoice**

```java
// Payment Service รอ REPAIR_COMPLETED event
@KafkaListener(topics = "repair-events")
public void onRepairCompleted(RepairCompletedEvent event) {
    // สร้าง invoice
    Invoice invoice = invoiceService.createInvoice(event.getRepairId());
    
    // Publish INVOICE_CREATED event
    PaymentEvent paymentEvent = new PaymentEvent(
        "INVOICE_CREATED", invoice.getId(), invoice.getAmount()
    );
    kafkaTemplate.send("payment-events", invoice.getId(), paymentEvent);
}
```

**2. ประมวลผลการชำระเงิน**

```
Customer → POST /api/v1/payments
           ↓
Payment Service (Producer) →
    ├─ เรียก Payment Gateway API
    ├─ บันทึกผล transaction ใน PostgreSQL
    └─ Publish "PAYMENT_COMPLETED" event
           ↓
Consumers:
    ├─ Notification Service → ส่งใบเสร็จทาง email
    ├─ Booking Service → อัพเดทสถานะเป็น "PAID"
    └─ Analytics Service → บันทึกรายได้
```


## Flow 4: การจัดการข้อผิดพลาด (Error Handling Flow)

### Saga Pattern with Compensation

**Distributed Transaction Management**

```java
// หากการชำระเงินล้มเหลว
@KafkaListener(topics = "payment-events")
public void handlePaymentFailed(PaymentFailedEvent event) {
    // Publish compensation events
    BookingCompensationEvent compensation = new BookingCompensationEvent(
        "ROLLBACK_BOOKING", event.getBookingId()
    );
    kafkaTemplate.send("compensation-events", compensation);
}

// Booking Service รับ compensation event
@KafkaListener(topics = "compensation-events")
public void compensateBooking(BookingCompensationEvent event) {
    // Rollback การจอง
    bookingService.cancelBooking(event.getBookingId());
    
    // คืนเงินหากมีการจ่ายล่วงหน้า
    // แจ้งเตือนลูกค้า
}
```


### Retry และ Dead Letter Queue

```java
// Configuration สำหรับ retry mechanism
@Bean
public ConcurrentKafkaListenerContainerFactory<String, Object> kafkaListenerContainerFactory() {
    factory.setCommonErrorHandler(new DefaultErrorHandler(
        new DeadLetterPublishingRecoverer(kafkaTemplate),
        new FixedBackOff(1000L, 3L) // Retry 3 ครั้ง ห่างกัน 1 วินาที
    ));
    return factory;
}

// Consumer with error handling
@KafkaListener(topics = "booking-events")
public void processBooking(BookingEvent event) {
    try {
        bookingService.process(event);
    } catch (RetryableException e) {
        // Kafka จะ retry อัตโนมัติ
        throw e;
    } catch (NonRetryableException e) {
        // ส่งไป DLQ ทันที
        log.error("Non-retryable error", e);
        throw e;
    }
}
```


## สรุปโครงสร้าง Event Flow

| Event Type | Producer | Topic | Consumers | Action |
| :-- | :-- | :-- | :-- | :-- |
| BOOKING_CREATED | Booking Service | booking-events | Notification, Analytics, Audit | ส่งการแจ้งเตือน, บันทึกสถิติ |
| REPAIR_STARTED | Repair Service | repair-events | Status Update, Notification | อัพเดทสถานะ real-time |
| REPAIR_COMPLETED | Repair Service | repair-events | Payment, Notification, Inventory | สร้าง invoice, อัพเดท stock |
| PAYMENT_COMPLETED | Payment Service | payment-events | Notification, Booking, Analytics | ส่งใบเสร็จ, ปิดงาน |
| NOTIFICATION_SENT | Notification Service | notification-events | Analytics, Audit | ติดตามประสิทธิภาพการแจ้งเตือน |

ระบบนี้รองรับ asynchronous communication, fault tolerance และ scalability สูง โดยแต่ละ microservice สามารถ scale แยกอิสระตามปริมาณ events ที่ต้องประมวลผล[^3][^1][^2]
<span style="display:none">[^10][^11][^5][^6][^7][^8][^9]</span>

<div align="center">⁂</div>

[^1]: https://www.redpanda.com/blog/build-event-driven-microservices-spring-boot

[^2]: https://www.linkedin.com/pulse/spring-boot-kafka-how-design-resilient-event-flow-bruno-vieira-q0vgf

[^3]: https://www.geeksforgeeks.org/apache-kafka/kafka-architecture/

[^4]: https://www.instaclustr.com/education/apache-kafka/apache-kafka-architecture-a-complete-guide-2025/

[^5]: https://developers.ascendcorp.com/สร้าง-event-driven-systems-ด้วย-spring-cloud-stream-และ-apache-kafka-ee9fa845a8ed

[^6]: https://blog.devops.dev/building-a-reactive-banking-microservices-architecture-with-spring-boot-kafka-redis-1bd4e25c5fd9

[^7]: https://stackoverflow.com/questions/42140285/how-to-implement-a-microservice-event-driven-architecture-with-spring-cloud-stre

[^8]: https://ably.com/blog/realtime-ticket-booking-solution-kafka-fastapi-ably

[^9]: https://www.tinybird.co/blog/event-sourcing-with-kafka

[^10]: https://www.linkedin.com/posts/rishi-jha-168719172_github-rishijha02kafka-taxiservice-activity-7384290161993592832-N3bi

[^11]: https://www.instaclustr.com/education/apache-kafka/spring-boot-with-apache-kafka-tutorial-and-best-practices/

