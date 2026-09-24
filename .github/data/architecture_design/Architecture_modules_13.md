
---

# 🚀 PART 13 — DEPLOYMENT & DEVOPS

> **ขนาด**: ใหญ่ — แยก 8 ตอนย่อย
> **Part 13A**: Deployment Topology
> **Part 13B**: Docker & Docker Compose
> **Part 13C**: Kubernetes (Production)
> **Part 13D**: CI/CD Pipeline
> **Part 13E**: Configuration Management
> **Part 13F**: Observability Stack
> **Part 13G**: Backup & Disaster Recovery
> **Part 13H**: Runbooks

---

## 🅰️ PART 13A — DEPLOYMENT TOPOLOGY

### A.1 Environments

| Env | Purpose | Scale | Domain |
|:---|:---|:---|:---|
| **local** | Dev | docker-compose | localhost |
| **dev** | Integration | 1 replica | dev.icmon.local |
| **staging** | Pre-prod | 2 replicas | staging.icmongolang.io |
| **prod** | Production | 3+ replicas | api.icmongolang.io |

### A.2 Production Topology

```
                      ┌──────────────┐
                      │  Cloudflare  │  (CDN + DDoS + WAF)
                      └──────┬───────┘
                             │
                      ┌──────▼───────┐
                      │     NLB      │  (Network Load Balancer)
                      └──────┬───────┘
                             │
                  ┌──────────┼──────────┐
                  │          │          │
             ┌────▼───┐ ┌────▼───┐ ┌────▼───┐
             │ Ingress│ │ Ingress│ │ Ingress│
             │  NGINX │ │  NGINX │ │  NGINX │
             └────┬───┘ └────┬───┘ └────┬───┘
                  │          │          │
       ┌──────────┴──────────┴──────────┴──────────┐
       │            Kubernetes Cluster              │
       │                                            │
       │  ┌──────────┐  ┌──────────┐  ┌──────────┐ │
       │  │ API pod  │  │ API pod  │  │ API pod  │ │  (3 replicas)
       │  └──────────┘  └──────────┘  └──────────┘ │
       │                                            │
       │  ┌──────────┐  ┌──────────┐               │
       │  │ MQTT GW  │  │MQTT GW   │               │  (2 replicas)
       │  └──────────┘  └──────────┘               │
       │                                            │
       │  ┌──────────┐  ┌──────────┐  ┌──────────┐ │
       │  │ Worker   │  │ Worker   │  │ Worker   │ │  (HPA)
       │  └──────────┘  └──────────┘  └──────────┘ │
       │                                            │
       │  ┌──────────┐  ┌──────────┐               │
       │  │Scheduler │  │Scheduler │               │  (leader only)
       │  └──────────┘  └──────────┘               │
       │                                            │
       └────────────────────────────────────────────┘
                             │
       ┌─────────────────────┼─────────────────────┐
       │                     │                     │
  ┌────▼────┐          ┌─────▼─────┐         ┌────▼────┐
  │Postgres │          │   Redis   │         │  Kafka  │
  │(Primary)│          │ (Sentinel)│         │(3 nodes)│
  │+Replica │          │           │         │         │
  └─────────┘          └───────────┘         └─────────┘
       │
  ┌────▼────┐          ┌───────────┐         ┌──────────┐
  │InfluxDB │          │Elasticsear│         │   S3     │
  │         │          │   ch      │         │ (backup) │
  └─────────┘          └───────────┘         └──────────┘
```

---

## 🅱️ PART 13B — DOCKER & DOCKER COMPOSE

### B.1 Multi-Stage Dockerfile

```dockerfile
# Dockerfile
# ============================================================
# Stage 1: Builder
# ============================================================
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /build

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source
COPY . .

# Build args
ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_TIME=unknown

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
      -ldflags="-w -s \
        -X main.Version=${VERSION} \
        -X main.Commit=${COMMIT} \
        -X main.BuildTime=${BUILD_TIME}" \
      -o /build/bin/api ./cmd/api

# ============================================================
# Stage 2: Runtime
# ============================================================
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata curl && \
    addgroup -S app && adduser -S app -G app

WORKDIR /app

COPY --from=builder /build/bin/api /app/api
COPY --from=builder /build/migrations /app/migrations

RUN chown -R app:app /app
USER app

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

ENTRYPOINT ["/app/api"]
```

