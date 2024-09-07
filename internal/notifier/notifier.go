package notifier

import (
	"rent-watcher/internal/models"
)

type Notifier interface {
	NotifyNewProperty(property *models.Property) error
	Close() error
}
