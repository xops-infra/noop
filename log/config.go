package log

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var zapLogger *zap.Logger

const DefaultFilename = "./app.log"

type Config struct {
	stdoutConfig  *StdoutConfig
	rollingConfig *FileConfig
	fieldsConfig  *FieldsConfig
}

type StdoutConfig struct {
	level Level
}

type FileConfig struct {
	level    Level
	encoding zapcore.EncoderConfig
	logger   *lumberjack.Logger
}

type FieldsConfig struct {
	fields []zapcore.Field
}

func (c *Config) Init() {
	consoleEncoder := zap.NewDevelopmentEncoderConfig()
	consoleEncoder.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(consoleEncoder),
		zapcore.AddSync(zapcore.Lock(os.Stdout)),
		zapcore.Level(c.stdoutConfig.level),
	)

	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename: c.rollingConfig.logger.Filename,
		MaxSize:  c.rollingConfig.logger.MaxSize, // megabytes
		MaxAge:   c.rollingConfig.logger.MaxAge,  // days
	})
	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(fileWriter),
		zapcore.Level(c.rollingConfig.level),
	)

	if c.fieldsConfig.fields != nil || len(c.fieldsConfig.fields) != 0 {
		consoleCore = consoleCore.With(c.fieldsConfig.fields)
		fileCore = fileCore.With(c.fieldsConfig.fields)
	}

	core := zapcore.NewTee(consoleCore, fileCore)
	zapLogger = zap.New(
		core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zap.ErrorLevel),
	)
	// zapLogger.Info("zap logger initialized")
	// zapLogger.Debug("zap logger debug level enabled")
}

func Default() *Config {
	return &Config{
		stdoutConfig: &StdoutConfig{
			level: DebugLevel,
		},
		rollingConfig: &FileConfig{
			encoding: zapcore.EncoderConfig{},
			logger: &lumberjack.Logger{
				Filename: getLogFilename(DefaultFilename),
				MaxSize:  500, // megabytes
				MaxAge:   30,  // days
			},
		},
		fieldsConfig: &FieldsConfig{},
	}
}

func (c *Config) WithFilename(filename string) *Config {
	c.rollingConfig.logger.Filename = filename
	return c
}

func (c *Config) WithLevel(level Level) *Config {
	c.stdoutConfig.level = level
	c.rollingConfig.level = level
	return c
}

func (c *Config) WithFields(fields []zapcore.Field) *Config {
	c.fieldsConfig.fields = fields
	return c
}

func getLogFilename(rawFilename string) string {
	if rawFilename == "" {
		return rawFilename
	}
	filename := filepath.Base(rawFilename)
	suffix := path.Ext(filename)
	filenameOnly := strings.TrimSuffix(filename, suffix)
	filenameOnly = fmt.Sprintf(filenameOnly+"_%v", time.Now().In(time.Local).Format("2006-01-01"))
	return strings.ReplaceAll(rawFilename, filename, filenameOnly+suffix)
}
