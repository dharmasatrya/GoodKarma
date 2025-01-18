package service

import (
	"context"
	"goodkarma-event-service/entity"
	"goodkarma-event-service/helpers"
	pb "goodkarma-event-service/proto"
	"goodkarma-event-service/src/repository"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type EventService struct {
	pb.UnimplementedEventServiceServer
	eventRepository repository.EventRepository
}

func NewEventService(eventRepository repository.EventRepository) *EventService {
	return &EventService{eventRepository: eventRepository}
}

func (s *EventService) CreateEvent(ctx context.Context, req *pb.EventRequest) (*pb.EventResponse, error) {
	_, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	event := entity.Event{
		UserID:      req.UserId,
		Name:        req.Name,
		Description: req.Description,
		DateStart:   helpers.ParseDate(req.DateStart),
		DateEnd:     helpers.ParseDate(req.DateEnd),
	}

	res, err := s.eventRepository.CreateEvent(event)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error creating event")
	}

	id := strconv.Itoa(res.ID)

	return &pb.EventResponse{
		Id:          id,
		UserId:      req.UserId,
		Name:        req.Name,
		Description: req.Description,
		DateStart:   req.DateStart,
		DateEnd:     req.DateEnd,
	}, nil
}

func (s *EventService) UpdateDescription(ctx context.Context, req *pb.UpdateDescriptionRequest) (*pb.UpdateDescriptionResponse, error) {
	_, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	id, err := strconv.Atoi(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error: id is not integer")
	}

	res, err := s.eventRepository.EditDescription(id, req.Description)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error creating wallet")
	}

	idString := strconv.Itoa(res.ID)

	return &pb.UpdateDescriptionResponse{
		Id:          idString,
		UserId:      res.UserID,
		Name:        res.Name,
		Description: res.Description,
		DateStart:   res.DateStart.Format("2006-01-02T15:04:05"), // Format time.Time as string
		DateEnd:     res.DateEnd.Format("2006-01-02T15:04:05"),   // Format time.Time as string
	}, nil
}

func (s *EventService) GetAllEvent(ctx context.Context, req *pb.Empty) (*pb.EventListResponse, error) {
	events, err := s.eventRepository.GetAllEvents()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error creating wallet")
	}

	// Transform events into the response format
	var eventResponses []*pb.EventResponse
	for _, event := range *events {
		eventResponses = append(eventResponses, &pb.EventResponse{
			Id:          strconv.Itoa(event.ID),
			UserId:      event.UserID,
			Name:        event.Name,
			Description: event.Description,
			DateStart:   event.DateStart.Format("2006-01-02T15:04:05"),
			DateEnd:     event.DateEnd.Format("2006-01-02T15:04:05"),
		})
	}

	// Return the response with the list of events
	return &pb.EventListResponse{
		Events: eventResponses,
	}, nil
}

func (s *EventService) GetEventById(ctx context.Context, req *pb.Id) (*pb.EventResponse, error) {
	id, err := strconv.Atoi(req.Id)

	event, err := s.eventRepository.GetEventById(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Errorf(codes.NotFound, "event with ID %s not found", req.Id)
		}
		return nil, status.Errorf(codes.Internal, "error retrieving event: %v", err)
	}

	return &pb.EventResponse{
		Id:          strconv.Itoa(event.ID),
		UserId:      event.UserID,
		Name:        event.Name,
		Description: event.Description,
		DateStart:   event.DateStart.Format("2006-01-02T15:04:05"),
		DateEnd:     event.DateEnd.Format("2006-01-02T15:04:05"),
	}, nil
}

func (s *EventService) GetEventByUserId(ctx context.Context, req *pb.UserId) (*pb.EventListResponse, error) {
	userId, err := strconv.Atoi(req.UserId)

	events, err := s.eventRepository.GetEventsByUserId(userId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error retrieving events: %v", err)
	}

	var eventResponses []*pb.EventResponse
	for _, event := range *events {
		eventResponses = append(eventResponses, &pb.EventResponse{
			Id:          strconv.Itoa(event.ID),
			UserId:      event.UserID,
			Name:        event.Name,
			Description: event.Description,
			DateStart:   event.DateStart.Format("2006-01-02T15:04:05"),
			DateEnd:     event.DateEnd.Format("2006-01-02T15:04:05"),
		})
	}

	return &pb.EventListResponse{
		Events: eventResponses,
	}, nil
}

func (s *EventService) GetEventByCategory(ctx context.Context, req *pb.Category) (*pb.EventListResponse, error) {
	events, err := s.eventRepository.GetEventsByCategory(req.Category)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error retrieving events: %v", err)
	}

	var eventResponses []*pb.EventResponse
	for _, event := range *events {
		eventResponses = append(eventResponses, &pb.EventResponse{
			Id:          strconv.Itoa(event.ID),
			UserId:      event.UserID,
			Name:        event.Name,
			Description: event.Description,
			DateStart:   event.DateStart.Format("2006-01-02T15:04:05"),
			DateEnd:     event.DateEnd.Format("2006-01-02T15:04:05"),
		})
	}

	return &pb.EventListResponse{
		Events: eventResponses,
	}, nil
}
