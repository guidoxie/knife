package log

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

// 按默认配置初始化
func init() {
	if err := New(gin.TestMode); err != nil {
		panic(err)
	}
}

// 自定义初始化
func New(mode string, config ...zap.Config) error {
	var conf zap.Config
	switch mode {
	case gin.ReleaseMode:
		conf = zap.NewProductionConfig()
	case gin.DebugMode, gin.TestMode:
		conf = zap.NewDevelopmentConfig()
	}

	if len(config) > 0 {
		conf = config[0]
	}

	l, err := conf.Build(zap.AddCallerSkip(1))
	if err != nil {
		return err
	}
	logger = l.Sugar()
	return err
}
