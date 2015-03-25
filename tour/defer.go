package main


import (
    "fmt"
)

func main() {
    i := 1
    defer fmt.Printf("world: %v\n", i)
    i++
    fmt.Printf("Hello: %v\n", i)

}

