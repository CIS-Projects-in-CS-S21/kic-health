package server_test

import (
	"context"
	"github.com/kic/health/internal/server"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"

	"github.com/kic/health/internal/setup"
	"github.com/kic/health/pkg/database"
	"github.com/kic/health/pkg/logging"
	pbcommon "github.com/kic/health/pkg/proto/common"
	pbhealth "github.com/kic/health/pkg/proto/health"
)

var log *zap.SugaredLogger
var client pbhealth.HealthTrackingClient

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

	repo, mongoClient := setup.DBRepositorySetup(log, "test-health")

	prepDBForTests(repo)

	ListenAddress := "localhost:50051"

	listener, err := net.Listen("tcp", ListenAddress)
	if err != nil {
		log.Fatalf("Unable to listen on %v: %v", ListenAddress, err)
	}

	grpcServer := grpc.NewServer()

	mediaService := server.NewHealthService(repo, log)
	pbhealth.RegisterHealthTrackingServer(grpcServer, mediaService)

	reflection.Register(grpcServer)

	go func() {
		defer listener.Close()
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	log.Infof("Server started on %v", ListenAddress)

	defer grpcServer.Stop()
	defer mongoClient.Disconnect(context.Background())

	conn, err := grpc.Dial(ListenAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client = pbhealth.NewHealthTrackingClient(conn)

	exitVal := m.Run()

	os.Exit(exitVal)
}








