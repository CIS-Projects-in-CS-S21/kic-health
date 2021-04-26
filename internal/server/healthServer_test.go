package server_test

import (
	"context"
	"github.com/kic/health/internal/server"
	"os"
	"testing"
	"time"

	"github.com/kic/health/pkg/database"
	"github.com/kic/health/pkg/logging"
	pbcommon "github.com/kic/health/pkg/proto/common"
	pbhealth "github.com/kic/health/pkg/proto/health"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.SugaredLogger
var healthService *server.HealthService

const testDataPath = "../../test_data"

func prepDBForTests(db database.Repository) {
	healthLogsToAdd :=[]*pbhealth.MentalHealthLog{
		{
			LogDate:     &pbcommon.Date{
				Year:  2021,
				Month: 4,
				Day:   5,
			},
			Score:       5,
			JournalName: "I am happy!",
			UserID:      1,
		},

		{
			LogDate:     &pbcommon.Date{
				Year:  2021,
				Month: 4,
				Day:   26,
			},
			Score:       2,
			JournalName: "I am ok",
			UserID:      1,
		},

		{
			LogDate:     &pbcommon.Date{
				Year:  2021,
				Month: 4,
				Day:   23,
			},
			Score:       -5,
			JournalName: "I am sad :(",
			UserID:      1,
		},
	}

	for _, file := range healthLogsToAdd {
		id, err := db.AddMentalHealthLog(context.Background(), file)
		log.Debugf("inserted id: %v", id)
		if err != nil {
			log.Debugf("insertion error: %v", err)
		}
	}
}

func TestMain(m *testing.M) {
	time.Sleep(1 * time.Second)
	log = logging.CreateLogger(zapcore.DebugLevel)

	mp := make(map[int]*pbhealth.MentalHealthLog)
	repo := database.NewMockRepository(mp, log)

	prepDBForTests(repo)

	healthService = server.NewHealthService(repo, log)

	exitVal := m.Run()

	os.Exit(exitVal)
}

func Test_ShouldUploadLog(t *testing.T) {
	resp, err := healthService.AddHealthDataForUser(context.Background(), &pbhealth.AddHealthDataForUserRequest{
		UserID:   1,
		NewEntry: &pbhealth.MentalHealthLog{
			LogDate:     &pbcommon.Date{
				Year:  2021,
				Month: 26,
				Day:   4,
			},
			Score:       0,
			JournalName: "I am neutral",
			UserID:      1,
		},
	})
	if err != nil || resp.Success == false {
		t.Errorf("Add Health Data should not fail")
	}
}

func Test_ShouldFailUploadLog(t *testing.T) {
	_, err := healthService.AddHealthDataForUser(context.Background(), &pbhealth.AddHealthDataForUserRequest{
		UserID:   -1,
		NewEntry: &pbhealth.MentalHealthLog{
			LogDate:     nil,
			Score:       0,
			JournalName: "",
			UserID:      -1,
		},
	})
	if err == nil {
		t.Errorf("Add Health Data should fail")
	}
}

func Test_ShouldDeleteLog(t *testing.T) {
	_, err := healthService.DeleteHealthDataForUser(context.Background(), &pbhealth.DeleteHealthDataForUserRequest{
		UserID: 1,
		Data:   &pbhealth.DeleteHealthDataForUserRequest_All{true},
	})

	if err != nil {
		t.Errorf("Delete Health Data should not fail")
	}
}

func Test_ShouldFailDeleteLog(t *testing.T) {
	_, err := healthService.DeleteHealthDataForUser(context.Background(), &pbhealth.DeleteHealthDataForUserRequest{
		UserID: -1,
		Data:   &pbhealth.DeleteHealthDataForUserRequest_All{true},
	})

	if err == nil {
		t.Errorf("Delete Health Data should fail")
	}
}








