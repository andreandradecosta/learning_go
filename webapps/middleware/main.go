package main

import (
	"log"
	"net/http"

	"github.com/codegangsta/negroni"
)

func main() {
	n := negroni.New(
		negroni.NewRecovery(),
		negroni.HandlerFunc(myMiddleware),
		negroni.NewLogger(),
		negroni.NewStatic(http.Dir("public")),
	)
	n.Run(":8081")
}

func myMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	log.Println("Init myMiddleware")

	if r.URL.Query().Get("username") == "myname" {
		next(w, r)
	} else {
		http.Error(w, "Not authorized", 401)
	}
	log.Println("End myMiddleware")
}
