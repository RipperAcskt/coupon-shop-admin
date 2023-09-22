package service

import "github.com/RipperAcskt/coupon-shop-admin/config"

type Service struct {
	SubscriptionService
	OrganizationService
	CouponService
}

type Repo interface {
	SubscriptionRepoInterface
	OrganizationRepoInterface
	CouponRepoInterface
}

func New(repo Repo, cfg config.Config) Service {
	return Service{
		SubscriptionService: NewSubscriptionService(repo),
		OrganizationService: NewOrganizationService(repo),
		CouponService:       NewCouponService(repo, cfg),
	}
}
