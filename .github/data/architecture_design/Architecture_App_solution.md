# Architectural Design IoT Solution
## เอกสารเชิงลึกตามโครงสร้าง 10 บท

---

# 1. บทนำ (Introduction)

## 1.1 ความเป็นมา

ในยุคที่อุปกรณ์ IoT (Internet of Things) ทวีจำนวนขึ้นอย่างรวดเร็ว องค์กรธุรกิจต้องเผชิญกับความท้าทายในการออกแบบระบบที่รองรับอุปกรณ์นับล้าน ส่งข้อมูลแบบ real-time ประมวลผลด้วย AI และรักษาความปลอดภัยระดับองค์กร การออกแบบ Software Architecture ที่ดีจึงไม่ใช่แค่เรื่องเทคนิค แต่เป็นเรื่อง стратегิc ที่กำหนดความอยู่รอดของธุรกิจ

เอกสารนี้จัดทำขึ้นเพื่อเป็น **คู่มือเชิงลึก** สำหรับสถาปนิกซอฟต์แวร์ วิศวกร และผู้จัดการเทคนิค ที่ต้องออกแบบระบบ IoT Solution ตั้งแต่ระดับ Device จนถึง Cloud ครอบคลุมตั้งแต่พื้นฐาน Computer Science ไปจนถึง Advanced Algorithm, Protocol, Database, Monitoring และ Business Model

## 1.2 วัตถุประสงค์

1. สร้างความเข้าใจร่วมกันเกี่ยวกับ **Software Architecture** สำหรับ IoT
2. ให้แนวทางเลือกใช้ **Protocol** ที่เหมาะสมกับ use case
3. ออกแบบ **Database Schema** ครบทุกตารางสำหรับ Logistics + IoT
4. อธิบาย **Advanced Algorithms** พร้อมโค้ดใช้งานจริง
5. วางแผน **Monitoring & Observability** ระดับ production
6. เสนอ **Business Canvas Model** สำหรับ Parking, Robotic Parking, Machine Vision
7. รวบรวม **ปัญหาและแนวทางแก้ไข** จากประสบการณ์จริง

## 1.3 ขอบเขต

- **Software Architecture**: Data Structures, Algorithms, Big O, FP, OOP, System Design, DDD, Clean Architecture, EDA
- **IoT Stack**: MQTT, REST, gRPC, WebSocket, SSE, GraphQL, Webhook, Message Queue
- **Database**: PostgreSQL, InfluxDB, MongoDB, Redis, Vector DB + S3
- **Monitoring**: Prometheus, ELK, OpenTelemetry
- **AI/ML**: Ollama, Embeddings, Vector Search
- **Business**: Logistics SaaS, ERP, CRM, Parking, Robotic Parking, Vision Camera
- **Algorithms**: Suffix Automaton, Palindromic Automaton, Mo's with Update, Centroid + Fenwick, HLD + Lazy SegTree, Fenwick 2D, Wavelet Tree 2D

## 1.4 กลุ่มเป้าหมาย

- Software Architect / Solution Architect
- Backend Engineer (Go, Java, Python)
- DevOps / SRE Engineer
- Data Engineer
- Technical Product Manager
- CTO / Engineering Manager

## 1.5 วิธีการอ่านเอกสาร

เอกสารแบ่งเป็น 3 ระดับ:
- **ระดับพื้นฐาน** (บทที่ 1-3): ทำความเข้าใจ concept
- **ระดับกลาง** (บทที่ 4-7): ออกแบบและประยุกต์
- **ระดับสูง** (บทที่ 8-10): Case study และประสบการณ์จริง

---

# 2. บทนิยาม (Definitions)

## 2.1 คำศัพท์ด้าน Software Architecture

| คำศัพท์ | ความหมาย |
|---|---|
| **Software Architecture** | โครงสร้างระดับสูงของระบบ ประกอบด้วยองค์ประกอบ ความสัมพันธ์ และหลักการกำกับ |
| **Clean Architecture** | แนวทางออกแบบที่แยก business logic ออกจาก framework ด้วย Dependency Rule |
| **Domain-Driven Design (DDD)** | แนวทางที่เน้น domain ธุรกิจเป็นศูนย์กลาง ผ่าน Ubiquitous Language และ Bounded Context |
| **Event-Driven Architecture (EDA)** | สถาปัตยกรรมที่ component สื่อสารผ่าน event แบบ asynchronous |
| **CQRS** | Command Query Responsibility Segregation — แยก read/write model |
| **Event Sourcing** | เก็บ state เป็นลำดับของ event ที่ replay ได้ |
| **Hexagonal Architecture** | Ports & Adapters — แยก core logic จาก infrastructure |
| **Microservices** | สถาปัตยกรรมที่แยก service เล็ก ๆ ตาม business capability |
| **CI/CD** | Continuous Integration / Continuous Delivery-Deployment |

## 2.2 คำศัพท์ด้าน IoT

| คำศัพท์ | ความหมาย |
|---|---|
| **MQTT** | Message Queuing Telemetry Transport — pub/sub protocol สำหรับ IoT |
| **QoS** | Quality of Service — ระดับการรับประกันการส่ง message (0, 1, 2) |
| **LWT** | Last Will and Testament — message ที่ broker ส่งเมื่อ device หลุด |
| **Edge Computing** | ประมวลผลที่ขอบ network ใกล้ device |
| **Digital Twin** | แบบจำลองเสมือนของ physical asset |
| **Telemetry** | ข้อมูลที่ device ส่งมาจากระยะไกล |
| **OTA Update** | Over-The-Air firmware update |
| **Provisioning** | กระบวนการตั้งค่า device ครั้งแรก |

## 2.3 คำศัพท์ด้าน Protocol

| คำศัพท์ | ความหมาย |
|---|---|
| **REST** | Representational State Transfer — architectural style บน HTTP |
| **gRPC** | Google RPC — framework ใช้ HTTP/2 + Protobuf |
| **GraphQL** | Query language ที่ client กำหนด response shape |
| **WebSocket** | Protocol full-duplex บน TCP |
| **SSE** | Server-Sent Events — server push one-way ผ่าน HTTP |
| **Webhook** | HTTP callback เมื่อเกิด event |
| **Protobuf** | Protocol Buffers — binary serialization format |

## 2.4 คำศัพท์ด้าน Database

| คำศัพท์ | ความหมาย |
|---|---|
| **OLTP** | Online Transaction Processing — ระบบธุรกรรม |
| **OLAP** | Online Analytical Processing — ระบบวิเคราะห์ |
| **Time-Series DB** | ฐานข้อมูลที่ optimize สำหรับ time-series data |
| **Vector DB** | ฐานข้อมูลที่เก็บ embedding vector สำหรับ similarity search |
| **Sharding** | แบ่งข้อมูลออกเป็น shard ตาม key |
| **Replication** | ทำสำเนาข้อมูลไปยัง node อื่น |
| **ACID** | Atomicity, Consistency, Isolation, Durability |
| **CAP Theorem** | Consistency, Availability, Partition tolerance — เลือกได้ 2 จาก 3 |

## 2.5 คำศัพท์ด้าน Algorithm

| คำศัพท์ | ความหมาย |
|---|---|
| **Big O** | สัญลักษณ์อธิบาย upper bound ของ complexity |
| **Suffix Automaton** | DFA ขนาดเล็กที่สุดที่ยอมรับ suffix ทุกตัว |
| **Palindromic Automaton** | Eertree — โครงสร้างเก็บ palindrome ทั้งหมด |
| **Mo's Algorithm** | Offline query algorithm ที่จัดเรียง query เพื่อลด pointer movement |
| **Centroid Decomposition** | แยก tree เป็น centroid ซ้ำ ๆ depth O(log n) |
| **HLD** | Heavy-Light Decomposition — แยก tree เป็น chain |
| **Fenwick Tree** | Binary Indexed Tree — โครงสร้าง query prefix sum O(log n) |
| **Wavelet Tree** | โครงสร้างสำหรับ range k-th smallest query |
| **Lazy Propagation** | เทคนิค segment tree ที่เลื่อนการ update |

## 2.6 คำศัพท์ด้าน Monitoring

| คำศัพท์ | ความหมาย |
|---|---|
| **Metrics** | ตัวเลขที่วัดได้ เช่น QPS, latency, error rate |
| **Logs** | ข้อความเหตุการณ์ที่เกิดขึ้น |
| **Traces** | การติดตาม request ข้าม service |
| **Prometheus** | Monitoring system + time-series DB |
| **ELK Stack** | Elasticsearch, Logstash, Kibana |
| **OpenTelemetry** | Standard สำหรับ traces, metrics, logs |
| **SLO/SLI/SLA** | Service Level Objective/Indicator/Agreement |
| **MTTR/MTBF** | Mean Time To Recovery / Between Failures |

## 2.7 คำศัพท์ด้าน Business

| คำศัพท์ | ความหมาย |
|---|---|
| **SaaS** | Software as a Service |
| **ERP** | Enterprise Resource Planning |
| **CRM** | Customer Relationship Management |
| **Business Model Canvas** | เครื่องมือวางแผนธุรกิจ 9 blocks |
| **TAM/SAM/SOM** | Total/Serviceable/Obtainable Market |
| **CAC/LTV** | Customer Acquisition Cost / Lifetime Value |
| **MRR/ARR** | Monthly/Annual Recurring Revenue |

---

# 3. บทหัวข้อ (Main Topics)

## 3.1 Software Architecture

### 3.1.1 คืออะไร

**Software Architecture** คือโครงสร้างระดับสูงของระบบซอฟต์แวร์ ที่กำหนด:
- **องค์ประกอบ** (Components): ส่วนประกอบหลักของระบบ
- **ความสัมพันธ์** (Relationships): การสื่อสารระหว่างองค์ประกอบ
- **หลักการกำกับ** (Principles): กฎที่ใช้ตัดสินใจ

Architecture ที่ดีต้องตอบคำถาม:
- ระบบรองรับ scale แค่ไหน?
- เมื่อ component หนึ่งล่ม ระบบยังทำงานได้ไหม?
- เพิ่ม feature ใหม่ใช้เวลานานแค่ไหน?
- ทีมใหม่เข้าใจระบบได้เร็วแค่ไหน?

### 3.1.2 ทำงานอย่างไร

Software Architecture ทำงานเป็น **ชุดของ decision** ที่มีผลระยะยาว:

```
[Business Requirements]
        ↓
[Quality Attributes] (scalability, security, maintainability)
        ↓
[Architecture Decision] (patterns, tech stack, boundaries)
        ↓
[Implementation Guidelines]
        ↓
[Code + Infrastructure]
        ↓
[Feedback Loop] → ปรับ architecture
```

