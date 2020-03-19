package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootConfig(t *testing.T) {
	os.Setenv("WAF_VERBOSE", "true")

	initConfig()

	assert.True(t, config.Log.Verbose)
}

func TestRootHelp(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{""})
	Execute()
	assert.Contains(t, buf.String(), "Usage:")
}

func TestVersionCmd(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"--version"})
	version := "test version"
	SetVersion(version)
	Execute()
	assert.Contains(t, buf.String(), "waf-tool version "+version)
}
