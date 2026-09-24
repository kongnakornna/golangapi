# ส่วนเพิ่มเติม - Infrastructure as Code + Operations

## 18. Terraform Modules (AWS)

### 18.1 Root Module Structure

```
deploy/terraform/
├── modules/
│   ├── pdpa-network/           # VPC, subnets, security groups
│   ├── pdpa-rds/               # PostgreSQL (RDS Multi-AZ)
│   ├── pdpa-elasticache/       # Redis (cluster mode)
│   ├── pdpa-msk/               # Kafka (MSK)
│   ├── pdpa-elasticsearch/     # OpenSearch
│   ├── pdpa-secrets/           # AWS Secrets Manager
│   └── pdpa-iam/               # IRSA roles
├── environments/
│   ├── dev/
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   ├── terraform.tfvars
│   │   └── backend.tf
│   ├── staging/
│   └── production/
└── README.md
```

### 18.2 RDS Module (PostgreSQL)

**`deploy/terraform/modules/pdpa-rds/main.tf`**

```hcl
terraform {
  required_version = ">= 1.6.0"
  required_providers {
    aws = { source = "hashicorp/aws", version = "~> 5.60" }
    random = { source = "hashicorp/random", version = "~> 3.6" }
  }
}

# ============================================================
# DB Subnet Group
# ============================================================
resource "aws_db_subnet_group" "this" {
  name       = "${var.name_prefix}-pdpa-db-subnet"
  subnet_ids = var.private_subnet_ids

  tags = merge(var.tags, {
    Name = "${var.name_prefix}-pdpa-db-subnet"
  })
}

# ============================================================
# Security Group — allow access from PDPA app SG only
# ============================================================
resource "aws_security_group" "rds" {
  name        = "${var.name_prefix}-pdpa-rds-sg"
  description = "PDPA RDS access"
  vpc_id      = var.vpc_id

  ingress {
    description     = "Postgres from app SG"
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = [var.app_security_group_id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = var.tags
}

# ============================================================
# KMS key for encryption at rest
# ============================================================
resource "aws_kms_key" "rds" {
  description             = "PDPA RDS encryption key"
  deletion_window_in_days = 30
  enable_key_rotation     = true
  tags                    = var.tags
}

resource "aws_kms_alias" "rds" {
  name          = "alias/${var.name_prefix}-pdpa-rds"
  target_key_id = aws_kms_key.rds.key_id
}

# ============================================================
# Random master password (rotated via Secrets Manager)
# ============================================================
resource "random_password" "master" {
  length           = 32
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}

resource "aws_secretsmanager_secret" "db" {
  name                    = "${var.name_prefix}/pdpa/db"
  description             = "PDPA PostgreSQL credentials"
  recovery_window_in_days = 7
  kms_key_id              = aws_kms_key.rds.arn
  tags                    = var.tags
}

resource "aws_secretsmanager_secret_version" "db" {
  secret_id = aws_secretsmanager_secret.db.id
  secret_string = jsonencode({
    username = var.master_username
    password = random_password.master.result
    host     = aws_db_instance.this.address
    port     = aws_db_instance.this.port
    dbname   = var.database_name
    dsn      = "host=${aws_db_instance.this.address} port=${aws_db_instance.this.port} user=${var.master_username} password=${random_password.master.result} dbname=${var.database_name} sslmode=require"
  })
}

# ============================================================
# Parameter Group — PDPA-tuned
# ============================================================
resource "aws_db_parameter_group" "this" {
  name   = "${var.name_prefix}-pdpa-pg16"
  family = "postgres16"

  parameter {
    name  = "log_statement"
    value = "ddl"
  }
  parameter {
    name  = "log_min_duration_statement"
    value = "1000" # log slow queries (>1s)
  }
  parameter {
    name         = "shared_preload_libraries"
    value        = "pg_stat_statements,pgcrypto,pg_trgm"
    apply_method = "pending-reboot"
  }
  parameter {
    name  = "rds.force_ssl"
    value = "1"
  }
  parameter {
    name  = "password_encryption"
    value = "scram-sha-256"
  }

  tags = var.tags
}

# ============================================================
# RDS Instance (Multi-AZ, encrypted, automated backups)
# ============================================================
resource "aws_db_instance" "this" {
  identifier     = "${var.name_prefix}-pdpa"
  engine         = "postgres"
  engine_version = "16.3"
  instance_class = var.instance_class

  allocated_storage     = var.allocated_storage
  max_allocated_storage = var.max_allocated_storage
  storage_type          = "gp3"
  storage_encrypted     = true
  kms_key_id            = aws_kms_key.rds.arn

  db_name  = var.database_name
  username = var.master_username
  password = random_password.master.result
  port     = 5432

  multi_az               = var.multi_az
  db_subnet_group_name   = aws_db_subnet_group.this.name
  vpc_security_group_ids = [aws_security_group.rds.id]
  parameter_group_name   = aws_db_parameter_group.this.name

  backup_retention_period   = var.backup_retention_days
  backup_window             = "03:00-04:00"    # UTC
  maintenance_window        = "sun:04:00-sun:05:00"
  copy_tags_to_snapshot     = true
  deletion_protection       = var.deletion_protection
  skip_final_snapshot       = var.skip_final_snapshot
  final_snapshot_identifier = var.skip_final_snapshot ? null : "${var.name_prefix}-pdpa-final-${formatdate("YYYY-MM-DD-hhmm", timestamp())}"

  performance_insights_enabled          = true
  performance_insights_retention_period = var.performance_insights_days
  monitoring_interval                   = 60
  monitoring_role_arn                   = aws_iam_role.rds_monitoring.arn

  enabled_cloudwatch_logs_exports = ["postgresql", "upgrade"]

  auto_minor_version_upgrade  = true
  apply_immediately           = false
  allow_major_version_upgrade = false

  tags = merge(var.tags, {
    Name = "${var.name_prefix}-pdpa"
    DataClassification = "PII"
    Compliance         = "PDPA"
  })

  lifecycle {
    ignore_changes = [password, final_snapshot_identifier]
  }
}

# ============================================================
# Enhanced monitoring role
# ============================================================
resource "aws_iam_role" "rds_monitoring" {
  name = "${var.name_prefix}-pdpa-rds-monitoring"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect    = "Allow"
      Principal = { Service = "monitoring.rds.amazonaws.com" }
      Action    = "sts:AssumeRole"
    }]
  })
}

resource "aws_iam_role_policy_attachment" "rds_monitoring" {
  role       = aws_iam_role.rds_monitoring.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonRDSEnhancedMonitoringRole"
}
```

**`deploy/terraform/modules/pdpa-rds/variables.tf`**

```hcl
variable "name_prefix"         { type = string }
variable "vpc_id"              { type = string }
variable "private_subnet_ids"  { type = list(string) }
variable "app_security_group_id" { type = string }

variable "instance_class" {
  type    = string
  default = "db.r6g.large"
}
variable "allocated_storage" {
  type    = number
  default = 100
}
variable "max_allocated_storage" {
  type    = number
  default = 1000
}
variable "multi_az" {
  type    = bool
  default = true
}
variable "backup_retention_days" {
  type    = number
  default = 30
}
variable "performance_insights_days" {
  type    = number
  default = 7
}
variable "database_name" {
  type    = string
  default = "icmongolang"
}
variable "master_username" {
  type    = string
  default = "pdpa_admin"
}
variable "deletion_protection" {
  type    = bool
  default = true
}
variable "skip_final_snapshot" {
  type    = bool
  default = false
}
variable "tags" {
  type    = map(string)
  default = {}
}
```

**`deploy/terraform/modules/pdpa-rds/outputs.tf`**

```hcl
output "endpoint" {
  value       = aws_db_instance.this.endpoint
  description = "RDS endpoint (host:port)"
}

output "address" {
  value = aws_db_instance.this.address
}

output "port" {
  value = aws_db_instance.this.port
}

output "database_name" {
  value = aws_db_instance.this.db_name
}

output "secret_arn" {
  value       = aws_secretsmanager_secret.db.arn
  description = "ARN of the Secrets Manager secret"
}

output "kms_key_arn" {
  value = aws_kms_key.rds.arn
}

output "security_group_id" {
  value = aws_security_group.rds.id
}
```

### 18.3 ElastiCache Module (Redis)

**`deploy/terraform/modules/pdpa-elasticache/main.tf`**

```hcl
terraform {
  required_providers {
    aws = { source = "hashicorp/aws", version = "~> 5.60" }
  }
}

resource "aws_elasticache_subnet_group" "this" {
  name       = "${var.name_prefix}-pdpa-redis-subnet"
  subnet_ids = var.private_subnet_ids
  tags       = var.tags
}

resource "aws_security_group" "redis" {
  name        = "${var.name_prefix}-pdpa-redis-sg"
  description = "PDPA Redis access"
  vpc_id      = var.vpc_id

  ingress {
    description     = "Redis from app SG"
    from_port       = 6379
    to_port         = 6379
    protocol        = "tcp"
    security_groups = [var.app_security_group_id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = var.tags
}

# Auth token stored in Secrets Manager
resource "random_password" "auth" {
  length  = 64
  special = false
}

resource "aws_secretsmanager_secret" "redis" {
  name                    = "${var.name_prefix}/pdpa/redis"
  description             = "PDPA Redis AUTH token"
  recovery_window_in_days = 7
  tags                    = var.tags
}

resource "aws_secretsmanager_secret_version" "redis" {
  secret_id = aws_secretsmanager_secret.redis.id
  secret_string = jsonencode({
    auth_token = random_password.auth.result
    host       = aws_elasticache_replication_group.this.primary_endpoint_address
    port       = 6379
    addr       = "${aws_elasticache_replication_group.this.primary_endpoint_address}:6379"
  })
}

# KMS for at-rest encryption
resource "aws_kms_key" "redis" {
  description         = "PDPA ElastiCache encryption key"
  enable_key_rotation = true
  tags                = var.tags
}

# Subnet group for replication group
resource "aws_elasticache_replication_group" "this" {
  replication_group_id = "${var.name_prefix}-pdpa-redis"
  description          = "PDPA Redis cluster"

  engine               = "redis"
  engine_version       = "7.1"
  node_type            = var.node_type
  port                 = 6379
  parameter_group_name = "default.redis7.cluster.on"

  # Cluster mode with 2 shards × 2 replicas (HA)
  num_node_groups         = var.num_shards
  replicas_per_node_group = var.replicas_per_shard

  automatic_failover_enabled = true
  multi_az_enabled           = true

  subnet_group_name  = aws_elasticache_subnet_group.this.name
  security_group_ids = [aws_security_group.redis.id]

  at_rest_encryption_enabled = true
  transit_encryption_enabled = true
  kms_key_id                 = aws_kms_key.redis.arn
  auth_token                 = random_password.auth.result

  snapshot_retention_limit = var.snapshot_retention_days
  snapshot_window          = "03:00-05:00"
  maintenance_window       = "sun:05:00-sun:06:00"

  auto_minor_version_upgrade = true
  apply_immediately          = false

  log_delivery_configuration {
    destination      = aws_cloudwatch_log_group.redis.name
    destination_type = "cloudwatch-logs"
    log_format       = "json"
    log_type         = "slow-log"
  }

  tags = merge(var.tags, {
    Name = "${var.name_prefix}-pdpa-redis"
    DataClassification = "PII"
    Compliance         = "PDPA"
  })

  lifecycle {
    ignore_changes = [auth_token]
  }
}

resource "aws_cloudwatch_log_group" "redis" {
  name              = "/aws/elasticache/${var.name_prefix}-pdpa-redis"
  retention_in_days = 30
  tags              = var.tags
}
```

