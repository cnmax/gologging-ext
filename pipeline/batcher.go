package pipeline

import (
	"time"

	"github.com/cnmax/gologging-ext/core"
	"github.com/cnmax/gologging-ext/internal"
)

type Batcher struct {
	writer        core.Writer
	ch            chan *core.Entry
	batchSize     int
	flushInterval time.Duration
}

func NewBatcher(
	writer core.Writer,
	batchSize int,
	flushInterval time.Duration,
) *Batcher {
	b := &Batcher{
		writer:        writer,
		ch:            make(chan *core.Entry, 1000),
		batchSize:     batchSize,
		flushInterval: flushInterval,
	}

	go b.run()
	return b
}

func (b *Batcher) Write(e *core.Entry) error {
	select {
	case b.ch <- e:
		return nil
	default:
		return &core.BackpressureError{
			Capacity:   cap(b.ch),
			CurrentLen: len(b.ch),
		}
	}
}

func (b *Batcher) run() {
	ticker := time.NewTicker(b.flushInterval)
	defer ticker.Stop()

	batch := make([]*core.Entry, 0, b.batchSize)

	for {
		select {
		case e := <-b.ch:
			batch = append(batch, e)

			if len(batch) >= b.batchSize {
				b.flush(batch)
				batch = batch[:0]
			}

		case <-ticker.C:
			if len(batch) > 0 {
				b.flush(batch)
				batch = batch[:0]
			}

		}
	}
}

func (b *Batcher) flush(batch []*core.Entry) {
	defer func() {
		if r := recover(); r != nil {
			internal.Error("batcher", "flush panic: %v", r)
		}
	}()

	var err error

	if bw, ok := b.writer.(BatchWriter); ok {
		err = bw.WriteBatch(batch)
	} else {
		for _, e := range batch {
			if e2 := b.writer.Write(e); e2 != nil {
				err = e2
				break
			}
		}
	}

	if err != nil {
		b.handleError(err, batch)
	}
}

func (b *Batcher) handleError(err error, batch []*core.Entry) {
	internal.Error("batcher", "flush failed: %v, batch=%d", err, len(batch))
}
