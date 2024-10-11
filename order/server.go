//go:generate protoc --go_out=./pb --go_opt=paths=source_relative --go-grpc_out=./pb --go-grpc_opt=paths=source_relative order.proto
package order

import (
	"context"
	"fmt"
	"net"
	"shop-graphql-demo/account"
	"shop-graphql-demo/catalog"
	"shop-graphql-demo/order/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	pb.UnimplementedOrderServiceServer
	service       Service
	accountClient *account.Client
	catalogClient *catalog.Client
}

func ListenGRPC(s Service, accountUrl, catalogUrl string, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	srv := grpc.NewServer()
	accountClient, err := account.NewClient(accountUrl)
	if err != nil {
		return fmt.Errorf("failed to create account client: %w", err)
	}
	catalogClient, err := catalog.NewClient(catalogUrl)
	if err != nil {
		accountClient.Close()
		return fmt.Errorf("failed to create catalog client: %w", err)
	}
	sv := &grpcServer{
		service:       s,
		accountClient: accountClient,
		catalogClient: catalogClient,
	}
	pb.RegisterOrderServiceServer(srv, sv)
	reflection.Register(srv)
	return srv.Serve(lis)
}

func (s *grpcServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	_, err := s.accountClient.GetAccount(ctx, req.AccountId)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}
	productIDs := []string{}
	for _, p := range req.Products {
		productIDs = append(productIDs, p.ProductId)
	}
	orderedProducts, err := s.catalogClient.GetProducts(ctx, 0, 0, "", productIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}
	products := []OrderedProduct{}
	for _, p := range orderedProducts {
		product := OrderedProduct{
			ID:          p.ID,
			Name:        p.Name,
			Price:       p.Price,
			Quantity:    0,
			Description: p.Description,
		}
		for _, rp := range req.Products {
			if rp.ProductId == p.ID {
				product.Quantity = rp.Quantity
				break
			}
		}
		if product.Quantity > 0 {
			products = append(products, product)
		}
	}
	order, err := s.service.CreateOrder(ctx, req.AccountId, products)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}
	orderProto := &pb.Order{
		Id:         order.ID,
		AccountId:  order.AccountID,
		TotalPrice: order.TotalPrice,
		Products:   []*pb.Order_OrderProduct{},
	}
	orderProto.CreatedAt, _ = order.CreatedAt.MarshalBinary()
	for _, p := range order.Products {
		orderProto.Products = append(
			orderProto.Products,
			&pb.Order_OrderProduct{
				Id:          p.ID,
				Name:        p.Name,
				Price:       p.Price,
				Quantity:    p.Quantity,
				Description: p.Description,
			},
		)
	}
	return &pb.CreateOrderResponse{Order: orderProto}, nil
}

func (s *grpcServer) GetOrderForAccount(ctx context.Context, req *pb.GetOrderForAccountRequest) (*pb.GetOrderForAccountResponse, error) {
	accountOrders, err := s.service.GetOrderForAccount(ctx, req.AccountId)
	if err != nil {
		return nil, fmt.Errorf("failed to get order for account: %w", err)
	}
	productIDMap := map[string]bool{}
	for _, o := range accountOrders {
		for _, p := range o.Products {
			productIDMap[p.ID] = true
		}
	}
	producIDs := []string{}
	for id := range productIDMap {
		producIDs = append(producIDs, id)
	}
	products, err := s.catalogClient.GetProducts(ctx, 0, 0, "", producIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}
	orders := []*pb.Order{}
	for _, o := range accountOrders {
		op := &pb.Order{
			AccountId:  o.AccountID,
			Id:         o.ID,
			TotalPrice: o.TotalPrice,
			Products:   []*pb.Order_OrderProduct{},
		}
		op.CreatedAt, _ = o.CreatedAt.MarshalBinary()
		for _, product := range o.Products {
			for _, p := range products {
				if p.ID == product.ID {
					product.Name = p.Name
					product.Description = p.Description
					product.Price = p.Price
					break
				}
			}
			op.Products = append(
				op.Products,
				&pb.Order_OrderProduct{
					Id:          product.ID,
					Name:        product.Name,
					Price:       product.Price,
					Quantity:    product.Quantity,
					Description: product.Description,
				},
			)
		}
		orders = append(orders, op)
	}
	return &pb.GetOrderForAccountResponse{Orders: orders}, nil
}
