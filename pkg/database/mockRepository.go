package database

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math"

	pbcommon "github.com/kic/health/pkg/proto/common"
	pbhealth "github.com/kic/health/pkg/proto/health"
)

type MockRepository struct {
	logCollection map[int]*pbhealth.MentalHealthLog

	idCounter int

	logger *zap.SugaredLogger
}

func NewMockRepository(logCollection map[int]*pbhealth.MentalHealthLog, logger *zap.SugaredLogger) *MockRepository {
	return &MockRepository{
		logCollection: logCollection,
		idCounter:     len(logCollection),
		logger:        logger,
	}
}

func (m *MockRepository) AddMentalHealthLog(ctx context.Context, healthLog *pbhealth.MentalHealthLog) (string, error) {
	if healthLog.UserID < 0 || healthLog.LogDate == nil {
		return "", status.Errorf(codes.InvalidArgument, "Invalid Argument for AddMentalHealthLog")
	}
	m.logCollection[m.idCounter] = healthLog
	var toReturn string
	toReturn = fmt.Sprint(m.idCounter)
	m.idCounter++

	return toReturn, nil
}

func (m *MockRepository) GetAllMentalHealthLogs(ctx context.Context, userID int64) ([]*pbhealth.MentalHealthLog, error) {
	toReturn := make([]*pbhealth.MentalHealthLog, 0)

	for _, val := range m.logCollection {
		if val.UserID == userID {
			toReturn = append(toReturn, val)
		}
	}

	if toReturn == nil {
		return nil, status.Errorf(codes.NotFound, "Health Log not found")
	}

	return toReturn, nil
}

func (m *MockRepository) GetAllMentalHealthLogsByDate(ctx context.Context, userID int64, date *pbcommon.Date) ([]*pbhealth.MentalHealthLog, error) {

	toReturn := make([]*pbhealth.MentalHealthLog, 0)

	for _, val := range m.logCollection {
		if val.UserID == userID && val.LogDate.Year == date.Year && val.LogDate.Month == date.Month && val.LogDate.Day == date.Day {
			toReturn = append(toReturn, val)
		}
	}

	if toReturn == nil {
		return nil, status.Errorf(codes.NotFound, "Health Log not found")
	}

	return toReturn, nil
}

func (m *MockRepository) DeleteMentalHealthLogs(ctx context.Context, userID int64, date *pbcommon.Date, all bool) (uint32, error) {
	var numDeleted uint32
	numDeleted = 0

	for key, val := range m.logCollection {
		if val.UserID == userID && val.LogDate.Year == date.Year && val.LogDate.Month == date.Month && val.LogDate.Day == date.Day {
			delete(m.logCollection, key)
			numDeleted++
		}
	}

	return numDeleted, nil
}

func (m *MockRepository) UpdateMentalHealthLogs(ctx context.Context, userID int64, healthLog *pbhealth.MentalHealthLog) error {
	for _, val := range m.logCollection {
		if val.UserID == userID && val.LogDate.Year == healthLog.LogDate.Year && val.LogDate.Month == healthLog.LogDate.Month && val.LogDate.Day == healthLog.LogDate.Day {
			val.Score = healthLog.Score
			val.JournalName = healthLog.JournalName
		}
	}

	return nil
}

func (m *MockRepository) GetOverallScore(ctx context.Context, userID int64) (int32, error) {
	logs, err := m.GetAllMentalHealthLogs(ctx, userID)
	m.logger.Infof("Logs fetched for user (ID = %v\n): %v\n", userID, logs)

	if err != nil {
		m.logger.Errorf("cannot get mental health logs for user: %v \n", err)
	}
	var totalScore float64
	totalScore = 0

	numLogs := 0

	for _, log := range logs {
		totalScore += float64(log.Score)
		numLogs++
	}

	var overallScore int32
	if numLogs == 0 {
		overallScore = 0
	} else {
		overallScore = int32(math.Round(totalScore / float64(numLogs)))
	}

	m.logger.Infof("Total sum of log scores for user (ID = %v\n): %v\n", userID, totalScore)
	m.logger.Infof("Number of total logs for user (ID = %v): %v\n:", userID, numLogs)
	m.logger.Infof("Average score for user (ID = %v): %v\n:", userID, overallScore)

	return overallScore, err
}