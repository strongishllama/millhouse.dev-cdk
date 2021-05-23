package handler

import (
	"context"

	"github.com/gofor-little/cfg"
	"github.com/gofor-little/env"
	"github.com/gofor-little/xerror"

	"github.com/strongishllama/millhouse.dev-cdk/pkg/db"
	"github.com/strongishllama/millhouse.dev-cdk/pkg/notification"
)

var (
	Cfg *Config
)

func Initialize(configSecretARN string) error {
	if err := cfg.Initialize("", ""); err != nil {
		return xerror.New("failed to initialize cfg package", err)
	}

	Cfg = &Config{
		WebsiteDomain: env.Get("WEBSITE_DOMAIN", "millhouse.dev"),
		APIDomain:     env.Get("API_DOMAIN", "api.millhouse.dev"),
	}
	if err := cfg.Load(context.Background(), configSecretARN, Cfg); err != nil {
		return xerror.New("failed to load config", err)
	}

	emailQueueURL, err := env.MustGet("EMAIL_QUEUE_URL")
	if err != nil {
		return xerror.New("failed to get environment variable", err)
	}
	if err := notification.Initialize("", "", emailQueueURL); err != nil {
		return xerror.New("failed to initialize the notifications package", err)
	}

	tableName, err := env.MustGet("TABLE_NAME")
	if err != nil {
		return xerror.New("failed to get environment variable", err)
	}
	if err := db.Initialize("", "", tableName); err != nil {
		return xerror.New("failed to initialize the db package", err)
	}

	return nil
}
