package repository

import (
	"context"
	"fmt"
	"github.com/RipperAcskt/coupon-shop-admin/internal/entities"
	"strings"
	"time"
)

func (r Repo) AddMembers(ctx context.Context, members []entities.Member) error {
	queryContext, cancel := context.WithTimeout(ctx, 2*time.Second)
	var errResponse error
	defer cancel()
	for i := range members {
		_, err := r.db.ExecContext(queryContext, "INSERT INTO members VALUES ($1, $2, $3, $4, $5)",
			members[i].ID, members[i].Email, members[i].FirstName, members[i].SecondName, members[i].OrganizationID)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				if errResponse == nil {
					errResponse = fmt.Errorf("creation of user with email %s blocked : %s", members[i].Email, "user already exists")
				} else {
					errResponse = fmt.Errorf("%w, creation of user with email %s blocked : %s", errResponse, members[i].Email, "user already exists")

				}
			} else {
				return err
			}
		}
	}
	return errResponse
}
