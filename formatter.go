package main

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"strings"

	"github.com/jedib0t/go-pretty/table"
	log "github.com/sirupsen/logrus"
)

type FormatParam struct {
	Input  string
	Format string
}

type FormattedResult struct {
	Result string
}

func Format(param FormatParam) FormattedResult {
	if param.Format == "raw" {
		return FormattedResult{Result: param.Input}
	}

	if param.Format == "csv" {
		return FormattedResult{Result: param.Input}
	}

	if param.Format == "json" {
		r := csv.NewReader(strings.NewReader(param.Input))

		var records []map[string]string
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

			record := make(map[string]string)

			for i, v := range row {
				record[header[i]] = v
			}

			records = append(records, record)
		}

		res, err := json.MarshalIndent(records, "", "  ")

		if err != nil {
			log.Fatal(err)
		}

		return FormattedResult{Result: string(res)}
	}

	if param.Format == "table" {
		t := table.NewWriter()
		r := csv.NewReader(strings.NewReader(param.Input))

		for {
			row, err := r.Read()

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatal(err)
			}

			rt := make(table.Row, len(row))
			for i, v := range row {
				rt[i] = v
			}

			t.AppendRow(rt)
		}

		res := t.Render()

		return FormattedResult{Result: res}
	}

	return FormattedResult{Result: param.Input}
}
