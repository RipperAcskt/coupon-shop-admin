package service

import (
	"context"
	"fmt"
	"github.com/RipperAcskt/coupon-shop-admin/config"
	"github.com/RipperAcskt/coupon-shop-admin/internal/entities"
	"os"
)

type CouponService struct {
	repo CouponRepoInterface
	cfg  config.Config
}

type CouponRepoInterface interface {
	CreateCoupon(ctx context.Context, coupon entities.Coupon) error
	GetCoupons(ctx context.Context) ([]entities.Coupon, error)
	GetCoupon(ctx context.Context, id string) (entities.Coupon, error)
	UpdateCoupon(ctx context.Context, id string, coupon entities.Coupon) error
	DeleteCoupon(ctx context.Context, id string) error
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
		coupons[i].ContentUrl = "http://parcus.shop/admin" + coupon.Media.Path
	}
	return coupons, err
}

func (svc CouponService) GetCoupon(ctx context.Context, id string) (entities.Coupon, error) {
	coupon, err := svc.repo.GetCoupon(ctx, id)
	if err != nil {
		return entities.NewCoupon(), fmt.Errorf("get coupon failed: %w", err)
	}

	coupon.ContentUrl = "http://parcus.shop/admin" + coupon.Media.Path

	return coupon, err
}

func (svc CouponService) UpdateCoupon(ctx context.Context, id string, coupon entities.Coupon) error {
	err := svc.repo.UpdateCoupon(ctx, id, coupon)
	return err
}

func (svc CouponService) DeleteCoupon(ctx context.Context, id string) error {
	err := os.Remove("./store/" + id + ".jpg")
	if err != nil {
		return fmt.Errorf("remove failed: %w", err)
	}
	err = svc.repo.DeleteCoupon(ctx, id)
	return err
}
