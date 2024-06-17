package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerFeedsGet(w http.ResponseWriter, r *http.Request) {
	feeds, err := cfg.DB.GetAllFeeds(r.Context())
	if err != nil {
		log.Printf("Error getting all feeds: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error getting all feeds")
		return
	}
	respondWithJSON(w, http.StatusOK, feeds)
}
