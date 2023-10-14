package entities

import "errors"

//type Members struct {
//	EmailAdminOrganization string   `json:"emailAdminOrganization"`
//	DataMembers            []Member `json:"data_members"`
//}

var (
	ErrMembersAlreadyAdded = errors.New("members are already added")
	ErrMembersDoesntExist  = errors.New("members doesn't exist")
)

type Member struct {
	ID             string `json:"id"`
	Email          string `json:"email"`
	FirstName      string `json:"name"`
	SecondName     string `json:"secondName"`
	OrganizationID string `json:"organizationID"`
}
