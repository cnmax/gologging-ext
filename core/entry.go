package core

import (
	"strconv"
	"time"
)

type Entry struct {
	Time     time.Time
	Level    string
	Message  string
	Function string
	File     string
	Line     int
	Stack    string
	Fields   map[string]interface{}
}

func (e *Entry) Location() string {
	buf := make([]byte, 0, len(e.Function)+len(e.File)+16)
	buf = append(buf, e.Function...)
	buf = append(buf, '(')
	buf = append(buf, e.File...)
	buf = append(buf, ':')
	buf = strconv.AppendInt(buf, int64(e.Line), 10)
	buf = append(buf, ')')
	return string(buf)
}
