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

// customCallerEncoder personaliza el formato del caller para mostrar solo la ruta desde src/
func customCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	fullPath := caller.TrimmedPath()

	// Buscar la posición de "src/" en la ruta
	srcIndex := strings.Index(fullPath, "src/")
	if srcIndex != -1 {
		// Si encontramos "src/", mostrar solo desde ahí
		shortPath := fullPath[srcIndex:]
		enc.AppendString(shortPath)
	} else {
		// Si no encontramos "src/", mostrar solo el nombre del archivo
		parts := strings.Split(fullPath, "/")
		if len(parts) > 0 {
			enc.AppendString(parts[len(parts)-1])
		} else {
			enc.AppendString(fullPath)
		}
	}
}

func NewLogger() (*Logger, error) {
	// Configurar timezone de CDMX
	cdmxLocation, err := time.LoadLocation("America/Mexico_City")
	if err != nil {
		// Si no se puede cargar la timezone de CDMX, usar UTC
		cdmxLocation = time.UTC
	}

	// Configuración personalizada del encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "ts"
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		// Convertir a timezone de CDMX y formatear como RFC3339
		cdmxTime := t.In(cdmxLocation)
		enc.AppendString(cdmxTime.Format(time.RFC3339))
	}
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = customCallerEncoder

	// Crear el logger con la configuración personalizada
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
		// Para logs de HTTP requests, no mostrar caller ya que siempre será desde el middleware
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
			SlowThreshold:             time.Second, // umbral para destacar consultas lentas
			LogLevel:                  gormlogger.Error,
			IgnoreRecordNotFoundError: true, // no loguear "record not found"
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
