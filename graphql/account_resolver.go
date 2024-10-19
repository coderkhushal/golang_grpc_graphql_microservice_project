package main

import (
	"context"
	"log"
	"time"
)

type Accountresolver struct {
	server *Server
}

func (r *Accountresolver) Orders(ctx context.Context, obj *Account) ([]*Order, error) {

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	ordersList, err := r.server.orderClient.GetOrdersForAccount(ctx, obj.ID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var orders []*Order
	for _, order := range ordersList {
		orderedProducts := []*OrderedProduct{}
		for _, product := range order.Products {
			orderedProducts = append(orderedProducts, &OrderedProduct{
				ID:          product.ID,
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
				Quantity:    product.Quantity,
			})
		}
		orders = append(orders, &Order{
			ID:         order.ID,
			CreatedAt:  order.CreatedAt,
			TotalPrice: order.TotalPrice,
			Products:   orderedProducts,
		})
	}
	return orders, nil
}
