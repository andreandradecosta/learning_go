package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/justinas/alice"
)

func loggingHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()
		log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
	}
	return http.HandlerFunc(fn)
}

func recoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "About page")
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome!")
}

func panicGenerator(w http.ResponseWriter, r *http.Request) {
	panic("error")
}

func main() {
	commonHandler := alice.New(loggingHandler, recoverHandler)
	http.Handle("/about", commonHandler.ThenFunc(aboutHandler))
	http.Handle("/", commonHandler.ThenFunc(indexHandler))
	http.Handle("/panic", commonHandler.ThenFunc(panicGenerator))
	http.ListenAndServe(":8080", nil)
}
