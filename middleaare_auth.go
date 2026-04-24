package main

import (
	"fmt"
	"net/http"

	"github.com/alinaqigit/rss-generator-project/internal/auth"
	"github.com/alinaqigit/rss-generator-project/internal/db"
)

type authHandler func(http.ResponseWriter, *http.Request, db.User)

func (apicfg *apiConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		apiKey, err := auth.GetAPIKey(req.Header);
		if err != nil { 
			responseWithError(res, 403, fmt.Sprintf("Auth error: %v", err))
			return
		}
		
		user, err := apicfg.DB.GetUserByAPIKey(req.Context(), apiKey)
		if err != nil {
			responseWithError(res, 400, fmt.Sprintf("Couldn't get user %v", err))
			return
		}

		handler(res, req, user)
	}
}