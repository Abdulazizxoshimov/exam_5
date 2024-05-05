package usecase

import (
	"exam_5/clientService/internal/entity"
	repository "exam_5/clientService/internal/infrastructure/repository"
	"log"

	"exam_5/clientService/internal/pkg/otlp"
	"context"
	"time"

	"github.com/k0kubun/pp"
)

const (
	serviceNameUserService = "userServiceUsecase"
	spanNameUserService    = "userSpanUsecase"
)

type User interface {
	Create(ctx context.Context, document *entity.Client) (*entity.Client, error)
	Get(ctx context.Context, params map[string]string) (*entity.Client, error)
	List(ctx context.Context, limit, offset uint64, filter map[string]string) ([]*entity.Client, error)
	Update(ctx context.Context, document *entity.Client) (*entity.Client, error)
	Delete(ctx context.Context, guid string) error
	IsUnique(ctx context.Context, params map[string]string)(bool, error)
    UpdateUserRefreshtoken(ctx context.Context, id, token string)error
}

type userService struct {
	BaseUseCase
	repo       repository.UserStorageI
	ctxTimeout time.Duration
}

func NewUserService(ctxTimeout time.Duration, repo repository.UserStorageI) userService {
	return userService{
		ctxTimeout: ctxTimeout,
		repo:       repo,
	}
}

func (u userService) Create(ctx context.Context, user *entity.Client) (*entity.Client, error) {
	

	ctx, span := otlp.Start(ctx, serviceNameUserService, spanNameUserService+"Create")
	defer span.End()

	u.beforeRequest(nil, &user.CreatedAt, &user.UpdatedAt)
	respUser, err := u.repo.Create(ctx, user)
	if err != nil{
		pp.Println("error while create user usecase layer :", err )
		return nil, err
	}
	return respUser, err
}


func (u userService) Get(ctx context.Context, params map[string]string) (*entity.Client, error) {
	
	ctx, span := otlp.Start(ctx, serviceNameUserService, spanNameUserService+"Get")
	defer span.End()
    
	respUser, err :=  u.repo.Get(ctx, params)
	if err != nil{
		log.Println("error while get user usercase layer")
		return nil, err
	}
	return respUser, nil
}


func (u userService) List(ctx context.Context, limit, offset uint64, filter map[string]string) ([]*entity.Client, error) {
	
	ctx, span := otlp.Start(ctx, serviceNameUserService, spanNameUserService+"List")
	defer span.End()
	
	return u.repo.GetAll(ctx, limit, offset, filter)
}


func (u userService) Update(ctx context.Context, user *entity.Client) (*entity.Client, error) {
	
	ctx, span := otlp.Start(ctx, serviceNameUserService, spanNameUserService+"Update")
	defer span.End()
	
	u.beforeRequest(nil, nil, &user.UpdatedAt)

	return u.repo.Update(ctx, user)
}



func (u userService) Delete(ctx context.Context, guid string) error {
	
	ctx, span := otlp.Start(ctx, serviceNameUserService, spanNameUserService+"Delete")
	defer span.End()
	
    filter := make(map[string]any)
    filter["deleted_at"] = time.Now().UTC()
	return u.repo.Delete(ctx, guid, filter)
}

func (u userService)IsUnique(ctx context.Context, params map[string]string)(bool, error){
	
	ctx, span := otlp.Start(ctx, serviceNameUserService, spanNameUserService+"Delete")
	defer span.End()
	

	status, err := u.repo.IsUnique(ctx, params)
    if err != nil{
		return false, err
	}
	return status, nil
}
func (u userService)UpdateUserRefreshtoken(ctx context.Context, id, token string)error{
	ctx, span := otlp.Start(ctx, serviceNameUserService, spanNameUserService+"updaterefresh token")
	defer span.End()

	err := u.repo.UpdateUserRefreshtoken(ctx, id, token)
	if err != nil{
		return err
	}
	return nil
}