// Package logger
// Created by Teocci.
// Author: teocci@yandex.com on 2021-Aug-12
package logger

import (
	"fmt"
	"io"
	"log"
	"os"
)

const callDepth = 3

const (
	levelDebugName = "debug"
	levelErrorName = "error"
)

type Level int

const (
	LevelDebug Level = iota
	LevelWarn
	LevelError
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return levelDebugName
	case LevelError:
		return levelErrorName
	}
	return ""
}

var (
	formatLogger *Logger
	levelMap     = map[string]Level{
		levelDebugName: LevelDebug,
		levelErrorName: LevelError,
	}
)

func InitLog(l string) {
	formatLogger = newLog(os.Stdout).setLevel(levelMap[l]).setFlags(log.Lshortfile)
}

type Logger struct {
	level Level
	log   *log.Logger
}

func newLog(w io.Writer) *Logger {
	return &Logger{
		log: log.New(w, "", 0),
	}
}

func (l *Logger) setFlags(flag int) *Logger {
	l.log.SetFlags(flag)
	return l
}

func (l *Logger) setLevel(level Level) *Logger {
	l.level = level
	return l
}

func (l *Logger) doLog(level Level, v ...interface{}) bool {
	if level < l.level {
		return false
	}
	_ = l.log.Output(callDepth, level.String()+" "+fmt.Sprintln(v...))
	return true
}

func (l *Logger) doLogf(level Level, format string, v ...interface{}) bool {
	if level < l.level {
		return false
	}
	_ = l.log.Output(callDepth, level.String()+" "+fmt.Sprintln(fmt.Sprintf(format, v...)))
	return true
}

func Debug(v ...interface{}) {
	formatLogger.doLog(LevelDebug, v...)
}

func Warn(v ...interface{}) {
	formatLogger.doLog(LevelWarn, v...)
}

func Error(v ...interface{}) {
	formatLogger.doLog(LevelError, v...)
}

func Errorf(format string, v ...interface{}) {
	formatLogger.doLogf(LevelError, format, v...)
}

func Warnf(format string, v ...interface{}) {
	formatLogger.doLogf(LevelWarn, format, v...)
}

func Debugf(format string, v ...interface{}) {
	formatLogger.doLogf(LevelDebug, format, v...)
}
