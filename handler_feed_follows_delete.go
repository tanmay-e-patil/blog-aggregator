package main

import (
	"github.com/google/uuid"
	"github.com/tanmay-e-patil/blog-aggregator/internal/database"
	"net/http"
)

func (cfg *apiConfig) handlerFeedFollowsDelete(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowIdString := r.PathValue("feedFollowId")
	feedFollowId, err := uuid.Parse(feedFollowIdString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid feed follow id format")
	}

	err = cfg.DB.DeleteFeedFollow(r.Context(), feedFollowId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusNoContent, struct{}{})
}
