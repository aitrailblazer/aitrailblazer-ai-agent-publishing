package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/aitrailblazer/aitrailblazer-ai-agent-publishing/internal/agent"
)

func TestHealth(t *testing.T) {
	t.Setenv("PUBLISHING_DEMO_API_KEY", "")
	mux := testMux()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	var body map[string]any
	decodeJSON(t, rec.Body, &body)
	if body["service"] != "aitrailblazer-ai-agent-publishing" {
		t.Fatalf("service = %v", body["service"])
	}
	if rec.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Fatalf("missing CORS header")
	}
}

func TestCORSPreflight(t *testing.T) {
	t.Setenv("PUBLISHING_DEMO_API_KEY", "secret")
	mux := testMux()
	req := httptest.NewRequest(http.MethodOptions, "/v1/judge-demo", nil)
	req.Header.Set("Origin", "https://aitrailblazer.github.io")
	req.Header.Set("Access-Control-Request-Headers", "Content-Type")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("preflight status = %d", rec.Code)
	}
	if rec.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Fatalf("missing preflight CORS header")
	}
}

func TestStaticPublishingSurface(t *testing.T) {
	t.Setenv("PUBLISHING_DEMO_API_KEY", "")
	temp := t.TempDir()
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(temp); err != nil {
		t.Fatalf("chdir temp: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Fatalf("restore cwd: %v", err)
		}
	})
	if err := os.WriteFile("index.html", []byte("<!doctype html><title>AITrailblazer</title>"), 0o600); err != nil {
		t.Fatalf("write index: %v", err)
	}
	if err := os.WriteFile("DEVPOST_SUBMISSION.html", []byte("<!doctype html><title>Devpost</title>"), 0o600); err != nil {
		t.Fatalf("write devpost: %v", err)
	}
	if err := os.Mkdir("img", 0o700); err != nil {
		t.Fatalf("mkdir img: %v", err)
	}
	if err := os.WriteFile("img/AITrailblazerAI.png", []byte("png"), 0o600); err != nil {
		t.Fatalf("write img: %v", err)
	}
	if err := os.WriteFile("go.mod", []byte("module private"), 0o600); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}

	mux := testMux()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "AITrailblazer") {
		t.Fatalf("root static = %d %q", rec.Code, rec.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, "/img/AITrailblazerAI.png", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK || rec.Body.String() != "png" {
		t.Fatalf("image static = %d %q", rec.Code, rec.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, "/DEVPOST_SUBMISSION.html", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "Devpost") {
		t.Fatalf("devpost static = %d %q", rec.Code, rec.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, "/go.mod", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("disallowed static status = %d", rec.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.URL.Path = "../secret"
	rec = httptest.NewRecorder()
	serveStaticPublishingSurface(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("traversal static status = %d", rec.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/health", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK || !strings.Contains(rec.Header().Get("Content-Type"), "application/json") {
		t.Fatalf("health route = %d content-type=%q", rec.Code, rec.Header().Get("Content-Type"))
	}
}

func TestArchiveBrief(t *testing.T) {
	t.Setenv("PUBLISHING_DEMO_API_KEY", "")
	mux := testMux()
	req := httptest.NewRequest(http.MethodPost, "/v1/archive-brief", bytes.NewBufferString(`{"publication_url":"https://deltasignal.substack.com/"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	var body agent.ArchiveBriefResponse
	decodeJSON(t, rec.Body, &body)
	if body.Publication != "DeltaSignal" {
		t.Fatalf("publication = %q", body.Publication)
	}
	if len(body.Findings) != 3 {
		t.Fatalf("findings length = %d", len(body.Findings))
	}
}

func TestArchiveBriefInvalidAndFailure(t *testing.T) {
	t.Setenv("PUBLISHING_DEMO_API_KEY", "")
	mux := testMux()
	req := httptest.NewRequest(http.MethodPost, "/v1/archive-brief", bytes.NewBufferString(`{`))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("invalid status = %d", rec.Code)
	}

	mux = newMux(discardLogger(), agent.Coordinator{Err: errors.New("coordinator failed")}, agent.Resolver{}, agent.NewMemoryStore(1), agent.CostTrackerFromEnv())
	req = httptest.NewRequest(http.MethodPost, "/v1/archive-brief", bytes.NewBufferString(`{}`))
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("failure status = %d", rec.Code)
	}
}

func TestTripCodeAndSessionMemory(t *testing.T) {
	t.Setenv("PUBLISHING_DEMO_API_KEY", "")
	mux := testMux()
	req := httptest.NewRequest(http.MethodPost, "/v1/tripcode", bytes.NewBufferString(`{"tripcode":"AIT-PUB-HUT-8-RERATE","session_id":"demo"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("tripcode status = %d, body=%s", rec.Code, rec.Body.String())
	}
	var first agent.TripCodeResponse
	decodeJSON(t, rec.Body, &first)
	if first.Memory == nil || !first.Memory.Available {
		t.Fatalf("expected memory snapshot, got %#v", first.Memory)
	}

	req = httptest.NewRequest(http.MethodPost, "/resolve?session_id=demo", bytes.NewBufferString("What should stay in context?"))
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("memory status = %d, body=%s", rec.Code, rec.Body.String())
	}
	var second agent.TripCodeResponse
	decodeJSON(t, rec.Body, &second)
	if second.Mode != "session-memory" {
		t.Fatalf("mode = %q", second.Mode)
	}
}

func TestTripCodeInvalidUnauthorizedAndResolverFailure(t *testing.T) {
	t.Setenv("PUBLISHING_DEMO_API_KEY", "")
	mux := testMux()
	req := httptest.NewRequest(http.MethodPost, "/v1/tripcode", bytes.NewBufferString(`{`))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("invalid tripcode status = %d", rec.Code)
	}

	mux = newMux(discardLogger(), agent.Coordinator{}, agent.Resolver{Err: errors.New("resolver failed")}, agent.NewMemoryStore(1), agent.CostTrackerFromEnv())
	req = httptest.NewRequest(http.MethodPost, "/v1/tripcode", bytes.NewBufferString(`{"tripcode":"HUT-RIVER-001"}`))
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadGateway {
		t.Fatalf("resolver status = %d", rec.Code)
	}

	t.Setenv("PUBLISHING_DEMO_API_KEY", "secret")
	mux = testMux()
	req = httptest.NewRequest(http.MethodPost, "/v1/tripcode", bytes.NewBufferString(`{}`))
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("unauthorized tripcode status = %d", rec.Code)
	}
}

func TestResolveGetPostAndParsing(t *testing.T) {
	t.Setenv("PUBLISHING_DEMO_API_KEY", "")
	mux := testMux()
	req := httptest.NewRequest(http.MethodGet, "/resolve?tripcode=HUT-RIVER-001&session_id=s1&question=q&publication_url=https://pub.test&include_river=yes&include_mongo_fit=true", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET resolve status = %d body=%s", rec.Code, rec.Body.String())
	}
	var first agent.TripCodeResponse
	decodeJSON(t, rec.Body, &first)
	if first.Memory == nil || first.Packet.Publication.URL != "https://pub.test" {
		t.Fatalf("GET resolve body = %#v", first)
	}

	req = httptest.NewRequest(http.MethodPost, "/resolve?tripcode=override&session_id=s2", bytes.NewBufferString(`{"tripcode":"body","question":"body question","publication_url":"https://body.test","include_river":true,"include_mongo_fit":true}`))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("POST json resolve status = %d body=%s", rec.Code, rec.Body.String())
	}
	var merged agent.TripCodeResponse
	decodeJSON(t, rec.Body, &merged)
	if merged.TripCode != "override" || merged.Packet.Publication.URL != "https://body.test" || merged.Memory == nil {
		t.Fatalf("merged response = %#v", merged)
	}

	req = httptest.NewRequest(http.MethodPost, "/resolve", bytes.NewBufferString("What changed?"))
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("POST text resolve status = %d", rec.Code)
	}

	req = httptest.NewRequest(http.MethodPost, "/resolve?tripcode=blank", strings.NewReader("   "))
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("POST blank resolve status = %d", rec.Code)
	}

	req = httptest.NewRequest(http.MethodPost, "/resolve", bytes.NewBufferString(`{`))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("POST invalid json resolve status = %d", rec.Code)
	}

	req = httptest.NewRequest(http.MethodPost, "/resolve?session_id=missing", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("missing memory status = %d", rec.Code)
	}

	t.Setenv("PUBLISHING_DEMO_API_KEY", "secret")
	mux = testMux()
	req = httptest.NewRequest(http.MethodGet, "/resolve?tripcode=HUT-RIVER-001", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("GET resolve unauthorized status = %d", rec.Code)
	}
	req = httptest.NewRequest(http.MethodPost, "/resolve", bytes.NewBufferString("question"))
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("POST resolve unauthorized status = %d", rec.Code)
	}
}

func TestJudgeDemo(t *testing.T) {
	t.Setenv("PUBLISHING_DEMO_API_KEY", "")
	mux := testMux()
	req := httptest.NewRequest(http.MethodGet, "/v1/judge-demo", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	var body agent.JudgeDemoResponse
	decodeJSON(t, rec.Body, &body)
	if body.Mode != "deterministic-judge-proof" {
		t.Fatalf("mode = %q", body.Mode)
	}
	if len(body.RequiredPanels) != 7 {
		t.Fatalf("required panels = %d", len(body.RequiredPanels))
	}
	if len(body.RequiredCollections) != 6 {
		t.Fatalf("required collections = %d", len(body.RequiredCollections))
	}
	if len(body.LiveFlows) != 3 {
		t.Fatalf("live flows = %d", len(body.LiveFlows))
	}
	if body.TripCodeResult.Packet.River.NodeCount != 3 {
		t.Fatalf("river node count = %d, want 3", body.TripCodeResult.Packet.River.NodeCount)
	}
	articleNodes := 0
	for _, node := range body.TripCodeResult.Packet.River.Nodes {
		if node.Kind == "article" {
			articleNodes++
		}
	}
	if articleNodes != 3 {
		t.Fatalf("article nodes = %d, want 3", articleNodes)
	}
	for _, flow := range body.LiveFlows {
		if !flow.Verified {
			t.Fatalf("flow %q was not verified", flow.Name)
		}
	}
	if body.FollowUpResult.Mode != "session-memory" {
		t.Fatalf("follow-up mode = %q", body.FollowUpResult.Mode)
	}
}

func TestJudgeDemoFailureAndUsageUnauthorized(t *testing.T) {
	t.Setenv("PUBLISHING_DEMO_API_KEY", "")
	mux := newMux(discardLogger(), agent.Coordinator{}, agent.Resolver{Err: errors.New("resolver failed")}, agent.NewMemoryStore(1), agent.CostTrackerFromEnv())
	req := httptest.NewRequest(http.MethodGet, "/v1/judge-demo", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("judge failure status = %d", rec.Code)
	}

	t.Setenv("PUBLISHING_DEMO_API_KEY", "secret")
	mux = testMux()
	req = httptest.NewRequest(http.MethodGet, "/v1/judge-demo", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("judge unauthorized status = %d", rec.Code)
	}
	req = httptest.NewRequest(http.MethodGet, "/v1/usage", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("usage unauthorized status = %d", rec.Code)
	}
	req = httptest.NewRequest(http.MethodGet, "/v1/usage", nil)
	req.Header.Set("Authorization", "Bearer secret")
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("usage authorized status = %d", rec.Code)
	}
}

func TestDemoKeyProtection(t *testing.T) {
	t.Setenv("PUBLISHING_DEMO_API_KEY", "secret")
	mux := testMux()
	req := httptest.NewRequest(http.MethodPost, "/v1/archive-brief", bytes.NewBufferString(`{}`))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status without key = %d", rec.Code)
	}

	req = httptest.NewRequest(http.MethodPost, "/v1/archive-brief", bytes.NewBufferString(`{}`))
	req.Header.Set("X-Demo-Key", "secret")
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status with key = %d, body=%s", rec.Code, rec.Body.String())
	}
}

