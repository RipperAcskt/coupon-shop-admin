package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/RipperAcskt/coupon-shop-admin/internal/entities"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func (r Repo) CreateSubscription(ctx context.Context, sub entities.Subscription) error {
	queryContext, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var id string
	err := r.db.QueryRowContext(queryContext, "SELECT id FROM subscriptions WHERE level = $1", sub.Level).Scan(&id)
	if err == nil {
		return entities.ErrSubscriptionAlreadyExists
	}

	_, err = r.db.ExecContext(queryContext, "INSERT INTO subscriptions VALUES ($1, $2, $3, $4, $5)",
		sub.ID, sub.Name, sub.Description, sub.Price, sub.Level)
	if err != nil {
		return fmt.Errorf("exec context failed: %w", err)
	}

	return nil
}
