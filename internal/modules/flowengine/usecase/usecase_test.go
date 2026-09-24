package usecase_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"icmongolang/internal/models"
	"icmongolang/internal/modules/flowengine/presenter"
	"icmongolang/internal/modules/flowengine/repository"
	"icmongolang/internal/modules/flowengine/usecase"
	notifierPresenter "icmongolang/internal/modules/notifier/presenter"
	realtimePresenter "icmongolang/internal/modules/realtime/presenter"
	"icmongolang/pkg/logger"
)

type memRepo struct {
	flows map[string]*models.FlowDefinition
	next  int
}

func (m *memRepo) Create(ctx context.Context, f *models.FlowDefinition) error {
	m.next++
	f.ID = "flow-" + itoa(m.next)
	if m.flows == nil {
		m.flows = map[string]*models.FlowDefinition{}
	}
	m.flows[f.ID] = f
	return nil
}

func (m *memRepo) GetByID(ctx context.Context, id string) (*models.FlowDefinition, error) {
	if f, ok := m.flows[id]; ok {
		return f, nil
	}
	return nil, errors.New("not found")
}

func (m *memRepo) Update(ctx context.Context, f *models.FlowDefinition) error {
	m.flows[f.ID] = f
	return nil
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	b := []byte{}
	for n > 0 {
		b = append([]byte{byte('0' + n%10)}, b...)
		n /= 10
	}
	return string(b)
}

type stubNotifier struct {
	calls []string
}

func (s *stubNotifier) Dispatch(ctx context.Context, req *notifierPresenter.DispatchRequest) (*notifierPresenter.DispatchResponse, error) {
	s.calls = append(s.calls, req.Channel)
	return &notifierPresenter.DispatchResponse{Channel: req.Channel, OK: true}, nil
}

type stubRealtime struct {
	calls []string
}

func (s *stubRealtime) Publish(ctx context.Context, req *realtimePresenter.PublishRequest) (*realtimePresenter.PublishResponse, error) {
	s.calls = append(s.calls, req.Event)
	return &realtimePresenter.PublishResponse{OK: true, Event: req.Event}, nil
}

func newUC(r repository.Repository, n usecase.Notifier, rt usecase.Realtime) usecase.FlowEngineUseCase {
	return usecase.NewFlowEngineUseCase(r, n, rt, logger.NewLogger("flow-test"))
}

func validRequest() *presenter.CreateRequest {
	return &presenter.CreateRequest{
		Name: "test",
		Nodes: []presenter.FlowNode{
			{ID: "a", Type: "alarm", Name: "alarm"},
			{ID: "b", Type: "email", Name: "mail"},
			{ID: "c", Type: "ws", Name: "push"},
		},
		Edges: []presenter.FlowEdge{
			{ID: "e1", Source: "a", Target: "b"},
			{ID: "e2", Source: "a", Target: "c"},
		},
	}
}

func TestValidate_Invalid(t *testing.T) {
	uc := newUC(&memRepo{}, nil, nil)
	if resp, _ := uc.Validate(context.Background(), nil); resp.Valid {
		t.Fatal("expected invalid for nil")
	}
	if resp, _ := uc.Validate(context.Background(), &presenter.CreateRequest{Name: ""}); resp.Valid {
		t.Fatal("expected invalid for empty name")
	}
	// dangling edge
	req := validRequest()
	req.Edges = append(req.Edges, presenter.FlowEdge{ID: "x", Source: "zz", Target: "a"})
	if resp, _ := uc.Validate(context.Background(), req); resp.Valid {
		t.Fatal("expected invalid for dangling edge")
	}
	// self-loop
	req2 := validRequest()
	req2.Edges = []presenter.FlowEdge{{ID: "s", Source: "a", Target: "a"}}
	if resp, _ := uc.Validate(context.Background(), req2); resp.Valid {
		t.Fatal("expected invalid for self-loop")
	}
	// duplicate node id
	req3 := validRequest()
	req3.Nodes = append(req3.Nodes, presenter.FlowNode{ID: "a", Type: "email"})
	if resp, _ := uc.Validate(context.Background(), req3); resp.Valid {
		t.Fatal("expected invalid for duplicate node id")
	}
}

