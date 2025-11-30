package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/manuelbamise/rssAggregator/internal/database"
)

func (apiCfg *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithERROR(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	newFeedFollow, err := apiCfg.DB.CreateFeedFollow(r.Context(),database.CreateFeedFollowParams{

		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID: params.FeedID,
	})
	if err != nil {
		respondWithERROR(w, 400, fmt.Sprintf("Error creating new feed: %v", err))
		return
	}

	respondWithJSON(w, 201,databaseFeedFollowToFeedFollow(newFeedFollow))
}

// func (apiCfg *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {

// 	feeds, err := apiCfg.DB.GetFeeds(r.Context())
// 	if err != nil {
// 		respondWithERROR(w, 400, fmt.Sprintf("Error getting all feeds: %v", err))
// 		return
// 	}

// 	respondWithJSON(w, 201, databaseFeedsToFeeds(feeds))
// }
