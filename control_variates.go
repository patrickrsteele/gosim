package gosim

/* A function that returns a value of intersest along with one control
/* variate. */
type TrialWithControl func() (float64, float64)

type MonteCarloWithControl struct {
	Simulation
	F TrialWithControl
}

func NewMonteCarloWithControl(f TrialWithControl) *MonteCarloWithControl {
	return &MonteCarloWithControl{F: f}
}

/* Simulate N trials and report a (1 - alpha)-level confidence interval */
func (m *MonteCarloWithControl) Simulate(N int, alpha float64) *Estimate {
	const (
		h = 1e-6
	)

	values := make([]float64, N)
	controls := make([]float64, N)
	for i, _ := range values {
		value[i], controls[i] = m.F()
	}

	value_mean, value_variance := Summary(values)
	control_mean, control_variance := Summary(controls)

	// Compute the correlation coefficient of the value of interest and
	// the control variate
	num := 0.0
	for i, _ := range values {
		num += (values[i] - value_mean) * (controls[i] - control_mean)
	}

	rho := num / ((float64(N) - 1) * math.Sqrt(value_variance*control_variance))
	variance := (1 - rho*rho) * value_variance

	coef := stats.InvStandardNormalCDF(h)(1 - alpha/2)
	coef *= math.Sqrt(variance / float64(N))

	return &Estimate{V: mean, C: &CI{Level: alpha, L: -coef, U: coef}}
}
