package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/RipperAcskt/coupon-shop-admin/internal/entities"
	"time"
)

func (r Repo) CreateCategory(ctx context.Context, category entities.Category) error {
	queryContext, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var id string
	err := r.db.QueryRowContext(queryContext, "SELECT id FROM categories WHERE name = $1", category.Name).Scan(&id)
	if err == nil {
		return entities.ErrCategoryAlreadyExists
	}

	_, err = r.db.ExecContext(queryContext, "INSERT INTO categories VALUES ($1, $2)", category.Id, category.Name)
	if err != nil {
		return fmt.Errorf("exec context failed: %w", err)
	}

	return nil
}

func (r Repo) GetCategories(ctx context.Context) ([]entities.Category, error) {
	queryContext, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(queryContext, "SELECT * FROM categories")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entities.ErrNoAnyCategory
		}
		return nil, fmt.Errorf("query context failed: %w", err)
	}
	defer rows.Close()

	categories := make([]entities.Category, 0)
	for rows.Next() {
		category := entities.Category{}
		err := rows.Scan(&category.Id, &category.Name)
		if err != nil {
			return nil, fmt.Errorf("rows scan failed: %w", err)
		}

		categories = append(categories, category)
	}
	return categories, nil
}

func (r Repo) GetCategory(ctx context.Context, id string) (entities.Category, error) {
	queryContext, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	row := r.db.QueryRowContext(queryContext, "SELECT * FROM categories WHERE id = $1", id)
	if row.Err() != nil {
		if errors.Is(row.Err(), sql.ErrNoRows) {
			return entities.Category{}, entities.ErrCategoryDoesNotExist
		}
		return entities.Category{}, fmt.Errorf("query row context failed: %w", row.Err())
	}

	category := entities.Category{}
	err := row.Scan(&category.Id, &category.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.Category{}, entities.ErrCategoryDoesNotExist
		}
		return entities.Category{}, fmt.Errorf("row scan failed: %w", err)
	}
	return category, nil
}

func (r Repo) UpdateCategory(ctx context.Context, id, name string) error {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := r.db.ExecContext(queryCtx, "UPDATE categories SET name = $1 WHERE id = $2", name, id)
	if err != nil {
		return fmt.Errorf("exec context failed: %w", err)
	}

	num, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected failed: %w", err)
	}
	if num == 0 {
		return entities.ErrCategoryDoesNotExist
	}
	return nil
}

func (r Repo) DeleteCategory(ctx context.Context, id string) error {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := r.db.ExecContext(queryCtx, "DELETE FROM categories WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("exec context failed: %w", err)
	}

	num, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected failed: %w", err)
	}
	if num == 0 {
		return entities.ErrCategoryDoesNotExist
	}
	return nil
}
