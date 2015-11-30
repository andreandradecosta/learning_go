package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

var books = []Book{
	{"Buildin Web Apps with Go", "Jeremy Saenz", 0},
	{"Go Lang", "Xxx Yyy Zzz", 1},
}

type Book struct {
	Title  string `json:"book-title"`
	Author string `json:"book-author"`
	ID     int
}

func main() {
	http.HandleFunc("/", showBooks)
	http.ListenAndServe(":8081", nil)
}

func showBooks(w http.ResponseWriter, r *http.Request) {
	js, err := json.Marshal(books)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	var out bytes.Buffer
	json.Indent(&out, js, "", "\t")
	out.WriteTo(os.Stdout)
}
