package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/urfave/cli/v2"
)

type queryRequest struct {
	Query string `json:"query"`
}

type responseBody struct {
	Result string `json:"result"`
}

type responseHeaders struct {
	ContentType string `json:"Content-Type"`
}

type responseError struct {
	Message string `json:"message"`
}

type QueryParams struct {
	Function  string
	Query     string
	Limit     int64
	InputFile string
}

type QueryResponse struct {
	Result string
}

func Query(c *cli.Context, q QueryParams) QueryResponse {
	ctx, cancel := context.WithTimeout(c.Context, time.Second*time.Duration(c.Int64("timeout")))
	defer cancel()

	sess := session.Must(session.NewSessionWithOptions(session.Options{SharedConfigState: session.SharedConfigEnable}))

	client := lambda.New(sess)

	if q.Limit > 0 {
		var offset int64
		var header string
		var temp string
		var res string

		for {
			query := fmt.Sprintf("%s LIMIT %v OFFSET %v", q.Query, q.Limit, offset)

			result := invokeRequest(client, query, q.Function)

			if result.Result == "" {
				break
			}

			if result.Result == "OK" {
				break
			}

			if result.Result == "None" {
				break
			}

			if result.Result == "No record found" {
				break
			}

			if result.Result != "" {
				res = getRecordsWithoutHeader(result.Result)
				temp = temp + res

				if header == "" {
					header = getHeaderFromRecords(result.Result)
				}
			}

			offset += q.Limit

			if ctx.Err() == context.DeadlineExceeded {
				log.Fatal("error: query timed out")
			}
		}

		temp = header + temp

		return QueryResponse{Result: temp}
	}

	return invokeRequest(client, q.Query, q.Function)
}

func invokeRequest(client *lambda.Lambda, query string, function string) QueryResponse {
	request := queryRequest{Query: query}

	payload, err := json.Marshal(request)

	if err != nil {
		fmt.Println("Error encoding request payload")
		os.Exit(0)
	}

	result, err := client.Invoke(&lambda.InvokeInput{FunctionName: aws.String(function), Payload: payload})
	if err != nil {
		log.Fatal(err)
		os.Exit(0)
	}

	if v := result.FunctionError; v != nil {
		log.Fatalf("%s", result.Payload)
	}

	var resp responseBody

	err = json.Unmarshal(result.Payload, &resp)

	if err != nil {
		fmt.Println("Error decoding response")
		os.Exit(0)
	}

	return QueryResponse{Result: resp.Result}
}

func getRecordsWithoutHeader(input string) string {
	r := csv.NewReader(strings.NewReader(input))

	var records [][]string
	var header []string

	for {
		row, err := r.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		if header == nil {
			header = row
			continue
		}

		records = append(records, row)
	}

	b := &bytes.Buffer{}
	w := csv.NewWriter(b)

	for _, v := range records {
		w.Write(v)
	}

	w.Flush()

	return b.String()
}

func getHeaderFromRecords(input string) string {
	r := csv.NewReader(strings.NewReader(input))

	var header []string

	for {
		row, err := r.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		if header == nil {
			header = row
			break
		}
	}

	return strings.Join(header, ",") + "\n"
}
