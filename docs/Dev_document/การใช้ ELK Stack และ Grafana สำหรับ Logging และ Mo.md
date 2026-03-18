<img src="https://r2cdn.perplexity.ai/pplx-full-logo-primary-dark%402x.png" style="height:64px;margin-right:32px"/>

## การใช้ ELK Stack และ Grafana สำหรับ Logging และ Monitoring

การผสมผสาน ELK Stack (Elasticsearch, Logstash, Kibana) สำหรับ centralized logging และ Grafana + Prometheus สำหรับ monitoring จะช่วยให้มี observability ที่สมบูรณ์สำหรับระบบศูนย์บริการรถยนต์[^1][^2]

## สถาปัตยกรรม Observability

```
┌─────────────────────────────────────────────────────────────┐
│                  APPLICATION LAYER                           │
│  NestJS API  │  Spring Boot API  │  React Frontend          │
└──────┬───────────────┬──────────────────┬───────────────────┘
       │               │                  │
       │ Logs          │ Logs             │ Logs
       ▼               ▼                  ▼
┌──────────────────────────────────────────────────────────────┐
│              FILEBEAT (Log Collector)                        │
│              Running as DaemonSet                            │
└──────┬───────────────────────────────────────────────────────┘
       │
       │ Forward Logs
       ▼
┌──────────────────────────────────────────────────────────────┐
│              LOGSTASH (Log Processing)                        │
│              Parse, Filter, Transform                         │
└──────┬───────────────────────────────────────────────────────┘
       │
       │ Store Logs
       ▼
┌──────────────────────────────────────────────────────────────┐
│              ELASTICSEARCH (Log Storage)                      │
│              Index & Search Engine                            │
└──────┬────────────────────────┬──────────────────────────────┘
       │                        │
       │ Query                  │ Query
       ▼                        ▼
┌──────────────┐         ┌─────────────────┐
│    KIBANA    │         │    GRAFANA      │
│ Log Analysis │         │ Unified Dashboards│
└──────────────┘         └────────┬─────────┘
                                  │
                         Also connects to:
                                  │
                         ┌────────▼─────────┐
                         │   PROMETHEUS     │
                         │  (Metrics Store) │
                         └──────────────────┘
```


## 1. ELK Stack Deployment

### Namespace และ Storage

**infrastructure/kubernetes/monitoring/namespace.yaml**

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: logging
  labels:
    name: logging
---
apiVersion: v1
kind: Namespace
metadata:
  name: monitoring
  labels:
    name: monitoring
```


### Elasticsearch Deployment

**infrastructure/kubernetes/elk/elasticsearch.yaml**

```yaml
apiVersion: v1
kind: Service
metadata:
  name: elasticsearch
  namespace: logging
  labels:
    app: elasticsearch
spec:
  selector:
    app: elasticsearch
  clusterIP: None
  ports:
  - port: 9200
    name: rest
  - port: 9300
    name: inter-node
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: elasticsearch
  namespace: logging
spec:
  serviceName: elasticsearch
  replicas: 3
  selector:
    matchLabels:
      app: elasticsearch
  template:
    metadata:
      labels:
        app: elasticsearch
    spec:
      containers:
      - name: elasticsearch
        image: docker.elastic.co/elasticsearch/elasticsearch:8.11.0
        resources:
          limits:
            memory: 4Gi
            cpu: 2000m
          requests:
            memory: 2Gi
            cpu: 1000m
        ports:
        - containerPort: 9200
          name: rest
          protocol: TCP
        - containerPort: 9300
          name: inter-node
          protocol: TCP
        volumeMounts:
        - name: data
          mountPath: /usr/share/elasticsearch/data
        env:
        - name: cluster.name
          value: k8s-logs
        - name: node.name
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        - name: discovery.seed_hosts
          value: "elasticsearch-0.elasticsearch,elasticsearch-1.elasticsearch,elasticsearch-2.elasticsearch"
        - name: cluster.initial_master_nodes
          value: "elasticsearch-0,elasticsearch-1,elasticsearch-2"
        - name: ES_JAVA_OPTS
          value: "-Xms2g -Xmx2g"
        - name: xpack.security.enabled
          value: "false"
      initContainers:
      - name: fix-permissions
        image: busybox
        command: ["sh", "-c", "chown -R 1000:1000 /usr/share/elasticsearch/data"]
        securityContext:
          privileged: true
        volumeMounts:
        - name: data
          mountPath: /usr/share/elasticsearch/data
      - name: increase-vm-max-map
        image: busybox
        command: ["sysctl", "-w", "vm.max_map_count=262144"]
        securityContext:
          privileged: true
  volumeClaimTemplates:
  - metadata:
      name: data
      labels:
        app: elasticsearch
    spec:
      accessModes: [ "ReadWriteOnce" ]
      storageClassName: gp2  # AWS EBS, adjust for your provider
      resources:
        requests:
          storage: 100Gi
