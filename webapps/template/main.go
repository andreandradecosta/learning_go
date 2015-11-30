package main

import (
	"html/template"
	"net/http"
)

type Book struct {
	Title  string
	Author string
}

var (
	templates = template.Must(template.ParseFiles("templates/index.html"))
)

func main() {
	http.HandleFunc("/", showBooks)
	http.ListenAndServe(":8081", nil)
}

func showBooks(w http.ResponseWriter, r *http.Request) {
	book := Book{"Building Web Apps with Go", "Jeremy Saenz"}
	if err := templates.Execute(w, book); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
