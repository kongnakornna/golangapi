<img src="https://r2cdn.perplexity.ai/pplx-full-logo-primary-dark%402x.png" style="height:64px;margin-right:32px"/>

## การจัดการ Transaction ระหว่าง Database กับ Kafka ใน Spring Boot

การจัดการ transaction ระหว่าง Database และ Kafka ใน Spring Boot ต้องเข้าใจว่า Kafka ไม่รองรับ XA transactions แบบ distributed ดังนั้น Spring จะจัดการ transactions แยกกันแต่ synchronize ให้ทำงานร่วมกัน  มี 3 วิธีหลักในการจัดการ transactions นี้[^1][^2][^3][^4]

## วิธีที่ 1: Database Commit First (Default)

### Configuration

```properties
# application.properties
spring.kafka.bootstrap-servers=localhost:9092
spring.kafka.producer.transaction-id-prefix=tx-
spring.kafka.consumer.enable-auto-commit=false
spring.kafka.consumer.properties.isolation.level=read_committed

# Database
spring.datasource.url=jdbc:postgresql://localhost:5432/carservice
spring.datasource.username=postgres
spring.datasource.password=password
```


### Transaction Manager Setup

```java
@Configuration
public class TransactionConfig {
    
    // Spring Boot สร้าง transaction managers อัตโนมัติ:
    // 1. transactionManager (JPA/DataSource)
    // 2. kafkaTransactionManager (Kafka)
    
    @Bean
    public DataSourceTransactionManager dstm(DataSource dataSource) {
        return new DataSourceTransactionManager(dataSource);
    }
}
```


### Implementation: Consumer Listener

```java
@Service
@RequiredArgsConstructor
@Slf4j
public class BookingEventListener {
    
    private final BookingRepository bookingRepository;
    private final KafkaTemplate<String, Object> kafkaTemplate;
    private final JdbcTemplate jdbcTemplate;
    
    /**
     * Kafka Listener container เริ่ม Kafka transaction
     * @Transactional เริ่ม DB transaction
     * 
     * ลำดับการ commit:
     * 1. DB commits first
     * 2. Kafka commits second
     * 
     * หาก Kafka commit ล้มเหลว: message จะถูก redelivered
     * ดังนั้น DB operation ต้องเป็น idempotent
     */
    @KafkaListener(
        id = "booking-listener",
        topics = "booking-events",
        groupId = "booking-consumer-group"
    )
    @Transactional("dstm")  // ใช้ DataSource Transaction Manager
    public void handleBookingEvent(BookingEvent event) {
        log.info("Received booking event: {}", event.getEventId());
        
        // 1. บันทึกข้อมูลลง PostgreSQL
        Booking booking = Booking.builder()
            .id(event.getBookingId())
            .customerId(event.getCustomerId())
            .vehicleId(event.getVehicleId())
            .bookingDate(event.getBookingDate())
            .status("CONFIRMED")
            .build();
        
        bookingRepository.save(booking);
        
        // 2. ส่ง notification event ไปยัง Kafka
        // KafkaTemplate จะ synchronize กับ DB transaction
        NotificationEvent notification = NotificationEvent.builder()
            .eventId(UUID.randomUUID().toString())
            .customerId(event.getCustomerId())
            .message("การจองของคุณได้รับการยืนยัน")
            .timestamp(Instant.now())
            .build();
        
        kafkaTemplate.send("notification-events", notification);
        
        // Method จบ: 
        // 1. DB transaction commits
        // 2. Kafka transaction commits
    }
}
```


### Implementation: Producer Service

