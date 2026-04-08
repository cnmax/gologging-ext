package dingtalk

import (
	"fmt"
	"strings"

	"github.com/cnmax/gologging-ext/core"
)

type DefaultDingTalkFormatter struct {
	AtMobiles []string
	AtUserIds []string
	IsAtAll   bool
}

func (f *DefaultDingTalkFormatter) Format(e *core.Entry) (map[string]any, error) {
	var buf strings.Builder

	fmt.Fprintf(
		&buf,
		"# %s\n%s\n",
		e.Level,
		e.Message,
	)

	if !f.IsAtAll && (len(f.AtMobiles) > 0 || len(f.AtUserIds) > 0) {
		fmt.Fprint(
			&buf,
			"\n",
		)

		for _, at := range f.AtMobiles {
			fmt.Fprintf(
				&buf,
				"@%s",
				at,
			)
		}

		for _, v := range f.AtUserIds {
			fmt.Fprintf(
				&buf,
				"@%s",
				v,
			)
		}
	}

	if len(e.Fields) > 0 {
		fmt.Fprint(
			&buf,
			"\n> ",
		)

		var items []string
		for k, v := range e.Fields {
			items = append(items, k+"="+fmt.Sprintf("%v", v))
		}

		fmt.Fprint(
			&buf,
			strings.Join(items, ", "),
		)
	}

	if e.File != "" {
		fmt.Fprintf(
			&buf,
			"\n###### Location: %s(%s:%d)",
			e.Function,
			e.File,
			e.Line,
		)
	}

	title := fmt.Sprintf("[%s] %s", e.Level, e.Message)

	return map[string]any{
		"msgtype": "markdown",
		"markdown": map[string]any{
			"title": title,
			"text":  buf.String(),
		},
		"at": map[string]any{
			"atMobiles": f.AtMobiles,
			"atUserIds": f.AtUserIds,
			"isAtAll":   f.IsAtAll,
		},
	}, nil
}
