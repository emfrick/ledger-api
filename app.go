package main

import (
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
	log.Printf("Started App on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/", a.indexHandler)
	a.Router.Handle("/users", TokenValidationHandler(a.Session, a.usersHandler))
	a.Router.HandleFunc("/auth", a.authHandler).Methods("POST")
	a.Router.Handle("/transactions", TokenValidationHandler(a.Session, a.postTransactions)).Methods("POST")
	a.Router.Handle("/transactions", TokenValidationHandler(a.Session, a.getTransactions)).Methods("GET")
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
