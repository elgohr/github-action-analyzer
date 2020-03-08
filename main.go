package main

import (
	"encoding/json"
	"fmt"
	"github.com/elgohr/action-analyzer/analyzer"
	"github.com/elgohr/action-analyzer/downloader"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		Name:  "Github Action Analyzer",
		Usage: "Helps analyzing the usage of Github Actions",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "name",
				Aliases: []string{"n"},
				Usage:   "the name of the repository to be analyzed (e.g. elgohr/action-analyzer)",
			},
			&cli.StringFlag{
				Name:    "access-token",
				Aliases: []string{"t"},
				Usage:   "personal access token to be used when searching for usages",
			},
			&cli.StringFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "verbose output",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "summary",
				Aliases: []string{"s"},
				Usage:   "summarizes the usage of the given action",
				Action: func(c *cli.Context) error {
					actionName := c.String("name")
					accessToken := c.String("access-token")
					d := downloader.NewDownloader()
					cs, errs := d.DownloadConfigurations(actionName, accessToken)
					go func(errs <-chan error) {
						for err := range errs {
							fmt.Fprintln(os.Stderr, err)
						}
					}(errs)
					res := analyzer.Analyze(actionName, cs)
					b, err := json.MarshalIndent(res, "", "  ")
					if err != nil {
						return err
					}
					fmt.Print(string(b))
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
