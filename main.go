package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-steputils/stepconf"
)

type config struct {
	PrivateToken  string `env:"private_token,required"`
	RepositoryURL string `env:"repository_url,required"`
	GitRef        string `env:"git_branch"`
	CommitHash    string `env:"commit_hash,required"`
	APIURL        string `env:"api_base_url,required"`

	Status      string  `env:"preset_status,opt[auto,pending,running,success,failed,canceled]"`
	TargetURL   string  `env:"target_url"`
	Context     string  `env:"context"`
	Description string  `env:"description"`
	Coverage    float64 `env:"coverage,range[0.0..100.0]"`
}

// getRepo parses the repository from a url
func getRepo(u string) string {
	r := regexp.MustCompile(`.*[:/](.+?\/.+?)(?:\.git|$|\/)`)
	if matches := r.FindStringSubmatch(u); len(matches) == 2 {
		return matches[1]
	}
	return ""
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
		return strings.Title(getState(state))
	}
	return desc
}

// sendStatus creates a commit status for the given commit.
// see also: https://docs.gitlab.com/ce/api/commits.html#post-the-build-status-to-a-commit
func sendStatus(cfg config) error {
	repo := url.PathEscape(getRepo(cfg.RepositoryURL))
	form := url.Values{
		"state":       {getState(cfg.Status)},
		"target_url":  {cfg.TargetURL},
		"description": {getDescription(cfg.Description, cfg.Status)},
		"context":     {cfg.Context},
		"coverage":    {fmt.Sprintf("%f", cfg.Coverage)},
	}

	if strings.TrimSpace(cfg.GitRef) != "" {
		form["ref"] = []string{strings.TrimSpace(cfg.GitRef)}
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := resp.Body.Close(); err != nil {
		return err
	}
	if 200 > resp.StatusCode || resp.StatusCode >= 300 {
		return fmt.Errorf("server error: %s url: %s code: %d body: %s", resp.Status, url, resp.StatusCode, string(body))
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
