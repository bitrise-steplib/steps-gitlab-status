package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-steputils/stepconf"
)

// Config ...
type Config struct {
	APIURL        string `env:"api_base_url,required"`
	PrivateToken  string `env:"private_token,required"`
	RepositoryURL string `env:"repository_url,required"`
	CommitHash    string `env:"commit_hash,required"`
	PresetStatus  string `env:"preset_status,opt[auto,pending,running,success,failed,canceled]"`
}

type projects []struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func getID(apiURL, token, repo string) (n int, err error) {
	client := &http.Client{}
	var req *http.Request
	url := fmt.Sprintf("%s/projects?simple=true&membership=true&search=%s", apiURL, repo)
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Add("PRIVATE-TOKEN", token)

	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to send the request: %s", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); err == nil {
			err = cerr
		}
	}()

	var ps projects
	if err := json.NewDecoder(resp.Body).Decode(&ps); err != nil {
		return 0, err
	}

	for _, p := range ps {
		if p.Name == repo {
			return p.ID, nil
		}
	}
	return 0, fmt.Errorf("id not found for repository (%s)", repo)
}

func getStatus(preset string) string {
	if preset != "auto" {
		return preset
	}
	if os.Getenv("BITRISE_BUILD_STATUS") == "0" {
		return "success"
	}
	return "failed"
}

// https://docs.gitlab.com/ce/api/commits.html#post-the-build-status-to-a-commit
func sendStatus(apiURL, token, commit, state string, id int) (err error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/projects/%d/statuses/%s?state=%s", apiURL, id, commit, state)
	var req *http.Request
	req, err = http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("PRIVATE-TOKEN", token)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send the request: %s", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); err == nil {
			err = cerr
		}
	}()

	return err
}

func main() {
	if os.Getenv("commit_hash") == "" {
		log.Warnf("GitLab requires a commit hash for build status reporting")
		os.Exit(1)
	}

	var cfg Config
	if err := stepconf.Parse(&cfg); err != nil {
		log.Errorf("Error: %s\n", err)
		os.Exit(1)
	}
	stepconf.Print(cfg)

	lastSlash := strings.LastIndex(cfg.RepositoryURL, "/")
	lastDot := strings.LastIndex(cfg.RepositoryURL, ".")
	repoName := cfg.RepositoryURL[lastSlash+1 : lastDot]

	id, err := getID(cfg.APIURL, cfg.PrivateToken, repoName)
	if err != nil {
		log.Errorf("Error: %s\n", err)
		os.Exit(1)
	}
	if err := sendStatus(cfg.APIURL, cfg.PrivateToken, cfg.CommitHash, getStatus(cfg.PresetStatus), id); err != nil {
		log.Errorf("Error: %s\n", err)
		os.Exit(1)
	}
}
