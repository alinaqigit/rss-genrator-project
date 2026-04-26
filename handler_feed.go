package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/alinaqigit/rss-generator-project/internal/db"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateUserFeed(res http.ResponseWriter, req *http.Request, user db.User) {
	type parameters struct {
		Name string `json:"name"`
		Url string `json:"url"`
	}

	decoder := json.NewDecoder(req.Body);

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responseWithError(res, 400, fmt.Sprint("Error Parsing json", err))
	}

	feed, err := apiCfg.DB.CreateFeed(req.Context(), db.CreateFeedParams{
		ID: uuid.New(),
		Name: params.Name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Url: params.Url,
		UserID: user.ID,
	})
	if err != nil {
		responseWithError(res, 400, fmt.Sprint("Couldn't create user", err))
		return
	}

	responseWithJson(res, 201, database_feed_to_Feed(feed))
}

func (apiCfg *apiConfig) handlerGetAllFeeds(res http.ResponseWriter, req *http.Request) {
	Feeds, err := apiCfg.DB.GetAllFeeds(req.Context())
	if(err != nil) {
		responseWithError(res, 500, "Failded to get Feeds")
	}

	responseWithJson(res, 200, database_feedSlice_to_FeedSlice(Feeds));
}

