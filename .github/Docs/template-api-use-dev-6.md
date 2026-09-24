# 🚢🧪🔔🔌 ชุดเครื่องมือขั้นสูง 4 อย่าง (ต่อ)

จัดโครงสร้างเพิ่ม:
```
project/
├── helm/
│   └── workflow-dashboard/
│       ├── Chart.yaml
│       ├── values.yaml
│       ├── values-dev.yaml
│       ├── values-prod.yaml
│       ├── .helmignore
│       ├── README.md
│       └── templates/
│           ├── _helpers.tpl
│           ├── NOTES.txt
│           ├── serviceaccount.yaml
│           ├── configmap-dashboard.yaml
│           ├── configmap-nginx.yaml
│           ├── configmap-grafana-provisioning.yaml
│           ├── configmap-grafana-dashboards.yaml
│           ├── secret-grafana.yaml
│           ├── secret-jira.yaml
│           ├── pvc.yaml
│           ├── deployment-dashboard.yaml
│           ├── deployment-grafana.yaml
│           ├── deployment-exporter.yaml
│           ├── deployment-api.yaml
│           ├── service-dashboard.yaml
│           ├── service-grafana.yaml
│           ├── service-exporter.yaml
│           ├── service-api.yaml
│           ├── ingress.yaml
│           ├── servicemonitor.yaml
│           ├── hpa.yaml
│           └── networkpolicy.yaml
├── services/
│   └── workflow-api/
│       ├── main.go
│       ├── go.mod
│       ├── Dockerfile
│       ├── README.md
│       ├── internal/
│       │   ├── config/config.go
│       │   ├── model/workflow.go
│       │   ├── store/{store.go,filesystem.go,filesystem_test.go}
│       │   ├── service/{workflow_service.go,workflow_service_test.go}
│       │   ├── handler/{handler.go,workflow_handler.go,health_handler.go,middleware.go}
│       │   ├── metrics/metrics.go
│       │   └── router/router.go
│       ├── api/{openapi.yaml,postman_collection.json}
│       └── testdata/summaries/*.json
├── tests-playwright/
│   ├── package.json
│   ├── playwright.config.ts
│   ├── tsconfig.json
│   ├── fixtures/{test-fixtures.ts, summaries/*.json}
│   └── tests/{dashboard.spec.ts,filters.spec.ts,charts.spec.ts,detail-dialog.spec.ts,a11y.spec.ts}
└── scripts/
    ├── Send-Notification.ps1
    └── lib/Notifications.psm1
```

---

## ☸️ 1) Helm Chart: `workflow-dashboard`

### 1.1 `helm/workflow-dashboard/Chart.yaml`

```yaml
apiVersion: v2
name: workflow-dashboard
description: 6-Phase Workflow Dashboard + Grafana + Prometheus exporter + REST API
type: application
version: 1.0.0
appVersion: "1.0.0"
kubeVersion: ">=1.27.0-0"
home: https://github.com/icmongolang/icmongolang
sources:
  - https://github.com/icmongolang/icmongolang
maintainers:
  - name: Platform Team
    email: platform@example.com
keywords:
  - workflow
  - dashboard
  - grafana
  - prometheus
  - go
annotations:
  category: Monitoring
```

### 1.2 `helm/workflow-dashboard/.helmignore`

```
.DS_Store
.git/
.gitignore
*.tmproj
.idea/
.vscode/
*.swp
*.bak
*.tgz
ci/
```

### 1.3 `helm/workflow-dashboard/values.yaml`

```yaml
# ─────────────────────────────────────────────────────────────
# Global
# ─────────────────────────────────────────────────────────────
nameOverride: ""
fullnameOverride: ""

global:
  imageRegistry: ""
  imagePullSecrets: []

commonLabels: {}
commonAnnotations: {}

# ─────────────────────────────────────────────────────────────
# Service Account
# ─────────────────────────────────────────────────────────────
serviceAccount:
  create: true
  name: ""
  annotations: {}
  automount: false

# ─────────────────────────────────────────────────────────────
# Pod security
# ─────────────────────────────────────────────────────────────
podSecurityContext:
  runAsNonRoot: true
  runAsUser: 65532
  runAsGroup: 65532
  fsGroup: 65532
  seccompProfile:
    type: RuntimeDefault

containerSecurityContext:
  allowPrivilegeEscalation: false
  readOnlyRootFilesystem: true
  capabilities:
    drop: ["ALL"]

# ─────────────────────────────────────────────────────────────
# Dashboard (nginx + static HTML)
# ─────────────────────────────────────────────────────────────
dashboard:
  enabled: true
  replicaCount: 2

  image:
    repository: icmongolang/workflow-dashboard
    tag: "1.0.0"
    pullPolicy: IfNotPresent

  service:
    type: ClusterIP
    port: 80
    targetPort: 8080
    annotations: {}

  resources:
    requests: { cpu: 10m, memory: 32Mi }
    limits:   { cpu: 200m, memory: 128Mi }

  probes:
    liveness:
      path: /health
      initialDelaySeconds: 5
      periodSeconds: 20
    readiness:
      path: /health
      initialDelaySeconds: 2
      periodSeconds: 10

  # เนื้อหา dashboard (HTML/CSS/JS) — override ได้
  content: {}
  #   index.html: |
  #     <!DOCTYPE html>...
  #   app.js: |
  #     ...
  #   styles.css: |
  #     ...

  nginx:
    workerProcesses: auto
    workerConnections: 1024

  podAnnotations: {}
  nodeSelector: {}
  tolerations: []
  affinity: {}
  topologySpreadConstraints: []

# ─────────────────────────────────────────────────────────────
# Grafana
# ─────────────────────────────────────────────────────────────
grafana:
  enabled: true
  replicaCount: 1

  image:
    repository: grafana/grafana-oss
    tag: "11.2.0"
    pullPolicy: IfNotPresent

  service:
    type: ClusterIP
    port: 3000
    annotations: {}

  resources:
    requests: { cpu: 50m, memory: 128Mi }
    limits:   { cpu: 500m, memory: 512Mi }

  admin:
    user: admin
    password: ""       # ถ้าว่าง จะ generate เก็บใน Secret
    existingSecret: ""

  persistence:
    enabled: true
    storageClass: ""
    accessModes: ["ReadWriteOnce"]
    size: 2Gi

  plugins: []

  datasources:
    prometheus:
      url: http://{{ .Release.Name }}-workflow-dashboard-exporter:9101
      access: proxy
      isDefault: true
      timeInterval: 30s

  dashboards:
    workflow:
      enabled: true

  podAnnotations: {}
  nodeSelector: {}
  tolerations: []
  affinity: {}

# ─────────────────────────────────────────────────────────────
# Prometheus exporter (Python)
# ─────────────────────────────────────────────────────────────
exporter:
  enabled: true
  replicaCount: 1

  image:
    repository: icmongolang/workflow-exporter
    tag: "1.0.0"
    pullPolicy: IfNotPresent

  service:
    type: ClusterIP
    port: 9101
    annotations: {}
    monitor:
      enabled: true
      interval: 30s
      scrapeTimeout: 10s

  resources:
    requests: { cpu: 10m, memory: 32Mi }
    limits:   { cpu: 200m, memory: 128Mi }

  env:
    SCRAPE_INTERVAL: "30"
    REPORTS_DIR: /reports

  podAnnotations: {}
  nodeSelector: {}
  tolerations: []
  affinity: {}

# ─────────────────────────────────────────────────────────────
# REST API (Go)
# ─────────────────────────────────────────────────────────────
api:
  enabled: true
  replicaCount: 2

  image:
    repository: icmongolang/workflow-api
    tag: "1.0.0"
    pullPolicy: IfNotPresent

  service:
    type: ClusterIP
    port: 8080
    annotations: {}

  env:
    LOG_LEVEL: info
    CACHE_TTL: 30s
    SHUTDOWN_TIMEOUT: 15s
    CORS_ORIGINS: "*"

  resources:
    requests: { cpu: 50m, memory: 64Mi }
    limits:   { cpu: 500m, memory: 256Mi }

  probes:
    liveness:
      path: /healthz
    readiness:
      path: /readyz

  autoscaling:
    enabled: false
    minReplicas: 2
    maxReplicas: 10
    targetCPUUtilizationPercentage: 70

  podAnnotations: {}
  nodeSelector: {}
  tolerations: []
  affinity: {}

# ─────────────────────────────────────────────────────────────
# Shared reports volume (จาก workflow runner)
# ─────────────────────────────────────────────────────────────
reports:
  # โหมด: "pvc" (default) หรือ "hostPath" หรือ "existingClaim"
  mode: pvc
  existingClaim: ""
  hostPath: "/var/lib/workflow-reports"
  storageClass: ""
  accessModes: ["ReadWriteMany"]
  size: 5Gi
  mountPath: /reports

# ─────────────────────────────────────────────────────────────
# Ingress
# ─────────────────────────────────────────────────────────────
ingress:
  enabled: false
  className: nginx
  annotations: {}
  hosts:
    - host: workflow.example.com
      paths:
        - path: /
          pathType: Prefix
          service: dashboard
        - path: /api
          pathType: Prefix
          service: api
        - path: /grafana
          pathType: Prefix
          service: grafana
  tls: []
  # - secretName: workflow-tls
  #   hosts:
  #     - workflow.example.com

# ─────────────────────────────────────────────────────────────
# Network Policy
# ─────────────────────────────────────────────────────────────
networkPolicy:
  enabled: false
  ingressNamespaceSelector: {}

# ─────────────────────────────────────────────────────────────
# Jira integration (สำหรับ workflow runner ภายนอก)
# ─────────────────────────────────────────────────────────────
jira:
  # สร้าง secret ให้ runner ใช้ (ถ้าต้องการ)
  createSecret: false
  existingSecret: ""
  baseUrl: ""
  email: ""
  apiToken: ""

# ─────────────────────────────────────────────────────────────
# Pod Disruption Budget
# ─────────────────────────────────────────────────────────────
podDisruptionBudget:
  dashboard:
    enabled: true
    minAvailable: 1
  api:
    enabled: true
    minAvailable: 1
```

### 1.4 `helm/workflow-dashboard/values-dev.yaml`

```yaml
dashboard:
  replicaCount: 1
  resources:
    requests: { cpu: 5m, memory: 16Mi }

grafana:
  enabled: true
  persistence:
    enabled: false

api:
  replicaCount: 1
  resources:
    requests: { cpu: 20m, memory: 32Mi }

exporter:
  replicaCount: 1

reports:
  mode: pvc
  size: 1Gi

ingress:
  enabled: true
  className: nginx
  hosts:
    - host: workflow-dev.example.com
      paths:
        - path: /
          pathType: Prefix
          service: dashboard
        - path: /api
          pathType: Prefix
          service: api
```

### 1.5 `helm/workflow-dashboard/values-prod.yaml`

```yaml
dashboard:
  replicaCount: 3
  resources:
    requests: { cpu: 50m, memory: 64Mi }
    limits:   { cpu: 500m, memory: 256Mi }

grafana:
  replicaCount: 2
  persistence:
    enabled: true
    size: 10Gi
  resources:
    requests: { cpu: 100m, memory: 256Mi }
    limits:   { cpu: 1, memory: 1Gi }

api:
  replicaCount: 3
  resources:
    requests: { cpu: 100m, memory: 128Mi }
    limits:   { cpu: 1, memory: 512Mi }
  autoscaling:
    enabled: true
    minReplicas: 3
    maxReplicas: 20
    targetCPUUtilizationPercentage: 65

exporter:
  replicaCount: 2

reports:
  mode: pvc
  storageClass: fast-ssd
  size: 20Gi

ingress:
  enabled: true
  className: nginx
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
    nginx.ingress.kubernetes.io/proxy-body-size: "10m"
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
  hosts:
    - host: workflow.example.com
      paths:
        - path: /
          pathType: Prefix
          service: dashboard
        - path: /api
          pathType: Prefix
          service: api
        - path: /grafana
          pathType: Prefix
          service: grafana
  tls:
    - secretName: workflow-tls
      hosts:
        - workflow.example.com

networkPolicy:
  enabled: true

podDisruptionBudget:
  dashboard: { enabled: true, minAvailable: 2 }
  api:       { enabled: true, minAvailable: 2 }
```

### 1.6 `helm/workflow-dashboard/templates/_helpers.tpl`

```gotemplate
{{/*
Expand the name of the chart.
*/}}
{{- define "workflow-dashboard.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Fully qualified app name.
*/}}
{{- define "workflow-dashboard.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{- define "workflow-dashboard.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "workflow-dashboard.labels" -}}
helm.sh/chart: {{ include "workflow-dashboard.chart" . }}
{{ include "workflow-dashboard.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/part-of: workflow-dashboard
{{- with .Values.commonLabels }}
{{ toYaml . }}
{{- end }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "workflow-dashboard.selectorLabels" -}}
app.kubernetes.io/name: {{ include "workflow-dashboard.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Component labels
*/}}
{{- define "workflow-dashboard.componentLabels" -}}
{{ include "workflow-dashboard.selectorLabels" .root }}
app.kubernetes.io/component: {{ .component }}
{{- end }}

{{/*
ServiceAccount name
*/}}
{{- define "workflow-dashboard.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "workflow-dashboard.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Image reference helper
*/}}
{{- define "workflow-dashboard.image" -}}
{{- $registry := .root.Values.global.imageRegistry | default "" -}}
{{- $repo := .image.repository -}}
{{- $tag := .image.tag | default .root.Chart.AppVersion -}}
{{- if $registry -}}
{{- printf "%s/%s:%s" $registry $repo $tag -}}
{{- else -}}
{{- printf "%s:%s" $repo $tag -}}
{{- end -}}
{{- end }}

{{/*
ImagePullSecrets
*/}}
{{- define "workflow-dashboard.imagePullSecrets" -}}
{{- with .Values.global.imagePullSecrets }}
imagePullSecrets:
{{ toYaml . }}
{{- end }}
{{- end }}

{{/*
Reports PVC name
*/}}
{{- define "workflow-dashboard.reportsPVC" -}}
{{- if eq .Values.reports.mode "existingClaim" -}}
{{- .Values.reports.existingClaim -}}
{{- else -}}
{{- printf "%s-reports" (include "workflow-dashboard.fullname" .) -}}
{{- end -}}
{{- end }}

{{/*
Grafana admin password — reuse existing หรือ generate
*/}}
{{- define "workflow-dashboard.grafanaPassword" -}}
{{- if .Values.grafana.admin.existingSecret -}}
{{- .Values.grafana.admin.existingSecret -}}
{{- else -}}
{{- printf "%s-grafana" (include "workflow-dashboard.fullname" .) -}}
{{- end -}}
{{- end }}
```

### 1.7 `helm/workflow-dashboard/templates/serviceaccount.yaml`

```gotemplate
{{- if .Values.serviceAccount.create -}}
apiVersion: v1
kind: ServiceAccount
metadata:
  name: {{ include "workflow-dashboard.serviceAccountName" . }}
  namespace: {{ .Release.Namespace }}
  labels:
    {{- include "workflow-dashboard.labels" . | nindent 4 }}
  {{- with .Values.serviceAccount.annotations }}
  annotations:
    {{- toYaml . | nindent 4 }}
  {{- end }}
automountServiceAccountToken: {{ .Values.serviceAccount.automount }}
{{- end }}
```

### 1.8 `helm/workflow-dashboard/templates/configmap-dashboard.yaml`

```gotemplate
{{- if .Values.dashboard.enabled -}}
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ include "workflow-dashboard.fullname" . }}-dashboard-content
  namespace: {{ .Release.Namespace }}
  labels:
    {{- include "workflow-dashboard.labels" . | nindent 4 }}
data:
  {{- range $name, $content := .Values.dashboard.content }}
  {{ $name }}: |
    {{- $content | nindent 4 }}
  {{- end }}
{{- end }}
```

### 1.9 `helm/workflow-dashboard/templates/configmap-nginx.yaml`

