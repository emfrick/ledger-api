package main

import (
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

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

func (a *App) usersHandler(w http.ResponseWriter, r *http.Request) {
	var users []User

	getAllObjectsFromTable(a.Session, USERS_TABLE, &users)

	writeJsonToHttp(w, users)
}

func (a *App) authHandler(w http.ResponseWriter, r *http.Request) {
	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"foo": "bar",
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(SECRET_KEY))

	log.Println(tokenString, err)

	if err != nil {
		writeErrorToHttp(w, http.StatusInternalServerError, "Error getting token")
		log.Println(err)
		return
	}

	writeJsonToHttp(w, tokenString)
}