func TestRequestHelpers(t *testing.T) {
	if !queryBool(mapValues("x", "1"), "x") || !queryBool(mapValues("x", "TRUE"), "x") || !queryBool(mapValues("x", "yes"), "x") {
		t.Fatal("truthy queryBool failed")
	}
	if queryBool(mapValues("x", "no"), "x") || queryBool(mapValues("x", ""), "x") {
		t.Fatal("false queryBool failed")
	}
	base := agent.TripCodeRequest{TripCode: "base", SessionID: "s", Question: "q", PublicationURL: "u"}
	merged := mergeTripCodeRequest(base, agent.TripCodeRequest{IncludeRiver: true, IncludeMongoFit: true})
	if merged.TripCode != "base" || !merged.IncludeRiver || !merged.IncludeMongoFit {
		t.Fatalf("merge bool = %#v", merged)
	}
	merged = mergeTripCodeRequest(agent.TripCodeRequest{}, agent.TripCodeRequest{TripCode: "t", SessionID: "s", Question: "q", PublicationURL: "u"})
	if merged.TripCode != "t" || merged.SessionID != "s" || merged.Question != "q" || merged.PublicationURL != "u" {
		t.Fatalf("merge strings = %#v", merged)
	}
	req := httptest.NewRequest(http.MethodPost, "/resolve?tripcode=t", errReader{})
	if _, err := tripCodeRequestFromResolvePost(req); err == nil {
		t.Fatal("expected read error")
	}
	rec := httptest.NewRecorder()
	writeTripCodeResult(rec, tripCodeResult{status: http.StatusCreated, body: map[string]string{"ok": "yes"}})
	if rec.Code != http.StatusCreated || !strings.Contains(rec.Body.String(), "yes") {
		t.Fatalf("writeTripCodeResult = %d %s", rec.Code, rec.Body.String())
	}
}

