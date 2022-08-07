package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	ExMethodStaticFile = "StaticFile"
	ExMethodStatic     = "Static"
	ExMethodStaticFS   = "StaticFS"
)

// 路由组
type IGroup interface {
	IRoute
	// 获取路由组中所有路由
	Routes() []*route
}

type Group struct {
	basePath string
	routes   []*route
	handlers []HandlerFunc
}

func (r *Group) Use(handlerFunc ...HandlerFunc) IRoute {
	r.handlers = append(r.handlers, handlerFunc...)
	return r
}

func (r *Group) Handle(httpMethod, relativePath string, handlers ...HandlerFunc) IRoute {
	r.routes = append(r.routes, &route{
		method:   httpMethod,
		path:     relativePath,
		handlers: handlers,
	})
	return r
}

func (r *Group) Any(relativePath string, handlers ...HandlerFunc) IRoute {
	r.Handle(http.MethodGet, relativePath, handlers...)
	r.Handle(http.MethodPost, relativePath, handlers...)
	r.Handle(http.MethodPut, relativePath, handlers...)
	r.Handle(http.MethodPatch, relativePath, handlers...)
	r.Handle(http.MethodHead, relativePath, handlers...)
	r.Handle(http.MethodOptions, relativePath, handlers...)
	r.Handle(http.MethodDelete, relativePath, handlers...)
	r.Handle(http.MethodConnect, relativePath, handlers...)
	r.Handle(http.MethodTrace, relativePath, handlers...)
	return r
}

func (r *Group) GET(relativePath string, handlers ...HandlerFunc) IRoute {
	r.routes = append(r.routes, &route{
		method:   http.MethodGet,
		path:     relativePath,
		handlers: handlers,
	})
	return r
}

func (r *Group) POST(relativePath string, handlers ...HandlerFunc) IRoute {
	r.routes = append(r.routes, &route{
		method:   http.MethodPost,
		path:     relativePath,
		handlers: handlers,
	})
	return r
}

func (r *Group) DELETE(relativePath string, handlers ...HandlerFunc) IRoute {
	r.routes = append(r.routes, &route{
		method:   http.MethodDelete,
		path:     relativePath,
		handlers: handlers,
	})
	return r
}

func (r *Group) PATCH(relativePath string, handlers ...HandlerFunc) IRoute {
	r.routes = append(r.routes, &route{
		method:   http.MethodPatch,
		path:     relativePath,
		handlers: handlers,
	})
	return r
}

func (r *Group) PUT(relativePath string, handlers ...HandlerFunc) IRoute {
	r.routes = append(r.routes, &route{
		method:   http.MethodPut,
		path:     relativePath,
		handlers: handlers,
	})
	return r
}

func (r *Group) OPTIONS(relativePath string, handlers ...HandlerFunc) IRoute {
	r.routes = append(r.routes, &route{
		method:   http.MethodPut,
		path:     relativePath,
		handlers: handlers,
	})
	return r
}

func (r *Group) HEAD(relativePath string, handlers ...HandlerFunc) IRoute {
	r.routes = append(r.routes, &route{
		method:   http.MethodHead,
		path:     relativePath,
		handlers: handlers,
	})
	return r
}

func (r *Group) StaticFile(relativePath, filepath string) IRoute {
	r.routes = append(r.routes, &route{
		method:   ExMethodStaticFile,
		path:     relativePath,
		filePath: filepath,
	})
	return r
}

func (r *Group) Static(relativePath, root string) IRoute {
	r.routes = append(r.routes, &route{
		method:   ExMethodStatic,
		path:     relativePath,
		filePath: root,
	})
	return r
}

func (r *Group) StaticFS(relativePath string, fs http.FileSystem) IRoute {
	r.routes = append(r.routes, &route{
		method: ExMethodStaticFS,
		path:   relativePath,
		fs:     fs,
	})
	return r
}

func (r *Group) Routes() []*route {
	return r.routes
}

func (r *Group) Handlers() []gin.HandlerFunc {
	handlers := make([]gin.HandlerFunc, len(r.handlers))
	for i, h := range r.handlers {
		handler := h
		handlers[i] = func(context *gin.Context) {
			handler(NewContext(context))
		}
	}
	return handlers
}
