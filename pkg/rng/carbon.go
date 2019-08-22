package rng

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
)

const carbonHostname = "https://api.carbonintensity.org.uk/intensity"

var _ APINoise = (*Carbon)(nil)

// Carbon encapsulates an API call to get the current carbon intensity in England
type Carbon struct {
	client *http.Client
}

// NewCarbon ...
func NewCarbon(c *http.Client) *Carbon {
	return &Carbon{
		client: c,
	}
}

// CarbonData is just a wrapper
type CarbonData struct {
	Data []*CarbonIntensity `json:"data"`
}

// CarbonIntensity encapsulates the data response payload
type CarbonIntensity struct {
	Intensity *CarbonResponse `json:"intensity"`
}

// CarbonResponse encapsulates the data from UK carbon
type CarbonResponse struct {
	Forecast int    `json:"forecast"`
	Actual   int    `json:"actual"`
	Index    string `json:"index"`
}

// Call returns the noise generated from the API call
func (c *Carbon) Call() *Noise {
	n := &Noise{}
	c.client = http.DefaultClient

	r, rerr := c.FetchCarbon()
	if rerr != nil {
		return nil
	}

	n.Contribution = simpleReduceInteger(c.ProcessCarbon(r))

	return n
}

// FetchCarbon returns the Carbon API data
func (c *Carbon) FetchCarbon() (*CarbonResponse, error) {
	cr := &CarbonData{}

	r, err := http.NewRequest("GET", carbonHostname, nil)
	if err != nil {
		return nil, err
	}

	r.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(r)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(cr)
	return cr.Data[0].Intensity, err

}

// ProcessCarbon turns the returned data into a contribution
func (c *Carbon) ProcessCarbon(cr *CarbonResponse) int {
	//var result int

	indexed, err := turnIndexIntoNumber(cr.Index)
	if err != nil {
		fmt.Println(err)
		return 0
	}

	merged := int(math.Abs(float64(cr.Forecast - cr.Actual)))

	final := int(math.Abs(float64(indexed + merged)))
	return final
}

// turnIndexIntoNumber returns an int representing a convoluted processing of index
func turnIndexIntoNumber(index string) (int, error) {
	var combined string

	for i := 0; i < len(index); i++ {
		combined = fmt.Sprintf("%s%x", combined, index[i])
	}

	n, err := strconv.ParseUint(combined, 16, 32)
	if err != nil {
		return 0, err
	}
	n2 := uint64(n)
	n2f := math.Float64frombits(n2)

	f1 := math.Abs(math.Log(n2f))
	f := int(simpleReduceIndex(f1))

	return f, nil
}
