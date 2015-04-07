package main
import (
    "golang.org/x/tour/tree"
    "fmt"
)

func walkImpl2(t *tree.Tree, ch chan int, quit chan int) {
    if t == nil {
        return
    }
    walkImpl2(t.Left, ch, quit)
    select {
        case ch <- t.Value:
        case <-quit:
            return
    }
    walkImpl2(t.Right, ch, quit)
}

func Walk2(t *tree.Tree, ch chan int, quit chan int) {
    walkImpl2(t, ch, quit)
    close(ch)
}


func Same2(t1, t2 *tree.Tree) bool {
    ch1, ch2 := make(chan int), make(chan int)
    quit := make(chan int)
    defer close(quit)

    go Walk2(t1, ch1, quit)
    go Walk2(t2, ch2, quit)
    for {
        v1, ok1 := <-ch1
        v2, ok2 := <-ch2
        if !ok1 || !ok2 {
            return ok1 == ok2
        }
        if v1 != v2 {
            return false
        }
    }
    return true
}

func main() {
    fmt.Println(Same2(tree.New(1), tree.New(1)))
    fmt.Println(Same2(tree.New(1), tree.New(2)))
}