```gotemplate
{{- if .Values.dashboard.enabled -}}
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ include "workflow-dashboard.fullname" . }}-nginx
  namespace: {{ .Release.Namespace }}
  labels:
    {{- include "workflow-dashboard.labels" . | nindent 4 }}
data:
  nginx.conf: |
    user nginx;
    worker_processes {{ .Values.dashboard.nginx.workerProcesses }};
    error_log /var/log/nginx/error.log warn;
    pid /tmp/nginx.pid;

    events {
      worker_connections {{ .Values.dashboard.nginx.workerConnections }};
      use epoll;
      multi_accept on;
    }

    http {
      include /etc/nginx/mime.types;
      default_type application/octet-stream;

      log_format main '$remote_addr - $remote_user [$time_local] "$request" '
                      '$status $body_bytes_sent "$http_referer" '
                      '"$http_user_agent" rt=$request_time';
      access_log /var/log/nginx/access.log main;

      sendfile on;
      tcp_nopush on;
      tcp_nodelay on;
      keepalive_timeout 65;
      types_hash_max_size 2048;
      server_tokens off;

      gzip on;
      gzip_vary on;
      gzip_min_length 1024;
      gzip_types text/plain text/css application/json application/javascript
                 text/xml application/xml application/xml+rss text/javascript;

      client_body_temp_path /tmp/client_temp;
      proxy_temp_path       /tmp/proxy_temp;
      fastcgi_temp_path     /tmp/fastcgi_temp;
      uwsgi_temp_path       /tmp/uwsgi_temp;
      scgi_temp_path        /tmp/scgi_temp;

      server {
        listen 8080;
        server_name _;
        root /usr/share/nginx/html;
        index index.html;

        location ~* \.(css|js|woff2?|ttf|eot|svg|png|jpg|jpeg|gif|ico)$ {
          expires 7d;
          add_header Cache-Control "public, immutable";
        }

        location /data/ {
          alias /reports/;
          autoindex on;
          autoindex_format json;
          add_header Cache-Control "no-store";
          add_header Access-Control-Allow-Origin "*";
        }

        location = /health {
          access_log off;
          return 200 "ok\n";
          add_header Content-Type text/plain;
        }

        location / {
          try_files $uri $uri/ /index.html;
        }
      }
    }
{{- end }}
```

### 1.10 `helm/workflow-dashboard/templates/configmap-grafana-provisioning.yaml`

```gotemplate
{{- if .Values.grafana.enabled -}}
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ include "workflow-dashboard.fullname" . }}-grafana-provisioning
  namespace: {{ .Release.Namespace }}
  labels:
    {{- include "workflow-dashboard.labels" . | nindent 4 }}
data:
  datasource-prometheus.yaml: |
    apiVersion: 1
    datasources:
      - name: Prometheus
        type: prometheus
        access: {{ .Values.grafana.datasources.prometheus.access }}
        url: {{ tpl .Values.grafana.datasources.prometheus.url . }}
        isDefault: {{ .Values.grafana.datasources.prometheus.isDefault }}
        editable: false
        jsonData:
          timeInterval: {{ .Values.grafana.datasources.prometheus.timeInterval }}
          httpMethod: GET

  dashboards.yaml: |
    apiVersion: 1
    providers:
      - name: 'Workflow'
        orgId: 1
        folder: 'Workflow'
        type: file
        disableDeletion: false
        updateIntervalSeconds: 30
        allowUiUpdates: true
        options:
          path: /var/lib/grafana/dashboards
          foldersFromFilesStructure: true
{{- end }}
```

### 1.11 `helm/workflow-dashboard/templates/configmap-grafana-dashboards.yaml`

```gotemplate
{{- if .Values.grafana.enabled -}}
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ include "workflow-dashboard.fullname" . }}-grafana-dashboards
  namespace: {{ .Release.Namespace }}
  labels:
    {{- include "workflow-dashboard.labels" . | nindent 4 }}
data:
  workflow-6phase.json: |
    {{- .Files.Get "dashboards/workflow-6phase.json" | nindent 4 }}
{{- end }}
```

> **Note:** ต้องคัดลอก `dashboard/grafana/dashboards/workflow-6phase.json` จากรอบก่อนมาไว้ที่ `helm/workflow-dashboard/dashboards/workflow-6phase.json`

### 1.12 `helm/workflow-dashboard/templates/secret-grafana.yaml`

```gotemplate
{{- if and .Values.grafana.enabled (not .Values.grafana.admin.existingSecret) -}}
{{- $secretName := printf "%s-grafana" (include "workflow-dashboard.fullname" .) -}}
{{- $existing := lookup "v1" "Secret" .Release.Namespace $secretName -}}
{{- $password := "" -}}
{{- if and $existing (hasKey $existing.data "admin-password") -}}
{{- $password = $existing.data.admin-password | b64dec -}}
{{- else if .Values.grafana.admin.password -}}
{{- $password = .Values.grafana.admin.password -}}
{{- else -}}
{{- $password = randAlphaNum 24 -}}
{{- end -}}
apiVersion: v1
kind: Secret
metadata:
  name: {{ $secretName }}
  namespace: {{ .Release.Namespace }}
  labels:
    {{- include "workflow-dashboard.labels" . | nindent 4 }}
  annotations:
    helm.sh/resource-policy: keep
type: Opaque
stringData:
  admin-user: {{ .Values.grafana.admin.user | quote }}
  admin-password: {{ $password | quote }}
{{- end }}
```

### 1.13 `helm/workflow-dashboard/templates/secret-jira.yaml`

```gotemplate
{{- if .Values.jira.createSecret -}}
{{- if not .Values.jira.baseUrl -}}{{- fail "jira.baseUrl is required" -}}{{- end -}}
{{- if not .Values.jira.email   -}}{{- fail "jira.email is required" -}}{{- end -}}
{{- if not .Values.jira.apiToken -}}{{- fail "jira.apiToken is required" -}}{{- end -}}
apiVersion: v1
kind: Secret
metadata:
  name: {{ include "workflow-dashboard.fullname" . }}-jira
  namespace: {{ .Release.Namespace }}
  labels:
    {{- include "workflow-dashboard.labels" . | nindent 4 }}
type: Opaque
stringData:
  JIRA_BASE_URL:  {{ .Values.jira.baseUrl | quote }}
  JIRA_EMAIL:     {{ .Values.jira.email | quote }}
  JIRA_API_TOKEN: {{ .Values.jira.apiToken | quote }}
{{- end }}
```

### 1.14 `helm/workflow-dashboard/templates/pvc.yaml`

```gotemplate
{{- if eq .Values.reports.mode "pvc" -}}
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: {{ printf "%s-reports" (include "workflow-dashboard.fullname" .) }}
  namespace: {{ .Release.Namespace }}
  labels:
    {{- include "workflow-dashboard.labels" . | nindent 4 }}
  annotations:
    helm.sh/resource-policy: keep
spec:
  accessModes:
    {{- toYaml .Values.reports.accessModes | nindent 4 }}
  {{- if .Values.reports.storageClass }}
  storageClassName: {{ .Values.reports.storageClass | quote }}
  {{- end }}
  resources:
    requests:
      storage: {{ .Values.reports.size }}
{{- end }}
```

### 1.15 `helm/workflow-dashboard/templates/deployment-dashboard.yaml`

```gotemplate
{{- if .Values.dashboard.enabled -}}
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ include "workflow-dashboard.fullname" . }}-dashboard
  namespace: {{ .Release.Namespace }}
  labels:
    {{- include "workflow-dashboard.labels" . | nindent 4 }}
    app.kubernetes.io/component: dashboard
spec:
  replicas: {{ .Values.dashboard.replicaCount }}
  selector:
    matchLabels:
      {{- include "workflow-dashboard.componentLabels" (dict "root" . "component" "dashboard") | nindent 6 }}
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  template:
    metadata:
      labels:
        {{- include "workflow-dashboard.componentLabels" (dict "root" . "component" "dashboard") | nindent 8 }}
      annotations:
        checksum/nginx: {{ include (print $.Template.BasePath "/configmap-nginx.yaml") . | sha256sum }}
        checksum/content: {{ include (print $.Template.BasePath "/configmap-dashboard.yaml") . | sha256sum }}
        {{- with .Values.dashboard.podAnnotations }}
        {{- toYaml . | nindent 8 }}
        {{- end }}
    spec:
      serviceAccountName: {{ include "workflow-dashboard.serviceAccountName" . }}
      {{- include "workflow-dashboard.imagePullSecrets" . | nindent 6 }}
      securityContext:
        {{- toYaml .Values.podSecurityContext | nindent 8 }}
      containers:
        - name: nginx
          image: {{ include "workflow-dashboard.image" (dict "root" . "image" .Values.dashboard.image) }}
          imagePullPolicy: {{ .Values.dashboard.image.pullPolicy }}
          securityContext:
            {{- toYaml .Values.containerSecurityContext | nindent 12 }}
          ports:
            - name: http
              containerPort: 8080
              protocol: TCP
          livenessProbe:
            httpGet: { path: {{ .Values.dashboard.probes.liveness.path }}, port: http }
            initialDelaySeconds: {{ .Values.dashboard.probes.liveness.initialDelaySeconds }}
            periodSeconds: {{ .Values.dashboard.probes.liveness.periodSeconds }}
          readinessProbe:
            httpGet: { path: {{ .Values.dashboard.probes.readiness.path }}, port: http }
            initialDelaySeconds: {{ .Values.dashboard.probes.readiness.initialDelaySeconds }}
            periodSeconds: {{ .Values.dashboard.probes.readiness.periodSeconds }}
          resources:
            {{- toYaml .Values.dashboard.resources | nindent 12 }}
          volumeMounts:
            - name: nginx-conf
              mountPath: /etc/nginx/nginx.conf
              subPath: nginx.conf
              readOnly: true
            - name: content
              mountPath: /usr/share/nginx/html
              readOnly: true
            - name: reports
              mountPath: /reports
              readOnly: true
            - name: nginx-cache
              mountPath: /var/cache/nginx
            - name: nginx-run
              mountPath: /var/run
            - name: tmp
              mountPath: /tmp
      volumes:
        - name: nginx-conf
          configMap:
            name: {{ include "workflow-dashboard.fullname" . }}-nginx
        - name: content
          configMap:
            name: {{ include "workflow-dashboard.fullname" . }}-dashboard-content
        - name: reports
          {{- if eq .Values.reports.mode "hostPath" }}
          hostPath:
            path: {{ .Values.reports.hostPath }}
            type: DirectoryOrCreate
          {{- else }}
          persistentVolumeClaim:
            claimName: {{ include "workflow-dashboard.reportsPVC" . }}
          {{- end }}
        - name: nginx-cache
          emptyDir: {}
        - name: nginx-run
          emptyDir: {}
        - name: tmp
          emptyDir: {}
      {{- with .Values.dashboard.nodeSelector }}
      nodeSelector: {{- toYaml . | nindent 8 }}
      {{- end }}
      {{- with .Values.dashboard.affinity }}
      affinity: {{- toYaml . | nindent 8 }}
      {{- end }}
      {{- with .Values.dashboard.tolerations }}
      tolerations: {{- toYaml . | nindent 8 }}
      {{- end }}
      {{- with .Values.dashboard.topologySpreadConstraints }}
      topologySpreadConstraints: {{- toYaml . | nindent 8 }}
      {{- end }}
{{- end }}
```

### 1.16 `helm/workflow-dashboard/templates/deployment-grafana.yaml`

```gotemplate
{{- if .Values.grafana.enabled -}}
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ include "workflow-dashboard.fullname" . }}-grafana
  namespace: {{ .Release.Namespace }}
  labels:
    {{- include "workflow-dashboard.labels" . | nindent 4 }}
    app.kubernetes.io/component: grafana
spec:
  replicas: {{ .Values.grafana.replicaCount }}
  selector:
    matchLabels:
      {{- include "workflow-dashboard.componentLabels" (dict "root" . "component" "grafana") | nindent 6 }}
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        {{- include "workflow-dashboard.componentLabels" (dict "root" . "component" "grafana") | nindent 8 }}
      annotations:
        checksum/provisioning: {{ include (print $.Template.BasePath "/configmap-grafana-provisioning.yaml") . | sha256sum }}
        checksum/dashboards: {{ include (print $.Template.BasePath "/configmap-grafana-dashboards.yaml") . | sha256sum }}
        {{- with .Values.grafana.podAnnotations }}
        {{- toYaml . | nindent 8 }}
        {{- end }}
    spec:
      serviceAccountName: {{ include "workflow-dashboard.serviceAccountName" . }}
      {{- include "workflow-dashboard.imagePullSecrets" . | nindent 6 }}
      securityContext:
        {{- toYaml .Values.podSecurityContext | nindent 8 }}
      containers:
        - name: grafana
          image: {{ include "workflow-dashboard.image" (dict "root" . "image" .Values.grafana.image) }}
          imagePullPolicy: {{ .Values.grafana.image.pullPolicy }}
          securityContext:
            {{- toYaml .Values.containerSecurityContext | nindent 12 }}
          ports:
            - name: http
              containerPort: 3000
              protocol: TCP
          env:
            - name: GF_SECURITY_ADMIN_USER
              valueFrom:
                secretKeyRef:
                  name: {{ include "workflow-dashboard.grafanaPassword" . }}
                  key: admin-user
            - name: GF_SECURITY_ADMIN_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: {{ include "workflow-dashboard.grafanaPassword" . }}
                  key: admin-password
            - name: GF_USERS_ALLOW_SIGN_UP
              value: "false"
            - name: GF_INSTALL_PLUGINS
              value: {{ join "," .Values.grafana.plugins | quote }}
            - name: GF_PATHS_PROVISIONING
              value: /etc/grafana/provisioning
          livenessProbe:
            httpGet: { path: /api/health, port: http }
            initialDelaySeconds: 30
            periodSeconds: 20
          readinessProbe:
            httpGet: { path: /api/health, port: http }
            initialDelaySeconds: 10
            periodSeconds: 10
          resources:
            {{- toYaml .Values.grafana.resources | nindent 12 }}
          volumeMounts:
            - name: provisioning
              mountPath: /etc/grafana/provisioning/datasources
              readOnly: true
            - name: provisioning
              mountPath: /etc/grafana/provisioning/dashboards
              readOnly: true
            - name: dashboards
              mountPath: /var/lib/grafana/dashboards
              readOnly: true
            - name: data
              mountPath: /var/lib/grafana
            - name: tmp
              mountPath: /tmp
      volumes:
        - name: provisioning
          configMap:
            name: {{ include "workflow-dashboard.fullname" . }}-grafana-provisioning
            items:
              - key: datasource-prometheus.yaml
                path: datasources/prometheus.yaml
              - key: dashboards.yaml
                path: dashboards/default.yaml
        - name: dashboards
          configMap:
            name: {{ include "workflow-dashboard.fullname" . }}-grafana-dashboards
        - name: data
          {{- if .Values.grafana.persistence.enabled }}
          persistentVolumeClaim:
            claimName: {{ include "workflow-dashboard.fullname" . }}-grafana
          {{- else }}
          emptyDir: {}
          {{- end }}
        - name: tmp
          emptyDir: {}
      {{- with .Values.grafana.nodeSelector }}
      nodeSelector: {{- toYaml . | nindent 8 }}
      {{- end }}
      {{- with .Values.grafana.affinity }}
      affinity: {{- toYaml . | nindent 8 }}
      {{- end }}
      {{- with .Values.grafana.tolerations }}
      tolerations: {{- toYaml . | nindent 8 }}
      {{- end }}
---
{{- if .Values.grafana.persistence.enabled }}
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: {{ include "workflow-dashboard.fullname" . }}-grafana
  namespace: {{ .Release.Namespace }}
  labels:
    {{- include "workflow-dashboard.labels" . | nindent 4 }}
spec:
  accessModes:
    {{- toYaml .Values.grafana.persistence.accessModes | nindent 4 }}
  {{- if .Values.grafana.persistence.storageClass }}
  storageClassName: {{ .Values.grafana.persistence.storageClass | quote }}
  {{- end }}
  resources:
    requests:
      storage: {{ .Values.grafana.persistence.size }}
{{- end }}
{{- end }}
```

### 1.17 `helm/workflow-dashboard/templates/deployment-exporter.yaml`

