package main

import (
	"github.com/guidoxie/knife/cmd/demo/internal/api"
	"github.com/guidoxie/knife/pkg/app"
	"github.com/guidoxie/knife/pkg/app/http"
	"github.com/guidoxie/knife/pkg/log"
	"time"
)

func main() {
	demo := app.New("demo", http.Options{
		WebAddr:      "0.0.0.0:8080",
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 120 * time.Second,
	})
	defer log.Recover()
	api.Init()
	demo.Run()
}
