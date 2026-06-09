package agent

import "time"

type ArchiveBriefRequest struct {
	PublicationURL string `json:"publication_url,omitempty"`
	Publication    string `json:"publication,omitempty"`
	Question       string `json:"question,omitempty"`
}

type TripCodeRequest struct {
	TripCode        string `json:"tripcode,omitempty"`
	SessionID       string `json:"session_id,omitempty"`
	Question        string `json:"question,omitempty"`
	PublicationURL  string `json:"publication_url,omitempty"`
	IncludeRiver    bool   `json:"include_river,omitempty"`
	IncludeMongoFit bool   `json:"include_mongo_fit,omitempty"`
}

type ArchiveBriefResponse struct {
	Publication    string             `json:"publication"`
	PublicationURL string             `json:"publication_url"`
	GeneratedAt    time.Time          `json:"generated_at"`
	Mode           string             `json:"mode"`
	Plan           []string           `json:"plan"`
	Findings       []SpecialistResult `json:"findings"`
	Brief          string             `json:"brief"`
	NextAction     string             `json:"next_action"`
	Disclosures    []string           `json:"disclosures"`
	Cost           *CostSnapshot      `json:"cost,omitempty"`
}

type TripCodeResponse struct {
	TripCode       string               `json:"tripcode"`
	GeneratedAt    time.Time            `json:"generated_at"`
	Mode           string               `json:"mode"`
	Packet         TripCodePacket       `json:"packet"`
	Summary        string               `json:"summary"`
	Disclosures    []string             `json:"disclosures"`
	Memory         *MemorySnapshot      `json:"memory,omitempty"`
	ExecutionTrace []ExecutionTraceStep `json:"execution_trace,omitempty"`
	RuntimeProof   []RuntimeProofItem   `json:"runtime_proof,omitempty"`
	Cost           *CostSnapshot        `json:"cost,omitempty"`
}

type SpecialistResult struct {
	Agent      string     `json:"agent"`
	Summary    string     `json:"summary"`
	Confidence string     `json:"confidence"`
	Evidence   []Evidence `json:"evidence"`
}

type Evidence struct {
	Source      string   `json:"source"`
	Title       string   `json:"title"`
	Observation string   `json:"observation"`
	URL         string   `json:"url,omitempty"`
	SourceDate  string   `json:"source_date,omitempty"`
	Caveats     []string `json:"caveats,omitempty"`
}

type TripCodePacket struct {
	Publication      PublicationProfile `json:"publication"`
	Article          ArticleObject      `json:"article"`
	River            River              `json:"river"`
	MongoDBFit       MongoDBFit         `json:"mongodb_fit,omitempty"`
	ResolverIdentity ResolverIdentity   `json:"resolver_identity"`
	Boundaries       []string           `json:"boundaries"`
}

type PublicationProfile struct {
	Name              string   `json:"name"`
	URL               string   `json:"url"`
	KnownArticleCount string   `json:"known_article_count"`
	Topics            []string `json:"topics"`
}

type ArticleObject struct {
	TripCode string   `json:"tripcode"`
	Title    string   `json:"title"`
	URL      string   `json:"url"`
	Claims   []string `json:"claims"`
}

type River struct {
	Name         string      `json:"name"`
	NodeCount    int         `json:"node_count"`
	Nodes        []RiverNode `json:"nodes"`
	MonitorNext  []string    `json:"monitor_next"`
	ReaderPrompt string      `json:"reader_prompt"`
}

type RiverNode struct {
	ID          string `json:"id"`
	Kind        string `json:"kind"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type MongoDBFit struct {
	WhyClose       string   `json:"why_close"`
	Collections    []string `json:"collections"`
	Capabilities   []string `json:"capabilities"`
	PartnerSources []string `json:"partner_sources"`
}

type ResolverIdentity struct {
	Product string `json:"product"`
	Version string `json:"version"`
	Mode    string `json:"mode"`
}

type ExecutionTraceStep struct {
	Order    int    `json:"order"`
	Actor    string `json:"actor"`
	Action   string `json:"action"`
	Evidence string `json:"evidence,omitempty"`
}

type MemoryEntry struct {
	TripCode    string    `json:"tripcode"`
	Title       string    `json:"title"`
	RiverName   string    `json:"river_name"`
	NodeCount   int       `json:"node_count"`
	UpdatedAt   time.Time `json:"updated_at"`
	MonitorNext []string  `json:"monitor_next,omitempty"`
}

type MemorySnapshot struct {
	SessionID     string        `json:"session_id"`
	Available     bool          `json:"available"`
	Turns         int           `json:"turns"`
	LastTripCode  string        `json:"last_tripcode,omitempty"`
	LastUpdatedAt time.Time     `json:"last_updated_at,omitempty"`
	Entries       []MemoryEntry `json:"entries,omitempty"`
}

type CostSnapshot struct {
	Enabled         bool    `json:"enabled"`
	Source          string  `json:"source"`
	Currency        string  `json:"currency"`
	RequestKind     string  `json:"request_kind,omitempty"`
	RequestCostUSD  float64 `json:"request_cost_usd,omitempty"`
	TrackedSpentUSD float64 `json:"tracked_spent_usd,omitempty"`
	BudgetUSD       float64 `json:"budget_usd,omitempty"`
	RemainingUSD    float64 `json:"remaining_usd,omitempty"`
	Note            string  `json:"note,omitempty"`
}

type JudgeDemoResponse struct {
	GeneratedAt         time.Time           `json:"generated_at"`
	Mode                string              `json:"mode"`
	OneLine             string              `json:"one_line"`
	DatasetBoundary     DemoDatasetBoundary `json:"dataset_boundary"`
	RequiredPanels      []JudgePanel        `json:"required_panels"`
	RequiredCollections []CollectionProof   `json:"required_collections"`
	RuntimeProof        []RuntimeProofItem  `json:"runtime_proof"`
	LiveFlows           []LiveFlowProof     `json:"live_flows"`
	TripCodeResult      TripCodeResponse    `json:"tripcode_result"`
	FollowUpResult      TripCodeResponse    `json:"follow_up_result"`
	SubmissionCuts      []string            `json:"submission_cuts"`
	Disclosures         []string            `json:"disclosures"`
	Cost                *CostSnapshot       `json:"cost,omitempty"`
}

type DemoDatasetBoundary struct {
	Articles      string `json:"articles"`
	TripCodePaths string `json:"tripcode_paths"`
	ClaimNodes    string `json:"claim_nodes"`
	RiverEdges    string `json:"river_edges"`
	Sessions      string `json:"sessions"`
	Reason        string `json:"reason"`
}

type JudgePanel struct {
	Name    string `json:"name"`
	Purpose string `json:"purpose"`
	Status  string `json:"status"`
}

type CollectionProof struct {
	Name    string   `json:"name"`
	Purpose string   `json:"purpose"`
	Fields  []string `json:"fields"`
	Seeded  bool     `json:"seeded"`
	Visible bool     `json:"visible"`
	Example string   `json:"example,omitempty"`
}

type RuntimeProofItem struct {
	System   string `json:"system"`
	Status   string `json:"status"`
	Evidence string `json:"evidence"`
}

type LiveFlowProof struct {
	Name     string               `json:"name"`
	UserAsk  string               `json:"user_ask"`
	Outcome  string               `json:"outcome"`
	Steps    []ExecutionTraceStep `json:"steps"`
	Verified bool                 `json:"verified"`
}
