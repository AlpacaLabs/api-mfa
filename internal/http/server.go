package http

import (
	"fmt"
	"net/http"

	"github.com/AlpacaLabs/mfa/internal/config"
	"github.com/AlpacaLabs/mfa/internal/services"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	config  config.Config
	service services.Service
}

func NewServer(config config.Config, service services.Service) Server {
	return Server{
		config:  config,
		service: service,
	}
}

func (s Server) Run() {
	r := mux.NewRouter()

	//r.HandleFunc("/mfa", s.SendCodeOptions).Methods(http.MethodPost)

	addr := fmt.Sprintf(":%d", s.config.HTTPPort)
	log.Infof("Listening for HTTP on %s...\n", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
