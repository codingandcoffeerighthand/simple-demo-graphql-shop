package account

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Repository interface {
	Close()
	Ping() error
	CreateAccount(ctx context.Context, account Account) error
	GetAccountByID(ctx context.Context, id string) (*Account, error)
	ListAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error)
}

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (Repository, error) {
	conn, err := sql.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %w", err)
	}
	err = conn.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping: %w", err)
	}
	return &postgresRepository{db: conn}, nil
}

func (r *postgresRepository) Close() {
	r.db.Close()
}
func (r *postgresRepository) Ping() error {
	return r.db.Ping()
}
func (r *postgresRepository) CreateAccount(ctx context.Context, account Account) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO accounts (id, name) VALUES ($1, $2)", account.ID, account.Name)
	return err
}
func (r *postgresRepository) GetAccountByID(ctx context.Context, id string) (*Account, error) {
	row, err := r.db.QueryContext(ctx, "SELECT id, name FROM accounts WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	var account Account
	defer row.Close()
	for row.Next() {
		if err := row.Scan(&account.ID, &account.Name); err != nil {
			return nil, err
		}
	}
	return &account, nil
}
func (r *postgresRepository) ListAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error) {
	rows, err := r.db.QueryContext(
		ctx,
		"SELECT id, name FROM accounts ORDER BY id DESC LIMIT $1 OFFSET $2",
		take, skip)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var accounts []Account
	for rows.Next() {
		var account Account
		if err := rows.Scan(&account.ID, &account.Name); err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}
