package main

import (
	"encoding/csv"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
)

func fetchDictFromFile() map[string][]string {

	path, _ := filepath.Abs("Files/Lughat.csv")
	file, _ := os.Open(path)
	reader := csv.NewReader(file)
	record, err := reader.ReadAll()
	dict := map[string][]string{}
	if err == nil {

		for i := range record {
			dict[record[i][0]] = append(dict[record[i][0]], record[i][1])
		}

	}
	return dict
}
