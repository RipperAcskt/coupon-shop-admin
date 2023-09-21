package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/RipperAcskt/coupon-shop-admin/internal/entities"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func (r Repo) CreateSubscription(ctx context.Context, sub entities.Subscription) error {
	queryContext, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var id int
	err := r.db.QueryRowContext(queryContext, "SELECT id FROM subscriptions WHERE name = $1 AND description = $2 AND price = $3",
		sub.Name, sub.Description, sub.Price).Scan(&id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("query row context failed: %w", err)
	}

	_, err = r.db.ExecContext(queryContext, "INSERT INTO subscriptions VALUES ($1, $2, $3, $4, $5)",
		sub.ID, sub.Name, sub.Description, sub.Price, sub.Level)
	if err != nil {
		return fmt.Errorf("exec context failed: %w", err)
	}

	return nil
}
