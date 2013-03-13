package gosim

/* A confidence interval associated with a value; it need not be
/* symmetric. */
type CI struct {
	// The [0, 1] confidence level
	Level float64

	// The width of the confidence interval
	Width float64
	
	// The lower and upper confidence bounds
	L float64
	U float64
}

/* A estimate of some value, and the associated confidence
/* interval. */
type Estimate struct {
	Value float64
	Bounds CI
}

	