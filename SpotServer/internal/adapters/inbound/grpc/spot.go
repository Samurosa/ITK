package grpc

import (
	"ITK_Code/m/v2/internal/adapters/inbound/grpc/mapper"
	"context"

	pb "github.com/Samurosa/exchange-contract/protobuf/gen/go/spot"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) CreateSpot(ctx context.Context, req *pb.CreateSpotRequest) (*pb.CreateSpotResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid argument error: "+err.Error())
	}

	reqSpot := mapper.FromProtoCreateSpot(req)

	idSpot, createdTo, err := s.spot.CreateSpot(ctx, reqSpot)
	if err != nil {
		return nil, mapper.ToGRPC(err)
	}

	return &pb.CreateSpotResponse{
		Id:        idSpot,
		CreatedAt: timestamppb.New(createdTo),
	}, nil
}

func (s *Server) GetSpot(ctx context.Context, req *pb.GetSpotRequest) (*pb.GetSpotResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid argument error: "+err.Error())
	}

	spot, err := s.spot.GetSpot(ctx, req.Id)
	if err != nil {
		return nil, mapper.ToGRPC(err)
	}

	return mapper.ToProtoSpot(spot), nil
}

func (s *Server) EnableSpot(ctx context.Context, req *pb.EnableSpotRequest) (*pb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid argument error: "+err.Error())
	}

	err := s.spot.EnableSpot(ctx, req.Id)
	if err != nil {
		return nil, mapper.ToGRPC(err)
	}

	return &pb.Empty{}, nil
}

func (s *Server) DisableSpot(ctx context.Context, req *pb.DisableSpotRequest) (*pb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid argument error: "+err.Error())
	}

	err := s.spot.DisableSpot(ctx, req.Id)
	if err != nil {
		return nil, mapper.ToGRPC(err)
	}

	return &pb.Empty{}, nil
}