**ตัวอย่าง decision hierarchy:**
1. **Strategic**: Monolith vs Microservices
2. **Tactical**: REST vs gRPC ระหว่าง service
3. **Operational**: Kubernetes vs ECS

### 3.1.3 การใช้งาน

**ใน IoT Logistics:**
- แยก **Fleet Service**, **Order Service**, **Telemetry Service** เป็น bounded context
- ใช้ **Event-Driven** สำหรับ real-time tracking
- ใช้ **Clean Architecture** ในแต่ละ service
- ใช้ **CQRS** แยก read model (dashboard) จาก write model (command)

**ใน Smart Parking:**
- **Edge Layer**: ประมวลผลที่กล้อง
- **Ingestion Layer**: MQTT + Kafka
- **Processing Layer**: State machine + Rule engine
- **Application Layer**: Dashboard + Mobile + LED

### 3.1.4 ข้อดี

- **Scalability**: ออกแบบให้รองรับ growth ได้
- **Maintainability**: code อ่านง่าย แก้ไขง่าย
- **Testability**: test ได้ทุกระดับ
- **Team Productivity**: ทีมทำงานขนานได้
- **Risk Reduction**: ระบุ single point of failure ก่อนเกิดจริง
- **Cost Optimization**: เลือก tech ให้เหมาะกับ workload

### 3.1.5 ข้อเสีย

- **Over-engineering**: ออกแบบซับซ้อนเกินความจำเป็น
- **Upfront Cost**: ใช้เวลาออกแบบนาน
- **Resistance to Change**: architecture แข็งตัว ยากต่อการปรับ
- **Learning Curve**: ทีมใหม่ต้องเรียนรู้
- **Analysis Paralysis**: ตัดสินใจช้าเพราะ options เยอะ

### 3.1.6 ข้อควรระวัง / ข้อห้าม / ข้อจำกัด

**ข้อควรระวัง:**
- อย่าเลือก tech ตาม hype
- อย่าออกแบบโดยไม่มี requirement ชัด
- อย่า ignore non-functional requirements
- อย่า copy architecture จากบริษัทอื่นโดยไม่ปรับ

**ข้อห้าม:**
- ❌ ห้าม share database ระหว่าง microservice
- ❌ ห้าม synchronous call chain ยาวเกิน 3 hops
- ❌ ห้ามเก็บ state ใน service ที่ต้อง scale
- ❌ ห้าม ignore security ในชั้น architecture
- ❌ ห้าม deploy โดยไม่มี rollback plan

**ข้อจำกัด:**
- CAP theorem: เลือกได้ 2 จาก 3
- Network latency: มีขีดจำกัดทางฟิสิกส์
- Team size: architecture ต้องเหมาะกับทีม
- Budget: architecture ต้องเหมาะกับงบ
- Legacy: ระบบเดิมมีข้อจำกัด

### 3.1.7 Workflow การออกแบบ Architecture

```
┌─────────────────────────────────────────────┐
│ 1. Gather Requirements                      │
│    - Functional / Non-functional            │
│    - Constraints (budget, time, team)       │
└──────────────────┬──────────────────────────┘
                   ↓
┌─────────────────────────────────────────────┐
│ 2. Identify Quality Attributes              │
│    - Scalability, Availability, Security    │
│    - Prioritize (ไม่สามารถ optimize ทั้งหมด)│
└──────────────────┬──────────────────────────┘
                   ↓
┌─────────────────────────────────────────────┐
│ 3. Define Boundaries                        │
│    - Bounded Contexts (DDD)                 │
│    - Service boundaries                     │
└──────────────────┬──────────────────────────┘
                   ↓
┌─────────────────────────────────────────────┐
│ 4. Choose Patterns                          │
│    - Layered, Hexagonal, EDA, CQRS          │
└──────────────────┬──────────────────────────┘
                   ↓
┌─────────────────────────────────────────────┐
│ 5. Design Data Flow                         │
│    - Sync vs Async                          │
│    - Protocol selection                     │
└──────────────────┬──────────────────────────┘
                   ↓
┌─────────────────────────────────────────────┐
│ 6. Design Data Model                        │
│    - Database per service                   │
│    - Schema evolution strategy              │
└──────────────────┬──────────────────────────┘
                   ↓
┌─────────────────────────────────────────────┐
│ 7. Address Cross-cutting Concerns           │
│    - Security, Observability, Resilience    │
└──────────────────┬──────────────────────────┘
                   ↓
┌─────────────────────────────────────────────┐
│ 8. Document Decisions (ADR)                 │
│    - Context, Decision, Consequences        │
└──────────────────┬──────────────────────────┘
                   ↓
┌─────────────────────────────────────────────┐
│ 9. Prototype & Validate                     │
│    - PoC critical paths                     │
└──────────────────┬──────────────────────────┘
                   ↓
┌─────────────────────────────────────────────┐
│ 10. Iterate                                 │
│    - Architecture ไม่ใช่ one-time decision  │
└─────────────────────────────────────────────┘
```

### 3.1.8 สรุป

Software Architecture คือ **ชุดการตัดสินใจ** ที่กำหนดคุณภาพระยะยาวของระบบ ต้องสมดุลระหว่าง:
- **Simplicity** vs **Capability**
- **Speed** vs **Quality**
- **Cost** vs **Scalability**
- **Innovation** vs **Stability**

ไม่มี architecture ที่ perfect — มีแต่ architecture ที่ **เหมาะกับ context** ที่สุด

---

## 3.2 IoT Solution Architecture

### 3.2.1 คืออะไร

**IoT Solution Architecture** คือสถาปัตยกรรมที่เชื่อมต่อ:
- **Physical World**: sensors, actuators, cameras
- **Edge Layer**: ประมวลผลใกล้ device
- **Network Layer**: สื่อสารข้อมูล
- **Cloud Layer**: ประมวลผล + เก็บข้อมูล
- **Application Layer**: user-facing

### 3.2.2 ทำงานอย่างไร

```
┌────────────────────────────────────────────────────┐
│  DEVICE LAYER                                       │
│  Sensors, Actuators, Cameras, PLC, Meters          │
└──────────────────────┬─────────────────────────────┘
                       │ (MQTT / CoAP / Modbus / OPC-UA)
                       ↓
┌────────────────────────────────────────────────────┐
│  EDGE LAYER                                         │
│  Gateway, Edge AI, Local Storage, Protocol Adapter │
└──────────────────────┬─────────────────────────────┘
                       │ (MQTT / HTTPS / gRPC)
                       ↓
┌────────────────────────────────────────────────────┐
│  INGESTION LAYER                                    │
│  MQTT Broker (EMQX), API Gateway, Load Balancer    │
└──────────────────────┬─────────────────────────────┘
                       │ (Kafka / NATS / RabbitMQ)
                       ↓
┌────────────────────────────────────────────────────┐
│  PROCESSING LAYER                                   │
│  Stream Processor (Flink), Rule Engine, ML Inference│
└──────────────────────┬─────────────────────────────┘
                       │
        ┌──────────────┼──────────────┐
        ↓              ↓              ↓
┌──────────────┐ ┌──────────────┐ ┌──────────────┐
│ TIME-SERIES  │ │  RELATIONAL  │ │   DOCUMENT   │
│  InfluxDB    │ │  PostgreSQL  │ │   MongoDB    │
└──────────────┘ └──────────────┘ └──────────────┘
        ↓              ↓              ↓
        └──────────────┼──────────────┘
                       ↓
┌────────────────────────────────────────────────────┐
│  APPLICATION LAYER                                  │
│  REST / gRPC / GraphQL / WebSocket / SSE           │
└──────────────────────┬─────────────────────────────┘
                       ↓
┌────────────────────────────────────────────────────┐
│  PRESENTATION LAYER                                 │
│  Web, Mobile, Dashboard, LED, Voice, AR            │
└────────────────────────────────────────────────────┘
```

### 3.2.3 การใช้งาน

**Use Case 1: Fleet Management**
- GPS tracker → MQTT → Kafka → InfluxDB (trajectory)
- Metadata → PostgreSQL
- Latest position → Redis Geo
- Alert → Rule engine → Notification

**Use Case 2: Smart Parking**
- Camera → Edge AI → MQTT → Cloud
- Slot state → Redis + PostgreSQL
- Real-time → WebSocket → Mobile
- Payment → Stripe/Omise

**Use Case 3: Industrial IoT**
- PLC → OPC-UA → Edge → MQTT → Cloud
- Anomaly detection → ML model
- Predictive maintenance → Alert

### 3.2.4 ข้อดี

- **Real-time visibility**: เห็นสถานะทันที
- **Automation**: ลดแรงงานคน
- **Predictive**: คาดการณ์ล่วงหน้า
- **Efficiency**: ใช้ทรัพยากรคุ้มค่า
- **New Business Model**: จาก product → service
- **Data-driven**: ตัดสินใจจากข้อมูล

### 3.2.5 ข้อเสีย

- **Complexity**: หลาย layer ต้องเชี่ยวชาญ
- **Security Risk**: attack surface กว้าง
- **Cost**: hardware + cloud + maintenance
- **Integration**: legacy system ยาก
- **Data Volume**: จัดการข้อมูลมหาศาล
- **Vendor Lock-in**: ผูกกับ vendor

### 3.2.6 ข้อควรระวัง / ข้อห้าม / ข้อจำกัด

**ข้อควรระวัง:**
- Device firmware update ต้องมี rollback
- Network partition ต้องจัดการได้
- Clock skew ระหว่าง device
- Data validation ที่ edge
- Certificate rotation

**ข้อห้าม:**
- ❌ ห้ามเก็บ secret ใน firmware
- ❌ ห้ามใช้ default password
- ❌ ห้าม ignore OTA security
- ❌ ห้าม trust device โดยไม่ authenticate
- ❌ ห้ามออกแบบโดยไม่มี offline mode

**ข้อจำกัด:**
- Bandwidth: device อาจมีข้อจำกัด
- Power: battery-powered device
- Compute: edge device มี RAM/CPU จำกัด
- Latency: network ไม่เสถียร
- Cost: data transfer แพง

### 3.2.7 Workflow

```
[Device Boot] → [Provisioning] → [Cert Exchange]
       ↓
[Connect MQTT] → [Subscribe Config Topic]
       ↓
[Read Sensor] → [Local Validation] → [Publish Telemetry]
       ↓
[Receive Command] → [Execute] → [Ack]
       ↓
[Periodic OTA Check] → [Download] → [Verify] → [Install]
       ↓
[Health Report] → [Alert ถ้าผิดปกติ]
```

