package repository

import (
	"context"
	"fmt"
	"github.com/RipperAcskt/coupon-shop-admin/internal/entities"
	"time"
)

func (r Repo) AddMembers(ctx context.Context, members []entities.Member) error {
	queryContext, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	for i := range members {
		_, err := r.db.ExecContext(queryContext, "INSERT INTO members VALUES ($1, $2, $3, $4, $5)",
			members[i].ID, members[i].Email, members[i].FirstName, members[i].SecondName, members[i].OrganizationID)
		if err != nil {
			return fmt.Errorf("exec context failed: %w", err)
		}

	}
	return nil
}
