package services

import (
	"context"
	"sync"
	"time"

	mapper "ITK_Code/m/v2/internal/mapper"
	db "ITK_Code/m/v2/internal/storage"

	pb "github.com/Truncklin/exchange-contract/protobuf/gen/go/user"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type User struct {
	log *zap.Logger
	pb.UnimplementedUserServiceServer
	mu   sync.RWMutex
	data *db.UserRepository
}

func NewUserService() *User {
	return &User{
		data: db.NewUserStorage(),
	}
}

func (h *User) Create(
	ctx context.Context,
	req *pb.CreateUserRequest,
) (*pb.CreateUserResponse, error) {

	id := uuid.New().String()
	now := time.Now()

	user := db.User{
		ID:         id,
		Name:       req.Name,
		Login:      req.Login,
		Balances:   make(map[string]*db.Balance),
		Role:       db.Role(req.Role),
		CreateTime: now,
		UpdateTime: now,
	}
	h.data.Create(user)

	return &pb.CreateUserResponse{
		User:      mapper.ToProto(user),
		CreatedAt: timestamppb.New(now),
	}, nil
}

func (h *UserService) Get(
	ctx context.Context,
	req *pb.GetRequest,
) (*pb.GetResponse, error) {

	user, ok := h.data.Get(req.Id)
	if !ok {
		return nil, status.Error(
			codes.NotFound,
			"user not found",
		)
	}

	return &pb.GetResponse{
		Id:        user.ID,
		Name:      user.Name,
		Login:     user.Login,
		Balances:  mapper.ToProtoBalances(user.Balances),
		Role:      mapper.ToProtoRole(user.Role),
		CreatedAt: timestamppb.New(user.CreateTime),
		UpdatedAt: timestamppb.New(user.UpdateTime),
	}, nil
}

func (h *UserService) Delete(
	ctx context.Context,
	req *pb.DeleteRequest,
) (*pb.DeleteResponse, error) {

	ok := h.data.Delete(req.Id)
	if !ok {
		return nil, status.Error(
			codes.NotFound,
			"user not found",
		)
	}

	return &pb.DeleteResponse{}, nil
}

func (h *UserService) Deposit(
	ctx context.Context,
	req *pb.DepositRequest,
) (*pb.DepositResponse,
	error,
) {
	bal, ok := h.data.Deposit(
		req.UserId,
		req.Asset,
		req.Amount,
	)
	if !ok {
		return nil, status.Error(
			codes.NotFound,
			"user not found",
		)
	}
	return &pb.DepositResponse{
		Balance: &pb.Balance{
			Asset:     bal.Asset,
			Available: bal.Available,
			Locked:    bal.Locked,
		},
	}, nil
}
