package app

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/config"
	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/logger"
	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/models"
	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/mq"
	"github.com/google/uuid"
)

type Scheduler struct {
	app      *App
	producer mq.Producer
	logger   *logger.Logger
	config   config.SchedulerConfig
}

func NewScheduler(app *App, producer mq.Producer, logger *logger.Logger, config config.SchedulerConfig) *Scheduler {
	return &Scheduler{
		app:      app,
		producer: producer,
		logger:   logger,
		config:   config,
	}
}

func (s *Scheduler) Run(ctx context.Context) error {
	ticker := time.NewTicker(s.config.Interval)
	defer ticker.Stop()

	if err := s.processNotifications(ctx); err != nil {
		s.logger.Error(fmt.Sprintf("Failed to process notifications: %v", err))
	}
	if err := s.cleanupOldEvents(ctx); err != nil {
		s.logger.Error(fmt.Sprintf("Failed to cleanup old events: %v", err))
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := s.processNotifications(ctx); err != nil {
				s.logger.Error(fmt.Sprintf("Failed to process notifications: %v", err))
			}
			if err := s.cleanupOldEvents(ctx); err != nil {
				s.logger.Error(fmt.Sprintf("Failed to cleanup old events: %v", err))
			}
		}
	}
}

func (s *Scheduler) processNotifications(ctx context.Context) error {
	now := time.Now()
	from := now
	to := now.Add(s.config.Interval)

	events, err := s.app.ListEvents(ctx, from, to)
	if err != nil {
		return fmt.Errorf("list events: %w", err)
	}

	sentCount := 0
	for _, event := range events {
		if s.shouldNotify(event, now) {
			notification := &models.Notification{
				ID:         uuid.New().String(),
				EventID:    strconv.FormatInt(event.ID, 10),
				EventTitle: event.Title,
				UserID:     "user", // TODO: заменить на реальный userID из event
				Message:    fmt.Sprintf("Напоминание: %s начинается в %s", event.Title, event.DateTime.Format("15:04")),
				NotifyAt:   time.Now(),
				CreatedAt:  time.Now(),
			}

			if err := s.producer.SendNotification(ctx, notification); err != nil {
				s.logger.Error(fmt.Sprintf("Failed to send notification for event %d: %v", event.ID, err))
				continue
			}

			sentCount++
			s.logger.Info(fmt.Sprintf("Notification sent for event: %s", event.Title))
		}
	}

	s.logger.Info(fmt.Sprintf("Processed %d events, sent %d notifications", len(events), sentCount))
	return nil
}

func (s *Scheduler) shouldNotify(event *models.Event, now time.Time) bool {
	if event.DateTime.IsZero() {
		return false
	}
	return event.DateTime.After(now) && event.DateTime.Before(now.Add(s.config.Interval))
}

func (s *Scheduler) cleanupOldEvents(ctx context.Context) error {
	cutoffTime := time.Now().Add(-s.config.CleanupOlderThan)

	oldEvents, err := s.app.ListEvents(ctx, time.Time{}, cutoffTime)
	if err != nil {
		return fmt.Errorf("list old events: %w", err)
	}

	deletedCount := 0
	for _, event := range oldEvents {
		if err := s.app.DeleteEvent(ctx, strconv.FormatInt(event.ID, 10)); err != nil {
			s.logger.Error(fmt.Sprintf("Failed to delete old event %d: %v", event.ID, err))
			continue
		}
		deletedCount++
		s.logger.Info(fmt.Sprintf("Deleted old event: %s", event.Title))
	}

	if deletedCount > 0 {
		s.logger.Info(fmt.Sprintf("Cleaned up %d old events", deletedCount))
	}

	return nil
}