### B.2 `.dockerignore`

```
.git
.github
.vscode
.idea
*.md
docs
test
tmp
bin
coverage.out
coverage.html
.env
.env.*
!.env.example
docker-compose*.yml
Makefile
```

### B.3 docker-compose.dev.yml

```yaml
version: '3.9'

services:
  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: icmon
      POSTGRES_PASSWORD: icmon_dev
      POSTGRES_DB: icmongolang
      TZ: Asia/Bangkok
    ports: ["5432:5432"]
    volumes:
      - pgdata:/var/lib/postgresql/data
      - ./scripts/init-db.sh:/docker-entrypoint-initdb.d/init.sh
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U icmon"]
      interval: 5s
      timeout: 3s
      retries: 10

  redis:
    image: redis:7-alpine
    command: redis-server --appendonly yes --requirepass ""
    ports: ["6379:6379"]
    volumes: [redisdata:/data]
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s

  kafka:
    image: confluentinc/cp-kafka:7.5.0
    depends_on: [zookeeper]
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_LISTENERS: PLAINTEXT://0.0.0.0:9092
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:9092
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
      KAFKA_AUTO_CREATE_TOPICS_ENABLE: 'true'
    ports: ["9092:9092"]
    healthcheck:
      test: ["CMD", "kafka-broker-api-versions", "--bootstrap-server", "localhost:9092"]
      interval: 10s
      timeout: 10s

  zookeeper:
    image: confluentinc/cp-zookeeper:7.5.0
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181
    ports: ["2181:2181"]

  influxdb:
    image: influxdb:2.7-alpine
    environment:
      DOCKER_INFLUXDB_INIT_MODE: setup
      DOCKER_INFLUXDB_INIT_USERNAME: admin
      DOCKER_INFLUXDB_INIT_PASSWORD: adminpass123
      DOCKER_INFLUXDB_INIT_ORG: icmon
      DOCKER_INFLUXDB_INIT_BUCKET: iot_telemetry
      DOCKER_INFLUXDB_INIT_RETENTION: 30d
      DOCKER_INFLUXDB_INIT_ADMIN_TOKEN: dev-token-change-me
    ports: ["8086:8086"]
    volumes: [influxdata:/var/lib/influxdb2]

  mosquitto:
    image: eclipse-mosquitto:2
    ports: ["1883:1883", "9001:9001"]
    volumes:
      - ./deploy/mosquitto/mosquitto.conf:/mosquitto/config/mosquitto.conf
      - mosquittodata:/mosquitto/data

  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:8.11.0
    environment:
      discovery.type: single-node
      xpack.security.enabled: "false"
      ES_JAVA_OPTS: "-Xms512m -Xmx512m"
    ports: ["9200:9200"]
    volumes: [esdata:/usr/share/elasticsearch/data]

  prometheus:
    image: prom/prometheus:latest
    ports: ["9090:9090"]
    volumes:
      - ./deploy/prometheus/prometheus.yml:/etc/prometheus/prometheus.yml

  grafana:
    image: grafana/grafana:latest
    ports: ["3000:3000"]
    environment:
      GF_SECURITY_ADMIN_PASSWORD: admin
    volumes: [grafanadata:/var/lib/grafana]

  loki:
    image: grafana/loki:latest
    ports: ["3100:3100"]

  jaeger:
    image: jaegertracing/all-in-one:latest
    ports: ["16686:16686", "4318:4318"]
    environment:
      COLLECTOR_OTLP_ENABLED: "true"

volumes:
  pgdata:
  redisdata:
  influxdata:
  esdata:
  grafanadata:
  mosquittodata:
```

### B.4 Makefile

