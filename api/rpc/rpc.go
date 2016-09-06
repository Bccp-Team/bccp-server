package rpc

import (
	"net"

	"google.golang.org/grpc"

	pb "github.com/Bccp-Team/bccp-server/proto/api"
)

type server struct{}

func SetupRpc(service string) error {
	lis, err := net.Listen("tcp", service)
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	pb.RegisterApiServer(s, &server{})
	s.Serve(lis)
	return nil
}
