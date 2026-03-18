<img src="https://r2cdn.perplexity.ai/pplx-full-logo-primary-dark%402x.png" style="height:64px;margin-right:32px"/>

## การเชื่อมต่อ Grafana กับ Elasticsearch สำหรับ Dashboard Logs

การเชื่อมต่อ Grafana กับ Elasticsearch ช่วยให้สามารถแสดง logs และ metrics ในที่เดียวกัน สร้าง unified observability platform[^1]

## 1. เพิ่ม Elasticsearch Data Source ใน Grafana

### วิธีที่ 1: ผ่าน UI

**ขั้นตอนการเพิ่ม Data Source**

1. เข้า Grafana → Configuration → Data Sources
2. คลิก "Add data source"
3. เลือก "Elasticsearch"
4. กรอกข้อมูลตามด้านล่าง

### วิธีที่ 2: ผ่าน Configuration File

**infrastructure/grafana/datasources/elasticsearch.yaml**

```yaml
apiVersion: 1

datasources:
  - name: Elasticsearch-Logs
    type: elasticsearch
    access: proxy
    url: http://elasticsearch.logging.svc.cluster.local:9200
    jsonData:
      # Elasticsearch version
      esVersion: "8.0.0"
      
      # Time field name
      timeField: "@timestamp"
      
      # Index pattern
      index: "car-service-*"
      interval: Daily
      
      # Log configuration
      logMessageField: "message"
      logLevelField: "log_level"
      
      # Maximum concurrent shard requests
      maxConcurrentShardRequests: 5
      
      # Timeout in seconds
      timeoutSeconds: 30
      
      # Include frozen indices
      includeFrozen: false
      
    # Basic auth (if enabled)
    # basicAuth: true
    # basicAuthUser: elastic
    # secureJsonData:
    #   basicAuthPassword: your-password
    
    # TLS/SSL configuration (if needed)
    # jsonData:
    #   tlsSkipVerify: false
    # secureJsonData:
    #   tlsCACert: |
    #     -----BEGIN CERTIFICATE-----
    #     ...
    #     -----END CERTIFICATE-----

  - name: Elasticsearch-Metrics
    type: elasticsearch
    access: proxy
    url: http://elasticsearch.logging.svc.cluster.local:9200
    jsonData:
      esVersion: "8.0.0"
      timeField: "@timestamp"
      index: "metricbeat-*"
      interval: Daily
```


### วิธีที่ 3: ผ่าน Kubernetes ConfigMap

**infrastructure/kubernetes/monitoring/grafana-datasource.yaml**

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: grafana-datasource-elasticsearch
  namespace: monitoring
  labels:
    grafana_datasource: "1"
data:
  elasticsearch-datasource.yaml: |
    apiVersion: 1
    datasources:
      - name: Elasticsearch-Logs
        type: elasticsearch
        access: proxy
        url: http://elasticsearch.logging.svc.cluster.local:9200
        isDefault: false
        editable: true
        jsonData:
          esVersion: "8.0.0"
          timeField: "@timestamp"
          index: "car-service-*"
          interval: Daily
          logMessageField: "message"
          logLevelField: "log_level"
          maxConcurrentShardRequests: 5
```

**Apply ConfigMap**

```bash
kubectl apply -f infrastructure/kubernetes/monitoring/grafana-datasource.yaml

