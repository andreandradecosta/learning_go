package main
import (
    "time"
    "fmt"
)


func main() {
    tick := time.Tick(1 * time.Second)
    boom := time.After(10 * time.Second)
    for {
        select {
            case <- tick:
                fmt.Println("tick.")
            case <- boom:
                fmt.Println("BOOM!")
                return
            default:
                fmt.Println("    .")
                time.Sleep(500 * time.Millisecond)
        }
    }
}
