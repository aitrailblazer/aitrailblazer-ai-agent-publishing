package main

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aitrailblazer/aitrailblazer-ai-agent-publishing/internal/agent"
)

var (
	listenAndServe           = http.ListenAndServe
	exitProcess              = os.Exit
	outputWriter   io.Writer = os.Stdout
)

func main() {
	exitProcess(run())
}

func run() int {
	logger := slog.New(slog.NewJSONHandler(outputWriter, nil))
	coordinator := agent.Coordinator{}
	resolver := agent.Resolver{Runtime: agent.RuntimeIntegrationsFromEnv()}
	memory := agent.NewMemoryStore(20)
	costTracker := agent.CostTrackerFromEnv()
	mux := newMux(logger, coordinator, resolver, memory, costTracker)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	logger.Info("starting AITrailblazer AI Agent Publishing", "port", port)
	if err := listenAndServe(":"+port, mux); err != nil {
		logger.Error("server stopped", "error", err)
		return 1
	}
	return 0
}

func newMux(
	logger *slog.Logger,
	coordinator agent.Coordinator,
	resolver agent.Resolver,
	memory *agent.MemoryStore,
	costTracker *agent.CostTracker,
) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("OPTIONS /", func(w http.ResponseWriter, _ *http.Request) {
		setCORSHeaders(w)
		w.WriteHeader(http.StatusNoContent)
	})
	healthHandler := func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{
			"ok":      true,
			"service": "aitrailblazer-ai-agent-publishing",
			"time":    time.Now().UTC(),
		})
	}
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /healthz", healthHandler)
	mux.HandleFunc("POST /v1/archive-brief", func(w http.ResponseWriter, r *http.Request) {
		if !authorizedDemoRequest(r) {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "missing or invalid demo key"})
			return
		}
		var req agent.ArchiveBriefRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON request"})
			return
		}
		resp, err := coordinator.BuildArchiveBrief(r.Context(), req)
		if err != nil {
			logger.Error("build archive brief failed", "error", err)
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to build archive brief"})
			return
		}
		resp.Cost = costTracker.Record("archive-brief")
		writeJSON(w, http.StatusOK, resp)
	})
	mux.HandleFunc("POST /v1/tripcode", func(w http.ResponseWriter, r *http.Request) {
		if !authorizedDemoRequest(r) {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "missing or invalid demo key"})
			return
		}
		var req agent.TripCodeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON request"})
			return
		}
		writeTripCodeResult(w, resolveTripCode(r.Context(), logger, req, resolver, memory, costTracker))
	})
	mux.HandleFunc("GET /resolve", func(w http.ResponseWriter, r *http.Request) {
		if !authorizedDemoRequest(r) {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "missing or invalid demo key"})
			return
		}
		req := tripCodeRequestFromQuery(r.URL.Query())
		writeTripCodeResult(w, resolveTripCode(r.Context(), logger, req, resolver, memory, costTracker))
	})
	mux.HandleFunc("POST /resolve", func(w http.ResponseWriter, r *http.Request) {
		if !authorizedDemoRequest(r) {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "missing or invalid demo key"})
			return
		}
		req, err := tripCodeRequestFromResolvePost(r)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid resolve request"})
			return
		}
		writeTripCodeResult(w, resolveTripCode(r.Context(), logger, req, resolver, memory, costTracker))
	})
	mux.HandleFunc("GET /v1/judge-demo", func(w http.ResponseWriter, r *http.Request) {
		if !authorizedDemoRequest(r) {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "missing or invalid demo key"})
			return
		}
		resp, err := agent.BuildJudgeDemo(r.Context(), resolver, memory)
		if err != nil {
			logger.Error("build judge demo failed", "error", err)
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to build judge demo"})
			return
		}
		resp.Cost = costTracker.Record("judge-demo")
		writeJSON(w, http.StatusOK, resp)
	})
	mux.HandleFunc("GET /v1/usage", func(w http.ResponseWriter, r *http.Request) {
		if !authorizedDemoRequest(r) {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "missing or invalid demo key"})
			return
		}
		writeJSON(w, http.StatusOK, costTracker.Snapshot())
	})
	mux.HandleFunc("GET /", serveStaticPublishingSurface)
	return mux
}

func serveStaticPublishingSurface(w http.ResponseWriter, r *http.Request) {
	clean := filepath.Clean(strings.TrimPrefix(r.URL.Path, "/"))
	if clean == "." {
		clean = "index.html"
	}
	if strings.HasPrefix(clean, "..") || strings.Contains(clean, string(filepath.Separator)+".."+string(filepath.Separator)) {
		http.NotFound(w, r)
		return
	}
	allowed := clean == "index.html" ||
		clean == "README.html" ||
		clean == "CHANGELOG.html" ||
		clean == "DEVPOST_SUBMISSION.html" ||
		clean == "EXECPLAN_FINAL_SUBMISSION_VIDEO.html" ||
		clean == "START_HERE.html" ||
		clean == "VIDEO_SLIDES.html" ||
		clean == "LICENSE" ||
		strings.HasPrefix(clean, "img"+string(filepath.Separator)) ||
		strings.HasPrefix(clean, "docs"+string(filepath.Separator))
	if !allowed {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, clean)
}

