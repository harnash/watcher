package logging

import (
	"log"
	"os"

	"github.com/rs/xlog"
)

// Logger instance
var logger xlog.Logger
var config *xlog.Config

// newLogger creates global Logger instance
func newLogger() xlog.Logger {
	return xlog.New(*config)
}

// GetLogger return the global Logger instance
func GetLogger() xlog.Logger {
	if logger == nil {
		if config == nil {
			log.Fatal("No logger initialized! Call InitConfig() first.")
			return nil
		}

		logger = newLogger()
	}

	return logger
}

// GetConfig return current logger configuration
func GetConfig() xlog.Config {
	return *config
}

// InitConfig provides the configuration for the logging subsytem
func InitConfig() {
	host, _ := os.Hostname()
	// logstashWriter, err := net.Dial("udp", "127.0.0.1:1410")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	newConfig := xlog.Config{
		// Log info level and higher
		Level: xlog.LevelInfo,
		// Set some global env fields
		Fields: xlog.F{
			"role": "watcher",
			"host": host,
		},
		// Setup output
		Output: xlog.NewOutputChannel(xlog.MultiOutput{
			0: xlog.NewConsoleOutput(),
			// 1: xlog.NewLogstashOutput(logstashWriter),
		}),
	}

	config = &newConfig
}
