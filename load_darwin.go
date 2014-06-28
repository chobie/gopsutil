// +build darwin
package gopsutil

/*
#include <stdlib.h>
 */
import "C"

func LoadAvg() (*LoadAvgStat, error) {
	avg := []C.double{0, 0, 0}

	C.getloadavg(&avg[0], C.int(len(avg)))

	ret := &LoadAvgStat{
		Load1:  float64(avg[0]),
		Load5:  float64(avg[1]),
		Load15: float64(avg[2]),
	}

	return ret, nil
}
