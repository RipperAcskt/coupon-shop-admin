package entities

import "errors"

//type Members struct {
//	EmailAdminOrganization string   `json:"emailAdminOrganization"`
//	DataMembers            []Member `json:"data_members"`
//}

var (
	ErrMembersAlreadyAdded = errors.New("organization already exists")
)

type Member struct {
	Email      string `json:"email"`
	FirstName  string `json:"name"`
	SecondName string `json:"secondName"`
}
