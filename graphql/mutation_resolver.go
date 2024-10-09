package main

import "context"

type mutationResolver struct {
	server *Server
}

func (m *mutationResolver) CreateProduct(ctx context.Context, product ProductInput) (*Product, error)
func (m *mutationResolver) CreataOrder(ctx context.Context, order OrderInput) (*Order, error)
func (m *mutationResolver) CreateAccount(ctx context.Context, accout AccountInput) (*Account, error)
