package main

import (
	"net/http"

	"gopkg.in/unrolled/render.v1"
)

type Action func(w http.ResponseWriter, r *http.Request) error

type AppController struct{}

func (c *AppController) Action(a Action) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := a(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

type MyController struct {
	AppController
	*render.Render
}

type Book struct {
	Title  string
	Author string
}

var myBook = Book{"Título", "André"}

func (c *MyController) Index(w http.ResponseWriter, r *http.Request) error {
	c.JSON(w, http.StatusOK, myBook)
	return nil
}

func main() {
	c := &MyController{Render: render.New(render.Options{})}
	http.ListenAndServe(":8081", c.Action(c.Index))
}
