package storage

import (
    "context"
    "time"
    "github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/models"
)

type Storage interface {
    ListEvents(ctx context.Context, from, to time.Time) ([]*models.Event, error)
    DeleteEvent(ctx context.Context, id string) error
    Close() error
}