### 3.2.8 สรุป

IoT Solution Architecture ต้องสมดุลระหว่าง:
- **Edge vs Cloud**: ประมวลผลที่ไหน
- **Sync vs Async**: สื่อสารอย่างไร
- **Push vs Pull**: ใครเริ่ม
- **Real-time vs Batch**: ต้องการเร็วแค่ไหน
- **Cost vs Performance**: คุ้มไหม

---

## 3.3 Protocol Comparison (เชิงลึก)

### 3.3.1 คืออะไร

**Protocol** คือข้อกำหนดการสื่อสารระหว่างระบบ ในบริบท IoT/API มีหลาย protocol ที่เหมาะกับ use case ต่างกัน

### 3.3.2 เปรียบเทียบเชิงลึก

#### 3.3.2.1 Fundamentals

| Protocol | Model | Transport | Architecture | Key Innovation |
|---|---|---|---|---|
| REST | Request-Response | HTTP/1.1, HTTP/2 | Stateless, resource-based | Standard HTTP verbs |
| gRPC | Unary + Streaming | HTTP/2 | RPC, contract-first | Multiplexing, header compression |
| GraphQL | Flexible Queries | HTTP/1.1 | Schema-first | Client-driven queries |
| WebSockets | Full-Duplex | TCP upgraded | Connection-oriented | Persistent, bidirectional |
| SSE | Server Push | HTTP/1.1 | One-way stream | Auto-reconnect |
| SOAP | Request-Response | HTTP, SMTP | Strict XML contract | WS-* standards |
| MQTT | Pub/Sub | TCP | Broker-based | Lightweight, QoS |
| Webhook | Callback | HTTP | Event-driven | Push notification |

#### 3.3.2.2 Data Handling

| Protocol | Format | Size | Human-readable |
|---|---|---|---|
| REST | JSON/XML | Medium | Yes |
| gRPC | Protobuf | Small | No |
| GraphQL | JSON | Medium | Yes |
| WebSocket | Any | Depends | Depends |
| SOAP | XML | Large | Yes |
| MQTT | Binary/JSON | Small | Depends |

#### 3.3.2.3 Complexity

| Protocol | Setup | Schema | Learning |
|---|---|---|---|
| REST | Low | Optional | Easy |
| gRPC | Medium | Strict | Moderate |
| GraphQL | Medium | Strong | Moderate |
| WebSocket | Medium | Custom | Moderate |
| SOAP | High | Strict | Steep |
| MQTT | Low | Topic | Easy |

#### 3.3.2.4 Performance

| Protocol | Latency | Throughput | CPU | Bandwidth |
|---|---|---|---|---|
| REST | Medium | Medium | Medium | High |
| gRPC | Low | High | Low | Low |
| GraphQL | Medium | Medium | High | Medium |
| WebSocket | Low | High | Medium | Low |
| MQTT | Low | High | Low | Low |
| SOAP | High | Low | High | High |

#### 3.3.2.5 State Management

| Protocol | State | Scaling |
|---|---|---|
| REST | Stateless | Easy |
| gRPC | Hybrid | Medium |
| GraphQL | Stateless | Easy |
| WebSocket | Stateful | Hard |
| MQTT | Stateful session | Medium |

#### 3.3.2.6 Security

| Protocol | Transport | Auth | Message-level |
|---|---|---|---|
| REST | HTTPS | OAuth/JWT/Key | No |
| gRPC | TLS | Token/mTLS | No |
| GraphQL | HTTPS | OAuth/JWT | No |
| WebSocket | WSS | Custom | No |
| MQTT | TLS | User/Pass/Cert | No |
| SOAP | HTTPS | WS-Security | Yes |

#### 3.3.2.7 Use Case Alignment

| Use Case | Best Option | Alternative |
|---|---|---|
| CRUD API | REST | GraphQL |
| Microservice-to-microservice | gRPC | REST |
| Real-time chat | WebSocket | SSE |
| Data aggregation | GraphQL | REST |
| Legacy enterprise | SOAP | REST |
| IoT device telemetry | MQTT | CoAP |
| Server push | SSE | WebSocket |
| Event notification | Webhook | MQTT |

#### 3.3.2.8 Database Integration

| Protocol | DB Fit |
|---|---|
| REST | Relational/NoSQL natural |
| gRPC | Internal service + low-latency |
| GraphQL | Federated + BFF |
| WebSocket | Backend-managed state |
| SOAP | Legacy + transactional |
| MQTT | Time-series + event |

#### 3.3.2.9 Versioning

| Protocol | Strategy |
|---|---|
| REST | URI /v1/, header, content negotiation |
| gRPC | Protobuf backward compat |
| GraphQL | Deprecate fields, additive |
| WebSocket | Message contract versioning |
| SOAP | WSDL namespace |
| MQTT | Topic versioning |

#### 3.3.2.10 Testing & Debugging

| Protocol | Tools | Difficulty |
|---|---|---|
| REST | curl, Postman, Swagger | Easy |
| gRPC | grpcurl, BloomRPC, Evans | Hard |
| GraphQL | GraphiQL, Apollo Studio | Medium |
| WebSocket | Browser DevTools, wscat | Medium |
| MQTT | MQTTX, mosquitto_pub | Easy |
| SOAP | SoapUI | Hard |

#### 3.3.2.11 Cost of Ownership

| Factor | REST | gRPC | GraphQL | WebSocket | MQTT |
|---|---|---|---|---|---|
| Training | Low | High | Medium | Medium | Low |
| Debugging | Easy | Hard | Medium | Medium | Easy |
| Infra | Medium | Low CPU | High CPU | Persistent | Low |
| Maintenance | Low | Medium | High | Medium-High | Low |

#### 3.3.2.12 TL;DR Summary

| Criteria | REST | gRPC | GraphQL | WebSocket | MQTT |
|---|---|---|---|---|---|
| Format | JSON | Protobuf | JSON | Custom | Binary |
| Performance | Medium | High | Medium | High | High |
| Real-time | No | Stream | Subscription | Yes | Yes |
| Browser | Yes | Proxy | Yes | Yes | No |
| Ideal | CRUD | Microservice | Data fetch | Real-time | IoT |

### 3.3.3 การใช้งานใน IoT

**เลือก Protocol ตาม Layer:**

```
Device → Edge:        MQTT, CoAP, Modbus, OPC-UA
Edge → Cloud:         MQTT, HTTPS, gRPC
Cloud → Cloud:        gRPC, Kafka, NATS
Cloud → App:          REST, GraphQL, gRPC
Cloud → App (real):   WebSocket, SSE
Cloud → External:     Webhook
```

### 3.3.4 ข้อดี / ข้อเสีย (โดยรวม)

**ข้อดีของการมีหลาย Protocol:**
- เลือกให้เหมาะกับ use case
- Optimize performance
- รองรับ client หลายประเภท
- Evolve ได้ตามเทคโนโลยี

**ข้อเสีย:**
- Complexity สูง
- ต้องเชี่ยวชาญหลายอย่าง
- Maintenance หลาย stack
- Debug ยากขึ้น
- Security surface กว้าง

### 3.3.5 ข้อควรระวัง

- **Over-engineering**: อย่าใช้ทุก protocol ในระบบเดียว
- **Consistency**: เลือกให้เหมาะกับทีม
- **Security**: ทุก protocol ต้องมี auth
- **Observability**: ต้อง trace ได้ข้าม protocol
- **Versioning**: ต้องมี plan ชัดเจน
- **Cost**: หลาย protocol = หลาย infra

### 3.3.6 Workflow การเลือก Protocol

```
┌──────────────────────────────┐
│ ระบุ Communication Pattern   │
│ - Request/Response?          │
│ - Stream?                    │
│ - Pub/Sub?                   │
│ - Real-time?                 │
└──────────────┬───────────────┘
               ↓
┌──────────────────────────────┐
│ ระบุ Constraint              │
│ - Bandwidth?                 │
│ - Latency?                   │
│ - Battery?                   │
│ - Browser support?           │
└──────────────┬───────────────┘
               ↓
┌──────────────────────────────┐
│ ระบุ Team Capability         │
│ - Expertise?                 │
│ - Tooling?                   │
└──────────────┬───────────────┘
               ↓
┌──────────────────────────────┐
│ เลือก Protocol               │
│ + Fallback                   │
└──────────────┬───────────────┘
               ↓
┌──────────────────────────────┐
│ Prototype & Benchmark        │
└──────────────┬───────────────┘
               ↓
┌──────────────────────────────┐
│ Document & Train Team        │
└──────────────────────────────┘
```

### 3.3.7 สรุป

ไม่มี protocol ที่ดีที่สุด — มีแต่ **เหมาะกับ use case** ที่สุด กฎง่าย ๆ:
- **IoT device** → MQTT
- **Microservice** → gRPC
- **Public API** → REST
- **Data aggregation** → GraphQL
- **Real-time** → WebSocket/SSE
- **Event notification** → Webhook
- **Async processing** → Message Queue

---

## 3.4 Database Design (PostgreSQL ครบทุกตาราง)

### 3.4.1 คืออะไร

**Database Design** คือกระบวนการออกแบบโครงสร้างข้อมูลให้:
- รองรับ query patterns
- รักษา data integrity
- Scale ได้
- Performance ดี
- Maintainable

### 3.4.2 ทำงานอย่างไร

**หลักการ:**
1. **Normalization**: ลด redundancy (1NF → 3NF)
2. **Denormalization**: เพิ่ม performance เมื่อจำเป็น
3. **Indexing**: เพิ่มความเร็ว query
4. **Partitioning**: แบ่งตารางใหญ่
5. **Sharding**: แบ่งข้าม server

**Schema ครบทุกตาราง (Logistics + IoT):**

#### 3.4.2.1 Multi-tenancy

```sql
CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    plan VARCHAR(50) NOT NULL DEFAULT 'free',
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    settings JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255),
    phone VARCHAR(50),
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    status VARCHAR(20) DEFAULT 'active',
    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(tenant_id, email)
);
CREATE INDEX idx_users_tenant ON users(tenant_id);
CREATE INDEX idx_users_email ON users(email);

CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    permissions JSONB NOT NULL DEFAULT '[]'
);

CREATE TABLE user_roles (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    role_id INT REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);
```

#### 3.4.2.2 Sites & Devices

