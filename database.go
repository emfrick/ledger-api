package main

// Database abstracts away perstistant storage
type Database struct {
	UAL UserAccessLayer
	TAL TransactionAccessLayer
}

// UserAccessLayer provides access to user objects
type UserAccessLayer interface {
	AddUser(user GoogleProfile) error
	DoesProfileExist(email string) bool
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id string) (*User, error)
	GetSharedUsersForProfile(profile User, out interface{})
	AddSharedUserToProfile(sharedUser User, profile User) error
	RemoveSharedUserFromProfile(sharedUser User, profile User) error
}

// TransactionAccessLayer provides access to transaction objects
type TransactionAccessLayer interface {
	StoreTransactions(t []Transaction) error
	GetTransactionsForProfile(profile User, out interface{}) error
}
