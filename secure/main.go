package main

import (
	"log"
	"net/http"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/unrolled/secure" // or "gopkg.in/unrolled/secure.v1"
)

var port = os.Getenv("PORT")

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("Teste"))
	})

	secureMiddleware := secure.New(secure.Options{
		AllowedHosts:          []string{"localhost:8080"},
		SSLRedirect:           true,
		SSLHost:               "localhost:8443",
		STSSeconds:            315360000,
		STSIncludeSubdomains:  true,
		STSPreload:            true,
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self'",
		PublicKey:             `pin-sha256="base64+primary=="; pin-sha256="base64+backup=="; max-age=5184000; includeSubdomains; report-uri="https://www.example.com/hpkp-report"`,
		IsDevelopment:         true,
	})
	n := negroni.Classic()
	n.Use(negroni.HandlerFunc(secureMiddleware.HandlerFuncWithNext))
	n.UseHandler(mux)

	addr := ":8080"
	httpsAddr := ":8443"
	l := log.New(os.Stdout, "[negroni] ", 0)
	l.Printf("listening on %s and %s", addr, httpsAddr)
	// HTTP
	go func() {
		log.Fatal(http.ListenAndServe(addr, n))
	}()

	l.Fatal(http.ListenAndServeTLS(httpsAddr, "cert.pem", "key.pem", n))
}
