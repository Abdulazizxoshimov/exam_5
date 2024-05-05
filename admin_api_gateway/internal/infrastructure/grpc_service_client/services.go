package grpc_service_clients

import (
	"fmt"

	pbu "exam_5/admin_api_gateway/genproto/clientProto"
	pb "exam_5/admin_api_gateway/genproto/jobProto"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"

	"exam_5/admin_api_gateway/internal/mock_data"
	"exam_5/admin_api_gateway/internal/pkg/config"
)

type ServiceClient interface {
	MockDataJob()   mock_data.MockServiceJob
	MockDataClient()   mock_data.MockServiceClient
	ClientService() pbu.UserServiceClient
	JobService()  pb.JobServiceClient
	Close()
}

type serviceClient struct {
	connections    []*grpc.ClientConn
	clientService pbu.UserServiceClient
	jobService pb.JobServiceClient
	mockJobService  mock_data.MockServiceJob
	mockClientService mock_data.MockServiceClient
}

func New(cfg *config.Config) (ServiceClient, error) {
	connUserService, err := grpc.Dial(
		fmt.Sprintf("%s%s", cfg.ClientService.Host, cfg.ClientService.Port),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
	if err != nil {
		return nil, err
	}
	connJobService, err := grpc.Dial(
		fmt.Sprintf("%s%s", cfg.JobService.Host, cfg.JobService.Port),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
	if err != nil {
		return nil, err
	}

	return &serviceClient{
		mockJobService: mock_data.NewMockServiceJob(),
		clientService: pbu.NewUserServiceClient(connUserService),
		jobService:  pb.NewJobServiceClient(connJobService),
		connections: []*grpc.ClientConn{
			connUserService,
		},
	}, nil
}

func (s *serviceClient)MockDataJob()mock_data.MockServiceJob{
	return s.mockJobService
}
func (s *serviceClient)MockDataClient()mock_data.MockServiceClient{
	return s.mockClientService
}
func (s *serviceClient) ClientService() pbu.UserServiceClient {
	return s.clientService
}
func (s *serviceClient) JobService() pb.JobServiceClient {
	return s.jobService
}
func (s *serviceClient) Close() {
	for _, conn := range s.connections {
		if err := conn.Close(); err != nil {
			// should be replaced by logger soon
			fmt.Printf("error while closing grpc connection: %v", err)
		}
	}
}
