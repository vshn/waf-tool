package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestTuneConfigFrom(t *testing.T) {
	url := "https://some-example.net/"
	os.Setenv("WAF_ES_URL", url)

	token := "token"
	os.Setenv("WAF_GITLAB_TOKEN", token)

	customCA := "PEMPEM"
	os.Setenv("WAF_ES_CUSTOM_CA", customCA)

	viper.Set("es-insecure-skip-tls-verify", true)

	caFile := "/some/path/to/cert.pem"
	viper.Set("es-custom-ca-file", caFile)

	initConfig()

	assert.Equal(t, url, config.ElasticSearch.URL)
	assert.Equal(t, customCA, config.ElasticSearch.CustomCA)
	assert.True(t, config.ElasticSearch.InsecureSkipVerify)
	assert.Equal(t, caFile, config.ElasticSearch.CustomCAFile)
	assert.Equal(t, token, config.GitLab.Token)
}

func TestTuneHelp(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"tune", "-h"})
	Execute()
	assert.Contains(t, buf.String(), "waf-tool tune [unique-id] [flags]")
}
