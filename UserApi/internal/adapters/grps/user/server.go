package user

import (
	"ITK_Code/m/v2/internal/core/user"

	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/user"

	"google.golang.org/grpc"
)

type serverApi struct {
	pb.UnimplementedUserServiceServer
	user user.Service
}

func RegisterUserService(grpc *grpc.Server, user user.Service) {
	pb.RegisterUserServiceServer(grpc, &serverApi{user: user})
}

// коды ошибок сделать маппер под ошибки и изменить ошибки в интерфейс методах
// добавить логи для внутренних ошибок
