package server

type Logger interface {
	Error(v ...interface{})
	Errorf(format string, v ...interface{})

	Warn(v ...interface{})
	Warnf(format string, v ...interface{})

	Info(v ...interface{})
	Infof(format string, v ...interface{})

	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
}

type nullLogger struct {
}

func (nullLogger) Error(v ...interface{}) {
}

func (nullLogger) Errorf(format string, v ...interface{}) {
}

func (nullLogger) Warn(v ...interface{}) {
}

func (nullLogger) Warnf(format string, v ...interface{}) {
}

func (nullLogger) Info(v ...interface{}) {
}

func (nullLogger) Infof(format string, v ...interface{}) {
}

func (nullLogger) Debug(v ...interface{}) {
}

func (nullLogger) Debugf(format string, v ...interface{}) {
}
