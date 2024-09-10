package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"rent-watcher/internal/config"
	"rent-watcher/internal/database"
	"rent-watcher/internal/discord"
	"rent-watcher/internal/geolocation"
	"rent-watcher/internal/scraper"
	"rent-watcher/internal/storage"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Received termination signal. Initiating graceful shutdown...")
		cancel()
	}()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db, err := database.Init(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	store := storage.NewSQLStorage(db)

	discord, err := discord.New(cfg.DiscordToken, cfg.DiscordChannel)
	if err != nil {
		log.Fatalf("Failed to initialize Discord bot: %v", err)
	}
	defer discord.Close()

	geoProvider := geolocation.NewGoogleMapsClient(cfg.GoogleMapsAPIKey)

	arantesScraper := scraper.NewArantesScraper(
		scraper.ArantesConfig(cfg.ArantesConfig),
		cfg.DestinationLat,
		cfg.DestinationLng,
		store,
		discord,
		geoProvider,
	)

	scraperCtx, scraperCancel := context.WithTimeout(ctx, 30*time.Minute)
	defer scraperCancel()

	if err := arantesScraper.Scrape(scraperCtx); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Println("Scraper timed out")
		} else if errors.Is(err, context.Canceled) {
			log.Println("Scraper was cancelled")
		} else {
			log.Printf("Failed to scrape Arantes: %v", err)
		}
	}

	log.Println("Scraping completed. Shutting down...")
}
