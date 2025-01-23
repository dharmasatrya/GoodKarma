package service

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/dharmasatrya/goodkarma/event-service/entity"
	"github.com/dharmasatrya/goodkarma/event-service/helpers"
	pb "github.com/dharmasatrya/goodkarma/event-service/proto"
	"github.com/dharmasatrya/goodkarma/event-service/src/repository"
	userpb "github.com/dharmasatrya/goodkarma/user-service/proto"
	"github.com/golang-jwt/jwt/v4"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type EventService struct {
	pb.UnimplementedEventServiceServer
	eventRepository repository.EventRepository
	userClient      userpb.UserServiceClient
}

func NewEventService(eventRepository repository.EventRepository) *EventService {
	userServiceURI := os.Getenv("USER_SERVICE_URI")

	grpcConn, err := grpc.NewClient(userServiceURI, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Error koneksi user service: %v", err)
	}

	userClient := userpb.NewUserServiceClient(grpcConn)

	return &EventService{eventRepository: eventRepository, userClient: userClient}
}

func (s *EventService) CreateEvent(ctx context.Context, req *pb.EventRequest) (*pb.EventResponse, error) {
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

	if err := s.ValidateCreateEvent(ctx, userID); err != nil {
		return nil, err
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
		log.Println(err)
		return nil, err
	}

	return &pb.EventResponse{
		Id:           uint32(res.ID),
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

	claims, ok := ctx.Value("claims").(jwt.MapClaims)
	if !ok {
		return nil, status.Errorf(codes.Internal, "failed to get user claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, status.Errorf(codes.Internal, "user_id not found in claims")
	}

	res, err := s.eventRepository.GetEventById(int(req.GetId()))
	if err == gorm.ErrRecordNotFound {
		log.Println("event with ID %s not found", req.Id)
		return nil, status.Errorf(codes.NotFound, "event with ID %s not found", req.Id)
	}

	if userID != res.UserID {
		log.Println("You are not the owner of this event!")
		return nil, status.Errorf(codes.Unauthenticated, "You are not the owner of this event!")
	}

	res, err = s.eventRepository.EditDescription(int(req.GetId()), req.Description)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error retrieving event: %v", err)
	}

	return &pb.UpdateDescriptionResponse{
		Id:           uint32(res.ID),
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

	var eventResponses []*pb.EventResponse
	for _, event := range *events {
		eventResponses = append(eventResponses, &pb.EventResponse{
			Id:           uint32(event.ID),
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

func (s *EventService) GetEventById(ctx context.Context, req *pb.Id) (*pb.EventResponse, error) {
	event, err := s.eventRepository.GetEventById(int(req.Id))

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Errorf(codes.NotFound, "event with ID %s not found", req.Id)
		}
		return nil, status.Errorf(codes.Internal, "error retrieving event: %v", err)
	}

	return &pb.EventResponse{
		Id:           uint32(event.ID),
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
			Id:           uint32(event.ID),
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
			Id:           uint32(event.ID),
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

func (s *EventService) ValidateCreateEvent(ctx context.Context, userID string) error {
	user, err := s.userClient.GetUserById(ctx, &userpb.GetUserByIdRequest{Id: userID})

	if err != nil {
		return err
	}

	if user.Role != "coordinator" {
		return fmt.Errorf("you do not have the necessary permissions to perform this action")
	}

	return nil
}
