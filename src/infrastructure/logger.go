package infrastructure

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	gormlogger "gorm.io/gorm/logger"
)

type Logger struct {
	Log *zap.Logger
}

// customCallerEncoder formats the caller path to show only the portion starting from src/
func customCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	fullPath := caller.TrimmedPath()

	// Find the position of "src/" in the path
	srcIndex := strings.Index(fullPath, "src/")
	if srcIndex != -1 {
		// If "src/" is found, keep only that part of the path
		shortPath := fullPath[srcIndex:]
		enc.AppendString(shortPath)
	} else {
		// If "src/" is not found, show only the file name
		parts := strings.Split(fullPath, "/")
		if len(parts) > 0 {
			enc.AppendString(parts[len(parts)-1])
		} else {
			enc.AppendString(fullPath)
		}
	}
}

func NewLogger() (*Logger, error) {
	// Configure Mexico City timezone
	cdmxLocation, err := time.LoadLocation("America/Mexico_City")
	if err != nil {
		// Use UTC if the Mexico City timezone cannot be loaded
		cdmxLocation = time.UTC
	}

	// Custom encoder configuration
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "ts"
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		// Convert to Mexico City timezone and format as RFC3339
		cdmxTime := t.In(cdmxLocation)
		enc.AppendString(cdmxTime.Format(time.RFC3339))
	}
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = customCallerEncoder

	// Create the logger with the custom configuration
	config := zap.NewProductionConfig()
	config.EncoderConfig = encoderConfig
	config.Development = false
	config.DisableCaller = false
	config.DisableStacktrace = true
	logger, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}
	return &Logger{Log: logger}, nil
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.Log.Info(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.Log.Error(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.Log.Warn(msg, fields...)
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.Log.Debug(msg, fields...)
}

func (l *Logger) GinZapLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		// Skip caller information for HTTP requests since it always originates from middleware
		l.Log.WithOptions(zap.AddCallerSkip(1)).Info("HTTP request", zap.String("method", c.Request.Method), zap.String("path", c.Request.URL.Path), zap.Int("status", c.Writer.Status()), zap.Duration("latency", latency), zap.String("client_ip", c.ClientIP()))
	}
}

type GormZapLogger struct {
	zap    *zap.SugaredLogger
	config gormlogger.Config
}

func NewGormLogger(base *zap.Logger) *GormZapLogger {
	sugar := base.Sugar()
	return &GormZapLogger{
		zap: sugar,
		config: gormlogger.Config{
			SlowThreshold:             time.Second, // threshold to highlight slow queries
			LogLevel:                  gormlogger.Error,
			IgnoreRecordNotFoundError: true, // do not log "record not found"
			Colorful:                  false,
		},
	}
}

func (l *GormZapLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newCfg := l.config
	newCfg.LogLevel = level
	return &GormZapLogger{zap: l.zap, config: newCfg}
}

func (l *GormZapLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.config.LogLevel >= gormlogger.Info {
		l.zap.Infof(msg, data...)
	}
}

func (l *GormZapLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.config.LogLevel >= gormlogger.Warn {
		l.zap.Warnf(msg, data...)
	}
}

func (l *GormZapLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.config.LogLevel >= gormlogger.Error &&
		(!l.config.IgnoreRecordNotFoundError || msg != gormlogger.ErrRecordNotFound.Error()) {
		l.zap.Errorf(msg, data...)
	}
}

func (l *GormZapLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)

	if err != nil {
		if l.config.IgnoreRecordNotFoundError && errors.Is(err, gormlogger.ErrRecordNotFound) {
			return
		}
		if l.config.LogLevel >= gormlogger.Error {
			sql, rows := fc()
			l.zap.Errorf("Error: %v | %.3fms | rows:%d | %s", err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
		return
	}

	if elapsed > l.config.SlowThreshold && l.config.LogLevel >= gormlogger.Warn {
		sql, rows := fc()
		l.zap.Warnf("SLOW ≥ %s | %.3fms | rows:%d | %s", l.config.SlowThreshold, float64(elapsed.Nanoseconds())/1e6, rows, sql)
	}
}
