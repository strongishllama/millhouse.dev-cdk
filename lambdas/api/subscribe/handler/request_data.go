package handler

import (
	"net/mail"

	"github.com/gofor-little/xerror"
)

type RequestData struct {
	EmailAddress            string `json:"emailAddress"`
	ReCaptchaChallengeToken string `json:"recaptchaChallengeToken"`
}

func (r *RequestData) Validate() error {
	if _, err := mail.ParseAddress(r.EmailAddress); err != nil {
		return xerror.New("failed to validate EmailAddress", err)
	}

	if r.ReCaptchaChallengeToken == "" {
		return xerror.Newf("ReCaptchaChallengeToken cannot be empty")
	}

	return nil
}
