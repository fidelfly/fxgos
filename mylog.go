package fxgos

import (
	"github.com/sirupsen/logrus"
	"os"
	"fmt"
	"io"
	"github.com/natefinch/lumberjack"
)

func init() {
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
	})
}

type LogConfig struct {
	LogLevel string
	LogPath string
	LogFile string
	MaxSize int
	Rotate string
	Stdout bool
}

func SetupLog(config LogConfig) {
	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		level = logrus.WarnLevel
	}
	logrus.SetLevel(level)

	if len(config.LogPath) == 0 {
		logrus.SetOutput(os.Stdout)
	} else {
		ljout := &lumberjack.Logger{
			Filename:   fmt.Sprintf("%s/%s", config.LogPath, config.LogFile),
			MaxSize:    config.MaxSize,
			MaxBackups: 3,
			MaxAge:     7,
			Compress:   true,
			LocalTime:  true,
		}

		if config.Stdout {
			logrus.SetOutput(io.MultiWriter(os.Stdout, ljout))
		} else {
			logrus.SetOutput(ljout)
		}
	}
	logrus.Info("Logrus is setup!")
}

