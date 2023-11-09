package entities

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
)

var (
	ErrNoAnyCoupons       = errors.New("there is not a single coupon")
	ErrCouponDoesNotExist = errors.New("coupon does not exist")
	ErrNoMedia            = fmt.Errorf("no media")
)

type Coupon struct {
	ID           string  `json:"id"`
	Name         string  `form:"name" json:"name"`
	Description  string  `form:"description" json:"description"`
	Price        int     `form:"price" json:"price"`
	Level        int     `form:"level" json:"level"`
	Percent      int     `form:"percent" json:"percent"`
	ContentUrl   string  `json:"content_url"`
	Media        Media   `json:"-"`
	Region       string  `form:"region" json:"region"`
	Category     string  `form:"category" json:"category"`
	Subcategory  *string `form:"subcategory" json:"subcategory"`
	Organization string  `form:"organization" json:"organization"`
}

type Media struct {
	ID   string
	Path string
}

func NewCoupon() Coupon {
	return Coupon{
		ID: uuid.NewString(),
	}
}

func NewMedia() Media {
	id := uuid.NewString()
	return Media{
		ID:   id,
		Path: "/store/" + id + ".jpg",
	}
}

func NewMediaId(id string) Media {
	return Media{
		ID:   id,
		Path: "/store/" + id + ".jpg",
	}
}
