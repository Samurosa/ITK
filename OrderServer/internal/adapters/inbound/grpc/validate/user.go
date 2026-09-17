package validate

import "ITK_Code/m/v2/internal/core/dto"

func AccessCheck(userRole string) bool {
	switch userRole {
	case string(dto.AdminRole):
		return false
	default:
		return true
	}
}
