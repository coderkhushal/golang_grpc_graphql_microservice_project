package order

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lib/pq"
)

type Repository interface {
	Close()
	PutOrder(ctx context.Context, o Order) error
	GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error)
}

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (Repository, error) {
	db, err := sql.Open("pgx", url)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &postgresRepository{
		db: db,
	}, nil
}

func (r *postgresRepository) Close() {
	r.db.Close()
}
func (r *postgresRepository) Ping() error {
	return r.db.Ping()
}

func (r *postgresRepository) PutOrder(ctx context.Context, o Order) error {
	tx, err := r.db.BeginTx(ctx, nil)
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
	_, err = tx.ExecContext(ctx, "INSERT INTO orders (id, created_at, account_id, total_price) VALUES ($1, $2, $3 ,$4)", o.ID, o.CreatedAt, o.AccountID, o.TotalPrice)
	if err != nil {
		return err
	}

	stmt, _ := tx.PrepareContext(ctx, pq.CopyIn("order_products", "order_id", "product_id", "quantity"))
	for _, p := range o.Products {
		_, err = stmt.ExecContext(ctx, o.ID, p.ID, p.Quantity)
		if err != nil {
			return nil
		}
	}
	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return nil
	}
	stmt.Close()
	return err
}
func (r *postgresRepository) GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT 
		o.id,
		o.created_id, 
		o.account_id, 
		o.total_price::money::float8::numeric,
		op.product_id,
		op.quantity 
		FROM orders o JOIN order_products op ON(o.id = op.order_id)
		WHERE o.account_id = $1,
		ORDER BY o.id`, accountID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	orders := []Order{}
	order := &Order{} // for scanning
	lastOrder := &Order{}
	products := []OrderedProduct{}
	orderedProduct := &OrderedProduct{} // for scanning

	for rows.Next() {
		if err := rows.Scan(
			&order.ID,
			&order.CreatedAt,
			&order.AccountID,
			&order.TotalPrice,
			&orderedProduct.ID,
			&orderedProduct.Quantity,
		); err != nil {
			return nil, err
		}
		// if a touple comes with new orderid that means the products in previous order are complete.
		// so make a order struct with all order details and products in completed order
		// and append in orders slice which finally has to be returned as answer
		if lastOrder.ID != "" && lastOrder.ID != order.ID {

			newOrder := Order{
				ID:         lastOrder.ID,
				AccountID:  lastOrder.AccountID,
				CreatedAt:  lastOrder.CreatedAt,
				TotalPrice: lastOrder.TotalPrice,
				Products:   products,
			}
			orders = append(orders, newOrder)
			products = []OrderedProduct{}
		}

		products = append(products, OrderedProduct{
			ID:       orderedProduct.ID,
			Quantity: orderedProduct.Quantity,
		})
		*lastOrder = *order
	}
	if lastOrder != nil {

		newOrder := Order{
			ID:         lastOrder.ID,
			AccountID:  lastOrder.AccountID,
			CreatedAt:  lastOrder.CreatedAt,
			TotalPrice: lastOrder.TotalPrice,
			Products:   products,
		}
		orders = append(orders, newOrder)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}
