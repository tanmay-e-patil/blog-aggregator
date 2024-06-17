package main

import (
	"github.com/tanmay-e-patil/blog-aggregator/internal/database"
	"net/http"
	"strings"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	///
	return func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		apiKey, found := strings.CutPrefix(authorizationHeader, "ApiKey ")
		if !found {
			respondWithError(w, http.StatusUnauthorized, "Invalid API key")
			return
		}

		user, err := cfg.DB.GetUserByApiKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		handler(w, r, user)
	}

}
