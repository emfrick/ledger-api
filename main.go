package main

const (
	DATABASE_HOST      = "localhost"
	DATABASE           = "ledger"
	USERS_TABLE        = "users"
	TRANSACTIONS_TABLE = "transactions"
)

func main() {

	// Instantiate the App
	app := NewApp(DATABASE_HOST)

	// Run the App
	app.Run(":3000")
}
