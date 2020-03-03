package cmd

import (
	"github.com/elastic/go-elasticsearch/v8"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/vshn/waf-tool/cfg"
	"gopkg.in/src-d/go-git.v4"
	"os"
)

func CreateDemoCommand() *cobra.Command {
	c := cfg.CreateDefaultConfig()
	cmd := &cobra.Command{
		Use: "demo",
		Run: func(cmd *cobra.Command, args []string) {
			GetConfigFromFlags(cmd.PersistentFlags(), &c)
			RunDemoCommand(c.ElasticSearch)
		},
	}
	cmd.PersistentFlags().String("elasticsearch.url", c.ElasticSearch.Url, "Elasticsearch target URL")
	return cmd
}

func RunDemoCommand(c cfg.ElasticSearchConfig) {

	log.WithFields(log.Fields{
		"version": elasticsearch.Version,
		"url":     c.Url,
	}).Debug("Connecting to Elasticsearch...")
	os.Setenv("ELASTICSEARCH_URL", c.Url)
	es, _ := elasticsearch.NewDefaultClient()
	log.Info(es.Info())
	url := "https://github.com/src-d/go-git"
	log.WithField("url", url).Info("Cloning Git repo...")
	if _, err := git.PlainClone("/tmp/foo", false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	}); err != nil {
		log.WithError(err).Error("Could not clone git repo.")
	}
}
