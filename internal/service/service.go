package service

import "github.com/RipperAcskt/coupon-shop-admin/config"

type Service struct {
	MembersService
	SubscriptionService
	OrganizationService
	CouponService
	AuthService
	RegionService
	CategoryService
	LinkService
	UserService
}

type Repo interface {
	MembersRepoInterface
	SubscriptionRepoInterface
	OrganizationRepoInterface
	CouponRepoInterface
	RegionRepoInterface
	CategoryRepoInterface
	LinkRepoInterface
	UserRepoInterface
}

type Cache interface {
	TokenRepo
}

func New(repo Repo, cache Cache, cfg config.Config) Service {
	return Service{
		SubscriptionService: NewSubscriptionService(repo),
		OrganizationService: NewOrganizationService(repo),
		CouponService:       NewCouponService(repo, cfg),
		MembersService:      NewMembersService(repo),
		AuthService:         NewAuthService(cache, cfg),
		RegionService:       NewRegionService(repo),
		CategoryService:     NewCategoryService(repo),
		LinkService:         NewLinkService(repo),
		UserService:         NewUserService(repo),
	}
}
