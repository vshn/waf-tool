package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vshn/waf-tool/pkg/tuner"
)

var (
	tuneCmd = &cobra.Command{
		Use:   "tune [unique-id]",
		Short: "Create ModSecurity rule exclusions for a given request unique ID",
		Long: `The tool will use the oc binary to start a port forward to the cluster's Elasticsearch.
Using the $KUBECONFIG token it will query ES for the given unique ID.`,
		Args:      cobra.ExactArgs(1),
		ValidArgs: []string{"unique-id"},
		RunE:      runTuneCommand,
	}
)

func init() {
	rootCmd.AddCommand(tuneCmd)
	tuneCmd.Flags().StringP("es-url", "u", config.ElasticSearch.URL, "Elasticsearch target URL")
	tuneCmd.Flags().BoolP("es-insecure-skip-tls-verify", "k", config.ElasticSearch.InsecureSkipVerify, "If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure")
	tuneCmd.Flags().String("es-custom-ca", config.ElasticSearch.CustomCA, "Custom CA certificate to trust (in PEM format)")
	tuneCmd.Flags().String("es-custom-ca-file", config.ElasticSearch.CustomCAFile, "Path to custom CA certificate to trust (in PEM format)")
	tuneCmd.Flags().StringP("repository", "r", config.GitLab.Repository, "Repository Local file system path or a URL")
	tuneCmd.Flags().StringP("gitlab-token", "t", config.GitLab.Token, "GitLab access token")
	tuneCmd.Flags().BoolP("create-merge-request", "m", config.GitLab.MergeRequest, "Create a merge request on GitLab")
	if err := viper.BindPFlags(tuneCmd.Flags()); err != nil {
		log.WithError(err).Fatal()
	}
}

// RunTuneCommand runs the tune command
func runTuneCommand(cmd *cobra.Command, args []string) error {
	return tuner.Tune(args[0], config)
}