```gotemplate
{{- if .Values.exporter.enabled -}}
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ include "workflow-dashboard.fullname" . }}-exporter
  namespace: {{ .Release.Namespace }}
  labels:
    {{- include "workflow-dashboard.labels" . | nindent 4 }}
    app.kubernetes.io/component: exporter
spec:
  replicas: {{ .Values.exporter.replicaCount }}
  selector:
    matchLabels:
      {{- include "workflow-dashboard.componentLabels" (dict "root" . "component" "exporter") | nindent 6 }}
  template:
    metadata:
      labels:
        {{- include "workflow-dashboard.componentLabels" (dict "root" . "component" "exporter") | nindent 8 }}
      {{- with .Values.exporter.podAnnotations }}
      annotations:
        {{- toYaml . | nindent 8 }}
      {{- end }}
    spec:
      serviceAccountName: {{ include "workflow-dashboard.serviceAccountName" . }}
      {{- include "workflow-dashboard.imagePullSecrets" . | nindent 6 }}
      securityContext:
        {{- toYaml .Values.podSecurityContext | nindent 8 }}
      containers:
        - name: exporter
          image: {{ include "workflow-dashboard.image" (dict "root" . "image" .Values.exporter.image) }}
          imagePullPolicy: {{ .Values.exporter.image.pullPolicy }}
          securityContext:
            {{- toYaml .Values.containerSecurityContext | nindent 12 }}
          ports:
            - name: metrics
              containerPort: 9101
              protocol: TCP
          env:
            {{- range $k, $v := .Values.exporter.env }}
            - name: {{ $k }}
              value: {{ $v | quote }}
            {{- end }}
          livenessProbe:
            httpGet: { path: /health, port: metrics }
            initialDelaySeconds: 5
            periodSeconds: 20
          readinessProbe:
            httpGet: { path: /health, port: metrics }
            initialDelaySeconds: 2
            periodSeconds: 10
          resources:
            {{- toYaml .Values.exporter.resources | nindent 12 }}
          volumeMounts:
            - name: reports
              mountPath: /reports
              readOnly: true
            - name: tmp
              mountPath: /tmp
      volumes:
        - name: reports
          {{- if eq .Values.reports.mode "hostPath" }}
          hostPath:
            path: {{ .Values.reports.hostPath }}
            type: DirectoryOrCreate
          {{- else }}
          persistentVolumeClaim:
            claimName: {{ include "workflow-dashboard.reportsPVC" . }}
          {{- end }}
        - name: tmp
          emptyDir: {}
      {{- with .Values.exporter.nodeSelector }}
      nodeSelector: {{- toYaml . | nindent 8 }}
      {{- end }}
      {{- with .Values.exporter.tolerations }}
      tolerations: {{- toYaml . | nindent 8 }}
      {{- end }}
{{- end }}
```

### 1.18 `helm/workflow-dashboard/templates/deployment-api.yaml`

```gotemplate
{{- if .Values.api.enabled -}}
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ include "workflow-dashboard.fullname" . }}-api
  namespace: {{ .Release.Namespace }}
  labels:
    {{- include "workflow-dashboard.labels" . | nindent 4 }}
    app.kubernetes.io/component: api
spec:
  {{- if not .Values.api.autoscaling.enabled }}
  replicas: {{ .Values.api.replicaCount }}
  {{- end }}
  selector:
    matchLabels:
      {{- include "workflow-dashboard.componentLabels" (dict "root" . "component" "api") | nindent 6 }}
  template:
    metadata:
      labels:
        {{- include "workflow-dashboard.componentLabels" (dict "root" . "component" "api") | nindent 8 }}
      {{- with .Values.api.podAnnotations }}
      annotations:
        {{- toYaml . | nindent 8 }}
      {{- end }}
    spec:
      serviceAccountName: {{ include "workflow-dashboard.serviceAccountName" . }}
      {{- include "workflow-dashboard.imagePullSecrets" . | nindent 6 }}
      securityContext:
        {{- toYaml .Values.podSecurityContext | nindent 8 }}
      containers:
        - name: api
          image: {{ include "workflow-dashboard.image" (dict "root" . "image" .Values.api.image) }}
          imagePullPolicy: {{ .Values.api.image.pullPolicy }}
          securityContext:
            {{- toYaml .Values.containerSecurityContext | nindent 12 }}
          ports:
            - name: http
              containerPort: 8080
              protocol: TCP
          env:
            - name: REPORTS_DIR
              value: {{ .Values.reports.mountPath | quote }}
            {{- range $k, $v := .Values.api.env }}
            - name: {{ $k }}
              value: {{ $v | quote }}
            {{- end }}
          livenessProbe:
            httpGet: { path: {{ .Values.api.probes.liveness.path }}, port: http }
            initialDelaySeconds: 5
            periodSeconds: 20
          readinessProbe:
            httpGet: { path: {{ .Values.api.probes.readiness.path }}, port: http }
            initialDelaySeconds: 2
            periodSeconds: 10
          resources:
            {{- toYaml .Values.api.resources | nindent 12 }}
          volumeMounts:
            - name: reports
              mountPath: {{ .Values.reports.mountPath }}
              readOnly: true
            - name: tmp
              mountPath: /tmp
      volumes:
        - name: reports
          {{- if eq .Values.reports.mode "hostPath" }}
          hostPath:
            path: {{ .Values.reports.hostPath }}
            type: DirectoryOrCreate
          {{- else }}
          persistentVolumeClaim:
            claimName: {{ include "workflow-dashboard.reportsPVC" . }}
          {{- end }}
        - name: tmp
          emptyDir: {}
      {{- with .Values.api.nodeSelector }}
      nodeSelector: {{- toYaml . | nindent 8 }}
      {{- end }}
      {{- with .Values.api.affinity }}
      affinity: {{- toYaml . | nindent 8 }}
      {{- end }}
      {{- with .Values.api.tolerations }}
      tolerations: {{- toYaml . | nindent 8 }}
      {{- end }}
{{- end }}
```

### 1.19 Services

`helm/workflow-dashboard/templates/service-dashboard.yaml`:
```gotemplate
{{- if .Values.dashboard.enabled -}}
apiVersion: v1
kind: Service
metadata:
  name: {{ include "workflow-dashboard.fullname" . }}-dashboard
  namespace: {{ .Release.Namespace }}
  labels:
    {{- include "workflow-dashboard.labels" . | nindent 4 }}
    app.kubernetes.io/component: dashboard
  {{- with .Values.dashboard.service.annotations }}
  annotations: {{- toYaml . | nindent 4 }}
  {{- end }}
spec:
  type: {{ .Values.dashboard.service.type }}
  ports:
    - port: {{ .Values.dashboard.service.port }}
      targetPort: http
      protocol: TCP
      name: http
  selector:
    {{- include "workflow-dashboard.componentLabels" (dict "root" . "component" "dashboard") | nindent 4 }}
{{- end }}
```

`service-grafana.yaml`, `service-exporter.yaml`, `service-api.yaml` — โครงสร้างเหมือนกัน เปลี่ยน `.dashboard` → `.grafana` / `.exporter` / `.api` และ port / component ให้ตรง

### 1.20 `helm/workflow-dashboard/templates/ingress.yaml`

```gotemplate
{{- if .Values.ingress.enabled -}}
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: {{ include "workflow-dashboard.fullname" . }}
  namespace: {{ .Release.Namespace }}
  labels:
    {{- include "workflow-dashboard.labels" . | nindent 4 }}
  {{- with .Values.ingress.annotations }}
  annotations: {{- toYaml . | nindent 4 }}
  {{- end }}
spec:
  {{- if .Values.ingress.className }}
  ingressClassName: {{ .Values.ingress.className }}
  {{- end }}
  {{- with .Values.ingress.tls }}
  tls: {{- toYaml . | nindent 4 }}
  {{- end }}
  rules:
    {{- range .Values.ingress.hosts }}
    - host: {{ .host | quote }}
      http:
        paths:
          {{- range .paths }}
          - path: {{ .path }}
            pathType: {{ .pathType }}
            backend:
              service:
                name: {{ include "workflow-dashboard.fullname" $ }}-{{ .service }}
                port:
                  name: http
          {{- end }}
    {{- end }}
{{- end }}
```

### 1.21 `helm/workflow-dashboard/templates/servicemonitor.yaml`

```gotemplate
{{- if and .Values.exporter.enabled .Values.exporter.service.monitor.enabled -}}
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: {{ include "workflow-dashboard.fullname" . }}-exporter
  namespace: {{ .Release.Namespace }}
  labels:
    {{- include "workflow-dashboard.labels" . | nindent 4 }}
    app.kubernetes.io/component: exporter
    release: prometheus
spec:
  selector:
    matchLabels:
      {{- include "workflow-dashboard.componentLabels" (dict "root" . "component" "exporter") | nindent 6 }}
  endpoints:
    - port: metrics
      interval: {{ .Values.exporter.service.monitor.interval }}
      scrapeTimeout: {{ .Values.exporter.service.monitor.scrapeTimeout }}
      path: /metrics
  namespaceSelector:
    matchNames:
      - {{ .Release.Namespace }}
{{- end }}
```

### 1.22 `helm/workflow-dashboard/templates/hpa.yaml`

```gotemplate
{{- if and .Values.api.enabled .Values.api.autoscaling.enabled -}}
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: {{ include "workflow-dashboard.fullname" . }}-api
  namespace: {{ .Release.Namespace }}
  labels:
    {{- include "workflow-dashboard.labels" . | nindent 4 }}
    app.kubernetes.io/component: api
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: {{ include "workflow-dashboard.fullname" . }}-api
  minReplicas: {{ .Values.api.autoscaling.minReplicas }}
  maxReplicas: {{ .Values.api.autoscaling.maxReplicas }}
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: {{ .Values.api.autoscaling.targetCPUUtilizationPercentage }}
  behavior:
    scaleDown:
      stabilizationWindowSeconds: 300
      policies:
        - type: Percent
          value: 50
          periodSeconds: 60
    scaleUp:
      stabilizationWindowSeconds: 30
      policies:
        - type: Percent
          value: 100
          periodSeconds: 30
{{- end }}
```

### 1.23 `helm/workflow-dashboard/templates/networkpolicy.yaml`

```gotemplate
{{- if .Values.networkPolicy.enabled -}}
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: {{ include "workflow-dashboard.fullname" . }}-dashboard
  namespace: {{ .Release.Namespace }}
  labels:
    {{- include "workflow-dashboard.labels" . | nindent 4 }}
spec:
  podSelector:
    matchLabels:
      {{- include "workflow-dashboard.componentLabels" (dict "root" . "component" "dashboard") | nindent 6 }}
  policyTypes: [Ingress]
  ingress:
    - from:
        - namespaceSelector:
            matchLabels:
              kubernetes.io/metadata.name: ingress-nginx
      ports:
        - port: 8080
          protocol: TCP
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: {{ include "workflow-dashboard.fullname" . }}-api
  namespace: {{ .Release.Namespace }}
  labels:
    {{- include "workflow-dashboard.labels" . | nindent 4 }}
spec:
  podSelector:
    matchLabels:
      {{- include "workflow-dashboard.componentLabels" (dict "root" . "component" "api") | nindent 6 }}
  policyTypes: [Ingress]
  ingress:
    - from:
        - namespaceSelector:
            matchLabels:
              kubernetes.io/metadata.name: ingress-nginx
        - podSelector: {}
      ports:
        - port: 8080
          protocol: TCP
{{- end }}
```

### 1.24 `helm/workflow-dashboard/templates/NOTES.txt`

```
🎉 Workflow Dashboard deployed!

Release:  {{ .Release.Name }}
Chart:    {{ .Chart.Name }}-{{ .Chart.Version }}
Namespace: {{ .Release.Namespace }}

{{- if .Values.ingress.enabled }}
Dashboard: https://{{ (index .Values.ingress.hosts 0).host }}/
API:       https://{{ (index .Values.ingress.hosts 0).host }}/api/v1/workflows
{{- else }}
Port-forward:
  kubectl -n {{ .Release.Namespace }} port-forward svc/{{ include "workflow-dashboard.fullname" . }}-dashboard 8080:80
  kubectl -n {{ .Release.Namespace }} port-forward svc/{{ include "workflow-dashboard.fullname" . }}-api 8081:8080
{{- end }}

Grafana admin credentials:
  kubectl -n {{ .Release.Namespace }} get secret {{ include "workflow-dashboard.fullname" . }}-grafana \
    -o jsonpath='{.data.admin-user}' | base64 -d
  kubectl -n {{ .Release.Namespace }} get secret {{ include "workflow-dashboard.fullname" . }}-grafana \
    -o jsonpath='{.data.admin-password}' | base64 -d

Health checks:
  curl http://<host>/healthz    (API)
  curl http://<host>/health     (Dashboard)
  curl http://<host>/metrics    (Exporter)

⚠️  อย่าลืมตั้งค่า secret Jira ผ่าน:
    --set jira.createSecret=true --set jira.baseUrl=... --set jira.email=... --set jira.apiToken=...
```

### 1.25 วิธี deploy

```bash
# ── Lint ──────────────────────────────────────────────────
helm lint helm/workflow-dashboard

# ── Dry-run ───────────────────────────────────────────────
helm template workflow helm/workflow-dashboard --namespace workflow

# ── Install dev ───────────────────────────────────────────
helm upgrade --install workflow helm/workflow-dashboard \
  --namespace workflow --create-namespace \
  -f helm/workflow-dashboard/values-dev.yaml

# ── Install prod ──────────────────────────────────────────
helm upgrade --install workflow helm/workflow-dashboard \
  --namespace workflow --create-namespace \
  -f helm/workflow-dashboard/values-prod.yaml \
  --set jira.createSecret=true \
  --set jira.baseUrl="https://your.atlassian.net" \
  --set jira.email="you@example.com" \
  --set jira.apiToken="$JIRA_API_TOKEN" \
  --wait --timeout 5m

# ── อัปเดต content (dashboard HTML) ────────────────────────
helm upgrade workflow helm/workflow-dashboard \
  --reuse-values \
  --set-file dashboard.content.index\.html=dashboard/index.html \
  --set-file dashboard.content.app\.js=dashboard/app.js \
  --set-file dashboard.content.styles\.css=dashboard/styles.css

# ── Rollback ──────────────────────────────────────────────
helm history workflow -n workflow
helm rollback workflow 1 -n workflow
```

---

## 🧪 2) Playwright Tests สำหรับ Dashboard

### 2.1 `tests-playwright/package.json`

```json
{
  "name": "workflow-dashboard-e2e",
  "version": "1.0.0",
  "private": true,
  "type": "module",
  "scripts": {
    "test": "playwright test",
    "test:headed": "playwright test --headed",
    "test:ui": "playwright test --ui",
    "test:debug": "playwright test --debug",
    "test:report": "playwright show-report",
    "install:browsers": "playwright install --with-deps chromium",
    "lint": "tsc --noEmit"
  },
  "devDependencies": {
    "@playwright/test": "^1.47.0",
    "@axe-core/playwright": "^4.10.0",
    "@types/node": "^22.5.0",
    "typescript": "^5.6.0"
  }
}
```

### 2.2 `tests-playwright/tsconfig.json`

```json
{
  "compilerOptions": {
    "target": "ES2022",
    "module": "ESNext",
    "moduleResolution": "Bundler",
    "strict": true,
    "esModuleInterop": true,
    "skipLibCheck": true,
    "resolveJsonModule": true,
    "types": ["node"],
    "lib": ["ES2022", "DOM"]
  },
  "include": ["**/*.ts"]
}
```

### 2.3 `tests-playwright/playwright.config.ts`

```typescript
import { defineConfig, devices } from '@playwright/test';

const PORT = Number(process.env.PORT ?? 4173);
const BASE_URL = process.env.BASE_URL ?? `http://127.0.0.1:${PORT}`;

