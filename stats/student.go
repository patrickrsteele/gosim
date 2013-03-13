package stats

import (
	"goint"
	"math"
)

/* Returns a pdf of Student's t-distribution with d degrees of
/* freedom. */
func Student(d int) func(float64) float64 {
	nu := float64(d)
	coef := math.Gamma((nu+1)/2) / (math.Gamma(nu/2) * math.Sqrt(nu*math.Pi))
	exp := -(float64(nu + 1)) / 2

	return func(x float64) float64 {
		return coef * math.Pow(1+x*x/nu, exp)
	}
}

/* The cdf of the Student's t-distribution with d degrees of
/* freedom. Compute using step size h. */
func StudentCDF(d int, x float64, h float64) float64 {
	// We know that this is a cdf with certain properties; exploit them
	// to avoid unnecessary (innacurate) numerical integration
	if x == 0 {
		return .5
	} else if x > 0 {
		return .5 + goint.SimpsonIntegration(Student(d), 0, x, h)
	}

	// We have that x < 0
	return .5 - goint.SimpsonIntegration(Student(d), x, 0, h)
}

/* Returns x such that F(x) = p, where F is the CDF of the Student's
/* t-distribution with d degrees of freedom */
func InvStudentCDF(d int, p float64, h float64) float64 {
	const perr = 1e-5  // Acceptable percentile error
	const berr = 1e-10 // Acceptable bounds error

	var L, U float64

	// Exploit known properties to avoid bisection, if we can
	if p == 0 {
		return math.Inf(-1)
	} else if p == 1 {
		return math.Inf(1)
	} else if p == .5 {
		return 0
	}

	// Choose bisection bounds somewhat intelligently
	if p > .5 {
		L = 0
		U = 1
		hh := h
		for StudentCDF(d, U, hh) < p {
			U *= 2
			hh *= 2
			if hh > 1 {
				hh = 1
			}
		}
	} else {
		U = 0
		L = -1
		hh := h
		for StudentCDF(d, L, hh) > p {
			L *= 2
			hh *= 2
			if hh > 1 {
				hh = 1
			}
		}
	}

	// Start bisection; halt when we're close enough to the proper percentile, or no progress is made
	pL := StudentCDF(d, L, h)
	pU := StudentCDF(d, U, h)
	for pU-pL > 2*perr && U-L > berr {
		M := (L + U) / 2
		pM := StudentCDF(d, M, h)

		if pM > p {
			U = M
			pU = pM
		} else if pM < p {
			L = M
			pL = pM
		} else {
			return L
		}
	}

	return (L + U) / 2
}
