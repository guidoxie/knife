package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type IRoute interface {
	Use(...HandlerFunc) IRoute

	Handle(string, string, ...HandlerFunc) IRoute
	Any(string, ...HandlerFunc) IRoute
	GET(string, ...HandlerFunc) IRoute
	POST(string, ...HandlerFunc) IRoute
	DELETE(string, ...HandlerFunc) IRoute
	PATCH(string, ...HandlerFunc) IRoute
	PUT(string, ...HandlerFunc) IRoute
	OPTIONS(string, ...HandlerFunc) IRoute
	HEAD(string, ...HandlerFunc) IRoute

	StaticFile(string, string) IRoute
	Static(string, string) IRoute
	StaticFS(string, http.FileSystem) IRoute

	Handlers() []gin.HandlerFunc
}

type route struct {
	method   string
	path     string
	handlers []HandlerFunc
	filePath string
	fs       http.FileSystem
}

func (r route) Handlers() []gin.HandlerFunc {
	handlers := make([]gin.HandlerFunc, len(r.handlers))
	for i, h := range r.handlers {
		handler := h
		handlers[i] = func(context *gin.Context) {
			handler(NewContext(context))
		}
	}
	return handlers
}
