package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"icmongolang/internal/modules/settings"
	"icmongolang/internal/modules/settings/presenter"
	redisDb "icmongolang/pkg/db/redis"
	"icmongolang/pkg/helpers"
	"icmongolang/pkg/httpErrors"
	"icmongolang/pkg/logger"
	"icmongolang/pkg/mqtt"
)

const mqttDataCachePrefix = "get_device_data_ALL"

type settingsUseCase struct {
	repo   settings.SettingsRepositoryI
	mqtt   mqtt.Client
	cache  redisDb.Cache
	logger logger.Logger
}

func CreateSettingsUseCaseI(
	repo settings.SettingsRepositoryI,
	mqttClient mqtt.Client,
	cache redisDb.Cache,
	logger logger.Logger,
) settings.SettingsUseCaseI {
	return &settingsUseCase{repo: repo, mqtt: mqttClient, cache: cache, logger: logger}
}

// ---- Generic operations ----

func (uc *settingsUseCase) ListPaginate(ctx context.Context, spec settings.ListSpec, q settings.ListQuery) (*presenter.ListResult, error) {
	items, total, err := uc.repo.ListPaginate(ctx, spec, q)
	if err != nil {
		uc.logger.Errorf("settings ListPaginate table=%s: %v", spec.Table, err)
		return nil, err
	}
	return presenter.NewListResult(items, total, q.Page, q.PageSize, q.FilterMap()), nil
}

func (uc *settingsUseCase) RowsAll(ctx context.Context, spec settings.ListSpec, q settings.ListQuery) ([]map[string]interface{}, error) {
	items, err := uc.repo.RowsAll(ctx, spec, q)
	if err != nil {
		uc.logger.Errorf("settings RowsAll table=%s: %v", spec.Table, err)
		return nil, err
	}
	return items, nil
}

func (uc *settingsUseCase) CreateRow(ctx context.Context, table string, m map[string]interface{}) error {
	if err := uc.repo.CreateMap(ctx, table, m); err != nil {
		uc.logger.Errorf("settings CreateRow table=%s: %v", table, err)
		return err
	}
	return nil
}

func (uc *settingsUseCase) GetWhere(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	err := uc.repo.GetWhere(ctx, dest, query, args...)
	if err != nil && !errors.Is(err, context.Canceled) {
		uc.logger.Errorf("settings GetWhere query=%s: %v", query, err)
	}
	return err
}

func (uc *settingsUseCase) ExistsWhere(ctx context.Context, table string, query string, args ...interface{}) (bool, error) {
	return uc.repo.ExistsWhere(ctx, table, query, args...)
}

func (uc *settingsUseCase) CountWhere(ctx context.Context, table string, query string, args ...interface{}) (int64, error) {
	return uc.repo.CountWhere(ctx, table, query, args...)
}

func (uc *settingsUseCase) UpdateFields(ctx context.Context, table string, idCol string, idVal interface{}, fields map[string]interface{}) (int64, error) {
	n, err := uc.repo.UpdateFieldsMap(ctx, table, idCol, idVal, fields)
	if err != nil {
		uc.logger.Errorf("settings UpdateFields table=%s idCol=%s: %v", table, idCol, err)
	}
	return n, err
}

func (uc *settingsUseCase) SetStatusExclusive(ctx context.Context, table string, idCol string, id interface{}, status int, resetAll bool) (int64, error) {
	if status == 1 && resetAll {
		// NestJS first resets every row to status 0 (with updateddate),
		// then activates the requested one.
		if _, err := uc.repo.UpdateAllMap(ctx, table, map[string]interface{}{"status": "0", "updateddate": time.Now().In(helpers.GetTimeLocation())}); err != nil {
			uc.logger.Errorf("settings SetStatusExclusive reset-all table=%s: %v", table, err)
			return 0, err
		}
	}
	n, err := uc.repo.UpdateFieldsMap(ctx, table, idCol, id, map[string]interface{}{
		"status":      status,
		"updateddate": time.Now().In(helpers.GetTimeLocation()),
	})
	if err != nil {
		uc.logger.Errorf("settings SetStatusExclusive table=%s idCol=%s: %v", table, idCol, err)
	}
	return n, err
}

func (uc *settingsUseCase) DeleteRow(ctx context.Context, table string, query string, args ...interface{}) (int64, error) {
	n, err := uc.repo.DeleteWhere(ctx, table, query, args...)
	if err != nil {
		uc.logger.Errorf("settings DeleteRow table=%s query=%s: %v", table, query, err)
	}
	return n, err
}

