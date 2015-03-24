package main

import "fmt"

const Pi = 3.14

func main() {
    const World = "abc"
    fmt.Println("Hello", World)
    fmt.Println("Happy", Pi, "Day")

    const Truth bool = true
    fmt.Println("Go rules?", Truth)
}
