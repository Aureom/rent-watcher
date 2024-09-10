package scraper

import (
	"context"
	"fmt"
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
	exists, err := bs.Storage.PropertyExists(property.ID)
	if err != nil {
		return fmt.Errorf("error checking if property exists: %w", err)
	}

	if !exists {
		if bs.GeolocationProvider != nil {
			distance, err := bs.GeolocationProvider.CalculateDistance(ctx, property, bs.DestinationLat, bs.DestinationLng)
			if err != nil {
				return fmt.Errorf("error calculating distance: %w", err)
			}
			property.DistanceMeters = distance
		}

		if err := bs.Notifier.NotifyNewProperty(property); err != nil {
			return fmt.Errorf("error notifying about new property: %w", err)
		}
	} else {
		existingProperty, err := bs.Storage.GetProperty(property.ID)
		if err != nil {
			return fmt.Errorf("error fetching existing property: %w", err)
		}
		property.DistanceMeters = existingProperty.DistanceMeters
	}

	err = bs.Storage.SaveOrUpdateProperty(property, rawData)
	if err != nil {
		return fmt.Errorf("error saving or updating property: %w", err)
	}

	return nil
}
