package main

import (
	"net/http"

	"gopkg.in/unrolled/render.v1"
)

func main() {
	rd := render.New(render.Options{
		IndentJSON: true,
	})
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome, visit sub pages now."))
	})

	mux.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		rd.Data(w, http.StatusOK, []byte("Some binary data."))
	})

	mux.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		rd.JSON(w, http.StatusOK, map[string]string{"chave": "valor"})
	})

	mux.HandleFunc("/json2", func(w http.ResponseWriter, r *http.Request) {
		type Book struct {
			Title  string
			Author string
		}
		book := Book{"Building Web Apps with Go", "Jeremy Saenz"}
		rd.JSON(w, http.StatusOK, book)
	})

	mux.HandleFunc("/html", func(w http.ResponseWriter, r *http.Request) {
		rd.HTML(w, http.StatusOK, "example", nil)
	})

	http.ListenAndServe(":8081", mux)
}
