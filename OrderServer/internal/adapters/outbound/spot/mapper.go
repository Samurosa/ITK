package spot

import (
	"ITK_Code/m/v2/internal/core/dto"

	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/user"
)

func FromProtoToRole(protoRole pb.Role) dto.Role {
	switch protoRole {
	case pb.Role_ROLE_UNSPECIFIED:
		return dto.UnspecifiedRole
	case pb.Role_ROLE_USER:
		return dto.UserRole
	case pb.Role_ROLE_GUEST:
		return dto.GuestRole
	case pb.Role_ROLE_PREMIUM:
		return dto.PremiumRole
	case pb.Role_ROLE_ADMIN:
		return dto.AdminRole
	default:
		return dto.UnspecifiedRole
	}
}

func FromProtoToRoles(protoRole []pb.Role) []dto.Role {
	result := make([]dto.Role, 0, len(protoRole))

	for _, role := range protoRole {
		result = append(result, FromProtoToRole(role))
	}

	return result
}
