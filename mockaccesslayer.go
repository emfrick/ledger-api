package main

import "errors"
import "gopkg.in/mgo.v2/bson"

// MockAccessLayer provides an in-memory database
type MockAccessLayer struct {
	Users        []GoogleProfile
	Transactions []Transaction
}

// AddUser adds a user to the database
func (mal *MockAccessLayer) AddUser(user GoogleProfile) error {
	mal.Users = append(mal.Users, user)

	return nil
}

// DoesProfileExist checks if a profile exists
func (mal *MockAccessLayer) DoesProfileExist(email string) bool {

	for _, u := range mal.Users {
		if u.Email == email {
			return true
		}
	}

	return false
}

// GetUserByEmail returns a user from the given email
func (mal *MockAccessLayer) GetUserByEmail(email string) (*User, error) {

	for _, u := range mal.Users {
		if u.Email == email {
			user := User{
				ID:         bson.ObjectIdHex(u.ID),
				FirstName:  u.FirstName,
				LastName:   u.LastName,
				Email:      u.Email,
				Gender:     u.Gender,
				SharedWith: []bson.ObjectId{},
			}
			return &user, nil
		}
	}

	return nil, errors.New("User Not Found")
}

// GetUserByID returns a user given the ID
func (mal *MockAccessLayer) GetUserByID(id string) (*User, error) {
	return &User{}, nil
}

// GetSharedUsersForProfile returns the users for the given profile
func (mal *MockAccessLayer) GetSharedUsersForProfile(profile User, out interface{}) {

}

// AddSharedUserToProfile adds a user to the given profile
func (mal *MockAccessLayer) AddSharedUserToProfile(sharedUser User, profile User) error {
	return nil
}

// RemoveSharedUserFromProfile removes a share
func (mal *MockAccessLayer) RemoveSharedUserFromProfile(sharedUser User, profile User) error {
	return nil
}

// StoreTransactions stores the given transactions
func (mal *MockAccessLayer) StoreTransactions(t []Transaction) error {
	return nil
}

// GetTransactionsForProfile gets all the transactions for a given profile
func (mal *MockAccessLayer) GetTransactionsForProfile(profile User, out interface{}) error {
	return nil
}
