package rng

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"

	"github.com/davecgh/go-spew/spew"
)

// Sunterfuge takes two given floats and returns an RN
// based on the sunrise/sunset for that lat/long
type Sunterfuge struct {
	Lat float64 // +/- 90
	Lon float64 // +/- 180
}

type SunResults struct {
	Data   *SunData `json:"results"`
	Status string   `json:"status"`
}

type SunData struct {
	Sunrise            string `json:"sunrise"`
	Sunset             string `json:"sunset"`
	SolarNoon          string `json:"solar_noon"`
	DayLength          string `json:"day_length"`
	CivilTwilightBegin string `json:"civil_twilight_begin"`
	CivilTwilightEnd   string `json:"civil_twilight_end"`
	NautTwiBegin       string `json:"nautical_twilight_begin"`
	NautTwiEnd         string `json:"nautical_twilight_end"`
	AstroTwiBegin      string `json:"astronomical_twilight_begin"`
	AstroTwiEnd        string `json:"astronomical_twilight_end"`
}

// NewSunterfuge returns a new NewSunterfuge object
func NewSunterfuge(lat float64, lon float64) *Sunterfuge {
	return &Sunterfuge{
		Lat: validLat(lat),
		Lon: validLon(lon),
	}
}

func validLat(lat float64) float64 {
	for lat > 90 || lat < -90 {
		lat = lat / math.Sqrt(lat)
	}
	return lat
}

func validLon(lon float64) float64 {
	for lon > 180 || lon < -180 {
		lon = lon / math.Sqrt(lon)
	}
	return lon
}

// Determine does the thing
func (s *Sunterfuge) Determine() error {
	r, rerr := s.apiRequest()
	if rerr != nil {
		return rerr
	}

	spew.Dump(r)
	return nil
}

func (s *Sunterfuge) formatRequest() string {
	return fmt.Sprintf("https://api.sunrise-sunset.org/json?lat=%.7f&lng=-%.7f&date=today", s.Lat, s.Lon)
}

func (s *Sunterfuge) apiRequest() (*SunData, error) {
	c := http.DefaultClient

	r, err := http.NewRequest("GET", s.formatRequest(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(r)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	data := &SunResults{}
	err = json.NewDecoder(resp.Body).Decode(data)

	if data.Status != "OK" {
		return nil, fmt.Errorf(data.Status)
	}

	return data.Data, err
}
