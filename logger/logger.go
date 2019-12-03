package logger

import (
	"fmt"
	"log"
	"os"
)

type LoggerConf struct {
	LogFile string
}

// Logger is responsible to write logs in every component in the system.
type Logger struct{
	f *os.File 
	log *log.Logger
	LoggerConf LoggerConf
}

// Creates new Logger object.
func NewLogger(conf LoggerConf)(*Logger, error){
	f, err := os.OpenFile(conf.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err !=nil{
		return nil, err
	}
	l := log.New(f, "", log.Ldate|log.Ltime|log.Lmicroseconds) 
	return &Logger{
		f : f,
		log : l,
		LoggerConf : conf,
	}, nil
}

func (l *Logger) Info(flow string, format string, a ...interface{}){
	m := fmt.Sprintf(format, a...)
	l.log.Printf(fmt.Sprintf("level=Info flow=%s %s", flow, m))
	l.f.Sync()
}

func (l *Logger) Error(flow string, format string, a ...interface{}){
	m := fmt.Sprintf(format, a...)
	l.log.Printf(fmt.Sprintf("level=Error flow=%s %s", flow, m))
	l.f.Sync()
}

func (l *Logger) Debug(flow string, format string, a ...interface{}){
	m := fmt.Sprintf(format, a...)
	l.log.Printf(fmt.Sprintf("level=Debug flow=%s %s", flow, m))
	l.f.Sync()
}
