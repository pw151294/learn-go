package main

import (
	bytes2 "bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
	searchUrl = "_search"
	addUrl    = "_doc"
	esUrl     = "http://127.0.0.1:9400"
)

type ReqSearchData struct {
	Hits HitsData `json:"hits"`
}

type HitsData struct {
	Total TotalData  `json:"total"`
	Hits  []HitsData `json:"hits"`
}

type TotalData struct {
	value    int
	Relation string
}
type HitsTwoData struct {
	Source json.RawMessage `json:"_source"`
}

// EsSearch ES的搜索功能
func EsSearch(indexName string, query map[string]interface{}, from, size int, sort []map[string]string) HitsData {
	params := map[string]interface{}{
		"query": query,
		"from":  from,
		"size":  size,
		"sort":  sort,
	}
	bytes, err := json.Marshal(params)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s/%s", esUrl, indexName, searchUrl), bytes2.NewBuffer(bytes))
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		log.Printf("invalid status code: %d", resp.StatusCode)
	}
	bytes, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	data := HitsData{}
	if err := json.Unmarshal(bytes, &data); err != nil {
		log.Fatal(err)
	}
	return data
}

func main() {

}
