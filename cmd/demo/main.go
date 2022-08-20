package main

import (
	"github.com/guidoxie/knife/cmd/demo/internal/api"
	"github.com/guidoxie/knife/pkg/app"
	"github.com/guidoxie/knife/pkg/log"
)

func main() {
	demo := app.New("demo")
	api.Init()
	log.New()
	demo.Run()
}
