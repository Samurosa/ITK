package handler

import (
	pb "ITK_Code/m/v1/proto"
	"context"
	"fmt"
)

type ChatApi struct {
	pb.UnimplementedChatAPIServer
}

func NewChatApi() *ChatApi {
	return &ChatApi{}
}

func (h *ChatApi) Create(
	ctx context.Context,
	req *pb.CreateRequest,
) (*pb.CreateResponse, error) {
	//сохраняем юзеров
	fmt.Println(req.Usernames)
	return &pb.CreateResponse{Id: 1}, nil
}

func (h *ChatApi) Delete(
	ctx context.Context,
	req *pb.DeleteRequest,
) (*pb.DeleteResponse, error) {
	//удаляем чат по id
	fmt.Println("delete user: ", req.Id)
	return &pb.DeleteResponse{}, nil
}

func (h *ChatApi) Send(
	ctx context.Context,
	req *pb.SendMessageRequest,
) (*pb.SendMessageResponse, error) {
	fmt.Printf("от %s,\n %s\n time: %v", req.From, req.Text, req.Timestamp.AsTime())

	return &pb.SendMessageResponse{}, nil
}
