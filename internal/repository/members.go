package repository

import (
	"context"
	"fmt"
	"github.com/RipperAcskt/coupon-shop-admin/internal/entities"
	"strings"
	"time"
)

type transferMember struct {
	Email      *string
	FirstName  *string
	SecondName *string
	Role       *entities.Role
}

func (r Repo) AddMembers(ctx context.Context, members []entities.Member) error {
	queryContext, cancel := context.WithTimeout(ctx, 2*time.Second)
	var errResponse error
	defer cancel()
	for i := range members {
		_, err := r.db.ExecContext(queryContext, "INSERT INTO members VALUES ($1, $2, $3, $4, $5, $6)",
			members[i].ID, members[i].Email, members[i].FirstName, members[i].SecondName, members[i].OrganizationID, members[i].Role)
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

func (r Repo) DeleteMembers(ctx context.Context, members []entities.Member, organizationID string) error {
	queryContext, cancel := context.WithTimeout(ctx, 2*time.Second)
	var errResponse error
	defer cancel()
	for i := range members {
		res, err := r.db.ExecContext(queryContext, "DELETE FROM members WHERE email = $1", members[i].Email)
		if err != nil {
			return fmt.Errorf("exec context failed: %w", err)
		}

		num, err := res.RowsAffected()
		if err != nil {
			return fmt.Errorf("rows affected failed: %w", err)
		}
		if num == 0 {
			if errResponse == nil {
				errResponse = fmt.Errorf("delete of user with email %s blocked : %s", members[i].Email, "user doesn't exists")
			} else {
				errResponse = fmt.Errorf("%w, delete of user with email %s blocked : %s", errResponse, members[i].Email, "user doesn't exists")

			}
		}

		res, err = r.db.ExecContext(queryContext, "DELETE FROM users WHERE email = $1", members[i].Email)
		if err != nil {
			return fmt.Errorf("exec context failed: %w", err)
		}

		num, err = res.RowsAffected()
		if err != nil {
			return fmt.Errorf("rows affected failed: %w", err)
		}
		if num == 0 {
			if errResponse == nil {
				errResponse = fmt.Errorf("delete of user with email %s blocked : %s", members[i].Email, "user doesn't exists")
			} else {
				errResponse = fmt.Errorf("%w, delete of user with email %s blocked : %s", errResponse, members[i].Email, "user doesn't exists")

			}
		}
	}
	return errResponse
}

func (r Repo) UpdateMembers(ctx context.Context, members []entities.Member, organizationID string) error {
	queryContext, cancel := context.WithTimeout(ctx, 3*time.Second)
	var errResponse error
	defer cancel()
	trMembers := make([]transferMember, len(members))
	for i := range members {
		if members[i].FirstName != "" {
			trMembers[i].FirstName = &members[i].FirstName
		}
		if members[i].SecondName != "" {
			trMembers[i].SecondName = &members[i].SecondName
		}
		if members[i].Role != "" {
			trMembers[i].Role = &members[i].Role
		}
		//if v.Email != "" {
		//	trMembers[i].FirstName = &v.FirstName
		//}
	}
	fmt.Println("organizationID : ", organizationID)
	for i := range trMembers {
		//fmt.Println(*trMembers[i].Email)
		res, err := r.db.ExecContext(queryContext, "UPDATE members SET first_name = COALESCE($1, first_name), second_name = COALESCE($2, second_name), role = COALESCE($3, role) WHERE email = $4 and organization_id = $5",
			trMembers[i].FirstName, trMembers[i].SecondName, trMembers[i].Role, members[i].Email, organizationID)
		if err != nil {
			return err
		}
		num, err := res.RowsAffected()
		if err != nil {
			return fmt.Errorf("rows affected failed: %w", err)
		}
		if num == 0 {
			if errResponse == nil {
				errResponse = fmt.Errorf("updating of user with email %s blocked : %s", members[i].Email, "user doesn't exists")
			} else {
				errResponse = fmt.Errorf("%w, updating of user with email %s blocked : %s", errResponse, members[i].Email, "user doesn't exists")

			}
		}
	}
	return errResponse
}
