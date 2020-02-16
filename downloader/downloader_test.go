package downloader_test

import (
	"fmt"
	"github.com/elgohr/action-analyzer/downloader"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDownloadsContent(t *testing.T) {
	var (
		searched      bool
		loadedRef     bool
		loadedDetails bool
	)
	details := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "token ACCESS_TOKEN", r.Header.Get("Authorization"))
		assert.Equal(t, "/cnrun/strava-x-api/be8cc384c5a3f136a2cac2cc5c561839db62f674/.github/workflows/dockerimage.yml", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)
		loadedDetails = true
	}))
	defer details.Close()
	ref := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "token ACCESS_TOKEN", r.Header.Get("Authorization"))
		assert.Equal(t, "/repositories/142542006/contents/.github/workflows/docker.yml", r.URL.Path)
		assert.Equal(t, "ref=136423fdf813227c87197369d69906b90731424a", r.URL.RawQuery)
		assert.Equal(t, http.MethodGet, r.Method)
		b, err := ioutil.ReadFile("testdata/ref_response.json")
		if err != nil {
			log.Fatalln(err)
		}
		withMockUrl := fmt.Sprintf(string(b), details.URL)
		if _, err := w.Write([]byte(withMockUrl)); err != nil {
			log.Fatalln(err)
		}
		loadedRef = true
	}))
	defer ref.Close()
	search := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "token ACCESS_TOKEN", r.Header.Get("Authorization"))
		assert.Equal(t, "/search/code", r.URL.Path)
		assert.Equal(t, "q=my-action+in:file+language:yaml", r.URL.RawQuery)
		assert.Equal(t, http.MethodGet, r.Method)
		b, err := ioutil.ReadFile("testdata/search_response.json")
		if err != nil {
			log.Fatalln(err)
		}
		withMockUrl := fmt.Sprintf(string(b), ref.URL)
		if _, err := w.Write([]byte(withMockUrl)); err != nil {
			log.Fatalln(err)
		}
		searched = true
	}))
	defer search.Close()
	a := downloader.NewDownloader()
	a.GithubApiRoot = search.URL
	configurations, err := a.DownloadConfigurations("my-action", "ACCESS_TOKEN")
	assert.Equal(t, 1, len(configurations))
	assert.NoError(t, err)
	assert.True(t, searched)
	assert.True(t, loadedRef)
	assert.True(t, loadedDetails)
}

func TestErrorsWhenFailingToDownloadContent(t *testing.T) {
	a := downloader.NewDownloader()
	a.GithubApiRoot = "http://localhost"
	_, err := a.DownloadConfigurations("my-action", "ACCESS_TOKEN")
	assert.Error(t, err)
}

func TestErrorsWhenSearchResponseIsNoJson(t *testing.T) {
	search := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("-"))
	}))
	defer search.Close()
	a := downloader.NewDownloader()
	a.GithubApiRoot = search.URL
	_, err := a.DownloadConfigurations("my-action", "ACCESS_TOKEN")
	assert.Error(t, err)
}

func TestErrorsWhenRefEndpointErrors(t *testing.T) {
	search := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("-")); err != nil {
			log.Fatalln(err)
		}
	}))
	defer search.Close()
	a := downloader.NewDownloader()
	a.GithubApiRoot = search.URL
	_, err := a.DownloadConfigurations("my-action", "ACCESS_TOKEN")
	assert.Error(t, err)
}

func TestErrorsWhenDownloadEndpointErrors(t *testing.T) {
	ref := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadFile("testdata/ref_response.json")
		if err != nil {
			log.Fatalln(err)
		}
		withMockUrl := fmt.Sprintf(string(b), "-")
		if _, err := w.Write([]byte(withMockUrl)); err != nil {
			log.Fatalln(err)
		}
	}))
	defer ref.Close()
	search := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadFile("testdata/search_response.json")
		if err != nil {
			log.Fatalln(err)
		}
		withMockUrl := fmt.Sprintf(string(b), ref.URL)
		if _, err := w.Write([]byte(withMockUrl)); err != nil {
			log.Fatalln(err)
		}
	}))
	defer search.Close()
	a := downloader.NewDownloader()
	a.GithubApiRoot = search.URL
	_, err := a.DownloadConfigurations("my-action", "ACCESS_TOKEN")
	assert.Error(t, err)
}

func TestNewAnalyzer(t *testing.T) {
	a := downloader.NewDownloader()
	assert.Equal(t, "https://api.github.com", a.GithubApiRoot)
	assert.NotNil(t, a.Client)
}
