package http

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/locales/zh_Hant_TW"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/guidoxie/knife/pkg/log"
	"time"
)

// 翻译
func Translations() HandlerFunc {
	return func(c *Context) {
		uni := ut.New(en.New(), zh.New(), zh_Hant_TW.New())
		locale := c.GetHeader("locale")
		trans, _ := uni.GetTranslator(locale)
		v, ok := binding.Validator.Engine().(*validator.Validate)
		if ok {
			switch locale {
			case "zh":
				_ = zh_translations.RegisterDefaultTranslations(v, trans)
				break
			case "en":
				_ = en_translations.RegisterDefaultTranslations(v, trans)
				break
			default:
				_ = zh_translations.RegisterDefaultTranslations(v, trans)
				break
			}
			c.Set("trans", trans)
		}
		c.Next()
	}
}

type AccessLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (a *AccessLogWriter) Write(p []byte) (int, error) {
	if n, err := a.body.Write(p); err != nil {
		return n, err
	}
	return a.ResponseWriter.Write(p)
}

// 访问日志
func AccessLog() HandlerFunc {
	return func(c *Context) {
		bodyWriter := &AccessLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = bodyWriter

		beginTime := time.Now().Unix()
		c.Next()
		endTime := time.Now().Unix()

		log.Infof("access log method: %s, path: %s, status_code: %d, begin_time: %d, end_time: %d , request: %s, response: %s",
			c.Request.Method,
			c.Request.URL.Path,
			bodyWriter.Status(),
			beginTime,
			endTime,
			c.Request.PostForm.Encode(),
			bodyWriter.body.String())
	}
}

// recovery
func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Errorf("panic recover err: %v", err)
				c.OutSysErr(nil)
				c.Abort()
			}
		}()
		c.Next()
	}
}
