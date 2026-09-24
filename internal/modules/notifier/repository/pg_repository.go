package repository

import (
	"context"

	"icmongolang/internal/models"
	"gorm.io/gorm"
)

// ChannelConfig is the resolved configuration for a notification channel.
// Only the fields relevant to a real dispatch are exposed.
type ChannelConfig struct {
	Channel   string // email|sms|line|discord|io
	Enabled   bool
	// email
	Host     string
	Port     int
	Username string
	Password string
	From     string
	To       string
	// sms
	APIKey     string
	Originator string
	// line
	AccessToken string
	// discord
	WebhookURL string
	// io (mqtt topic)
	ControlTopic string
}

// Repository provides channel configuration and alarm log persistence.
type Repository interface {
	GetChannelConfig(ctx context.Context, channel string) (ChannelConfig, error)
	GetAllChannelConfigs(ctx context.Context) map[string]ChannelConfig
	SaveLog(ctx context.Context, log *models.SdAlarmProcessLogEmail) error
}

type pgRepository struct {
	db *gorm.DB
}

// NewPgRepository creates a DB-backed notifier repository.
func NewPgRepository(db *gorm.DB) Repository {
	return &pgRepository{db: db}
}

func (r *pgRepository) GetChannelConfig(ctx context.Context, channel string) (ChannelConfig, error) {
	cfg := ChannelConfig{Channel: channel}
	switch channel {
	case "email":
		var email models.SdIotEmail
		if err := r.db.WithContext(ctx).Where("status = 1").Order("createddate DESC").First(&email).Error; err != nil {
			return cfg, err
		}
		cfg.Enabled = true
		cfg.Host = email.Host
		if email.Port != nil {
			cfg.Port = *email.Port
		}
		cfg.Username = email.Username
		cfg.Password = email.Password
	case "sms":
		var sms models.SdIotSms
		if err := r.db.WithContext(ctx).Where("status = 1").Order("createddate DESC").First(&sms).Error; err != nil {
			return cfg, err
		}
		cfg.Enabled = true
		cfg.Host = sms.Host
		cfg.Username = sms.Username
		cfg.Password = sms.Password
		cfg.APIKey = sms.Apikey
		cfg.Originator = sms.Originator
	case "line":
		var line models.SdIotLine
		if err := r.db.WithContext(ctx).Where("status = 1").Order("createddate DESC").First(&line).Error; err != nil {
			return cfg, err
		}
		cfg.Enabled = true
		if line.Accesstoken != nil {
			cfg.AccessToken = *line.Accesstoken
		}
	case "nodered":
		var nr models.SdIotNodered
		if err := r.db.WithContext(ctx).Where("status = 1").Order("createddate DESC").First(&nr).Error; err != nil {
			return cfg, err
		}
		cfg.Enabled = true
		cfg.Host = nr.Host
		cfg.WebhookURL = nr.Host + ":" + nr.Port
	case "discord", "io":
		// Discord webhook URL and IO control topic are provided per-request
		// (or via env/config) since they are site-specific routing targets.
		cfg.Enabled = true
	default:
		cfg.Enabled = false
	}
	return cfg, nil
}

func (r *pgRepository) GetAllChannelConfigs(ctx context.Context) map[string]ChannelConfig {
	out := make(map[string]ChannelConfig)
	for _, ch := range []string{"email", "sms", "line", "nodered", "discord", "io"} {
		if c, err := r.GetChannelConfig(ctx, ch); err == nil {
			out[ch] = c
		}
	}
	return out
}

func (r *pgRepository) SaveLog(ctx context.Context, log *models.SdAlarmProcessLogEmail) error {
	if r.db == nil {
		return nil
	}
	return r.db.WithContext(ctx).Create(log).Error
}
