package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Returns all documents from the given table
func getAllObjectsFromTable(session *mgo.Session, table string, out interface{}) {

	// Copy the mongo session and defer its close
	cSession := session.Copy()
	defer cSession.Close()

	log.Printf("Getting Objects From Database: %v, and table: %v\n", Database, table)

	c := cSession.DB(Database).C(table)
	c.Find(bson.M{}).All(out)
}

// Grabs the valid users given a profile (returns the given profile along with
// any accounts that are shred with that profile)
func getSharedUsersForProfile(session *mgo.Session, profile User, out interface{}) {

	// Copy the mongo session and defer its close
	cSession := session.Copy()
	defer cSession.Close()

	c := cSession.DB(Database).C(UsersTable)

	// Return the user OR and accounts that share with the user
	query := bson.M{"$or": []bson.M{bson.M{"shared_with": profile.ID}, bson.M{"_id": profile.ID}}}

	c.Find(query).All(out)
}

// Inserts a document into the given table
func insertObjectIntoTable(session *mgo.Session, table string, obj interface{}) error {

	// Copy the mongo session and defer its close
	cSession := session.Copy()
	defer cSession.Close()

	log.Printf("Inserting %v into table %s\n", obj, table)

	c := cSession.DB(Database).C(table)
	err := c.Insert(obj)

	return err
}

// Takes an access token and returns the Google Profile
func getProfileFromGoogle(accessToken string) (*GoogleProfile, error) {

	// Create the URL and run an http GET
	url := fmt.Sprintf("%s?access_token=%s", GoogleProfileURL, accessToken)
	response, err := http.Get(url)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Decode the response body
	decoder := json.NewDecoder(response.Body)

	// Make sure Google didn't return 401
	if response.StatusCode == http.StatusUnauthorized {
		var googleOauthError GoogleOauthError
		err = decoder.Decode(&googleOauthError)
		return nil, googleOauthError
	}

	// Decode the JSON
	var profile GoogleProfile
	err = decoder.Decode(&profile)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &profile, nil
}

// Returns the profile for the given email
func getUserByEmail(session *mgo.Session, email string) (*User, error) {
	var user User

	// Copy the mongo session and defer its close
	cSession := session.Copy()
	defer cSession.Close()

	c := cSession.DB(Database).C(UsersTable)
	err := c.Find(bson.M{"email": email}).One(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Stores the given transactions in the Mongo Database
func storeTransactions(session *mgo.Session, t []Transaction) error {

	// Copy the mongo session and defer its close
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

// Returns all transactions submitted by the given User as well as any transactions
// that are shared with the user
func getTransactionsForProfile(session *mgo.Session, profile User, out interface{}) error {

	// Copy the mongo session and defer its close
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

// Writes JSON to the response
func writeJSONToHTTP(w http.ResponseWriter, code int, objects interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	encoder := json.NewEncoder(w)
	encoder.Encode(objects)
}

// Custom Error
func (e GoogleOauthError) Error() string {
	return e.Err.Message
}

func addSharedUserToProfile(session *mgo.Session, sharedUser User, profile User) error {
	cSession := session.Copy()
	defer cSession.Close()

	usersCol := cSession.DB(Database).C(UsersTable)

	who := bson.M{"_id": profile.ID}
	what := bson.M{"$push": bson.M{"shared_with": sharedUser.ID}}

	return usersCol.Update(who, what)
}
