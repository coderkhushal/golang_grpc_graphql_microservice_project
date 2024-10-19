package main

import "context"

type mutationresolver struct {
	server *Server
}

func (r *mutationresolver) createAccount(ctx context.Context, in AccountInput) (*Account, error) {
}
func (r *mutationresolver) createProduct(ctx context.Context, in ProductInput) (*Product, error) {
}

func (r *mutationresolver) createOrder(ctx context.Context, in OrderInput) (*Order, error) {

}
