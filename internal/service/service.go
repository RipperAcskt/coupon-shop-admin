package service

type Service struct {
	SubscriptionService
	OrganizationService
}

type Repo interface {
	SubscriptionRepoInterface
	OrganizationRepoInterface
}

func New(repo Repo) Service {
	return Service{
		SubscriptionService: NewSubscriptionService(repo),
		OrganizationService: NewOrganizationService(repo),
	}
}
