package main

type Database struct {
	UAL UserAccessLayer
	TAL TransactionAccessLayer
}

type UserAccessLayer interface {
	AddUser(user GoogleProfile) error
	DoesProfileExist(email string) bool
	GetUserByEmail(email string) (*User, error)
	GetUserById(id string) (*User, error)
	GetSharedUsersForProfile(profile User, out interface{})
	AddSharedUserToProfile(sharedUser User, profile User) error
	RemoveSharedUserFromProfile(sharedUser User, profile User) error
}

type TransactionAccessLayer interface {
	StoreTransactions(t []Transaction) error
	GetTransactionsForProfile(profile User, out interface{}) error
}
