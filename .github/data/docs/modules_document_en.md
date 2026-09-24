# คู่มือการพัฒนา, Deploy, บำรุงรักษา และ Monitoring – ภาษาไทย

## 📚 สารบัญ

1. [บทนำ](#1-บทนำ)
2. [การพัฒนา (Development)](#2-การพัฒนา-development)
3. [การ Deploy](#3-การ-deploy)
4. [การบำรุงรักษา (Maintenance)](#4-การบำรุงรักษา-maintenance)
5. [การขยายระบบ (System Expansion)](#5-การขยายระบบ-system-expansion)
6. [การ Monitoring (การตรวจสอบ)](#6-การ-monitoring-การตรวจสอบ)
7. [การแก้ไขปัญหาเบื้องต้น (Troubleshooting)](#7-การแก้ไขปัญหาเบื้องต้น-troubleshooting)

---

## 1. บทนำ

### 1.1 วัตถุประสงค์
คู่มือนี้จัดทำขึ้นเพื่อเป็นแนวทางในการพัฒนา, ติดตั้ง, บำรุงรักษา, ขยายระบบ และตรวจสอบการทำงานของระบบบริหารจัดการอู่ซ่อมรถ Auto Repair Shop Management System

### 1.2 ผู้ใช้งานคู่มือ
- นักพัฒนาระบบ (Developers)
- วิศวกร DevOps
- ผู้ดูแลระบบ (System Administrators)
- ทีมงานบำรุงรักษา

### 1.3 สถาปัตยกรรมโดยรวม

```
┌─────────────────────────────────────────────────────────────────────┐
│                     Frontend (React / Angular / Mobile)             │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                     API Gateway (Spring Cloud Gateway / NGINX)      │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                     Spring Boot Application (REST API)              │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐│
│  │  Auth    │ │  Job     │ │  Quot    │ │  Invoice │ │  Payment ││
│  │  Module  │ │  Module  │ │  Module  │ │  Module  │ │  Module  ││
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘ └──────────┘│
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐│
│  │Inventory │ │  PO      │ │Customer  │ │Document  │ │   IoT    ││
│  │  Module  │ │  Module  │ │  Module  │ │  Module  │ │  Module  ││
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘ └──────────┘│
└─────────────────────────────────────────────────────────────────────┘
                                    │
            ┌───────────────────────┼───────────────────────┐
            ▼                       ▼                       ▼
┌───────────────────┐   ┌───────────────────┐   ┌───────────────────┐
│    PostgreSQL     │   │      Redis        │   │      Kafka        │
│   (Main Database) │   │     (Cache)       │   │   (Event Bus)     │
└───────────────────┘   └───────────────────┘   └───────────────────┘
            │                       │                       │
            ▼                       ▼                       ▼
┌───────────────────┐   ┌───────────────────┐   ┌───────────────────┐
│    MongoDB        │   │     InfluxDB      │   │   Elasticsearch   │
│   (Logging)       │   │    (IoT Data)     │   │   (Search/Logs)   │
└───────────────────┘   └───────────────────┘   └───────────────────┘
```

### 1.4 แผนภาพการ Deploy

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         AWS CLOUD (EC2 / ECS)                          │
│                                                                         │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                     Load Balancer (ALB)                          │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                    │                                     │
│              ┌─────────────────────┼─────────────────────┐              │
│              ▼                     ▼                     ▼              │
│  ┌───────────────────┐ ┌───────────────────┐ ┌───────────────────┐    │
│  │   Instance 1      │ │   Instance 2      │ │   Instance 3      │    │
│  │  (App Server)     │ │  (App Server)     │ │  (App Server)     │    │
│  └───────────────────┘ └───────────────────┘ └───────────────────┘    │
│              │                     │                     │              │
│              └─────────────────────┼─────────────────────┘              │
│                                    ▼                                     │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                     RDS (PostgreSQL)                            │   │
│  │                     ElastiCache (Redis)                         │   │
│  │                     MSK (Kafka)                                 │   │
│  │                     S3 (Document Storage)                       │   │
│  └─────────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 2. การพัฒนา (Development)

### 2.1 ข้อกำหนดเบื้องต้น (Prerequisites)

| รายการ | เวอร์ชัน |
|--------|----------|
| Java | 21+ |
| Maven | 3.8+ |
| Docker | 20.10+ |
| Docker Compose | 2.0+ |
| Git | 2.30+ |
| IntelliJ IDEA (แนะนำ) | 2023.2+ |

### 2.2 ขั้นตอนการติดตั้งสำหรับ Development

#### 2.2.1 Clone Project
```bash
git clone https://github.com/your-org/auto-repair-system.git
cd auto-repair-system
```

#### 2.2.2 ตั้งค่า Environment Variables
```bash
cp .env.example .env
# แก้ไขไฟล์ .env ให้ตรงกับสภาพแวดล้อม
```

#### 2.2.3 เริ่มต้น Database ด้วย Docker Compose
```bash
docker-compose up -d
```

บริการที่ถูก:
- PostgreSQL (port 5432)
- Redis (port 6379)
- Kafka + Zookeeper (port 9092)
- MongoDB (port 27017)
- InfluxDB (port 8086)
- Elasticsearch (port 9200)
- Kibana (port 5601)
- Grafana (port 3000)

#### 2.2.4 รัน Application
```bash
# Development Profile
mvn spring-boot:run -Dspring-boot.run.profiles=dev

# หรือรันด้วย IntelliJ IDEA โดยเลือก Profile "dev"
```

#### 2.2.5 เข้าถึง Swagger UI
```
http://localhost:5000/api/swagger-ui.html
```

### 2.3 การทดสอบ (Testing)

#### 2.3.1 รัน Unit Tests
```bash
mvn test
```

#### 2.3.2 รัน Tests เฉพาะโมดูล
```bash
mvn test -Dtest="com.icmon.module.job.**"
```

#### 2.3.3 ดู Coverage Report
```bash
mvn jacoco:report
# เปิดไฟล์ target/site/jacoco/index.html
```

#### 2.3.4 Test Coverage Target
| Layer | Target |
|-------|--------|
| Domain | ≥ 95% |
| Application | ≥ 90% |
| Infrastructure | ≥ 70% |
| Presentation | ≥ 80% |
| **Overall** | **≥ 80%** |

### 2.4 Coding Standards

#### 2.4.1 การตั้งชื่อ (Naming Conventions)
| ประเภท | รูปแบบ | ตัวอย่าง |
|--------|--------|----------|
| Package | lowercase | `com.icmon.module.job` |
| Class | PascalCase | `JobServiceImpl` |
| Interface | PascalCase | `JobService` |
| Method | camelCase | `createJob()` |
| Variable | camelCase | `jobId` |
| Constant | UPPER_SNAKE_CASE | `MAX_RETRY_COUNT` |

#### 2.4.2 โครงสร้าง Package
```
com.icmon/
├── configuration/          # Spring Configurations
├── module/
│   └── {module_name}/
│       ├── application/    # Service Layer
│       │   ├── interfaces/
│       │   └── impl/
│       ├── domain/         # Domain Layer
│       ├── infrastructure/ # Infrastructure Layer
│       │   ├── repository/
│       │   ├── entity/
│       │   └── mapper/
│       └── presentation/   # Presentation Layer
│           ├── controller/
│           └── dto/
├── exception/              # Global Exception
├── logging/                # Logging System
└── utils/                  # Utilities
```

### 2.5 Git Workflow

```
main (Production)
  │
  ├── develop (Integration)
  │   ├── feature/job-management    ← พัฒนา Feature
  │   ├── feature/quotation-module
  │   └── feature/inventory
  │
  ├── release/v2.0.0
  │
  └── hotfix/critical-bug
```

#### คำสั่งที่ใช้บ่อย
```bash
# สร้าง Branch ใหม่สำหรับพัฒนา Feature
git checkout -b feature/your-feature-name develop

# Commit และ Push
git add .
git commit -m "feat: add new feature description"
git push origin feature/your-feature-name

# สร้าง Pull Request บน GitHub/GitLab
# เมื่อ Merge แล้วให้ลบ Branch
```

#### Commit Message Format
```
<type>(<scope>): <subject>

[optional body]

[optional footer]
```

**Type:** `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`

**ตัวอย่าง:**
```
feat(job): add job status history tracking

- Add TJobStatusHistory entity
- Create JobStatusHistoryRepository
- Implement status change logging

Closes #123
```

---

## 3. การ Deploy

### 3.1 Environment Profiles

| Profile | ไฟล์ | การใช้งาน |
|---------|------|----------|
| **local** | `application-local.yml` | พัฒนาท้องถิ่น (Local Development) |
| **dev** | `application-dev.yml` | Environment พัฒนา (Development) |
| **test** | `application-test.yml` | Environment ทดสอบ (Test/UAT) |
| **prod** | `application-prod.yml` | Environment ผลิตจริง (Production) |

### 3.2 ระบบ Database

#### 3.2.1 Development (Local)
```yaml
spring:
  datasource:
    url: jdbc:postgresql://localhost:5432/autorepair
    username: admin
    password: admin123
  jpa:
    hibernate:
      ddl-auto: update
    show-sql: true
```

#### 3.2.2 UAT/Test
```yaml
spring:
  datasource:
    url: ${DB_URL}
    username: ${DB_USERNAME}
    password: ${DB_PASSWORD}
  jpa:
    hibernate:
      ddl-auto: validate
    show-sql: false
```

#### 3.2.3 Production
```yaml
spring:
  datasource:
    url: ${DB_URL}
    username: ${DB_USERNAME}
    password: ${DB_PASSWORD}
  jpa:
    hibernate:
      ddl-auto: validate
    show-sql: false
    properties:
      hibernate:
        jdbc:
          batch_size: 50
        order_inserts: true
        order_updates: true
```

### 3.3 การ Build

#### 3.3.1 Build สำหรับ Development
```bash
mvn clean compile
```

#### 3.3.2 Build สำหรับ Production
```bash
mvn clean package -Pprod
```

#### 3.3.3 Build Docker Image
```bash
# สร้าง Docker Image
docker build -t auto-repair-system:${VERSION} .

# Push ไปยัง Registry
docker tag auto-repair-system:${VERSION} ${REGISTRY}/auto-repair-system:${VERSION}
docker push ${REGISTRY}/auto-repair-system:${VERSION}
```

### 3.4 การ Deploy แต่ละ Environment

#### 3.4.1 Local Development
```bash
mvn spring-boot:run -Dspring-boot.run.profiles=local
```

#### 3.4.2 Development (Dev) Server
```bash
# สร้างไฟล์ .env.dev
# รันด้วย Docker Compose
docker-compose -f docker-compose.dev.yml up -d
```

#### 3.4.3 UAT/Test Server
```bash
# 1. Build JAR
mvn clean package -Ptest

# 2. Deploy JAR
java -jar target/app-1.0.0-SNAPSHOT.jar \
  --spring.profiles.active=test \
  --server.port=5000
```

#### 3.4.4 Production Server
```bash
# 1. Build JAR with Production Profile
mvn clean package -Pprod

# 2. Deploy with Production Config
java -jar target/app-1.0.0-SNAPSHOT.jar \
  --spring.profiles.active=prod \
  --server.port=8443 \
  --server.ssl.enabled=true \
  --server.ssl.key-store=/path/to/keystore.p12
```

### 3.5 CI/CD Pipeline (Jenkins)

```groovy
pipeline {
    agent any
    
    environment {
        REGISTRY = 'docker.io/your-registry'
        IMAGE_NAME = 'auto-repair-system'
    }
    
    stages {
        stage('Checkout') {
            steps {
                git branch: 'main', url: 'https://github.com/your-org/auto-repair-system.git'
            }
        }
        
        stage('Build') {
            steps {
                sh 'mvn clean compile'
            }
        }
        
        stage('Unit Tests') {
            steps {
                sh 'mvn test'
            }
        }
        
        stage('Integration Tests') {
            steps {
                sh 'mvn test -Dgroups="integration"'
            }
        }
        
        stage('Build JAR') {
            steps {
                sh 'mvn package -Pprod'
            }
        }
        
        stage('Build Docker Image') {
            steps {
                sh '''
                    docker build -t ${REGISTRY}/${IMAGE_NAME}:${BUILD_NUMBER} .
                    docker tag ${REGISTRY}/${IMAGE_NAME}:${BUILD_NUMBER} ${REGISTRY}/${IMAGE_NAME}:latest
                '''
            }
        }
        
        stage('Push to Registry') {
            steps {
                sh '''
                    docker push ${REGISTRY}/${IMAGE_NAME}:${BUILD_NUMBER}
                    docker push ${REGISTRY}/${IMAGE_NAME}:latest
                '''
            }
        }
        
        stage('Deploy to DEV') {
            when { branch 'develop' }
            steps {
                sh 'kubectl set image deployment/auto-repair-system-dev app=${REGISTRY}/${IMAGE_NAME}:${BUILD_NUMBER} -n dev'
            }
        }
        
        stage('Deploy to UAT') {
            when { branch 'release/*' }
            steps {
                sh 'kubectl set image deployment/auto-repair-system-uat app=${REGISTRY}/${IMAGE_NAME}:${BUILD_NUMBER} -n uat'
            }
        }
        
        stage('Deploy to Production') {
            when { branch 'main' }
            steps {
                input message: 'Deploy to Production?', ok: 'Yes'
                sh 'kubectl set image deployment/auto-repair-system-prod app=${REGISTRY}/${IMAGE_NAME}:${BUILD_NUMBER} -n prod'
            }
        }
    }
}
```

### 3.6 Docker Compose สำหรับ Production

```yaml
version: '3.8'

services:
  app:
    image: ${REGISTRY}/auto-repair-system:latest
    restart: always
    ports:
      - "8443:8443"
    environment:
      SPRING_PROFILES_ACTIVE: prod
      DB_URL: jdbc:postgresql://postgres:5432/autorepair
      DB_USERNAME: ${DB_USERNAME}
      DB_PASSWORD: ${DB_PASSWORD}
      REDIS_HOST: redis
      REDIS_PORT: 6379
      KAFKA_BOOTSTRAP_SERVERS: kafka:9092
    volumes:
      - ./logs:/app/logs
      - ./uploads:/app/uploads
    depends_on:
      - postgres
      - redis
      - kafka
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8443/actuator/health"]
      interval: 30s
      timeout: 10s
      retries: 5

  postgres:
    image: postgres:15
    restart: always
    environment:
      POSTGRES_DB: autorepair
      POSTGRES_USER: ${DB_USERNAME}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"

  redis:
    image: redis:7-alpine
    restart: always
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

  kafka:
    image: confluentinc/cp-kafka:latest
    restart: always
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:9092
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
    ports:
      - "9092:9092"
    depends_on:
      - zookeeper

  zookeeper:
    image: confluentinc/cp-zookeeper:latest
    restart: always
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181
      ZOOKEEPER_TICK_TIME: 2000
    ports:
      - "2181:2181"

volumes:
  postgres_data:
  redis_data:
```

---

## 4. การบำรุงรักษา (Maintenance)

### 4.1 การ Backup Database

#### 4.1.1 Backup PostgreSQL
```bash
# Backup ทั้ง Database
pg_dump -U admin -h localhost -d autorepair > backup_$(date +%Y%m%d).sql

# Backup เฉพาะ Schema
pg_dump -U admin -h localhost -d autorepair --schema-only > schema_backup.sql

# Backup เฉพาะ Data
pg_dump -U admin -h localhost -d autorepair --data-only > data_backup.sql
```

#### 4.1.2 Restore Database
```bash
# Restore จากไฟล์ Backup
psql -U admin -h localhost -d autorepair < backup_20260704.sql
```

#### 4.1.3 กำหนดการ Backup อัตโนมัติ (Cron Job)
```bash
# เพิ่มใน crontab (ทุกวัน เวลา 02:00)
0 2 * * * /usr/local/bin/backup-script.sh

# backup-script.sh
#!/bin/bash
BACKUP_DIR="/backups/postgres"
DATE=$(date +%Y%m%d_%H%M%S)
pg_dump -U admin -h localhost -d autorepair > ${BACKUP_DIR}/backup_${DATE}.sql
gzip ${BACKUP_DIR}/backup_${DATE}.sql
# ลบไฟล์ที่เก่ากว่า 7 วัน
find ${BACKUP_DIR} -name "*.sql.gz" -mtime +7 -delete
```

### 4.2 การ Cleanup Logs

#### 4.2.1 การล้าง Log เก่าใน MongoDB
```javascript
// ล้าง Log ที่เก่ากว่า 90 วัน
db.error_log.deleteMany({
    timestamp: { $lt: new Date(Date.now() - 90 * 24 * 60 * 60 * 1000) }
});

db.method_call_log.deleteMany({
    timestamp: { $lt: new Date(Date.now() - 90 * 24 * 60 * 60 * 1000) }
});
```

#### 4.2.2 การล้าง Application Logs
```bash
# ล้าง Log เก่า (เก็บไว้ 30 วัน)
find /var/log/auto-repair -name "*.log" -mtime +30 -delete
```

### 4.3 การ Update Dependencies

#### 4.3.1 ตรวจสอบ Dependencies Version
```bash
mvn versions:display-dependency-updates
```

#### 4.3.2 อัปเดต Dependencies
```bash
mvn versions:update-properties
mvn clean compile
mvn test
```

### 4.4 การจัดการ SSL Certificate (Production)

```bash
# สร้าง Keystore
keytool -genkey -alias autorepair -keyalg RSA -keystore keystore.p12 -storetype PKCS12

# อัปเดต Certificate (Let's Encrypt)
certbot renew
# แล้วแปลงเป็น PKCS12
openssl pkcs12 -export -in fullchain.pem -inkey privkey.pem -out keystore.p12 -name autorepair
```

### 4.5 การจัดตารางบำรุงรักษา

| งาน | ความถี่ | เวลา | ผู้รับผิดชอบ |
|-----|---------|------|-------------|
| Backup Database | ทุกวัน | 02:00 | DevOps |
| Cleanup Logs | ทุกสัปดาห์ | 03:00 | System Admin |
| Update Dependencies | ทุกเดือน | - | Developer |
| Security Patch | ตามประกาศ | - | DevOps |
| Performance Review | ทุกไตรมาส | - | Architect |

---

## 5. การขยายระบบ (System Expansion)

### 5.1 การเพิ่มโมดูลใหม่

#### 5.1.1 ขั้นตอนการเพิ่มโมดูล
1. สร้างโครงสร้างโฟลเดอร์:
```
src/main/java/com/icmon/module/{new_module}/
├── application/
│   ├── interfaces/
│   └── impl/
├── domain/
├── infrastructure/
│   ├── repository/
│   ├── entity/
│   └── mapper/
└── presentation/
    ├── controller/
    └── dto/
```

2. สร้าง Entity และ Repository ตาม DDD
3. สร้าง Service และ Implementation
4. สร้าง Controller และ DTO
5. เพิ่ม API Documentation (Swagger)
6. เขียน Unit Tests
7. อัปเดต Documentation

#### 5.1.2 ตัวอย่างการเพิ่มโมดูลใหม่
```bash
# สร้างโครงสร้างโมดูลใหม่
mkdir -p src/main/java/com/icmon/module/report
mkdir -p src/main/java/com/icmon/module/report/{application/{interfaces,impl},domain,infrastructure/{repository,entity,mapper},presentation/{controller,dto}}

# สร้างไฟล์หลัก
touch src/main/java/com/icmon/module/report/application/interfaces/ReportService.java
touch src/main/java/com/icmon/module/report/application/impl/ReportServiceImpl.java
touch src/main/java/com/icmon/module/report/presentation/controller/ReportController.java
```

### 5.2 การปรับขนาดระบบ (Horizontal Scaling)

#### 5.2.1 ระบบที่ Scale ได้
- **Application Server**: Scale ออกแนวนอน (เพิ่ม Instance)
- **PostgreSQL**: ใช้ Read Replica
- **Redis**: ใช้ Redis Cluster
- **Kafka**: เพิ่ม Partition

#### 5.2.2 การ Scale Application
```yaml
# docker-compose.prod.yml
services:
  app:
    deploy:
      replicas: 3
      resources:
        limits:
          cpus: '2'
          memory: 2G
```

#### 5.2.3 การ Scale Database
```sql
-- PostgreSQL Read Replica
CREATE SUBSCRIPTION autorepair_subscription 
CONNECTION 'host=master-db port=5432 dbname=autorepair user=repl password=repl' 
PUBLICATION autorepair_publication;
```

### 5.3 การเพิ่ม API Endpoints

#### 5.3.1 ตัวอย่างการเพิ่ม Endpoint ใหม่
```java
@RestController
@RequestMapping("/api/v1/reports")
@Tag(name = "Reports", description = "Report APIs")
public class ReportController {

    @GetMapping("/daily-sales")
    @Operation(summary = "Get daily sales report")
    @RateLimit(limit = 30, duration = 60)
    public ResponseEntity<DailySalesReportDTO> getDailySalesReport(
            @RequestParam LocalDate date) {
        return ResponseEntity.ok(reportService.getDailySales(date));
    }
}
```

### 5.4 การเพิ่ม Table ใหม่

```sql
-- ตัวอย่าง Migration File
-- V15__add_report_table.sql

CREATE TABLE IF NOT EXISTS t_daily_report (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    report_date DATE NOT NULL,
    total_revenue DECIMAL(15,2),
    total_jobs INTEGER,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    whitelabel_id UUID NOT NULL
);

CREATE INDEX idx_t_daily_report_date ON t_daily_report(report_date);
```

---

## 6. การ Monitoring (การตรวจสอบ)

### 6.1 ระบบ Monitoring

| ระบบ | เทคโนโลยี | หน้าที่ |
|------|----------|--------|
| **Metrics** | Micrometer + Prometheus | เก็บ Metrics (CPU, Memory, Request) |
| **Dashboard** | Grafana | แสดงผล Real-time |
| **Logging** | ELK Stack | เก็บและค้นหา Logs |
| **Tracing** | Zipkin / Jaeger | Distributed Tracing |
| **Alerting** | AlertManager | แจ้งเตือนผิดปกติ |

### 6.2 Metrics ที่ต้องตรวจสอบ

#### 6.2.1 System Metrics
| Metric | คำอธิบาย | Threshold |
|--------|----------|-----------|
| `system.cpu.usage` | CPU Usage | > 80% |
| `system.memory.usage` | Memory Usage | > 85% |
| `system.disk.usage` | Disk Usage | > 80% |

#### 6.2.2 Application Metrics
| Metric | คำอธิบาย | Threshold |
|--------|----------|-----------|
| `http.server.requests` | HTTP Request Count | - |
| `http.server.requests.duration` | Response Time | > 2000ms |
| `http.server.requests.error` | Error Rate | > 5% |
| `jdbc.connections.active` | Active DB Connections | > 80% |
| `jdbc.connections.idle` | Idle DB Connections | < 10% |
| `cache.hit.ratio` | Cache Hit Ratio | < 80% |

#### 6.2.3 Business Metrics
| Metric | คำอธิบาย | Threshold |
|--------|----------|-----------|
| `jobs.created` | จำนวน Job ที่สร้าง | - |
| `jobs.completed` | จำนวน Job ที่เสร็จ | - |
| `revenue.daily` | รายได้รายวัน | - |
| `revenue.monthly` | รายได้รายเดือน | - |

### 6.3 การตั้งค่า Prometheus

```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'spring-boot-app'
    metrics_path: '/actuator/prometheus'
    static_configs:
      - targets: ['localhost:5000']
        labels:
          application: 'auto-repair-system'
          environment: 'prod'
```

### 6.4 Grafana Dashboard

#### 6.4.1 Dashboard Variables
```json
{
  "templating": {
    "list": [
      {
        "name": "environment",
        "type": "query",
        "query": "label_values(environment)",
        "current": {"value": "prod"}
      },
      {
        "name": "instance",
        "type": "query",
        "query": "label_values(instance)",
        "current": {"value": "app-01"}
      }
    ]
  }
}
```

#### 6.4.2 Dashboard Panels
1. **System Overview**: CPU, Memory, Disk
2. **Request Metrics**: Request Rate, Response Time, Error Rate
3. **Database Metrics**: Connection Pool, Query Performance
4. **Business Metrics**: Jobs Created/Completed, Revenue
5. **Cache Metrics**: Hit/Miss Ratio, Size

### 6.5 การตั้งค่า Alerting (AlertManager)

```yaml
# alertmanager.yml
route:
  group_by: ['alertname']
  group_wait: 30s
  group_interval: 5m
  repeat_interval: 1h
  receiver: 'slack-notifications'

receivers:
  - name: 'slack-notifications'
    slack_configs:
      - channel: '#alerts'
        api_url: 'https://hooks.slack.com/services/...'
        title: '{{ .GroupLabels.alertname }}'
        text: '{{ range .Alerts }}{{ .Annotations.summary }}\n{{ .Annotations.description }}{{ end }}'
```

### 6.6 Alert Rules (Prometheus)

```yaml
groups:
  - name: application_alerts
    rules:
      - alert: HighErrorRate
        expr: rate(http_server_requests_seconds_count{status=~"5.."}[5m]) > 0.05
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High error rate detected"
          description: "Error rate is {{ $value }}%"

      - alert: HighResponseTime
        expr: histogram_quantile(0.95, rate(http_server_requests_seconds_bucket[5m])) > 2
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High response time detected"
          description: "95th percentile response time is {{ $value }}s"

      - alert: LowCacheHitRatio
        expr: cache_hit_ratio < 0.8
        for: 10m
        labels:
          severity: warning
        annotations:
          summary: "Low cache hit ratio"
          description: "Cache hit ratio is {{ $value }}%"

      - alert: HighDBConnection
        expr: jdbc_connections_active / jdbc_connections_max > 0.8
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High database connections"
          description: "Active DB connections: {{ $value }}%"
```

### 6.7 Logging (ELK Stack)

#### 6.7.1 Logback Configuration
```xml
<!-- logback-spring.xml -->
<configuration>
    <property name="LOG_PATH" value="${LOG_PATH:-/var/log/auto-repair}"/>
    <property name="LOG_FILE" value="${LOG_PATH}/application.log"/>

    <appender name="FILE" class="ch.qos.logback.core.rolling.RollingFileAppender">
        <file>${LOG_FILE}</file>
        <rollingPolicy class="ch.qos.logback.core.rolling.TimeBasedRollingPolicy">
            <fileNamePattern>${LOG_PATH}/application.%d{yyyy-MM-dd}.log.gz</fileNamePattern>
            <maxHistory>30</maxHistory>
        </rollingPolicy>
        <encoder>
            <pattern>%d{yyyy-MM-dd HH:mm:ss.SSS} [%thread] %-5level %logger{36} - %msg%n</pattern>
        </encoder>
    </appender>

    <appender name="CONSOLE" class="ch.qos.logback.core.ConsoleAppender">
        <encoder>
            <pattern>%d{yyyy-MM-dd HH:mm:ss.SSS} [%thread] %-5level %logger{36} - %msg%n</pattern>
        </encoder>
    </appender>

    <root level="INFO">
        <appender-ref ref="CONSOLE"/>
        <appender-ref ref="FILE"/>
    </root>

    <logger name="com.icmon" level="DEBUG"/>
    <logger name="org.springframework" level="WARN"/>
    <logger name="org.hibernate" level="WARN"/>
</configuration>
```

### 6.8 Application Health Check

```java
@Component
public class CustomHealthIndicator implements HealthIndicator {
    
    private final DataSource dataSource;
    private final RedisTemplate<String, String> redisTemplate;
    
    @Override
    public Health health() {
        try {
            // ตรวจสอบ Database
            dataSource.getConnection().isValid(3);
            
            // ตรวจสอบ Redis
            redisTemplate.opsForValue().get("health_check");
            
            return Health.up()
                    .withDetail("database", "connected")
                    .withDetail("redis", "connected")
                    .withDetail("status", "UP")
                    .build();
        } catch (Exception e) {
            return Health.down()
                    .withDetail("error", e.getMessage())
                    .build();
        }
    }
}
```

---

## 7. การแก้ไขปัญหาเบื้องต้น (Troubleshooting)

### 7.1 ปัญหาทั่วไปและแนวทางแก้ไข

| ปัญหา | สาเหตุ | วิธีแก้ไข |
|-------|--------|----------|
| **Application ไม่ Start** | Circular Dependency | ตรวจสอบ Dependency Cycle, ใช้ @Lazy |
| **Database Connection Error** | PostgreSQL ไม่ทำงาน | ตรวจสอบ Docker Compose, Restart PostgreSQL |
| **OutOfMemoryError** | Memory Leak | เพิ่ม Heap Size, ตรวจสอบ Memory Usage |
| **Slow Response Time** | Query Performance | ตรวจสอบ Index, Optimize Query |
| **Cache ไม่ทำงาน** | Redis ไม่ทำงาน | ตรวจสอบ Redis Connection |
| **Kafka Consumer Lag** | Consumer ช้า | เพิ่ม Partition, เพิ่ม Consumer |

### 7.2 การ Debug

#### 7.2.1 Remote Debug
```bash
# รัน Application ด้วย Remote Debug
java -agentlib:jdwp=transport=dt_socket,server=y,suspend=n,address=5005 -jar app.jar

# เชื่อมต่อจาก IDE ด้วย Port 5005
```

#### 7.2.2 การดู Log
```bash
# ดู Log แบบ Real-time
tail -f /var/log/auto-repair/application.log

# ค้นหาคำ Error
grep "ERROR" /var/log/auto-repair/application.log

# ดู Log เฉพาะวันที่
cat /var/log/auto-repair/application.2026-07-04.log
```

### 7.3 การ Rollback

```bash
# Rollback Docker Image
docker pull ${REGISTRY}/auto-repair-system:${PREVIOUS_VERSION}
kubectl set image deployment/auto-repair-system-prod app=${REGISTRY}/auto-repair-system:${PREVIOUS_VERSION} -n prod

# Rollback JAR
cp previous-version/app.jar current-version/app.jar
systemctl restart auto-repair
```

---

## 📞 ติดต่อทีมสนับสนุน

| บทบาท | ชื่อ | อีเมล |
|-------|------|-------|
| DevOps Lead | - | devops@company.com |
| Backend Lead | - | backend@company.com |
| System Admin | - | sysadmin@company.com |

---

**ผู้เขียน:** Kongnakorn Jantakun  
**วันที่:** 2026-07-04  
**เวอร์ชัน:** 1.0  
**สถานะ:** ฉบับสมบูรณ์ ✅

---
# Development, Deployment, Maintenance & Monitoring Manual – English

## 📚 Table of Contents

1. [Introduction](#1-introduction)
2. [Development](#2-development)
3. [Deployment](#3-deployment)
4. [Maintenance](#4-maintenance)
5. [System Expansion](#5-system-expansion)
6. [Monitoring](#6-monitoring)
7. [Troubleshooting](#7-troubleshooting)

---

## 1. Introduction

### 1.1 Purpose
This manual provides guidelines for developing, deploying, maintaining, expanding, and monitoring the Auto Repair Shop Management System.

### 1.2 Audience
- Developers
- DevOps Engineers
- System Administrators
- Maintenance Team

### 1.3 Overall Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                     Frontend (React / Angular / Mobile)             │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                     API Gateway (Spring Cloud Gateway / NGINX)      │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                     Spring Boot Application (REST API)              │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐│
│  │  Auth    │ │  Job     │ │  Quot    │ │  Invoice │ │  Payment ││
│  │  Module  │ │  Module  │ │  Module  │ │  Module  │ │  Module  ││
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘ └──────────┘│
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐│
│  │Inventory │ │  PO      │ │Customer  │ │Document  │ │   IoT    ││
│  │  Module  │ │  Module  │ │  Module  │ │  Module  │ │  Module  ││
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘ └──────────┘│
└─────────────────────────────────────────────────────────────────────┘
                                    │
            ┌───────────────────────┼───────────────────────┐
            ▼                       ▼                       ▼
┌───────────────────┐   ┌───────────────────┐   ┌───────────────────┐
│    PostgreSQL     │   │      Redis        │   │      Kafka        │
│   (Main Database) │   │     (Cache)       │   │   (Event Bus)     │
└───────────────────┘   └───────────────────┘   └───────────────────┘
            │                       │                       │
            ▼                       ▼                       ▼
┌───────────────────┐   ┌───────────────────┐   ┌───────────────────┐
│    MongoDB        │   │     InfluxDB      │   │   Elasticsearch   │
│   (Logging)       │   │    (IoT Data)     │   │   (Search/Logs)   │
└───────────────────┘   └───────────────────┘   └───────────────────┘
```

### 1.4 Deployment Diagram

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         AWS CLOUD (EC2 / ECS)                          │
│                                                                         │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                     Load Balancer (ALB)                          │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                    │                                     │
│              ┌─────────────────────┼─────────────────────┐              │
│              ▼                     ▼                     ▼              │
│  ┌───────────────────┐ ┌───────────────────┐ ┌───────────────────┐    │
│  │   Instance 1      │ │   Instance 2      │ │   Instance 3      │    │
│  │  (App Server)     │ │  (App Server)     │ │  (App Server)     │    │
│  └───────────────────┘ └───────────────────┘ └───────────────────┘    │
│              │                     │                     │              │
│              └─────────────────────┼─────────────────────┘              │
│                                    ▼                                     │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                     RDS (PostgreSQL)                            │   │
│  │                     ElastiCache (Redis)                         │   │
│  │                     MSK (Kafka)                                 │   │
│  │                     S3 (Document Storage)                       │   │
│  └─────────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Development

### 2.1 Prerequisites

| Item | Version |
|------|---------|
| Java | 21+ |
| Maven | 3.8+ |
| Docker | 20.10+ |
| Docker Compose | 2.0+ |
| Git | 2.30+ |
| IntelliJ IDEA (recommended) | 2023.2+ |

### 2.2 Setup for Development

#### 2.2.1 Clone Project
```bash
git clone https://github.com/your-org/auto-repair-system.git
cd auto-repair-system
```

#### 2.2.2 Configure Environment Variables
```bash
cp .env.example .env
# Edit .env file to match your environment
```

#### 2.2.3 Start Databases with Docker Compose
```bash
docker-compose up -d
```

Services started:
- PostgreSQL (port 5432)
- Redis (port 6379)
- Kafka + Zookeeper (port 9092)
- MongoDB (port 27017)
- InfluxDB (port 8086)
- Elasticsearch (port 9200)
- Kibana (port 5601)
- Grafana (port 3000)

#### 2.2.4 Run Application
```bash
# Development Profile
mvn spring-boot:run -Dspring-boot.run.profiles=dev

# Or run with IntelliJ IDEA selecting "dev" profile
```

#### 2.2.5 Access Swagger UI
```
http://localhost:5000/api/swagger-ui.html
```

### 2.3 Testing

#### 2.3.1 Run Unit Tests
```bash
mvn test
```

#### 2.3.2 Run Specific Module Tests
```bash
mvn test -Dtest="com.icmon.module.job.**"
```

#### 2.3.3 View Coverage Report
```bash
mvn jacoco:report
# Open target/site/jacoco/index.html
```

#### 2.3.4 Test Coverage Target
| Layer | Target |
|-------|--------|
| Domain | ≥ 95% |
| Application | ≥ 90% |
| Infrastructure | ≥ 70% |
| Presentation | ≥ 80% |
| **Overall** | **≥ 80%** |

### 2.4 Coding Standards

#### 2.4.1 Naming Conventions
| Type | Format | Example |
|------|--------|---------|
| Package | lowercase | `com.icmon.module.job` |
| Class | PascalCase | `JobServiceImpl` |
| Interface | PascalCase | `JobService` |
| Method | camelCase | `createJob()` |
| Variable | camelCase | `jobId` |
| Constant | UPPER_SNAKE_CASE | `MAX_RETRY_COUNT` |

#### 2.4.2 Package Structure
```
com.icmon/
├── configuration/          # Spring Configurations
├── module/
│   └── {module_name}/
│       ├── application/    # Service Layer
│       │   ├── interfaces/
│       │   └── impl/
│       ├── domain/         # Domain Layer
│       ├── infrastructure/ # Infrastructure Layer
│       │   ├── repository/
│       │   ├── entity/
│       │   └── mapper/
│       └── presentation/   # Presentation Layer
│           ├── controller/
│           └── dto/
├── exception/              # Global Exception
├── logging/                # Logging System
└── utils/                  # Utilities
```

### 2.5 Git Workflow

```
main (Production)
  │
  ├── develop (Integration)
  │   ├── feature/job-management    ← Feature Development
  │   ├── feature/quotation-module
  │   └── feature/inventory
  │
  ├── release/v2.0.0
  │
  └── hotfix/critical-bug
```

#### Common Commands
```bash
# Create new branch for feature development
git checkout -b feature/your-feature-name develop

# Commit and Push
git add .
git commit -m "feat: add new feature description"
git push origin feature/your-feature-name

# Create Pull Request on GitHub/GitLab
# After merge, delete the branch
```

#### Commit Message Format
```
<type>(<scope>): <subject>

[optional body]

[optional footer]
```

**Type:** `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`

**Example:**
```
feat(job): add job status history tracking

- Add TJobStatusHistory entity
- Create JobStatusHistoryRepository
- Implement status change logging

Closes #123
```

---

## 3. Deployment

### 3.1 Environment Profiles

| Profile | File | Usage |
|---------|------|-------|
| **local** | `application-local.yml` | Local Development |
| **dev** | `application-dev.yml` | Development Environment |
| **test** | `application-test.yml` | Test/UAT Environment |
| **prod** | `application-prod.yml` | Production Environment |

### 3.2 Database Configuration

#### 3.2.1 Development (Local)
```yaml
spring:
  datasource:
    url: jdbc:postgresql://localhost:5432/autorepair
    username: admin
    password: admin123
  jpa:
    hibernate:
      ddl-auto: update
    show-sql: true
```

#### 3.2.2 UAT/Test
```yaml
spring:
  datasource:
    url: ${DB_URL}
    username: ${DB_USERNAME}
    password: ${DB_PASSWORD}
  jpa:
    hibernate:
      ddl-auto: validate
    show-sql: false
```

#### 3.2.3 Production
```yaml
spring:
  datasource:
    url: ${DB_URL}
    username: ${DB_USERNAME}
    password: ${DB_PASSWORD}
  jpa:
    hibernate:
      ddl-auto: validate
    show-sql: false
    properties:
      hibernate:
        jdbc:
          batch_size: 50
        order_inserts: true
        order_updates: true
```

### 3.3 Building

#### 3.3.1 Build for Development
```bash
mvn clean compile
```

#### 3.3.2 Build for Production
```bash
mvn clean package -Pprod
```

#### 3.3.3 Build Docker Image
```bash
# Build Docker Image
docker build -t auto-repair-system:${VERSION} .

# Push to Registry
docker tag auto-repair-system:${VERSION} ${REGISTRY}/auto-repair-system:${VERSION}
docker push ${REGISTRY}/auto-repair-system:${VERSION}
```

### 3.4 Deploy per Environment

#### 3.4.1 Local Development
```bash
mvn spring-boot:run -Dspring-boot.run.profiles=local
```

#### 3.4.2 Development (Dev) Server
```bash
# Create .env.dev file
# Run with Docker Compose
docker-compose -f docker-compose.dev.yml up -d
```

#### 3.4.3 UAT/Test Server
```bash
# 1. Build JAR
mvn clean package -Ptest

# 2. Deploy JAR
java -jar target/app-1.0.0-SNAPSHOT.jar \
  --spring.profiles.active=test \
  --server.port=5000
```

#### 3.4.4 Production Server
```bash
# 1. Build JAR with Production Profile
mvn clean package -Pprod

# 2. Deploy with Production Config
java -jar target/app-1.0.0-SNAPSHOT.jar \
  --spring.profiles.active=prod \
  --server.port=8443 \
  --server.ssl.enabled=true \
  --server.ssl.key-store=/path/to/keystore.p12
```

### 3.5 CI/CD Pipeline (Jenkins)

```groovy
pipeline {
    agent any
    
    environment {
        REGISTRY = 'docker.io/your-registry'
        IMAGE_NAME = 'auto-repair-system'
    }
    
    stages {
        stage('Checkout') {
            steps {
                git branch: 'main', url: 'https://github.com/your-org/auto-repair-system.git'
            }
        }
        
        stage('Build') {
            steps {
                sh 'mvn clean compile'
            }
        }
        
        stage('Unit Tests') {
            steps {
                sh 'mvn test'
            }
        }
        
        stage('Integration Tests') {
            steps {
                sh 'mvn test -Dgroups="integration"'
            }
        }
        
        stage('Build JAR') {
            steps {
                sh 'mvn package -Pprod'
            }
        }
        
        stage('Build Docker Image') {
            steps {
                sh '''
                    docker build -t ${REGISTRY}/${IMAGE_NAME}:${BUILD_NUMBER} .
                    docker tag ${REGISTRY}/${IMAGE_NAME}:${BUILD_NUMBER} ${REGISTRY}/${IMAGE_NAME}:latest
                '''
            }
        }
        
        stage('Push to Registry') {
            steps {
                sh '''
                    docker push ${REGISTRY}/${IMAGE_NAME}:${BUILD_NUMBER}
                    docker push ${REGISTRY}/${IMAGE_NAME}:latest
                '''
            }
        }
        
        stage('Deploy to DEV') {
            when { branch 'develop' }
            steps {
                sh 'kubectl set image deployment/auto-repair-system-dev app=${REGISTRY}/${IMAGE_NAME}:${BUILD_NUMBER} -n dev'
            }
        }
        
        stage('Deploy to UAT') {
            when { branch 'release/*' }
            steps {
                sh 'kubectl set image deployment/auto-repair-system-uat app=${REGISTRY}/${IMAGE_NAME}:${BUILD_NUMBER} -n uat'
            }
        }
        
        stage('Deploy to Production') {
            when { branch 'main' }
            steps {
                input message: 'Deploy to Production?', ok: 'Yes'
                sh 'kubectl set image deployment/auto-repair-system-prod app=${REGISTRY}/${IMAGE_NAME}:${BUILD_NUMBER} -n prod'
            }
        }
    }
}
```

### 3.6 Docker Compose for Production

```yaml
version: '3.8'

services:
  app:
    image: ${REGISTRY}/auto-repair-system:latest
    restart: always
    ports:
      - "8443:8443"
    environment:
      SPRING_PROFILES_ACTIVE: prod
      DB_URL: jdbc:postgresql://postgres:5432/autorepair
      DB_USERNAME: ${DB_USERNAME}
      DB_PASSWORD: ${DB_PASSWORD}
      REDIS_HOST: redis
      REDIS_PORT: 6379
      KAFKA_BOOTSTRAP_SERVERS: kafka:9092
    volumes:
      - ./logs:/app/logs
      - ./uploads:/app/uploads
    depends_on:
      - postgres
      - redis
      - kafka
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8443/actuator/health"]
      interval: 30s
      timeout: 10s
      retries: 5

  postgres:
    image: postgres:15
    restart: always
    environment:
      POSTGRES_DB: autorepair
      POSTGRES_USER: ${DB_USERNAME}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"

  redis:
    image: redis:7-alpine
    restart: always
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

  kafka:
    image: confluentinc/cp-kafka:latest
    restart: always
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:9092
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
    ports:
      - "9092:9092"
    depends_on:
      - zookeeper

  zookeeper:
    image: confluentinc/cp-zookeeper:latest
    restart: always
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181
      ZOOKEEPER_TICK_TIME: 2000
    ports:
      - "2181:2181"

volumes:
  postgres_data:
  redis_data:
```

---

## 4. Maintenance

### 4.1 Database Backup

#### 4.1.1 Backup PostgreSQL
```bash
# Backup entire database
pg_dump -U admin -h localhost -d autorepair > backup_$(date +%Y%m%d).sql

# Backup only schema
pg_dump -U admin -h localhost -d autorepair --schema-only > schema_backup.sql

# Backup only data
pg_dump -U admin -h localhost -d autorepair --data-only > data_backup.sql
```

#### 4.1.2 Restore Database
```bash
# Restore from backup file
psql -U admin -h localhost -d autorepair < backup_20260704.sql
```

#### 4.1.3 Automatic Backup (Cron Job)
```bash
# Add to crontab (every day at 02:00)
0 2 * * * /usr/local/bin/backup-script.sh

# backup-script.sh
#!/bin/bash
BACKUP_DIR="/backups/postgres"
DATE=$(date +%Y%m%d_%H%M%S)
pg_dump -U admin -h localhost -d autorepair > ${BACKUP_DIR}/backup_${DATE}.sql
gzip ${BACKUP_DIR}/backup_${DATE}.sql
# Delete files older than 7 days
find ${BACKUP_DIR} -name "*.sql.gz" -mtime +7 -delete
```

### 4.2 Log Cleanup

#### 4.2.1 Clean Old Logs in MongoDB
```javascript
// Delete logs older than 90 days
db.error_log.deleteMany({
    timestamp: { $lt: new Date(Date.now() - 90 * 24 * 60 * 60 * 1000) }
});

db.method_call_log.deleteMany({
    timestamp: { $lt: new Date(Date.now() - 90 * 24 * 60 * 60 * 1000) }
});
```

#### 4.2.2 Clean Application Logs
```bash
# Delete old logs (keep 30 days)
find /var/log/auto-repair -name "*.log" -mtime +30 -delete
```

### 4.3 Dependency Updates

#### 4.3.1 Check Dependency Updates
```bash
mvn versions:display-dependency-updates
```

#### 4.3.2 Update Dependencies
```bash
mvn versions:update-properties
mvn clean compile
mvn test
```

### 4.4 SSL Certificate Management (Production)

```bash
# Create Keystore
keytool -genkey -alias autorepair -keyalg RSA -keystore keystore.p12 -storetype PKCS12

# Update Certificate (Let's Encrypt)
certbot renew
# Convert to PKCS12
openssl pkcs12 -export -in fullchain.pem -inkey privkey.pem -out keystore.p12 -name autorepair
```

### 4.5 Maintenance Schedule

| Task | Frequency | Time | Responsible |
|------|-----------|------|-------------|
| Database Backup | Daily | 02:00 | DevOps |
| Log Cleanup | Weekly | 03:00 | System Admin |
| Dependency Updates | Monthly | - | Developer |
| Security Patches | As needed | - | DevOps |
| Performance Review | Quarterly | - | Architect |

---

## 5. System Expansion

### 5.1 Adding New Modules

#### 5.1.1 Steps to Add a Module
1. Create folder structure:
```
src/main/java/com/icmon/module/{new_module}/
├── application/
│   ├── interfaces/
│   └── impl/
├── domain/
├── infrastructure/
│   ├── repository/
│   ├── entity/
│   └── mapper/
└── presentation/
    ├── controller/
    └── dto/
```

2. Create Entity and Repository following DDD
3. Create Service and Implementation
4. Create Controller and DTO
5. Add API Documentation (Swagger)
6. Write Unit Tests
7. Update Documentation

#### 5.1.2 Example: Adding New Module
```bash
# Create new module structure
mkdir -p src/main/java/com/icmon/module/report
mkdir -p src/main/java/com/icmon/module/report/{application/{interfaces,impl},domain,infrastructure/{repository,entity,mapper},presentation/{controller,dto}}

# Create main files
touch src/main/java/com/icmon/module/report/application/interfaces/ReportService.java
touch src/main/java/com/icmon/module/report/application/impl/ReportServiceImpl.java
touch src/main/java/com/icmon/module/report/presentation/controller/ReportController.java
```

### 5.2 Horizontal Scaling

#### 5.2.1 Scalable Components
- **Application Server**: Scale out horizontally (add instances)
- **PostgreSQL**: Use Read Replica
- **Redis**: Use Redis Cluster
- **Kafka**: Increase Partitions

#### 5.2.2 Scaling Application
```yaml
# docker-compose.prod.yml
services:
  app:
    deploy:
      replicas: 3
      resources:
        limits:
          cpus: '2'
          memory: 2G
```

#### 5.2.3 Scaling Database
```sql
-- PostgreSQL Read Replica
CREATE SUBSCRIPTION autorepair_subscription 
CONNECTION 'host=master-db port=5432 dbname=autorepair user=repl password=repl' 
PUBLICATION autorepair_publication;
```

### 5.3 Adding API Endpoints

#### 5.3.1 Example: Adding New Endpoint
```java
@RestController
@RequestMapping("/api/v1/reports")
@Tag(name = "Reports", description = "Report APIs")
public class ReportController {

    @GetMapping("/daily-sales")
    @Operation(summary = "Get daily sales report")
    @RateLimit(limit = 30, duration = 60)
    public ResponseEntity<DailySalesReportDTO> getDailySalesReport(
            @RequestParam LocalDate date) {
        return ResponseEntity.ok(reportService.getDailySales(date));
    }
}
```

### 5.4 Adding New Tables

```sql
-- Example Migration File
-- V15__add_report_table.sql

CREATE TABLE IF NOT EXISTS t_daily_report (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    report_date DATE NOT NULL,
    total_revenue DECIMAL(15,2),
    total_jobs INTEGER,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    whitelabel_id UUID NOT NULL
);

CREATE INDEX idx_t_daily_report_date ON t_daily_report(report_date);
```

---

## 6. Monitoring

### 6.1 Monitoring Systems

| System | Technology | Purpose |
|--------|------------|---------|
| **Metrics** | Micrometer + Prometheus | Collect metrics (CPU, Memory, Request) |
| **Dashboard** | Grafana | Real-time display |
| **Logging** | ELK Stack | Store and search logs |
| **Tracing** | Zipkin / Jaeger | Distributed Tracing |
| **Alerting** | AlertManager | Alert on abnormalities |

### 6.2 Metrics to Monitor

#### 6.2.1 System Metrics
| Metric | Description | Threshold |
|--------|-------------|-----------|
| `system.cpu.usage` | CPU Usage | > 80% |
| `system.memory.usage` | Memory Usage | > 85% |
| `system.disk.usage` | Disk Usage | > 80% |

#### 6.2.2 Application Metrics
| Metric | Description | Threshold |
|--------|-------------|-----------|
| `http.server.requests` | HTTP Request Count | - |
| `http.server.requests.duration` | Response Time | > 2000ms |
| `http.server.requests.error` | Error Rate | > 5% |
| `jdbc.connections.active` | Active DB Connections | > 80% |
| `jdbc.connections.idle` | Idle DB Connections | < 10% |
| `cache.hit.ratio` | Cache Hit Ratio | < 80% |

#### 6.2.3 Business Metrics
| Metric | Description | Threshold |
|--------|-------------|-----------|
| `jobs.created` | Number of Jobs Created | - |
| `jobs.completed` | Number of Jobs Completed | - |
| `revenue.daily` | Daily Revenue | - |
| `revenue.monthly` | Monthly Revenue | - |

### 6.3 Prometheus Configuration

```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'spring-boot-app'
    metrics_path: '/actuator/prometheus'
    static_configs:
      - targets: ['localhost:5000']
        labels:
          application: 'auto-repair-system'
          environment: 'prod'
```

### 6.4 Grafana Dashboard

#### 6.4.1 Dashboard Variables
```json
{
  "templating": {
    "list": [
      {
        "name": "environment",
        "type": "query",
        "query": "label_values(environment)",
        "current": {"value": "prod"}
      },
      {
        "name": "instance",
        "type": "query",
        "query": "label_values(instance)",
        "current": {"value": "app-01"}
      }
    ]
  }
}
```

#### 6.4.2 Dashboard Panels
1. **System Overview**: CPU, Memory, Disk
2. **Request Metrics**: Request Rate, Response Time, Error Rate
3. **Database Metrics**: Connection Pool, Query Performance
4. **Business Metrics**: Jobs Created/Completed, Revenue
5. **Cache Metrics**: Hit/Miss Ratio, Size

### 6.5 AlertManager Configuration

```yaml
# alertmanager.yml
route:
  group_by: ['alertname']
  group_wait: 30s
  group_interval: 5m
  repeat_interval: 1h
  receiver: 'slack-notifications'

receivers:
  - name: 'slack-notifications'
    slack_configs:
      - channel: '#alerts'
        api_url: 'https://hooks.slack.com/services/...'
        title: '{{ .GroupLabels.alertname }}'
        text: '{{ range .Alerts }}{{ .Annotations.summary }}\n{{ .Annotations.description }}{{ end }}'
```

### 6.6 Alert Rules (Prometheus)

```yaml
groups:
  - name: application_alerts
    rules:
      - alert: HighErrorRate
        expr: rate(http_server_requests_seconds_count{status=~"5.."}[5m]) > 0.05
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High error rate detected"
          description: "Error rate is {{ $value }}%"

      - alert: HighResponseTime
        expr: histogram_quantile(0.95, rate(http_server_requests_seconds_bucket[5m])) > 2
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High response time detected"
          description: "95th percentile response time is {{ $value }}s"

      - alert: LowCacheHitRatio
        expr: cache_hit_ratio < 0.8
        for: 10m
        labels:
          severity: warning
        annotations:
          summary: "Low cache hit ratio"
          description: "Cache hit ratio is {{ $value }}%"

      - alert: HighDBConnection
        expr: jdbc_connections_active / jdbc_connections_max > 0.8
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High database connections"
          description: "Active DB connections: {{ $value }}%"
```

### 6.7 Logging (ELK Stack)

#### 6.7.1 Logback Configuration
```xml
<!-- logback-spring.xml -->
<configuration>
    <property name="LOG_PATH" value="${LOG_PATH:-/var/log/auto-repair}"/>
    <property name="LOG_FILE" value="${LOG_PATH}/application.log"/>

    <appender name="FILE" class="ch.qos.logback.core.rolling.RollingFileAppender">
        <file>${LOG_FILE}</file>
        <rollingPolicy class="ch.qos.logback.core.rolling.TimeBasedRollingPolicy">
            <fileNamePattern>${LOG_PATH}/application.%d{yyyy-MM-dd}.log.gz</fileNamePattern>
            <maxHistory>30</maxHistory>
        </rollingPolicy>
        <encoder>
            <pattern>%d{yyyy-MM-dd HH:mm:ss.SSS} [%thread] %-5level %logger{36} - %msg%n</pattern>
        </encoder>
    </appender>

    <appender name="CONSOLE" class="ch.qos.logback.core.ConsoleAppender">
        <encoder>
            <pattern>%d{yyyy-MM-dd HH:mm:ss.SSS} [%thread] %-5level %logger{36} - %msg%n</pattern>
        </encoder>
    </appender>

    <root level="INFO">
        <appender-ref ref="CONSOLE"/>
        <appender-ref ref="FILE"/>
    </root>

    <logger name="com.icmon" level="DEBUG"/>
    <logger name="org.springframework" level="WARN"/>
    <logger name="org.hibernate" level="WARN"/>
</configuration>
```

### 6.8 Application Health Check

```java
@Component
public class CustomHealthIndicator implements HealthIndicator {
    
    private final DataSource dataSource;
    private final RedisTemplate<String, String> redisTemplate;
    
    @Override
    public Health health() {
        try {
            // Check Database
            dataSource.getConnection().isValid(3);
            
            // Check Redis
            redisTemplate.opsForValue().get("health_check");
            
            return Health.up()
                    .withDetail("database", "connected")
                    .withDetail("redis", "connected")
                    .withDetail("status", "UP")
                    .build();
        } catch (Exception e) {
            return Health.down()
                    .withDetail("error", e.getMessage())
                    .build();
        }
    }
}
```

---

## 7. Troubleshooting

### 7.1 Common Issues and Solutions

| Issue | Cause | Solution |
|-------|-------|----------|
| **Application doesn't start** | Circular Dependency | Check Dependency Cycle, use @Lazy |
| **Database Connection Error** | PostgreSQL not running | Check Docker Compose, Restart PostgreSQL |
| **OutOfMemoryError** | Memory Leak | Increase Heap Size, Check Memory Usage |
| **Slow Response Time** | Query Performance | Check Indexes, Optimize Queries |
| **Cache not working** | Redis not running | Check Redis Connection |
| **Kafka Consumer Lag** | Slow Consumer | Increase Partitions, Add Consumers |

### 7.2 Debugging

#### 7.2.1 Remote Debug
```bash
# Run Application with Remote Debug
java -agentlib:jdwp=transport=dt_socket,server=y,suspend=n,address=5005 -jar app.jar

# Connect from IDE on Port 5005
```

#### 7.2.2 Viewing Logs
```bash
# View logs in real-time
tail -f /var/log/auto-repair/application.log

# Search for errors
grep "ERROR" /var/log/auto-repair/application.log

# View specific date log
cat /var/log/auto-repair/application.2026-07-04.log
```

### 7.3 Rollback

```bash
# Rollback Docker Image
docker pull ${REGISTRY}/auto-repair-system:${PREVIOUS_VERSION}
kubectl set image deployment/auto-repair-system-prod app=${REGISTRY}/auto-repair-system:${PREVIOUS_VERSION} -n prod

# Rollback JAR
cp previous-version/app.jar current-version/app.jar
systemctl restart auto-repair
```

---

## 📞 Support Team Contacts

| Role | Name | Email |
|------|------|-------|
| DevOps Lead | - | devops@company.com |
| Backend Lead | - | backend@company.com |
| System Admin | - | sysadmin@company.com |

---

**Author:** Kongnakorn Jantakun  
**Date:** 2026-07-04  
**Version:** 1.0  
**Status:** Complete ✅