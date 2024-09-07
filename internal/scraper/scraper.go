package scraper

import (
	"rent-watcher/internal/models"
	"rent-watcher/internal/notifier"
	"rent-watcher/internal/storage"
)

type Scraper interface {
	Scrape() error
}

type BaseScraper struct {
	Storage  storage.Storage
	Notifier notifier.Notifier
}

func (bs *BaseScraper) ProcessProperty(property *models.Property, rawData string) error {
	isNew, err := bs.Storage.SaveProperty(property, rawData)
	if err != nil {
		return err
	}
	if isNew {
		return bs.Notifier.NotifyNewProperty(property)
	}
	return nil
}
