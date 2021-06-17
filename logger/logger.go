package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log *zap.Logger
)

func init() {
	logConfig := zap.Config{
		Level:    zap.NewAtomicLevelAt(zap.InfoLevel),
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			LevelKey:     "level",
			TimeKey:      "time",
			MessageKey:   "msg",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeLevel:  zapcore.LowercaseLevelEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
		OutputPaths: []string{"stdout"},

		/*		{
				"level": "info"
				"time": "2006-01-02T15:04:05.000Z0700"
				"msg": "This is a log"
			},*/
	}

	var err error
	if log, err = logConfig.Build(); err != nil {
		panic(err.Error())
	}
}

func Info(msg string, tags ...zap.Field) {
	log.Info(msg, tags...)
	err := log.Sync()
	if err != nil {
		return
	}
}

func Error(msg string, err error, tags ...zap.Field) {
	tags = append(tags, zap.NamedError("error", err))

	log.Info(msg, tags...)
	_err := log.Sync()
	if _err != nil {
		return
	}
}

func GetLogger() *zap.Logger {
	return log
}
