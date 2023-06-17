package logger

import (
	"fmt"
	"os"

	"github.com/go-playground/validator"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger interface {
	Debug(string)
	Info(string)
	Warn(string)
	Error(error)
	Panic(error)
	Fatal(error)
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
	Panicf(string, ...interface{})
	Fatalf(string, ...interface{})
	Close()
}

type logger struct {
	logger *zap.Logger
	file   *lumberjack.Logger
}

type LoggerOption struct {
	LogLevel    LogLevel
	LogName     string `validate:"required"`
	SkipCaller  int    `validate:"required"`
	LogSize     int    `validate:"required"`
	LogAge      int    `validate:"required"`
	LogBackup   int    `validate:"required"`
	LogCompress bool
}

func validateLoggerOption(option *LoggerOption) *LoggerOption {
	validate := validator.New()
	if err := validate.Struct(option); err != nil {
		panic(err)
	}

	return option
}

func NewLogger(option *LoggerOption) Logger {
	option = validateLoggerOption(option)

	productionEncoderConfig := zap.NewProductionEncoderConfig()
	productionEncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	productionEncoderConfig.TimeKey = "datetime"
	fileEncoder := zapcore.NewJSONEncoder(productionEncoderConfig)
	consoleEncoder := zapcore.NewConsoleEncoder(productionEncoderConfig)

	priority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.Level(option.LogLevel)
	})

	lumberJackLogger := &lumberjack.Logger{
		Filename:   option.LogName,
		MaxSize:    option.LogSize,
		MaxBackups: option.LogBackup,
		MaxAge:     option.LogAge,
		Compress:   option.LogCompress,
	}

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(lumberJackLogger), priority),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), priority),
	)

	l := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	return &logger{logger: l, file: lumberJackLogger}
}

func (l logger) Debug(msg string) {
	l.logger.Debug(msg)
}

func (l logger) Info(msg string) {
	l.logger.Info(msg)
}

func (l logger) Warn(msg string) {
	l.logger.Warn(msg)
}

func (l logger) Error(err error) {
	l.logger.Error(err.Error())
}

func (l logger) Panic(err error) {
	l.logger.Panic(err.Error())
}

func (l logger) Fatal(err error) {
	l.logger.Fatal(err.Error())
}

func (l logger) Debugf(format string, v ...interface{}) {
	l.logger.Debug(fmt.Sprintf(format, v...))
}

func (l logger) Infof(format string, v ...interface{}) {
	l.logger.Info(fmt.Sprintf(format, v...))
}

func (l logger) Warnf(format string, v ...interface{}) {
	l.logger.Warn(fmt.Sprintf(format, v...))
}

func (l logger) Errorf(format string, v ...interface{}) {
	l.logger.Error(fmt.Sprintf(format, v...))
}

func (l logger) Panicf(format string, v ...interface{}) {
	l.logger.Panic(fmt.Sprintf(format, v...))
}

func (l logger) Fatalf(format string, v ...interface{}) {
	l.logger.Fatal(fmt.Sprintf(format, v...))
}

func (l *logger) Close() {
	l.logger.Sync()
	l.file.Close()
	fmt.Println("Logger Closed")
}