```makefile
.PHONY: docker-build
docker-build:
	docker build -t icmongolang/api:$(VERSION) .

.PHONY: docker-build-all
docker-build-all:
	docker build --target api -t icmongolang/api:$(VERSION) .
	docker build --target mqtt-gateway -t icmongolang/mqtt-gateway:$(VERSION) .
	docker build --target workers -t icmongolang/workers:$(VERSION) .

.PHONY: dev-up
dev-up:
	docker compose -f docker-compose.dev.yml up -d

.PHONY: dev-down
dev-down:
	docker compose -f docker-compose.dev.yml down

.PHONY: dev-logs
dev-logs:
	docker compose -f docker-compose.dev.yml logs -f --tail=100

.PHONY: dev-reset
dev-reset:
	docker compose -f docker-compose.dev.yml down -v
	docker compose -f docker-compose.dev.yml up -d
```

---

## 🅲 PART 13C — KUBERNETES

### C.1 Namespace + ConfigMap + Secret

```yaml
# k8s/base/namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: icmon

---
# k8s/base/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: icmon-config
  namespace: icmon
data:
  APP_ENV: "production"
  APP_PORT: "8080"
  LOG_LEVEL: "info"
  DB_HOST: "postgres.icmon.svc.cluster.local"
  DB_PORT: "5432"
  DB_NAME: "icmongolang"
  REDIS_ADDR: "redis.icmon.svc.cluster.local:6379"
  KAFKA_BROKERS: "kafka-0.kafka:9092,kafka-1.kafka:9092,kafka-2.kafka:9092"
  INFLUX_URL: "http://influxdb.icmon.svc:8086"
  INFLUX_ORG: "icmon"
  INFLUX_BUCKET: "iot_telemetry"
  ELASTICSEARCH_URL: "http://elasticsearch.icmon.svc:9200"
  MQTT_BROKER: "tcp://mosquitto.icmon.svc:1883"

---
# k8s/base/secret.yaml (ใช้ External Secrets หรือ SealedSecrets ใน prod)
apiVersion: v1
kind: Secret
metadata:
  name: icmon-secrets
  namespace: icmon
type: Opaque
stringData:
  DB_PASSWORD: "REPLACE_ME"
  JWT_SECRET: "REPLACE_ME"
  INFLUX_TOKEN: "REPLACE_ME"
```

### C.2 API Deployment

```yaml
# k8s/base/api-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: icmon-api
  namespace: icmon
  labels: { app: icmon-api }
spec:
  replicas: 3
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  selector:
    matchLabels: { app: icmon-api }
  template:
    metadata:
      labels: { app: icmon-api }
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "8080"
        prometheus.io/path: "/metrics"
    spec:
      serviceAccountName: icmon-api
      securityContext:
        runAsNonRoot: true
        runAsUser: 10001
        fsGroup: 10001
      containers:
        - name: api
          image: icmongolang/api:v1.0.0
          imagePullPolicy: IfNotPresent
          ports:
            - containerPort: 8080
              name: http
          envFrom:
            - configMapRef: { name: icmon-config }
            - secretRef: { name: icmon-secrets }
          env:
            - name: DB_USER
              value: icmon
          resources:
            requests: { cpu: "200m", memory: "256Mi" }
            limits:   { cpu: "1000m", memory: "1Gi" }
          livenessProbe:
            httpGet: { path: /health, port: http }
            initialDelaySeconds: 10
            periodSeconds: 10
          readinessProbe:
            httpGet: { path: /readyz, port: http }
            initialDelaySeconds: 5
            periodSeconds: 5
          startupProbe:
            httpGet: { path: /health, port: http }
            failureThreshold: 30
            periodSeconds: 2
          lifecycle:
            preStop:
              exec: { command: ["sh", "-c", "sleep 10"] }
          securityContext:
            allowPrivilegeEscalation: false
            readOnlyRootFilesystem: true
            capabilities: { drop: ["ALL"] }
          volumeMounts:
            - { name: tmp, mountPath: /tmp }
      volumes:
        - name: tmp
          emptyDir: {}
      terminationGracePeriodSeconds: 30

---
apiVersion: v1
kind: Service
metadata:
  name: icmon-api
  namespace: icmon
spec:
  selector: { app: icmon-api }
  ports:
    - name: http
      port: 80
      targetPort: 8080
  type: ClusterIP
```

