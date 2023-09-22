package entities

import "github.com/google/uuid"

type Coupon struct {
	ID          string `json:"id"`
	Name        string `form:"name"`
	Description string `form:"description"`
	Price       int    `form:"price"`
	Level       int    `form:"level"`
	Media       Media
}

type Media struct {
	ID   string
	Path string
}

func NewCoupon() Coupon {
	return Coupon{}
}

func NewMedia() Media {
	id := uuid.NewString()
	return Media{
		ID:   id,
		Path: "/store/" + id,
	}
}