### 18.4 MSK Module (Kafka)

**`deploy/terraform/modules/pdpa-msk/main.tf`**

```hcl
terraform {
  required_providers {
    aws = { source = "hashicorp/aws", version = "~> 5.60" }
  }
}

resource "aws_security_group" "msk" {
  name        = "${var.name_prefix}-pdpa-msk-sg"
  description = "PDPA MSK access"
  vpc_id      = var.vpc_id

  ingress {
    description     = "Kafka from app SG"
    from_port       = 9096
    to_port         = 9096
    protocol        = "tcp"
    security_groups = [var.app_security_group_id]
  }
  ingress {
    description     = "Kafka TLS from app SG"
    from_port       = 9094
    to_port         = 9094
    protocol        = "tcp"
    security_groups = [var.app_security_group_id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = var.tags
}

resource "aws_kms_key" "msk" {
  description         = "PDPA MSK encryption key"
  enable_key_rotation = true
  tags                = var.tags
}

resource "aws_msk_cluster" "this" {
  cluster_name           = "${var.name_prefix}-pdpa"
  kafka_version          = "3.6.0"
  number_of_broker_nodes = var.broker_count

  broker_node_group_info {
    instance_type   = var.instance_type
    client_subnets  = var.private_subnet_ids
    security_groups = [aws_security_group.msk.id]

    storage_info {
      ebs_storage_info {
        volume_size = var.ebs_volume_size
      }
    }
  }

  encryption_info {
    encryption_at_rest_kms_key_arn = aws_kms_key.msk.arn
    encryption_in_transit {
      client_broker = "TLS"
      in_cluster    = true
    }
  }

  client_authentication {
    tls {
      certificate_authority_arns = [aws_acmpca_certificate_authority.this.arn]
    }
    sasl {
      scram = true
    }
  }

  configuration_info {
    arn      = aws_msk_configuration.this.arn
    revision = aws_msk_configuration.this.latest_revision
  }

  logging_info {
    broker_logs {
      cloudwatch_logs {
        enabled   = true
        log_group = aws_cloudwatch_log_group.msk.name
      }
      s3 {
        enabled = true
        bucket  = aws_s3_bucket.logs.id
        prefix  = "msk-logs/"
      }
    }
  }

  open_monitoring {
    prometheus {
      jmx_exporter  { enabled_in_broker = true }
      node_exporter { enabled_in_broker = true }
    }
  }

  tags = merge(var.tags, {
    Name = "${var.name_prefix}-pdpa-msk"
    DataClassification = "PII"
    Compliance         = "PDPA"
  })
}

# Custom MSK config tuned for PDPA workloads
resource "aws_msk_configuration" "this" {
  name           = "${var.name_prefix}-pdpa-msk-config"
  kafka_versions = ["3.6.0"]

  server_properties = <<PROPERTIES
auto.create.topics.enable=false
default.replication.factor=3
min.insync.replicas=2
num.partitions=6
log.retention.hours=168
log.retention.bytes=107374182400
log.segment.bytes=1073741824
compression.type=producer
unclean.leader.election.enable=false
transaction.state.log.replication.factor=3
transaction.state.log.min.isr=2
PROPERTIES
}

# ACM Private CA for TLS auth
resource "aws_acmpca_certificate_authority" "this" {
  certificate_authority_configuration {
    key_algorithm     = "RSA_4096"
    signing_algorithm = "SHA512WITHRSA"
    subject {
      common_name = "${var.name_prefix}-pdpa-msk-ca"
    }
  }
  permanent_deletion_time_in_days = 30
  type                            = "SUBORDINATE"
  tags                            = var.tags
}

resource "aws_cloudwatch_log_group" "msk" {
  name              = "/aws/msk/${var.name_prefix}-pdpa"
  retention_in_days = 30
  tags              = var.tags
}

resource "aws_s3_bucket" "logs" {
  bucket        = "${var.name_prefix}-pdpa-msk-logs-${random_id.bucket.hex}"
  force_destroy = false
  tags          = var.tags
}

resource "random_id" "bucket" { byte_length = 4 }
```

### 18.5 Secrets Module

**`deploy/terraform/modules/pdpa-secrets/main.tf`**

```hcl
# Application-level secrets (ไม่รวม DSN ที่มาจาก RDS module)
resource "random_password" "anon_salt" {
  length  = 48
  special = false
}

resource "aws_secretsmanager_secret" "app" {
  name                    = "${var.name_prefix}/pdpa/app"
  description             = "PDPA application secrets"
  recovery_window_in_days = 7
  kms_key_id              = var.kms_key_arn
  tags                    = var.tags
}

resource "aws_secretsmanager_secret_version" "app" {
  secret_id = aws_secretsmanager_secret.app.id
  secret_string = jsonencode({
    anonSalt       = random_password.anon_salt.result
    llmApiKey      = var.llm_api_key
    blockchainKey  = var.blockchain_api_key
    smtpUser       = var.smtp_user
    smtpPassword   = var.smtp_password
  })
}

output "app_secret_arn" {
  value = aws_secretsmanager_secret.app.arn
}
```

### 18.6 IRSA Module (IAM Roles for Service Accounts)

**`deploy/terraform/modules/pdpa-iam/main.tf`**

```hcl
# IRSA role allowing PDPA pods to read Secrets Manager
data "aws_iam_policy_document" "assume" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRoleWithWebIdentity"]
    principals {
      type        = "Federated"
      identifiers = [var.oidc_provider_arn]
    }
    condition {
      test     = "StringEquals"
      variable = "${var.oidc_provider_url}:sub"
      values   = ["system:serviceaccount:${var.namespace}:${var.service_account_name}"]
    }
  }
}

resource "aws_iam_role" "pdpa" {
  name               = "${var.name_prefix}-pdpa-irsa"
  assume_role_policy = data.aws_iam_policy_document.assume.json
  tags               = var.tags
}

data "aws_iam_policy_document" "pdpa" {
  statement {
    sid    = "ReadSecrets"
    effect = "Allow"
    actions = [
      "secretsmanager:GetSecretValue",
      "secretsmanager:DescribeSecret",
    ]
    resources = var.secret_arns
  }

  statement {
    sid    = "KMSSDecrypt"
    effect = "Allow"
    actions = [
      "kms:Decrypt",
      "kms:DescribeKey",
    ]
    resources = var.kms_key_arns
  }

  statement {
    sid    = "MSKIAMAuth"
    effect = "Allow"
    actions = ["kafka-cluster:*"]
    resources = ["*"]
  }
}

resource "aws_iam_role_policy" "pdpa" {
  name   = "${var.name_prefix}-pdpa-policy"
  role   = aws_iam_role.pdpa.id
  policy = data.aws_iam_policy_document.pdpa.json
}

output "role_arn" { value = aws_iam_role.pdpa.arn }
```

### 18.7 Production Environment

**`deploy/terraform/environments/production/main.tf`**

```hcl
terraform {
  required_version = ">= 1.6.0"
  required_providers {
    aws = { source = "hashicorp/aws", version = "~> 5.60" }
  }
}

provider "aws" {
  region = var.region

  default_tags {
    tags = {
      Project    = "icmongolang"
      Module     = "pdpa"
      Environment = var.environment
      ManagedBy  = "terraform"
    }
  }
}

# ============================================================
# VPC (ใช้ VPC ที่มีอยู่ หรือสร้างใหม่)
# ============================================================
data "aws_vpc" "this" {
  id = var.vpc_id
}

data "aws_subnets" "private" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.this.id]
  }
  filter {
    name   = "tag:Tier"
    values = ["private"]
  }
}

# ============================================================
# RDS
# ============================================================
module "rds" {
  source = "../../modules/pdpa-rds"

  name_prefix           = "${var.name_prefix}-prod"
  vpc_id                = data.aws_vpc.this.id
  private_subnet_ids    = data.aws_subnets.private.ids
  app_security_group_id = var.app_security_group_id

  instance_class           = "db.r6g.2xlarge"
  allocated_storage        = 500
  max_allocated_storage    = 5000
  multi_az                 = true
  backup_retention_days    = 90
  performance_insights_days = 731
  deletion_protection      = true
  skip_final_snapshot      = false
}

# ============================================================
# ElastiCache
# ============================================================
module "redis" {
  source = "../../modules/pdpa-elasticache"

  name_prefix           = "${var.name_prefix}-prod"
  vpc_id                = data.aws_vpc.this.id
  private_subnet_ids    = data.aws_subnets.private.ids
  app_security_group_id = var.app_security_group_id

  node_type              = "cache.r7g.xlarge"
  num_shards             = 3
  replicas_per_shard     = 2
  snapshot_retention_days = 30
}

# ============================================================
# MSK
# ============================================================
module "msk" {
  source = "../../modules/pdpa-msk"

  name_prefix           = "${var.name_prefix}-prod"
  vpc_id                = data.aws_vpc.this.id
  private_subnet_ids    = data.aws_subnets.private.ids
  app_security_group_id = var.app_security_group_id

  broker_count    = 6
  instance_type   = "kafka.m7g.xlarge"
  ebs_volume_size = 500
}

# ============================================================
# Secrets
# ============================================================
module "secrets" {
  source = "../../modules/pdpa-secrets"

  name_prefix = "${var.name_prefix}-prod"
  kms_key_arn = module.rds.kms_key_arn
  llm_api_key = var.llm_api_key
  smtp_user   = var.smtp_user
  smtp_password = var.smtp_password
}

# ============================================================
# IAM / IRSA
# ============================================================
module "iam" {
  source = "../../modules/pdpa-iam"

  name_prefix          = "${var.name_prefix}-prod"
  namespace            = "pdpa"
  service_account_name = "pdpa"
  oidc_provider_arn    = var.eks_oidc_provider_arn
  oidc_provider_url    = var.eks_oidc_provider_url

  secret_arns  = [module.rds.secret_arn, module.redis.secret_arn, module.secrets.app_secret_arn]
  kms_key_arns = [module.rds.kms_key_arn]
}

# ============================================================
# Outputs (ใช้สำหรับ Helm values)
# ============================================================
output "rds_endpoint"    { value = module.rds.endpoint    }
output "redis_endpoint"  { value = module.redis_endpoint  }
output "msk_brokers"     { value = module.msk.bootstrap_brokers_tls }
output "irsa_role_arn"   { value = module.iam.role_arn }
```

**`deploy/terraform/environments/production/backend.tf`**

