package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/alinaqigit/rss-generator-project/internal/db"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateFeedFollow(res http.ResponseWriter, req *http.Request, user db.User) {

	type parameters struct {
		FeedID	uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if(err != nil) {
		responseWithError(res, 500, "Error Parsing Json")
		return
	}

	FeedFollow, err := apiCfg.DB.CreateFeedFollow(req.Context(), db.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FeedID: params.FeedID,
		UserID: user.ID,
	})
	if(err != nil) {
		responseWithError(res, 500, "Error Creating a Feed Follow")
		return
	}

	responseWithJson(res, 201, database_feedFollow_to_FeedFollow(FeedFollow))
}

func (apiCfg *apiConfig) handlerGetFeedFollows(res http.ResponseWriter, req *http.Request, user db.User) {
	FeedFollows, err := apiCfg.DB.GetFeedFollowByUserID(req.Context(), user.ID)
	if(err != nil) {
		responseWithError(res, 400, "Holy smokes! something went wrong")
		return
	}

	responseWithJson(res, 200, database_feedFollowSlice_to_FeedFollowSlice(FeedFollows))
}

func (apiCfg *apiConfig) handlerDeleteFeedFollows(res http.ResponseWriter, req *http.Request, user db.User) {

	idstr := chi.URLParam(req, "feed-follow-id")
	id, err := uuid.Parse(idstr)
	if err != nil {
		responseWithError(res, 400, "Could not parse id")
		return
	}

	err = apiCfg.DB.DeleteFeedFollow(req.Context(), db.DeleteFeedFollowParams{
		ID: id,
		UserID: user.ID,
	})
	if(err != nil) {
		responseWithError(res, 400, "Holy smokes! something went wrong")
		return
	}

	responseWithJson(res, 204, struct{}{})
}