```


### Logstash Deployment

**infrastructure/kubernetes/elk/logstash.yaml**

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: logstash-config
  namespace: logging
data:
  logstash.yml: |
    http.host: "0.0.0.0"
    xpack.monitoring.elasticsearch.hosts: [ "http://elasticsearch:9200" ]
  
  logstash.conf: |
    input {
      beats {
        port => 5044
      }
    }
    
    filter {
      # Parse JSON logs
      if [message] =~ /^\{.*\}$/ {
        json {
          source => "message"
        }
      }
      
      # Parse NestJS logs
      if [kubernetes][container][name] == "nestjs-api" {
        grok {
          match => { "message" => "\[%{DATA:log_level}\] %{NUMBER:pid} - %{TIMESTAMP_ISO8601:timestamp} %{GREEDYDATA:log_message}" }
        }
      }
      
      # Parse Spring Boot logs
      if [kubernetes][container][name] == "spring-api" {
        grok {
          match => { "message" => "%{TIMESTAMP_ISO8601:timestamp} %{LOGLEVEL:log_level} %{NUMBER:pid} --- \[%{DATA:thread}\] %{DATA:class} : %{GREEDYDATA:log_message}" }
        }
      }
      
      # Add custom fields
      mutate {
        add_field => {
          "environment" => "${ENVIRONMENT:production}"
          "application" => "car-service"
        }
      }
      
      # Parse timestamps
      date {
        match => [ "timestamp", "ISO8601" ]
        target => "@timestamp"
      }
      
      # Remove unnecessary fields
      mutate {
        remove_field => [ "host", "agent", "ecs", "input" ]
      }
    }
    
    output {
      elasticsearch {
        hosts => ["http://elasticsearch:9200"]
        index => "car-service-%{[kubernetes][namespace]}-%{+YYYY.MM.dd}"
        manage_template => true
        template_name => "car-service"
        template_overwrite => true
      }
      
      # Debug output (remove in production)
      # stdout { codec => rubydebug }
    }
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: logstash
  namespace: logging
spec:
  replicas: 2
  selector:
    matchLabels:
      app: logstash
  template:
    metadata:
      labels:
        app: logstash
    spec:
      containers:
      - name: logstash
        image: docker.elastic.co/logstash/logstash:8.11.0
        ports:
        - containerPort: 5044
          name: beats
        volumeMounts:
        - name: config
          mountPath: /usr/share/logstash/config/logstash.yml
          subPath: logstash.yml
        - name: pipeline
          mountPath: /usr/share/logstash/pipeline/logstash.conf
          subPath: logstash.conf
        env:
        - name: LS_JAVA_OPTS
          value: "-Xmx1g -Xms1g"
        - name: ENVIRONMENT
          value: "production"
        resources:
          limits:
            memory: 2Gi
            cpu: 1000m
          requests:
            memory: 1Gi
            cpu: 500m
      volumes:
      - name: config
        configMap:
          name: logstash-config
          items:
          - key: logstash.yml
            path: logstash.yml
      - name: pipeline
        configMap:
          name: logstash-config
          items:
          - key: logstash.conf
            path: logstash.conf
---
apiVersion: v1
kind: Service
metadata:
  name: logstash
  namespace: logging
spec:
  selector:
    app: logstash
  ports:
  - port: 5044
    targetPort: 5044
    name: beats
```


### Filebeat DaemonSet

**infrastructure/kubernetes/elk/filebeat.yaml**

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: filebeat-config
  namespace: logging
