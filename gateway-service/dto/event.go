package dto

import "time"

type Event struct {
	ID           int       `json:"id"`
	UserID       string    `json:"user_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	DateStart    time.Time `json:"date_start`
	DateEnd      time.Time `json:"date_end"`
	DonationType string    `json:"donation_type"`
}

type EventRequest struct {
	UserID       string `json:"user_id" validate:"required"`
	Name         string `json:"name" validate:"required"`
	Description  string `json:"description" validate:"required"`
	DateStart    string `json:"date_start" validate:"required"`
	DateEnd      string `json:"date_end" validate:"required"`
	DonationType string `json:"donation_type" validate:"required"`
}

type EventResponse struct {
	ID           int       `json:"id"`
	UserID       string    `json:"user_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	DateStart    time.Time `json:"date_start"`
	DateEnd      time.Time `json:"date_end"`
	DonationType string    `json:"donation_type"`
}

type UpdateDescriptionRequest struct {
	Description string `json:"description"`
}
