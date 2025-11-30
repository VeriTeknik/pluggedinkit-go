package pluggedinkit

import (
	"context"
	"fmt"
)

// Agent represents a PAP agent
type Agent struct {
	UUID                   string                 `json:"uuid"`
	Name                   string                 `json:"name"`
	DNSName                string                 `json:"dns_name"`
	State                  string                 `json:"state"` // NEW, PROVISIONED, ACTIVE, DRAINING, TERMINATED, KILLED
	KubernetesNamespace    *string                `json:"kubernetes_namespace,omitempty"`
	KubernetesDeployment   *string                `json:"kubernetes_deployment,omitempty"`
	CreatedAt              string                 `json:"created_at"`
	ProvisionedAt          *string                `json:"provisioned_at,omitempty"`
	ActivatedAt            *string                `json:"activated_at,omitempty"`
	TerminatedAt           *string                `json:"terminated_at,omitempty"`
	LastHeartbeatAt        *string                `json:"last_heartbeat_at,omitempty"`
	Metadata               map[string]interface{} `json:"metadata,omitempty"`
}

// CreateAgentRequest represents a request to create a new agent
type CreateAgentRequest struct {
	Name        string                  `json:"name"`
	Description *string                 `json:"description,omitempty"`
	Image       *string                 `json:"image,omitempty"`
	Resources   *ResourceRequirements   `json:"resources,omitempty"`
}

// ResourceRequirements represents Kubernetes resource requirements
type ResourceRequirements struct {
	CPURequest    *string `json:"cpu_request,omitempty"`
	MemoryRequest *string `json:"memory_request,omitempty"`
	CPULimit      *string `json:"cpu_limit,omitempty"`
	MemoryLimit   *string `json:"memory_limit,omitempty"`
}

// CreateAgentResponse represents the response from creating an agent
type CreateAgentResponse struct {
	Agent      Agent              `json:"agent"`
	Deployment DeploymentResult   `json:"deployment"`
}

// DeploymentResult represents Kubernetes deployment result
type DeploymentResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Heartbeat represents an agent heartbeat (liveness only - zombie prevention)
type Heartbeat struct {
	Mode          string  `json:"mode"` // EMERGENCY, IDLE, SLEEP
	UptimeSeconds float64 `json:"uptime_seconds"`
	Timestamp     string  `json:"timestamp"`
}

// Metrics represents agent metrics (resource telemetry - separate from heartbeat)
type Metrics struct {
	CPUPercent      float64                `json:"cpu_percent"`
	MemoryMB        float64                `json:"memory_mb"`
	RequestsHandled int                    `json:"requests_handled"`
	Timestamp       string                 `json:"timestamp"`
	CustomMetrics   map[string]interface{} `json:"custom_metrics,omitempty"`
}

// LifecycleEvent represents an agent lifecycle event
type LifecycleEvent struct {
	EventType string                 `json:"event_type"`
	FromState *string                `json:"from_state"`
	ToState   string                 `json:"to_state"`
	Timestamp string                 `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// AgentDetails represents detailed agent information
type AgentDetails struct {
	Agent             Agent                  `json:"agent"`
	RecentHeartbeats  []Heartbeat            `json:"recentHeartbeats"`
	RecentMetrics     []Metrics              `json:"recentMetrics"`
	LifecycleEvents   []LifecycleEvent       `json:"lifecycleEvents"`
	KubernetesStatus  map[string]interface{} `json:"kubernetesStatus,omitempty"`
}

// ExportAgentRequest represents a request to export agent data
type ExportAgentRequest struct {
	IncludeTelemetry bool `json:"include_telemetry"`
	TelemetryLimit   int  `json:"telemetry_limit"`
}

// HeartbeatRequest represents a heartbeat submission
type HeartbeatRequest struct {
	Mode          string  `json:"mode"` // EMERGENCY, IDLE, SLEEP
	UptimeSeconds float64 `json:"uptime_seconds"`
}

// MetricsRequest represents a metrics submission
type MetricsRequest struct {
	CPUPercent      float64                `json:"cpu_percent"`
	MemoryMB        float64                `json:"memory_mb"`
	RequestsHandled int                    `json:"requests_handled"`
	CustomMetrics   map[string]interface{} `json:"custom_metrics,omitempty"`
}

// StandardResponse represents a standard API response
type StandardResponse struct {
	Message string `json:"message"`
}

// DeleteAgentResponse represents the response from deleting an agent
type DeleteAgentResponse struct {
	Message    string           `json:"message"`
	Kubernetes DeploymentResult `json:"kubernetes"`
}

// AgentsService handles agent operations
type AgentsService struct {
	client *Client
}

// List retrieves all PAP agents
func (s *AgentsService) List(ctx context.Context) ([]Agent, error) {
	var agents []Agent
	err := s.client.get(ctx, "/api/agents", &agents)
	return agents, err
}

// Create creates a new PAP agent
func (s *AgentsService) Create(ctx context.Context, request CreateAgentRequest) (*CreateAgentResponse, error) {
	var response CreateAgentResponse
	err := s.client.post(ctx, "/api/agents", request, &response)
	return &response, err
}

// Get retrieves details for a specific agent
func (s *AgentsService) Get(ctx context.Context, agentID string) (*AgentDetails, error) {
	var details AgentDetails
	path := fmt.Sprintf("/api/agents/%s", agentID)
	err := s.client.get(ctx, path, &details)
	return &details, err
}

// Delete terminates and deletes an agent
func (s *AgentsService) Delete(ctx context.Context, agentID string) (*DeleteAgentResponse, error) {
	var response DeleteAgentResponse
	path := fmt.Sprintf("/api/agents/%s", agentID)
	err := s.client.delete(ctx, path, &response)
	return &response, err
}

// Export exports agent data including telemetry
func (s *AgentsService) Export(ctx context.Context, agentID string, includeTelemetry bool, telemetryLimit int) (map[string]interface{}, error) {
	var response map[string]interface{}
	path := fmt.Sprintf("/api/agents/%s/export", agentID)
	request := ExportAgentRequest{
		IncludeTelemetry: includeTelemetry,
		TelemetryLimit:   telemetryLimit,
	}
	err := s.client.post(ctx, path, request, &response)
	return response, err
}

// Heartbeat submits a heartbeat for an agent
//
// CRITICAL: Heartbeats are liveness-only (PAP zombie prevention).
// Never include resource data (CPU, memory) in heartbeats - use Metrics() instead.
func (s *AgentsService) Heartbeat(ctx context.Context, agentID string, mode string, uptimeSeconds float64) (*StandardResponse, error) {
	var response StandardResponse
	path := fmt.Sprintf("/api/agents/%s/heartbeat", agentID)
	request := HeartbeatRequest{
		Mode:          mode,
		UptimeSeconds: uptimeSeconds,
	}
	err := s.client.post(ctx, path, request, &response)
	return &response, err
}

// Metrics submits metrics for an agent
//
// CRITICAL: Metrics are separate from heartbeats (PAP zombie prevention).
// Heartbeats are liveness-only, metrics are resource telemetry.
func (s *AgentsService) Metrics(ctx context.Context, agentID string, cpuPercent, memoryMB float64, requestsHandled int, customMetrics map[string]interface{}) (*StandardResponse, error) {
	var response StandardResponse
	path := fmt.Sprintf("/api/agents/%s/metrics", agentID)
	request := MetricsRequest{
		CPUPercent:      cpuPercent,
		MemoryMB:        memoryMB,
		RequestsHandled: requestsHandled,
		CustomMetrics:   customMetrics,
	}
	err := s.client.post(ctx, path, request, &response)
	return &response, err
}
