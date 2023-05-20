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
const HumanTime = "human_time"

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
	fields map[string]any
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
		fileCore = fileCore.With(c.transformFields())
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

func (c *Config) WithFields(fields map[string]any) *Config {
	c.addFields(fields)
	return c
}

func (c *Config) WithHumanTime(location *time.Location) *Config {
	if location == nil {
		location = time.Local
	}
	c.addFields(map[string]any{
		HumanTime: time.Now().In(location).Format("2006-01-02 15:04:05.000"),
	})
	return c
}

func getLogFilename(rawFilename string) string {
	if rawFilename == "" {
		return rawFilename
	}
	filename := filepath.Base(rawFilename)
	suffix := path.Ext(filename)
	filenameOnly := strings.TrimSuffix(filename, suffix)
	filenameOnly = fmt.Sprintf(filenameOnly + "_" + time.Now().In(time.Local).Format("2006-01-02"))
	return strings.ReplaceAll(rawFilename, filename, filenameOnly+suffix)
}

func (c *Config) transformFields() []zapcore.Field {
	var zapFields []zapcore.Field
	for k, v := range c.fieldsConfig.fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return zapFields
}

func (c *Config) addFields(fields map[string]any) {
	if c.fieldsConfig.fields == nil {
		c.fieldsConfig.fields = make(map[string]any)
	}

	for k, v := range fields {
		c.fieldsConfig.fields[k] = v
	}
}
