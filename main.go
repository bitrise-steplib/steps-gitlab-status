package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/bitrise-io/go-utils/log"
	"github.com/kdobmayer/input"
)

// Config ...
type Config struct {
	APIURL        string `env:"api_base_url" validate:"required"`
	PrivateToken  string `env:"private_token" validate:"required"`
	RepositoryURL string `env:"repository_url" validate:"required"`
	CommitHash    string `env:"commit_hash" validate:"required"`
	PresetStatus  string `env:"preset_status"`
}

type projects []struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func getID(apiURL, token, repo string) (n int, err error) {
	client := &http.Client{}
	var req *http.Request
	req, err = http.NewRequest("GET", apiURL+"/projects?simple=true", nil)
	if err != nil {
		return 0, err
	}
	req.Header.Add("PRIVATE-TOKEN", token)

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
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
	return 0, fmt.Errorf("id not found for repository %q", repo)
}

func getStatus(preset string) string {
	if preset != "" {
		return preset
	}
	if os.Getenv("STEPLIB_BUILD_STATUS") == "0" {
		return "success"
	}
	return "failed"
}

func sendStatus(apiURL, token, commit, state string, id int) (err error) {
	// https://docs.gitlab.com/ce/api/commits.html#post-the-build-status-to-a-commit
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
		return err
	}
	defer func() {
		if cerr := resp.Body.Close(); err == nil {
			err = cerr
		}
	}()
	return err
}

func main() {
	var c Config
	if err := input.New(&c); err != nil {
		log.Errorf("Couldn't create config: %v\n", err)
		os.Exit(1)
	}
	input.Print(c)

	lastSlash, lastDot := strings.LastIndex(c.RepositoryURL, "/"), strings.LastIndex(c.RepositoryURL, ".")
	repoName := c.RepositoryURL[lastSlash+1 : lastDot]
	id, err := getID(c.APIURL, c.PrivateToken, repoName)
	if err != nil {
		log.Errorf("error: %s", err)
		os.Exit(1)
	}
	if err := sendStatus(c.APIURL, c.PrivateToken, c.CommitHash, getStatus(c.PresetStatus), id); err != nil {
		log.Errorf("error: %s", err)
		os.Exit(1)
	}
}
