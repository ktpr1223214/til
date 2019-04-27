package main

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

func OpenFile(filePath string) error {
	_, err := os.Open(filePath)
	if err != nil {
		return errors.Wrapf(err, "failed to open file %s", filePath)
	}
	return nil
}

func InitLogger() (*zap.Logger, error) {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "time"
	encoderCfg.MessageKey = "message"
	encoderCfg.EncodeTime = LocalTimeEncoder

	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig = encoderCfg

	return cfg.Build()
}

func LocalTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006/01/02 15:04:05 MST"))
}

func main() {
	// logger, _ := zap.NewProduction()
	logger, _ := InitLogger()

	// field 追加
	logger.Info("Hello zap", zap.String("key", "value"), zap.Time("now", time.Now()))

	// error field
	logger.Error("error test", zap.Error(OpenFile("not_exist_file.txt")))

	// 現在 logger をクローンし、指定したフィールドを保持した新しいロガーを返す
	log := logger.With(zap.String("userId", "fuga"), zap.String("requestId", "piyo"))
	log.Info("hello world")
}
