// Package models provides functionality for models.
// Package models contains data structures.
package models

import "time"

type Event struct {
	ID       int64     `db:"id"`
	Title    string    `db:"title"`
	DateTime time.Time `db:"date_time"`
}
