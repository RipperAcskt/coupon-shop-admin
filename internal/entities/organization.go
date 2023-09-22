package entities

import "errors"

var (
	ErrOrganizationAlreadyExists = errors.New("organization already exists")
	ErrNoAnyOrganization         = errors.New("there is not a single organization")
	ErrOrganizationnDoesNotExist = errors.New("organization does not exist")
)

type Organization struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	EmailAdmin        string `json:"email_admin"`
	LevelSubscription int    `json:"levelSubscription"`
}

func NewOrganization() Organization {
	return Organization{}
}

const (
	LowLevel = iota + 1
	MidLevel
	HighLevel
)
