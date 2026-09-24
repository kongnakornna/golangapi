package config

import (
	"errors"
	"log"
	"net"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
)

var (
	cfg *Config
)

type Config struct {
	Server         ServerConfig
	Postgres       PostgresConfig
	Redis          RedisConfig
	Jwt            JwtConfig
	FirstSuperUser FirstSuperUserConfig
	Logger         Logger
	SmtpEmail      SmtpEmailConfig
	Email          EmailConfig
	TaskRedis      TaskRedisConfig
	MQTT           MQTTConfig          `mapstructure:"mqtt"`
	InfluxDB       InfluxDBConfig      `mapstructure:"influxdb"`
	Kafka          KafkaConfig         `mapstructure:"kafka"` // ✅ เพิ่ม
	Elasticsearch  ElasticsearchConfig `mapstructure:"elasticsearch"`
	LLM            LLMConfig           `mapstructure:"llm"`
	VectorDB       VectorDBConfig      `mapstructure:"vectorDb"`
	Upload         UploadConfig        `mapstructure:"upload"`
}

type ServerConfig struct {
	AppVersion     string
	Port           string
	Mode           string
	ProcessTimeout int
	ReadTimeout    int
	WriteTimeout   int
	MigrateOnStart bool
	Timezone       string `mapstructure:"timezone"`
	BaseUrl        string
}

type Logger struct {
	Encoding string
	Level    string
}

type PostgresConfig struct {
	Host              string
	Port              string
	User              string
	Password          string
	Dbname            string
	ConnectionTimeout int
	MaxOpenConns      int
	MaxIdleConns      int
	ConnMaxLifetimeMin int
}

type RedisConfig struct {
	Addr         string
	Password     string
	Db           int
	MinIdleConns int
	PoolSize     int
	PoolTimeout  int
}

type TaskRedisConfig struct {
	Addr        string
	Db          int
	PoolTimeout int
}

type JwtConfig struct {
	SecretKey                  string
	DownloadTokenSecret        string
	Issuer                     string
	AccessTokenExpireDuration  int64
	AccessTokenPrivateKey      string
	AccessTokenPublicKey       string
	RefreshTokenExpireDuration int64
	RefreshTokenPrivateKey     string
	RefreshTokenPublicKey      string
}

type FirstSuperUserConfig struct {
	Email    string
	Name     string
	Password string
}

type LoggerConfig struct {
	Encoding string
	Level    string
}

type EmailConfig struct {
	From                string
	Name                string
	Link                string
	LogoLink            string
	Copyright           string
	VerificationSubject string
	ResetSubject        string
}

type SmtpEmailConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	UseTls   bool
	UseSsl   bool
	Timeout  int
}

type MQTTConfig struct {
	// ICMON_MQTT_BROKER is used instead of the generic MQTT_BROKER so an
	// ambient shell/env MQTT_BROKER (e.g. pointing at a docker-published
	// port) can't override the app's own broker config on native runs.
	Broker string `mapstructure:"broker" env:"ICMON_MQTT_BROKER"`
	// Port overrides the port embedded in Broker when set (yml mqtt.port or
	// ICMON_MQTT_PORT), so the MQTT port can be configured via .env without
	// editing the broker URL.
	Port              string `mapstructure:"port" env:"ICMON_MQTT_PORT"`
	ClientID          string `mapstructure:"client_id"`
	Username          string `mapstructure:"username"`
	Password          string `mapstructure:"password"`
	QOS               byte   `mapstructure:"qos"`
	ConnectionTimeout int    `mapstructure:"connection_timeout"`
}

