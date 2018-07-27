package system

import (
	"github.com/sirupsen/logrus"
	"os"
	"github.com/natefinch/lumberjack"
	"fmt"
	"io"
)

func init() {
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
	})
}

func SetupLog() (err error) {
	parseLogConfig(Runtime.LogConfig)
	logrus.AddHook(&TraceHook{})
	return
}

func parseLogConfig(config LogConfig) {
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

type TraceHook struct {

}
var traceLevel = []logrus.Level{
	logrus.PanicLevel,
	logrus.FatalLevel,
	logrus.ErrorLevel,
	logrus.WarnLevel,
	logrus.InfoLevel,
}

func (trace *TraceHook) Levels() []logrus.Level {
	return traceLevel
}

func (trace *TraceHook) Fire(entry *logrus.Entry) (err error) {
	traceField, isPresent := entry.Data["trace"]
	if isPresent {
		if traceFlag, ok := traceField.(bool); ok {
			if traceFlag {
				traceLog := TraceLog{
					Code: entry.Data["code"].(string),
					Type: entry.Level.String(),
					RequestUrl:entry.Data["requestUrl"].(string),
					Message: entry.Message,
					Info: entry.Data["info"].(string),
				}
				if entry.Data["userId"] != nil {
					traceLog.UserId = entry.Data["userId"].(int64);
				}
				if entry.Data["user"] != nil {
					traceLog.User = entry.Data["user"].(string);
				}
				if entry.Data["tenantId"] != nil {
					traceLog.TenantId = entry.Data["tenantId"].(int64);
				}
				if entry.Data["tenant"] != nil {
					traceLog.Tenant = entry.Data["tenant"].(string);
				}
				go logToDatabase(&traceLog)
			}
		}
	}

	return
}

func logToDatabase(log *TraceLog)  {
	_, err := DbEngine.Insert(log)
	if err != nil {
		logrus.Error(err)
	} else {
		logrus.Infof("Log ID = %d is insert", log.Id)
	}
}

