package service

import (
	"context"
	"fmt"
	"github.com/RipperAcskt/coupon-shop-admin/internal/entities"
	"github.com/google/uuid"
)

type RegionService struct {
	repo RegionRepoInterface
}

type RegionRepoInterface interface {
	CreateRegion(ctx context.Context, region entities.Region) error
	GetRegions(ctx context.Context) ([]entities.Region, error)
	UpdateRegion(ctx context.Context, id string, region entities.Region) error
	DeleteRegion(ctx context.Context, id string) error
}

func NewRegionService(repo RegionRepoInterface) RegionService {
	return RegionService{
		repo: repo,
	}
}

func (svc RegionService) CreateRegion(ctx context.Context, region entities.Region) error {
	region.Id = uuid.NewString()
	err := svc.repo.CreateRegion(ctx, region)
	return err
}

func (svc RegionService) GetRegions(ctx context.Context) ([]entities.Region, error) {
	regions, err := svc.repo.GetRegions(ctx)
	if err != nil {
		return nil, fmt.Errorf("get coupons failed: %w", err)
	}
	return regions, err
}

func (svc RegionService) UpdateRegion(ctx context.Context, id string, region entities.Region) error {
	err := svc.repo.UpdateRegion(ctx, id, region)
	return err
}

func (svc RegionService) DeleteRegion(ctx context.Context, id string) error {
	err := svc.repo.DeleteRegion(ctx, id)
	return err
}
