package cfg

import (
	log "github.com/sirupsen/logrus"
	"os"
)

type (
	ConfigMap struct {
		Log           LogConfig
		ElasticSearch ElasticSearchConfig
	}
	LogConfig struct {
		Level string
	}
	ElasticSearchConfig struct {
		Url string
	}
)

func CreateDefaultConfig() ConfigMap {
	return ConfigMap{
		Log: LogConfig{
			Level: "info",
		},
		ElasticSearch: ElasticSearchConfig{
			Url: "http://localhost:9200/",
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
