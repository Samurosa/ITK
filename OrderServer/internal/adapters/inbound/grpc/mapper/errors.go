package mapper

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ToGRPC(err error) error {
	switch {

	default:
		return status.Error(codes.Internal, "internal SpotServer error")
	}
}
