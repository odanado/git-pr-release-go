package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"

	"github.com/cbroglie/mustache"
	"github.com/google/go-github/v60/github"
)

//go:embed git-pr-release.mustache
var defaultTemplate string

func readTemplate(filename *string) (string, error) {
	if filename == nil || *filename == "" {
		return defaultTemplate, nil
	}

	data, err := os.ReadFile(*filename)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

type RenderTemplateData struct {
	PullRequests     []github.PullRequest      `json:"pull_requests"`
	Commits          []github.RepositoryCommit `json:"commits"`
	Date             string                    `json:"date"`
	From             string                    `json:"from"`
	To               string                    `json:"to"`
	CustomParameters any                       `json:"custom_parameters"`
}

func convertJson(data RenderTemplateData) (any, error) {
	var jsonData any

	jsonByte, err := json.Marshal(data)

	json.Unmarshal(jsonByte, &jsonData)

	if err != nil {
		return nil, err
	}

	return jsonData, nil

}

func getRunUrl() string {
	serverUrl := os.Getenv("GITHUB_SERVER_URL")
	repo := os.Getenv("GITHUB_REPOSITORY")
	runId := os.Getenv("GITHUB_RUN_ID")
	runAttempt := os.Getenv("GITHUB_RUN_ATTEMPT")

	if serverUrl != "" && repo != "" && runId != "" && runAttempt != "" {
		return serverUrl + "/" + repo + "/actions/runs/" + runId + "/attempts/" + runAttempt
	}
	return ""
}

func RenderTemplate(filename *string, data RenderTemplateData, disableGeneratedByMessage bool) (string, error) {
	template, err := readTemplate(filename)

	if err != nil {
		return "", err
	}

	jsonData, err := convertJson(data)

	if err != nil {
		return "", err
	}

	text, err := mustache.Render(template, jsonData)

	if err != nil {
		return "", err
	}

	if !disableGeneratedByMessage {
		withinText := ""
		runUrl := getRunUrl()
		if runUrl != "" {
			withinText = " within [GitHub Actions workflow](" + runUrl + ")"
		}
		footer := `
---
*Automatically generated by [git-pr-release-go](https://github.com/odanado/git-pr-release-go)%s.*
`
		text += fmt.Sprintf(footer, withinText)
	}

	return text, nil
}
