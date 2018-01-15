package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"gopkg.in/mgo.v2/bson"

	jwt "github.com/dgrijalva/jwt-go"
)

// Handler for the / route
func (a *App) indexHandler(w http.ResponseWriter, r *http.Request) {

	response := map[string]string{}
	response["message"] = "This is the index route"

	writeJSONToHTTP(w, http.StatusOK, response)
}

// Handler for the /error route
func (a *App) errorHandler(w http.ResponseWriter, r *http.Request) {
	writeJSONToHTTP(w, http.StatusInternalServerError, ResponseError{"This is the error route"})
}

func (a *App) postSharedHandler(profile *User, w http.ResponseWriter, r *http.Request) {
	var decodedJSON map[string]string
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&decodedJSON)
	sharedEmail := decodedJSON["share_with"]

	if err != nil || sharedEmail == "" {
		log.Println(err)
		writeJSONToHTTP(w, http.StatusBadRequest, ResponseError{"Unable to get shared email (share_with)"})
		return
	}

	var sharedProfile *User
	sharedProfile, err = a.Database.UAL.GetUserByEmail(decodedJSON["share_with"])

	if err != nil {
		log.Println(err)
		formattedResponse := fmt.Sprintf("Unable to find user '%s'", sharedEmail)
		writeJSONToHTTP(w, http.StatusBadRequest, ResponseError{formattedResponse})
		return
	}

	err = a.Database.UAL.AddSharedUserToProfile(*sharedProfile, *profile)

	if err != nil {
		log.Println(err)
		formattedResponse := fmt.Sprintf("Unable to share with user user '%s'", sharedEmail)
		writeJSONToHTTP(w, http.StatusInternalServerError, ResponseError{formattedResponse})
		return
	}

	successMsg := fmt.Sprintf("Successfully shared with user '%s'", sharedEmail)
	writeJSONToHTTP(w, http.StatusCreated, map[string]string{"message": successMsg})
}

// Handler for the GET shared route
// Returns a list of users that are shared with the current user
func (a *App) getSharedHandler(profile *User, w http.ResponseWriter, r *http.Request) {
	var users []User

	a.Database.UAL.GetSharedUsersForProfile(*profile, &users)

	writeJSONToHTTP(w, http.StatusOK, users)
}

func (a *App) deleteSharedHandler(profile *User, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	if id == "" {
		log.Println("Missing ID for DELETE /shared")
		writeJSONToHTTP(w, http.StatusBadRequest, ResponseError{"Missing ID"})
		return
	}

	if bson.ObjectIdHex(id) == profile.ID {
		log.Println("Attempt to remove self from /shared")
		writeJSONToHTTP(w, http.StatusBadRequest, ResponseError{"You cannot unshare with yourself"})
		return
	}

	sharedUser, err := a.Database.UAL.GetUserById(id)

	if err != nil {
		log.Println(err)
		writeJSONToHTTP(w, http.StatusBadRequest, ResponseError{"Unable to find user"})
		return
	}

	err = a.Database.UAL.RemoveSharedUserFromProfile(*sharedUser, *profile)

	if err != nil {
		log.Println(err)
		writeJSONToHTTP(w, http.StatusInternalServerError, ResponseError{"Unable to remove user from shares"})
		return
	}

	successMsg := fmt.Sprintf("Successfully removed user '%s' from your shared users", sharedUser.Email)
	writeJSONToHTTP(w, http.StatusOK, map[string]string{"message": successMsg})
}

// Handler for the GET /transactions route
func (a *App) getTransactions(profile *User, w http.ResponseWriter, r *http.Request) {
	var transactions []Transaction

	err := a.Database.TAL.GetTransactionsForProfile(*profile, &transactions)

	if err != nil {
		writeJSONToHTTP(w, http.StatusInternalServerError, ResponseError{"Error getting transactions"})
	}

	writeJSONToHTTP(w, http.StatusOK, transactions)
}

// Handler for the POST /transactions route
// Takes a JSON array of transaction objects
func (a *App) postTransactions(profile *User, w http.ResponseWriter, r *http.Request) {
	var t []Transaction
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&t)

	if err != nil {
		log.Println(err)
		writeJSONToHTTP(w, http.StatusInternalServerError, ResponseError{"Unable to decode JSON"})
		return
	}

	for index := range t {
		t[index].SubmittedBy = bson.ObjectId(profile.ID)
	}

	if err = a.Database.TAL.StoreTransactions(t); err != nil {
		writeJSONToHTTP(w, http.StatusInternalServerError, ResponseError{"Unable to save transactions"})
	}

	writeJSONToHTTP(w, http.StatusCreated, t)
}

// Handles authentication against Google and creates a JWT
func (a *App) authHandler(w http.ResponseWriter, r *http.Request) {

	var data map[string]string

	// Decode the request body
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&data)

	// Make sure there is an access token (this would come from the OAuth provider)
	if data["access_token"] == "" {
		writeJSONToHTTP(w, http.StatusBadRequest, ResponseError{"Google Access Token Required"})
		return
	}

	// Get the profile from Google
	profile, err := getProfileFromGoogle(data["access_token"])

	if err != nil {
		writeJSONToHTTP(w, http.StatusInternalServerError, ResponseError{err.Error()})
		return
	}

	// Check if the user exists (by email)
	if a.Database.UAL.DoesProfileExist(profile.Email) {
		// Login
		log.Println("Login")
	} else {
		// Register!
		err = a.Database.UAL.AddUser(*profile)

		if err != nil {
			writeJSONToHTTP(w, http.StatusInternalServerError, ResponseError{"Unable to add user"})
			return
		}
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
		writeJSONToHTTP(w, http.StatusInternalServerError, ResponseError{"Error getting token"})
		log.Println(err)
		return
	}

	writeJSONToHTTP(w, http.StatusOK, map[string]string{"token": tokenString})
}
