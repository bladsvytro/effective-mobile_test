package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger глобальный экземпляр логгера
var Logger *zap.Logger

// Init инициализирует логгер на основе конфигурации
func Init(level, encoding string) error {
	var cfg zap.Config

	if encoding == "json" {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
	}

	// Устанавливаем уровень логирования
	var lvl zapcore.Level
	if err := lvl.UnmarshalText([]byte(level)); err != nil {
		lvl = zapcore.InfoLevel
	}
	cfg.Level = zap.NewAtomicLevelAt(lvl)

	// Настраиваем вывод
	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stderr"}

	logger, err := cfg.Build()
	if err != nil {
		return err
	}

	Logger = logger
	return nil
}

// Get возвращает инициализированный логгер
func Get() *zap.Logger {
	if Logger == nil {
		// fallback на логгер по умолчанию
		Logger, _ = zap.NewProduction()
	}
	return Logger
}

// Sync выполняет синхронизацию логгера (вызывать при завершении приложения)
func Sync() {
	_ = Logger.Sync()
}
