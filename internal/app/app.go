package app

import (
	"sync"

	"github.com/AlpacaLabs/api-mfa/internal/grpc"

	"github.com/AlpacaLabs/api-mfa/internal/configuration"
	"github.com/AlpacaLabs/api-mfa/internal/db"
	"github.com/AlpacaLabs/api-mfa/internal/http"
	"github.com/AlpacaLabs/api-mfa/internal/service"
	log "github.com/sirupsen/logrus"
	grpcGo "google.golang.org/grpc"
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
	dbConn := db.Connect(a.config.DBUser, a.config.DBPass, a.config.DBHost, a.config.DBName)
	dbClient := db.NewClient(dbConn)
	accountConn, err := grpcGo.Dial(a.config.AccountGRPCAddress)
	if err != nil {
		log.Fatalf("failed to dial Account service: %v", err)
	}
	svc := service.NewService(a.config, dbClient, accountConn)

	var wg sync.WaitGroup

	wg.Add(1)
	httpServer := http.NewServer(a.config, svc)
	httpServer.Run()

	wg.Add(1)
	grpcServer := grpc.NewServer(a.config, svc)
	grpcServer.Run()

	wg.Wait()
}
