package config

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// ข้อผิดพลาดการตรวจสอบคอนฟิกูเรชัน
var (
	ErrInvalidPort         = errors.New("หมายเลขพอร์ตเซิร์ฟเวอร์ไม่ถูกต้อง")
	ErrMissingDatabaseHost = errors.New("ไม่พบโฮสต์ของฐานข้อมูล")
)

// AppConfig โครงสร้างคอนฟิกูเรชันระดับบนสุด ให้ตรงกับคีย์ app ในไฟล์ yaml
type AppConfig struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `structure:"redis"`
	Log      LogConfig      `mapstructure:"log"`
	JWT      JWTConfig      `mapstructure:"jwt"`
}

// Config โครงสร้างคอนฟิกูเรชันของแอปพลิเคชัน
type Config struct {
	App AppConfig `mapstructure:"app"`
}

// ServerConfig คอนฟิกูเรชันเซิร์ฟเวอร์
type ServerConfig struct {
	Port         int           `mapstructure:"port" env:"SERVER_PORT"`
	Timeout      time.Duration `mapstructure:"timeout" env:"SERVER_TIMEOUT"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout" env:"SERVER_READ_TIMEOUT"`
	WriteTimeout time.Duration `mapstructure:"write_timeout" env:"SERVER_WRITE_TIMEOUT"`
}

// DatabaseConfig คอนฟิกูเรชันฐานข้อมูล
type DatabaseConfig struct {
	Driver          string        `mapstructure:"driver" env:"DB_DRIVER"`
	Host            string        `mapstructure:"host" env:"DB_HOST"`
	Port            int           `mapstructure:"port" env:"DB_PORT"`
	Username        string        `mapstructure:"username" env:"DB_USERNAME"`
	Password        string        `mapstructure:"password" env:"DB_PASSWORD"`
	DBName          string        `mapstructure:"dbname" env:"DB_NAME"`
	SSLMode         string        `mapstructure:"sslmode" env:"DB_SSLMODE"`
	MaxOpenConns    int           `mapstructure:"max_open_conns" env:"DB_MAX_OPEN_CONNS"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns" env:"DB_MAX_IDLE_CONNS"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime" env:"DB_CONN_MAX_LIFETIME"`
}

// RedisConfig คอนฟิกูเรชัน Redis
type RedisConfig struct {
	Enabled  bool   `mapstructure:"enabled" env:"REDIS_ENABLED"`
	Host     string `mapstructure:"host" env:"REDIS_HOST"`
	Port     int    `mapstructure:"port" env:"REDIS_PORT"`
	Password string `mapstructure:"password" env:"REDIS_PASSWORD"`
	DB       int    `mapstructure:"db" env:"REDIS_DB"`
}

// LogConfig คอนฟิกูเรชันบันทึก日志
type LogConfig struct {
	Level   string `mapstructure:"level" env:"LOG_LEVEL"`
	File    string `mapstructure:"file" env:"LOG_FILE"`
	Console bool   `mapstructure:"console" env:"LOG_CONSOLE"`
}

// JWTConfig คอนฟิกูเรชัน JWT
type JWTConfig struct {
	Secret          string        `mapstructure:"secret" env:"JWT_SECRET"`
	AccessTokenExp  time.Duration `mapstructure:"access_token_exp" env:"JWT_ACCESS_TOKEN_EXP"`
	RefreshTokenExp time.Duration `mapstructure:"refresh_token_exp" env:"JWT_REFRESH_TOKEN_EXP"`
	Issuer          string        `mapstructure:"issuer" env:"JWT_ISSUER"`
}

// LoadConfig โหลดคอนฟิกูเรชัน
func LoadConfig(path string) (*AppConfig, error) {
	// เริ่มต้น viper
	viper.SetConfigFile(path)

	// ตั้งค่าตัวแปรสภาพแวดล้อมและตัวคั่น
	viper.SetEnvPrefix("APP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// ผูกตัวแปรสภาพแวดล้อม
	bindEnvVariables()

	// ตั้งค่าค่าเริ่มต้น (เพื่อหลีกเลี่ยงการถูกแทนที่ด้วยค่าเริ่มต้นของ bool)
	viper.SetDefault("app.redis.enabled", true)

	// เปิดใช้งานการสนับสนุนตัวแปรสภาพแวดล้อม
	viper.AutomaticEnv()

	// อ่านไฟล์คอนฟิกูเรชันก่อน
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("ไม่สามารถอ่านไฟล์คอนฟิกูเรชันได้: %w", err)
	}

	// แยกวิเคราะห์คอนฟิกูเรชันไปยังโครงสร้าง
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("ไม่สามารถแยกวิเคราะห์คอนฟิกูเรชันได้: %w", err)
	}

	// ตั้งค่าค่าเริ่มต้น
	setDefaults(&config.App)

	return &config.App, nil
}

