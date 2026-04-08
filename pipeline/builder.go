package pipeline

import (
	"time"

	"github.com/cnmax/gologging-ext/core"
)

type BuildOptions struct {
	EnableAsync bool
	EnableBatch bool
	EnableRetry bool

	BatchSize     int
	FlushInterval time.Duration

	MaxRetries int
	BaseDelay  time.Duration
}

func DefaultBuildOptions() *BuildOptions {
	return &BuildOptions{
		EnableAsync: true,
		EnableBatch: true,
		EnableRetry: true,

		BatchSize:     100,
		FlushInterval: 2 * time.Second,

		MaxRetries: 3,
		BaseDelay:  100 * time.Millisecond,
	}
}

func Build(writer core.Writer, opts ...*BuildOptions) core.Writer {
	opt := DefaultBuildOptions()
	if len(opts) > 0 && opts[0] != nil {
		opt = opts[0]
	}

	w := writer

	if opt.EnableRetry {
		w = NewRetry(w, opt.MaxRetries, opt.BaseDelay)
	}

	if opt.EnableBatch {
		if _, ok := writer.(BatchWriter); ok {
			w = NewBatcher(w, opt.BatchSize, opt.FlushInterval)
		}
	}

	if opt.EnableAsync {
		if _, ok := w.(AsyncCapable); !ok {
			w = NewAsync(w)
		}
	}

	return w
}
