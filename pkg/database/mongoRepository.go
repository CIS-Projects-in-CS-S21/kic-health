package database

import (
	"context"
	_ "errors"
	pbcommon "github.com/kic/health/pkg/proto/common"
	pbhealth "github.com/kic/health/pkg/proto/health"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	_ "strings"
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

func (m *MongoRepository) AddMentalHealthLog(ctx context.Context, healthLog *pbhealth.MentalHealthLog) (string, error) {
	res, err := m.fileCollection.InsertOne(context.TODO(), healthLog)
	if err != nil {
		m.logger.Infof("%v", err)
		return "", err
	}
	var toReturn string
	toReturn = res.InsertedID.(primitive.ObjectID).Hex()

	return toReturn, err

}

func (m *MongoRepository) GetAllMentalHealthLogs(ctx context.Context, userID int64) ([]*pbhealth.MentalHealthLog, error) {
	toReturn := make([]*pbhealth.MentalHealthLog, 0)

	filter := bson.M{"userid": userID}

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
		toReturn = append(toReturn, healthLog)
	}

	return toReturn, err
}


func (m *MongoRepository) GetOverallScore(ctx context.Context, userID int64) (int, error) {


	return 0, nil
}

func (m *MongoRepository) GetAllMentalHealthLogsByDate(ctx context.Context, userID int64, date *pbcommon.Date) ([]*pbhealth.MentalHealthLog, error) {
	toReturn := make([]*pbhealth.MentalHealthLog, 0)

	filter := bson.M{"userid": userID, "logdate": date}

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
		toReturn = append(toReturn, healthLog)
	}

	return toReturn, err
}



