package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type RuntimeIntegrations struct {
	HTTPClient       *http.Client
	TokenURL         string
	MetadataHost     string
	AccessToken      string
	ProjectID        string
	Location         string
	GeminiModel      string
	UseGemini        bool
	MCPEndpoint      string
	MCPMethod        string
	MCPTool          string
	MCPSessionID     string
	MCPAuthorization string
	MCPAPIKey        string
	AgentEndpoint    string
	AgentID          string
	AgentToken       string
	AgentMode        string
}

type RuntimeProbeRequest struct {
	Question       string
	TripCodeResult TripCodeResponse
}

type integrationResult struct {
	System   string
	Status   string
	Evidence string
	Text     string
}

func RuntimeIntegrationsFromEnv() RuntimeIntegrations {
	return RuntimeIntegrations{
		ProjectID:        firstEnv("GOOGLE_CLOUD_PROJECT", "CLOUDSDK_CORE_PROJECT"),
		Location:         envDefault("GOOGLE_CLOUD_LOCATION", "global"),
		GeminiModel:      envDefault("GEMINI_MODEL", "gemini-2.5-flash"),
		UseGemini:        envBool("PUBLISHING_USE_GEMINI"),
		AccessToken:      strings.TrimSpace(os.Getenv("GOOGLE_OAUTH_ACCESS_TOKEN")),
		TokenURL:         strings.TrimSpace(os.Getenv("GOOGLE_METADATA_TOKEN_URL")),
		MetadataHost:     strings.TrimSpace(os.Getenv("GCE_METADATA_HOST")),
		MCPEndpoint:      normalizeMCPEndpoint(os.Getenv("MCP_SERVER_URL")),
		MCPMethod:        envDefault("MCP_METHOD", "tools/call"),
		MCPTool:          envDefault("MCP_TOOL_NAME", "find"),
		MCPSessionID:     envDefault("MCP_SESSION_ID", "aitrailblazer-judge-proof"),
		MCPAuthorization: strings.TrimSpace(os.Getenv("MCP_AUTHORIZATION")),
		MCPAPIKey:        strings.TrimSpace(os.Getenv("MCP_API_KEY")),
		AgentEndpoint:    firstEnv("AGENT_BUILDER_ENDPOINT", "AGENT_BUILDER_AGENT_URL"),
		AgentID:          strings.TrimSpace(os.Getenv("AGENT_BUILDER_AGENT_ID")),
		AgentToken:       firstEnv("AGENT_BUILDER_TOKEN", "GOOGLE_OAUTH_ACCESS_TOKEN"),
		AgentMode:        strings.TrimSpace(os.Getenv("AGENT_BUILDER_MODE")),
	}
}

func (r RuntimeIntegrations) ApplyToTripCode(ctx context.Context, req RuntimeProbeRequest, resp *TripCodeResponse) {
	if resp == nil {
		return
	}
	result := r.invokeGemini(ctx, req)
	resp.RuntimeProof = append(resp.RuntimeProof, RuntimeProofItem{System: result.System, Status: result.Status, Evidence: result.Evidence})
	if result.Text != "" {
		resp.Summary = result.Text
		resp.Mode = "gemini-grounded"
		resp.ExecutionTrace = append(resp.ExecutionTrace, ExecutionTraceStep{
			Order:    len(resp.ExecutionTrace) + 1,
			Actor:    "gemini",
			Action:   "Generated grounded synthesis from TripCode/River packet.",
			Evidence: result.Evidence,
		})
	}
}

