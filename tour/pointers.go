package main



import (
    "fmt"
)


func main() {
    var pt *int
    fmt.Println(pt)

    i, j := 42, 2701

    p := &i
    fmt.Println(*p, p)
    *p = 21
    fmt.Println(i)


    p = &j
    *p = *p / 37
    fmt.Println(j)
}
