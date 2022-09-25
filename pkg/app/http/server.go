package http

import (
	"github.com/gin-gonic/gin"
	"github.com/guidoxie/knife/pkg/log"
	"time"

	"net/http"
	"path"
)

type IServer interface {
	Name() string
	// 获取 gin.Engine 对象
	GetEngine() *gin.Engine
	// 获取默认中间件
	GetDefaultMiddle() []gin.HandlerFunc
	// 设置默认中间件
	SetDefaultMiddle(middle ...HandlerFunc) IServer
	// 添加默认中间件。如果中间件已经注册过则会忽略。
	AddDefaultMiddle(middle ...HandlerFunc) IServer
	// 获取已经设置的 RouteGroup
	GetGroups() []*Group
	// 路由组
	Group(relativePath string, middle ...HandlerFunc) *Group
	// 运行
	Run() error
}

type Server struct {
	http.Server
	name          string
	engine        *gin.Engine
	defaultMiddle []HandlerFunc
	routes        []*route
	routeGroups   []*Group
}

type Options struct {
	WebAddr        string // Web Server 监听地址
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxHeaderBytes int
}

func NewServer(name string, options ...Options) *Server {
	// 关闭gin自身的日志
	gin.SetMode(gin.ReleaseMode)
	s := &Server{name: name, engine: gin.New()}
	if len(options) > 0 {
		s.Addr = options[0].WebAddr
		s.ReadTimeout = options[0].ReadTimeout
		s.WriteTimeout = options[0].WriteTimeout
		s.MaxHeaderBytes = options[0].MaxHeaderBytes
	}
	return s
}

func (s *Server) Name() string {
	return s.name
}

func (s *Server) GetEngine() *gin.Engine {
	return s.engine
}

func (s *Server) GetDefaultMiddle() []gin.HandlerFunc {
	handlers := make([]gin.HandlerFunc, len(s.defaultMiddle))
	for i, h := range s.defaultMiddle {
		handler := h
		handlers[i] = func(context *gin.Context) {
			handler(NewContext(context))
		}
	}
	return handlers
}

func (s *Server) SetDefaultMiddle(middle ...HandlerFunc) IServer {
	s.defaultMiddle = middle
	return s
}

func (s *Server) AddDefaultMiddle(middle ...HandlerFunc) IServer {
	s.defaultMiddle = append(s.defaultMiddle, middle...)
	return s
}

func (s *Server) GetGroups() []*Group {
	return s.routeGroups
}

func (s *Server) Group(relativePath string, middle ...HandlerFunc) *Group {
	group := &Group{
		basePath: relativePath,
		routes:   nil,
		handlers: middle,
	}
	s.routeGroups = append(s.routeGroups, group)
	return group
}

func (s *Server) Run() error {
	// 注册默认中间键
	s.engine.Use(s.GetDefaultMiddle()...)
	// 注册路由
	for _, g := range s.GetGroups() {
		group := s.GetEngine().Group(s.joinPaths(g.basePath), g.Handlers()...)
		for _, r := range g.routes {
			switch r.method {
			case ExMethodStatic:
				group.Static(r.path, r.filePath)
			case ExMethodStaticFS:
				group.StaticFS(r.path, r.fs)
			case ExMethodStaticFile:
				group.StaticFile(r.path, r.filePath)
			default:
				group.Handle(r.method, r.path, r.Handlers()...)
				log.Infof("%s\t/%s\t--> %v (%d handlers)", r.method, path.Join(s.joinPaths(g.basePath), r.path),
					r.HandlerNames()[r.HandlerNumber()-1], r.HandlerNumber()+g.HandlerNumber()+len(s.defaultMiddle))
			}
		}
	}
	s.Handler = s.GetEngine()
	log.Infof("Listening and serving HTTP on %s", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

// api/${name}/
func (s *Server) joinPaths(relativePath string) string {
	if relativePath == "" {
		return path.Join("api", s.name)
	}

	finalPath := path.Join(path.Join("api", s.name), relativePath)
	if lastChar(relativePath) == '/' && lastChar(finalPath) != '/' {
		return finalPath + "/"
	}
	return finalPath
}

func lastChar(str string) uint8 {
	if str == "" {
		panic("The length of the string can't be 0")
	}
	return str[len(str)-1]
}
