package requestContext

type Principal struct {
	UserID string
	Role   string
}

type RequestMetadata struct {
	ClientIP string
	DeviceID string
}

type RequestContext struct {
	Principal Principal
	Metadata  RequestMetadata

	JTI string
}
