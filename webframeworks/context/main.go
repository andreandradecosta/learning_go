package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/context"
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
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
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

func getUser(db *sql.DB, authToken string) (string, error) {
	return "user", nil
}

func main() {
	app := new(appContext)
	//app.db = sql.Open("postgres", "...")
	commonHandlers := alice.New(context.ClearHandler, loggingHandler, recoverHandler)
	http.Handle("/admin", commonHandlers.Append(app.authHandler).ThenFunc(app.adminHandler))
	http.Handle("/about", commonHandlers.ThenFunc(aboutHandler))
	http.Handle("/", commonHandlers.ThenFunc(indexHandler))
	http.Handle("/panic", commonHandlers.ThenFunc(panicGenerator))
	http.ListenAndServe(":8080", nil)
}
