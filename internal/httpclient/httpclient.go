package httpclient

import (
	"bytes"
	"context"
	"io"
	"net"
	"net/http"
	"time"
)

var defaultClient *http.Client

func init() {
	transport := &http.Transport{
		// 连接池
		MaxIdleConns:        1000,
		MaxIdleConnsPerHost: 100,
		MaxConnsPerHost:     200,

		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 5 * time.Second,

		ExpectContinueTimeout: 1 * time.Second,

		// TCP 超时
		DialContext: (&net.Dialer{
			Timeout:   3 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,

		ForceAttemptHTTP2: true,
	}

	defaultClient = &http.Client{
		Transport: transport,
		Timeout:   5 * time.Second, // 整体超时
	}
}

func Do(req *http.Request) (*http.Response, error) {
	return defaultClient.Do(req)
}

func Get(ctx context.Context, url string) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, err
	}

	resp, err := defaultClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	return body, resp.StatusCode, nil
}

func PostJSON(ctx context.Context, url string, body []byte) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := defaultClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	return respBody, resp.StatusCode, nil
}
