package services

import (
	"context"

	"github.com/wenchy/grpcio/internal/corepb"
)

// ClientTestService implements the protobuf interface
type ClientTestService struct {
	// Embed the unimplemented server
	corepb.UnimplementedClientTestServiceServer
}

// NewClientTestService initializes a new ClientTestService struct.
func NewClientTestService() *ClientTestService {
	return &ClientTestService{}
}

// Echo transfer back the request.
func (s *ClientTestService) Echo(ctx context.Context, req *corepb.EchoRequest) (*corepb.EchoResponse, error) {
	rsp := &corepb.EchoResponse{
		Msg: req.Msg,
	}
	return rsp, nil
}

// Greet greet back the request.
func (s *ClientTestService) Greet(ctx context.Context, req *corepb.GreetRequest) (*corepb.GreetResponse, error) {
	rsp := &corepb.GreetResponse{
		Name: req.Name,
		Age:  req.Age,
	}
	return rsp, nil
}
