package agent

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func fixedTime() time.Time {
	return time.Date(2026, 6, 8, 12, 0, 0, 0, time.UTC)
}

func TestCoordinatorBuildArchiveBriefDefaultAndCustom(t *testing.T) {
	clock := func() time.Time { return fixedTime() }
	resp, err := (Coordinator{Clock: clock}).BuildArchiveBrief(context.Background(), ArchiveBriefRequest{})
	if err != nil {
		t.Fatalf("BuildArchiveBrief default error: %v", err)
	}
	if resp.Publication != "DeltaSignal" || resp.PublicationURL != "https://deltasignal.substack.com/" {
		t.Fatalf("default publication = %q %q", resp.Publication, resp.PublicationURL)
	}
	if !resp.GeneratedAt.Equal(fixedTime()) {
		t.Fatalf("generated = %s", resp.GeneratedAt)
	}
	if len(resp.Plan) != 5 || len(resp.Findings) != 3 || len(resp.Disclosures) != 3 {
		t.Fatalf("unexpected plan/findings/disclosures lengths: %d %d %d", len(resp.Plan), len(resp.Findings), len(resp.Disclosures))
	}

	custom, err := (Coordinator{}).BuildArchiveBrief(context.Background(), ArchiveBriefRequest{
		Publication:    "Research Desk",
		PublicationURL: "https://example.test/archive",
	})
	if err != nil {
		t.Fatalf("BuildArchiveBrief custom error: %v", err)
	}
	if custom.Publication != "Research Desk" || custom.PublicationURL != "https://example.test/archive" {
		t.Fatalf("custom publication = %q %q", custom.Publication, custom.PublicationURL)
	}
	if custom.Mode != "deterministic-demo" || custom.NextAction == "" || custom.Brief == "" {
		t.Fatalf("custom response incomplete: %#v", custom)
	}
}

func TestCoordinatorBuildArchiveBriefError(t *testing.T) {
	errExpected := errors.New("boom")
	_, err := (Coordinator{Err: errExpected}).BuildArchiveBrief(context.Background(), ArchiveBriefRequest{})
	if !errors.Is(err, errExpected) {
		t.Fatalf("error = %v, want %v", err, errExpected)
	}
	if errCoordinatorUnavailable == nil {
		t.Fatal("expected package sentinel to be initialized")
	}
}

func TestPublicationFromRequest(t *testing.T) {
	defaults := publicationFromRequest(ArchiveBriefRequest{Publication: "  ", PublicationURL: " "})
	if defaults.Name != "DeltaSignal" || defaults.URL != "https://deltasignal.substack.com/" {
		t.Fatalf("defaults = %#v", defaults)
	}
	custom := publicationFromRequest(ArchiveBriefRequest{Publication: "  Pub  ", PublicationURL: " https://pub.test "})
	if custom.Name != "Pub" || custom.URL != "https://pub.test" {
		t.Fatalf("custom = %#v", custom)
	}
	if custom.KnownArticleCount != "300+" || len(custom.Topics) != 5 {
		t.Fatalf("profile metadata = %#v", custom)
	}
}

