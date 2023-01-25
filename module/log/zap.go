package log

import (
	"fmt"
	"time"
	"os"
	. "wechatbot/config"
	"wechatbot/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)
var LOG *zap.Logger
var z = new(_zap)
type _zap struct{}
func InitZap() {
	if ok, _ := utils.PathExists(GlobalConfig.LogDirectory); !ok { // 判断是否有Director文件夹
		fmt.Printf("create %v directory\n", GlobalConfig.LogDirectory)
		_ = os.Mkdir(GlobalConfig.LogDirectory, os.ModePerm)
	}

	cores := z.getZapCores()
	logger := zap.New(zapcore.NewTee(cores...))

	if GlobalConfig.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}
	LOG = logger
	// return logger
}
// GetEncoder 获取 zapcore.Encoder
func (z *_zap) getEncoder() zapcore.Encoder {
	if GlobalConfig.LogFormat == "json" {
		return zapcore.NewJSONEncoder(z.getEncoderConfig())
	}
	return zapcore.NewConsoleEncoder(z.getEncoderConfig())
}

// GetEncoderConfig 获取zapcore.EncoderConfig
func (z *_zap) getEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "tracekey",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseColorLevelEncoder, // 小写编码器带颜色
		EncodeTime:     z.CustomTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
}
// GetEncoderCore 获取Encoder的 zapcore.Core
func (z *_zap) getEncoderCore(l zapcore.Level, level zap.LevelEnablerFunc) zapcore.Core {
	writer, err := FileRotatelogs.getWriteSyncer(l.String()) // 使用file-rotatelogs进行日志分割
	if err != nil {
		fmt.Printf("Get Write Syncer Failed err:%v", err.Error())
		return nil
	}

	return zapcore.NewCore(z.getEncoder(), writer, level)
}

// CustomTimeEncoder 自定义日志输出时间格式
func (z *_zap) CustomTimeEncoder(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
	encoder.AppendString(/* "wechatBot" +  */"[" + t.Format("2006/01/02 - 15:04:05.000") + "]")
}

// GetZapCores 根据配置文件的Level获取 []zapcore.Core
func (z *_zap) getZapCores() []zapcore.Core {
	cores := make([]zapcore.Core, 0, 7)
	for level := GetLogLevel(); level <= zapcore.FatalLevel; level++ {
		cores = append(cores, z.getEncoderCore(level, z.getLevelPriority(level)))
	}
	return cores
}

// GetLevelPriority 根据 zapcore.Level 获取 zap.LevelEnablerFunc
func (z *_zap) getLevelPriority(level zapcore.Level) zap.LevelEnablerFunc {
	switch level {
	case zapcore.DebugLevel:
		return func(level zapcore.Level) bool { // 调试级别
			return level == zap.DebugLevel
		}
	case zapcore.InfoLevel:
		return func(level zapcore.Level) bool { // 日志级别
			return level == zap.InfoLevel
		}
	case zapcore.WarnLevel:
		return func(level zapcore.Level) bool { // 警告级别
			return level == zap.WarnLevel
		}
	case zapcore.ErrorLevel:
		return func(level zapcore.Level) bool { // 错误级别
			return level == zap.ErrorLevel
		}
	case zapcore.DPanicLevel:
		return func(level zapcore.Level) bool { // dpanic级别
			return level == zap.DPanicLevel
		}
	case zapcore.PanicLevel:
		return func(level zapcore.Level) bool { // panic级别
			return level == zap.PanicLevel
		}
	case zapcore.FatalLevel:
		return func(level zapcore.Level) bool { // 终止级别
			return level == zap.FatalLevel
		}
	default:
		return func(level zapcore.Level) bool { // 调试级别
			return level == zap.DebugLevel
		}
	}
}
