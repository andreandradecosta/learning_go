package main


import (
    "fmt"
)

func main() {
    var a [10]int
    fmt.Println(a)
    var b [2]string
    b[0] = "Hello"
    b[1] = "World"
    fmt.Println(b[0], b[1])
    fmt.Println(b)
}

