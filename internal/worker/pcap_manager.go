package worker

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// PcapManager manages the PCAP files
type PcapManager struct {
	pcapDir       string
	pcapTTL       time.Duration
	storageLimit  int64
	retainedPcaps map[string]time.Time
	mu            sync.Mutex
}

// NewPcapManager creates a new PCAP manager
func NewPcapManager(pcapDir string, pcapTTL time.Duration, storageLimit int64) *PcapManager {
	return &PcapManager{
		pcapDir:       pcapDir,
		pcapTTL:       pcapTTL,
		storageLimit:  storageLimit,
		retainedPcaps: make(map[string]time.Time),
	}
}

// CleanupExpiredPcaps removes expired PCAP files
func (pm *PcapManager) CleanupExpiredPcaps() error {
	files, err := os.ReadDir(pm.pcapDir)
	if err != nil {
		return err
	}

	now := time.Now()
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		info, err := file.Info()
		if err != nil {
			log.Error().Err(err).Str("file", file.Name()).Msg("Failed to get file info")
			continue
		}

		// Skip files that are marked for retention
		if pm.isRetained(file.Name()) {
			continue
		}

		// Delete files older than pcapTTL
		if now.Sub(info.ModTime()) > pm.pcapTTL {
			filePath := filepath.Join(pm.pcapDir, file.Name())
			if err := os.Remove(filePath); err != nil {
				log.Error().Err(err).Str("file", filePath).Msg("Failed to remove expired PCAP file")
			} else {
				log.Debug().Str("file", filePath).Msg("Removed expired PCAP file")
			}
		}
	}

	return nil
}

// RetainPcap marks a PCAP file for retention
func (pm *PcapManager) RetainPcap(pcapName string, duration time.Duration) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pm.retainedPcaps[pcapName] = time.Now().Add(duration)
	log.Debug().Str("pcapName", pcapName).Dur("duration", duration).Msg("PCAP file marked for retention")
}

// isRetained checks if a PCAP file is marked for retention
func (pm *PcapManager) isRetained(pcapName string) bool {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	expireTime, exists := pm.retainedPcaps[pcapName]
	if !exists {
		return false
	}

	// If the retention period has expired, remove it from the map
	if time.Now().After(expireTime) {
		delete(pm.retainedPcaps, pcapName)
		return false
	}

	return true
}

// GetStorageUsage returns the current storage usage of PCAP files
func (pm *PcapManager) GetStorageUsage() (int64, error) {
	var totalSize int64

	err := filepath.Walk(pm.pcapDir, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			totalSize += info.Size()
		}

		return nil
	})

	return totalSize, err
}

// GetPcapDir returns the PCAP directory
func (pm *PcapManager) GetPcapDir() string {
	return pm.pcapDir
}

// IsRetained is an exported wrapper for isRetained (for testing)
func (pm *PcapManager) IsRetained(pcapName string) bool {
	return pm.isRetained(pcapName)
}
