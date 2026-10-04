package spot

import (
	"context"
	"fmt"

	"ITK_Code/m/v2/internal/core/corerrors"
	"ITK_Code/m/v2/internal/core/dto"

	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/spot"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Client struct {
	client pb.SpotInstrumentServiceClient
}

func NewClient(conn *grpc.ClientConn) *Client {
	return &Client{
		client: pb.NewSpotInstrumentServiceClient(conn),
	}
}

func (c *Client) GetSpot(ctx context.Context, spotID string) (dto.Spot, error) {
	resp, err := c.client.GetSpot(ctx,
		&pb.GetSpotRequest{
			Id: spotID,
		})
	if status.Code(err) == codes.NotFound {
		return dto.Spot{}, corerrors.ErrSpotNotFound
	}
	if status.Code(err) == codes.Canceled {
		return dto.Spot{}, fmt.Errorf("get spot: %w", context.Canceled)
	}
	if status.Code(err) == codes.DeadlineExceeded && ctx.Err() != nil {
		return dto.Spot{}, fmt.Errorf("get spot: %w", ctx.Err())
	}
	if err != nil {
		return dto.Spot{}, fmt.Errorf("get spot: %w", err)
	}

	return dto.Spot{
		ID:                resp.Id,
		BaseAsset:         resp.BaseAsset,
		QuoteAsset:        resp.QuoteAsset,
		PricePrecision:    resp.PricePrecision,
		QuantityPrecision: resp.QuantityPrecision,
		MinOrderSize:      resp.MinOrderSize,
		MaxOrderSize:      resp.MaxOrderSize,
		Status:            dto.SpotStatus(resp.Status.String()),
		AllowedRoles:      FromProtoToRoles(resp.AllowedRoles),
	}, nil
}
