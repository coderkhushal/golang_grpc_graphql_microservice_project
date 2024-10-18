package catalog

import (
	"context"
	"log"

	"github.com/coderkhushal/go-grpc-graphql-microservices/catalog/pb"
	"google.golang.org/grpc"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.CatalogServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.NewClient(url, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Client connection failed in catalog: %v", err)
	}

	c := pb.NewCatalogServiceClient(conn)

	return &Client{
		conn:    conn,
		service: c,
	}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) PostProduct(ctx context.Context, name string, description string, price float64) (*Product, error) {
	res, err := c.service.PostProduct(ctx, &pb.PostProductRequest{
		Name:        name,
		Description: description,
		Price:       price,
	})
	if err != nil {
		return nil, err
	}
	return &Product{
		ID:          res.Product.Id,
		Name:        res.Product.Name,
		Description: res.Product.Description,
		Price:       res.Product.Price,
	}, nil
}

func (c *Client) GetProduct(ctx context.Context, id string) (*Product, error) {
	res, err := c.service.GetProduct(ctx, &pb.GetProductRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return &Product{
		ID:          res.Product.Id,
		Name:        res.Product.Name,
		Description: res.Product.Description,
		Price:       res.Product.Price,
	}, nil
}

func (c *Client) GetProducts(ctx context.Context, query string, skip int, take int, ids []string) ([]Product, error) {
	res, err := c.service.GetProducts(ctx, &pb.GetProductsRequest{
		Skip:  int64(skip),
		Take:  int64(take),
		Ids:   ids,
		Query: query,
	})
	if err != nil {
		return nil, err
	}
	var products []Product
	for _, value := range res.Products {
		products = append(products, Product{
			ID:          value.Id,
			Name:        value.Name,
			Description: value.Description,
			Price:       value.Price,
		})
	}
	return products, nil
}
