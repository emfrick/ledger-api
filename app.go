package main

import (
	"encoding/json"
	"log"
	"net/http"

	"gopkg.in/mgo.v2/bson"

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
	a.Router.Handle("/users", TokenValidationHandler(a.usersHandler))
	a.Router.HandleFunc("/auth", a.authHandler).Methods("POST")
	a.Router.HandleFunc("/error", a.errorHandler)
}

func (a *App) DoesProfileExist(p GoogleProfile) bool {

	session := a.Session.Copy()
	col := session.DB(Database).C(UsersTable)

	count, err := col.Find(bson.M{"email": p.Email}).Count()

	if err != nil {
		log.Println(err)
	}

	if count > 0 {
		return true
	}

	return false
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
