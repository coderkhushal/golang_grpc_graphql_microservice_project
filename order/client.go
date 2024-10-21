package order

import (
	"context"
	"log"

	"github.com/coderkhushal/go-grpc-graphql-microservices/order/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.OrderServiceClient
}

func NewClient(url string) (*Client, error) {
	// conn, err := grpc.NewClient(url, grpc.WithInsecure())
	// conn, err := grpc.Dial(url)

	conn, err := grpc.Dial(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("error connecting order client")
		return nil, err
	}
	c := pb.NewOrderServiceClient(conn)
	return &Client{
		conn:    conn,
		service: c,
	}, nil
}
func (c *Client) Close() {
	c.conn.Close()
}
func (c *Client) PostOrder(ctx context.Context, accountID string, products []OrderedProduct) (*Order, error) {
	productsProto := []*pb.PostOrderRequest_OrderProduct{}
	for _, value := range products {
		productsProto = append(productsProto, &pb.PostOrderRequest_OrderProduct{
			ProductId: value.ID,
			Quantity:  uint32(value.Quantity),
		})
	}
	res, err := c.service.PostOrder(ctx, &pb.PostOrderRequest{
		AccountId: accountID,
		Products:  productsProto,
	})
	if err != nil {
		return nil, err
	}
	neworder := Order{
		ID:         res.Order.Id,
		TotalPrice: res.Order.TotalPrice,
		AccountID:  res.Order.AccountId,
		Products:   products,
	}

	neworder.CreatedAt.UnmarshalBinary(res.Order.CreatedAt)
	return &neworder, nil

}

func (c *Client) GetOrdersForAccount(ctx context.Context, accountId string) ([]Order, error) {
	res, err := c.service.GetOrdersForAccount(ctx, &pb.GetOrdersForAccountRequest{
		AccountId: accountId,
	})
	if err != nil {
		return nil, err
	}
	orders := []Order{}
	for _, protoorder := range res.Orders {
		products := []OrderedProduct{}
		for _, protoproduct := range protoorder.Products {
			products = append(products, OrderedProduct{
				ID:          protoproduct.Id,
				Name:        protoproduct.Name,
				Description: protoproduct.Description,
				Price:       protoproduct.Price,
				Quantity:    int(protoproduct.Quantity),
			})
		}
		order := Order{
			ID:         protoorder.Id,
			TotalPrice: protoorder.TotalPrice,
			AccountID:  protoorder.AccountId,
			Products:   products,
		}
		order.CreatedAt.UnmarshalBinary(protoorder.CreatedAt)
		orders = append(orders, order)
	}
	return orders, nil
}
