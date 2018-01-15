package main

import (
	"log"

	"gopkg.in/mgo.v2"
)

// Constants
const (
	DatabaseHost      = "localhost"
	DatabaseName      = "ledger"
	UsersTable        = "users"
	TransactionsTable = "transactions"
	SecretKey         = "SECRETKEYISHARDTOGUESS"
	GoogleProfileURL  = "https://www.googleapis.com/userinfo/v2/me"
)

func main() {

	// Create a MongoDB Session
	session, err := mgo.Dial(DatabaseHost)

	if err != nil {
		log.Fatal(err)
	}

	// Create a Mongo Access Layer
	mal := MongoAccessLayer{session}

	// Create a Database
	db := Database{
		UAL: mal,
		TAL: mal,
	}

	// Instantiate the App
	app := NewApp(&db)

	// Run the App
	app.Run(":3000")
}
