package order

import (
	"context"
	"time"

	"github.com/segmentio/ksuid"
)

type Service interface {
	PostOrder(ctx context.Context, accountID string, products []OrderedProduct) (*Order, error)
	GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error)
}

type Order struct {
	ID         string
	CreatedAt  time.Time
	TotalPrice float64
	AccountID  string
	Products   []OrderedProduct
}
type OrderedProduct struct {
	ID          string
	Name        string
	Description string
	Price       float64
	Quantity    int
}

type orderservice struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &orderservice{
		repository: r,
	}
}

func (s *orderservice) PostOrder(ctx context.Context, accountID string, products []OrderedProduct) (*Order, error) {
	o := Order{
		ID:        ksuid.New().String(),
		CreatedAt: time.Now().UTC(),
		AccountID: accountID,
		Products:  products,
	}
	o.TotalPrice = 0.0
	for _, value := range products {
		o.TotalPrice = o.TotalPrice + (value.Price)*(float64(value.Quantity))
	}
	if err := s.repository.PutOrder(ctx, o); err != nil {
		return nil, err
	}
	return &o, nil
}

func (s *orderservice) GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error) {
	return s.repository.GetOrdersForAccount(ctx, accountID)
}
