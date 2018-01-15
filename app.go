package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// App holds the mux router and a database
type App struct {
	Router   *mux.Router
	Database *Database
}

// NewApp returns a new instance of the App connected to the given database hose
func NewApp(db *Database) *App {
	var app = App{}

	app.Router = mux.NewRouter()
	app.Database = db

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
	a.Router.Handle("/shared", TokenValidationHandler(a.Database, a.postSharedHandler)).Methods("POST")
	a.Router.Handle("/shared", TokenValidationHandler(a.Database, a.getSharedHandler)).Methods("GET")
	a.Router.Handle("/shared/{id}", TokenValidationHandler(a.Database, a.deleteSharedHandler)).Methods("DELETE")
	a.Router.Handle("/transactions", TokenValidationHandler(a.Database, a.postTransactions)).Methods("POST")
	a.Router.Handle("/transactions", TokenValidationHandler(a.Database, a.getTransactions)).Methods("GET")
	a.Router.HandleFunc("/error", a.errorHandler)
}
