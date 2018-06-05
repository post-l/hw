package main

import (
	"fmt"
	"time"
)

func main() {
	for j := 0; j < 3; j++ {
		// Determine time.Now() monotonic resolution.
		if j > 0 {
			fmt.Println("time.Sleep(3ms)")
			time.Sleep(3 * time.Millisecond)
		}
		for i := 0; i < 5; i++ {
			start := time.Now()
			stop := time.Now()
			fmt.Printf("monotonic resolution: %v\n", stop.Sub(start))
		}
	}
	v := 0
	size := 1000
	for i := 0; i < size; i++ {
		start := time.Now()
		stop := time.Now()
		v += int(stop.Sub(start))
	}
	v /= size
	fmt.Printf("\n\n\nmonotonic resolution avg:  %v\n\n\n\n", time.Duration(v))

	vs := make([]int, 11)
	for x := uint(0); x < 11; x++ {
		vs[x] = 130 << uint(x)
	}

	vreal := []int{
		53,    // 130
		120,   // 260
		250,   // 520
		510,   // 1040
		1000,  // 2080
		2100,  // 4160
		4800,  // 8320
		10000, // 16640
		30000, // 33280
		60000, // 66560
		37000, // 133120
	}

	vs2 := make([]int, len(vs))
	size = 5000
	for x, v := range vs {
		if v <= 66560 {
			continue
		}
		start := time.Now()
		d := time.Duration(vreal[x])
		for i := 0; i < size; i++ {
			time.Sleep(d)
		}
		vs2[x] = int(time.Since(start)) / size
	}

	vs3 := make([]int, len(vs))
	for x, v := range vs {
		if v > 66560 {
			continue
		}
		start := time.Now()
		for i := 0; i < size; i++ {
			for j := vreal[x]; j != 0; j-- {
			}
		}
		vs3[x] = int(time.Since(start)) / size
	}

	fmt.Printf("Wanted     Got\n")
	for i := range vs {
		fmt.Printf("%6d\t%6d\t%6d\n", vs[i], vs2[i], vs3[i])
	}
}
