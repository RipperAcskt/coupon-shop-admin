package service

type OrganizationService struct {
	repo OrganizationRepoInterface
}

type OrganizationRepoInterface interface {
	CreateOrganization()
}

func NewOrganizationService(repo OrganizationRepoInterface) OrganizationService {
	return OrganizationService{repo: repo}
}

func (svc OrganizationService) CreateOrganization() {
}
