package entities

import "errors"

var (
	ErrSubscriptionAlreadyExists = errors.New("subscription already exists")
	ErrNoAnySubscription         = errors.New("there is not a single subscription")
	ErrSubscriptionDoesNotExist  = errors.New("subscription does not exist")
)

type Subscription struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Level       int    `json:"level"`
}

func NewSubscription() Subscription {
	return Subscription{}
}
