package mapper

import (
	"ITK_Code/m/v2/internal/core/dto"

	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/user"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToProtoTokens(tokens dto.TokensModel) *pb.TokenPairResponse {
	return &pb.TokenPairResponse{
		AccessToken:      tokens.AccessToken,
		RefreshToken:     tokens.RefreshToken,
		AccessExpiresAt:  timestamppb.New(tokens.AccessExpiresAt),
		AccessIssuedAt:   timestamppb.New(tokens.AccessIssuedAt),
		RefreshExpiresAt: timestamppb.New(tokens.RefreshExpiresAt),
		RefreshIssuedAt:  timestamppb.New(tokens.RefreshIssuedAt),
	}
}
