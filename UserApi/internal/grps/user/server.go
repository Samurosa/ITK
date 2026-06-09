package user

import (
	"context"
	"time"

	mapper "ITK_Code/m/v2/internal/mapper"
	data "ITK_Code/m/v2/internal/storage"

	pb "github.com/Truncklin/exchange-contract/protobuf/gen/go/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type UserServer interface {
	Register(ctx context.Context,
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
		balances map[string]*data.Balance,
		role data.Role,
		createdAt time.Time,
		updatedAt time.Time,
		err error,
	)

	UpdateUser(ctx context.Context,
		id string,
		name *wrapperspb.StringValue,
		login *wrapperspb.StringValue,
		password *wrapperspb.StringValue,
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
		balance *data.Balance,
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
	userServer UserServer
}

func NewServer(grps *grpc.Server, userServer UserServer) {
	pb.RegisterUserServiceServer(grps, &serverApi{userServer: userServer})
}

func (s *serverApi) Register(
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

	id, createdAt, err := s.userServer.Register(ctx, req.Login, req.Password)
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
		err := s.userServer.GetUser(ctx, req.Id)
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

	updated, updatedAt, err := s.userServer.UpdateUser(
		ctx,
		req.Id,
		req.Name,
		req.Login,
		req.Password,
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

	success, err := s.userServer.DeleteUser(ctx, req.Id)
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

	success, balance, err := s.userServer.Deposit(
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

	token, err := s.userServer.Authorization(ctx, req.Login, req.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to authorize user")
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
