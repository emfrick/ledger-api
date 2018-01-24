package main_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"."
)

var app *main.App

type route struct {
	method string
	path   string
}

var protectedRoutes = []route{
	{"GET", "/shared"},
	{"POST", "/shared"},
	{"DELETE", "/shared/1"},
	{"GET", "/transactions"},
	{"POST", "/transactions"},
}

func TestMain(m *testing.M) {
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
	app = main.NewApp(&db)

	result := m.Run()

	os.Exit(result)
}

func TestProtectedRoutes(t *testing.T) {

	// Define the expectred response code
	expected := http.StatusUnauthorized

	// Check that all protected routes return the expected response code
	for _, route := range protectedRoutes {
		req, _ := http.NewRequest(route.method, route.path, nil)
		rr := httptest.NewRecorder()
		app.Router.ServeHTTP(rr, req)

		if expected != rr.Code {
			t.Errorf("Expected route %s to return response code %d. Got %d\n", route.path, expected, rr.Code)
		}
	}
}
