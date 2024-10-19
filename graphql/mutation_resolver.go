package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/coderkhushal/go-grpc-graphql-microservices/order"
)

var (
	ErrInvalidParameter = errors.New("Invalid parameter ")
)

type Mutationresolver struct {
	server *Server
}

func (r *Mutationresolver) CreateAccount(ctx context.Context, in *AccountInput) (*Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	a, err := r.server.accoutClient.PostAccount(ctx, in.Name)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &Account{
		ID:   a.ID,
		Name: a.Name,
	}, nil
}
func (r *Mutationresolver) CreateProduct(ctx context.Context, in *ProductInput) (*Product, error) {

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	p, err := r.server.catalogClient.PostProduct(ctx, in.Name, in.Description, in.Price)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &Product{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
	}, nil
}

func (r *Mutationresolver) CreateOrder(ctx context.Context, in *OrderInput) (*Order, error) {

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	var products []order.OrderedProduct
	for _, value := range in.Products {
		if value.Quantity <= 0 {
			return nil, ErrInvalidParameter
		}
		products = append(products, order.OrderedProduct{
			ID:       value.ID,
			Quantity: value.Quantity,
		})
	}
	o, err := r.server.orderClient.PostOrder(ctx, in.AccountID, products)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	order := Order{
		ID:         o.ID,
		CreatedAt:  o.CreatedAt,
		TotalPrice: o.TotalPrice,
		Products:   []*OrderedProduct{},
	}
	for _, value := range o.Products {
		order.Products = append(order.Products, &OrderedProduct{
			ID:          value.ID,
			Name:        value.Name,
			Description: value.Description,
			Price:       value.Price,
			Quantity:    value.Quantity,
		})
	}
	return &order, nil
}
