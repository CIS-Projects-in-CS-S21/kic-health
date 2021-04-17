package server

import (
	"context"
	"github.com/kic/health/pkg/database"
	pbhealth "github.com/kic/health/pkg/proto/health"
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
		UnimplementedHealthTrackingServer: pbhealth.UnimplementedHealthTrackingServer{},
		logger:  logger,
		db: db,
	}
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

func (h *HealthService) GetHealthDataForUser(
	ctx context.Context,
	req *pbhealth.GetHealthDataForUserRequest,
) (*pbhealth.GetHealthDataForUserResponse, error) {

	logs, err := h.db.GetAllMentalHealthLogs(ctx, req.UserID)

	if err != nil {
		h.logger.Infof("%v", err)
		return &pbhealth.GetHealthDataForUserResponse{
			HealthData: nil,
		}, status.Errorf(codes.Internal, "Error getting health data by date")
	}

	h.logger.Infof("Successfully got mental health log by date: %v\n", logs)

	successRes := &pbhealth.GetHealthDataForUserResponse{HealthData: logs}


	return successRes, err
}

func (h *HealthService) GetHealthDataByDate(
	ctx context.Context,
	req *pbhealth.GetHealthDataByDateRequest,
) (*pbhealth.GetHealthDataByDateResponse, error) {
	logs, err := h.db.GetAllMentalHealthLogsByDate(ctx, req.UserID, req.LogDate)

	if err != nil {
		h.logger.Infof("%v", err)
		return &pbhealth.GetHealthDataByDateResponse{
			HealthData: nil,
		}, status.Errorf(codes.Internal, "Error getting health data by date")
	}

	h.logger.Infof("Successfully got mental health log by date: %v\n", logs)

	successRes := &pbhealth.GetHealthDataByDateResponse{HealthData: logs}


	return successRes, err
}

func (h *HealthService) GetMentalHealthScoreForUser(
	ctx context.Context,
	req *pbhealth.GetMentalHealthScoreForUserRequest,
) (*pbhealth.GetMentalHealthScoreForUserResponse, error) {
	score, err := h.db.GetOverallScore(ctx, req.UserID)

	if err != nil {
		h.logger.Errorf("cannot get mental overall score for user: %v \n", err)
	}

	res := &pbhealth.GetMentalHealthScoreForUserResponse{Score: score}

	return res, err
}

func (h *HealthService) DeleteHealthDataForUser(
	ctx context.Context,
	req *pbhealth.DeleteHealthDataForUserRequest,
) (*pbhealth.DeleteHealthDataForUserResponse, error) {
	var err error
	var numDeleted uint32

	switch x := req.Data.(type) {
	case *pbhealth.DeleteHealthDataForUserRequest_All:
		numDeleted, err = h.db.DeleteMentalHealthLogs(ctx, req.UserID, nil, x.All)
		break
	case *pbhealth.DeleteHealthDataForUserRequest_DateToRemove:
		numDeleted, err = h.db.DeleteMentalHealthLogs(ctx, req.UserID, x.DateToRemove, false)

	}

	res := &pbhealth.DeleteHealthDataForUserResponse{EntriesDeleted: numDeleted}

	return res, err
}

func (h *HealthService) UpdateHealthDataForDate(
	ctx context.Context,
	req *pbhealth.UpdateHealthDataForDateRequest,
) (*pbhealth.UpdateHealthDataForDateResponse, error) {

	err := h.db.UpdateMentalHealthLogs(ctx, req.UserID, req.DesiredLogInfo)
	if err != nil {
		h.logger.Errorf("%v", err)
		return &pbhealth.UpdateHealthDataForDateResponse{
			Success: false,
		}, status.Errorf(codes.InvalidArgument, "Error updating mental health logs")
	}

	h.logger.Infof("Successfully updated mental health log.")

	successRes := &pbhealth.UpdateHealthDataForDateResponse{Success: true}

	return successRes, err

	return nil, nil
}


