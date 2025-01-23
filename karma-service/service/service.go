package service

import (
	"context"

	"github.com/dharmasatrya/goodkarma/karma-service/entity"
	pb "github.com/dharmasatrya/goodkarma/karma-service/proto"
	"github.com/dharmasatrya/goodkarma/karma-service/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type KarmaService struct {
	KarmaRepository repository.KarmaRepository
	pb.UnimplementedKarmaServiceServer
}

func NewKarmaService(karmaRepository repository.KarmaRepository) *KarmaService {
	return &KarmaService{
		KarmaRepository: karmaRepository,
	}
}

func (s *KarmaService) CreateKarma(ctx context.Context, req *pb.CreateKarmaRequest) (*pb.CreateKarmaResponse, error) {
	karma := entity.CreateKarmaRequest{
		UserID: req.UserId,
		Amount: req.Amount,
	}

	reult, err := s.KarmaRepository.CreateKarma(ctx, karma)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "error creating karma")
	}

	return &pb.CreateKarmaResponse{
		Id:     reult.ID.Hex(),
		UserId: reult.UserID.Hex(),
		Amount: reult.Amount,
	}, nil
}

func (s *KarmaService) GetReferralCount(ctx context.Context, req *pb.GetReferralCountRequest) (*pb.GetReferralCountResponse, error) {
	count, err := s.KarmaRepository.GetReferralCount(ctx, req.ReferralCode)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error getting referral count")
	}

	return &pb.GetReferralCountResponse{
		Count: count,
	}, nil
}

func (s *KarmaService) CreateReferralLog(ctx context.Context, req *pb.CreateReferralLogRequest) (*pb.Empty, error) {
	referral := entity.ReferralLog{
		UserID:       req.UserId,
		ReferralCode: req.ReferralCode,
	}

	err := s.KarmaRepository.CreateReferralLog(ctx, referral)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error creating referral log")
	}

	return &pb.Empty{}, nil
}
