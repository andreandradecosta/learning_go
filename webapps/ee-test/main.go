package main

import (
	"fmt"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
)

func HelloWorld(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprint(w, "Hello World")
}

func App() http.Handler {
	n := negroni.Classic()

	m := func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		fmt.Fprint(w, "Before...")
		next(w, r)
		fmt.Fprint(w, "...After")
	}
	n.Use(negroni.HandlerFunc(m))
	r := httprouter.New()
	r.GET("/", HelloWorld)
	n.UseHandler(r)
	return n
}

func main() {
	http.ListenAndServe(":8081", App())
}