// ---- Special flows ----

// MqttData ports GET /mqttdata: Redis cache -> MQTT request-and-wait fallback.
func (uc *settingsUseCase) MqttData(ctx context.Context, topic string) (*presenter.MqttDataResult, error) {
	topic = strings.TrimSpace(topic)
	if topic == "" {
		return nil, httpErrors.NewError(http.StatusUnprocessableEntity, "mqttdata is null.")
	}
	cacheKey := mqttDataCachePrefix + topic

	var cached []map[string]interface{}
	if uc.cache != nil {
		if err := uc.cache.Get(ctx, cacheKey, &cached); err == nil && cached != nil {
			return &presenter.MqttDataResult{GetDataFrom: "Cache", Payload: cached}, nil
		} else if err != nil && !errors.Is(err, context.DeadlineExceeded) {
			uc.logger.Warnf("settings MqttData cache get %s: %v", cacheKey, err)
		}
	}

	if uc.mqtt == nil {
		return nil, httpErrors.NewError(http.StatusServiceUnavailable, "mqtt client is not available")
	}
	raw, err := uc.mqtt.GetDataFromTopic(ctx, topic, 10*time.Second)
	if err != nil {
		uc.logger.Errorf("settings MqttData GetDataFromTopic topic=%s: %v", topic, err)
		return nil, err
	}
	payload := []map[string]interface{}{}
	if len(raw) > 0 {
		if jerr := json.Unmarshal(raw, &payload); jerr != nil {
			// non-JSON payloads are surfaced as a single raw record
			payload = []map[string]interface{}{{"raw": string(raw)}}
		}
	}
	if uc.cache != nil {
		if serr := uc.cache.Set(ctx, cacheKey, payload, 30*time.Second); serr != nil {
			uc.logger.Warnf("settings MqttData cache set %s: %v", cacheKey, serr)
		}
	}
	return &presenter.MqttDataResult{GetDataFrom: "MQTT", Payload: payload}, nil
}

// SendEmailStub ports GET /sendemail: the NestJS service never actually sends
// (the nodemailer call is commented out), so we return the same success shape.
func (uc *settingsUseCase) SendEmailStub(to, subject, content string) *presenter.SendEmailResult {
	if strings.TrimSpace(to) == "" {
		to = "cmoniots@gmail.com"
	}
	if strings.TrimSpace(subject) == "" {
		subject = "CmonIoT test send email"
	}
	if strings.TrimSpace(content) == "" {
		content = "test send email"
	}
	return &presenter.SendEmailResult{
		Success:   true,
		Code:      http.StatusOK,
		To:        to,
		Subject:   subject,
		Content:   content,
		Message:   "finally send email",
		MessageTh: "finally ส่ง email",
	}
}

// TestGmailConnection ports testGmailConnection: verifies smtp.gmail.com on
// ports 465 (implicit TLS) and 587 (STARTTLS). Credentials come from env
// GMAIL_USERNAME/GMAIL_APP_PASSWORD instead of hardcoded secrets.
func (uc *settingsUseCase) TestGmailConnection(ctx context.Context) *presenter.SendEmailResult {
	user := os.Getenv("GMAIL_USERNAME")
	pass := os.Getenv("GMAIL_APP_PASSWORD")
	to := os.Getenv("GMAIL_TEST_TO")
	if strings.TrimSpace(to) == "" {
		to = "icmon0955@gmail.com"
	}
	subject := "Alarm Test"
	content := "Alarm Test"

	fail := func(code int, msg string) *presenter.SendEmailResult {
		return &presenter.SendEmailResult{
			Success: false, Code: code, To: to, Subject: subject, Content: content,
			Error: msg, Message: msg, MessageTh: msg,
		}
	}
	if user == "" || pass == "" {
		return fail(http.StatusOK, "GMAIL_USERNAME/GMAIL_APP_PASSWORD not configured")
	}

	for _, port := range []int{465, 587} {
		if err := smtpVerifyAndSend(user, pass, to, subject, content, port); err != nil {
			uc.logger.Warnf("TestGmailConnection port %d failed: %v", port, err)
			continue
		}
		return &presenter.SendEmailResult{
			Success: true, Code: http.StatusOK, To: to, Subject: subject, Content: content,
			Message:   fmt.Sprintf("Gmail connection successful on port %d", port),
			MessageTh: fmt.Sprintf("เชื่อมต่อ Gmail สำเร็จบนพอร์ต %d", port),
		}
	}
	return fail(http.StatusOK, "Gmail connection failed on ports 465 and 587")
}
