package grpc

import (
	"ITK_Code/m/v2/internal/core/dto"
	"ITK_Code/m/v2/internal/core/user"
	"context"
	"time"

	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/user"

	"google.golang.org/grpc"
)

type UserService interface {
	GetUser(context.Context, string) (user.User, error)
	UpdateUserInfo(context.Context, string, string) error
	DeleteUser(context.Context, string) error
	ChangePassword(context.Context, string, string, string) error
}

type AuthService interface {
	Registration(context.Context, string, string, string, string, string) (string, time.Time, error)
	Login(context.Context, string, string, string, string) (dto.TokensModel, error)
	Logout(context.Context, string, string) error
	LogoutAllDevices(context.Context, string) error
	RefreshToken(context.Context, string) (dto.TokensModel, error)
}

type UserServer struct {
	pb.UnimplementedUserServiceServer
	user UserService
	auth AuthService
}

func RegisterUserService(grpc *grpc.Server, user UserService, auth AuthService) {
	pb.RegisterUserServiceServer(grpc, &UserServer{user: user, auth: auth})
}