type InfluxDBConfig struct {
	URL      string `mapstructure:"url"`
	Token    string `mapstructure:"token"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Org      string `mapstructure:"org"`
	Bucket   string `mapstructure:"bucket"`
	Timeout  int    `mapstructure:"timeout"`
}

// ✅ KafkaConfig – ใหม่
type KafkaConfig struct {
	Brokers []string `mapstructure:"brokers"`
	Topic   string   `mapstructure:"topic"`
	GroupID string   `mapstructure:"group_id"`
	Timeout int      `mapstructure:"timeout"`
}

// ElasticsearchConfig – Elasticsearch / Vector Database
type ElasticsearchConfig struct {
	Addresses []string `mapstructure:"addresses"`
	Username  string   `mapstructure:"username"`
	Password  string   `mapstructure:"password"`
	Index     string   `mapstructure:"index"`
	Timeout   int      `mapstructure:"timeout"`
}

// LLMConfig – OpenAI-compatible embedding provider
type LLMConfig struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
	Model   string `mapstructure:"model"`
	Timeout int    `mapstructure:"timeout"`
	Dims    int    `mapstructure:"dims"`
}

// VectorDBConfig – abstract vector database backend
// Provider: "elasticsearch" (dense_vector) หรือ "pgvector" (PostgreSQL + pgvector)
type VectorDBConfig struct {
	Provider string `mapstructure:"provider"` // elasticsearch | pgvector
	Index    string `mapstructure:"index"`    // ES index name / PG table name
	Dims     int    `mapstructure:"dims"`     // vector dimension (fallback 768)
}

type UploadConfig struct {
	Dir            string `mapstructure:"dir"`              // root dir ของไฟล์อัปโหลด
	MaxAvatarBytes int64  `mapstructure:"max_avatar_bytes"` // default 2MB
}

func ToSnakeCase(str string) string {
	snake := regexp.MustCompile("(.)([A-Z][a-z]+)").ReplaceAllString(str, "${1}_${2}")
	snake = regexp.MustCompile("([a-z0-9])([A-Z])").ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func BindEnvs(vp *viper.Viper, iface interface{}, partsKey []string, partsEnvKey []string) {
	ifv := reflect.ValueOf(iface)
	ift := reflect.TypeOf(iface)
	for i := 0; i < ift.NumField(); i++ {
		v := ifv.Field(i)
		t := ift.Field(i)

		seg := t.Name
		if tag := t.Tag.Get("mapstructure"); tag != "" {
			seg = tag
		}

		tv := strings.ToUpper(ToSnakeCase(seg))

		switch v.Kind() {
		case reflect.Struct:
			BindEnvs(vp, v.Interface(), append(partsKey, seg), append(partsEnvKey, tv))
		default:
			key := strings.ToLower(strings.Join(append(partsKey, seg), "."))
			envKey := strings.ToUpper(strings.Join(append(partsEnvKey, tv), "_"))
			// An explicit env tag (e.g. ICMON_MQTT_BROKER) is a full override
			// and must not be joined with the parent section's env prefix.
			if envTag := t.Tag.Get("env"); envTag != "" {
				envKey = strings.ToUpper(envTag)
			}

			vp.BindEnv(key, envKey) //nolint:errcheck
		}
	}
}

// Load config file from given path
func LoadConfig() (*viper.Viper, error) {
	// Load .env into the process environment (only sets vars that are not
	// already present, so Docker's env_file values always win). Ignore a
	// missing file so non-.env deployments keep working.
	_ = gotenv.Load(".env")

	v := viper.New()

	v.AddConfigPath(".")
	v.SetConfigName("config/config.default")
	v.SetConfigType("yml")

	BindEnvs(v, Config{}, []string{}, []string{})

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	var c Config
	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}
	return v, nil
}

// applyMQTTPort overrides the port embedded in the broker URL when an explicit
// port is configured (yml mqtt.port or ICMON_MQTT_PORT).
//
// Precedence:
//   - ICMON_MQTT_PORT (env) always wins over any port in the broker URL.
//   - A port already present in the broker URL (e.g. ICMON_MQTT_BROKER with its
//     own port, as docker-compose sets) wins over the yml mqtt.port default.
//   - Otherwise the configured port (yml mqtt.port) is appended to the broker.
func applyMQTTPort(m *MQTTConfig) {
	if m.Port == "" || m.Broker == "" {
		return
	}
	u, err := url.Parse(m.Broker)
	if err != nil || u.Host == "" {
		return
	}
	if u.Port() != "" && os.Getenv("ICMON_MQTT_PORT") == "" {
		return
	}
	u.Host = net.JoinHostPort(u.Hostname(), m.Port)
	m.Broker = u.String()
}

// Parse config file
func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	applyMQTTPort(&c.MQTT)

	cfg = &c

	return &c, nil
}

func GetCfg() *Config {
	if cfg == nil {
		cfg = new(Config)
	}
	return cfg
}