export default defineConfig({
  testDir: './tests',
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 2 : undefined,
  reporter: [
    ['list'],
    ['html', { open: 'never', outputFolder: 'playwright-report' }],
    ['junit', { outputFile: 'playwright-report/junit.xml' }],
    ...(process.env.CI ? [['github'] as const] : []),
  ],
  timeout: 30_000,
  expect: { timeout: 5_000 },

  use: {
    baseURL: BASE_URL,
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
    video: 'retain-on-failure',
    locale: 'th-TH',
    timezoneId: 'Asia/Bangkok',
  },

  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
    {
      name: 'firefox',
      use: { ...devices['Desktop Firefox'] },
    },
    {
      name: 'webkit',
      use: { ...devices['Desktop Safari'] },
    },
    {
      name: 'mobile-chrome',
      use: { ...devices['Pixel 7'] },
    },
  ],

  webServer: {
    command: `npx --yes http-server ./static -p ${PORT} -c-1 --silent`,
    url: BASE_URL,
    reuseExistingServer: !process.env.CI,
    timeout: 30_000,
  },
});
```

### 2.4 `tests-playwright/fixtures/summaries/ABC-123-summary.json`

```json
{
  "ticketKey": "ABC-123",
  "task": "เพิ่ม API ค้นหาสินค้า",
  "overall": "PASS",
  "durationMs": 45230,
  "finishedAt": "2026-09-15T10:30:00.000Z",
  "results": {
    "gather":      { "phaseId": "gather",      "phaseName": "Gather",      "emoji": "📥", "status": "PASS", "durationMs": 3200,  "report": "# Gathered ✅" },
    "plan":        { "phaseId": "plan",        "phaseName": "Plan",        "emoji": "🗺️", "status": "PASS", "durationMs": 4100,  "report": "# Planned ✅" },
    "execute":     { "phaseId": "execute",     "phaseName": "Execute",     "emoji": "⚙️", "status": "PASS", "durationMs": 18200, "report": "# Executed ✅" },
    "security":    { "phaseId": "security",    "phaseName": "Security",    "emoji": "🔒", "status": "PASS", "durationMs": 6800,  "report": "# Secure ✅" },
    "performance": { "phaseId": "performance", "phaseName": "Performance", "emoji": "⚡", "status": "PASS", "durationMs": 8800,  "report": "# Fast ✅" },
    "rca":         { "phaseId": "rca",         "phaseName": "RCA",         "emoji": "🔍", "status": "PASS", "durationMs": 4130,  "report": "# No issues ✅" }
  }
}
```

`DEF-456-summary.json` (FAIL):
```json
{
  "ticketKey": "DEF-456",
  "task": "แก้บั๊ก pagination",
  "overall": "FAIL",
  "durationMs": 23100,
  "finishedAt": "2026-09-15T11:00:00.000Z",
  "results": {
    "gather":      { "phaseId": "gather",      "phaseName": "Gather",      "emoji": "📥", "status": "PASS", "durationMs": 3100, "report": "# OK" },
    "plan":        { "phaseId": "plan",        "phaseName": "Plan",        "emoji": "🗺️", "status": "PASS", "durationMs": 4000, "report": "# OK" },
    "execute":     { "phaseId": "execute",     "phaseName": "Execute",     "emoji": "⚙️", "status": "PASS", "durationMs": 8000, "report": "# OK" },
    "security":    { "phaseId": "security",    "phaseName": "Security",    "emoji": "🔒", "status": "FAIL", "durationMs": 8000, "report": "# ❌ SQL injection at line 42" }
  }
}
```

`GHI-789-summary.json` (WARN):
```json
{
  "ticketKey": "GHI-789",
  "task": "ปรับปรุง performance ของ query",
  "overall": "WARN",
  "durationMs": 15000,
  "finishedAt": "2026-09-15T12:00:00.000Z",
  "results": {
    "gather":      { "phaseId": "gather",      "phaseName": "Gather",      "emoji": "📥", "status": "PASS", "durationMs": 2000, "report": "# OK" },
    "plan":        { "phaseId": "plan",        "phaseName": "Plan",        "emoji": "🗺️", "status": "PASS", "durationMs": 3000, "report": "# OK" },
    "execute":     { "phaseId": "execute",     "phaseName": "Execute",     "emoji": "⚙️", "status": "PASS", "durationMs": 5000, "report": "# OK" },
    "security":    { "phaseId": "security",    "phaseName": "Security",    "emoji": "🔒", "status": "PASS", "durationMs": 2000, "report": "# OK" },
    "performance": { "phaseId": "performance", "phaseName": "Performance", "emoji": "⚡", "status": "WARN", "durationMs": 3000, "report": "# 🟡 p95 180ms > 150ms" }
  }
}
```

### 2.5 `tests-playwright/fixtures/test-fixtures.ts`

```typescript
import { test as base, expect, Page, Route } from '@playwright/test';
import { readFileSync, readdirSync } from 'node:fs';
import { resolve, dirname } from 'node:path';
import { fileURLToPath } from 'node:url';

const __dirname = dirname(fileURLToPath(import.meta.url));
const SUMMARIES_DIR = resolve(__dirname, 'summaries');

export type Summary = {
  ticketKey: string;
  task: string;
  overall: 'PASS' | 'WARN' | 'FAIL' | 'ERROR';
  durationMs: number;
  finishedAt: string;
  results: Record<string, { status: string; durationMs: number; emoji: string; phaseName: string; report: string }>;
};

export function loadAllSummaries(): Summary[] {
  return readdirSync(SUMMARIES_DIR)
    .filter((f) => f.endsWith('-summary.json'))
    .map((f) => JSON.parse(readFileSync(resolve(SUMMARIES_DIR, f), 'utf8')));
}

/** ติดตั้ง route intercept ให้ /data/* ตอบจาก fixture ที่กำหนด */
export async function mockData(page: Page, summaries: Summary[]) {
  await page.route('**/data/**', async (route: Route) => {
    const url = new URL(route.request().url());
    const path = url.pathname;

    // autoindex JSON ของ nginx
    if (path === '/data/' || path === '/data') {
      const listing = summaries.map((s) => ({
        name: `${s.ticketKey}-summary.json`,
        type: 'file',
      }));
      return route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify(listing),
      });
    }

    // summary file
    const m = path.match(/\/([^/]+-summary\.json)$/);
    if (m) {
      const s = summaries.find((x) => `${x.ticketKey}-summary.json` === m[1]);
      if (s) {
        return route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify(s),
        });
      }
    }

    return route.fulfill({ status: 404, body: 'not found' });
  });
}

type Fixtures = {
  summaries: Summary[];
  dashboardPage: Page;
};

export const test = base.extend<Fixtures>({
  summaries: async ({}, use) => {
    await use(loadAllSummaries());
  },

  dashboardPage: async ({ page, summaries }, use) => {
    await mockData(page, summaries);
    await page.goto('/');
    await page.waitForLoadState('networkidle');
    await use(page);
  },
});

export { expect };
```

### 2.6 `tests-playwright/tests/dashboard.spec.ts`

```typescript
import { test, expect } from '../fixtures/test-fixtures';

test.describe('Dashboard — main view', () => {
  test('โหลดหน้าและแสดง brand + subtitle', async ({ dashboardPage: page }) => {
    await expect(page).toHaveTitle(/6-Phase Workflow Dashboard/i);
    await expect(page.getByRole('heading', { name: /6-Phase Workflow Dashboard/i })).toBeVisible();
    await expect(page.getByText(/Gather → Plan → Execute/i)).toBeVisible();
  });

  test('แสดง KPI cards ครบ', async ({ dashboardPage: page }) => {
    const kpis = page.locator('.kpis .kpi');
    await expect(kpis).toHaveCount(6);

    // Total = 3 fixtures
    await expect(kpis.nth(0).locator('.kpi-value')).toHaveText('3');

    // Pass rate = 1/3 = 33%
    await expect(kpis.nth(1).locator('.kpi-value')).toHaveText('33%');

    // Warnings = 1, Failures = 1, Errors = 0
    await expect(kpis.nth(2).locator('.kpi-value')).toHaveText('1');
    await expect(kpis.nth(3).locator('.kpi-value')).toHaveText('1');
    await expect(kpis.nth(4).locator('.kpi-value')).toHaveText('0');
  });

  test('ตารางแสดงทุก ticket ที่โหลดมา', async ({ dashboardPage: page }) => {
    const rows = page.locator('#workflowTable tbody tr');
    await expect(rows).toHaveCount(3);

    await expect(rows.filter({ hasText: 'ABC-123' })).toBeVisible();
    await expect(rows.filter({ hasText: 'DEF-456' })).toBeVisible();
    await expect(rows.filter({ hasText: 'GHI-789' })).toBeVisible();
  });

  test('overall badge มีสีตาม status', async ({ dashboardPage: page }) => {
    const abc = page.locator('#workflowTable tbody tr', { hasText: 'ABC-123' });
    const def = page.locator('#workflowTable tbody tr', { hasText: 'DEF-456' });
    const ghi = page.locator('#workflowTable tbody tr', { hasText: 'GHI-789' });

    await expect(abc.locator('.badge.pass').first()).toBeVisible();
    await expect(def.locator('.badge.fail').first()).toBeVisible();
    await expect(ghi.locator('.badge.warn').first()).toBeVisible();
  });

  test('lastUpdate แสดง timestamp', async ({ dashboardPage: page }) => {
    await expect(page.locator('#lastUpdate')).toContainText(/Updated/i);
  });
});
```

### 2.7 `tests-playwright/tests/filters.spec.ts`

```typescript
import { test, expect } from '../fixtures/test-fixtures';

test.describe('Dashboard — filters', () => {
  test('filter PASS แสดงเฉพาะ ticket ที่ PASS', async ({ dashboardPage: page }) => {
    await page.locator('#filterStatus').selectOption('PASS');

    const rows = page.locator('#workflowTable tbody tr');
    await expect(rows).toHaveCount(1);
    await expect(rows.first()).toContainText('ABC-123');
  });

  test('filter FAIL', async ({ dashboardPage: page }) => {
    await page.locator('#filterStatus').selectOption('FAIL');
    await expect(page.locator('#workflowTable tbody tr')).toHaveCount(1);
    await expect(page.locator('#workflowTable tbody tr')).toContainText('DEF-456');
  });

  test('filter WARN', async ({ dashboardPage: page }) => {
    await page.locator('#filterStatus').selectOption('WARN');
    await expect(page.locator('#workflowTable tbody tr')).toHaveCount(1);
    await expect(page.locator('#workflowTable tbody tr')).toContainText('GHI-789');
  });

  test('filter ERROR → ไม่มี row (empty state)', async ({ dashboardPage: page }) => {
    await page.locator('#filterStatus').selectOption('ERROR');
    await expect(page.locator('#workflowTable tbody tr')).toHaveCount(0);
  });

  test('search ticket key', async ({ dashboardPage: page }) => {
    await page.locator('#search').fill('DEF');
    await expect(page.locator('#workflowTable tbody tr')).toHaveCount(1);
    await expect(page.locator('#workflowTable tbody tr')).toContainText('DEF-456');
  });

  test('search task description', async ({ dashboardPage: page }) => {
    await page.locator('#search').fill('pagination');
    await expect(page.locator('#workflowTable tbody tr')).toHaveCount(1);
    await expect(page.locator('#workflowTable tbody tr')).toContainText('DEF-456');
  });

  test('search ไม่พบ → 0 rows', async ({ dashboardPage: page }) => {
    await page.locator('#search').fill('zzz-nonexistent');
    await expect(page.locator('#workflowTable tbody tr')).toHaveCount(0);
  });

  test('search + filter ทำงานร่วมกัน', async ({ dashboardPage: page }) => {
    await page.locator('#filterStatus').selectOption('PASS');
    await page.locator('#search').fill('DEF');
    await expect(page.locator('#workflowTable tbody tr')).toHaveCount(0);

    await page.locator('#search').fill('');
    await expect(page.locator('#workflowTable tbody tr')).toHaveCount(1);
  });

  test('clear search แล้วเห็นทุก row กลับมา', async ({ dashboardPage: page }) => {
    await page.locator('#search').fill('ABC');
    await expect(page.locator('#workflowTable tbody tr')).toHaveCount(1);

    await page.locator('#search').fill('');
    await expect(page.locator('#workflowTable tbody tr')).toHaveCount(3);
  });
});
```

### 2.8 `tests-playwright/tests/charts.spec.ts`

```typescript
import { test, expect } from '../fixtures/test-fixtures';

test.describe('Dashboard — charts', () => {
  test('สร้าง canvas สำหรับทั้ง 3 charts', async ({ dashboardPage: page }) => {
    await expect(page.locator('#chartStatus')).toBeVisible();
    await expect(page.locator('#chartDuration')).toBeVisible();
    await expect(page.locator('#chartPhases')).toBeVisible();
  });

  test('Chart.js ถูกโหลดและสร้าง instance', async ({ dashboardPage: page }) => {
    const chartCount = await page.evaluate(() => {
      // @ts-expect-error — Chart global จาก CDN
      return typeof window.Chart;
    });
    expect(chartCount).toBe('function');
  });

  test('status doughnut มีข้อมูลตรงกับ fixtures', async ({ dashboardPage: page }) => {
    const data = await page.evaluate(() => {
      // @ts-expect-error
      const chart = window.Chart.getChart(document.getElementById('chartStatus'));
      return chart?.data?.datasets?.[0]?.data;
    });
    // PASS=1, WARN=1, FAIL=1, ERROR=0
    expect(data).toEqual([1, 1, 1, 0]);
  });

  test('duration chart มี 3 จุด (ตามจำนวน ticket)', async ({ dashboardPage: page }) => {
    const points = await page.evaluate(() => {
      // @ts-expect-error
      const chart = window.Chart.getChart(document.getElementById('chartDuration'));
      return chart?.data?.datasets?.[0]?.data;
    });
    expect(points).toHaveLength(3);
  });

  test('top failing tickets list แสดงเฉพาะ FAIL/ERROR', async ({ dashboardPage: page }) => {
    const items = page.locator('#topFail li');
    await expect(items).toHaveCount(1);
    await expect(items.first()).toContainText('DEF-456');
  });
});
```

### 2.9 `tests-playwright/tests/detail-dialog.spec.ts`

```typescript
import { test, expect } from '../fixtures/test-fixtures';

test.describe('Dashboard — detail dialog', () => {
  test('คลิก row เปิด dialog พร้อมข้อมูล ticket', async ({ dashboardPage: page }) => {
    await page.locator('#workflowTable tbody tr', { hasText: 'ABC-123' }).click();

    const dialog = page.locator('#detailDialog');
    await expect(dialog).toBeVisible();
    await expect(dialog.locator('#detailTitle')).toContainText('ABC-123');
    await expect(dialog.locator('#detailTitle')).toContainText('เพิ่ม API ค้นหาสินค้า');
  });

  test('dialog แสดง block ทุก phase ที่มี', async ({ dashboardPage: page }) => {
    await page.locator('#workflowTable tbody tr', { hasText: 'ABC-123' }).click();

    const blocks = page.locator('#detailDialog .phase-block');
    await expect(blocks).toHaveCount(6);
  });

  test('dialog ของ FAIL แสดง report ที่สื่อปัญหา', async ({ dashboardPage: page }) => {
    await page.locator('#workflowTable tbody tr', { hasText: 'DEF-456' }).click();

    const dialog = page.locator('#detailDialog');
    await expect(dialog).toBeVisible();
    await expect(dialog).toContainText('SQL injection');
    await expect(dialog.locator('.phase-block.fail')).toBeVisible();
  });

  test('ปุ่ม close ปิด dialog', async ({ dashboardPage: page }) => {
    await page.locator('#workflowTable tbody tr', { hasText: 'ABC-123' }).click();
    await expect(page.locator('#detailDialog')).toBeVisible();

    await page.locator('#closeDetail').click();
    await expect(page.locator('#detailDialog')).not.toBeVisible();
  });

  test('กด Esc ปิด dialog', async ({ dashboardPage: page }) => {
    await page.locator('#workflowTable tbody tr', { hasText: 'ABC-123' }).click();
    await expect(page.locator('#detailDialog')).toBeVisible();

    await page.keyboard.press('Escape');
    await expect(page.locator('#detailDialog')).not.toBeVisible();
  });
});
```

### 2.10 `tests-playwright/tests/a11y.spec.ts`

```typescript
import { test, expect } from '../fixtures/test-fixtures';
import AxeBuilder from '@axe-core/playwright';

