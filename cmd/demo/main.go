package main

import (
	"github.com/guidoxie/knife/cmd/demo/internal/api"
	"github.com/guidoxie/knife/pkg/app"
)

func main() {
	demo := app.New("demo")
	api.Init()
	demo.Run()
}
