# ภาคผนวก: การเริ่มต้นโปรเจกต์อย่างรวดเร็ว (Quick Start) และเทมเพลตที่ดาวน์โหลดได้

## สารบัญรวมเอกสารทั้ง 4 เล่ม

| เล่ม | ชื่อ | บทที่ | เนื้อหาหลัก |
|------|------|-------|--------------|
| **เล่ม 1** | ภาคทฤษฎี | 1-3 | บทนำ, นิยามศัพท์, การวิเคราะห์ความต้องการและกรณีศึกษา |
| **เล่ม 2** | สถาปัตยกรรมโครงสร้างระบบ | 1-4 | โครงสร้างโปรเจกต์, Config/Logging/Middleware, Database/Repository/Transaction, Authentication/JWT/RBAC |
| **เล่ม 3** | การพัฒนาเชิงปฏิบัติ | 1-3 | MQTT & Real-time ingestion, WebSocket Dashboard & Visualization, Scheduler & Automation |
| **เล่ม 4** | การปรับใช้และการบำรุงรักษา | 1 | Deployment (Docker/K8s), Graceful shutdown, Monitoring, Backup, Scaling |

---

## Quick Start: การรันระบบ CMON บนเครื่อง Developer ใน 15 นาที

### ขั้นตอนที่ 0: สิ่งที่ต้องติดตั้ง
- Go 1.21+
- Docker และ Docker Compose
- Git
- Make (optional)

### ขั้นตอนที่ 1: Clone โปรเจกต์และสร้าง environment

```bash
git clone https://github.com/your-org/cmon-backend.git
cd cmon-backend
cp .env.example .env
# แก้ไข .env ให้ตรงกับ environment ของคุณ (หรือใช้ค่าตั้งต้น)
```

**ไฟล์ `.env.example`**
```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres123
DB_NAME=cmon_db

# Redis
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# MQTT
MQTT_BROKER=tcp://localhost:1883
MQTT_CLIENT_ID=cmon-backend-dev

# JWT (RSA keys) – ต้องสร้างก่อน
JWT_PRIVATE_KEY_BASE64=LS0tLS1CRUdJTi... (base64 ของ private key)
JWT_PUBLIC_KEY_BASE64=LS0tLS1CRUdJTi... (base64 ของ public key)

# SMTP (สำหรับ email alert)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM=alerts@cmon.local

# Line Notify Token
LINE_TOKEN=your_line_notify_token

# Server
SERVER_PORT=8080
ENVIRONMENT=development
```

### ขั้นตอนที่ 2: สร้าง RSA keys สำหรับ JWT

```bash
# สร้าง private key
openssl genrsa -out private.pem 2048
# สร้าง public key
openssl rsa -in private.pem -pubout -out public.pem
# แปลงเป็น base64 (เอาไปใส่ใน .env)
cat private.pem | base64 | tr -d '\n'
cat public.pem | base64 | tr -d '\n'
```

### ขั้นตอนที่ 3: รัน infrastructure dependencies ด้วย Docker Compose

```bash
docker-compose -f docker-compose.dev.yml up -d postgres redis mqtt
# ตรวจสอบว่า service เริ่มทำงาน
docker-compose ps
```

### ขั้นตอนที่ 4: Run database migration

```bash
go run cmd/migrate.go up
```

### ขั้นตอนที่ 5: รัน Go backend

```bash
# โหมด development (hot-reload)
air

# หรือรันตรงๆ
go run cmd/api/main.go serve
```

### ขั้นตอนที่ 6: ทดสอบ API

```bash
# Health check
curl http://localhost:8080/health/live

# Register user
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"admin123","full_name":"Admin"}'

# Login
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"admin123"}'
# จะได้ access_token และ refresh_token

# เรียก API ที่ต้องใช้ token
curl -X GET http://localhost:8080/profile \
  -H "Authorization: Bearer <access_token>"
```

### ขั้นตอนที่ 7: ทดสอบ MQTT (จำลองเซนเซอร์)

ใช้ MQTT client (เช่น mosquitto_pub หรือ MQTT Explorer)

```bash
# ส่งข้อมูลอุณหภูมิ (จำลอง)
mosquitto_pub -h localhost -p 1883 -t "cmom/dc/bkk01/temperature/rack_a1" \
  -m '{"device_id":"sensor01","sensor_type":"temperature","value":36.5,"unit":"C","location":"rack_a1","timestamp":"2025-04-04T10:00:00Z"}'
```