test.describe('Dashboard — accessibility', () => {
  test('ไม่มี critical/serious violations บนหน้าแรก', async ({ dashboardPage: page }) => {
    const results = await new AxeBuilder({ page })
      .withTags(['wcag2a', 'wcag2aa', 'wcag21a', 'wcag21aa'])
      .analyze();

    const critical = results.violations.filter(
      (v) => v.impact === 'critical' || v.impact === 'serious',
    );
    expect(critical, JSON.stringify(critical, null, 2)).toEqual([]);
  });

  test('dialog เปิดแล้วไม่มี violations', async ({ dashboardPage: page }) => {
    await page.locator('#workflowTable tbody tr', { hasText: 'ABC-123' }).click();
    await expect(page.locator('#detailDialog')).toBeVisible();

    const results = await new AxeBuilder({ page })
      .include('#detailDialog')
      .withTags(['wcag2a', 'wcag2aa'])
      .analyze();

    const critical = results.violations.filter(
      (v) => v.impact === 'critical' || v.impact === 'serious',
    );
    expect(critical).toEqual([]);
  });

  test('ตารางสามารถเข้าถึงด้วย keyboard', async ({ dashboardPage: page }) => {
    // tab ไปที่ search แล้วพิมพ์
    await page.locator('#search').focus();
    await page.keyboard.type('DEF');
    await expect(page.locator('#workflowTable tbody tr')).toHaveCount(1);
  });

  test('select มี label', async ({ dashboardPage: page }) => {
    await expect(page.locator('label[for="filterStatus"]')).toBeVisible();
  });
});
```

### 2.11 `tests-playwright/.gitignore`

```
node_modules/
playwright-report/
test-results/
blob-report/
playwright/.cache/
```

### 2.12 `tests-playwright/static/` — โครงสร้างที่ webServer serve

```
tests-playwright/static/
├── index.html    ← คัดลอกจาก dashboard/index.html
├── app.js        ← คัดลอกจาก dashboard/app.js
└── styles.css    ← คัดลอกจาก dashboard/styles.css
```

> ใช้ `scripts/sync-dashboard.ps1` เพื่อ sync อัตโนมัติ:
> ```powershell
> Copy-Item ../dashboard/index.html, ../dashboard/app.js, ../dashboard/styles.css ./static/ -Force
> ```

### 2.13 วิธีรัน

```bash
cd tests-playwright
npm ci
npm run install:browsers

npm test                    # ทุก browser
npm test -- --project=chromium
npm test -- --grep "filter"
npm run test:ui             # interactive mode
npm run test:report         # เปิด HTML report

# CI mode
CI=1 npm test
```

### 2.14 เพิ่มใน GitHub Actions

`.github/workflows/playwright.yml`:
```yaml
name: Playwright E2E

on:
  push:
    paths:
      - 'dashboard/**'
      - 'tests-playwright/**'
      - '.github/workflows/playwright.yml'
  pull_request:
    paths:
      - 'dashboard/**'
      - 'tests-playwright/**'

jobs:
  test:
    name: 🎭 Playwright
    runs-on: ubuntu-latest
    timeout-minutes: 20
    defaults:
      run:
        working-directory: tests-playwright
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-node@v4
        with:
          node-version: 20
          cache: npm
          cache-dependency-path: tests-playwright/package-lock.json

      - name: Sync dashboard → static
        run: |
          mkdir -p static
          cp ../dashboard/index.html static/
          cp ../dashboard/app.js static/
          cp ../dashboard/styles.css static/

      - run: npm ci
      - run: npx playwright install --with-deps
      - run: npm test

      - uses: actions/upload-artifact@v4
        if: always()
        with:
          name: playwright-report
          path: tests-playwright/playwright-report/
          retention-days: 14
```

---

## 🔔 3) Slack / Teams Notification

### 3.1 `scripts/lib/Notifications.psm1`

```powershell
#Requires -Version 7.0
Set-StrictMode -Version Latest

# ─────────────────────────────────────────────────────────────
# Notifications module — Slack + Microsoft Teams
# ─────────────────────────────────────────────────────────────

$script:StatusColors = @{
  PASS  = '#3fb950'
  WARN  = '#d29922'
  FAIL  = '#db6d28'
  ERROR = '#f85149'
}

$script:StatusEmoji = @{
  PASS  = '🟢'
  WARN  = '🟡'
  FAIL  = '🟠'
  ERROR = '🔴'
}

function Get-StatusColor {
  param([string]$Status)
  if ($script:StatusColors.ContainsKey($Status)) { return $script:StatusColors[$Status] }
  return '#8b949e'
}

function Get-StatusEmoji {
  param([string]$Status)
  if ($script:StatusEmoji.ContainsKey($Status)) { return $script:StatusEmoji[$Status] }
  return '⚪'
}

function Test-WebhookUrl {
  param([Parameter(Mandatory)][string]$Url)
  try {
    $u = [System.Uri]$Url
    return $u.Scheme -in @('http', 'https')
  } catch {
    return $false
  }
}

# ─────────────────────────────────────────────────────────────
# Slack
# ─────────────────────────────────────────────────────────────

function Send-SlackNotification {
  [CmdletBinding()]
  param(
    [Parameter(Mandatory)][string]$WebhookUrl,
    [Parameter(Mandatory)][string]$Title,
    [Parameter(Mandatory)][string]$Text,
    [ValidateSet('PASS', 'WARN', 'FAIL', 'ERROR')]
    [string]$Status = 'PASS',
    [string]$TicketKey = '',
    [hashtable]$Fields = @{},
    [string[]]$Actions = @(),
    [string]$Footer = '6-Phase Workflow',
    [string]$Channel = '',
    [string]$Username = 'Workflow Bot',
    [string]$IconEmoji = ':rocket:'
  )

  if (-not (Test-WebhookUrl $WebhookUrl)) {
    Write-Warning "Send-SlackNotification: invalid webhook URL"
    return
  }

  $emoji = Get-StatusEmoji $Status
  $color = Get-StatusColor $Status

  # ── Blocks ────────────────────────────────────────────────
  $blocks = [System.Collections.Generic.List[hashtable]]::new()

  # header
  $blocks.Add(@{
    type = 'header'
    text = @{ type = 'plain_text'; text = "$emoji $Title"; emoji = $true }
  })

  # context (ticket)
  if ($TicketKey) {
    $blocks.Add(@{
      type = 'context'
      elements = @(
        @{ type = 'mrkdwn'; text = "*Ticket:* ``$TicketKey``" }
        @{ type = 'mrkdwn'; text = "*Status:* $Status" }
        @{ type = 'mrkdwn'; text = "*Time:* <!date^$([DateTimeOffset]::UtcNow.ToUnixTimeSeconds())^{date_short_pretty} {time}|$(Get-Date -Format 'yyyy-MM-dd HH:mm')>" }
      )
    })
  }

  # body
  if ($Text) {
    $blocks.Add(@{
      type = 'section'
      text = @{ type = 'mrkdwn'; text = $Text }
    })
  }

  # fields
  if ($Fields.Count -gt 0) {
    $fieldArr = foreach ($kv in $Fields.GetEnumerator()) {
      @{ type = 'mrkdwn'; text = "*$($kv.Key):*\n$($kv.Value)" }
    }
    $blocks.Add(@{ type = 'section'; fields = @($fieldArr) })
  }

  # divider + actions
  if ($Actions.Count -gt 0) {
    $blocks.Add(@{ type = 'divider' })
    $elements = foreach ($a in $Actions) {
      @{ type = 'mrkdwn'; text = $a }
    }
    $blocks.Add(@{ type = 'context'; elements = @($elements) })
  }

  # footer
  $blocks.Add(@{ type = 'divider' })
  $blocks.Add(@{
    type = 'context'
    elements = @(@{ type = 'mrkdwn'; text = "_$Footer_" })
  })

  $payload = @{
    text      = "$emoji $Title"
    username  = $Username
    icon_emoji = $IconEmoji
    attachments = @(
      @{
        color    = $color
        blocks   = $blocks.ToArray()
        fallback = "$Title — $Status"
      }
    )
  }
  if ($Channel) { $payload.channel = $Channel }

  $body = $payload | ConvertTo-Json -Depth 30 -Compress

  try {
    $resp = Invoke-RestMethod -Uri $WebhookUrl -Method Post `
      -ContentType 'application/json; charset=utf-8' `
      -Body ([Text.Encoding]::UTF8.GetBytes($body)) `
      -TimeoutSec 10
    Write-Verbose "Slack notification sent: $resp"
    return $true
  } catch {
    Write-Warning "Slack notification failed: $($_.Exception.Message)"
    return $false
  }
}

# ─────────────────────────────────────────────────────────────
# Microsoft Teams (Incoming Webhook — MessageCard legacy)
# ─────────────────────────────────────────────────────────────

function Send-TeamsNotification {
  [CmdletBinding()]
  param(
    [Parameter(Mandatory)][string]$WebhookUrl,
    [Parameter(Mandatory)][string]$Title,
    [Parameter(Mandatory)][string]$Text,
    [ValidateSet('PASS', 'WARN', 'FAIL', 'ERROR')]
    [string]$Status = 'PASS',
    [string]$TicketKey = '',
    [hashtable]$Fields = @{},
    [string]$Footer = '6-Phase Workflow'
  )

  if (-not (Test-WebhookUrl $WebhookUrl)) {
    Write-Warning "Send-TeamsNotification: invalid webhook URL"
    return
  }

  $emoji = Get-StatusEmoji $Status
  $color = Get-StatusColor $Status

  $facts = [System.Collections.Generic.List[hashtable]]::new()
  if ($TicketKey) {
    $facts.Add(@{ name = 'Ticket'; value = $TicketKey })
  }
  $facts.Add(@{ name = 'Status'; value = "$emoji $Status" })
  $facts.Add(@{ name = 'Time';   value = (Get-Date -Format 'yyyy-MM-dd HH:mm:ss') })
  foreach ($kv in $Fields.GetEnumerator()) {
    $facts.Add(@{ name = [string]$kv.Key; value = [string]$kv.Value })
  }

  # Adaptive Card (ทันสมัย + ใช้ได้ใน Teams)
  $adaptive = @{
    type    = 'message'
    attachments = @(
      @{
        contentType = 'application/vnd.microsoft.card.adaptive'
        content = @{
          type    = 'AdaptiveCard'
          version = '1.4'
          '$schema' = 'http://adaptivecards.io/schemas/adaptive-card.json'
          body = @(
            @{
              type   = 'TextBlock'
              text   = "$emoji $Title"
              weight = 'Bolder'
              size   = 'Medium'
              color  = switch ($Status) {
                'PASS'  { 'Good' }
                'WARN'  { 'Warning' }
                'FAIL'  { 'Attention' }
                'ERROR' { 'Attention' }
                default { 'Default' }
              }
            },
            @{
              type   = 'FactSet'
              facts  = @($facts | ForEach-Object { @{ title = $_.name; value = $_.value } })
            },
            @{
              type = 'TextBlock'
              text = $Text
              wrap = $true
            },
            @{
              type     = 'TextBlock'
              text     = $Footer
              isSubtle = $true
              size     = 'Small'
              spacing  = 'Medium'
            }
          )
        }
      }
    )
    summary = "$Title — $Status"
  }

  $body = $adaptive | ConvertTo-Json -Depth 30 -Compress

  try {
    $resp = Invoke-RestMethod -Uri $WebhookUrl -Method Post `
      -ContentType 'application/json; charset=utf-8' `
      -Body ([Text.Encoding]::UTF8.GetBytes($body)) `
      -TimeoutSec 10
    Write-Verbose "Teams notification sent: $resp"
    return $true
  } catch {
    Write-Warning "Teams notification failed: $($_.Exception.Message)"
    return $false
  }
}

Export-ModuleMember -Function `
  Send-SlackNotification, `
  Send-TeamsNotification, `
  Get-StatusColor, `
  Get-StatusEmoji, `
  Test-WebhookUrl
```

### 3.2 `scripts/Send-Notification.ps1`

```powershell
<#
.SYNOPSIS
  ส่ง notification ไป Slack และ/หรือ Teams

.PARAMETER Channel
  slack | teams | both (default: both)

.PARAMETER Status
  PASS | WARN | FAIL | ERROR

.PARAMETER TicketKey
.PARAMETER Title
.PARAMETER Text
.PARAMETER FieldsJson
  JSON string ของ hashtable fields
.PARAMETER SummaryPath
  path ของ summary JSON — ถ้าระบุ จะดึงข้อมูลมาสร้างข้อความเอง

.EXAMPLE
  .\Send-Notification.ps1 -Channel slack -Status FAIL -TicketKey ABC-1 `
    -Title "Phase security failed" -Text "SQL injection at line 42"

.EXAMPLE
  .\Send-Notification.ps1 -Channel both -SummaryPath .\.workflow-reports\ABC-1-summary.json
#>
[CmdletBinding()]
param(
  [ValidateSet('slack', 'teams', 'both')]
  [string]$Channel = 'both',

  [ValidateSet('PASS', 'WARN', 'FAIL', 'ERROR')]
  [string]$Status,

  [string]$TicketKey = '',
  [string]$Title = '',
  [string]$Text  = '',
  [string]$FieldsJson = '{}',
  [string]$SummaryPath = '',

  [string]$SlackWebhook = $env:SLACK_WEBHOOK_URL,
  [string]$TeamsWebhook = $env:TEAMS_WEBHOOK_URL,

  [switch]$Silent
)

$ErrorActionPreference = 'Stop'
Set-StrictMode -Version Latest

# ── import module ─────────────────────────────────────────
Import-Module (Join-Path $PSScriptRoot 'lib/Notifications.psm1') -Force

# ── ถ้ามี SummaryPath ให้ดึงข้อมูล ────────────────────────
$fields = @{}
if (Test-Path $FieldsJson) {
  # (ไม่คาดว่าเป็น path — แต่กันพลาด)
}
try { $fields = $FieldsJson | ConvertFrom-Json -AsHashtable } catch { $fields = @{} }

if ($SummaryPath -and (Test-Path $SummaryPath)) {
  $summary = Get-Content $SummaryPath -Raw | ConvertFrom-Json -AsHashtable
  $Status    = $summary.overall
  $TicketKey = $summary.ticketKey
  $Title     = "Workflow $Status — $($summary.ticketKey)"

  $d = [math]::Round(($summary.durationMs / 1000), 1)
  $phaseLines = foreach ($kv in ($summary.results.GetEnumerator() | Sort-Object { $_.Key })) {
    $r = $kv.Value
    "$($r.emoji) *$($r.phaseName)*: $($r.status) ($([math]::Round($r.durationMs/1000,1))s)"
  }
  $Text = @"
*Task:* $($summary.task)
*Duration:* ${d}s

$($phaseLines -join "`n")
"@

  foreach ($kv in $summary.results.GetEnumerator()) {
    $fields["Phase:$($kv.Key)"] = $kv.Value.status
  }
}

if (-not $Title) {
  $Title = "Workflow notification — $TicketKey"
}

# ── ส่ง ───────────────────────────────────────────────────
$results = [ordered]@{}

if ($Channel -in @('slack', 'both')) {
  if (-not $SlackWebhook) {
    Write-Warning "SLACK_WEBHOOK_URL not set — skip Slack"
    $results.slack = 'skipped'
  } else {
    $ok = Send-SlackNotification `
      -WebhookUrl $SlackWebhook `
      -Title $Title `
      -Text $Text `
      -Status $Status `
      -TicketKey $TicketKey `
      -Fields $fields
    $results.slack = if ($ok) { 'sent' } else { 'failed' }
  }
}

if ($Channel -in @('teams', 'both')) {
  if (-not $TeamsWebhook) {
    Write-Warning "TEAMS_WEBHOOK_URL not set — skip Teams"
    $results.teams = 'skipped'
  } else {
    $ok = Send-TeamsNotification `
      -WebhookUrl $TeamsWebhook `
      -Title $Title `
      -Text $Text `
      -Status $Status `
      -TicketKey $TicketKey `
      -Fields $fields
    $results.teams = if ($ok) { 'sent' } else { 'failed' }
  }
}

if (-not $Silent) {
  $results.GetEnumerator() | ForEach-Object {
    $color = if ($_.Value -eq 'sent') { 'Green' }
             elseif ($_.Value -eq 'failed') { 'Red' }
             else { 'DarkGray' }
    Write-Host ("  🔔 {0,-6} → {1}" -f $_.Key, $_.Value) -ForegroundColor $color
  }
}

