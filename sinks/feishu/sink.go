package feishu

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
	secret    string
	formatter core.JsonMessageFormatter
}

type Option func(*sink)

func WithFormatter(formatter core.JsonMessageFormatter) Option {
	return func(s *sink) {
		s.formatter = formatter
	}
}

func New(ctx context.Context, webhook, secret string, opts ...Option) core.Writer {
	s := &sink{
		ctx:     ctx,
		webhook: webhook,
		secret:  secret,
	}

	for _, opt := range opts {
		opt(s)
	}

	if s.formatter == nil {
		s.formatter = &DefaultFeiShuFormatter{}
	}

	return s
}

func (s *sink) Write(entry *core.Entry) error {
	msg, err := s.formatter.Format(entry)
	if err != nil {
		return fmt.Errorf("%w: %v", core.ErrFormat, err)
	}

	if s.secret != "" {
		ts := time.Now().Unix()
		sign, err := genSign(s.secret, ts)
		if err != nil {
			return fmt.Errorf("sign error: %w", err)
		}

		msg["timestamp"] = ts
		msg["sign"] = sign
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
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		return fmt.Errorf("%w: %v", core.ErrSerialize, err)
	}

	if result.Code != 0 {
		return fmt.Errorf("%w: %v", core.ErrSend, &APIError{
			Status:  status,
			Code:    result.Code,
			Message: result.Msg,
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
	return fmt.Sprintf("feishu error status=%d code=%d msg=%s", e.Status, e.Code, e.Message)
}

func (e *APIError) Retryable() bool {
	switch e.Code {
	case 1500, // 内部错误，请稍后重试
		1503, // 内部错误，更新token，但无任何鉴权结果返回，请检查后重试
		1642, // 内部错误，更新session失败，请稍后重试
		1663, // 内部错误，请稍后重试
		1665, // 内部错误
		1668, // 内部错误，请稍后重试
		2200, // 内部服务错误，通常会在频繁调用接口的情况下出现。请降低请求速率或增加重试机制。
		4006, // 服务异常，检查服务可用性
		5000: // 内部错误，减少调用频率，稍后再试
		return true
	}
	return e.Status > 500 || e.Code == 429
}