สังเกต log ของ Go backend จะแสดงการรับข้อมูล และถ้าค่า > 35 จะมีการแจ้งเตือน

### ขั้นตอนที่ 8: เปิด Dashboard

เปิด browser ไปที่ `http://localhost:8080` (ถ้ามี static files ใน `./static`) หรือใช้หน้า HTML ที่สร้างในบทที่ 3.2

---

## เทมเพลตที่ดาวน์โหลดได้ (Markdown)

### เทมเพลต 1: Task List (ใช้สำหรับติดตามความคืบหน้า)

```markdown
# CMON IoT Project Task List

## Milestone 1: Infrastructure & Setup (Week 1-2)
- [ ] MQTT broker (EMQX) installed and configured
- [ ] PostgreSQL database created with migrations
- [ ] Redis instance running
- [ ] Go project structure initialized
- [ ] CI pipeline (GitHub Actions) configured

## Milestone 2: Core Backend (Week 3-5)
- [ ] JWT authentication (login, refresh, logout)
- [ ] User CRUD with role (admin/user)
- [ ] Middleware: logger, rate limit, CORS, recovery
- [ ] MQTT subscriber (receive sensor data)
- [ ] Rule engine (basic threshold)
- [ ] Sensor data storage in PostgreSQL

## Milestone 3: Real-time & Dashboard (Week 6-8)
- [ ] WebSocket hub and client broadcast
- [ ] Dashboard frontend (HTML/JS + Chart.js)
- [ ] Device control API (turn on/off fan/AC)
- [ ] Email and Line notification

## Milestone 4: Automation & Reports (Week 9-10)
- [ ] Scheduler (cron jobs) with distributed lock
- [ ] Report generator (PDF/Excel)
- [ ] Automated email reports

## Milestone 5: Production Readiness (Week 11-12)
- [ ] Docker multi-stage build
- [ ] Docker Compose for production
- [ ] Health checks and graceful shutdown
- [ ] Prometheus metrics + Grafana dashboard
- [ ] Load testing (k6)
- [ ] Documentation and runbook
```

### เทมเพลต 2: Deployment Checklist

```markdown
# Production Deployment Checklist

## Pre-Deployment
- [ ] All tests passed (unit + integration)
- [ ] No secrets in code (use environment variables)
- [ ] Database migrations tested on staging
- [ ] JWT keys are stored in secrets manager
- [ ] CORS origins restricted to known domains
- [ ] Rate limiting configured (e.g., 100 req/sec per IP)

## Infrastructure
- [ ] PostgreSQL max_connections set (>= 100)
- [ ] Redis maxmemory and eviction policy set (allkeys-lru)
- [ ] MQTT broker authentication enabled (if exposed to internet)
- [ ] Load balancer with SSL termination
- [ ] Firewall rules: allow only necessary ports (80, 443, 1883 optional)

## Application Configuration
- [ ] Environment variables: DB, Redis, MQTT, SMTP, Line token
- [ ] Log level = "info" (not debug)
- [ ] JWT access token expiry <= 15 minutes
- [ ] HTTP timeouts: read=15s, write=15s, idle=60s
- [ ] Graceful shutdown timeout >= 30 seconds

## Monitoring & Alerting
- [ ] Prometheus metrics endpoint (/metrics) exposed
- [ ] Grafana dashboard imported (or created)
- [ ] Alert rules: error rate > 1%, high latency, down services
- [ ] Log aggregation (Loki/ELK) configured

## Backup & Disaster Recovery
- [ ] PostgreSQL backup script (pg_dump) scheduled daily
- [ ] Backups stored offsite (S3, separate disk)
- [ ] Restore procedure tested
- [ ] Runbook includes rollback steps

## Security
- [ ] HTTPS enforced (HSTS header)
- [ ] Security headers: X-Frame-Options, X-XSS-Protection, CSP
- [ ] Rate limiting on login endpoint
- [ ] Audit log for critical actions (login, control device)

## Post-Deployment Validation
- [ ] Health endpoints (/live, /ready) return 200
- [ ] WebSocket connection works from dashboard
- [ ] MQTT sensor data flows to dashboard
- [ ] Alert emails received and contain correct info
```

### เทมเพลต 3: Timeline Project (Gantt ใน Markdown)

