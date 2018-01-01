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

	log.Printf("Getting Objects From Database: %v, and table: %v\n", Database, table)

	c := cSession.DB(Database).C(table)
	c.Find(bson.M{}).All(out)
}

func getValidUsersForProfile(session *mgo.Session, profile User, out interface{}) {
	cSession := session.Copy()
	defer cSession.Close()

	c := cSession.DB(Database).C(UsersTable)

	query := bson.M{"$or": []bson.M{bson.M{"shared_with": profile.ID}, bson.M{"_id": profile.ID}}}

	c.Find(query).All(out)
}

func insertObjectIntoTable(session *mgo.Session, table string, obj interface{}) error {
	cSession := session.Copy()
	defer cSession.Close()

	log.Printf("Inserting %v into table %s\n", obj, table)

	c := cSession.DB(Database).C(table)
	err := c.Insert(obj)

	return err
}

func getProfileFromGoogle(accessToken string) (*GoogleProfile, error) {
	url := fmt.Sprintf("%s?access_token=%s", GoogleProfileURL, accessToken)
	response, err := http.Get(url)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	decoder := json.NewDecoder(response.Body)

	if response.StatusCode == http.StatusUnauthorized {
		var googleOauthError GoogleOauthError
		err = decoder.Decode(&googleOauthError)
		return nil, googleOauthError
	}

	var profile GoogleProfile
	err = decoder.Decode(&profile)

	log.Println(err)

	return &profile, nil
}

func getUserByEmail(session *mgo.Session, email string) (*User, error) {
	var user User

	cSession := session.Copy()
	defer cSession.Close()

	c := cSession.DB(Database).C(UsersTable)
	err := c.Find(bson.M{"email": email}).One(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func storeTransactions(session *mgo.Session, t []Transaction) error {
	cSession := session.Copy()
	defer cSession.Close()

	c := cSession.DB(Database).C(TransactionsTable)

	for _, transaction := range t {
		if err := c.Insert(transaction); err != nil {
			log.Printf("Error: %v\n", err)
			return err
		}
	}

	return nil
}

func getTransactionsForProfile(session *mgo.Session, profile User, out interface{}) error {
	cSession := session.Copy()
	defer cSession.Close()

	transactionsCol := cSession.DB(Database).C(TransactionsTable)
	usersCol := cSession.DB(Database).C(UsersTable)

	// Grab all the IDs from the "SharedWith" property
	// There has to be a better way to do this
	var sharedUsers []User
	queryFindSharedIds := bson.M{"shared_with": profile.ID}

	err := usersCol.Find(queryFindSharedIds).All(&sharedUsers)

	// Loop and create a list of IDs
	// There has to be a better way to do this
	var sharedIds []bson.ObjectId
	for _, u := range sharedUsers {
		sharedIds = append(sharedIds, u.ID)
	}

	query := bson.M{"$or": []bson.M{bson.M{"submitted_by": bson.M{"$in": []bson.ObjectId(sharedIds)}}, bson.M{"submitted_by": profile.ID}}}
	err = transactionsCol.Find(query).All(out)

	return err
}

func writeJSONToHTTP(w http.ResponseWriter, code int, objects interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	encoder := json.NewEncoder(w)
	encoder.Encode(objects)
}

func (e GoogleOauthError) Error() string {
	return e.Err.Message
}