# Restart Grafana to load new datasource
kubectl rollout restart deployment/prometheus-grafana -n monitoring
```


## 2. สร้าง Log Dashboard

### Dashboard สำหรับ Application Logs

**infrastructure/grafana-dashboards/application-logs-dashboard.json**

```json
{
  "dashboard": {
    "title": "Car Service Application Logs",
    "timezone": "browser",
    "panels": [
      {
        "id": 1,
        "title": "Log Volume Over Time",
        "type": "graph",
        "gridPos": {
          "x": 0,
          "y": 0,
          "w": 24,
          "h": 8
        },
        "targets": [
          {
            "datasource": "Elasticsearch-Logs",
            "metrics": [
              {
                "type": "count",
                "id": "1"
              }
            ],
            "bucketAggs": [
              {
                "type": "date_histogram",
                "field": "@timestamp",
                "id": "2",
                "settings": {
                  "interval": "auto",
                  "min_doc_count": "0",
                  "trimEdges": "0"
                }
              },
              {
                "type": "terms",
                "field": "log_level.keyword",
                "id": "3",
                "settings": {
                  "size": "10",
                  "order": "desc",
                  "orderBy": "_count"
                }
              }
            ],
            "query": "kubernetes.namespace:production",
            "timeField": "@timestamp"
          }
        ],
        "yaxes": [
          {
            "label": "Log Count",
            "format": "short"
          }
        ],
        "legend": {
          "show": true,
          "values": true,
          "current": true
        }
      },
      {
        "id": 2,
        "title": "Error Logs",
        "type": "logs",
        "gridPos": {
          "x": 0,
          "y": 8,
          "w": 12,
          "h": 10
        },
        "targets": [
          {
            "datasource": "Elasticsearch-Logs",
            "query": "log_level:ERROR AND kubernetes.namespace:production",
            "timeField": "@timestamp",
            "metrics": [
              {
                "type": "logs"
              }
            ]
          }
        ],
        "options": {
          "showTime": true,
          "showLabels": false,
          "wrapLogMessage": true,
          "sortOrder": "Descending"
        }
      },
      {
        "id": 3,
        "title": "Warning Logs",
        "type": "logs",
        "gridPos": {
          "x": 12,
          "y": 8,
          "w": 12,
          "h": 10
        },
        "targets": [
          {
            "datasource": "Elasticsearch-Logs",
            "query": "log_level:WARN AND kubernetes.namespace:production",
            "timeField": "@timestamp",
            "metrics": [
              {
                "type": "logs"
              }
            ]
          }
        ]
      },
      {
        "id": 4,
        "title": "Logs by Service",
        "type": "piechart",
        "gridPos": {
          "x": 0,
          "y": 18,
          "w": 8,
          "h": 8
        },
        "targets": [
          {
            "datasource": "Elasticsearch-Logs",
            "metrics": [
              {
                "type": "count",
                "id": "1"
              }
            ],
            "bucketAggs": [
              {
                "type": "terms",
                "field": "kubernetes.container.name.keyword",
                "id": "2",
                "settings": {
                  "size": "10",
                  "order": "desc",
                  "orderBy": "_count"
                }
              }
            ],
            "query": "kubernetes.namespace:production"
          }
        ]
      },
      {
        "id": 5,
        "title": "Top Error Messages",
        "type": "table",
        "gridPos": {
          "x": 8,
          "y": 18,
          "w": 16,
          "h": 8
        },
        "targets": [
          {
            "datasource": "Elasticsearch-Logs",
            "metrics": [
              {
                "type": "count",
                "id": "1"
              }
            ],
            "bucketAggs": [
              {
                "type": "terms",
                "field": "log_message.keyword",
                "id": "2",
                "settings": {
                  "size": "10",
                  "order": "desc",
                  "orderBy": "_count"
                }
              }
            ],
            "query": "log_level:ERROR AND kubernetes.namespace:production"
          }
        ],
        "transformations": [
          {
            "id": "organize",
            "options": {
              "excludeByName": {},
              "indexByName": {},
              "renameByName": {
                "log_message.keyword": "Error Message",
                "Count": "Occurrences"
              }
            }
          }
        ]
      }
    ],
    "templating": {
      "list": [
        {
          "name": "namespace",
          "type": "query",
          "datasource": "Elasticsearch-Logs",
          "query": {
            "query": "*",
            "field": "kubernetes.namespace.keyword"
          },
          "multi": false,
          "includeAll": false,
          "current": {
            "value": "production"
          }
        },
        {
          "name": "service",
          "type": "query",
          "datasource": "Elasticsearch-Logs",
          "query": {
            "query": "kubernetes.namespace:$namespace",
            "field": "kubernetes.container.name.keyword"
          },
          "multi": true,
          "includeAll": true
        },
        {
          "name": "log_level",
          "type": "custom",
          "multi": true,
          "includeAll": true,
          "options": [
            { "text": "ERROR", "value": "ERROR" },
            { "text": "WARN", "value": "WARN" },
            { "text": "INFO", "value": "INFO" },
            { "text": "DEBUG", "value": "DEBUG" }
          ]
        }
      ]
    },
    "time": {
      "from": "now-1h",
      "to": "now"
    },
    "refresh": "30s"
  }
}
```


## 3. ตัวอย่าง Query Patterns

### Basic Queries

```lucene
# All logs from production namespace
kubernetes.namespace:production

# Error logs only
log_level:ERROR

# Logs from specific service
kubernetes.container.name:nestjs-api