```sql
CREATE TABLE sites (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    name VARCHAR(255) NOT NULL,
    address TEXT,
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    timezone VARCHAR(50) DEFAULT 'Asia/Bangkok',
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_sites_tenant ON sites(tenant_id);

CREATE TABLE device_types (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    category VARCHAR(50),
    schema JSONB NOT NULL DEFAULT '{}'
);

CREATE TABLE devices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    site_id UUID REFERENCES sites(id),
    device_type_id INT REFERENCES device_types(id),
    serial_number VARCHAR(100) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    model VARCHAR(100),
    firmware_version VARCHAR(50),
    status VARCHAR(20) DEFAULT 'offline',
    last_seen_at TIMESTAMPTZ,
    config JSONB DEFAULT '{}',
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_devices_tenant ON devices(tenant_id);
CREATE INDEX idx_devices_site ON devices(site_id);
CREATE INDEX idx_devices_status ON devices(status);
CREATE INDEX idx_devices_last_seen ON devices(last_seen_at);

CREATE TABLE device_credentials (
    device_id UUID PRIMARY KEY REFERENCES devices(id) ON DELETE CASCADE,
    mqtt_username VARCHAR(100) UNIQUE,
    mqtt_password_hash VARCHAR(255),
    api_key_hash VARCHAR(255),
    client_cert TEXT,
    rotated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE device_groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(id),
    name VARCHAR(255) NOT NULL,
    description TEXT
);

CREATE TABLE device_group_members (
    group_id UUID REFERENCES device_groups(id) ON DELETE CASCADE,
    device_id UUID REFERENCES devices(id) ON DELETE CASCADE,
    PRIMARY KEY (group_id, device_id)
);

CREATE TABLE telemetry_schemas (
    id SERIAL PRIMARY KEY,
    device_type_id INT REFERENCES device_types(id),
    metric_name VARCHAR(100) NOT NULL,
    unit VARCHAR(20),
    data_type VARCHAR(20) NOT NULL,
    min_value DECIMAL,
    max_value DECIMAL,
    UNIQUE(device_type_id, metric_name)
);

CREATE TABLE device_latest_state (
    device_id UUID PRIMARY KEY REFERENCES devices(id) ON DELETE CASCADE,
    state JSONB NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

#### 3.4.2.3 Customers & Orders

```sql
CREATE TABLE customers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    code VARCHAR(50),
    name VARCHAR(255) NOT NULL,
    tax_id VARCHAR(50),
    email VARCHAR(255),
    phone VARCHAR(50),
    billing_address JSONB,
    shipping_address JSONB,
    credit_limit DECIMAL(15,2) DEFAULT 0,
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_customers_tenant ON customers(tenant_id);

CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    order_no VARCHAR(50) UNIQUE NOT NULL,
    customer_id UUID REFERENCES customers(id),
    status VARCHAR(30) NOT NULL DEFAULT 'pending',
    priority INT DEFAULT 0,
    pickup_address JSONB NOT NULL,
    pickup_lat DECIMAL(10,8),
    pickup_lng DECIMAL(11,8),
    pickup_at TIMESTAMPTZ,
    delivery_address JSONB NOT NULL,
    delivery_lat DECIMAL(10,8),
    delivery_lng DECIMAL(11,8),
    delivery_at TIMESTAMPTZ,
    weight_kg DECIMAL(10,2),
    volume_m3 DECIMAL(10,4),
    total_amount DECIMAL(15,2),
    notes TEXT,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_orders_tenant_status ON orders(tenant_id, status);
CREATE INDEX idx_orders_customer ON orders(customer_id);
CREATE INDEX idx_orders_created ON orders(created_at DESC);

CREATE TABLE order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID REFERENCES orders(id) ON DELETE CASCADE,
    sku VARCHAR(100),
    description TEXT,
    quantity INT NOT NULL DEFAULT 1,
    unit_price DECIMAL(15,2),
    total_price DECIMAL(15,2)
);

CREATE TABLE order_status_history (
    id BIGSERIAL PRIMARY KEY,
    order_id UUID REFERENCES orders(id) ON DELETE CASCADE,
    from_status VARCHAR(30),
    to_status VARCHAR(30) NOT NULL,
    changed_by UUID REFERENCES users(id),
    reason TEXT,
    changed_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_order_history_order ON order_status_history(order_id, changed_at DESC);
```

#### 3.4.2.4 Fleet & Trips

```sql
CREATE TABLE vehicles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    plate_no VARCHAR(20) UNIQUE NOT NULL,
    type VARCHAR(50),
    brand VARCHAR(100),
    model VARCHAR(100),
    year INT,
    capacity_kg DECIMAL(10,2),
    capacity_m3 DECIMAL(10,4),
    fuel_type VARCHAR(20),
    status VARCHAR(20) DEFAULT 'available',
    device_id UUID REFERENCES devices(id),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE drivers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    user_id UUID REFERENCES users(id),
    license_no VARCHAR(50) UNIQUE,
    license_expiry DATE,
    phone VARCHAR(50),
    status VARCHAR(20) DEFAULT 'available',
    rating DECIMAL(3,2) DEFAULT 0
);

