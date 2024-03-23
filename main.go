package main

import (
	"context"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"
)

type Options struct {
	// from flag
	from string
	to   string

	// from env
	owner       string
	repo        string
	gitHubToken string
	apiUrl      *url.URL
}

func getOptions() (Options, error) {
	from := flag.String("from", "", "source branch")
	to := flag.String("to", "", "target branch")
	flag.Parse()

	githubToken := os.Getenv("GITHUB_TOKEN")
	repository := strings.Split(os.Getenv("GITHUB_REPOSITORY"), "/")
	owner := repository[0]
	repo := repository[1]
	rawApiUrl := os.Getenv("GITHUB_API_URL")

	apiUrl, _ := url.Parse(rawApiUrl)

	return Options{
		from:        *from,
		to:          *to,
		owner:       owner,
		repo:        repo,
		gitHubToken: githubToken,
		apiUrl:      apiUrl,
	}, nil
}

func run(options Options) error {
	logger = GetLogger()
	ctx := context.Background()

	from := options.from
	to := options.to

	client := NewClient(GithubClientOptions{owner: options.owner, repo: options.repo, githubToken: options.gitHubToken, apiUrl: options.apiUrl})

	prNumbers, err := client.FetchPullRequestNumbers(ctx, from, to)
	if err != nil {
		return err
	}

	if len(prNumbers) == 0 {
		logger.Println("No pull requests were found for the release. Nothing to do.")
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

	logger.Println("Title of pull request:  ", title)

	pr, created, err := client.CreatePullRequest(ctx, title, body, from, to)
	if err != nil {
		return err
	}

	if created {
		logger.Println("Created new a pull request.", pr.GetNumber())
	} else {
		_, err := client.UpdatePullRequest(ctx, pr.GetNumber(), title, body)
		if err != nil {
			return err
		}
		logger.Println("The pull request already exists. The body was updated.", pr.GetNumber())
	}

	return nil
}

func main() {
	options, err := getOptions()

	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	err = run(options)

	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}
