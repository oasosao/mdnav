package zap

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger = *zap.Logger
type Field = zapcore.Field

// 创建生产环境Logger
func NewLogger() Logger {

	// 日志轮转配置
	lumberjackLogger := &lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    10,   // MB
		MaxBackups: 30,   // 保留30个备份
		MaxAge:     90,   // 保留90天
		Compress:   true, // 压缩旧日志
	}

	// 编码器配置
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.StacktraceKey = "stacktrace"

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zap.InfoLevel)

	// 核心配置：同时输出到文件和控制台
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(lumberjackLogger),
			// zapcore.AddSync(os.Stdout), // 开发时也输出到控制台
		),
		atomicLevel,
	)

	// 创建Logger
	return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(0))
}

func Error(err error) Field {
	return zap.Error(err)
}

func String(key, val string) Field {
	return zap.String(key, val)
}

func Int(key string, val int) Field {
	return zap.Int(key, val)
}

func Duration(key string, val time.Duration) Field {
	return zap.Duration(key, val)
}
