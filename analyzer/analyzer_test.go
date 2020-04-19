package analyzer_test

import (
	"github.com/elgohr/github-action-analyzer/analyzer"
	"github.com/elgohr/github-action-analyzer/downloader"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"testing"
)

func TestAnalyze(t *testing.T) {
	b, err := ioutil.ReadFile("testdata/details_response.yml")
	if err != nil {
		log.Fatalln(err)
	}
	configs := make(chan downloader.ActionConfiguration, 1)
	configs <- downloader.ActionConfiguration{Configuration: b}
	go func() {
		r := analyzer.Analyze("elgohr/Publish-Docker-Github-Action", configs)
		assert.Equal(t, 4, r.With["name"])
		assert.Equal(t, 4, r.With["username"])
		assert.Equal(t, 4, r.With["password"])
		assert.Equal(t, 2, r.With["registry"])
		assert.Equal(t, 4, r.With["dockerfile"])
	}()
}
