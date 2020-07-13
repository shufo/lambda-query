package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// Version provides ecs-fargate version
var Version = "default"

func main() {

	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)

	cli.VersionFlag = &cli.BoolFlag{
		Name: "version", Aliases: []string{"V"},
		Usage: "print only the version",
	}

	app := &cli.App{
		Name:     "lambda-query",
		Version:  Version,
		Compiled: time.Now(),
		Usage:    "run db query from lambda function",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "function, f",
				Value:    "",
				Aliases:  []string{"f"},
				Usage:    "lambda function name",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "query, q",
				Value:    "",
				Aliases:  []string{"q"},
				Usage:    "query",
				Required: true,
			},
			&cli.Int64Flag{
				Name:     "limit, l",
				Value:    0,
				Aliases:  []string{"l"},
				Usage:    "limit",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "format",
				Usage:    "specify format [table, csv]",
				Required: false,
			},
			&cli.Int64Flag{
				Name:     "timeout, t",
				Aliases:  []string{"t"},
				Usage:    "max execution time seconds",
				Value:    60,
				Required: false,
			},
			&cli.StringFlag{
				Name:     "output, o",
				Aliases:  []string{"o"},
				Usage:    "output file path",
				Required: false,
			},
			&cli.BoolFlag{
				Name:     "verbose",
				Aliases:  []string{"v"},
				Usage:    "Show verbose logs",
				Required: false,
			},
		},
		Action: func(c *cli.Context) error {
			params := QueryParams{
				Function: c.String("function"),
				Query:    c.String("query"),
				Limit:    c.Int64("limit"),
			}

			res := Query(c, params)
			formatted := Format(FormatParam{Input: res.Result, Format: c.String("format")})

			if output := c.String("output"); output != "" {
				err := ioutil.WriteFile(output, []byte(formatted.Result), 0755)

				if err != nil {
					log.Fatal(err)
				}

				return nil
			}

			fmt.Println(formatted.Result)

			return nil
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}