```java
@Service
@RequiredArgsConstructor
@Slf4j
public class RepairService {
    
    private final RepairRepository repairRepository;
    private final KafkaTemplate<String, Object> kafkaTemplate;
    
    /**
     * Producer-only transaction
     * KafkaTemplate synchronizes กับ DB transaction
     */
    @Transactional("dstm")
    public RepairResponse startRepair(RepairRequest request) {
        // 1. บันทึกข้อมูลการซ่อมใน database
        Repair repair = Repair.builder()
            .id(UUID.randomUUID().toString())
            .bookingId(request.getBookingId())
            .technicianId(request.getTechnicianId())
            .status("IN_PROGRESS")
            .estimatedCost(request.getEstimatedCost())
            .createdAt(LocalDateTime.now())
            .build();
        
        Repair savedRepair = repairRepository.save(repair);
        
        // 2. Publish event ไปยัง Kafka
        // KafkaTemplate จะรอจน DB commit ก่อน
        RepairEvent event = RepairEvent.builder()
            .eventId(UUID.randomUUID().toString())
            .eventType("REPAIR_STARTED")
            .repairId(savedRepair.getId())
            .timestamp(Instant.now())
            .build();
        
        kafkaTemplate.send("repair-events", savedRepair.getId(), event);
        
        // ลำดับการ commit:
        // 1. PostgreSQL transaction commits
        // 2. Kafka transaction commits
        
        return RepairResponse.from(savedRepair);
    }
}
```


## วิธีที่ 2: Kafka Commit First (Nested Transactions)

### Implementation

```java
@Service
@RequiredArgsConstructor
public class PaymentService {
    
    private final PaymentRepository paymentRepository;
    private final JdbcTemplate jdbcTemplate;
    
    /**
     * Outer method ใช้ DataSource Transaction Manager
     * Inner method ใช้ Kafka Transaction Manager
     * 
     * ลำดับการ commit:
     * 1. Kafka commits first
     * 2. DB commits only if Kafka succeeds
     */
    @Transactional("dstm")
    public void processPayment(PaymentRequest request) {
        // 1. บันทึกข้อมูลใน database
        jdbcTemplate.execute(
            "INSERT INTO payments (id, booking_id, amount, status) " +
            "VALUES ('" + request.getId() + "', '" + request.getBookingId() + 
            "', " + request.getAmount() + ", 'PENDING')"
        );
        
        // 2. เรียก inner method เพื่อส่ง Kafka message
        sendPaymentEvent(request);
        
        // Kafka commits ก่อน (จาก inner method)
        // จากนั้น DB commits
    }
    
    @Transactional("kafkaTransactionManager")
    public void sendPaymentEvent(PaymentRequest request) {
        PaymentEvent event = PaymentEvent.builder()
            .eventId(UUID.randomUUID().toString())
            .eventType("PAYMENT_PROCESSED")
            .paymentId(request.getId())
            .amount(request.getAmount())
            .build();
        
        kafkaTemplate.send("payment-events", event);
        
        // Kafka transaction commits เมื่อออกจาก method นี้
    }
}
```


## วิธีที่ 3: Transactional Outbox Pattern

### วิธีการทำงาน

Outbox Pattern แก้ปัญหา dual-write โดยเขียนข้อมูลทั้ง business data และ outbox message ลงใน database ในหนึ่ง transaction จากนั้นใช้กระบวนการแยกต่างหากในการอ่าน outbox และส่งไปยัง Kafka[^3][^5]

### Database Schema

```sql
-- Business table
CREATE TABLE bookings (
    id VARCHAR(36) PRIMARY KEY,
    customer_id VARCHAR(36) NOT NULL,
    vehicle_id VARCHAR(36) NOT NULL,
    booking_date TIMESTAMP NOT NULL,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Outbox table
CREATE TABLE outbox_events (
    id VARCHAR(36) PRIMARY KEY,
    aggregate_type VARCHAR(50) NOT NULL,
    aggregate_id VARCHAR(36) NOT NULL,
    event_type VARCHAR(50) NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    processed BOOLEAN DEFAULT FALSE,
    processed_at TIMESTAMP
);

CREATE INDEX idx_outbox_processed ON outbox_events(processed, created_at);
```


### Implementation: Write to Outbox

