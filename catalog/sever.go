//go:generate protoc --go_out=./pb --go_opt=paths=source_relative --go-grpc_out=./pb --go-grpc_opt=paths=source_relative catalog.proto
package catalog

import (
	"context"
	"fmt"
	"net"
	"shop-graphql-demo/catalog/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	pb.UnimplementedCatalogServiceServer
	service Serice
}

func ListenGRPC(s Serice, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	srv := grpc.NewServer()
	pb.RegisterCatalogServiceServer(srv, &grpcServer{service: s})
	reflection.Register(srv)
	return srv.Serve(lis)
}
func (s *grpcServer) CreateProduct(ctx context.Context, r *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	product, err := s.service.CreateProduct(context.Background(), r.Name, r.Description, r.Price)
	if err != nil {
		return nil, err
	}
	return &pb.CreateProductResponse{Product: &pb.Product{Id: product.ID, Name: product.Name}}, nil
}
func (s *grpcServer) GetProduct(ctx context.Context, r *pb.GetProductRequests) (*pb.GetProductResponses, error) {
	p, err := s.service.GetProductById(context.Background(), r.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetProductResponses{Product: &pb.Product{Id: p.ID, Name: p.Name}}, nil
}
func (s *grpcServer) GetProducts(context.Context, *pb.GetProductsRequest) (*pb.GetProductsResponse, error) {
}
