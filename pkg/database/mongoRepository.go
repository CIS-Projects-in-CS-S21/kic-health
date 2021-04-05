package database

import (
	"context"
	"errors"
	pbcommon "github.com/kic/health/pkg/proto/common"
	pbhealth "github.com/kic/health/pkg/proto/health"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"strings"
)

const (
	fileCollectionName = "health"
)

type MongoRepository struct {
	client         *mongo.Client
	fileCollection *mongo.Collection

	logger *zap.SugaredLogger
}

func NewMongoRepository(client *mongo.Client, logger *zap.SugaredLogger) *MongoRepository {
	return &MongoRepository{
		client: client,
		logger: logger,
	}
}

func (m *MongoRepository) SetCollections(databaseName string) {
	m.fileCollection = m.client.Database(databaseName).Collection(fileCollectionName)
}

func (m *MongoRepository) AddMentalHealthLog(ctx context.Context, file *pbhealth.MentalHealthLog) error {

}

func (m *MongoRepository) GetAllMentalHealthLog(ctx context.Context, userID int64) ([]*pbhealth.MentalHealthLog, error) {
	toReturn := make([]*pbhealth.MentalHealthLog, 0)

	filter := bson.M{"userID": userID}

	cur, err := m.fileCollection.Find(ctx, filter)
	if err != nil {
		m.logger.Errorf("Error finding mental health logs: %v", err)
	}

	for cur.Next(context.Background()) {
		healthLog := &pbhealth.MentalHealthLog{}
		err = cur.Decode(healthLog)
		if err != nil {
			m.logger.Errorf("Error decoding file: %v", err)
			return toReturn, err
		}

	}

	return nil, nil
}

func (m *MongoRepository) GetAllMentalHealthLogsByDate(ctx context.Context, userID int64) (*pbhealth.MentalHealthLog, error) {
	return nil, nil
}

func (m *MongoRepository) GetOverallScore(ctx context.Context, userID int64) (int, error) {


	return 0, nil
}