### C.3 HPA (Horizontal Pod Autoscaler)

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: icmon-api-hpa
  namespace: icmon
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: icmon-api
  minReplicas: 3
  maxReplicas: 20
  metrics:
    - type: Resource
      resource:
        name: cpu
        target: { type: Utilization, averageUtilization: 70 }
    - type: Resource
      resource:
        name: memory
        target: { type: Utilization, averageUtilization: 80 }
    - type: Pods
      pods:
        metric: { name: http_requests_per_second }
        target: { type: AverageValue, averageValue: "1000" }
  behavior:
    scaleUp:
      stabilizationWindowSeconds: 30
      policies:
        - { type: Percent, value: 100, periodSeconds: 30 }
    scaleDown:
      stabilizationWindowSeconds: 300
      policies:
        - { type: Percent, value: 10, periodSeconds: 60 }
```

### C.4 Ingress

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: icmon-ingress
  namespace: icmon
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
    nginx.ingress.kubernetes.io/proxy-body-size: "10m"
    nginx.ingress.kubernetes.io/proxy-read-timeout: "3600"
    nginx.ingress.kubernetes.io/proxy-send-timeout: "3600"
    nginx.ingress.kubernetes.io/websocket-services: icmon-api
    nginx.ingress.kubernetes.io/rate-limit: "1000"
    nginx.ingress.kubernetes.io/rate-limit-window: "1m"
spec:
  ingressClassName: nginx
  tls:
    - hosts: [api.icmongolang.io]
      secretName: icmon-tls
  rules:
    - host: api.icmongolang.io
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: icmon-api
                port: { number: 80 }
```

### C.5 Worker Deployment (with KEDA)

```yaml
apiVersion: keda.sh/v1alpha1
kind: ScaledObject
metadata:
  name: telemetry-worker
  namespace: icmon
spec:
  scaleTargetRef:
    name: icmon-telemetry-worker
  minReplicaCount: 2
  maxReplicaCount: 30
  triggers:
    - type: kafka
      metadata:
        bootstrapServers: kafka.icmon:9092
        consumerGroup: telemetry-worker
        topic: icmon.device.telemetry.ingested
        lagThreshold: "1000"
```

---

## 🅳 PART 13D — CI/CD PIPELINE

### D.1 GitLab CI

```yaml
# .gitlab-ci.yml
stages:
  - test
  - build
  - deploy-dev
  - deploy-staging
  - deploy-prod

variables:
  DOCKER_REGISTRY: registry.gitlab.com/icmongolang
  IMAGE_TAG: $CI_COMMIT_SHORT_SHA

# ============================================================
# Test
# ============================================================
unit-test:
  stage: test
  image: golang:1.23-alpine
  script:
    - go mod download
    - go test -short -race -coverprofile=coverage.out ./...
    - go tool cover -func=coverage.out
  artifacts:
    reports: { coverage_report: { coverage_format: cobertura, path: coverage.xml } }

integration-test:
  stage: test
  image: golang:1.23-alpine
  services:
    - postgres:16-alpine
    - redis:7-alpine
  script:
    - go test -tags=integration -timeout=15m ./test/integration/...

# ============================================================
# Build
# ============================================================
build:
  stage: build
  image: docker:24
  services: [docker:24-dind]
  script:
    - docker login -u $CI_REGISTRY_USER -p $CI_REGISTRY_PASSWORD $CI_REGISTRY
    - docker build
        --build-arg VERSION=$CI_COMMIT_TAG
        --build-arg COMMIT=$CI_COMMIT_SHA
        -t $DOCKER_REGISTRY/api:$IMAGE_TAG
        -t $DOCKER_REGISTRY/api:latest .
    - docker push $DOCKER_REGISTRY/api:$IMAGE_TAG
    - docker push $DOCKER_REGISTRY/api:latest
  rules:
    - if: $CI_COMMIT_BRANCH == "main"

# ============================================================
# Deploy
# ============================================================
.deploy-template: &deploy
  image: bitnami/kubectl:latest
  script:
    - kubectl config use-context $KUBE_CONTEXT
    - kubectl -n icmon set image deployment/icmon-api api=$DOCKER_REGISTRY/api:$IMAGE_TAG
    - kubectl -n icmon rollout status deployment/icmon-api --timeout=5m

deploy-dev:
  <<: *deploy
  stage: deploy-dev
  variables: { KUBE_CONTEXT: dev }
  environment: { name: dev }
  rules: [{ if: $CI_COMMIT_BRANCH == "develop" }]

deploy-staging:
  <<: *deploy
  stage: deploy-staging
  variables: { KUBE_CONTEXT: staging }
  environment: { name: staging }
  when: manual
  rules: [{ if: $CI_COMMIT_BRANCH == "main" }]

deploy-prod:
  <<: *deploy
  stage: deploy-prod
  variables: { KUBE_CONTEXT: prod }
  environment: { name: production }
  when: manual
  rules: [{ if: $CI_COMMIT_TAG =~ /^v\d+\.\d+\.\d+$/ }]
```

