package main

import (
	"os"
	"strings"
	"testing"

	"github.com/google/go-github/v60/github"
)

func TestRenderTemplate(t *testing.T) {
	data := RenderTemplateData{
		PullRequests: []github.PullRequest{
			{
				Number: github.Int(1),
			},
			{
				Number: github.Int(2),
			},
		},
		Date: "2021-01-01",
	}

	template, err := RenderTemplate(nil, data)

	if err != nil {
		t.Errorf("RenderTemplate returned error: %v", err)
	}

	if !strings.Contains(template, "Release 2021-01-01") {
		t.Errorf("RenderTemplate returned %v, want %v", template, "2021-01-01")
	}

	if !strings.Contains(template, "#1") {
		t.Errorf("RenderTemplate returned %v, want %v", template, "#1")
	}
}

func TestRenderTemplateWithFilename(t *testing.T) {
	data := RenderTemplateData{
		PullRequests: []github.PullRequest{
			{
				Number: github.Int(1),
			},
			{
				Number: github.Int(2),
			},
		},
		Date: "2021-01-01",
	}

	tmpFile, err := os.CreateTemp("", "custom.mustache")
	if err != nil {
		panic(err)
	}
	defer os.Remove(tmpFile.Name())

	filename := tmpFile.Name()
	_, err = tmpFile.Write([]byte("This is custom template"))
	if err != nil {
		panic(err)
	}

	template, err := RenderTemplate(&filename, data)

	if err != nil {
		t.Errorf("RenderTemplate returned error: %v", err)
	}

	want := "This is custom template"
	if !strings.Contains(template, want) {
		t.Errorf("RenderTemplate returned %v, want %v", template, want)
	}
}
