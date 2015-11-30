package main

import (
	"net/http"
	"os"
	"github.com/andreandradecosta/webapps/fileserver/Godeps/_workspace/src/github.com/russross/blackfriday"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	http.HandleFunc("/markdown", generateMarkdown)
	http.Handle("/", http.FileServer(http.Dir("public")))
	http.ListenAndServe(":"+port, nil)
}

func generateMarkdown(w http.ResponseWriter, r *http.Request) {
	markdown := blackfriday.MarkdownCommon([]byte(r.FormValue("body")))
	w.Write(markdown)
}