func TestCostTrackerDisabledNilAndEnabled(t *testing.T) {
	t.Setenv("PUBLISHING_COST_TRACKING", "")
	t.Setenv("PUBLISHING_COST_SOURCE", "")
	disabled := CostTrackerFromEnv()
	snap := disabled.Record("tripcode")
	if snap == nil || snap.Enabled || snap.Source != "local-estimate" || snap.RequestKind != "tripcode" {
		t.Fatalf("disabled snapshot = %#v", snap)
	}
	if nilRecord := (*CostTracker)(nil).Record("tripcode"); nilRecord != nil {
		t.Fatalf("nil Record = %#v", nilRecord)
	}
	if nilSnapshot := (*CostTracker)(nil).Snapshot(); nilSnapshot.Source != "unavailable" || nilSnapshot.Enabled {
		t.Fatalf("nil Snapshot = %#v", nilSnapshot)
	}

	t.Setenv("PUBLISHING_COST_TRACKING", "true")
	t.Setenv("PUBLISHING_COST_SOURCE", "unit-test")
	t.Setenv("PUBLISHING_GOOGLE_CREDIT_BUDGET_USD", "0.03")
	t.Setenv("PUBLISHING_ESTIMATED_TRIPCODE_COST_USD", "0.02")
	t.Setenv("PUBLISHING_ESTIMATED_JUDGE_DEMO_COST_USD", "bad")
	enabled := CostTrackerFromEnv()
	first := enabled.Record("tripcode")
	second := enabled.Record("tripcode")
	unknown := enabled.Record("unknown")
	total := enabled.Snapshot()
	if !first.Enabled || first.Source != "unit-test" || first.RequestCostUSD != 0.02 || math.Abs(first.RemainingUSD-0.01) > 0.000001 {
		t.Fatalf("first = %#v", first)
	}
	if second.RemainingUSD != 0 || total.TrackedSpentUSD != 0.04 || unknown.RequestCostUSD != 0 {
		t.Fatalf("second/unknown/total = %#v %#v %#v", second, unknown, total)
	}
	if envFloat("missing", 7) != 7 || envFloat("PUBLISHING_ESTIMATED_JUDGE_DEMO_COST_USD", 9) != 9 {
		t.Fatal("envFloat fallback failed")
	}
	t.Setenv("PUBLISHING_ESTIMATED_JUDGE_DEMO_COST_USD", "1.25")
	if envFloat("PUBLISHING_ESTIMATED_JUDGE_DEMO_COST_USD", 9) != 1.25 {
		t.Fatal("envFloat parse failed")
	}
}

func TestMemoryStoreCapacitySnapshotAndCopies(t *testing.T) {
	store := NewMemoryStore(0)
	if store.capacity != 10 {
		t.Fatalf("default capacity = %d", store.capacity)
	}
	if empty := store.Remember(" ", TripCodeResponse{}); empty.Available || empty.SessionID != "" {
		t.Fatalf("empty remember = %#v", empty)
	}
	if snap := store.Snapshot(" missing "); snap.Available || snap.SessionID != "missing" {
		t.Fatalf("missing snapshot = %#v", snap)
	}

	capped := NewMemoryStore(2)
	resp := func(code, title string, at time.Time) TripCodeResponse {
		return TripCodeResponse{
			TripCode:    code,
			GeneratedAt: at,
			Packet: TripCodePacket{
				Article: ArticleObject{Title: title},
				River:   River{Name: "River", NodeCount: 3, MonitorNext: []string{"a", "b"}},
			},
		}
	}
	capped.Remember("s", resp("old", "Old", fixedTime().Add(2*time.Hour)))
	capped.Remember("s", resp("newer", "Newer", fixedTime().Add(1*time.Hour)))
	snapshot := capped.Remember("s", resp("latest", "Latest", fixedTime().Add(3*time.Hour)))
	if snapshot.Turns != 2 || snapshot.LastTripCode != "latest" || snapshot.Entries[0].TripCode != "newer" {
		t.Fatalf("snapshot = %#v", snapshot)
	}
	snapshot.Entries[0].TripCode = "mutated"
	if again := capped.Snapshot("s"); again.Entries[0].TripCode == "mutated" {
		t.Fatalf("snapshot was not copied: %#v", again)
	}
	snapshot.Entries[0].MonitorNext[0] = "changed"
	if again := capped.Snapshot("s"); again.Entries[0].MonitorNext[0] == "changed" {
		t.Fatalf("monitor-next was not copied: %#v", again)
	}
}

