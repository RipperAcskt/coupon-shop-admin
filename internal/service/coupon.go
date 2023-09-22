package service

import (
	"context"
	"github.com/RipperAcskt/coupon-shop-admin/internal/entities"
)

type CouponService struct {
	repo CouponRepoInterface
}

type CouponRepoInterface interface {
	CreateCoupon(ctx context.Context, coupon entities.Coupon) error
}

func NewCouponService(repo CouponRepoInterface) CouponService {
	return CouponService{repo: repo}
}

func (svc CouponService) CreateCoupon(ctx context.Context, coupon entities.Coupon) error {
	err := svc.repo.CreateCoupon(ctx, coupon)
	return err
}
