package downloader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Downloader struct {
	Client        *http.Client
	GithubApiRoot string
}

func NewDownloader() *Downloader {
	return &Downloader{
		Client:        http.DefaultClient,
		GithubApiRoot: "https://api.github.com",
	}
}

func (d *Downloader) DownloadConfigurations(actionName string, personalAccessToken string) ([]ActionConfiguration, error) {
	url := fmt.Sprintf("%s/search/code?q=%s+in:file+language:yaml", d.GithubApiRoot, actionName)
	authHeader := "token " + personalAccessToken
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("Authorization", authHeader)
	searchRes, err := d.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer searchRes.Body.Close()
	searchBody, err := ioutil.ReadAll(searchRes.Body)
	if err != nil {
		return nil, err
	}
	var searchResponse SearchResponse
	if err := json.Unmarshal(searchBody, &searchResponse); err != nil {
		return nil, err
	}

	var configurations []ActionConfiguration
	for _, item := range searchResponse.Items {
		req, _ := http.NewRequest(http.MethodGet, item.Url, nil)
		req.Header.Add("Authorization", authHeader)
		refRes, err := d.Client.Do(req)
		if err != nil {
			return nil, err
		}
		defer refRes.Body.Close()
		refBody, err := ioutil.ReadAll(refRes.Body)
		if err != nil {
			return nil, err
		}
		var refResponse RefResponse
		if err := json.Unmarshal(refBody, &refResponse); err != nil {
			return nil, err
		}

		downloadReq, _ := http.NewRequest(http.MethodGet, refResponse.DownloadUrl, nil)
		downloadReq.Header.Add("Authorization", authHeader)
		downloadRes, err := d.Client.Do(downloadReq)
		if err != nil {
			return nil, err
		}
		defer downloadRes.Body.Close()
		downloadBody, err := ioutil.ReadAll(downloadRes.Body)
		if err != nil {
			return nil, err
		}
		configuration := ActionConfiguration{Configuration: downloadBody}
		configurations = append(configurations, configuration)
	}

	return configurations, nil
}

type SearchResponse struct {
	Items []SearchItem `json:"items"`
}

type SearchItem struct {
	Url string `json:"url"`
}

type RefResponse struct {
	DownloadUrl string `json:"download_url"`
}

type ActionConfiguration struct {
	Configuration []byte
}
