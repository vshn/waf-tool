package cfg

type (
	// Configuration combines all configs
	Configuration struct {
		Log           LogConfig           `mapstructure:",squash"`
		ElasticSearch ElasticSearchConfig `mapstructure:",squash"`
		GitLab        GitLabConfig        `mapstructure:",squash"`
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

	// GitLabConfig configures GitLab
	GitLabConfig struct {
		Repository   string `mapstructure:"repository"`
		MergeRequest bool   `mapstructure:"create-merge-request"`
		Token        string `mapstructure:"gitlab-token"`
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
		GitLab: GitLabConfig{
			Repository:   "./",
			MergeRequest: false,
		},
	}
}
