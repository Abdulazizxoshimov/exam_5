package repository




import (
	"context"
	"exam_5/jobsService/internal/entity"
)


type JobStorageI interface {
Create(ctx context.Context, client *entity.Jobs) (*entity.Jobs, error)
Get(ctx context.Context, params map[string]string) (*entity.Jobs, error)
GetAll(ctx context.Context, limit, offset uint64, filter map[string]string) ([]*entity.Jobs, error)
Delete(ctx context.Context, guid string, filter map[string]any) error
Update(ctx context.Context, job *entity.Jobs) (*entity.Jobs, error)
}