```markdown
# Project Timeline – CMON IoT Implementation

| Week | Start Date | Activities | Deliverables | Owner |
|------|------------|------------|--------------|-------|
| W1 | 2025-04-07 | Requirements gathering, architecture design, sensor selection | PRD, Arch Diagram | PM + SA |
| W2 | 2025-04-14 | Setup dev environment, Git, CI, base Go project | Repo, docker-compose dev | Backend |
| W3 | 2025-04-21 | JWT auth, user management, middleware | Login API | Backend |
| W4 | 2025-04-28 | MQTT subscriber, store sensor data, basic rule engine | Data ingestion | Backend |
| W5 | 2025-05-05 | Rule engine advanced (multi-condition), notifiers (Email, Line) | Alert system | Backend |
| W6 | 2025-05-12 | WebSocket hub, dashboard frontend (real-time charts) | Live dashboard | Frontend + Backend |
| W7 | 2025-05-19 | Device control API, MQTT publish, integrate with dashboard | Control panel | Fullstack |
| W8 | 2025-05-26 | Scheduler (cron), distributed lock, report generator (PDF) | Automation | Backend |
| W9 | 2025-06-02 | Integration testing, load testing (k6), security audit | Test report | QA |
| W10 | 2025-06-09 | Docker production build, staging deployment, UAT | Staging live | DevOps |
| W11 | 2025-06-16 | UAT feedback, bug fixes, documentation | Sign-off | Team |
| W12 | 2025-06-23 | Production deployment, go-live, handover | Production system | DevOps + Team |
```

---

## ตัวอย่างไฟล์ docker-compose เต็มรูปแบบสำหรับ Production

**docker-compose.prod.yml** (ขยายจากบทที่ 4)

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    container_name: cmon-postgres
    restart: always
    environment:
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
      POSTGRES_DB: ${DB_NAME}
    volumes:
      - pg_data:/var/lib/postgresql/data
      - ./backup:/backup
    networks:
      - cmon_network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER}"]
      interval: 10s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    container_name: cmon-redis
    restart: always
    command: redis-server --appendonly yes --maxmemory 512mb --maxmemory-policy allkeys-lru
    volumes:
      - redis_data:/data
    networks:
      - cmon_network
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 3s
      retries: 3

  emqx:
    image: emqx/emqx:5.0.26
    container_name: cmon-emqx
    restart: always
    environment:
      EMQX_NAME: emqx
      EMQX_HOST: 0.0.0.0
    ports:
      - "1883:1883"   # MQTT
      - "8083:8083"   # WebSocket over MQTT
      - "18083:18083" # Dashboard
    volumes:
      - emqx_data:/opt/emqx/data
    networks:
      - cmon_network

  api:
    image: ${REGISTRY:-cmon}/api:${TAG:-latest}
    container_name: cmon-api
    restart: always
    environment:
      DB_HOST: postgres
      DB_PORT: 5432
      DB_USER: ${DB_USER}
      DB_PASSWORD: ${DB_PASSWORD}
      DB_NAME: ${DB_NAME}
      REDIS_ADDR: redis:6379
      MQTT_BROKER: tcp://emqx:1883
      JWT_PRIVATE_KEY_BASE64: ${JWT_PRIVATE_KEY_BASE64}
      JWT_PUBLIC_KEY_BASE64: ${JWT_PUBLIC_KEY_BASE64}
      SMTP_HOST: ${SMTP_HOST}
      SMTP_PORT: ${SMTP_PORT}
      SMTP_USER: ${SMTP_USER}
      SMTP_PASSWORD: ${SMTP_PASSWORD}
      LINE_TOKEN: ${LINE_TOKEN}
      SERVER_PORT: 8080
      ENVIRONMENT: production
    ports:
      - "8080:8080"
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
      emqx:
        condition: service_started
    networks:
      - cmon_network
    deploy:
      replicas: 2
      resources:
        limits:
          cpus: '1'
          memory: 1G
    healthcheck:
      test: ["CMD", "wget", "--spider", "http://localhost:8080/health/ready"]
      interval: 30s
      timeout: 5s
      retries: 3

  nginx:
    image: nginx:alpine
    container_name: cmon-nginx
    restart: always
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
      - ./ssl:/etc/nginx/ssl:ro
    depends_on:
      - api
    networks:
      - cmon_network

  prometheus:
    image: prom/prometheus:latest
    container_name: cmon-prometheus
    restart: always
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
      - prom_data:/prometheus
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
    ports:
      - "9090:9090"
    networks:
      - cmon_network

  grafana:
    image: grafana/grafana:latest
    container_name: cmon-grafana
    restart: always
    environment:
      GF_SECURITY_ADMIN_PASSWORD: ${GRAFANA_PASSWORD:-admin}
    volumes:
      - grafana_data:/var/lib/grafana
    ports:
      - "3000:3000"
    depends_on:
      - prometheus
    networks:
      - cmon_network

