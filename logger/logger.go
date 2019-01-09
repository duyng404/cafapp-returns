package logger

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

// Info :
func Info(args ...interface{}) {
	fmt.Print(a(), "[INFO] ")
	fmt.Println(args...)
}

// Fatal :
func Fatal(args ...interface{}) {
	fmt.Print(a(), "[FATAL] ")
	fmt.Println(args...)
	os.Exit(1)
}

// Panic :
func Panic(args ...interface{}) {
	fmt.Print(a(), "[PANIC] ")
	fmt.Println(args...)
	panic("panicking!")
}

// Error :
func Error(args ...interface{}) {
	fmt.Print(a(), "[ERROR] ")
	fmt.Println(args...)
}

func a() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return fmt.Sprintf("[ %s:%d ] ", file, line)
}
