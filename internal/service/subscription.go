package service

import (
	"context"

	"github.com/RipperAcskt/coupon-shop-admin/internal/entities"

	"github.com/google/uuid"
)

type SubscriptionService struct {
	repo SubscriptionRepoInterface
}

type SubscriptionRepoInterface interface {
	CreateSubscription(ctx context.Context, sub entities.Subscription) error
}

func NewSubscriptionService(repo SubscriptionRepoInterface) SubscriptionService {
	return SubscriptionService{repo: repo}
}

func (svc SubscriptionService) CreateSubscription(ctx context.Context, sub entities.Subscription) error {
	sub.ID = uuid.NewString()
	err := svc.repo.CreateSubscription(ctx, sub)
	return err
}
