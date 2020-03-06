package cmd

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestTuneConfigFrom(t *testing.T) {
	url := "https://some-example.net/"
	os.Setenv("WAF_ES_URL", url)

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
}
