package main

import (
	"context"
	"fmt"
	"os"
)

func exitWithError(err error) {
	fmt.Fprintln(os.Stderr, "Error: ", err)
	os.Exit(1)
}

func main() {

	fmt.Println("Hello, world!")

	owner := "odanado"
	repo := "git-pr-release-go"
	githubToken := os.Getenv("GH_TOKEN")

	client := NewClient(GithubClientOptions{owner: owner, repo: repo, githubToken: githubToken})

	ctx := context.Background()
	from := "main"
	to := "release/cli"

	totalCommits, pullRequests, commits, err := client.FetchChanges(ctx, from, to)

	if err != nil {
		exitWithError(err)
	}

	fmt.Println("Total commits: ", totalCommits)
	fmt.Println("Pull requests: ", len(pullRequests))
	for _, pr := range pullRequests {
		fmt.Println("  ", pr.GetNumber(), pr.GetTitle())
	}

	fmt.Println("Commits: ")
	for _, commit := range commits {
		fmt.Println("  ", commit.GetSHA(), commit.GetCommit().GetMessage())
	}

}