func (r RuntimeIntegrations) RuntimeProof(ctx context.Context, req RuntimeProbeRequest) []RuntimeProofItem {
	results := []integrationResult{
		r.invokeGemini(ctx, req),
		r.invokeAgentBuilder(ctx, req),
		r.invokeMCP(ctx, req),
	}
	return []RuntimeProofItem{
		{System: "Gemini", Status: results[0].Status, Evidence: results[0].Evidence},
		{System: "Agent Builder", Status: results[1].Status, Evidence: results[1].Evidence},
		{System: "MongoDB", Status: mongodbAtlasStatus(results[2]), Evidence: mongodbAtlasEvidence(results[2])},
		{System: "MongoDB MCP", Status: results[2].Status, Evidence: results[2].Evidence},
		{System: "Cloud Run", Status: "deployment target", Evidence: "Container-ready HTTP service with health and judge-demo endpoints."},
	}
}

func (r RuntimeIntegrations) invokeGemini(ctx context.Context, req RuntimeProbeRequest) integrationResult {
	if !r.UseGemini {
		return integrationResult{System: "Gemini", Status: "deterministic fallback", Evidence: "Set PUBLISHING_USE_GEMINI=true with Google Cloud project, location, and model to invoke Gemini."}
	}
	if strings.TrimSpace(r.ProjectID) == "" {
		return integrationResult{System: "Gemini", Status: "configuration missing", Evidence: "GOOGLE_CLOUD_PROJECT is required for Vertex AI Gemini invocation."}
	}
	token, err := r.accessToken(ctx)
	if err != nil {
		return integrationResult{System: "Gemini", Status: "token unavailable", Evidence: err.Error()}
	}
	payload := map[string]any{
		"contents": []map[string]any{{
			"role": "user",
			"parts": []map[string]string{{
				"text": geminiPrompt(req),
			}},
		}},
		"generationConfig": map[string]any{"temperature": 0.2, "maxOutputTokens": 512},
	}
	var out struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}
	if err := r.postJSON(ctx, r.geminiURL(), token, payload, &out, nil); err != nil {
		return integrationResult{System: "Gemini", Status: "invoke failed", Evidence: err.Error()}
	}
	text := strings.TrimSpace(firstGeminiText(out.Candidates))
	if text == "" {
		return integrationResult{System: "Gemini", Status: "empty response", Evidence: "Gemini returned no text candidate."}
	}
	return integrationResult{System: "Gemini", Status: "live invoked", Evidence: "model=" + r.GeminiModel + " project=" + r.ProjectID, Text: text}
}

func (r RuntimeIntegrations) invokeMCP(ctx context.Context, req RuntimeProbeRequest) integrationResult {
	if strings.TrimSpace(r.MCPEndpoint) == "" {
		return integrationResult{System: "MongoDB MCP", Status: "not configured", Evidence: "Set MCP_SERVER_URL to invoke the selected partner MCP server at runtime."}
	}
	payload := map[string]any{
		"jsonrpc": "2.0",
		"id":      "aitrailblazer-judge-proof",
		"method":  r.MCPMethod,
		"params": map[string]any{
			"name": r.MCPTool,
			"arguments": map[string]any{
				"database":   envDefault("MONGODB_DATABASE", "aitrailblazer_demo"),
				"collection": "tripcodes",
				"filter": map[string]any{
					"tripcode": req.TripCodeResult.TripCode,
				},
				"projection": map[string]any{"_id": 0},
				"limit":      1,
			},
		},
	}
	headers := map[string]string{}
	sessionID := strings.TrimSpace(r.MCPSessionID)
	if sessionID == "" {
		sessionID = envDefault("MCP_SESSION_ID", "aitrailblazer-judge-proof")
	}
	if sessionID != "" {
		headers["mcp-session-id"] = sessionID
	}
	if r.MCPAuthorization != "" {
		headers["Authorization"] = r.MCPAuthorization
	}
	if r.MCPAPIKey != "" {
		headers["X-API-Key"] = r.MCPAPIKey
	}
	var out map[string]any
	headers["Accept"] = "application/json, text/event-stream"
	if err := r.postJSON(ctx, normalizeMCPEndpoint(r.MCPEndpoint), "", payload, &out, headers); err != nil {
		return integrationResult{System: "MongoDB MCP", Status: "invoke failed", Evidence: err.Error()}
	}
	if _, hasError := out["error"]; hasError {
		return integrationResult{System: "MongoDB MCP", Status: "tool error", Evidence: compactJSON(out)}
	}
	return integrationResult{System: "MongoDB MCP", Status: "live invoked", Evidence: "method=" + r.MCPMethod + " tool=" + r.MCPTool + " database=" + envDefault("MONGODB_DATABASE", "aitrailblazer_demo") + " collection=tripcodes tripcode=" + req.TripCodeResult.TripCode}
}

