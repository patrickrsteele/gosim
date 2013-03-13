package gosim

/* A useful collection of simulation objects, such as queues. */

import (
	"bytes"
	"fmt"
	"math/rand"
)

type Queue interface {
	// Simulate the Queue until time T
	Until(T float64)

	// Simulate until the next event; store and return the resulting state
	Advance() *QueueState

	// The states of the queue at the event times
	States() []QueueState
}

/* The number of customers N in service in a Queue at time T */
type QueueState struct {
	T float64
	N int
}

func (qs *QueueState) String() string {
	return fmt.Sprintf("[%.2f, %d]", qs.T, qs.N)
}

/* A Queue with Poisson arrivals, exponential service times, and c
/* servers. */
type MMCQueue struct {
	Rand    *rand.Rand
	Arrival float64
	Service float64
	C       int
	states  []QueueState
}

func NewMMCQueue(r *rand.Rand, arrival float64, service float64, c int) *MMCQueue {
	queue := &MMCQueue{Arrival: arrival, Service: service, C: c, Rand: r}

	// Start with [0, 0] as a state
	queue.states = []QueueState{QueueState{T: 0, N: 0}}

	return queue
}

func (q *MMCQueue) Advance() *QueueState {
	t := Exponential(q.Rand, q.Arrival)
	arrival := true

	cur_state := q.states[len(q.states)-1]

	// There are only C servers
	in_service := cur_state.N
	if in_service > q.C {
		in_service = q.C
	}

	// See if any service completes before the next arrival
	for i := 0; i < in_service; i++ {
		if s := Exponential(q.Rand, q.Service); s < t {
			t = s
			arrival = false
		}
	}

	new_state := QueueState{T: cur_state.T + t, N: cur_state.N}
	if arrival {
		new_state.N += 1
	} else {
		new_state.N -= 1
	}

	q.states = append(q.states, new_state)
	return &new_state
}

func (q *MMCQueue) Until(T float64) {
	t := 0.0
	if len(q.states) > 0 {
		t = q.states[len(q.states)-1].T
	}

	for t < T {
		state := q.Advance()
		t = state.T
	}
}

func (q *MMCQueue) String() string {
	var buffer bytes.Buffer

	for i, qs := range q.states {
		if i > 0 {
			buffer.WriteString(", ")
		}
		buffer.WriteString(qs.String())
	}

	return buffer.String()
}

func (q *MMCQueue) States() []QueueState {
	return q.states
}
