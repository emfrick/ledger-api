package main

const (
	DATABASE_HOST      = "localhost"
	DATABASE           = "ledger"
	USERS_TABLE        = "users"
	TRANSACTIONS_TABLE = "transactions"
	SECRET_KEY         = "SECRETKEYISHARDTOGUESS"
	GOOGLE_PROFILE_URL = "https://www.googleapis.com/userinfo/v2/me"
)

func main() {

	// Instantiate the App
	app := NewApp(DATABASE_HOST)

	// Run the App
	app.Run(":3000")
}
