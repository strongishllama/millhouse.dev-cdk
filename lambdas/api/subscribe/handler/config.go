package handler

var (
	Cfg *Config
)

type Config struct {
	RecaptchaSecret string `json:"recaptchaSecret"`
	QueueURL        string `json:"queueUrl"`
	From            string `json:"from"`
}
