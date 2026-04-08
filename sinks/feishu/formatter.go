package feishu

import (
	"fmt"

	"github.com/cnmax/gologging-ext/core"
)

type DefaultFeiShuFormatter struct {
	AtUsers []string
}

func (f *DefaultFeiShuFormatter) Format(e *core.Entry) (map[string]any, error) {
	var content [][]map[string]any

	content = append(content, []map[string]any{
		{
			"tag":  "text",
			"text": e.Message,
		},
	})

	if len(f.AtUsers) > 0 {
		var ats []map[string]any

		for _, at := range f.AtUsers {
			ats = append(ats, map[string]any{
				"tag":     "at",
				"user_id": at,
			})
		}

		content = append(content, ats)
	}

	if len(e.Fields) > 0 {
		for k, v := range e.Fields {
			content = append(content, []map[string]any{
				{
					"tag":  "text",
					"text": k + "=" + fmt.Sprintf("%v", v),
				},
			})
		}
	}

	if e.File != "" {
		content = append(content, []map[string]any{
			{
				"tag":  "text",
				"text": fmt.Sprintf("Location: %s(%s:%d)", e.Function, e.File, e.Line),
			},
		})
	}

	return map[string]any{
		"msg_type": "post",
		"content": map[string]any{
			"post": map[string]any{
				"zh_cn": map[string]any{
					"title":   e.Level,
					"content": content,
				},
			},
		},
	}, nil
}