func TestResolverResolveDefaultCustomAndError(t *testing.T) {
	clock := func() time.Time { return fixedTime() }
	resp, err := (Resolver{Clock: clock}).Resolve(context.Background(), TripCodeRequest{PublicationURL: "https://pub.test"})
	if err != nil {
		t.Fatalf("Resolve error: %v", err)
	}
	if resp.TripCode != "AIT-PUB-HUT-8-RERATE" || !resp.GeneratedAt.Equal(fixedTime()) {
		t.Fatalf("default response = %#v", resp)
	}
	if resp.Packet.Publication.URL != "https://pub.test" || resp.Packet.River.NodeCount != 3 || len(resp.Packet.River.Nodes) != 3 {
		t.Fatalf("packet = %#v", resp.Packet)
	}
	custom, err := (Resolver{}).Resolve(context.Background(), TripCodeRequest{TripCode: " HUT-RIVER-001 "})
	if err != nil || custom.TripCode != "HUT-RIVER-001" || custom.Packet.Article.TripCode != "HUT-RIVER-001" {
		t.Fatalf("custom response = %#v err=%v", custom, err)
	}
	errExpected := errors.New("resolver boom")
	_, err = (Resolver{Err: errExpected}).Resolve(context.Background(), TripCodeRequest{})
	if !errors.Is(err, errExpected) {
		t.Fatalf("error = %v, want %v", err, errExpected)
	}
	if errResolverUnavailable == nil {
		t.Fatal("expected resolver sentinel to be initialized")
	}
}

func TestSessionMemoryResponseAndHelpers(t *testing.T) {
	empty := NewSessionMemoryResponse(TripCodeRequest{TripCode: "fallback", SessionID: "s"}, MemorySnapshot{SessionID: "s"})
	if empty.TripCode != "fallback" || empty.Packet.Article.Title != "" || empty.Mode != "session-memory" {
		t.Fatalf("empty session response = %#v", empty)
	}
	snapshot := MemorySnapshot{
		SessionID:    "s",
		LastTripCode: "HUT-RIVER-001",
		Entries: []MemoryEntry{{
			TripCode:    "HUT-RIVER-001",
			Title:       "Anchor",
			RiverName:   "River",
			NodeCount:   3,
			MonitorNext: []string{"next"},
		}},
		Turns: 1,
	}
	resp := NewSessionMemoryResponse(TripCodeRequest{SessionID: "s"}, snapshot)
	if resp.TripCode != "HUT-RIVER-001" || resp.Packet.Article.Title != "Anchor" || resp.Packet.River.MonitorNext[0] != "next" {
		t.Fatalf("session response = %#v", resp)
	}
	resp.Packet.River.MonitorNext[0] = "changed"
	if snapshot.Entries[0].MonitorNext[0] != "next" {
		t.Fatalf("source snapshot changed: %#v", snapshot)
	}
	if intString(0) != "0" || intString(12034) != "12034" {
		t.Fatal("intString failed")
	}
	if NowUTC().Location() != time.UTC {
		t.Fatal("NowUTC did not return UTC")
	}
}

func TestBuildJudgeDemoSuccessAndError(t *testing.T) {
	store := NewMemoryStore(5)
	resp, err := BuildJudgeDemo(context.Background(), Resolver{Clock: fixedTime}, store)
	if err != nil {
		t.Fatalf("BuildJudgeDemo error: %v", err)
	}
	if resp.Mode != "deterministic-judge-proof" || !resp.GeneratedAt.Equal(fixedTime()) {
		t.Fatalf("demo metadata = %#v", resp)
	}
	if resp.DatasetBoundary.Articles != "3 seeded River articles" || len(resp.RequiredPanels) != 7 || len(resp.RequiredCollections) != 6 || len(resp.LiveFlows) != 3 {
		t.Fatalf("demo shape = %#v", resp)
	}
	if resp.TripCodeResult.Packet.River.NodeCount != 3 || resp.FollowUpResult.Mode != "session-memory" {
		t.Fatalf("demo flow = %#v %#v", resp.TripCodeResult.Packet.River, resp.FollowUpResult.Mode)
	}
	if snap := store.Snapshot("judge-demo"); !snap.Available || snap.Turns != 1 {
		t.Fatalf("stored snapshot = %#v", snap)
	}

	errExpected := errors.New("demo resolver boom")
	_, err = BuildJudgeDemo(context.Background(), Resolver{Err: errExpected}, NewMemoryStore(1))
	if !errors.Is(err, errExpected) {
		t.Fatalf("error = %v, want %v", err, errExpected)
	}
}

