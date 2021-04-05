package database

import (
	"context"
	pbcommon "github.com/kic/health/pkg/proto/common"

	pbhealth "github.com/kic/health/pkg/proto/health"
)

// Repository - interface for a data provider that interfaces between the database backend and the grpc server
// enables the repository pattern so that we can swap out the database backend easily
type Repository interface {
	GetOverallScore(ctx context.Context, userID int64) (int, error)
	GetAllMentalHealthLogs(ctx context.Context, userID int64) ([]*pbhealth.MentalHealthLog, error)
	GetAllMentalHealthLogsByDate(ctx context.Context, userID int64, date *pbcommon.Date) ([]*pbhealth.MentalHealthLog, error)
	AddMentalHealthLog(ctx context.Context, healthLog *pbhealth.MentalHealthLog) (string, error)
}
