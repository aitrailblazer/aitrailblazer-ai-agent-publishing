package agent

import (
	"context"
	"errors"
	"strings"
	"time"
)

type Resolver struct {
	Clock   func() time.Time
	Err     error
	Runtime RuntimeIntegrations
}

func (r Resolver) Resolve(ctx context.Context, req TripCodeRequest) (TripCodeResponse, error) {
	if r.Err != nil {
		return TripCodeResponse{}, r.Err
	}
	tripcode := strings.TrimSpace(req.TripCode)
	if tripcode == "" {
		tripcode = "AIT-PUB-HUT-8-RERATE"
	}
	now := time.Now().UTC()
	if r.Clock != nil {
		now = r.Clock().UTC()
	}
	packet := demoPacket(tripcode, req)
	resp := TripCodeResponse{
		TripCode:    tripcode,
		GeneratedAt: now,
		Mode:        "deterministic-demo",
		Packet:      packet,
		Summary: "One TripCode resolves an article into publication identity, article claims, a three-article River, MongoDB memory fit, boundaries, and monitor-next prompts. " +
			"The generalized product lets other publishers mint the same kind of resolvable agent memory over their own archives.",
		Disclosures: []string{
			"TripCodes are publication memory resolver keys, not official evidence identifiers.",
			"Rivers preserve editorial and analytical continuity; they do not prove every claim without linked evidence.",
			"Production ingestion should preserve source URLs, timestamps, author identity, permissions, and publisher ownership boundaries.",
		},
		ExecutionTrace: []ExecutionTraceStep{
			{Order: 1, Actor: "judge-harness", Action: "Submitted a TripCode or archive-memory request.", Evidence: "tripcode=" + tripcode},
			{Order: 2, Actor: "publishing-agent", Action: "Loaded the DeltaSignal proof-case packet.", Evidence: packet.Article.URL},
			{Order: 3, Actor: "publishing-agent", Action: "Mapped three article records into one River path.", Evidence: "article_nodes=" + intString(packet.River.NodeCount)},
			{Order: 4, Actor: "publishing-agent", Action: "Mapped article, River, and MongoDB memory fit into a generalized publisher workflow.", Evidence: "collections=publications,articles,tripcodes,rivers,embeddings,sessions"},
		},
	}
	r.Runtime.ApplyToTripCode(ctx, RuntimeProbeRequest{Question: req.Question, TripCodeResult: resp}, &resp)
	return resp, nil
}

var errResolverUnavailable = errors.New("resolver unavailable")

func NewSessionMemoryResponse(req TripCodeRequest, snapshot MemorySnapshot) TripCodeResponse {
	tripcode := snapshot.LastTripCode
	if tripcode == "" {
		tripcode = req.TripCode
	}
	last := MemoryEntry{}
	if len(snapshot.Entries) > 0 {
		last = snapshot.Entries[len(snapshot.Entries)-1]
	}
	return TripCodeResponse{
		TripCode:    tripcode,
		GeneratedAt: time.Now().UTC(),
		Mode:        "session-memory",
		Packet: TripCodePacket{
			Publication: PublicationProfile{
				Name:              "DeltaSignal",
				URL:               "https://deltasignal.substack.com/",
				KnownArticleCount: "300+",
				Topics:            []string{"publication memory", "TripCodes", "Rivers"},
			},
			Article: ArticleObject{
				TripCode: tripcode,
				Title:    last.Title,
				URL:      "https://deltasignal.substack.com/p/hut-8-the-re-rating-has-a-deadline",
			},
			River: River{
				Name:        last.RiverName,
				NodeCount:   last.NodeCount,
				MonitorNext: append([]string(nil), last.MonitorNext...),
			},
			ResolverIdentity: ResolverIdentity{Product: "AITrailblazer AI Agent Publishing", Version: "v0.1", Mode: "session-memory"},
			Boundaries: []string{
				"Session memory is a compact continuity snapshot, not a full archive re-ingestion.",
				"Follow-up answers should point back to stored TripCode and River context.",
			},
		},
		Summary:     "Loaded prior TripCode and River context from the session memory snapshot.",
		Disclosures: []string{"Session memory is in-memory for the local demo and resets when the process restarts."},
		Memory:      &snapshot,
		ExecutionTrace: []ExecutionTraceStep{
			{Order: 1, Actor: "judge-harness", Action: "Submitted a follow-up without a TripCode.", Evidence: "session_id=" + req.SessionID},
			{Order: 2, Actor: "publishing-agent", Action: "Loaded compact publication River memory.", Evidence: "turns=" + intString(snapshot.Turns)},
		},
	}
}