```java
@Service
@RequiredArgsConstructor
public class BookingService {
    
    private final BookingRepository bookingRepository;
    private final OutboxEventRepository outboxEventRepository;
    private final ObjectMapper objectMapper;
    
    /**
     * เขียนทั้ง booking และ outbox event ใน transaction เดียว
     * รับประกัน atomicity
     */
    @Transactional
    public BookingResponse createBooking(BookingRequest request) {
        // 1. สร้าง booking
        Booking booking = Booking.builder()
            .id(UUID.randomUUID().toString())
            .customerId(request.getCustomerId())
            .vehicleId(request.getVehicleId())
            .bookingDate(request.getBookingDate())
            .status("PENDING")
            .build();
        
        Booking savedBooking = bookingRepository.save(booking);
        
        // 2. สร้าง outbox event
        BookingEvent event = BookingEvent.builder()
            .eventId(UUID.randomUUID().toString())
            .eventType("BOOKING_CREATED")
            .bookingId(savedBooking.getId())
            .customerId(savedBooking.getCustomerId())
            .timestamp(Instant.now())
            .build();
        
        OutboxEvent outboxEvent = OutboxEvent.builder()
            .id(event.getEventId())
            .aggregateType("Booking")
            .aggregateId(savedBooking.getId())
            .eventType(event.getEventType())
            .payload(objectMapper.writeValueAsString(event))
            .processed(false)
            .build();
        
        outboxEventRepository.save(outboxEvent);
        
        // Transaction commit: ทั้ง booking และ outbox ถูกบันทึกพร้อมกัน
        
        return BookingResponse.from(savedBooking);
    }
}
```


### Implementation: Outbox Publisher

```java
@Component
@RequiredArgsConstructor
@Slf4j
public class OutboxEventPublisher {
    
    private final OutboxEventRepository outboxEventRepository;
    private final KafkaTemplate<String, String> kafkaTemplate;
    private final ObjectMapper objectMapper;
    
    /**
     * Scheduled task ดึง unprocessed events และส่งไปยัง Kafka
     */
    @Scheduled(fixedDelay = 5000) // ทุก 5 วินาที
    @Transactional
    public void publishOutboxEvents() {
        List<OutboxEvent> events = outboxEventRepository
            .findTop100ByProcessedFalseOrderByCreatedAtAsc();
        
        for (OutboxEvent event : events) {
            try {
                // ส่งไปยัง Kafka
                String topic = determineTopicByEventType(event.getEventType());
                
                kafkaTemplate.send(topic, event.getAggregateId(), event.getPayload())
                    .addCallback(
                        result -> {
                            // สำเร็จ: mark as processed
                            event.setProcessed(true);
                            event.setProcessedAt(Instant.now());
                            outboxEventRepository.save(event);
                            log.info("Published outbox event: {}", event.getId());
                        },
                        ex -> {
                            log.error("Failed to publish outbox event: {}", 
                                event.getId(), ex);
                            // Retry ในรอบถัดไป
                        }
                    );
                
            } catch (Exception e) {
                log.error("Error processing outbox event: {}", event.getId(), e);
            }
        }
    }
    
    private String determineTopicByEventType(String eventType) {
        if (eventType.startsWith("BOOKING")) return "booking-events";
        if (eventType.startsWith("REPAIR")) return "repair-events";
        if (eventType.startsWith("PAYMENT")) return "payment-events";
        return "default-events";
    }
}
```


## การจัดการ Error และ Rollback

### Retry Mechanism

