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

type transferSubscription struct {
	Name        *string
	Description *string
	Price       *int
}

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

func (r Repo) GetSubscriptions(ctx context.Context) ([]entities.Subscription, error) {
	queryContext, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(queryContext, "SELECT * FROM subscriptions")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entities.ErrNoAnySubscription
		}
		return nil, fmt.Errorf("query context failed: %w", err)
	}
	defer rows.Close()

	subs := make([]entities.Subscription, 0)
	for rows.Next() {
		sub := entities.NewSubscription()
		err := rows.Scan(&sub.ID, &sub.Name, &sub.Description, &sub.Price, &sub.Level)
		if err != nil {
			return nil, fmt.Errorf("rows scan failed: %w", err)
		}

		subs = append(subs, sub)
	}
	return subs, nil
}

func (r Repo) GetSubscription(ctx context.Context, id string) (entities.Subscription, error) {
	queryContext, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	row := r.db.QueryRowContext(queryContext, "SELECT * FROM subscriptions WHERE id = $1", id)
	if row.Err() != nil {
		if errors.Is(row.Err(), sql.ErrNoRows) {
			return entities.NewSubscription(), entities.ErrSubscriptionDoesNotExist
		}
		return entities.NewSubscription(), fmt.Errorf("query row context failed: %w", row.Err())
	}

	sub := entities.NewSubscription()
	err := row.Scan(&sub.ID, &sub.Name, &sub.Description, &sub.Price, &sub.Level)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.NewSubscription(), entities.ErrSubscriptionDoesNotExist
		}
		return entities.NewSubscription(), fmt.Errorf("row scan failed: %w", err)
	}
	return sub, nil
}

func (r Repo) UpdateSubscription(ctx context.Context, id string, subscription entities.Subscription) error {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var transfer transferSubscription
	if subscription.Name != "" {
		transfer.Name = &subscription.Name
	}
	if subscription.Description != "" {
		transfer.Description = &subscription.Description
	}
	if subscription.Price != 0 {
		transfer.Price = &subscription.Price
	}

	res, err := r.db.ExecContext(queryCtx, "UPDATE subscriptions SET name = COALESCE($1, name), description = COALESCE($2, description), price = COALESCE($3, price) WHERE id = $4",
		transfer.Name, transfer.Description, transfer.Price, id)
	if err != nil {
		return fmt.Errorf("exec context failed: %w", err)
	}

	num, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected failed: %w", err)
	}
	if num == 0 {
		return entities.ErrSubscriptionDoesNotExist
	}
	return nil
}

func (r Repo) DeleteSubscription(ctx context.Context, id string) error {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := r.db.ExecContext(queryCtx, "DELETE FROM subscriptions WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("exec context failed: %w", err)
	}

	num, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected failed: %w", err)
	}
	if num == 0 {
		return entities.ErrSubscriptionDoesNotExist
	}
	return nil
}
