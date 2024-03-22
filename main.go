package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v60/github"
)

type Options struct {
	// from flag
	DryRun bool
	From   string
	To     string

	// from env
	Owner   string
	Repo    string
	GhToken string
}

func getOptions() (Options, error) {
	dryRunFlag := flag.Bool("dry-run", false, "perform a dry run; does not update PR")
	from := flag.String("from", "", "source branch")
	to := flag.String("to", "", "target branch")
	flag.Parse()

	ghToken := os.Getenv("GITHUB_TOKEN")
	repository := strings.Split(os.Getenv("GITHUB_REPOSITORY"), "/")
	owner := repository[0]
	repo := repository[1]

	return Options{
		DryRun:  *dryRunFlag,
		From:    *from,
		To:      *to,
		Owner:   owner,
		Repo:    repo,
		GhToken: ghToken,
	}, nil
}

func run(options Options) error {
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

	fmt.Println("prNumbers: ", prNumbers)

	pullRequests, err := client.FetchPullRequests(ctx, prNumbers)

	if err != nil {
		return err
	}

	fmt.Println("pullRequests: ", pullRequests[0].Number)

	data, err := RenderTemplate(nil, RenderTemplateData{pullRequests})
	parts := strings.SplitN(data, "\n", 2)

	title := parts[0]
	body := parts[1]

	if err != nil {
		return err
	}

	fmt.Println("data: ", data)

	pr, err := client.CreatePullRequest(ctx, title, body, from, to)
	if err != nil {
		return err
	}

	fmt.Println("created: ", pr)

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
