package analyzer_test

import (
	"github.com/elgohr/action-analyzer/analyzer"
	"github.com/elgohr/action-analyzer/downloader"
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
		assert.Equal(t, 4, r.WithUsages["name"])
		assert.Equal(t, 4, r.WithUsages["username"])
		assert.Equal(t, 4, r.WithUsages["password"])
		assert.Equal(t, 2, r.WithUsages["registry"])
		assert.Equal(t, 4, r.WithUsages["dockerfile"])
	}()
}
