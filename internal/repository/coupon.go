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

type transferCoupon struct {
	Name        *string
	Description *string
	Price       *int
	Percent     *int
	Region      *string
	Category    *string
	Subcategory *string
}

func (r Repo) CreateCoupon(ctx context.Context, coupon entities.Coupon) error {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	row := r.db.QueryRowContext(queryCtx, "SELECT id FROM regions WHERE name = $1", coupon.Region)
	if row.Err() != nil {
		return fmt.Errorf("query row context region failed: %w", row.Err())
	}

	var id string
	err := row.Scan(&id)
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	row = r.db.QueryRowContext(queryCtx, "SELECT id FROM organization WHERE name = $1", coupon.Organization)
	if row.Err() != nil {
		return fmt.Errorf("query row context region failed: %w", row.Err())
	}

	var org_id string
	err = row.Scan(&org_id)
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	row = r.db.QueryRowContext(queryCtx, "INSERT INTO coupons VALUES($1, $2, $3, $4, $5, $6, $7, $8)", coupon.ID, coupon.Name, coupon.Description, coupon.Price, coupon.Percent, coupon.Level, id, org_id)
	if row.Err() != nil {
		return fmt.Errorf("query row context coupon failed: %w", row.Err())
	}

	row = r.db.QueryRowContext(queryCtx, "INSERT INTO media VALUES($1, $2, $3)", coupon.Media.ID, coupon.ID, coupon.Media.Path)
	if row.Err() != nil {
		return fmt.Errorf("query row context failed: %w", row.Err())
	}

	var idCategory, idSubcategory string

	row = r.db.QueryRowContext(queryCtx, "SELECT id FROM categories WHERE name = $1", coupon.Category)
	if row.Err() != nil {
		return fmt.Errorf("query row context region failed: %w", row.Err())
	}
	err = row.Scan(&idCategory)
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	if coupon.Subcategory != nil {
		if *coupon.Subcategory != "" {
			row = r.db.QueryRowContext(queryCtx, "SELECT id FROM subcategories WHERE name = $1", coupon.Subcategory)
			if row.Err() != nil {
				return fmt.Errorf("query row context region failed: %w", row.Err())
			}
			err = row.Scan(&idSubcategory)
			if err != nil {
				return fmt.Errorf("scan failed: %w", err)
			}

			row = r.db.QueryRowContext(queryCtx, "INSERT INTO categories_coupons VALUES($1, $2, $3)", idCategory, idSubcategory, coupon.ID)
			if row.Err() != nil {
				return fmt.Errorf("query row context failed: %w", row.Err())
			}

			return nil
		}
	}

	row = r.db.QueryRowContext(queryCtx, "INSERT INTO categories_coupons (id_category, id_coupon) VALUES($1, $2)", idCategory, coupon.ID)
	if row.Err() != nil {
		return fmt.Errorf("query row context failed: %w", row.Err())
	}

	return nil
}

