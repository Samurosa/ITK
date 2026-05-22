package handler

import (
	db "ITK_Code/m/v2/internal/storage"
	pb "ITK_Code/m/v2/proto"
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	mu   sync.RWMutex
	data *db.UserData
}

func NewUserService() *UserService {
	return &UserService{
		data: db.NewUserStorage(),
	}
}

func (h *UserService) Create(
	ctx context.Context,
	req *pb.CreateUserRequest,
) (*pb.CreateUserResponse, error) {
	if _, err := fmt.Printf("Name: %s\nEmail: %s\nPassword: %s\n PasswordConfirm: %s\nRole: %s\n",
		req.Name,
		req.Email,
		req.Password,
		req.PasswordConfirm,
		req.Role.String(),
	); err != nil {
		return nil, err
	}

	id := uuid.New().String()
	h.mu.Lock()
	defer h.mu.Unlock()
	h.data.Data[id] = db.User{
		ID:              id,
		Name:            req.Name,
		Email:           req.Email,
		Password:        req.Password,
		PasswordConfirm: req.PasswordConfirm,
		Role:            req.Role,
	}
	return &pb.CreateUserResponse{
		Id: id,
	}, nil
}

func (h *UserService) Get(
	ctx context.Context,
	req *pb.GetRequest,
) (*pb.GetResponse, error) {
	h.mu.RLock()
	obj, ok := h.data.Data[req.Id]
	h.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("user not found")
	}

	return &pb.GetResponse{
		Id:    req.Id,
		Name:  obj.Name,
		Email: obj.Email,
		Role:  obj.Role,
	}, nil
}

func (h *UserService) Update(
	ctx context.Context,
	req *pb.UpdateRequest,
) (*pb.UpdateResponse, error) {

	if _, err := fmt.Println("ID user updated: ", req.Id); err != nil {
		return nil, err
	}

	return &pb.UpdateResponse{}, nil
}

func (h *UserService) Delete(
	ctx context.Context,
	req *pb.DeleteRequest,
) (*pb.DeleteResponse, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.data.Data[req.Id]; !ok {
		return nil, fmt.Errorf("object not found, deletion error")
	}

	delete(h.data.Data, req.Id)

	if _, err := fmt.Println("ID user deleted: ", req.Id); err != nil {
		return nil, err
	}

	return &pb.DeleteResponse{}, nil
}
