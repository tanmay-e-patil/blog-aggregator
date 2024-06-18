package main

import (
	"database/sql"
	"github.com/joho/godotenv"
	"github.com/tanmay-e-patil/blog-aggregator/internal/database"
	"log"
	"net/http"
	"os"
	"time"
)
import _ "github.com/lib/pq"

type apiConfig struct {
	DB *database.Queries
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		log.Fatal("$DATABASE_URL must be set")
	}

	db, err := sql.Open("postgres", databaseUrl)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)

	apiCfg := apiConfig{
		dbQueries,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/healthz", handlerReadiness)
	mux.HandleFunc("GET /v1/err", handlerErr)

	mux.HandleFunc("POST /v1/users", apiCfg.handlerUsersCreate)
	mux.HandleFunc("GET /v1/users", apiCfg.handlerUserGet)

	mux.HandleFunc("POST /v1/feeds", apiCfg.middlewareAuth(apiCfg.handlerFeedsCreate))
	mux.HandleFunc("GET /v1/feeds", apiCfg.handlerFeedsGet)

	mux.HandleFunc("POST /v1/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerFeedFollowsCreate))
	mux.HandleFunc("GET /v1/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerFeedFollowsGetByUser))
	mux.HandleFunc("DELETE /v1/feed_follows/{feedFollowId}", apiCfg.middlewareAuth(apiCfg.handlerFeedFollowsDelete))

	mux.HandleFunc("GET /v1/posts", apiCfg.middlewareAuth(apiCfg.handlerPostsGet))

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	const collectionConcurrency = 10
	const collectionInterval = time.Minute
	go startScraping(dbQueries, collectionConcurrency, collectionInterval)

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
