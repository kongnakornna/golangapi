package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"icmongolang/internal/models"
	"icmongolang/internal/modules/flowengine/presenter"
	"icmongolang/internal/modules/flowengine/repository"
	notifierPresenter "icmongolang/internal/modules/notifier/presenter"
	realtimePresenter "icmongolang/internal/modules/realtime/presenter"
	"icmongolang/pkg/logger"
)

// ErrInvalidFlow is returned when a flow graph fails validation.
var ErrInvalidFlow = errors.New("flowengine: invalid flow")

// ErrFlowNotFound is returned when a flow id does not exist.
var ErrFlowNotFound = errors.New("flowengine: flow not found")

// Notifier is the subset of the notifier use case the runtime dispatches to.
type Notifier interface {
	Dispatch(ctx context.Context, req *notifierPresenter.DispatchRequest) (*notifierPresenter.DispatchResponse, error)
}

// Realtime is the subset of the realtime use case the runtime pushes to.
type Realtime interface {
	Publish(ctx context.Context, req *realtimePresenter.PublishRequest) (*realtimePresenter.PublishResponse, error)
}

// FlowEngineUseCase defines the workflow graph + runtime contract.
type FlowEngineUseCase interface {
	Create(ctx context.Context, req *presenter.CreateRequest) (*presenter.Flow, error)
	Get(ctx context.Context, id string) (*presenter.Flow, error)
	Validate(ctx context.Context, req *presenter.CreateRequest) (*presenter.ValidateResponse, error)
	Trigger(ctx context.Context, req *presenter.TriggerRequest) (*presenter.TriggerResponse, error)
}

type flowEngineUseCase struct {
	repo     repository.Repository
	notifier Notifier
	realtime Realtime
	logger   logger.Logger
}

// NewFlowEngineUseCase builds the flow engine use case. repo/notifier/realtime
// may be nil — trigger then runs the graph but skips unwired dispatchers.
func NewFlowEngineUseCase(repo repository.Repository, notifier Notifier, realtime Realtime, log logger.Logger) FlowEngineUseCase {
	return &flowEngineUseCase{repo: repo, notifier: notifier, realtime: realtime, logger: log}
}

func (u *flowEngineUseCase) Create(ctx context.Context, req *presenter.CreateRequest) (*presenter.Flow, error) {
	if req == nil {
		return nil, errors.New("flowengine: nil create request")
	}
	if resp := u.validateGraph(req); !resp.Valid {
		return nil, ErrInvalidFlow
	}
	nodes, err := json.Marshal(req.Nodes)
	if err != nil {
		return nil, err
	}
	edges, err := json.Marshal(req.Edges)
	if err != nil {
		return nil, err
	}
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	rec := &models.FlowDefinition{Name: req.Name, Nodes: nodes, Edges: edges, Enabled: enabled}
	if u.repo != nil {
		if err := u.repo.Create(ctx, rec); err != nil {
			u.logger.Errorf("flowengine: create flow error: %v", err)
			return nil, err
		}
	}
	u.logger.Infof("flowengine: created flow id=%s name=%s", rec.ID, req.Name)
	return &presenter.Flow{ID: rec.ID, Name: rec.Name, Nodes: req.Nodes, Edges: req.Edges, Enabled: enabled}, nil
}

func (u *flowEngineUseCase) Get(ctx context.Context, id string) (*presenter.Flow, error) {
	if u.repo == nil {
		return nil, ErrFlowNotFound
	}
	rec, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrFlowNotFound
	}
	return toPresenterFlow(rec)
}

func (u *flowEngineUseCase) Validate(ctx context.Context, req *presenter.CreateRequest) (*presenter.ValidateResponse, error) {
	resp := u.validateGraph(req)
	return resp, nil
}

func (u *flowEngineUseCase) Trigger(ctx context.Context, req *presenter.TriggerRequest) (*presenter.TriggerResponse, error) {
	if req == nil || req.FlowID == "" {
		return nil, errors.New("flowengine: nil trigger request")
	}
	if u.repo == nil {
		return &presenter.TriggerResponse{OK: false, Message: "flow store not wired"}, nil
	}
	rec, err := u.repo.GetByID(ctx, req.FlowID)
	if err != nil {
		return nil, ErrFlowNotFound
	}
	flow, err := toPresenterFlow(rec)
	if err != nil {
		return nil, err
	}
	order, err := topologicalOrder(flow)
	if err != nil {
		return nil, ErrInvalidFlow
	}
	var executed []string
	for _, node := range order {
		if u.dispatchNode(ctx, node, req.Payload) {
			executed = append(executed, node.ID)
		}
	}
	u.logger.Infof("flowengine: triggered flow=%s nodes=%d", req.FlowID, len(executed))
	return &presenter.TriggerResponse{OK: true, Executed: executed, Message: "ok"}, nil
}

// dispatchNode routes an action node to the wired dispatcher (notifier/realtime).
// Returns true when the node was a dispatcher type that ran (or was skipped
// gracefully due to a nil dependency).
func (u *flowEngineUseCase) dispatchNode(ctx context.Context, node presenter.FlowNode, payload map[string]any) bool {
	switch strings.ToLower(node.Type) {
	case "email", "sms", "line", "discord", "io", "notify":
		if u.notifier == nil {
			u.logger.Warnf("flowengine: notifier nil, skipping node=%s", node.ID)
			return true
		}
		_, _ = u.notifier.Dispatch(ctx, flowNotifyRequest(node, payload))
		return true
	case "ws", "realtime", "monitor":
		if u.realtime == nil {
			u.logger.Warnf("flowengine: realtime nil, skipping node=%s", node.ID)
			return true
		}
		_, _ = u.realtime.Publish(ctx, &realtimePresenter.PublishRequest{
			Event: confString(node.Conf, "event", "monitor"),
			Data:  payload,
		})
		return true
	default:
		// trigger/source nodes and unknown nodes do not dispatch.
		return false
	}
}

