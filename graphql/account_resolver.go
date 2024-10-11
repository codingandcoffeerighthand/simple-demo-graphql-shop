package main

import (
	"context"
	"time"
)

type accountResolver struct {
	server *Server
}

func (a *accountResolver) Orders(ctx context.Context, obj *Account) ([]*Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	orderList, err := a.server.orderClient.GetOrderForAccount(ctx, obj.ID)
	if err != nil {
		return nil, err
	}

	var order []*Order
	for _, o := range orderList {
		var products []*OrderedProduct
		for _, p := range o.Products {
			products = append(products, &OrderedProduct{
				ID:          p.ID,
				Name:        p.Name,
				Price:       p.Price,
				Quantity:    int(p.Quantity),
				Description: p.Description,
			})
		}
		order = append(order, &Order{
			ID:         o.ID,
			TotalPrice: o.TotalPrice,
			CreatedAt:  time.Time{},
			Products:   products,
		})
	}
	return order, nil
}
