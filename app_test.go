package main_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"

	"."
	mgo "gopkg.in/mgo.v2"
)

func TestErrorRouteReturns500(t *testing.T) {
	a := main.App{}

	a.Router = mux.NewRouter()
	a.Session, _ = mgo.Dial("localhost")

	rr := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/error", nil)
	a.Router.ServeHTTP(rr, req)
}
