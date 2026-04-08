package internal

import (
	"fmt"
	"time"
)

const module = "gologging-ext"

func log(level string, component string, msg string, args ...interface{}) {
	ts := time.Now().Format("2006-01-02 15:04:05")

	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}

	fmt.Printf(
		"%s [%s] %-5s %s: %s\n",
		ts,
		module,
		level,
		component,
		msg,
	)
}

func Warn(component, msg string, args ...interface{}) {
	log("WARN", component, msg, args...)
}

func Error(component, msg string, args ...interface{}) {
	log("ERROR", component, msg, args...)
}