func TestRuntimeIntegrationsFromEnvAndHelpers(t *testing.T) {
	t.Setenv("GOOGLE_CLOUD_PROJECT", "")
	t.Setenv("CLOUDSDK_CORE_PROJECT", "fallback-project")
	t.Setenv("GOOGLE_CLOUD_LOCATION", "")
	t.Setenv("GEMINI_MODEL", "")
	t.Setenv("PUBLISHING_USE_GEMINI", "yes")
	t.Setenv("GOOGLE_OAUTH_ACCESS_TOKEN", "google-token")
	t.Setenv("GOOGLE_METADATA_TOKEN_URL", "http://metadata.test")
	t.Setenv("GCE_METADATA_HOST", "metadata-host")
	t.Setenv("MCP_SERVER_URL", "http://mcp.test")
	t.Setenv("MCP_METHOD", "")
	t.Setenv("MCP_TOOL_NAME", "")
	t.Setenv("MCP_SESSION_ID", "")
	t.Setenv("MCP_AUTHORIZATION", "Bearer mcp")
	t.Setenv("MCP_API_KEY", "mcp-key")
	t.Setenv("AGENT_BUILDER_ENDPOINT", "")
	t.Setenv("AGENT_BUILDER_AGENT_URL", "http://agent.test")
	t.Setenv("AGENT_BUILDER_AGENT_ID", "agent-1")
	t.Setenv("AGENT_BUILDER_TOKEN", "")
	t.Setenv("AGENT_BUILDER_MODE", "discoveryengine-search")

	cfg := RuntimeIntegrationsFromEnv()
	if cfg.ProjectID != "fallback-project" || cfg.Location != "global" || cfg.GeminiModel != "gemini-2.5-flash" || !cfg.UseGemini {
		t.Fatalf("env config = %#v", cfg)
	}
	if cfg.AccessToken != "google-token" || cfg.MCPEndpoint != "http://mcp.test/mcp" || cfg.MCPMethod != "tools/call" || cfg.MCPTool != "find" || cfg.MCPSessionID != "aitrailblazer-judge-proof" || cfg.AgentEndpoint != "http://agent.test" || cfg.AgentMode != "discoveryengine-search" {
		t.Fatalf("env endpoints = %#v", cfg)
	}
	if envDefault("missing", "fallback") != "fallback" || firstEnv("missing-one", "missing-two") != "" || !envBool("PUBLISHING_USE_GEMINI") {
		t.Fatal("env helper fallback failed")
	}
	t.Setenv("RUNTIME_TEST_ENV_DEFAULT", "configured")
	if envDefault("RUNTIME_TEST_ENV_DEFAULT", "fallback") != "configured" {
		t.Fatal("envDefault configured value failed")
	}
	t.Setenv("PUBLISHING_USE_GEMINI", "no")
	if envBool("PUBLISHING_USE_GEMINI") {
		t.Fatal("envBool false failed")
	}
	if (RuntimeIntegrations{ProjectID: "p", Location: "us-central1", GeminiModel: "m"}).geminiURL() != "https://us-central1-aiplatform.googleapis.com/v1/projects/p/locations/us-central1/publishers/google/models/m:generateContent" {
		t.Fatal("regional gemini URL failed")
	}
	if (RuntimeIntegrations{ProjectID: "p"}).geminiURL() != "https://aiplatform.googleapis.com/v1/projects/p/locations/global/publishers/google/models/gemini-2.5-flash:generateContent" {
		t.Fatal("global gemini URL failed")
	}
	if (RuntimeIntegrations{}).httpClient() == nil || (RuntimeIntegrations{HTTPClient: http.DefaultClient}).httpClient() != http.DefaultClient {
		t.Fatal("httpClient helper failed")
	}
	if normalizeMCPEndpoint(" http://mcp.test/ ") != "http://mcp.test/mcp" || normalizeMCPEndpoint("http://mcp.test/mcp") != "http://mcp.test/mcp" || normalizeMCPEndpoint(" ") != "" {
		t.Fatal("normalizeMCPEndpoint failed")
	}
	if compactJSON(func() {}) != "unserializable JSON" || !strings.Contains(compactJSON(map[string]string{"x": "y"}), `"x":"y"`) {
		t.Fatal("compactJSON failed")
	}
}

