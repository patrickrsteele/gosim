package gosim

import (
	"math"
	"math/rand"
	"sort"
)

type Process interface {
	At(float64) float64
}

type Brownian struct {
	Process

	// A slice of (t, v) pairs, where t is the time and v is the value
	// of the Brownian motion
	states []brownianState

	Rand *rand.Rand
}

func NewBrownian(rand *rand.Rand) *Brownian {
	return &Brownian{Rand: rand}
}

type brownianState struct {
	t float64
	v float64
}

type brownianStateSlice []brownianState

func (b *Brownian) At(t float64) float64 {
	// Make sure (0, 0) is in the points
	if len(b.states) == 0 {
		b.states = append(b.states, brownianState{t: 0, v: 0})
	}

	// We'll eventually add this state if time t hasn't been simulated
	new_state := brownianState{t: t}

	// Perform a binary search to find the nearest indices to t
	L := 0
	U := len(b.states) - 1
	for L+1 < U {
		M := (L + U) / 2

		if b.states[M].t <= t {
			L = M
		} else {
			U = M
		}
	}

	if b.states[U].t < t {
		// The process hasn't been simulated to this point; extend it
		std_dev := math.Sqrt(t - b.states[U].t)
		delta := b.Rand.NormFloat64() * std_dev
		new_state.v = b.states[U].v + delta
	} else if b.states[L].t == t || b.states[U].t == t {
		// We've simulated this time before, simply return it
		return b.states[L].v
	} else {
		// Use a Brownian bridge construction to compute the value at t
		s0 := b.states[L].t
		v0 := b.states[L].v
		s1 := b.states[U].t
		v1 := b.states[U].v

		mean := (s1-t)/(s1-s0)*v0 + (t-s0)/(s1-s0)*v1
		variance := (s1 - t) * (t - s0) / (s1 - s0)
		std_dev := math.Sqrt(variance)

		new_state.v = mean + b.Rand.NormFloat64()*std_dev
	}

	// Add the new state to the slice of states, in sorted order
	b.states = append(b.states, new_state)
	sort.Sort(brownianStateSlice(b.states))

	// Finally, return the actual value
	return new_state.v
}

// Implement the sort.Interface interface method Len
func (b brownianStateSlice) Len() int {
	return len(b)
}

// Implement the sort.Interface interface method Less
func (b brownianStateSlice) Less(i, j int) bool {
	return b[i].t <= b[j].t
}

// Implement the sort.Interface interface method Swap
func (b brownianStateSlice) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

type GeometricBrownian struct {
	// Inherit from standard Brownian Motion
	Brownian

	// The drift of the process
	Drift float64

	// The volatility of the process
	Volatility float64

	// The square root of the volatility
	std_dev float64

	// The scaling coefficient
	Scale float64
}

func NewGeometricBrownian(rand *rand.Rand, scale, drift, volatility float64) *GeometricBrownian {
	std_dev := math.Sqrt(volatility)
	return &GeometricBrownian{Brownian: Brownian{Rand: rand},
		Drift: drift, Volatility: volatility, std_dev: std_dev, Scale: scale}
}

func (b *GeometricBrownian) At(t float64) float64 {
	v := b.Brownian.At(t)

	return b.Scale * math.Exp((b.Drift-b.Volatility/2.0)*t+b.std_dev*v)
}
