package app

import (
    "context"
    "time"

    "github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/logger"
    "github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/models"
    "github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/storage"
)

type App struct {
    logger  *logger.Logger
    storage storage.Storage
}

func New(logger *logger.Logger, storage storage.Storage) *App {
    return &App{
        logger:  logger,
        storage: storage,
    }
}

func (a *App) ListEvents(ctx context.Context, from, to time.Time) ([]*models.Event, error) {
    return a.storage.ListEvents(ctx, from, to)
}

func (a *App) DeleteEvent(ctx context.Context, id string) error {
    return a.storage.DeleteEvent(ctx, id)
}
