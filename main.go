package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v60/github"
)

type Options struct {
	// from flag
	from     string
	to       string
	labels   []string
	template *string
	json     bool

	// from env
	owner       string
	repo        string
	gitHubToken string
	apiUrl      *url.URL
}

func getOptions() (Options, error) {
	from := flag.String("from", "", "The base branch name.")
	to := flag.String("to", "", "The target branch name.")
	labelsFlag := flag.String("labels", "", "Specify the labels to add to the pull request as a comma-separated list of strings.")
	template := flag.String("template", "", "The path to the template file.")
	json := flag.Bool("json", false, "Output the result as JSON.")
	flag.Parse()

	githubToken := os.Getenv("GITHUB_TOKEN")
	repository := strings.Split(os.Getenv("GITHUB_REPOSITORY"), "/")
	owner := repository[0]
	repo := repository[1]
	rawApiUrl := os.Getenv("GITHUB_API_URL")

	apiUrl, _ := url.Parse(rawApiUrl)

	var labels []string
	if *labelsFlag != "" {
		labels = strings.Split(*labelsFlag, ",")
	}

	return Options{
		from:        *from,
		to:          *to,
		labels:      labels,
		template:    template,
		json:        *json,
		owner:       owner,
		repo:        repo,
		gitHubToken: githubToken,
		apiUrl:      apiUrl,
	}, nil
}

type Result struct {
	IsCreated          bool                `json:"is_created,omitempty"`
	ReleasePullRequest *github.PullRequest `json:"release_pull_request,omitempty"`
}

func getResultJson(result Result) (string, error) {
	resultJson, err := json.Marshal(result)
	if err != nil {
		return "", err
	}
	return string(resultJson), nil
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
	data, err := RenderTemplate(options.template, RenderTemplateData{pullRequests, date})

	if err != nil {
		return err
	}

	parts := strings.SplitN(data, "\n", 2)
	title := parts[0]
	body := parts[1]

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

	if len(options.labels) > 0 {
		err := client.AddLabelsToPullRequest(ctx, pr.GetNumber(), options.labels)
		if err != nil {
			return err
		}
		logger.Println("Added labels to the pull request.", pr.GetNumber())
	}

	if options.json {
		result := Result{IsCreated: created, ReleasePullRequest: pr}
		resultJson, err := getResultJson(result)
		if err != nil {
			return err
		}

		fmt.Println(resultJson)
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
