package thingmocker

import (
	"log"

	"go.uber.org/zap"
)

var (
	sugar *zap.SugaredLogger
	std   *log.Logger
)

func init() {
	var logger *zap.Logger

	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	config.OutputPaths = []string{"stdout"}
	logger, _ = config.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	defer logger.Sync()
	sugar = logger.Sugar()
	std = zap.NewStdLog(logger)
}

func Print(args ...interface{}) {
	std.Print(args...)
}

func Println(args ...interface{}) {
	std.Println(args...)
}

func Printf(template string, args ...interface{}) {
	std.Printf(template, args...)
}

func Info(args ...interface{}) {
	sugar.Info(args...)
}

func Infof(template string, args ...interface{}) {
	sugar.Infof(template, args...)
}

func Debug(args ...interface{}) {
	sugar.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	sugar.Debugf(template, args...)
}

func Warn(args ...interface{}) {
	sugar.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	sugar.Warnf(template, args...)
}

func Error(args ...interface{}) {
	sugar.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	sugar.Errorf(template, args...)
}

func Fatal(args ...interface{}) {
	sugar.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	sugar.Fatalf(template, args...)
}

func Panic(args ...interface{}) {
	sugar.Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	sugar.Panicf(template, args...)
}
