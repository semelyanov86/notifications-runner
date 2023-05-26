package logger

import "strconv"

type Logger interface {
	Debug(msg LogMessage)
	Info(msg LogMessage)
	Warn(msg LogMessage)
	Error(msg LogMessage)
}

type LogMessage struct {
	Msg        string
	Code       string
	Properties map[string]string
}

func ConvertErrorToStruct(err error, code int, properties map[string]string) LogMessage {
	if code == 0 {
		code = 500
	}
	return LogMessage{
		Msg:        err.Error(),
		Code:       strconv.Itoa(code),
		Properties: properties,
	}
}

func GenerateErrorMessageFromString(message string) LogMessage {
	return LogMessage{
		Msg:        message,
		Code:       "500",
		Properties: nil,
	}
}
