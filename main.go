package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// Version provides lambda-query
var Version = "0.2.0"

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
				Required: false,
			},
			&cli.Int64Flag{
				Name:     "limit, l",
				Value:    0,
				Aliases:  []string{"l"},
				Usage:    "limit",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "input file, i",
				Value:    "",
				Aliases:  []string{"i", "inputfile"},
				Usage:    "e.g. lambda-query -i my_query.sql",
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
			if c.String("inputfile") == "" && c.String("query") == "" {
				log.Fatal("warning: `input file` name or `query` string is needed to run lambda-query")
			}

			if c.String("inputfile") != "" && c.String("query") != "" {
				log.Fatal("warning: either `input file` name or `query` string is allowed at once")
			}

			var query string
			var inputQuery []byte
			var err error

			if c.String("inputfile") != "" {
				inputQuery, err = ioutil.ReadFile(c.String("inputfile"))
				if err != nil {
					log.Fatal(err)
				}

				query = string(inputQuery)
			}

			if c.String("query") != "" {
				query = c.String("query")
			}

			params := QueryParams{
				Function:  c.String("function"),
				Query:     query,
				Limit:     c.Int64("limit"),
				InputFile: c.String("inputfile"),
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
