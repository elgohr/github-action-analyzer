package analyzer

import (
	"fmt"
	"github.com/elgohr/github-action-analyzer/downloader"
	"gopkg.in/yaml.v2"
	"strings"
)

func GetSummary(actionName string, configs <-chan downloader.ActionConfiguration) *Summary {
	summary := Summary{
		TotalRepositories: 0,
		TotalSteps:        0,
		With:              map[string]int{},
	}
	for config := range configs {
		fmt.Println(fmt.Sprintf("analyzing usage in %s", config.Name))
		summary.TotalRepositories += 1
		var parsedConfig Configuration
		if err := yaml.Unmarshal(config.Configuration, &parsedConfig); err != nil {
			fmt.Println(err)
		}
		for _, build := range parsedConfig.Jobs {
			for _, step := range build.Steps {
				if strings.HasPrefix(step.Uses, actionName) {
					summary.TotalSteps += 1
					for key := range step.With {
						if count, exists := summary.With[key]; exists {
							summary.With[key] = count + 1
						} else {
							summary.With[key] = 1
						}
					}
				}
			}
		}
	}
	return &summary
}

type Summary struct {
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
