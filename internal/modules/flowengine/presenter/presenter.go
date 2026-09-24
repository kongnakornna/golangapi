package presenter

// FlowNode is a single node in a workflow graph (Node-RED-like).
type FlowNode struct {
	ID   string         `json:"id"`
	Type string         `json:"type"` // mqtt-in|alarm|delay|email|line|discord|io-control|http
	Name string         `json:"name,omitempty"`
	Conf map[string]any `json:"conf,omitempty"`
}

// FlowEdge is a wire connecting two nodes.
type FlowEdge struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
}

// Flow is a persisted workflow definition.
type Flow struct {
	ID      string     `json:"id"`
	Name    string     `json:"name"`
	Nodes   []FlowNode `json:"nodes"`
	Edges   []FlowEdge `json:"edges"`
	Enabled bool       `json:"enabled"`
}

// CreateRequest is used to create or update a flow.
type CreateRequest struct {
	Name    string     `json:"name"`
	Nodes   []FlowNode `json:"nodes"`
	Edges   []FlowEdge `json:"edges"`
	Enabled *bool      `json:"enabled,omitempty"`
}

// ValidateResponse reports graph validation results.
type ValidateResponse struct {
	Valid  bool     `json:"valid"`
	Errors []string `json:"errors,omitempty"`
}

// TriggerRequest triggers a flow execution with an input payload.
type TriggerRequest struct {
	FlowID  string         `json:"flow_id"`
	Payload map[string]any `json:"payload,omitempty"`
}

// TriggerResponse reports execution status.
type TriggerResponse struct {
	OK       bool     `json:"ok"`
	Executed []string `json:"executed_nodes,omitempty"`
	Message  string   `json:"message,omitempty"`
}
