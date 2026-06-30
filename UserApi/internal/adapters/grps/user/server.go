package user

import (
	"ITK_Code/m/v2/internal/core/auth"
	"ITK_Code/m/v2/internal/core/user"
	"ITK_Code/m/v2/internal/core/wallet"

	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/user"
	"go.uber.org/zap"

	"google.golang.org/grpc"
)

type ServerApi struct {
	pb.UnimplementedUserServiceServer
	user   user.Service
	auth   auth.Service
	wallet wallet.Service
	log    *zap.Logger
}

func RegisterUserService(grpc *grpc.Server, user user.Service, auth auth.Service, wallet wallet.Service, log *zap.Logger) {
	pb.RegisterUserServiceServer(grpc, &ServerApi{user: user, auth: auth, wallet: wallet, log: log})
}
