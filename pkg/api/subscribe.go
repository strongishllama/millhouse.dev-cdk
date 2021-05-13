package api

import (
	"github.com/gofor-little/log"
)

func Subscribe(emailAddress string) error {
	log.Info(log.Fields{"emailAddress": emailAddress})
	return nil
}
