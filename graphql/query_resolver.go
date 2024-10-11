package main

import (
	"context"
	"time"
)

type queryResolver struct {
	server *Server
}

func (r *queryResolver) Accounts(ctx context.Context, pagination *PaginationInput, id *string) ([]*Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if id != nil {
		r, err := r.server.accountClient.GetAccount(ctx, *id)
		if err != nil {
			return nil, err
		}
		return []*Account{&Account{
			ID:   r.ID,
			Name: r.Name,
		}}, nil
	}
	skip, take := uint64(0), uint64(10)
	if pagination != nil {
		skip = uint64(*pagination.Skip)
		take = uint64(*pagination.Take)
	}
	rs, err := r.server.accountClient.GetAccounts(ctx, skip, take)
	if err != nil {
		return nil, err
	}
	accounts := make([]*Account, len(rs))
	for i, a := range rs {
		accounts[i] = &Account{
			ID:   a.ID,
			Name: a.Name,
		}
	}

	return accounts, err
}

func (r *queryResolver) Products(ctx context.Context, pagination *PaginationInput, query *string, id *string) ([]*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if id != nil {
		r, err := r.server.catalogClient.GetProduct(ctx, *id)
		if err != nil {
			return nil, err
		}
		return []*Product{&Product{
			ID:          r.ID,
			Name:        r.Name,
			Description: r.Description,
			Price:       r.Price,
		}}, nil
	}
	skip, take := uint64(0), uint64(10)
	if pagination != nil {
		skip = uint64(*pagination.Skip)
		take = uint64(*pagination.Take)
	}
	q := ""
	if query != nil {
		q = *query
	}
	productList, err := r.server.catalogClient.GetProducts(ctx, skip, take, q, nil)
	if err != nil {
		return nil, err
	}
	products := make([]*Product, len(productList))
	for i, p := range productList {
		products[i] = &Product{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		}
	}
	return products, nil
}
