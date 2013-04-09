package stats

import (
	"goint"
	"math"
)

/* Returns the pdf of a N(mu, sigma^2) random variable. */
func Normal(mu float64, sigma float64) func(float64) float64 {
	coef := 1.0 / (math.Sqrt(2*math.Pi) * sigma)
	sigmasq := sigma * sigma

	return func(x float64) float64 {
		return coef * math.Exp(-(x-mu)*(x-mu)/(2.0*sigmasq))
	}
}

/* Returns the cdf of a N(mu, sigma^2) random variable, accurate to
/* within h. */
func NormalCDF(mu, sigma, h float64) func(float64) float64 {
	pdf := Normal(mu, sigma)

	return func(x float64) float64 {
		// We know that this is a cdf with certain properties; exploit them
		// to avoid unnecessary (innacurate) numerical integration
		if x == mu {
			return .5
		} else if x > mu {
			return .5 + goint.Integrate(pdf, mu, x, h)
		}

		// We have that x < mu
		return .5 - goint.Integrate(pdf, x, mu, h)
	}
}

func InvStandardNormalCDF(h float64) func(float64) float64 {
	const perr = 1e-5  // Acceptable percentile error
	const berr = 1e-10 // Acceptable bounds error

	cdf := NormalCDF(0, 1, h)

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
			for cdf(U) < p {
				U *= 2
			}
		} else {
			U = 0
			L = -1
			for cdf(L) > p {
				L *= 2
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
