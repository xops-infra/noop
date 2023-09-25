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
const HumanTime = "_human_time"

type Config struct {
	stdoutConfig          *StdoutConfig
	rollingConfig         *FileConfig
	fieldsConfig          *FieldsConfig
	levelFilterFileConfig *LevelFilterFileConfig
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

type LevelFilterFileConfig struct {
	warnLevelEnable  bool
	warnLogFilename  string
	errorLevelEnable bool
	errorLogFilename string
}

func (c *Config) Init() {
	var cores []zapcore.Core
	consoleEncoder := zap.NewDevelopmentEncoderConfig()
	consoleEncoder.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(consoleEncoder),
		zapcore.AddSync(zapcore.Lock(os.Stdout)),
		zapcore.Level(c.stdoutConfig.level),
	)

	fileCore := c.getCore(c.rollingConfig.logger.Filename, c.getSmallestLevelEnable())
	fileCore = c.setFields(fileCore)

	cores = append(cores, consoleCore, fileCore)

	if c.levelFilterFileConfig.warnLevelEnable {
		levelEnablerFunc := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			if c.levelFilterFileConfig.errorLevelEnable {
				return lvl == zapcore.WarnLevel
			}
			return lvl >= zapcore.WarnLevel
		})
		warnFileCore := c.getCore(c.levelFilterFileConfig.warnLogFilename, levelEnablerFunc)
		warnFileCore = c.setFields(warnFileCore)
		cores = append(cores, warnFileCore)
	}

	if c.levelFilterFileConfig.errorLevelEnable {
		errorFileCore := c.getCore(c.levelFilterFileConfig.errorLogFilename, func(lvl zapcore.Level) bool {
			return lvl >= zapcore.ErrorLevel
		})
		errorFileCore = c.setFields(errorFileCore)
		cores = append(cores, errorFileCore)
	}

	core := zapcore.NewTee(cores...)
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
				Filename: getLogFilename(DefaultFilename, ""),
				MaxSize:  500, // megabytes
				MaxAge:   30,  // days
			},
		},
		fieldsConfig:          &FieldsConfig{},
		levelFilterFileConfig: &LevelFilterFileConfig{},
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
		HumanTime: location,
	})
	return c
}

func (c *Config) WithWarnLog(optionWarnLogFilename string) *Config {
	c.levelFilterFileConfig.warnLevelEnable = true
	c.levelFilterFileConfig.warnLogFilename = getLogLevelFilterFilename(optionWarnLogFilename, "warn")
	return c
}

func (c *Config) WithErrorLog(optionErrorLogFilename string) *Config {
	c.levelFilterFileConfig.errorLevelEnable = true
	c.levelFilterFileConfig.errorLogFilename = getLogLevelFilterFilename(optionErrorLogFilename, "error")
	return c
}

func getLogFilename(rawFilename string, level string) string {
	if rawFilename == "" {
		return rawFilename
	}
	filename := filepath.Base(rawFilename)
	suffix := path.Ext(filename)
	filenameOnly := strings.TrimSuffix(filename, suffix)
	filenameOnly = fmt.Sprintf(filenameOnly + "_" + time.Now().In(time.Local).Format("2006-01-02"))
	if level != "" {
		level = "_" + level
	}
	return strings.ReplaceAll(rawFilename, filename, filenameOnly+level+suffix)
}

func getLogLevelFilterFilename(optionLogFilterFilename string, level string) string {
	if optionLogFilterFilename == "" {
		optionLogFilterFilename = getLogFilename(DefaultFilename, level)
	}
	return optionLogFilterFilename
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

func (c *Config) setFields(core zapcore.Core) zapcore.Core {
	if c.fieldsConfig.fields != nil && len(c.fieldsConfig.fields) != 0 {
		core = core.With(c.transformFields())
	}
	return core
}

func (c *Config) getCore(logFilename string, levelEnablerFunc zap.LevelEnablerFunc) zapcore.Core {
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename: logFilename,
		MaxSize:  c.rollingConfig.logger.MaxSize, // megabytes
		MaxAge:   c.rollingConfig.logger.MaxAge,  // days
	})
	fileEncoder := zap.NewProductionEncoderConfig()
	if timeLocation, ok := c.fieldsConfig.fields[HumanTime]; ok {
		fileEncoder.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.In(timeLocation.(*time.Location)).Format("2006-01-02 15:04:05.000"))
		}
		delete(c.fieldsConfig.fields, HumanTime)
	}
	return zapcore.NewCore(
		zapcore.NewJSONEncoder(fileEncoder),
		zapcore.AddSync(fileWriter),
		levelEnablerFunc,
	)
}

func (c *Config) getSmallestLevelEnable() zap.LevelEnablerFunc {
	return func(lvl zapcore.Level) bool {
		if c.levelFilterFileConfig.errorLevelEnable && c.levelFilterFileConfig.warnLevelEnable {
			return lvl < zapcore.WarnLevel
		} else if c.levelFilterFileConfig.errorLevelEnable {
			return lvl < zapcore.ErrorLevel
		} else if c.levelFilterFileConfig.warnLevelEnable {
			return lvl < zapcore.WarnLevel
		}
		return lvl >= zapcore.Level(c.rollingConfig.level)
	}
}
