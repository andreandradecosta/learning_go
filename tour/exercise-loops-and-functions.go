package main

import (
    "fmt"
    "math"
)

func Sqrt(x float64) float64 {
    z := 1.0
    for i := 0; i < 10; i++ {
        z = z - (math.Pow(z, 2) - x)/(2*z)
    }
    return z
}


func Sqrt2(x float64) float64 {
    minDelta := math.Pow(10, -10)
    z2 := 1.0
    for delta := x; delta > minDelta; {
        z1 := z2
        z2 = z1 - (math.Pow(z1, 2) - x) / (2 * z1)
        delta = math.Abs(z2 - z1)
    }
    return z2
}

func main() {

    fmt.Printf("Mine: %v, His: %v, Diff: %v\n", Sqrt2(1), math.Sqrt(1), Sqrt2(1) - math.Sqrt(1))
    fmt.Printf("Mine: %v, His: %v, Diff: %v\n", Sqrt2(2), math.Sqrt(2), Sqrt2(2) - math.Sqrt(2))
    fmt.Printf("Mine: %v, His: %v, Diff: %v\n", Sqrt2(3), math.Sqrt(3), Sqrt2(3) - math.Sqrt(3))
    fmt.Printf("Mine: %v, His: %v, Diff: %v\n", Sqrt2(4), math.Sqrt(4), Sqrt2(4) - math.Sqrt(4))
    fmt.Printf("Mine: %v, His: %v, Diff: %v\n", Sqrt2(9), math.Sqrt(9), Sqrt2(9) - math.Sqrt(9))
    fmt.Printf("Mine: %v, His: %v, Diff: %v\n", Sqrt2(2345235), math.Sqrt(2345235), Sqrt2(2345235) - math.Sqrt(2345235))
}

