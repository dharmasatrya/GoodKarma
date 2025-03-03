package repository

import (
	"errors"
	"fmt"

	"github.com/dharmasatrya/goodkarma/event-service/entity"

	"gorm.io/gorm" // ORM (Object Relational Mapping) Gorm untuk interaksi dengan database.
)

type EventRepository interface {
	CreateEvent(event entity.Event) (*entity.Event, error)
	EditDescription(id int, description string) (*entity.Event, error)
	GetAllEvents() (*[]entity.Event, error)
	GetEventById(id int) (*entity.Event, error)
	GetEventsByUserId(id string) (*[]entity.Event, error)
	GetEventsByCategory(category string) (*[]entity.Event, error)
}

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) *eventRepository {
	return &eventRepository{db}
}

func (r *eventRepository) CreateEvent(event entity.Event) (*entity.Event, error) {
	// Menyimpan data order ke database menggunakan GORM.
	if err := r.db.Create(&event).Error; err != nil {
		return nil, err
	}

	return &event, nil
}

func (r *eventRepository) EditDescription(id int, description string) (*entity.Event, error) {
	result := r.db.Model(&entity.Event{}).
		Where("id = ?", id).
		Update("description", gorm.Expr("?", description))

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	var updatedEvent entity.Event
	if err := r.db.Where("id = ?", id).First(&updatedEvent).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("event with id %d not found: %w", id, err)
		}
		return nil, err
	}

	return &updatedEvent, nil
}

func (r *eventRepository) GetAllEvents() (*[]entity.Event, error) {
	var events []entity.Event
	if err := r.db.Find(&events).Error; err != nil {
		return nil, err
	}

	return &events, nil
}

func (r *eventRepository) GetEventById(id int) (*entity.Event, error) {
	var event entity.Event
	if err := r.db.Where("id = ?", id).First(&event).Error; err != nil {
		return nil, err
	}

	return &event, nil
}

func (r *eventRepository) GetEventsByUserId(id string) (*[]entity.Event, error) {
	var event []entity.Event
	if err := r.db.Where("user_id = ?", id).Find(&event).Error; err != nil {
		return nil, err
	}

	return &event, nil
}

func (r *eventRepository) GetEventsByCategory(category string) (*[]entity.Event, error) {
	var event []entity.Event
	query := "%" + category + "%"
	if err := r.db.Where("donation_type LIKE ?", query).Find(&event).Error; err != nil {
		return nil, err
	}

	return &event, nil
}
