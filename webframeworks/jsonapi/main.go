package main

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"runtime/debug"
	"time"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Repo

type Tea struct {
	ID       bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string        `json:"name"`
	Category string        `json:"category"`
}

type TeasCollection struct {
	Data []Tea `json:"data"`
}

type TeaResource struct {
	Data Tea `json:"data"`
}

type TeaRepo struct {
	coll *mgo.Collection
}

func (r *TeaRepo) All() (TeasCollection, error) {
	result := TeasCollection{[]Tea{}}
	err := r.coll.Find(nil).All(&result.Data)
	return result, err
}

func (r *TeaRepo) Find(id string) (TeaResource, error) {
	result := TeaResource{}
	err := r.coll.FindId(bson.ObjectIdHex(id)).One(&result.Data)
	return result, err
}

func (r *TeaRepo) Create(tea *Tea) error {
	id := bson.NewObjectId()
	_, err := r.coll.UpsertId(id, tea)
	if err != nil {
		return err
	}
	tea.ID = id
	return nil
}

func (r *TeaRepo) Update(tea *Tea) error {
	return r.coll.UpdateId(tea.ID, tea)
}

func (r *TeaRepo) Delete(id string) error {
	return r.coll.RemoveId(bson.ObjectIdHex(id))
}

//Errors

type Errors struct {
	Errors []*Error `json:"errors"`
}

type Error struct {
	ID     string `json:"id"`
	Status int    `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

var (
	ErrBadRequest           = &Error{"bad_request", http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Request not valid"}
	ErrNotAcceptable        = &Error{"not_acceptable", http.StatusNotAcceptable, http.StatusText(http.StatusNotAcceptable), "Accept header must be set to 'application/vnd.api+json'."}
	ErrUnsupportedMediaType = &Error{"unsupported_media_type", http.StatusUnsupportedMediaType, http.StatusText(http.StatusUnsupportedMediaType), "Content-Type header must be set to: 'application/vnd.api+json'."}
	ErrInternalServer       = &Error{"internal_server_error", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Erro desconhecido"}
)

func WriteError(w http.ResponseWriter, err *Error) {
	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.WriteHeader(err.Status)
	json.NewEncoder(w).Encode(Errors{[]*Error{err}})
}

//Middlewares

func recoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %v", err)
				debug.PrintStack()
				WriteError(w, ErrInternalServer)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func loggingHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), time.Since(t1))
	}
	return http.HandlerFunc(fn)
}

func acceptHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept") != "application/vnd.api+json" {
			WriteError(w, ErrNotAcceptable)
			return
		}
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func contentTypeHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/vnd.api+json" {
			WriteError(w, ErrUnsupportedMediaType)
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func bodyParseHandler(v interface{}) func(http.Handler) http.Handler {
	t := reflect.TypeOf(v)

	m := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			val := reflect.New(t).Interface()
			err := json.NewDecoder(r.Body).Decode(val)
			if err != nil {
				WriteError(w, ErrBadRequest)
				return
			}

			if next != nil {
				context.Set(r, "body", val)
				next.ServeHTTP(w, r)
			}
		}
		return http.HandlerFunc(fn)

	}
	return m
}

// Main handlers

type appContext struct {
	db *mgo.Database
}

func (c *appContext) teasHandler(w http.ResponseWriter, r *http.Request) {
	repo := TeaRepo{c.db.C("teas")}
	teas, err := repo.All()
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/vnd.api+json")
	json.NewEncoder(w).Encode(teas)
}

func (c *appContext) teaHandler(w http.ResponseWriter, r *http.Request) {
	params := context.Get(r, "params").(httprouter.Params)
	repo := TeaRepo{c.db.C("teas")}
	tea, err := repo.Find(params.ByName("id"))
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/vnd.api+json")

	json.NewEncoder(w).Encode(tea)
}

func (c *appContext) createTeaHandler(w http.ResponseWriter, r *http.Request) {
	body := context.Get(r, "body").(*TeaResource)
	repo := TeaRepo{c.db.C("teas")}
	err := repo.Create(&body.Data)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(body)
}

func (c *appContext) updateTeaHandler(w http.ResponseWriter, r *http.Request) {
	params := context.Get(r, "params").(httprouter.Params)
	body := context.Get(r, "body").(*TeaResource)
	body.Data.ID = bson.ObjectIdHex(params.ByName("id"))
	repo := TeaRepo{c.db.C("teas")}
	err := repo.Update(&body.Data)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte("\n"))
}

func (c *appContext) deleteTeaHandler(w http.ResponseWriter, r *http.Request) {
	params := context.Get(r, "params").(httprouter.Params)
	repo := TeaRepo{c.db.C("teas")}
	err := repo.Delete(params.ByName("id"))
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte("\n"))
}

// Router

type router struct {
	*httprouter.Router
}

func (r *router) Get(path string, handler http.Handler) {
	r.GET(path, wrapHandler(handler))
}
func (r *router) Post(path string, handler http.Handler) {
	r.POST(path, wrapHandler(handler))
}
func (r *router) Put(path string, handler http.Handler) {
	r.PUT(path, wrapHandler(handler))
}
func (r *router) Delete(path string, handler http.Handler) {
	r.DELETE(path, wrapHandler(handler))
}

func NewRouter() *router {
	return &router{httprouter.New()}
}

func wrapHandler(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		context.Set(r, "params", p)
		h.ServeHTTP(w, r)
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	appC := appContext{session.DB("test")}
	commonHandlers := alice.New(context.ClearHandler, loggingHandler, recoverHandler, acceptHandler)
	router := NewRouter()
	router.Get("/teas/:id", commonHandlers.ThenFunc(appC.teaHandler))
	router.Put("/teas/:id", commonHandlers.Append(contentTypeHandler, bodyParseHandler(TeaResource{})).ThenFunc(appC.updateTeaHandler))
	router.Delete("/teas/:id", commonHandlers.ThenFunc(appC.deleteTeaHandler))
	router.Get("/teas", commonHandlers.ThenFunc(appC.teasHandler))
	router.Post("/teas", commonHandlers.Append(contentTypeHandler, bodyParseHandler(TeaResource{})).ThenFunc(appC.createTeaHandler))

	http.ListenAndServe(":8080", router)
}
