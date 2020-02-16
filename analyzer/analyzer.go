package analyzer

import (
	"github.com/elgohr/action-analyzer/downloader"
	"gopkg.in/yaml.v2"
	"strings"
)

func Analyze(actionName string, configs []downloader.ActionConfiguration) (*Result, error) {
	var result Result
	result.WithResult = map[string]int{}
	for _, config := range configs {
		var parsedConfig Configuration
		if err := yaml.Unmarshal(config.Configuration, &parsedConfig); err != nil {
			return nil, err
		}
		for _, step := range parsedConfig.Jobs.Build.Steps {
			if strings.HasPrefix(step.Uses, actionName) {
				for key := range step.With {
					if count, exists := result.WithResult[key]; exists {
						result.WithResult[key] = count + 1
					} else {
						result.WithResult[key] = 1
					}
				}
			}
		}
	}
	return &result, nil
}

type Result struct {
	WithResult map[string]int
}

type Configuration struct {
	Jobs struct {
		Build struct {
			Steps []struct {
				Uses string
				With map[string]string
			}
		}
	}
}
