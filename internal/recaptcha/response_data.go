package recaptcha

import "time"

type ResponseData struct {
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	Score       float32   `json:"score"`
	Success     bool      `json:"success"`
	ErrorCodes  []string  `json:"error-codes"`
}
