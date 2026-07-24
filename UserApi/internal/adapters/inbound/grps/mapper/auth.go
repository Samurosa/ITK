package mapper

import (
	"ITK_Code/m/v2/internal/core/auth"

	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/user"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToProtoTokens(tokens auth.TokensModel) *pb.TokenPairResponse {
	return &pb.TokenPairResponse{
		AccessToken:      tokens.AccessToken,
		RefreshToken:     tokens.RefreshToken,
		AccessExpiresAt:  timestamppb.New(tokens.AccessExpiresAt),
		AccessCreatesAt:  timestamppb.New(tokens.AccessCreatedAt),
		RefreshExpiresAt: timestamppb.New(tokens.RefreshExpiresAt),
		RefreshCreatesAt: timestamppb.New(tokens.RefreshCreatedAt),
	}
}
