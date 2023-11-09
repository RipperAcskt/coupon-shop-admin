package entities

import "time"

type UserEntity struct {
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Roles            []Role
	ID               string
	Email            string
	Phone            string
	Code             string
	OrganizationID   *int64
	Subscription     *string
	SubscriptionTime time.Time
}
