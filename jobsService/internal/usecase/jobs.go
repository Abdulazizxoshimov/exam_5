package usecase

import (
	"exam_5/jobsService/internal/entity"
	repository "exam_5/jobsService/internal/infrastructure/repository"
	"log"

	"exam_5/jobsService/internal/pkg/otlp"
	"context"
	"time"

	"github.com/k0kubun/pp"
)

const (
	serviceNameUserService = "jobServiceUsecase"
	spanNameUserService    = "jobSpanUsecase"
)

type Job interface {
	Create(ctx context.Context, document *entity.Jobs) (*entity.Jobs, error)
	Get(ctx context.Context, params map[string]string) (*entity.Jobs, error)
	List(ctx context.Context, limit, offset uint64, filter map[string]string) ([]*entity.Jobs, error)
	Update(ctx context.Context, document *entity.Jobs) (*entity.Jobs, error)
	Delete(ctx context.Context, guid string) error
	}

type jobService struct {
	BaseUseCase
	repo       repository.JobStorageI
	ctxTimeout time.Duration
}

func NewJobService(ctxTimeout time.Duration, repo repository.JobStorageI) jobService {
	return jobService{
		ctxTimeout: ctxTimeout,
		repo:       repo,
	}
}

func (u jobService) Create(ctx context.Context, job *entity.Jobs) (*entity.Jobs, error) {
	

	ctx, span := otlp.Start(ctx, serviceNameUserService, spanNameUserService+"Create")
	defer span.End()

	u.beforeRequest(&job.Id, &job.Created_at, &job.Updated_at)
	respJob, err := u.repo.Create(ctx, job)
	if err != nil{
		pp.Println("error while create user usecase layer :", err )
		return nil, err
	}
	return respJob, err
}


func (u jobService) Get(ctx context.Context, params map[string]string) (*entity.Jobs, error) {
	
	ctx, span := otlp.Start(ctx, serviceNameUserService, spanNameUserService+"Get")
	defer span.End()
    
	respJob, err :=  u.repo.Get(ctx, params)
	if err != nil{
		log.Println("error while get user usercase layer")
		return nil, err
	}
	return respJob, nil
}


func (u jobService) List(ctx context.Context, limit, offset uint64, filter map[string]string) ([]*entity.Jobs, error) {
	
	ctx, span := otlp.Start(ctx, serviceNameUserService, spanNameUserService+"List")
	defer span.End()
	
	return u.repo.GetAll(ctx, limit, offset, filter)
}


func (u jobService) Update(ctx context.Context, job *entity.Jobs) (*entity.Jobs, error) {
	
	ctx, span := otlp.Start(ctx, serviceNameUserService, spanNameUserService+"Update")
	defer span.End()
	
	u.beforeRequest(nil, nil, &job.Updated_at)

	return u.repo.Update(ctx, job)
}



func (u jobService) Delete(ctx context.Context, guid string) error {
	
	ctx, span := otlp.Start(ctx, serviceNameUserService, spanNameUserService+"Delete")
	defer span.End()
	
    filter := make(map[string]any)
    filter["deleted_at"] = time.Now().UTC()
	return u.repo.Delete(ctx, guid, filter)
}

