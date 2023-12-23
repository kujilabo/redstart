package config

import (
	"fmt"
	"log/slog"
	"strings"

	libdomain "github.com/kujilabo/redstart/lib/domain"
	liblog "github.com/kujilabo/redstart/lib/log"
)

type LogConfig struct {
	Level map[string]string `yaml:"level"`
}

func InitLog(cfg *LogConfig) error {
	defaultLogLevel := slog.LevelWarn
	if rootLevel, ok := cfg.Level["default"]; ok {
		defaultLogLevel = stringToLogLevel(rootLevel)
	}

	liblog.DefaultLogLevel = defaultLogLevel

	for name, level := range cfg.Level {
		logLevel := stringToLogLevel(level)
		liblog.Loggers[libdomain.ContextKey(name)] = slog.New(liblog.LogHandlers[logLevel])
	}

	return nil
}

func stringToLogLevel(value string) slog.Level {
	switch strings.ToLower(value) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		slog.Info(fmt.Sprintf("Unsupported log level: %s", value))
		return slog.LevelWarn
	}
}
