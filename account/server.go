//go:generate protoc --go_out=./pb --go_opt=paths=source_relative --go-grpc_out=./pb --go-grpc_opt=paths=source_relative account.proto
package account

import (
	"context"
	"fmt"
	"net"
	"shop-graphql-demo/account/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	pb.UnimplementedAccountServiceServer
	service Service
}

func ListenGRPC(s Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	srv := grpc.NewServer()
	pb.RegisterAccountServiceServer(srv, &grpcServer{service: s})
	reflection.Register(srv)
	return srv.Serve(lis)
}

func (s *grpcServer) CreateAccount(ctx context.Context, r *pb.CreeateAccountRequest) (*pb.CreateAccountResponse, error) {
	a, err := s.service.CreateAccount(ctx, r.Name)
	if err != nil {
		return nil, err
	}
	var account = pb.Account{Id: a.ID, Name: a.Name}
	return &pb.CreateAccountResponse{Account: &account}, nil
}

func (s *grpcServer) GetAccount(ctx context.Context, r *pb.GetAccountRequest) (*pb.GetAccountResponse, error) {
	a, err := s.service.GetAccount(ctx, r.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetAccountResponse{Account: &pb.Account{Id: a.ID, Name: a.Name}}, nil
}

func (s *grpcServer) GetAccounts(ctx context.Context, r *pb.GetAccountsRequest) (*pb.GetAccountsResponse, error) {
	accounts, err := s.service.GetAccounts(ctx, r.Skip, r.Take)
	if err != nil {
		return nil, err
	}
	rs := make([]*pb.Account, len(accounts))
	for _, a := range accounts {
		rs = append(rs, &pb.Account{Id: a.ID, Name: a.Name})
	}
	return &pb.GetAccountsResponse{Accounts: rs}, nil
}
