package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var zapLogger *zap.Logger

const DefaultFilename = "./app.log"

type Config struct {
	stdoutConfig  *StdoutConfig
	rollingConfig *FileConfig
}

type StdoutConfig struct {
}

type FileConfig struct {
	encoding zapcore.EncoderConfig
	logger   *lumberjack.Logger
}

func (c *Config) Init() {
	consoleEncoder := zap.NewDevelopmentEncoderConfig()
	consoleEncoder.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(consoleEncoder),
		zapcore.AddSync(zapcore.Lock(os.Stdout)),
		zapcore.DebugLevel,
	)

	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename: c.rollingConfig.logger.Filename,
		MaxSize:  c.rollingConfig.logger.MaxSize, // megabytes
		MaxAge:   c.rollingConfig.logger.MaxAge,  // days
	})
	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(fileWriter),
		zapcore.DebugLevel,
	)

	core := zapcore.NewTee(consoleCore, fileCore)
	zapLogger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	zapLogger.Info("zap logger initialized")
	zapLogger.Debug("zap logger debug level enabled")
}

func Default() *Config {
	return &Config{
		rollingConfig: &FileConfig{
			encoding: zapcore.EncoderConfig{},
			logger: &lumberjack.Logger{
				Filename: DefaultFilename,
				MaxSize:  500, // megabytes
				MaxAge:   30,  // days
			},
		},
	}
}

func (c *Config) WithFilename(filename string) *Config {
	c.rollingConfig.logger.Filename = filename
	return c
}
