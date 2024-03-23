package main

import (
	"context"
	"slices"
	"time"

	"github.com/google/go-github/v60/github"
)

type GithubClientOptions struct {
	owner string
	repo  string
}

type GithubClient struct {
	client *github.Client

	options GithubClientOptions
}

func NewClient(client *github.Client, options GithubClientOptions) *GithubClient {

	return &GithubClient{
		client:  client,
		options: options,
	}
}

func (c *GithubClient) FetchPullRequestNumbers(ctx context.Context, from string, to string) ([]int, error) {
	commitsComparison, _, err := c.client.Repositories.CompareCommits(ctx, c.options.owner, c.options.repo, to, from, nil)

	if err != nil {
		return nil, err
	}

	prNumbers := []int{}
	for i := 0; i < len(commitsComparison.Commits); i++ {
		commit := commitsComparison.Commits[i]

		pulls, _, err := c.client.PullRequests.ListPullRequestsWithCommit(ctx, c.options.owner, c.options.repo, commit.GetSHA(), nil)
		if err != nil {
			return nil, err
		}

		for j := 0; j < len(pulls); j++ {
			prNumbers = append(prNumbers, pulls[j].GetNumber())
		}
	}

	slices.Sort(prNumbers)

	return slices.Compact(prNumbers), nil
}

type PullRequest struct {
	Number   int
	Title    string
	Assignee string
	MergedAt time.Time
}

func (c *GithubClient) FetchPullRequests(ctx context.Context, prNumbers []int) ([]github.PullRequest, error) {
	pullRequests := []github.PullRequest{}

	for i := 0; i < len(prNumbers); i++ {
		prNumber := prNumbers[i]
		pr, _, err := c.client.PullRequests.Get(ctx, c.options.owner, c.options.repo, prNumber)
		if err != nil {
			return nil, err
		}

		pullRequests = append(pullRequests, *pr)
	}

	slices.SortFunc(pullRequests, func(a, b github.PullRequest) int {
		return a.MergedAt.Compare(b.MergedAt.Time)
	})

	return pullRequests, nil
}

func (c *GithubClient) CreatePullRequest(ctx context.Context, title, body, from, to string) (*github.PullRequest, bool, error) {
	prs, _, err := c.client.PullRequests.List(ctx, c.options.owner, c.options.repo, &github.PullRequestListOptions{
		Base:  to,
		Head:  from,
		State: "open",
	})

	if err != nil {
		return nil, false, err
	}

	if len(prs) > 0 {
		return prs[0], false, nil
	}

	pr, _, err := c.client.PullRequests.Create(ctx, c.options.owner, c.options.repo, &github.NewPullRequest{
		Title: &title,
		Body:  &body,
		Base:  &to,
		Head:  &from,
	})

	if err != nil {
		return nil, false, err
	}
	return pr, true, nil
}

func (c *GithubClient) UpdatePullRequest(ctx context.Context, prNumber int, title, body string) (*github.PullRequest, error) {
	pr, _, err := c.client.PullRequests.Edit(ctx, c.options.owner, c.options.repo, prNumber, &github.PullRequest{
		Title: &title,
		Body:  &body,
	})

	if err != nil {
		return nil, err
	}

	return pr, nil
}
