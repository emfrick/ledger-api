package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func getAllObjectsFromTable(session *mgo.Session, table string, out interface{}) {
	cSession := session.Copy()

	defer cSession.Close()

	fmt.Printf("Getting Objects From Database: %v, and table: %v", Database, table)

	c := cSession.DB(Database).C(table)

	c.Find(bson.M{}).All(out)
	fmt.Printf("Objects: %v\n", &out)
}

func getValidUsersForProfile(session *mgo.Session, profile User, out interface{}) {
	cSession := session.Copy()

	defer cSession.Close()

	c := cSession.DB(Database).C(UsersTable)

	c.Find(bson.M{"sharesWith": profile.ID}).All(out)
}

func insertObjectIntoTable(session *mgo.Session, table string, obj interface{}) {
	cSession := session.Copy()

	defer cSession.Close()

	fmt.Printf("Inserting %v into table %s", obj, table)

	c := cSession.DB(Database).C(table)

	err := c.Insert(obj)

	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func getProfileFromGoogle(accessToken string) (*GoogleProfile, error) {
	url := fmt.Sprintf("%s?access_token=%s", GoogleProfileURL, accessToken)
	response, err := http.Get(url)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	var profile GoogleProfile

	decoder := json.NewDecoder(response.Body)
	decoder.Decode(&profile)

	return &profile, nil
}

func getUserByEmail(session *mgo.Session, email string) *User {
	cSession := session.Copy()

	defer cSession.Close()

	c := cSession.DB(Database).C(UsersTable)

	var user User
	err := c.Find(bson.M{"email": email}).One(&user)

	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	return &user
}
