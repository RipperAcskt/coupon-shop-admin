package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/RipperAcskt/coupon-shop-admin/internal/entities"
	"time"
)

type transerOrganization struct {
	Name              *string
	EmailAdmin        *string
	LevelSubscription *int
	ORGN              *string
	KPP               *string
	INN               *string
	Address           *string
}

func (r Repo) CreateOrganization(ctx context.Context, org entities.Organization) error {
	queryContext, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var id string
	err := r.db.QueryRowContext(queryContext, "SELECT id FROM organization WHERE email_adimin = $1", org.EmailAdmin).Scan(&id)
	if err == nil {
		return entities.ErrOrganizationAlreadyExists
	}
	subs, err := r.GetSubscriptions(ctx)
	if err != nil {
		return err
	}
	isLevelCorrect := false
	for _, v := range subs {
		if v.Level == org.LevelSubscription {
			isLevelCorrect = true
		}
	}
	if isLevelCorrect == false {
		return fmt.Errorf("error : level of subscription isn't correct")
	}
	_, err = r.db.ExecContext(queryContext, "INSERT INTO organization VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		org.ID, org.Name, org.EmailAdmin, org.LevelSubscription, org.ORGN, org.KPP, org.INN, org.Address)
	if err != nil {
		return fmt.Errorf("exec context failed: %w", err)
	}
	row := r.db.QueryRowContext(ctx, "INSERT INTO images VALUES($1, $2, $3)", org.OrgImage.ID, org.ID, org.OrgImage.Path)
	if row.Err() != nil {
		return fmt.Errorf("query row context failed: %w", row.Err())
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
		err := rows.Scan(&org.ID, &org.Name, &org.EmailAdmin, &org.LevelSubscription, &org.ORGN, &org.KPP, &org.INN, &org.Address)
		if err != nil {
			return nil, fmt.Errorf("rows scan failed: %w", err)
		}

		media, err := r.getMyImage(ctx, org.ID)
		if err != nil {
			return nil, fmt.Errorf("get media failed: %w", err)
		}
		org.OrgImage = media
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

func (r Repo) GetOrganization(ctx context.Context, organizationID string) (entities.Organization, error) {
	queryContext, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	row := r.db.QueryRowContext(queryContext, "SELECT * FROM organization WHERE id = $1", organizationID)
	if row.Err() != nil {
		if errors.Is(row.Err(), sql.ErrNoRows) {
			return entities.NewOrganization(), entities.ErrOrganizationnDoesNotExist
		}
		return entities.NewOrganization(), fmt.Errorf("query row context failed: %w", row.Err())
	}

	org := entities.NewOrganization()
	err := row.Scan(&org.ID, &org.Name, &org.EmailAdmin, &org.LevelSubscription, &org.ORGN, &org.KPP, &org.INN, &org.Address)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.NewOrganization(), entities.ErrSubscriptionDoesNotExist
		}
		return entities.NewOrganization(), fmt.Errorf("row scan failed: %w", err)
	}

	rows, err := r.db.QueryContext(queryContext, "SELECT * FROM members WHERE organization_id = $1", organizationID)
	if err != nil {
		return org, nil
	}
	defer rows.Close()

	members := make([]entities.Member, 0)
	for rows.Next() {
		member := entities.Member{}
		err := rows.Scan(&member.ID, &member.Email, &member.FirstName, &member.SecondName, &member.OrganizationID, &member.Role)
		if err != nil {
			return org, fmt.Errorf("rows scan failed: %w", err)
		}

		members = append(members, member)
	}
	org.Members = members

	image, err := r.getMyImage(ctx, org.ID)
	if err != nil {
		return org, fmt.Errorf("get media failed: %w", err)
	}

	org.OrgImage = image

	return org, nil
}

func (r Repo) getMyImage(ctx context.Context, id string) (entities.Image, error) {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(queryCtx, "SELECT m.id, m.path FROM images m INNER JOIN organization c ON m.organization_id = c.id WHERE c.id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.NewImage(), entities.ErrNoMedia
		}

		return entities.NewImage(), fmt.Errorf("query context media failed: %w", err)
	}
	rows.Next()
	media := entities.NewImage()
	err = rows.Scan(&media.ID, &media.Path)
	if err != nil {
		return entities.NewImage(), fmt.Errorf("scan failed: %w", err)
	}

	return media, nil
}

func (r Repo) UpdateOrganization(ctx context.Context, org entities.Organization, organizationID string) error {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var transfer transerOrganization
	if org.Name != "" {
		transfer.Name = &org.Name
	}
	if org.EmailAdmin != "" {
		transfer.EmailAdmin = &org.EmailAdmin
	}
	if org.LevelSubscription != 0 {
		transfer.LevelSubscription = &org.LevelSubscription
	}
	if org.INN != "" {
		transfer.INN = &org.INN
	}
	if org.ORGN != "" {
		transfer.ORGN = &org.ORGN
	}
	if org.KPP != "" {
		transfer.KPP = &org.KPP
	}
	if org.ORGN != "" {
		transfer.Address = &org.Address
	}
	res, err := r.db.ExecContext(queryCtx, "UPDATE organization SET name = COALESCE($1, name), email_admin = COALESCE($2, email_admin), level_subscription = COALESCE($3, level_subscription), orgn = COALESCE($4, orgn) , kpp = COALESCE($5, kpp), inn = COALESCE($6, inn), address = COALESCE($7, address) WHERE id = $8",
		transfer.Name, transfer.EmailAdmin, transfer.LevelSubscription, transfer.ORGN, transfer.KPP, transfer.INN, transfer.Address, organizationID)
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

	if org.OrgImage.Path != "" {
		res, err := r.db.ExecContext(queryCtx, "UPDATE images SET path = COALESCE($1, path) WHERE organization_id = $2",
			org.OrgImage.Path, organizationID)
		if err != nil {
			return fmt.Errorf("exec context media failed: %w", err)
		}

		num, err := res.RowsAffected()
		if err != nil {
			return fmt.Errorf("rows affected media failed: %w", err)
		}
		if num == 0 {
			return entities.ErrOrganizationnDoesNotExist
		}
	}
	return nil
}
