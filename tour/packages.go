package main

import (
    "fmt"
    "math/rand"
)


func main() {
    rand.Seed(10)
    fmt.Println("My number is", rand.Intn(10))
    fmt.Println("My number is", rand.Intn(10))
}