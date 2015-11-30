package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	r := httprouter.New()
	r.GET("/", homeHandler)

	r.GET("/posts", postIndexHandler)
	r.POST("/posts", postCreateHandler)

	r.GET("/posts/:id", postShowHandler)
	r.PUT("/posts/:id", postUpdateHandler)
	r.GET("/posts/:id/edit", postEditHandler)

	r.NotFound = http.FileServer(http.Dir("public"))

	fmt.Println("Server starting on :8081")
	http.ListenAndServe(":8081", r)
}

func homeHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintln(w, "Home")
}

func postIndexHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintln(w, "posts Index")
}

func postCreateHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintln(w, "posts Create")
}

func postShowHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintln(w, "Show Post:", p.ByName("id"))
}

func postUpdateHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintln(w, "Post Update")
}

func postEditHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintln(w, "Post Edit", p.ByName("id"))
}
