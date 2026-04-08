package zap

import (
	"encoding/json"

	"github.com/cnmax/gologging-ext/core"
	"go.uber.org/zap/zapcore"
)

type Core struct {
	writer  core.Writer
	level   zapcore.Level
	fields  []zapcore.Field
	encoder zapcore.Encoder
}

func NewCore(writer core.Writer, level zapcore.Level) *Core {
	return &Core{
		writer:  writer,
		level:   level,
		encoder: zapcore.NewJSONEncoder(zapcore.EncoderConfig{}),
	}
}

func (c *Core) Enabled(level zapcore.Level) bool {
	return level >= c.level
}

func (c *Core) With(fields []zapcore.Field) zapcore.Core {
	newFields := make([]zapcore.Field, 0, len(c.fields)+len(fields))
	newFields = append(newFields, c.fields...)
	newFields = append(newFields, fields...)

	return &Core{
		writer:  c.writer,
		level:   c.level,
		fields:  newFields,
		encoder: c.encoder,
	}
}

func (c *Core) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}

func (c *Core) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	entry := &core.Entry{
		Time:    ent.Time,
		Level:   ent.Level.CapitalString(),
		Message: ent.Message,
		Stack:   ent.Stack,
	}

	if ent.Caller.Defined {
		entry.Function = ent.Caller.Function
		entry.File = ent.Caller.File
		entry.Line = ent.Caller.Line
	}

	allFields := append(c.fields, fields...)

	if len(fields) > 0 {
		entry.Fields = fieldsToMap(c.encoder, ent, allFields)
	}

	return c.writer.Write(entry)
}

func (c *Core) Sync() error {
	return nil
}

func fieldsToMap(enc zapcore.Encoder, ent zapcore.Entry, fields []zapcore.Field) map[string]any {
	clone := enc.Clone()

	buf, err := clone.EncodeEntry(ent, fields)
	if err != nil {
		return nil
	}
	defer buf.Free()

	var m map[string]any
	_ = json.Unmarshal(buf.Bytes(), &m)

	return m
}
