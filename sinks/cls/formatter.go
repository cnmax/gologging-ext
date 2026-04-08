package cls

import (
	"fmt"
	"time"

	"github.com/cnmax/gologging-ext/core"
	cls "github.com/tencentcloud/tencentcloud-cls-sdk-go"
)

type DefaultCLSFormatter struct{}

func (f *DefaultCLSFormatter) Format(entry *core.Entry) (*cls.Log, error) {
	contents := make([]*cls.Log_Content, 0, 5+len(entry.Fields))

	contents = append(contents, &cls.Log_Content{
		Key:   new("level"),
		Value: new(entry.Level),
	})

	contents = append(contents, &cls.Log_Content{
		Key:   new("message"),
		Value: new(entry.Message),
	})

	contents = append(contents, &cls.Log_Content{
		Key:   new("time"),
		Value: new(entry.Time.Format(time.RFC3339)),
	})

	if entry.File != "" {
		contents = append(contents, &cls.Log_Content{
			Key:   new("location"),
			Value: new(entry.Location()),
		})
	}

	if entry.Stack != "" {
		contents = append(contents, &cls.Log_Content{
			Key:   new("stack"),
			Value: new(entry.Stack),
		})
	}

	for k, v := range entry.Fields {
		contents = append(contents, &cls.Log_Content{
			Key:   new(k),
			Value: new(fmt.Sprint(v)),
		})
	}

	return &cls.Log{
		Time:     new(entry.Time.UnixMilli()),
		Contents: contents,
	}, nil
}