// ผูกตัวแปรสภาพแวดล้อม
func bindEnvVariables() {
	// ตัวแปรสภาพแวดล้อมคอนฟิกูเรชันเซิร์ฟเวอร์
	viper.BindEnv("app.server.port", "APP_SERVER_PORT")
	viper.BindEnv("app.server.timeout", "APP_SERVER_TIMEOUT")
	viper.BindEnv("app.server.read_timeout", "APP_SERVER_READ_TIMEOUT")
	viper.BindEnv("app.server.write_timeout", "APP_SERVER_WRITE_TIMEOUT")

	// ตัวแปรสภาพแวดล้อมคอนฟิกูเรชันฐานข้อมูล
	viper.BindEnv("app.database.driver", "APP_DB_DRIVER")
	viper.BindEnv("app.database.host", "APP_DB_HOST")
	viper.BindEnv("app.database.port", "APP_DB_PORT")
	viper.BindEnv("app.database.username", "APP_DB_USERNAME")
	viper.BindEnv("app.database.password", "APP_DB_PASSWORD")
	viper.BindEnv("app.database.dbname", "APP_DB_NAME")
	viper.BindEnv("app.database.sslmode", "APP_DB_SSLMODE")
	viper.BindEnv("app.database.max_open_conns", "APP_DB_MAX_OPEN_CONNS")
	viper.BindEnv("app.database.max_idle_conns", "APP_DB_MAX_IDLE_CONNS")
	viper.BindEnv("app.database.conn_max_lifetime", "APP_DB_CONN_MAX_LIFETIME")

	// ตัวแปรสภาพแวดล้อมคอนฟิกูเรชัน Redis
	viper.BindEnv("app.redis.enabled", "APP_REDIS_ENABLED")
	viper.BindEnv("app.redis.host", "APP_REDIS_HOST")
	viper.BindEnv("app.redis.port", "APP_REDIS_PORT")
	viper.BindEnv("app.redis.password", "APP_REDIS_PASSWORD")
	viper.BindEnv("app.redis.db", "APP_REDIS_DB")

	// ตัวแปรสภาพแวดล้อมคอนฟิกูเรชันบันทึก日志
	viper.BindEnv("app.log.level", "APP_LOG_LEVEL")
	viper.BindEnv("app.log.file", "APP_LOG_FILE")
	viper.BindEnv("app.log.console", "APP_LOG_CONSOLE")

	// ตัวแปรสภาพแวดล้อมคอนฟิกูเรชัน JWT
	viper.BindEnv("app.jwt.secret", "APP_JWT_SECRET")
	viper.BindEnv("app.jwt.access_token_exp", "APP_JWT_ACCESS_TOKEN_EXP")
	viper.BindEnv("app.jwt.refresh_token_exp", "APP_JWT_REFRESH_TOKEN_EXP")
	viper.BindEnv("app.jwt.issuer", "APP_JWT_ISSUER")
}

// ตั้งค่าค่าเริ่มต้น
func setDefaults(config *AppConfig) {
	// ค่าเริ่มต้นเซิร์ฟเวอร์
	if config.Server.Port == 0 {
		config.Server.Port = 7001
	}
	if config.Server.Timeout == 0 {
		config.Server.Timeout = 30 * time.Second
	}
	if config.Server.ReadTimeout == 0 {
		config.Server.ReadTimeout = 15 * time.Second
	}
	if config.Server.WriteTimeout == 0 {
		config.Server.WriteTimeout = 15 * time.Second
	}

	// ค่าเริ่มต้นพูลการเชื่อมต่อฐานข้อมูล
	if config.Database.MaxOpenConns == 0 {
		config.Database.MaxOpenConns = 20
	}
	if config.Database.MaxIdleConns == 0 {
		config.Database.MaxIdleConns = 5
	}
	if config.Database.ConnMaxLifetime == 0 {
		config.Database.ConnMaxLifetime = 1 * time.Hour
	}

	// ค่าเริ่มต้น JWT
	if config.JWT.AccessTokenExp == 0 {
		config.JWT.AccessTokenExp = 24 * time.Hour
	}
	if config.JWT.RefreshTokenExp == 0 {
		config.JWT.RefreshTokenExp = 7 * 24 * time.Hour
	}
	if config.JWT.Issuer == "" {
		config.JWT.Issuer = "go-rest-starter"
	}
}

// GetDSN รับสตริงการเชื่อมต่อฐานข้อมูล
func (c *DatabaseConfig) GetDSN() string {
	// สร้าง DSN สำหรับ PostgreSQL - ตรวจสอบให้แน่ใจว่าพารามิเตอร์ dbname ถูกต้อง
	if c.Password == "" {
		return fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=%s",
			c.Host, c.Port, c.Username, c.DBName, c.SSLMode)
	} else {
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			c.Host, c.Port, c.Username, c.Password, c.DBName, c.SSLMode)
	}
}