```java
@Configuration
public class KafkaErrorHandlingConfig {
    
    @Bean
    public ConcurrentKafkaListenerContainerFactory<String, Object> 
        kafkaListenerContainerFactory(
            ConsumerFactory<String, Object> consumerFactory,
            KafkaTemplate<String, Object> kafkaTemplate) {
        
        ConcurrentKafkaListenerContainerFactory<String, Object> factory =
            new ConcurrentKafkaListenerContainerFactory<>();
        factory.setConsumerFactory(consumerFactory);
        
        // Error handler with retry and DLQ
        DefaultErrorHandler errorHandler = new DefaultErrorHandler(
            new DeadLetterPublishingRecoverer(kafkaTemplate),
            new FixedBackOff(2000L, 3L) // Retry 3 ครั้ง ห่างกัน 2 วินาที
        );
        
        // กำหนด exceptions ที่ควร retry
        errorHandler.addRetryableExceptions(
            TransientDataAccessException.class,
            TimeoutException.class
        );
        
        // กำหนด exceptions ที่ไม่ควร retry (ส่งไป DLQ ทันที)
        errorHandler.addNotRetryableExceptions(
            IllegalArgumentException.class,
            JsonProcessingException.class
        );
        
        factory.setCommonErrorHandler(errorHandler);
        return factory;
    }
}
```


### Idempotent Consumer

```java
@Service
@RequiredArgsConstructor
public class IdempotentBookingConsumer {
    
    private final BookingRepository bookingRepository;
    private final ProcessedMessageRepository processedMessageRepository;
    
    @KafkaListener(topics = "booking-events")
    @Transactional
    public void handleBooking(BookingEvent event) {
        // ตรวจสอบว่าเคย process message นี้แล้วหรือไม่
        if (processedMessageRepository.existsById(event.getEventId())) {
            log.warn("Duplicate message detected: {}", event.getEventId());
            return; // Skip duplicate
        }
        
        // Process booking
        Booking booking = bookingRepository.findById(event.getBookingId())
            .orElse(new Booking());
        booking.setStatus("CONFIRMED");
        bookingRepository.save(booking);
        
        // บันทึกว่าได้ process message นี้แล้ว
        ProcessedMessage processed = new ProcessedMessage(
            event.getEventId(), 
            Instant.now()
        );
        processedMessageRepository.save(processed);
    }
}
```


## สรุปการเปรียบเทียบ

| แนวทาง | ข้อดี | ข้อเสีย | Use Case |
| :-- | :-- | :-- | :-- |
| DB Commit First | ง่าย, built-in support | Message อาจ redelivered, ต้อง idempotent | Consumer ที่ไม่ซับซ้อน |
| Kafka Commit First | Kafka guarantees ก่อน | DB อาจล้มเหลวหลัง Kafka commit | Critical event publishing |
| Outbox Pattern | Guaranteed delivery, true atomicity | ซับซ้อนกว่า, latency สูงขึ้น | Mission-critical systems |

สำหรับระบบศูนย์บริการรถยนต์ ควรใช้ Outbox Pattern สำหรับ critical operations เช่น payment และใช้ DB Commit First สำหรับ operations ทั่วไปเพื่อความเรียบง่าย[^3][^5]
<span style="display:none">[^10][^11][^6][^7][^8][^9]</span>

<div align="center">⁂</div>

[^1]: https://www.instaclustr.com/education/apache-kafka/spring-boot-with-apache-kafka-tutorial-and-best-practices/

[^2]: https://piotrminkowski.com/2022/10/29/kafka-transactions-with-spring-boot/

[^3]: https://spring.io/blog/2023/10/24/a-use-case-for-transactions-adapting-to-transactional-outbox-pattern

[^4]: https://docs.spring.io/spring-kafka/reference/tips.html

[^5]: https://axual.com/blog/transactional-outbox-pattern-kafka

[^6]: https://docs.spring.io/spring-kafka/reference/kafka/transactions.html

[^7]: https://stackoverflow.com/questions/52207598/synchronising-transactions-between-database-and-kafka-producer

[^8]: https://www.youtube.com/watch?v=cORDvFYE1og

[^9]: https://www.confluent.io/blog/transactions-apache-kafka/

[^10]: https://spring.io/blog/2023/09/27/introduction-to-transactions-in-spring-cloud-stream-kafka-applications

[^11]: https://www.linkedin.com/pulse/spring-boot-kafka-how-design-resilient-event-flow-bruno-vieira-q0vgf

