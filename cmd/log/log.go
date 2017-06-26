package log

import (
	"fmt"
	"strings"
	"time"
)

type LogLevel int

const (
	ERROR LogLevel = iota
	WARNING
	INFO
	DEBUG
)

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

type _Log struct {
	logLevel   LogLevel
	timeFormat string
}

func NewLogger(logLevel LogLevel) Logger {
	return &_Log{
		logLevel:   logLevel,
		timeFormat: "15:04:05.000",
	}
}

func (self *_Log) Error(v ...interface{}) {
	if self.logLevel >= ERROR {
		self.out(" Err ", fmt.Sprint(v...))
	}
}

func (self *_Log) Errorf(format string, v ...interface{}) {
	if self.logLevel >= ERROR {
		self.out(" Err ", format, v...)
	}
}

func (self *_Log) Warn(v ...interface{}) {
	if self.logLevel >= WARNING {
		self.out(" Wrn ", fmt.Sprint(v...))
	}
}

func (self *_Log) Warnf(format string, v ...interface{}) {
	if self.logLevel >= WARNING {
		self.out(" Wrn ", format, v...)
	}
}

func (self *_Log) Info(v ...interface{}) {
	if self.logLevel >= INFO {
		self.out(" Inf ", fmt.Sprint(v...))
	}
}

func (self *_Log) Infof(format string, v ...interface{}) {
	if self.logLevel >= INFO {
		self.out(" Inf ", format, v...)
	}
}

func (self *_Log) Debug(v ...interface{}) {
	if self.logLevel >= DEBUG {
		self.out(" Dbg ", fmt.Sprint(v...))
	}
}

func (self *_Log) Debugf(format string, v ...interface{}) {
	if self.logLevel >= DEBUG {
		self.out(" Dbg ", format, v...)
	}
}

func (self *_Log) out(level string, format string, v ...interface{}) {
	fmt.Printf(
		strings.Join(
			[]string{
				time.Now().Format(self.timeFormat),
				level,
				format,
				"\n",
			},
			""),
		v...,
	)
}
