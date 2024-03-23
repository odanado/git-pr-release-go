package main

import (
	"context"
	"net/url"
	"slices"
	"time"

	"github.com/google/go-github/v60/github"
)

type GithubClientOptions struct {
	owner       string
	repo        string
	githubToken string
	apiUrl      *url.URL
}

type GithubClient struct {
	client *github.Client

	owner string
	repo  string
}

func NewClient(options GithubClientOptions) *GithubClient {
	githubClient := github.NewClient(nil).WithAuthToken(options.githubToken)
	if options.apiUrl != nil {
		githubClient.BaseURL = options.apiUrl
	}

	return &GithubClient{
		client: githubClient,
		owner:  options.owner,
		repo:   options.repo,
	}
}

func (c *GithubClient) FetchPullRequestNumbers(ctx context.Context, from string, to string) ([]int, error) {
	commitsComparison, _, err := c.client.Repositories.CompareCommits(ctx, c.owner, c.repo, to, from, nil)

	if err != nil {
		return nil, err
	}

	prNumbers := []int{}
	for i := 0; i < len(commitsComparison.Commits); i++ {
		commit := commitsComparison.Commits[i]

		pulls, _, err := c.client.PullRequests.ListPullRequestsWithCommit(ctx, c.owner, c.repo, commit.GetSHA(), nil)
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
		pr, _, err := c.client.PullRequests.Get(ctx, c.owner, c.repo, prNumber)
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
	prs, _, err := c.client.PullRequests.List(ctx, c.owner, c.repo, &github.PullRequestListOptions{
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

	pr, _, err := c.client.PullRequests.Create(ctx, c.owner, c.repo, &github.NewPullRequest{
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
	pr, _, err := c.client.PullRequests.Edit(ctx, c.owner, c.repo, prNumber, &github.PullRequest{
		Title: &title,
		Body:  &body,
	})

	if err != nil {
		return nil, err
	}

	return pr, nil
}
