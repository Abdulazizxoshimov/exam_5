package grpc_service_clients

import (
	pb "exam_5/jobsService/genproto/clientProto"
	"exam_5/jobsService/internal/pkg/config"
	"fmt"

	"google.golang.org/grpc"
)

type ServiceClients interface {
	ClientService()  pb.UserServiceClient
	Close()
}

type serviceClients struct {
	clientService pb.UserServiceClient
	services []*grpc.ClientConn
}

func New(config *config.Config) (ServiceClients, error) {
	clientConn, err := grpc.Dial(
		fmt.Sprintf("%s:%s", config.ClientService.Host, config.ClientService.Port),
		grpc.WithInsecure())
	fmt.Println(err)
	if err != nil {
		return nil, fmt.Errorf("post service dial host: %s port: %s", config.ClientService.Host, config.ClientService.Port)
	}
	return &serviceClients{
		clientService: pb.NewUserServiceClient(clientConn),
		services: []*grpc.ClientConn{},
	}, nil
}

func (s *serviceClients) Close() {
	for _, conn := range s.services {
		conn.Close()
	}
}
func (s serviceClients)ClientService()pb.UserServiceClient{
	return s.clientService
}
