package services

import (
	"context"
	pb "exam_5/clientService/genproto/clientProto"
	"exam_5/clientService/internal/delivery/grpc"
	"exam_5/clientService/internal/entity"
	"exam_5/clientService/internal/pkg/otlp"
	"exam_5/clientService/internal/usecase"

	"go.opentelemetry.io/otel/attribute"

	"go.uber.org/zap"
)

const (
	serviceNameUserService = "userServiceDelivery"
	spanNameUserService    = "userSpanUsecaseDelivery"
)

type contentRPC struct {
	logger *zap.Logger
	user   usecase.User
}

func NewRPC(logger *zap.Logger,
	user usecase.User,
) pb.UserServiceServer {
	return &contentRPC{
		logger: logger,
		user:   user,
	}
}
func (s contentRPC) Create(ctx context.Context, in *pb.User) (*pb.User, error) {
	ctx, span := otlp.Start(ctx, serviceNameUserService, spanNameUserService+"Create")
	span.SetAttributes(
		attribute.Key("key").String(in.Name),
	)
	defer span.End()

	var user entity.Client
	user.Id  = in.Id
	user.Name = in.Name
	user.LastName = in.LastName
	user.Email = in.Email
	user.Password = in.Password
	user.RefreshToken = in.RefreshToken
	user.Role = in.Role

	_, err := s.user.Create(ctx, &user)
	if err != nil {
		s.logger.Error("user.Create", zap.Error(err))
		return nil, grpc.Error(ctx, err)
	}

	resp := &pb.User{
		Id:           user.Id,
		Name:         user.Name,
		LastName:     user.LastName,
		Email:        user.Email,
		Password:     user.Password,
		RefreshToken: user.RefreshToken,
		Role:         user.Role,
		CreatedAt:    user.CreatedAt.String(),
		UpdatedAt:    user.UpdatedAt.String(),
	}
	return resp, nil
}
func (s contentRPC) Delete(ctx context.Context, in *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	ctx, span := otlp.Start(ctx, serviceNameUserService, spanNameUserService+"delete")
	defer span.End()

	err := s.user.Delete(ctx, in.Value)

	if err != nil {
		s.logger.Error("user.Delete", zap.Error(err))
		return nil, grpc.Error(ctx, err)
	}

	return &pb.DeleteResponse{}, nil
}
func (s contentRPC) Get(ctx context.Context, in *pb.GetRequest) (*pb.User, error) {
	ctx, span := otlp.Start(ctx, serviceNameUserService, spanNameUserService+"get")
	defer span.End()

	result := map[string]string{
		in.Field: in.Useroremail,
	}

	user, err := s.user.Get(ctx, result)

	if err != nil {
		s.logger.Error("user.Get", zap.Error(err))
		return nil, grpc.Error(ctx, err)
	}

	return &pb.User{
		Id:           user.Id,
		Name:         user.Name,
		LastName:     user.LastName,
		Email:        user.Email,
		Password:     user.Password,
		RefreshToken: user.RefreshToken,
		Role:         user.Role,
		CreatedAt:    user.CreatedAt.String(),
		UpdatedAt:    user.UpdatedAt.String(),
	}, nil
}
func (s contentRPC) List(ctx context.Context, in *pb.GetAllUsersRequest) (*pb.GetAllUsersResponse, error) {
	ctx, span := otlp.Start(ctx, serviceNameUserService, spanNameUserService+"list")
	defer span.End()

	offset := (in.Page - 1) * in.Limit

	users, err := s.user.List(ctx, uint64(in.Limit), uint64(offset), map[string]string{})
	if err != nil {
		s.logger.Error("user.List", zap.Error(err))
		return nil, grpc.Error(ctx, err)
	}

	respUsers := []*pb.User{}
	for _, user := range users {
		respUser := &pb.User{
			Id:           user.Id,
			Name:         user.Name,
			LastName:     user.LastName,
			Email:        user.Email,
			Password:     user.Password,
			RefreshToken: user.RefreshToken,
			Role:         user.Role,
			CreatedAt:    user.CreatedAt.String(),
			UpdatedAt:    user.UpdatedAt.String(),
		}

		respUsers = append(respUsers, respUser)
	}
	return &pb.GetAllUsersResponse{
		Users: respUsers,
	}, nil
}
func (s contentRPC) Update(ctx context.Context, in *pb.User) (*pb.User, error) {
	ctx, span := otlp.Start(ctx, serviceNameUserService, spanNameUserService+"update")
	defer span.End()

	respUser, err := s.user.Update(ctx, &entity.Client{
		Id:       in.Id,
		Email:    in.Email,
		Name:     in.Name,
		LastName: in.LastName,
		Password: in.Password,
		Role:     in.Role,
	})

	if err != nil {
		s.logger.Error("user.update", zap.Error(err))
		return nil, grpc.Error(ctx, err)
	}

	return &pb.User{
		Id:           respUser.Id,
		Name:         respUser.Name,
		LastName:     respUser.LastName,
		Email:        respUser.Email,
		Password:     respUser.Password,
		Role:         respUser.Role,
		RefreshToken: respUser.RefreshToken,
		CreatedAt:    respUser.CreatedAt.String(),
		UpdatedAt:    respUser.UpdatedAt.GoString(),
	}, nil
}
func (s contentRPC) IsUnique(ctx context.Context, in *pb.IsUniqueRequest) (*pb.IsUniqueResponse, error) {
	ctx, span := otlp.Start(ctx, serviceNameUserService, spanNameUserService+"unique")
	defer span.End()
	params := map[string]string{
		in.Field: in.Value,
	}
	statrus, err := s.user.IsUnique(ctx, params)
	if err != nil {
		s.logger.Error("is.unique", zap.Error(err))
		return nil, grpc.Error(ctx, err)
	}
	return &pb.IsUniqueResponse{
		IsUnique: statrus,
	}, nil
}
func (s contentRPC) UpdateUserRefreshToken(ctx context.Context, in *pb.UpdateRefreshToken) (*pb.UpdateRefreshTokenResponse, error) {
	err := s.user.UpdateUserRefreshtoken(ctx, in.Id, in.RefreshToken)
	if err != nil {
		s.logger.Error("update.refresh.token", zap.Error(err))
		return nil, grpc.Error(ctx, err)
	}

	return &pb.UpdateRefreshTokenResponse{}, nil
}
