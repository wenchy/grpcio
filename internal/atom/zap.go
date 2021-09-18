package atom

import (
	"errors"
	"log"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var levelMap = map[string]zapcore.Level{
	"debug": zapcore.DebugLevel,
	"info":  zapcore.InfoLevel,
	"warn":  zapcore.WarnLevel,
	"error": zapcore.ErrorLevel,
}

var zaplogger *zap.Logger

func GetZapLogger() *zap.Logger {
	return zaplogger
}

func (atom *Atom) InitZap(level string, dir string) error {
	zapLevel, ok := levelMap[level]
	if !ok {
		log.Fatalf("illegal log level: %s", level)
		return errors.New("illegal log level")
	}
	writeSyncer := getLogWriter(dir)
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapLevel)

	zaplogger = zap.New(core, zap.AddCaller())
	sugar := zaplogger.Sugar()

	// sugar.Infow("sugar log test1",
	// 	"url", "http://example.com",
	// 	"attempt", 3,
	// 	"backoff", time.Second,
	// )

	// sugar.Infof("sugar log test2: %s", "http://example.com")
	atom.Log = sugar

	return nil
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.FunctionKey = "func"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.ConsoleSeparator = "|"
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(dir string) zapcore.WriteSyncer {
	// lumberJackLogger := &lumberjack.Logger{
	// 	Filename: dir + "/log" + time.Now().Format("20060102"),
	// 	MaxSize:  100, // megabytes
	// 	// MaxBackups: 70, // The default is not to remove old log files based on age.
	// 	MaxAge:    7,     // days
	// 	LocalTime: true,  // The default is to user UTC time.
	// 	Compress:  false, // disabled by default
	// }
	// return zapcore.AddSync(lumberJackLogger)
	hook, err := rotatelogs.New(
		dir+"/log"+".%Y%m%d",
		rotatelogs.WithLinkName(dir+"/log"),
		rotatelogs.WithMaxAge(time.Hour*24*7),
		rotatelogs.WithRotationTime(time.Hour*24),
	)

	if err != nil {
		log.Printf("failed to create rotatelogs: %s", err)
	}
	return zapcore.AddSync(hook)
}
