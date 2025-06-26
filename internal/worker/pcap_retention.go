package worker

import (
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// PcapRetention manages the retention of PCAP files for scripting
type PcapRetention struct {
	retainedPcaps map[string]time.Time
	mu            sync.Mutex
	defaultTTL    time.Duration
}

// NewPcapRetention creates a new PCAP retention manager
func NewPcapRetention(defaultTTL time.Duration) *PcapRetention {
	return &PcapRetention{
		retainedPcaps: make(map[string]time.Time),
		defaultTTL:    defaultTTL,
	}
}

// RetainPcap marks a PCAP file for extended retention
func (pr *PcapRetention) RetainPcap(pcapPath string, ttl time.Duration) {
	pr.mu.Lock()
	defer pr.mu.Unlock()
	
	retentionTime := time.Now().Add(ttl)
	pr.retainedPcaps[pcapPath] = retentionTime
	log.Debug().
		Str("pcapPath", pcapPath).
		Time("retainedUntil", retentionTime).
		Msg("Extended PCAP retention for scripting")
}

// ShouldRetain checks if a PCAP file should be retained
func (pr *PcapRetention) ShouldRetain(pcapPath string) bool {
	pr.mu.Lock()
	defer pr.mu.Unlock()
	
	expiryTime, exists := pr.retainedPcaps[pcapPath]
	if !exists {
		return false
	}
	
	if time.Now().After(expiryTime) {
		delete(pr.retainedPcaps, pcapPath)
		return false
	}
	
	return true
}

// GetDefaultTTL returns the default retention TTL
func (pr *PcapRetention) GetDefaultTTL() time.Duration {
	return pr.defaultTTL
}

// CleanupExpired removes expired entries from the retention map
func (pr *PcapRetention) CleanupExpired() {
	pr.mu.Lock()
	defer pr.mu.Unlock()
	
	now := time.Now()
	for path, expiry := range pr.retainedPcaps {
		if now.After(expiry) {
			delete(pr.retainedPcaps, path)
			log.Debug().
				Str("pcapPath", path).
				Msg("Removed expired PCAP retention entry")
		}
	}
}
