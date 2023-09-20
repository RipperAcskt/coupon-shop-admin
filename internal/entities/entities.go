package entities

type Organization struct {
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
