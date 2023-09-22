package service

import (
	"context"

	"github.com/RipperAcskt/coupon-shop-admin/internal/entities"
)

type SubscriptionService struct {
	repo SubscriptionRepoInterface
}

type SubscriptionRepoInterface interface {
	CreateSubscription(ctx context.Context, sub entities.Subscription) error
	GetSubscriptions(ctx context.Context) ([]entities.Subscription, error)
	GetSubscription(ctx context.Context, id string) (entities.Subscription, error)
	UpdateSubscription(ctx context.Context, id string, subscription entities.Subscription) error
	DeleteSubscription(ctx context.Context, id string) error
}

func NewSubscriptionService(repo SubscriptionRepoInterface) SubscriptionService {
	return SubscriptionService{repo: repo}
}

func (svc SubscriptionService) CreateSubscription(ctx context.Context, sub entities.Subscription) error {
	err := svc.repo.CreateSubscription(ctx, sub)
	return err
}

func (svc SubscriptionService) GetSubscriptions(ctx context.Context) ([]entities.Subscription, error) {
	subs, err := svc.repo.GetSubscriptions(ctx)
	return subs, err
}

func (svc SubscriptionService) GetSubscription(ctx context.Context, id string) (entities.Subscription, error) {
	sub, err := svc.repo.GetSubscription(ctx, id)
	return sub, err
}

func (svc SubscriptionService) UpdateSubscription(ctx context.Context, id string, subscription entities.Subscription) error {
	err := svc.repo.UpdateSubscription(ctx, id, subscription)
	return err
}

func (svc SubscriptionService) DeleteSubscription(ctx context.Context, id string) error {
	err := svc.repo.DeleteSubscription(ctx, id)
	return err
}
