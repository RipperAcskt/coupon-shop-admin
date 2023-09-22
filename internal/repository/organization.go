package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/RipperAcskt/coupon-shop-admin/internal/entities"
	"time"
)

func (r Repo) CreateOrganization(ctx context.Context, org entities.Organization) error {
	queryContext, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var id string
	err := r.db.QueryRowContext(queryContext, "SELECT id FROM organization WHERE email_adimin = $1", org.EmailAdmin).Scan(&id)
	if err == nil {
		return entities.ErrOrganizationAlreadyExists
	}
	//TODO : firstly we should check, if the level is listed in tables subscription
	_, err = r.db.ExecContext(queryContext, "INSERT INTO organization VALUES ($1, $2, $3, $4)",
		org.ID, org.Name, org.EmailAdmin, org.LevelSubscription)
	if err != nil {
		return fmt.Errorf("exec context failed: %w", err)
	}

	return nil
}

func (r Repo) GetOrganizations(ctx context.Context) ([]entities.Organization, error) {
	queryContext, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(queryContext, "SELECT * FROM organization")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entities.ErrNoAnyOrganization
		}
		return nil, fmt.Errorf("query context failed: %w", err)
	}
	defer rows.Close()

	orgs := make([]entities.Organization, 0)
	for rows.Next() {
		org := entities.NewOrganization()
		err := rows.Scan(&org.ID, &org.Name, &org.EmailAdmin, &org.LevelSubscription)
		if err != nil {
			return nil, fmt.Errorf("rows scan failed: %w", err)
		}

		orgs = append(orgs, org)
	}
	return orgs, nil
}

func (r Repo) DeleteOrganization(ctx context.Context, id string) error {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := r.db.ExecContext(queryCtx, "DELETE FROM organization WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("exec context failed: %w", err)
	}

	num, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected failed: %w", err)
	}
	if num == 0 {
		return entities.ErrOrganizationnDoesNotExist
	}
	return nil
}
