
// LogsConfig represents configuration for the logging system
type LogsConfig struct {
	// Level sets the minimum log level to output
	Level string `yaml:"level" json:"level" default:"info"`
	
	// File is the path to the log file
	File string `yaml:"file" json:"file"`
	
	// MaxSize is the maximum size in megabytes of the log file before it gets rotated
	MaxSize int `yaml:"maxSize" json:"maxSize" default:"100"`
	
	// MaxBackups is the maximum number of old log files to retain
	MaxBackups int `yaml:"maxBackups" json:"maxBackups" default:"5"`
	
	// MaxAge is the maximum number of days to retain old log files
	MaxAge int `yaml:"maxAge" json:"maxAge" default:"28"`
	
	// Compress determines if the rotated log files should be compressed
	Compress bool `yaml:"compress" json:"compress" default:"true"`
}
// LogsConfig defines configuration options for the logs subsystem
type LogsConfig struct {
	// Namespaces is the list of namespaces from which to collect logs
	Namespaces []string `yaml:"namespaces" json:"namespaces" default:"[]"`
	
	// PodRegexStr is a regex pattern to match pods for log collection
	PodRegexStr string `yaml:"regex" json:"regex" default:".*"`
	
	// Since is the starting point for logs collection (e.g., "10m", "1h")
	Since string `yaml:"since" json:"since" default:"1h"`
	
	// Follow determines if log streaming should continue
	Follow bool `yaml:"follow" json:"follow" default:"false"`
	
	// Tail specifies the number of lines to show from the end of logs
	Tail int `yaml:"tail" json:"tail" default:"10"`
	
	// Timestamps indicates whether to include timestamps in logs output
	Timestamps bool `yaml:"timestamps" json:"timestamps" default:"true"`
}
import (
	"fmt"
	"os"
	"path"
	"github.com/kubeshark/kubeshark/misc"
)

const (
	FileLogsName = "file"
	GrepLogsName = "grep"
)

type LogsConfig struct {
	FileStr string `yaml:"file" json:"file"`
	Grep    string `yaml:"grep" json:"grep"`
}

func (config *LogsConfig) Validate() error {
	if config.FileStr == "" {
		_, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get PWD, %v (try using `%s logs -f <full path dest zip file>)`", err, misc.Program)
		}
	}
	return nil
}

func (config *LogsConfig) FilePath() string {
	if config.FileStr == "" {
		pwd, _ := os.Getwd()
		return path.Join(pwd, fmt.Sprintf("%s_logs.zip", misc.Program))
	}

	return config.FileStr
}

// LogsConfig defines logging configuration
type LogsConfig struct {
	Console bool   `yaml:"console" default:"true"`
	File    string `yaml:"file" default:""`
}