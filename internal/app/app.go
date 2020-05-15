package app

import (
	"github.com/AlpacaLabs/go-kontext"

	"github.com/AlpacaLabs/api-mfa/internal/grpc"

	"github.com/AlpacaLabs/api-mfa/internal/configuration"
	"github.com/AlpacaLabs/api-mfa/internal/db"
	"github.com/AlpacaLabs/api-mfa/internal/http"
	"github.com/AlpacaLabs/api-mfa/internal/service"
	log "github.com/sirupsen/logrus"
)

type App struct {
	config configuration.Config
}

func NewApp(c configuration.Config) App {
	return App{
		config: c,
	}
}

func (a App) Run() {
	dbConn, err := db.Connect(a.config.SQLConfig)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	dbClient := db.NewClient(dbConn)
	accountConn, err := kontext.Dial(a.config.AccountGRPCAddress)
	if err != nil {
		log.Fatalf("failed to dial Account service: %v", err)
	}
	svc := service.NewService(a.config, dbClient, accountConn)

	httpServer := http.NewServer(a.config, svc)
	go httpServer.Run()

	grpcServer := grpc.NewServer(a.config, svc)
	go grpcServer.Run()
}
