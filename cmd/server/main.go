package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/XAER/osu-ha-integration/internal/api"
	"github.com/XAER/osu-ha-integration/internal/config"
	"github.com/XAER/osu-ha-integration/internal/osu"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	if errEnv := godotenv.Load(); errEnv != nil {
		log.Fatal("No .env file found")
	}

	conf := config.LoadConfig("/config.yaml")

	clientID := os.Getenv("OSU_CLIENT_ID")
	clientSecret := os.Getenv("OSU_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		log.Fatal("OSU_CLIENT_ID and OSU_CLIENT_SECRET must be set")
	}

	logger := &osu.StdLogger{}

	client := osu.NewClient(clientID, clientSecret)
	cache := osu.NewCache(time.Duration(conf.Cache.Duration)*time.Second, logger)
	handler := api.NewHandler(client, cache)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/stats", handler.GetUserStats)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Printf("Starting osu-ha-integration on :%s \n", conf.Server.Port)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", conf.Server.Port), r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
