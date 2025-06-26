package config

import (
	"time"

	"github.com/kubeshark/kubeshark/config/configStructs"
	"github.com/rs/zerolog/log"
)

// ValidateConfig validates the configuration
func ValidateConfig(config *configStructs.ConfigStruct) error {
	// Check for minimum PCAP TTL when scripting is enabled
	if config.Scripting.Enabled {
		// Parse the PCAP TTL duration
		ttlDuration, err := time.ParseDuration(config.Tap.Misc.PcapTTL)
		if err == nil && ttlDuration < 60*time.Second {
			log.Warn().Msgf("PCAP TTL is set to %s which may be too short for scripting. Consider setting it to at least 60s.", config.Tap.Misc.PcapTTL)
		}
	}

	// Validate tap configuration
	return nil
}
