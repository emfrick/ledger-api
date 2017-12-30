package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	mgo "gopkg.in/mgo.v2"
)

type App struct {
	Router  *mux.Router
	Session *mgo.Session
}

func NewApp(dbHost string) *App {
	var app = App{}

	var err error

	app.Router = mux.NewRouter()
	app.Session, err = mgo.Dial(dbHost)

	if err != nil {
		log.Fatal(err)
	}

	app.initializeRoutes()

	return &app
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))

	log.Printf("Started App on %s", addr)
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/", a.indexHandler)
	a.Router.Handle("/users", TokenValidationHandler(http.HandlerFunc(a.usersHandler)))
	a.Router.HandleFunc("/auth", a.authHandler)
	a.Router.HandleFunc("/error", a.errorHandler)
}

func writeErrorToHttp(w http.ResponseWriter, code int, message string) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	encoder := json.NewEncoder(w)
	encoder.Encode(map[string]string{"error": message})
}

func writeJsonToHttp(w http.ResponseWriter, objects interface{}) {
	respBody, err := json.MarshalIndent(objects, "", "  ")

	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}
