package models

import "time"

type Event struct {
	ID       int64     `db:"id"`
	Title    string    `db:"title"`
	DateTime time.Time `db:"date_time"` // время проведения события
	UserID   string    `db:"user_id"`   // владелец события
	NotifyAt time.Time `db:"notify_at"` // время отправки уведомления
}
