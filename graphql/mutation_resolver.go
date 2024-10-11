package main

import (
	"context"
	"fmt"
	"shop-graphql-demo/order"
	"time"
)

type mutationResolver struct {
	server *Server
}

func (m *mutationResolver) CreateProduct(ctx context.Context, product ProductInput) (*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	p, err := m.server.catalogClient.CreateProduct(ctx, product.Name, product.Description, product.Price)
	if err != nil {
		return nil, err
	}
	return &Product{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
	}, nil

}
func (m *mutationResolver) CreataOrder(ctx context.Context, o OrderInput) (*Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	products := []order.OrderedProduct{}
	for _, p := range o.Products {
		if p.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity")
		}
		products = append(products, order.OrderedProduct{
			ID:       p.ID,
			Quantity: uint32(p.Quantity),
		})
	}
	rs, err := m.server.orderClient.CreateOrder(ctx, o.AccountID, products)
	if err != nil {
		return nil, err
	}
	return &Order{
		ID:         rs.ID,
		TotalPrice: rs.TotalPrice,
		CreatedAt:  rs.CreatedAt,
	}, nil
}
func (m *mutationResolver) CreateAccount(ctx context.Context, accout AccountInput) (*Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	a, err := m.server.accountClient.CreateAccount(ctx, accout.Name)
	if err != nil {
		return nil, err
	}
	return &Account{
		ID:   a.ID,
		Name: a.Name,
	}, nil
}
