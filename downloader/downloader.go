package downloader

import (
	"encoding/json"
	"fmt"
	"github.com/tomnomnom/linkheader"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type Downloader struct {
	Client        *http.Client
	GithubApiRoot string
	CacheDir      string
}

func NewDownloader() *Downloader {
	cacheDir := "/tmp/action-analyzer/"
	if err := os.MkdirAll(cacheDir, os.ModePerm); err != nil {
		log.Fatalln(err)
	}
	return &Downloader{
		Client:        http.DefaultClient,
		GithubApiRoot: "https://api.github.com",
		CacheDir:      cacheDir,
	}
}

func (d Downloader) DownloadConfigurations(actionName string, personalAccessToken string) (<-chan ActionConfiguration, <-chan error) {
	authHeader := "token " + personalAccessToken
	configurations := make(chan ActionConfiguration)
	errs := make(chan error)

	go func(configurations chan ActionConfiguration) {
		url := fmt.Sprintf("%s/search/code?q=%s+in:file+language:yaml&page=1", d.GithubApiRoot, actionName)
		wg := sync.WaitGroup{}
		d.search(url, authHeader, configurations, &wg, errs)
		wg.Wait()
		close(configurations)
	}(configurations)

	return configurations, errs
}

func (d Downloader) search(url string, authHeader string, configurations chan ActionConfiguration, wg *sync.WaitGroup, errs chan error) {
	fmt.Println("indexing " + url)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("Authorization", authHeader)
	searchRes, err := d.Client.Do(req)
	if err != nil {
		errs <- err
		return
	}
	defer searchRes.Body.Close()
	searchBody, err := ioutil.ReadAll(searchRes.Body)
	if err != nil {
		errs <- err
		return
	}
	var searchResponse SearchResponse
	if err := json.Unmarshal(searchBody, &searchResponse); err != nil {
		errs <- err
		return
	}

	wg.Add(len(searchResponse.Items))
	go func(wg *sync.WaitGroup) {
		d.downloadConfiguration(searchResponse, authHeader, configurations, wg, errs)
	}(wg)

	paginationLink := searchRes.Header.Get("Link")
	links := linkheader.Parse(paginationLink)
	for _, link := range links {
		if link.Rel == "next" {
			waitWhenRateLimited(searchRes, errs)
			d.search(link.URL, authHeader, configurations, wg, errs)
		}
	}
}

func (d Downloader) downloadConfiguration(searchResponse SearchResponse, authHeader string, configurations chan ActionConfiguration, wg *sync.WaitGroup, errs chan error) {
	for _, item := range searchResponse.Items {
		fmt.Println(fmt.Sprintf("downloading configuration for %s", item.Repository.FullName))
		req, _ := http.NewRequest(http.MethodGet, item.Url, nil)
		req.Header.Add("Authorization", authHeader)
		refRes, err := d.Client.Do(req)
		if err != nil {
			errs <- err
			return
		}
		defer refRes.Body.Close()
		refBody, err := ioutil.ReadAll(refRes.Body)
		if err != nil {
			errs <- err
			return
		}
		var refResponse RefResponse
		if err := json.Unmarshal(refBody, &refResponse); err != nil {
			errs <- err
			return
		}
		waitWhenRateLimited(refRes, errs)

		downloadReq, _ := http.NewRequest(http.MethodGet, refResponse.DownloadUrl, nil)
		downloadReq.Header.Add("Authorization", authHeader)
		downloadRes, err := d.Client.Do(downloadReq)
		if err != nil {
			errs <- err
			return
		}
		defer downloadRes.Body.Close()
		body, err := ioutil.ReadAll(downloadRes.Body)
		if err != nil {
			errs <- err
			return
		}
		configurations <- ActionConfiguration{
			Name:          item.Repository.FullName,
			Configuration: body,
		}
		wg.Done()
	}
}

func waitWhenRateLimited(res *http.Response, errs chan error) {
	remainingHeader := res.Header.Get("X-RateLimit-Remaining")
	if remainingHeader != "" {
		remaining, err := strconv.ParseInt(remainingHeader, 10, 32)
		if err != nil {
			errs <- err
			return
		}

		if remaining < 3 {
			resetTimeEpoch, err := strconv.ParseInt(res.Header.Get("X-RateLimit-Reset"), 10, 64)
			if err != nil {
				errs <- err
				return
			}
			resetTime := time.Unix(resetTimeEpoch, 0)
			waitTime := time.Until(resetTime)
			fmt.Println(fmt.Sprintf("Waiting %v seconds due to rate-limit", waitTime.Seconds()))
			time.Sleep(waitTime)
		}
	}
}

type SearchResponse struct {
	Items []SearchItem `json:"items"`
}

type SearchItem struct {
	Url        string `json:"url"`
	Repository struct {
		FullName string `json:"full_name"`
	}
}

type RefResponse struct {
	DownloadUrl string `json:"download_url"`
}

type ActionConfiguration struct {
	Name          string
	Configuration []byte
}
