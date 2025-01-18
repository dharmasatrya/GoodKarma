package service

import (
	"context"
	"gateway-service/dto"
	"log"
	"net/http"

	pb "github.com/dharmasatrya/goodkarma/event-service/proto"
)

type EventService interface {
	CreateEvent(input dto.EventRequest) (int, *dto.Event)
	EditEvent(input dto.EventRequest) (int, *dto.EventResponse)
	GetAllEvents() (int, *[]dto.EventResponse)
	GetEventById(id int) (int, *dto.EventResponse)
	GetEventByUserId(user_id int) (int, *dto.EventResponse)
	GetEventByCategory(category string) (int, *[]dto.EventResponse)
}

type eventService struct {
	Client pb.EventServiceClient
}

func NewEventService(eventClient pb.EventServiceClient) *eventService {
	return &eventService{eventClient}
}

func (s *eventService) CreateEvent(input dto.EventRequest) (int, *dto.Event) {
	res, err := s.Client.CreateEvent(context.Background(), &pb.CreateEventRequest{})
	if err != nil {
		log.Fatalf("error while create request %v", err)
	}

	response := dto.Event{
		ID:          res.Id,
		UserID:      res.UserId,
		Name:        res.Name,
		Description: res.Description,
		DateStart:   res.DateStart,
		DateEnd:     res.DateEnd,
	}

	return http.StatusCreated, &response
}

func (s *eventService) EditEvent(input dto.EventRequest) (int, *dto.EventResponse) {
	res, err := s.Client.CreateEvent(context.Background(), &pb.UpdateDescriptionRequst{})
	if err != nil {
		log.Fatalf("error while create request %v", err)
	}

	response := dto.Event{
		ID:          res.Id,
		UserID:      res.UserId,
		Name:        res.Name,
		Description: res.Description,
		DateStart:   res.DateStart,
		DateEnd:     res.DateEnd,
	}

	return http.StatusCreated, &response
}

func (s *eventService) GetAllEvents() (int, *[]dto.EventResponse) {
	res, err := s.Client.GetAllEvent(context.Background(), &pb.Empty{})
	if err != nil {
		log.Printf("error while creating request: %v", err)
		return http.StatusInternalServerError, nil
	}

	var events []dto.EventResponse
	for _, event := range res.Events {
		events = append(events, dto.EventResponse{
			ID:          event.Id,
			UserID:      event.UserId,
			Name:        event.Name,
			Description: event.Description,
			DateStart:   event.DateStart,
			DateEnd:     event.DateEnd,
		})
	}

	return http.StatusOK, &events
}

func (s *eventService) GetEventById(id int) (int, *dto.EventResponse) {
	res, err := s.Client.GetEventById(context.Background(), &pb.Id{})
	if err != nil {
		log.Fatalf("error while create request %v", err)
	}

	response := dto.Event{
		ID:          res.Id,
		UserID:      res.UserId,
		Name:        res.Name,
		Description: res.Description,
		DateStart:   res.DateStart,
		DateEnd:     res.DateEnd,
	}

	return http.StatusCreated, &response
}

func (s *eventService) GetEventByUserId(user_id int) (int, *[]dto.EventResponse) {
	res, err := s.Client.GetEventByUserId(context.Background(), &pb.UserId{
		Id: int32(user_id),
	})
	if err != nil {
		log.Printf("error while creating request: %v", err)
		return http.StatusInternalServerError, nil
	}

	var events []dto.EventResponse
	for _, event := range res.Events {
		events = append(events, dto.EventResponse{
			ID:          event.Id,
			UserID:      event.UserId,
			Name:        event.Name,
			Description: event.Description,
			DateStart:   event.DateStart,
			DateEnd:     event.DateEnd,
		})
	}

	return http.StatusOK, &events
}

func (s *eventService) GetEventByCategory(category string) (int, *[]dto.EventResponse) {
	res, err := s.Client.GetEventByCategory(context.Background(), &pb.Category{
		Name: category,
	})
	if err != nil {
		log.Printf("error while creating request: %v", err)
		return http.StatusInternalServerError, nil
	}

	var events []dto.EventResponse
	for _, event := range res.Events {
		events = append(events, dto.EventResponse{
			ID:          event.Id,
			UserID:      event.UserId,
			Name:        event.Name,
			Description: event.Description,
			DateStart:   event.DateStart,
			DateEnd:     event.DateEnd,
		})
	}

	return http.StatusOK, &events
}
