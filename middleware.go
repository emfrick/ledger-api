package main

import (
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

// TokenValidationHandler is used to protect routes and ensure tokens are valid
func TokenValidationHandler(db *Database, h AuthorizedHTTPHandlerFunc) http.Handler {

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

		// Make sure the token is valid
		if err == nil && token.Valid {
			claims := token.Claims.(jwt.MapClaims)
			profile, err := db.UAL.GetUserByEmail(claims["email"].(string))

			if err != nil {
				writeJSONToHTTP(w, http.StatusInternalServerError, ResponseError{"Unable to get profile from token"})
				return
			}

			// Forward to the handler with the profile
			h(profile, w, r)
		} else {
			writeJSONToHTTP(w, http.StatusUnauthorized, ResponseError{"Invalid Token"})
		}
	})
}