# Logs containing specific text
log_message:*booking*

# Combined query
log_level:ERROR AND kubernetes.container.name:nestjs-api AND log_message:*payment*

# Time range query (handled by Grafana UI)
@timestamp:[now-1h TO now]

# Exclude system logs
NOT kubernetes.namespace:(kube-system OR logging OR monitoring)

# Multiple services
kubernetes.container.name:(nestjs-api OR spring-api)

# Regular expression
log_message:/.*timeout.*/

# Field exists
_exists_:error_code

# Field does not exist
NOT _exists_:error_code

# Numeric range
response_time:[200 TO 500]
```


### Advanced Queries

```lucene
# HTTP 5xx errors with high response time
log_level:ERROR AND http_status:[500 TO 599] AND response_time:>1000

# Database query errors
log_message:*database* AND log_level:ERROR

# Specific user activity
user_id:"user-123" AND (log_message:*login* OR log_message:*logout*)

# API endpoint errors
kubernetes.container.name:nestjs-api AND log_message:*/api/v1/bookings* AND log_level:ERROR

# Slow queries
log_message:*query* AND execution_time:>5000

# Payment failures
kubernetes.container.name:spring-api AND log_message:*payment* AND log_level:ERROR
```


## 4. Dashboard สำหรับ Business Metrics

**infrastructure/grafana-dashboards/business-logs-dashboard.json**

```json
{
  "dashboard": {
    "title": "Car Service Business Metrics from Logs",
    "panels": [
      {
        "id": 1,
        "title": "Booking Events",
        "type": "graph",
        "targets": [
          {
            "datasource": "Elasticsearch-Logs",
            "metrics": [
              {
                "type": "count",
                "id": "1"
              }
            ],
            "bucketAggs": [
              {
                "type": "date_histogram",
                "field": "@timestamp",
                "id": "2"
              },
              {
                "type": "terms",
                "field": "event_type.keyword",
                "id": "3"
              }
            ],
            "query": "log_message:*BOOKING* AND (event_type:CREATED OR event_type:CONFIRMED OR event_type:CANCELLED)"
          }
        ]
      },
      {
        "id": 2,
        "title": "Payment Transactions",
        "type": "stat",
        "targets": [
          {
            "datasource": "Elasticsearch-Logs",
            "metrics": [
              {
                "type": "count",
                "id": "1"
              }
            ],
            "bucketAggs": [
              {
                "type": "terms",
                "field": "payment_status.keyword",
                "id": "2"
              }
            ],
            "query": "log_message:*PAYMENT* AND kubernetes.container.name:spring-api"
          }
        ]
      },
      {
        "id": 3,
        "title": "User Login Events",
        "type": "logs",
        "targets": [
          {
            "datasource": "Elasticsearch-Logs",
            "query": "log_message:*LOGIN* OR log_message:*AUTHENTICATION*",
            "timeField": "@timestamp"
          }
        ]
      },
      {
        "id": 4,
        "title": "Failed Transactions",
        "type": "table",
        "targets": [
          {
            "datasource": "Elasticsearch-Logs",
            "metrics": [
              {
                "type": "count",
                "id": "1"
              }
            ],
            "bucketAggs": [
              {
                "type": "terms",
                "field": "error_code.keyword",
                "id": "2",
                "settings": {
                  "size": "10",
                  "order": "desc"
                }
              },
              {
                "type": "terms",
                "field": "transaction_type.keyword",
                "id": "3"
              }
            ],
            "query": "log_level:ERROR AND (log_message:*PAYMENT* OR log_message:*BOOKING*)"
          }
        ]
      }
    ]
  }
}
```


## 5. การ Import Dashboard ผ่าน Kubernetes

**infrastructure/kubernetes/monitoring/grafana-dashboard-configmap.yaml**

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: grafana-dashboard-logs
  namespace: monitoring
  labels:
    grafana_dashboard: "1"
data:
  application-logs.json: |
    {
      "dashboard": {
        "title": "Application Logs",
        ...
      }
    }
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: grafana-dashboard-provider
  namespace: monitoring
data:
  dashboards.yaml: |
    apiVersion: 1
    providers:
      - name: 'Logs'
        orgId: 1
        folder: 'Logs'
        type: file
        disableDeletion: false
        updateIntervalSeconds: 30
        allowUiUpdates: true
        options:
          path: /var/lib/grafana/dashboards
```

