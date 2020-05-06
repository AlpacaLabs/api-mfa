package main

import (
	"sync"

	"github.com/AlpacaLabs/api-mfa/internal/app"
	"github.com/AlpacaLabs/api-mfa/internal/configuration"
)

func main() {
	c := configuration.LoadConfig()
	a := app.NewApp(c)

	var wg sync.WaitGroup

	wg.Add(1)
	go a.Run()

	wg.Wait()
}
