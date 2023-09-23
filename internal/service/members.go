package service

import (
	"context"
	"github.com/RipperAcskt/coupon-shop-admin/internal/entities"
)

type MembersService struct {
	repo MembersRepoInterface
}

type MembersRepoInterface interface {
	AddMembers(ctx context.Context, members []entities.Member, id string) error
	GetMembers(ctx context.Context, organizationID string) ([]entities.Member, error)
	DeleteMembers(ctx context.Context, membersToDelete []entities.Member) error
}

func NewMembersService(repo MembersRepoInterface) MembersService {
	return MembersService{repo: repo}
}

func (svc OrganizationService) AddMembers(ctx context.Context, org entities.Organization) error {
	err := svc.repo.CreateOrganization(ctx, org)
	return err
}
