package entities

import (
	"errors"
	"github.com/google/uuid"
)

var (
	ErrOrganizationAlreadyExists = errors.New("organization already exists")
	ErrNoAnyOrganization         = errors.New("there is not a single organization")
	ErrOrganizationnDoesNotExist = errors.New("organization does not exist")
)

type Organization struct {
	ID                string   `form:"id" json:"id"`
	Name              string   `form:"name" json:"name"`
	EmailAdmin        string   `form:"email_admin" json:"email_admin"`
	LevelSubscription int      `form:"level_subscription" json:"level_subscription"`
	ORGN              string   `form:"orgn" json:"orgn"`
	KPP               string   `form:"kpp" json:"kpp"`
	INN               string   `form:"inn" json:"inn"`
	Address           string   `form:"address" json:"address"`
	Members           []Member `json:"members"`
	OrgImage          Image    `json:"-"`
	ContentUrl        string   `json:"content_url"`
}

func NewOrganization() Organization {
	return Organization{
		ID: uuid.NewString(),
	}

}

type Image struct {
	ID   string
	Path string
}

func NewImageID(ID string) Image {
	return Image{
		ID:   ID,
		Path: "/store/organization/" + ID + ".jpg",
	}
}

func NewImage() Image {
	id := uuid.NewString()
	return Image{
		ID:   id,
		Path: "/store/organization/" + id + ".jpg",
	}
}
