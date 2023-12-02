package webutils

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "gopkg.in/natefinch/lumberjack.v2"
    "os"
    "time"
)

var (
    TimeFmtMS = "2006-01-02 15:04:05.000"
    TimeFmt   = "2006-01-02 15:04:05"
)

func createLog(logfile string) *zap.SugaredLogger {
    encoderConfig := zap.NewProductionEncoderConfig()
    customTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
        enc.AppendString(t.Format(TimeFmtMS))
    }
    encoderConfig.EncodeTime = customTimeEncoder
    encoder := zapcore.NewConsoleEncoder(encoderConfig)

    infoLog := &lumberjack.Logger{
        Filename:   logfile,
        MaxSize:    10,
        MaxAge:     30,
        MaxBackups: 10,
        Compress:   true,
    }
    infoWriteSyncer := zapcore.AddSync(infoLog)
    core := zapcore.NewTee(
        zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.InfoLevel),
        zapcore.NewCore(encoder, infoWriteSyncer, zapcore.InfoLevel),
    )

    return zap.New(core, zap.AddCaller()).Sugar()
}
