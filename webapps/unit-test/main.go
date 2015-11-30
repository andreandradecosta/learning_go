package main

import (
	"fmt"
	"net/http"
)

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World:", r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/", HelloWorld)
	http.ListenAndServe(":8081", nil)
}