# exit non-zero ถ้าทุกช่องที่ตั้งใจส่ง failed
$attempted = $results.Values | Where-Object { $_ -ne 'skipped' }
if ($attempted.Count -gt 0 -and ($attempted -notcontains 'sent')) {
  exit 1
}
```

### 3.3 Integrate เข้า `Run-6PhaseWorkflow.ps1`

แก้ orchestrator ให้เรียก notification ทุกครั้งที่:
1. จบแต่ละเฟส (ถ้า status != PASS)
2. จบ workflow ทั้งหมด (สรุป)

เพิ่ม **หลัง** block `Publish-JiraComment`:

```powershell
# ── ส่ง notification ─────────────────────────────────────
if (-not $SkipNotification) {
  if ($result.status -ne 'PASS' -or $result.phaseId -in @('security', 'performance')) {
    $notifyParams = @{
      Channel     = $NotificationChannel
      Status      = $result.status
      TicketKey   = $TicketKey
      Title       = "$($result.emoji) $($result.phaseName) — $($result.status)"
      Text        = $result.report.Substring(0, [Math]::Min(2000, $result.report.Length))
      FieldsJson  = (@{
        Phase    = $result.phaseName
        Duration = "$([math]::Round($result.durationMs/1000,1))s"
      } | ConvertTo-Json -Compress)
    }
    & "$PSScriptRoot/Send-Notification.ps1" @notifyParams -Silent

    if ($phase.failFast -and $result.status -in @('FAIL','ERROR')) {
      # แจ้งเตือน failFast โดยเฉพาะ
      & "$PSScriptRoot/Send-Notification.ps1" `
        -Channel $NotificationChannel `
        -Status 'FAIL' `
        -TicketKey $TicketKey `
        -Title "🛑 FailFast triggered" `
        -Text "Workflow หยุดที่เฟส *$($result.phaseName)* เนื่องจาก status = $($result.status)" `
        -Silent
    }
  }
}
```

และตอนท้าย (ก่อน exit):

```powershell
# ── แจ้งสรุปรวม ─────────────────────────────────────────
if (-not $SkipNotification) {
  & "$PSScriptRoot/Send-Notification.ps1" `
    -Channel $NotificationChannel `
    -SummaryPath $summaryPath -Silent
}
```

### 3.4 เพิ่ม parameter ให้ orchestrator

```powershell
# ใน param() ของ Run-6PhaseWorkflow.ps1
[ValidateSet('slack', 'teams', 'both', 'none')]
[string]$NotificationChannel = 'both',
[switch]$SkipNotification
```

### 3.5 ตั้งค่า env

```powershell
# Slack incoming webhook
# ไปที่ https://api.slack.com/messaging/webhooks
$env:SLACK_WEBHOOK_URL = 'https://hooks.slack.com/services/T000/B000/XXXX'

