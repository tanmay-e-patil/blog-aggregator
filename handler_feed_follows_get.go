package main

import (
	"github.com/tanmay-e-patil/blog-aggregator/internal/database"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerFeedFollowsGetByUser(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := cfg.DB.GetFeedsFollowsByUser(r.Context(), user.ID)
	if err != nil {
		log.Printf("Error getting all feeds: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error getting all feeds")
		return
	}
	respondWithJSON(w, http.StatusOK, databaseFeedFollowsToFeedFollows(feedFollows))
}
