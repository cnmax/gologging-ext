package core

import (
	"errors"
	"fmt"
)

var (
	ErrFormat    = errors.New("log format error")
	ErrSerialize = errors.New("log serialize error")
	ErrSend      = errors.New("log send error")
)

type BackpressureError struct {
	Capacity   int // 队列容量
	CurrentLen int // 当前队列长度
}

func (e *BackpressureError) Error() string {
	return fmt.Sprintf("backpressure: %d/%d logs in queue", e.CurrentLen, e.Capacity)
}
