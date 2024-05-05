package services

import (
	"context"
	client "exam_5/jobsService/genproto/clientProto"
	pb "exam_5/jobsService/genproto/jobProto"
	"exam_5/jobsService/internal/delivery/grpc"
	"exam_5/jobsService/internal/entity"
	"exam_5/jobsService/internal/infrastructure/grpc_service_clients"
	"exam_5/jobsService/internal/pkg/otlp"
	"exam_5/jobsService/internal/usecase"

	"go.opentelemetry.io/otel/attribute"

	"go.uber.org/zap"
)

const (
	serviceNameJobService = "jobServiceDelivery"
	spanNameJobService    = "jobSpanUsecaseDelivery"
)

type contentRPC struct {
	logger *zap.Logger
	job    usecase.Job
	client grpc_service_clients.ServiceClients
}

func NewRPC(logger *zap.Logger,  job usecase.Job, client grpc_service_clients.ServiceClients) pb.JobServiceServer {
	return &contentRPC{
		logger: logger,
		job:    job,
		client: client,
	}
}
func (s contentRPC) Create(ctx context.Context, in *pb.Job) (*pb.Job, error) {
	ctx, span := otlp.Start(ctx, serviceNameJobService, spanNameJobService+"Create")
	span.SetAttributes(
		attribute.Key("key").String(in.Name),
	)
	defer span.End()

	job := entity.Jobs{
		Id:        in.Id,
		Client_id: in.ClientId,
		Name:      in.Name,
		Comp_name: in.CompName,
		Status:    in.Status,
		Location:  in.Location,
		StartDate: in.StartDate,
		EndDate:   in.EndDate,
	}

	_, err := s.job.Create(ctx, &job)
	if err != nil {
		s.logger.Error("user.Create", zap.Error(err))
		return nil, grpc.Error(ctx, err)
	}

	resp := &pb.Job{
		Id :       job.Id,
		ClientId: job.Client_id,
		Name:      job.Name,
		CompName: job.Comp_name,
		Status:    job.Status,
		Location:  job.Location,
		StartDate: job.StartDate,
		EndDate:   job.EndDate,
		CreatedAt: job.Created_at.String(),
		UpdatedAt: job.Updated_at.String(),
	}
	return resp, nil
}
func (s contentRPC) Delete(ctx context.Context, in *pb.JobDeleteRequest) (*pb.JobDeleteResponse, error) {
	ctx, span := otlp.Start(ctx, serviceNameJobService, spanNameJobService+"delete")
	defer span.End()

	err := s.job.Delete(ctx, in.Id)

	if err != nil {
		s.logger.Error("user.Delete", zap.Error(err))
		return nil, grpc.Error(ctx, err)
	}

	return &pb.JobDeleteResponse{}, nil
}
func (s contentRPC) Get(ctx context.Context, in *pb.JobGetRequest) (*pb.Job, error) {
	ctx, span := otlp.Start(ctx, serviceNameJobService, spanNameJobService+"get")
	defer span.End()

	result := map[string]string{
		"id": in.Id,
	}

	job, err := s.job.Get(ctx, result)

	if err != nil {
		s.logger.Error("user.Get", zap.Error(err))
		return nil, grpc.Error(ctx, err)
	}

	return &pb.Job{
		Id :       job.Id,
		ClientId: job.Client_id,
		Name:      job.Name,
		CompName: job.Comp_name,
		Status:    job.Status,
		Location:  job.Location,
		StartDate: job.StartDate,
		EndDate:   job.EndDate,
		CreatedAt: job.Created_at.String(),
		UpdatedAt: job.Updated_at.String(),
	}, nil
}

func (s contentRPC) GetAll(ctx context.Context, in *pb.JobGetAllRequest) (*pb.JobGetAllResponse, error) {
	ctx, span := otlp.Start(ctx, serviceNameJobService, spanNameJobService+"list")
	defer span.End()

	offset := (in.Page - 1) * in.Limit

	jobs, err := s.job.List(ctx, uint64(in.Limit), uint64(offset), in.Filter)
	if err != nil {
		s.logger.Error("user.List", zap.Error(err))
		return nil, grpc.Error(ctx, err)
	}

	respJobs := []*pb.Job{}
	for _, job := range jobs {
		respUser := &pb.Job{
			Id :       job.Id,
			ClientId: job.Client_id,
			Name:      job.Name,
			CompName: job.Comp_name,
			Status:    job.Status,
			Location:  job.Location,
			StartDate: job.StartDate,
			EndDate:   job.EndDate,
			CreatedAt: job.Created_at.String(),
			UpdatedAt: job.Updated_at.String(),
		}

		respJobs = append(respJobs, respUser)
	}
	return &pb.JobGetAllResponse{
		Jobs: respJobs,
	}, nil
}
func (s contentRPC) Update(ctx context.Context, in *pb.Job) (*pb.Job, error) {
	ctx, span := otlp.Start(ctx, serviceNameJobService, spanNameJobService+"update")
	defer span.End()

	job, err := s.job.Update(ctx, &entity.Jobs{
		Id:        in.Id,
		Client_id: in.ClientId,
		Name:      in.ClientId,
		Comp_name: in.CompName,
		Status:    in.Status,
		Location:  in.Location,
		StartDate: in.StartDate,
		EndDate:   in.EndDate,
	})

	if err != nil {
		s.logger.Error("user.update", zap.Error(err))
		return nil, grpc.Error(ctx, err)
	}

	return &pb.Job{
		Id :       job.Id,
			ClientId: job.Client_id,
			Name:      job.Name,
			CompName: job.Comp_name,
			Status:    job.Status,
			Location:  job.Location,
			StartDate: job.StartDate,
			EndDate:   job.EndDate,
			CreatedAt: job.Created_at.String(),
			UpdatedAt: job.Updated_at.String(),
	}, nil
}

func (s contentRPC)GetAllJobWithOwner(ctx context.Context, in *pb.JobGetAllRequest)(*pb.GetAllJobByClientIdResponse, error){
	ctx, span := otlp.Start(ctx, serviceNameJobService, spanNameJobService+"list")
	defer span.End()

	offset := (in.Page - 1) * in.Limit

	jobs, err := s.job.List(ctx, uint64(in.Limit), uint64(offset), in.Filter)
	if err != nil {
		s.logger.Error("user.List", zap.Error(err))
		return nil, grpc.Error(ctx, err)
	}

	respJobs := []*pb.JobWithOwner{}
	for _, job := range jobs {
		client, err := s.client.ClientService().Get(ctx, &client.GetRequest{
			Field: "id",
			Useroremail: job.Client_id,
		})
		if err != nil{
			s.logger.Error("user.List", zap.Error(err))
		return nil, grpc.Error(ctx, err)
		}
		owner := pb.Owner{
			Id: client.Id,
			Name: client.Name,
			LastName: client.LastName,
			Email: client.Email,
			Password: client.Password,
			Role: client.Role,
			RefreshToken: client.RefreshToken,
			UpdatedAt: client.UpdatedAt,
			CreatedAt: client.CreatedAt,
		}
		respUser := &pb.JobWithOwner{
			Id :       job.Id,
			ClientId: job.Client_id,
			Name:      job.Name,
			CompName: job.Comp_name,
			Status:    job.Status,
			Location:  job.Location,
			StartDate: job.StartDate,
			EndDate:   job.EndDate,
			CreatedAt: job.Created_at.String(),
			UpdatedAt: job.Updated_at.String(),
			Owner: &owner,
		}

		respJobs = append(respJobs, respUser)
	}
	return &pb.GetAllJobByClientIdResponse{
		Jobs: respJobs,
	}, nil
}