package main

import (
	"log"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MongoAccessLayer provides access to a Mongo DB
type MongoAccessLayer struct {
	Session *mgo.Session
}

// AddUser adds a user to the database
func (mal MongoAccessLayer) AddUser(user GoogleProfile) error {
	// Copy the mongo session and defer its close
	cSession := mal.Session.Copy()
	defer cSession.Close()

	log.Printf("Adding User: %v\n", user)

	c := cSession.DB(DatabaseName).C(UsersTable)
	err := c.Insert(user)

	return err
}

// DoesProfileExist checks the database for the given email
func (mal MongoAccessLayer) DoesProfileExist(email string) bool {

	// Copy the mongo session and defer its close
	cSession := mal.Session.Copy()
	defer cSession.Close()

	col := cSession.DB(DatabaseName).C(UsersTable)
	count, err := col.Find(bson.M{"email": email}).Count()

	if err != nil {
		log.Println(err)
	}

	if count > 0 {
		return true
	}

	return false
}

// GetUserByEmail will return a user object given an email
func (mal MongoAccessLayer) GetUserByEmail(email string) (*User, error) {
	var user User

	// Copy the mongo session and defer its close
	cSession := mal.Session.Copy()
	defer cSession.Close()

	c := cSession.DB(DatabaseName).C(UsersTable)
	err := c.Find(bson.M{"email": email}).One(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByID returns a user given the ID
func (mal MongoAccessLayer) GetUserByID(id string) (*User, error) {
	var user User

	// Copy the mongo session and defer its close
	cSession := mal.Session.Copy()
	defer cSession.Close()

	c := cSession.DB(DatabaseName).C(UsersTable)
	err := c.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetSharedUsersForProfile returns the users for the given profile
func (mal MongoAccessLayer) GetSharedUsersForProfile(profile User, out interface{}) {
	// Copy the mongo session and defer its close
	cSession := mal.Session.Copy()
	defer cSession.Close()

	c := cSession.DB(DatabaseName).C(UsersTable)

	// Return the user OR and accounts that share with the user
	query := bson.M{"$or": []bson.M{bson.M{"shared_with": profile.ID}, bson.M{"_id": profile.ID}}}

	c.Find(query).All(out)
}

// AddSharedUserToProfile adds a user to the given profile
func (mal MongoAccessLayer) AddSharedUserToProfile(sharedUser User, profile User) error {
	cSession := mal.Session.Copy()
	defer cSession.Close()

	usersCol := cSession.DB(DatabaseName).C(UsersTable)

	who := bson.M{"_id": profile.ID}
	what := bson.M{"$push": bson.M{"shared_with": sharedUser.ID}}

	return usersCol.Update(who, what)
}

// RemoveSharedUserFromProfile removes a share
func (mal MongoAccessLayer) RemoveSharedUserFromProfile(sharedUser User, profile User) error {
	cSession := mal.Session.Copy()
	defer cSession.Close()

	usersCol := cSession.DB(DatabaseName).C(UsersTable)

	who := bson.M{"_id": profile.ID}
	what := bson.M{"$pull": bson.M{"shared_with": sharedUser.ID}}

	return usersCol.Update(who, what)
}

// StoreTransactions stores the given transactions
func (mal MongoAccessLayer) StoreTransactions(t []Transaction) error {

	// Copy the mongo session and defer its close
	cSession := mal.Session.Copy()
	defer cSession.Close()

	c := cSession.DB(DatabaseName).C(TransactionsTable)

	for _, transaction := range t {
		if err := c.Insert(transaction); err != nil {
			log.Printf("Error: %v\n", err)
			return err
		}
	}

	return nil
}

// GetTransactionsForProfile gets all the transactions for a given profile
func (mal MongoAccessLayer) GetTransactionsForProfile(profile User, out interface{}) error {

	// Copy the mongo session and defer its close
	cSession := mal.Session.Copy()
	defer cSession.Close()

	transactionsCol := cSession.DB(DatabaseName).C(TransactionsTable)
	usersCol := cSession.DB(DatabaseName).C(UsersTable)

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
