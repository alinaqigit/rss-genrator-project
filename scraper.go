package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/alinaqigit/rss-generator-project/internal/db"
	"github.com/google/uuid"
)

func startScraping(
	db *db.Queries,
	concurrency int, 
	timeBetweenRequest time.Duration,
) {
	log.Printf("Scraping on %v goroutines every %s duration", concurrency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)
	for ; ; <- ticker.C {
		feeds, err := db.GetNextFeedsToBeFetched(
			context.Background(),
			int32(concurrency),
		)
		if err != nil {
			log.Println("error fetching feeds: ", err)
			continue
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)

			go scrapFeed(db, wg, feed)
		}
		wg.Wait()
	}
}

func scrapFeed(dbQuries *db.Queries, wg *sync.WaitGroup, feed db.Feed) {
	defer wg.Done()

	_, err := dbQuries.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error marking a feed as fetched: ", err);
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("Error Fetching feed: ", err)
		return
	}

	for _, item := range rssFeed.Channel.Item {

		// description
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}

		// Published date
		pubAt, err := time.Parse(time.RFC1123, item.PubDate)
		if err != nil {
			log.Printf("Could not parse date %v with error %v", item.PubDate, err)
			continue
		}

		_, errr := dbQuries.CreatePost(context.Background(), db.CreatePostParams{
			ID: uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Title: item.Title,
			Description: description,
			PublishedAt: pubAt,
			Url: item.Link,
			FeedID: feed.ID,
		})
		if errr != nil {
			if strings.Contains(errr.Error(), "duplicate key") {
				continue;
			}
			log.Println("Failed to create post", errr)
		}
	}

	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}