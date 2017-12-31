package main

import (
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

func TokenValidationHandler(h AuthorizedHttpHandlerFunc) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		// Make sure the Authorization Header exists
		if authHeader == "" {
			writeErrorToHTTP(w, http.StatusUnauthorized, "Missing Token")
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
			h(claims["email"].(string), w, r)
		} else {
			writeErrorToHTTP(w, http.StatusUnauthorized, "Invalid Token")
		}
	})
}
