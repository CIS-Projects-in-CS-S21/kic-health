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
	conn, err := grpc.Dial("test.api.keeping-it-casual.com:50051", grpc.WithInsecure())
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

	// creating auth context
	md := metadata.Pairs("Authorization", fmt.Sprintf("Bearer %v", tokRes.Token))
	authCtx := metadata.NewOutgoingContext(context.Background(), md)

	// -----------------

	// Adding health log for user

	addReq :=&pbhealth.AddHealthDataForUserRequest{
		UserID:   1,
		NewEntry: &pbhealth.MentalHealthLog{
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

	addRes, err := client.AddHealthDataForUser(authCtx, addReq)
	if err != nil {
		log.Fatal("cannot upload mental health log: ", err)
	}
	log.Printf("addRess: %v\n", addRes)


	// --------------------------------------

}
