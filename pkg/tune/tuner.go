package tune

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os/exec"
	"strings"

	"github.com/elastic/go-elasticsearch/v5"
	"github.com/elastic/go-elasticsearch/v5/esapi"
	log "github.com/sirupsen/logrus"
	"github.com/vshn/waf-tool/cfg"
	"github.com/vshn/waf-tool/pkg/forward"
)

// Tune creates exclusion rules for a given uniqe ID
func Tune(uniqueID string, config cfg.Configuration) (returnError error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: config.ElasticSearch.InsecureSkipVerify,
		},
	}

	if len(config.ElasticSearch.CustomCAFile) > 0 {
		caCert, err := ioutil.ReadFile(config.ElasticSearch.CustomCAFile)
		if err != nil {
			return fmt.Errorf("Could not read file with CA cert: %w", err)
		}
		if err := addRootCA(caCert, transport.TLSClientConfig); err != nil {
			return err
		}
	}

	if len(config.ElasticSearch.CustomCA) > 0 {
		pem := []byte(config.ElasticSearch.CustomCA)
		if err := addRootCA(pem, transport.TLSClientConfig); err != nil {
			return err
		}
	}

	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{config.ElasticSearch.URL},
		Transport: transport,
	})

	if err != nil {
		return fmt.Errorf("Could not create ES client: %w", err)
	}

	url, err := url.Parse(config.ElasticSearch.URL)
	if err != nil {
		return err
	}

	forwarder := forward.NewPortForwarder(url.Port())

	defer func() {
		if err := forwarder.Stop(); err != nil {
			if returnError == nil {
				returnError = err
			} else {
				// Already encountered an error so we log this one here
				log.WithError(err).Error("Could not stop port forwarding")
			}
		}
	}()
	if err := forwarder.Start(); err != nil {
		return err
	}

	out, err := exec.Command("oc", "whoami", "--show-token").Output()
	if err != nil {
		return fmt.Errorf("Could not get token: %w", err)
	}

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

func addRootCA(customCA []byte, tlsConfig *tls.Config) error {
	if tlsConfig.RootCAs == nil {
		tlsConfig.RootCAs = x509.NewCertPool()
	}
	if ok := tlsConfig.RootCAs.AppendCertsFromPEM(customCA); !ok {
		return fmt.Errorf("Could not read CA cert: %s", customCA)
	}
	return nil
}
