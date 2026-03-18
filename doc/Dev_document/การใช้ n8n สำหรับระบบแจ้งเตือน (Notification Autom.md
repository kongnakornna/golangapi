

## การใช้ n8n สำหรับระบบแจ้งเตือน (Notification Automation)

n8n เป็น workflow automation tool แบบ open-source ที่ช่วยให้สร้าง automated workflows สำหรับแจ้งเตือนและ integrate กับระบบต่างๆ ได้อย่างง่ายดาย

## 1. ภาพรวม n8n Integration Architecture

```
┌─────────────────────────────────────────────────────────────┐
│              APPLICATION EVENTS                              │
│  Booking Created │ Payment Success │ Repair Complete        │
└─────────┬────────────────┬────────────────┬─────────────────┘
          │                │                │
          ▼                ▼                ▼
┌─────────────────────────────────────────────────────────────┐
│                    KAFKA TOPICS                              │
│  booking-events  │  payment-events  │  repair-events        │
└─────────┬────────────────┬────────────────┬─────────────────┘
          │                │                │
          │ Webhooks       │                │
          ▼                ▼                ▼
┌─────────────────────────────────────────────────────────────┐
│                      n8n WORKFLOWS                           │
│  • Notification Router                                       │
│  • Customer Notifications                                    │
│  • Admin Alerts                                             │
│  • Technician Assignments                                   │
└─────────┬────────────┬──────────┬──────────┬────────────────┘
          │            │          │          │
          ▼            ▼          ▼          ▼
┌─────────────┬────────────┬──────────┬──────────────────────┐
│   Slack     │   Email    │   SMS    │  LINE/WhatsApp       │
│  (DevOps)   │ (Customer) │ (Urgent) │  (Customer)          │
└─────────────┴────────────┴──────────┴──────────────────────┘
```


## 2. การติดตั้ง n8n บน Kubernetes

### n8n Deployment

**infrastructure/kubernetes/n8n/deployment.yaml**

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: n8n
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: n8n-data
  namespace: n8n
spec:
  accessModes:
    - ReadWriteOnce
  storageClassName: gp2
  resources:
    requests:
      storage: 10Gi
---
apiVersion: v1
kind: Secret
metadata:
  name: n8n-secrets
  namespace: n8n
type: Opaque
stringData:
  N8N_BASIC_AUTH_USER: admin
  N8N_BASIC_AUTH_PASSWORD: changeme123
  N8N_ENCRYPTION_KEY: "your-encryption-key-32-chars"
  WEBHOOK_URL: "https://n8n.car-service.com"
  
  # Email SMTP
  SMTP_HOST: smtp.gmail.com
  SMTP_PORT: "587"
  SMTP_USER: notifications@car-service.com
  SMTP_PASSWORD: your-app-password
  
  # Slack
  SLACK_WEBHOOK_URL: https://hooks.slack.com/services/YOUR/WEBHOOK/URL
  
  # LINE Notify
  LINE_NOTIFY_TOKEN: your-line-notify-token
  
  # Twilio SMS
  TWILIO_ACCOUNT_SID: your-account-sid
  TWILIO_AUTH_TOKEN: your-auth-token
  TWILIO_PHONE_NUMBER: "+1234567890"
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: n8n
  namespace: n8n
spec:
  replicas: 1
  selector:
    matchLabels:
      app: n8n
  template:
    metadata:
      labels:
        app: n8n
    spec:
      containers:
      - name: n8n
        image: n8nio/n8n:latest
        ports:
        - containerPort: 5678
          name: http
        env:
        - name: N8N_HOST
          value: "n8n.car-service.com"
        - name: N8N_PORT
          value: "5678"
        - name: N8N_PROTOCOL
          value: "https"
        - name: NODE_ENV
          value: "production"
        - name: EXECUTIONS_PROCESS
          value: "main"
        - name: EXECUTIONS_MODE
          value: "regular"
        - name: N8N_BASIC_AUTH_ACTIVE
          value: "true"
        - name: GENERIC_TIMEZONE
          value: "Asia/Bangkok"
        - name: N8N_DIAGNOSTICS_ENABLED
          value: "false"
        - name: N8N_LOG_LEVEL
          value: "info"
        
        envFrom:
        - secretRef:
            name: n8n-secrets
        
        volumeMounts:
        - name: data
          mountPath: /home/node/.n8n
        
        resources:
          requests:
            memory: "512Mi"
            cpu: "250m"
          limits:
            memory: "1Gi"
            cpu: "500m"
        
        livenessProbe:
          httpGet:
            path: /healthz
            port: 5678
          initialDelaySeconds: 30
          periodSeconds: 10
        
        readinessProbe:
          httpGet:
            path: /healthz
            port: 5678
          initialDelaySeconds: 10
          periodSeconds: 5
      
      volumes:
      - name: data
        persistentVolumeClaim:
          claimName: n8n-data
---
apiVersion: v1
kind: Service
metadata:
  name: n8n
  namespace: n8n
spec:
  selector:
    app: n8n
  ports:
  - port: 5678
    targetPort: 5678
    name: http
  type: ClusterIP
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: n8n-ingress
  namespace: n8n
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
    nginx.ingress.kubernetes.io/proxy-body-size: "50m"
spec:
  ingressClassName: nginx
  tls:
  - hosts:
    - n8n.car-service.com
    secretName: n8n-tls
  rules:
  - host: n8n.car-service.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: n8n
            port:
              number: 5678
```


### Deploy n8n

```bash
# Deploy n8n
kubectl apply -f infrastructure/kubernetes/n8n/deployment.yaml

# Check status
kubectl get pods -n n8n

# Get n8n URL
echo "n8n URL: https://$(kubectl get ingress n8n-ingress -n n8n -o jsonpath='{.spec.rules[0].host}')"

# Get credentials
echo "Username: $(kubectl get secret n8n-secrets -n n8n -o jsonpath='{.data.N8N_BASIC_AUTH_USER}' | base64 -d)"
echo "Password: $(kubectl get secret n8n-secrets -n n8n -o jsonpath='{.data.N8N_BASIC_AUTH_PASSWORD}' | base64 -d)"
```


## 3. n8n Workflows สำหรับแจ้งเตือน

### Workflow 1: Booking Notification

**จุดประสงค์:** แจ้งเตือนลูกค้าเมื่อมีการจอง, ยืนยัน, หรือยกเลิก

**n8n Workflow JSON:**

```json
{
  "name": "Booking Notification Workflow",
  "nodes": [
    {
      "parameters": {
        "httpMethod": "POST",
        "path": "booking-webhook",
        "responseMode": "onReceived",
        "options": {}
      },
      "name": "Webhook - Booking Event",
      "type": "n8n-nodes-base.webhook",
      "position": [250, 300]
    },
    {
      "parameters": {
        "conditions": {
          "string": [
            {
              "value1": "={{$json.eventType}}",
              "operation": "equals",
              "value2": "BOOKING_CREATED"
            }
          ]
        }
      },
      "name": "If Booking Created",
      "type": "n8n-nodes-base.if",
      "position": [450, 200]
    },
    {
      "parameters": {
        "conditions": {
          "string": [
            {
              "value1": "={{$json.eventType}}",
              "operation": "equals",
              "value2": "BOOKING_CONFIRMED"
            }
          ]
        }
      },
      "name": "If Booking Confirmed",
      "position": [450, 300]
    },
    {
      "parameters": {
        "conditions": {
          "string": [
            {
              "value1": "={{$json.eventType}}",
              "operation": "equals",
              "value2": "BOOKING_CANCELLED"
            }
          ]
        }
      },
      "name": "If Booking Cancelled",
      "position": [450, 400]
    },
    {
      "parameters": {
        "authentication": "genericCredentialType",
        "genericAuthType": "httpHeaderAuth",
        "url": "https://api.car-service.com/api/v1/users/={{$json.customerId}}",
        "options": {}
      },
      "name": "Get Customer Details",
      "type": "n8n-nodes-base.httpRequest",
      "position": [650, 250]
    },
    {
      "parameters": {
        "fromEmail": "notifications@car-service.com",
        "toEmail": "={{$node['Get Customer Details'].json.email}}",
        "subject": "การจองของคุณได้รับการยืนยันแล้ว",
        "text": "เรียน คุณ {{$node['Get Customer Details'].json.name}}\n\nการจองของคุณได้รับการยืนยันแล้ว\n\nรายละเอียด:\n- หมายเลขจอง: {{$json.bookingId}}\n- วันที่: {{$json.bookingDate}}\n- ประเภทบริการ: {{$json.serviceType}}\n\nขอบคุณที่ใช้บริการ\nศูนย์บริการรถยนต์",
        "options": {
          "allowUnauthorizedCerts": false
        }
      },
      "name": "Send Email - Confirmed",
      "type": "n8n-nodes-base.emailSend",
      "credentials": {
        "smtp": {
          "id": "1",
          "name": "SMTP Account"
        }
      },
      "position": [850, 200]
    },
    {
      "parameters": {
        "message": "=📅 *การจองใหม่*\n\nลูกค้า: {{$node['Get Customer Details'].json.name}}\nหมายเลขจอง: {{$json.bookingId}}\nวันที่: {{$json.bookingDate}}\nบริการ: {{$json.serviceType}}",
        "otherOptions": {}
      },
      "name": "Send LINE Notify",
      "type": "n8n-nodes-base.line",
      "credentials": {
        "lineNotifyOAuth2Api": {
          "id": "2",
          "name": "LINE Notify"
        }
      },
      "position": [850, 300]
    },
    {
      "parameters": {
        "channel": "#bookings",
        "text": "=:calendar: *New Booking*\n\nCustomer: {{$node['Get Customer Details'].json.name}}\nBooking ID: {{$json.bookingId}}\nDate: {{$json.bookingDate}}\nService: {{$json.serviceType}}",
        "otherOptions": {
          "username": "Booking Bot"
        }
      },
      "name": "Send Slack Notification",
      "type": "n8n-nodes-base.slack",
      "credentials": {
        "slackApi": {
          "id": "3",
          "name": "Slack"
        }
      },
      "position": [850, 400]
    }
  ],
  "connections": {
    "Webhook - Booking Event": {
      "main": [
        [
          {
            "node": "If Booking Created",
            "type": "main",
            "index": 0
          },
          {
            "node": "If Booking Confirmed",
            "type": "main",
            "index": 0
          },
          {
            "node": "If Booking Cancelled",
            "type": "main",
            "index": 0
          }
        ]
      ]
    },
    "If Booking Created": {
      "main": [
        [
          {
            "node": "Get Customer Details",
            "type": "main",
            "index": 0
          }
        ]
      ]
    },
    "If Booking Confirmed": {
      "main": [
        [
          {
            "node": "Get Customer Details",
            "type": "main",
            "index": 0
          }
        ]
      ]
    },
    "Get Customer Details": {
      "main": [
        [
          {
            "node": "Send Email - Confirmed",
            "type": "main",
            "index": 0
          },
          {
            "node": "Send LINE Notify",
            "type": "main",
            "index": 0
          },
          {
            "node": "Send Slack Notification",
            "type": "main",
            "index": 0
          }
        ]
      ]
    }
  },
  "active": true
}
```


### Workflow 2: System Alert Notification

**จุดประสงค์:** แจ้งเตือนทีม DevOps เมื่อมีปัญหาระบบ

```json
{
  "name": "System Alert Workflow",
  "nodes": [
    {
      "parameters": {
        "httpMethod": "POST",
        "path": "alertmanager-webhook",
        "responseMode": "onReceived"
      },
      "name": "Webhook - Alertmanager",
      "type": "n8n-nodes-base.webhook",
      "position": [250, 300]
    },
    {
      "parameters": {
        "conditions": {
          "string": [
            {
              "value1": "={{$json.alerts[0].labels.severity}}",
              "operation": "equals",
              "value2": "critical"
            }
          ]
        }
      },
      "name": "If Critical Alert",
      "type": "n8n-nodes-base.if",
      "position": [450, 250]
    },
    {
      "parameters": {
        "conditions": {
          "string": [
            {
              "value1": "={{$json.alerts[0].labels.severity}}",
              "operation": "equals",
              "value2": "warning"
            }
          ]
        }
      },
      "name": "If Warning Alert",
      "position": [450, 400]
    },
    {
      "parameters": {
        "channel": "#critical-alerts",
        "text": "=🚨 *CRITICAL ALERT*\n\n*Alert:* {{$json.alerts[0].labels.alertname}}\n*Description:* {{$json.alerts[0].annotations.description}}\n*Service:* {{$json.alerts[0].labels.service}}\n*Environment:* {{$json.alerts[0].labels.environment}}\n\n@channel - Immediate action required!",
        "otherOptions": {
          "username": "Alert Bot",
          "icon_emoji": ":rotating_light:"
        }
      },
      "name": "Slack - Critical",
      "type": "n8n-nodes-base.slack",
      "position": [650, 150]
    },
    {
      "parameters": {
        "to": "+66812345678,+66898765432",
        "message": "=🚨 CRITICAL: {{$json.alerts[0].labels.alertname}}\n{{$json.alerts[0].annotations.description}}",
        "options": {}
      },
      "name": "SMS - Critical (Twilio)",
      "type": "n8n-nodes-base.twilio",
      "credentials": {
        "twilioApi": {
          "id": "4",
          "name": "Twilio"
        }
      },
      "position": [650, 250]
    },
    {
      "parameters": {
        "fromEmail": "alerts@car-service.com",
        "toEmail": "devops-team@car-service.com",
        "subject": "=🚨 CRITICAL ALERT: {{$json.alerts[0].labels.alertname}}",
        "html": "=<h2 style='color: red;'>Critical Alert</h2>\n<p><strong>Alert:</strong> {{$json.alerts[0].labels.alertname}}</p>\n<p><strong>Description:</strong> {{$json.alerts[0].annotations.description}}</p>\n<p><strong>Service:</strong> {{$json.alerts[0].labels.service}}</p>\n<p><strong>Environment:</strong> {{$json.alerts[0].labels.environment}}</p>\n<p><strong>Started:</strong> {{$json.alerts[0].startsAt}}</p>"
      },
      "name": "Email - Critical",
      "type": "n8n-nodes-base.emailSend",
      "position": [650, 350]
    },
    {
      "parameters": {
        "channel": "#warnings",
        "text": "=⚠️ *Warning Alert*\n\n*Alert:* {{$json.alerts[0].labels.alertname}}\n*Description:* {{$json.alerts[0].annotations.description}}\n*Service:* {{$json.alerts[0].labels.service}}",
        "otherOptions": {
          "username": "Alert Bot"
        }
      },
      "name": "Slack - Warning",
      "type": "n8n-nodes-base.slack",
      "position": [650, 400]
    }
  ],
  "connections": {
    "Webhook - Alertmanager": {
      "main": [
        [
          {
            "node": "If Critical Alert",
            "type": "main",
            "index": 0
          },
          {
            "node": "If Warning Alert",
            "type": "main",
            "index": 0
          }
        ]
      ]
    },
    "If Critical Alert": {
      "main": [
        [
          {
            "node": "Slack - Critical",
            "type": "main",
            "index": 0
          },
          {
            "node": "SMS - Critical (Twilio)",
            "type": "main",
            "index": 0
          },
          {
            "node": "Email - Critical",
            "type": "main",
            "index": 0
          }
        ]
      ]
    },
    "If Warning Alert": {
      "main": [
        [
          {
            "node": "Slack - Warning",
            "type": "main",
            "index": 0
          }
        ]
      ]
    }
  },
  "active": true
}
```


## 4. การเชื่อมต่อ Application กับ n8n

### NestJS Integration

**backend/nestjs-api/src/notifications/n8n.service.ts**

```typescript
import { Injectable, HttpService } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';

@Injectable()
export class N8nService {
  private readonly n8nUrl: string;
  
  constructor(
    private httpService: HttpService,
    private configService: ConfigService,
  ) {
    this.n8nUrl = this.configService.get('N8N_WEBHOOK_URL');
  }

  async sendBookingNotification(bookingEvent: any): Promise<void> {
    try {
      await this.httpService
        .post(`${this.n8nUrl}/webhook/booking-webhook`, bookingEvent)
        .toPromise();
      
      console.log('Booking notification sent to n8n');
    } catch (error) {
      console.error('Failed to send notification to n8n:', error);
    }
  }

  async sendRepairNotification(repairEvent: any): Promise<void> {
    try {
      await this.httpService
        .post(`${this.n8nUrl}/webhook/repair-webhook`, repairEvent)
        .toPromise();
      
      console.log('Repair notification sent to n8n');
    } catch (error) {
      console.error('Failed to send notification to n8n:', error);
    }
  }

  async sendPaymentNotification(paymentEvent: any): Promise<void> {
    try {
      await this.httpService
        .post(`${this.n8nUrl}/webhook/payment-webhook`, paymentEvent)
        .toPromise();
      
      console.log('Payment notification sent to n8n');
    } catch (error) {
      console.error('Failed to send notification to n8n:', error);
    }
  }
}
```

**Usage in Service:**

```typescript
@Injectable()
export class BookingService {
  constructor(
    private bookingRepository: BookingRepository,
    private n8nService: N8nService,
  ) {}

  async createBooking(createBookingDto: CreateBookingDto): Promise<Booking> {
    const booking = await this.bookingRepository.save(createBookingDto);
    
    // Send notification to n8n
    await this.n8nService.sendBookingNotification({
      eventType: 'BOOKING_CREATED',
      bookingId: booking.id,
      customerId: booking.customerId,
      bookingDate: booking.bookingDate,
      serviceType: booking.serviceType,
      timestamp: new Date().toISOString(),
    });
    
    return booking;
  }
}
```


### Alertmanager Integration

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

    route:
      group_by: ['alertname', 'cluster', 'service']
      group_wait: 10s
      group_interval: 10s
      repeat_interval: 12h
      receiver: 'n8n-webhook'

    receivers:
    - name: 'n8n-webhook'
      webhook_configs:
      - url: 'https://n8n.car-service.com/webhook/alertmanager-webhook'
        send_resolved: true
        http_config:
          bearer_token: 'your-bearer-token'
```


## 5. n8n Workflow Templates

### Template: Repair Status Update

```javascript
// Workflow: Repair Status Notification
// Trigger: Webhook
// Steps:
// 1. Receive repair update
// 2. Get customer info from API
// 3. Format message based on status
// 4. Send via LINE/Email/SMS

{
  "webhook": {
    "eventType": "REPAIR_STATUS_UPDATE",
    "repairId": "repair-123",
    "customerId": "customer-456",
    "status": "IN_PROGRESS",
    "progress": 60,
    "estimatedCompletion": "2025-11-26T15:00:00Z"
  }
}

// Message template:
"สวัสดีค่ะคุณ {{customerName}}

📋 สถานะการซ่อมของคุณ
รหัสงาน: {{repairId}}
สถานะ: {{status_thai}}
ความคืบหน้า: {{progress}}%
เสร็จโดยประมาณ: {{estimatedTime}}

ขอบคุณที่ใช้บริการค่ะ"
```


## 6. Monitoring n8n

**Prometheus ServiceMonitor for n8n:**

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: n8n-monitor
  namespace: monitoring
spec:
  selector:
    matchLabels:
      app: n8n
  namespaceSelector:
    matchNames:
    - n8n
  endpoints:
  - port: http
    interval: 30s
    path: /metrics
```


## Summary: n8n Notification System

| ฟีเจอร์ | รายละเอียด | ช่องทาง |
| :-- | :-- | :-- |
| **Booking Notifications** | แจ้งเตือนการจอง/ยืนยัน/ยกเลิก | Email, LINE, SMS |
| **Repair Updates** | แจ้งสถานะการซ่อมแบบ real-time | LINE, Push Notification |
| **Payment Confirmations** | ยืนยันการชำระเงิน | Email, SMS |
| **System Alerts** | แจ้งเตือนปัญหาระบบ | Slack, SMS, Email |
| **Marketing** | ส่งโปรโมชั่น/ข่าวสาร | Email, LINE Broadcast |
| **Reminders** | เตือนนัดหมาย 1 วันก่อน | SMS, LINE |

n8n ทำให้การจัดการ notification แบบ centralized ง่ายขึ้น ไม่ต้อง hard-code ใน application สามารถปรับเปลี่ยน workflow ได้ผ่าน UI โดยไม่ต้อง redeploy code