CREATE TABLE trips (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    trip_no VARCHAR(50) UNIQUE NOT NULL,
    vehicle_id UUID REFERENCES vehicles(id),
    driver_id UUID REFERENCES drivers(id),
    status VARCHAR(30) DEFAULT 'planned',
    planned_start TIMESTAMPTZ,
    actual_start TIMESTAMPTZ,
    planned_end TIMESTAMPTZ,
    actual_end TIMESTAMPTZ,
    total_distance_km DECIMAL(10,2),
    total_stops INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_trips_vehicle ON trips(vehicle_id, actual_start DESC);

CREATE TABLE trip_stops (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    trip_id UUID REFERENCES trips(id) ON DELETE CASCADE,
    sequence INT NOT NULL,
    order_id UUID REFERENCES orders(id),
    stop_type VARCHAR(20) NOT NULL,
    address JSONB,
    lat DECIMAL(10,8),
    lng DECIMAL(11,8),
    planned_arrival TIMESTAMPTZ,
    actual_arrival TIMESTAMPTZ,
    status VARCHAR(20) DEFAULT 'pending',
    UNIQUE(trip_id, sequence)
);

CREATE TABLE trip_routes (
    id BIGSERIAL PRIMARY KEY,
    trip_id UUID REFERENCES trips(id) ON DELETE CASCADE,
    recorded_at TIMESTAMPTZ NOT NULL,
    lat DECIMAL(10,8) NOT NULL,
    lng DECIMAL(11,8) NOT NULL,
    speed_kmh DECIMAL(6,2),
    heading DECIMAL(5,2)
);
CREATE INDEX idx_trip_routes_trip ON trip_routes(trip_id, recorded_at DESC);
```

#### 3.4.2.5 Alerts & Rules

```sql
CREATE TABLE alert_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    name VARCHAR(255) NOT NULL,
    device_id UUID REFERENCES devices(id),
    device_type_id INT REFERENCES device_types(id),
    metric VARCHAR(100) NOT NULL,
    operator VARCHAR(10) NOT NULL,
    threshold DECIMAL(15,4),
    duration_sec INT DEFAULT 0,
    severity VARCHAR(20) NOT NULL DEFAULT 'warning',
    enabled BOOLEAN DEFAULT true,
    actions JSONB DEFAULT '[]',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE alerts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    rule_id UUID REFERENCES alert_rules(id),
    device_id UUID REFERENCES devices(id),
    severity VARCHAR(20) NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT,
    value JSONB,
    status VARCHAR(20) DEFAULT 'open',
    acknowledged_by UUID REFERENCES users(id),
    acknowledged_at TIMESTAMPTZ,
    resolved_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_alerts_tenant_status ON alerts(tenant_id, status, created_at DESC);
CREATE INDEX idx_alerts_device ON alerts(device_id, created_at DESC);
```

#### 3.4.2.6 Audit & Billing

```sql
CREATE TABLE audit_logs (
    id BIGSERIAL PRIMARY KEY,
    tenant_id UUID,
    user_id UUID,
    action VARCHAR(100) NOT NULL,
    resource_type VARCHAR(50),
    resource_id VARCHAR(100),
    changes JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_audit_tenant ON audit_logs(tenant_id, created_at DESC);
CREATE INDEX idx_audit_resource ON audit_logs(resource_type, resource_id);

CREATE TABLE invoices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    invoice_no VARCHAR(50) UNIQUE NOT NULL,
    customer_id UUID REFERENCES customers(id),
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    subtotal DECIMAL(15,2) NOT NULL,
    tax DECIMAL(15,2) DEFAULT 0,
    total DECIMAL(15,2) NOT NULL,
    status VARCHAR(20) DEFAULT 'draft',
    due_date DATE,
    paid_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE invoice_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id UUID REFERENCES invoices(id) ON DELETE CASCADE,
    description TEXT NOT NULL,
    quantity DECIMAL(15,4),
    unit_price DECIMAL(15,2),
    amount DECIMAL(15,2)
);
```

#### 3.4.2.7 Row-Level Security

```sql
ALTER TABLE devices ENABLE ROW LEVEL SECURITY;

CREATE POLICY devices_tenant_isolation ON devices
    USING (tenant_id = current_setting('app.tenant_id')::UUID);
```

### 3.4.3 การใช้งาน

- **OLTP**: PostgreSQL สำหรับ transaction
- **Time-series**: InfluxDB สำหรับ telemetry
- **Document**: MongoDB สำหรับ logs/events
- **Cache**: Redis สำหรับ state
- **Vector**: สำหรับ embedding search

### 3.4.4 ข้อดี

- **ACID**: มั่นใจได้ใน transaction
- **Rich types**: JSONB, array, geo
- **Extensible**: extension เยอะ (PostGIS, TimescaleDB)
- **Mature**: community ใหญ่
- **Free**: open source

### 3.4.5 ข้อเสีย

- **Vertical scaling**: scale up มีขีดจำกัด
- **Sharding ยาก**: ต้องใช้ Citus หรือ app-level
- **Write throughput**: น้อยกว่า NoSQL
- **Schema change**: ALTER TABLE ใหญ่ใช้เวลานาน

### 3.4.6 ข้อควรระวัง

- ใช้ `EXPLAIN ANALYZE` ทุก query หนัก
- ระวัง N+1 query
- ตั้ง `statement_timeout`
- ใช้ connection pool (PgBouncer)
- Vacuum/Analyze สม่ำเสมอ
- Backup + PITR
- Monitor replication lag

### 3.4.7 Workflow การออกแบบ Schema

```
[Identify Entities] → [Define Relationships] → [Normalize to 3NF]
       ↓
[Identify Query Patterns] → [Add Indexes]
       ↓
[Estimate Volume] → [Partition/Shard Plan]
       ↓
[Security] → [RLS + Encryption]
       ↓
[Test with Production-like Data] → [Optimize]
       ↓
[Document + Migration Strategy]
```

### 3.4.8 สรุป

Database design ที่ดีต้อง:
- Model ตรงกับ domain
- รองรับ query patterns
- Scale ได้ตาม growth
- ปลอดภัย
- Maintainable

**Rule of thumb**: เริ่มจาก normalized, denormalize เมื่อจำเป็น, วัดผลก่อน optimize

---

## 3.5 Monitoring & Observability

### 3.5.1 คืออะไร

**Observability** คือความสามารถในการเข้าใจสถานะภายในของระบบจาก output ภายนอก ประกอบด้วย 3 pillars:
- **Metrics**: ตัวเลข
- **Logs**: เหตุการณ์
- **Traces**: การเดินทางของ request

### 3.5.2 ทำงานอย่างไร

```
[Application] → [Instrumentation] → [Collector]
                                          ↓
                    ┌─────────────────────┼─────────────────────┐
                    ↓                     ↓                     ↓
              [Prometheus]          [Elasticsearch]      [Jaeger/Tempo]
              (metrics)             (logs)               (traces)
                    ↓                     ↓                     ↓
              [Grafana]              [Kibana]              [Grafana]
                    ↓                     ↓                     ↓
              [Alertmanager]         [Alert]              [Alert]
                    ↓
              [LINE/Slack/PagerDuty]
```

**Metrics (Prometheus):**
```go
var httpRequests = prometheus.NewCounterVec(
    prometheus.CounterOpts{
        Name: "http_requests_total",
        Help: "Total HTTP requests",
    },
    []string{"method", "endpoint", "status"},
)

var httpDuration = prometheus.NewHistogramVec(
    prometheus.HistogramOpts{
        Name:    "http_request_duration_seconds",
        Buckets: prometheus.DefBuckets,
    },
    []string{"method", "endpoint"},
)
```

**Logs (ELK):**
```go
logger.Info("device connected",
    zap.String("device_id", "d-001"),
    zap.String("tenant_id", "t-001"),
    zap.String("ip", "10.0.0.1"),
    zap.Duration("latency", 23*time.Millisecond),
)
```

**Traces (OpenTelemetry):**
```go
ctx, span := otel.Tracer("device").Start(c.Request.Context(), "GetDevice")
defer span.End()
span.SetAttributes(
    attribute.String("device.id", c.Param("id")),
    attribute.String("tenant.id", c.GetString("tenant_id")),
)
```

### 3.5.3 การใช้งาน

- **Metrics**: QPS, latency, error rate, saturation
- **Logs**: debugging, audit, security
- **Traces**: bottleneck, dependency map
- **Alerting**: แจ้งเตือนเมื่อผิดปกติ
- **Dashboard**: ภาพรวมสำหรับทีม

### 3.5.4 ข้อดี

- **Visibility**: เห็นสถานะระบบ
- **Fast Debug**: หาสาเหตุเร็ว
- **Proactive**: แจ้งเตือนก่อนผู้ใช้บ่น
- **Capacity Planning**: วางแผน resource
- **SLO Tracking**: วัด SLA

### 3.5.5 ข้อเสีย

- **Cost**: เก็บ data เยอะแพง
- **Complexity**: หลาย component
- **Cardinality**: metrics สูงเกินไป
- **Noise**: alert เยอะเกิน
- **Learning curve**: ต้องเข้าใจ tooling

### 3.5.6 ข้อควรระวัง

- อย่าเก็บ high-cardinality labels
- Sampling traces เพื่อลด cost
- Retention policy ชัดเจน
- Alert ต้อง actionable
- อย่า log sensitive data
- Correlation ID ทุก request
- Monitor monitoring system

### 3.5.7 Workflow

```
[Define SLI] → [Instrument] → [Collect] → [Visualize]
       ↓
[Define SLO] → [Alert Rule] → [Runbook]
       ↓
[Incident] → [Detect] → [Diagnose] → [Mitigate] → [Post-mortem]
```

### 3.5.8 สรุป

Observability ที่ดีต้อง:
- Instrument ตั้งแต่ต้น
- Correlation ID ทุก layer
- Alert ที่ actionable
- Dashboard ต่อ role
- Cost-aware

---

# 4. ออกแบบ Workflow (Workflow Design)

## 4.1 Workflow หลักของระบบ IoT Logistics

```
┌─────────────────────────────────────────────────────────────┐
│ PHASE 1: DEVICE ONBOARDING                                   │
└─────────────────────────────────────────────────────────────┘
[Manufacture] → [Flash Firmware] → [Generate Cert]
       ↓
[Ship to Customer] → [Install] → [Power On]
       ↓
[Auto Provision] → [Register to Cloud] → [Assign to Tenant]
       ↓
[Verify] → [Active]

┌─────────────────────────────────────────────────────────────┐
│ PHASE 2: TELEMETRY INGESTION                                 │
└─────────────────────────────────────────────────────────────┘
[Read Sensor] → [Validate] → [Buffer Local]
       ↓
[Connect MQTT] → [Publish QoS 1]
       ↓
[EMQX Broker] → [Auth Check] → [Route to Kafka]
       ↓
[Kafka Topic: iot.telemetry] → [Partition by tenant+device]
       ↓
[Stream Processor] → [Parse] → [Enrich] → [Validate]
       ↓
   ┌───┴───┬────────┬────────┐
   ↓       ↓        ↓        ↓
[Influx][Postgres][Redis][Rule Engine]
   ↓       ↓        ↓        ↓
[Query][Metadata][Cache][Alert]

┌─────────────────────────────────────────────────────────────┐
│ PHASE 3: COMMAND & CONTROL                                   │
└─────────────────────────────────────────────────────────────┘
[User Action] → [REST API] → [Validate] → [Publish Command]
       ↓
[MQTT Topic: iot/{tenant}/{device}/command]
       ↓
[Device Receives] → [Validate] → [Execute]
       ↓
[Publish Ack] → [Update DB] → [Notify User]
       ↓
[Timeout?] → [Retry] → [Dead Letter]

┌─────────────────────────────────────────────────────────────┐
│ PHASE 4: ORDER LIFECYCLE                                     │
└─────────────────────────────────────────────────────────────┘
[Create Order] → [Validate] → [Save DB]
       ↓
[Assign Vehicle + Driver] → [Create Trip]
       ↓
[Optimize Route] → [Push to Driver App]
       ↓
[Track Real-time] → [Update Status]
       ↓
[Pickup] → [In Transit] → [Delivered]
       ↓
[Proof of Delivery] → [Invoice] → [Close]

┌─────────────────────────────────────────────────────────────┐
│ PHASE 5: ALERT & INCIDENT                                    │
└─────────────────────────────────────────────────────────────┘
[Rule Match] → [Create Alert] → [Severity Classify]
       ↓
[Notify: WebSocket + Push + LINE]
       ↓
[Acknowledge] → [Investigate] → [Resolve]
       ↓
[Post-mortem] → [Update Rule]
```

## 4.2 Workflow Diagram (Swimlane)

```
┌──────────┬──────────┬──────────┬──────────┬──────────┐
│ Device   │ Edge     │ Cloud    │ App      │ User     │
├──────────┼──────────┼──────────┼──────────┼──────────┤
│ Read     │          │          │          │          │
│ Sensor   │          │          │          │          │
│    ↓     │          │          │          │          │
│ Publish  │ Receive  │          │          │          │
│ MQTT ────┼──→       │          │          │          │
│          │ Validate │          │          │          │
│          │    ↓     │          │          │          │
│          │ Forward ─┼──→ Kafka │          │          │
│          │          │    ↓     │          │          │
│          │          │ Process  │          │          │
│          │          │    ↓     │          │          │
│          │          │ Store    │          │          │
│          │          │    ↓     │          │          │
│          │          │ Push ────┼──→ WS    │          │
│          │          │          │    ↓     │          │
│          │          │          │ Display  │ View     │
│          │          │          │          │    ↓     │
│          │          │          │          │ Action   │
│          │          │ Receive ─┼──← REST  │    │     │
│          │          │    ↓     │          │    │     │
│          │          │ Command ─┼──→ MQTT  │    │     │
│          │ Receive  │          │          │    │     │
│ Execute  │          │          │          │    │     │
│    ↓     │          │          │          │    │     │
│ Ack ─────┼──→       │ Update   │          │    │     │
│          │          │    ↓     │          │    │     │
│          │          │ Notify ──┼──→ WS ───┼────┘     │
└──────────┴──────────┴──────────┴──────────┴──────────┘
```

---

# 5. Case Study

## 5.1 Case Study 1: Logistics SaaS — Fleet 5,000 คัน

### 5.1.1 บริบท
- ลูกค้า: บริษัท logistics ขนาดกลางใน SEA
- Fleet: 5,000 คัน, 3,000 คนขับ
- Orders: 50,000/วัน
- ต้องการ: real-time tracking, ETA prediction, alert

### 5.1.2 ความท้าทาย
- GPS จาก 5 vendor ต่างกัน format
- Peak 80,000 msg/s
- ต้อง multi-tenant
- ต้องการ 99.95% uptime
- งบจำกัด

### 5.1.3 Architecture ที่เลือก

```
[GPS Device] → MQTT (EMQX Cluster 3 nodes)
                     ↓
              Kafka (10 partitions, RF=3)
                     ↓
              Go Stream Processor (10 pods)
                     ↓
    ┌────────────────┼────────────────┐
    ↓                ↓                ↓
[InfluxDB]      [PostgreSQL]      [Redis Geo]
(trajectory)    (metadata)        (latest)
    ↓                ↓                ↓
    └────────────────┼────────────────┘
                     ↓
              API Gateway (Kong)
                     ↓
    ┌────────────────┼────────────────┐
    ↓                ↓                ↓
[REST]          [WebSocket]      [gRPC]
    ↓                ↓                ↓
[Web App]       [Live Map]       [Mobile]
```

### 5.1.4 ผลลัพธ์
| Metric | Before | After |
|---|---|---|
| Latency P95 | 5s | 180ms |
| Throughput | 5k msg/s | 80k msg/s |
| Uptime | 99.5% | 99.97% |
| MTTR | 30 min | 5 min |
| Cost/msg | $0.001 | $0.0004 |

### 5.1.5 บทเรียน
- Partition Kafka ตาม tenant + device
- Redis Geo สำหรับ "nearby" query
- Protobuf ลด bandwidth 60%
- Tiered storage ลด cost 40%
- Circuit breaker ป้องกัน cascade

---

## 5.2 Case Study 2: Smart Parking 20,000 ช่องจอด

### 5.2.1 บริบท
- ลูกค้า: ห้างสรรพสินค้า 5 แห่ง
- ช่องจอด: 20,000
- กล้อง: 500 ตัว
- Sensor: 20,000 ตัว
- ต้องการ: real-time, accuracy 98%+, payment

### 5.2.2 Architecture

```
[Camera 500] → [Edge AI Jetson] → MQTT
[Sensor 20k] ──────────────────────┘
                                     ↓
                              [EMQX Cluster]
                                     ↓
                              [Kafka]
                                     ↓
                          [Flink Stream Processor]
                                     ↓
                    ┌────────────────┼────────────────┐
                    ↓                ↓                ↓
              [Redis State]   [PostgreSQL]     [InfluxDB]
              (slot state)    (history)        (occupancy)
                    ↓                ↓                ↓
                    └────────────────┼────────────────┘
                                     ↓
                              [API Gateway]
                                     ↓
    ┌────────────────┬───────────────┼───────────────┐
    ↓                ↓               ↓               ↓
[Mobile App]    [Kiosk]         [LED Sign]      [Dashboard]
    ↓                ↓               ↓               ↓
[WebSocket]     [REST]          [REST]          [REST]
```

### 5.2.3 ML Pipeline
```
[Camera Frame] → [YOLOv8 Detection] → [OCR Plate]
       ↓
[Slot Classification] → [Confidence Check]
       ↓
[Publish to MQTT] → [Aggregate]
       ↓
[Redis Update] → [WebSocket Broadcast]
```

### 5.2.4 ผลลัพธ์
| Metric | Value |
|---|---|
| Detection Accuracy | 98.5% |
| OCR Accuracy | 96.2% |
| Real-time Latency | < 2s |
| Uptime | 99.95% |
| เวลาหาที่จอด | ลด 45% |
| รายได้ที่จอด | เพิ่ม 22% |

### 5.2.5 บทเรียน
- Edge inference ลด bandwidth 90%
- Fallback rule เมื่อ ML ต่ำ
- Calibration กล้องทุก 3 เดือน
- Handle occlusion (รถบัง)
- Graceful degradation เมื่อ edge ล่ม

---

## 5.3 Case Study 3: Industrial IoT — Predictive Maintenance

### 5.3.1 บริบท
- โรงงานผลิตชิ้นส่วนยานยนต์
- เครื่องจักร: 200 ตัว
- Sensors: vibration, temperature, current
- ต้องการ: พยากรณ์ panne ล่วงหน้า

### 5.3.2 Architecture

```
[Vibration Sensor] → [Edge Gateway] → OPC-UA
[Temperature]  ────────────────────────┘
[Current]      ────────────────────────┘
                          ↓
                   [Edge Processing]
                   (FFT, Feature Extract)
                          ↓
                   [MQTT] → [Cloud]
                          ↓
                   [InfluxDB]
                          ↓
                   [ML Model (Python)]
                   (LSTM, Isolation Forest)
                          ↓
                   [Alert if Anomaly]
                          ↓
                   [CMMS Integration]
```

### 5.3.3 ML Model
- **Features**: RMS, kurtosis, FFT peaks, temperature trend
- **Model**: LSTM + Autoencoder
- **Training**: 6 เดือน historical
- **Inference**: ทุก 5 นาที
- **Threshold**: Dynamic (ไม่ใช่ fixed)

### 5.3.4 ผลลัพธ์
- ลด unplanned downtime 65%
- ประหยัดค่าซ่อม 40%
- เพิ่ม OEE 12%
- ROI ภายใน 8 เดือน

---

# 6. โครงสร้าง (Structure)

## 6.1 โครงสร้างคืออะไร

**โครงสร้าง (Structure)** ของระบบ IoT Solution หมายถึงการจัดองค์ประกอบของระบบออกเป็น layer, module, service ที่มีความสัมพันธ์กันอย่างชัดเจน โดยกำหนด:
- **ขอบเขต** ของแต่ละส่วน
- **ความรับผิดชอบ** ของแต่ละส่วน
- **การสื่อสาร** ระหว่างส่วน
- **การพึ่งพา** ระหว่างส่วน

โครงสร้างที่ดีทำให้ระบบ:
- เข้าใจง่าย
- แก้ไขง่าย
- Test ง่าย
- Scale ได้
- Evolve ได้

## 6.2 โครงสร้างทำงานอย่างไร

### 6.2.1 Layered Structure

```
┌────────────────────────────────────────────────────────┐
│ LAYER 7: PRESENTATION                                   │
│ Web, Mobile, Dashboard, AR, Voice                      │
├────────────────────────────────────────────────────────┤
│ LAYER 6: API GATEWAY                                    │
│ REST, GraphQL, gRPC, WebSocket, SSE                    │
├────────────────────────────────────────────────────────┤
│ LAYER 5: APPLICATION (Use Cases)                        │
│ Order Service, Fleet Service, Alert Service           │
├────────────────────────────────────────────────────────┤
│ LAYER 4: DOMAIN (Business Logic)                        │
│ Entities, Value Objects, Domain Services              │
├────────────────────────────────────────────────────────┤
│ LAYER 3: INFRASTRUCTURE                                 │
│ Repository Impl, External Clients, Message Bus         │
├────────────────────────────────────────────────────────┤
│ LAYER 2: DATA                                           │
│ PostgreSQL, InfluxDB, MongoDB, Redis, S3               │
├────────────────────────────────────────────────────────┤
│ LAYER 1: INTEGRATION                                    │
│ MQTT, Kafka, Webhook, gRPC                             │
├────────────────────────────────────────────────────────┤
│ LAYER 0: EDGE                                           │
│ IoT Gateway, Edge AI, Protocol Adapter                 │
└────────────────────────────────────────────────────────┘
```

### 6.2.2 Service Structure (Microservices)

```
┌─────────────────────────────────────────────────────────┐
│                    API GATEWAY                           │
│              (Kong / Traefik / Envoy)                    │
└────┬──────┬──────┬──────┬──────┬──────┬──────┬─────────┘
     ↓      ↓      ↓      ↓      ↓      ↓      ↓
┌────────┐┌──────┐┌──────┐┌──────┐┌──────┐┌──────┐┌──────┐
│Identity││Device││Order ││Fleet ││Alert ││Report││Notify│
│Service ││Svc   ││Svc   ││Svc   ││Svc   ││Svc   ││Svc   │
└───┬────┘└──┬───┘└──┬───┘└──┬───┘└──┬───┘└──┬───┘└──┬───┘
    ↓        ↓       ↓       ↓       ↓       ↓       ↓
┌────────┐┌──────┐┌──────┐┌──────┐┌──────┐┌──────┐┌──────┐
│ PG     ││ PG   ││ PG   ││ PG   ││ PG   ││ PG   ││Redis │
│(users) ││(dev) ││(ord) ││(fleet)││(alert)││(read) ││(sess)│
└────────┘└──────┘└──────┘└──────┘└──────┘└──────┘└──────┘
    ↓        ↓       ↓       ↓       ↓       ↓       ↓
    └────────┴───────┴───────┴───────┴───────┴───────┘
                            ↓
                    [Kafka Event Bus]
```

### 6.2.3 Clean Architecture Structure (ต่อ Service)

```
┌─────────────────────────────────────────────────────┐
│ Frameworks & Drivers                                 │
│ ┌─────────────────────────────────────────────────┐ │
│ │ Interface Adapters                               │ │
│ │ ┌─────────────────────────────────────────────┐ │ │
│ │ │ Application (Use Cases)                      │ │ │
│ │ │ ┌─────────────────────────────────────────┐ │ │ │
│ │ │ │ Domain (Entities)                        │ │ │ │
│ │ │ │                                          │ │ │ │
│ │ │ └─────────────────────────────────────────┘ │ │ │
│ │ └─────────────────────────────────────────────┘ │ │
│ └─────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────┘
         ↑ Dependency Rule: ชี้เข้าด้านในเสมอ ↑
```

## 6.3 โครงสร้างการใช้งาน

### 6.3.1 Directory Structure (Go)

```
iot-logistics/
├── cmd/
│   ├── api/main.go
│   ├── worker/main.go
│   └── migrator/main.go
├── internal/
│   ├── domain/
│   │   ├── device/
│   │   │   ├── entity.go
│   │   │   ├── value_object.go
│   │   │   ├── repository.go
│   │   │   └── service.go
│   │   ├── order/
│   │   ├── fleet/
│   │   └── alert/
│   ├── application/
│   │   ├── device/
│   │   │   ├── register.go
│   │   │   ├── command.go
│   │   │   └── query.go
│   │   ├── order/
│   │   └── fleet/
│   ├── infrastructure/
│   │   ├── persistence/
│   │   │   ├── postgres/
│   │   │   ├── influx/
│   │   │   ├── mongo/
│   │   │   └── redis/
│   │   ├── messaging/
│   │   │   ├── kafka/
│   │   │   ├── mqtt/
│   │   │   └── nats/
│   │   └── external/
│   │       ├── payment/
│   │       └── map/
│   ├── interface/
│   │   ├── http/
│   │   │   ├── handler/
│   │   │   ├── middleware/
│   │   │   └── router.go
│   │   ├── grpc/
│   │   └── websocket/
│   └── shared/
│       ├── config/
│       ├── logger/
│       ├── errors/
│       └── telemetry/
├── pkg/
│   ├── algorithm/
│   │   ├── suffix_automaton.go
│   │   ├── hld.go
│   │   └── wavelet.go
│   └── utils/
├── api/
│   ├── proto/
│   └── openapi/
├── migrations/
├── deploy/
│   ├── docker/
│   ├── k8s/
│   └── terraform/
├── test/
│   ├── unit/
│   ├── integration/
│   └── e2e/
├── docs/
│   ├── adr/
│   └── architecture/
├── go.mod
└── Makefile
```

### 6.3.2 Data Structure ในแต่ละ Layer

| Layer | Data Structure | ตัวอย่าง |
|---|---|---|
| Domain | Entity, VO, Aggregate | `Device`, `DeviceID` |
| Application | DTO, Command, Query | `RegisterDeviceCmd` |
| Interface | Request, Response | `RegisterDeviceRequest` |
| Infrastructure | Row, Document | `deviceRow` |
| Shared | Error, Result | `DomainError` |

## 6.4 สรุป

โครงสร้างที่ดีต้อง:
- **ชัดเจน**: ใครทำอะไร
- **สอดคล้อง**: ใช้ pattern เดียวกัน
- **ยืดหยุ่น**: เปลี่ยนได้โดยไม่กระทบมาก
- **ทดสอบได้**: test แยก layer
- **Documented**: มี ADR + diagram

---

# 7. แนวทางการประยุกต์ใช้ (Application Guidelines)

## 7.1 การเลือก Architecture Pattern

| สถานการณ์ | Pattern ที่แนะนำ |
|---|---|
| Startup, MVp | Modular Monolith |
| Team เล็ก (2-5) | Modular Monolith |
| Team กลาง (5-15) | Microservices (บางส่วน) |
| Team ใหญ่ (15+) | Microservices + DDD |
| Real-time | Event-Driven + CQRS |
| Legacy integration | Hexagonal + Adapter |
| High read | CQRS + Read Replica |
| Audit required | Event Sourcing |
| Burst workload | Serverless |

## 7.2 การเลือก Database

| ข้อมูล | Database | เหตุผล |
|---|---|---|
| Transaction | PostgreSQL | ACID |
| Time-series | InfluxDB | Optimize สำหรับ time |
| Document | MongoDB | Schema flexible |
| Cache | Redis | In-memory |
| Search | Elasticsearch | Full-text |
| Vector | Qdrant/Milvus | Similarity |
| Graph | Neo4j | Relationship |
| Wide-column | Cassandra | High write |

## 7.3 การเลือก Protocol

```
Device → Cloud:      MQTT (lightweight, QoS)
Service → Service:    gRPC (fast, typed)
App → API:           REST (simple) / GraphQL (flexible)
Real-time:           WebSocket (bidirectional) / SSE (one-way)
Event:               Webhook / Kafka
```

## 7.4 การออกแบบ API

**REST Best Practices:**
- ใช้ noun ไม่ใช่ verb: `/devices` ไม่ใช่ `/getDevices`
- Version: `/v1/devices`
- Pagination: `?page=1&limit=20`
- Filtering: `?status=active&type=sensor`
- Sorting: `?sort=-created_at`
- Error format: มาตรฐาน
- Idempotency key สำหรับ POST
- Rate limit header: `X-RateLimit-Remaining`

**gRPC Best Practices:**
- ใช้ `proto3`
- Version package: `iot.v1`
- ใช้ streaming เมื่อเหมาะสม
- Error ด้วย `google.rpc.Status`
- Deadline ทุก call
- Metadata สำหรับ auth

**GraphQL Best Practices:**
- DataLoader แก้ N+1
- Query complexity limit
- Depth limit
- Persisted queries
- Cache directive

## 7.5 การออกแบบ Event

**Event Design:**
```json
{
  "id": "evt_01H...",
  "type": "device.telemetry.received",
  "source": "iot/device/d-001",
  "time": "2026-01-15T10:00:00Z",
  "tenant_id": "t-001",
  "data": {
    "device_id": "d-001",
    "metric": "temperature",
    "value": 25.3,
    "unit": "celsius"
  },
  "metadata": {
    "version": "1.0",
    "correlation_id": "req_abc"
  }
}
```

**Event Best Practices:**
- ใช้ past tense: `OrderCreated` ไม่ใช่ `CreateOrder`
- Immutable
- Include metadata
- Version schema
- Idempotency
- Ordering key

## 7.6 การออกแบบ Security

**Layers:**
1. **Network**: VPC, Security Group, WAF
2. **Transport**: TLS 1.3, mTLS
3. **Application**: OAuth2, JWT, API Key
4. **Data**: Encryption at rest, RLS
5. **Device**: Certificate, secure boot
6. **Audit**: ทุก action

**Zero Trust:**
- Never trust, always verify
- Least privilege
- Assume breach
- Micro-segmentation

## 7.7 การออกแบบ Observability

**ทุก Service ต้องมี:**
- `/health` (liveness)
- `/ready` (readiness)
- `/metrics` (Prometheus)
- Structured logs (JSON)
- Trace propagation
- Correlation ID

**SLI/SLO ตัวอย่าง:**
| SLI | SLO | Error Budget |
|---|---|---|
| Availability | 99.95% | 21.6 min/เดือน |
| Latency P95 | < 200ms | 5% request |
| Error rate | < 0.1% | 0.1% |
| Throughput | > 10k rps | - |

## 7.8 การออกแบบ Cost

- **Compute**: Spot instance สำหรับ batch, Reserved สำหรับ baseline
- **Storage**: S3 Intelligent-Tiering, Lifecycle policy
- **Network**: CDN, compression, batch
- **Database**: Read replica, connection pool
- **Monitoring**: Sampling, retention policy
- **ML**: Batch inference, model quantization

---

# 8. กรณีตัวอย่างจากประสบการณ์

## 8.1 ประสบการณ์ที่ 1: MQTT Broker ล่มช่วง Peak

**สถานการณ์:**
- EMQX broker node 1 ใน 3 ล่ม
- เกิดตอน peak 60k msg/s
- ลูกค้าบ่นว่า tracking หยุด

**สาเหตุ:**
- Node รับ connection เกิน capacity
- ไม่มี connection limit
- ไม่มี backpressure

**แนวทางแก้ไข:**
- ตั้ง `max_connections` ต่อ node
- ใช้ `max_message_size` limit
- ใช้ shared subscription กระจาย load
- เพิ่ม node เป็น 5
- ใช้ load balancer แบบ least connection
- Monitor connection count

**บทเรียน:**
- Capacity planning สำคัญ
- ต้องมี backpressure
- Load test ก่อน peak season

---

## 8.2 ประสบการณ์ที่ 2: Cache Stampede ตอน Deploy

**สถานการณ์:**
- Deploy service ใหม่ → cache หาย
- Database รับ load พุ่ง 50x
- API timeout

**สาเหตุ:**
- Cache invalidation พร้อมกัน
- ไม่มี singleflight
- ไม่มี circuit breaker

**แนวทางแก้ไข:**
- ใช้ singleflight (golang.org/x/sync/singleflight)
- Jitter TTL
- Warm-up cache หลัง deploy
- Circuit breaker ที่ DB client
- Read replica สำหรับ fallback

**โค้ด:**
```go
var group singleflight.Group

func GetDevice(ctx context.Context, id string) (*Device, error) {
    v, err, _ := group.Do("device:"+id, func() (any, error) {
        // check cache
        if d, err := cache.Get(ctx, id); err == nil {
            return d, nil
        }
        // query DB (single request)
        d, err := repo.FindByID(ctx, id)
        if err != nil { return nil, err }
        cache.Set(ctx, id, d, jitterTTL(5*time.Minute))
        return d, nil
    })
    if err != nil { return nil, err }
    return v.(*Device), nil
}
```

**บทเรียน:**
- Singleflight แก้ stampede
- Jitter ป้องกัน synchronized expiry
- Circuit breaker ป้องกัน cascade

---

## 8.3 ประสบการณ์ที่ 3: Clock Skew ทำให้ Data ผิด

**สถานการณ์:**
- Device บางตัว clock เร็ว/ช้า 5-10 นาที
- Time-series data เรียงผิด
- Alert trigger ผิดเวลา

**สาเหตุ:**
- Device ไม่มี NTP
- ใช้ local clock
- ไม่มี validation

**แนวทางแก้ไข:**
- บังคับ NTP sync ตอน boot
- Server-side timestamp เป็นหลัก
- เก็บ `device_time` และ `server_time`
- Reject ถ้า skew > 1 นาที
- Alert ถ้า skew เกิน

**บทเรียน:**
- อย่า trust device clock
- Server-side timestamp
- Monitor clock skew

---

## 8.4 ประสบการณ์ที่ 4: Kafka Hot Partition

**สถานการณ์:**
- Partition หนึ่งรับ load 80% ของทั้งหมด
- Consumer lag สูง
- Latency เพิ่ม

**สาเหตุ:**
- Key = tenant_id
- Tenant ใหญ่ 1 รายส่ง 80% ของ msg
- Partition ไม่ balance

**แนวทางแก้ไข:**
- Key = tenant_id + device_id
- Custom Partitioner
- เพิ่ม partition
- Separate topic สำหรับ tenant ใหญ่
- Monitor partition skew

**บทเรียน:**
- อย่าใช้ key ที่ cardinality ต่ำ
- Monitor partition distribution
- Plan สำหรับ tenant ใหญ่

---

## 8.5 ประสบการณ์ที่ 5: WebSocket Connection Leak

**สถานการณ์:**
- หลังรัน 3 วัน memory เพิ่ม 5GB
- Server crash
- ต้อง restart

**สาเหตุ:**
- Goroutine leak
- ไม่มี read/write deadline
- ไม่มี ping/pong
- ไม่ cleanup connection ที่ตาย

**แนวทางแก้ไข:**
```go
const (
    writeWait  = 10 * time.Second
    pongWait   = 60 * time.Second
    pingPeriod = (pongWait * 9) / 10
)

func (c *Client) readPump() {
    defer func() {
        c.hub.unregister <- c
        c.conn.Close()
    }()
    c.conn.SetReadDeadline(time.Now().Add(pongWait))
    c.conn.SetPongHandler(func(string) error {
        c.conn.SetReadDeadline(time.Now().Add(pongWait))
        return nil
    })
    for {
        _, msg, err := c.conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err) {
                log.Error("unexpected close", zap.Error(err))
            }
            break
        }
        // handle msg
    }
}

