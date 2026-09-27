package spot

import (
	"ITK_Code/m/v2/internal/core/dto"

	"github.com/Samurosa/exchange-contract/protobuf/gen/go/shared"
)

func FromProtoToRole(protoRole shared.Role) dto.Role {
	switch protoRole {
	case shared.Role_ROLE_UNSPECIFIED:
		return dto.UnspecifiedRole
	case shared.Role_ROLE_USER:
		return dto.UserRole
	case shared.Role_ROLE_GUEST:
		return dto.GuestRole
	case shared.Role_ROLE_PREMIUM:
		return dto.PremiumRole
	case shared.Role_ROLE_ADMIN:
		return dto.AdminRole
	default:
		return dto.UnspecifiedRole
	}
}

func FromProtoToRoles(protoRole []shared.Role) []dto.Role {
	result := make([]dto.Role, 0, len(protoRole))

	for _, role := range protoRole {
		result = append(result, FromProtoToRole(role))
	}

	return result
}
