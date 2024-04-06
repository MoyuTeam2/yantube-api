package services

import (
	pb "api/server/rpc/model"
	"context"
)

type HelloService struct {
	pb.UnimplementedHelloServiceServer
}

func (h *HelloService) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Message: "Hello " + req.Name}, nil
}
