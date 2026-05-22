package main

import (
	pb "ITK_Code/m/v1/proto"

	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	client := pb.NewChatAPIClient(conn)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	createResponse, err := client.Create(
		ctx,
		&pb.CreateRequest{
			Usernames: []string{"giil", "fiiil", "liil"},
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(createResponse.Id)

	deleteResponse, err := client.Delete(
		ctx,
		&pb.DeleteRequest{
			Id: 1,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("deleted ", deleteResponse)

	sendResponse, err := client.Send(
		ctx,
		&pb.SendMessageRequest{
			From:      "Igor",
			Text:      "hello world",
			Timestamp: timestamppb.Now(),
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("message delivered", sendResponse)
}
