package telegram

import (
	"fmt"
	"strings"

	"github.com/cnmax/gologging-ext/core"
)

type DefaultTelegramFormatter struct{}

func (f *DefaultTelegramFormatter) Format(e *core.Entry) (map[string]any, error) {
	var buf strings.Builder

	fmt.Fprintf(
		&buf,
		"<b>%s</b>\n%s",
		e.Level,
		e.Message,
	)

	if len(e.Fields) > 0 {
		fmt.Fprint(
			&buf,
			"<pre>",
		)

		for k, v := range e.Fields {
			fmt.Fprintf(
				&buf,
				"%s=%v\n",
				k,
				v,
			)
		}

		fmt.Fprint(
			&buf,
			"</pre>",
		)
	}

	if e.File != "" {
		fmt.Fprintf(
			&buf,
			"<i>Location: %s(%s:%d)</i>",
			e.Function,
			e.File,
			e.Line,
		)
	}

	return map[string]any{
		"text":                     buf.String(),
		"parse_mode":               "HTML",
		"disable_web_page_preview": true,
	}, nil
}
