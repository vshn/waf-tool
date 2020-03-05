package cfg

import (
	"os"

	log "github.com/sirupsen/logrus"
)

type (
	// ConfigMap combines all configs
	ConfigMap struct {
		Log           LogConfig
		ElasticSearch ElasticSearchConfig
	}
	// LogConfig configures the log level
	LogConfig struct {
		Level string
	}
	// ElasticSearchConfig configures ES
	ElasticSearchConfig struct {
		URL string
	}
)

// CreateDefaultConfig creates a default config
func CreateDefaultConfig() ConfigMap {
	return ConfigMap{
		Log: LogConfig{
			Level: "info",
		},
		ElasticSearch: ElasticSearchConfig{
			URL: "https://localhost:9200/",
		},
	}
}

// SetupLogging initializes logging framework
func SetupLogging(cfg LogConfig) {

	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	level, err := log.ParseLevel(cfg.Level)
	if err != nil {
		log.WithError(err).Warn("Using info level.")
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(level)
	}
}
