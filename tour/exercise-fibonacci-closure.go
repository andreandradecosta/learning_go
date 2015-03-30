package main

import (
    "fmt"
)

func fibonacci() func() int {
    x, y := 0, 1
    return func() int {
        x, y = y, y+x
        return x
    }
}


func main() {
    f := fibonacci()
    for i := 0; i < 10; i++ {
        fmt.Println(f())
    }
}
