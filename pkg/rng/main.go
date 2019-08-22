// Package rng generates random numbers
package rng

import "fmt"

// New returns a random number between 0 and 9
func New() int {
	//g := &Generator{}
	c := &Carbon{}
	n := c.Call()

	fmt.Println(n)

	return 0
}