func TestRunAndMain(t *testing.T) {
	origListen := listenAndServe
	origExit := exitProcess
	origOutput := outputWriter
	defer func() {
		listenAndServe = origListen
		exitProcess = origExit
		outputWriter = origOutput
	}()

	outputWriter = io.Discard
	t.Setenv("PORT", "")
	listenAndServe = func(addr string, _ http.Handler) error {
		if addr != ":8080" {
			t.Fatalf("addr = %q", addr)
		}
		return nil
	}
	if code := run(); code != 0 {
		t.Fatalf("run success code = %d", code)
	}

	t.Setenv("PORT", "9090")
	listenAndServe = func(addr string, _ http.Handler) error {
		if addr != ":9090" {
			t.Fatalf("addr = %q", addr)
		}
		return errors.New("listen failed")
	}
	if code := run(); code != 1 {
		t.Fatalf("run failure code = %d", code)
	}

	exitCode := -1
	exitProcess = func(code int) { exitCode = code }
	listenAndServe = func(string, http.Handler) error { return nil }
	main()
	if exitCode != 0 {
		t.Fatalf("main exit = %d", exitCode)
	}
}

func testMux() *http.ServeMux {
	return newMux(discardLogger(), agent.Coordinator{}, agent.Resolver{}, agent.NewMemoryStore(20), agent.CostTrackerFromEnv())
}

func decodeJSON(t *testing.T, reader io.Reader, out any) {
	t.Helper()
	if err := json.NewDecoder(reader).Decode(out); err != nil {
		t.Fatalf("decode json: %v", err)
	}
}

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func mapValues(key, value string) map[string][]string {
	return map[string][]string{key: {value}}
}

type errReader struct{}

func (errReader) Read([]byte) (int, error) {
	return 0, errors.New("read failed")
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
