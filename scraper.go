package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/manuelbamise/rssAggregator/internal/database"
)

func startScraping(db *database.Queries, concurrency int, timeBetweenRequests time.Duration) {
	log.Printf("Logging on %v goroutines every %v time duration", concurrency, timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)

	for ; ; <-ticker.C {
		nextFeedsToFetch, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Println("Error fetching feeds: ", err)
			continue
		}

		wg := &sync.WaitGroup{}
		for _, feed := range nextFeedsToFetch {
			wg.Add(1)

			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error marking feed as fetched", err)
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("Error converting url to feed", err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		description := sql.NullString{}

		if item.Title == "" {
			log.Println("Skipping posts without titles")
			continue
		}

		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}

		postUrl := item.Link
		if postUrl == "" {
			// Generate a unique URL using the feed ID and item title or use a UUID
			postUrl = fmt.Sprintf("%s/post/%s", feed.Url, uuid.New().String())
		}

		pubAt := time.Now().UTC()
		if item.PubDate != "" {
			parsedTime, err := time.Parse(time.RFC1123Z, item.PubDate)
			if err == nil {
				pubAt = parsedTime
			} else {
				log.Printf("Error parsing date '%s': %v", item.PubDate, err)
			}
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Description: description,
			Url:         postUrl,
			PublishedAt: pubAt,
			FeedID:      feed.ID,
		})

		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Println("Error creating new post", err)
		}

	}

	log.Printf("Feeds %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))

}
