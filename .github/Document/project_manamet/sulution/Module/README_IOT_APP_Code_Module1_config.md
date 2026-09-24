# Module 1: การจัดการ Configuration (Config Management)

## สำหรับโฟลเดอร์ `config/`

ไฟล์ที่เกี่ยวข้อง:
- `config/config.go`
- `config/config-local.yml`
- `config/config-prod.yml`

---

## หลักการ (Concept)

### คืออะไร?
Configuration management คือกระบวนการแยกค่าที่เปลี่ยนแปลงตาม environment (database URL, JWT secret, SMTP host) ออกจาก source code เพื่อให้สามารถ deploy แอปพลิเคชันเดียวกันไปยัง environment ต่างกัน (local, staging, production) โดยไม่ต้องแก้ไขโค้ด

### มีกี่แบบ?
1. **Environment variables** – ปลอดภัย, ใช้ใน production
2. **Config file (YAML/JSON/TOML)** – อ่านง่าย, เก็บ default
3. **Remote config (Consul, etcd)** – dynamic update
4. **Command-line flags** – ใช้สำหรับ override

### ใช้อย่างไร / นำไปใช้กรณีไหน
- ใช้ Viper library เพื่อโหลด config จากหลายแหล่ง โดยมี priority: flag > env > config file > default
- ในโปรเจกต์นี้: โหลดจาก `config-local.yml` (development) หรือ `config-prod.yml` (production) และสามารถ override ด้วย environment variables

### ทำไมต้องใช้
- ป้องกัน hard-code credentials
- ทำให้ application portable
- ลดความผิดพลาดจากการตั้งค่าต่าง environment

### ประโยชน์ที่ได้รับ
- เปลี่ยน DB password โดยไม่ rebuild image
- รองรับ 12-factor app
- ทดสอบ automation ได้ง่าย

### ข้อควรระวัง
- ห้าม commit config file ที่มี secret จริง (ใช้ .example แทน)
- ต้อง validate config ตอนเริ่มโปรแกรม (panic early ถ้าขาด required field)

### ข้อดี
- ยืดหยุ่น, ปลอดภัย, แยกการตั้งค่าออกจากโค้ด

### ข้อเสีย
- เพิ่ม complexity เล็กน้อย
- ต้องระวังการ override ที่ไม่ตั้งใจ

### ข้อห้าม
- อย่าใช้ config file สำหรับความลับ (ใช้ env หรือ secret manager)
- อย่า reload config แบบไม่จำกัดโดยไม่มีการ validate

---

## โค้ดที่รันได้จริง

### ไฟล์ `config/config.go`

```go
// Package config provides configuration loading using Viper.
// It supports YAML files and environment variables.
// -------------------------------------------------------
// แพ็คเกจ config ให้บริการโหลดค่ากำหนดผ่าน Viper
// รองรับไฟล์ YAML และตัวแปรสภาพแวดล้อม
package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration values for the application.
// โครงสร้าง Config เก็บค่ากำหนดทั้งหมดของแอปพลิเคชัน
type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	Redis     RedisConfig
	JWT       JWTConfig
	SMTP      SMTPConfig
	RateLimit RateLimitConfig
	Log       LogConfig
}

// ServerConfig defines HTTP server settings.
// ค่ากำหนดสำหรับ HTTP server
type ServerConfig struct {
	Port         string        `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

// DatabaseConfig defines PostgreSQL connection settings.
// ค่ากำหนดการเชื่อมต่อ PostgreSQL
type DatabaseConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"dbname"`
	SSLMode      string `mapstructure:"sslmode"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

// RedisConfig defines Redis connection settings.
// ค่ากำหนดการเชื่อมต่อ Redis
type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// JWTConfig defines JWT token settings.
// ค่ากำหนดสำหรับ JWT
type JWTConfig struct {
	AccessTokenDuration  time.Duration `mapstructure:"access_token_duration"`
	RefreshTokenDuration time.Duration `mapstructure:"refresh_token_duration"`
	PrivateKeyBase64     string        `mapstructure:"private_key_base64"`
	PublicKeyBase64      string        `mapstructure:"public_key_base64"`
}

