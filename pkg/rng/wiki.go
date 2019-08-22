package rng

import (
	"encoding/json"
	"net/http"
)

const wikiURL = "https://en.wikipedia.org/w/api.php?action=query&format=json&list=random&rnlimit=5"

// Wiki ...
type Wiki struct {
	client *http.Client
}

// NewWiki returns a new wiki fetcher
func NewWiki(c *http.Client) *Wiki {
	return &Wiki{
		client: c,
	}
}

// WikiQuery ...
type WikiQuery struct {
	Query *WikiRandom `json:"query"`
}

// WikiRandom ...
type WikiRandom struct {
	Random []*WikiEntry `json:"random"`
}

// WikiEntry ...
type WikiEntry struct {
	ID    int    `json:"id"`
	NS    int    `json:"ns"`
	Title string `json:"title"`
}

// WikiContribution ...
type WikiContribution struct {
	Lat float64
	Lon float64
	Sum int
}

// Wik does the work
func (w *Wiki) Wik() (*WikiContribution, error) {
	wi, err := w.apiRequest()
	if err != nil {
		return nil, err
	}
	wc := w.Process(wi)
	return wc, nil
}

// Process ...
func (w *Wiki) Process(chunk *WikiRandom) *WikiContribution {
	wc := &WikiContribution{}

	// So we need two floats
	lat1 := chunk.Random[0]
	lon1 := chunk.Random[1]

	wc.Lat = float64(lat1.ID) * 0.000001
	wc.Lon = float64(lon1.ID) * 0.000001

	var tempSum int
	var count int

	for entry := range chunk.Random {
		count++
		tempSum += chunk.Random[entry].NS
	}

	wc.Sum = tempSum / count

	return wc
}

func (w *Wiki) apiRequest() (*WikiRandom, error) {
	r, err := http.NewRequest("GET", wikiURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := w.client.Do(r)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	data := &WikiQuery{}
	err = json.NewDecoder(resp.Body).Decode(data)

	return data.Query, err
}
