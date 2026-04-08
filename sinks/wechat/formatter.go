package wechat

import (
	"fmt"
	"strings"

	"github.com/cnmax/gologging-ext/core"
)

type DefaultWechatFormatter struct{}

func (f *DefaultWechatFormatter) Format(e *core.Entry) (map[string]any, error) {
	var buf strings.Builder

	fmt.Fprintf(
		&buf,
		"# %s\n%s\n",
		e.Level,
		e.Message,
	)

	for k, v := range e.Fields {
		fmt.Fprintf(
			&buf,
			"\n> <font color=\"info\">%s=%v</font>",
			k,
			v,
		)
	}

	if e.File != "" {
		fmt.Fprintf(
			&buf,
			"\n<font color=\"warning\">Location: %s(%s:%d)</font>",
			e.Function,
			e.File,
			e.Line,
		)
	}

	return map[string]any{
		"msgtype": "markdown",
		"markdown": map[string]any{
			"content": buf.String(),
		},
	}, nil
}
