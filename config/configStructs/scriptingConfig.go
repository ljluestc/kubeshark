package configStructs
package configStructs

// ScriptingConfig represents configuration for the scripting system
type ScriptingConfig struct {
	// Enabled indicates whether scripting is enabled
	Enabled bool `yaml:"enabled" json:"enabled" default:"true"`
	
	// Timeout is the maximum execution time for scripts in milliseconds
	Timeout int `yaml:"timeout" json:"timeout" default:"5000"`
	
	// MaxMemory is the maximum memory usage for scripts in megabytes
	MaxMemory int `yaml:"maxMemory" json:"maxMemory" default:"100"`
	
	// AllowFileSystem indicates whether scripts can access the file system
	AllowFileSystem bool `yaml:"allowFileSystem" json:"allowFileSystem" default:"false"`
	
	// AllowNetwork indicates whether scripts can access the network
	AllowNetwork bool `yaml:"allowNetwork" json:"allowNetwork" default:"false"`
}

// ScriptingPermissions represents user permissions for scripting
type ScriptingPermissions struct {
	// CanSave indicates whether the user can save scripts
	CanSave bool `yaml:"canSave" json:"canSave" default:"false"`
	
	// CanActivate indicates whether the user can activate scripts
	CanActivate bool `yaml:"canActivate" json:"canActivate" default:"false"`
	
	// CanDelete indicates whether the user can delete scripts
	CanDelete bool `yaml:"canDelete" json:"canDelete" default:"false"`
}
// ScriptingConfig defines configuration options for the scripting subsystem
type ScriptingConfig struct {
	// Enabled determines if scripting functionality is available
	Enabled bool `yaml:"enabled" json:"enabled" default:"false"`
	
	// RetentionLimit sets the maximum duration (in seconds) that scripts can retain PCAP files
	// Setting to 0 means no limit
	RetentionLimit int `yaml:"retentionLimit" json:"retentionLimit" default:"3600"`
	
	// ScriptsDir is the directory where scripts are stored
	ScriptsDir string `yaml:"scriptsDir" json:"scriptsDir" default:"scripts"`
	
	// MaxConcurrency is the maximum number of scripts that can run concurrently
	MaxConcurrency int `yaml:"maxConcurrency" json:"maxConcurrency" default:"10"`
	
	// Timeout is the maximum execution time for a script in milliseconds
	Timeout int `yaml:"timeout" json:"timeout" default:"5000"`
}
import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/kubeshark/kubeshark/misc"
	"github.com/rs/zerolog/log"
)

type ScriptingConfig struct {
	Env          map[string]interface{} `yaml:"env" json:"env" default:"{}"`
	Source       string                 `yaml:"source" json:"source" default:""`
	Sources      []string               `yaml:"sources" json:"sources" default:"[]"`
	WatchScripts bool                   `yaml:"watchScripts" json:"watchScripts" default:"true"`
	Active       []string               `yaml:"active" json:"active" default:"[]"`
	Console      bool                   `yaml:"console" json:"console" default:"true"`
}

