package main

import (
	"log"
	"rent-watcher/internal/config"
	"rent-watcher/internal/database"
	"rent-watcher/internal/discord"
	"rent-watcher/internal/scraper"
	"rent-watcher/internal/storage"
)

func main() {
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

	arantesScraper := scraper.NewArantesScraper(scraper.ArantesConfig(cfg.ArantesConfig), store, discord)

	if err := arantesScraper.Scrape(); err != nil {
		log.Fatalf("Failed to scrape Arantes: %v", err)
	}
}
