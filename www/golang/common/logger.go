package common

import (
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func NewLogger() *logrus.Logger {
	if Log != nil {
		return Log
	}

	path := "./logs/go.log"
	writer, err := rotatelogs.New(
		path+".%Y%m%d",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(6*30*24*time.Hour),       // 6 月, default 7 day
		rotatelogs.WithRotationTime(86400*time.Second), // 1 天, default 86400 sec
	)

	if err != nil {
		panic(err)
	}

	pathMap := lfshook.WriterMap{
		logrus.InfoLevel:  writer,
		logrus.FatalLevel: writer,
		logrus.ErrorLevel: writer,
		logrus.DebugLevel: writer,
		logrus.WarnLevel:  writer,
		logrus.PanicLevel: writer,
		logrus.TraceLevel: writer,
	}

	logrus.AddHook(lfshook.NewHook(
		pathMap,
		&logrus.TextFormatter{},
	))

	Log = logrus.New()
	Log.Hooks.Add(lfshook.NewHook(
		pathMap,
		&logrus.TextFormatter{},
	))
	return Log
}
