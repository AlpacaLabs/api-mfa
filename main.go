package main

import (
	"sync"

	"github.com/AlpacaLabs/mfa/internal/app"
	"github.com/AlpacaLabs/mfa/internal/config"
)

func main() {
	c := config.LoadConfig()
	a := app.NewApp(c)

	var wg sync.WaitGroup

	wg.Add(1)
	go a.Run()

	wg.Wait()
}