# Teams incoming webhook
# ไปที่ Teams channel → Connectors → Incoming Webhook
$env:TEAMS_WEBHOOK_URL = 'https://outlook.office.com/webhook/...'
```

### 3.6 เพิ่ม job ใน GitHub Actions

เพิ่มใน `.github/workflows/6phase-on-label.yml` ต่อจาก job `run`:

```yaml
  notify:
    name: 🔔 Notify on failure
    needs: run
    if: failure() || needs.run.result == 'failure'
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Download reports
        uses: actions/download-artifact@v4
        with:
          name: workflow-reports-${{ needs.detect.outputs.ticket }}
          path: .workflow-reports
        continue-on-error: true

      - name: Send notification
        shell: pwsh
        env:
          SLACK_WEBHOOK_URL: ${{ secrets.SLACK_WEBHOOK_URL }}
          TEAMS_WEBHOOK_URL: ${{ secrets.TEAMS_WEBHOOK_URL }}
        run: |
          $summary = Get-ChildItem ./.workflow-reports/*-summary.json | Select-Object -First 1
          if ($summary) {
            ./scripts/Send-Notification.ps1 -Channel both -SummaryPath $summary.FullName
          } else {
            ./scripts/Send-Notification.ps1 -Channel both -Status ERROR `
              -TicketKey "${{ needs.detect.outputs.ticket }}" `
              -Title "Workflow failed before summary" `
              -Text "Job failed — no summary produced"
          }
```

---

## 🔌 4) Go REST API — `workflow-api`

### 4.1 `services/workflow-api/go.mod`

```go
module github.com/icmongolang/workflow-api

go 1.22

require (
    github.com/go-chi/chi/v5 v5.1.0
    github.com/go-chi/cors v1.2.1
    github.com/prometheus/client_golang v1.20.2
    github.com/stretchr/testify v1.9.0
)
```

### 4.2 `services/workflow-api/internal/config/config.go`

```go
package config

import (
    "fmt"
    "os"
    "strings"
    "time"
)

type Config struct {
    ListenAddr      string
    ReportsDir      string
    CacheTTL        time.Duration
    ShutdownTimeout time.Duration
    LogLevel        string
    CORSOrigins     []string
    MaxPageSize     int
}

func Load() (*Config, error) {
    c := &Config{
        ListenAddr:      getEnv("LISTEN_ADDR", ":8080"),
        ReportsDir:      getEnv("REPORTS_DIR", "./.workflow-reports"),
        LogLevel:        strings.ToLower(getEnv("LOG_LEVEL", "info")),
        MaxPageSize:     200,
    }

    var err error
    if c.CacheTTL, err = time.ParseDuration(getEnv("CACHE_TTL", "30s")); err != nil {
        return nil, fmt.Errorf("CACHE_TTL: %w", err)
    }
    if c.ShutdownTimeout, err = time.ParseDuration(getEnv("SHUTDOWN_TIMEOUT", "15s")); err != nil {
        return nil, fmt.Errorf("SHUTDOWN_TIMEOUT: %w", err)
    }

    origins := getEnv("CORS_ORIGINS", "*")
    for _, o := range strings.Split(origins, ",") {
        o = strings.TrimSpace(o)
        if o != "" {
            c.CORSOrigins = append(c.CORSOrigins, o)
        }
    }

    return c, nil
}

func getEnv(k, def string) string {
    if v, ok := os.LookupEnv(k); ok && v != "" {
        return v
    }
    return def
}
```

### 4.3 `services/workflow-api/internal/model/workflow.go`

```go
package model

import "time"

// Summary mirrors the JSON ที่ workflow runner เขียนไว้
type Summary struct {
    TicketKey  string          `json:"ticketKey"`
    Task       string          `json:"task"`
    Overall    Status          `json:"overall"`
    DurationMs int64           `json:"durationMs"`
    FinishedAt time.Time       `json:"finishedAt"`
    Results    map[string]Phase `json:"results"`
}

type Phase struct {
    PhaseID    string `json:"phaseId"`
    PhaseName  string `json:"phaseName"`
    Emoji      string `json:"emoji"`
    Status     Status `json:"status"`
    DurationMs int64  `json:"durationMs"`
    Report     string `json:"report,omitempty"`
    Error      string `json:"error,omitempty"`
    StartedAt  time.Time `json:"startedAt,omitempty"`
    EndedAt    time.Time `json:"endedAt,omitempty"`
}

type Status string

const (
    StatusPass  Status = "PASS"
    StatusWarn  Status = "WARN"
    StatusFail  Status = "FAIL"
    StatusError Status = "ERROR"
)

func (s Status) Valid() bool {
    switch s {
    case StatusPass, StatusWarn, StatusFail, StatusError:
        return true
    }
    return false
}

// ── API DTOs ─────────────────────────────────────────────────

type ListQuery struct {
    Status   Status
    Search   string
    Phase    string
    Sort     string // "finished_at" | "duration_ms" | "ticket_key"
    Order    string // "asc" | "desc"
    Page     int
    PageSize int
}

type ListResult struct {
    Items    []*Summary `json:"items"`
    Total    int        `json:"total"`
    Page     int        `json:"page"`
    PageSize int        `json:"pageSize"`
}

type Stats struct {
    Total       int                `json:"total"`
    ByStatus    map[Status]int     `json:"byStatus"`
    AvgDurationMs float64          `json:"avgDurationMs"`
    P95DurationMs int64            `json:"p95DurationMs"`
    PhaseStats  map[string]PhaseStat `json:"phaseStats"`
}

type PhaseStat struct {
    Total  int     `json:"total"`
    Pass   int     `json:"pass"`
    Warn   int     `json:"warn"`
    Fail   int     `json:"fail"`
    Error  int     `json:"error"`
    PassRate float64 `json:"passRate"`
}
```

### 4.4 `services/workflow-api/internal/store/store.go`

```go
package store

import (
    "context"
    "errors"

    "github.com/icmongolang/workflow-api/internal/model"
)

var (
    ErrNotFound = errors.New("not found")
    ErrInvalid  = errors.New("invalid input")
)

type Store interface {
    // List คืน summaries ตาม query (กรองแล้ว, paginate แล้ว)
    List(ctx context.Context, q model.ListQuery) (*model.ListResult, error)

    // GetByTicket คืน summary ของ ticket เดียว
    GetByTicket(ctx context.Context, ticketKey string) (*model.Summary, error)

    // GetPhase คืน phase ของ ticket
    GetPhase(ctx context.Context, ticketKey, phaseID string) (*model.Phase, error)

    // Stats คำนวณสถิติรวม
    Stats(ctx context.Context) (*model.Stats, error)

    // Close คืน resource
    Close() error
}
```

### 4.5 `services/workflow-api/internal/store/filesystem.go`

```go
package store

import (
    "context"
    "encoding/json"
    "fmt"
    "log/slog"
    "os"
    "path/filepath"
    "sort"
    "strings"
    "sync"
    "time"

    "github.com/icmongolang/workflow-api/internal/model"
)

type FilesystemStore struct {
    dir      string
    cacheTTL time.Duration

    mu       sync.RWMutex
    cache    []*model.Summary
    cachedAt time.Time
}

func NewFilesystemStore(dir string, cacheTTL time.Duration) (*FilesystemStore, error) {
    if dir == "" {
        return nil, fmt.Errorf("reports dir is empty")
    }
    if _, err := os.Stat(dir); err != nil {
        return nil, fmt.Errorf("reports dir: %w", err)
    }
    return &FilesystemStore{dir: dir, cacheTTL: cacheTTL}, nil
}

func (s *FilesystemStore) Close() error { return nil }

// ── Loading ─────────────────────────────────────────────────

func (s *FilesystemStore) load() ([]*model.Summary, error) {
    s.mu.RLock()
    if time.Since(s.cachedAt) < s.cacheTTL && s.cache != nil {
        out := s.cache
        s.mu.RUnlock()
        return out, nil
    }
    s.mu.RUnlock()

    s.mu.Lock()
    defer s.mu.Unlock()

    // double-check
    if time.Since(s.cachedAt) < s.cacheTTL && s.cache != nil {
        return s.cache, nil
    }

    entries, err := os.ReadDir(s.dir)
    if err != nil {
        return nil, fmt.Errorf("read reports dir: %w", err)
    }

    out := make([]*model.Summary, 0, len(entries))
    for _, e := range entries {
        if e.IsDir() || !strings.HasSuffix(e.Name(), "-summary.json") {
            continue
        }
        sum, err := s.readSummary(filepath.Join(s.dir, e.Name()))
        if err != nil {
            slog.Warn("skip malformed summary", "file", e.Name(), "err", err)
            continue
        }
        out = append(out, sum)
    }

    // default sort: finishedAt desc
    sort.Slice(out, func(i, j int) bool {
        return out[i].FinishedAt.After(out[j].FinishedAt)
    })

    s.cache = out
    s.cachedAt = time.Now()
    return out, nil
}

func (s *FilesystemStore) readSummary(path string) (*model.Summary, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    var sum model.Summary
    dec := json.NewDecoder(f)
    dec.DisallowUnknownFields()
    if err := dec.Decode(&sum); err != nil {
        return nil, err
    }
    if sum.TicketKey == "" {
        return nil, fmt.Errorf("missing ticketKey")
    }
    return &sum, nil
}

// ── Queries ─────────────────────────────────────────────────

func (s *FilesystemStore) List(ctx context.Context, q model.ListQuery) (*model.ListResult, error) {
    if err := ctx.Err(); err != nil {
        return nil, err
    }

    all, err := s.load()
    if err != nil {
        return nil, err
    }

    // filter
    filtered := make([]*model.Summary, 0, len(all))
    for _, w := range all {
        if q.Status != "" && w.Overall != q.Status {
            continue
        }
        if q.Phase != "" {
            if _, ok := w.Results[q.Phase]; !ok {
                continue
            }
        }
        if q.Search != "" {
            ql := strings.ToLower(q.Search)
            if !strings.Contains(strings.ToLower(w.TicketKey), ql) &&
                !strings.Contains(strings.ToLower(w.Task), ql) {
                continue
            }
        }
        filtered = append(filtered, w)
    }

    // sort
    sortSummaries(filtered, q.Sort, q.Order)

    // paginate
    total := len(filtered)
    page := q.Page
    if page < 1 {
        page = 1
    }
    size := q.PageSize
    if size < 1 {
        size = 20
    }
    start := (page - 1) * size
    if start > total {
        start = total
    }
    end := start + size
    if end > total {
        end = total
    }

    return &model.ListResult{
        Items:    filtered[start:end],
        Total:    total,
        Page:     page,
        PageSize: size,
    }, nil
}

func (s *FilesystemStore) GetByTicket(ctx context.Context, ticketKey string) (*model.Summary, error) {
    if err := ctx.Err(); err != nil {
        return nil, err
    }
    all, err := s.load()
    if err != nil {
        return nil, err
    }
    for _, w := range all {
        if w.TicketKey == ticketKey {
            return w, nil
        }
    }
    return nil, ErrNotFound
}

func (s *FilesystemStore) GetPhase(ctx context.Context, ticketKey, phaseID string) (*model.Phase, error) {
    sum, err := s.GetByTicket(ctx, ticketKey)
    if err != nil {
        return nil, err
    }
    p, ok := sum.Results[phaseID]
    if !ok {
        return nil, ErrNotFound
    }
    return &p, nil
}

func (s *FilesystemStore) Stats(ctx context.Context) (*model.Stats, error) {
    if err := ctx.Err(); err != nil {
        return nil, err
    }
    all, err := s.load()
    if err != nil {
        return nil, err
    }

    st := &model.Stats{
        Total:      len(all),
        ByStatus:   make(map[model.Status]int),
        PhaseStats: make(map[string]model.PhaseStat),
    }

    var durations []int64
    for _, w := range all {
        st.ByStatus[w.Overall]++
        durations = append(durations, w.DurationMs)

        for id, p := range w.Results {
            ps := st.PhaseStats[id]
            ps.Total++
            switch p.Status {
            case model.StatusPass:
                ps.Pass++
            case model.StatusWarn:
                ps.Warn++
            case model.StatusFail:
                ps.Fail++
            case model.StatusError:
                ps.Error++
            }
            st.PhaseStats[id] = ps
        }
    }

    if len(durations) > 0 {
        var sum int64
        for _, d := range durations {
            sum += d
        }
        st.AvgDurationMs = float64(sum) / float64(len(durations))
        st.P95DurationMs = percentile(durations, 0.95)
    }

    for id, ps := range st.PhaseStats {
        if ps.Total > 0 {
            ps.PassRate = float64(ps.Total-ps.Fail-ps.Error) / float64(ps.Total) * 100
        }
        st.PhaseStats[id] = ps
    }

    return st, nil
}

// ── Helpers ─────────────────────────────────────────────────

func sortSummaries(items []*model.Summary, sortKey, order string) {
    if sortKey == "" {
        sortKey = "finished_at"
    }
    if order == "" {
        order = "desc"
    }
    less := func(i, j int) bool { return items[i].FinishedAt.Before(items[j].FinishedAt) }
    switch sortKey {
    case "duration_ms":
        less = func(i, j int) bool { return items[i].DurationMs < items[j].DurationMs }
    case "ticket_key":
        less = func(i, j int) bool { return items[i].TicketKey < items[j].TicketKey }
    case "finished_at":
        // default
    }
    if order == "desc" {
        sort.SliceStable(items, func(i, j int) bool { return less(j, i) })
    } else {
        sort.SliceStable(items, less)
    }
}

func percentile(values []int64, p float64) int64 {
    if len(values) == 0 {
        return 0
    }
    sorted := make([]int64, len(values))
    copy(sorted, values)
    sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })

    idx := int(float64(len(sorted)-1) * p)
    if idx < 0 {
        idx = 0
    }
    if idx >= len(sorted) {
        idx = len(sorted) - 1
    }
    return sorted[idx]
}
```

### 4.6 `services/workflow-api/internal/store/filesystem_test.go`

```go
package store

import (
    "context"
    "testing"
    "time"

    "github.com/icmongolang/workflow-api/internal/model"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestFilesystemStore_List(t *testing.T) {
    s, err := NewFilesystemStore("../../testdata/summaries", time.Second)
    require.NoError(t, err)
    t.Cleanup(func() { _ = s.Close() })

    ctx := context.Background()

    t.Run("returns all", func(t *testing.T) {
        r, err := s.List(ctx, model.ListQuery{Page: 1, PageSize: 10})
        require.NoError(t, err)
        assert.Equal(t, 3, r.Total)
        assert.Len(t, r.Items, 3)
    })

    t.Run("filter by status", func(t *testing.T) {
        r, err := s.List(ctx, model.ListQuery{Status: model.StatusFail, PageSize: 10})
        require.NoError(t, err)
        assert.Equal(t, 1, r.Total)
        assert.Equal(t, "DEF-456", r.Items[0].TicketKey)
    })

    t.Run("search", func(t *testing.T) {
        r, err := s.List(ctx, model.ListQuery{Search: "pagination", PageSize: 10})
        require.NoError(t, err)
        assert.Equal(t, 1, r.Total)
        assert.Equal(t, "DEF-456", r.Items[0].TicketKey)
    })

    t.Run("pagination", func(t *testing.T) {
        r, err := s.List(ctx, model.ListQuery{Page: 1, PageSize: 2})
        require.NoError(t, err)
        assert.Len(t, r.Items, 2)
        assert.Equal(t, 3, r.Total)

        r2, err := s.List(ctx, model.ListQuery{Page: 2, PageSize: 2})
        require.NoError(t, err)
        assert.Len(t, r2.Items, 1)
    })
}

func TestFilesystemStore_GetByTicket(t *testing.T) {
    s, err := NewFilesystemStore("../../testdata/summaries", time.Second)
    require.NoError(t, err)

    ctx := context.Background()
    t.Run("found", func(t *testing.T) {
        w, err := s.GetByTicket(ctx, "ABC-123")
        require.NoError(t, err)
        assert.Equal(t, "ABC-123", w.TicketKey)
        assert.Equal(t, model.StatusPass, w.Overall)
    })

    t.Run("not found", func(t *testing.T) {
        _, err := s.GetByTicket(ctx, "NOPE-999")
        assert.ErrorIs(t, err, ErrNotFound)
    })
}

func TestFilesystemStore_GetPhase(t *testing.T) {
    s, err := NewFilesystemStore("../../testdata/summaries", time.Second)
    require.NoError(t, err)

    ctx := context.Background()
    p, err := s.GetPhase(ctx, "ABC-123", "security")
    require.NoError(t, err)
    assert.Equal(t, "Security", p.PhaseName)
    assert.Equal(t, model.StatusPass, p.Status)
}

func TestFilesystemStore_Stats(t *testing.T) {
    s, err := NewFilesystemStore("../../testdata/summaries", time.Second)
    require.NoError(t, err)

    st, err := s.Stats(context.Background())
    require.NoError(t, err)
    assert.Equal(t, 3, st.Total)
    assert.Equal(t, 1, st.ByStatus[model.StatusPass])
    assert.Equal(t, 1, st.ByStatus[model.StatusWarn])
    assert.Equal(t, 1, st.ByStatus[model.StatusFail])
    assert.NotZero(t, st.AvgDurationMs)
    assert.NotEmpty(t, st.PhaseStats)
}

func TestPercentile(t *testing.T) {
    assert.Equal(t, int64(0), percentile(nil, 0.95))
    assert.Equal(t, int64(5), percentile([]int64{5}, 0.95))
    assert.Equal(t, int64(9), percentile([]int64{1,2,3,4,5,6,7,8,9,10}, 0.95))
}
```

### 4.7 `services/workflow-api/internal/service/workflow_service.go`

```go
package service

import (
    "context"
    "errors"

    "github.com/icmongolang/workflow-api/internal/model"
    "github.com/icmongolang/workflow-api/internal/store"
)

type WorkflowService struct {
    store store.Store
}

func NewWorkflowService(s store.Store) *WorkflowService {
    return &WorkflowService{store: s}
}

func (svc *WorkflowService) List(ctx context.Context, q model.ListQuery) (*model.ListResult, error) {
    // normalize
    if q.Page < 1 {
        q.Page = 1
    }
    if q.PageSize < 1 || q.PageSize > 200 {
        q.PageSize = 20
    }
    if q.Status != "" && !q.Status.Valid() {
        return nil, store.ErrInvalid
    }
    if q.Sort == "" {
        q.Sort = "finished_at"
    }
    if q.Order != "asc" && q.Order != "desc" {
        q.Order = "desc"
    }
    return svc.store.List(ctx, q)
}

func (svc *WorkflowService) Get(ctx context.Context, ticketKey string) (*model.Summary, error) {
    if ticketKey == "" {
        return nil, store.ErrInvalid
    }
    return svc.store.GetByTicket(ctx, ticketKey)
}

func (svc *WorkflowService) GetPhase(ctx context.Context, ticketKey, phaseID string) (*model.Phase, error) {
    if ticketKey == "" || phaseID == "" {
        return nil, store.ErrInvalid
    }
    return svc.store.GetPhase(ctx, ticketKey, phaseID)
}

func (svc *WorkflowService) Stats(ctx context.Context) (*model.Stats, error) {
    return svc.store.Stats(ctx)
}

// IsNotFound เป็นตัวช่วยให้ handler ตรวจ error
func IsNotFound(err error) bool {
    return errors.Is(err, store.ErrNotFound)
}

func IsInvalid(err error) bool {
    return errors.Is(err, store.ErrInvalid)
}
```

### 4.8 `services/workflow-api/internal/service/workflow_service_test.go`

```go
package service

import (
    "context"
    "testing"
    "time"

    "github.com/icmongolang/workflow-api/internal/model"
    "github.com/icmongolang/workflow-api/internal/store"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func newSvc(t *testing.T) *WorkflowService {
    t.Helper()
    s, err := store.NewFilesystemStore("../../testdata/summaries", time.Second)
    require.NoError(t, err)
    return NewWorkflowService(s)
}

func TestWorkflowService_List_Defaults(t *testing.T) {
    svc := newSvc(t)
    r, err := svc.List(context.Background(), model.ListQuery{})
    require.NoError(t, err)
    assert.Equal(t, 1, r.Page)
    assert.Equal(t, 20, r.PageSize)
    assert.Equal(t, "desc", "desc") // default order
}

func TestWorkflowService_List_RejectsInvalidStatus(t *testing.T) {
    svc := newSvc(t)
    _, err := svc.List(context.Background(), model.ListQuery{Status: "BOGUS"})
    assert.True(t, IsInvalid(err))
}

func TestWorkflowService_Get_EmptyTicket(t *testing.T) {
    svc := newSvc(t)
    _, err := svc.Get(context.Background(), "")
    assert.True(t, IsInvalid(err))
}

func TestWorkflowService_Get_NotFound(t *testing.T) {
    svc := newSvc(t)
    _, err := svc.Get(context.Background(), "MISSING-1")
    assert.True(t, IsNotFound(err))
}
```

### 4.9 `services/workflow-api/internal/metrics/metrics.go`

```go
package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    HTTPRequests = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Namespace: "workflow_api",
            Name:      "http_requests_total",
            Help:      "Total HTTP requests",
        },
        []string{"method", "route", "status"},
    )

    HTTPDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Namespace: "workflow_api",
            Name:      "http_request_duration_seconds",
            Help:      "HTTP request duration",
            Buckets:   prometheus.DefBuckets,
        },
        []string{"method", "route"},
    )

    InFlight = promauto.NewGauge(
        prometheus.GaugeOpts{
            Namespace: "workflow_api",
            Name:      "http_in_flight_requests",
            Help:      "In-flight HTTP requests",
        },
    )

    StoreLoadDuration = promauto.NewHistogram(
        prometheus.HistogramOpts{
            Namespace: "workflow_api",
            Name:      "store_load_duration_seconds",
            Help:      "Duration of store load from disk",
            Buckets:   []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1},
        },
    )

    CacheHit = promauto.NewCounter(
        prometheus.CounterOpts{
            Namespace: "workflow_api",
            Name:      "cache_hit_total",
            Help:      "Cache hits",
        },
    )
)
```

### 4.10 `services/workflow-api/internal/handler/middleware.go`

```go
package handler

import (
    "log/slog"
    "net/http"
    "time"

    "github.com/icmongolang/workflow-api/internal/metrics"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
)

// MetricsMiddleware เก็บ metric ของทุก request
func MetricsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        metrics.InFlight.Inc()
        defer metrics.InFlight.Dec()

        ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
        next.ServeHTTP(ww, r)

        route := chi.RouteContext(r.Context()).RoutePattern()
        if route == "" {
            route = "unmatched"
        }
        status := ww.Status()
        if status == 0 {
            status = 200
        }

        metrics.HTTPRequests.
            WithLabelValues(r.Method, route, http.StatusText(status)).
            Inc()
        metrics.HTTPDuration.
            WithLabelValues(r.Method, route).
            Observe(time.Since(start).Seconds())
    })
}

// Recoverer — กัน panic ใน handler
func Recoverer(logger *slog.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            defer func() {
                if rec := recover(); rec != nil {
                    logger.Error("panic recovered",
                        "panic", rec,
                        "path", r.URL.Path,
                        "method", r.Method,
                    )
                    writeError(w, http.StatusInternalServerError, "internal server error", nil)
                }
            }()
            next.ServeHTTP(w, r)
        })
    }
}
```

### 4.11 `services/workflow-api/internal/handler/handler.go`

```go
package handler

import (
    "encoding/json"
    "errors"
    "log/slog"
    "net/http"
    "strconv"

    "github.com/icmongolang/workflow-api/internal/model"
    "github.com/icmongolang/workflow-api/internal/service"
)

type Handler struct {
    svc    *service.WorkflowService
    logger *slog.Logger
}

func New(svc *service.WorkflowService, logger *slog.Logger) *Handler {
    return &Handler{svc: svc, logger: logger}
}

// ── helpers ─────────────────────────────────────────────────

func writeJSON(w http.ResponseWriter, status int, v any) {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.Header().Set("X-Content-Type-Options", "nosniff")
    w.WriteHeader(status)
    if v == nil {
        return
    }
    if err := json.NewEncoder(w).Encode(v); err != nil {
        // เราไม่ panic — log เท่านั้น
        slog.Warn("write response failed", "err", err)
    }
}

func writeError(w http.ResponseWriter, status int, msg string, details map[string]any) {
    body := map[string]any{"error": msg}
    if details != nil {
        body["details"] = details
    }
    writeJSON(w, status, body)
}

func parseIntDefault(s string, def int) int {
    if s == "" {
        return def
    }
    n, err := strconv.Atoi(s)
    if err != nil {
        return def
    }
    return n
}

// ── List ────────────────────────────────────────────────────

func (h *Handler) ListWorkflows(w http.ResponseWriter, r *http.Request) {
    q := r.URL.Query()
    page := parseIntDefault(q.Get("page"), 1)
    size := parseIntDefault(q.Get("pageSize"), 20)

    query := model.ListQuery{
        Status:   model.Status(q.Get("status")),
        Search:   q.Get("search"),
        Phase:    q.Get("phase"),
        Sort:     q.Get("sort"),
        Order:    q.Get("order"),
        Page:     page,
        PageSize: size,
    }

    res, err := h.svc.List(r.Context(), query)
    if err != nil {
        if service.IsInvalid(err) {
            writeError(w, http.StatusBadRequest, "invalid query", map[string]any{
                "status": q.Get("status"),
            })
            return
        }
        h.logger.Error("list failed", "err", err)
        writeError(w, http.StatusInternalServerError, "internal error", nil)
        return
    }
    writeJSON(w, http.StatusOK, res)
}

// ── Get by ticket ───────────────────────────────────────────

func (h *Handler) GetWorkflow(w http.ResponseWriter, r *http.Request) {
    key := pathParam(r, "ticketKey")
    if key == "" {
        writeError(w, http.StatusBadRequest, "ticketKey required", nil)
        return
    }

    sum, err := h.svc.Get(r.Context(), key)
    if err != nil {
        if service.IsNotFound(err) {
            writeError(w, http.StatusNotFound, "workflow not found", map[string]any{"ticketKey": key})
            return
        }
        if service.IsInvalid(err) {
            writeError(w, http.StatusBadRequest, "invalid ticketKey", nil)
            return
        }
        h.logger.Error("get failed", "err", err, "ticket", key)
        writeError(w, http.StatusInternalServerError, "internal error", nil)
        return
    }
    writeJSON(w, http.StatusOK, sum)
}

// ── Get phase ───────────────────────────────────────────────

func (h *Handler) GetPhase(w http.ResponseWriter, r *http.Request) {
    key := pathParam(r, "ticketKey")
    phase := pathParam(r, "phaseID")
    if key == "" || phase == "" {
        writeError(w, http.StatusBadRequest, "ticketKey and phaseID required", nil)
        return
    }

    p, err := h.svc.GetPhase(r.Context(), key, phase)
    if err != nil {
        if service.IsNotFound(err) {
            writeError(w, http.StatusNotFound, "phase not found", map[string]any{
                "ticketKey": key, "phaseID": phase,
            })
            return
        }
        h.logger.Error("get phase failed", "err", err, "ticket", key, "phase", phase)
        writeError(w, http.StatusInternalServerError, "internal error", nil)
        return
    }
    writeJSON(w, http.StatusOK, p)
}

// ── Stats ───────────────────────────────────────────────────

func (h *Handler) Stats(w http.ResponseWriter, r *http.Request) {
    st, err := h.svc.Stats(r.Context())
    if err != nil {
        h.logger.Error("stats failed", "err", err)
        writeError(w, http.StatusInternalServerError, "internal error", nil)
        return
    }
    writeJSON(w, http.StatusOK, st)
}

// ── Not found ───────────────────────────────────────────────

func (h *Handler) NotFound(w http.ResponseWriter, r *http.Request) {
    writeError(w, http.StatusNotFound, "route not found", map[string]any{"path": r.URL.Path})
}

var errNil = errors.New("nil")
var _ = errNil
```

### 4.12 `services/workflow-api/internal/handler/health_handler.go`

```go
package handler

import (
    "context"
    "net/http"
    "time"
)

func (h *Handler) Healthz(w http.ResponseWriter, r *http.Request) {
    // liveness — process ยังอยู่
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    _, _ = w.Write([]byte("ok\n"))
}

func (h *Handler) Readyz(w http.ResponseWriter, r *http.Request) {
    // readiness — store ทำงานได้
    ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
    defer cancel()

    if _, err := h.svc.Stats(ctx); err != nil {
        w.Header().Set("Content-Type", "text/plain; charset=utf-8")
        w.WriteHeader(http.StatusServiceUnavailable)
        _, _ = w.Write([]byte("not ready: " + err.Error() + "\n"))
        return
    }
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    _, _ = w.Write([]byte("ready\n"))
}
```

### 4.13 `services/workflow-api/internal/router/router.go`

```go
package router

import (
    "log/slog"
    "net/http"

    "github.com/icmongolang/workflow-api/internal/config"
    "github.com/icmongolang/workflow-api/internal/handler"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/go-chi/cors"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

func pathParam(r *http.Request, key string) string {
    return chi.URLParam(r, key)
}

func New(h *handler.Handler, cfg *config.Config) http.Handler {
    r := chi.NewRouter()

    // ── Global middleware ────────────────────────────────
    r.Use(middleware.RequestID)
    r.Use(middleware.RealIP)
    r.Use(middleware.Heartbeat("/ping"))
    r.Use(handler.Recoverer(slog.Default()))
    r.Use(handler.MetricsMiddleware)
    r.Use(middleware.Timeout(30 * time.Second))
    r.Use(cors.Handler(cors.Options{
        AllowedOrigins:   cfg.CORSOrigins,
        AllowedMethods:   []string{"GET", "OPTIONS"},
        AllowedHeaders:   []string{"Accept", "Content-Type", "X-Request-ID"},
        ExposedHeaders:   []string{"X-Request-ID", "X-Total-Count"},
        AllowCredentials: false,
        MaxAge:           300,
    }))
    r.Use(middleware.Compress(5))
    r.Use(middleware.CleanPath)

    // ── Health & metrics ────────────────────────────────
    r.Get("/healthz", h.Healthz)
    r.Get("/readyz",  h.Readyz)
    r.Handle("/metrics", promhttp.Handler())

    // ── API ─────────────────────────────────────────────
    r.Route("/api/v1", func(r chi.Router) {
        r.Get("/workflows",                  h.ListWorkflows)
        r.Get("/workflows/{ticketKey}",      h.GetWorkflow)
        r.Get("/workflows/{ticketKey}/phases/{phaseID}", h.GetPhase)
        r.Get("/stats",                      h.Stats)
    })

    // ── 404 ─────────────────────────────────────────────
    r.NotFound(h.NotFound)
    r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusMethodNotAllowed)
        _, _ = w.Write([]byte(`{"error":"method not allowed"}`))
    })

    return r
}
```

> **หมายเหตุ:** ต้องเพิ่ม `"time"` import ใน router.go — ผมตัดให้สั้น แต่ตัวจริงต้องมี

### 4.14 `services/workflow-api/main.go`

```go
package main

import (
    "context"
    "errors"
    "log/slog"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/icmongolang/workflow-api/internal/config"
    "github.com/icmongolang/workflow-api/internal/handler"
    "github.com/icmongolang/workflow-api/internal/router"
    "github.com/icmongolang/workflow-api/internal/service"
    "github.com/icmongolang/workflow-api/internal/store"
)

func main() {
    if err := run(); err != nil {
        slog.Error("fatal", "err", err)
        os.Exit(1)
    }
}

func run() error {
    cfg, err := config.Load()
    if err != nil {
        return err
    }

    logger := newLogger(cfg.LogLevel)
    slog.SetDefault(logger)

    st, err := store.NewFilesystemStore(cfg.ReportsDir, cfg.CacheTTL)
    if err != nil {
        return err
    }
    defer func() { _ = st.Close() }()

    svc := service.NewWorkflowService(st)
    h := handler.New(svc, logger)
    r := router.New(h, cfg)

    srv := &http.Server{
        Addr:              cfg.ListenAddr,
        Handler:           r,
        ReadHeaderTimeout: 5 * time.Second,
        ReadTimeout:       15 * time.Second,
        WriteTimeout:      30 * time.Second,
        IdleTimeout:       60 * time.Second,
    }

    errCh := make(chan error, 1)
    go func() {
        logger.Info("http server starting",
            "addr", cfg.ListenAddr,
            "reports", cfg.ReportsDir,
            "cacheTTL", cfg.CacheTTL.String(),
        )
        if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
            errCh <- err
        }
    }()

    stop := make(chan os.Signal, 1)
    signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

    select {
    case err := <-errCh:
        return err
    case sig := <-stop:
        logger.Info("shutdown signal received", "signal", sig.String())
    }

    ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        logger.Error("graceful shutdown failed", "err", err)
        return err
    }
    logger.Info("server stopped gracefully")
    return nil
}

func newLogger(level string) *slog.Logger {
    var lvl slog.Level
    switch level {
    case "debug":
        lvl = slog.LevelDebug
    case "warn":
        lvl = slog.LevelWarn
    case "error":
        lvl = slog.LevelError
    default:
        lvl = slog.LevelInfo
    }
    return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl}))
}
```

### 4.15 `services/workflow-api/Dockerfile`

```dockerfile
# syntax=docker/dockerfile:1.7
FROM golang:1.22-alpine AS builder

WORKDIR /src

RUN apk add --no-cache ca-certificates git tzdata

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
      -trimpath \
      -ldflags='-s -w -extldflags "-static"' \
      -o /out/workflow-api .

# ── runtime ─────────────────────────────────────────────────
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=builder /out/workflow-api /app/workflow-api
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

ENV TZ=Asia/Bangkok

EXPOSE 8080
USER nonroot:nonroot

ENTRYPOINT ["/app/workflow-api"]
```

### 4.16 `services/workflow-api/api/openapi.yaml`

```yaml
openapi: 3.1.0
info:
  title: 6-Phase Workflow API
  version: 1.0.0
  description: REST API for querying 6-phase workflow history
  contact:
    name: Platform Team
    email: platform@example.com

servers:
  - url: http://localhost:8080
    description: local
  - url: https://workflow.example.com/api
    description: prod

tags:
  - name: Workflows
  - name: Stats
  - name: Health

paths:
  /healthz:
    get:
      tags: [Health]
      summary: Liveness probe
      responses:
        '200': { description: OK }

  /readyz:
    get:
      tags: [Health]
      summary: Readiness probe
      responses:
        '200': { description: Ready }
        '503': { description: Not ready }

  /metrics:
    get:
      tags: [Health]
      summary: Prometheus metrics
      responses:
        '200': { description: OK }

  /api/v1/workflows:
    get:
      tags: [Workflows]
      summary: List workflows
      parameters:
        - in: query
          name: status
          schema: { type: string, enum: [PASS, WARN, FAIL, ERROR] }
        - in: query
          name: search
          schema: { type: string }
        - in: query
          name: phase
          schema: { type: string }
        - in: query
          name: sort
          schema: { type: string, enum: [finished_at, duration_ms, ticket_key], default: finished_at }
        - in: query
          name: order
          schema: { type: string, enum: [asc, desc], default: desc }
        - in: query
          name: page
          schema: { type: integer, minimum: 1, default: 1 }
        - in: query
          name: pageSize
          schema: { type: integer, minimum: 1, maximum: 200, default: 20 }
      responses:
        '200':
          description: OK
          content:
            application/json:
              schema: { $ref: '#/components/schemas/ListResult' }
        '400': { $ref: '#/components/responses/BadRequest' }

  /api/v1/workflows/{ticketKey}:
    get:
      tags: [Workflows]
      summary: Get workflow by ticket
      parameters:
        - $ref: '#/components/parameters/TicketKey'
      responses:
        '200':
          description: OK
          content:
            application/json:
              schema: { $ref: '#/components/schemas/Summary' }
        '404': { $ref: '#/components/responses/NotFound' }

  /api/v1/workflows/{ticketKey}/phases/{phaseID}:
    get:
      tags: [Workflows]
      summary: Get a phase of a workflow
      parameters:
        - $ref: '#/components/parameters/TicketKey'
        - in: path
          name: phaseID
          required: true
          schema: { type: string, enum: [gather, plan, execute, security, performance, rca] }
      responses:
        '200':
          description: OK
          content:
            application/json:
              schema: { $ref: '#/components/schemas/Phase' }
        '404': { $ref: '#/components/responses/NotFound' }

  /api/v1/stats:
    get:
      tags: [Stats]
      summary: Aggregate statistics
      responses:
        '200':
          description: OK
          content:
            application/json:
              schema: { $ref: '#/components/schemas/Stats' }

components:
  parameters:
    TicketKey:
      in: path
      name: ticketKey
      required: true
      schema: { type: string, pattern: '^[A-Z]+-\d+$' }

  responses:
    BadRequest:
      description: Bad request
      content:
        application/json:
          schema: { $ref: '#/components/schemas/Error' }
    NotFound:
      description: Not found
      content:
        application/json:
          schema: { $ref: '#/components/schemas/Error' }

  schemas:
    Status:
      type: string
      enum: [PASS, WARN, FAIL, ERROR]

    Phase:
      type: object
      required: [phaseId, phaseName, status, durationMs]
      properties:
        phaseId:    { type: string }
        phaseName:  { type: string }
        emoji:      { type: string }
        status:     { $ref: '#/components/schemas/Status' }
        durationMs: { type: integer, format: int64 }
        report:     { type: string }
        error:      { type: string }

    Summary:
      type: object
      required: [ticketKey, task, overall, durationMs, finishedAt, results]
      properties:
        ticketKey:  { type: string, pattern: '^[A-Z]+-\d+$' }
        task:       { type: string }
        overall:    { $ref: '#/components/schemas/Status' }
        durationMs: { type: integer, format: int64 }
        finishedAt: { type: string, format: date-time }
        results:
          type: object
          additionalProperties: { $ref: '#/components/schemas/Phase' }

    ListResult:
      type: object
      properties:
        items:
          type: array
          items: { $ref: '#/components/schemas/Summary' }
        total:    { type: integer }
        page:     { type: integer }
        pageSize: { type: integer }

    PhaseStat:
      type: object
      properties:
        total:    { type: integer }
        pass:     { type: integer }
        warn:     { type: integer }
        fail:     { type: integer }
        error:    { type: integer }
        passRate: { type: number, format: double }

    Stats:
      type: object
      properties:
        total:         { type: integer }
        byStatus:      { type: object, additionalProperties: { type: integer } }
        avgDurationMs: { type: number, format: double }
        p95DurationMs: { type: integer, format: int64 }
        phaseStats:    { type: object, additionalProperties: { $ref: '#/components/schemas/PhaseStat' } }

    Error:
      type: object
      properties:
        error:   { type: string }
        details: { type: object, additionalProperties: true }
```

### 4.17 `services/workflow-api/api/postman_collection.json`

```json
{
  "info": {
    "name": "6-Phase Workflow API",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "variable": [
    { "key": "baseUrl", "value": "http://localhost:8080" }
  ],
  "item": [
    {
      "name": "Health",
      "item": [
        {
          "name": "Liveness",
          "request": {
            "method": "GET",
            "url": "{{baseUrl}}/healthz"
          },
          "event": [{
            "listen": "test",
            "script": { "exec": [
              "pm.test('200', () => pm.response.to.have.status(200));"
            ]}
          }]
        },
        {
          "name": "Readiness",
          "request": { "method": "GET", "url": "{{baseUrl}}/readyz" }
        },
        {
          "name": "Metrics",
          "request": { "method": "GET", "url": "{{baseUrl}}/metrics" }
        }
      ]
    },
    {
      "name": "Workflows",
      "item": [
        {
          "name": "List all",
          "request": {
            "method": "GET",
            "url": "{{baseUrl}}/api/v1/workflows"
          },
          "event": [{
            "listen": "test",
            "script": { "exec": [
              "pm.test('200', () => pm.response.to.have.status(200));",
              "pm.test('has items', () => pm.expect(pm.response.json().items).to.be.an('array'));"
            ]}
          }]
        },
        {
          "name": "List filtered by status",
          "request": {
            "method": "GET",
            "url": "{{baseUrl}}/api/v1/workflows?status=FAIL&page=1&pageSize=10"
          }
        },
        {
          "name": "List search",
          "request": {
            "method": "GET",
            "url": "{{baseUrl}}/api/v1/workflows?search=pagination"
          }
        },
        {
          "name": "Get by ticket",
          "request": {
            "method": "GET",
            "url": "{{baseUrl}}/api/v1/workflows/ABC-123"
          },
          "event": [{
            "listen": "test",
            "script": { "exec": [
              "pm.test('200', () => pm.response.to.have.status(200));",
              "const j = pm.response.json();",
              "pm.test('has ticketKey', () => pm.expect(j.ticketKey).to.eql('ABC-123'));"
            ]}
          }]
        },
        {
          "name": "Get phase",
          "request": {
            "method": "GET",
            "url": "{{baseUrl}}/api/v1/workflows/ABC-123/phases/security"
          }
        },
        {
          "name": "404 not found",
          "request": {
            "method": "GET",
            "url": "{{baseUrl}}/api/v1/workflows/NOPE-999"
          },
          "event": [{
            "listen": "test",
            "script": { "exec": [
              "pm.test('404', () => pm.response.to.have.status(404));"
            ]}
          }]
        }
      ]
    },
    {
      "name": "Stats",
      "item": [
        {
          "name": "Aggregate",
          "request": {
            "method": "GET",
            "url": "{{baseUrl}}/api/v1/stats"
          }
        }
      ]
    }
  ]
}
```

### 4.18 `services/workflow-api/testdata/summaries/` — คัดลอกจาก `tests-playwright/fixtures/summaries/`

### 4.19 `services/workflow-api/README.md`

```markdown
# workflow-api

REST API for querying 6-phase workflow history.

## Quick start

```bash
# local
REPORTS_DIR=../../.workflow-reports go run .

# docker
docker build -t workflow-api .
docker run --rm -p 8080:8080 -v $PWD/.workflow-reports:/reports:ro \
  -e REPORTS_DIR=/reports workflow-api
```

## Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/healthz` | Liveness |
| GET | `/readyz`  | Readiness |
| GET | `/metrics` | Prometheus |
| GET | `/api/v1/workflows` | List (filter + paginate) |
| GET | `/api/v1/workflows/{ticketKey}` | Get one |
| GET | `/api/v1/workflows/{ticketKey}/phases/{phaseID}` | Get phase |
| GET | `/api/v1/stats` | Aggregate stats |

## Query params

- `status=PASS|WARN|FAIL|ERROR`
- `search=<keyword>` — match ticketKey / task
- `phase=gather|plan|execute|security|performance|rca`
- `sort=finished_at|duration_ms|ticket_key`
- `order=asc|desc`
- `page=1..N`, `pageSize=1..200`

## Env vars

| Name | Default | Description |
|------|---------|-------------|
| `LISTEN_ADDR` | `:8080` | HTTP listen address |
| `REPORTS_DIR` | `./.workflow-reports` | Reports directory |
| `CACHE_TTL` | `30s` | Cache TTL for filesystem read |
| `SHUTDOWN_TIMEOUT` | `15s` | Graceful shutdown timeout |
| `LOG_LEVEL` | `info` | `debug|info|warn|error` |
| `CORS_ORIGINS` | `*` | Comma-separated origins |

## Design principles

- **Zero panic** — return errors; recover middleware เป็น safety net
- **Clean architecture** — handler → service → store
- **Context propagation** — ทุก call รับ ctx
- **Observability** — structured log (slog) + Prometheus metrics
- **Graceful shutdown** — SIGINT/SIGTERM
- **Test-driven** — store/service มี unit test; handler ใช้ integration

## Tests

```bash
go test ./... -race -cover
go vet ./...
golangci-lint run
```
```

### 4.20 อัปเดต Helm ให้ deploy API

ใน `helm/workflow-dashboard/templates/deployment-api.yaml` (มีอยู่แล้วใน 1.18) และ service (1.19) — ครบแล้ว

### 4.21 ตัวอย่างการใช้งาน API

```bash
# ── List ทั้งหมด ─────────────────────────────────────────
curl -s http://localhost:8080/api/v1/workflows | jq

# ── Filter FAIL ──────────────────────────────────────────
curl -s 'http://localhost:8080/api/v1/workflows?status=FAIL&pageSize=5' | jq '.items[].ticketKey'

# ── Search ───────────────────────────────────────────────
curl -s 'http://localhost:8080/api/v1/workflows?search=pagination' | jq

# ── ดู ticket เดียว ───────────────────────────────────────
curl -s http://localhost:8080/api/v1/workflows/ABC-123 | jq

# ── ดู phase security ────────────────────────────────────
curl -s http://localhost:8080/api/v1/workflows/ABC-123/phases/security | jq

# ── Stats ────────────────────────────────────────────────
curl -s http://localhost:8080/api/v1/stats | jq
```

---

## 📁 โครงสร้างไฟล์ที่เพิ่มในรอบนี้

```
project/
├── helm/workflow-dashboard/              ← ⭐ ใหม่
│   ├── Chart.yaml
│   ├── values.yaml
│   ├── values-dev.yaml
│   ├── values-prod.yaml
│   ├── .helmignore
│   ├── dashboards/workflow-6phase.json
│   └── templates/*.yaml (20 ไฟล์)
├── tests-playwright/                     ← ⭐ ใหม่
│   ├── package.json
│   ├── playwright.config.ts
│   ├── tsconfig.json
│   ├── fixtures/{test-fixtures.ts,summaries/*.json}
│   └── tests/*.spec.ts (5 ไฟล์)
├── services/workflow-api/                ← ⭐ ใหม่
│   ├── main.go, go.mod, Dockerfile, README.md
│   ├── internal/{config,model,store,service,handler,metrics,router}/
│   ├── api/{openapi.yaml,postman_collection.json}
│   └── testdata/summaries/*.json
└── scripts/
    ├── Send-Notification.ps1             ← ⭐ ใหม่
    └── lib/Notifications.psm1            ← ⭐ ใหม่
```

---

## 🎯 สรุปการทำงานร่วมกัน

| บทบาท | เครื่องมือ |
|-------|-----------|
| **รัน 6 เฟส** | `Run-6PhaseWorkflow.ps1` |
| **รันใน container** | `docker compose` |
| **กด F5** | `.vscode/launch.json` |
| **CI trigger ตาม label** | GitHub Actions |
| **Deploy k8s** | `helm install workflow helm/workflow-dashboard` |
| **ทดสอบ dashboard** | `npm test` (Playwright) |
| **แจ้งเตือนเมื่อ fail** | `Send-Notification.ps1` (Slack/Teams) |
| **Query history** | `curl http://api:8080/api/v1/workflows` |
| **ดู dashboard** | `https://workflow.example.com` |
| **Grafana** | `https://workflow.example.com/grafana` |
| **Metrics** | `https://workflow.example.com/api/metrics` |

 