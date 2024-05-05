package postgres

import (
	"context"
	"exam_5/jobsService/internal/entity"

	"exam_5/jobsService/internal/pkg/otlp"
	postgres "exam_5/jobsService/internal/pkg/postgres"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/k0kubun/pp"
)

const (
	jobServiceTableName   = "jobs"
	serviceNameJobService = "JobServiceRepo"
	spanNameJobService    = "JobSpanRepo"
)

type clientRepo struct {
	tableName string
	db        *postgres.PostgresDB
}

// NewUserRepo ...
func NewJobRepo(db *postgres.PostgresDB) *clientRepo {
	return &clientRepo{
		tableName: jobServiceTableName,
		db:        db,
	}
}


func (p *clientRepo) userSelectQueryPrefix() squirrel.SelectBuilder {
	return p.db.Sq.Builder.Select(
		"id",
		"client_id",
		"name",
		"company_name",
		"location",
		"start_date",
		"end_date",
		"status",
		"created_at",
		"updated_at",
	).From(p.tableName)
}

func (p clientRepo) Create(ctx context.Context, job *entity.Jobs) (*entity.Jobs, error) {
	ctx, span := otlp.Start(ctx, serviceNameJobService, spanNameJobService+"Create")
	defer span.End()
	data := map[string]any{
		"id": job.Id,
		"client_id": job.Client_id,
		"name": job.Name,
		"company_name" : job.Comp_name,
		"location": job.Location,
		"start_date": job.StartDate,
		"end_date": job.EndDate,
		"status": job.Status,
		"created_at": job.Created_at,
		"updated_at":job.Updated_at,
	}
	query, args, err := p.db.Sq.Builder.Insert(p.tableName).SetMap(data).ToSql()
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "create"))
	}

	_, err = p.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, p.db.Error(err)
	}

	return job, nil
}



func (p clientRepo) Get(ctx context.Context, params map[string]string) (*entity.Jobs, error) {
	ctx, span := otlp.Start(ctx, serviceNameJobService, spanNameJobService+"get")
	defer span.End()
	var (
		job entity.Jobs
	)

	queryBuilder := p.userSelectQueryPrefix()

	for key, value := range params {
		if key == "id" {
			queryBuilder = queryBuilder.Where(p.db.Sq.Equal(key, value)).Where("deleted_at is null")
		}
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "get"))
	}

	if err = p.db.QueryRow(ctx, query, args...).Scan(
		&job.Id,
		&job.Client_id,
		&job.Name,
		&job.Comp_name,
		&job.Location,
		&job.StartDate,
		&job.EndDate,
		&job.Status,
		&job.Created_at,
		&job.Updated_at,
	); err != nil {
		return nil, p.db.Error(err)
	}

	return &job, nil
}
func (p clientRepo) Update(ctx context.Context, job *entity.Jobs) (*entity.Jobs, error) {
	ctx, span := otlp.Start(ctx, serviceNameJobService, spanNameJobService+"Update")
	defer span.End()
	clauses := map[string]any{
		"name": job.Name,
		"company_name" : job.Comp_name,
		"location": job.Location,
		"start_date": job.StartDate,
		"end_date": job.EndDate,
		"status": job.Status,
		"updated_at": job.Updated_at,
	}
	sqlStr, args, err := p.db.Sq.Builder.
		Update(p.tableName).
		SetMap(clauses).
		Where(p.db.Sq.Equal("id", job.Id)).
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
	params["id"] = job.Id
	respjob, err := p.Get(ctx, params)
	if err != nil {
		return nil, err
	}
	return respjob, nil
}

func (p clientRepo) GetAll(ctx context.Context, limit, offset uint64, filter map[string]string) ([]*entity.Jobs, error) {
	ctx, span := otlp.Start(ctx, serviceNameJobService, spanNameJobService+"List")
	defer span.End()
	var (
		jobs []*entity.Jobs
	)
	queryBuilder := p.userSelectQueryPrefix()

	if limit != 0 {
		queryBuilder = queryBuilder.Limit(limit).Offset(offset)
	}

	for key, value := range filter {
		if key == "client_id" {
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
    pp.Println(query)
	rows, err := p.db.Query(ctx, query, args...)
	if err != nil {
		return nil, p.db.Error(err)
	}
	defer rows.Close()
	jobs = make([]*entity.Jobs, 0)
	for rows.Next() {
		var job entity.Jobs
		if err = rows.Scan(
			&job.Id,
			&job.Client_id,
			&job.Name,
			&job.Comp_name,
			&job.Location,
			&job.StartDate,
			&job.EndDate,
			&job.Status,
			&job.Created_at,
			&job.Updated_at,
		); err != nil {
			return nil, p.db.Error(err)
		}
		jobs = append(jobs, &job)
	}

	return jobs, nil
}

func (p clientRepo) Delete(ctx context.Context, guid string, filter map[string]any) error {
	ctx, span := otlp.Start(ctx, serviceNameJobService, spanNameJobService+"Delete")
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
// func (p clientRepo) GetAllByClientId(ctx context.Context, limit, offset uint64, filter map[string]string) ([]*entity.Jobs, error) {
// 	ctx, span := otlp.Start(ctx, serviceNameJobService, spanNameJobService+"List")
// 	defer span.End()
// 	var (
// 		jobs []*entity.Jobs
// 	)
// 	queryBuilder := p.userSelectQueryPrefix()

// 	if limit != 0 {
// 		queryBuilder = queryBuilder.Limit(limit).Offset(offset)
// 	}

// 	for key, value := range filter {
// 		if key == "client_id" {
// 			queryBuilder = queryBuilder.Where(p.db.Sq.Equal(key, value))
// 			continue
// 		}
// 		if key == "created_at" {
// 			queryBuilder = queryBuilder.Where("created_at=?", value)
// 			continue
// 		}
// 	}

// 	query, args, err := queryBuilder.ToSql()
// 	if err != nil {
// 		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "list"))
// 	}

// 	rows, err := p.db.Query(ctx, query, args...)
// 	if err != nil {
// 		return nil, p.db.Error(err)
// 	}
// 	defer rows.Close()
// 	jobs = make([]*entity.Jobs, 0)
// 	for rows.Next() {
// 		var job entity.Jobs
// 		if err = rows.Scan(
// 			&job.Id,
// 			&job.Client_id,
// 			&job.Name,
// 			&job.Comp_name,
// 			&job.Location,
// 			&job.StartDate,
// 			&job.EndDate,
// 			&job.Status,
// 			&job.Created_at,
// 			&job.Updated_at,
// 		); err != nil {
// 			return nil, p.db.Error(err)
// 		}
// 		jobs = append(jobs, &job)
// 	}

// 	return jobs, nil
// }

