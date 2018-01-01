package main

import (
	"net/http"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// GoogleProfile is a struct that contains the attributes returned from Google
type GoogleProfile struct {
	ID            string `json:"id"`
	Email         string `json:"email" bson:"email"`
	VerifiedEmail bool   `json:"verified_email" bson:"verified_email"`
	FullName      string `json:"name" bson:"full_name"`
	FirstName     string `json:"given_name" bson:"first_name"`
	LastName      string `json:"family_name" bson:"last_name"`
	ProfileLink   string `json:"link" bson:"link"`
	Picture       string `json:"picture" bson:"picture"`
	Gender        string `json:"gender" bson:"gender"`
	Locale        string `json:"locale" bson:"locale"`
}

// GoogleOauthError is a custom Error object
type GoogleOauthError struct {
	Err struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// User is a representation of a stored user
type User struct {
	ID         bson.ObjectId   `json:"id,omitempty" bson:"_id,omitempty"`
	FirstName  string          `json:"first_name" bson:"first_name"`
	LastName   string          `json:"last_name" bson:"last_name"`
	Email      string          `json:"email" bson:"email"`
	Gender     string          `json:"gender" bson:"gender"`
	SharedWith []bson.ObjectId `json:"-" bson:"shared_with"`
}

// Transaction represents a transaction
type Transaction struct {
	ID          bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Title       string        `json:"title" bson:"title"`
	Description string        `json:"description" bson:"description"`
	Date        time.Time     `json:"date" bson:"date"`
	Amount      int64         `json:"amount" bson:"amount"`
	SubmittedBy bson.ObjectId `json:"submitted_by" bson:"submitted_by"`
	Tags        []string      `json:"categories" bson:"categories"`
}

// AuthorizedHTTPHandlerFunc is passed in middleware
type AuthorizedHTTPHandlerFunc func(*User, http.ResponseWriter, *http.Request)

// ResponseError is meant to be returned by
type ResponseError struct {
	Error string `json:"error"`
}
