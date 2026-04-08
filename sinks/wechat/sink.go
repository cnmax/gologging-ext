package wechat

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cnmax/gologging-ext/core"
	"github.com/cnmax/gologging-ext/internal/httpclient"
)

type sink struct {
	ctx       context.Context
	webhook   string
	formatter core.JsonMessageFormatter
}

type Option func(*sink)

func WithFormatter(formatter core.JsonMessageFormatter) Option {
	return func(s *sink) {
		s.formatter = formatter
	}
}

func New(ctx context.Context, webhook string, opts ...Option) core.Writer {
	s := &sink{
		ctx:     ctx,
		webhook: webhook,
	}

	for _, opt := range opts {
		opt(s)
	}

	if s.formatter == nil {
		s.formatter = &DefaultWechatFormatter{}
	}

	return s
}

func (s *sink) Write(entry *core.Entry) error {
	msg, err := s.formatter.Format(entry)
	if err != nil {
		return fmt.Errorf("%w: %v", core.ErrFormat, err)
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("%w: %v", core.ErrSerialize, err)
	}

	resp, status, err := s.doRequest(body)

	if err != nil {
		return fmt.Errorf("%w: %v", core.ErrSend, err)
	}

	if status != 200 {
		return fmt.Errorf("%w: %v", core.ErrSend, &APIError{
			Status:  status,
			Code:    status,
			Message: string(resp),
		})
	}

	return s.parseResponse(resp, status)
}

func (s *sink) doRequest(body []byte) ([]byte, int, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 3*time.Second)
	defer cancel()

	return httpclient.PostJSON(ctx, s.webhook, body)
}

func (s *sink) parseResponse(resp []byte, status int) error {
	var result struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		return fmt.Errorf("%w: %v", core.ErrSerialize, err)
	}

	if result.ErrCode != 0 {
		return fmt.Errorf("%w: %v", core.ErrSend, &APIError{
			Status:  status,
			Code:    result.ErrCode,
			Message: result.ErrMsg,
		})
	}

	return nil
}

type APIError struct {
	Status  int
	Code    int
	Message string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("wechat error status=%d code=%d msg=%s", e.Status, e.Code, e.Message)
}

func (e *APIError) Retryable() bool {
	return e.Status > 500 || e.Status == 429 || e.Code == -1
}
