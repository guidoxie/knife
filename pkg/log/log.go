package log

import "go.uber.org/zap"

var logger *zap.SugaredLogger

func New(options ...zap.Option) error {
	l, err := zap.NewDevelopment(zap.AddCallerSkip(1))
	if err != nil {
		return err
	}
	logger = l.Sugar()
	return err
}
