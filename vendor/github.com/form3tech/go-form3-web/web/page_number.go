package web

import "math"

func lastPageNumber(totalRecords int, pageSize int) int {
	d := float64(totalRecords) / float64(pageSize)
	if d == 0 {
		return 0
	}
	return int(math.Ceil(d)) - 1
}
