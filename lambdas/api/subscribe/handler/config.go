package handler

type Config struct {
	RecaptchaSecret string `json:"recaptchaSecret"`
	To              string `json:"to"`
	From            string `json:"from"`
	WebsiteDomain   string `json:"websiteDomain"`
	APIDomain       string `json:"apiDomain"`
}
