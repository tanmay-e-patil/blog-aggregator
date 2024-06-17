package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/tanmay-e-patil/blog-aggregator/internal/database"
	"net/http"
	"time"
)

func (cfg *apiConfig) handlerFeedsCreate(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	id, err := uuid.NewUUID()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	feed, err := cfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		Name:      params.Name,
		ID:        id,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Url:       params.Url,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseFeedToFeed(feed))

}
