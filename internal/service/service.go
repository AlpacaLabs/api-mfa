package service

import (
	"github.com/AlpacaLabs/mfa/internal/configuration"
	"github.com/AlpacaLabs/mfa/internal/db"
)

type Service struct {
	config   configuration.Config
	dbClient db.Client
}

func NewService(config configuration.Config, dbClient db.Client) Service {
	return Service{
		config:   config,
		dbClient: dbClient,
	}
}
