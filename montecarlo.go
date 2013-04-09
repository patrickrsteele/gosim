package gosim

import (
	"github.com/skelterjohn/go.matrix"
	"gosim/stats"
	"math"
)

type Simulation interface {
	/* Simulate N trials and report a (1 - alpha)-level confidence
	/* interval */
	Simulate(N int, alpha float64) *Estimate
}

/* Return an estimate of a parameter along with a slice of control variates;
/* the slice can be nil. */
type Trial func() (float64, []float64)

type MonteCarloSimulation struct {
	Simulation
	F Trial
}

func NewMonteCarlo(f Trial) *MonteCarloSimulation {
	return &MonteCarloSimulation{F: f}
}

/* Simulate N trials and report a (1 - alpha)-level confidence interval */
func (m *MonteCarloSimulation) Simulate(N int, alpha float64) *Estimate {
	// We require at least one trial
	if N <= 1 {
		panic("N must be greater than 1")
	}

	n := float64(N)

	// Run all trials, storing the results in ys and the controls in xs
	ys := make([]float64, N)
	xs := make([][]float64, N)
	for i, _ := range ys {
		ys[i], xs[i] = m.F()
	}

	d := len(xs[0])
	// If there are no control variates, just do standard Monte Carlo
	if d == 0 {
		return create_estimate(ys, alpha)
	}

	// Otherwise return an estimate based off the control variates
	return control_estimate(ys, xs, alpha)
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

/* Return a Monte Carlo estimate of the parameter, using the data in
/* xs[i] as a control variate for observation ys[i]. */
func control_estimate(ys []float64, xs [][]float64) *Estimate {
	d := len(xs[0])

	// Compute the means of the control variates
	x_means := make([]float64, d)
	for _, X := range xs {
		for j, _ := range x_means {
			x_means[j] += X[j]
		}
	}
	for j, _ := range x_means {
		x_means[j] /= n
	}

	// Compute the mean of the target
	y_mean := Mean(ys)

	// Compute the sample covariance of the control variates
	Sx := matrix.MakeDenseMatrix(make([]float64, d*d), d, d)
	for j := 0; j < d; j++ {
		for k := 0; k < d; k++ {
			entry := 0.0
			for i, _ := range xs {
				entry += xs[i][j]*xs[i][k] - n*x_means[j]*x_means[k]
			}
			entry /= n - 1

			Sx.Set(j, k, entry)
		}
	}

	// Compute the sample covariance of the control variates and the
	// response
	Sxy := matrix.MakeDenseMatrix(make([]float64, d), d, 1)
	for j := 0; j < d; j++ {
		entry := 0.0
		for i, _ := range ys {
			entry += xs[i][j]*ys[i] - n*x_means[j]*y_mean
		}
		entry /= n - 1

		Sxy.Set(j, 0, entry)
	}

	// Compute the control variate weights
	b := matrix.Product(matrix.Inverse(Sx), Sxy)

	// Compute the offset estimates
	yoff := make([]float64, N)
	for i, _ := range yoff {
		yoff[i] = ys[i]

		for j := 0; j < d; j++ {
			yoff[j] -= b.Get(j, 0) * (xs[i][j] - x_means[j])
		}
	}

	// Compute the point estimate and sample variance of the offset
	// values
	return create_estimate(yoff, alpha)
}
