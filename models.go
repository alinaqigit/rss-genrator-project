package main

import (
	"time"

	"github.com/alinaqigit/rss-generator-project/internal/db"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Username  string    `json:"name"`
	ApiKey		string		`json:"api_key"`
}

type Feed struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name 			string		`json:"name"`
	Url				string 		`json:"url"`
	UserID		uuid.UUID	`json:"user_id"`
}

type FeedFollow struct {
	ID        uuid.UUID
  CreatedAt time.Time
  UpdatedAt time.Time
  UserID    uuid.UUID
  FeedID    uuid.UUID
}

func database_user_to_User(databaseUser db.User) User {
	return User{
		ID: databaseUser.ID,
		CreatedAt: databaseUser.CreatedAt,
    UpdatedAt: databaseUser.UpdatedAt,
    Username: databaseUser.Username,
		ApiKey: databaseUser.ApiKey,
	}
}

func database_feed_to_Feed(databaseFeed db.Feed) Feed {
	return Feed{
		ID: databaseFeed.ID,
		CreatedAt: databaseFeed.CreatedAt,
    UpdatedAt: databaseFeed.UpdatedAt,
    Name: databaseFeed.Name,
		Url: databaseFeed.Url,
		UserID: databaseFeed.UserID,
	}
}

func database_feedSlice_to_FeedSlice(databaseFeed []db.Feed) []Feed {
	feeds := []Feed{}

	for _, dbFeed := range databaseFeed {
		feeds = append(feeds, database_feed_to_Feed(dbFeed))
	}

	return feeds
}

func database_feedFollow_to_FeedFollow(dbFeedFollow db.FeedFollow) FeedFollow {
	return FeedFollow{
		ID: dbFeedFollow.ID,
		CreatedAt: dbFeedFollow.CreatedAt,
		UpdatedAt: dbFeedFollow.UpdatedAt,
		UserID: dbFeedFollow.UserID,
		FeedID: dbFeedFollow.FeedID,
	}
}

func database_feedFollowSlice_to_FeedFollowSlice(databaseFeedFollow []db.FeedFollow) []FeedFollow {
	feedFollows := []FeedFollow{}

	for _, dbFeed := range databaseFeedFollow {
		feedFollows = append(feedFollows, database_feedFollow_to_FeedFollow(dbFeed))
	}

	return feedFollows
}
