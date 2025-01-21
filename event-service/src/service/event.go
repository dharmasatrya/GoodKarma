package service

import (
	"context"
	"log"
	"strconv"

	"github.com/dharmasatrya/goodkarma/event-service/entity"
	"github.com/dharmasatrya/goodkarma/event-service/helpers"
	pb "github.com/dharmasatrya/goodkarma/event-service/proto"
	"github.com/dharmasatrya/goodkarma/event-service/src/repository"
	"github.com/golang-jwt/jwt/v4"

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

	claims, ok := ctx.Value("claims").(jwt.MapClaims)
	if !ok {
		log.Printf("Type of ctx.Value(\"claims\"): %T", ctx.Value("claims"))   // Log the type
		log.Printf("Value of ctx.Value(\"claims\"): %+v", ctx.Value("claims")) // Log the value
		return nil, status.Errorf(codes.Internal, "failed to get user claims")
	}

	// Proceed with claims as jwt.MapClaims
	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, status.Errorf(codes.Internal, "user_id not found in claims")
	}

	if err := validateCreateEventRequest(req); err != nil {
		return nil, err
	}

	event := entity.Event{
		UserID:       userID,
		Name:         req.Name,
		Description:  req.Description,
		DateStart:    helpers.ParseDate(req.DateStart),
		DateEnd:      helpers.ParseDate(req.DateEnd),
		DonationType: req.DonationType,
	}

	res, err := s.eventRepository.CreateEvent(event)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error creating event")
	}

	id := strconv.Itoa(res.ID)

	return &pb.EventResponse{
		Id:           id,
		UserId:       res.UserID,
		Name:         req.Name,
		Description:  req.Description,
		DateStart:    req.DateStart,
		DateEnd:      req.DateEnd,
		DonationType: req.DonationType,
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
		Id:           idString,
		UserId:       res.UserID,
		Name:         res.Name,
		Description:  res.Description,
		DateStart:    res.DateStart.Format("2006-01-02T15:04:05"), // Format time.Time as string
		DateEnd:      res.DateEnd.Format("2006-01-02T15:04:05"),   // Format time.Time as string
		DonationType: res.DonationType,
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
			Id:           strconv.Itoa(event.ID),
			UserId:       event.UserID,
			Name:         event.Name,
			Description:  event.Description,
			DateStart:    event.DateStart.Format("2006-01-02T15:04:05"),
			DateEnd:      event.DateEnd.Format("2006-01-02T15:04:05"),
			DonationType: event.DonationType,
		})
	}

	// Return the response with the list of events
	return &pb.EventListResponse{
		Events: eventResponses,
	}, nil
}

func (s *EventService) GetEventById(ctx context.Context, req *pb.Id) (*pb.EventResponse, error) {
	id, err := strconv.Atoi(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error: id is not integer")
	}

	event, err := s.eventRepository.GetEventById(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Errorf(codes.NotFound, "event with ID %s not found", req.Id)
		}
		return nil, status.Errorf(codes.Internal, "error retrieving event: %v", err)
	}

	return &pb.EventResponse{
		Id:           strconv.Itoa(event.ID),
		UserId:       event.UserID,
		Name:         event.Name,
		Description:  event.Description,
		DateStart:    event.DateStart.Format("2006-01-02T15:04:05"),
		DateEnd:      event.DateEnd.Format("2006-01-02T15:04:05"),
		DonationType: event.DonationType,
	}, nil
}

func (s *EventService) GetEventByUserId(ctx context.Context, req *pb.Empty) (*pb.EventListResponse, error) {
	_, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	claims, ok := ctx.Value("claims").(jwt.MapClaims)
	if !ok {
		return nil, status.Errorf(codes.Internal, "failed to get user claims")
	}

	// Proceed with claims as jwt.MapClaims
	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, status.Errorf(codes.Internal, "user_id not found in claims")
	}

	events, err := s.eventRepository.GetEventsByUserId(userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error retrieving events: %v", err)
	}

	var eventResponses []*pb.EventResponse
	for _, event := range *events {
		eventResponses = append(eventResponses, &pb.EventResponse{
			Id:           strconv.Itoa(event.ID),
			UserId:       event.UserID,
			Name:         event.Name,
			Description:  event.Description,
			DateStart:    event.DateStart.Format("2006-01-02T15:04:05"),
			DateEnd:      event.DateEnd.Format("2006-01-02T15:04:05"),
			DonationType: event.DonationType,
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
			Id:           strconv.Itoa(event.ID),
			UserId:       event.UserID,
			Name:         event.Name,
			Description:  event.Description,
			DateStart:    event.DateStart.Format("2006-01-02T15:04:05"),
			DateEnd:      event.DateEnd.Format("2006-01-02T15:04:05"),
			DonationType: event.DonationType,
		})
	}

	return &pb.EventListResponse{
		Events: eventResponses,
	}, nil
}

func validateCreateEventRequest(req *pb.EventRequest) error {
	dateStart := helpers.ParseDate(req.DateStart)
	dateEnd := helpers.ParseDate(req.DateEnd)

	if req.UserId == "" {
		return status.Error(codes.InvalidArgument, "user id is required")
	}

	if req.Name == "" {
		return status.Error(codes.InvalidArgument, "name is required")
	}

	if len(req.Name) < 10 {
		return status.Error(codes.InvalidArgument, "event name must be at least 10 characters")
	}

	if req.Description == "" {
		return status.Error(codes.InvalidArgument, "description is required")
	}

	if len(req.Description) < 20 {
		return status.Error(codes.InvalidArgument, "event description must be at least 20 characters")
	}

	if req.DateStart == "" {
		return status.Error(codes.InvalidArgument, "date start is required")
	}

	if req.DateEnd == "" {
		return status.Error(codes.InvalidArgument, "date end is required")
	}

	if dateEnd.Before(dateStart) {
		return status.Error(codes.InvalidArgument, "date end must be after date start")
	}

	if req.DonationType == "" {
		return status.Error(codes.InvalidArgument, "donation type is required")
	}

	return nil
}
