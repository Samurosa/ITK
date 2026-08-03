package dto

type ContextKey string

const (
	UserIDContextKey ContextKey = "user_id"
	RoleContextKey   ContextKey = "role"
	DeviceContextKey ContextKey = "device"
	JtiContextKey    ContextKey = "jti"
	ClientIPKey      ContextKey = "client_ip"
)
