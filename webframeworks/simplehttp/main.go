package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t1 := time.Now()
	fmt.Fprintf(w, "Welcome!")
	t2 := time.Now()
	log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	t1 := time.Now()
	fmt.Fprintf(w, "About page")
	t2 := time.Now()
	log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/about", aboutHandler)
	http.ListenAndServe(":8080", nil)
}
