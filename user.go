package main

import "gopkg.in/mgo.v2/bson"

type User struct {
	ID        bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	FirstName string        `json:"first_name" bson:"first_name"`
	LastName  string        `json:"last_name" bson:"last_name"`
	Email     string        `json:"email" bson:"email"`
}
