package services

import (
	"api/config"
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
	logger := h.log.With().Str("api", "Register").Str("host", in.GetHost()).Logger()
	logger.Info().Msg("received stream server register request")

	if in.GetHost() == "" {
		logger.Error().Msg("invalid host")
		return &pb.StreamServerRegisterResponse{
			Success: false,
			Message: "invalid host",
		}, nil
	}

	if config.Config.StreamServer.Secret != "" && in.GetSecret() != config.Config.StreamServer.Secret {
		logger.Error().Msg("invalid secret")
		return &pb.StreamServerRegisterResponse{
			Success: false,
			Message: "invalid secret",
		}, nil
	}

	err := db.Get().RegisterStreamServer(in.GetHost())
	if err != nil {
		logger.Error().Err(err).Msg("failed to register stream server")
		return &pb.StreamServerRegisterResponse{
			Success: false,
			Message: "failed to register stream server",
		}, nil
	}

	return &pb.StreamServerRegisterResponse{Success: true}, nil
}

func (h *StreamServerService) Unregister(ctx context.Context, in *pb.StreamServerRegisterRequest) (*pb.StreamServerRegisterResponse, error) {
	logger := h.log.With().Str("api", "Unregister").Str("host", in.GetHost()).Logger()
	logger.Info().Msg("received stream server unregister request")

	if in.GetHost() == "" {
		logger.Error().Msg("invalid host")
		return &pb.StreamServerRegisterResponse{
			Success: false,
			Message: "invalid host",
		}, nil
	}

	if config.Config.StreamServer.Secret != "" && in.GetSecret() != config.Config.StreamServer.Secret {
		logger.Error().Msg("invalid secret")
		return &pb.StreamServerRegisterResponse{
			Success: false,
			Message: "invalid secret",
		}, nil
	}

	err := db.Get().UnregisterStreamServer(in.GetHost())
	if err != nil {
		logger.Error().Err(err).Msg("failed to unregister stream server")
		return &pb.StreamServerRegisterResponse{
			Success: false,
			Message: "failed to unregister stream server",
		}, nil
	}

	return &pb.StreamServerRegisterResponse{Success: true}, nil
}

func (h *StreamServerService) KeepAlive(ctx context.Context, in *pb.StreamServerKeepAliveRequest) (*pb.StreamServerKeepAliveResponse, error) {
	logger := h.log.With().Str("api", "KeepAlive").Str("host", in.GetHost()).Logger()
	logger.Info().Msg("received stream server keep alive request")

	if in.GetHost() == "" {
		logger.Error().Msg("invalid host")
		return &pb.StreamServerKeepAliveResponse{
			Success: false,
			Message: "invalid host",
		}, nil
	}

	if config.Config.StreamServer.Secret != "" && in.GetSecret() != config.Config.StreamServer.Secret {
		logger.Error().Msg("invalid secret")
		return &pb.StreamServerKeepAliveResponse{
			Success: false,
			Message: "invalid secret",
		}, nil
	}

	err := db.Get().StreamServerKeepLive(in.GetHost())
	if err != nil {
		logger.Error().Err(err).Msg("failed to keep alive stream server")
		return &pb.StreamServerKeepAliveResponse{
			Success: false,
			Message: "failed to keep alive stream server",
		}, nil
	}

	return &pb.StreamServerKeepAliveResponse{Success: true}, nil
}
