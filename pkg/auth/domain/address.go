package auth_domain

type Address struct {
	ID      string `json:"addressId"`
	Country string `json:"country"`
	Region  string `json:"region"`
	City    string `json:"city"`
	SubCity string `json:"subCity"`
	Woreda  int    `json:"woreda"`
	HouseNo string `json:"houseNo"`
}
