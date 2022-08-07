package http

import (
	"github.com/gin-gonic/gin"
	"github.com/guidoxie/knife/pkg/errcode"
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`    // 返回码，0=正常
	Message string      `json:"message"` // 信息，出错时存储错误信息
	Data    interface{} `json:"data"`    // 返回数据
}

type Context struct {
	*gin.Context
}
type HandlerFunc func(c *Context)

func NewContext(c *gin.Context) *Context {
	return &Context{
		c,
	}
}

func (c *Context) OutSuccess(data interface{}) {
	c.JSON(http.StatusOK, &Response{Data: data})
}

// 参数错误
func (c *Context) OutParamErr(err error) {
	c.OutErr(errcode.InvalidParams.WithDetails(err.Error()))
}

// 系统错误
func (c *Context) OutSysErr(err error) {

}

func (c *Context) OutErr(err *errcode.Error) {
	c.JSON(http.StatusOK, &Response{
		Code:    err.Code(),
		Message: err.Msg(),
		Data:    struct{}{},
	})
}

// 绑定参数参数
func (c *Context) BindAndValid(v interface{}) {

}

//
