package entities

type Subscription struct {
	ID          string `json:"-"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Level       int    `json:"level"`
}

func NewSubscription() Subscription {
	return Subscription{}
}
