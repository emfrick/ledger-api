package main_test

import (
	"fmt"
	"testing"

	"."
)

func TestGetUserByEmail(t *testing.T) {

	// Create a Mock Access Layer
	mal := &main.MockAccessLayer{
		Users:        []main.GoogleProfile{},
		Transactions: []main.Transaction{},
	}

	// Create a Database
	db := main.Database{
		UAL: mal,
		TAL: mal,
	}

	// Instantiate the App
	app := main.NewApp(&db)

	// Create a fake user
	emfrick := main.GoogleProfile{
		ID:            "000000000000000000000000",
		Email:         "emfrick@gmail.com",
		VerifiedEmail: true,
		FullName:      "Eric Frick",
		FirstName:     "Eric",
		LastName:      "Frick",
		ProfileLink:   "http://www.example.com",
		Picture:       "http://www.example",
		Gender:        "male",
		Locale:        "en",
	}

	joeuser := main.GoogleProfile{
		ID:            "111111111111111111111111",
		Email:         "joeuser@gmail.com",
		VerifiedEmail: true,
		FullName:      "Joe User",
		FirstName:     "Joe",
		LastName:      "User",
		ProfileLink:   "http://www.example.com",
		Picture:       "http://www.example",
		Gender:        "male",
		Locale:        "en",
	}

	err := app.Database.UAL.AddUser(emfrick)
	err = app.Database.UAL.AddUser(joeuser)

	profile, err := app.Database.UAL.GetUserByEmail(joeuser.Email)

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	fmt.Printf("USER: %v\n", profile)
}
