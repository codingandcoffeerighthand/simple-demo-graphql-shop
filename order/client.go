package order

import (
	"context"
	"shop-graphql-demo/order/pb"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.OrderServiceClient
}

func NewClient(orderUrl string) (*Client, error) {
	conn, err := grpc.NewClient(orderUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &Client{
		conn:    conn,
		service: pb.NewOrderServiceClient(conn),
	}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) CreateOrder(ctx context.Context, userId string, products []OrderedProduct) (*Order, error) {
	protoProduct := []*pb.CreateOrderRequest_OrderProduct{}
	for _, p := range products {
		protoProduct = append(protoProduct, &pb.CreateOrderRequest_OrderProduct{
			ProductId: p.ID,
			Quantity:  p.Quantity,
		})
	}
	r, err := c.service.CreateOrder(ctx, &pb.CreateOrderRequest{
		AccountId: userId,
		Products:  protoProduct,
	})
	if err != nil {
		return nil, err
	}
	newOrder := r.Order
	newOrderCreatedAt := time.Time{}
	newOrderCreatedAt.UnmarshalBinary(newOrder.CreatedAt)
	return &Order{
		ID:         newOrder.Id,
		CreatedAt:  newOrderCreatedAt,
		TotalPrice: newOrder.TotalPrice,
		AccountID:  newOrder.AccountId,
		Products:   products,
	}, nil
}

func (c *Client) GetOrderForAccount(ctx context.Context, accountID string) ([]Order, error) {
	r, err := c.service.GetOrderForAccount(ctx, &pb.GetOrderForAccountRequest{AccountId: accountID})
	if err != nil {
		return nil, err
	}
	orders := []Order{}
	for _, o := range r.Orders {
		ord := Order{
			ID:         o.Id,
			TotalPrice: o.TotalPrice,
			AccountID:  o.AccountId,
		}
		ord.CreatedAt = time.Time{}
		ord.CreatedAt.UnmarshalBinary(o.CreatedAt)
		products := []OrderedProduct{}
		for _, p := range o.Products {
			products = append(products, OrderedProduct{
				ID:          p.Id,
				Name:        p.Name,
				Price:       p.Price,
				Quantity:    p.Quantity,
				Description: p.Description,
			})
		}
		ord.Products = products
		orders = append(orders, ord)
	}
	return orders, nil
}
