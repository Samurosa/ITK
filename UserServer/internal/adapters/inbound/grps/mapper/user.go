package mapper

import (
	"ITK_Code/m/v2/internal/core/dto"

	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/user"
)

func ToProtoRole(role dto.Role) pb.Role {
	switch role {
	case dto.UserRole:
		return pb.Role_ROLE_USER
	case dto.GuestRole:
		return pb.Role_ROLE_GUEST
	case dto.PremiumRole:
		return pb.Role_ROLE_PREMIUM
	case dto.AdminRole:
		return pb.Role_ROLE_ADMIN
	default:
		return pb.Role_ROLE_UNSPECIFIED
	}
}
