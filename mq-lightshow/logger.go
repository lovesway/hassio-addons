package main

import (
	"encoding/json"
	"os"
	"strings"
	"syscall"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// GetLogger sets up the logger.
func GetLogger(logLevel string) *zap.SugaredLogger {
	// Set up the logger
	rawJSON := []byte(`{
			"encoding": "console",
			"outputPaths": ["stdout"],
			"errorOutputPaths": ["stderr"],
			"encoderConfig": {
			  "timeKey": "time",
			  "timeEncoder": "ISO8601",
			  "messageKey": "message",
			  "levelKey": "level",
			  "levelEncoder": "capital"
			}
		  }`)

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}

	switch strings.ToLower(logLevel) {
	case "debug":
		cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		// Add caller when in debug mode
		cfg.EncoderConfig.CallerKey = "caller"
		cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	case "info":
		cfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warning":
		cfg.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		cfg.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	default:
		cfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	return logger.Sugar()
}

func loggerSync() {
	err := log.Sync()
	// Currently Sync() on stdout and stderr return errors on Linux and macOS respectively:
	// - sync /dev/stdout: invalid argument
	// - sync /dev/stdout: inappropriate ioctl for device
	// Since these are not actionable ignore them.
	if osErr, ok := err.(*os.PathError); ok {
		wrappedErr := osErr.Unwrap()
		switch wrappedErr {
		case syscall.EINVAL, syscall.ENOTSUP, syscall.ENOTTY:
			err = nil
		}
	}

	if err != nil {
		log.Error(err)
	}
}