### D.2 ArgoCD (GitOps)

```yaml
# argocd/applications/icmon-prod.yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: icmon-prod
  namespace: argocd
spec:
  project: icmon
  source:
    repoURL: https://github.com/icmongolang/deploy.git
    targetRevision: main
    path: overlays/prod
  destination:
    server: https://kubernetes.default.svc
    namespace: icmon
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
      allowEmpty: false
    syncOptions:
      - CreateNamespace=false
      - PrunePropagationPolicy=foreground
    retry:
      limit: 5
      backoff:
        duration: 5s
        factor: 2
        maxDuration: 3m
  revisionHistoryLimit: 10
```

### D.3 Release Strategy

| Strategy | ใช้เมื่อ | Rollback |
|:---|:---|:---|
| **Rolling Update** | default | kubectl rollout undo |
| **Blue-Green** | major release | switch service selector |
| **Canary** | risky change | ลด traffic กลับ |
| **Feature Flag** | เปิด/ปิดฟีเจอร์ | toggle flag |

---

## 🅴 PART 13E — CONFIGURATION MANAGEMENT

### E.1 12-Factor Config

```go
// internal/shared/config/config.go
package config

type Config struct {
    Env         string `env:"APP_ENV" envDefault:"development"`
    Port        string `env:"APP_PORT" envDefault:"8080"`
    LogLevel    string `env:"LOG_LEVEL" envDefault:"info"`

    DB          DBConfig
    Redis       RedisConfig
    Kafka       KafkaConfig
    Influx      InfluxConfig
    ES          ESConfig
    MQTT        MQTTConfig
    JWT         JWTConfig
    LLM         LLMConfig
    Payment     PaymentConfig
    Storage     StorageConfig
    RateLimit   RateLimitConfig
    CORS        CORSConfig
}

func Load() *Config {
    cfg := &Config{}
    if err := env.Parse(cfg); err != nil {
        log.Fatal(err)
    }
    validate(cfg)
    return cfg
}
```

### E.2 Environment Variable Convention

```env
# Core
APP_ENV=production
APP_PORT=8080
LOG_LEVEL=info
APP_VERSION=v1.0.0

# Database
DB_HOST=postgres.icmon.svc
DB_PORT=5432
DB_USER=icmon
DB_PASSWORD=${VAULT:prod_db_password}
DB_NAME=icmongolang
DB_MAX_OPEN=25
DB_MAX_IDLE=10
DB_CONN_TIMEOUT=5s

# Redis
REDIS_ADDR=redis.icmon.svc:6379
REDIS_PASSWORD=${VAULT:prod_redis_password}
REDIS_DB=0
REDIS_POOL_SIZE=100

# Kafka
KAFKA_BROKERS=kafka-0:9092,kafka-1:9092,kafka-2:9092
KAFKA_PREFIX=icmon
KAFKA_GROUP_ID=icmon-api

# MQTT
MQTT_BROKER=tcp://mosquitto.icmon.svc:1883
MQTT_CLIENT_ID=icmon-gateway
MQTT_USERNAME=${VAULT:mqtt_user}
MQTT_PASSWORD=${VAULT:mqtt_pass}
MQTT_USE_TLS=true

# InfluxDB
INFLUX_URL=http://influxdb.icmon.svc:8086
INFLUX_TOKEN=${VAULT:influx_token}
INFLUX_ORG=icmon
INFLUX_BUCKET=iot_telemetry

# Elasticsearch
ES_URL=http://elasticsearch.icmon.svc:9200
ES_PREFIX=icmon

# JWT
JWT_SECRET=${VAULT:jwt_secret}
JWT_ACCESS_TTL=15m
JWT_REFRESH_TTL=168h

# Rate Limit
RATE_LIMIT_GLOBAL=10000
RATE_LIMIT_PER_TENANT=1000
RATE_LIMIT_WINDOW=1m

# External
LLM_PROVIDER=ollama
LLM_API_URL=http://ollama.ai.svc:11434
LLM_MODEL=llama3
STRIPE_SECRET_KEY=${VAULT:stripe_key}
STRIPE_WEBHOOK_SECRET=${VAULT:stripe_webhook}
```

