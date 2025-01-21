package service

import (
	"errors"
	"gateway-service/dto"
	"gateway-service/helpers"
	"log"
	"net/http"
	"time"

	pb "github.com/dharmasatrya/goodkarma/event-service/proto"
	"google.golang.org/grpc/status"
)

func parseDate(dateStr string) time.Time {
	// Define the layout based on the full date format
	const layout = "2006-01-02T15:04:05" // Adjusted to handle date, time, and timezone
	parsedDate, err := time.Parse(layout, dateStr)
	if err != nil {
		log.Printf("Failed to parse date: %v", err)
	}
	return parsedDate
}

type EventService interface {
	CreateEvent(token string, input dto.EventRequest) (int, *dto.Event, error)
	EditEvent(token string, id int, input dto.UpdateDescriptionRequest) (int, *dto.Event)
	GetAllEvents() (int, *[]dto.Event)
	GetEventById(id int) (int, *dto.Event)
	GetEventByUserLogin(token string) (int, *[]dto.Event)
	GetEventByCategory(category string) (int, *[]dto.Event)
}

type eventService struct {
	Client pb.EventServiceClient
}

func NewEventService(eventClient pb.EventServiceClient) *eventService {
	return &eventService{eventClient}
}

func (s *eventService) CreateEvent(token string, input dto.EventRequest) (int, *dto.Event, error) {
	ctx, cancel, err := helpers.NewServiceContext(token)
	if err != nil {
		log.Printf("Error creating context %v", err)
		return http.StatusInternalServerError, nil, err
	}
	defer cancel()

	res, err := s.Client.CreateEvent(ctx, &pb.EventRequest{
		Name:         input.Name,
		Description:  input.Description,
		DateStart:    input.DateStart,
		DateEnd:      input.DateEnd,
		DonationType: input.DonationType,
	})

	if err != nil {
		grpcErr, ok := status.FromError(err)
		if ok {
			return http.StatusInternalServerError, nil, errors.New(grpcErr.Message())
		}

		return http.StatusInternalServerError, nil, err
	}

	const layout = "2006-01-02" // Matches dates like "2024-01-01"
	dateStart, err := time.Parse(layout, res.DateStart)
	if err != nil {
		log.Printf("Failed to parse date: %v", err)
	}

	dateEnd, err := time.Parse(layout, res.DateEnd)
	if err != nil {
		log.Printf("Failed to parse date: %v", err)
	}

	response := dto.Event{
		ID:           int(res.Id),
		UserID:       res.UserId,
		Name:         res.Name,
		Description:  res.Description,
		DateStart:    dateStart,
		DateEnd:      dateEnd,
		DonationType: res.DonationType,
	}

	return http.StatusCreated, &response, nil
}

func (s *eventService) EditEvent(token string, id int, input dto.UpdateDescriptionRequest) (int, *dto.Event) {
	ctx, cancel, err := helpers.NewServiceContext(token)
	if err != nil {
		log.Printf("Error creating context %v", err)
		return http.StatusInternalServerError, nil
	}
	defer cancel()

	res, err := s.Client.UpdateDescription(ctx, &pb.UpdateDescriptionRequest{Id: uint32(id), Description: input.Description})
	if err != nil {
		log.Printf("Error updating desciprtion: %v", err)
		status := helpers.AssertGrpcStatus(err)
		return status, nil
	}

	dateStart := parseDate(res.DateStart)
	dateEnd := parseDate(res.DateEnd)

	response := dto.Event{
		ID:           int(res.Id),
		UserID:       res.UserId,
		Name:         res.Name,
		Description:  res.Description,
		DateStart:    dateStart,
		DateEnd:      dateEnd,
		DonationType: res.DonationType,
	}

	return http.StatusOK, &response
}

func (s *eventService) GetAllEvents() (int, *[]dto.Event) {
	ctx, cancel, err := helpers.NewServiceWithoutTokenContext()
	if err != nil {
		log.Printf("Error updating desciprtion: %v", err)
		return http.StatusInternalServerError, nil
	}
	defer cancel()

	res, err := s.Client.GetAllEvent(ctx, &pb.Empty{})
	if err != nil {
		log.Printf("Error get all events: %v", err)
		status := helpers.AssertGrpcStatus(err)
		return status, nil
	}

	var events []dto.Event
	for _, event := range res.Events {
		dateStart := parseDate(event.DateStart)
		dateEnd := parseDate(event.DateEnd)

		events = append(events, dto.Event{
			ID:           int(event.Id),
			UserID:       event.UserId,
			Name:         event.Name,
			Description:  event.Description,
			DateStart:    dateStart,
			DateEnd:      dateEnd,
			DonationType: event.DonationType,
		})
	}

	return http.StatusOK, &events
}

func (s *eventService) GetEventById(id int) (int, *dto.Event) {
	ctx, cancel, err := helpers.NewServiceWithoutTokenContext()
	if err != nil {
		return http.StatusInternalServerError, nil
	}
	defer cancel()

	res, err := s.Client.GetEventById(ctx, &pb.Id{Id: uint32(id)})
	if err != nil {
		log.Printf("Error get event by id: %v", err)
		status := helpers.AssertGrpcStatus(err)
		return status, nil
	}

	dateStart := parseDate(res.DateStart)
	dateEnd := parseDate(res.DateEnd)

	response := dto.Event{
		ID:           id,
		UserID:       res.UserId,
		Name:         res.Name,
		Description:  res.Description,
		DateStart:    dateStart,
		DateEnd:      dateEnd,
		DonationType: res.DonationType,
	}

	return http.StatusOK, &response
}

func (s *eventService) GetEventByUserLogin(token string) (int, *[]dto.Event) {
	ctx, cancel, err := helpers.NewServiceContext(token)
	if err != nil {
		log.Printf("Error creating context %v", err)
		return http.StatusInternalServerError, nil
	}
	defer cancel()

	res, err := s.Client.GetEventByUserId(ctx, &pb.Empty{})

	if err != nil {
		log.Printf("Error get event by user id: %v", err)
		status := helpers.AssertGrpcStatus(err)
		return status, nil
	}

	var events []dto.Event
	for _, event := range res.Events {
		dateStart := parseDate(event.DateStart)
		dateEnd := parseDate(event.DateEnd)

		events = append(events, dto.Event{
			ID:           int(event.Id),
			UserID:       event.UserId,
			Name:         event.Name,
			Description:  event.Description,
			DateStart:    dateStart,
			DateEnd:      dateEnd,
			DonationType: event.DonationType,
		})
	}

	return http.StatusOK, &events
}

func (s *eventService) GetEventByCategory(category string) (int, *[]dto.Event) {
	ctx, cancel, err := helpers.NewServiceWithoutTokenContext()
	if err != nil {
		log.Printf("Error creating context %v", err)
		return http.StatusInternalServerError, nil
	}
	defer cancel()

	res, err := s.Client.GetEventByCategory(ctx, &pb.Category{
		Category: category,
	})
	if err != nil {
		log.Printf("Error get event by category: %v", err)
		status := helpers.AssertGrpcStatus(err)
		return status, nil
	}

	var events []dto.Event
	for _, event := range res.Events {
		dateStart := parseDate(event.DateStart)
		dateEnd := parseDate(event.DateEnd)

		events = append(events, dto.Event{
			ID:           int(event.Id),
			UserID:       event.UserId,
			Name:         event.Name,
			Description:  event.Description,
			DateStart:    dateStart,
			DateEnd:      dateEnd,
			DonationType: event.DonationType,
		})
	}

	return http.StatusOK, &events
}
