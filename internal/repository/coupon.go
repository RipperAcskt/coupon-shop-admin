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

func (r Repo) CreateCoupon(ctx context.Context, coupon entities.Coupon) error {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	row := r.db.QueryRowContext(queryCtx, "INSERT INTO coupons VALUES($1, $2, $3, $4, $5)", coupon.ID, coupon.Name, coupon.Description, coupon.Price, coupon.Level)
	if row.Err() != nil {
		return fmt.Errorf("query row context order failed: %w", row.Err())
	}

	row = r.db.QueryRowContext(queryCtx, "INSERT INTO media VALUES($1, $2, $3)", coupon.Media.ID, coupon.ID, coupon.Media.Path)
	if row.Err() != nil {
		return fmt.Errorf("query row context failed: %w", row.Err())
	}

	return nil
}

func (r Repo) GetCoupons(ctx context.Context) ([]entities.Coupon, error) {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(queryCtx, "SELECT * FROM coupons")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entities.ErrNoAnyCoupons
		}
		return nil, fmt.Errorf("query context failed: %w", err)
	}
	defer rows.Close()

	coupons := make([]entities.Coupon, 0)

	for rows.Next() {
		coupon := entities.NewCoupon()
		err := rows.Scan(&coupon.ID, &coupon.Name, &coupon.Description, &coupon.Price, &coupon.Level)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}

		media, err := r.getMyMedia(ctx, coupon.ID)
		if err != nil {
			return nil, fmt.Errorf("get media failed: %w", err)
		}

		coupon.Media = media

		coupons = append(coupons, coupon)
	}

	return coupons, nil
}

func (r Repo) getMyMedia(ctx context.Context, id string) (entities.Media, error) {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(queryCtx, "SELECT m.id, m.path FROM media m INNER JOIN coupons c ON m.coupon_id = c.id WHERE c.id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.NewMedia(), entities.ErrNoMedia
		}

		return entities.NewMedia(), fmt.Errorf("query context media failed: %w", err)
	}
	rows.Next()
	media := entities.NewMedia()
	err = rows.Scan(&media.ID, &media.Path)
	if err != nil {
		return entities.NewMedia(), fmt.Errorf("scan failed: %w", err)
	}

	return media, nil
}
