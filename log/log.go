package log

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

)

const (
    INFO  = "INFO "
    TRACE = "TRACE"
    DEBUG = "DEBUG"
    ERROR = "ERROR"
    formatter = "date time \tlevel\tmethod(line number) - message"
)

func showSimpleLog(level, out string) { 
    t := time.Now().Format(time.DateTime) 
    fmt.Printf("%v %s \t- %s\n", t, level, out)
}
func showLogStatement(level string, out string, f string, l int) {
    wd, err := os.Getwd()
    if err != nil {
        showSimpleLog(level, out)
    }

    t := time.Now().Format(time.DateTime) 
    dirs := strings.Split(f, wd) 
    fmt.Printf("%v %s [ %s (line: %d) ] - %s\n", t, level, dirs[1], l, out)    
}

func Info(out string) {
    _, file, line, ok := runtime.Caller(1)
    if ok {
        showLogStatement(INFO, out, file, line)
    } else {
        showSimpleLog(INFO, out)
    } 
}

func Trace(out string) {
    _, file, line, ok := runtime.Caller(1)
    if ok {
        showLogStatement(TRACE, out, file, line)
    } else {
        showSimpleLog(TRACE, out)
    } 
}

func Debug(out string) {
    _, file, line, ok := runtime.Caller(1)
    if ok {
        showLogStatement(DEBUG, out, file, line)
    } else {
        showSimpleLog(DEBUG, out)
    } 

}

func Error(out string) {
    _, file, line, ok := runtime.Caller(1)
    if ok {
        showLogStatement(ERROR, out, file, line)
    } else {
        showSimpleLog(ERROR, out)
    } 
}
