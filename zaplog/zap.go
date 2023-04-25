package zaplog

import (
	"log"
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Writer 日志输出类型
type Writer string

const (
	// Console 输出到控制台
	Console Writer = "console"

	// File 输出到文件
	File Writer = "file"

	// ConsoleFile 同时输出到控制台和文件
	ConsoleFile Writer = "console_file"
)

// Conf 日志相关配置
type Conf struct {
	Writer     Writer `yaml:"writer"`     // 日志输出方式（console-控制台输出、file-文件输出、console_file-文件控制台都输出）
	Level      string `yaml:"level"`      // 日志级别（debug、info、warn、error、dpanic、panic、fatal）
	Filename   string `yaml:"filename"`   // 日志文件的位置
	MaxBackups int    `yaml:"maxBackups"` // 保留旧文件的最大个数
	MaxAge     int    `yaml:"aaxAge"`     // 保留旧文件的最大天数
}

// NewZapLogger 新建建zap.SugaredLogger
func NewZapLogger(conf Conf) *zap.Logger {
	levelName := conf.Level

	var level = zapcore.DebugLevel
	if err := level.Set(levelName); err != nil {
		log.Fatalf("logger level name [%s] is wrong", levelName)
	}

	core := zapcore.NewCore(getEncoder(), getLogWriter(conf), level)
	return zap.New(core, zap.AddCaller())
}

func NewSugerZapLogger(conf Conf) *zap.SugaredLogger {
	return NewZapLogger(conf).Sugar()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(conf Conf) zapcore.WriteSyncer {
	hook, err := rotatelogs.New(
		conf.Filename+".%Y%m%d",
		rotatelogs.WithLinkName(conf.Filename),
		rotatelogs.WithMaxAge(time.Duration(conf.MaxAge*24)*time.Hour),
		rotatelogs.WithRotationCount(uint(conf.MaxBackups)),
	)
	if err != nil {
		log.Fatalf("logger rotate init fail, err: %+v", err)
	}

	if conf.Writer == ConsoleFile {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(hook))
	}

	if conf.Writer == File {
		return zapcore.AddSync(hook)
	}
	// 默认输出控制台
	return zapcore.AddSync(os.Stdout)
}
