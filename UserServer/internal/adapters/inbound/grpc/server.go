package grpc

import (
	"ITK_Code/m/v2/internal/application"

	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/user"
	"go.uber.org/zap"

	"google.golang.org/grpc"
)

type UserServer struct {
	pb.UnimplementedUserServiceServer
	user application.User
	auth application.Auth
	log  *zap.Logger
}

func RegisterUserService(grpc *grpc.Server, user application.User, auth application.Auth, log *zap.Logger) {
	pb.RegisterUserServiceServer(grpc, &UserServer{user: user, auth: auth, log: log})
}
