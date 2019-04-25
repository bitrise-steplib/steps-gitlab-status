package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-steputils/stepconf"
)

type config struct {
	PrivateToken  string `env:"private_token,required"`
	RepositoryURL string `env:"repository_url,required"`
	GitBranch     string `env:"git_branch,required"`
	CommitHash    string `env:"commit_hash,required"`
	APIURL        string `env:"api_base_url,required"`

	Status      string `env:"preset_status,opt[auto,pending,running,success,failed,canceled]"`
	TargetURL   string `env:"target_url"`
	Context     string `env:"context"`
	Description string `env:"description"`
}

// getRepo parses the repository from a url. Possible url formats:
// - https://hostname/owner/repository.git
// - git@hostname:owner/repository.git
func getRepo(url string) string {
	url = strings.TrimPrefix(strings.TrimPrefix(url, "https://"), "git@")
	return url[strings.IndexAny(url, ":/")+1 : strings.Index(url, ".git")]
}

func getState(preset string) string {
	if preset != "auto" {
		return preset
	}
	if os.Getenv("BITRISE_BUILD_STATUS") == "0" {
		return "success"
	}
	return "failed"
}

func getDescription(desc, state string) string {
	if desc == "" {
		strings.Title(getState(state))
	}
	return desc
}

// sendStatus creates a commit status for the given commit.
// see also: https://docs.gitlab.com/ce/api/commits.html#post-the-build-status-to-a-commit
func sendStatus(cfg config) error {
	repo := url.PathEscape(getRepo(cfg.RepositoryURL))
	form := url.Values{
		"state":       {getState(cfg.Status)},
		"ref":         {cfg.GitBranch},
		"target_url":  {cfg.TargetURL},
		"description": {getDescription(cfg.Description, cfg.Status)},
		"context":     {cfg.Context},
	}

	url := fmt.Sprintf("%s/projects/%s/statuses/%s", cfg.APIURL, repo, cfg.CommitHash)
	req, err := http.NewRequest("POST", url, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("PRIVATE-TOKEN", cfg.PrivateToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send the request: %s", err)
	}
	if err := resp.Body.Close(); err != nil {
		return err
	}
	if 200 > resp.StatusCode || resp.StatusCode >= 300 {
		return fmt.Errorf("server error: %s", resp.Status)
	}

	return err
}

func main() {
	if os.Getenv("commit_hash") == "" {
		log.Warnf("GitLab requires a commit hash for build status reporting")
		os.Exit(1)
	}

	var cfg config
	if err := stepconf.Parse(&cfg); err != nil {
		log.Errorf("Error: %s\n", err)
		os.Exit(1)
	}
	stepconf.Print(cfg)

	if err := sendStatus(cfg); err != nil {
		log.Errorf("Error: %s\n", err)
		os.Exit(1)
	}
}
