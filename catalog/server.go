package catalog

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/coderkhushal/go-grpc-graphql-microservices/catalog/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	service Service
	pb.UnimplementedCatalogServiceServer
}

func ListenGRPC(s Service, port int) error {
	list, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	serv := grpc.NewServer()
	pb.RegisterCatalogServiceServer(serv, &grpcServer{service: s})

	reflection.Register(serv)

	return serv.Serve(list)

}
func (s *grpcServer) PostProduct(ctx context.Context, r *pb.PostProductRequest) (*pb.PostProductResponse, error) {
	p, err := s.service.PostProduct(ctx, r.Name, r.Description, float64(r.Price))
	if err != nil {
		return nil, err
	}
	return &pb.PostProductResponse{
		Product: &pb.Product{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		},
	}, nil
}

func (s *grpcServer) GetProduct(ctx context.Context, r *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	p, err := s.service.GetProduct(ctx, r.Id)

	if err != nil {
		return nil, err
	}
	return &pb.GetProductResponse{
		Product: &pb.Product{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		},
	}, nil
}

func (s *grpcServer) GetProducts(ctx context.Context, r *pb.GetProductsRequest) (*pb.GetProductsResponse, error) {
	var res []Product
	var err error

	if r.Query != "" {
		res, err = s.service.SearchProducts(ctx, r.Query, int(r.Skip), int(r.Take))
	} else if len(r.Ids) != 0 {
		res, err = s.service.GetProductsByIDs(ctx, r.Ids)
	} else {
		res, err = s.service.GetProducts(ctx, int(r.Skip), int(r.Take))
	}

	if err != nil {
		log.Println(err)
		return nil, err
	}
	products := []*pb.Product{}
	for _, value := range res {
		products = append(products, &pb.Product{
			Id:          value.ID,
			Name:        value.Name,
			Description: value.Description,
			Price:       value.Price,
		})
	}

	return &pb.GetProductsResponse{
		Products: products,
	}, nil
}
