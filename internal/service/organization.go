package service

import (
	"context"
	"github.com/RipperAcskt/coupon-shop-admin/internal/entities"
)

type OrganizationService struct {
	repo OrganizationRepoInterface
}

type OrganizationRepoInterface interface {
	CreateOrganization(ctx context.Context, org entities.Organization) error
	GetOrganizations(ctx context.Context) ([]entities.Organization, error)
	DeleteOrganization(ctx context.Context, id string) error
	GetOrganization(ctx context.Context, organizationID string) (entities.Organization, error)
}

func NewOrganizationService(repo OrganizationRepoInterface) OrganizationService {
	return OrganizationService{repo: repo}
}

func (svc OrganizationService) CreateOrganization(ctx context.Context, org entities.Organization) error {
	err := svc.repo.CreateOrganization(ctx, org)
	return err
}

func (svc OrganizationService) GetOrganizations(ctx context.Context) ([]entities.Organization, error) {
	orgs, err := svc.repo.GetOrganizations(ctx)
	if err != nil {
		return nil, err
	}
	for i, org := range orgs {
		orgs[i].ContentUrl = "http://parcus.shop/admin" + org.OrgImage.Path
	}
	return orgs, err
}

func (svc OrganizationService) DeleteOrganization(ctx context.Context, id string) error {
	err := svc.repo.DeleteOrganization(ctx, id)
	return err
}

func (svc OrganizationService) GetOrganization(ctx context.Context, organizationID string) (entities.Organization, error) {
	org, err := svc.repo.GetOrganization(ctx, organizationID)
	if err != nil {
		return entities.Organization{}, err
	}
	org.ContentUrl = "http://parcus.shop/admin" + org.OrgImage.Path
	return org, nil
}
