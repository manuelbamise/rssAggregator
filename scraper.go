package main

import (
	"context"
	"log"
	"sync"
	"time"

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
		log.Println("Found post: ", item.Title)
	}

	log.Printf("Feeds %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))

}
