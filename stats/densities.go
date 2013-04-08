package stats

import (
	"goint"
	"math"
)

/* These functions provide a general framework for manipulating pdfs,
/* cdfs, and inverse cdfs. */

type Function func(float64) float64

/* Compute a cdf by integrating a pdf with step size h */
func CDFFromPDF(pdf Function, h float64) Function {
	return func(x float64) float64 {
		return goint.Integrate((func(float64) float64)(pdf), math.Inf(-1), x, h)
	}
}

/* Returns a function that is the inverse of the given cdf, computed
/* with a step size of h. */
func InvCDFFromCDF(cdf Function, h float64) Function {
	const perr = 1e-5  // Acceptable percentile error
	const berr = 1e-10 // Acceptable bounds error

	return func(p float64) float64 {

		L := -1.0
		U := 1.0

		for cdf(U) < p {
			U *= 2
		}

		for cdf(L) > p {
			L *= 2
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
