package main

import (
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	mgo "gopkg.in/mgo.v2"
)

func TokenValidationHandler(session *mgo.Session, h AuthorizedHttpHandlerFunc) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		// Make sure the Authorization Header exists
		if authHeader == "" {
			writeJSONToHTTP(w, http.StatusUnauthorized, ResponseError{"Missing Token"})
			return
		}

		// Parse the token from Bearer xxxx
		keys := strings.Split(authHeader, " ")
		tokenString := keys[1]

		// Get the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		})

		if err == nil && token.Valid {
			claims := token.Claims.(jwt.MapClaims)
			profile := getUserByEmail(session, claims["email"].(string))
			h(profile, w, r)
		} else {
			writeJSONToHTTP(w, http.StatusUnauthorized, ResponseError{"Invalid Token"})
		}
	})
}
