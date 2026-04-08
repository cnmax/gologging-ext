package zerolog

import (
	"encoding/json"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/cnmax/gologging-ext/core"
	"github.com/rs/zerolog"
	"github.com/tidwall/gjson"
)

type Writer struct {
	writer core.Writer
	level  zerolog.Level
}

func NewWriter(writer core.Writer, level zerolog.Level) *Writer {
	return &Writer{
		writer: writer,
		level:  level,
	}
}

func (w *Writer) Write(p []byte) (int, error) {
	levelStr := gjson.GetBytes(p, "level").String()

	if fastParseLevel(levelStr) < w.level {
		return 0, nil
	}

	var log map[string]any
	if err := json.Unmarshal(p, &log); err != nil {
		return 0, err
	}

	fields := make(map[string]any, len(log))
	for k, v := range log {
		if k == "level" {
			continue
		}
		fields[k] = v
	}

	entry := &core.Entry{
		Time:   time.Now(),
		Level:  strings.ToUpper(levelStr),
		Fields: fields,
	}

	if v, ok := fields["message"].(string); ok {
		entry.Message = v
		delete(fields, "message")
	}

	if v, ok := fields["time"].(string); ok {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			entry.Time = t
		}
		delete(fields, "time")
	}

	if v, ok := fields["caller"].(string); ok {
		entry.File, entry.Line, entry.Function = getFullCaller(v, 10)
		delete(fields, "caller")
	}

	return len(p), w.writer.Write(entry)
}

func fastParseLevel(lv string) zerolog.Level {
	switch lv {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		return zerolog.InfoLevel
	}
}

func getFullCaller(caller string, maxDepth int) (file string, line int, function string) {
	file0, line0 := splitCaller(caller)

	for i := 2; i < maxDepth; i++ {
		pc, f, l, ok := runtime.Caller(i)
		if !ok {
			break
		}

		if f != file0 {
			continue
		}

		if l != line0 {
			continue
		}

		if fn := runtime.FuncForPC(pc); fn != nil {
			function = fn.Name()
		}

		return f, l, function
	}

	return file0, line0, ""
}

func splitCaller(caller string) (file string, line int) {
	idx := strings.LastIndex(caller, ":")
	if idx == -1 {
		return caller, 0
	}

	file = caller[:idx]

	l, err := strconv.Atoi(caller[idx+1:])
	if err != nil {
		return file, 0
	}

	return file, l
}
