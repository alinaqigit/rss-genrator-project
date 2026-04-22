package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alinaqigit/rss-generator-project/internal/db"
	"github.com/alinaqigit/rss-generator-project/internal/db/auth"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateUser(res http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(req.Body);

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responseWithError(res, 400, fmt.Sprint("Error Parsing json", err))
	}

	user, err := apiCfg.DB.CreateUser(req.Context(), db.CreateUserParams{
		ID: uuid.New(),
		Username: params.Name,
	})
	if err != nil {
		responseWithError(res, 400, fmt.Sprint("Couldn't create user", err))
		return
	}

	responseWithJson(res, 201, user)
}

func (apiCfg *apiConfig) handlerGetUserByAPI(res http.ResponseWriter, req *http.Request) {
	apiKey, err := auth.GetAPIKey(req.Header);
	if err != nil {
		responseWithError(res, 403, fmt.Sprintf("Auth error: %v", err))
		return
	}
	
	user, err := apiCfg.DB.GetUserByAPIKey(req.Context(), apiKey)
	if err != nil {
		responseWithError(res, 400, fmt.Sprintf("Couldn't get user"))
		return
	}

	responseWithJson(res, 200, database_user_to_User(user));
}