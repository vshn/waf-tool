package git

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/vshn/waf-tool/cfg"
	"github.com/vshn/waf-tool/pkg/file"
	"github.com/vshn/waf-tool/pkg/model"
	"github.com/xanzy/go-gitlab"
	"io"
	"io/ioutil"
	"net/url"
	"strings"
	"time"
)

const (
	author        = "waf-tuner"
	commitMessage = "Update exception rules"
	remote        = "origin"
	protocol      = "https"
)

type WorkingBranch struct {
	Worktree   *git.Worktree
	Head       *plumbing.Reference
	BranchName plumbing.ReferenceName
	RuleFile   file.RuleFile
}

func (wb *WorkingBranch) Create() error {
	err := wb.Worktree.Checkout(&git.CheckoutOptions{
		Branch: wb.BranchName,
		Create: true,
	})
	if err != nil {
		return fmt.Errorf("could not create branch: %w", err)
	}
	return nil
}

func (wb *WorkingBranch) Update(ruleFunc model.RuleIdFunc) error {
	_, err := wb.RuleFile.Open(wb.Worktree)
	if err != nil {
		return err
	}
	id, err := wb.RuleFile.GetLastId()
	if err != nil {
		return err
	}
	rule, err := ruleFunc(id)
	if err != nil {
		return err
	}
	err = wb.RuleFile.Write(rule)
	if err != nil {
		return err
	}
	err = wb.RuleFile.Close()
	if err != nil {
		return err
	}
	_, err = wb.Worktree.Add(wb.RuleFile.Path)
	if err != nil {
		return fmt.Errorf("could not add updated rules to git worktree: %w", err)
	}
	return nil
}

func (wb *WorkingBranch) Save() (plumbing.Hash, error) {
	commit, err := wb.Worktree.Commit(commitMessage, &git.CommitOptions{
		All: true,
		Author: &object.Signature{
			Name: author,
			When: time.Now(),
		},
	})
	if err != nil {
		return commit, fmt.Errorf("could not commit to worketree: %w", err)
	}
	return commit, nil
}

func CreateMergeRequest(config *cfg.Configuration, featureBranchName string) error {
	uri, err := url.ParseRequestURI(config.GitLab.Repository)
	client, err := gitlab.NewOAuthClient(config.GitLab.Token, gitlab.WithBaseURL(uri.Scheme+"://"+uri.Host))
	if err != nil {
		return fmt.Errorf("could not create new gitlab client: %w", err)
	}

	_, response, err := client.MergeRequests.CreateMergeRequest(uri.Path[1:len(uri.Path)], &gitlab.CreateMergeRequestOptions{
		Title:              gitlab.String("New exception rules"),
		Description:        gitlab.String("This request has been created automatically from waf-tuner"),
		SourceBranch:       gitlab.String(featureBranchName),
		TargetBranch:       gitlab.String("master"),
		RemoveSourceBranch: gitlab.Bool(true),
	})
	if err != nil {
		return fmt.Errorf("cannot create merge request: %w", err)
	}
	if response == nil {
		return fmt.Errorf("cannot create merge request, response is nil")
	}

	if response.Response.StatusCode != 201 {
		buffer := new(strings.Builder)
		_, err := io.Copy(buffer, response.Response.Body)
		if err != nil {
			return fmt.Errorf("could not create merge request, cannot read response body")
		}
		return fmt.Errorf("could not create merge request: %s", buffer.String())
	}

	return nil
}

func GetRepository(config *cfg.Configuration) (*git.Repository, error) {
	if strings.HasPrefix(config.GitLab.Repository, protocol) {
		pathDir, err := ioutil.TempDir("", "project")
		if err != nil {
			return nil, fmt.Errorf("could not create temporary clone folder: %w", err)
		}
		repository, err := git.PlainClone(pathDir, false, &git.CloneOptions{
			Auth: &http.BasicAuth{
				Username: "user", // yes, this can be anything except an empty string
				Password: config.GitLab.Token,
			},
			URL: config.GitLab.Repository,
		})
		if err != nil {
			return nil, fmt.Errorf("could not clone repository from URL: %w", err)
		}
		return repository, nil
	} else {
		repository, err := git.PlainOpen(config.GitLab.Repository)
		if err != nil {
			return nil, fmt.Errorf("could not open repository from path: %w", err)
		}
		gitRemote, err := repository.Remote(remote)
		if err != nil {
			return nil, fmt.Errorf("could not get remote url from repository: %w", err)
		}
		remoteUrl := gitRemote.Config().URLs[0]
		if !strings.HasPrefix(remoteUrl, protocol) {
			return nil, fmt.Errorf("invalid %s remote repository: %s", protocol, remoteUrl)
		}
		config.GitLab.Repository = strings.Replace(remoteUrl, ".git", "", 1)
		return repository, nil
	}
}

func SaveRepository(repository *git.Repository, commit plumbing.Hash, token string) error {
	_, err := repository.CommitObject(commit)
	if err != nil {
		return fmt.Errorf("could not commit updated rules: %w", err)
	}

	err = repository.Push(&git.PushOptions{
		RemoteName: remote,
		Auth: &http.BasicAuth{
			Username: "user", // yes, this can be anything except an empty string
			Password: token,
		},
		Force: true,
	})
	if err != nil {
		return fmt.Errorf("could not push to remote repository: %w", err)
	}
	return nil
}