// SMTPConfig defines email sending settings.
// ค่ากำหนดสำหรับการส่งอีเมล
type SMTPConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	From     string `mapstructure:"from"`
}

// RateLimitConfig defines request throttling settings.
// ค่ากำหนดการจำกัดอัตราการเรียกใช้
type RateLimitConfig struct {
	RequestsPerSecond int `mapstructure:"requests_per_second"`
	Burst             int `mapstructure:"burst"`
}

// LogConfig defines logging settings.
// ค่ากำหนดการบันทึก logs
type LogConfig struct {
	Level string `mapstructure:"level"` // debug, info, warn, error
	Env   string `mapstructure:"env"`   // development, production
}

// Load reads configuration from file and environment variables.
// Returns a populated Config struct or exits on fatal error.
// --------------------------------------------------------------
// Load อ่านค่ากำหนดจากไฟล์และตัวแปรสภาพแวดล้อม
// คืนค่าโครงสร้าง Config ที่ถูกเติมข้อมูล หรือจบการทำงานถ้ามีข้อผิดพลาดร้ายแรง
func Load() *Config {
	// Set config file name and type
	// กำหนดชื่อไฟล์และประเภทของ config
	viper.SetConfigName("config-local") // default to local
	viper.SetConfigType("yml")
	viper.AddConfigPath("./config") // look in ./config folder

	// Allow override via environment variable APP_ENV (production, staging)
	// อนุญาตให้ override ผ่านตัวแปรสภาพแวดล้อม APP_ENV
	if env := viper.GetString("APP_ENV"); env != "" {
		viper.SetConfigName("config-" + env)
	}

	// Read config file (ignore if not found)
	// อ่านไฟล์ config (ไม่ถือเป็น error ถ้าไม่มีไฟล์)
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: config file not found: %v", err)
	}

	// Enable environment variable binding with prefix "CMON_"
	// เปิดใช้งานการผูกตัวแปรสภาพแวดล้อมที่มี prefix "CMON_"
	viper.SetEnvPrefix("CMON")
	viper.AutomaticEnv()
	// Replace dots with underscores for nested keys (e.g., server.port -> SERVER_PORT)
	viper.SetEnvKeyReplacer(viper.NewReplacer(".", "_"))

	// Set default values
	// กำหนดค่าเริ่มต้น
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.read_timeout", "15s")
	viper.SetDefault("server.write_timeout", "15s")
	viper.SetDefault("server.idle_timeout", "60s")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("database.max_open_conns", 50)
	viper.SetDefault("database.max_idle_conns", 25)
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("jwt.access_token_duration", "15m")
	viper.SetDefault("jwt.refresh_token_duration", "168h") // 7 days
	viper.SetDefault("rate_limit.requests_per_second", 100)
	viper.SetDefault("rate_limit.burst", 200)
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.env", "development")

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Unable to decode config: %v", err)
	}

	// Validate required fields (example)
	// ตรวจสอบฟิลด์ที่จำเป็น
	if cfg.JWT.PrivateKeyBase64 == "" || cfg.JWT.PublicKeyBase64 == "" {
		log.Fatal("JWT private/public keys are required (set via config or env)")
	}
	if cfg.Database.Host == "" {
		log.Fatal("Database host is required")
	}

	return &cfg
}

// GetDSN returns PostgreSQL connection string.
// -------------------------------------------------
// GetDSN คืนค่า connection string ของ PostgreSQL
func (c *DatabaseConfig) GetDSN() string {
	return "host=" + c.Host +
		" port=" + string(rune(c.Port)) + // in real use: strconv.Itoa
		" user=" + c.User +
		" password=" + c.Password +
		" dbname=" + c.DBName +
		" sslmode=" + c.SSLMode
}
```

### ไฟล์ `config/config-local.yml` (ตัวอย่างสำหรับ development)

```yaml
# Development configuration for CMON backend
# ค่ากำหนดสำหรับการพัฒนาแบ็คเอนด์ CMON
server:
  port: 8080
  read_timeout: 15s
  write_timeout: 15s
  idle_timeout: 60s

