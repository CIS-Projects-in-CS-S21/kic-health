/*
This is for running integration tests in a production like environment
*/

package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"

	pbcommon "github.com/kic/health/pkg/proto/common"
	pbusers "github.com/kic/health/pkg/proto/users"
	pbhealth "github.com/kic/health/pkg/proto/health"
)

func main() {
	conn, err := grpc.Dial("api.keeping-it-casual.com:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pbhealth.NewHealthTrackingClient(conn)


	// User client for auth
	usersClient := pbusers.NewUsersClient(conn)

	// get JWT
	tokRes, err := usersClient.GetJWTToken(context.Background(), &pbusers.GetJWTTokenRequest{
		Username: "testuser",
		Password: "testpass",
	})
	// ------------

	// creating auth context
	md := metadata.Pairs("Authorization", fmt.Sprintf("Bearer %v", tokRes.Token))
	authCtx := metadata.NewOutgoingContext(context.Background(), md)
	// -----------------

	// Getting userID
	userRes, err := usersClient.GetUserByUsername(authCtx, &pbusers.GetUserByUsernameRequest{Username: "testuser"})
	userID := userRes.User.UserID

	log.Printf("UserID is %v\n", userID)
	// ----------------------

	// Adding health log for user
	shouldAddLog(authCtx, client, userID)
	// --------------------------------------

	// Getting all mental health logs for a user
	shouldGetAllLogs(authCtx, client, userID)
	// -----------------------

	// Getting all mental health logs for a user for a specific date
	shouldGetLogsForDate(authCtx, client, userID)
	// -----------------------

	// Getting overall score for user
	shouldGetScore(authCtx, client, userID)
	// ----------------------

	// Updating mental health logs for a user
	shouldUpdateLogs(authCtx, client, userID)
	// -----------------------

	// Deleting mental health logs for a user for a specific date
	shouldDeleteLogsForDate(authCtx, client, userID)
	// ----------------------

	// Adding another health log for user
	shouldAddLog(authCtx, client, userID)
	// --------------------------------------

	// Deleting all mental logs for a user, regardless of date
	shouldDeleteAllLogs(authCtx, client, userID)
	// -----------------------
}

func shouldAddLog(authCtx context.Context, client pbhealth.HealthTrackingClient, userID int64) {
	addReq :=&pbhealth.AddHealthDataForUserRequest{
		UserID:   userID,
		NewEntry: &pbhealth.MentalHealthLog{
			LogDate:     &pbcommon.Date{
				Year:  2021,
				Month: 4,
				Day:   26,
			},
			Score:       5,
			JournalName: "I am sad!",
			UserID:      userID,
		},
	}

	addRes, err := client.AddHealthDataForUser(authCtx, addReq)
	if err != nil {
		log.Fatal("cannot upload mental health log: ", err)
	}
	log.Printf("addRes: %v\n", addRes)

}

func shouldGetAllLogs(authCtx context.Context, client pbhealth.HealthTrackingClient, userID int64) {
	getAllReq := &pbhealth.GetHealthDataForUserRequest{UserID: userID}

	getAllRes, err := client.GetHealthDataForUser(authCtx, getAllReq)
	if err != nil {
		log.Fatal("cannot get mental health logs for user: ", err)
	}
	log.Printf("getAllRes: %v\n", getAllRes)
	log.Printf("Logs retrieved: %v\n", getAllRes.HealthData)
}

func shouldGetLogsForDate(authCtx context.Context, client pbhealth.HealthTrackingClient, userID int64) {
	getAllByDateReq := &pbhealth.GetHealthDataByDateRequest{UserID: userID, LogDate: &pbcommon.Date{
		Year:  2021,
		Month: 4,
		Day:   26,
	},
	}

	getAllByDateRes, err := client.GetHealthDataByDate(authCtx, getAllByDateReq)
	if err != nil {
		log.Fatal("cannot get mental health logs for user by date: ", err)
	}
	log.Printf("getAllByDateRes: %v\n", getAllByDateRes)
	log.Printf("Logs retrieved: %v\n", getAllByDateRes.HealthData)
}

func shouldGetScore(authCtx context.Context, client pbhealth.HealthTrackingClient, userID int64) {
	getScoreReq := &pbhealth.GetMentalHealthScoreForUserRequest{UserID: userID}
	getScoreRes, err := client.GetMentalHealthScoreForUser(authCtx, getScoreReq)

	if err != nil {
		log.Fatal("cannot get mental health score for user: ", err)
	}
	log.Printf("getScoreRes: %v\n", getScoreRes)
	log.Printf("Mental Health Score: %v\n", getScoreRes.Score)
}

func shouldDeleteAllLogs(authCtx context.Context, client pbhealth.HealthTrackingClient, userID int64) {
	deleteReq := &pbhealth.DeleteHealthDataForUserRequest{
		UserID: userID,
		Data:   &pbhealth.DeleteHealthDataForUserRequest_All{true},
	}

	deleteRes, err := client.DeleteHealthDataForUser(authCtx, deleteReq)

	if err != nil {
		log.Fatal("cannot delete mental health score for user: ", err)
	}

	log.Printf("deleteRes: %v\n", deleteRes)
}

func shouldDeleteLogsForDate(authCtx context.Context, client pbhealth.HealthTrackingClient, userID int64) {
	deleteReq2 := &pbhealth.DeleteHealthDataForUserRequest{
		UserID: userID,
		Data:   &pbhealth.DeleteHealthDataForUserRequest_DateToRemove{&pbcommon.Date{
			Year:  2021,
			Month: 4,
			Day:   26,
		}},
	}

	deleteRes2, err := client.DeleteHealthDataForUser(authCtx, deleteReq2)

	if err != nil {
		log.Fatal("cannot delete mental health score for user: ", err)
	}

	log.Printf("deleteRes: %v\n", deleteRes2)
}

func shouldUpdateLogs(authCtx context.Context, client pbhealth.HealthTrackingClient, userID int64) {
	updateReq := &pbhealth.UpdateHealthDataForDateRequest{
		UserID:         userID,
		DesiredLogInfo: &pbhealth.MentalHealthLog{
			LogDate:     &pbcommon.Date{
				Year:  2021,
				Month: 4,
				Day:   26,
			},
			Score:       5,
			JournalName: "I am so happy!!!!!!",
			UserID:      userID,
		},
	}

	updateRes, err := client.UpdateHealthDataForDate(authCtx, updateReq)
	if err != nil {
		log.Fatal("cannot update mental health score for user: ", err)
	}
	log.Printf("updateRes: %v\n", updateRes)
}