package auth

import (
	"errors"
	"net/http"
	"strings"
)

// GetAPIKey extracts API key
// from the headers of an HTTP request
//
// Example:
//
// Authorization: ApiKey {insert-api-key-here}
func GetAPIKey(headers http.Header) (string, error) {
	
	// Extract the api key, reject in case of invalid authorizatoin token
	user_API_key := headers.Get("Authorization")
	if user_API_key == "" {
		return "", errors.New("No Authentication info found");
	}

	vals := strings.Split(user_API_key, " ");
	if len(vals) != 2 {
		return "", errors.New("Malformed auth Header");
	}
	if vals[0] != "ApiKey" {
		return "", errors.New("Malformed auth Header");
	}

	// Successfully return 
	return vals[1], nil
}