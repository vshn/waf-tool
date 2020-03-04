package main

import (
	"bytes"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_main_version(t *testing.T) {
	type field struct {
		arguments []string
		expected  string
	}
	tests := []struct {
		name   string
		fields field
	}{
		{
			name: "Main_ShouldParseVersionFlag",
			fields: field{
				arguments: []string{
					"--version",
				},
				expected: "waf-tool version latest, commit dirty, date today",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cobra.OnInitialize(initialize)
			buf := new(bytes.Buffer)
			rootCmd.SetOut(buf)
			rootCmd.SetArgs(tt.fields.arguments)
			main()
			assert.Contains(t, buf.String(), tt.fields.expected)
		})
	}
}
