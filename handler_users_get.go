package main

import (
	"net/http"
	"strings"
)

func (cfg *apiConfig) handlerUserGet(w http.ResponseWriter, r *http.Request) {
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
	respondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}
