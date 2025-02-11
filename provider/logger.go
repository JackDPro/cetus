package provider

import (
	"github.com/JackDPro/cetus/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"sync"
)

var loggerInstance *zap.SugaredLogger
var loggerOnce sync.Once

func GetLogger() *zap.SugaredLogger {
	loggerOnce.Do(func() {
		logConf := config.GetLogConfig()
		appConf := config.GetAppConfig()
		cfg := zap.NewProductionEncoderConfig()
		cfg.EncodeTime = zapcore.ISO8601TimeEncoder
		cfg.EncodeLevel = zapcore.CapitalLevelEncoder

		// 添加输出通道
		var writer io.Writer
		if logConf.FileOut {
			file, err := os.OpenFile(logConf.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				panic(err)
			}
			if logConf.ConsoleOut {
				writer = io.MultiWriter(file, os.Stdout)
			} else {
				writer = io.MultiWriter(file)
			}
		} else {
			writer = io.MultiWriter(os.Stdout)
		}
		writeSyncer := zapcore.AddSync(writer)

		// 修改输出格式
		encoder := zapcore.NewConsoleEncoder(cfg)
		if logConf.Format == "json" {
			encoder = zapcore.NewJSONEncoder(cfg)
		}

		var level zapcore.Level
		switch logConf.Level {
		case "debug":
			level = zap.DebugLevel
		case "info":
			level = zap.InfoLevel
		case "warn":
			level = zap.WarnLevel
		case "error":
			level = zap.ErrorLevel
		}
		core := zapcore.NewCore(encoder, writeSyncer, level)

		// 添加统一字段
		zapLogger := zap.New(core, zap.AddCaller()).With(zap.String("app", appConf.Name), zap.String("env", appConf.Env))
		loggerInstance = zapLogger.Sugar()
	})
	return loggerInstance
}
