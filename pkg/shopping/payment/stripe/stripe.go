package stripe

// Card symbolizes a user card.
type Card struct {
	Number   string `json:"number"`
	ExpMonth string `json:"exp_month"`
	ExpYear  string `json:"exp_year" validate:"len=4"`
	CVC      string `json:"cvc" validate:"len=3"`
}
