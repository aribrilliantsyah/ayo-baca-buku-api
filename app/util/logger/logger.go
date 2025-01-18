package logger

import (
	"context"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/gorm/logger"
)

var (
	globalLogger *zap.Logger
	once         sync.Once
)

func NewLogger() *zap.Logger {
	once.Do(func() {
		logFile := &lumberjack.Logger{
			Filename:   "logs/app-" + time.Now().Format("2006-01-02") + ".log",
			MaxSize:    10,
			MaxBackups: 30,
			MaxAge:     7,
			Compress:   true,
		}

		// Encoder configuration
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.TimeKey = "timestamp"
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		// Create core for file logging
		fileCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(logFile),
			zap.InfoLevel,
		)

		// Create core for console logging
		consoleCore := zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			zap.InfoLevel,
		)

		// Combine both cores
		core := zapcore.NewTee(fileCore, consoleCore)

		globalLogger = zap.New(core)
	})
	return globalLogger
}

func GetLogger() *zap.Logger {
	if globalLogger == nil {
		return NewLogger()
	}
	return globalLogger
}

type GormLogger struct {
	ZapLogger *zap.Logger
	LogLevel  logger.LogLevel
}

func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return &GormLogger{
		ZapLogger: l.ZapLogger,
		LogLevel:  level,
	}
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		l.ZapLogger.Sugar().Infof(msg, data...)
	}
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		l.ZapLogger.Sugar().Warnf(msg, data...)
	}
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		l.ZapLogger.Sugar().Errorf(msg, data...)
	}
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel > 0 {
		elapsed := time.Since(begin)
		sql, rows := fc()

		switch {
		case err != nil && l.LogLevel >= logger.Error:
			l.ZapLogger.Error("SQL error",
				zap.Error(err),
				zap.String("sql", sql),
				zap.Int64("rows", rows),
				zap.Duration("elapsed", elapsed),
			)
		case l.LogLevel >= logger.Info:
			l.ZapLogger.Info("SQL executed",
				zap.String("sql", sql),
				zap.Int64("rows", rows),
				zap.Duration("elapsed", elapsed),
			)
		}
	}
}
