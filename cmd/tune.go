package cmd

import (
	"crypto/tls"
	"net/http"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/elastic/go-elasticsearch/v5"
	"github.com/elastic/go-elasticsearch/v5/esapi"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	tuneCmd = &cobra.Command{
		Use:   "tune",
		Short: "Create ModSecurity rule exclusions for a given request unique ID",
		Long: `The tool will use the oc binary to start a port forward to the cluster's Elasticsearch.
Using the $KUBECONFIG token it will query ES for the given unique ID.`,
		RunE: runTuneCommand,
	}
)

func init() {
	rootCmd.AddCommand(tuneCmd)
	tuneCmd.Flags().StringP("es-url", "u", config.ElasticSearch.URL, "Elasticsearch target URL")
	tuneCmd.Flags().BoolP("es-insecure-skip-tls-verify", "k", config.ElasticSearch.InsecureSkipVerify, "If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure")
	tuneCmd.Flags().String("es-custom-ca", config.ElasticSearch.CustomCA, "Custom CA certificate to trust (in PEM format)")
	tuneCmd.Flags().String("es-custom-ca-file", config.ElasticSearch.CustomCAFile, "Path to custom CA certificate to trust (in PEM format)")
	if err := viper.BindPFlags(tuneCmd.Flags()); err != nil {
		log.WithError(err).Fatal()
	}
}

// RunTuneCommand runs the tune command
func runTuneCommand(cmd *cobra.Command, args []string) error {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{config.ElasticSearch.URL},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: config.ElasticSearch.InsecureSkipVerify,
			},
		},
	})
	if err != nil {
		return err
	}

	out, err := exec.Command("oc", "whoami", "--show-token").Output()
	if err != nil {
		return err
	}
	port := "9200"
	portForward := exec.Command("oc", "port-forward", "-n", "logging", "svc/logging-es", port)
	defer func() {
		if err := portForward.Process.Signal(syscall.SIGTERM); err != nil {
			log.WithError(err).Error()
		}
		portForward.Wait()
	}()
	log.WithField("port", port).Info("Starting port forward...")
	err = portForward.Start()
	if err != nil {
		return err
	}

	time.Sleep(1 * time.Second)

	log.WithFields(log.Fields{
		"client_version": elasticsearch.Version,
		"url":            config.ElasticSearch.URL,
	}).Debug("Connecting to Elasticsearch...")

	res, err := es.Info(func(req *esapi.InfoRequest) {
		if req.Header == nil {
			req.Header = http.Header{}
		}
		token := strings.TrimSpace(string(out))
		req.Header.Add("Authorization", "Bearer "+token)
	})
	if err != nil {
		return err
	}
	log.Info(res)

	return nil
}
