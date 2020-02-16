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
	r, err := analyzer.Analyze("elgohr/Publish-Docker-Github-Action", []downloader.ActionConfiguration{{Configuration: b}})
	assert.NoError(t, err)
	assert.Equal(t, 4, r.WithResult["name"])
	assert.Equal(t, 4, r.WithResult["username"])
	assert.Equal(t, 4, r.WithResult["password"])
	assert.Equal(t, 2, r.WithResult["registry"])
	assert.Equal(t, 4, r.WithResult["dockerfile"])
}
