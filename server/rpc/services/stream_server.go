package services

import (
	"api/db"
	pb "api/server/rpc/model"
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type StreamServerService struct {
	log zerolog.Logger
	pb.UnimplementedStreamServerServer
}

func NewStreamServerService() *StreamServerService {
	return &StreamServerService{
		log: log.With().Str("service", "StreamServerService").Logger(),
	}
}

func (h *StreamServerService) Register(ctx context.Context, in *pb.StreamServerRegisterRequest) (*pb.StreamServerRegisterResponse, error) {
	h.log.Info().Str("host", in.GetHost()).Msg("received stream server register request")

	err := db.Get().RegisterStreamServer(in.GetHost())
	if err != nil {
		h.log.Error().Err(err).Str("host", in.GetHost()).Msg("failed to register stream server")
		return &pb.StreamServerRegisterResponse{Success: false}, nil
	}

	return &pb.StreamServerRegisterResponse{Success: true}, nil
}

func (h *StreamServerService) Unregister(ctx context.Context, in *pb.StreamServerRegisterRequest) (*pb.StreamServerRegisterResponse, error) {
	h.log.Info().Str("host", in.GetHost()).Msg("received stream server unregister request")

	err := db.Get().UnregisterStreamServer(in.GetHost())
	if err != nil {
		h.log.Error().Err(err).Str("host", in.GetHost()).Msg("failed to unregister stream server")
		return &pb.StreamServerRegisterResponse{Success: false}, nil
	}

	return &pb.StreamServerRegisterResponse{Success: true}, nil
}

func (h *StreamServerService) KeepAlive(ctx context.Context, in *pb.StreamServerKeepAliveRequest) (*pb.StreamServerKeepAliveResponse, error) {
	h.log.Info().Str("host", in.GetHost()).Msg("received stream server keep alive request")

	err := db.Get().StreamServerKeepLive(in.GetHost())
	if err != nil {
		h.log.Error().Err(err).Str("host", in.GetHost()).Msg("failed to keep alive stream server")
		return &pb.StreamServerKeepAliveResponse{Success: false}, nil
	}

	return &pb.StreamServerKeepAliveResponse{Success: true}, nil
}
