package handler

import (
	"context"

	"github.com/gofor-little/cfg"
	"github.com/gofor-little/xerror"
)

var (
	config *Config
)

type Config struct {
	RecaptchaSecret string `json:"recaptchaSecret"`
}

func Initialize(ctx context.Context, configSecretARN string) error {
	config = &Config{}
	if err := cfg.Load(ctx, configSecretARN, config); err != nil {
		return xerror.New("failed to load config", err)
	}

	return nil
}
