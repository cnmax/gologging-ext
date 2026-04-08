package pipeline

import (
	"fmt"

	"github.com/cnmax/gologging-ext/core"
)

type Dispatcher struct {
	sinks []core.Writer
}

func NewDispatcher(sinks ...core.Writer) *Dispatcher {
	return &Dispatcher{sinks: sinks}
}

func (d *Dispatcher) Write(entry *core.Entry) error {
	var errs []error

	for _, s := range d.sinks {
		if err := s.Write(entry); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("dispatcher: %v", errs)
	}

	return nil
}
