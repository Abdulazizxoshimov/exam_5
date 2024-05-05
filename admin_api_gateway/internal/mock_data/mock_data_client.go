package mock_data

import (
	"context"
	pbc "exam_5/admin_api_gateway/genproto/clientProto"
)

type MockServiceClient interface {
	CreateClient(ctx context.Context, kyc *pbc.User) (*pbc.User, error)
	GetClient(ctx context.Context, params map[string]string) (*pbc.User, error)
	ListClient(ctx context.Context, limit, offset uint64, filter map[string]string) ([]*pbc.User, error)
	UpdateClient(ctx context.Context, kyc *pbc.User) (*pbc.User, error)
	DeleteClient(ctx context.Context, id string) (*pbc.DeleteResponse, error)
}

type mockServiceClient struct {
}

func NewMockServiceClient() MockServiceClient {
	return &mockServiceClient{}
}

func (c *mockServiceClient) CreateClient(ctx context.Context, in *pbc.User) (*pbc.User, error) {
	return in, nil
}

func (c *mockServiceClient) UpdateClient(ctx context.Context, kyc *pbc.User) (*pbc.User, error) {
	return &pbc.User{
		Id:           "mock update",
		Name:         "mock update",
		LastName:     "mock update",
		Email:        "mock update",
		Role:         "mock update",
		Password:     "mock update",
		RefreshToken: "mock update",
		CreatedAt:    "mock update",
		UpdatedAt:    "mock update",
		DeletedAt:    "mock update",
	}, nil
}

func (c *mockServiceClient) DeleteClient(ctx context.Context, in string) (*pbc.DeleteResponse, error) {
	return &pbc.DeleteResponse{

	}, nil
}

func (c *mockServiceClient) GetClient(ctx context.Context, params map[string]string) (*pbc.User, error) {
	return &pbc.User{
		Id:           "mock id",
		Name:         "mock name",
		LastName:     "mock last_name",
		Email:        "mock email",
		Password:     "mock password",
		Role:          "mock Role",
		RefreshToken: "mock refresh_token",
		CreatedAt:    "mock time",
		UpdatedAt:    "mock time",
		DeletedAt:    "mock time",
	}, nil
}

func (c *mockServiceClient) ListClient(ctx context.Context, limit, offset uint64, filter map[string]string) ([]*pbc.User, error) {
	return []*pbc.User{
		{
			Id:           "mock 1",
			Name:         "mock 1",
			LastName:     "mock 1",
			Email:        "mock 1",
			Password:     "mock 1",
			RefreshToken: "mock 1",
			CreatedAt:    "mock 1",
			UpdatedAt:    "mock 1",
			DeletedAt:    "mock 1",
			Role:          "mock 1",
		},
		{
			Id:           "mock 2",
			Name:         "mock 2",
			LastName:     "mock 2",
			Email:        "mock 2",
			Role:          "mock 2",
			Password:     "mock 2",
			RefreshToken: "mock 2",
			CreatedAt:    "mock 2",
			UpdatedAt:    "mock 2",
			DeletedAt:    "mock 2",
		},
	}, nil

}