func (r RuntimeIntegrations) invokeAgentBuilder(ctx context.Context, req RuntimeProbeRequest) integrationResult {
	if strings.TrimSpace(r.AgentEndpoint) == "" {
		return integrationResult{System: "Agent Builder", Status: "not configured", Evidence: "Set AGENT_BUILDER_ENDPOINT or AGENT_BUILDER_AGENT_URL to invoke the orchestration lane at runtime."}
	}
	if r.agentBuilderSearchMode() {
		token := strings.TrimSpace(r.AgentToken)
		if token == "" {
			var err error
			token, err = r.accessToken(ctx)
			if err != nil {
				return integrationResult{System: "Agent Builder", Status: "token unavailable", Evidence: err.Error()}
			}
		}
		payload := map[string]any{
			"query":    strings.TrimSpace(req.TripCodeResult.TripCode + " " + req.Question),
			"pageSize": 3,
			"contentSearchSpec": map[string]any{
				"summarySpec": map[string]any{"summaryResultCount": 3, "includeCitations": true},
			},
		}
		var out map[string]any
		if err := r.postJSON(ctx, r.AgentEndpoint, token, payload, &out, nil); err != nil {
			return integrationResult{System: "Agent Builder", Status: "invoke failed", Evidence: err.Error()}
		}
		return integrationResult{System: "Agent Builder", Status: "live invoked", Evidence: "discovery_engine_search total_size=" + fmt.Sprint(out["totalSize"])}
	}
	payload := map[string]any{
		"agent_id": r.AgentID,
		"input": map[string]any{
			"tripcode": req.TripCodeResult.TripCode,
			"question": req.Question,
			"river":    req.TripCodeResult.Packet.River.Name,
		},
	}
	headers := map[string]string{}
	if r.AgentToken != "" {
		headers["Authorization"] = "Bearer " + r.AgentToken
	}
	var out map[string]any
	if err := r.postJSON(ctx, r.AgentEndpoint, "", payload, &out, headers); err != nil {
		return integrationResult{System: "Agent Builder", Status: "invoke failed", Evidence: err.Error()}
	}
	return integrationResult{System: "Agent Builder", Status: "live invoked", Evidence: "agent_id=" + strings.TrimSpace(r.AgentID)}
}

func (r RuntimeIntegrations) agentBuilderSearchMode() bool {
	mode := strings.ToLower(strings.TrimSpace(r.AgentMode))
	return mode == "discoveryengine-search" || strings.Contains(r.AgentEndpoint, "discoveryengine.googleapis.com")
}

func (r RuntimeIntegrations) accessToken(ctx context.Context) (string, error) {
	if strings.TrimSpace(r.AccessToken) != "" {
		return r.AccessToken, nil
	}
	tokenURL := r.TokenURL
	if tokenURL == "" {
		host := r.MetadataHost
		if host == "" {
			host = "metadata.google.internal"
		}
		tokenURL = "http://" + host + "/computeMetadata/v1/instance/service-accounts/default/token"
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, tokenURL, nil)
	if err != nil {
		return "", err
	}
	request.Header.Set("Metadata-Flavor", "Google")
	resp, err := r.httpClient().Do(request)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return "", fmt.Errorf("metadata token status=%d", resp.StatusCode)
	}
	var body struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&body); err != nil {
		return "", err
	}
	if strings.TrimSpace(body.AccessToken) == "" {
		return "", errors.New("metadata token response missing access_token")
	}
	return body.AccessToken, nil
}

