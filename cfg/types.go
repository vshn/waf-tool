package cfg

type (
	// Configuration combines all configs
	Configuration struct {
		Log           LogConfig           `mapstructure:",squash"`
		ElasticSearch ElasticSearchConfig `mapstructure:",squash"`
	}
	// LogConfig configures the log level
	LogConfig struct {
		Verbose bool
	}
	// ElasticSearchConfig configures ES
	ElasticSearchConfig struct {
		URL                string `mapstructure:"es-url"`
		InsecureSkipVerify bool   `mapstructure:"es-insecure-skip-tls-verify"`
		CustomCA           string `mapstructure:"es-custom-ca"`
		CustomCAFile       string `mapstructure:"es-custom-ca-file"`
	}
)

// NewDefaultConfig creates a default configuration
func NewDefaultConfig() Configuration {
	return Configuration{
		Log: LogConfig{
			Verbose: false,
		},
		ElasticSearch: ElasticSearchConfig{
			URL:                "https://localhost:9200/",
			InsecureSkipVerify: false,
		},
	}
}
