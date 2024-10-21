package main

import (
	"context"
	"log"
	"time"
)

type queryResolver struct {
	server *Server
}

func (r *queryResolver) Accounts(ctx context.Context, pagination *PaginationInput, id *string) ([]*Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if id != nil {
		r, err := r.server.accountClient.GetAccount(ctx, *id)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return []*Account{
			{
				ID:   r.ID,
				Name: r.Name,
			}}, nil
	}
	skip, take := 0, 0
	if pagination != nil {

		skip, take = *pagination.Skip, *pagination.Take
	}

	accountList, err := r.server.accountClient.GetAccounts(ctx, skip, take)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	var accounts []*Account
	for _, value := range accountList {
		accounts = append(accounts, &Account{
			ID:   value.ID,
			Name: value.Name,
		})
	}
	return accounts, nil

}
func (r *queryResolver) Products(ctx context.Context, pagination *PaginationInput, query *string, id *string) ([]*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if id != nil {
		res, err := r.server.catalogClient.GetProduct(ctx, *id)

		if err != nil {
			log.Println(err)
			return nil, err
		}
		return []*Product{{
			ID:          res.ID,
			Description: res.Description,
			Name:        res.Name,
			Price:       res.Price,
		}}, nil
	}
	skip, take := 0, 0
	if pagination != nil && pagination.Skip != nil && pagination.Take != nil {
		skip, take = *pagination.Skip, *pagination.Take
	}
	q := ""
	if query != nil {
		q = *query
	}
	productList, err := r.server.catalogClient.GetProducts(ctx, q, skip, take, nil)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	var products []*Product
	for _, value := range productList {
		products = append(products, &Product{
			ID:          value.ID,
			Name:        value.Name,
			Description: value.Description,
			Price:       value.Price,
		})
	}
	return products, nil
}
