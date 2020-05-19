package app

import (
	"sync"

	"github.com/AlpacaLabs/api-mfa/internal/async"

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
	config := a.config
	dbConn, err := db.Connect(config.SQLConfig)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	dbClient := db.NewClient(dbConn)
	accountConn, err := kontext.Dial(config.AccountGRPCAddress)
	if err != nil {
		log.Fatalf("failed to dial Account service: %v", err)
	}
	svc := service.NewService(config, dbClient, accountConn)

	var wg sync.WaitGroup

	wg.Add(1)
	httpServer := http.NewServer(config, svc)
	go httpServer.Run()

	wg.Add(1)
	grpcServer := grpc.NewServer(config, svc)
	go grpcServer.Run()

	wg.Add(1)
	go async.RelayMessagesForSend(config, dbClient)

	wg.Wait()
}
