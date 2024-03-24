package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-github/v60/github"
)

func TestFetchPullRequestNumbers(t *testing.T) {
	ctx := context.Background()

	mux := http.NewServeMux()

	mux.HandleFunc(
		"/repos/owner/repo/compare/to...from",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `{"commits": [{"sha": "sha1"}, {"sha": "sha2"}, {"sha": "sha3"}]}`)
		},
	)
	mux.HandleFunc(
		"/repos/owner/repo/commits/sha1/pulls",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `[{"number": 1}]`)
		},
	)
	mux.HandleFunc(
		"/repos/owner/repo/commits/sha2/pulls",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `[{"number": 1}]`)
		},
	)
	mux.HandleFunc(
		"/repos/owner/repo/commits/sha3/pulls",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `[{"number": 2}]`)
		},
	)

	ts := httptest.NewServer(mux)
	defer ts.Close()

	apiUrl, _ := url.Parse(ts.URL)
	client := NewClient(GithubClientOptions{owner: "owner", repo: "repo", githubToken: "githubToken", apiUrl: apiUrl})

	prNumbers, err := client.FetchPullRequestNumbers(ctx, "from", "to")

	if err != nil {
		t.Errorf("PullRequests.Get returned error: %v", err)
	}

	want := []int{1, 2}
	if !cmp.Equal(prNumbers, want) {
		t.Errorf("PullRequests.List returned %+v, want %+v", prNumbers, want)
	}
}

func TestFetchPullRequests(t *testing.T) {
	ctx := context.Background()

	mux := http.NewServeMux()

	mux.HandleFunc(
		"/repos/owner/repo/pulls/1",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `{"number": 1, "merged_at": "2021-02-01T00:00:00Z"}`)
		},
	)
	mux.HandleFunc(
		"/repos/owner/repo/pulls/2",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `{"number": 2, "merged_at": "2021-01-01T00:00:00Z"}`)
		},
	)

	ts := httptest.NewServer(mux)
	defer ts.Close()

	apiUrl, _ := url.Parse(ts.URL)
	client := NewClient(GithubClientOptions{owner: "owner", repo: "repo", githubToken: "githubToken", apiUrl: apiUrl})

	prNumbers := []int{1, 2}
	prs, err := client.FetchPullRequests(ctx, prNumbers)

	if err != nil {
		t.Errorf("PullRequests.Get returned error: %v", err)
	}

	time1, _ := time.Parse("2006-01-02T15:04:05Z", "2021-01-01T00:00:00Z")
	time2, _ := time.Parse("2006-01-02T15:04:05Z", "2021-02-01T00:00:00Z")
	want := []github.PullRequest{
		{Number: github.Int(2), MergedAt: &github.Timestamp{Time: time1}},
		{Number: github.Int(1), MergedAt: &github.Timestamp{Time: time2}},
	}

	if !cmp.Equal(prs, want) {
		t.Errorf("PullRequests.List returned %+v, want %+v", prs, want)
	}
}

func TestCreatePullRequest(t *testing.T) {
	ctx := context.Background()

	mux := http.NewServeMux()

	mux.HandleFunc(
		"/repos/owner/repo/pulls",
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" {
				fmt.Fprint(w, `[]`)
			}
			if r.Method == "POST" {
				fmt.Fprint(w, `{"number": 1}`)
			}
		},
	)

	ts := httptest.NewServer(mux)
	defer ts.Close()

	apiUrl, _ := url.Parse(ts.URL)
	client := NewClient(GithubClientOptions{owner: "owner", repo: "repo", githubToken: "githubToken", apiUrl: apiUrl})

	pr, created, err := client.CreatePullRequest(ctx, "title", "body", "from", "to")

	if err != nil {
		t.Errorf("PullRequests.Get returned error: %v", err)
	}

	want := &github.PullRequest{Number: github.Int(1)}
	if !cmp.Equal(pr, want) {
		t.Errorf("PullRequests.List returned %+v, want %+v", pr, want)
	}

	if !created {
		t.Errorf("PullRequests.List returned %+v, want %+v", created, true)
	}
}

func TestCreatePullRequest_alreadyExists(t *testing.T) {
	ctx := context.Background()

	mux := http.NewServeMux()

	mux.HandleFunc(
		"/repos/owner/repo/pulls",
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" {
				fmt.Fprint(w, `[{"number": 1}]`)
			}
		},
	)

	ts := httptest.NewServer(mux)
	defer ts.Close()

	apiUrl, _ := url.Parse(ts.URL)
	client := NewClient(GithubClientOptions{owner: "owner", repo: "repo", githubToken: "githubToken", apiUrl: apiUrl})

	pr, created, err := client.CreatePullRequest(ctx, "title", "body", "from", "to")

	if err != nil {
		t.Errorf("PullRequests.Get returned error: %v", err)
	}

	want := &github.PullRequest{Number: github.Int(1)}
	if !cmp.Equal(pr, want) {
		t.Errorf("PullRequests.List returned %+v, want %+v", pr, want)
	}

	if created {
		t.Errorf("PullRequests.List returned %+v, want %+v", created, false)
	}
}