### E.3 Secret Management (Vault / AWS Secrets Manager)

```go
// internal/shared/config/vault.go
package config

import (
    vault "github.com/hashicorp/vault/api"
)

func LoadFromVault(path string) (map[string]interface{}, error) {
    cfg := vault.DefaultConfig()
    client, err := vault.NewClient(cfg)
    if err != nil {
        return nil, err
    }
    client.SetToken(os.Getenv("VAULT_TOKEN"))

    secret, err := client.Logical().Read(path)
    if err != nil {
        return nil, err
    }
    return secret.Data, nil
}
```

### E.4 Feature Flags

```go
// internal/shared/infrastructure/feature_flags/flags.go
package featureflags

type Client interface {
    IsEnabled(ctx context.Context, flag string, tenantID uuid.UUID) bool
    Value(ctx context.Context, flag string, tenantID uuid.UUID) interface{}
}

// Simple DB-backed implementation
type DBClient struct {
    repo repository.FeatureFlagRepository
    cache application.Cache
}

func (c *DBClient) IsEnabled(ctx context.Context, flag string, tenantID uuid.UUID) bool {
    // 1. Check tenant-specific override
    // 2. Check global flag
    // 3. Check percentage rollout
    // 4. Default: false
    return false
}

// Usage
if c.features.IsEnabled(ctx, "ai_insights", tenantID) {
    // enable AI insights
}
```

---

## 🅵 PART 13F — OBSERVABILITY STACK

### F.1 Prometheus Config

```yaml
# deploy/prometheus/prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: icmon-api
    kubernetes_sd_configs:
      - role: pod
        namespaces: { names: [icmon] }
    relabel_configs:
      - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_scrape]
        action: keep
        regex: "true"
      - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_port]
        action: replace
        target_label: __address__
        regex: (.+)
        replacement: ${1}:8080

  - job_name: kafka
    static_configs:
      - targets: [kafka-exporter:9308]

  - job_name: postgres
    static_configs:
      - targets: [postgres-exporter:9187]

  - job_name: redis
    static_configs:
      - targets: [redis-exporter:9121]
```

### F.2 Grafana Dashboards (JSON)

```json
// deploy/grafana/dashboards/api-overview.json
{
  "title": "icmon API Overview",
  "panels": [
    {
      "title": "Request Rate",
      "targets": [{ "expr": "sum(rate(http_requests_total[1m])) by (path)" }]
    },
    {
      "title": "p95 Latency",
      "targets": [{ "expr": "histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))" }]
    },
    {
      "title": "Error Rate",
      "targets": [{ "expr": "sum(rate(http_requests_total{status=~\"5..\"}[5m])) / sum(rate(http_requests_total[5m]))" }]
    },
    {
      "title": "Kafka Lag",
      "targets": [{ "expr": "kafka_consumer_group_lag" }]
    },
    {
      "title": "Device Online",
      "targets": [{ "expr": "sum(device_status{status=\"online\"})" }]
    }
  ]
}
```

### F.3 Alert Rules

```yaml
# deploy/prometheus/alerts.yml
groups:
  - name: icmon-alerts
    rules:
      - alert: HighErrorRate
        expr: sum(rate(http_requests_total{status=~"5.."}[5m])) / sum(rate(http_requests_total[5m])) > 0.05
        for: 5m
        labels: { severity: critical }
        annotations:
          summary: "Error rate > 5%"
          runbook: "https://wiki/runbooks/high-error"

      - alert: HighLatency
        expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 1
        for: 5m
        labels: { severity: warning }

      - alert: KafkaLagHigh
        expr: kafka_consumer_group_lag > 10000
        for: 10m
        labels: { severity: warning }

      - alert: DeviceOfflineRate
        expr: sum(device_status{status="offline"}) / sum(device_status) > 0.3
        for: 10m
        labels: { severity: warning }

      - alert: DatabaseConnectionsHigh
        expr: pg_stat_activity_count > 80
        for: 5m
```

