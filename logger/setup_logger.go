package logger

import (
	"GinTalk/settings"
	"go.uber.org/zap"
)

// SetupGlobalLogger 用于初始化全局日志
func SetupGlobalLogger(cfg *settings.LoggerConfig) error {
	logger, err := newLogger(cfg)
	if err != nil {
		return err
	}
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			logger.Error("Failed to sync logger", zap.Error(err))
		}
	}(logger)

	zap.ReplaceGlobals(logger)
	return nil
}
