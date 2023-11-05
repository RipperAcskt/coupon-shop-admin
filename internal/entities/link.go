package entities

import "errors"

var (
	ErrLinkAlreadyExists = errors.New("category link exists")
	ErrNoAnyLink         = errors.New("there is not a single link")
	ErrLinkDoesNotExist  = errors.New("link does not exist")
)

type Link struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Link   string `json:"link"`
	Region string `json:"region"`
}
