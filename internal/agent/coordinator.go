package agent

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Coordinator struct {
	Clock func() time.Time
	Err   error
}

func (c Coordinator) BuildArchiveBrief(_ context.Context, req ArchiveBriefRequest) (ArchiveBriefResponse, error) {
	if c.Err != nil {
		return ArchiveBriefResponse{}, c.Err
	}
	profile := publicationFromRequest(req)
	plan := []string{
		"Inventory the publication archive and normalize article URLs into stable article objects.",
		"Mint TripCodes for articles, claim clusters, topic clusters, and follow-up prompts.",
		"Build Rivers that connect articles, topics, claims, evidence, and thesis evolution.",
		"Persist operational records plus vector and semantic indexes for agent retrieval.",
		"Return a publication-ready agent workflow that can resolve one TripCode into context.",
	}
	findings := []SpecialistResult{
		{
			Agent:      "archive-ingestion-agent",
			Confidence: "deterministic-demo",
			Summary:    "DeltaSignal proves the seed pattern with a 300+ article macro and financial Substack archive.",
			Evidence: []Evidence{{
				Source:      "deltasignal-substack",
				Title:       "DeltaSignal Substack",
				URL:         "https://deltasignal.substack.com/",
				Observation: "The archive is the proof case for converting a serious publication into agent-readable memory.",
				Caveats:     []string{"Article count is user-provided context and should be verified before final submission copy."},
			}},
		},
		{
			Agent:      "tripcode-river-agent",
			Confidence: "deterministic-demo",
			Summary:    "A TripCode acts as a stable resolver key; a River preserves continuity across prior articles, claims, evidence, and monitor-next items.",
			Evidence: []Evidence{{
				Source:      "deltasignal-proof-case",
				Title:       "Hut 8: The Re-Rating Has A Deadline",
				URL:         "https://deltasignal.substack.com/p/hut-8-the-re-rating-has-a-deadline",
				Observation: "The HUT article is the visible proof case for TripCode and River-style research continuity.",
			}},
		},
		{
			Agent:      "mongodb-memory-agent",
			Confidence: "devpost-resource-grounded",
			Summary:    "MongoDB is the closest Rapid Agent partner fit for persistent publication memory because Atlas combines operational, vector, and semantic data for AI workloads.",
			Evidence: []Evidence{{
				Source:      "rapid-agent-devpost-resources",
				Title:       "MongoDB resources",
				URL:         "https://rapid-agent.devpost.com/details/mongodb-resources",
				Observation: "The partner page frames MongoDB Atlas as a unified operational foundation and persistent memory layer for AI and agentic workloads.",
			}},
		},
	}
	brief := fmt.Sprintf(
		"AITrailblazer AI Agent Publishing turns %s into a resolvable agent memory layer. The demo path starts with DeltaSignal's 300+ article archive, mints TripCodes for important articles, builds Rivers for topic continuity, and uses MongoDB-style persistent semantic storage so Gemini or MCP-aware agents can answer with article context instead of one-off search.",
		profile.Name,
	)
	now := time.Now().UTC()
	if c.Clock != nil {
		now = c.Clock().UTC()
	}
	return ArchiveBriefResponse{
		Publication:    profile.Name,
		PublicationURL: profile.URL,
		GeneratedAt:    now,
		Mode:           "deterministic-demo",
		Plan:           plan,
		Findings:       findings,
		Brief:          brief,
		NextAction:     "Implement an ingestion demo: load several DeltaSignal URLs, mint TripCodes, store River edges, and resolve one TripCode into a judge-facing packet.",
		Disclosures: []string{
			"This is a competition publishing demo, not a production archive crawler.",
			"DeltaSignal is the proof case; the product thesis is generalized for other publications.",
			"MongoDB fit is grounded in the Rapid Agent partner resource page and should be verified again before final Devpost submission.",
		},
	}, nil
}

func publicationFromRequest(req ArchiveBriefRequest) PublicationProfile {
	name := strings.TrimSpace(req.Publication)
	url := strings.TrimSpace(req.PublicationURL)
	if name == "" {
		name = "DeltaSignal"
	}
	if url == "" {
		url = "https://deltasignal.substack.com/"
	}
	return PublicationProfile{
		Name:              name,
		URL:               url,
		KnownArticleCount: "300+",
		Topics:            []string{"macro", "financial markets", "public-company research", "crypto infrastructure", "issuer intelligence"},
	}
}

var errCoordinatorUnavailable = errors.New("coordinator unavailable")