**Update Grafana Deployment**

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: prometheus-grafana
  namespace: monitoring
spec:
  template:
    spec:
      containers:
      - name: grafana
        volumeMounts:
        - name: dashboard-logs
          mountPath: /var/lib/grafana/dashboards/logs
        - name: dashboard-provider
          mountPath: /etc/grafana/provisioning/dashboards
      volumes:
      - name: dashboard-logs
        configMap:
          name: grafana-dashboard-logs
      - name: dashboard-provider
        configMap:
          name: grafana-dashboard-provider
```


## 6. การใช้งาน Explore Mode

### ตัวอย่างการ Query ใน Explore

```bash
# เปิด Grafana Explore
# 1. Click "Explore" icon (compass) ใน sidebar
# 2. เลือก "Elasticsearch-Logs" data source
# 3. ใส่ query และ select time range
```

**Example Explore Queries:**

```lucene
# Quick error investigation
log_level:ERROR AND kubernetes.namespace:production

# Trace specific request
request_id:"req-123456"

# Performance issues
log_message:*slow* OR execution_time:>3000

# Security events
log_message:*authentication* OR log_message:*authorization*
```


## 7. Alert Rules จาก Elasticsearch Logs

**infrastructure/kubernetes/monitoring/elasticsearch-alert-rules.yaml**

```yaml
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: elasticsearch-log-alerts
  namespace: monitoring
spec:
  groups:
  - name: elasticsearch-logs
    interval: 1m
    rules:
    # Alert on error rate from logs
    - alert: HighErrorRateInLogs
      expr: |
        sum(rate(elasticsearch_query_count{log_level="ERROR"}[5m])) 
        / 
        sum(rate(elasticsearch_query_count[5m])) > 0.1
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "High error rate in application logs"
        description: "Error log rate is {{ $value | humanizePercentage }}"

    # Alert on specific error patterns
    - alert: PaymentFailureSpike
      expr: |
        sum(rate(elasticsearch_query_count{
          log_message=~".*payment.*failed.*"
        }[5m])) > 5
      for: 2m
      labels:
        severity: critical
      annotations:
        summary: "Payment failure spike detected"
        description: "{{ $value }} payment failures per second"
```


## 8. Testing Connection

**Test script: scripts/test-grafana-elasticsearch.sh**

```bash
#!/bin/bash

# Test Elasticsearch connectivity
echo "Testing Elasticsearch connection..."
ES_URL="http://elasticsearch.logging.svc.cluster.local:9200"

# Check cluster health
curl -X GET "$ES_URL/_cluster/health?pretty"

# Check indices
curl -X GET "$ES_URL/_cat/indices/car-service-*?v"

# Test query
curl -X GET "$ES_URL/car-service-*/_search?pretty" -H 'Content-Type: application/json' -d'
{
  "query": {
    "match": {
      "log_level": "ERROR"
    }
  },
  "size": 5
}
'

# Test Grafana API
echo ""
echo "Testing Grafana API..."
GRAFANA_URL="http://prometheus-grafana.monitoring.svc.cluster.local"
GRAFANA_USER="admin"
GRAFANA_PASS=$(kubectl get secret -n monitoring prometheus-grafana -o jsonpath='{.data.admin-password}' | base64 -d)

# List data sources
curl -u "$GRAFANA_USER:$GRAFANA_PASS" "$GRAFANA_URL/api/datasources"

echo ""
echo "✅ Tests complete"
```


## Summary: Grafana + Elasticsearch Integration

| Feature | Configuration | Purpose |
| :-- | :-- | :-- |
| **Data Source** | URL: `elasticsearch.logging:9200` | Connect to ES |
| **Index Pattern** | `car-service-*` | Query logs |
| **Time Field** | `@timestamp` | Time series data |
| **Log Fields** | `message`, `log_level` | Display logs |
| **Dashboard Types** | Logs, Graph, Table, Pie | Visualize data |
| **Query Language** | Lucene | Search logs |
| **Refresh Rate** | 30s - 5m | Real-time updates |

การเชื่อมต่อ Grafana กับ Elasticsearch ทำให้สามารถดู metrics จาก Prometheus และ logs จาก Elasticsearch ในที่เดียวกัน สร้าง unified observability platform ที่ทรงพลัง[^1]

<div align="center">⁂</div>

[^1]: https://last9.io/blog/how-to-connect-elk-stack-with-grafana/

