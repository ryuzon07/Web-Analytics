package types

import "time"

// Event defines the structure of an analytics event
type Event struct {
	SiteID    string    `json:"site_id" binding:"required"`
	EventType string    `json:"event_type" binding:"required"`
	Path      string    `json:"path"`
	UserID    string    `json:"user_id"`
	Timestamp time.Time `json:"timestamp" binding:"required"`
}