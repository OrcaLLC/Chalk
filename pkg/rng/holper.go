package rng

import "math"

const maxInt = 9
const minInt = 0
const maxFlt = 9.00
const minFlt = 0.00

func reduceInteger(i int) int {
	for i > maxInt || i < minInt {
		f := float64(i) / math.Sqrt(float64(i))
		//f = math.Abs(math.Log(f))
		i = int(f)
	}
	return i
}

func simpleReduceInteger(i int) int {
	for i > maxInt || i < minInt {
		i = i / 2
	}
	return i
}

func simpleReduceIndex(f float64) float64 {
	for f > maxFlt || f < minFlt {
		f = f / 2
	}
	return f
}

func average(in ...int) int {
	var sum, count int
	for i := range in {
		sum += in[i]
		count++
	}
	final := sum / count
	return final
}
