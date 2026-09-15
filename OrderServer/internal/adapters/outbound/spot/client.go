package spot

import (
	"context"

	"ITK_Code/m/v2/internal/core/dto"

	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/spot"

	"google.golang.org/grpc"
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
	if err != nil {
		return dto.Spot{}, err
	}

	return dto.Spot{
		ID:           resp.Id,
		BaseAsset:    resp.BaseAsset,
		QuoteAsset:   resp.QuoteAsset,
		Status:       dto.SpotStatus(resp.Status.String()),
		AllowedRoles: FromProtoToRoles(resp.AllowedRoles),
	}, nil
}
