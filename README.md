# git-pr-release-go

Go implementation of [git-pr-release](https://github.com/x-motemen/git-pr-release).

This command creates "Release Pull Request" on GitHub. The body of "Release Pull Request" lists the pull requests included in that release.

![](./images/screenshot.png)

# Usage

```bash
$ git-pr-release-go --from main --to release/production
```

## Usage in GitHub Actions

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

      - name: Download git-pr-release-go
        run: |
          mkdir -p tmp
          gh release download --repo odanado/git-pr-release-go --pattern "*Linux_x86_64*" --output - | tar zxvf - -C tmp
          mkdir -p /opt/git-pr-release-go/bin
          mv tmp/git-pr-release-go /opt/git-pr-release-go/bin
          rm -rf tmp

          echo "/opt/git-pr-release-go/bin" >> $GITHUB_PATH
        env:
          GH_TOKEN: ${{ github.token }}

      - run: ./git-pr-release-go --from main --to release/production
        working-directory: git-pr-release-go
        env:
          GITHUB_TOKEN: ${{ steps.app-token.outputs.token }}
```

## Options

- `--from`: The base branch name. Required.
- `--to`: The target branch name. Required.

## Environment Variables

- `GITHUB_TOKEN`: GitHub API token. Required.
- `GITHUB_API_URL`: GitHub API URL. Optional.
- `GITHUB_REPOSITORY`: GitHub repository name. Required.

If you are using GitHub Actions, `GITHUB_API_URL` and `GITHUB_REPOSITORY` are automatically set by the runner and you do not need to specify them.

# Compare with git-pr-release

This tool is developed in Go, eliminating the need for Ruby, as it operates entirely through a binary file.

While inspired by git-pr-release, this tool pays homage to its predecessor yet introduces several distinct features:

- By default, the pull request description is overwritten.
- Squash merging is supported without the need for additional options.
- A config file is not supported.
- Templates use Mustache files instead of ERB files.

# TODO
- [ ] Support a custom template
- [ ] Support custom labels
- [ ] Add more testing

# Release flow

- Create new tag `git tag -a vx.y.z -m ""`
- Push the tag `git push origin vx.y.z`