```hcl
terraform {
  backend "s3" {
    bucket         = "icmongolang-terraform-state"
    key            = "pdpa/production/terraform.tfstate"
    region         = "ap-southeast-1"
    encrypt        = true
    kms_key_id     = "arn:aws:kms:ap-southeast-1:123456789012:key/xxxxxxxx"
    dynamodb_table = "terraform-locks"
  }
}
```

**`deploy/terraform/environments/production/terraform.tfvars`**

```hcl
region                = "ap-southeast-1"
environment           = "production"
name_prefix           = "icmongolang"
vpc_id                = "vpc-xxxxxxxx"
app_security_group_id = "sg-xxxxxxxx"
eks_oidc_provider_arn = "arn:aws:iam::123456789012:oidc-provider/oidc.eks.ap-southeast-1.amazonaws.com/id/xxxxxxxx"
eks_oidc_provider_url = "oidc.eks.ap-southeast-1.amazonaws.com/id/xxxxxxxx"
```

---

## 19. ArgoCD Application Manifests (GitOps)

### 19.1 Repository Structure

```
deploy/gitops/
├── argocd/
│   ├── projects/
│   │   └── pdpa-project.yaml
│   ├── applications/
│   │   ├── pdpa-staging.yaml
│   │   ├── pdpa-production.yaml
│   │   └── pdpa-monitoring.yaml
│   └── app-of-apps.yaml
├── overlays/
│   ├── staging/
│   │   └── values-patch.yaml
│   └── production/
│       └── values-patch.yaml
└── README.md
```

### 19.2 AppProject

**`deploy/gitops/argocd/projects/pdpa-project.yaml`**

```yaml
apiVersion: argoproj.io/v1alpha1
kind: AppProject
metadata:
  name: pdpa
  namespace: argocd
  finalizers:
    - resources-finalizer.argocd.argoproj.io
spec:
  description: PDPA Compliance Module

  sourceRepos:
    - 'https://github.com/icmongolang/icmongolang.git'
    - 'https://github.com/icmongolang/icmongolang-helm.git'

  destinations:
    - namespace: 'pdpa-*'
      server: https://kubernetes.default.svc
    - namespace: 'pdpa'
      server: https://kubernetes.default.svc
    - namespace: 'observability'
      server: https://kubernetes.default.svc

  # อนุญาตให้ ArgoCD สร้าง resources เหล่านี้ได้
  clusterResourceWhitelist:
    - group: ''
      kind: Namespace
    - group: 'rbac.authorization.k8s.io'
      kind: ClusterRole
    - group: 'rbac.authorization.k8s.io'
      kind: ClusterRoleBinding
    - group: 'networking.k8s.io'
      kind: NetworkPolicy
    - group: 'policy'
      kind: PodDisruptionBudget

  namespaceResourceWhitelist:
    - group: 'apps'
      kind: Deployment
    - group: 'apps'
      kind: StatefulSet
    - group: 'batch'
      kind: CronJob
    - group: ''
      kind: Service
    - group: ''
      kind: ServiceAccount
    - group: ''
      kind: ConfigMap
    - group: ''
      kind: Secret
    - group: 'networking.k8s.io'
      kind: Ingress
    - group: 'autoscaling'
      kind: HorizontalPodAutoscaler
    - group: 'policy'
      kind: PodDisruptionBudget
    - group: 'monitoring.coreos.com'
      kind: ServiceMonitor
    - group: 'monitoring.coreos.com'
      kind: PrometheusRule
    - group: 'networking.k8s.io'
      kind: NetworkPolicy

  # ป้องกันไม่ให้ ArgoCD ลบ resources ที่อาจยังต้องใช้
  orphanedResources:
    warn: true
```

### 19.3 App-of-Apps (Bootstrap)

**`deploy/gitops/argocd/app-of-apps.yaml`**

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: pdpa-bootstrap
  namespace: argocd
  finalizers:
    - resources-finalizer.argocd.argoproj.io
spec:
  project: default
  source:
    repoURL: https://github.com/icmongolang/icmongolang.git
    targetRevision: main
    path: deploy/gitops/argocd/applications
    directory:
      recurse: true
  destination:
    server: https://kubernetes.default.svc
    namespace: argocd
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
    syncOptions:
      - CreateNamespace=true
      - ServerSideApply=true
```

### 19.4 Staging Application

**`deploy/gitops/argocd/applications/pdpa-staging.yaml`**

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: pdpa-staging
  namespace: argocd
  finalizers:
    - resources-finalizer.argocd.argoproj.io
  annotations:
    notifications.argoproj.io/subscribe.on-sync-succeeded.slack: pdpa-deployments
    notifications.argoproj.io/subscribe.on-sync-failed.slack: pdpa-deployments
    notifications.argoproj.io/subscribe.on-health-degraded.slack: pdpa-alerts
spec:
  project: pdpa

  source:
    repoURL: https://github.com/icmongolang/icmongolang.git
    targetRevision: develop
    path: deploy/helm/pdpa
    helm:
      releaseName: pdpa
      valueFiles:
        - values.yaml
        - ../../gitops/overlays/staging/values-patch.yaml
      parameters:
        - name: global.imageTag
          value: "develop"
        - name: ingress.hosts[0].host
          value: "staging-api.icmongolang.dev"

  destination:
    server: https://kubernetes.default.svc
    namespace: pdpa-staging

  syncPolicy:
    automated:
      prune: true
      selfHeal: true
      allowEmpty: false
    syncOptions:
      - CreateNamespace=true
      - PrunePropagationPolicy=foreground
      - PruneLast=true
      - ServerSideApply=true
      - ApplyOutOfSyncOnly=true
    retry:
      limit: 5
      backoff:
        duration: 10s
        factor: 2
        maxDuration: 5m

  revisionHistoryLimit: 10

  ignoreDifferences:
    - group: apps
      kind: Deployment
      jsonPointers:
        - /spec/replicas   # HPA manages this
```

### 19.5 Production Application

**`deploy/gitops/argocd/applications/pdpa-production.yaml`**

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: pdpa-production
  namespace: argocd
  finalizers:
    - resources-finalizer.argocd.argoproj.io
  annotations:
    notifications.argoproj.io/subscribe.on-sync-succeeded.slack: pdpa-deployments
    notifications.argoproj.io/subscribe.on-sync-failed.slack: pdpa-alerts
    notifications.argoproj.io/subscribe.on-health-degraded.slack: pdpa-alerts
    argocd.argoproj.io/managed-by: argocd
spec:
  project: pdpa

  source:
    repoURL: https://github.com/icmongolang/icmongolang.git
    targetRevision: main
    path: deploy/helm/pdpa
    helm:
      releaseName: pdpa
      valueFiles:
        - values.yaml
        - ../../gitops/overlays/production/values-patch.yaml
      parameters:
        - name: global.imageTag
          value: "pdpa/v1.0.0"
        - name: ingress.hosts[0].host
          value: "api.icmongolang.dev"

  destination:
    server: https://kubernetes.default.svc
    namespace: pdpa

  # ⚠️ Production: manual sync only
  syncPolicy:
    syncOptions:
      - CreateNamespace=false
      - ServerSideApply=true
      - ApplyOutOfSyncOnly=true
      - RespectIgnoreDifferences=true
    retry:
      limit: 3
      backoff:
        duration: 30s
        factor: 2
        maxDuration: 5m

  revisionHistoryLimit: 20

  ignoreDifferences:
    - group: apps
      kind: Deployment
      jsonPointers: [/spec/replicas]
```

### 19.6 Production Overlay

**`deploy/gitops/overlays/production/values-patch.yaml`**

```yaml
api:
  replicaCount: 5
  autoscaling:
    enabled: true
    minReplicas: 5
    maxReplicas: 50
  resources:
    requests: { cpu: 500m, memory: 512Mi }
    limits:   { cpu: 2000m, memory: 2Gi }

worker:
  replicaCount: 4
  autoscaling:
    minReplicas: 4
    maxReplicas: 24
  resources:
    requests: { cpu: 500m, memory: 1Gi }
    limits:   { cpu: 2000m, memory: 2Gi }

externalDependencies:
  postgres:
    host: "icmongolang-prod-pdpa.xxxxx.ap-southeast-1.rds.amazonaws.com"
    database: "icmongolang"
    existingSecret: "pdpa-postgres-prod"
  redis:
    host: "icmongolang-prod-pdpa-redis.xxxxx.clustercfg.apse1.cache.amazonaws.com"
    existingSecret: "pdpa-redis-prod"
  kafka:
    brokers: "b-1.icmongolang-prod.xxxxx.c6.kafka.ap-southeast-1.amazonaws.com:9096,b-2...,b-3..."

ingress:
  hosts:
    - host: api.icmongolang.dev
      paths:
        - path: /api/v1/pdpa
          pathType: Prefix
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
    nginx.ingress.kubernetes.io/limit-rps: "200"
    nginx.ingress.kubernetes.io/limit-burst-multiplier: "3"

config:
  logLevel: warn
```

### 19.7 ExternalSecrets Integration

**`deploy/gitops/overlays/production/externalsecrets.yaml`**

```yaml
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: pdpa-postgres-prod
  namespace: pdpa
spec:
  refreshInterval: 1h
  secretStoreRef:
    name: aws-secrets-manager
    kind: ClusterSecretStore
  target:
    name: pdpa-postgres-prod
    creationPolicy: Owner
  dataFrom:
    - extract:
        key: icmongolang-prod/pdpa/db
---
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: pdpa-redis-prod
  namespace: pdpa
spec:
  refreshInterval: 1h
  secretStoreRef:
    name: aws-secrets-manager
    kind: ClusterSecretStore
  target:
    name: pdpa-redis-prod
  dataFrom:
    - extract:
        key: icmongolang-prod/pdpa/redis
---
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: pdpa-app-prod
  namespace: pdpa
spec:
  refreshInterval: 1h
  secretStoreRef:
    name: aws-secrets-manager
    kind: ClusterSecretStore
  target:
    name: pdpa-app-prod
  dataFrom:
    - extract:
        key: icmongolang-prod/pdpa/app
```

### 19.8 ArgoCD Notifications ConfigMap

**`deploy/gitops/argocd/notifications-cm.yaml`**

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: argocd-notifications-cm
  namespace: argocd
data:
  service.slack: |
    token: $slack-token

  template.app-sync-status: |
    message: |
      *PDPA {{.app.metadata.name}}* — {{.app.status.sync.status}}
      *Health*: {{.app.status.health.status}}
      *Revision*: {{.app.status.sync.revision}}
      *URL*: {{.context.argocdUrl}}/applications/{{.app.metadata.name}}
    slack:
      attachments: |
        [{
          "title": "{{.app.metadata.name}}",
          "title_link": "{{.context.argocdUrl}}/applications/{{.app.metadata.name}}",
          "color": "#18be52",
          "fields": [
            {"title":"Sync","value":"{{.app.status.sync.status}}","short":true},
            {"title":"Health","value":"{{.app.status.health.status}}","short":true}
          ]
        }]

  trigger.on-sync-succeeded: |
    - when: app.status.sync.status == 'Synced' and app.status.health.status == 'Healthy'
      send: [app-sync-status]

  trigger.on-sync-failed: |
    - when: app.status.operationState.phase in ['Error', 'Failed']
      send: [app-sync-status]

  trigger.on-health-degraded: |
    - when: app.status.health.status == 'Degraded'
      send: [app-sync-status]
```

