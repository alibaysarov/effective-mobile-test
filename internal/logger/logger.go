package logger

import (
	"io"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
)

var Logger *slog.Logger

func InitLogger(service string, file *os.File) (*slog.Logger, error) {

	logger := slog.New(slog.NewJSONHandler(io.MultiWriter(file, os.Stdout), nil))

	Logger = slog.With("gin_mode", gin.EnvGinMode)
	slog.SetDefault(logger)

	return logger, nil
}