func TestValidate_Cycle(t *testing.T) {
	uc := newUC(&memRepo{}, nil, nil)
	req := &presenter.CreateRequest{
		Name: "cycle",
		Nodes: []presenter.FlowNode{
			{ID: "a", Type: "email"},
			{ID: "b", Type: "email"},
			{ID: "c", Type: "email"},
		},
		Edges: []presenter.FlowEdge{
			{ID: "1", Source: "a", Target: "b"},
			{ID: "2", Source: "b", Target: "c"},
			{ID: "3", Source: "c", Target: "a"},
		},
	}
	if resp, _ := uc.Validate(context.Background(), req); resp.Valid {
		t.Fatal("expected invalid for cycle")
	}
}

func TestValidate_Valid(t *testing.T) {
	uc := newUC(&memRepo{}, nil, nil)
	resp, err := uc.Validate(context.Background(), validRequest())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.Valid {
		t.Fatalf("expected valid, got %+v", resp)
	}
}

func TestCreate_Valid_Persists(t *testing.T) {
	repo := &memRepo{}
	uc := newUC(repo, &stubNotifier{}, &stubRealtime{})
	flow, err := uc.Create(context.Background(), validRequest())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if flow.ID == "" || flow.Name != "test" {
		t.Fatalf("unexpected flow: %+v", flow)
	}
	if len(repo.flows) != 1 {
		t.Fatalf("expected 1 persisted flow, got %d", len(repo.flows))
	}
}

func TestCreate_Invalid_Raises(t *testing.T) {
	uc := newUC(&memRepo{}, nil, nil)
	if _, err := uc.Create(context.Background(), &presenter.CreateRequest{Name: ""}); err == nil {
		t.Fatal("expected ErrInvalidFlow")
	}
}

func TestCreate_NilRepo_Graceful(t *testing.T) {
	uc := newUC(nil, nil, nil)
	flow, err := uc.Create(context.Background(), validRequest())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if flow.Name != "test" {
		t.Fatal("expected flow even without repo")
	}
}

func TestGet_NotFound(t *testing.T) {
	uc := newUC(&memRepo{}, nil, nil)
	if _, err := uc.Get(context.Background(), "missing"); err == nil {
		t.Fatal("expected not-found error")
	}
}

func TestTrigger_DispatchesToNotifierAndRealtime(t *testing.T) {
	repo := &memRepo{}
	notif := &stubNotifier{}
	rt := &stubRealtime{}
	uc := newUC(repo, notif, rt)
	flow, err := uc.Create(context.Background(), validRequest())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp, err := uc.Trigger(context.Background(), &presenter.TriggerRequest{
		FlowID:  flow.ID,
		Payload: map[string]any{"message": "hello"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.OK {
		t.Fatalf("expected ok, %s", resp.Message)
	}
	if len(notif.calls) != 1 || notif.calls[0] != "email" {
		t.Fatalf("expected 1 email dispatch, got %v", notif.calls)
	}
	if len(rt.calls) != 1 || rt.calls[0] != "monitor" {
		t.Fatalf("expected 1 realtime publish, got %v", rt.calls)
	}
}

func TestTrigger_NilDeps_SkipGracefully(t *testing.T) {
	repo := &memRepo{}
	uc := newUC(repo, nil, nil)
	flow, _ := uc.Create(context.Background(), validRequest())
	resp, err := uc.Trigger(context.Background(), &presenter.TriggerRequest{FlowID: flow.ID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.OK {
		t.Fatalf("expected ok with nil deps, %s", resp.Message)
	}
}

func TestTrigger_FlowNotFound(t *testing.T) {
	uc := newUC(&memRepo{}, nil, nil)
	if _, err := uc.Trigger(context.Background(), &presenter.TriggerRequest{FlowID: "nope"}); err == nil {
		t.Fatal("expected not-found error")
	}
}

var _ = json.Marshal
