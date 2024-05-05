package mock_data

import (
	"context"
	pbj "exam_5/admin_api_gateway/genproto/jobProto"
)

type MockServiceJob interface {
	CreateJob(ctx context.Context, kyc *pbj.Job) (*pbj.Job, error)
	GetJob(ctx context.Context, params map[string]string) (*pbj.Job, error)
	ListJob(ctx context.Context, limit, offset uint64, filter map[string]string) ([]*pbj.Job, error)
	UpdateJob(ctx context.Context, kyc *pbj.Job) (*pbj.Job, error)
	DeleteJob(ctx context.Context, id string) (*pbj.JobDeleteResponse, error)
}

type mockServiceJob struct {
}

func NewMockServiceJob() MockServiceJob {
	return &mockServiceJob{}
}

func (c *mockServiceJob) CreateJob(ctx context.Context, kyc *pbj.Job) (*pbj.Job, error) {
	return kyc, nil
}

func (c *mockServiceJob) UpdateJob(ctx context.Context, kyc *pbj.Job) (*pbj.Job, error) {
	return &pbj.Job{
		Id:        "update job",
		Name:      "update job",
		CompName:  "update job",
		StartDate: "2006-02-01",
		EndDate:   "2006-02-01",
		Status:    true,
		ClientId:  "dewrty-erf43-sb03-bdsf-bgbsdf",
		CreatedAt: "update time job",
		UpdatedAt: "update time job",
	}, nil
}

func (c *mockServiceJob) DeleteJob(ctx context.Context, id string) (*pbj.JobDeleteResponse, error) {
	return &pbj.JobDeleteResponse{}, nil
}

func (c *mockServiceJob) GetJob(ctx context.Context, params map[string]string) (*pbj.Job, error) {
	return &pbj.Job{
		Id:        "mock id",
		Name:      "mock name",
		CompName:  "mock Job",
		StartDate: "2006-02-01",
		EndDate:   "2006-02-01",
		Status:    true,
		ClientId:  "ca-daskj3b-akj3u-ASJKC3-jjjd2i1n",
		CreatedAt: "mock time",
		UpdatedAt: "mock time",
	}, nil
}

func (c *mockServiceJob) ListJob(ctx context.Context, limit, offset uint64, filter map[string]string) ([]*pbj.Job, error) {
	return []*pbj.Job{
		{
			Id:        "mock id",
			Name:      "mock name",
			CompName:  "mock Job",
			StartDate: "2006-02-01",
			EndDate:   "2006-02-01",
			Status:    true,
			ClientId:  "ca-daskj3b-akj3u-ASJKC3-jjjd2i1n",
			CreatedAt: "mock time",
			UpdatedAt: "mock time",
		},
		{
			Id:        "mock id",
			Name:      "mock name",
			CompName:  "mock Job",
			StartDate: "2006-02-01",
			EndDate:   "2006-02-01",
			Status:    true,
			ClientId:  "ca-daskj3b-akj3u-ASJKC3-jjjd2i1n",
			CreatedAt: "mock time",
			UpdatedAt: "mock time",
		},
	}, nil

}
