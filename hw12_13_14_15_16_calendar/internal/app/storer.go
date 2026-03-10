package app

import (
	"context"
	"fmt"

	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/logger"
	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/models"
	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/mq"
	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/storage/notifications"
)

type Storer struct {
	storage  notifications.NotificationStorage
	consumer mq.Consumer
	logger   *logger.Logger
}

func NewStorer(storage notifications.NotificationStorage, consumer mq.Consumer, logger *logger.Logger) *Storer {
	return &Storer{
		storage:  storage,
		consumer: consumer,
		logger:   logger,
	}
}

func (s *Storer) Run(ctx context.Context) error {
	notificationsCh, err := s.consumer.Consume(ctx)
	if err != nil {
		return err
	}
	s.logger.Info("Started consuming notifications")
	for {
		select {
		case <-ctx.Done():
			return nil
		case notification, ok := <-notificationsCh:
			if !ok {
				s.logger.Info("Notifications channel closed")
				return nil
			}
			if err := s.processNotification(ctx, notification); err != nil {
				s.logger.Error(fmt.Sprintf("Failed to process notification: %v", err))
			}
		}
	}
}

func (s *Storer) processNotification(ctx context.Context, notification *models.Notification) error {
	if err := s.storage.SaveNotification(ctx, notification); err != nil {
		return err
	}
	s.logger.Info(fmt.Sprintf("Stored notification for event: %s", notification.EventTitle))
	return nil
}