func TestRuntimeIntegrationsGeminiBranches(t *testing.T) {
	req := RuntimeProbeRequest{Question: "What changed?", TripCodeResult: tripCodeFixture()}
	disabled := RuntimeIntegrations{}.invokeGemini(context.Background(), req)
	if disabled.Status != "deterministic fallback" {
		t.Fatalf("disabled = %#v", disabled)
	}
	missingProject := RuntimeIntegrations{UseGemini: true}.invokeGemini(context.Background(), req)
	if missingProject.Status != "configuration missing" {
		t.Fatalf("missing project = %#v", missingProject)
	}
	tokenErr := RuntimeIntegrations{UseGemini: true, ProjectID: "p", TokenURL: "://bad"}.invokeGemini(context.Background(), req)
	if tokenErr.Status != "token unavailable" {
		t.Fatalf("token err = %#v", tokenErr)
	}

	postStatus := RuntimeIntegrations{
		UseGemini:   true,
		ProjectID:   "p",
		AccessToken: "token",
		HTTPClient:  httpClientReturning(http.StatusBadGateway, `nope`),
	}.invokeGemini(context.Background(), req)
	if postStatus.Status != "invoke failed" {
		t.Fatalf("post status = %#v", postStatus)
	}
	empty := RuntimeIntegrations{
		UseGemini:   true,
		ProjectID:   "p",
		AccessToken: "token",
		HTTPClient:  httpClientReturning(http.StatusOK, `{"candidates":[{"content":{"parts":[{"text":"  "} ]}}]}`),
	}.invokeGemini(context.Background(), req)
	if empty.Status != "empty response" {
		t.Fatalf("empty = %#v", empty)
	}
	success := RuntimeIntegrations{
		UseGemini:   true,
		ProjectID:   "p",
		GeminiModel: "gemini-test",
		AccessToken: "token",
		HTTPClient: httpClientFunc(func(r *http.Request) (*http.Response, error) {
			if r.Header.Get("Authorization") != "Bearer token" {
				t.Fatalf("authorization = %q", r.Header.Get("Authorization"))
			}
			return jsonResponse(http.StatusOK, `{"candidates":[{"content":{"parts":[{"text":"Live synthesis"}]}}]}`), nil
		}),
	}
	result := success.invokeGemini(context.Background(), req)
	if result.Status != "live invoked" || result.Text != "Live synthesis" {
		t.Fatalf("success = %#v", result)
	}
	resp := tripCodeFixture()
	resp.RuntimeProof = nil
	success.ApplyToTripCode(context.Background(), req, &resp)
	if resp.Mode != "gemini-grounded" || resp.Summary != "Live synthesis" || len(resp.RuntimeProof) != 1 {
		t.Fatalf("applied = %#v", resp)
	}
	success.ApplyToTripCode(context.Background(), req, nil)

	var candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	}
	if firstGeminiText(candidates) != "" {
		t.Fatal("empty candidates returned text")
	}
}

