package main

import (
	"net/http"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type GoogleProfile struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	FullName      string `json:"name"`
	FirstName     string `json:"given_name"`
	LastName      string `json:"family_name"`
	ProfileLink   string `json:"link"`
	Picture       string `json:"picture"`
	Gender        string `json:"gender"`
	Locale        string `json:"locale"`
}

type User struct {
	ID        bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	FirstName string        `json:"first_name" bson:"first_name"`
	LastName  string        `json:"last_name" bson:"last_name"`
	Email     string        `json:"email" bson:"email"`
	Gender    string        `json:"gender" bson:"gender"`
}

type Transaction struct {
	ID          bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	Title       string        `json:"title" bson:"title"`
	Description string        `json:"description" bson:"description"`
	Date        time.Time     `json:"date" bson:"date"`
	Amount      int64         `json:"amount" bson:"amount"`
	User        bson.ObjectId `json:"user_id" bson:"user_id"`
	Tags        []string      `json:"categories" bson:"categories"`
}

type AuthorizedHttpHandlerFunc func(string, http.ResponseWriter, *http.Request)
