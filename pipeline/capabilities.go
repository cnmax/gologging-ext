package pipeline

import "github.com/cnmax/gologging-ext/core"

type BatchWriter interface {
	WriteBatch([]*core.Entry) error
}

type AsyncCapable interface {
	Async() bool
}

type Retryable interface {
	Retryable() bool
}
