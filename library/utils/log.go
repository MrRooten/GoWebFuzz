package utils

import (
	"errors"
	"fmt"
	"github.com/gookit/color"
	"path"
	"runtime"
	"strings"
	"time"
)

//日志级别
type LogLevel uint16

//Logger 日志结构体
type Logger struct {
	Level LogLevel
}

const (
	//定义日志级别
	UNKNOWN LogLevel = iota
	DEBUG
	INFO
	WARNING
	ERROR
	SUCC
)

var (
	ColorRedStart    = fmt.Sprintf("%c[%d;%d;%dm", 0x1B, 0, 40, 31)
	ColorYellowStart = fmt.Sprintf("%c[%d;%d;%dm", 0x1B, 0, 40, 33)
	ColorGreenStart  = fmt.Sprintf("%c[%d;%d;%dm", 0x1B, 0, 40, 32)
	ColorBlueStart   = fmt.Sprintf("%c[%d;%d;%dm", 0x1B, 0, 40, 34)
	ColorEnd         = fmt.Sprintf("%c[0m", 0x1B)
)

func parseLogLevel(s string) (LogLevel, error) {
	s = strings.ToLower(s)
	switch s {
	case "debug":
		return DEBUG, nil
	case "info":
		return INFO, nil
	case "warning":
		return WARNING, nil
	case "error":
		return ERROR, nil
	case "success":
		return SUCC, nil
	default:
		return UNKNOWN, errors.New("not found right LogLevel")
	}
}

//获取行号等信息
func getInfo(skip int) (funcName, fileName string, lineNo int) {
	pc, file, lineNo, ok := runtime.Caller(skip)
	if !ok {
		fmt.Printf("logging get line number error\n")
		return
	}
	funcName = runtime.FuncForPC(pc).Name()
	fileName = path.Base(file)
	return
}

//向终端写日志相关内容
//NewLog 构造函数
func NewLog(levelStr string) Logger {
	level, err := parseLogLevel(levelStr)
	if err != nil {
		panic(err)
	}
	return Logger{
		Level: level,
	}
}
func (l *Logger) SetLevel(levelStr string) {
	level, err := parseLogLevel(levelStr)
	if err != nil {
		panic(err)
	}
	l.Level = level
}
func (l *Logger) Debug(msg string) {
	if l.Level <= DEBUG {
		now := time.Now()
		_, fileName, lineNo := getInfo(2)
		//fmt.Printf("%c[%d;%d;%dm[DEB]%c[0m %s [%s:%d] %s\n", 0x1B, 0, 40, 36, 0x1B, now.Format("2006-01-02 15:04:05"), fileName, lineNo, msg)
		fmt.Printf("%s %s [%s:%d] %s\n", color.FgBlue.Render("[DEB]"), now.Format("2006-01-02 15:04:05"), fileName, lineNo, msg)
	}
}
func (l *Logger) Debugf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if l.Level <= DEBUG {
		now := time.Now()
		_, fileName, lineNo := getInfo(2)
		//fmt.Printf("%c[%d;%d;%dm[DEB]%c[0m %s [%s:%d] %s\n", 0x1B, 0, 40, 36, 0x1B, now.Format("2006-01-02 15:04:05"), fileName, lineNo, msg)
		fmt.Printf("%s %s [%s:%d] %s\n", color.FgBlue.Render("[DEB]"), now.Format("2006-01-02 15:04:05"), fileName, lineNo, msg)
	}
}
func (l *Logger) Info(msg string) {
	if l.Level <= INFO {
		now := time.Now()
		_, fileName, lineNo := getInfo(2)
		fmt.Printf("%c[%d;%d;%dm[INF]%c[0m %s [%s:%d] %s\n", 0x1B, 0, 40, 32, 0x1B, now.Format("2006-01-02 15:04:05"), fileName, lineNo, msg)
	}
}
func (l *Logger) Infof(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if l.Level <= INFO {
		now := time.Now()
		_, fileName, lineNo := getInfo(2)
		fmt.Printf("%c[%d;%d;%dm[INF]%c[0m %s [%s:%d] %s\n", 0x1B, 0, 40, 32, 0x1B, now.Format("2006-01-02 15:04:05"), fileName, lineNo, msg)
	}
}
func (l *Logger) Warning(msg string) {
	if l.Level <= WARNING {
		now := time.Now()
		_, fileName, lineNo := getInfo(2)
		fmt.Printf("%c[%d;%d;%dm[WRN]%c[0m %s [%s:%d] %s\n", 0x1B, 0, 40, 33, 0x1B, now.Format("2006-01-02 15:04:05"), fileName, lineNo, msg)
	}
}
func (l *Logger) Warningf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if l.Level <= WARNING {
		now := time.Now()
		_, fileName, lineNo := getInfo(2)
		fmt.Printf("%c[%d;%d;%dm[WRN]%c[0m %s [%s:%d] %s\n", 0x1B, 0, 40, 33, 0x1B, now.Format("2006-01-02 15:04:05"), fileName, lineNo, msg)
	}
}
func (l *Logger) Error(msg string) {
	if l.Level <= ERROR {
		now := time.Now()
		_, fileName, lineNo := getInfo(2)
		fmt.Printf("%c[%d;%d;%dm[ERR] %s [%s:%d] %s\n%c[0m", 0x1B, 0, 40, 31, now.Format("2006-01-02 15:04:05"), fileName, lineNo, msg, 0x1B)
	}
}
func (l *Logger) Errorf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if l.Level <= ERROR {
		now := time.Now()
		_, fileName, lineNo := getInfo(2)
		fmt.Printf("%c[%d;%d;%dm[ERR] %s [%s:%d] %s\n%c[0m", 0x1B, 0, 40, 31, now.Format("2006-01-02 15:04:05"), fileName, lineNo, msg, 0x1B)
	}
}
func (l *Logger) Success(msg string) {
	if l.Level <= SUCC {
		now := time.Now()
		_, fileName, lineNo := getInfo(2)
		fmt.Printf("%c[%d;%d;%dm[SUC] %s [%s:%d] %s\n%c[0m", 0x1B, 0, 40, 31, now.Format("2006-01-02 15:04:05"), fileName, lineNo, msg, 0x1B)
	}
}
func (l *Logger) Successf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if l.Level <= SUCC {
		//now := time.Now()
		//_, fileName, lineNo := getInfo(2)
		fmt.Printf("%c[%d;%d;%dm%s\n%c[0m", 0x1B, 0, 40, 31, msg, 0x1B)
	}
}

var Log = NewLog("info")
