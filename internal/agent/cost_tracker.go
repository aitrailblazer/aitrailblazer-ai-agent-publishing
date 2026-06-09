package agent

import (
	"os"
	"strconv"
	"strings"
	"sync"
)

type CostTracker struct {
	mu       sync.Mutex
	enabled  bool
	source   string
	budget   float64
	spent    float64
	estimate map[string]float64
}

func CostTrackerFromEnv() *CostTracker {
	source := strings.TrimSpace(os.Getenv("PUBLISHING_COST_SOURCE"))
	if source == "" {
		source = "local-estimate"
	}
	return &CostTracker{
		enabled: strings.EqualFold(os.Getenv("PUBLISHING_COST_TRACKING"), "true"),
		source:  source,
		budget:  envFloat("PUBLISHING_GOOGLE_CREDIT_BUDGET_USD", 500),
		estimate: map[string]float64{
			"archive-brief":  envFloat("PUBLISHING_ESTIMATED_ARCHIVE_BRIEF_COST_USD", 0),
			"tripcode":       envFloat("PUBLISHING_ESTIMATED_TRIPCODE_COST_USD", 0),
			"judge-demo":     envFloat("PUBLISHING_ESTIMATED_JUDGE_DEMO_COST_USD", 0),
			"session-memory": envFloat("PUBLISHING_ESTIMATED_SESSION_MEMORY_COST_USD", 0),
		},
	}
}

func (c *CostTracker) Record(kind string) *CostSnapshot {
	if c == nil {
		return nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.enabled {
		return &CostSnapshot{
			Enabled:     false,
			Source:      c.source,
			Currency:    "USD",
			RequestKind: kind,
			Note:        "Cost tracking disabled. Enable PUBLISHING_COST_TRACKING=true for local estimate metadata.",
		}
	}
	requestCost := c.estimate[kind]
	c.spent += requestCost
	return c.snapshotLocked(kind, requestCost)
}

func (c *CostTracker) Snapshot() CostSnapshot {
	if c == nil {
		return CostSnapshot{Enabled: false, Source: "unavailable", Currency: "USD"}
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	return *c.snapshotLocked("", 0)
}

func (c *CostTracker) snapshotLocked(kind string, requestCost float64) *CostSnapshot {
	remaining := c.budget - c.spent
	if remaining < 0 {
		remaining = 0
	}
	return &CostSnapshot{
		Enabled:         c.enabled,
		Source:          c.source,
		Currency:        "USD",
		RequestKind:     kind,
		RequestCostUSD:  requestCost,
		TrackedSpentUSD: c.spent,
		BudgetUSD:       c.budget,
		RemainingUSD:    remaining,
		Note:            "Local estimate only. Official Google Cloud spend requires Billing console or export.",
	}
}

func envFloat(key string, fallback float64) float64 {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	value, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return fallback
	}
	return value
}
