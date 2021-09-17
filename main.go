package main

import (
	"fmt"
)

func Sqrt(x float64) float64 {
		z := 1.0
	for i := 0; i < 10; i++ {
		z -= (z*z - x) / (2 * z)
	}

	return z

}

func main() {

	/*res := math.Sqrt(4)
	fmt.Println(res)*/
	fmt.Println(Sqrt(36))
}

