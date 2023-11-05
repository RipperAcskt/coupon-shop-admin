package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/RipperAcskt/coupon-shop-admin/internal/entities"
	"time"
)

func (r Repo) CreateLink(ctx context.Context, link entities.Link) error {
	queryContext, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var id string
	err := r.db.QueryRowContext(queryContext, "SELECT id FROM regions WHERE name = $1", link.Region).Scan(&id)
	if err != nil {
		return fmt.Errorf("query context failed: %w", err)
	}

	_, err = r.db.ExecContext(queryContext, "INSERT INTO links VALUES ($1, $2, $3, $4)", link.Id, link.Name, link.Link, id)
	if err != nil {
		return fmt.Errorf("exec context failed: %w", err)
	}
	return nil
}

func (r Repo) GetLinks(ctx context.Context) ([]entities.Link, error) {
	queryContext, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(queryContext, "SELECT links.id, links.name, links.link, regions.name FROM links JOIN regions ON links.region_id = regions.id")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entities.ErrNoAnyLink
		}
		return nil, fmt.Errorf("query context failed: %w", err)
	}
	defer rows.Close()

	links := make([]entities.Link, 0)
	for rows.Next() {
		link := entities.Link{}
		err := rows.Scan(&link.Id, &link.Name, &link.Link, link.Region)
		if err != nil {
			return nil, fmt.Errorf("rows scan failed: %w", err)
		}

		links = append(links, link)
	}
	return links, nil
}

func (r Repo) GetLinkByRegion(ctx context.Context, region string) (entities.Link, error) {
	queryContext, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	row := r.db.QueryRowContext(queryContext, "SELECT links.id, links.name, links.link, regions.name FROM links JOIN regions ON links.region_id = regions.id WHERE regions.name = $1", region)
	if row.Err() != nil {
		if errors.Is(row.Err(), sql.ErrNoRows) {
			return entities.Link{}, entities.ErrLinkDoesNotExist
		}
		return entities.Link{}, fmt.Errorf("query row context failed: %w", row.Err())
	}

	link := entities.Link{}
	err := row.Scan(&link.Id, &link.Name, &link.Link, link.Region)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.Link{}, entities.ErrLinkDoesNotExist
		}
		return entities.Link{}, fmt.Errorf("row scan failed: %w", err)
	}
	return link, nil
}

func (r Repo) UpdateLink(ctx context.Context, id, link string) error {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := r.db.ExecContext(queryCtx, "UPDATE links SET link = $1 WHERE id = $2", link, id)
	if err != nil {
		return fmt.Errorf("exec context failed: %w", err)
	}

	num, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected failed: %w", err)
	}
	if num == 0 {
		return entities.ErrLinkDoesNotExist
	}
	return nil
}

func (r Repo) DeleteLink(ctx context.Context, id string) error {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := r.db.ExecContext(queryCtx, "DELETE FROM links WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("exec context failed: %w", err)
	}

	num, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected failed: %w", err)
	}
	if num == 0 {
		return entities.ErrLinkDoesNotExist
	}
	return nil
}
