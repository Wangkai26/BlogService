package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"runtime"
	"time"
)

type Level int8
type Fields map[string]interface{}

const(
	LevelDubug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelPanic
)

func (l Level) String() string {
	switch l {
	case LevelDubug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	case LevelFatal:
		return "fatal"
	case LevelPanic:
		return "panic"
	}
	return ""
}


type Logger struct {
	newLogger *log.Logger
	ctx 	  context.Context
	fields    Fields
	callers   []string
}

func NewLogger(w io.Writer,prefix string,flag int) *Logger  {
	l := log.New(w,prefix,flag)
	return &Logger{newLogger: l}
}

func (l *Logger) clone() *Logger {
	nl := *l
	return &nl
}

func (l *Logger) WithFields(f Fields) *Logger {
	ll := l.clone()
	if ll.fields == nil{
		ll.fields = make(Fields)
	}
	for k,v := range f{
		ll.fields[k] = v
	}
	return ll
}

func (l *Logger) WithContext(ctx context.Context) *Logger  {
	ll := l.clone()
	ll.ctx = ctx
	return ll
}

func (l *Logger) WithCaller(skip int) *Logger  {
	ll := l.clone()
	pc,file,line,ok := runtime.Caller(skip)
	if ok{
		f := runtime.FuncForPC(pc)
		ll.callers = []string{fmt.Sprintf("%s: %d %s",file,line,f.Name())}
	}

	return ll
}

func (l *Logger) WithCallersFrames() *Logger {
	maxCallerDepth := 25
	minCallerDepth := 1
	callers := []string{}
	pcs := make([]uintptr,maxCallerDepth)
	depth := runtime.Callers(minCallerDepth,pcs)
	frames := runtime.CallersFrames(pcs[:depth])
	for frame,more := frames.Next();more;frame,more=frames.Next(){
		callers = append(callers,fmt.Sprintf("%s: %d %s",frame.File,frame.Line,frame.Function))
		if !more {
			break
		}
	}
	ll := l.clone()
	ll.callers = callers
	return ll
}

// 日志内容格式化
func (l *Logger) JSONFormat(level Level,message string) map[string]interface{} {
	data := make(Fields,len(l.fields)+4)
	data["level"] = level.String()
	data["time"] = time.Now().Local().UnixNano()
	data["message"] = message
	data["callers"] = l.callers
	if len(l.fields) > 0{
		for k,v := range l.fields{
			if _,ok := data[k];!ok{
				data[k] = v
			}
		}
	}

	return data
}

// 日志输出动作
func (l *Logger) Output(level Level,message string)  {
	body,_ := json.Marshal(l.JSONFormat(level,message))
	content := string(body)
	switch level{
	case LevelDubug:
		l.newLogger.Print(content)
	case LevelInfo:
		l.newLogger.Print(content)
	case LevelWarn:
		l.newLogger.Print(content)
	case LevelError:
		l.newLogger.Print(content)
	case LevelFatal:
		l.newLogger.Fatal(content)
	case LevelPanic:
		l.newLogger.Panic(content)
	}
}

// 根据先前定义的日志分级，编写对应的日志输出的外部方法，
// 就是根据日志的六个等级，编写对应的方法
// 1.Info
func (l *Logger) Info(v ...interface{})  {
	l.Output(LevelInfo,fmt.Sprint(v...))
}
func (l *Logger) Infof(format string,v ...interface{})  {
	l.Output(LevelInfo,fmt.Sprintf(format,v...))
}

// 2.Fatal
func (l *Logger) Fatal(v ...interface{})  {
	l.Output(LevelFatal,fmt.Sprint(v...))
}

func (l *Logger) Fatalf(format string,v ...interface{})  {
	l.Output(LevelFatal,fmt.Sprintf(format,v...))
}

// 3.Debug
func (l *Logger) Debug(v ...interface{})  {
	l.Output(LevelDubug,fmt.Sprint(v...))
}
func (l *Logger) Debugf(format string,v ...interface{})  {
	l.Output(LevelDubug,fmt.Sprintf(format,v...))
}

// 4.Warn
func (l *Logger) Warn(v ...interface{})  {
	l.Output(LevelWarn,fmt.Sprint(v...))
}
func (l *Logger) Warnf(format string,v ...interface{})  {
	l.Output(LevelWarn,fmt.Sprintf(format,v...))
}

// 5.Error
func (l *Logger) Error(v ...interface{})  {
	l.Output(LevelError,fmt.Sprint(v...))
}

func (l *Logger) Errorf(format string,v ...interface{})  {
	l.Output(LevelError,fmt.Sprintf(format,v...))
}

// 6.Panic
func (l *Logger) Panic(v ...interface{})  {
	l.Output(LevelPanic,fmt.Sprint(v...))
}
func (l *Logger) Panicf(format string,v ...interface{})  {
	l.Output(LevelPanic,fmt.Sprintf(format,v...))
}