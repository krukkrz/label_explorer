package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/krukkrz/label_explorer/pkg/api"
	"github.com/krukkrz/label_explorer/pkg/api/handlers"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

func main() {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	port := os.Getenv("DB_PORT")

	connStr := fmt.Sprintf("user=%s port=%s password=%s dbname=labels sslmode=disable host=db", dbUser, port, dbPassword)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}
	defer db.Close()

	h := handlers.New(db)
	router := api.NewRouter(h)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"}, // Replace with your frontend URL
		AllowedMethods: []string{"GET", "POST"},           // Allowed methods
		AllowedHeaders: []string{"Content-Type"},          // Allowed headers
	})
	handlerWithCors := c.Handler(router)

	log.Println("Started a server on localhost:8081")
	if err = http.ListenAndServe(":8081", handlerWithCors); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
