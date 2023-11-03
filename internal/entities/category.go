package entities

import "errors"

var (
	ErrCategoryAlreadyExists = errors.New("category already exists")
	ErrNoAnyCategory         = errors.New("there is not a single category")
	ErrCategoryDoesNotExist  = errors.New("category does not exist")
)

type Category struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Subcategory bool   `json:"subcategory"`
}
