package main

import (
	"testing"
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
