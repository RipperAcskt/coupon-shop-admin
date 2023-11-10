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

func (r Repo) CreateRegion(ctx context.Context, region entities.Region) error {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	row := r.db.QueryRowContext(queryCtx, "INSERT INTO regions VALUES($1, $2, $3, $4)", region.Id, region.Name, region.Tg, region.Vk)
	if row.Err() != nil {
		return fmt.Errorf("query row context order failed: %w", row.Err())
	}

	return nil
}

func (r Repo) GetRegions(ctx context.Context) ([]entities.Region, error) {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(queryCtx, "SELECT * FROM regions")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entities.ErrNoAnyCoupons
		}
		return nil, fmt.Errorf("query context failed: %w", err)
	}
	defer rows.Close()

	regions := make([]entities.Region, 0)

	for rows.Next() {
		region := entities.Region{}
		err := rows.Scan(&region.Id, &region.Name, &region.Tg, &region.Vk)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}

		regions = append(regions, region)
	}

	return regions, nil
}

func (r Repo) UpdateRegion(ctx context.Context, id string, region entities.Region) error {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := r.db.ExecContext(queryCtx, "UPDATE regions SET name = $1, tg = $2, vk = $3 WHERE id = $4", region.Name, region.Tg, region.Vk, id)
	if err != nil {
		return fmt.Errorf("exec context failed: %w", err)
	}

	num, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected failed: %w", err)
	}
	if num == 0 {
		return entities.ErrCouponDoesNotExist
	}
	return nil
}

func (r Repo) DeleteRegion(ctx context.Context, id string) error {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := r.db.ExecContext(queryCtx, "DELETE FROM regions WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("exec context failed: %w", err)
	}

	num, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected failed: %w", err)
	}
	if num == 0 {
		return entities.ErrCouponDoesNotExist
	}
	return nil
}
