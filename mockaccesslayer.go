package main

import "errors"
import "gopkg.in/mgo.v2/bson"

type MockAccessLayer struct {
	Users        []GoogleProfile
	Transactions []Transaction
}

func (mal *MockAccessLayer) AddUser(user GoogleProfile) error {
	mal.Users = append(mal.Users, user)

	return nil
}

func (mal *MockAccessLayer) DoesProfileExist(email string) bool {
	return true
}

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

func (mal *MockAccessLayer) GetUserById(id string) (*User, error) {
	return &User{}, nil
}

func (mal *MockAccessLayer) GetSharedUsersForProfile(profile User, out interface{}) {

}

func (mal *MockAccessLayer) AddSharedUserToProfile(sharedUser User, profile User) error {
	return nil
}

func (mal *MockAccessLayer) RemoveSharedUserFromProfile(sharedUser User, profile User) error {
	return nil
}

func (mal *MockAccessLayer) StoreTransactions(t []Transaction) error {
	return nil
}

func (mal *MockAccessLayer) GetTransactionsForProfile(profile User, out interface{}) error {
	return nil
}
