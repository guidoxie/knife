package app

import "github.com/guidoxie/knife/pkg/app/http"

func HttpServer() http.IServer {
	return Current().httpServer
}