func TestRuntimeIntegrationsMCPAgentAndProof(t *testing.T) {
	req := RuntimeProbeRequest{Question: "Resolve", TripCodeResult: tripCodeFixture()}
	none := RuntimeIntegrations{}
	proof := none.RuntimeProof(context.Background(), req)
	if len(proof) != 5 || proof[0].Status != "deterministic fallback" || proof[2].Status != "schema ready, MCP not configured" {
		t.Fatalf("fallback proof = %#v", proof)
	}

	mcpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/mcp" {
			t.Fatalf("mcp path = %q", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer mcp" || r.Header.Get("X-API-Key") != "key" || r.Header.Get("mcp-session-id") != "aitrailblazer-judge-proof" || r.Header.Get("Accept") != "application/json, text/event-stream" {
			t.Fatalf("mcp headers = auth:%q key:%q session:%q accept:%q", r.Header.Get("Authorization"), r.Header.Get("X-API-Key"), r.Header.Get("mcp-session-id"), r.Header.Get("Accept"))
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode mcp: %v", err)
		}
		if body["method"] != "tools/call" {
			t.Fatalf("mcp body = %#v", body)
		}
		params, ok := body["params"].(map[string]any)
		if !ok || params["name"] != "find" {
			t.Fatalf("mcp params = %#v", body["params"])
		}
		args, ok := params["arguments"].(map[string]any)
		if !ok || args["database"] != "aitrailblazer_demo" || args["collection"] != "tripcodes" || args["limit"] != float64(1) {
			t.Fatalf("mcp arguments = %#v", params["arguments"])
		}
		filter, ok := args["filter"].(map[string]any)
		if !ok || filter["tripcode"] != req.TripCodeResult.TripCode {
			t.Fatalf("mcp filter = %#v", args["filter"])
		}
		_, _ = w.Write([]byte(`{"jsonrpc":"2.0","id":"aitrailblazer-judge-proof","result":{"ok":true}}`))
	}))
	defer mcpServer.Close()
	mcp := RuntimeIntegrations{MCPEndpoint: mcpServer.URL, MCPMethod: "tools/call", MCPTool: "find", MCPAuthorization: "Bearer mcp", MCPAPIKey: "key"}
	mcpResult := mcp.invokeMCP(context.Background(), req)
	if mcpResult.Status != "live invoked" || mongodbAtlasStatus(mcpResult) != "live via MongoDB MCP" || !strings.Contains(mongodbAtlasEvidence(mcpResult), "proved") {
		t.Fatalf("mcp success = %#v", mcpResult)
	}
	t.Setenv("MONGODB_DEPLOYMENT_KIND", "atlas")
	if !strings.Contains(mongodbAtlasEvidence(mcpResult), "Atlas collection access") {
		t.Fatalf("atlas evidence = %q", mongodbAtlasEvidence(mcpResult))
	}
	t.Setenv("MONGODB_DEPLOYMENT_KIND", "")

	mcpErrorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"error":{"message":"tool failed"}}`))
	}))
	defer mcpErrorServer.Close()
	if result := (RuntimeIntegrations{MCPEndpoint: mcpErrorServer.URL}).invokeMCP(context.Background(), req); result.Status != "tool error" {
		t.Fatalf("mcp tool error = %#v", result)
	}
	if result := (RuntimeIntegrations{MCPEndpoint: "://bad"}).invokeMCP(context.Background(), req); result.Status != "invoke failed" {
		t.Fatalf("mcp invoke error = %#v", result)
	}

	agentServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer agent-token" {
			t.Fatalf("agent auth = %q", r.Header.Get("Authorization"))
		}
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer agentServer.Close()
	agent := RuntimeIntegrations{AgentEndpoint: agentServer.URL, AgentID: "agent-1", AgentToken: "agent-token"}
	agentResult := agent.invokeAgentBuilder(context.Background(), req)
	if agentResult.Status != "live invoked" || !strings.Contains(agentResult.Evidence, "agent-1") {
		t.Fatalf("agent success = %#v", agentResult)
	}
	searchAgent := RuntimeIntegrations{AgentEndpoint: "https://discoveryengine.googleapis.com/v1/projects/p/locations/global/collections/default_collection/engines/e/servingConfigs/default_search:search", AgentToken: "agent-token", HTTPClient: httpClientFunc(func(r *http.Request) (*http.Response, error) {
		if r.Header.Get("Authorization") != "Bearer agent-token" {
			t.Fatalf("search agent auth = %q", r.Header.Get("Authorization"))
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode search body: %v", err)
		}
		if !strings.Contains(body["query"].(string), "HUT-RIVER-001") {
			t.Fatalf("search query = %#v", body["query"])
		}
		return jsonResponse(http.StatusOK, `{"totalSize":3}`), nil
	})}
	searchResult := searchAgent.invokeAgentBuilder(context.Background(), req)
	if searchResult.Status != "live invoked" || !strings.Contains(searchResult.Evidence, "total_size=3") || !searchAgent.agentBuilderSearchMode() {
		t.Fatalf("search agent = %#v", searchResult)
	}
	searchWithMetadataToken := RuntimeIntegrations{AgentEndpoint: "http://agent-search.test", AgentMode: "discoveryengine-search", TokenURL: agentServer.URL, HTTPClient: httpClientFunc(func(r *http.Request) (*http.Response, error) {
		if r.Method == http.MethodGet {
			return jsonResponse(http.StatusOK, `{"access_token":"metadata-agent-token"}`), nil
		}
		if r.Header.Get("Authorization") != "Bearer metadata-agent-token" {
			t.Fatalf("metadata search auth = %q", r.Header.Get("Authorization"))
		}
		return jsonResponse(http.StatusOK, `{"totalSize":1}`), nil
	})}
	if result := searchWithMetadataToken.invokeAgentBuilder(context.Background(), req); result.Status != "live invoked" {
		t.Fatalf("metadata search result = %#v", result)
	}
	searchTokenError := RuntimeIntegrations{AgentEndpoint: "http://agent-search.test", AgentMode: "discoveryengine-search", TokenURL: "://bad"}
	if result := searchTokenError.invokeAgentBuilder(context.Background(), req); result.Status != "token unavailable" {
		t.Fatalf("search token error = %#v", result)
	}
	searchPostError := RuntimeIntegrations{AgentEndpoint: "http://agent-search.test", AgentMode: "discoveryengine-search", AgentToken: "token", HTTPClient: httpClientReturning(http.StatusInternalServerError, "bad")}
	if result := searchPostError.invokeAgentBuilder(context.Background(), req); result.Status != "invoke failed" {
		t.Fatalf("search post error = %#v", result)
	}
	if result := (RuntimeIntegrations{}).invokeAgentBuilder(context.Background(), req); result.Status != "not configured" {
		t.Fatalf("agent missing = %#v", result)
	}
	if result := (RuntimeIntegrations{AgentEndpoint: "://bad"}).invokeAgentBuilder(context.Background(), req); result.Status != "invoke failed" {
		t.Fatalf("agent invoke error = %#v", result)
	}
	if mongodbAtlasStatus(integrationResult{Status: "failed"}) != "MCP evidence unavailable" || !strings.Contains(mongodbAtlasEvidence(integrationResult{Evidence: "x"}), "Expected collections") {
		t.Fatal("mongo atlas fallback helpers failed")
	}
}

func TestRuntimeIntegrationHTTPAndTokenErrors(t *testing.T) {
	ctx := context.Background()
	rt := RuntimeIntegrations{HTTPClient: httpClientFunc(func(*http.Request) (*http.Response, error) {
		return nil, errors.New("transport failed")
	})}
	if err := rt.postJSON(ctx, "http://example.test", "", map[string]any{"ok": true}, nil, nil); err == nil {
		t.Fatal("expected transport error")
	}
	if err := (RuntimeIntegrations{}).postJSON(ctx, "://bad", "", map[string]any{"ok": true}, nil, nil); err == nil {
		t.Fatal("expected request error")
	}
	if err := (RuntimeIntegrations{}).postJSON(ctx, "http://example.test", "", map[string]any{"bad": func() {}}, nil, nil); err == nil {
		t.Fatal("expected marshal error")
	}
	readErr := RuntimeIntegrations{HTTPClient: httpClientFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Body: errReadCloser{}}, nil
	})}
	if err := readErr.postJSON(ctx, "http://example.test", "", map[string]any{"ok": true}, nil, nil); err == nil {
		t.Fatal("expected read error")
	}
	empty := RuntimeIntegrations{HTTPClient: httpClientReturning(http.StatusOK, "")}
	if err := empty.postJSON(ctx, "http://example.test", "", map[string]any{"ok": true}, &map[string]any{}, nil); err != nil {
		t.Fatalf("empty body err = %v", err)
	}
	okNoOutput := RuntimeIntegrations{HTTPClient: httpClientReturning(http.StatusOK, `{"ok":true}`)}
	if err := okNoOutput.postJSON(ctx, "http://example.test", "", map[string]any{"ok": true}, nil, nil); err != nil {
		t.Fatalf("nil output err = %v", err)
	}
	badJSON := RuntimeIntegrations{HTTPClient: httpClientReturning(http.StatusOK, "{")}
	if err := badJSON.postJSON(ctx, "http://example.test", "", map[string]any{"ok": true}, &map[string]any{}, nil); err == nil {
		t.Fatal("expected json error")
	}

	directToken := RuntimeIntegrations{AccessToken: "direct"}
	if token, err := directToken.accessToken(ctx); err != nil || token != "direct" {
		t.Fatalf("direct token = %q err=%v", token, err)
	}
	tokenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Metadata-Flavor") != "Google" {
			t.Fatalf("metadata header = %q", r.Header.Get("Metadata-Flavor"))
		}
		_, _ = w.Write([]byte(`{"access_token":"metadata-token"}`))
	}))
	defer tokenServer.Close()
	metadata := RuntimeIntegrations{TokenURL: tokenServer.URL}
	if token, err := metadata.accessToken(ctx); err != nil || token != "metadata-token" {
		t.Fatalf("metadata token = %q err=%v", token, err)
	}
	statusToken := RuntimeIntegrations{TokenURL: tokenServer.URL, HTTPClient: httpClientReturning(http.StatusForbidden, "no")}
	if _, err := statusToken.accessToken(ctx); err == nil {
		t.Fatal("expected token status error")
	}
	decodeToken := RuntimeIntegrations{TokenURL: tokenServer.URL, HTTPClient: httpClientReturning(http.StatusOK, "{")}
	if _, err := decodeToken.accessToken(ctx); err == nil {
		t.Fatal("expected token decode error")
	}
	missingToken := RuntimeIntegrations{TokenURL: tokenServer.URL, HTTPClient: httpClientReturning(http.StatusOK, "{}")}
	if _, err := missingToken.accessToken(ctx); err == nil {
		t.Fatal("expected missing token error")
	}
	doToken := RuntimeIntegrations{TokenURL: tokenServer.URL, HTTPClient: httpClientFunc(func(*http.Request) (*http.Response, error) {
		return nil, errors.New("no metadata")
	})}
	if _, err := doToken.accessToken(ctx); err == nil {
		t.Fatal("expected metadata transport error")
	}
	hostToken := RuntimeIntegrations{MetadataHost: tokenServer.Listener.Addr().String()}
	if token, err := hostToken.accessToken(ctx); err != nil || token != "metadata-token" {
		t.Fatalf("host token = %q err=%v", token, err)
	}
	defaultHostToken := RuntimeIntegrations{HTTPClient: httpClientFunc(func(r *http.Request) (*http.Response, error) {
		if !strings.Contains(r.URL.String(), "metadata.google.internal") {
			t.Fatalf("default metadata url = %s", r.URL.String())
		}
		return jsonResponse(http.StatusOK, `{"access_token":"default-host-token"}`), nil
	})}
	if token, err := defaultHostToken.accessToken(ctx); err != nil || token != "default-host-token" {
		t.Fatalf("default host token = %q err=%v", token, err)
	}
	if strings.Contains(geminiPrompt(RuntimeProbeRequest{Question: " q ", TripCodeResult: tripCodeFixture()}), "  q ") {
		t.Fatal("prompt did not trim question")
	}
}

func tripCodeFixture() TripCodeResponse {
	resp, err := (Resolver{Clock: fixedTime}).Resolve(context.Background(), TripCodeRequest{TripCode: "HUT-RIVER-001", Question: "Resolve"})
	if err != nil {
		panic(err)
	}
	return resp
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func httpClientFunc(f func(*http.Request) (*http.Response, error)) *http.Client {
	return &http.Client{Transport: roundTripFunc(f)}
}

func httpClientReturning(status int, body string) *http.Client {
	return httpClientFunc(func(*http.Request) (*http.Response, error) {
		return jsonResponse(status, body), nil
	})
}

func jsonResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

type errReadCloser struct{}

func (errReadCloser) Read([]byte) (int, error) {
	return 0, errors.New("read failed")
}

func (errReadCloser) Close() error {
	return nil
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