func (c *Client) writePump() {
    ticker := time.NewTicker(pingPeriod)
    defer func() {
        ticker.Stop()
        c.conn.Close()
    }()
    for {
        select {
        case msg, ok := <-c.send:
            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if !ok { return }
            c.conn.WriteMessage(websocket.TextMessage, msg)
        case <-ticker.C:
            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
    }
}
```

**บทเรียน:**
- ตั้ง deadline ทุก read/write
- Ping/pong สำหรับ keepalive
- ใช้ pprof หา leak
- Graceful shutdown

---

# 9. การนำไปใช้งานจริง (Production Implementation)

## 9.1 Infrastructure

```yaml
# Kubernetes deployment
apiVersion: apps/v1
kind: Deployment
metadata:
  name: device-service
spec:
  replicas: 5
  selector:
    matchLabels: { app: device-service }
  template:
    metadata:
      labels: { app: device-service }
    spec:
      containers:
      - name: app
        image: registry/device-service:v1.2.3
        ports:
        - containerPort: 8080
        - containerPort: 9090  # metrics
        env:
        - name: DB_URL
          valueFrom: { secretKeyRef: { name: db, key: url } }
        resources:
          requests: { cpu: 500m, memory: 512Mi }
          limits:   { cpu: 2000m, memory: 2Gi }
        livenessProbe:
          httpGet: { path: /health, port: 8080 }
          initialDelaySeconds: 10
        readinessProbe:
          httpGet: { path: /ready, port: 8080 }
          initialDelaySeconds: 5
```

## 9.2 CI/CD Pipeline

```yaml
# .github/workflows/deploy.yml
name: Deploy
on:
  push:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    - uses: actions/setup-go@v5
      with: { go-version: '1.22' }
    - run: go test ./... -race -cover
    - run: go vet ./...
    - uses: golangci/golangci-lint-action@v4

  build:
    needs: test
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    - uses: docker/build-push-action@v5
      with:
        push: true
        tags: registry/device-service:${{ github.sha }}

  deploy:
    needs: build
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    - run: |
        kubectl set image deployment/device-service \
          app=registry/device-service:${{ github.sha }}
        kubectl rollout status deployment/device-service
```

## 9.3 Migration Strategy

```go
// migrations/20260115_add_device_groups.sql
-- +migrate Up
CREATE TABLE device_groups (...);
CREATE INDEX ...;

-- +migrate Down
DROP TABLE device_groups;
```

**Zero-downtime migration:**
1. **Expand**: เพิ่ม column/table ใหม่
2. **Migrate**: copy data + dual write
3. **Contract**: ลบ column/table เก่า

## 9.4 Deployment Strategy

| Strategy | ข้อดี | ข้อเสีย |
|---|---|---|
| Recreate | ง่าย | Downtime |
| Rolling | No downtime | ช้า, 2 version พร้อมกัน |
| Blue-Green | เร็ว, rollback ง่าย | ใช้ resource 2x |
| Canary | ความเสี่ยงต่ำ | ซับซ้อน |
| A/B | Test feature | ซับซ้อน |

## 9.5 Rollback Plan

```bash
# Quick rollback
kubectl rollout undo deployment/device-service

# Database rollback
migrate -path ./migrations -database $DB_URL down 1

# Feature flag disable
curl -X POST $FF_URL/toggle -d '{"flag":"new_feature","enabled":false}'
```

## 9.6 Runbook ตัวอย่าง

```markdown
# Runbook: High Error Rate

## Detection
- Alert: `HighErrorRate` > 5%
- Dashboard: Error rate panel

## Diagnosis
1. ดู error log: `kubectl logs -l app=device-service --tail=100`
2. ดู trace: Jaeger → filter error
3. ดู metric: Grafana → error by endpoint
4. ดู dependency: status page ของ DB, Kafka

## Mitigation
1. ถ้า code bug → rollback: `kubectl rollout undo`
2. ถ้า DB overload → scale read replica
3. ถ้า dependency down → enable circuit breaker
4. ถ้า traffic spike → scale up

## Post-incident
- เขียน post-mortem ภายใน 48h
- Update runbook
- Add test
```

---

# 10. ปัญหาและแนวทางแก้ไข

## 10.1 ปัญหาด้าน Architecture

| ปัญหา | อาการ | แนวทางแก้ไข |
|---|---|---|
| Monolith แข็งตัว | Deploy ช้า, test นาน | แยก module → service |
| Distributed monolith | Service เยอะ แต่ deploy พร้อมกัน | แยก data + async |
| Over-engineering | ระบบซับซ้อนเกิน | Simplify, เริ่มจาก monolith |
| Tight coupling | เปลี่ยน 1 ที่ พังหลายที่ | Interface + DI |
| Shared database | Service ผูกกัน | Database per service |
| No clear boundary | งงว่าอะไรอยู่ไหน | DDD Bounded Context |

## 10.2 ปัญหาด้าน Performance

| ปัญหา | อาการ | แนวทางแก้ไข |
|---|---|---|
| N+1 query | DB load สูง | Eager loading, DataLoader |
| Slow query | API ตอบช้า | Index, EXPLAIN, cache |
| Memory leak | RAM เพิ่มเรื่อย ๆ | pprof, fix goroutine leak |
| CPU spike | Response ช้า | Profile, optimize algorithm |
| Network latency | ตอบช้า | Cache, CDN, edge |
| Lock contention | Deadlock | ลด transaction scope |

## 10.3 ปัญหาด้าน Scalability

| ปัญหา | อาการ | แนวทางแก้ไข |
|---|---|---|
| Vertical limit | Scale up ไม่ได้ | Horizontal + sharding |
| Hot partition | Node หนึ่ง overload | Better key, salt |
| Stateful service | Scale ยาก | Stateless + external state |
| DB bottleneck | Write ช้า | Sharding, CQRS, async |
| Connection limit | Connection refused | Pool, PgBouncer |

## 10.4 ปัญหาด้าน Reliability

| ปัญหา | อาการ | แนวทางแก้ไข |
|---|---|---|
| Single point of failure | ล่มทั้งระบบ | Redundancy |
| Cascade failure | ล่มเป็นลูกโซ่ | Circuit breaker |
| Thundering herd | Load พุ่ง | Jitter, singleflight |
| Retry storm | Load ทวีคูณ | Exponential backoff + jitter |
| Partial failure | บางส่วนพัง | Timeout, bulkhead |
| Data loss | ข้อมูลหาย | Replication, backup |

## 10.5 ปัญหาด้าน Security

| ปัญหา | อาการ | แนวทางแก้ไข |
|---|---|---|
| Weak auth | ถูก brute force | MFA, rate limit |
| Secret leak | Secret ใน code | Vault, scanning |
| Injection | SQL/NoSQL injection | Parameterized query |
| XSS | Script injection | CSP, sanitize |
| MITM | ดักข้อมูล | TLS, mTLS |
| DDoS | Server ล่ม | WAF, rate limit |
| Insider threat | พนักงานทำร้าย | Least privilege, audit |

## 10.6 ปัญหาด้าน Data

| ปัญหา | อาการ | แนวทางแก้ไข |
|---|---|---|
| Data inconsistency | ข้อมูลไม่ตรง | Transaction, saga |
| Schema drift | Schema ไม่ตรง | Migration, contract |
| Data loss | ข้อมูลหาย | Backup, PITR |
| Duplicate | ข้อมูลซ้ำ | Idempotency, dedup |
| Clock skew | เวลาผิด | NTP, server timestamp |
| GDPR/PDPA | ผิดกฎหมาย | Data residency, consent |

## 10.7 ปัญหาด้าน Operations

| ปัญหา | อาการ | แนวทางแก้ไข |
|---|---|---|
| Manual deploy | ผิดพลาดบ่อย | CI/CD |
| No rollback | แก้ไม่ได้ | Blue-green, feature flag |
| Config drift | Env ไม่ตรง | IaC, GitOps |
| No monitoring | ไม่รู้สถานะ | Observability |
| Alert fatigue | ไม่สนใจ alert | Actionable alert, SLO |
| Knowledge silo | คนเดียวรู้ | Doc, runbook, rotation |

## 10.8 ปัญหาด้าน Cost

| ปัญหา | อาการ | แนวทางแก้ไข |
|---|---|---|
| Over-provisioning | จ่ายเกิน | Right-sizing, auto-scale |
| Data transfer | ค่า network แพง | CDN, compression |
| Storage growth | ค่า storage พุ่ง | Lifecycle, tiering |
| Idle resource | จ่ายเปล่า | Spot, schedule |
| Log/metric volume | ค่า observability แพง | Sampling, retention |
| Vendor lock-in | ย้ายยาก | Abstraction, multi-cloud |

## 10.9 ปัญหาด้าน Team

| ปัญหา | อาการ | แนวทางแก้ไข |
|---|---|---|
| Conway's law | Architecture = org chart | Team topology |
| Knowledge gap | ทำไม่ได้ | Training, pair |
| Communication | งงกัน | Doc, ADR, RFC |
| Burnout | คนลาออก | Sustainable pace |
| On-call fatigue | เครียด | Rotation, runbook |

## 10.10 Best Practices สรุป

### 10.10.1 Architecture
- เริ่มจาก simplicity
- Domain-first, tech-second
- Document decisions (ADR)
- Evolve ไม่ใช่ big bang
- Measure before optimize

### 10.10.2 Code
- Clean code
- Test pyramid
- Code review
- Static analysis
- Continuous refactor

### 10.10.3 Operations
- Automate everything
- Infrastructure as Code
- Immutable infrastructure
- GitOps
- Observability-first

### 10.10.4 Security
- Zero trust
- Least privilege
- Defense in depth
- Encrypt everywhere
- Audit everything

### 10.10.5 Data
- Schema as contract
- Migration discipline
- Backup + test restore
- Data quality
- Privacy by design

### 10.10.6 Team
- Blameless post-mortem
- Continuous learning
- Documentation
- On-call rotation
- Sustainable pace

---

# บทสรุปสุดท้าย

เอกสารนี้ครอบคลุม **10 บท** ตามโครงสร้างที่กำหนด:

1. **บทนำ** — ความเป็นมา วัตถุประสงค์ ขอบเขต
2. **บทนิยาม** — คำศัพท์ 7 หมวด
3. **บทหัวข้อ** — Software Architecture, IoT, Protocol, Database, Monitoring (แต่ละหัวข้อมี 3.1-3.8 ครบ)
4. **ออกแบบ Workflow** — 5 phases + swimlane
5. **Case Study** — 3 cases (Logistics, Parking, Industrial)
6. **โครงสร้าง** — Layered, Microservices, Clean Architecture, Directory
7. **แนวทางการประยุกต์ใช้** — Pattern, DB, Protocol, API, Security, Observability, Cost
8. **กรณีตัวอย่างจากประสบการณ์** — 5 ประสบการณ์จริง
9. **การนำไปใช้งานจริง** — K8s, CI/CD, Migration, Deployment, Runbook
10. **ปัญหาและแนวทางแก้ไข** — 10 หมวด + Best Practices

**หลักคิดสำคัญ:**
- ไม่มี silver bullet — มีแต่ trade-off
- Architecture ต้อง evolve ตาม business
- Observability ไม่ใช่ optional
- Security ต้อง built-in ไม่ใช่ bolt-on
- Team คือ factor สำคัญที่สุด

**ขั้นตอนต่อไป:**
1. เลือก scope ที่เหมาะกับธุรกิจ
2. สร้าง PoC สำหรับ critical path
3. วัดผลก่อน scale
4. Iterate ตาม feedback
5. Document ทุก decision

---
 