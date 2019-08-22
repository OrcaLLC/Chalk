package rng

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
)

// Sunterfuge takes two given floats and returns an RN
// based on the sunrise/sunset for that lat/long
type Sunterfuge struct {
	Lat    float64 // +/- 90
	Lon    float64 // +/- 180
	client *http.Client
}

// SunResults ...
type SunResults struct {
	Data   *SunData `json:"results"`
	Status string   `json:"status"`
}

// SunData ...
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

// SunCont ...
type SunCont struct {
	Modifiers []int
}

// NewSunterfuge returns a new NewSunterfuge object
func NewSunterfuge(c *http.Client, lat float64, lon float64) *Sunterfuge {
	return &Sunterfuge{
		Lat:    validLat(lat),
		Lon:    validLon(lon),
		client: c,
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
func (s *Sunterfuge) Determine() (*SunCont, error) {
	r, rerr := s.apiRequest()
	if rerr != nil {
		return nil, rerr
	}
	sc := s.Process(r)

	return sc, nil
}

// Process ...
func (s *Sunterfuge) Process(dat *SunData) *SunCont {
	var int1, int2, int3 int
	int1 += extractNumFromDate(dat.Sunrise, 0)
	int1 += extractNumFromDate(dat.Sunset, 0)
	int1 += extractNumFromDate(dat.SolarNoon, 0)
	int1 += extractNumFromDate(dat.DayLength, 0)
	int1 += extractNumFromDate(dat.CivilTwilightBegin, 0)
	int1 += extractNumFromDate(dat.CivilTwilightEnd, 0)
	int1 += extractNumFromDate(dat.NautTwiBegin, 0)
	int1 += extractNumFromDate(dat.NautTwiEnd, 0)
	int1 += extractNumFromDate(dat.AstroTwiBegin, 0)
	int1 += extractNumFromDate(dat.AstroTwiEnd, 0)
	int1 = int1 / 10
	int1 = reduceInteger(int1)

	int2 += extractNumFromDate(dat.Sunrise, 1)
	int2 += extractNumFromDate(dat.Sunset, 1)
	int2 += extractNumFromDate(dat.SolarNoon, 1)
	int2 += extractNumFromDate(dat.DayLength, 1)
	int2 += extractNumFromDate(dat.CivilTwilightBegin, 1)
	int2 += extractNumFromDate(dat.CivilTwilightEnd, 1)
	int2 += extractNumFromDate(dat.NautTwiBegin, 1)
	int2 += extractNumFromDate(dat.NautTwiEnd, 1)
	int2 += extractNumFromDate(dat.AstroTwiBegin, 1)
	int2 += extractNumFromDate(dat.AstroTwiEnd, 1)
	int2 = int2 / 10
	int2 = reduceInteger(int2)

	int3 += extractNumFromDate(dat.Sunrise, 2)
	int3 += extractNumFromDate(dat.Sunset, 2)
	int3 += extractNumFromDate(dat.SolarNoon, 2)
	int3 += extractNumFromDate(dat.DayLength, 2)
	int3 += extractNumFromDate(dat.CivilTwilightBegin, 2)
	int3 += extractNumFromDate(dat.CivilTwilightEnd, 2)
	int3 += extractNumFromDate(dat.NautTwiBegin, 2)
	int3 += extractNumFromDate(dat.NautTwiEnd, 2)
	int3 += extractNumFromDate(dat.AstroTwiBegin, 2)
	int3 += extractNumFromDate(dat.AstroTwiEnd, 2)
	int3 = int3 / 10
	int3 = reduceInteger(int3)

	sc := &SunCont{}
	sc.Modifiers = append(sc.Modifiers, int1)
	sc.Modifiers = append(sc.Modifiers, int2)
	sc.Modifiers = append(sc.Modifiers, int3)

	return sc
}

// extractNumFromDate returns field int from date after split on ":"
func extractNumFromDate(date string, field int) int {
	cleaned := strings.Split(date, " ")[0]
	split := strings.Split(cleaned, ":")
	num, err := strconv.Atoi(split[field])
	if err != nil {
		// fuck it right?
		return 1
	}

	return num
}

func (s *Sunterfuge) formatRequest() string {
	return fmt.Sprintf("https://api.sunrise-sunset.org/json?lat=%.7f&lng=-%.7f&date=today", s.Lat, s.Lon)
}

func (s *Sunterfuge) apiRequest() (*SunData, error) {
	r, err := http.NewRequest("GET", s.formatRequest(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(r)
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
