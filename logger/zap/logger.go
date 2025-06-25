package zap

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"time"
)

const logPath = "middleware/nacos/logs/nacos.log"

var log *Logger

type Logger struct {
	_logger *zap.Logger
}

func InitLogger() {
	encoder := NewEncoder()
	writerSyncer := NewLogWriter()
	core := zapcore.NewCore(encoder, writerSyncer, zapcore.InfoLevel)

	log = &Logger{_logger: zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))}
}

func NewLogWriter() zapcore.WriteSyncer {
	if err := os.MkdirAll(filepath.Dir(logPath), os.ModePerm); err != nil {
		panic(err)
	}
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	return zapcore.AddSync(file)
}

func NewEncoder() zapcore.Encoder {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(time.Format("2006-01-02 15:04:05"))
	}
	encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder

	return zapcore.NewConsoleEncoder(encoderCfg)
}

func (l *Logger) Info(msg string, kv ...interface{}) {
	l._logger.Info(msg, newFields(kv)...)
}

func (l *Logger) Debug(msg string, kv ...interface{}) {
	l._logger.Debug(msg, newFields(kv)...)
}

func (l *Logger) Warn(msg string, kv ...interface{}) {
	l._logger.Warn(msg, newFields(kv)...)
}

func (l *Logger) Error(msg string, kv ...interface{}) {
	l._logger.Error(msg, newFields(kv)...)
}

func newFields(kv []interface{}) []zap.Field {
	if len(kv)%2 != 0 {
		kv = append(kv, "<nil>")
	}
	fields := make([]zap.Field, 0, len(kv)/2)
	for i := 0; i < len(kv); i += 2 {
		k := fmt.Sprintf("%v", kv[i])
		fields = append(fields, zap.Any(k, kv[i+1]))
	}
	return fields
}

func GetLogger() *Logger {
	return log
}
