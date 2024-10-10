package order

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Repository interface {
	Close()
	CreateOrder(ctx context.Context, o Order) error
	GetOrderForAccount(ctx context.Context, accountID string) ([]Order, error)
}

type orderRepository struct {
	db *sql.DB
}

// Close implements Repository.
func (o *orderRepository) Close() {
	o.db.Close()
}

// CreateOrder implements Repository.
func (o *orderRepository) CreateOrder(ctx context.Context, order Order) (err error) {
	tx, err := o.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()
	_, err = tx.ExecContext(ctx,
		"INSERT INTO orders (id, created_at, account_id, total_price) VALUES ($1, $2, $3, $4)",
		order.ID, order.CreatedAt, order.AccountID, order.TotalPrice)
	if err != nil {
		return err
	}
	stmt, _ := tx.PrepareContext(ctx, pq.CopyIn("ordered_products", "order_id", "product_id", "quantity"))
	for _, p := range order.Products {
		_, err := stmt.ExecContext(ctx, order.ID, p.ID, p.Quantity)
		if err != nil {
			return err
		}
	}
	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return err
	}
	defer stmt.Close()
	return nil
}

// GetOrderForAccount implements Repository.
func (o *orderRepository) GetOrderForAccount(ctx context.Context, accountID string) ([]Order, error) {
	rows, err := o.db.QueryContext(ctx,
		`SELECT o.id, o.created_at, o.account_id , o.total_price::money::numeric::float
		FROM orders o INNER JOIN ordered_products op ON o.id = op.order_id 
		WHERE o.account_id = $1
		ORDERBY o.id`, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	orders := []Order{}
	order := &Order{}
	lastOrder := &Order{}
	orderedProducts := &OrderedProduct{}
	products := []OrderedProduct{}
	for rows.Next() {
		if err = rows.Scan(
			&order.ID,
			&order.CreatedAt,
			&order.AccountID,
			&order.TotalPrice,
			&orderedProducts.ID,
			&orderedProducts.Quantity,
		); err != nil {
			return nil, err
		}
		if lastOrder.ID != "" && lastOrder.ID != order.ID {
			newOrder := Order{
				ID:         lastOrder.ID,
				TotalPrice: lastOrder.TotalPrice,
				AccountID:  lastOrder.AccountID,
				CreatedAt:  lastOrder.CreatedAt,
				Products:   products,
			}
			orders = append(orders, newOrder)
			products = []OrderedProduct{}
		}
		products = append(products, OrderedProduct{
			ID:       orderedProducts.ID,
			Quantity: orderedProducts.Quantity,
		})
		*lastOrder = *order
	}
	if lastOrder.ID != "" {
		newOrder := Order{
			ID:         lastOrder.ID,
			TotalPrice: lastOrder.TotalPrice,
			AccountID:  lastOrder.AccountID,
			CreatedAt:  lastOrder.CreatedAt,
			Products:   products,
		}
		orders = append(orders, newOrder)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}

// NewOrderRepository implements Repository.
func NewOrderRepository(url string) (Repository, error) {
	conn, err := sql.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %w", err)
	}
	err = conn.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping: %w", err)
	}
	return &orderRepository{db: conn}, nil
}
