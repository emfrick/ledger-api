package main

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Transaction struct {
	ID          bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	Title       string        `json:"title" bson:"title"`
	Description string        `json:"description" bson:"description"`
	Date        time.Time     `json:"date" bson:"date"`
	Amount      int64         `json:"amount" bson:"amount"`
	User
	Tags []string `json:"categories" bson:"categories"`
}