### F.4 Structured Logging (Loki)

```go
// internal/shared/infrastructure/logger/zerolog.go
package logger

import (
    "github.com/rs/zerolog"
    "icmongolang/internal/shared/application"
)

type zerologAdapter struct {
    logger zerolog.Logger
}

func New(level, env string) application.Logger {
    zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
    var zl zerolog.Logger
    if env == "production" {
        zl = zerolog.New(os.Stdout).With().Timestamp().Logger()
    } else {
        zl = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
    }
    lvl, _ := zerolog.ParseLevel(level)
    zl = zl.Level(lvl)
    return &zerologAdapter{logger: zl}
}

func (l *zerologAdapter) Info(msg string, fields ...application.Field) {
    ev := l.logger.Info()
    for _, f := range fields {
        ev = ev.Interface(f.Key, f.Value)
    }
    ev.Msg(msg)
}
```

### F.5 Distributed Tracing

```go
// HTTP middleware inject trace
func TracingMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx := otel.GetTextMapPropagator().Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))
        tracer := otel.Tracer("icmon-api")
        ctx, span := tracer.Start(ctx, c.FullPath())
        defer span.End()
        span.SetAttributes(
            attribute.String("http.method", c.Request.Method),
            attribute.String("http.path", c.Request.URL.Path),
            attribute.String("http.request_id", c.GetString("request_id")),
        )
        c.Request = c.Request.WithContext(ctx)
        c.Next()
        span.SetAttributes(attribute.Int("http.status_code", c.Writer.Status()))
    }
}
```

---

## 🅶 PART 13G — BACKUP & DISASTER RECOVERY

### G.1 Backup Strategy

| Component | Method | Frequency | Retention | RPO | RTO |
|:---|:---|:---|:---|:---|:---|
| **PostgreSQL** | WAL + pg_dump | Continuous + daily | 30 days | 5 min | 30 min |
| **Redis** | RDB snapshots | hourly | 7 days | 1 hr | 5 min |
| **InfluxDB** | Snapshot | daily | 30 days | 24 hr | 1 hr |
| **Kafka** | Topic replication | continuous | 7 days | 0 | 5 min |
| **Elasticsearch** | Snapshot | daily | 30 days | 24 hr | 2 hr |
| **S3** | Versioning + replication | continuous | 365 days | 0 | 5 min |

### G.2 PostgreSQL Backup Script

```bash
#!/bin/bash
# scripts/backup-postgres.sh
set -euo pipefail

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR=/backups/postgres
S3_BUCKET=s3://icmon-backups/postgres
RETENTION_DAYS=30

mkdir -p "$BACKUP_DIR"

# Backup
pg_dump -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" \
    --format=custom --compress=9 \
    --file="$BACKUP_DIR/icmon_${DATE}.dump"

# Upload to S3
aws s3 cp "$BACKUP_DIR/icmon_${DATE}.dump" "$S3_BUCKET/" --storage-class STANDARD_IA

# Cleanup local
find "$BACKUP_DIR" -name "icmon_*.dump" -mtime +$RETENTION_DAYS -delete

# Verify
pg_restore --list "$BACKUP_DIR/icmon_${DATE}.dump" > /dev/null
echo "✓ Backup completed: icmon_${DATE}.dump"
```

### G.3 Restore Procedure

```bash
#!/bin/bash
# scripts/restore-postgres.sh
set -euo pipefail

BACKUP_FILE=$1
TARGET_DB=${2:-icmongolang_restore}

echo "Restoring $BACKUP_FILE to $TARGET_DB"

# 1. Create DB
psql -h "$DB_HOST" -U "$DB_USER" -c "CREATE DATABASE $TARGET_DB;"

# 2. Restore
pg_restore -h "$DB_HOST" -U "$DB_USER" -d "$TARGET_DB" \
    --no-owner --no-acl --jobs=4 "$BACKUP_FILE"

# 3. Verify
psql -h "$DB_HOST" -U "$DB_USER" -d "$TARGET_DB" -c "\dt"
echo "✓ Restore completed"
```

