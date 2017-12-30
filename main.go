package main

// Constants
const (
	DatabaseHost      = "localhost"
	Database          = "ledger"
	UsersTable        = "users"
	TransactionsTable = "transactions"
	SecretKey         = "SECRETKEYISHARDTOGUESS"
	GoogleProfileURL  = "https://www.googleapis.com/userinfo/v2/me"
)

func main() {

	// Instantiate the App
	app := NewApp(DatabaseHost)

	// Run the App
	app.Run(":3000")
}
