package scraper

import (
	"context"
	"rent-watcher/internal/models"
	"rent-watcher/internal/notifier"
	"rent-watcher/internal/storage"
)

type Scraper interface {
	Scrape(ctx context.Context) error
}

type GeolocationProvider interface {
	CalculateDistance(ctx context.Context, property *models.Property, destLat, destLng float64) (int, error)
}

type BaseScraper struct {
	Storage             storage.Storage
	Notifier            notifier.Notifier
	GeolocationProvider GeolocationProvider
	DestinationLat      float64
	DestinationLng      float64
}

func (bs *BaseScraper) ProcessProperty(ctx context.Context, property *models.Property, rawData string) error {
	if bs.GeolocationProvider != nil {
		distance, err := bs.GeolocationProvider.CalculateDistance(ctx, property, bs.DestinationLat, bs.DestinationLng)
		if err != nil {
			println("Error calculating distance:", err.Error())
		} else {
			property.DistanceMeters = distance
		}
	}

	isNew, err := bs.Storage.SaveProperty(property, rawData)
	if err != nil {
		return err
	}

	if isNew {
		return bs.Notifier.NotifyNewProperty(property)
	}

	return nil
}
