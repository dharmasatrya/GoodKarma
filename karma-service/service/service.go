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

func (s *KarmaService) UpdateKarmaAmount(ctx context.Context, req *pb.UpdateKarmaAmountRequest) (*pb.Empty, error) {
	karma := entity.UpdateKarmaRequest{
		UserID: req.UserId,
		Amount: req.Amount,
	}

	err := s.KarmaRepository.UpdateKarmaAmount(ctx, karma)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error updating karma")
	}

	return &pb.Empty{}, nil
}

func (s *KarmaService) GetUserByReferralCode(ctx context.Context, req *pb.GetUserByReferralCodeRequest) (*pb.GetUserByReferralCodeResponse, error) {
	userID, err := s.KarmaRepository.GetUserByReferralCode(ctx, req.ReferralCode)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	return &pb.GetUserByReferralCodeResponse{
		UserId: userID,
	}, nil
}

func (s *KarmaService) ExchangeReward(ctx context.Context, req *pb.ExchangeRewardRequest) (*pb.Empty, error) {
	exchange := entity.ExchangeRewardRequest{
		UserID:        req.UserId,
		KarmaRewardID: req.KarmaRewardId,
	}

	err := s.KarmaRepository.ExchangeReward(ctx, exchange)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error exchanging reward")
	}

	return &pb.Empty{}, nil
}

func (s *KarmaService) GetKarmaReward(ctx context.Context, req *pb.Empty) (*pb.GetKarmaRewardResponse, error) {
	rewards, err := s.KarmaRepository.GetKarmaReward(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error getting karma rewards")
	}

	var res []*pb.KarmaReward
	for _, reward := range rewards {
		res = append(res, &pb.KarmaReward{
			Id:          reward.ID.Hex(),
			Name:        reward.Name,
			Amount:      reward.Amount,
			Description: reward.Description,
		})
	}

	return &pb.GetKarmaRewardResponse{
		Rewards: res,
	}, nil
}
