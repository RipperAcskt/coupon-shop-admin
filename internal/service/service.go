package service

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

func New(repo Repo) Service {
	return Service{
		SubscriptionService: NewSubscriptionService(repo),
		OrganizationService: NewOrganizationService(repo),
		CouponService:       NewCouponService(repo),
	}
}
