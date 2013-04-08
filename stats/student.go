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

/* Returns the cdf of the Student's t-distribution with d degrees of
/* freedom. The cdf requires a second argument h as the step size of
/* the integration used to compute the cdf. */
func StudentCDF(d int, h float64) func(float64) float64 {
	pdf := Student(d)

	return func(x float64) float64 {
		// We know that this is a cdf with certain properties; exploit them
		// to avoid unnecessary (innacurate) numerical integration
		if x == 0 {
			return .5
		} else if x > 0 {
			return .5 + goint.Integrate(pdf, 0, x, h)
		}

		// We have that x < 0
		return .5 - goint.Integrate(pdf, x, 0, h)
	}
}

/* Returns an inverse of the function returned by StudentCDF with
/* identical argumentss. */
func InvStudentCDF(d int, h float64) func(float64) float64 {
	const perr = 1e-5  // Acceptable percentile error
	const berr = 1e-10 // Acceptable bounds error

	cdf := StudentCDF(d, h)

	return func(p float64) float64 {

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
			for StudentCDF(d, hh)(U) < p {
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
			for StudentCDF(d, hh)(L) > p {
				L *= 2
				hh *= 2
				if hh > 1 {
					hh = 1
				}
			}
		}

		// Start bisection; halt when we're close enough to the proper percentile, or no progress is made
		pL := cdf(L)
		pU := cdf(U)
		for pU-pL > 2*perr && U-L > berr {
			M := (L + U) / 2
			pM := cdf(M)

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
}
