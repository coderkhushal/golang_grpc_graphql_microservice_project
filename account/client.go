package account

import (
	"context"
	"log"

	"github.com/coderkhushal/go-grpc-graphql-microservices/account/pb"
	"google.golang.org/grpc"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.AccountServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.NewClient(url, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Connection failed %v", err)

	}

	c := pb.NewAccountServiceClient(conn)

	return &Client{
		conn:    conn,
		service: c,
	}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) PostAccount(ctx context.Context, name string) (*Account, error) {
	// INVOKING GRPC FUNCTIONS FROM CLIENT
	res, err := c.service.PostAccount(ctx, &pb.PostAccountRequest{Name: name})
	if err != nil {
		return nil, err
	}
	return &Account{
		ID: res.Account.
			Id, Name: res.Account.Name,
	}, nil

}

func (c *Client) GetAccount(ctx context.Context, id string) (*Account, error) {
	res, err := c.service.GetAccount(ctx, &pb.GetAccountRequest{
		Id: id,
	})

	if err != nil {
		return nil, err
	}
	return &Account{
		ID: res.Account.
			Id, Name: res.Account.Name,
	}, nil
}
func (c *Client) GetAccounts(ctx context.Context, skip int, take int) ([]Account, error) {
	res, err := c.service.GetAccounts(ctx, &pb.GetAccountsRequest{
		Skip: int64(skip),
		Take: int64(take),
	})

	if err != nil {
		return nil, err
	}
	accounts := []Account{}
	for _, account := range res.Accounts {
		accounts = append(accounts, Account{
			ID:   account.Id,
			Name: account.Name,
		})
	}
	return accounts, nil

}
