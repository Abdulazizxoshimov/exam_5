package repository




import (
	"context"
	"exam_5/clientService/internal/entity"
)


type UserStorageI interface {
Create(ctx context.Context, client *entity.Client) (*entity.Client, error)
Get(ctx context.Context, params map[string]string) (*entity.Client, error)
Update(ctx context.Context, client *entity.Client) (*entity.Client, error)
GetAll(ctx context.Context, limit, offset uint64, filter map[string]string) ([]*entity.Client, error)
Delete(ctx context.Context, guid string, filter map[string]any) error 
IsUnique(ctx context.Context, params map[string]string)(bool, error)
UpdateUserRefreshtoken(ctx context.Context, id, token string)error
}