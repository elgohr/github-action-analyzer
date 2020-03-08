# Github-Action-Analyzer
Analyzer for the usage of Github Actions

## Usage
At the moment there is just a summary
```
$ anlyzer -name=my-repository -access-token=my-token summary
{
  "TotalRepositories": 1020,
  "TotalSteps": 1270,
  "WithUsages": {
    "actions_token": 1,
    "auto_tag": 1,
    "buildargs": 90,
    "buildoptions": 8,
    "cache": 75,
    "context": 26,
    "dockerfile": 332,
    "name": 1270,
    "password": 1270,
    "registry": 346,
    "snapshot": 163,
    "source-url": 2,
    "tag": 10,
    "tag-names": 1,
    "tag_names": 177,
    "tag_semver": 6,
    "tagging": 9,
    "tags": 168,
    "username": 1270,
    "workDir": 2,
    "workdir": 193
  }

```