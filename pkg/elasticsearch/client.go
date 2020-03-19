package elasticsearch

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/elastic/go-elasticsearch/v5"
	"github.com/vshn/waf-tool/cfg"
	"github.com/vshn/waf-tool/pkg/model"
)

const (
	indexWildcard = "project.*"
)

// Client is used to communicate to Elasticsearch
type Client interface {
	SearchUniqueID(uniqueID string) (model.SearchResult, error)
}

type client struct {
	es    *elasticsearch.Client
	token string
}

// New creates a new Elasticsearch client
func New(config cfg.ElasticSearchConfig, token string) (Client, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: config.InsecureSkipVerify,
		},
	}

	if len(config.CustomCAFile) > 0 {
		caCert, err := ioutil.ReadFile(config.CustomCAFile)
		if err != nil {
			return nil, fmt.Errorf("could not read file with CA cert: %w", err)
		}
		if err := addRootCA(caCert, transport.TLSClientConfig); err != nil {
			return nil, err
		}
	}

	if len(config.CustomCA) > 0 {
		pem := []byte(config.CustomCA)
		if err := addRootCA(pem, transport.TLSClientConfig); err != nil {
			return nil, err
		}
	}

	cl, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{config.URL},
		Transport: transport,
	})

	return &client{
		es:    cl,
		token: token,
	}, err
}

func addRootCA(customCA []byte, tlsConfig *tls.Config) error {
	if tlsConfig.RootCAs == nil {
		tlsConfig.RootCAs = x509.NewCertPool()
	}
	if ok := tlsConfig.RootCAs.AppendCertsFromPEM(customCA); !ok {
		return fmt.Errorf("could not read CA cert: %s", customCA)
	}
	return nil
}
