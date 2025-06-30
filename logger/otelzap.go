package logger

import (
	"context"
	"os"
	"skytakeout/global"
	"sync"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	once sync.Once
)

// Init初始化全局 logger
func Init(isDev, path string) {
	once.Do(func() {
		// 控制台编码配置
		consoleEncoderCfg := zap.NewDevelopmentEncoderConfig()
		consoleEncoderCfg.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
		consoleEncoderCfg.TimeKey = "time"
		consoleEncoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder // 彩色

		// 文件编码配置
		fileEncoderCfg := zap.NewProductionEncoderConfig()
		fileEncoderCfg.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
		fileEncoderCfg.TimeKey = "time"
		fileEncoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder // 普通大写，不加颜色

		// 日志文件
		file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		fileCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(fileEncoderCfg),
			zapcore.AddSync(file),
			zapcore.ErrorLevel,
		)

		// 控制台输出
		consoleCore := zapcore.NewCore(
			zapcore.NewConsoleEncoder(consoleEncoderCfg),
			zapcore.AddSync(os.Stdout),
			zapcore.DebugLevel,
		)

		tee := zapcore.NewTee(consoleCore, fileCore)
		l := zap.New(tee, zap.AddCaller(), zap.AddCallerSkip(1))

		global.ZapLog = otelzap.New(l)
	})
}

// Logger 返回带 context 的 otelzap logger
func Logger(ctx context.Context) otelzap.LoggerWithCtx {
	if global.ZapLog == nil {
		panic("logger not initialized, please call InitLogger first")
	}
	return global.ZapLog.Ctx(ctx)
}
