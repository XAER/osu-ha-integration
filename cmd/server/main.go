package main

import (
	"log"
	"net/http"
	"os"

	"github.com/XAER/osu-ha-integration/internal/api"
	"github.com/XAER/osu-ha-integration/internal/osu"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	if errEnv := godotenv.Load(); errEnv != nil {
		log.Fatal("No .env file found")
	}

	token := os.Getenv("OSU_API_TOKEN")
	if token == "" {
		log.Fatal("OSU_API_TOKEN environment variable is required")
	}

	client := osu.NewClient(token)
	handler := api.NewHandler(client)

	r := chi.NewRouter()
	r.Get("/stats", handler.GetUserStats)

	log.Println("Starting osu-ha-integration on :8087")

	if err := http.ListenAndServe(":8087", r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
