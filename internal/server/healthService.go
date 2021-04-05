package server

import (
	"context"
	pbhealth "github.com/kic/health/pkg/proto/health"
	"github.com/kic/health/pkg/database"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (h *HealthService) AddHealthDataForUser(
	ctx context.Context,
	req *pbhealth.AddHealthDataForUserRequest,
) (*pbhealth.AddHealthDataForUserResponse, error) {

	id, err := h.db.AddMentalHealthLog(ctx, req.NewEntry)
	if err != nil {
		h.logger.Infof("%v", err)
		return &pbhealth.AddHealthDataForUserResponse{
			Success: false,
		}, status.Errorf(codes.Internal, "Error adding mental health log to database")
	}

	h.logger.Infof("Successfully added new mental health log. ID of user: %v\n", id)

	successRes := &pbhealth.AddHealthDataForUserResponse{Success: true}

	return successRes, err
}