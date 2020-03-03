package cmd

import (
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/vshn/waf-tool/cfg"
	"os"
	"testing"
)

func TestGetConfigFromFlags(t *testing.T) {
	type args struct {
		flags  *pflag.FlagSet
		cfg    cfg.ConfigMap
		getter func(configMap cfg.ConfigMap) interface{}
	}
	tests := []struct {
		name     string
		args     args
		key      string
		expected string
	}{
		{
			name: "ShouldParseEnvironmentVariable",
			args: args{
				flags: CreateDemoCommand().PersistentFlags(),
				cfg:   cfg.CreateDefaultConfig(),
				getter: func(configMap cfg.ConfigMap) interface{} {
					return configMap.ElasticSearch.Url
				},
			},
			key:      "ELASTICSEARCH_URL",
			expected: "http://elk:9200/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv(tt.key, tt.expected)
			GetConfigFromFlags(tt.args.flags, &tt.args.cfg)
			assert.Equal(t, tt.args.getter(tt.args.cfg), tt.expected)
		})
	}
}
