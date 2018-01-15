package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Takes an access token and returns the Google Profile
func getProfileFromGoogle(accessToken string) (*GoogleProfile, error) {

	// Create the URL and run an http GET
	url := fmt.Sprintf("%s?access_token=%s", GoogleProfileURL, accessToken)
	response, err := http.Get(url)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Decode the response body
	decoder := json.NewDecoder(response.Body)

	// Make sure Google didn't return 401
	if response.StatusCode == http.StatusUnauthorized {
		var googleOauthError GoogleOauthError
		err = decoder.Decode(&googleOauthError)
		return nil, googleOauthError
	}

	// Decode the JSON
	var profile GoogleProfile
	err = decoder.Decode(&profile)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &profile, nil
}

// Writes JSON to the response
func writeJSONToHTTP(w http.ResponseWriter, code int, objects interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	encoder := json.NewEncoder(w)
	encoder.Encode(objects)
}

// Custom Error
func (e GoogleOauthError) Error() string {
	return e.Err.Message
}