---

## 20. Load Test Scripts (k6 + Locust)

### 20.1 k6 Test Suite

**`test/load/k6/lib/config.js`**

```javascript
export const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
export const TEST_TOKEN = __ENV.TEST_TOKEN || '';
export const USER_IDS = (__ENV.USER_IDS || '').split(',').filter(Boolean);

export const thresholds = {
  // 95% ของ requests ต้องเสร็จภายใน 500ms
  http_req_duration: ['p(95)<500', 'p(99)<1500'],
  // error rate < 1%
  http_req_failed: ['rate<0.01'],
  // ตรวจสอบ checks
  checks: ['rate>0.99'],
};

export function headers() {
  return {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${TEST_TOKEN}`,
    'X-Trace-ID': `k6-${__VU}-${__ITER}-${Date.now()}`,
  };
}

export function randomUser() {
  if (USER_IDS.length === 0) return crypto.randomUUID();
  return USER_IDS[Math.floor(Math.random() * USER_IDS.length)];
}

import crypto from 'k6/crypto';
```

**`test/load/k6/scenarios/consent-write.js`**

```javascript
import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend, Counter } from 'k6/metrics';
import { BASE_URL, headers, thresholds, randomUser } from '../lib/config.js';

const consentSuccess = new Rate('pdpa_consent_success');
const consentDuration = new Trend('pdpa_consent_duration');
const consentErrors = new Counter('pdpa_consent_errors');

export const options = {
  scenarios: {
    // Warm-up + sustained load
    constant_rate: {
      executor: 'constant-arrival-rate',
      rate: 500,                    // 500 req/s
      timeUnit: '1s',
      duration: '5m',
      preAllocatedVUs: 200,
      maxVUs: 1000,
    },
    // Spike test
    spike: {
      executor: 'ramping-arrival-rate',
      startRate: 100,
      timeUnit: '1s',
      preAllocatedVUs: 200,
      maxVUs: 2000,
      stages: [
        { target: 500,  duration: '2m' },
        { target: 500,  duration: '3m' },
        { target: 2000, duration: '30s' }, // spike
        { target: 2000, duration: '1m' },
        { target: 500,  duration: '30s' },
        { target: 0,    duration: '1m' },
      ],
    },
  },
  thresholds,
};

export default function () {
  const userID = randomUser();
  const payload = JSON.stringify({
    purposes: {
      NECESSARY: true,
      ANALYTICS: Math.random() > 0.3,
      MARKETING: Math.random() > 0.7,
    },
    session_id: `sess-${userID}-${Date.now()}`,
  });

  const start = Date.now();
  const res = http.post(`${BASE_URL}/api/v1/pdpa/consent`, payload, { headers: headers() });
  consentDuration.add(Date.now() - start);

  const ok = check(res, {
    'status is 200': (r) => r.status === 200,
    'has purposes': (r) => r.json('purposes') !== undefined,
    'response time < 500ms': (r) => r.timings.duration < 500,
  });

  consentSuccess.add(ok);
  if (!ok) consentErrors.add(1);

  sleep(0.1);
}
```

**`test/load/k6/scenarios/dsar-submit.js`**

```javascript
import http from 'k6/http';
import { check, sleep } from 'k6';
import { BASE_URL, headers, randomUser } from '../lib/config.js';

export const options = {
  scenarios: {
    dsar: {
      executor: 'ramping-vus',
      startVUs: 10,
      stages: [
        { duration: '1m', target: 50 },
        { duration: '5m', target: 100 },
        { duration: '1m', target: 0 },
      ],
    },
  },
  thresholds: {
    http_req_duration: ['p(95)<800'],
    http_req_failed: ['rate<0.02'],
    'checks{scenario:dsar}': ['rate>0.95'],
  },
};

const TYPES = ['ACCESS', 'ERASURE', 'WITHDRAW_CONSENT'];

export default function () {
  const userID = randomUser();
  const type = TYPES[Math.floor(Math.random() * TYPES.length)];

  const res = http.post(
    `${BASE_URL}/api/v1/pdpa/dsar`,
    JSON.stringify({ request_type: type }),
    { headers: headers() },
  );

  check(res, {
    'status is 201': (r) => r.status === 201,
    'has id': (r) => r.json('id') !== undefined,
    'status is PENDING': (r) => r.json('status') === 'PENDING',
    // rate limit (429) ถือว่า acceptable
    'no 5xx': (r) => r.status < 500,
  });

  sleep(Math.random() * 2 + 1);
}
```

**`test/load/k6/scenarios/consent-read.js`**

```javascript
import http from 'k6/http';
import { check, sleep } from 'k6';
import { BASE_URL, headers, randomUser } from '../lib/config.js';

export const options = {
  scenarios: {
    read_heavy: {
      executor: 'constant-vus',
      vus: 200,
      duration: '5m',
    },
  },
  thresholds: {
    http_req_duration: ['p(95)<200', 'p(99)<500'],
    http_req_failed: ['rate<0.005'],
  },
};

export default function () {
  const userID = randomUser();

  // History
  let res = http.get(`${BASE_URL}/api/v1/pdpa/consent/history`, { headers: headers() });
  check(res, { 'history 200': (r) => r.status === 200 });

  sleep(0.05);

  // Active policy (cached)
  res = http.get(`${BASE_URL}/api/v1/pdpa/policy`, { headers: headers() });
  check(res, { 'policy 200 or 404': (r) => r.status === 200 || r.status === 404 });
}
```

**`test/load/k6/scenarios/mixed-soak.js`**

```javascript
import http from 'k6/http';
import { check, sleep } from 'k6';
import { Trend } from 'k6/metrics';
import { BASE_URL, headers, randomUser } from '../lib/config.js';

const writeLatency = new Trend('pdpa_write_latency');

// 4-hour soak test — verify memory leaks / connection leaks
export const options = {
  scenarios: {
    mixed: {
      executor: 'ramping-vus',
      startVUs: 50,
      stages: [
        { duration: '10m', target: 200 },
        { duration: '3h40m', target: 200 },  // sustained
        { duration: '10m', target: 0 },
      ],
    },
  },
  thresholds: {
    http_req_duration: ['p(95)<600', 'p(99)<2000'],
    http_req_failed: ['rate<0.01'],
    pdpa_write_latency: ['p(95)<500'],
  },
};

export default function () {
  const roll = Math.random();

  if (roll < 0.6) {
    // 60% write
    const start = Date.now();
    const res = http.post(
      `${BASE_URL}/api/v1/pdpa/consent`,
      JSON.stringify({ purposes: { NECESSARY: true, ANALYTICS: true } }),
      { headers: headers() },
    );
    writeLatency.add(Date.now() - start);
    check(res, { 'consent 200': (r) => r.status === 200 });
  } else if (roll < 0.9) {
    // 30% read
    const res = http.get(`${BASE_URL}/api/v1/pdpa/consent/history`, { headers: headers() });
    check(res, { 'history 200': (r) => r.status === 200 });
  } else {
    // 10% DSAR submit (rate-limited)
    const res = http.post(
      `${BASE_URL}/api/v1/pdpa/dsar`,
      JSON.stringify({ request_type: 'ACCESS' }),
      { headers: headers() },
    );
    check(res, { 'dsar 200/201/429': (r) => [200, 201, 429].includes(r.status) });
  }

  sleep(Math.random() * 0.5);
}
```

**`test/load/k6/scenarios/websocket.js`**

```javascript
import ws from 'k6/ws';
import { check } from 'k6';

const BASE_WS = __ENV.BASE_WS || 'ws://localhost:8080';
const TEST_TOKEN = __ENV.TEST_TOKEN || '';

export const options = {
  scenarios: {
    ws_clients: {
      executor: 'ramping-vus',
      startVUs: 100,
      stages: [
        { duration: '1m', target: 500 },
        { duration: '5m', target: 1000 },
        { duration: '1m', target: 0 },
      ],
    },
  },
};

export default function () {
  const url = `${BASE_WS}/api/v1/pdpa/ws/dsar-status?token=${TEST_TOKEN}`;
  const res = ws.connect(url, {}, function (socket) {
    socket.on('open', () => {
      // Keep connection alive, receive messages
    });
    socket.on('message', (data) => {
      check(data, { 'received message': (d) => d.length > 0 });
    });
    socket.on('error', (e) => console.error('ws error', e));
    socket.setTimeout(() => socket.close(), 30000);
  });

  check(res, { 'status 101': (r) => r && r.status === 101 });
}
```

### 20.2 Running k6

**`test/load/k6/run.sh`**

```bash
#!/usr/bin/env bash
# ============================================================================
# k6 Load Test Runner
# ============================================================================
set -euo pipefail

SCENARIO="${1:-consent-write}"
ENV="${2:-local}"

case "$ENV" in
  local)      BASE_URL="http://localhost:8080" ;;
  staging)    BASE_URL="https://staging-api.icmongolang.dev" ;;
  production) BASE_URL="https://api.icmongolang.dev" ;;
esac

export BASE_URL
export TEST_TOKEN="${TEST_TOKEN:-$(generate-jwt)}"
export USER_IDS="${USER_IDS:-$(cat test/load/user-ids.txt 2>/dev/null | tr '\n' ',' | sed 's/,$//')}"

echo "▶ Running k6 scenario=$SCENARIO env=$ENV base=$BASE_URL"

case "$SCENARIO" in
  consent-write) k6 run test/load/k6/scenarios/consent-write.js ;;
  consent-read)  k6 run test/load/k6/scenarios/consent-read.js ;;
  dsar)          k6 run test/load/k6/scenarios/dsar-submit.js ;;
  soak)          k6 run test/load/k6/scenarios/mixed-soak.js ;;
  websocket)     k6 run test/load/k6/scenarios/websocket.js ;;
  all)
    for s in consent-read consent-write dsar; do
      k6 run "test/load/k6/scenarios/$s.js"
    done
    ;;
  *) echo "Unknown scenario: $SCENARIO"; exit 1 ;;
esac
```

### 20.3 Locust Alternative (Python)

**`test/load/locust/locustfile.py`**

```python
"""
Locust load test for PDPA module.
Run:
    locust -f test/load/locust/locustfile.py --host=http://localhost:8080
"""
import os
import random
import time

from locust import HttpUser, task, between, events, tag


USER_IDS = [u for u in os.environ.get("USER_IDS", "").split(",") if u]
TEST_TOKEN = os.environ.get("TEST_TOKEN", "")


