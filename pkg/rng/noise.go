package rng

var _ APINoise = (*Noise)(nil)

// APINoise is an interface that defines fetching remote data to use in generation
type APINoise interface {
	Call() *Noise
}

// Noise encapsulates this API call's contribution to the RNG
type Noise struct {
	Contribution int
}

// Call returns the noise generated from the API call
func (n *Noise) Call() *Noise {
	return n
}