func (r Repo) GetCoupons(ctx context.Context) ([]entities.Coupon, error) {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(queryCtx, "WITH regionInfo(coupon_id, region_name) AS (SELECT coupons.id, regions.name FROM coupons JOIN regions ON coupons.region = regions.id), categoryInfo(coupon_id, category_name, subcategory_name) AS (SELECT categories_coupons.id_coupon, categories.name, subcategories.name FROM categories_coupons JOIN categories ON categories_coupons.id_category = categories.id left JOIN subcategories ON categories_coupons.id_subcategory = subcategories.id), organizationInfo(coupon_id, org_name) AS (SELECT coupons.id, organization.name FROM coupons JOIN organization ON coupons.organization_id = organization.id) SELECT coupons.id, coupons.name, coupons.description, coupons.price, coupons.percent, coupons.level, regionInfo.region_name, categoryInfo.category_name, categoryInfo.subcategory_name, organizationInfo.org_name FROM coupons JOIN regionInfo ON regionInfo.coupon_id = coupons.id JOIN categoryInfo ON categoryInfo.coupon_id = coupons.id JOIN organizationInfo ON organizationInfo.coupon_id = coupons.id")
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
		err := rows.Scan(&coupon.ID, &coupon.Name, &coupon.Description, &coupon.Price, &coupon.Percent, &coupon.Level, &coupon.Region, &coupon.Category, &coupon.Subcategory, &coupon.Organization)
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

func (r Repo) GetCouponsSearch(ctx context.Context, s string) ([]entities.Coupon, error) {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	q := fmt.Sprintf(`WITH regionInfo(coupon_id, region_name) AS (SELECT coupons.id, regions.name FROM coupons JOIN regions ON coupons.region = regions.id), categoryInfo(coupon_id, category_name, subcategory_name) AS (SELECT categories_coupons.id_coupon, categories.name, subcategories.name FROM categories_coupons JOIN categories ON categories_coupons.id_category = categories.id left JOIN subcategories ON categories_coupons.id_subcategory = subcategories.id), organizationInfo(coupon_id, org_name) AS (SELECT coupons.id, organization.name FROM coupons JOIN organization ON coupons.organization_id = organization.id) SELECT coupons.id, coupons.name, coupons.description, coupons.price, coupons.percent, coupons.level, regionInfo.region_name, categoryInfo.category_name, categoryInfo.subcategory_name, organizationInfo.org_name FROM coupons JOIN regionInfo ON regionInfo.coupon_id = coupons.id JOIN categoryInfo ON categoryInfo.coupon_id = coupons.id JOIN organizationInfo ON organizationInfo.coupon_id = coupons.id WHERE coupons.name SIMILAR TO '%s%%'`, s)
	rows, err := r.db.QueryContext(queryCtx, q)
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
		err := rows.Scan(&coupon.ID, &coupon.Name, &coupon.Description, &coupon.Price, &coupon.Percent, &coupon.Level, &coupon.Region, &coupon.Category, &coupon.Subcategory, &coupon.Organization)
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

func (r Repo) GetCouponsByRegion(ctx context.Context, region string) ([]entities.Coupon, error) {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(queryCtx, "WITH regionInfo(coupon_id, region_name) AS (SELECT coupons.id, regions.name FROM coupons JOIN regions ON coupons.region = regions.id), categoryInfo(coupon_id, category_name, subcategory_name) AS (SELECT categories_coupons.id_coupon, categories.name, subcategories.name FROM categories_coupons JOIN categories ON categories_coupons.id_category = categories.id left JOIN subcategories ON categories_coupons.id_subcategory = subcategories.id), organizationInfo(coupon_id, org_name) AS (SELECT coupons.id, organization.name FROM coupons JOIN organization ON coupons.organization_id = organization.id) SELECT coupons.id, coupons.name, coupons.description, coupons.price, coupons.percent, coupons.level, regionInfo.region_name, categoryInfo.category_name, categoryInfo.subcategory_name, organizationInfo.org_name FROM coupons JOIN regionInfo ON regionInfo.coupon_id = coupons.id JOIN categoryInfo ON categoryInfo.coupon_id = coupons.id JOIN organizationInfo ON organizationInfo.coupon_id = coupons.id WHERE regionInfo.region_name = $1", region)
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
		err := rows.Scan(&coupon.ID, &coupon.Name, &coupon.Description, &coupon.Price, &coupon.Percent, &coupon.Level, &coupon.Region, &coupon.Category, &coupon.Subcategory, &coupon.Organization)
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

func (r Repo) GetCouponsByCategory(ctx context.Context, category string) ([]entities.Coupon, error) {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(queryCtx, "WITH regionInfo(coupon_id, region_name) AS (SELECT coupons.id, regions.name FROM coupons JOIN regions ON coupons.region = regions.id), categoryInfo(coupon_id, category_name, subcategory_name) AS (SELECT categories_coupons.id_coupon, categories.name, subcategories.name FROM categories_coupons JOIN categories ON categories_coupons.id_category = categories.id left JOIN subcategories ON categories_coupons.id_subcategory = subcategories.id), organizationInfo(coupon_id, org_name) AS (SELECT coupons.id, organization.name FROM coupons JOIN organization ON coupons.organization_id = organization.id) SELECT coupons.id, coupons.name, coupons.description, coupons.price, coupons.percent, coupons.level, regionInfo.region_name, categoryInfo.category_name, categoryInfo.subcategory_name, organizationInfo.org_name FROM coupons JOIN regionInfo ON regionInfo.coupon_id = coupons.id JOIN categoryInfo ON categoryInfo.coupon_id = coupons.id JOIN organizationInfo ON organizationInfo.coupon_id = coupons.id WHERE categoryInfo.category_name = $1", category)
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
		err := rows.Scan(&coupon.ID, &coupon.Name, &coupon.Description, &coupon.Price, &coupon.Percent, &coupon.Level, &coupon.Region, &coupon.Category, &coupon.Subcategory, &coupon.Organization)
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

func (r Repo) GetCouponsBySubcategory(ctx context.Context, category string) ([]entities.Coupon, error) {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(queryCtx, "WITH regionInfo(coupon_id, region_name) AS (SELECT coupons.id, regions.name FROM coupons JOIN regions ON coupons.region = regions.id), categoryInfo(coupon_id, category_name, subcategory_name) AS (SELECT categories_coupons.id_coupon, categories.name, subcategories.name FROM categories_coupons JOIN categories ON categories_coupons.id_category = categories.id left JOIN subcategories ON categories_coupons.id_subcategory = subcategories.id), organizationInfo(coupon_id, org_name) AS (SELECT coupons.id, organization.name FROM coupons JOIN organization ON coupons.organization_id = organization.id) SELECT coupons.id, coupons.name, coupons.description, coupons.price, coupons.percent, coupons.level, regionInfo.region_name, categoryInfo.category_name, categoryInfo.subcategory_name, organizationInfo.org_name FROM coupons JOIN regionInfo ON regionInfo.coupon_id = coupons.id JOIN categoryInfo ON categoryInfo.coupon_id = coupons.id JOIN organizationInfo ON organizationInfo.coupon_id = coupons.id WHERE categoryInfo.subcategory_name = $1", category)
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
		err := rows.Scan(&coupon.ID, &coupon.Name, &coupon.Description, &coupon.Price, &coupon.Percent, &coupon.Level, &coupon.Region, &coupon.Category, &coupon.Subcategory, &coupon.Organization)
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

func (r Repo) GetCoupon(ctx context.Context, id string) (entities.Coupon, error) {
	queryContext, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	row := r.db.QueryRowContext(queryContext, "WITH regionInfo(coupon_id, region_name) AS (SELECT coupons.id, regions.name FROM coupons JOIN regions ON coupons.region = regions.id), categoryInfo(coupon_id, category_name, subcategory_name) AS (SELECT categories_coupons.id_coupon, categories.name, subcategories.name FROM categories_coupons JOIN categories ON categories_coupons.id_category = categories.id left JOIN subcategories ON categories_coupons.id_subcategory = subcategories.id), organizationInfo(coupon_id, org_name) AS (SELECT coupons.id, organization.name FROM coupons JOIN organization ON coupons.organization_id = organization.id) SELECT coupons.id, coupons.name, coupons.description, coupons.price, coupons.percent, coupons.level, regionInfo.region_name, categoryInfo.category_name, categoryInfo.subcategory_name, organizationInfo.org_name FROM coupons JOIN regionInfo ON regionInfo.coupon_id = coupons.id JOIN categoryInfo ON categoryInfo.coupon_id = coupons.id JOIN organizationInfo ON organizationInfo.coupon_id = coupons.id WHERE coupons.id = $1", id)
	if row.Err() != nil {
		if errors.Is(row.Err(), sql.ErrNoRows) {
			return entities.NewCoupon(), entities.ErrSubscriptionDoesNotExist
		}
		return entities.NewCoupon(), fmt.Errorf("query row context failed: %w", row.Err())
	}

	coupon := entities.NewCoupon()
	err := row.Scan(&coupon.ID, &coupon.Name, &coupon.Description, &coupon.Price, &coupon.Percent, &coupon.Level, &coupon.Region, &coupon.Category, &coupon.Subcategory, &coupon.Organization)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.NewCoupon(), entities.ErrCouponDoesNotExist
		}
		return entities.NewCoupon(), fmt.Errorf("row scan failed: %w", err)
	}

	media, err := r.getMyMedia(ctx, coupon.ID)
	if err != nil {
		return entities.NewCoupon(), fmt.Errorf("get media failed: %w", err)
	}

	coupon.Media = media

	return coupon, nil
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

func (r Repo) UpdateCoupon(ctx context.Context, id string, coupon entities.Coupon) error {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var transfer transferCoupon
	if coupon.Name != "" {
		transfer.Name = &coupon.Name
	}
	if coupon.Description != "" {
		transfer.Description = &coupon.Description
	}
	if coupon.Price != 0 {
		transfer.Price = &coupon.Price
	}
	if coupon.Percent != 0 {
		transfer.Percent = &coupon.Percent
	}
	if coupon.Region != "" {
		transfer.Region = &coupon.Region
	}
	if coupon.Category != "" {
		transfer.Category = &coupon.Category
	}
	if *coupon.Subcategory != "" {
		transfer.Subcategory = coupon.Subcategory
	}

	res, err := r.db.ExecContext(queryCtx, "UPDATE coupons SET name = COALESCE($1, name), description = COALESCE($2, description), price = COALESCE($3, price), percent = COALESCE($4, percent), region = COALESCE($5, region) WHERE id = $6",
		transfer.Name, transfer.Description, transfer.Price, transfer.Percent, transfer.Region, id)
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

	res, err = r.db.ExecContext(queryCtx, "UPDATE categories_coupons SET id_category = COALESCE($1, id_category), id_subcategory = COALESCE($2, id_subcategory) WHERE id_coupon = $3",
		transfer.Category, transfer.Subcategory, id)
	if err != nil {
		return fmt.Errorf("exec context failed: %w", err)
	}

	num, err = res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected failed: %w", err)
	}
	if num == 0 {
		return entities.ErrCouponDoesNotExist
	}

	if coupon.Media.Path != "" {
		res, err := r.db.ExecContext(queryCtx, "UPDATE media SET path = COALESCE($1, path) WHERE coupon_id = $2",
			coupon.Media.Path, id)
		if err != nil {
			return fmt.Errorf("exec context media failed: %w", err)
		}

		num, err := res.RowsAffected()
		if err != nil {
			return fmt.Errorf("rows affected media failed: %w", err)
		}
		if num == 0 {
			return entities.ErrCouponDoesNotExist
		}
	}
	return nil
}

func (r Repo) DeleteCoupon(ctx context.Context, id string) error {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := r.db.ExecContext(queryCtx, "DELETE FROM coupons WHERE id = $1", id)
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
