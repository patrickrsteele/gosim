package gosim

import (
	"math"
	"math/rand"
	"sort"
)

type Brownian struct {
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

// Implement the sort.Interface interface
func (b *Brownian) Len() int {
	return len(b.states)
}

func (b *Brownian) Less(i, j int) bool {
	return b.states[i].t <= b.states[j].t
}

func (b *Brownian) Swap(i, j int) {
	b.states[i], b.states[j] = b.states[j], b.states[i]
}

type Process interface {
	At(float64) float64
}

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
	} else if L == U {
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
	sort.Sort(b)

	// Finally, return the actual value
	return new_state.v
}
