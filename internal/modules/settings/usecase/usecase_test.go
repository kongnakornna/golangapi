package usecase

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"icmongolang/internal/modules/settings"
	redisDb "icmongolang/pkg/db/redis"
	"icmongolang/pkg/httpErrors"
)

// ---- fakes ----

type fakeRepo struct {
	settings.SettingsRepositoryI

	listItems []map[string]interface{}
	listTotal int64
	listErr   error
}

func (f *fakeRepo) ListPaginate(ctx context.Context, spec settings.ListSpec, q settings.ListQuery) ([]map[string]interface{}, int64, error) {
	return f.listItems, f.listTotal, f.listErr
}

// cacheHit returns a cached payload; cacheMiss always misses.
type cacheHit struct{ key string }

func (c *cacheHit) Get(ctx context.Context, key string, dst interface{}) error {
	p := dst.(*[]map[string]interface{})
	*p = []map[string]interface{}{{"temperature": 25.5}}
	return nil
}

func (c *cacheHit) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	c.key = key
	return nil
}

type cacheMiss struct{}

func (cacheMiss) Get(ctx context.Context, key string, dst interface{}) error {
	return errors.New("miss")
}
func (cacheMiss) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return nil
}

func newUC(repo settings.SettingsRepositoryI, cache redisDb.Cache) settings.SettingsUseCaseI {
	return CreateSettingsUseCaseI(repo, nil, cache, testLogger{})
}

type testLogger struct{}

func (testLogger) InitLogger()                        {}
func (testLogger) Debug(args ...interface{})          {}
func (testLogger) Debugf(t string, a ...interface{})  {}
func (testLogger) Info(args ...interface{})           {}
func (testLogger) Infof(t string, a ...interface{})   {}
func (testLogger) Warn(args ...interface{})           {}
func (testLogger) Warnf(t string, a ...interface{})   {}
func (testLogger) Error(args ...interface{})          {}
func (testLogger) Errorf(t string, a ...interface{})  {}
func (testLogger) DPanic(args ...interface{})         {}
func (testLogger) DPanicf(t string, a ...interface{}) {}
func (testLogger) Fatal(args ...interface{})          {}
func (testLogger) Fatalf(t string, a ...interface{})  {}
func (testLogger) Sync() error                        { return nil }

// ---- tests ----

func TestListPaginateWrapsMeta(t *testing.T) {
	repo := &fakeRepo{listItems: []map[string]interface{}{{"setting_id": int64(1)}}, listTotal: 250}
	uc := newUC(repo, cacheMiss{})

	res, err := uc.ListPaginate(context.Background(), settings.SettingSpec,
		settings.NewListQuery(2, 100, "", nil))
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if res.Total != 250 || res.Page != 2 || res.PageSize != 100 || res.TotalPages != 3 {
		t.Fatalf("bad meta: %+v", res)
	}
	if res.CurrentPage != 2 || res.Filter == nil {
		t.Fatalf("bad nestjs-compatible meta: %+v", res)
	}
	if len(res.Data) != 1 {
		t.Fatalf("expected 1 item")
	}
}

func TestListPaginateEmptyBecomesEmptySlice(t *testing.T) {
	uc := newUC(&fakeRepo{}, cacheMiss{})
	res, err := uc.ListPaginate(context.Background(), settings.SettingSpec, settings.NewListQuery(1, 10, "", nil))
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if res.Data == nil || len(res.Data) != 0 {
		t.Fatal("data must be non-nil empty slice")
	}
	if res.TotalPages < 1 {
		t.Fatal("totalPages must be at least 1 for empty lists")
	}
}

func TestMqttDataRejectsEmptyTopic(t *testing.T) {
	uc := newUC(&fakeRepo{}, cacheMiss{})
	_, err := uc.MqttData(context.Background(), "  ")
	resp := httpErrors.ParseErrors(err)
	if resp.GetStatus() != http.StatusUnprocessableEntity {
		t.Fatalf("want 422, got %d", resp.GetStatus())
	}
}

func TestMqttDataNilMqttClientAfterCacheMiss(t *testing.T) {
	uc := newUC(&fakeRepo{}, cacheMiss{})
	_, err := uc.MqttData(context.Background(), "device/1/data")
	resp := httpErrors.ParseErrors(err)
	if resp.GetStatus() != http.StatusServiceUnavailable {
		t.Fatalf("want 503 when mqtt client is nil, got %d", resp.GetStatus())
	}
}

func TestMqttDataCacheHit(t *testing.T) {
	uc := newUC(&fakeRepo{}, &cacheHit{})
	res, err := uc.MqttData(context.Background(), "device/1/data")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if res.GetDataFrom != "Cache" || len(res.Payload) != 1 {
		t.Fatalf("expected cache hit payload: %+v", res)
	}
}

func TestSendEmailStubDefaults(t *testing.T) {
	uc := newUC(&fakeRepo{}, cacheMiss{})
	res := uc.SendEmailStub("", "", "")
	if !res.Success || res.Code != 200 {
		t.Fatalf("stub must succeed: %+v", res)
	}
	if res.To != "cmoniots@gmail.com" || res.Subject == "" || res.Content == "" {
		t.Fatalf("defaults not applied: %+v", res)
	}
}

func TestTestGmailConnectionWithoutEnvFailsGracefully(t *testing.T) {
	t.Setenv("GMAIL_USERNAME", "")
	t.Setenv("GMAIL_APP_PASSWORD", "")
	uc := newUC(&fakeRepo{}, cacheMiss{})
	res := uc.TestGmailConnection(context.Background())
	if res.Success {
		t.Fatal("must fail without credentials configured")
	}
}

func TestSensorTypesLang(t *testing.T) {
	uc := newUC(&fakeRepo{}, cacheMiss{})
	en := uc.SensorTypes("")
	th := uc.SensorTypes("th")
	if len(en) == 0 || len(th) != len(en) {
		t.Fatal("sensor type lists must be parallel")
	}
	if en[0]["sensor_name"] == th[0]["sensor_name"] {
		t.Fatal("th list should differ from en list")
	}
}
