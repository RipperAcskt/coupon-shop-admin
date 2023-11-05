package service

import (
	"context"
	"github.com/google/uuid"

	"github.com/RipperAcskt/coupon-shop-admin/internal/entities"
)

type LinkService struct {
	repo LinkRepoInterface
}

type LinkRepoInterface interface {
	CreateLink(ctx context.Context, link entities.Link) error
	GetLinks(ctx context.Context) ([]entities.Link, error)
	GetLinkByRegion(ctx context.Context, region string) (entities.Link, error)
	UpdateLink(ctx context.Context, id, link string) error
	DeleteLink(ctx context.Context, id string) error
}

func NewLinkService(repo LinkRepoInterface) LinkService {
	return LinkService{repo: repo}
}

func (svc LinkService) CreateLink(ctx context.Context, link entities.Link) error {
	link.Id = uuid.NewString()
	err := svc.repo.CreateLink(ctx, link)
	return err
}

func (svc LinkService) GetLinks(ctx context.Context) ([]entities.Link, error) {
	subs, err := svc.repo.GetLinks(ctx)
	return subs, err
}

func (svc LinkService) GetLinkByRegion(ctx context.Context, region string) (entities.Link, error) {
	sub, err := svc.repo.GetLinkByRegion(ctx, region)
	return sub, err
}

func (svc LinkService) UpdateLink(ctx context.Context, id, link string) error {
	err := svc.repo.UpdateLink(ctx, id, link)
	return err
}

func (svc LinkService) DeleteLink(ctx context.Context, id string) error {
	err := svc.repo.DeleteLink(ctx, id)
	return err
}