data:
  filebeat.yml: |
    filebeat.inputs:
    - type: container
      paths:
        - /var/log/containers/*.log
      processors:
      - add_kubernetes_metadata:
          host: ${NODE_NAME}
          matchers:
          - logs_path:
              logs_path: "/var/log/containers/"
      - drop_event:
          when:
            or:
            - equals:
                kubernetes.namespace: "kube-system"
            - equals:
                kubernetes.namespace: "logging"
            - equals:
                kubernetes.namespace: "monitoring"
    
    output.logstash:
      hosts: ["logstash:5044"]
      loadbalance: true
    
    logging.level: info
    logging.to_files: false
    logging.to_syslog: false
---
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: filebeat
  namespace: logging
  labels:
    app: filebeat
spec:
  selector:
    matchLabels:
      app: filebeat
  template:
    metadata:
      labels:
        app: filebeat
    spec:
      serviceAccountName: filebeat
      terminationGracePeriodSeconds: 30
      hostNetwork: true
      dnsPolicy: ClusterFirstWithHostNet
      containers:
      - name: filebeat
        image: docker.elastic.co/beats/filebeat:8.11.0
        args: [
          "-c", "/etc/filebeat.yml",
          "-e",
        ]
        env:
        - name: NODE_NAME
          valueFrom:
            fieldRef:
              fieldPath: spec.nodeName
        securityContext:
          runAsUser: 0
        resources:
          limits:
            memory: 200Mi
            cpu: 100m
          requests:
            memory: 100Mi
            cpu: 50m
        volumeMounts:
        - name: config
          mountPath: /etc/filebeat.yml
          readOnly: true
          subPath: filebeat.yml
        - name: data
          mountPath: /usr/share/filebeat/data
        - name: varlibdockercontainers
          mountPath: /var/lib/docker/containers
          readOnly: true
        - name: varlog
          mountPath: /var/log
          readOnly: true
      volumes:
      - name: config
        configMap:
          defaultMode: 0640
          name: filebeat-config
      - name: varlibdockercontainers
        hostPath:
          path: /var/lib/docker/containers
      - name: varlog
        hostPath:
          path: /var/log
      - name: data
        hostPath:
          path: /var/lib/filebeat-data
          type: DirectoryOrCreate
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: filebeat
rules:
- apiGroups: [""]
  resources:
  - namespaces
  - pods
  - nodes
  verbs:
  - get
  - watch
  - list
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: filebeat
  namespace: logging
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: filebeat
subjects:
- kind: ServiceAccount
  name: filebeat
  namespace: logging
roleRef:
  kind: ClusterRole
  name: filebeat
  apiGroup: rbac.authorization.k8s.io
```


### Kibana Deployment

**infrastructure/kubernetes/elk/kibana.yaml**

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kibana
  namespace: logging
spec:
  replicas: 1
  selector:
    matchLabels:
      app: kibana
  template:
    metadata:
      labels:
        app: kibana
    spec:
      containers:
      - name: kibana
        image: docker.elastic.co/kibana/kibana:8.11.0
        resources:
          limits:
            memory: 2Gi
            cpu: 1000m
          requests:
            memory: 1Gi
            cpu: 500m
        env:
        - name: ELASTICSEARCH_HOSTS
          value: "http://elasticsearch:9200"
        - name: SERVER_NAME
          value: "kibana"
        - name: SERVER_BASEPATH
          value: "/kibana"
        ports:
        - containerPort: 5601
---
apiVersion: v1
kind: Service
metadata:
  name: kibana
  namespace: logging
spec:
  selector:
    app: kibana
  ports:
  - port: 5601
    targetPort: 5601
  type: LoadBalancer
```


## 2. Prometheus + Grafana Monitoring Stack

### Prometheus Deployment with Helm

**Deploy Prometheus using Prometheus Operator**

```bash
# Add Helm repository
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

# Install Prometheus Stack (includes Grafana)
helm install prometheus prometheus-community/kube-prometheus-stack \
  --namespace monitoring \
  --create-namespace \
  --set prometheus.prometheusSpec.retention=30d \
  --set prometheus.prometheusSpec.storageSpec.volumeClaimTemplate.spec.resources.requests.storage=100Gi \
  --set grafana.adminPassword=admin123 \
  --set grafana.persistence.enabled=true \
  --set grafana.persistence.size=10Gi
```


### Custom ServiceMonitor for Applications

**infrastructure/kubernetes/monitoring/servicemonitor.yaml**

```yaml
apiVersion: v1
kind: Service
metadata:
  name: nestjs-api-metrics
  namespace: production
  labels:
    app: nestjs-api
spec:
  selector:
    app: nestjs-api
  ports:
  - name: metrics
    port: 9090
    targetPort: 9090
---
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: nestjs-api-monitor
  namespace: monitoring
  labels:
    release: prometheus
spec:
  selector:
    matchLabels:
      app: nestjs-api
  namespaceSelector:
    matchNames:
    - production
  endpoints:
  - port: metrics
    interval: 30s
    path: /metrics
---
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: spring-api-monitor
  namespace: monitoring
  labels:
    release: prometheus
spec:
  selector:
    matchLabels:
      app: spring-api
  namespaceSelector:
    matchNames:
    - production
  endpoints:
  - port: metrics
    interval: 30s
    path: /actuator/prometheus
```


### Application Metrics Configuration

**NestJS with Prometheus**

**backend/nestjs-api/src/metrics/metrics.module.ts**

```typescript
import { Module } from '@nestjs/common';
import { PrometheusModule } from '@willsoto/nestjs-prometheus';
import { MetricsController } from './metrics.controller';

@Module({
  imports: [
    PrometheusModule.register({
      path: '/metrics',
      defaultMetrics: {
        enabled: true,
      },
    }),
  ],
  controllers: [MetricsController],
})
export class MetricsModule {}
```

**Custom Metrics**

```typescript
import { Injectable } from '@nestjs/common';
import { Counter, Histogram, makeCounterProvider, makeHistogramProvider } from '@willsoto/nestjs-prometheus';

@Injectable()
export class MetricsService {
  constructor(
    private readonly bookingCounter: Counter<string>,
    private readonly requestDuration: Histogram<string>,
  ) {}

  incrementBookings(status: string) {
    this.bookingCounter.inc({ status });
  }

  recordRequestDuration(method: string, path: string, duration: number) {
    this.requestDuration.observe({ method, path }, duration);
  }
}

export const bookingCounterProvider = makeCounterProvider({
  name: 'booking_total',
  help: 'Total number of bookings',
  labelNames: ['status'],
});

export const requestDurationProvider = makeHistogramProvider({
  name: 'http_request_duration_seconds',
  help: 'HTTP request duration in seconds',
  labelNames: ['method', 'path'],
  buckets: [0.1, 0.5, 1, 2, 5],
});
```

**Spring Boot with Micrometer**

**pom.xml**

```xml
<dependency>
    <groupId>io.micrometer</groupId>
    <artifactId>micrometer-registry-prometheus</artifactId>
</dependency>
<dependency>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-actuator</artifactId>
</dependency>
```

**application.yml**

```yaml
management:
  endpoints:
    web:
      exposure:
        include: health,info,prometheus,metrics
  metrics:
    export:
      prometheus:
        enabled: true
    tags:
      application: car-service-spring
      environment: production
```


## 3. Grafana Dashboards

### Connect Elasticsearch to Grafana

**Add Elasticsearch Data Source in Grafana**

```json
{
  "name": "Elasticsearch-Logs",
  "type": "elasticsearch",
  "access": "proxy",
  "url": "http://elasticsearch.logging.svc.cluster.local:9200",
  "database": "car-service-*",
  "jsonData": {
    "timeField": "@timestamp",
    "esVersion": "8.0.0",
    "logMessageField": "message",
    "logLevelField": "log_level"
  }
}
```


### Application Overview Dashboard

**infrastructure/grafana-dashboards/car-service-overview.json**

```json
{
  "dashboard": {
    "title": "Car Service System Overview",
    "panels": [
      {
        "title": "Request Rate",
        "targets": [
          {
            "expr": "sum(rate(http_requests_total[5m])) by (service)",
            "legendFormat": "{{service}}"
          }
        ],
        "type": "graph"
      },
      {
        "title": "Error Rate",
        "targets": [
          {
            "expr": "sum(rate(http_requests_total{status=~\"5..\"}[5m])) / sum(rate(http_requests_total[5m])) * 100",
            "legendFormat": "Error %"
          }
        ],
        "type": "graph"
      },
      {
        "title": "Response Time (P95)",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le, service))",
            "legendFormat": "{{service}} P95"
          }
        ],
        "type": "graph"
      },
      {
        "title": "Active Bookings",
        "targets": [
          {
            "expr": "booking_total{status=\"confirmed\"}",
            "legendFormat": "Active"
          }
        ],
        "type": "stat"
      },
      {
        "title": "Database Connections",
        "targets": [
          {
            "expr": "pg_stat_activity_count",
            "legendFormat": "{{datname}}"
          }
        ],
        "type": "graph"
      },
      {
        "title": "Redis Hit Rate",
        "targets": [
          {
            "expr": "redis_keyspace_hits_total / (redis_keyspace_hits_total + redis_keyspace_misses_total) * 100",
            "legendFormat": "Hit Rate %"
          }
        ],
        "type": "gauge"
      }
    ]
  }
}
```


### Logs Dashboard

**Query logs from Elasticsearch in Grafana**

```json
{
  "title": "Application Logs",
  "panels": [
    {
      "title": "Error Logs",
      "targets": [
        {
          "datasource": "Elasticsearch-Logs",
          "query": "log_level:ERROR",
          "timeField": "@timestamp"
        }
      ],
      "type": "logs"
    },
    {
      "title": "Booking Events",
      "targets": [
        {
          "datasource": "Elasticsearch-Logs",
          "query": "kubernetes.container.name:nestjs-api AND log_message:*booking*",
          "timeField": "@timestamp"
        }
      ],
      "type": "logs"
    }
  ]
}
```


## 4. Alerting Configuration

### Prometheus Alert Rules

**infrastructure/kubernetes/monitoring/prometheus-rules.yaml**

```yaml
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: car-service-alerts
  namespace: monitoring
  labels:
    release: prometheus
spec:
  groups:
  - name: car-service
    interval: 30s
    rules:
    # High Error Rate
    - alert: HighErrorRate
      expr: |
        (
          sum(rate(http_requests_total{status=~"5.."}[5m]))
          /
          sum(rate(http_requests_total[5m]))
        ) > 0.05
      for: 5m
      labels:
        severity: critical
      annotations:
        summary: "High error rate detected"
        description: "Error rate is {{ $value | humanizePercentage }} (threshold: 5%)"

    # High Response Time
    - alert: HighResponseTime
      expr: |
        histogram_quantile(0.95,
          sum(rate(http_request_duration_seconds_bucket[5m])) by (le)
        ) > 2
      for: 10m
      labels:
        severity: warning
      annotations:
        summary: "High response time"
        description: "P95 latency is {{ $value }}s (threshold: 2s)"

    # Database Connection Pool Exhaustion
    - alert: DatabaseConnectionPoolHigh
      expr: |
        pg_stat_activity_count > 80
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "Database connection pool usage high"
        description: "{{ $value }} connections active (threshold: 80)"

    # Pod Restart
    - alert: PodRestarting
      expr: |
        rate(kube_pod_container_status_restarts_total[15m]) > 0
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "Pod {{ $labels.pod }} is restarting"
        description: "Pod has restarted {{ $value }} times in the last 15 minutes"

    # Disk Space
    - alert: DiskSpaceRunningOut
      expr: |
        (
          node_filesystem_avail_bytes{mountpoint="/"}
          /
          node_filesystem_size_bytes{mountpoint="/"}
        ) < 0.1
      for: 5m
      labels:
        severity: critical
      annotations:
        summary: "Disk space running out"
        description: "Only {{ $value | humanizePercentage }} disk space remaining"
```


### Alertmanager Configuration

**infrastructure/kubernetes/monitoring/alertmanager-config.yaml**

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: alertmanager-prometheus-kube-prometheus-alertmanager
  namespace: monitoring
stringData:
  alertmanager.yaml: |
    global:
      resolve_timeout: 5m
      slack_api_url: 'YOUR_SLACK_WEBHOOK_URL'

    route:
      group_by: ['alertname', 'cluster', 'service']
      group_wait: 10s
      group_interval: 10s
      repeat_interval: 12h
      receiver: 'slack-notifications'
      routes:
      - match:
          severity: critical
        receiver: 'slack-critical'
        continue: true
      - match:
          severity: warning
        receiver: 'slack-warnings'

    receivers:
    - name: 'slack-notifications'
      slack_configs:
      - channel: '#devops-alerts'
        title: 'Car Service Alert'
        text: '{{ range .Alerts }}{{ .Annotations.description }}{{ end }}'

    - name: 'slack-critical'
      slack_configs:
      - channel: '#critical-alerts'
        title: '🚨 CRITICAL ALERT'
        text: '{{ range .Alerts }}{{ .Annotations.description }}{{ end }}'
        send_resolved: true

    - name: 'slack-warnings'
      slack_configs:
      - channel: '#warnings'
        title: '⚠️ Warning'
        text: '{{ range .Alerts }}{{ .Annotations.description }}{{ end }}'
```


## 5. Deployment Script

**scripts/deploy-observability.sh**

```bash
#!/bin/bash

set -e

echo "🚀 Deploying Observability Stack"

# Create namespaces
echo "📦 Creating namespaces..."
kubectl apply -f infrastructure/kubernetes/monitoring/namespace.yaml

# Deploy ELK Stack
echo "📊 Deploying Elasticsearch..."
kubectl apply -f infrastructure/kubernetes/elk/elasticsearch.yaml
kubectl rollout status statefulset/elasticsearch -n logging --timeout=10m

echo "📊 Deploying Logstash..."
kubectl apply -f infrastructure/kubernetes/elk/logstash.yaml
kubectl rollout status deployment/logstash -n logging --timeout=5m

echo "📊 Deploying Filebeat..."
kubectl apply -f infrastructure/kubernetes/elk/filebeat.yaml
kubectl rollout status daemonset/filebeat -n logging --timeout=5m

echo "📊 Deploying Kibana..."
kubectl apply -f infrastructure/kubernetes/elk/kibana.yaml
kubectl rollout status deployment/kibana -n logging --timeout=5m

# Deploy Prometheus + Grafana
echo "📈 Deploying Prometheus Stack..."
helm upgrade --install prometheus prometheus-community/kube-prometheus-stack \
  --namespace monitoring \
  --values infrastructure/helm/prometheus-values.yaml

# Apply custom ServiceMonitors
echo "📈 Applying ServiceMonitors..."
kubectl apply -f infrastructure/kubernetes/monitoring/servicemonitor.yaml

# Apply AlertRules
echo "🔔 Applying Alert Rules..."
kubectl apply -f infrastructure/kubernetes/monitoring/prometheus-rules.yaml

# Get access information
echo ""
echo "✅ Deployment Complete!"
echo ""
echo "📊 Kibana: http://$(kubectl get svc kibana -n logging -o jsonpath='{.status.loadBalancer.ingress[^0].hostname}'):5601"
echo "📈 Grafana: http://$(kubectl get svc prometheus-grafana -n monitoring -o jsonpath='{.status.loadBalancer.ingress[^0].hostname}')"
echo ""
echo "🔐 Grafana Credentials:"
echo "   Username: admin"
echo "   Password: $(kubectl get secret -n monitoring prometheus-grafana -o jsonpath='{.data.admin-password}' | base64 -d)"
```


## Architecture Benefits Summary

| Component | Purpose | Retention | Query Performance |
| :-- | :-- | :-- | :-- |
| **Elasticsearch** | Log storage \& search | 30 days | Fast text search |
| **Kibana** | Log exploration | N/A | Deep log analysis |
| **Prometheus** | Metrics storage | 30 days | Time-series queries |
| **Grafana** | Unified dashboards | N/A | Real-time visualization |

การผสมผสาน ELK และ Grafana ให้ observability ที่สมบูรณ์ โดย ELK จัดการ logs และ Grafana แสดงทั้ง metrics และ logs ในที่เดียว[^2][^1]
<span style="display:none">[^10][^3][^4][^5][^6][^7][^8][^9]</span>

<div align="center">⁂</div>

[^1]: https://last9.io/blog/how-to-connect-elk-stack-with-grafana/

[^2]: https://www.refontelearning.com/blog/monitoring-logging-prometheus-grafana-elk-stack-loki

[^3]: https://www.youtube.com/watch?v=OLHpnPqV3-k

[^4]: https://www.elastic.co/docs/deploy-manage/deploy/cloud-on-k8s

[^5]: https://aws.plainenglish.io/one-minute-elk-stack-on-kubernetes-full-logging-setup-with-a-single-script-ba92aecb4379

[^6]: https://logz.io/learn/complete-guide-elk-stack/

[^7]: https://www.youtube.com/watch?v=5IlFTWsgMzU

[^8]: https://www.site24x7.com/learn/prometheus-grafana-monitoring.html

[^9]: https://dev.to/akhil_mittal/monitoring-logging-setup-of-application-deployed-in-eks-45fg

[^10]: https://navendu.me/posts/introduction-to-monitoring-microservices/

