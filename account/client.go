package account

import (
	"context"
	"shop-graphql-demo/account/pb"

	"google.golang.org/grpc"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.AccountServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.NewClient(url)
	if err != nil {
		return nil, err
	}
	return &Client{
		conn:    conn,
		service: pb.NewAccountServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) CreateAccount(ctx context.Context, name string) (*Account, error) {
	resp, err := c.service.CreateAccount(ctx, &pb.CreeateAccountRequest{Name: name})
	if err != nil {
		return nil, err
	}
	return &Account{
		ID:   resp.Account.Id,
		Name: resp.Account.Name,
	}, nil
}

func (c *Client) GetAccount(ctx context.Context, id string) (*Account, error) {
	resp, err := c.service.GetAccount(ctx, &pb.GetAccountRequest{Id: id})
	if err != nil {
		return nil, err
	}
	return &Account{
		ID:   resp.Account.Id,
		Name: resp.Account.Name,
	}, nil
}

func (c *Client) GetAccounts(ctx context.Context, skip, take uint64) ([]Account, error) {
	resp, err := c.service.GetAccounts(ctx, &pb.GetAccountsRequest{Skip: skip, Take: take})
	if err != nil {
		return nil, err
	}
	accounts := make([]Account, 0, len(resp.Accounts))
	for _, a := range resp.Accounts {
		accounts = append(accounts, Account{
			ID:   a.Id,
			Name: a.Name,
		})
	}
	return accounts, nil
}
