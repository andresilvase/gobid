package services

import (
	"context"
	"errors"

	"github.com/andresilvase/gobid/internal/store/pgstore"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductService struct {
	pool    *pgxpool.Pool
	queries *pgstore.Queries
}

func NewProductService(pool *pgxpool.Pool) ProductService {
	return ProductService{
		pool:    pool,
		queries: pgstore.New(pool),
	}
}

var ErrProductNotFound = errors.New("product not found")

func (ps *ProductService) CreateProduct(ctx context.Context, createProductParams pgstore.CreateProductParams) (uuid.UUID, error) {
	return ps.queries.CreateProduct(ctx, createProductParams)
}

func (ps *ProductService) GetProductById(ctx context.Context, productId uuid.UUID) (pgstore.Product, error) {
	product, err := ps.queries.GetProductById(ctx, productId)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgstore.Product{}, ErrProductNotFound
		}
		return pgstore.Product{}, err
	}

	return product, nil
}
