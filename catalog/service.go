package catalog

import (
	"context"

	"github.com/segmentio/ksuid"
)

type Serice interface {
	CreateProduct(ctx context.Context, name, description string, price float64) (*Product, error)
	GetProductById(ctx context.Context, id string) (*Product, error)
	GetListProducts(ctx context.Context, skip, take uint64) ([]Product, error)
	GetListProductsByIDs(ctx context.Context, ids []string) ([]Product, error)
	SearchProduct(ctx context.Context, query string, skip, take uint64) ([]Product, error)
}

type catalogService struct {
	repository Repository
}

func NewService(r Repository) (Serice, error) {
	return &catalogService{repository: r}, nil
}

// CreateProduct implements Serice.
func (c *catalogService) CreateProduct(ctx context.Context, name, description string, price float64) (*Product, error) {
	p := Product{
		ID:          ksuid.New().String(),
		Name:        name,
		Description: description,
		Price:       price,
	}
	if err := c.repository.CreateProduct(ctx, p); err != nil {
		return nil, err
	}
	return &p, nil
}

// GetListProducts implements Serice.
func (c *catalogService) GetListProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}
	return c.repository.GetListProducts(ctx, skip, take)
}

// GetListProductsByIDs implements Serice.
func (c *catalogService) GetListProductsByIDs(ctx context.Context, ids []string) ([]Product, error) {
	return c.repository.GetListProductWithIDs(ctx, ids)
}

// GetProductById implements Serice.
func (c *catalogService) GetProductById(ctx context.Context, id string) (*Product, error) {
	return c.repository.GetProductById(ctx, id)
}

// SearchProduct implements Serice.
func (c *catalogService) SearchProduct(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}
	return c.repository.SearchProducts(ctx, query, skip, take)
}
