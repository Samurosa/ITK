package mapper

import (
	"ITK_Code/m/v2/internal/core/user"

	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/user"
)

func ToProtoRole(role user.Role) pb.Role {
	switch role {
	case user.UserRole:
		return pb.Role_ROLE_USER
	case user.GuestRole:
		return pb.Role_ROLE_GUEST
	case user.PremiumRole:
		return pb.Role_ROLE_PREMIUM
	case user.AdminRole:
		return pb.Role_ROLE_ADMIN
	default:
		return pb.Role_ROLE_UNSPECIFIED
	}
}
