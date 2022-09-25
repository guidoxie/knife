package user

import "github.com/guidoxie/knife/pkg/app"

func Init() {
	p := controller{}
	group := app.HttpServer().Group("v1/user/")
	group.POST("add", p.Add)
	group.GET("add", p.Add)
}