func (r RuntimeIntegrations) postJSON(ctx context.Context, url, bearer string, payload any, out any, headers map[string]string) error {
	raw, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(raw))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	if bearer != "" {
		request.Header.Set("Authorization", "Bearer "+bearer)
	}
	for key, value := range headers {
		request.Header.Set(key, value)
	}
	resp, err := r.httpClient().Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return err
	}
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	if out == nil {
		return nil
	}
	if len(bytes.TrimSpace(body)) == 0 {
		return nil
	}
	return json.Unmarshal(body, out)
}

func (r RuntimeIntegrations) geminiURL() string {
	location := strings.TrimSpace(r.Location)
	if location == "" {
		location = "global"
	}
	endpoint := "https://" + location + "-aiplatform.googleapis.com"
	if location == "global" {
		endpoint = "https://aiplatform.googleapis.com"
	}
	model := strings.TrimSpace(r.GeminiModel)
	if model == "" {
		model = "gemini-2.5-flash"
	}
	return endpoint + "/v1/projects/" + r.ProjectID + "/locations/" + location + "/publishers/google/models/" + model + ":generateContent"
}

func (r RuntimeIntegrations) httpClient() *http.Client {
	if r.HTTPClient != nil {
		return r.HTTPClient
	}
	return &http.Client{Timeout: 10 * time.Second}
}

func geminiPrompt(req RuntimeProbeRequest) string {
	packet := req.TripCodeResult.Packet
	return "Generate a concise, source-grounded publishing-agent answer for TripCode " + req.TripCodeResult.TripCode +
		". Article: " + packet.Article.Title +
		". River: " + packet.River.Name +
		". User question: " + strings.TrimSpace(req.Question) +
		". Include what changed and what to monitor next."
}

func firstGeminiText(candidates []struct {
	Content struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"content"`
}) string {
	for _, candidate := range candidates {
		for _, part := range candidate.Content.Parts {
			if strings.TrimSpace(part.Text) != "" {
				return part.Text
			}
		}
	}
	return ""
}

func mongodbAtlasStatus(mcp integrationResult) string {
	if mcp.Status == "live invoked" {
		return "live via MongoDB MCP"
	}
	if mcp.Status == "not configured" {
		return "schema ready, MCP not configured"
	}
	return "MCP evidence unavailable"
}

func mongodbAtlasEvidence(mcp integrationResult) string {
	if mcp.Status == "live invoked" {
		if strings.ToLower(envDefault("MONGODB_DEPLOYMENT_KIND", "")) == "atlas" {
			return "MongoDB Atlas collection access proved through MCP tool invocation."
		}
		return "MongoDB collection access proved through official MCP tool invocation."
	}
	return "Expected collections: articles, tripcodes, claims, river_edges, reader_sessions, agent_runs. MCP evidence: " + mcp.Evidence
}

func compactJSON(value any) string {
	raw, err := json.Marshal(value)
	if err != nil {
		return "unserializable JSON"
	}
	return string(raw)
}

func envDefault(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func firstEnv(keys ...string) string {
	for _, key := range keys {
		if value := strings.TrimSpace(os.Getenv(key)); value != "" {
			return value
		}
	}
	return ""
}

func normalizeMCPEndpoint(value string) string {
	endpoint := strings.TrimSpace(value)
	if endpoint == "" {
		return ""
	}
	if strings.HasSuffix(endpoint, "/") {
		endpoint = strings.TrimRight(endpoint, "/")
	}
	if !strings.HasSuffix(endpoint, "/mcp") {
		endpoint += "/mcp"
	}
	return endpoint
}

func envBool(key string) bool {
	value := strings.ToLower(strings.TrimSpace(os.Getenv(key)))
	return value == "1" || value == "true" || value == "yes"
}
