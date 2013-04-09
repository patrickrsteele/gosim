package main

import (
	"fmt"
	"gosim/stats"
)

/* Print a t-table to compare to
/* http://www.sjsu.edu/faculty/gerstman/StatPrimer/t-table.pdf */
func main() {
	dfs := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
		11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
		21, 22, 23, 24, 25, 26, 27, 28, 29, 30,
		40, 60, 80, 100}

	cls := []float64{.5, .75, .8, .85, .9, .95, .975, .99, .995, .999, .9995}
	//cls := []float64{.999, .9995}

	for _, df := range dfs {
		icdf := stats.InvStudentCDF(df, 1e-5)
		fmt.Printf("%d  ", df)
		for _, cl := range cls {
			fmt.Printf("%.3f   ", icdf(cl))
		}
		fmt.Println()
	}
}
