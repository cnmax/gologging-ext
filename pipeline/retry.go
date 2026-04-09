package pipeline

import (
	"errors"
	"time"

	"github.com/cnmax/gologging-ext/core"
)

type Retry struct {
	writer     core.Writer
	maxRetries int
	baseDelay  time.Duration
	maxDelay   time.Duration
}

func NewRetry(
	writer core.Writer,
	maxRetries int,
	baseDelay time.Duration,
) *Retry {
	if r, ok := writer.(*Retry); ok {
		return r
	}

	if baseDelay <= 0 {
		baseDelay = 500 * time.Millisecond
	}

	return &Retry{
		writer:     writer,
		maxRetries: maxRetries,
		baseDelay:  baseDelay,
		maxDelay:   3 * time.Second,
	}
}

func (r *Retry) Write(e *core.Entry) error {
	var err error

	delay := r.baseDelay

	for i := 0; i <= r.maxRetries; i++ {
		err = r.writer.Write(e)
		if err == nil {
			return nil
		}

		var retryable Retryable
		if !errors.As(err, &retryable) {
			return err
		}

		if !retryable.Retryable() {
			return err
		}

		if i == r.maxRetries {
			break
		}

		time.Sleep(delay)

		// 指数退避
		delay *= 2
		if delay > r.maxDelay {
			delay = r.maxDelay
		}
	}

	return err
}
