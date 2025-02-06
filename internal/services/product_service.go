package services

import (
	"context"

	"github.com/andresilvase/gobid/internal/store/pgstore"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductService struct {
	pool    *pgxpool.Pool
	queries *pgstore.Queries
}

func NewProductService(pool *pgxpool.Pool) *ProductService {
	return &ProductService{
		pool:    pool,
		queries: pgstore.New(pool),
	}
}

func (ps *ProductService) CreateProduct(ctx context.Context, createProductParams pgstore.CreateProductParams) (uuid.UUID, error) {
	return ps.queries.CreateProduct(ctx, createProductParams)
}
