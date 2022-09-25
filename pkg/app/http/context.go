package http

import (
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/guidoxie/knife/pkg/errcode"
)

type Response struct {
	Code    int         `json:"code"`              // 返回码，0=正常
	Message string      `json:"message,omitempty"` // 信息，出错时存储错误信息
	Data    interface{} `json:"data,omitempty"`    // 返回数据
	Details []string    `json:"details,omitempty"` // 详细错误信息
}

type Context struct {
	*gin.Context
}

type HandlerFunc func(c *Context)

func NewContext(c *gin.Context) *Context {
	return &Context{c}
}

// 成功
func (c *Context) OutSuccess(data interface{}) {
	c.toResponse(data, errcode.Success)
}

// 业务错误
func (c *Context) OutErr(err *errcode.Error) {
	c.toResponse(nil, err)
}

// 参数错误
func (c *Context) OutParamErr(err ValidErrors) {
	c.toResponse(nil, errcode.InvalidParams.WithDetails(err.Error()))
}

// 系统错误
func (c *Context) OutSysErr() {
	c.toResponse(nil, errcode.ServerError)
}

func (c *Context) toResponse(data interface{}, err *errcode.Error) {
	c.JSON(err.StatusCode(), &Response{
		Code:    err.Code(),
		Message: err.Msg(),
		Data:    data,
		Details: err.Details(),
	})
}

// 绑定参数参数
func (c *Context) BindAndValid(v interface{}) ValidErrors {
	var errs ValidErrors
	err := c.ShouldBind(v)
	if err != nil {
		v := c.Value("trans")
		trans, _ := v.(ut.Translator)
		verrs, ok := err.(validator.ValidationErrors)
		if !ok {
			return nil
		}
		for key, value := range verrs.Translate(trans) {
			errs = append(errs, &ValidError{
				Key:     key,
				Message: value,
			})
		}
	}
	return errs
}
