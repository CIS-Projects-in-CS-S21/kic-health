package server

import (
	"context"
	pbhealth "github.com/kic/health/pkg/proto/health"
	"github.com/kic/health/pkg/database"
	"go.uber.org/zap"
)

type HealthService struct {
	pbhealth.UnimplementedHealthTrackingServer
	db         database.Repository

	logger  *zap.SugaredLogger
}

func NewHealthService(db database.Repository, logger *zap.SugaredLogger) *HealthService {
	return &HealthService{
		logger:  logger,
		db: db,
	}
}

func (h *HealthService) GetHealthDataForUser(
	ctx context.Context,
	req *pbhealth.GetHealthDataForUserRequest,
) (*pbhealth.GetHealthDataForUserResponse, error) {


	return nil, nil
}

func (h *HealthService) GetHealthDataByDate(
	ctx context.Context,
	req *pbhealth.GetHealthDataByDateRequest,
) (*pbhealth.GetHealthDataByDateResponse, error) {


	return nil, nil
}

func (h *HealthService) GetMentalHealthScoreForUser(
	ctx context.Context,
	req *pbhealth.GetMentalHealthScoreForUserRequest,
) (*pbhealth.GetMentalHealthScoreForUserResponse, error) {
	score, err := h.db.GetOverallScore(ctx, req.UserID)

	res := &pbhealth.GetMentalHealthScoreForUserResponse{Score: int32(score)}

	return res, err
}