func demoPacket(tripcode string, req TripCodeRequest) TripCodePacket {
	publicationURL := strings.TrimSpace(req.PublicationURL)
	if publicationURL == "" {
		publicationURL = "https://deltasignal.substack.com/"
	}
	return TripCodePacket{
		Publication: PublicationProfile{
			Name:              "DeltaSignal",
			URL:               publicationURL,
			KnownArticleCount: "300+",
			Topics:            []string{"macro", "financial markets", "issuer intelligence", "crypto infrastructure"},
		},
		Article: ArticleObject{
			TripCode: tripcode,
			Title:    "Hut 8: The Re-Rating Has A Deadline",
			URL:      "https://deltasignal.substack.com/p/hut-8-the-re-rating-has-a-deadline",
			Claims: []string{
				"A publication article can be promoted into a stable TripCode object.",
				"A TripCode can retrieve prior River context instead of treating the article as isolated prose.",
				"Agents can compare current article claims to prior claims, evidence links, and monitor-next prompts.",
			},
		},
		River: River{
			Name:      "HUT re-rating research River",
			NodeCount: 3,
			Nodes: []RiverNode{
				{ID: "article.prior", Kind: "article", Title: "HUT infrastructure setup note", Description: "Seeded prior article record that captures the starting thesis, evidence watchlist, and original assumptions."},
				{ID: "article.anchor", Kind: "article", Title: "Hut 8: The Re-Rating Has A Deadline", Description: "Anchor article for HUT-RIVER-001; the visible TripCode resolves here first."},
				{ID: "article.followup", Kind: "article", Title: "HUT monitor-next follow-up note", Description: "Seeded follow-up article record that turns the River into watch items and next research questions."},
			},
			MonitorNext: []string{
				"Which claims changed since the prior article?",
				"Which assumptions weakened or need new evidence?",
				"What should the reader monitor before the next update?",
			},
			ReaderPrompt: "Resolve this article, reconstruct its River, compare prior claims with current evidence, and show what changed.",
		},
		MongoDBFit: MongoDBFit{
			WhyClose: "MongoDB Atlas maps to the publication-memory layer because it can keep operational article records, vector embeddings, semantic search indexes, River edges, and session memory together.",
			Collections: []string{
				"publications",
				"articles",
				"tripcodes",
				"river_edges",
				"claim_clusters",
				"embeddings",
				"reader_sessions",
			},
			Capabilities: []string{
				"Atlas Search for archive search",
				"Vector Search for article and claim retrieval",
				"Aggregations for River reconstruction",
				"MongoDB MCP Server for agent/database access",
				"Voyage AI or Google embeddings for semantic memory",
			},
			PartnerSources: []string{
				"https://rapid-agent.devpost.com/details/mongodb-resources",
				"https://www.mongodb.com/docs/mcp-server/get-started/",
			},
		},
		ResolverIdentity: ResolverIdentity{
			Product: "AITrailblazer AI Agent Publishing",
			Version: "v0.1",
			Mode:    "deterministic-demo",
		},
		Boundaries: []string{
			"TripCodes identify publisher-owned memory objects.",
			"Rivers explain continuity and context; source evidence must remain separately cited.",
			"User publications stay owned by the publisher; production ingestion requires permission and access controls.",
		},
	}
}

func intString(value int) string {
	if value == 0 {
		return "0"
	}
	digits := make([]byte, 0, 8)
	for value > 0 {
		digits = append([]byte{byte('0' + value%10)}, digits...)
		value /= 10
	}
	return string(digits)
}
