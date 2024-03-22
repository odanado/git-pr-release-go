package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v60/github"
)

type Options struct {
	// from flag
	From string
	To   string

	// from env
	Owner   string
	Repo    string
	GhToken string
}

func getOptions() (Options, error) {
	from := flag.String("from", "", "source branch")
	to := flag.String("to", "", "target branch")
	flag.Parse()

	ghToken := os.Getenv("GITHUB_TOKEN")
	repository := strings.Split(os.Getenv("GITHUB_REPOSITORY"), "/")
	owner := repository[0]
	repo := repository[1]

	return Options{
		From:    *from,
		To:      *to,
		Owner:   owner,
		Repo:    repo,
		GhToken: ghToken,
	}, nil
}

func run(options Options) error {
	logger = GetLogger()
	ctx := context.Background()

	githubClient := github.NewClient(nil).WithAuthToken(options.GhToken)

	owner := options.Owner
	repo := options.Repo
	from := options.From
	to := options.To

	client := NewClient(githubClient, GithubClientOptions{owner, repo})

	prNumbers, err := client.FetchPullRequestNumbers(ctx, from, to)
	if err != nil {
		return err
	}

	if len(prNumbers) == 0 {
		logger.Println("No pull requests found")
		return nil
	}

	logger.Println("Found pull requests: ", prNumbers)

	pullRequests, err := client.FetchPullRequests(ctx, prNumbers)

	if err != nil {
		return err
	}

	currentTime := time.Now()
	date := currentTime.Format("2006-01-02")
	data, err := RenderTemplate(nil, RenderTemplateData{pullRequests, date})
	parts := strings.SplitN(data, "\n", 2)

	title := parts[0]
	body := parts[1]

	if err != nil {
		return err
	}

	logger.Println("Title: ", title)

	pr, created, err := client.CreatePullRequest(ctx, title, body, from, to)
	if err != nil {
		return err
	}

	if !created {
		_, err := client.UpdatePullRequest(ctx, pr.GetNumber(), title, body)
		if err != nil {
			return err
		}
	}

	if created {
		logger.Println("created: ", pr.GetNumber())
	} else {
		logger.Println("updated: ", pr.GetNumber())
	}

	return nil
}

func main() {
	options, err := getOptions()

	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	fmt.Println("Options: ", options)

	err = run(options)

	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}
