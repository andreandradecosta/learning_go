package main

import (
	"io"
	"log"
	"net/http"
)

func HelloServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello, world!\n")
}

//openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout key.pem -out cert.pem
func main() {
	http.HandleFunc("/", HelloServer)
	err := http.ListenAndServeTLS(":8443", "cert.pem", "key.pem", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
