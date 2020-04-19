package analyzer

import (
	"fmt"
	"github.com/elgohr/github-action-analyzer/downloader"
	"gopkg.in/yaml.v2"
	"strings"
)

func Analyze(actionName string, configs <-chan downloader.ActionConfiguration) *Result {
	result := Result{
		TotalRepositories: 0,
		TotalSteps:        0,
		With:              map[string]int{},
	}
	for config := range configs {
		fmt.Println(fmt.Sprintf("analyzing usage in %s", config.Name))
		result.TotalRepositories += 1
		var parsedConfig Configuration
		if err := yaml.Unmarshal(config.Configuration, &parsedConfig); err != nil {
			fmt.Println(err)
		}
		for _, build := range parsedConfig.Jobs {
			for _, step := range build.Steps {
				if strings.HasPrefix(step.Uses, actionName) {
					result.TotalSteps += 1
					for key := range step.With {
						if count, exists := result.With[key]; exists {
							result.With[key] = count + 1
						} else {
							result.With[key] = 1
						}
					}
				}
			}
		}
	}
	return &result
}

type Result struct {
	TotalRepositories int
	TotalSteps        int
	With              map[string]int
}

type Configuration struct {
	Jobs map[string]struct {
		Steps []struct {
			Uses string
			With map[string]string
		}
	}
}
