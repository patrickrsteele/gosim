package gosim

import (
	"math/rand"
)

// Returns an exponential random variable with rate lambda.
func Exponential(r *rand.Rand, lambda float64) float64 {
	return r.ExpFloat64() / lambda
}

// Simulates a Poisson process with rate lambda until the first event
// occurs at or after time T.
func PoissonProcess(r *rand.Rand, lambda float64, T float64) []float64 {
	process := make([]float64, 0)

	t := 0.0
	for t < T {
		step := Exponential(r, lambda)
		t += step
		process = append(process, t)
	}

	return process
}