func (config *ScriptingConfig) GetScripts() (scripts []*misc.Script, err error) {
	// Check if both Source and Sources are empty
	if config.Source == "" && len(config.Sources) == 0 {
		return nil, nil
	}

	var allFiles []struct {
		Source string
		File   fs.DirEntry
	}
package configStructs

// ScriptingConfig defines configuration options for the scripting subsystem
type ScriptingConfig struct {
	// Enabled determines if scripting functionality is available
	Enabled bool `yaml:"enabled" json:"enabled" default:"false"`
	
	// RetentionLimit sets the maximum duration (in seconds) that scripts can retain PCAP files
	// Setting to 0 means no limit
	RetentionLimit int `yaml:"retentionLimit" json:"retentionLimit" default:"3600"`
	
	// ScriptsDir is the directory where scripts are stored
	ScriptsDir string `yaml:"scriptsDir" json:"scriptsDir" default:"scripts"`
	
	// MaxConcurrency is the maximum number of scripts that can run concurrently
	MaxConcurrency int `yaml:"maxConcurrency" json:"maxConcurrency" default:"10"`
	
	// Timeout is the maximum execution time for a script in milliseconds
	Timeout int `yaml:"timeout" json:"timeout" default:"5000"`
}
	// Handle single Source directory
	if config.Source != "" {
		files, err := os.ReadDir(config.Source)
		if err != nil {
			return nil, fmt.Errorf("failed to read directory %s: %v", config.Source, err)
		}
		for _, file := range files {
			allFiles = append(allFiles, struct {
				Source string
				File   fs.DirEntry
			}{Source: config.Source, File: file})
		}
	}

	// Handle multiple Sources directories
	if len(config.Sources) > 0 {
		for _, source := range config.Sources {
			files, err := os.ReadDir(source)
			if err != nil {
				return nil, fmt.Errorf("failed to read directory %s: %v", source, err)
			}
			for _, file := range files {
				allFiles = append(allFiles, struct {
					Source string
					File   fs.DirEntry
				}{Source: source, File: file})
			}
		}
	}
package configStructs

// ScriptingConfig represents configuration for scripting capabilities
type ScriptingConfig struct {
	// Enabled indicates whether scripting is enabled
	Enabled bool `yaml:"enabled" json:"enabled" default:"true"`
	
	// Timeout is the maximum execution time for scripts in milliseconds
	Timeout int `yaml:"timeout" json:"timeout" default:"5000"`
	
	// MaxMemory is the maximum memory usage for scripts in megabytes
	MaxMemory int `yaml:"maxMemory" json:"maxMemory" default:"100"`
}
	// Iterate over all collected files
	for _, f := range allFiles {
		if f.File.IsDir() {
			continue
		}

		// Construct the full path based on the relevant source directory
		path := filepath.Join(f.Source, f.File.Name())
		if !strings.HasSuffix(f.File.Name(), ".js") { // Use file name suffix for skipping non-JS files
			log.Info().Str("path", path).Msg("Skipping non-JS file")
			continue
		}

		// Read the script file
		var script *misc.Script
		script, err = misc.ReadScriptFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read script file %s: %v", path, err)
		}
package configStructs

// ScriptingConfig defines settings for the scripting engine
type ScriptingConfig struct {
	Enabled       bool   `yaml:"enabled" default:"false"`
	TimeoutMs     int    `yaml:"timeoutMs" default:"5000"`
	DefaultScript string `yaml:"defaultScript" default:""`
}
		// Append the valid script to the scripts slice
		scripts = append(scripts, script)

		log.Debug().Str("path", path).Msg("Found script:")
	}

	// Return the collected scripts and nil error if successful
	return scripts, nil
}
package configStructs
package configStructs

// ScriptingConfig represents configuration for the scripting system
type ScriptingConfig struct {
	// Enabled indicates whether scripting is enabled
	Enabled bool `yaml:"enabled" json:"enabled" default:"true"`
	
	// Timeout is the maximum execution time for scripts in milliseconds
	Timeout int `yaml:"timeout" json:"timeout" default:"5000"`
	
	// MaxMemory is the maximum memory usage for scripts in megabytes
	MaxMemory int `yaml:"maxMemory" json:"maxMemory" default:"100"`
}

// ScriptingPermissions represents user permissions for scripting
type ScriptingPermissions struct {
	// CanSave indicates whether the user can save scripts
	CanSave bool `yaml:"canSave" json:"canSave" default:"false"`
	
	// CanActivate indicates whether the user can activate scripts
	CanActivate bool `yaml:"canActivate" json:"canActivate" default:"false"`
	
	// CanDelete indicates whether the user can delete scripts
	CanDelete bool `yaml:"canDelete" json:"canDelete" default:"false"`
}
// ScriptingConfig represents configuration for the scripting system
type ScriptingConfig struct {
	// Enabled indicates whether scripting is enabled
	Enabled bool `yaml:"enabled" json:"enabled" default:"true"`
	
	// Timeout is the maximum execution time for scripts in milliseconds
	Timeout int `yaml:"timeout" json:"timeout" default:"5000"`
	
	// MaxMemory is the maximum memory usage for scripts in megabytes
	MaxMemory int `yaml:"maxMemory" json:"maxMemory" default:"100"`
	
	// AllowFileSystem indicates whether scripts can access the file system
	AllowFileSystem bool `yaml:"allowFileSystem" json:"allowFileSystem" default:"false"`
	
	// AllowNetwork indicates whether scripts can access the network
	AllowNetwork bool `yaml:"allowNetwork" json:"allowNetwork" default:"false"`
}

// ScriptingPermissions represents user permissions for scripting
type ScriptingPermissions struct {
	// CanSave indicates whether the user can save scripts
	CanSave bool `yaml:"canSave" json:"canSave" default:"false"`
	
	// CanActivate indicates whether the user can activate scripts
	CanActivate bool `yaml:"canActivate" json:"canActivate" default:"false"`
	
	// CanDelete indicates whether the user can delete scripts
	CanDelete bool `yaml:"canDelete" json:"canDelete" default:"false"`
}