### G.4 DR Drill Checklist

```
[ ] ทดสอบ restore จาก backup (รายเดือน)
[ ] ทดสอบ failover PostgreSQL primary → replica
[ ] ทดสอบ Kafka broker restart
[ ] ทดสอบ Redis sentinel failover
[ ] ทดสอบ multi-region failover (รายปี)
[ ] วัด RTO/RPO จริง vs target
[ ] อัปเดต runbook
```

---

## 🅷 PART 13H — RUNBOOKS

### H.1 Common Incidents

```markdown
## Runbook: High Error Rate

### Symptoms
- Alert `HighErrorRate` firing
- Grafana: error rate > 5%

### Diagnosis
1. ดู log: `kubectl logs -n icmon -l app=icmon-api --tail=1000 | grep ERROR`
2. ดู trace: Jaeger → filter by status code 5xx
3. ตรวจ dependency status: `/readyz`
4. ดู recent deploy: `kubectl rollout history deployment/icmon-api`

### Mitigation
1. ถ้า deploy ใหม่ → rollback: `kubectl rollout undo deployment/icmon-api`
2. ถ้า DB slow → scale up Postgres / add index
3. ถ้า external service down → enable circuit breaker
4. แจ้ง stakeholders

### Post-mortem
- สร้าง incident report ภายใน 48 ชม.
- Root cause analysis (5 Whys)
- Action items + owner + due date

---

## Runbook: Kafka Consumer Lag

### Symptoms
- Alert `KafkaLagHigh`
- Telemetry processing ช้า

### Diagnosis
```
kafka-consumer-groups --bootstrap-server kafka:9092 \
    --describe --group telemetry-worker
```

### Mitigation
1. Scale workers: `kubectl scale deployment icmon-telemetry-worker --replicas=10`
2. ตรวจ slow consumer: profiler / trace
3. ถ้า partition 太少 → เพิ่ม partitions
4. ถ้า consumer stuck → restart pod

---

## Runbook: Device Offline Spike

### Symptoms
- Alert `DeviceOfflineRate`
- > 30% devices offline

### Diagnosis
1. ตรวจ MQTT broker: `systemctl status mosquitto` หรือ EMQX dashboard
2. ตรวจ network: `ping iot-gateway`
3. ตรวจ device firmware version distribution
4. ตรวจ MQTT topic: `mosquitto_sub -t 'iot/+/status' -v`

### Mitigation
1. ถ้า broker down → restart / failover
2. ถ้า network issue → escalate to infra team
3. ถ้า firmware bug → rollback firmware
```

### H.2 Maintenance Windows

| Activity | Window | Duration | Downtime |
|:---|:---|:---|:---|
| **DB migration** | Sun 02:00-04:00 ICT | 30 min | none (online) |
| **K8s upgrade** | Sun 03:00-05:00 ICT | 1 hr | none |
| **InfluxDB compaction** | daily 01:00 | varies | none |
| **Certificate renewal** | auto (cert-manager) | - | none |

### H.3 Rollback Procedures

```bash
# 1. Rollback deployment
kubectl -n icmon rollout undo deployment/icmon-api
kubectl -n icmon rollout status deployment/icmon-api

# 2. Rollback DB migration (golang-migrate)
migrate -path ./migrations -database "$DB_DSN" down 1

# 3. Rollback via ArgoCD
argocd app rollback icmon-prod

# 4. Feature flag kill switch
curl -X POST https://api.icmongolang.io/admin/flags/ai_insights/disable \
    -H "Authorization: Bearer $ADMIN_TOKEN"
```

---

## 🎯 PART 13 — SUMMARY

| Deliverable | Count |
|---|:-:|
| **Dockerfiles** | 3 (api, worker, mqtt-gateway) |
| **Docker Compose** | 2 (dev, seed) |
| **K8s Manifests** | 20+ |
| **CI/CD Pipelines** | 3 (GitHub, GitLab, ArgoCD) |
| **Grafana Dashboards** | 5+ |
| **Prometheus Alerts** | 15+ |
| **Runbooks** | 10+ |

---