type tripCodeResult struct {
	status int
	body   any
}

func resolveTripCode(
	ctx context.Context,
	logger *slog.Logger,
	req agent.TripCodeRequest,
	resolver agent.Resolver,
	memory *agent.MemoryStore,
	costTracker *agent.CostTracker,
) tripCodeResult {
	if strings.TrimSpace(req.TripCode) == "" && strings.TrimSpace(req.SessionID) != "" {
		snapshot := memory.Snapshot(req.SessionID)
		if !snapshot.Available {
			return tripCodeResult{status: http.StatusNotFound, body: map[string]string{"error": "no TripCode session memory found"}}
		}
		resp := agent.NewSessionMemoryResponse(req, snapshot)
		resp.Cost = costTracker.Record("session-memory")
		return tripCodeResult{status: http.StatusOK, body: resp}
	}
	resp, err := resolver.Resolve(ctx, req)
	if err != nil {
		logger.Error("resolve tripcode failed", "error", err)
		return tripCodeResult{status: http.StatusBadGateway, body: map[string]string{"error": err.Error()}}
	}
	if strings.TrimSpace(req.SessionID) != "" {
		snapshot := memory.Remember(req.SessionID, resp)
		resp.Memory = &snapshot
		resp.ExecutionTrace = append(resp.ExecutionTrace, agent.ExecutionTraceStep{
			Order:    len(resp.ExecutionTrace) + 1,
			Actor:    "publishing-agent",
			Action:   "Stored compact TripCode/River continuity for a follow-up turn.",
			Evidence: "session_id=" + req.SessionID,
		})
	}
	resp.Cost = costTracker.Record("tripcode")
	return tripCodeResult{status: http.StatusOK, body: resp}
}

func writeTripCodeResult(w http.ResponseWriter, result tripCodeResult) {
	writeJSON(w, result.status, result.body)
}

func tripCodeRequestFromResolvePost(r *http.Request) (agent.TripCodeRequest, error) {
	req := tripCodeRequestFromQuery(r.URL.Query())
	raw, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		return agent.TripCodeRequest{}, err
	}
	if strings.TrimSpace(string(raw)) == "" {
		return req, nil
	}
	if strings.Contains(strings.ToLower(r.Header.Get("Content-Type")), "application/json") {
		var bodyReq agent.TripCodeRequest
		if err := json.Unmarshal(raw, &bodyReq); err != nil {
			return agent.TripCodeRequest{}, err
		}
		return mergeTripCodeRequest(bodyReq, req), nil
	}
	req.Question = strings.TrimSpace(string(raw))
	return req, nil
}

func tripCodeRequestFromQuery(values url.Values) agent.TripCodeRequest {
	return agent.TripCodeRequest{
		TripCode:        values.Get("tripcode"),
		SessionID:       values.Get("session_id"),
		Question:        values.Get("question"),
		PublicationURL:  values.Get("publication_url"),
		IncludeRiver:    queryBool(values, "include_river"),
		IncludeMongoFit: queryBool(values, "include_mongo_fit"),
	}
}

func mergeTripCodeRequest(base, override agent.TripCodeRequest) agent.TripCodeRequest {
	if strings.TrimSpace(override.TripCode) != "" {
		base.TripCode = override.TripCode
	}
	if strings.TrimSpace(override.SessionID) != "" {
		base.SessionID = override.SessionID
	}
	if strings.TrimSpace(override.Question) != "" {
		base.Question = override.Question
	}
	if strings.TrimSpace(override.PublicationURL) != "" {
		base.PublicationURL = override.PublicationURL
	}
	if override.IncludeRiver {
		base.IncludeRiver = true
	}
	if override.IncludeMongoFit {
		base.IncludeMongoFit = true
	}
	return base
}

func queryBool(values url.Values, key string) bool {
	value := strings.ToLower(strings.TrimSpace(values.Get(key)))
	return value == "1" || value == "true" || value == "yes"
}

func authorizedDemoRequest(r *http.Request) bool {
	expected := strings.TrimSpace(os.Getenv("PUBLISHING_DEMO_API_KEY"))
	if expected == "" {
		return true
	}
	if r.Header.Get("X-Demo-Key") == expected {
		return true
	}
	auth := strings.TrimSpace(r.Header.Get("Authorization"))
	return strings.EqualFold(auth, "Bearer "+expected)
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	setCORSHeaders(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func setCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Demo-Key")
	w.Header().Set("Access-Control-Max-Age", "3600")
}
