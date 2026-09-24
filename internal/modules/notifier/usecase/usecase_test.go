package usecase_test

import (
	"context"
	"testing"

	"icmongolang/config"
	"icmongolang/internal/models"
	"icmongolang/internal/modules/notifier/presenter"
	"icmongolang/internal/modules/notifier/repository"
	"icmongolang/internal/modules/notifier/usecase"
	"icmongolang/pkg/logger"
)

type stubRepo struct {
	configs map[string]repository.ChannelConfig
}

func (s *stubRepo) GetChannelConfig(ctx context.Context, channel string) (repository.ChannelConfig, error) {
	return s.configs[channel], nil
}

func (s *stubRepo) GetAllChannelConfigs(ctx context.Context) map[string]repository.ChannelConfig {
	return s.configs
}

func (s *stubRepo) SaveLog(ctx context.Context, _ *models.SdAlarmProcessLogEmail) error {
	return nil
}

func newUC(repo repository.Repository) usecase.NotifierUseCase {
	return usecase.NewNotifierUseCase(repo, &config.Config{}, nil, logger.NewLogger("notifier-test"))
}

func TestDispatch_InvalidRequest(t *testing.T) {
	uc := newUC(&stubRepo{})
	if _, err := uc.Dispatch(context.Background(), nil); err == nil {
		t.Fatal("expected error for nil request")
	}
	if _, err := uc.Dispatch(context.Background(), &presenter.DispatchRequest{}); err == nil {
		t.Fatal("expected error for empty channel")
	}
}

func TestDispatch_UnknownChannel(t *testing.T) {
	uc := newUC(&stubRepo{})
	_, err := uc.Dispatch(context.Background(), &presenter.DispatchRequest{Channel: "fax"})
	if err == nil {
		t.Fatal("expected unknown-channel error")
	}
}

func TestDispatch_EmailNotConfigured_ReturnsGracefully(t *testing.T) {
	uc := newUC(&stubRepo{configs: map[string]repository.ChannelConfig{
		"email": {Channel: "email", Enabled: true},
	}})
	resp, err := uc.Dispatch(context.Background(), &presenter.DispatchRequest{
		Channel: "email",
		Subject: "test",
		Content: "body",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Empty SMTP config (no from/to) -> sender returns (false, nil): graceful,
	// never a hard failure in existing flows.
	if resp.OK {
		t.Fatalf("expected not-sent response")
	}
}

func TestHealth(t *testing.T) {
	uc := newUC(&stubRepo{})
	h := uc.Health(context.Background())
	if !h.Enabled {
		t.Fatalf("expected enabled=true")
	}
}
