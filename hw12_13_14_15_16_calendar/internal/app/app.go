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

// Временные заглушки для совместимости с HTTP-хендлерами
func (a *App) FindEventsByDay(ctx context.Context, day time.Time) ([]*models.Event, error) {
	return a.ListEvents(ctx, day, day.Add(24*time.Hour))
}

func (a *App) FindEventsByWeek(ctx context.Context, start time.Time) ([]*models.Event, error) {
	return a.ListEvents(ctx, start, start.Add(7*24*time.Hour))
}

func (a *App) FindEventsByMonth(ctx context.Context, start time.Time) ([]*models.Event, error) {
	return a.ListEvents(ctx, start, start.AddDate(0, 1, 0))
}


// GetEvent возвращает событие по ID
func (a *App) GetEvent(ctx context.Context, id string) (*models.Event, error) {
    return a.storage.GetEvent(ctx, id)
}

// UpdateEvent обновляет существующее событие
func (a *App) UpdateEvent(ctx context.Context, event *models.Event) error {
    return a.storage.UpdateEvent(ctx, event)
}

// CreateEvent создаёт новое событие (реализуйте вызов storage)
func (a *App) CreateEvent(ctx context.Context, event *models.Event) error {
    // Здесь можно добавить валидацию, логирование и т.п.
    return a.storage.CreateEvent(ctx, event)
}
