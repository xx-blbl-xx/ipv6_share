package xlog

import (
	"fmt"
	"log"
	"os"
)

const (
	prefixErr  = "[ERROR] "
	prefixWarn = "[WARN] "
	prefixInfo = "[INFO] "
)

type myLogger struct {
	*log.Logger
	objFormat string
}

var (
	errLogger = &myLogger{
		Logger:    log.New(os.Stderr, prefixErr, log.LstdFlags),
		objFormat: " %d:OBJ:[%#v]",
	}
	warnLogger = &myLogger{
		Logger:    log.New(os.Stderr, prefixWarn, log.LstdFlags),
		objFormat: " %d:OBJ:[%#v]",
	}
	infoLogger = &myLogger{
		Logger:    log.New(os.Stderr, prefixInfo, log.LstdFlags),
		objFormat: " %d:OBJ:[%+v]",
	}
)

func Error(info string, v ...interface{}) {
	errLogger.out(info, v)
}

func Warn(info string, v ...interface{}) {
	warnLogger.out(info, v)
}

func Info(info string, v ...interface{}) {
	infoLogger.out(info, v)
}

func (l *myLogger) out(info string, objs []interface{}) {
	str := info
	if len(objs) > 0 {
		for k, v := range objs {
			str += fmt.Sprintf(l.objFormat, k, v)
		}
	}
	_ = l.Output(2, str)
}
