package main

import (
	"crypto/rand"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/tanmay-e-patil/blog-aggregator/internal/database"
	"log"
	"net/http"
	"time"
)

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
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

	apiKeyBytes := make([]byte, 32)
	_, err = rand.Read(apiKeyBytes)
	if err != nil {
		return
	}

	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		Name:      params.Name,
		ID:        id,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		log.Printf("Error creating user: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseUserToUser(user))
}
