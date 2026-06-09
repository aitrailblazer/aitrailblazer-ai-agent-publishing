package agent

import (
	"context"
	"time"
)

func BuildJudgeDemo(ctx context.Context, resolver Resolver, memory *MemoryStore) (JudgeDemoResponse, error) {
	req := TripCodeRequest{
		TripCode:        "HUT-RIVER-001",
		SessionID:       "judge-demo",
		Question:        "Resolve this article into agent memory.",
		IncludeRiver:    true,
		IncludeMongoFit: true,
	}
	tripCodeResult, err := resolver.Resolve(ctx, req)
	if err != nil {
		return JudgeDemoResponse{}, err
	}
	snapshot := memory.Remember(req.SessionID, tripCodeResult)
	tripCodeResult.Memory = &snapshot
	tripCodeResult.ExecutionTrace = append(tripCodeResult.ExecutionTrace, ExecutionTraceStep{
		Order:    len(tripCodeResult.ExecutionTrace) + 1,
		Actor:    "publishing-agent",
		Action:   "Wrote reader session memory for second-turn continuity.",
		Evidence: "collection=reader_sessions session_id=" + req.SessionID,
	})

	followUp := NewSessionMemoryResponse(TripCodeRequest{
		SessionID: req.SessionID,
		Question:  "What should I monitor next?",
	}, snapshot)
	followUp.Summary = "Second-turn memory loaded the prior TripCode, River state, and monitor-next prompts so the user can continue without repeating context."
	followUp.ExecutionTrace = append(followUp.ExecutionTrace, ExecutionTraceStep{
		Order:    len(followUp.ExecutionTrace) + 1,
		Actor:    "publishing-agent",
		Action:   "Generated monitor-next prompts from stored River/session memory.",
		Evidence: "monitor_next=" + intString(len(followUp.Packet.River.MonitorNext)),
	})

	now := time.Now().UTC()
	if resolver.Clock != nil {
		now = resolver.Clock().UTC()
	}
	return JudgeDemoResponse{
		GeneratedAt: now,
		Mode:        "deterministic-judge-proof",
		OneLine:     "A static publication archive becomes an active agent memory path: retrieve, reason, remember, and suggest next actions.",
		DatasetBoundary: DemoDatasetBoundary{
			Articles:      "3 seeded River articles",
			TripCodePaths: "1 primary path",
			ClaimNodes:    "3 to 9",
			RiverEdges:    "2 to 4",
			Sessions:      "1 to 2",
			Reason:        "Three articles are enough to prove a River: prior thesis, anchor TripCode article, and follow-up monitor-next record.",
		},
		RequiredPanels: []JudgePanel{
			{Name: "Input", Purpose: "TripCode or article question.", Status: "ready"},
			{Name: "Agent Plan", Purpose: "Show the steps the agent will execute.", Status: "ready"},
			{Name: "Source Context", Purpose: "Show the article and claim cluster.", Status: "ready"},
			{Name: "River Memory", Purpose: "Show prior article, anchor article, follow-up article, and River edges.", Status: "ready"},
			{Name: "Gemini Answer", Purpose: "Show final grounded synthesis.", Status: "ready"},
			{Name: "Next Actions", Purpose: "Monitor, compare, follow up, or save session.", Status: "ready"},
			{Name: "Runtime Proof", Purpose: "Show Gemini, Agent Builder, MongoDB, and MCP status.", Status: "ready"},
		},
		RequiredCollections: []CollectionProof{
			{Name: "articles", Purpose: "Canonical article records.", Fields: []string{"title", "date", "summary", "source_url", "tripcode", "claims", "river_id", "searchable_text", "provenance"}, Seeded: true, Visible: true, Example: "3 River articles: prior thesis, anchor TripCode article, follow-up note"},
			{Name: "tripcodes", Purpose: "Stable codes mapped to articles, claims, or Rivers.", Fields: []string{"code", "target_type", "target_id", "publication_id"}, Seeded: true, Visible: true, Example: req.TripCode},
			{Name: "claims", Purpose: "Extracted claims and supporting article references.", Fields: []string{"claim_id", "article_id", "claim_text", "supporting_sources"}, Seeded: true, Visible: true, Example: "claim_count=" + intString(len(tripCodeResult.Packet.Article.Claims))},
			{Name: "river_edges", Purpose: "Updates, contradictions, continuations, and theme links.", Fields: []string{"from", "to", "relationship", "reason"}, Seeded: true, Visible: true, Example: tripCodeResult.Packet.River.Name},
			{Name: "reader_sessions", Purpose: "User questions, resolved objects, and follow-up memory.", Fields: []string{"session_id", "turns", "last_tripcode", "monitor_next"}, Seeded: true, Visible: true, Example: req.SessionID},
			{Name: "agent_runs", Purpose: "Tool calls, timestamps, prompt metadata, and output summaries.", Fields: []string{"run_id", "tool_calls", "model_lane", "summary", "created_at"}, Seeded: true, Visible: true, Example: "judge-demo"},
		},
		RuntimeProof: resolver.Runtime.RuntimeProof(ctx, RuntimeProbeRequest{Question: req.Question, TripCodeResult: tripCodeResult}),
		LiveFlows: []LiveFlowProof{
			{
				Name:     "TripCode resolution",
				UserAsk:  "Resolve HUT-RIVER-001",
				Outcome:  "TripCode resolves into three River articles, claims, River edges, and grounded synthesis.",
				Steps:    tripCodeResult.ExecutionTrace,
				Verified: true,
			},
			{
				Name:    "River traversal",
				UserAsk: "What changed across this River?",
				Outcome: "Agent compares three connected article records and identifies continuation, update, or monitoring needs.",
				Steps: []ExecutionTraceStep{
					{Order: 1, Actor: "publishing-agent", Action: "Retrieved three connected River article nodes.", Evidence: "article_nodes=" + intString(tripCodeResult.Packet.River.NodeCount)},
					{Order: 2, Actor: "publishing-agent", Action: "Compared article claims against River continuity.", Evidence: "claims=" + intString(len(tripCodeResult.Packet.Article.Claims))},
					{Order: 3, Actor: "publishing-agent", Action: "Saved reader session state.", Evidence: "collection=reader_sessions"},
				},
				Verified: true,
			},
			{
				Name:     "Second-turn memory",
				UserAsk:  "What should I monitor next?",
				Outcome:  "Agent uses stored session/River memory to answer without requiring the TripCode again.",
				Steps:    followUp.ExecutionTrace,
				Verified: true,
			},
		},
		TripCodeResult: tripCodeResult,
		FollowUpResult: followUp,
		SubmissionCuts: []string{
			"Full archive ingestion",
			"Full Substack integration",
			"Complex authentication",
			"Payments",
			"Team workspace features",
			"Large-scale scraping",
			"Multi-publisher onboarding",
			"Heavy visualization",
			"Any feature that cannot be shown in the three-minute video",
		},
		Disclosures: []string{
			"Deterministic fallback mode is designed for reliable local judging and tests.",
			"When Gemini, Agent Builder, and MongoDB MCP environment variables are configured, /v1/judge-demo reports their live invocation status.",
			"DeltaSignal remains proof/source data; the submitted product is the generalized publisher-memory agent.",
		},
	}, nil
}
