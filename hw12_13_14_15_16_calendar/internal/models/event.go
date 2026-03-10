package models

import "time"

type Event struct {
    ID        string
    Title     string
    UserID    string
    StartTime time.Time
}
