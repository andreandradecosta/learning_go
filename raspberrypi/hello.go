package main

//env GOOS=linux GOARCH=arm GOARM=7 go build hello.go

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Printf("Hello. I'm running %s on %s arch.\n", runtime.GOOS, runtime.GOARCH)
}
