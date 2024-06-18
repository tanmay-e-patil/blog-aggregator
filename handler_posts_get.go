package main

import (
	"github.com/tanmay-e-patil/blog-aggregator/internal/database"
	"net/http"
	"strconv"
)

func (cfg *apiConfig) handlerPostsGet(w http.ResponseWriter, r *http.Request, user database.User) {
	limitString := r.URL.Query().Get("limit")
	var limit int
	if limitString != "" {
		limit, _ = strconv.Atoi(limitString)
	} else {
		limit = 10
	}

	posts, err := cfg.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, databasePostsToPosts(posts))
}
