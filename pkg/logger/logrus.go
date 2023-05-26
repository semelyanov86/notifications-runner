package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

func Debug(msg LogMessage) {
	logrus.Debug(msg)
}

func Debugf(format string, msg LogMessage) {
	logrus.Debugf(format, msg)
}

func Info(msg LogMessage) {
	logrus.Info(msg.Msg)
}

func Infof(format string, msg LogMessage) {
	logrus.Infof(format, msg)
}

func Warn(msg LogMessage) {
	logrus.Warn(msg)
}

func Warnf(format string, msg LogMessage) {
	logrus.Warnf(format, msg)
}

func Error(msg LogMessage) {
	logrus.Error(msg)
}

func WrapError(msg string, err error) {
	error := fmt.Errorf("%s: %w", msg, err)
	logrus.Error(ConvertErrorToStruct(error, 0, nil))
}

func Errorf(format string, msg LogMessage) {
	logrus.Errorf(format, msg)
}
