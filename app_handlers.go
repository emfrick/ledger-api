package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"gopkg.in/mgo.v2/bson"

	jwt "github.com/dgrijalva/jwt-go"
)

func (a *App) indexHandler(w http.ResponseWriter, r *http.Request) {
	session := a.Session.Copy()
	defer session.Close()

	dbNames, err := session.DatabaseNames()

	if err != nil {
		writeErrorToHTTP(w, http.StatusInternalServerError, "Unable to Get Database Names")
		return
	}

	writeJSONToHTTP(w, http.StatusOK, dbNames)
}

func (a *App) errorHandler(w http.ResponseWriter, r *http.Request) {
	session := a.Session.Copy()
	defer session.Close()

	writeErrorToHTTP(w, http.StatusInternalServerError, "This is the error route")
}

func (a *App) usersHandler(email string, w http.ResponseWriter, r *http.Request) {
	var users []User

	profile := getUserByEmail(a.Session, email)

	getValidUsersForProfile(a.Session, *profile, &users)

	writeJSONToHTTP(w, http.StatusOK, users)
}

func (a *App) getTransactions(email string, w http.ResponseWriter, r *http.Request) {
	profile := getUserByEmail(a.Session, email)

	var transactions []Transaction

	err := getTransactionsForProfile(a.Session, *profile, &transactions)

	if err != nil {
		writeErrorToHTTP(w, http.StatusInternalServerError, "Error getting transactions")
	}

	writeJSONToHTTP(w, http.StatusOK, transactions)
}

func (a *App) postTransactions(email string, w http.ResponseWriter, r *http.Request) {
	profile := getUserByEmail(a.Session, email)

	var t []Transaction
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&t)

	if err != nil {
		log.Println(err)
		writeErrorToHTTP(w, http.StatusInternalServerError, "Unable to decode JSON")
		return
	}

	for index := range t {
		t[index].User = bson.ObjectId(profile.ID)
	}

	if err = storeTransactions(a.Session, t); err != nil {
		writeErrorToHTTP(w, http.StatusInternalServerError, "Unable to save transactions")
	}

	writeJSONToHTTP(w, http.StatusCreated, t)
}

func (a *App) authHandler(w http.ResponseWriter, r *http.Request) {

	var data map[string]string

	// Decode the request body
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&data)

	// Make sure there is an access token (this would come from the OAuth provider)
	if data["access_token"] == "" {
		writeErrorToHTTP(w, http.StatusBadRequest, "Google Access Token Required")
		return
	}

	// Get the profile from Google
	profile, err := getProfileFromGoogle(data["access_token"])

	if err != nil {
		writeErrorToHTTP(w, http.StatusInternalServerError, "Error Getting Google Profile")
		return
	}

	// Check if the user exists (by email)
	if a.DoesProfileExist(*profile) {
		// Login
		log.Println("Login")
	} else {
		// Register!
		insertObjectIntoTable(a.Session, UsersTable, profile)
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

	if err != nil {
		writeErrorToHTTP(w, http.StatusInternalServerError, "Error getting token")
		log.Println(err)
		return
	}

	writeJSONToHTTP(w, http.StatusOK, tokenString)
}
