package main


import (
    "golang.org/x/tour/wc"
    "strings"
)

func WordCount(s string) map[string]int {
    count := make(map[string]int)
    for _, element := range strings.Fields(s) {
        count[element]++
    }
    return count
}

func main() {
    wc.Test(WordCount)
}