func flowNotifyRequest(node presenter.FlowNode, payload map[string]any) *notifierPresenter.DispatchRequest {
	channel := confString(node.Conf, "channel", node.Type)
	req := &notifierPresenter.DispatchRequest{
		Channel:      channel,
		Title:        confString(node.Conf, "title", "Flow notification"),
		Subject:      confString(node.Conf, "subject", node.Name),
		Content:      confString(node.Conf, "content", messageFromPayload(payload)),
		ControlTopic: confString(node.Conf, "control_topic", ""),
		WebhookURL:   confString(node.Conf, "webhook_url", ""),
	}
	if v, ok := node.Conf["mqtt_control_on"].(string); ok {
		req.MqttControlOn = v
	}
	if v, ok := node.Conf["mqtt_control_off"].(string); ok {
		req.MqttControlOff = v
	}
	return req
}

func confString(m map[string]any, key, fallback string) string {
	if s, ok := m[key].(string); ok && s != "" {
		return s
	}
	return fallback
}

func messageFromPayload(payload map[string]any) string {
	if payload == nil {
		return ""
	}
	if s, ok := payload["message"].(string); ok && s != "" {
		return s
	}
	if s, ok := payload["content"].(string); ok {
		return s
	}
	return ""
}

// validateGraph checks node/edge consistency and that the graph is acyclic.
func (u *flowEngineUseCase) validateGraph(req *presenter.CreateRequest) *presenter.ValidateResponse {
	resp := &presenter.ValidateResponse{Valid: true}
	if req == nil {
		resp.Valid = false
		resp.Errors = append(resp.Errors, "nil request")
		return resp
	}
	if strings.TrimSpace(req.Name) == "" {
		resp.Valid = false
		resp.Errors = append(resp.Errors, "name is required")
	}
	if len(req.Nodes) == 0 {
		resp.Valid = false
		resp.Errors = append(resp.Errors, "at least one node is required")
	}
	nodeIDs := map[string]bool{}
	for _, n := range req.Nodes {
		if n.ID == "" {
			resp.Valid = false
			resp.Errors = append(resp.Errors, "node id is required")
			continue
		}
		if nodeIDs[n.ID] {
			resp.Valid = false
			resp.Errors = append(resp.Errors, "duplicate node id: "+n.ID)
		}
		nodeIDs[n.ID] = true
	}
	for _, e := range req.Edges {
		if e.Source == e.Target {
			resp.Valid = false
			resp.Errors = append(resp.Errors, "self-loop edge: "+e.ID)
			continue
		}
		if !nodeIDs[e.Source] {
			resp.Valid = false
			resp.Errors = append(resp.Errors, "edge source not found: "+e.Source)
		}
		if !nodeIDs[e.Target] {
			resp.Valid = false
			resp.Errors = append(resp.Errors, "edge target not found: "+e.Target)
		}
	}
	// Cycle detection (Kahn's algorithm).
	if resp.Valid {
		if inCycle(req.Nodes, req.Edges) {
			resp.Valid = false
			resp.Errors = append(resp.Errors, "graph contains a cycle")
		}
	}
	return resp
}

func inCycle(nodes []presenter.FlowNode, edges []presenter.FlowEdge) bool {
	indeg := map[string]int{}
	adj := map[string][]string{}
	for _, n := range nodes {
		indeg[n.ID] = 0
	}
	for _, e := range edges {
		indeg[e.Target]++
		adj[e.Source] = append(adj[e.Source], e.Target)
	}
	queue := []string{}
	for id, d := range indeg {
		if d == 0 {
			queue = append(queue, id)
		}
	}
	visited := 0
	for len(queue) > 0 {
		id := queue[0]
		queue = queue[1:]
		visited++
		for _, next := range adj[id] {
			indeg[next]--
			if indeg[next] == 0 {
				queue = append(queue, next)
			}
		}
	}
	return visited < len(nodes)
}

// topologicalOrder returns nodes in dependency order, or an error on a cycle.
func topologicalOrder(flow *presenter.Flow) ([]presenter.FlowNode, error) {
	if inCycle(flow.Nodes, flow.Edges) {
		return nil, ErrInvalidFlow
	}
	byID := map[string]presenter.FlowNode{}
	indeg := map[string]int{}
	adj := map[string][]string{}
	for _, n := range flow.Nodes {
		byID[n.ID] = n
		indeg[n.ID] = 0
	}
	for _, e := range flow.Edges {
		indeg[e.Target]++
		adj[e.Source] = append(adj[e.Source], e.Target)
	}
	queue := []string{}
	for id, d := range indeg {
		if d == 0 {
			queue = append(queue, id)
		}
	}
	out := []presenter.FlowNode{}
	for len(queue) > 0 {
		id := queue[0]
		queue = queue[1:]
		out = append(out, byID[id])
		for _, next := range adj[id] {
			indeg[next]--
			if indeg[next] == 0 {
				queue = append(queue, next)
			}
		}
	}
	return out, nil
}

func toPresenterFlow(rec *models.FlowDefinition) (*presenter.Flow, error) {
	var nodes []presenter.FlowNode
	var edges []presenter.FlowEdge
	if len(rec.Nodes) > 0 {
		if err := json.Unmarshal(rec.Nodes, &nodes); err != nil {
			return nil, err
		}
	}
	if len(rec.Edges) > 0 {
		if err := json.Unmarshal(rec.Edges, &edges); err != nil {
			return nil, err
		}
	}
	return &presenter.Flow{ID: rec.ID, Name: rec.Name, Nodes: nodes, Edges: edges, Enabled: rec.Enabled}, nil
}
