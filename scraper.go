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

			go scrapeFeed(wg)
		}
		wg.Wait()
	}
}

func scrapeFeed(wg *sync.WaitGroup) {
	defer wg.Done()
}
