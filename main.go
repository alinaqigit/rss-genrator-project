package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"

	"github.com/alinaqigit/rss-generator-project/internal/db"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *db.Queries
}

func main() {

	feed, err := urlToFeed("https://wagslane.dev/index.xml")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(feed)

	// load env
	godotenv.Load()

	// extract env
	PORT := os.Getenv("PORT")
	HOST := os.Getenv("HOST")
	DB_URL := os.Getenv("DATABASE_URL")
	
	// validate env
	if PORT == "" {
		log.Fatal("PORT is not set in the environment variables")
	}
	if DB_URL == "" {
		log.Fatal("DATABASE_URL is not set in the environment variables")
	}
	

	// Opening connection to database
	conn, err := sql.Open("postgres", DB_URL);
	if(err != nil){
		log.Fatal("Failed to connect to database. Err:", err);
	}

	// Getting qureies
	queries := db.New(conn)

	apiCfg := apiConfig{
		DB: queries,
	}

	// Server

	router := chi.NewRouter();
	// Cors
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	}))

	// 1. v1
	v1router := chi.NewRouter()
	v1router.Get("/healthz", handlerReadiness)
	v1router.Get("/err", handlerErr)

	// User
	v1router.Post("/user", apiCfg.handlerCreateUser)
	v1router.Get("/user", apiCfg.middlewareAuth(apiCfg.handlerGetUserByAPI))
	v1router.Delete("/user", apiCfg.middlewareAuth(apiCfg.handlerDeactivateUser))

	// Feed
	v1router.Post("/feed", apiCfg.middlewareAuth(apiCfg.handlerCreateUserFeed))
	v1router.Get("/feed", apiCfg.handlerGetAllFeeds)
	
	// Feed Follow 
	v1router.Post("/feed-follow", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
	v1router.Get("/feed-follow", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollows))
	v1router.Delete("/feed-follow/{feed-follow-id}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollows))


	// Server Stuff

	router.Mount("/v1", v1router)

	server := &http.Server{
		Handler: router,
		Addr: HOST + ":" + PORT,
	}

	log.Printf("Application Starting on http://%v:%v", HOST, PORT)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}

}