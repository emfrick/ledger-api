package main

import (
	"log"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

func TokenValidationHandler(h http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			writeErrorToHttp(w, http.StatusUnauthorized, "Missing Token")
			return
		}

		keys := strings.Split(authHeader, " ")
		tokenString := keys[1]

		log.Println(tokenString)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		})

		if err == nil && token.Valid {
			h.ServeHTTP(w, r)
		} else {
			writeErrorToHttp(w, http.StatusUnauthorized, "Invalid Token")
		}
	})
}
