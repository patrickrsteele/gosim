package gosim

import (
	"gosim/stats"
	"math"
)

type Simulation interface {
	/* Simulate N trials and report a (1 - alpha)-level confidence
	/* interval */
	Simulate(N int, alpha float64) *Estimate
}

type Trial func() float64

type MonteCarloSimulation struct {
	Simulation
	F Trial
}

func NewMonteCarlo(f Trial) *MonteCarloSimulation {
	return &MonteCarloSimulation{F: f}
}

/* Simulate N trials and report a (1 - alpha)-level confidence interval */
func (m *MonteCarloSimulation) Simulate(N int, alpha float64) *Estimate {

	values := make([]float64, N)
	for i, _ := range values {
		values[i] = m.F()
	}

	return create_estimate(values, alpha)
}

func create_estimate(data []float64, alpha float64) *Estimate {
	const (
		h = 1e-6
	)

	mean, variance := Summary(data)

	coef := stats.InvStandardNormalCDF(h)(1 - alpha/2)
	coef *= math.Sqrt(variance / float64(len(data)))

	return &Estimate{V: mean, C: &CI{Level: alpha, L: -coef, U: coef}}
}
