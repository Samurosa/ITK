package user

import (
	"context"
	"time"

	"ITK_Code/m/v2/internal/domain/models"
	"ITK_Code/m/v2/internal/mapper"

	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type User interface {
	Registration(ctx context.Context,
		login string,
		password string,
	) (
		id string,
		createdAt time.Time,
		err error,
	)

	GetUser(ctx context.Context,
		id string,
	) (
		name string,
		login string,
		balances map[string]*models.Balance,
		role models.Role,
		createdAt time.Time,
		updatedAt time.Time,
		err error,
	)

	UpdateUser(ctx context.Context,
		id string,
		name string,
		login string,
		password string,
	) (
		updated bool,
		updatedAt time.Time,
		err error,
	)

	DeleteUser(ctx context.Context,
		id string,
	) (
		success bool,
		err error,
	)

	Deposit(ctx context.Context,
		id string,
		asset string,
		amount string,
	) (
		success bool,
		balance *models.Balance,
		err error,
	)

	Authorization(ctx context.Context,
		login string,
		password string,
	) (
		token string,
		err error,
	)
}

type serverApi struct {
	pb.UnimplementedUserServiceServer
	user User
}

func RegisterUserService(grpc *grpc.Server, user User) {
	pb.RegisterUserServiceServer(grpc, &serverApi{user: user})
}

func (s *serverApi) Registration(
	ctx context.Context,
	req *pb.RegisterUserRequest,
) (
	*pb.RegisterUserResponse,
	error,
) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := ValidateLoginPassword(req.Login, req.Password); err != nil {
		return nil, err
	}

	id, createdAt, err := s.user.Registration(ctx, req.Login, req.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to register user")
	}

	return &pb.RegisterUserResponse{
		Id:        id,
		CreatedAt: timestamppb.New(createdAt),
	}, nil
}

func (s *serverApi) GetUser(
	ctx context.Context,
	req *pb.GetUserRequest,
) (
	*pb.GetUserResponse,
	error,
) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := ValidateUserId(req.Id); err != nil {
		return nil, err
	}

	name, login, balances, role, createdAt, updatedAt,
		err := s.user.GetUser(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	return &pb.GetUserResponse{
		Name:      name,
		Login:     login,
		Balances:  mapper.ToProtoBalances(balances),
		Role:      mapper.ToProtoRole(role),
		CreatedAt: timestamppb.New(createdAt),
		UpdatedAt: timestamppb.New(updatedAt),
	}, nil
}

func (s *serverApi) UpdateUser(
	ctx context.Context,
	req *pb.UpdateUserRequest,
) (
	*pb.UpdateUserResponse,
	error,
) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := ValidateUserId(req.Id); err != nil {
		return nil, err
	}

	name := ""
	login := ""
	password := ""
	if req.Name != nil {
		name = req.Name.Value
	}
	if req.Login != nil {
		login = req.Login.Value
	}
	if req.Password != nil {
		password = req.Password.Value
	}

	updated, updatedAt, err := s.user.UpdateUser(
		ctx,
		req.Id,
		name,
		login,
		password,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to update user")
	}

	return &pb.UpdateUserResponse{
		Updated:   updated,
		UpdatedAt: timestamppb.New(updatedAt),
	}, nil
}

func (s *serverApi) DeleteUser(
	ctx context.Context,
	req *pb.DeleteUserRequest,
) (
	*pb.DeleteUserResponse,
	error,
) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := ValidateUserId(req.Id); err != nil {
		return nil, err
	}

	success, err := s.user.DeleteUser(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to delete user")
	}
	return &pb.DeleteUserResponse{Success: success}, nil
}

func (s *serverApi) Deposit(
	ctx context.Context,
	req *pb.DepositRequest,
) (
	*pb.DepositResponse,
	error,
) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err := ValidateUserId(req.Id); err != nil {
		return nil, err
	}
	if err := ValidateDepositRequest(req); err != nil {
		return nil, err
	}

	success, balance, err := s.user.Deposit(
		ctx,
		req.Id,
		req.Asset,
		req.Amount,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to deposit funds")
	}
	return &pb.DepositResponse{
		Success: success,
		Balance: mapper.ToProtoBalance(balance),
	}, nil
}

func (s *serverApi) Authorization(
	ctx context.Context,
	req *pb.AuthorizationRequest,
) (
	*pb.AuthorizationResponse,
	error,
) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := ValidateLoginPassword(req.Login, req.Password); err != nil {
		return nil, err
	}

	token, err := s.user.Authorization(ctx, req.Login, req.Password)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "failed to authorize user")
	}

	return &pb.AuthorizationResponse{Token: token}, nil
}

func ValidateUserId(id string) error {
	if id == "" {
		return status.Error(codes.InvalidArgument, "Id is required")
	}
	return nil
}

func ValidateDepositRequest(req *pb.DepositRequest) error {
	if req.GetId() == "" {
		return status.Error(codes.InvalidArgument, "UserId is required")
	}
	if req.GetAsset() == "" {
		return status.Error(codes.InvalidArgument, "Asset is required")
	}
	if req.GetAmount() == "" {
		return status.Error(codes.InvalidArgument, "Amount is required")
	}
	return nil
}

func ValidateLoginPassword(login string, password string) error {
	if login == "" {
		return status.Error(codes.InvalidArgument, "Login is required")
	}
	if password == "" {
		return status.Error(codes.InvalidArgument, "Password is required")
	}
	return nil
}