networks:
  cmon_network:
    driver: bridge

volumes:
  pg_data:
  redis_data:
  emqx_data:
  prom_data:
  grafana_data:
```

---

## บทสรุปส่งท้าย (Conclusion)

เอกสารชุด **CMON IoT Solution – Go Backend API** นี้ได้ออกแบบและอธิบายอย่างครบถ้วนตั้งแต่:

1. **เล่ม 1: ภาคทฤษฎี** – ปูพื้นฐานความจำเป็นในการใช้ระบบตรวจสอบ Data Center, นิยามศัพท์เทคโนโลยีสำคัญ (MQTT, JWT, WebSocket, Rule Engine, Scheduler, Clean Architecture) และวิเคราะห์ความต้องการพร้อมกรณีศึกษา

2. **เล่ม 2: สถาปัตยกรรมโครงสร้างระบบ** – อธิบายโครงสร้างโฟลเดอร์แบบ 3-layer, การจัดการ Config/Logging/Middleware, Repository Pattern + Transaction + GORM, และ Authentication (JWT + Refresh Token + RBAC) พร้อมตัวอย่างโค้ดที่รันได้จริง

3. **เล่ม 3: การพัฒนาเชิงปฏิบัติ** – สร้าง MQTT subscriber สำหรับรับข้อมูลเซนเซอร์, Rule Engine และ Notifiers (Email, Line, WebSocket), WebSocket Hub + Dashboard Real-time (Chart.js), Scheduler + Distributed Lock + Report Generator

4. **เล่ม 4: การปรับใช้และการบำรุงรักษา** – Docker multi-stage, Docker Compose (dev/prod), Graceful shutdown, Health checks, Prometheus + Grafana, Backup/Restore, Horizontal scaling guidelines, พร้อมเทมเพลต Task List, Checklist, Timeline สำหรับการดำเนินโครงการจริง

**สิ่งที่ได้รับจากเอกสารชุดนี้:**
- ความเข้าใจเชิงลึกในการออกแบบระบบ Go API สำหรับ IoT Monitoring
- โค้ดตัวอย่างที่นำไปปรับใช้ได้ทันที (copy-paste แล้วรัน)
- แนวทางปฏิบัติที่ดีที่สุด (best practices) สำหรับ production
- เทมเพลตสำหรับบริหารโครงการ (Task list, Checklist, Timeline)

**ข้อเสนอแนะในการนำไปใช้ต่อ:**
1. หากเป็นโปรเจกต์ขนาดเล็ก เริ่มจากเล่ม 2 และ 3 แล้ว deploy ด้วย docker-compose (เล่ม 4)
2. หากมีทีมหลายคน ควรใช้ Git flow และ CI/CD ตามที่แนะนำ
3. สำหรับการขยายระบบไปยังหลาย Data Center ให้ใช้ MQTT bridge และ Redis replication
4. ควรทำ load testing ก่อน production เสมอ (ใช้ k6 หรือ wrk)

**การสนับสนุนเพิ่มเติม:**
- หากต้องการให้จัดทำเป็นไฟล์ PDF แต่ละเล่ม สามารถใช้ pandoc หรือเครื่องมือ markdown to PDF
- หากต้องการปรับแต่งให้เข้ากับองค์กร (โลโก้, ชื่อโครงการ) สามารถแก้ไขได้ตามต้องการ
- สำหรับคำถามหรือปัญหาขณะ implement สามารถเปิด issue ใน repository หรือปรึกษาชุมชน Go และ MQTT

---

**ผู้จัดทำ:** Solution Architect Team  
**วันที่แล้วเสร็จ:** 4 เมษายน 2026  
**ลิขสิทธิ์:** สามารถนำไปใช้และปรับแต่งได้ตามต้องการ เพื่อประโยชน์ในการพัฒนา IoT Monitoring Systems

---

*จบเอกสารทั้ง 4 เล่ม และภาคผนวก*