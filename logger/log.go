package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/mysll/toolkit"
)

const (
	LOG_DEBUG = iota
	LOG_INFO
	LOG_WARN
	LOG_ERR
	LOG_FATAL
)

type Logger interface {
	// 日志函数
	LogDebug(v ...interface{})
	// 日志函数
	LogInfo(v ...interface{})
	// 日志函数
	LogWarn(v ...interface{})
	// 日志函数
	LogErr(v ...interface{})
	// 日志函数
	LogFatal(v ...interface{})
	// 日志函数
	LogDebugf(format string, v ...interface{})
	// 日志函数
	LogInfof(format string, v ...interface{})
	// 日志函数
	LogWarnf(format string, v ...interface{})
	// 日志函数
	LogErrf(format string, v ...interface{})
	// 日志函数
	LogFatalf(format string, v ...interface{})
}

type Log struct {
	logger   *log.Logger
	logfile  *os.File
	logLevel int
}

func New(file string, log_level int) *Log {
	l := &Log{}
	p, _ := path.Split(file)
	if ok, err := toolkit.PathExists(p); err != nil {
		panic(err)
	} else {
		if !ok {
			os.Mkdir(p, os.ModePerm)
		}
	}

	l.logLevel = log_level
	var err error
	l.logfile, err = os.Create(file)
	if err != nil {
		panic("create log file failed " + file)
	}
	w := io.MultiWriter(l.logfile, os.Stdout)
	l.logger = log.New(w, "", log.Ldate|log.Ltime|log.Lshortfile)
	return l
}

func (l *Log) CloseLog() {
	if l.logfile != nil {
		l.logfile.Close()
	}
}

func (l *Log) SetPrefix(p string) {
	if l.logger == nil {
		log.SetPrefix(p)
		return
	}
	l.logger.SetPrefix(p)
}

func (l *Log) Outputf(depth int, prefix string, format string, v ...interface{}) {
	if l.logger == nil {
		log.SetPrefix(prefix)
		log.Output(depth, fmt.Sprintf(format, v...))
		return
	}
	l.logger.SetPrefix(prefix)
	l.logger.Output(depth, fmt.Sprintf(format, v...))
}

func (l *Log) Output(depth int, prefix string, v ...interface{}) {
	if l.logger == nil {
		log.SetPrefix(prefix)
		log.Output(depth, fmt.Sprint(v...))
		return
	}
	l.logger.SetPrefix(prefix)
	l.logger.Output(depth, fmt.Sprint(v...))
}

func (l *Log) LogDebugf(format string, v ...interface{}) {
	if l.logLevel > LOG_DEBUG {
		return
	}
	if l.logger == nil {
		log.SetPrefix("[D] ")
		log.Output(2, fmt.Sprintf(format, v...))
		return
	}

	l.logger.SetPrefix("[D] ")
	l.logger.Output(2, fmt.Sprintf(format, v...))
}

func (l *Log) LogInfof(format string, v ...interface{}) {
	if l.logLevel > LOG_INFO {
		return
	}
	if l.logger == nil {
		log.SetPrefix("[I] ")
		log.Output(2, fmt.Sprintf(format, v...))
		return
	}

	l.logger.SetPrefix("[I] ")
	l.logger.Output(2, fmt.Sprintf(format, v...))
}

func (l *Log) LogWarnf(format string, v ...interface{}) {
	if l.logLevel > LOG_WARN {
		return
	}
	if l.logger == nil {
		log.SetPrefix("[W] ")
		log.Output(2, fmt.Sprintf(format, v...))
		return
	}

	l.logger.SetPrefix("[W] ")
	l.logger.Output(2, fmt.Sprintf(format, v...))
}

func (l *Log) LogErrf(format string, v ...interface{}) {
	if l.logLevel > LOG_ERR {
		return
	}
	if l.logger == nil {
		log.SetPrefix("[E] ")
		log.Output(2, fmt.Sprintf(format, v...))
		return
	}

	l.logger.SetPrefix("[E] ")
	l.logger.Output(2, fmt.Sprintf(format, v...))
}

func (l *Log) LogFatalf(format string, v ...interface{}) {
	if l.logLevel > LOG_FATAL {
		return
	}
	if l.logger == nil {
		log.SetPrefix("[F] ")
		log.Output(2, fmt.Sprintf(format, v...))
		os.Exit(1)
		return
	}

	l.logger.SetPrefix("[F] ")
	l.logger.Output(2, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func (l *Log) LogDebug(v ...interface{}) {
	if l.logLevel > LOG_DEBUG {
		return
	}
	if l.logger == nil {
		log.SetPrefix("[D] ")
		log.Output(2, fmt.Sprint(v...))
		return
	}

	l.logger.SetPrefix("[D] ")
	l.logger.Output(2, fmt.Sprint(v...))
}

func (l *Log) LogInfo(v ...interface{}) {
	if l.logLevel > LOG_INFO {
		return
	}
	if l.logger == nil {
		log.SetPrefix("[I] ")
		log.Output(2, fmt.Sprint(v...))
		return
	}

	l.logger.SetPrefix("[I] ")
	l.logger.Output(2, fmt.Sprint(v...))
}

func (l *Log) LogWarn(v ...interface{}) {
	if l.logLevel > LOG_WARN {
		return
	}
	if l.logger == nil {
		log.SetPrefix("[W] ")
		log.Output(2, fmt.Sprint(v...))
		return
	}

	l.logger.SetPrefix("[W] ")
	l.logger.Output(2, fmt.Sprint(v...))
}

func (l *Log) LogErr(v ...interface{}) {
	if l.logLevel > LOG_ERR {
		return
	}
	if l.logger == nil {
		log.SetPrefix("[E] ")
		log.Output(2, fmt.Sprint(v...))
		return
	}

	l.logger.SetPrefix("[E] ")
	l.logger.Output(2, fmt.Sprint(v...))
}

func (l *Log) LogFatal(v ...interface{}) {
	if l.logLevel > LOG_FATAL {
		return
	}
	if l.logger == nil {
		log.SetPrefix("[F] ")
		log.Output(2, fmt.Sprint(v...))
		os.Exit(1)
		return
	}

	l.logger.SetPrefix("[F] ")
	l.logger.Output(2, fmt.Sprint(v...))
	os.Exit(1)
}
