package telegram

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
	token     string
	chatID    int
	endpoint  string
	formatter core.JsonMessageFormatter
}

type Option func(*sink)

func WithEndpoint(endpoint string) Option {
	return func(s *sink) {
		s.endpoint = endpoint
	}
}

func WithFormatter(formatter core.JsonMessageFormatter) Option {
	return func(s *sink) {
		s.formatter = formatter
	}
}

func New(ctx context.Context, token string, chatID int, opts ...Option) core.Writer {
	s := &sink{
		ctx:    ctx,
		token:  token,
		chatID: chatID,
	}

	for _, opt := range opts {
		opt(s)
	}

	if s.formatter == nil {
		s.formatter = &DefaultTelegramFormatter{}
	}

	if s.endpoint == "" {
		s.endpoint = "https://api.telegram.org"
	}

	return s
}

func (s *sink) Write(entry *core.Entry) error {
	msg, err := s.formatter.Format(entry)
	if err != nil {
		return fmt.Errorf("%w: %v", core.ErrFormat, err)
	}

	// ChatID
	msg["chat_id"] = s.chatID

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
			Message: string(body),
		})
	}

	return s.parseResponse(resp, status)
}

func (s *sink) doRequest(body []byte) ([]byte, int, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 3*time.Second)
	defer cancel()

	return httpclient.PostJSON(
		ctx,
		fmt.Sprintf("%s/bot%s/sendMessage", s.endpoint, s.token),
		body,
	)
}

func (s *sink) parseResponse(resp []byte, status int) error {
	var result struct {
		Ok          bool   `json:"ok"`
		ErrorCode   int    `json:"error_code"`
		Description string `json:"description"`
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		return fmt.Errorf("%w: %v", core.ErrSerialize, err)
	}

	if !result.Ok {
		return fmt.Errorf("%w: %v", core.ErrSend, &APIError{
			Status:  status,
			Code:    result.ErrorCode,
			Message: result.Description,
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
	return fmt.Sprintf("telegram error status=%d code=%d msg=%s", e.Status, e.Code, e.Message)
}

func (e *APIError) Retryable() bool {
	return e.Status > 500 || e.Status == 429 || e.Code == 429
}
