package service

import (
	"context"
	"fmt"
	"github.com/RipperAcskt/coupon-shop-admin/config"
	"github.com/RipperAcskt/coupon-shop-admin/internal/entities"
)

type CouponService struct {
	repo CouponRepoInterface
	cfg  config.Config
}

type CouponRepoInterface interface {
	CreateCoupon(ctx context.Context, coupon entities.Coupon) error
	GetCoupons(ctx context.Context) ([]entities.Coupon, error)
}

func NewCouponService(repo CouponRepoInterface, cfg config.Config) CouponService {
	return CouponService{
		repo: repo,
		cfg:  cfg,
	}
}

func (svc CouponService) CreateCoupon(ctx context.Context, coupon entities.Coupon) error {
	err := svc.repo.CreateCoupon(ctx, coupon)
	return err
}

func (svc CouponService) GetCoupons(ctx context.Context) ([]entities.Coupon, error) {
	coupons, err := svc.repo.GetCoupons(ctx)
	if err != nil {
		return nil, fmt.Errorf("get coupons failed: %w", err)
	}

	for i, coupon := range coupons {
		coupons[i].ContentUrl = svc.cfg.ServerHost + coupon.Media.Path
	}
	return coupons, err
}
