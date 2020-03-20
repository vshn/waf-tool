package tuner

import (
	"fmt"
	"net/url"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/vshn/waf-tool/cfg"
	"github.com/vshn/waf-tool/pkg/elasticsearch"
	"github.com/vshn/waf-tool/pkg/forwarder"
	"github.com/vshn/waf-tool/pkg/model"
	"github.com/vshn/waf-tool/pkg/rules"
)

const (
	baseID   = 10100
	ocBinary = "oc"
)

// Tune creates exclusion rules for a given uniqe ID
func Tune(uniqueID string, config cfg.Configuration) (returnError error) {
	es, fwd, err := prepareEsClient(config)
	if err != nil {
		return err
	}
	defer func() {
		if err := fwd.Stop(); err != nil {
			if returnError == nil {
				returnError = err
			} else {
				// Already encountered an error so we log this one here
				log.WithError(err).Error("Could not stop port forwarding")
			}
		}
	}()

	result, err := es.SearchUniqueID(uniqueID)
	if err != nil {
		return err
	}

	if len(result.Hits.Hits) == 0 {
		log.WithField("unique-id", uniqueID).Warnf("Could not find any request")
		return nil
	}

	var alerts []model.ModsecAlert
	var access *model.ApacheAccess
	for _, result := range result.Hits.Hits {
		if result.Source.ApacheAccess != nil {
			if access != nil {
				return fmt.Errorf("found multiple access logs for same unique id: %s", uniqueID)
			}
			log.WithField("access", result.Source.ApacheAccess).Debug("Found apache access log")
			access = result.Source.ApacheAccess
		}
		if result.Source.ModsecAlert != nil {
			log.WithField("alert", result.Source.ModsecAlert).Debug("Found ModSec alert")
			alerts = append(alerts, *result.Source.ModsecAlert)
		}
	}
	if access == nil {
		return fmt.Errorf("could not find any access log for unique id: %s", uniqueID)
	}
	if len(alerts) == 0 {
		log.WithField("unique-id", uniqueID).Warnf("Could not find any ModSecurity alerts for request")
		return nil
	}
	rule, err := rules.CreateByIDExclusion(alerts, baseID)
	if err != nil {
		return err
	}
	fmt.Print(rule)

	return nil
}

func prepareEsClient(config cfg.Configuration) (elasticsearch.Client, forwarder.PortForwarder, error) {
	out, err := exec.Command(ocBinary, "whoami", "--show-token").Output()
	if err != nil {
		return nil, nil, fmt.Errorf("could not get token: %w", err)
	}
	token := strings.TrimSpace(string(out))
	es, err := elasticsearch.New(config.ElasticSearch, token)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create ES client: %w", err)
	}

	url, err := url.Parse(config.ElasticSearch.URL)
	if err != nil {
		return nil, nil, err
	}

	fwd := forwarder.New(url.Port())

	if err := fwd.Start(); err != nil {
		return nil, nil, err
	}
	return es, fwd, nil
}
