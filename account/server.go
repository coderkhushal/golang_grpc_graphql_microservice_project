package account

import (
	"context"
	"fmt"
	"net"

	"github.com/coderkhushal/go-grpc-graphql-microservices/account/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// mustEmbedUnimplementedAccountServiceServer implements pb.AccountServiceServer.

type grpcServer struct {
	service Service
	pb.UnimplementedAccountServiceServer
}

func ListenGRPC(s Service, port int) error {
	// listening server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	// creating new grpc server
	serv := grpc.NewServer()
	// register service in grpc server
	pb.RegisterAccountServiceServer(serv, &grpcServer{
		service: s,
	})
	// exposing endpoints , better for development ( i think auto completion )
	reflection.Register(serv)
	// running server
	return serv.Serve(lis)
}

func (s *grpcServer) PostAccount(
	ctx context.Context,
	r *pb.PostAccountRequest,
) (*pb.PostAccountResponse, error) {
	a, err := s.service.PostAccount(ctx, r.Name)
	if err != nil {
		return nil, err
	}

	return &pb.PostAccountResponse{
		Account: &pb.Account{
			Id:   a.ID,
			Name: a.Name,
		},
	}, nil
}

func (s *grpcServer) GetAccount(ctx context.Context, r *pb.GetAccountRequest) (*pb.GetAccountResponse, error) {
	a, err := s.service.GetAccount(ctx, r.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetAccountResponse{
		Account: &pb.Account{
			Id:   a.ID,
			Name: a.Name,
		},
	}, nil
}
func (s *grpcServer) GetAccounts(ctx context.Context, r *pb.GetAccountsRequest) (*pb.GetAccountsResponse, error) {
	res, err := s.service.GetAccounts(ctx, int(r.Skip), int(r.Take))
	if err != nil {
		return nil, err
	}
	accounts := []*pb.Account{}

	for _, a := range res {
		accounts = append(accounts,
			&pb.Account{
				Id:   a.ID,
				Name: a.Name,
			})
	}
	return &pb.GetAccountsResponse{
		Accounts: accounts,
	}, nil
}
