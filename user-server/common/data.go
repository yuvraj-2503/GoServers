package common

type AuthenticatedUser struct {
	UserId      string       `json:"userId"`
	EmailId     string       `json:"emailId"`
	PhoneNumber *PhoneNumber `json:"phoneNumber"`
	ApiKey      string       `json:"apiKey"`
}

type PhoneNumber struct {
	CountryCode string `json:"countryCode" bson:"countryCode"`
	Number      string `json:"number" bson:"number"`
}

type Device struct {
	SerialNumber string     `json:"serialNumber"`
	Name         string     `json:"name"`
	OS           string     `json:"os"`
	FingerPrint  string     `json:"fingerPrint"`
	DeviceType   DeviceType `json:"deviceType"`
}

type DeviceType string

const (
	ANDROID DeviceType = "ANDROID"
	IOS     DeviceType = "IOS"
	DESKTOP DeviceType = "DESKTOP"
)
