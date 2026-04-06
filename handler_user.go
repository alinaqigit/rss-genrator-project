package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alinaqigit/rss-generator-project/internal/db"
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

	responseWithJson(res, 200, user)
}