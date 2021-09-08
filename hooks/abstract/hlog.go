package abstract

import (
	"github.com/adminhmi/hlog"
)

type Logger interface {
	DebugLogger
	InfoLogger
	WarnLogger
	ErrorLogger
	FatalLogger
	PanicLogger
	LevelLogger(level Level) LevelLogger
}
type DebugLogger interface {
	Debug(msg string, fields ...Field)
}

type ErrorLogger interface {
	Error(msg string, fields ...Field)
}

type FatalLogger interface {
	Fatal(msg string, fields ...Field)
}

type InfoLogger interface {
	Info(msg string, fields ...Field)
}

type PanicLogger interface {
	Panic(msg string, fields ...Field)
}

type WarnLogger interface {
	Warn(msg string, fields ...Field)
}
type LevelLogger interface {
	Println(v ...interface{})
	Printf(format string, v ...interface{})
}

func NewHlogLogger(l *hlog.Logger, level Level) *hlogLogger {
	return &hlogLogger{
		l:          l,
		levelCheck: NewLevelCheck(level),
	}
}

//  implements the Logger frontend using the popular  library as a backend
// It makes use of the LevelCheck helper to increase performance
type hlogLogger struct {
	l          *hlog.Logger
	levelCheck LevelCheck
}

func (l *hlogLogger) LevelLogger(level Level) LevelLogger {
	return &hlogLevelLogger{
		l:     l.l,
		level: level,
	}
}

func (l *hlogLogger) fields(fields []Field) hlog.Fields {
	out := make(hlog.Fields, len(fields))
	for i := range fields {
		switch fields[i].kind {
		case StringField:
			out[fields[i].key] = fields[i].stringValue
		case ByteStringField:
			out[fields[i].key] = string(fields[i].byteValue)
		case IntField:
			out[fields[i].key] = fields[i].intValue
		case BoolField:
			out[fields[i].key] = fields[i].intValue != 0
		case ErrorField, NamedErrorField:
			out[fields[i].key] = fields[i].errorValue
		case StringsField:
			out[fields[i].key] = fields[i].stringsValue
		default:
			out[fields[i].key] = fields[i].interfaceValue
		}
	}
	return out
}

func (l *hlogLogger) Debug(msg string, fields ...Field) {
	if !l.levelCheck.Check(DebugLevel) {
		return
	}
	l.l.WithFields(l.fields(fields)).Debug(msg)
}

func (l *hlogLogger) Info(msg string, fields ...Field) {
	if !l.levelCheck.Check(InfoLevel) {
		return
	}
	l.l.WithFields(l.fields(fields)).Info(msg)
}

func (l *hlogLogger) Warn(msg string, fields ...Field) {
	if !l.levelCheck.Check(WarnLevel) {
		return
	}
	l.l.WithFields(l.fields(fields)).Warn(msg)
}

func (l *hlogLogger) Error(msg string, fields ...Field) {
	if !l.levelCheck.Check(ErrorLevel) {
		return
	}
	l.l.WithFields(l.fields(fields)).Error(msg)
}

func (l *hlogLogger) Fatal(msg string, fields ...Field) {
	if !l.levelCheck.Check(FatalLevel) {
		return
	}
	l.l.WithFields(l.fields(fields)).Fatal(msg)
}

func (l *hlogLogger) Panic(msg string, fields ...Field) {
	if !l.levelCheck.Check(PanicLevel) {
		return
	}
	l.l.WithFields(l.fields(fields)).Panic(msg)
}

type hlogLevelLogger struct {
	l     *hlog.Logger
	level Level
}

func (s *hlogLevelLogger) Println(v ...interface{}) {
	switch s.level {
	case DebugLevel:
		s.l.Debug(v...)
	case InfoLevel:
		s.l.Info(v...)
	case WarnLevel:
		s.l.Warn(v...)
	case ErrorLevel:
		s.l.Error(v...)
	case FatalLevel:
		s.l.Fatal(v...)
	case PanicLevel:
		s.l.Panic(v...)
	}
}

func (s *hlogLevelLogger) Printf(format string, v ...interface{}) {
	switch s.level {
	case DebugLevel:
		s.l.Debugf(format, v...)
	case InfoLevel:
		s.l.Infof(format, v...)
	case WarnLevel:
		s.l.Warnf(format, v...)
	case ErrorLevel:
		s.l.Errorf(format, v...)
	case FatalLevel:
		s.l.Fatalf(format, v...)
	case PanicLevel:
		s.l.Panicf(format, v...)
	}
}
