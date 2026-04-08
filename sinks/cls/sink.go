package cls

import (
	"context"
	"fmt"
	"time"

	"github.com/cnmax/gologging-ext/core"
	cls "github.com/tencentcloud/tencentcloud-cls-sdk-go"
)

type sink struct {
	ctx       context.Context
	client    *cls.CLSClient
	topicID   string
	formatter core.MessageFormatter[*cls.Log]
}

type Option func(*sink)

func WithFormatter(formatter core.MessageFormatter[*cls.Log]) Option {
	return func(s *sink) {
		s.formatter = formatter
	}
}

func New(ctx context.Context, endpoint, secretID, secretKey, topicID string, opts ...Option) core.Writer {
	clsOpts := &cls.Options{
		Host: "https://" + endpoint,
		Credentials: cls.Credentials{
			SecretID:  secretID,
			SecretKEY: secretKey,
		},
	}

	client, err := cls.NewCLSClient(clsOpts)
	if err != nil {
		panic(APIError{CLSError: err})
	}

	s := &sink{
		ctx:     ctx,
		client:  client,
		topicID: topicID,
	}

	for _, opt := range opts {
		opt(s)
	}

	if s.formatter == nil {
		s.formatter = &DefaultCLSFormatter{}
	}

	return s
}

func (s *sink) Write(entry *core.Entry) error {
	log, err := s.formatter.Format(entry)
	if err != nil {
		return fmt.Errorf("%w: %v", core.ErrFormat, err)
	}

	ctx, cancel := context.WithTimeout(s.ctx, 3*time.Second)
	defer cancel()

	group := &cls.LogGroup{Logs: []*cls.Log{log}}

	if err := s.client.Send(ctx, s.topicID, group); err != nil {
		return fmt.Errorf("%w: %v", core.ErrSend, &APIError{CLSError: err})
	}

	return nil
}

func (s *sink) WriteBatch(entries []*core.Entry) error {
	if len(entries) == 0 {
		return nil
	}

	logs := make([]*cls.Log, 0, len(entries))
	var dropped int

	for _, entry := range entries {
		log, err := s.formatter.Format(entry)
		if err != nil {
			dropped++
			continue
		}
		logs = append(logs, log)
	}

	if len(logs) == 0 {
		return fmt.Errorf("%w: all logs dropped (%d)",
			core.ErrFormat, dropped)
	}

	ctx, cancel := context.WithTimeout(s.ctx, 3*time.Second)
	defer cancel()

	group := &cls.LogGroup{Logs: logs}

	if err := s.client.Send(ctx, s.topicID, group); err != nil {
		return fmt.Errorf("%w: %v", core.ErrSend, &APIError{CLSError: err})
	}

	if dropped > 0 {
		return fmt.Errorf("%w: partial dropped (%d)",
			core.ErrFormat, dropped)
	}

	return nil
}

type APIError struct {
	*cls.CLSError
}

func (e *APIError) Error() string {
	return fmt.Sprintf("cls error status=%d code=%s msg=%s", e.HTTPCode, e.Code, e.Message)
}

func (e *APIError) Retryable() bool {
	return e.HTTPCode >= 500 || e.HTTPCode == 429
}
