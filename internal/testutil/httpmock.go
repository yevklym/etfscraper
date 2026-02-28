package testutil

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

// MockHTTPClient is a simple programmable HTTP client for tests.
type MockHTTPClient struct {
	ResponseBody string
	StatusCode   int
	Error        error
	Delay        time.Duration
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	select {
	case <-req.Context().Done():
		return nil, req.Context().Err()
	default:
	}

	if m.Delay > 0 {
		select {
		case <-time.After(m.Delay):
		case <-req.Context().Done():
			return nil, req.Context().Err()
		}
	}

	if m.Error != nil {
		return nil, m.Error
	}

	response := &http.Response{
		StatusCode: m.StatusCode,
		Body:       io.NopCloser(bytes.NewReader([]byte(m.ResponseBody))),
		Header:     make(http.Header),
	}
	response.Header.Set("Content-Type", "application/json")

	return response, nil
}

// CountingMockHTTPClient wraps MockHTTPClient and counts calls.
type CountingMockHTTPClient struct {
	MockHTTPClient
	CallCount int
}

func (c *CountingMockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	c.CallCount++
	return c.MockHTTPClient.Do(req)
}
