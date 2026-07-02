package products

import (
	"context"
	repo "gary/ecom/internal/adapters/postgres/sqlc"
)

type Service interface {
	ListProducts(ctx context.Context) ([]repo.Product, error)
	ListProductById(ctx context.Context, id int64) (repo.Product, error)
}

type svc struct {
	//repository
	repo repo.Querier
}

func NewService(repo repo.Querier) Service {
	return &svc{
		repo: repo,
	}
}

func (s *svc) ListProducts(ctx context.Context) ([]repo.Product, error) {
	return s.repo.ListProducts(ctx)
}

func (s *svc) ListProductById(ctx context.Context, id int64) (repo.Product, error) {
	return s.repo.ListProductById(ctx, id)
}
