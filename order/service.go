package order

import (
	"context"
	"time"

	"github.com/segmentio/ksuid"
)

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
	Quantity    uint32
}
type Service interface {
	CreateOrder(ctx context.Context, accountID string, products []OrderedProduct) (*Order, error)
	GetOrderForAccount(ctx context.Context, accountID string) ([]Order, error)
}

type orderService struct {
	repo Repository
}

func (s *orderService) CreateOrder(ctx context.Context, accountID string, products []OrderedProduct) (*Order, error) {
	o := Order{
		ID:        ksuid.New().String(),
		CreatedAt: time.Now().UTC(),
		AccountID: accountID,
		Products:  products,
	}
	o.TotalPrice = 0.0
	for _, p := range o.Products {
		o.TotalPrice += p.Price * float64(p.Quantity)
	}
	err := s.repo.CreateOrder(ctx, o)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

// GetOrderForAccount implements Service.
func (o *orderService) GetOrderForAccount(ctx context.Context, accountID string) ([]Order, error) {
	return o.repo.GetOrderForAccount(ctx, accountID)
}

// New Service
func NewOrderService(repo Repository) Service {
	return &orderService{repo: repo}
}
