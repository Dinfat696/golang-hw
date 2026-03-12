package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/models"
	"github.com/google/uuid"
)

type MemoryStorage struct {
    mu     sync.RWMutex
    events map[string]*models.Event
}

func NewStorage() *MemoryStorage {
    return &MemoryStorage{
        events: make(map[string]*models.Event),
    }
}


func (s *MemoryStorage) CreateEvent(ctx context.Context, event *models.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, e := range s.events {
		if e.StartTime.Equal(event.StartTime) && e.UserID == event.UserID {
			return models.ErrDateBusy
		}
	}

	event.ID = uuid.New().String()
	s.events[event.ID] = event
	return nil
}

func (s *MemoryStorage) UpdateEvent(ctx context.Context, event *models.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.events[event.ID]; !exists {
		return models.ErrEventNotFound
	}

	for _, e := range s.events {
		if e.ID != event.ID && e.StartTime.Equal(event.StartTime) && e.UserID == event.UserID {
			return models.ErrDateBusy
		}
	}

	s.events[event.ID] = event
	return nil
}

func (s *MemoryStorage) DeleteEvent(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.events[id]; !exists {
		return models.ErrEventNotFound
	}

	delete(s.events, id)
	return nil
}

func (s *MemoryStorage) GetEvent(ctx context.Context, id string) (*models.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	event, exists := s.events[id]
	if !exists {
		return nil, models.ErrEventNotFound
	}

	return event, nil
}

func (s *MemoryStorage) ListEvents(ctx context.Context, from, to time.Time) ([]*models.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var events []*models.Event
	for _, event := range s.events {
		if (event.StartTime.After(from) || event.StartTime.Equal(from)) &&
			(event.StartTime.Before(to) || event.StartTime.Equal(to)) {
			events = append(events, event)
		}
	}

	return events, nil
}

func (s *MemoryStorage) Close() error {
	return nil
}
