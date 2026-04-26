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

	responseWithJson(res, 201, database_user_to_User(user))
}

func (apiCfg *apiConfig) handlerGetUserByAPI(res http.ResponseWriter, req *http.Request, user db.User) {
	responseWithJson(res, 200, database_user_to_User(user));
}

func (apiCfg *apiConfig) handlerGetUserPosts(res http.ResponseWriter, req *http.Request, user db.User) {
	posts, err := apiCfg.DB.GetPostsForUser(req.Context(), db.GetPostsForUserParams{
		UserID: user.ID,
		Limit: 10,
	})
	if err != nil {
		responseWithError(res, 400, fmt.Sprint("Couldn't get posts for user", err))
		return
	}

	responseWithJson(res, 200, database_posts_to_Posts(posts))

}

func (apiCfg *apiConfig) handlerDeactivateUser (res http.ResponseWriter, req *http.Request, user db.User) {

	apiCfg.DB.DeleteUserByAPIKey(req.Context(), user.ApiKey)

	responseWithJson(res, 204, struct{}{})
}