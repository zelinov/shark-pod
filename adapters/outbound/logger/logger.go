package logger

import (
	"log"
	conf "sharks/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func init() {
	var (
		level zapcore.Level
		err   error
	)

	if err = level.UnmarshalText([]byte(conf.GetConfig().LogLevel)); err != nil {
		log.Fatal(err)
	}

	config := zap.NewDevelopmentConfig()
	config.Level = zap.NewAtomicLevelAt(level)

	if Log, err = config.Build(); err != nil {
		log.Fatal(err)
	}
}
