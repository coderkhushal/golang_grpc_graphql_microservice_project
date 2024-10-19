package order

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/coderkhushal/go-grpc-graphql-microservices/account"
	"github.com/coderkhushal/go-grpc-graphql-microservices/catalog"
	"github.com/coderkhushal/go-grpc-graphql-microservices/order/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	service       Service
	accountClient *account.Client
	catalogClient *catalog.Client
	pb.UnimplementedOrderServiceServer
}

func ListenGRPC(s Service, accountURL, catalogURL string, port int) error {
	accountClient, err := account.NewClient(accountURL)

	if err != nil {

		return err
	}
	catalogClient, err := catalog.NewClient(catalogURL)

	if err != nil {
		accountClient.Close()
		return err
	}
	lis, err := net.Listen("tpc", fmt.Sprintf("%s", port))
	if err != nil {
		accountClient.Close()
		catalogClient.Close()
		return err
	}

	serv := grpc.NewServer()
	pb.RegisterOrderServiceServer(serv, &grpcServer{
		service:       s,
		accountClient: accountClient,
		catalogClient: catalogClient,
	})

	reflection.Register(serv)

	return serv.Serve(lis)
}

func (s *grpcServer) PostOrder(ctx context.Context, r *pb.PostOrderRequest) (*pb.PostOrderResponse, error) {
	_, err := s.accountClient.GetAccount(ctx, r.AccountId)

	if err != nil {
		log.Printf("Error getting account %v", err)

		return nil, errors.New("account not found ")
	}

	productIds := []string{}
	orderedProducts, err := s.catalogClient.GetProducts(ctx, "", 0, 0, productIds)
	if err != nil {
		log.Println("error getting products", err)
		return nil, errors.New("Products not found")

	}
	// converting to service readable format
	products := []OrderedProduct{}
	for _, value := range orderedProducts {
		product := OrderedProduct{
			ID:          value.ID,
			Name:        value.Name,
			Description: value.Description,
			Price:       value.Price,
			Quantity:    0,
		}
		for _, rp := range r.Products {
			if rp.ProductId == product.ID {
				product.Quantity = int(rp.Quantity)
				break
			}
		}
		if product.Quantity != 0 {
			products = append(products, product)
		}
	}
	// creating order
	res, err := s.service.PostOrder(ctx, r.AccountId, products)
	if err != nil {
		return nil, errors.New("Could not post order")
	}
	// converting to response
	orderedProto := &pb.Order{
		Id:         res.ID,
		AccountId:  res.AccountID,
		TotalPrice: res.TotalPrice,
		Products:   []*pb.Order_OrderProduct{},
	}
	orderedProto.CreatedAt, _ = res.CreatedAt.MarshalBinary()
	for _, value := range products {
		orderedProto.Products = append(orderedProto.Products, &pb.Order_OrderProduct{
			Id:          value.ID,
			Name:        value.Name,
			Description: value.Description,
			Price:       value.Price,
			Quantity:    uint64(value.Quantity),
		})

	}
	return &pb.PostOrderResponse{
		Order: orderedProto,
	}, nil
}
func (s *grpcServer) GetOrdersForAccount(ctx context.Context, r *pb.GetOrdersForAccountRequest) (*pb.GetOrdersForAccountResponse, error) {
	_, err := s.accountClient.GetAccount(ctx, r.AccountId)
	if err != nil {
		log.Println("Erorr in getting account")
		return nil, errors.New("account not found")
	}

	accountOrders, err := s.service.GetOrdersForAccount(ctx, r.AccountId)

	if err != nil {
		return nil, errors.New("cannot get orders")
	}
	// fetching product name ,description and other things from elastic search as we doent have it in accountorders products
	productIdMap := map[string]bool{}
	for _, order := range accountOrders {
		for _, product := range order.Products {
			productIdMap[product.ID] = true
		}
	}
	productIds := []string{}
	for id := range productIdMap {
		productIds = append(productIds, id)
	}
	products, err := s.catalogClient.GetProducts(ctx, "", 0, 0, productIds)
	if err != nil {
		log.Println("error finding product")
		return nil, errors.New("products not found")
	}

	// converting to proto response
	ordersProto := []*pb.Order{}

	for _, order := range accountOrders {
		// appending product information from elastic search into proto
		orderedproductsProto := []*pb.Order_OrderProduct{}
		for _, p := range order.Products {
			for _, product := range products {
				if product.ID == p.ID {
					p.Name = product.Name
					p.Description = product.Description
					p.Price = product.Price
				}
				break
			}
			orderedproductsProto = append(orderedproductsProto, &pb.Order_OrderProduct{
				Id:          p.ID,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
				Quantity:    uint64(p.Price),
			})
		}
		o := pb.Order{

			Id:         order.ID,
			AccountId:  order.AccountID,
			TotalPrice: order.TotalPrice,
			Products:   orderedproductsProto,
		}
		o.CreatedAt, _ = order.CreatedAt.MarshalBinary()
		ordersProto = append(ordersProto, &o)
	}
	return &pb.GetOrdersForAccountResponse{
		Orders: ordersProto,
	}, nil
}
