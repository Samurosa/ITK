package mapper

import (
	"ITK_Code/m/v2/internal/core/user"

	"github.com/Samurosa/exchange-contract/protobuf/gen/go/shared"
)

func ToProtoRole(role user.Role) shared.Role {
	switch role {
	case user.UserRole:
		return shared.Role_ROLE_USER
	case user.GuestRole:
		return shared.Role_ROLE_GUEST
	case user.PremiumRole:
		return shared.Role_ROLE_PREMIUM
	case user.AdminRole:
		return shared.Role_ROLE_ADMIN
	default:
		return shared.Role_ROLE_UNSPECIFIED
	}
}
