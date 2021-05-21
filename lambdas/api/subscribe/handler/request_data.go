package handler

import (
	"net/mail"

	"github.com/gofor-little/xerror"
)

type RequestData struct {
	EmailAddress            string `json:"emailAddress"`
	ReCaptchaChallengeToken string `json:"recaptchaChallengeToken"`
}

func (s *RequestData) Validate() error {
	if _, err := mail.ParseAddress(s.EmailAddress); err != nil {
		return xerror.New("failed to validate EmailAddress", err)
	}

	if s.ReCaptchaChallengeToken == "" {
		return xerror.Newf("ReCaptchaChallengeToken cannot be empty")
	}

	return nil
}
