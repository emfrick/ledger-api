package main

import "net/http"

func (a *App) indexHandler(w http.ResponseWriter, r *http.Request) {
	session := a.Session.Copy()
	defer session.Close()

	dbNames, err := session.DatabaseNames()

	if err != nil {
		writeErrorToHttp(w, http.StatusInternalServerError, "Unable to Get Database Names")
		return
	}

	writeJsonToHttp(w, dbNames)
}

func (a *App) errorHandler(w http.ResponseWriter, r *http.Request) {
	session := a.Session.Copy()
	defer session.Close()

	writeErrorToHttp(w, http.StatusInternalServerError, "This is the error route")
}

func (a *App) usersHandler(w http.ResponseWriter, r *http.Request) {
	var users []User

	getAllObjectsFromTable(a.Session, USERS_TABLE, &users)

	writeJsonToHttp(w, users)
}