class PDPAUser(HttpUser):
    wait_time = between(0.5, 2.0)

    def on_start(self):
        self.user_id = random.choice(USER_IDS) if USER_IDS else None
        self.headers = {
            "Content-Type": "application/json",
            "Authorization": f"Bearer {TEST_TOKEN}",
        }

    # ---------------------------------------------------------------
    # Consent — write
    # ---------------------------------------------------------------
    @tag("consent", "write")
    @task(6)
    def record_consent(self):
        payload = {
            "purposes": {
                "NECESSARY": True,
                "ANALYTICS": random.random() > 0.3,
                "MARKETING": random.random() > 0.7,
            },
            "session_id": f"locust-{random.randint(10000, 99999)}",
        }
        with self.client.post(
            "/api/v1/pdpa/consent",
            json=payload,
            headers=self.headers,
            catch_response=True,
            name="POST /consent",
        ) as r:
            if r.status_code != 200:
                r.failure(f"unexpected {r.status_code}")

    # ---------------------------------------------------------------
    # Consent — read
    # ---------------------------------------------------------------
    @tag("consent", "read")
    @task(3)
    def consent_history(self):
        self.client.get(
            "/api/v1/pdpa/consent/history",
            headers=self.headers,
            name="GET /consent/history",
        )

    # ---------------------------------------------------------------
    # DSAR
    # ---------------------------------------------------------------
    @tag("dsar")
    @task(1)
    def submit_dsar(self):
        payload = {"request_type": random.choice(["ACCESS", "ERASURE", "WITHDRAW_CONSENT"])}
        with self.client.post(
            "/api/v1/pdpa/dsar",
            json=payload,
            headers=self.headers,
            catch_response=True,
            name="POST /dsar",
        ) as r:
            # ยอมรับ 429 (rate limited)
            if r.status_code in (200, 201):
                r.success()
            elif r.status_code == 429:
                r.success()
            else:
                r.failure(f"unexpected {r.status_code}")

    # ---------------------------------------------------------------
    # Policy (public)
    # ---------------------------------------------------------------
    @tag("policy", "read")
    @task(2)
    def get_policy(self):
        self.client.get("/api/v1/pdpa/policy", name="GET /policy")


@events.test_start.add_listener
def on_test_start(environment, **kwargs):
    print("📊 PDPA Locust test starting...")
    print(f"   Target: {environment.host}")
    print(f"   Users configured: {len(USER_IDS)}")


@events.test_stop.add_listener
def on_test_stop(environment, **kwargs):
    stats = environment.runner.stats.total
    print("\n📈 Test results:")
    print(f"   Requests: {stats.num_requests}")
    print(f"   Failures: {stats.num_failures}")
    print(f"   RPS: {stats.total_rps:.2f}")
    print(f"   p95: {stats.get_response_time_percentile(0.95)} ms")
    print(f"   p99: {stats.get_response_time_percentile(0.99)} ms")


# unused import guard
_ = time
```

**`test/load/locust/run.sh`**

```bash
#!/usr/bin/env bash
set -euo pipefail

MODE="${1:-web}"       # web | headless
HOST="${2:-http://localhost:8080}"

if [ "$MODE" = "headless" ]; then
  locust -f test/load/locust/locustfile.py \
    --host="$HOST" \
    --headless \
    --users 500 \
    --spawn-rate 50 \
    --run-time 10m \
    --only-summary \
    --html report.html
else
  locust -f test/load/locust/locustfile.py --host="$HOST"
  # เปิด http://localhost:8089
fi
```

### 20.4 Load Test CI Integration

**`.github/workflows/pdpa-loadtest.yml`**

```yaml
name: PDPA Load Test

on:
  workflow_dispatch:
    inputs:
      scenario: { required: true, default: consent-write, type: choice, options: [consent-write, consent-read, dsar, soak, all] }
      environment: { required: true, default: staging, type: choice, options: [staging, production] }
  schedule:
    - cron: '0 4 * * 0'   # weekly Sunday 04:00 UTC (staging only)

jobs:
  load-test:
    name: k6 ${{ github.event.inputs.scenario || 'consent-write' }}
    runs-on: ubuntu-latest
    timeout-minutes: 60
    steps:
      - uses: actions/checkout@v4

      - name: Install k6
        run: |
          sudo gpg -k
          sudo gpg --no-default-keyring --keyring /usr/share/keyrings/k6-archive-keyring.gpg --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
          echo "deb [signed-by=/usr/share/keyrings/k6-archive-keyring.gpg] https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
          sudo apt-get update && sudo apt-get install k6

      - name: Run k6
        env:
          BASE_URL: ${{ github.event.inputs.environment == 'production' && 'https://api.icmongolang.dev' || 'https://staging-api.icmongolang.dev' }}
          TEST_TOKEN: ${{ secrets.LOAD_TEST_TOKEN }}
          USER_IDS: ${{ secrets.LOAD_TEST_USER_IDS }}
        run: |
          bash test/load/k6/run.sh "${{ github.event.inputs.scenario || 'consent-write' }}" "${{ github.event.inputs.environment || 'staging' }}"

      - uses: actions/upload-artifact@v4
        if: always()
        with:
          name: k6-results
          path: |
            summary.json
            *.html
