package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
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

// CreateTokenForProfile returns a JWT for the given profile
func (a *App) CreateTokenForProfile(profile GoogleProfile) (string, error) {

	// Ensure profile exists
	if !a.Database.UAL.DoesProfileExist(profile.Email) {
		return "", errors.New("profile does not exist")
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"first_name": profile.FirstName,
		"last_name":  profile.LastName,
		"email":      profile.Email,
		"exp":        time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(SecretKey))

	return tokenString, err
}
