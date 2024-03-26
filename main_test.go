package main

import (
	"testing"

	"github.com/google/go-github/v60/github"
)

func TestGetResultJson(t *testing.T) {
	result := Result{IsCreated: true, ReleasePullRequest: &github.PullRequest{Number: github.Int(1)}}

	resultJson, err := getResultJson(result)

	if err != nil {
		t.Errorf("outputResult returned error: %v", err)
	}

	want := `{"is_created":true,"release_pull_request":{"number":1}}`
	if resultJson != want {
		t.Errorf("outputResult returned %v, want %v", resultJson, want)
	}
}
