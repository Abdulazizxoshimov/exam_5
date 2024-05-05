package postgres

import (
	"context"
	"exam_5/clientService/internal/entity"
	"exam_5/clientService/internal/pkg/otlp"
	postgres "exam_5/clientService/internal/pkg/postgres"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/k0kubun/pp"
)

const (
	clientServiceTableName   = "clients"
	serviceNameClientService = "userServiceRepo"
	spanNameClientService    = "userSpanRepo"
)

type clientRepo struct {
	tableName string
	db        *postgres.PostgresDB
}

// NewUserRepo ...
func NewUserRepo(db *postgres.PostgresDB) *clientRepo {
	return &clientRepo{
		tableName: clientServiceTableName,
		db:        db,
	}
}

func (p *clientRepo) userSelectQueryPrefix() squirrel.SelectBuilder {
	return p.db.Sq.Builder.Select(
		"id",
		"name",
		"last_name",
		"email",
		"refresh_token",
		"password",
		"created_at",
		"updated_at",
		"role",
	).From(p.tableName)
}

func (p clientRepo) Create(ctx context.Context, client *entity.Client) (*entity.Client, error) {
	ctx, span := otlp.Start(ctx, serviceNameClientService, spanNameClientService+"Create")
	defer span.End()
	data := map[string]any{
		"id":            client.Id,
		"name":          client.Name,
		"last_name":     client.LastName,
		"email":         client.Email,
		"refresh_token": client.RefreshToken,
		"password":      client.Password,
		"created_at":    client.CreatedAt,
		"updated_at":    client.UpdatedAt,
		"role":          client.Role,
	}
	query, args, err := p.db.Sq.Builder.Insert(p.tableName).SetMap(data).ToSql()
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "create"))
	}

	_, err = p.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, p.db.Error(err)
	}

	return client, nil
}

func (r *clientRepo) IsUnique(ctx context.Context, params map[string]string) (bool, error) {
	queryBuilder := r.db.Sq.Builder.Select("COUNT(1)").
		From("clients").
		Where(squirrel.Eq{"email": params["email"]}).
		Where("deleted_at is null")

	query, args, err := queryBuilder.ToSql()

	if err != nil {
		return false, r.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", clientServiceTableName, "isUnique"))
	}

	var count int

	if err = r.db.QueryRow(ctx, query, args...).Scan(&count); err != nil {
		return false, r.db.Error(err)

	}
	if count != 0 {
		return true, nil
	}
	return false, nil
}

func (p clientRepo) Get(ctx context.Context, params map[string]string) (*entity.Client, error) {
	ctx, span := otlp.Start(ctx, serviceNameClientService, spanNameClientService+"get")
	defer span.End()
	var (
		client entity.Client
	)

	queryBuilder := p.userSelectQueryPrefix()

	for key, value := range params {
		queryBuilder = queryBuilder.Where(p.db.Sq.Equal(key, value)).Where("deleted_at is null")
		
	}
    
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "get"))
	}
	
	if err = p.db.QueryRow(ctx, query, args...).Scan(
		&client.Id,
		&client.Name,
		&client.LastName,
		&client.Email,
		&client.RefreshToken,
		&client.Password,
		&client.CreatedAt,
		&client.UpdatedAt,
		&client.Role,
	); err != nil {
		return nil, p.db.Error(err)
	}
    pp.Println(client)
	return &client, nil
}


func (p clientRepo) Update(ctx context.Context, client *entity.Client) (*entity.Client, error) {
	ctx, span := otlp.Start(ctx, serviceNameClientService, spanNameClientService+"Update")
	defer span.End()
	clauses := map[string]any{
		"name":       client.Name,
		"last_name":  client.LastName,
		"email":      client.Email,
		"password":   client.Password,
		"updated_at": client.UpdatedAt,
	}
	sqlStr, args, err := p.db.Sq.Builder.
		Update(p.tableName).
		SetMap(clauses).
		Where(p.db.Sq.Equal("id", client.Id)).
		ToSql()
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, p.tableName+" Update")
	}

	commandTag, err := p.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return nil, p.db.Error(err)
	}

	if commandTag.RowsAffected() == 0 {
		return nil, p.db.Error(fmt.Errorf("no sql rows"))
	}
	params := make(map[string]string)
	params["id"] = client.Id
	respUser, err := p.Get(ctx, params)
	if err != nil {
		return nil, err
	}
	return respUser, nil
}

func (p clientRepo) GetAll(ctx context.Context, limit, offset uint64, filter map[string]string) ([]*entity.Client, error) {
	ctx, span := otlp.Start(ctx, serviceNameClientService, spanNameClientService+"List")
	defer span.End()
	var (
		users []*entity.Client
	)
	queryBuilder := p.userSelectQueryPrefix()

	if limit != 0 {
		queryBuilder = queryBuilder.Limit(limit).Offset(offset)
	}

	for key, value := range filter {
		if key == "id" {
			queryBuilder = queryBuilder.Where(p.db.Sq.Equal(key, value))
			continue
		}
		if key == "created_at" {
			queryBuilder = queryBuilder.Where("created_at=?", value)
			continue
		}
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "list"))
	}

	rows, err := p.db.Query(ctx, query, args...)
	if err != nil {
		return nil, p.db.Error(err)
	}
	defer rows.Close()
	users = make([]*entity.Client, 0)
	for rows.Next() {
		var client entity.Client
		if err = rows.Scan(
			&client.Id,
			&client.Name,
			&client.LastName,
			&client.Email,
			&client.RefreshToken,
			&client.Password,
			&client.CreatedAt,
			&client.UpdatedAt,
			&client.Role,
		); err != nil {
			return nil, p.db.Error(err)
		}
		users = append(users, &client)
	}

	return users, nil
}

func (p clientRepo) Delete(ctx context.Context, guid string, filter map[string]any) error {
	ctx, span := otlp.Start(ctx, serviceNameClientService, spanNameClientService+"Delete")
	defer span.End()

	sqlStr, args, err := p.db.Sq.Builder.
		Update(p.tableName).
		SetMap(filter).
		Where(p.db.Sq.Equal("id", guid)).
		Where("deleted_at is null").
		ToSql()
	if err != nil {
		return p.db.ErrSQLBuild(err, p.tableName+" delete")
	}

	commandTag, err := p.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return p.db.Error(err)
	}

	if commandTag.RowsAffected() == 0 {
		return p.db.Error(fmt.Errorf("no sql rows"))
	}

	return nil
}

func (r * clientRepo) UpdateUserRefreshtoken(ctx context.Context, id, token string)(err error){
    ctx, span := otlp.Start(ctx, serviceNameClientService, spanNameClientService+"Update")
	defer span.End()
	clauses := map[string]any{
		"refresh_token": token,
	}
	sqlStr, args, err := r.db.Sq.Builder.
		Update(r.tableName).
		SetMap(clauses).
		Where(r.db.Sq.Equal("id", id)).
		ToSql()
	if err != nil {
		return  r.db.ErrSQLBuild(err, r.tableName+" Update")
	}

	commandTag, err := r.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return  r.db.Error(err)
	}

	if commandTag.RowsAffected() == 0 {
		return  r.db.Error(fmt.Errorf("no sql rows"))
	}
	
	return nil
}
