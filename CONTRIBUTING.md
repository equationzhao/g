# Contribution Guidelines

This Contribution Guidelines is modified from the [GitHub doc](https://github.com/github/docs/) project

Feel free to ask any questions in the [Discussions QA](https://github.com/Equationzhao/g/discussions/categories/q-a)

## Issues

### Create a new issue

If you have the following ideas

- spot a problem with the project 🐛
- have any suggestions 📈
- want a new style 💄

[search if an issue already exists](https://docs.github.com/en/github/searching-for-information-on-github/searching-on-github/searching-issues-and-pull-requests#search-by-the-title-body-or-comments). \
If a related issue doesn't exist, you can open a new issue using a relevant [issue form](https://github.com/Equationzhao/g/issues/new/choose).


### Solve an issue

Scan through our [existing issues](https://github.com/Equationzhao/g/issues) to find one that interests you. You can narrow down the search using `labels` as filters.
If you find an issue to work on, you are welcome to open a PR with a fix.


## Make changes

1. Fork and Clone the Repository
```shell
git clone git@github.com:yourname/g.git
```
2. Make sure you have installed **go**, and go version >= 1.26.1 .
```shell
go version
```
3. Make your changes!

## Commit messages
It's recommended to follow the commit style below:
```text
<type>[optional scope]: <description>

[optional body]
```

```text
fix:
feat:
build:
chore:
ci:
docs:
style:
refactor:
perf:
test:
...
```
Also, you can use [gitmoji](https://gitmoji.dev) in the commit message

Examples are provided to illustrate the recommended style:
```text
style: :lipstick: change color for 'readme' file
    
change from BrightYellow to Yellow
```

## Tests

Please refer to the [TestWorkflow](docs/TestWorkflow.md)

## PR

When you're finished with the changes, create a pull request, also known as a PR.\
Before you submit your pr, please check out the following stuffs

- [ ] you have run `go mod tidy` and `gofumpt --extra -w -l .`
- [ ] your code has the necessary comments and documentation (if needed)
- [ ] you have written tests for your changes (if needed)
- [ ] Pass other tests, ***if you're breaking the tests, please explain why in the PR description***

and when you're submitting your pr, make sure to:
- Fill the "Ready for review" template so that we can review your PR. This template helps reviewers understand your changes as well as the purpose of your pull request.
- Don't forget to [link PR to issue](https://docs.github.com/en/issues/tracking-your-work-with-issues/linking-a-pull-request-to-an-issue) if you are solving one.
- Enable the checkbox to [allow maintainer edits](https://docs.github.com/en/github/collaborating-with-issues-and-pull-requests/allowing-changes-to-a-pull-request-branch-created-from-a-fork) so the branch can be updated for a merge.
  Once you submit your PR, maintainers will review your proposal. We may ask questions or request additional information.
- We may ask for changes to be made before a PR can be merged, either using [suggested changes](https://docs.github.com/en/github/collaborating-with-issues-and-pull-requests/incorporating-feedback-in-your-pull-request) or pull request comments. You can apply suggested changes directly through the UI. You can make any other changes in your fork, then commit them to your branch.
- As you update your PR and apply changes, mark each conversation as [resolved](https://docs.github.com/en/github/collaborating-with-issues-and-pull-requests/commenting-on-a-pull-request#resolving-conversations).
- If you run into any merge issues, checkout this [git tutorial](https://github.com/skills/resolve-merge-conflicts) to help you resolve merge conflicts and other issues.

## Release

Please refer to the [ReleaseWorkflow](docs/ReleaseWorkflow.md)
