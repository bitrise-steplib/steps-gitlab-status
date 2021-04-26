package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getRepo(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		want    string
		wantErr bool
	}{
		{"https", "https://github.com/owner/repository.git", "owner/repository", false},
		{"https with subgroup", "https://github.com/group/subgroup/subsubgroup/repository.git", "group/subgroup/subsubgroup/repository", false},
		{"https with port", "https://github.com:443/owner/repository.git", "owner/repository", false},
		{"https slash at end", "https://github.com/owner/repository.git/", "owner/repository", false},
		{"https no .git", "https://github.com/owner/repository", "owner/repository", false},
		{"https no .git with slash at end", "https://github.com/owner/repository/", "owner/repository", false},
		{"https custom domain", "https://gitlab.custom.com/owner/repository.git", "owner/repository", false},
		{"https with basic auth", "https://username:token@github.com/owner/repository.git", "owner/repository", false},
		{"https with port and basic auth", "https://username:token@github.com:443/owner/repository.git", "owner/repository", false},
		{"ssh without scheme", "user@github.com:owner/repository.git", "owner/repository", false},
		{"ssh without scheme, no user", "github.com:owner/repository.git", "owner/repository", false},
		{"ssh without scheme with subgroups", "user@gitlab.com:group/subgroup/subsubgroup/repository.git", "group/subgroup/subsubgroup/repository", false},
		{"ssh with scheme without port", "ssh://user@github.com/owner/repository.git", "owner/repository", false},
		{"ssh with scheme without port, no user", "ssh://github.com/owner/repository.git", "owner/repository", false},
		{"ssh with scheme with port", "ssh://user@github.com:22/owner/repository.git", "owner/repository", false},
		{"ssh with scheme with mutilple path", "ssh://gitlab.company.com:category/project-name/subproject/repository.git", "project-name/subproject/repository", false},
		{"ssh with scheme with subgroups", "ssh://user@gitlab.com/group/subgroup/subsubgroup/repository.git", "group/subgroup/subsubgroup/repository", false},
		{"ssh with scheme, invalid repo", "ssh://user@gitlab.com/repository.git", "", false},
		{"ssh without scheme, invalid repo", "user@gitlab.com:repository.git", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getRepo(tt.url); got != tt.want {
				t.Errorf("getRepo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseConfig(t *testing.T) {
	testCases := []struct {
		desc             string
		envVars          map[string]string
		expectedError    error
		expectedCoverage float64
	}{
		{
			desc: "when config can be parsed",
			envVars: map[string]string{
				"coverage":       "0.1",
				"private_token":  "asd123",
				"repository_url": "http://repo.url",
				"commit_hash":    "aaa111",
				"api_base_url":   "http://api.baseurl",
				"preset_status":  "success",
			},
			expectedError:    nil,
			expectedCoverage: 0.1,
		},
		{
			desc: "when coverage status field has an extra dot character, it removes the extra and succeed",
			envVars: map[string]string{
				"coverage":       "0.1.",
				"private_token":  "asd123",
				"repository_url": "http://repo.url",
				"commit_hash":    "aaa111",
				"api_base_url":   "http://api.baseurl",
				"preset_status":  "success",
			},
			expectedError:    nil,
			expectedCoverage: 0.1,
		},
		{
			desc: "when coverage status field has an extra whitespace character, it removes the extra and succeed",
			envVars: map[string]string{
				"coverage":       "0.1      ",
				"private_token":  "asd123",
				"repository_url": "http://repo.url",
				"commit_hash":    "aaa111",
				"api_base_url":   "http://api.baseurl",
				"preset_status":  "success",
			},
			expectedError:    nil,
			expectedCoverage: 0.1,
		},
		{
			desc: "when coverage status field has extra characters, it removes the extra and succeed",
			envVars: map[string]string{
				"coverage":       "0.1asdsdasdf34,asd.eerv5.3",
				"private_token":  "asd123",
				"repository_url": "http://repo.url",
				"commit_hash":    "aaa111",
				"api_base_url":   "http://api.baseurl",
				"preset_status":  "success",
			},
			expectedError:    nil,
			expectedCoverage: 0.1,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			// set up
			for k, v := range tC.envVars {
				os.Setenv(k, v)
			}

			conf, err := fixAndParseConfig()

			assert.Equal(t, tC.expectedError, err)
			assert.Equal(t, tC.expectedCoverage, conf.Coverage)

			// tear down
			for k := range tC.envVars {
				os.Setenv(k, "")
			}
		})
	}
}