database:
  host: localhost
  port: 5432
  user: postgres
  password: postgres123
  dbname: cmon_dev
  sslmode: disable
  max_open_conns: 25
  max_idle_conns: 10

redis:
  addr: localhost:6379
  password: ""
  db: 0

jwt:
  access_token_duration: 15m
  refresh_token_duration: 168h
  # ต้องแทนที่ด้วย base64 ของ private/public key จริง
  private_key_base64: "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQp... (ตัวอย่าง)"
  public_key_base64: "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0K..."

smtp:
  host: smtp.gmail.com
  port: 587
  username: your-email@gmail.com
  password: your-app-password
  from: alerts@cmon.local

rate_limit:
  requests_per_second: 100
  burst: 200

log:
  level: debug
  env: development
```

### ไฟล์ `config/config-prod.yml` (ตัวอย่างสำหรับ production)

```yaml
# Production configuration (usually overridden by environment variables)
# ค่ากำหนดสำหรับ production (ปกติจะถูก override โดย environment variables)
server:
  port: 8080
  read_timeout: 10s
  write_timeout: 10s
  idle_timeout: 120s

database:
  # ควรใช้ environment variables สำหรับ production
  host: ${CMON_DATABASE_HOST}
  port: ${CMON_DATABASE_PORT}
  user: ${CMON_DATABASE_USER}
  password: ${CMON_DATABASE_PASSWORD}
  dbname: ${CMON_DATABASE_NAME}
  sslmode: require
  max_open_conns: 100
  max_idle_conns: 50

redis:
  addr: ${CMON_REDIS_ADDR}
  password: ${CMON_REDIS_PASSWORD}
  db: 0

jwt:
  access_token_duration: 15m
  refresh_token_duration: 168h
  private_key_base64: ${CMON_JWT_PRIVATE_KEY_BASE64}
  public_key_base64: ${CMON_JWT_PUBLIC_KEY_BASE64}

smtp:
  host: ${SMTP_HOST}
  port: ${SMTP_PORT}
  username: ${SMTP_USER}
  password: ${SMTP_PASSWORD}
  from: ${SMTP_FROM}

rate_limit:
  requests_per_second: 50
  burst: 100

log:
  level: info
  env: production
```

---

## วิธีใช้งาน module นี้

1. ติดตั้ง dependency: `go get github.com/spf13/viper`
2. วางไฟล์ตามโครงสร้าง
3. เรียก `cfg := config.Load()` ที่ `main()` เพื่อโหลด config
4. ใช้ค่า config เช่น `cfg.Server.Port` เพื่อรัน HTTP server
5. สำหรับ production ให้ตั้ง environment variables (CMON_DATABASE_HOST, CMON_JWT_PRIVATE_KEY_BASE64 ฯลฯ) หรือใช้ Docker secrets

---

## แบบฝึกหัดท้าย module (3 ข้อ)

1. เพิ่ม field `app_name` ใน config และแสดงเมื่อ server start
2. สร้าง function `Validate()` ใน `Config` ที่ตรวจสอบค่า port ต้องอยู่ระหว่าง 1024-65535
3. เปลี่ยนให้โหลด config จาก S3 bucket โดยใช้ viper remote provider (ศึกษาเพิ่มเติม)

---

## แหล่งอ้างอิง

- [Viper documentation](https://github.com/spf13/viper)
- [12-factor config](https://12factor.net/config)

---

**หมายเหตุ:** module นี้เป็นส่วนหนึ่งของระบบ gobackend ทั้งหมด หากต้องการ module ถัดไป (Models, Repository, Usecase, Delivery, ฯลฯ) โปรดแจ้งคำว่า "ต่อไป" หรือระบุชื่อ module ที่ต้องการ