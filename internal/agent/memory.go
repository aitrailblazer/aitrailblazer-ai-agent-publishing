package agent

import (
	"sort"
	"strings"
	"sync"
	"time"
)

type MemoryStore struct {
	mu       sync.Mutex
	capacity int
	entries  map[string][]MemoryEntry
}

func NewMemoryStore(capacity int) *MemoryStore {
	if capacity <= 0 {
		capacity = 10
	}
	return &MemoryStore{capacity: capacity, entries: map[string][]MemoryEntry{}}
}

func (s *MemoryStore) Remember(sessionID string, resp TripCodeResponse) MemorySnapshot {
	sessionID = strings.TrimSpace(sessionID)
	if sessionID == "" {
		return MemorySnapshot{}
	}
	entry := MemoryEntry{
		TripCode:    resp.TripCode,
		Title:       resp.Packet.Article.Title,
		RiverName:   resp.Packet.River.Name,
		NodeCount:   resp.Packet.River.NodeCount,
		UpdatedAt:   resp.GeneratedAt,
		MonitorNext: append([]string(nil), resp.Packet.River.MonitorNext...),
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[sessionID] = append(s.entries[sessionID], entry)
	if len(s.entries[sessionID]) > s.capacity {
		s.entries[sessionID] = s.entries[sessionID][len(s.entries[sessionID])-s.capacity:]
	}
	return s.snapshotLocked(sessionID)
}

func (s *MemoryStore) Snapshot(sessionID string) MemorySnapshot {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.snapshotLocked(strings.TrimSpace(sessionID))
}

func (s *MemoryStore) snapshotLocked(sessionID string) MemorySnapshot {
	entries := append([]MemoryEntry(nil), s.entries[sessionID]...)
	for i := range entries {
		entries[i].MonitorNext = append([]string(nil), entries[i].MonitorNext...)
	}
	if len(entries) == 0 {
		return MemorySnapshot{SessionID: sessionID, Available: false}
	}
	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].UpdatedAt.Before(entries[j].UpdatedAt)
	})
	last := entries[len(entries)-1]
	return MemorySnapshot{
		SessionID:     sessionID,
		Available:     true,
		Turns:         len(entries),
		LastTripCode:  last.TripCode,
		LastUpdatedAt: last.UpdatedAt,
		Entries:       entries,
	}
}

func NowUTC() time.Time {
	return time.Now().UTC()
}
