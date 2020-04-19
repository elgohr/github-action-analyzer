# Github-Action-Analyzer
Analyzes the usage of Github Actions

## Install
`go get github.com/elgohr/github-action-analyzer`  
Make sure to have `GOBIN` in the `PATH`.

## Usage
At the moment there is just a summary
```
$ github-action-analyzer -name=my-repository -access-token=my-token summary
{
  "TotalRepositories": 1020,
  "TotalSteps": 1270,
  "With": {
    "myWithFirstOption": 20
    "myWithSecondOption": 1
  }
```