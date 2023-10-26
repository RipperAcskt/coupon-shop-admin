package service

import (
	"context"
	"github.com/RipperAcskt/coupon-shop-admin/internal/entities"
	"github.com/google/uuid"
)

type MembersService struct {
	repo MembersRepoInterface
}

type MembersRepoInterface interface {
	AddMembers(ctx context.Context, members []entities.Member) error
	DeleteMembers(ctx context.Context, membersToDelete []entities.Member, organizationID string) error
	UpdateMembers(ctx context.Context, members []entities.Member, organizationID string) error
}

func NewMembersService(repo MembersRepoInterface) MembersService {
	return MembersService{repo: repo}
}

func (svc MembersService) AddMembers(ctx context.Context, members []entities.Member, organizationID string) error {
	for i := range members {
		if members[i].Role != entities.User && members[i].Role != entities.Owner && members[i].Role != entities.Editor {
			members[i].Role = entities.User
		}
		members[i].ID = uuid.NewString()
		members[i].OrganizationID = organizationID

	}
	err := svc.repo.AddMembers(ctx, members)
	return err
}

func (svc MembersService) DeleteMembers(ctx context.Context, members []entities.Member, organizationID string) error {
	for i := range members {
		members[i].ID = uuid.NewString()
		members[i].OrganizationID = organizationID

	}
	err := svc.repo.DeleteMembers(ctx, members, organizationID)
	return err
}

func (svc MembersService) UpdateMembers(ctx context.Context, members []entities.Member, organizationID string) error {
	for i := range members {
		if members[i].Role != entities.User && members[i].Role != entities.Owner && members[i].Role != entities.Editor {
			members[i].Role = ""
		}
	}
	err := svc.repo.UpdateMembers(ctx, members, organizationID)
	return err
}
