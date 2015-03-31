package main

import (
    "math"
    "fmt"
)

type ErrNegativeSqrt float64

func (e ErrNegativeSqrt) Error() string {
    return fmt.Sprintf("cannot Sqrt negative number: %v", float64(e))
}

func Sqrt3(x float64) (float64, error) {
    if x < 0 {
        return 0, ErrNegativeSqrt(x)
    }
    delta := 1e-6
    z := x
    n := 0.0
    for math.Abs(n-z) > delta {
        n, z = z, z-(z*z-x)/(2*z)
    }
    return z, nil
}

func main() {
    fmt.Println(Sqrt3(2))
    fmt.Println(Sqrt3(-2))
}