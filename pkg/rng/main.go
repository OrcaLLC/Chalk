// Package rng generates random numbers
package rng

import (
	"fmt"
	"net/http"
)

// Talker is the client wrapper for all of RNG -I'm sure it's fine
type Talker struct {
	client *http.Client
}

// New returns a random number between 0 and 9
func New() int {
	t := &Talker{
		client: http.DefaultClient,
	}

	c := NewCarbon(t.client)
	// Carbon noise
	cn := c.Call()
	//fmt.Printf("Carbon: %v\n", cn.Contribution)

	w := NewWiki(t.client)
	wc, err := w.Wik()
	if err != nil {
		fmt.Println(err)
		return 0
	}

	st := NewSunterfuge(t.client, wc.Lat, wc.Lon)
	sc, err := st.Determine()
	if err != nil {
		fmt.Println(err)
		return 0
	}

	//spew.Dump(sc)

	crazy1 := sc.Modifiers[0] + cn.Contribution
	crazy1 = crazy1 / 3

	crazy2 := sc.Modifiers[1] + cn.Contribution
	crazy2 = crazy2 / 3

	crazy3 := sc.Modifiers[1] + cn.Contribution
	crazy3 = crazy3 / 3

	avg1 := average(crazy1, cn.Contribution)
	avg2 := average(crazy2, cn.Contribution)
	avg3 := average(crazy3, cn.Contribution)

	final1 := reduceInteger((avg1 + avg2/avg3) + (avg1 / avg3) + (avg2 / avg3))
	fmt.Println(final1)

	return 0
}

/*
Add These

https://en.wikipedia.org/w/api.php?action=query&format=json&list=random&rnlimit=5

use to generate lat/long
*/
