package main

import (
	pb "ITK_Code/m/v2/proto"

	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	createResp, err := client.Create(
		ctx,
		&pb.CreateUserRequest{
			Name: "Alex",

			Email: "alex@mail.ru",

			Password: "123",

			PasswordConfirm: "123",

			Role: pb.Role_USER,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(
		"created id:",
		createResp.Id,
	)

	getResp, err := client.Get(
		ctx,
		&pb.GetRequest{
			Id: createResp.Id,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(
		"%+v\n",
		getResp,
	)

	_, err = client.Update(
		ctx,
		&pb.UpdateRequest{
			Id: createResp.Id,
		},
	)

	_, err = client.Delete(
		ctx,
		&pb.DeleteRequest{
			Id: createResp.Id,
		},
	)
}
