package main

import (
	"log"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"
	mgo "gopkg.in/mgo.v2"
)

// App holds the mux router and mongo session
type App struct {
	Router  *mux.Router
	Session *mgo.Session
}

// NewApp returns a new instance of the App connected to the given database hose
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

// Run starts up the app listening on the given address
func (a *App) Run(addr string) {
	log.Printf("Started App on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

// Define the handlers for the routes
func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/", a.indexHandler)
	a.Router.HandleFunc("/auth", a.authHandler).Methods("POST")
	a.Router.Handle("/shared", TokenValidationHandler(a.Session, a.getSharedHandler)).Methods("GET")
	a.Router.Handle("/transactions", TokenValidationHandler(a.Session, a.postTransactions)).Methods("POST")
	a.Router.Handle("/transactions", TokenValidationHandler(a.Session, a.getTransactions)).Methods("GET")
	a.Router.HandleFunc("/error", a.errorHandler)
}

// DoesProfileExist checks the database for the given email
func (a *App) DoesProfileExist(p GoogleProfile) bool {

	session := a.Session.Copy()
	defer session.Close()

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
