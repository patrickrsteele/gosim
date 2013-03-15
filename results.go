package gosim

import (
	"fmt"
)

/* A confidence interval associated with a value; it need not be
/* symmetric. */
type CI struct {
	// The [0, 1] confidence level
	Level float64

	// The lower- and upper-deviations
	L float64
	U float64
}

/* A estimate of some value, and the associated confidence
/* interval. */
type Estimate struct {
	// The value
	V float64

	// The associated confidence interval
	C *CI
}

func (est *Estimate) String() string {
	return fmt.Sprintf("[%.3f - %.3f, %.3f + %.3f]",
		est.V, est.C.L, est.V, est.C.U)
}

/* Computes the mean and sample variance of the data provided. */
func Summary(data []float64) (float64, float64) {
	mean := 0.0
	for _, v := range data {
		mean += v
	}
	mean /= float64(len(data))

	variance := 0.0
	for _, v := range data {
		variance += (v - mean) * (v - mean)
	}
	variance /= float64(len(data) - 1)

	return mean, variance
}
