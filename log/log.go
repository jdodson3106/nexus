package log

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

var (
	Black   = "\033[1;30m%s\033[0m"
	Red     = "\033[1;31m%s\033[0m"
	Green   = "\033[1;32m%s\033[0m"
	Yellow  = "\033[1;33m%s\033[0m"
	Purple  = "\033[1;34m%s\033[0m"
	Magenta = "\033[1;35m%s\033[0m"
	Teal    = "\033[1;36m%s\033[0m"
	White   = "\033[1;37m%s\033[0m"
)

var (
	INFO      = fmt.Sprintf(Teal, "INFO ")
	TRACE     = fmt.Sprintf(Magenta, "TRACE")
	DEBUG     = fmt.Sprintf(Green, "DEBUG")
	ERROR     = fmt.Sprintf(Red, "ERROR")
	WARN      = fmt.Sprintf(Yellow, "WARN ")
	formatter = "date time \tlevel\tmethod(line number) - message"
)

func showSimpleLog(level, out string) {
	t := time.Now().Format(time.DateTime)
	fmt.Printf("%v %s \t- %s\n", t, level, fmt.Sprintf(Purple, out))
}
func showLogStatement(level string, out string, f string, l int) {
	wd, err := os.Getwd()
	if err != nil {
		showSimpleLog(level, fmt.Sprintf(Purple, out))
	}

	t := time.Now().Format(time.DateTime)
	dirs := strings.Split(f, wd)
	loc := strings.Split(dirs[len(dirs)-1], "/")
	fmt.Printf("%v %s [ %s (line: %d) ] - %s\n", t, level, loc[len(loc)-1], l, out)
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