```

---

## 21. Postman Collection

### 21.1 Collection

**`api/postman/PDPA.postman_collection.json`**

```json
{
  "info": {
    "_postman_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "name": "PDPA Module",
    "description": "Personal Data Protection Act (PDPA) — API collection\n\n**Environments:** Local · Staging · Production",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "auth": {
    "type": "bearer",
    "bearer": [{ "key": "token", "value": "{{access_token}}", "type": "string" }]
  },
  "event": [
    {
      "listen": "prerequest",
      "script": {
        "type": "text/javascript",
        "exec": [
          "// Generate X-Trace-ID for every request",
          "if (!pm.variables.get('trace_id')) {",
          "  pm.variables.set('trace_id', 'pm-' + Date.now() + '-' + Math.random().toString(36).slice(2, 10));",
          "}",
          "pm.request.headers.upsert({ key: 'X-Trace-ID', value: pm.variables.get('trace_id') });"
        ]
      }
    },
    {
      "listen": "test",
      "script": {
        "type": "text/javascript",
        "exec": [
          "// Common assertions",
          "pm.test('Response time < 2s', () => pm.expect(pm.response.responseTime).to.be.below(2000));",
          "pm.test('Has content-type json', () => {",
          "  if (pm.response.code !== 204) {",
          "    pm.expect(pm.response.headers.get('Content-Type')).to.include('json');",
          "  }",
          "});"
        ]
      }
    }
  ],
  "variable": [
    { "key": "base_url", "value": "http://localhost:8080", "type": "string" },
    { "key": "access_token", "value": "", "type": "string" },
    { "key": "user_id", "value": "", "type": "string" },
    { "key": "dsar_id", "value": "", "type": "string" }
  ],
  "item": [
    {
      "name": "🔐 Auth",
      "item": [
        {
          "name": "Login (get token)",
          "event": [
            {
              "listen": "test",
              "script": {
                "type": "text/javascript",
                "exec": [
                  "pm.test('200 OK', () => pm.response.to.have.status(200));",
                  "const json = pm.response.json();",
                  "pm.collectionVariables.set('access_token', json.access_token);",
                  "if (json.user_id) pm.collectionVariables.set('user_id', json.user_id);"
                ]
              }
            }
          ],
          "request": {
            "method": "POST",
            "header": [{ "key": "Content-Type", "value": "application/json" }],
            "body": {
              "mode": "raw",
              "raw": "{\n  \"email\": \"test@example.com\",\n  \"password\": \"{{test_password}}\"\n}"
            },
            "url": { "raw": "{{base_url}}/api/v1/auth/login", "host": ["{{base_url}}"], "path": ["api","v1","auth","login"] }
          }
        }
      ]
    },
    {
      "name": "✅ Consent — Write",
      "item": [
        {
          "name": "Record consent",
          "event": [
            {
              "listen": "test",
              "script": {
                "type": "text/javascript",
                "exec": [
                  "pm.test('200 OK', () => pm.response.to.have.status(200));",
                  "const json = pm.response.json();",
                  "pm.test('Has purposes array', () => pm.expect(json.purposes).to.be.an('array'));",
                  "pm.test('Has ids array', () => pm.expect(json.ids).to.be.an('array'));",
                  "pm.test('Contains NECESSARY', () => pm.expect(json.purposes).to.include('NECESSARY'));"
                ]
              }
            }
          ],
          "request": {
            "method": "POST",
            "header": [{ "key": "Content-Type", "value": "application/json" }],
            "body": {
              "mode": "raw",
              "raw": "{\n  \"purposes\": {\n    \"NECESSARY\": true,\n    \"ANALYTICS\": true,\n    \"MARKETING\": false\n  },\n  \"session_id\": \"postman-{{$timestamp}}\"\n}"
            },
            "url": { "raw": "{{base_url}}/api/v1/pdpa/consent", "host": ["{{base_url}}"], "path": ["api","v1","pdpa","consent"] }
          }
        },
        {
          "name": "Revoke consent (ANALYTICS)",
          "event": [
            {
              "listen": "test",
              "script": {
                "type": "text/javascript",
                "exec": [
                  "pm.test('200 OK', () => pm.response.to.have.status(200));"
                ]
              }
            }
          ],
          "request": {
            "method": "DELETE",
            "header": [{ "key": "Content-Type", "value": "application/json" }],
            "body": { "mode": "raw", "raw": "{\n  \"purpose\": \"ANALYTICS\"\n}" },
            "url": { "raw": "{{base_url}}/api/v1/pdpa/consent", "host": ["{{base_url}}"], "path": ["api","v1","pdpa","consent"] }
          }
        },
        {
          "name": "Revoke NECESSARY (expect 400)",
          "event": [
            {
              "listen": "test",
              "script": {
                "type": "text/javascript",
                "exec": [
                  "pm.test('400 Bad Request', () => pm.response.to.have.status(400));",
                  "pm.test('Error mentions mandatory', () => pm.expect(pm.response.json().error).to.match(/mandatory/i));"
                ]
              }
            }
          ],
          "request": {
            "method": "DELETE",
            "header": [{ "key": "Content-Type", "value": "application/json" }],
            "body": { "mode": "raw", "raw": "{\n  \"purpose\": \"NECESSARY\"\n}" },
            "url": { "raw": "{{base_url}}/api/v1/pdpa/consent", "host": ["{{base_url}}"], "path": ["api","v1","pdpa","consent"] }
          }
        }
      ]
    },
    {
      "name": "📖 Consent — Read",
      "item": [
        {
          "name": "Consent history",
          "event": [
            {
              "listen": "test",
              "script": {
                "type": "text/javascript",
                "exec": [
                  "pm.test('200 OK', () => pm.response.to.have.status(200));",
                  "pm.test('Is array', () => pm.expect(pm.response.json()).to.be.an('array'));"
                ]
              }
            }
          ],
          "request": {
            "method": "GET",
            "url": { "raw": "{{base_url}}/api/v1/pdpa/consent/history", "host": ["{{base_url}}"], "path": ["api","v1","pdpa","consent","history"] }
          }
        },
        {
          "name": "Consent status",
          "request": {
            "method": "GET",
            "url": { "raw": "{{base_url}}/api/v1/pdpa/consent/status", "host": ["{{base_url}}"], "path": ["api","v1","pdpa","consent","status"] }
          }
        }
      ]
    },
    {
      "name": "📝 DSAR",
      "item": [
        {
          "name": "Submit DSAR (ACCESS)",
          "event": [
            {
              "listen": "test",
              "script": {
                "type": "text/javascript",
                "exec": [
                  "pm.test('201 Created', () => pm.response.to.have.status(201));",
                  "const json = pm.response.json();",
                  "pm.test('Has id', () => pm.expect(json.id).to.be.a('string'));",
                  "pm.test('Status is PENDING', () => pm.expect(json.status).to.eql('PENDING'));",
                  "pm.collectionVariables.set('dsar_id', json.id);"
                ]
              }
            }
          ],
          "request": {
            "method": "POST",
            "header": [{ "key": "Content-Type", "value": "application/json" }],
            "body": { "mode": "raw", "raw": "{\n  \"request_type\": \"ACCESS\"\n}" },
            "url": { "raw": "{{base_url}}/api/v1/pdpa/dsar", "host": ["{{base_url}}"], "path": ["api","v1","pdpa","dsar"] }
          }
        },
        {
          "name": "Get DSAR status",
          "request": {
            "method": "GET",
            "url": { "raw": "{{base_url}}/api/v1/pdpa/dsar/{{dsar_id}}", "host": ["{{base_url}}"], "path": ["api","v1","pdpa","dsar","{{dsar_id}}"] }
          }
        },
        {
          "name": "List DSAR",
          "request": {
            "method": "GET",
            "url": { "raw": "{{base_url}}/api/v1/pdpa/dsar?limit=20", "host": ["{{base_url}}"], "path": ["api","v1","pdpa","dsar"], "query": [{ "key": "limit", "value": "20" }] }
          }
        },
        {
          "name": "Verify OTP",
          "request": {
            "method": "POST",
            "header": [{ "key": "Content-Type", "value": "application/json" }],
            "body": { "mode": "raw", "raw": "{\n  \"otp_code\": \"123456\"\n}" },
            "url": { "raw": "{{base_url}}/api/v1/pdpa/dsar/{{dsar_id}}/verify-otp", "host": ["{{base_url}}"], "path": ["api","v1","pdpa","dsar","{{dsar_id}}","verify-otp"] }
          }
        },
        {
          "name": "Process DSAR (admin)",
          "request": {
            "method": "POST",
            "header": [{ "key": "Content-Type", "value": "application/json" }],
            "body": { "mode": "raw", "raw": "{\n  \"action\": \"APPROVE\"\n}" },
            "url": { "raw": "{{base_url}}/api/v1/pdpa/dsar/{{dsar_id}}/process", "host": ["{{base_url}}"], "path": ["api","v1","pdpa","dsar","{{dsar_id}}","process"] }
          }
        }
      ]
    },
    {
      "name": "👤 Account (Admin)",
      "item": [
        {
          "name": "Suspend account",
          "request": {
            "method": "POST",
            "header": [{ "key": "Content-Type", "value": "application/json" }],
            "body": { "mode": "raw", "raw": "{\n  \"user_id\": \"{{user_id}}\",\n  \"suspended_at\": \"{{$isoTimestamp}}\",\n  \"retention_years\": 1\n}" },
            "url": { "raw": "{{base_url}}/api/v1/pdpa/account/suspend", "host": ["{{base_url}}"], "path": ["api","v1","pdpa","account","suspend"] }
          }
        },
        {
          "name": "Terminate account",
          "request": {
            "method": "POST",
            "header": [{ "key": "Content-Type", "value": "application/json" }],
            "body": { "mode": "raw", "raw": "{\n  \"user_id\": \"{{user_id}}\",\n  \"terminated_at\": \"{{$isoTimestamp}}\"\n}" },
            "url": { "raw": "{{base_url}}/api/v1/pdpa/account/terminate", "host": ["{{base_url}}"], "path": ["api","v1","pdpa","account","terminate"] }
          }
        },
        {
          "name": "Confirm deletion",
          "request": {
            "method": "POST",
            "header": [{ "key": "Content-Type", "value": "application/json" }],
            "body": { "mode": "raw", "raw": "{\n  \"user_id\": \"{{user_id}}\"\n}" },
            "url": { "raw": "{{base_url}}/api/v1/pdpa/account/confirm-deletion", "host": ["{{base_url}}"], "path": ["api","v1","pdpa","account","confirm-deletion"] }
          }
        }
      ]
    },
    {
      "name": "📜 Policy",
      "item": [
        {
          "name": "Get active policy",
          "request": {
            "method": "GET",
            "auth": { "type": "noauth" },
            "url": { "raw": "{{base_url}}/api/v1/pdpa/policy", "host": ["{{base_url}}"], "path": ["api","v1","pdpa","policy"] }
          }
        },
        {
          "name": "List policy versions",
          "request": {
            "method": "GET",
            "auth": { "type": "noauth" },
            "url": { "raw": "{{base_url}}/api/v1/pdpa/policy/versions", "host": ["{{base_url}}"], "path": ["api","v1","pdpa","policy","versions"] }
          }
        },
        {
          "name": "Publish policy (admin)",
          "request": {
            "method": "POST",
            "header": [{ "key": "Content-Type", "value": "application/json" }],
            "body": { "mode": "raw", "raw": "{\n  \"version\": \"v1.2.0\",\n  \"title\": \"Privacy Policy v1.2.0\",\n  \"content\": \"...\",\n  \"effective_date\": \"{{$isoTimestamp}}\"\n}" },
            "url": { "raw": "{{base_url}}/api/v1/pdpa/policy", "host": ["{{base_url}}"], "path": ["api","v1","pdpa","policy"] }
          }
        }
      ]
    },
    {
      "name": "📊 Admin",
      "item": [
        {
          "name": "Reports",
          "request": {
            "method": "GET",
            "url": {
              "raw": "{{base_url}}/api/v1/pdpa/admin/reports?from=2026-01-01T00:00:00Z&to=2026-01-31T23:59:59Z",
              "host": ["{{base_url}}"],
              "path": ["api","v1","pdpa","admin","reports"],
              "query": [
                { "key": "from", "value": "2026-01-01T00:00:00Z" },
                { "key": "to",   "value": "2026-01-31T23:59:59Z" }
              ]
            }
          }
        },
        {
          "name": "Audit trails",
          "request": {
            "method": "GET",
            "url": { "raw": "{{base_url}}/api/v1/pdpa/admin/audit-trails?limit=100", "host": ["{{base_url}}"], "path": ["api","v1","pdpa","admin","audit-trails"], "query": [{ "key": "limit", "value": "100" }] }
          }
        },
        {
          "name": "Statistics",
          "request": {
            "method": "GET",
            "url": { "raw": "{{base_url}}/api/v1/pdpa/admin/statistics", "host": ["{{base_url}}"], "path": ["api","v1","pdpa","admin","statistics"] }
          }
        }
      ]
    }
  ]
}
```

### 21.2 Environments

**`api/postman/PDPA-local.postman_environment.json`**

```json
{
  "name": "PDPA Local",
  "values": [
    { "key": "base_url", "value": "http://localhost:8080", "enabled": true },
    { "key": "test_password", "value": "Test@1234", "type": "secret", "enabled": true }
  ],
  "_postman_variable_scope": "environment"
}
```

**`api/postman/PDPA-production.postman_environment.json`**

```json
{
  "name": "PDPA Production",
  "values": [
    { "key": "base_url", "value": "https://api.icmongolang.dev", "enabled": true },
    { "key": "test_password", "value": "", "type": "secret", "enabled": false }
  ],
  "_postman_variable_scope": "environment"
}
```

### 21.3 Newman CI Runner

**`.github/workflows/pdpa-postman.yml`**

```yaml
name: PDPA Postman Smoke Tests

on:
  workflow_dispatch:
  schedule:
    - cron: '0 */6 * * *'   # every 6 hours

jobs:
  smoke:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-node@v4
        with: { node-version: '20' }

      - name: Install Newman
        run: npm i -g newman newman-reporter-htmlextra

      - name: Run collection (staging)
        env:
          BASE_URL: https://staging-api.icmongolang.dev
          TEST_PASSWORD: ${{ secrets.TEST_PASSWORD }}
        run: |
          newman run api/postman/PDPA.postman_collection.json \
            -e api/postman/PDPA-local.postman_environment.json \
            --env-var "base_url=$BASE_URL" \
            --env-var "test_password=$TEST_PASSWORD" \
            --reporters cli,htmlextra,junit \
            --reporter-htmlextra-export newman-report.html \
            --reporter-junit-export newman-report.xml \
            --bail \
            --timeout-request 5000

      - uses: actions/upload-artifact@v4
        if: always()
        with:
          name: newman-report
          path: |
            newman-report.html
            newman-report.xml
```

---

## 22. Runbook (SRE Operations Guide)

**`docs/runbook/pdpa-runbook.md`**

```markdown
# PDPA Module — SRE Runbook

> Operational runbook สำหรับ PDPA Compliance Module
> Last updated: 2026-01-15
> Owners: PDPA Team · SRE Team · Security Team

---

## 📞 Escalation Matrix

| Severity | Primary | Secondary | Response SLA |
|----------|---------|-----------|--------------|
| **P0** (down/SLA breach) | On-call SRE | PDPA Lead | 5 min |
| **P1** (degraded) | On-call SRE | PDPA Engineer | 15 min |
| **P2** (partial) | PDPA Engineer | — | 1 hour |
| **P3** (cosmetic) | PDPA Engineer | — | Next business day |

**Slack channels:** `#pdpa-alerts` · `#pdpa-oncall` · `#pdpa-deployments`
**PagerDuty:** `pdpa-p0`, `pdpa-p1`

---

## 1. Service Overview

### 1.1 Components

| Component | Purpose | Replicas (prod) | Criticality |
|-----------|---------|-----------------|-------------|
| `pdpa-api` | HTTP API + WebSocket | 5 (HPA 5-50) | 🔴 Critical |
| `pdpa-worker` | Kafka consumers | 4 (HPA 4-24) | 🔴 Critical |
| `pdpa-scheduler` | CronJob (daily cleanup) | — | 🟡 Important |
| Postgres RDS | Primary DB | Multi-AZ | 🔴 Critical |
| ElastiCache Redis | Cache + idempotency | 3 shards × 2 replicas | 🔴 Critical |
| MSK Kafka | Event bus | 6 brokers | 🔴 Critical |
| OpenSearch | Search index | 3 nodes | 🟢 Optional |

### 1.2 SLA

- **Availability:** 99.95% monthly (≤ 22 min/month downtime)
- **API p95 latency:** < 500ms
- **DSAR response SLA:** ≤ 30 days (PDPA requirement)
- **Data deletion SLA:** ≤ 24 hours after confirmation

### 1.3 Key URLs

- **API:** https://api.icmongolang.dev
- **Swagger:** https://api.icmongolang.dev/api/v1/pdpa/docs
- **Grafana:** https://grafana.icmongolang.dev/d/pdpa-overview
- **ArgoCD:** https://argocd.icmongolang.dev/applications/pdpa-production
- **Kafka UI:** https://kafka-ui.icmongolang.dev
- **Runbook:** https://runbooks.icmongolang.dev/pdpa

---

## 2. Health Checks

### 2.1 Quick Health

```bash
# API liveness/readiness
curl -sS https://api.icmongolang.dev/healthz
curl -sS https://api.icmongolang.dev/readyz

# K8s pods
kubectl -n pdpa get pods -l app.kubernetes.io/name=pdpa

# ArgoCD app status
argocd app get pdpa-production
```

### 2.2 Deep Health

```bash
# DB
kubectl -n pdpa exec deploy/pdpa-api -- \
  pg_isready -h $DB_HOST -p 5432

# Redis
kubectl -n pdpa exec deploy/pdpa-api -- \
  redis-cli -h $REDIS_HOST -a $REDIS_AUTH PING

# Kafka
kubectl -n pdpa exec deploy/pdpa-worker -- \
  kafka-broker-api-versions --bootstrap-server $KAFKA_BROKERS

# Outbox backlog
kubectl -n pdpa exec deploy/pdpa-api -- sh -c \
  'psql "$DB_DSN" -c "SELECT count(*) FROM pdpa_outbox WHERE status='"'"'PENDING'"'"';"'
```

---

## 3. Common Incidents & Remediation

### 3.1 🔴 P0: API 5xx Spike

**Symptoms:** `PDPA_High_Error_Rate` alert fires (>5% error rate for 5m)

**Diagnosis:**
```bash
# 1. Check error breakdown
kubectl -n pdpa logs -l app.kubernetes.io/component=api --tail=500 \
  | jq -r 'select(.level=="error") | .error' | sort | uniq -c | sort -rn

# 2. Check dependency status
kubectl -n pdpa exec deploy/pdpa-api -- sh -c \
  'psql "$DB_DSN" -c "SELECT 1"'
kubectl -n pdpa exec deploy/pdpa-api -- redis-cli -h $REDIS_HOST PING

# 3. Check recent deployments
argocd app history pdpa-production
```

**Remediation:**
```bash
# Rollback if recent deploy
argocd app rollback pdpa-production <previous-revision>

# Scale up if resource-bound
kubectl -n pdpa scale deploy/pdpa-api --replicas=10

# Restart if memory leak suspected
kubectl -n pdpa rollout restart deploy/pdpa-api
```

### 3.2 🔴 P0: Outbox Stuck (no events published)

**Symptoms:** `PDPA_Outbox_Stuck` alert — pending > 0 but no publish rate for 5m

**Diagnosis:**
```bash
# 1. Check worker pods
kubectl -n pdpa get pods -l app.kubernetes.io/component=worker

# 2. Check pending count
kubectl -n pdpa exec deploy/pdpa-api -- sh -c \
  'psql "$DB_DSN" -c "SELECT status, count(*) FROM pdpa_outbox GROUP BY status;"'

# 3. Check oldest pending
kubectl -n pdpa exec deploy/pdpa-api -- sh -c \
  'psql "$DB_DSN" -c "SELECT id, topic, retry_count, last_error, created_at FROM pdpa_outbox WHERE status='"'"'PENDING'"'"' ORDER BY created_at LIMIT 5;"'

# 4. Check DLQ topics
kafka-console-consumer --bootstrap-server $KAFKA_BROKERS \
  --topic pdpa.dlq.pdpa.consent.granted --from-beginning --max-messages 5

# 5. Check worker logs
kubectl -n pdpa logs -l app.kubernetes.io/component=worker --tail=200 \
  | grep -i "outbox\|error"
```

**Remediation:**
```bash
# 1. Restart worker if hung
kubectl -n pdpa rollout restart deploy/pdpa-worker

# 2. Manually re-queue failed events (if status=FAILED and error is transient)
kubectl -n pdpa exec deploy/pdpa-api -- sh -c \
  'psql "$DB_DSN" -c "UPDATE pdpa_outbox SET status='"'"'PENDING'"'"', retry_count=0, last_error=NULL WHERE status='"'"'FAILED'"'"' AND retry_count < 5;"'

# 3. If DB connection pool exhausted
kubectl -n pdpa exec deploy/pdpa-api -- sh -c \
  'psql "$DB_DSN" -c "SELECT count(*), state FROM pg_stat_activity WHERE datname = current_database() GROUP BY state;"'
```

### 3.3 🟠 P1: Consumer Lag High

**Symptoms:** `PDPA_Consumer_Lag_High` (>5000 messages lag)

**Diagnosis:**
```bash
# 1. Which consumer group is lagging?
kubectl -n pdpa exec deploy/pdpa-worker -- \
  kafka-consumer-groups --bootstrap-server $KAFKA_BROKERS --all-groups --describe \
  | grep pdpa | awk '$5>1000'

# 2. Are consumers alive?
kubectl -n pdpa get pods -l app.kubernetes.io/component=worker

# 3. Check for poison messages
kubectl -n pdpa logs -l app.kubernetes.io/component=worker --tail=500 \
  | grep "handler failed"
```

**Remediation:**
```bash
# 1. Scale workers (Kafka จะ rebalance partitions อัตโนมัติ)
kubectl -n pdpa scale deploy/pdpa-worker --replicas=12

# 2. If specific handler is slow — check external deps
#    e.g., LLM consumer → LLM API rate limit
kubectl -n pdpa logs -l app.kubernetes.io/component=worker \
  | grep "llm" | tail -50

# 3. Increase topic partitions (long-term)
kafka-topics --bootstrap-server $KAFKA_BROKERS \
  --alter --topic pdpa.llm.analysis.requested --partitions 12
```

### 3.4 🟠 P1: Redis Down

**Symptoms:** Cache miss rate >90%, `pdpa_consent_cache_hits_total` drops

**Diagnosis:**
```bash
# Ping Redis
kubectl -n pdpa exec deploy/pdpa-api -- \
  redis-cli -h $REDIS_HOST -p 6379 -a $REDIS_AUTH ping

# Check ElastiCache cluster health
aws elasticache describe-replication-groups \
  --replication-group-id icmongolang-prod-pdpa-redis \
  | jq '.ReplicationGroups[0].NodeGroups[].NodeGroupMembers[] | {id: .CacheNodeId, status: .CurrentRole}'
```

**Remediation:**
- Redis is **cache + idempotency** — API ยังทำงานได้แต่ช้า
- If failover in progress → รอ 1-2 นาที (multi-AZ auto-failover)
- If all nodes down → **critical**: contact AWS Support, consider disable cache temporarily via feature flag

### 3.5 🟠 P1: DSAR Backlog

**Symptoms:** `PDPA_DSAR_Backlog` or `PDPA_DSAR_SLA_Breach`

**Diagnosis:**
```bash
# Pending DSARs
kubectl -n pdpa exec deploy/pdpa-api -- sh -c \
  'psql "$DB_DSN" -c "SELECT status, count(*), min(requested_at) FROM pdpa_dsar_requests GROUP BY status;"'

# SLA breach candidates (>30 days pending)
kubectl -n pdpa exec deploy/pdpa-api -- sh -c \
  'psql "$DB_DSN" -c "SELECT id, user_id, request_type, requested_at FROM pdpa_dsar_requests WHERE status='"'"'PENDING'"'"' AND requested_at < NOW() - INTERVAL '"'"'30 days'"'"';"'
```

**Remediation:**
- **Business escalation** → notify DPO (Data Protection Officer)
- Manual processing via admin API if automated pipeline is stuck
- Document root cause in incident report (PDPA compliance requirement)

### 3.6 🟡 P2: Kafka DLQ Growing

**Diagnosis:**
```bash
# DLQ counts per original topic
kubectl -n pdpa exec deploy/pdpa-worker -- \
  kafka-run-class kafka.tools.GetOffsetShell \
  --broker-list $KAFKA_BROKERS --topic pdpa.dlq.pdpa.consent.granted \
  --time -1

# Peek messages
kubectl -n pdpa exec deploy/pdpa-worker -- \
  kafka-console-consumer --bootstrap-server $KAFKA_BROKERS \
  --topic pdpa.dlq.pdpa.consent.granted --from-beginning --max-messages 10
```

**Remediation:**
```bash
# 1. Identify common error
# 2. Fix root cause (bug, schema mismatch, external dep)
# 3. Replay from DLQ back to original topic
kafka-console-consumer --bootstrap-server $KAFKA_BROKERS \
  --topic pdpa.dlq.pdpa.consent.granted --from-beginning \
  | kafka-console-producer --bootstrap-server $KAFKA_BROKERS \
      --topic pdpa.consent.granted
```

### 3.7 🟡 P2: WebSocket disconnect burst

**Symptoms:** `PDPA_WS_Disconnected_Burst` (>50/s)

**Diagnosis:**
```bash
# Check ingress pod (nginx) logs
kubectl -n ingress-nginx logs -l app.kubernetes.io/name=ingress-nginx --tail=500 \
  | grep "pdpa/ws" | grep -E "upstream|timeout"
```

**Remediation:**
```yaml
# Ensure ingress has WebSocket-friendly timeouts
nginx.ingress.kubernetes.io/proxy-read-timeout: "3600"
nginx.ingress.kubernetes.io/proxy-send-timeout: "3600"
nginx.ingress.kubernetes.io/proxy-connect-timeout: "60"
```

---

## 4. Routine Operations

### 4.1 Deploy (GitOps)

```bash
# 1. Dev: promote develop → staging (auto via ArgoCD)
git push origin develop

# 2. Staging: verify
argocd app wait pdpa-staging --health --timeout 300

# 3. Prod: create tag
git tag pdpa/v1.2.3 && git push origin pdpa/v1.2.3
# → GitHub Actions builds image + signs
# → ArgoCD detects main branch change
# → Manual sync

argocd app sync pdpa-production
argocd app wait pdpa-production --health --timeout 600
```

### 4.2 Rollback

```bash
# 1. List history
argocd app history pdpa-production

# 2. Rollback
argocd app rollback pdpa-production <REVISION>

# 3. Verify
kubectl -n pdpa rollout status deploy/pdpa-api
curl -sS https://api.icmongolang.dev/readyz
```

### 4.3 Database Migration

```bash
# ⚠️ Always take a snapshot first
aws rds create-db-snapshot \
  --db-instance-identifier icmongolang-prod-pdpa \
  --db-snapshot-identifier pdpa-pre-migration-$(date +%Y%m%d-%H%M)

# Apply migration (idempotent — uses IF NOT EXISTS)
kubectl -n pdpa create job --from=cronjob/pdpa-migrate migrate-$(date +%s)

# Verify
kubectl -n pdpa exec deploy/pdpa-api -- sh -c \
  'psql "$DB_DSN" -c "\dt pdpa_*"'
```

### 4.4 Scaling

```bash
# Manual scale (HPA จะ override ในภายหลัง)
kubectl -n pdpa scale deploy/pdpa-api --replicas=10

# Or update HPA min temporarily
kubectl -n pdpa patch hpa pdpa-api --type merge \
  -p '{"spec":{"minReplicas":10}}'
```

### 4.5 Secret Rotation

```bash
# 1. Update in AWS Secrets Manager
aws secretsmanager update-secret \
  --secret-id icmongolang-prod/pdpa/db \
  --secret-string '{"username":"...","password":"..."}'

# 2. External Secrets Operator auto-syncs within 1h (refreshInterval)
# Force refresh:
kubectl -n pdpa annotate externalsecret pdpa-postgres-prod \
  force-sync=$(date +%s) --overwrite

# 3. Rolling restart to pick up new secret
kubectl -n pdpa rollout restart deploy/pdpa-api deploy/pdpa-worker
```

### 4.6 Kafka Topic Management

```bash
# Create topic
kafka-topics --bootstrap-server $KAFKA_BROKERS \
  --create --topic pdpa.new.topic --partitions 6 --replication-factor 3

# Increase partitions
kafka-topics --bootstrap-server $KAFKA_BROKERS \
  --alter --topic pdpa.consent.granted --partitions 12

# Describe
kafka-topics --bootstrap-server $KAFKA_BROKERS --describe \
  --topic pdpa.consent.granted

# Reset consumer group offset (⚠️ DANGER)
kafka-consumer-groups --bootstrap-server $KAFKA_BROKERS \
  --group pdpa-consent-granted-cg --reset-offsets --to-earliest \
  --topic pdpa.consent.granted --execute
```

### 4.7 DSAR Manual Processing

```bash
# 1. List pending DSARs
curl -sS https://api.icmongolang.dev/api/v1/pdpa/dsar?status=PENDING \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq .

# 2. Approve (admin)
curl -X POST https://api.icmongolang.dev/api/v1/pdpa/dsar/$DSAR_ID/process \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"action":"APPROVE"}'

# 3. Verify
curl -sS https://api.icmongolang.dev/api/v1/pdpa/dsar/$DSAR_ID \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq .status
```

---

## 5. Monitoring & Alerts

### 5.1 Key Metrics

| Metric | Description | Alert Threshold |
|--------|-------------|-----------------|
| `pdpa_consent_recorded_total` | Consent granted rate | — |
| `pdpa_dsar_submitted_total` | DSAR submit rate | — |
| `pdpa_outbox_pending_count` | Unpublished events | >1000 warn, >5000 critical |
| `kafka_consumergroup_lag` | Consumer lag | >5000 warn, >50000 crit |
| `pdpa_dlq_total` | DLQ ingress | >10/15m warn |
| `pdpa_dsar_pending_total{status="PENDING"}` | DSAR backlog | >100 warn |
| `http_request_duration_seconds` | API latency | p95 > 1s warn |
| `pdpa_ws_connected_users` | WS connections | — |

### 5.2 Dashboards

- **Overview:** https://grafana.icmongolang.dev/d/pdpa-overview
- **Consumers:** https://grafana.icmongolang.dev/d/pdpa-consumers
- **Business:** https://grafana.icmongolang.dev/d/pdpa-business

### 5.3 Log Queries (Loki)

```logql
# All errors in last hour
{namespace="pdpa"} |= "error" | json | line_format "{{.error}}"

# Slow requests
{namespace="pdpa", app="pdpa-api"} | json | http_request_duration > 1.0

# Consent grants by user
{namespace="pdpa"} | json | event_type="CONSENT_GRANTED" | line_format "{{.user_id}}"

# DSAR flow trace
{namespace="pdpa"} | json | correlation_id="<CORR-ID>"
```

---

## 6. Disaster Recovery

### 6.1 RTO / RPO

| Component | RTO | RPO | Backup |
|-----------|-----|-----|--------|
| Postgres | 15 min | 5 min | RDS PITR (90 days) + snapshots |
| Redis | 10 min | 1 min | ElastiCache snapshot (daily) |
| Kafka | 30 min | — | MSK retention 7 days |
| OpenSearch | 1 h | 24 h | Manual snapshot |

### 6.2 Full Region Failure (ap-southeast-1 → ap-northeast-1)

```bash
# 1. Restore RDS from snapshot to DR region
aws rds restore-db-instance-from-db-snapshot \
  --db-instance-identifier icmongolang-dr-pdpa \
  --db-snapshot-identifier <latest-snapshot> \
  --region ap-northeast-1

# 2. Update DNS (Route53 failover)
aws route53 change-resource-record-sets --hosted-zone-id $ZONE_ID --change-batch file://dr-failover.json

# 3. Deploy app to DR cluster
argocd app sync pdpa-dr --server https://dr-cluster.example.com
```

### 6.3 Data Corruption (bad migration)

```bash
# 1. Stop writes — scale workers to 0
kubectl -n pdpa scale deploy/pdpa-worker --replicas=0

# 2. Restore from PITR to a new instance
aws rds restore-db-instance-to-point-in-time \
  --source-db-instance-identifier icmongolang-prod-pdpa \
  --target-db-instance-identifier pdpa-recovery \
  --restore-time 2026-01-15T09:00:00Z

# 3. Verify data
psql -h pdpa-recovery.xxx.rds.amazonaws.com -c "SELECT count(*) FROM pdpa_consents;"

# 4. Swap endpoint / update secret / rolling restart
# 5. Scale workers back up
kubectl -n pdpa scale deploy/pdpa-worker --replicas=4
```

---

## 7. Security Incidents

### 7.1 Suspected Data Breach

**Immediate (0-1h):**
1. Isolate affected pods: `kubectl -n pdpa scale deploy/pdpa-api --replicas=0`
2. Snapshot forensics: `kubectl logs --since=24h > incident-logs.txt`
3. Preserve evidence: take RDS snapshot `pdpa-incident-$(date)`
4. Notify **PDPC within 72h** (PDPA Section 37(4))

**Investigation (1-24h):**
- Query audit trail: `SELECT * FROM pdpa_audit_trails WHERE created_at > NOW() - INTERVAL '24 hours'`
- Check blockchain records (immutable hash verification)
- Review IAM/IRSA role usage in CloudTrail

**Remediation:**
- Rotate all secrets
- Patch vulnerability
- Post-incident report to PDPC + affected users

### 7.2 Secret Leaked

```bash
# 1. Rotate immediately
aws secretsmanager rotate-secret --secret-id icmongolang-prod/pdpa/app

# 2. Force sync + rolling restart
kubectl -n pdpa annotate externalsecret pdpa-app-prod force-sync=$(date +%s) --overwrite
kubectl -n pdpa rollout restart deploy/pdpa-api deploy/pdpa-worker

# 3. Revoke compromised JWT tokens (invalidate all)
kubectl -n pdpa exec deploy/pdpa-api -- sh -c \
  'redis-cli -h $REDIS_HOST -a $REDIS_AUTH --scan --pattern "auth:token:*" | xargs redis-cli DEL'
```

---

## 8. Compliance & Audit

### 8.1 Monthly PDPA Report

```bash
# Generate report
curl -sS "https://api.icmongolang.dev/api/v1/pdpa/admin/reports?from=$(date -d '1 month ago' -Iseconds)&to=$(date -Iseconds)" \
  -H "Authorization: Bearer $ADMIN_TOKEN" > monthly-report.json

# Required sections:
# - Consent stats by purpose
# - DSAR submitted/completed/rejected
# - Data deletions (immediate/auto)
# - Policy versions published
# - Audit events
```

### 8.2 DSAR SLA Compliance

```sql
-- DSARs approaching SLA (25-30 days)
SELECT id, user_id, request_type, requested_at,
       EXTRACT(DAY FROM (NOW() - requested_at)) AS days_pending
FROM pdpa_dsar_requests
WHERE status IN ('PENDING', 'PROCESSING')
  AND requested_at < NOW() - INTERVAL '25 days'
ORDER BY requested_at ASC;
```

### 8.3 Audit Trail Export (for PDPC)

```sql
-- Export audit trail for a specific user
COPY (
  SELECT action, details, ip_address, created_at
  FROM pdpa_audit_trails
  WHERE user_id = '<UUID>'
  ORDER BY created_at
) TO '/tmp/audit-<UUID>.csv' WITH CSV HEADER;
```

---

## 9. Contacts & References

| Resource | Link |
|----------|------|
| PDPA Team | #pdpa-team |
| On-call SRE | #pdpa-oncall |
| DPO (Data Protection Officer) | dpo@icmongolang.dev |
| Security Team | security@icmongolang.dev |
| AWS Support | console.aws.amazon.com/support |
| PDPC (Thai regulator) | https://www.pdpc.or.th |

**Runbook repository:** https://github.com/icmongolang/runbooks/tree/main/pdpa
**Postmortem template:** https://github.com/icmongolang/runbooks/tree/main/templates
```

---

## สรุปส่วนที่ 5

| # | Component | ไฟล์ | รายละเอียด |
|---|-----------|------|-----------|
| **18** | Terraform modules | ~15 ไฟล์ | RDS (Multi-AZ + encryption), ElastiCache (cluster mode), MSK (mTLS + SCRAM), Secrets, IRSA |
| **19** | ArgoCD GitOps | ~8 ไฟล์ | AppProject, App-of-Apps, Staging (auto), Production (manual), ExternalSecrets, Notifications |
| **20** | Load tests | ~10 ไฟล์ | k6 (5 scenarios: write/read/dsar/soak/ws) + Locust + CI integration |
| **21** | Postman | 3 ไฟล์ | Collection (24 requests, auto-trace-id, test scripts) + Environments + Newman CI |
| **22** | Runbook | 1 ไฟล์ | P0-P3 incidents, routine ops, DR, security, compliance |

**รวมทั้ง 5 ส่วน:** ~175 ไฟล์

### Final Production Readiness Matrix

| ด้าน | สถานะ | หลักฐาน |
|------|-------|--------|
| **Architecture** | ✅ | Clean Arch + DDD + EDA + Outbox |
| **Security** | ✅ | Distroless + non-root + cosign + Trivy + IRSA + mTLS |
| **Observability** | ✅ | Metrics + Logs + Traces + 3 Grafana dashboards + 12 alerts |
| **Reliability** | ✅ | Multi-AZ + PDB + HPA + retry + circuit breaker + DLQ |
| **Scalability** | ✅ | HPA (API 5-50, Worker 4-24) + Kafka partitions |
| **IaC** | ✅ | Terraform (RDS/MSK/Redis/Secrets/IAM) |
| **GitOps** | ✅ | ArgoCD App-of-Apps + auto-sync staging + manual prod |
| **Testing** | ✅ | Unit + Integration + E2E + k6 + Locust + Postman |
| **Compliance** | ✅ | Consent log + audit trail + blockchain + retention policy + DSAR SLA |
| **Operations** | ✅ | Runbook + escalation + DR + secret rotation |

 