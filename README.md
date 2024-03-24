# git-pr-release-go

git-pr-release-go is a Go-based reimagination of the original Ruby CLI tool, [git-pr-release](https://github.com/x-motemen/git-pr-release).

Designed to streamline the development workflow, this tool automates the creation of "Release Pull Requests" on GitHub. Each "Release Pull Request" generated compiles a comprehensive list of pull requests slated for the upcoming release, facilitating a clear overview and seamless integration process.


![](./images/screenshot.png)

## Usage

```bash
$ git-pr-release-go --from main --to release/production
```

### GitHub Actions Usage

For this CLI to function within GitHub Actions, it requires the following permissions:

- `contents: read`
- `pull-requests: write`

Here's a sample workflow:

```yaml
name: Create Release Pull Request
on:
  push:
    branches:
      - main

jobs:
  create-release-pr:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      pull-requests: write
    steps:
      - name: Setup git-pr-release-go
        uses: KeisukeYamashita/setup-release@v1.0.2
        with:
          repository: odanado/git-pr-release-go
          arch: x86_64
          platform: "Linux"

      - run: git-pr-release-go --from main --to release/production
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

#### Using GitHub Apps Tokens

To authenticate using a GitHub Apps token, incorporate [actions/create-github-app-token](https://github.com/actions/create-github-app-token) in your workflow.

```yaml
name: Create Release Pull Request
on:
  push:
    branches:
      - main

jobs:
  create-release-pr:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/create-github-app-token@v1
        id: app-token
        with:
          app-id: ${{ vars.APP_ID }}
          private-key: ${{ secrets.PRIVATE_KEY }}

      - name: Setup git-pr-release-go
        uses: KeisukeYamashita/setup-release@v1.0.2
        with:
          repository: odanado/git-pr-release-go
          arch: x86_64
          platform: Linux

      - run: git-pr-release-go --from main --to release/production
        env:
          GITHUB_TOKEN: ${{ steps.app-token.outputs.token }}
```

### Options

- `--from`: The base branch name. Required.
- `--to`: The target branch name. Required.
- `--labels`: Specify the labels to add to the pull request as a comma-separated list of strings. Optional.
- `--template`: Specify the Mustache template file. Optional.

### Environment Variables

- `GITHUB_TOKEN`: GitHub API token. Required.
- `GITHUB_API_URL`: GitHub API URL. Optional.
- `GITHUB_REPOSITORY`: GitHub repository name. Required.

If you are using GitHub Actions, `GITHUB_API_URL` and `GITHUB_REPOSITORY` are automatically set by the runner and you do not need to specify them.

### Mustache template customization
Customize your pull request description with Mustache templates, leveraging variables like:

```json5
{
  // Execution date of the CLI
  "date": "yyyy-MM-dd",
  // Array of pull requests for the release, using fields from the GitHub REST API response.
  // https://docs.github.com/ja/rest/pulls/pulls?apiVersion=2022-11-28#list-pull-requests
  "pull_requests": []
}
```

For a practical example, refer to our [default template file](./git-pr-release.mustache).

## Compare with git-pr-release

This tool is developed in Go, eliminating the need for Ruby, as it operates entirely through a binary file.

While inspired by git-pr-release, this tool pays homage to its predecessor yet introduces several distinct features:

- By default, the pull request description is overwritten.
- Squash merging is supported without the need for additional options.
- A config file is not supported.
- Templates use Mustache files instead of ERB files.

## TODO
- [ ] Add more testing

## Release flow

- Create new tag `git tag -a vx.y.z -m ""`
- Push the tag `git push origin vx.y.z`
