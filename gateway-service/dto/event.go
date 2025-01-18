package dto

import "time"

type Event struct {
	ID          int       `json:"id"`
	UserID      string    `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	DateStart   time.Time `json:"date_start`
	DateEnd     time.Time `json:"date_end"`
}

type EventRequest struct {
	UserID      string    `json:"user_id" validate:"required"`
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description" validate:"required"`
	DateStart   time.Time `json:"date_start" validate:"required"`
	DateEnd     time.Time `json:"date_end" validate:"required"`
}

type EventResponse struct {
	ID          int       `json:"id"`
	UserID      string    `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	DateStart   time.Time `json:"date_start"`
	DateEnd     time.Time `json:"date_end"`
}
