package main

import (
	_ "embed"
	"os"

	"github.com/cbroglie/mustache"
)

//go:embed git-pr-release.mustache
var defaultTemplate string

func readTemplate(filename *string) (string, error) {
	if filename == nil {
		return defaultTemplate, nil
	}

	data, err := os.ReadFile(*filename)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

type RenderTemplateData struct {
	PullRequests []PullRequest
	Date         string
}

func RenderTemplate(filename *string, data RenderTemplateData) (string, error) {
	template, err := readTemplate(filename)

	if err != nil {
		return "", err
	}

	return mustache.Render(template, map[string]interface{}{
		"pullRequests": data.PullRequests,
		"date":         data.Date,
	})
}
