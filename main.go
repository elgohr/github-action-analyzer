package main

import (
	"fmt"
	"github.com/elgohr/action-analyzer/analyzer"
	"github.com/elgohr/action-analyzer/downloader"
	"log"
	"os"
)

func main() {
	args := os.Args[1:]
	actionName := args[0]
	accessToken := args[1]
	d := downloader.NewDownloader()
	cs, err := d.DownloadConfigurations(actionName, accessToken)
	if err != nil {
		log.Fatalln(err)
	}
	res, err := analyzer.Analyze(actionName, cs)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(res)
}
