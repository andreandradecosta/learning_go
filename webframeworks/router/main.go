package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
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

type appContext struct {
	db *sql.DB
}

func (c *appContext) authHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get("Authorization")
		user, err := getUser(c.db, authToken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		context.Set(r, "user", user)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func (c *appContext) adminHandler(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user")
	json.NewEncoder(w).Encode(user)
}

func (c *appContext) teaHandler(w http.ResponseWriter, r *http.Request) {
	params := context.Get(r, "params").(httprouter.Params)
	tea := getTea(c.db, params.ByName("id"))
	json.NewEncoder(w).Encode(tea)
}

func getUser(db *sql.DB, authToken string) (string, error) {
	return "user", errors.New("No user")
}

func getTea(db *sql.DB, id string) string {
	return "Tea" + id
}

func wrapHandler(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		context.Set(r, "params", p)
		h.ServeHTTP(w, r)
	}
}

func main() {
	db, _ := sql.Open("postgres", "...")
	app := appContext{db}
	commonHandlers := alice.New(context.ClearHandler, loggingHandler, recoverHandler)
	router := httprouter.New()
	router.GET("/admin", wrapHandler(commonHandlers.Append(app.authHandler).ThenFunc(app.adminHandler)))
	router.GET("/about", wrapHandler(commonHandlers.ThenFunc(aboutHandler)))
	router.GET("/", wrapHandler(commonHandlers.ThenFunc(indexHandler)))
	router.GET("/tea/:id", wrapHandler(commonHandlers.ThenFunc(app.teaHandler)))
	router.GET("/panic", wrapHandler(commonHandlers.ThenFunc(panicGenerator)))
	http.ListenAndServe(":8080", router)
}
