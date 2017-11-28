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

func (a *App) Initialize(dbHost string) {
	var err error

	a.Router = mux.NewRouter()
	a.Session, err = mgo.Dial(dbHost)

	if err != nil {
		log.Fatal(err)
	}

	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/", a.indexHandler)
	a.Router.HandleFunc("/error", a.errorHandler)
}

func (a *App) indexHandler(w http.ResponseWriter, r *http.Request) {
	session := a.Session.Copy()
	defer session.Close()

	dbNames, err := session.DatabaseNames()

	if err != nil {
		writeErrorToHttp(w, http.StatusInternalServerError, "Unable to Get Database Names")
		return
	}

	writeJsonToHttp(w, dbNames)
}

func (a *App) errorHandler(w http.ResponseWriter, r *http.Request) {
	session := a.Session.Copy()
	defer session.Close()

	writeErrorToHttp(w, http.StatusInternalServerError, "This is the error route")
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
