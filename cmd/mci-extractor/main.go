package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/daneharrigan/mci/scanner"
)

func main() {
	filenames := os.Args[1:]
	var results []scanner.Result

	for _, filename := range filenames {
		f, err := os.Open(filename)
		if err != nil {
			log.Printf("fn=Open error=%q", err)
			os.Exit(1)
		}
		defer f.Close()

		response := new(scanner.Response)
		if err := json.NewDecoder(f).Decode(&response); err != nil {
			log.Printf("fn=Decode error=%q", err)
			os.Exit(1)
		}

		for _, r := range response.Data.Results {
			t, err := scanner.FindUnlimitedDate(r.Dates)
			if err != nil || t.IsZero() {
				continue
			}

			if contains(results, r) {
				continue
			}

			results = append(results, r)
		}
	}

	if err := json.NewEncoder(os.Stdout).Encode(results); err != nil {
		log.Printf("fn=Encode error=%q", err)
		os.Exit(1)
	}
}

func contains(results []scanner.Result, result scanner.Result) bool {
	for _, r := range results {
		if r.Title == result.Title && r.Series.Name == result.Series.Name {
			return true
		}
	}

	return false
}
