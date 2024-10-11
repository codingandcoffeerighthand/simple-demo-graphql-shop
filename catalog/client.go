package catalog

import (
	"context"
	"shop-graphql-demo/catalog/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.CatalogServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &Client{
		conn:    conn,
		service: pb.NewCatalogServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) CreateProduct(ctx context.Context, name, description string, price float64) (*Product, error) {
	resp, err := c.service.CreateProduct(ctx, &pb.CreateProductRequest{
		Name:        name,
		Description: description,
		Price:       price,
	})
	if err != nil {
		return nil, err
	}
	return &Product{
		ID:          resp.Product.Id,
		Name:        resp.Product.Name,
		Description: resp.Product.Description,
		Price:       resp.Product.Price,
	}, nil
}

func (c *Client) GetProduct(ctx context.Context, id string) (*Product, error) {
	resp, err := c.service.GetProduct(ctx, &pb.GetProductRequests{Id: id})
	if err != nil {
		return nil, err
	}
	return &Product{
		ID:          resp.Product.Id,
		Name:        resp.Product.Name,
		Description: resp.Product.Description,
		Price:       resp.Product.Price,
	}, nil
}

func (c *Client) GetProducts(ctx context.Context, skip, take uint64, query string, ids []string) ([]Product, error) {
	resp, err := c.service.GetProducts(ctx, &pb.GetProductsRequest{Skip: skip, Take: take, Query: query, Ids: ids})
	if err != nil {
		return nil, err
	}
	products := make([]Product, 0, len(resp.Products))
	for _, p := range resp.Products {
		products = append(products, Product{
			ID:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		})
	}
	return products, nil
}
