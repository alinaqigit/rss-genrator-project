package main

import (
	"database/sql"
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

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	}))

	//v1
	v1router := chi.NewRouter()
	v1router.Get("/healthz", handlerReadiness)
	v1router.Get("/err", handlerErr)

	v1router.Post("/user", apiCfg.handlerCreateUser)

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