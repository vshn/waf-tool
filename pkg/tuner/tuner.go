package tuner

import (
	"fmt"
	"github.com/vshn/waf-tool/pkg/file"
	"github.com/vshn/waf-tool/pkg/rules"
	"net/url"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/vshn/waf-tool/cfg"
	"github.com/vshn/waf-tool/pkg/elasticsearch"
	"github.com/vshn/waf-tool/pkg/forwarder"
	"github.com/vshn/waf-tool/pkg/git"
	"github.com/vshn/waf-tool/pkg/model"
)

const (
	baseID   = 10100
	ocBinary = "oc"
)
const (
	ruleFilePath      = "deployment/base/rules/before-crs/rules.conf"
	ref               = "refs/heads/"
	featureBranchName = "update-waf-rules"
)

// Tune creates exclusion rules for a given uniqe ID
func Tune(uniqueID string, configuration cfg.Configuration) (returnError error) {
	es, fwd, err := prepareEsClient(configuration)
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
	ruleFunc := func(id int) (string, error) { return rules.CreateByIDExclusion(alerts, id) }
	if configuration.GitLab.MergeRequest {
		log.WithField("unique-id", uniqueID).Info("Prepare creating new merge request")
		err = manageGit(configuration, ruleFunc)
		if err != nil {
			return err
		}
	} else {
		fmt.Print(ruleFunc(baseID))
	}
	return nil
}

func manageGit(configuration cfg.Configuration, ruleFunc model.RuleIdFunc) error {
	log.WithField("repository", configuration.GitLab.Repository).Info("Set GitLab repository")
	repository, err := git.GetRepository(&configuration)
	if err != nil {
		return err
	}
	worktree, err := repository.Worktree()
	if err != nil {
		return fmt.Errorf("could not get worktree: %w", err)
	}
	head, err := repository.Head()
	if err != nil {
		return fmt.Errorf("could not get head reference from repository: %w", err)
	}

	workingBranch := git.WorkingBranch{
		Worktree:   worktree,
		Head:       head,
		BranchName: ref + featureBranchName,
		RuleFile:   file.RuleFile{Path: ruleFilePath},
	}
	log.WithField("branch", featureBranchName).Info("Manage new branch")
	err = workingBranch.Create()
	if err != nil {
		return err
	}
	err = workingBranch.Update(ruleFunc)
	if err != nil {
		return err
	}
	commit, err := workingBranch.Save()
	if err != nil {
		return err
	}
	log.WithField("branch", featureBranchName).Info("Save new branch")
	err = git.SaveRepository(repository, commit, configuration.GitLab.Token)
	if err != nil {
		return err
	}
	log.Info("Create merge request")
	err = git.CreateMergeRequest(&configuration, featureBranchName)
	if err != nil {
		return err
	}
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
