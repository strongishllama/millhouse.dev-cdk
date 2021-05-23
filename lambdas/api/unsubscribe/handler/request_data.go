package handler

import (
	"net/mail"

	"github.com/gofor-little/xerror"
)

type RequestData struct {
	ID           string `json:"id"`
	EmailAddress string `json:"emailAddress"`
}

func (r *RequestData) Validate() error {
	if r.ID == "" {
		return xerror.Newf("ID cannot be empty")
	}

	if _, err := mail.ParseAddress(r.EmailAddress); err != nil {
		return xerror.New("failed to validate EmailAddress", err)
	}

	return nil
}
