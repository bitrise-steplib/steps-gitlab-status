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
		{"https with basic auth", "https://username:token@github.com/owner/repository.git", "owner/repository", false},
		{"ssh without scheme", "user@github.com:owner/repository.git", "owner/repository", false},
		{"ssh with scheme without port", "ssh://user@github.com/owner/repository.git", "owner/repository", false},
		{"ssh with scheme with port", "ssh://user@github.com:22/owner/repository.git", "owner/repository", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getRepo(tt.url); got != tt.want {
				t.Errorf("getRepo() = %v, want %v", got, tt.want)
			}
		})
	}
}
