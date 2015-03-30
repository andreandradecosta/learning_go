package main

import (
    "fmt"
)


func main() {
    var a []int
    printSlice2("a", a)

    a = append(a, 0)
    printSlice2("a", a)

    a = append(a, 1)
    printSlice2("a", a)

    a = append(a, 2, 3, 4)
    printSlice2("a", a)
}


func printSlice2(s string, x []int) {
    fmt.Printf("%s len=%d cap=%d %v\n",
        s, len(x), cap(x), x)
}