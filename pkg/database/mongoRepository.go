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
	"log"
	"math"
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
		m.logger.Infof("Error adding mental health log: %v", err)
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

func (m *MongoRepository) GetOverallScore(ctx context.Context, userID int64) (int32, error) {
	logs, err := m.GetAllMentalHealthLogs(ctx, userID)
	log.Printf("Logs fetched for user (ID = %v\n): %v\n", userID, logs)

	if err != nil {
		log.Fatal("cannot get mental health logs for user: ", err)
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

	log.Printf("Total sum of log scores for user (ID = %v\n): %v\n", userID, totalScore)
	log.Printf("Number of total logs for user (ID = %v): %v\n:", userID, numLogs)
	log.Printf("Average score for user (ID = %v): %v\n:", userID, overallScore)

	return overallScore, err
}

func (m *MongoRepository) DeleteMentalHealthLogs(ctx context.Context, userID int64, date *pbcommon.Date, all bool) (uint32, error) {

	var filter bson.M // declaring vbariable

	if all {
		filter = bson.M{"userid": userID} // filtering by user id and date
	} else {
		filter = bson.M{"userid": userID, "logdate": date} // filtering by user id and date
	}

	res, err := m.fileCollection.DeleteMany(ctx, filter) // deleting all health logs with the given date

	if err != nil {
		log.Fatal("cannot delete mental health logs for user: ", err)
	}

	numDeleted := uint32(res.DeletedCount) // getting number of entries deleted

	return numDeleted, err
}

func (m *MongoRepository) UpdateMentalHealthLogs(ctx context.Context, userID int64, healthLog *pbhealth.MentalHealthLog) error {

	filter := bson.M{"userid": userID, "logdate": healthLog.LogDate}

	_, err := m.fileCollection.UpdateOne(ctx, filter, healthLog)
	if err != nil {
		m.logger.Infof("Error updating mentalh health log: %v", err)
	}

	return err
}


