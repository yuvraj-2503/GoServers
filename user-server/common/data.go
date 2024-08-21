package common

type PhoneNumber struct {
	CountryCode string `json:"countryCode" bson:"countryCode"`
	Number      string `json:"number" bson:"number"`
}
