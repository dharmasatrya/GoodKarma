package service

import (
	"context"

	entity "github.com/dharmasatrya/goodkarma/karma-service/entity"
	pb "github.com/dharmasatrya/goodkarma/karma-service/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type KarmaService interface {
	GetKarmaReward() ([]entity.KarmaReward, error)
}

type karmaService struct {
	karmaClient pb.KarmaServiceClient
}

func NewKarmaService(karmaClient pb.KarmaServiceClient) *karmaService {
	return &karmaService{karmaClient}
}

func (s *karmaService) GetKarmaReward() ([]entity.KarmaReward, error) {
	res, err := s.karmaClient.GetKarmaReward(context.Background(), &pb.Empty{})

	var rewards []entity.KarmaReward

	if err != nil {
		return nil, err
	}

	for _, reward := range res.Rewards {
		rewardId, err := primitive.ObjectIDFromHex(reward.Id)
		if err != nil {
			return nil, err
		}

		rewards = append(rewards, entity.KarmaReward{
			ID:          rewardId,
			Name:        reward.Name,
			Description: reward.Description,
			Amount:      reward.Amount,
		})
	}

	return rewards, nil
}
